<script setup>
import { computed, h, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import ContentToolbar from '../ContentToolbar.vue'
import AddLink from '../icons/AddLink.vue'
import { NButton, NCode, NIcon, NInput, useMessage } from 'naive-ui'
import { size } from 'lodash'
import useConnectionStore from '../../stores/connection.js'
import useDialogStore from '../../stores/dialog.js'
import { types, types as redisTypes } from '../../consts/support_redis_type.js'
import EditableTableColumn from '../EditableTableColumn.vue'

const i18n = useI18n()
const props = defineProps({
    name: String,
    db: Number,
    keyPath: String,
    ttl: {
        type: Number,
        default: -1,
    },
    value: Array,
})

const connectionStore = useConnectionStore()
const dialogStore = useDialogStore()
const keyType = redisTypes.LIST
const currentEditRow = ref({
    no: 0,
    value: null,
})

const valueColumn = reactive({
    key: 'value',
    title: i18n.t('value'),
    align: 'center',
    titleAlign: 'center',
    filterOptionValue: null,
    filter(value, row) {
        return !!~row.value.indexOf(value.toString())
    },
    render: (row) => {
        const isEdit = currentEditRow.value.no === row.no
        if (isEdit) {
            return h(NInput, {
                value: currentEditRow.value.value,
                type: 'textarea',
                autosize: { minRow: 2, maxRows: 5 },
                style: 'text-align: left;',
                'onUpdate:value': (val) => {
                    currentEditRow.value.value = val
                },
            })
        } else {
            return h(NCode, { language: 'plaintext', wordWrap: true }, { default: () => row.value })
        }
    },
})
const actionColumn = {
    key: 'action',
    title: i18n.t('action'),
    width: 100,
    align: 'center',
    titleAlign: 'center',
    fixed: 'right',
    render: (row) => {
        return h(EditableTableColumn, {
            editing: currentEditRow.value.no === row.no,
            bindKey: row.value,
            onEdit: () => {
                currentEditRow.value.no = row.no
                currentEditRow.value.key = row.key
                currentEditRow.value.value = row.value
            },
            onDelete: async () => {
                try {
                    const { success, msg } = await connectionStore.removeSetItem(
                        props.name,
                        props.db,
                        props.keyPath,
                        row.value
                    )
                    if (success) {
                        connectionStore.loadKeyValue(props.name, props.db, props.keyPath).then((r) => {})
                        message.success(i18n.t('delete_key_succ', { key: row.value }))
                        // update display value
                        // props.value.splice(row.no - 1, 1)
                    } else {
                        message.error(msg)
                    }
                } catch (e) {
                    message.error(e.message)
                }
            },
            onSave: async () => {
                try {
                    const { success, msg } = await connectionStore.updateSetItem(
                        props.name,
                        props.db,
                        props.keyPath,
                        row.value,
                        currentEditRow.value.value
                    )
                    if (success) {
                        connectionStore.loadKeyValue(props.name, props.db, props.keyPath).then((r) => {})
                        message.success(i18n.t('save_value_succ'))
                        // update display value
                        // props.value[row.no - 1] = currentEditRow.value.value
                    } else {
                        message.error(msg)
                    }
                } catch (e) {
                    message.error(e.message)
                } finally {
                    currentEditRow.value.no = 0
                }
            },
            onCancel: () => {
                currentEditRow.value.no = 0
            },
        })
    },
}
const columns = computed(() => {
    return [
        {
            key: 'no',
            title: '#',
            width: 80,
            align: 'center',
            titleAlign: 'center',
        },
        valueColumn,
        actionColumn,
    ]
})

const tableData = computed(() => {
    const data = []
    const len = size(props.value)
    for (let i = 0; i < len; i++) {
        data.push({
            no: i + 1,
            value: props.value[i],
        })
    }
    return data
})

const message = useMessage()
const onAddValue = (value) => {
    dialogStore.openAddFieldsDialog(props.name, props.db, props.keyPath, types.SET)
}

const filterValue = ref('')
const onFilterInput = (val) => {
    valueColumn.filterOptionValue = val
}

const clearFilter = () => {
    valueColumn.filterOptionValue = null
}

const onUpdateFilter = (filters, sourceColumn) => {
    valueColumn.filterOptionValue = filters[sourceColumn.key]
}
</script>

<template>
    <div class="content-wrapper flex-box-v">
        <content-toolbar :db="props.db" :key-path="props.keyPath" :key-type="keyType" :server="props.name" :ttl="ttl" />
        <div class="tb2 flex-box-h">
            <div class="flex-box-h">
                <n-input
                    v-model:value="filterValue"
                    :placeholder="$t('search')"
                    clearable
                    @clear="clearFilter"
                    @update:value="onFilterInput"
                />
            </div>
            <div class="flex-item-expand"></div>
            <n-button plain @click="onAddValue">
                <template #icon>
                    <n-icon :component="AddLink" size="18" />
                </template>
                {{ $t('add_row') }}
            </n-button>
        </div>
        <div class="fill-height flex-box-h" style="user-select: text">
            <n-data-table
                :key="(row) => row.no"
                :columns="columns"
                :data="tableData"
                :single-column="true"
                :single-line="false"
                flex-height
                max-height="100%"
                size="small"
                striped
                virtual-scroll
                @update:filters="onUpdateFilter"
            />
        </div>
    </div>
</template>

<style lang="scss" scoped></style>
