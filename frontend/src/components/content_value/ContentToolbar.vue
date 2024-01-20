<script setup>
import { validType } from '@/consts/support_redis_type.js'
import useDialog from 'stores/dialog.js'
import Delete from '@/components/icons/Delete.vue'
import Edit from '@/components/icons/Edit.vue'
import Refresh from '@/components/icons/Refresh.vue'
import Timer from '@/components/icons/Timer.vue'
import RedisTypeTag from '@/components/common/RedisTypeTag.vue'
import { useI18n } from 'vue-i18n'
import IconButton from '@/components/common/IconButton.vue'
import Copy from '@/components/icons/Copy.vue'
import { ClipboardSetText } from 'wailsjs/runtime/runtime.js'
import { computed, onUnmounted, reactive, watch } from 'vue'
import { NIcon, useThemeVars } from 'naive-ui'
import { timeout } from '@/utils/promise.js'
import AutoRefreshForm from '@/components/common/AutoRefreshForm.vue'
import dayjs from 'dayjs'

const props = defineProps({
    server: String,
    db: Number,
    keyType: {
        type: String,
        validator(value) {
            return validType(value)
        },
        default: 'STRING',
    },
    keyPath: String,
    keyCode: {
        type: Array,
        default: null,
    },
    ttl: {
        type: Number,
        default: -1,
    },
    loading: Boolean,
})

const emit = defineEmits(['reload', 'rename', 'delete'])

const autoRefresh = reactive({
    on: false,
    interval: 2,
})

const themeVars = useThemeVars()
const dialogStore = useDialog()
const i18n = useI18n()

const binaryKey = computed(() => {
    return !!props.keyCode
})

const ttlString = computed(() => {
    if (props.ttl > 0) {
        const dur = dayjs.duration(props.ttl, 'seconds')
        const days = dur.days()
        if (days > 0) {
            return days + i18n.t('common.unit_day') + ' ' + dur.format('HH:mm:ss')
        } else {
            return dur.format('HH:mm:ss')
        }
    } else if (props.ttl < 0) {
        return i18n.t('interface.forever')
    } else {
        return '00:00:00'
    }
})

const startAutoRefresh = async () => {
    let lastExec = Date.now()
    do {
        if (!autoRefresh.on) {
            break
        }
        await timeout(100)
        if (props.loading || Date.now() - lastExec < autoRefresh.interval * 1000) {
            continue
        }
        lastExec = Date.now()
        emit('reload')
    } while (true)
    stopAutoRefresh()
}

const stopAutoRefresh = () => {
    autoRefresh.on = false
}

watch(
    () => props.keyPath,
    () => {
        stopAutoRefresh()
        autoRefresh.interval = props.interval
    },
)

onUnmounted(() => stopAutoRefresh())

const onToggleRefresh = (on) => {
    if (on) {
        startAutoRefresh()
    } else {
        stopAutoRefresh()
    }
}

const onCopyKey = () => {
    ClipboardSetText(props.keyPath)
        .then((succ) => {
            if (succ) {
                $message.success(i18n.t('interface.copy_succ'))
            }
        })
        .catch((e) => {
            $message.error(e.message)
        })
}

const onTTL = () => {
    dialogStore.openTTLDialog({
        server: props.server,
        db: props.db,
        key: binaryKey.value ? props.keyCode : props.keyPath,
        ttl: props.ttl,
    })
}
</script>

<template>
    <div class="content-toolbar flex-box-h">
        <n-input-group>
            <redis-type-tag :binary-key="binaryKey" :type="props.keyType" size="large" />
            <n-input v-model:value="props.keyPath" :title="props.keyPath" readonly @dblclick="onCopyKey">
                <template #suffix>
                    <n-popover :delay="500" keep-alive-on-hover placement="bottom" trigger="hover">
                        <template #trigger>
                            <icon-button
                                :loading="props.loading"
                                size="18"
                                @click="emit('reload')"
                                @dblclick.stop="() => {}">
                                <n-icon :size="props.size">
                                    <refresh
                                        :class="{ 'auto-rotate': autoRefresh.on }"
                                        :color="autoRefresh.on ? themeVars.primaryColor : undefined"
                                        :stroke-width="autoRefresh.on ? 6 : 3" />
                                </n-icon>
                            </icon-button>
                        </template>
                        <auto-refresh-form
                            v-model:interval="autoRefresh.interval"
                            v-model:on="autoRefresh.on"
                            :default-value="2"
                            :loading="props.loading"
                            @toggle="onToggleRefresh" />
                    </n-popover>
                </template>
            </n-input>
            <icon-button :icon="Copy" border size="18" t-tooltip="interface.copy_key" @click="onCopyKey" />
        </n-input-group>
        <n-button-group>
            <n-tooltip>
                <template #trigger>
                    <n-button :focusable="false" @click="onTTL">
                        <template #icon>
                            <n-icon :component="Timer" size="18" />
                        </template>
                        {{ ttlString }}
                    </n-button>
                </template>
                TTL{{ `${ttl > 0 ? ': ' + ttl + $t('common.second') : ''}` }}
            </n-tooltip>
            <icon-button
                :disabled="binaryKey"
                :icon="Edit"
                :t-tooltip="binaryKey ? 'dialogue.rename_binary_key_fail' : 'interface.rename_key'"
                border
                size="18"
                @click="emit('rename')" />
        </n-button-group>
        <n-tooltip :show-arrow="false">
            <template #trigger>
                <n-button :focusable="false" @click="emit('delete')">
                    <template #icon>
                        <n-icon :component="Delete" size="18" />
                    </template>
                </n-button>
            </template>
            {{ $t('interface.delete_key') }}
        </n-tooltip>
    </div>
</template>

<style lang="scss" scoped>
.content-toolbar {
    align-items: center;
    gap: 5px;
}
</style>
