<script setup>
import { reactive, watch } from 'vue'
import useDialog from 'stores/dialog'
import { useI18n } from 'vue-i18n'
import useConnectionStore from 'stores/connections.js'
import { isEmpty, size } from 'lodash'

const deleteForm = reactive({
    server: '',
    db: 0,
    key: '',
    showAffected: false,
    loadingAffected: false,
    affectedKeys: [],
})

const dialogStore = useDialog()
const connectionStore = useConnectionStore()
watch(
    () => dialogStore.deleteKeyDialogVisible,
    (visible) => {
        if (visible) {
            const { server, db, key } = dialogStore.deleteKeyParam
            deleteForm.server = server
            deleteForm.db = db
            deleteForm.key = key
            deleteForm.showAffected = false
            deleteForm.loadingAffected = false
            deleteForm.affectedKeys = []
        }
    },
)

const scanAffectedKey = async () => {
    try {
        deleteForm.loadingAffected = true
        const { keys = [] } = await connectionStore.scanKeys(deleteForm.server, deleteForm.db, deleteForm.key)
        deleteForm.affectedKeys = keys || []
        deleteForm.showAffected = true
    } finally {
        deleteForm.loadingAffected = false
    }
}

const resetAffected = () => {
    deleteForm.showAffected = false
    deleteForm.affectedKeys = []
}

const i18n = useI18n()
const onConfirmDelete = async () => {
    try {
        const { server, db, key } = deleteForm
        const success = await connectionStore.deleteKeyPrefix(server, db, key)
        if (success) {
            $message.success(i18n.t('dialogue.handle_succ'))
        }
    } catch (e) {
        $message.error(e.message)
    }
    dialogStore.closeDeleteKeyDialog()
}

const onClose = () => {
    dialogStore.closeDeleteKeyDialog()
}
</script>

<template>
    <n-modal
        v-model:show="dialogStore.deleteKeyDialogVisible"
        :closable="false"
        :close-on-esc="false"
        :mask-closable="false"
        :show-icon="false"
        :title="$t('interface.batch_delete_key')"
        preset="dialog"
        transform-origin="center">
        <n-form :model="deleteForm" :show-require-mark="false" label-placement="top">
            <n-form-item :label="$t('dialogue.key.server')">
                <n-input :value="deleteForm.server" readonly />
            </n-form-item>
            <n-form-item :label="$t('dialogue.key.db_index')">
                <n-input :value="deleteForm.db.toString()" readonly />
            </n-form-item>
            <n-form-item :label="$t('dialogue.key.key_expression')" required>
                <n-input v-model:value="deleteForm.key" placeholder="" @input="resetAffected" />
            </n-form-item>
            <n-card v-if="deleteForm.showAffected" :title="$t('dialogue.key.affected_key')" size="small">
                <n-skeleton v-if="deleteForm.loadingAffected" :repeat="10" text />
                <n-log
                    v-else
                    :line-height="1.5"
                    :lines="deleteForm.affectedKeys"
                    :rows="10"
                    style="user-select: text; cursor: text" />
            </n-card>
        </n-form>

        <template #action>
            <div class="flex-item n-dialog__action">
                <n-button :focusable="false" @click="onClose">{{ $t('common.cancel') }}</n-button>
                <n-button v-if="!deleteForm.showAffected" type="primary" :focusable="false" @click="scanAffectedKey">
                    {{ $t('dialogue.key.show_affected_key') }}
                </n-button>
                <n-button
                    v-else
                    :focusable="false"
                    :disabled="isEmpty(deleteForm.affectedKeys)"
                    type="error"
                    @click="onConfirmDelete">
                    {{ $t('dialogue.key.confirm_delete_key', { num: size(deleteForm.affectedKeys) }) }}
                </n-button>
            </div>
        </template>
    </n-modal>
</template>

<style lang="scss" scoped></style>
