<script setup>
import { computed, reactive, ref, watchEffect } from 'vue'
import useDialog from 'stores/dialog'
import useBrowserStore from 'stores/browser.js'
import FileSaveInput from '@/components/common/FileSaveInput.vue'
import { isEmpty, map, size } from 'lodash'
import { decodeRedisKey } from '@/utils/key_convert.js'
import dayjs from 'dayjs'

const exportKeyForm = reactive({
    server: '',
    db: 0,
    expire: false,
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
        exportKeyForm.ttl = false
        exportKeyForm.keys = keys
        exportKeyForm.file = ''
        exporting.value = false
    }
})

const keyLines = computed(() => {
    return map(exportKeyForm.keys, (k) => decodeRedisKey(k))
})

const exportEnable = computed(() => {
    return !isEmpty(exportKeyForm.keys) && !isEmpty(exportKeyForm.file)
})

const onConfirmExport = async () => {
    try {
        exporting.value = true
        const { server, db, keys, file, expire } = exportKeyForm
        browserStore.exportKeys(server, db, keys, file, expire).catch((e) => {})
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
        :mask-closable="false"
        :show-icon="false"
        :title="$t('dialogue.export.name')"
        close-on-esc
        preset="dialog"
        transform-origin="center"
        @esc="onClose">
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
                <n-form-item :label="$t('dialogue.export.export_expire_title')">
                    <n-checkbox v-model:checked="exportKeyForm.expire" :autofocus="false">
                        {{ $t('dialogue.export.export_expire') }}
                    </n-checkbox>
                </n-form-item>
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
                    <n-virtual-list :item-size="25" :items="keyLines" class="list-wrapper">
                        <template #default="{ item }">
                            <div class="line-item content-value">
                                {{ item }}
                            </div>
                        </template>
                    </n-virtual-list>
                </n-card>
            </n-form>
        </n-spin>

        <template #action>
            <div class="flex-item n-dialog__action">
                <n-button :disabled="loading" :focusable="false" @click="onClose">
                    {{ $t('common.cancel') }}
                </n-button>
                <n-button
                    :disabled="!exportEnable"
                    :focusable="false"
                    :loading="loading"
                    type="primary"
                    @click="onConfirmExport">
                    {{ $t('dialogue.export.export') }}
                </n-button>
            </div>
        </template>
    </n-modal>
</template>

<style lang="scss" scoped>
.line-item {
    line-height: 1.6;
}

.list-wrapper {
    box-sizing: border-box;
    max-height: 180px;
    user-select: text;
    cursor: text;
}
</style>
