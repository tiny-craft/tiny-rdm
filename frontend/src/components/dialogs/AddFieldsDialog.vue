<script setup>
import { computed, reactive, watchEffect } from 'vue'
import { types } from '@/consts/support_redis_type.js'
import useDialog from 'stores/dialog'
import NewStringValue from '@/components/new_value/NewStringValue.vue'
import NewSetValue from '@/components/new_value/NewSetValue.vue'
import { useI18n } from 'vue-i18n'
import AddListValue from '@/components/new_value/AddListValue.vue'
import AddHashValue from '@/components/new_value/AddHashValue.vue'
import AddZSetValue from '@/components/new_value/AddZSetValue.vue'
import NewStreamValue from '@/components/new_value/NewStreamValue.vue'
import { get, isEmpty, size, slice } from 'lodash'
import useBrowserStore from 'stores/browser.js'
import useTabStore from 'stores/tab.js'

const i18n = useI18n()
const newForm = reactive({
    server: '',
    db: 0,
    key: '',
    keyCode: null,
    type: '',
    opType: 0,
    value: null,
    reload: false,
})

const addValueComponent = {
    [types.STRING]: NewStringValue,
    [types.HASH]: AddHashValue,
    [types.LIST]: AddListValue,
    [types.SET]: NewSetValue,
    [types.ZSET]: AddZSetValue,
    [types.STREAM]: NewStreamValue,
}
const defaultValue = {
    [types.STRING]: '',
    [types.HASH]: [],
    [types.LIST]: [],
    [types.SET]: [],
    [types.ZSET]: [],
    [types.STREAM]: ['*'],
}

/**
 * dialog title
 * @type {ComputedRef<string>}
 */
const title = computed(() => {
    switch (newForm.type) {
        case types.LIST:
            return 'dialogue.field.new_item'
        case types.HASH:
            return 'dialogue.field.new'
        case types.SET:
            return 'dialogue.field.new'
        case types.ZSET:
            return 'dialogue.field.new'
        case types.STREAM:
            return 'dialogue.field.new'
    }
    return ''
})

const dialogStore = useDialog()
watchEffect(() => {
    if (dialogStore.addFieldsDialogVisible) {
        const { server, db, key, keyCode, type } = dialogStore.addFieldParam
        newForm.server = server
        newForm.db = db
        newForm.key = key
        newForm.keyCode = keyCode
        newForm.type = type
        newForm.opType = 0
        newForm.value = null
    }
})

const browserStore = useBrowserStore()
const tab = useTabStore()
const onAdd = async () => {
    try {
        const { server, db, key, keyCode, type } = newForm
        let { value } = newForm
        if (value == null) {
            value = defaultValue[type]
        }
        const keyName = isEmpty(keyCode) ? key : keyCode
        let success = false
        let msg = ''
        switch (type) {
            case types.LIST:
                {
                    let data
                    if (newForm.opType === 1) {
                        data = await browserStore.prependListItem({
                            server,
                            db,
                            key: keyName,
                            values: value,
                            reload: newForm.reload,
                        })
                    } else {
                        data = await browserStore.appendListItem({
                            server,
                            db,
                            key: keyName,
                            values: value,
                            reload: newForm.reload,
                        })
                    }
                    success = get(data, 'success')
                    msg = get(data, 'msg')
                }
                break

            case types.HASH:
                {
                    const data = await browserStore.addHashField({
                        server,
                        db,
                        key: keyName,
                        action: newForm.opType,
                        fieldItems: value,
                        reload: newForm.reload,
                    })
                    success = get(data, 'success')
                    msg = get(data, 'msg')
                }
                break

            case types.SET:
                {
                    const data = await browserStore.addSetItem({
                        server,
                        db,
                        key: keyName,
                        value,
                        reload: newForm.reload,
                    })
                    success = get(data, 'success')
                    msg = get(data, 'msg')
                }
                break

            case types.ZSET:
                {
                    const data = await browserStore.addZSetItem({
                        server,
                        db,
                        key: keyName,
                        action: newForm.opType,
                        vs: value,
                        reload: newForm.reload,
                    })
                    success = get(data, 'success')
                    msg = get(data, 'msg')
                }
                break

            case types.STREAM:
                {
                    if (size(value) > 2) {
                        const data = await browserStore.addStreamValue({
                            server,
                            db,
                            key: keyName,
                            id: value[0],
                            values: slice(value, 1),
                            reload: newForm.reload,
                        })
                        success = get(data, 'success')
                        msg = get(data, 'msg')
                    }
                }
                break
        }

        if (success) {
            $message.success(i18n.t('dialogue.handle_succ'))
        } else if (!isEmpty(msg)) {
            $message.error(msg)
        }

        dialogStore.closeAddFieldsDialog()
    } catch (e) {
        $message.error(e.message)
    }
}

const onClose = () => {
    dialogStore.closeAddFieldsDialog()
}
</script>

<template>
    <n-modal
        v-model:show="dialogStore.addFieldsDialogVisible"
        :closable="false"
        :mask-closable="false"
        :negative-button-props="{ size: 'medium' }"
        :negative-text="$t('common.cancel')"
        :positive-button-props="{ size: 'medium' }"
        :positive-text="$t('common.confirm')"
        :show-icon="false"
        :title="title ? $t(title) : ''"
        close-on-esc
        preset="dialog"
        style="width: 600px"
        transform-origin="center"
        @esc="onClose"
        @positive-click="onAdd"
        @negative-click="onClose">
        <n-scrollbar style="max-height: 500px">
            <n-form :model="newForm" :show-require-mark="false" label-placement="top" style="padding-right: 15px">
                <n-form-item :label="$t('common.key')" path="key" required>
                    <n-input v-model:value="newForm.key" placeholder="" readonly />
                </n-form-item>
                <component
                    :is="addValueComponent[newForm.type]"
                    v-model:type="newForm.opType"
                    v-model:value="newForm.value" />
                <n-form-item :show-label="false" path="key" required>
                    <n-checkbox v-model:checked="newForm.reload">
                        {{ $t('dialogue.field.reload_when_succ') }}
                    </n-checkbox>
                </n-form-item>
            </n-form>
        </n-scrollbar>
    </n-modal>
</template>

<style lang="scss" scoped></style>
