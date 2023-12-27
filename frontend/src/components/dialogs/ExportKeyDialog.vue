<script setup>
import { computed, reactive, ref, watchEffect } from 'vue'
import useDialog from 'stores/dialog'
import { useI18n } from 'vue-i18n'
import useBrowserStore from 'stores/browser.js'
import FileSaveInput from '@/components/common/FileSaveInput.vue'
import { isEmpty, map, size } from 'lodash'
import { decodeRedisKey } from '@/utils/key_convert.js'
import dayjs from 'dayjs'

const exportKeyForm = reactive({
    server: '',
    db: 0,
    keys: [],
    file: '',
})

const dialogStore = useDialog()
const browserStore = useBrowserStore()
const loading = ref(false)
const exporting = ref(false)
watchEffect(() => {
    if (dialogStore.exportKeyDialogVisible) {
        const { server, db, keys } = dialogStore.exportKeyParam
        exportKeyForm.server = server
        exportKeyForm.db = db
        exportKeyForm.keys = keys
        exporting.value = false
    }
})

const keyLines = computed(() => {
    return map(exportKeyForm.keys, (k) => decodeRedisKey(k))
})

const exportEnable = computed(() => {
    return !isEmpty(exportKeyForm.keys) && !isEmpty(exportKeyForm.file)
})

const i18n = useI18n()
const onConfirmExport = async () => {
    try {
        exporting.value = true
        const { server, db, keys, file } = exportKeyForm
        browserStore.exportKeys(server, db, keys, file).catch((e) => {})
    } catch (e) {
        $message.error(e.message)
        return
    } finally {
        exporting.value = false
    }
    dialogStore.closeExportKeyDialog()
}

const onClose = () => {
    dialogStore.closeExportKeyDialog()
}
</script>

<template>
    <n-modal
        v-model:show="dialogStore.exportKeyDialogVisible"
        :closable="false"
        :close-on-esc="false"
        :mask-closable="false"
        :show-icon="false"
        :title="$t('dialogue.export.name')"
        preset="dialog"
        transform-origin="center">
        <n-spin :show="loading">
            <n-form :model="exportKeyForm" :show-require-mark="false" label-placement="top">
                <n-grid :x-gap="10">
                    <n-form-item-gi :label="$t('dialogue.key.server')" :span="12">
                        <n-input :autofocus="false" :value="exportKeyForm.server" readonly />
                    </n-form-item-gi>
                    <n-form-item-gi :label="$t('dialogue.key.db_index')" :span="12">
                        <n-input :autofocus="false" :value="exportKeyForm.db.toString()" readonly />
                    </n-form-item-gi>
                </n-grid>
                <n-form-item :label="$t('dialogue.export.save_file')" required>
                    <file-save-input
                        v-model:value="exportKeyForm.file"
                        :default-path="`export_${dayjs().format('YYYYMMDDHHmmss')}.csv`"
                        :placeholder="$t('dialogue.export.save_file_tip')" />
                </n-form-item>
                <n-card
                    :title="$t('dialogue.key.affected_key') + `(${size(exportKeyForm.keys)})`"
                    embedded
                    size="small">
                    <n-log :line-height="1.5" :lines="keyLines" :rows="10" style="user-select: text; cursor: text" />
                </n-card>
            </n-form>
        </n-spin>

        <template #action>
            <div class="flex-item n-dialog__action">
                <n-button :disabled="loading" :focusable="false" @click="onClose">{{ $t('common.cancel') }}</n-button>
                <n-button
                    :disabled="!exportEnable"
                    :focusable="false"
                    :loading="loading"
                    type="error"
                    @click="onConfirmExport">
                    {{ $t('dialogue.export.export') }}
                </n-button>
            </div>
        </template>
    </n-modal>
</template>

<style lang="scss" scoped></style>
