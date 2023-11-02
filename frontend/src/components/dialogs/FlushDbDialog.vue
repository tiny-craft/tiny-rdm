<script setup>
import { reactive, watch } from 'vue'
import useDialog from 'stores/dialog'
import { useI18n } from 'vue-i18n'
import useConnectionStore from 'stores/connections.js'

const flushForm = reactive({
    server: '',
    db: 0,
    key: '',
    async: false,
    confirm: false,
})

const dialogStore = useDialog()
const connectionStore = useConnectionStore()
watch(
    () => dialogStore.flushDBDialogVisible,
    (visible) => {
        if (visible) {
            const { server, db } = dialogStore.flushDBParam
            flushForm.server = server
            flushForm.db = db
            flushForm.async = true
            flushForm.confirm = false
        }
    },
)

const i18n = useI18n()
const onConfirmFlush = async () => {
    try {
        const { server, db, async } = flushForm
        const success = await connectionStore.flushDatabase(server, db, async)
        if (success) {
            $message.success(i18n.t('dialogue.handle_succ'))
        }
    } catch (e) {
        $message.error(e.message)
    }
    dialogStore.closeFlushDBDialog()
}

const onClose = () => {
    dialogStore.closeFlushDBDialog()
}
</script>

<template>
    <n-modal
        v-model:show="dialogStore.flushDBDialogVisible"
        :closable="false"
        :close-on-esc="false"
        :mask-closable="false"
        :show-icon="false"
        :title="$t('interface.flush_db')"
        preset="dialog"
        transform-origin="center">
        <n-form :model="flushForm" :show-require-mark="false" label-placement="top">
            <n-form-item :label="$t('dialogue.key.server')">
                <n-input :value="flushForm.server" readonly />
            </n-form-item>
            <n-form-item :label="$t('dialogue.key.db_index')">
                <n-input :value="flushForm.db.toString()" readonly />
            </n-form-item>
            <n-form-item :label="$t('dialogue.key.async_delete')" required>
                <n-checkbox v-model:checked="flushForm.async">{{ $t('dialogue.key.async_delete_title') }}</n-checkbox>
            </n-form-item>
            <n-form-item :label="$t('common.warning')" required>
                <n-checkbox v-model:checked="flushForm.confirm">
                    <span style="color: red; font-weight: bold">{{ $t('dialogue.key.confirm_flush') }}</span>
                </n-checkbox>
            </n-form-item>
        </n-form>

        <template #action>
            <n-button :focusable="false" @click="onClose">{{ $t('common.cancel') }}</n-button>
            <n-button :disabled="!!!flushForm.confirm" :focusable="false" type="primary" @click="onConfirmFlush">
                {{ $t('dialogue.key.confirm_flush_db') }}
            </n-button>
        </template>
    </n-modal>
</template>

<style lang="scss" scoped></style>
