<script setup>
import { computed, reactive, ref, watchEffect } from 'vue'
import useDialog from 'stores/dialog'
import { useI18n } from 'vue-i18n'
import useBrowserStore from 'stores/browser.js'
import { isEmpty } from 'lodash'
import FileOpenInput from '@/components/common/FileOpenInput.vue'

const importKeyForm = reactive({
    server: '',
    db: 0,
    expire: true,
    reload: true,
    file: '',
    type: 0,
    conflict: 0,
})

const dialogStore = useDialog()
const browserStore = useBrowserStore()
const loading = ref(false)
const importing = ref(false)
watchEffect(() => {
    if (dialogStore.importKeyDialogVisible) {
        const { server, db } = dialogStore.importKeyParam
        importKeyForm.server = server
        importKeyForm.db = db
        importKeyForm.expire = true
        importKeyForm.reload = true
        importKeyForm.file = ''
        importKeyForm.type = 0
        importKeyForm.conflict = 0
        importing.value = false
    }
})

const i18n = useI18n()
const conflictOption = computed(() => [
    {
        value: 0,
        label: i18n.t('dialogue.import.conflict_overwrite'),
    },
    {
        value: 1,
        label: i18n.t('dialogue.import.conflict_ignore'),
    },
])

const importEnable = computed(() => {
    return !isEmpty(importKeyForm.file)
})

const onConfirmImport = async () => {
    try {
        importing.value = true
        const { server, db, file, conflict, expire, reload } = importKeyForm
        browserStore.importKeysFromCSVFile(server, db, file, conflict, expire, reload).catch((e) => {})
    } catch (e) {
        $message.error(e.message)
        return
    } finally {
        importing.value = false
    }
    dialogStore.closeImportKeyDialog()
}

const onClose = () => {
    dialogStore.closeImportKeyDialog()
}
</script>

<template>
    <n-modal
        v-model:show="dialogStore.importKeyDialogVisible"
        :closable="false"
        :close-on-esc="false"
        :mask-closable="false"
        :show-icon="false"
        :title="$t('dialogue.import.name')"
        preset="dialog"
        transform-origin="center">
        <n-spin :show="loading">
            <n-form :model="importKeyForm" :show-require-mark="false" label-placement="top">
                <n-grid :x-gap="10">
                    <n-form-item-gi :label="$t('dialogue.key.server')" :span="12">
                        <n-input :autofocus="false" :value="importKeyForm.server" readonly />
                    </n-form-item-gi>
                    <n-form-item-gi :label="$t('dialogue.key.db_index')" :span="12">
                        <n-input :autofocus="false" :value="importKeyForm.db.toString()" readonly />
                    </n-form-item-gi>
                </n-grid>
                <n-form-item :label="$t('dialogue.import.open_csv_file')" required>
                    <file-open-input
                        v-model:value="importKeyForm.file"
                        :placeholder="$t('dialogue.import.open_csv_file_tip')"
                        ext="csv" />
                </n-form-item>
                <n-form-item :label="$t('dialogue.import.conflict_handle')">
                    <n-radio-group v-model:value="importKeyForm.conflict">
                        <n-radio-button
                            v-for="(op, i) in conflictOption"
                            :key="i"
                            :label="op.label"
                            :value="op.value" />
                    </n-radio-group>
                </n-form-item>
                <n-form-item :label="$t('dialogue.import.import_expire_title')" :show-label="false">
                    <n-space :wrap-item="false">
                        <n-checkbox v-model:checked="importKeyForm.expire" :autofocus="false">
                            {{ $t('dialogue.import.import_expire') }}
                        </n-checkbox>
                        <n-checkbox v-model:checked="importKeyForm.reload" :autofocus="false">
                            {{ $t('dialogue.import.reload') }}
                        </n-checkbox>
                    </n-space>
                </n-form-item>
            </n-form>
        </n-spin>

        <template #action>
            <div class="flex-item n-dialog__action">
                <n-button :disabled="loading" :focusable="false" @click="onClose">{{ $t('common.cancel') }}</n-button>
                <n-button
                    :disabled="!importEnable"
                    :focusable="false"
                    :loading="loading"
                    type="primary"
                    @click="onConfirmImport">
                    {{ $t('dialogue.export.export') }}
                </n-button>
            </div>
        </template>
    </n-modal>
</template>

<style lang="scss" scoped></style>
