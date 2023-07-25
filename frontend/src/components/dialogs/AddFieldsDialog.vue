<script setup>
import { computed, reactive, watch } from 'vue'
import { types } from '../../consts/support_redis_type'
import useDialog from '../../stores/dialog'
import NewStringValue from '../new_value/NewStringValue.vue'
import NewSetValue from '../new_value/NewSetValue.vue'
import { useI18n } from 'vue-i18n'
import { useMessage } from 'naive-ui'
import AddListValue from '../new_value/AddListValue.vue'
import AddHashValue from '../new_value/AddHashValue.vue'
import AddZSetValue from '../new_value/AddZSetValue.vue'
import useConnectionStore from '../../stores/connections.js'
import NewStreamValue from '../new_value/NewStreamValue.vue'
import { size, slice } from 'lodash'

const i18n = useI18n()
const newForm = reactive({
    server: '',
    db: 0,
    key: '',
    type: '',
    opType: 0,
    value: null,
    reload: true,
})

const formLabelWidth = '60px'
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
            return i18n.t('new_item')
        case types.HASH:
            return i18n.t('new_field')
        case types.SET:
            return i18n.t('new_field')
        case types.ZSET:
            return i18n.t('new_field')
        case types.STREAM:
            return i18n.t('new_field')
    }
    return ''
})

const dialogStore = useDialog()
watch(
    () => dialogStore.addFieldsDialogVisible,
    (visible) => {
        if (visible) {
            const { server, db, key, type } = dialogStore.addFieldParam
            newForm.server = server
            newForm.db = db
            newForm.key = key
            newForm.type = type
            newForm.opType = 0
            newForm.value = null
        }
    },
)

const connectionStore = useConnectionStore()
const message = useMessage()
const onAdd = async () => {
    try {
        const { server, db, key, type } = newForm
        let { value } = newForm
        if (value == null) {
            value = defaultValue[type]
        }
        switch (type) {
            case types.LIST:
                {
                    let data
                    if (newForm.opType === 1) {
                        data = await connectionStore.prependListItem(server, db, key, value)
                    } else {
                        data = await connectionStore.appendListItem(server, db, key, value)
                    }
                    const { success, msg } = data
                    if (success) {
                        if (newForm.reload) {
                            connectionStore.loadKeyValue(server, db, key).then(() => {})
                        }
                        message.success(i18n.t('handle_succ'))
                    } else {
                        message.error(msg)
                    }
                }
                break

            case types.HASH:
                {
                    const { success, msg } = await connectionStore.addHashField(server, db, key, newForm.opType, value)
                    if (success) {
                        if (newForm.reload) {
                            connectionStore.loadKeyValue(server, db, key).then(() => {})
                        }
                        message.success(i18n.t('handle_succ'))
                    } else {
                        message.error(msg)
                    }
                }
                break

            case types.SET:
                {
                    const { success, msg } = await connectionStore.addSetItem(server, db, key, value)
                    if (success) {
                        if (newForm.reload) {
                            connectionStore.loadKeyValue(server, db, key).then(() => {})
                        }
                        message.success(i18n.t('handle_succ'))
                    } else {
                        message.error(msg)
                    }
                }
                break

            case types.ZSET:
                {
                    const { success, msg } = await connectionStore.addZSetItem(server, db, key, newForm.opType, value)
                    if (success) {
                        if (newForm.reload) {
                            connectionStore.loadKeyValue(server, db, key).then(() => {})
                        }
                        message.success(i18n.t('handle_succ'))
                    } else {
                        message.error(msg)
                    }
                }
                break

            case types.STREAM:
                {
                    if (size(value) > 2) {
                        const { success, msg } = await connectionStore.addStreamValue(
                            server,
                            db,
                            key,
                            value[0],
                            slice(value, 1),
                        )
                        if (success) {
                            if (newForm.reload) {
                                connectionStore.loadKeyValue(server, db, key).then(() => {})
                            }
                            message.success(i18n.t('handle_succ'))
                        } else {
                            message.error(msg)
                        }
                    }
                }
                break
        }
        dialogStore.closeAddFieldsDialog()
    } catch (e) {
        message.error(e.message)
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
        :close-on-esc="false"
        :mask-closable="false"
        :negative-button-props="{ size: 'medium' }"
        :negative-text="$t('cancel')"
        :positive-button-props="{ size: 'medium' }"
        :positive-text="$t('confirm')"
        :show-icon="false"
        :title="title"
        preset="dialog"
        style="width: 600px"
        transform-origin="center"
        @positive-click="onAdd"
        @negative-click="onClose"
    >
        <n-scrollbar style="max-height: 500px">
            <n-form
                :label-width="formLabelWidth"
                :model="newForm"
                :show-require-mark="false"
                label-align="right"
                label-placement="left"
                style="padding-right: 15px"
            >
                <n-form-item :label="$t('key')" path="key" required>
                    <n-input v-model:value="newForm.key" placeholder="" readonly />
                </n-form-item>
                <component
                    :is="addValueComponent[newForm.type]"
                    v-model:type="newForm.opType"
                    v-model:value="newForm.value"
                />
                <n-form-item label=" " path="key" required>
                    <n-checkbox v-model:checked="newForm.reload">
                        {{ $t('reload_when_succ') }}
                    </n-checkbox>
                </n-form-item>
            </n-form>
        </n-scrollbar>
    </n-modal>
</template>

<style lang="scss" scoped></style>
