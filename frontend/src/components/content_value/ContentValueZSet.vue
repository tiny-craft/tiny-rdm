<script setup>
import { computed, h, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import ContentToolbar from './ContentToolbar.vue'
import AddLink from '@/components/icons/AddLink.vue'
import { NButton, NCode, NIcon, NInput, NInputNumber, useThemeVars } from 'naive-ui'
import { types, types as redisTypes } from '@/consts/support_redis_type.js'
import EditableTableColumn from '@/components/common/EditableTableColumn.vue'
import { isEmpty, size } from 'lodash'
import useDialogStore from 'stores/dialog.js'
import bytes from 'bytes'
import { decodeTypes, formatTypes } from '@/consts/value_view_type.js'
import useBrowserStore from 'stores/browser.js'
import LoadList from '@/components/icons/LoadList.vue'
import LoadAll from '@/components/icons/LoadAll.vue'
import IconButton from '@/components/common/IconButton.vue'

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
    end: Boolean,
    loading: Boolean,
})

const emit = defineEmits(['loadmore', 'loadall', 'reload', 'rename', 'delete'])

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

const browserStore = useBrowserStore()
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
    ellipsis: {
        tooltip: true,
    },
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

const cancelEdit = () => {
    currentEditRow.value.no = 0
}

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
                    const { success, msg } = await browserStore.removeZSetItem(
                        props.name,
                        props.db,
                        keyName.value,
                        row.value,
                    )
                    if (success) {
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
                    const { success, msg } = await browserStore.updateZSetItem(
                        props.name,
                        props.db,
                        keyName.value,
                        row.value,
                        newValue,
                        currentEditRow.value.score,
                    )
                    if (success) {
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
            onCancel: cancelEdit,
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
    if (!isEmpty(props.value)) {
        let index = 0
        for (const elem of props.value) {
            data.push({
                no: ++index,
                value: elem.value,
                score: elem.score,
            })
        }
    }
    return data
})

const entries = computed(() => {
    const len = size(tableData.value)
    return `${len} / ${Math.max(len, props.length)}`
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

defineExpose({
    reset: () => {
        clearFilter()
        cancelEdit()
    },
})
</script>

<template>
    <div class="content-wrapper flex-box-v">
        <content-toolbar
            :db="props.db"
            :key-code="props.keyCode"
            :key-path="props.keyPath"
            :key-type="keyType"
            :loading="props.loading"
            :server="props.name"
            :ttl="ttl"
            class="value-item-part"
            @delete="emit('delete')"
            @reload="emit('reload')"
            @rename="emit('rename')" />
        <div class="tb2 value-item-part flex-box-h">
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
            <div class="flex-item-expand"></div>
            <n-button-group>
                <icon-button
                    :disabled="props.end || props.loading"
                    :icon="LoadList"
                    border
                    size="18"
                    t-tooltip="interface.load_more_entries"
                    @click="emit('loadmore')" />
                <icon-button
                    :disabled="props.end || props.loading"
                    :icon="LoadAll"
                    border
                    size="18"
                    t-tooltip="interface.load_all_entries"
                    @click="emit('loadall')" />
            </n-button-group>
            <n-button :focusable="false" plain @click="onAddRow">
                <template #icon>
                    <n-icon :component="AddLink" size="18" />
                </template>
                {{ $t('interface.add_row') }}
            </n-button>
        </div>
        <div class="value-wrapper value-item-part flex-box-v flex-item-expand">
            <n-data-table
                :key="(row) => row.no"
                :bordered="false"
                :bottom-bordered="false"
                :columns="columns"
                :data="tableData"
                :loading="props.loading"
                :single-column="true"
                :single-line="false"
                class="flex-item-expand"
                flex-height
                size="small"
                striped
                virtual-scroll
                @update:filters="onUpdateFilter" />
        </div>
        <div class="value-footer flex-box-h">
            <n-text v-if="!isNaN(props.length)">{{ $t('interface.entries') }}: {{ entries }}</n-text>
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
