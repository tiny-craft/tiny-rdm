<script setup>
import { computed, h, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import ContentToolbar from './ContentToolbar.vue'
import AddLink from '@/components/icons/AddLink.vue'
import { NButton, NCode, NIcon, NInput, NInputNumber } from 'naive-ui'
import { types, types as redisTypes } from '@/consts/support_redis_type.js'
import EditableTableColumn from '@/components/common/EditableTableColumn.vue'
import { isEmpty } from 'lodash'
import useDialogStore from 'stores/dialog.js'
import useConnectionStore from 'stores/connections.js'

const i18n = useI18n()
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
        label: i18n.t('common.value'),
    },
    {
        value: 2,
        label: i18n.t('interface.score'),
    },
]
const filterType = ref(1)

const connectionStore = useConnectionStore()
const dialogStore = useDialogStore()
const keyType = redisTypes.ZSET
const currentEditRow = ref({
    no: 0,
    score: 0,
    value: null,
})
const scoreColumn = reactive({
    key: 'score',
    title: i18n.t('interface.score'),
    align: 'center',
    titleAlign: 'center',
    resizable: true,
    filterOptionValue: null,
    filter(value, row) {
        const score = parseFloat(row.score)
        if (isNaN(score)) {
            return true
        }

        const regex = /^(>=|<=|>|<|=|!=)?(\d+(\.\d*)?)?$/
        const matches = value.match(regex)
        if (matches) {
            const operator = matches[1] || ''
            const filterScore = parseFloat(matches[2] || '')
            if (!isNaN(filterScore)) {
                switch (operator) {
                    case '>=':
                        return score >= filterScore
                    case '<=':
                        return score <= filterScore
                    case '>':
                        return score > filterScore
                    case '<':
                        return score < filterScore
                    case '=':
                        return score === filterScore
                    case '!=':
                        return score !== filterScore
                }
            }
        } else {
            return !!~row.value.indexOf(value.toString())
        }

        return true
    },
    render: (row) => {
        const isEdit = currentEditRow.value.no === row.no
        if (isEdit) {
            return h(NInputNumber, {
                value: currentEditRow.value.score,
                'onUpdate:value': (val) => {
                    currentEditRow.value.score = val
                },
            })
        } else {
            return row.score
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
            bindKey: row.value,
            onEdit: () => {
                currentEditRow.value.no = row.no
                currentEditRow.value.value = row.value
                currentEditRow.value.score = row.score
            },
            onDelete: async () => {
                try {
                    const { success, msg } = await connectionStore.removeZSetItem(
                        props.name,
                        props.db,
                        keyName.value,
                        row.value,
                    )
                    if (success) {
                        connectionStore.loadKeyValue(props.name, props.db, keyName.value).then((r) => {})
                        $message.success(i18n.t('dialogue.delete_key_succ', { key: row.value }))
                    } else {
                        $message.error(msg)
                    }
                } catch (e) {
                    $message.error(e.message)
                }
            },
            onSave: async () => {
                try {
                    const newValue = currentEditRow.value.value
                    if (isEmpty(newValue)) {
                        $message.error(i18n.t('dialogue.spec_field_required', { key: i18n.t('common.value') }))
                        return
                    }
                    const { success, msg } = await connectionStore.updateZSetItem(
                        props.name,
                        props.db,
                        keyName.value,
                        row.value,
                        newValue,
                        currentEditRow.value.score,
                    )
                    if (success) {
                        connectionStore.loadKeyValue(props.name, props.db, keyName.value).then((r) => {})
                        $message.success(i18n.t('dialogue.save_value_succ'))
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
    valueColumn,
    scoreColumn,
    actionColumn,
])

const tableData = computed(() => {
    const data = []
    let index = 0
    for (const elem of props.value) {
        data.push({
            no: ++index,
            value: elem.value,
            score: elem.score,
        })
    }
    return data
})

const onAddRow = () => {
    dialogStore.openAddFieldsDialog(props.name, props.db, props.keyPath, props.keyCode, types.ZSET)
}

const filterValue = ref('')
const onFilterInput = (val) => {
    switch (filterType.value) {
        case filterOption[0].value:
            // filter value
            scoreColumn.filterOptionValue = null
            valueColumn.filterOptionValue = val
            break
        case filterOption[1].value:
            // filter score
            valueColumn.filterOptionValue = null
            scoreColumn.filterOptionValue = val
            break
    }
}

const onChangeFilterType = (type) => {
    onFilterInput(filterValue.value)
}

const clearFilter = () => {
    valueColumn.filterOptionValue = null
    scoreColumn.filterOptionValue = null
}

const onUpdateFilter = (filters, sourceColumn) => {
    switch (filterType.value) {
        case filterOption[0].value:
            // filter value
            valueColumn.filterOptionValue = filters[sourceColumn.key]
            break
        case filterOption[1].value:
            // filter score
            scoreColumn.filterOptionValue = filters[sourceColumn.key]
            break
    }
}
</script>

<template>
    <div class="content-wrapper flex-box-v">
        <content-toolbar
            :db="props.db"
            :key-path="props.keyPath"
            :key-code="props.keyCode"
            :key-type="keyType"
            :server="props.name"
            :ttl="ttl" />
        <div class="tb2 flex-box-h">
            <div class="flex-box-h">
                <n-input-group>
                    <n-select
                        v-model:value="filterType"
                        :consistent-menu-width="false"
                        :options="filterOption"
                        style="width: 120px"
                        @update:value="onChangeFilterType" />
                    <n-tooltip :delay="500" :disabled="filterType !== 2">
                        <template #trigger>
                            <n-input
                                v-model:value="filterValue"
                                :placeholder="$t('interface.search')"
                                clearable
                                @clear="clearFilter"
                                @update:value="onFilterInput" />
                        </template>
                        <div class="text-block">{{ $t('interface.score_filter_tip') }}</div>
                    </n-tooltip>
                </n-input-group>
            </div>
            <div class="tb2-extra-info flex-item-expand">
                <n-tag size="large">{{ $t('interface.total', { size: props.size }) }}</n-tag>
            </div>
            <n-button plain :focusable="false" @click="onAddRow">
                <template #icon>
                    <n-icon :component="AddLink" size="18" />
                </template>
                {{ $t('interface.add_row') }}
            </n-button>
        </div>
        <div class="value-wrapper fill-height flex-box-h">
            <n-data-table
                :key="(row) => row.no"
                :columns="columns"
                :data="tableData"
                :single-column="true"
                :single-line="false"
                :bordered="false"
                :bottom-bordered="false"
                flex-height
                max-height="100%"
                size="small"
                striped
                virtual-scroll
                @update:filters="onUpdateFilter" />
        </div>
    </div>
</template>

<style lang="scss" scoped></style>
