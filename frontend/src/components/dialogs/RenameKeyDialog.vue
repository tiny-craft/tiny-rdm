<script setup>
import { reactive, watch } from 'vue'
import useDialog from 'stores/dialog'
import { useI18n } from 'vue-i18n'
import useBrowserStore from 'stores/browser.js'
import useTabStore from 'stores/tab.js'

const renameForm = reactive({
    server: '',
    db: 0,
    key: '',
    newKey: '',
})

const dialogStore = useDialog()
const browserStore = useBrowserStore()
const tab = useTabStore()
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
const onRename = async () => {
    try {
        const { server, db, key, newKey } = renameForm
        const { success, msg, nodeKey } = await browserStore.renameKey(server, db, key, newKey)
        if (success) {
            tab.setSelectedKeys(server, nodeKey)
            browserStore.loadKeySummary({ server, db, key: newKey })
            $message.success(i18n.t('dialogue.handle_succ'))
        } else {
            $message.error(msg)
        }
    } catch (e) {
        $message.error(e.message)
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
        :negative-text="$t('common.cancel')"
        :positive-button-props="{ size: 'medium' }"
        :positive-text="$t('common.confirm')"
        :show-icon="false"
        :title="$t('interface.rename_key')"
        preset="dialog"
        transform-origin="center"
        @positive-click="onRename"
        @negative-click="onClose">
        <n-form
            :model="renameForm"
            :show-label="false"
            :show-require-mark="false"
            label-align="left"
            label-placement="top">
            <n-form-item :label="$t('dialogue.key.new_name')" required>
                <n-input v-model:value="renameForm.newKey" :input-props="{ spellcheck: 'false' }" />
            </n-form-item>
        </n-form>
    </n-modal>
</template>

<style lang="scss" scoped></style>
