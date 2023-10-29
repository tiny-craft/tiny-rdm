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
import useConnectionStore from 'stores/connections.js'
import Copy from '@/components/icons/Copy.vue'
import { ClipboardSetText } from 'wailsjs/runtime/runtime.js'
import { computed } from 'vue'
import { isEmpty, padStart } from 'lodash'

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
})

const dialogStore = useDialog()
const connectionStore = useConnectionStore()
const i18n = useI18n()

const binaryKey = computed(() => {
    return !!props.keyCode
})

/**
 *
 * @type {ComputedRef<string|number[]>}
 */
const keyName = computed(() => {
    return !isEmpty(props.keyCode) ? props.keyCode : props.keyPath
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

const onReloadKey = () => {
    connectionStore.loadKeyValue(props.server, props.db, keyName.value)
}

const onCopyKey = () => {
    ClipboardSetText(props.keyPath)
        .then((succ) => {
            if (succ) {
                $message.success(i18n.t('dialogue.copy_succ'))
            }
        })
        .catch((e) => {
            $message.error(e.message)
        })
}

const onRenameKey = () => {
    if (binaryKey.value) {
        $message.error(i18n.t('dialogue.rename_binary_key_fail'))
    } else {
        dialogStore.openRenameKeyDialog(props.server, props.db, props.keyPath)
    }
}

const onDeleteKey = () => {
    $dialog.warning(i18n.t('dialogue.remove_tip', { name: props.keyPath }), () => {
        connectionStore.deleteKey(props.server, props.db, keyName.value).then((success) => {
            if (success) {
                $message.success(i18n.t('dialogue.delete_key_succ', { key: props.keyPath }))
            }
        })
    })
}
</script>

<template>
    <div class="content-toolbar flex-box-h">
        <n-input-group>
            <redis-type-tag :binary-key="binaryKey" :type="props.keyType" size="large" />
            <n-input v-model:value="props.keyPath">
                <template #suffix>
                    <icon-button :icon="Refresh" size="18" t-tooltip="interface.reload" @click="onReloadKey" />
                </template>
            </n-input>
            <icon-button :icon="Copy" border size="18" t-tooltip="interface.copy_key" @click="onCopyKey" />
        </n-input-group>
        <n-button-group>
            <n-tooltip>
                <template #trigger>
                    <n-button :focusable="false" @click="dialogStore.openTTLDialog(props.ttl)">
                        <template #icon>
                            <n-icon :component="Timer" size="18" />
                        </template>
                        {{ ttlString }}
                    </n-button>
                </template>
                TTL{{ `${ttl > 0 ? ': ' + ttl + $t('common.second') : ''}` }}
            </n-tooltip>
            <icon-button :icon="Edit" border size="18" t-tooltip="interface.rename_key" @click="onRenameKey" />
        </n-button-group>
        <n-tooltip :show-arrow="false">
            <template #trigger>
                <n-button :focusable="false" @click="onDeleteKey">
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
