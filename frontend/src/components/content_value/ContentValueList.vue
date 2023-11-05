<script setup>
import { computed, h, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import ContentToolbar from './ContentToolbar.vue'
import AddLink from '@/components/icons/AddLink.vue'
import { NButton, NCode, NIcon, NInput, useThemeVars } from 'naive-ui'
import { isEmpty, size } from 'lodash'
import { types, types as redisTypes } from '@/consts/support_redis_type.js'
import EditableTableColumn from '@/components/common/EditableTableColumn.vue'
import useDialogStore from 'stores/dialog.js'
import bytes from 'bytes'
import { decodeTypes, formatTypes } from '@/consts/value_view_type.js'
import useBrowserStore from 'stores/browser.js'

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

const browserStore = useBrowserStore()
const dialogStore = useDialogStore()
const keyType = redisTypes.LIST
const currentEditRow = ref({
    no: 0,
    value: null,
})
const valueColumn = reactive({
    key: 'value',
    title: i18n.t('common.value'),
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
    title: i18n.t('interface.action'),
    width: 100,
    align: 'center',
    titleAlign: 'center',
    fixed: 'right',
    render: (row) => {
        return h(EditableTableColumn, {
            editing: currentEditRow.value.no === row.no,
            bindKey: '#' + row.no,
            onEdit: () => {
                currentEditRow.value.no = row.no
                currentEditRow.value.value = row.value
            },
            onDelete: async () => {
                try {
                    const { success, msg } = await browserStore.removeListItem(
                        props.name,
                        props.db,
                        keyName.value,
                        row.no - 1,
                    )
                    if (success) {
                        browserStore.loadKeyValue(props.name, props.db, keyName.value).then((r) => {})
                        $message.success(i18n.t('dialogue.delete_key_succ', { key: '#' + row.no }))
                        // update display value
                        // if (!isEmpty(removed)) {
                        //     props.value.splice(removed[0], 1)
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
                    const { success, msg } = await browserStore.updateListItem(
                        props.name,
                        props.db,
                        keyName.value,
                        currentEditRow.value.no - 1,
                        currentEditRow.value.value,
                    )
                    if (success) {
                        browserStore.loadKeyValue(props.name, props.db, keyName.value).then((r) => {})
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

const onAddValue = (value) => {
    dialogStore.openAddFieldsDialog(props.name, props.db, props.keyPath, props.keyCode, types.LIST)
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
                <n-input
                    v-model:value="filterValue"
                    :placeholder="$t('interface.search')"
                    clearable
                    @clear="clearFilter"
                    @update:value="onFilterInput" />
            </div>
            <div class="flex-item-expand"></div>
            <n-button :focusable="false" plain @click="onAddValue">
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
