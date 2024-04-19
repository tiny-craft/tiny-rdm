<script setup>
import { computed, nextTick, reactive, ref, watchEffect } from 'vue'
import useDialog from 'stores/dialog'
import { useI18n } from 'vue-i18n'
import { isEmpty, map, size } from 'lodash'
import useBrowserStore from 'stores/browser.js'
import { decodeRedisKey } from '@/utils/key_convert.js'

const deleteForm = reactive({
    server: '',
    db: 0,
    key: '',
    showAffected: false,
    loadingAffected: false,
    affectedKeys: [],
    async: true,
})

const dialogStore = useDialog()
const browserStore = useBrowserStore()

watchEffect(() => {
    if (dialogStore.deleteKeyDialogVisible) {
        const { server, db, key } = dialogStore.deleteKeyParam
        deleteForm.server = server
        deleteForm.db = db
        deleteForm.key = key
        deleteForm.loadingAffected = false
        // deleteForm.async = true
        loading.value = false
        deleting.value = false
        if (key instanceof Array) {
            deleteForm.showAffected = true
            deleteForm.affectedKeys = key
        } else {
            deleteForm.showAffected = false
            deleteForm.affectedKeys = []
        }
    }
})

const loading = ref(false)
const deleting = ref(false)
const scanAffectedKey = async () => {
    try {
        loading.value = true
        deleteForm.loadingAffected = true
        const { keys = [] } = await browserStore.scanKeys({
            server: deleteForm.server,
            db: deleteForm.db,
            match: deleteForm.key,
            loadType: 2,
        })
        deleteForm.affectedKeys = keys || []
        deleteForm.showAffected = true
    } finally {
        deleteForm.loadingAffected = false
        loading.value = false
    }
}

const resetAffected = () => {
    deleteForm.showAffected = false
    deleteForm.affectedKeys = []
}

const keyLines = computed(() => {
    return map(deleteForm.affectedKeys, (k) => decodeRedisKey(k))
})

const i18n = useI18n()
const onConfirmDelete = async () => {
    try {
        deleting.value = true
        const { server, db, key, affectedKeys } = deleteForm
        await nextTick()
        browserStore.deleteKeys(server, db, affectedKeys).catch((e) => {})
    } catch (e) {
        $message.error(e.message)
        return
    } finally {
        deleting.value = false
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
        :mask-closable="false"
        :show-icon="false"
        :title="$t('interface.batch_delete_key')"
        close-on-esc
        preset="dialog"
        transform-origin="center"
        @esc="onClose">
        <n-spin :show="loading">
            <n-form :model="deleteForm" :show-require-mark="false" label-placement="top">
                <n-grid :x-gap="10">
                    <n-form-item-gi :label="$t('dialogue.key.server')" :span="12">
                        <n-input :autofocus="false" :value="deleteForm.server" readonly />
                    </n-form-item-gi>
                    <n-form-item-gi :label="$t('dialogue.key.db_index')" :span="12">
                        <n-input :autofocus="false" :value="deleteForm.db.toString()" readonly />
                    </n-form-item-gi>
                </n-grid>
                <n-form-item
                    v-if="!(deleteForm.key instanceof Array)"
                    :label="$t('dialogue.key.key_expression')"
                    required>
                    <n-input v-model:value="deleteForm.key" placeholder="" @input="resetAffected" />
                </n-form-item>
                <!--                <n-checkbox v-model:checked="deleteForm.async">-->
                <!--                    {{ $t('dialogue.key.silent') }}-->
                <!--                </n-checkbox>-->
                <n-card
                    v-if="deleteForm.showAffected"
                    :title="$t('dialogue.key.affected_key') + `(${size(deleteForm.affectedKeys)})`"
                    embedded
                    size="small">
                    <n-skeleton v-if="deleteForm.loadingAffected" :repeat="10" text />
                    <n-virtual-list v-else :item-size="25" :items="keyLines" class="list-wrapper">
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
                <n-button :disabled="loading" :focusable="false" @click="onClose">{{ $t('common.cancel') }}</n-button>
                <n-button
                    v-if="!deleteForm.showAffected"
                    :focusable="false"
                    :loading="loading"
                    type="primary"
                    @click="scanAffectedKey">
                    {{ $t('dialogue.key.show_affected_key') }}
                </n-button>
                <n-button
                    v-else
                    :disabled="isEmpty(deleteForm.affectedKeys)"
                    :focusable="false"
                    :loading="loading"
                    type="primary"
                    @click="onConfirmDelete">
                    {{ $t('dialogue.key.confirm_delete_key', { num: size(deleteForm.affectedKeys) }) }}
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
