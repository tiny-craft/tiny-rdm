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
import { computed } from 'vue'
import { padStart } from 'lodash'

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

const dialogStore = useDialog()
const i18n = useI18n()

const binaryKey = computed(() => {
    return !!props.keyCode
})

const ttlString = computed(() => {
    let s = ''
    if (props.ttl > 0) {
        const hours = Math.floor(props.ttl / 3600)
        s += padStart(hours + ':', 3, '0')
        const minutes = Math.floor((props.ttl % 3600) / 60)
        s += padStart(minutes + ':', 3, '0')
        const seconds = Math.floor(props.ttl % 60)
        s += padStart(seconds + '', 2, '0')
    } else if (props.ttl < 0) {
        s = i18n.t('interface.forever')
    } else {
        s = '00:00:00'
    }
    return s
})

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
            <n-input v-model:value="props.keyPath" :title="props.keyPath" readonly>
                <template #suffix>
                    <icon-button
                        :icon="Refresh"
                        :loading="props.loading"
                        size="18"
                        t-tooltip="interface.reload"
                        @click="emit('reload')" />
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
