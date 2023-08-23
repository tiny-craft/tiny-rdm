<script setup>
import { reactive, watch } from 'vue'
import useDialog from 'stores/dialog'
import { useMessage } from '@/utils/message.js'
import { useI18n } from 'vue-i18n'
import useConnectionStore from 'stores/connections.js'

const renameForm = reactive({
    server: '',
    db: 0,
    key: '',
    newKey: '',
})

const dialogStore = useDialog()
const connectionStore = useConnectionStore()
watch(
    () => dialogStore.renameDialogVisible,
    (visible) => {
        if (visible) {
            const { server, db, key } = dialogStore.renameKeyParam
            renameForm.server = server
            renameForm.db = db
            renameForm.key = key
            renameForm.newKey = key
        }
    },
)

const i18n = useI18n()
const message = useMessage()
const onRename = async () => {
    try {
        const { server, db, key, newKey } = renameForm
        const { success, msg } = await connectionStore.renameKey(server, db, key, newKey)
        if (success) {
            await connectionStore.loadKeyValue(server, db, newKey)
            message.success(i18n.t('handle_succ'))
        } else {
            message.error(msg)
        }
    } catch (e) {
        message.error(e.message)
    }
    dialogStore.closeRenameKeyDialog()
}

const onClose = () => {
    dialogStore.closeRenameKeyDialog()
}
</script>

<template>
    <n-modal
        v-model:show="dialogStore.renameDialogVisible"
        :closable="false"
        :close-on-esc="false"
        :mask-closable="false"
        :negative-button-props="{ size: 'medium' }"
        :negative-text="$t('cancel')"
        :positive-button-props="{ size: 'medium' }"
        :positive-text="$t('confirm')"
        :show-icon="false"
        :title="$t('rename_key')"
        preset="dialog"
        transform-origin="center"
        @positive-click="onRename"
        @negative-click="onClose"
    >
        <n-form
            :model="renameForm"
            :show-require-mark="false"
            :show-label="false"
            label-align="left"
            label-placement="top"
        >
            <n-form-item :label="$t('new_key_name')" required>
                <n-input v-model:value="renameForm.newKey" />
            </n-form-item>
        </n-form>
    </n-modal>
</template>

<style lang="scss" scoped></style>
