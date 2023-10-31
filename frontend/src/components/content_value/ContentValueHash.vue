<script setup>
import { computed, h, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import ContentToolbar from './ContentToolbar.vue'
import AddLink from '@/components/icons/AddLink.vue'
import { NButton, NCode, NIcon, NInput, useThemeVars } from 'naive-ui'
import { types, types as redisTypes } from '@/consts/support_redis_type.js'
import EditableTableColumn from '@/components/common/EditableTableColumn.vue'
import useDialogStore from 'stores/dialog.js'
import useConnectionStore from 'stores/connections.js'
import { isEmpty } from 'lodash'
import bytes from 'bytes'
import { decodeTypes, formatTypes } from '@/consts/value_view_type.js'

const i18n = useI18n()
const themeVars = useThemeVars()

const props = defineProps({
    name: String,
    db: Number,
    keyPath: String,
    keyCode: {
        type: Array,
        default: null,
    },
    ttl: {
        type: Number,
        default: -1,
    },
    value: Object,
    size: Number,
    length: Number,
    viewAs: {
        type: String,
        default: formatTypes.PLAIN_TEXT,
    },
    decode: {
        type: String,
        default: decodeTypes.NONE,
    },
})

/**
 *
 * @type {ComputedRef<string|number[]>}
 */
const keyName = computed(() => {
    return !isEmpty(props.keyCode) ? props.keyCode : props.keyPath
})

const filterOption = [
    {
        value: 1,
        label: i18n.t('common.field'),
    },
    {
        value: 2,
        label: i18n.t('common.value'),
    },
]
const filterType = ref(1)

const connectionStore = useConnectionStore()
const dialogStore = useDialogStore()
const keyType = redisTypes.HASH
const currentEditRow = ref({
    no: 0,
    key: '',
    value: null,
})
const fieldColumn = reactive({
    key: 'key',
    title: i18n.t('common.field'),
    align: 'center',
    titleAlign: 'center',
    resizable: true,
    filterOptionValue: null,
    filter(value, row) {
        return !!~row.key.indexOf(value.toString())
    },
    // sorter: (row1, row2) => row1.key - row2.key,
    render: (row) => {
        const isEdit = currentEditRow.value.no === row.no
        if (isEdit) {
            return h(NInput, {
                value: currentEditRow.value.key,
                'onUpdate:value': (val) => {
                    currentEditRow.value.key = val
                },
            })
        } else {
            return row.key
        }
    },
})
const valueColumn = reactive({
    key: 'value',
    title: i18n.t('common.value'),
    align: 'center',
    titleAlign: 'center',
    resizable: true,
    filterOptionValue: null,
    filter(value, row) {
        return !!~row.value.indexOf(value.toString())
    },
    // sorter: (row1, row2) => row1.value - row2.value,
    // ellipsis: {
    //     tooltip: true
    // },
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
    title: i18n.t('interface.action'),
    width: 100,
    align: 'center',
    titleAlign: 'center',
    fixed: 'right',
    render: (row) => {
        return h(EditableTableColumn, {
            editing: currentEditRow.value.no === row.no,
            bindKey: row.key,
            onEdit: () => {
                currentEditRow.value.no = row.no
                currentEditRow.value.key = row.key
                currentEditRow.value.value = row.value
            },
            onDelete: async () => {
                try {
                    const { success, msg } = await connectionStore.removeHashField(
                        props.name,
                        props.db,
                        keyName.value,
                        row.key,
                    )
                    if (success) {
                        connectionStore.loadKeyValue(props.name, props.db, keyName.value).then((r) => {})
                        $message.success(i18n.t('dialogue.delete_key_succ', { key: row.key }))
                        // update display value
                        // if (!isEmpty(removed)) {
                        //     for (const elem of removed) {
                        //         delete props.value[elem]
                        //     }
                        // }
                    } else {
                        $message.error(msg)
                    }
                } catch (e) {
                    $message.error(e.message)
                }
            },
            onSave: async () => {
                try {
                    const { success, msg } = await connectionStore.setHash(
                        props.name,
                        props.db,
                        keyName.value,
                        row.key,
                        currentEditRow.value.key,
                        currentEditRow.value.value,
                    )
                    if (success) {
                        connectionStore.loadKeyValue(props.name, props.db, keyName.value).then((r) => {})
                        $message.success(i18n.t('dialogue.save_value_succ'))
                        // update display value
                        // if (!isEmpty(updated)) {
                        //     for (const key in updated) {
                        //         props.value[key] = updated[key]
                        //     }
                        // }
                    } else {
                        $message.error(msg)
                    }
                } catch (e) {
                    $message.error(e.message)
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
const columns = reactive([
    {
        key: 'no',
        title: '#',
        width: 80,
        align: 'center',
        titleAlign: 'center',
    },
    fieldColumn,
    valueColumn,
    actionColumn,
])

const tableData = computed(() => {
    const data = []
    let index = 0
    for (const key in props.value) {
        data.push({
            no: ++index,
            key,
            value: props.value[key],
        })
    }
    return data
})
const onAddRow = () => {
    dialogStore.openAddFieldsDialog(props.name, props.db, props.keyPath, props.keyCode, types.HASH)
}

const filterValue = ref('')
const onFilterInput = (val) => {
    switch (filterType.value) {
        case filterOption[0].value:
            // filter field
            valueColumn.filterOptionValue = null
            fieldColumn.filterOptionValue = val
            break
        case filterOption[1].value:
            // filter value
            fieldColumn.filterOptionValue = null
            valueColumn.filterOptionValue = val
            break
    }
}

const onChangeFilterType = (type) => {
    onFilterInput(filterValue.value)
}

const clearFilter = () => {
    fieldColumn.filterOptionValue = null
    valueColumn.filterOptionValue = null
}

const onUpdateFilter = (filters, sourceColumn) => {
    switch (filterType.value) {
        case filterOption[0].value:
            fieldColumn.filterOptionValue = filters[sourceColumn.key]
            break
        case filterOption[1].value:
            valueColumn.filterOptionValue = filters[sourceColumn.key]
            break
    }
}
</script>

<template>
    <div class="content-wrapper flex-box-v">
        <content-toolbar
            :db="props.db"
            :decode="props.decode"
            :key-code="props.keyCode"
            :key-path="props.keyPath"
            :key-type="keyType"
            :server="props.name"
            :ttl="ttl"
            :view-as="props.viewAs"
            class="value-item-part" />
        <div class="tb2 value-item-part flex-box-h">
            <div class="flex-box-h">
                <n-input-group>
                    <n-select
                        v-model:value="filterType"
                        :consistent-menu-width="false"
                        :options="filterOption"
                        style="width: 120px"
                        @update:value="onChangeFilterType" />
                    <n-input
                        v-model:value="filterValue"
                        :placeholder="$t('interface.search')"
                        clearable
                        @clear="clearFilter"
                        @update:value="onFilterInput" />
                </n-input-group>
            </div>
            <div class="flex-item-expand"></div>
            <n-button :focusable="false" plain @click="onAddRow">
                <template #icon>
                    <n-icon :component="AddLink" size="18" />
                </template>
                {{ $t('interface.add_row') }}
            </n-button>
        </div>
        <div class="value-wrapper value-item-part fill-height flex-box-h">
            <n-data-table
                :key="(row) => row.no"
                :bordered="false"
                :bottom-bordered="false"
                :columns="columns"
                :data="tableData"
                :single-column="true"
                :single-line="false"
                flex-height
                max-height="100%"
                size="small"
                striped
                virtual-scroll
                @update:filters="onUpdateFilter" />
        </div>
        <div class="value-footer flex-box-h">
            <n-text v-if="!isNaN(props.length)">{{ $t('interface.entries') }}: {{ props.length }}</n-text>
            <n-divider v-if="!isNaN(props.length)" vertical />
            <n-text v-if="!isNaN(props.size)">{{ $t('interface.memory_usage') }}: {{ bytes(props.size) }}</n-text>
            <div class="flex-item-expand"></div>
        </div>
    </div>
</template>

<style lang="scss" scoped>
.value-footer {
    border-top: v-bind('themeVars.borderColor') 1px solid;
    background-color: v-bind('themeVars.tableHeaderColor');
}
</style>
