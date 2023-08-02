<script setup>
import { validType } from '@/consts/support_redis_type.js'
import useDialog from 'stores/dialog.js'
import Delete from '@/components/icons/Delete.vue'
import Edit from '@/components/icons/Edit.vue'
import Refresh from '@/components/icons/Refresh.vue'
import Timer from '@/components/icons/Timer.vue'
import RedisTypeTag from '@/components/common/RedisTypeTag.vue'
import { useI18n } from 'vue-i18n'
import { useMessage } from 'naive-ui'
import IconButton from '@/components/common/IconButton.vue'
import useConnectionStore from 'stores/connections.js'
import { useConfirmDialog } from '@/utils/confirm_dialog.js'
import Copy from '@/components/icons/Copy.vue'
import { ClipboardSetText } from 'wailsjs/runtime/runtime.js'

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
    ttl: {
        type: Number,
        default: -1,
    },
})

const dialogStore = useDialog()
const connectionStore = useConnectionStore()
const message = useMessage()
const i18n = useI18n()

const onReloadKey = () => {
    connectionStore.loadKeyValue(props.server, props.db, props.keyPath)
}

const onCopyKey = () => {
    ClipboardSetText(props.keyPath)
        .then((succ) => {
            if (succ) {
                message.success(i18n.t('copy_succ'))
            }
        })
        .catch((e) => {
            message.error(e.message)
        })
}

const confirmDialog = useConfirmDialog()
const onDeleteKey = () => {
    confirmDialog.warning(i18n.t('remove_tip', { name: props.keyPath }), () => {
        connectionStore.deleteKey(props.server, props.db, props.keyPath).then((success) => {
            if (success) {
                message.success(i18n.t('delete_key_succ', { key: props.keyPath }))
            }
        })
    })
}
</script>

<template>
    <div class="content-toolbar flex-box-h">
        <n-input-group>
            <redis-type-tag :type="props.keyType" size="large" />
            <n-input v-model:value="props.keyPath">
                <template #suffix>
                    <icon-button :icon="Refresh" t-tooltip="reload" size="18" @click="onReloadKey" />
                </template>
            </n-input>
            <icon-button :icon="Copy" t-tooltip="copy_key" size="18" border @click="onCopyKey" />
        </n-input-group>
        <n-button-group>
            <n-tooltip>
                <template #trigger>
                    <n-button @click="dialogStore.openTTLDialog(props.ttl)">
                        <template #icon>
                            <n-icon :component="Timer" size="18" />
                        </template>
                        <template v-if="ttl < 0">
                            {{ $t('forever') }}
                        </template>
                        <template v-else> {{ ttl }} {{ $t('second') }}</template>
                    </n-button>
                </template>
                TTL
            </n-tooltip>
            <icon-button
                border
                :icon="Edit"
                t-tooltip="rename_key"
                size="18"
                @click="dialogStore.openRenameKeyDialog(props.server, props.db, props.keyPath)"
            />
            <!--            <n-button @click="dialogStore.openRenameKeyDialog(props.server, props.db, props.keyPath)">-->
            <!--                <template #icon>-->
            <!--                    <n-icon :component="Edit" size="18" />-->
            <!--                </template>-->
            <!--                {{ $t('rename_key') }}-->
            <!--            </n-button>-->
        </n-button-group>
        <n-tooltip>
            <template #trigger>
                <n-button>
                    <template #icon>
                        <n-icon :component="Delete" size="18" @click="onDeleteKey" />
                    </template>
                </n-button>
            </template>
            {{ $t('delete_key') }}
        </n-tooltip>
    </div>
</template>

<style lang="scss" scoped>
.content-toolbar {
    align-items: center;
    gap: 5px;
}
</style>
