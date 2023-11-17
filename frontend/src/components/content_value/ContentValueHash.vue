<script setup>
import { computed, h, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import ContentToolbar from './ContentToolbar.vue'
import AddLink from '@/components/icons/AddLink.vue'
import { NButton, NCode, NIcon, NInput, useThemeVars } from 'naive-ui'
import { types, types as redisTypes } from '@/consts/support_redis_type.js'
import EditableTableColumn from '@/components/common/EditableTableColumn.vue'
import useDialogStore from 'stores/dialog.js'
import { isEmpty, size } from 'lodash'
import bytes from 'bytes'
import { decodeTypes, formatTypes } from '@/consts/value_view_type.js'
import useBrowserStore from 'stores/browser.js'
import LoadList from '@/components/icons/LoadList.vue'
import LoadAll from '@/components/icons/LoadAll.vue'
import IconButton from '@/components/common/IconButton.vue'
import ContentEntryEditor from '@/components/content_value/ContentEntryEditor.vue'
import Edit from '@/components/icons/Edit.vue'
import FormatSelector from '@/components/content_value/FormatSelector.vue'
import { decodeRedisKey } from '@/utils/key_convert.js'

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
    value: {
        type: Array,
        default: () => [],
    },
    size: Number,
    length: Number,
    format: {
        type: String,
        default: formatTypes.RAW,
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
        label: i18n.t('common.field'),
    },
    {
        value: 2,
        label: i18n.t('common.value'),
    },
]
const filterType = ref(1)

const browserStore = useBrowserStore()
const dialogStore = useDialogStore()
const keyType = redisTypes.HASH
const currentEditRow = reactive({
    no: 0,
    key: '',
    value: null,
    format: formatTypes.RAW,
    decode: decodeTypes.NONE,
})

const inEdit = computed(() => {
    return currentEditRow.no > 0
})
const fullEdit = ref(false)
const inFullEdit = computed(() => {
    return inEdit.value && fullEdit.value
})

const tableRef = ref(null)
const fieldFilterOption = ref(null)
const fieldColumn = computed(() => ({
    key: 'key',
    title: i18n.t('common.field'),
    align: 'center',
    titleAlign: 'center',
    resizable: true,
    ellipsis: {
        tooltip: true,
    },
    filterOptionValue: fieldFilterOption.value,
    className: inEdit.value ? 'clickable' : '',
    filter: (value, row) => {
        return !!~row.k.indexOf(value.toString())
    },
    render: (row) => {
        return decodeRedisKey(row.k)
    },
}))

const displayCode = computed(() => {
    return props.format === formatTypes.JSON
})
const valueFilterOption = ref(null)
const valueColumn = computed(() => ({
    key: 'value',
    title: i18n.t('common.value'),
    align: displayCode.value ? 'left' : 'center',
    titleAlign: 'center',
    resizable: true,
    ellipsis: displayCode.value
        ? false
        : {
              tooltip: true,
          },
    filterOptionValue: valueFilterOption.value,
    className: inEdit.value ? 'clickable' : '',
    filter: (value, row) => {
        if (row.dv) {
            return !!~row.dv.indexOf(value.toString())
        }
        return !!~row.v.indexOf(value.toString())
    },
    render: (row) => {
        if (displayCode.value) {
            return h(NCode, { language: 'json', wordWrap: true, code: row.dv || row.v })
        }
        return row.dv || row.v
    },
}))

const startEdit = async (no, key, value) => {
    currentEditRow.no = no
    currentEditRow.key = key
    currentEditRow.value = value
}

const saveEdit = async (field, value, decode, format) => {
    try {
        const row = props.value[currentEditRow.no - 1]
        if (row == null) {
            throw new Error('row not exists')
        }

        const { success, msg } = await browserStore.setHash({
            server: props.name,
            db: props.db,
            key: keyName.value,
            field: row.k,
            newField: field,
            value,
            decode,
            format,
            retDecode: props.decode,
            retFormat: props.format,
            index: [currentEditRow.no - 1],
        })
        if (success) {
            $message.success(i18n.t('dialogue.save_value_succ'))
        } else {
            $message.error(msg)
        }
    } catch (e) {
        $message.error(e.message)
    }
}

const resetEdit = () => {
    currentEditRow.no = 0
    currentEditRow.key = ''
    currentEditRow.value = null
    currentEditRow.format = formatTypes.RAW
    currentEditRow.decode = decodeTypes.NONE
}

const actionColumn = {
    key: 'action',
    title: i18n.t('interface.action'),
    width: 100,
    align: 'center',
    titleAlign: 'center',
    fixed: 'right',
    render: (row, index) => {
        return h(EditableTableColumn, {
            editing: false,
            bindKey: row.k,
            onEdit: () => startEdit(index + 1, row.k, row.v),
            onDelete: async () => {
                try {
                    const { removed, success, msg } = await browserStore.removeHashField(
                        props.name,
                        props.db,
                        keyName.value,
                        row.k,
                    )
                    if (success) {
                        props.value.splice(index, 1)
                        $message.success(i18n.t('dialogue.delete_key_succ', { key: row.k }))
                    } else {
                        $message.error(msg)
                    }
                } catch (e) {
                    $message.error(e.message)
                }
            },
        })
    },
}

const columns = computed(() => {
    if (!inEdit.value) {
        return [
            {
                key: 'no',
                title: '#',
                width: 80,
                align: 'center',
                titleAlign: 'center',
                render: (row, index) => {
                    return index + 1
                },
            },
            fieldColumn.value,
            valueColumn.value,
            actionColumn,
        ]
    } else {
        return [
            {
                key: 'no',
                title: '#',
                width: 80,
                align: 'center',
                titleAlign: 'center',
                render: (row, index) => {
                    if (index + 1 === currentEditRow.no) {
                        // editing row, show edit state
                        return h(NIcon, { size: 16, color: 'red' }, () => h(Edit, { strokeWidth: 5 }))
                    } else {
                        return index + 1
                    }
                },
            },
            fieldColumn.value,
        ]
    }
})

const rowProps = (row, index) => {
    return {
        onClick: () => {
            // in edit mode, switch edit row by click
            if (inEdit.value) {
                startEdit(index + 1, row.k, row.v)
            }
        },
    }
}

const entries = computed(() => {
    const len = size(props.value)
    return `${len} / ${Math.max(len, props.length)}`
})

const onAddRow = () => {
    dialogStore.openAddFieldsDialog(props.name, props.db, props.keyPath, props.keyCode, types.HASH)
}

const filterValue = ref('')
const onFilterInput = (val) => {
    switch (filterType.value) {
        case filterOption[0].value:
            // filter field
            valueFilterOption.value = null
            fieldFilterOption.value = val
            break
        case filterOption[1].value:
            // filter value
            fieldFilterOption.value = null
            valueFilterOption.value = val
            break
    }
}

const onChangeFilterType = (type) => {
    onFilterInput(filterValue.value)
}

const clearFilter = () => {
    fieldFilterOption.value = null
    valueFilterOption.value = null
}

const onUpdateFilter = (filters, sourceColumn) => {
    switch (filterType.value) {
        case filterOption[0].value:
            fieldFilterOption.value = filters[sourceColumn.key]
            break
        case filterOption[1].value:
            valueFilterOption.value = filters[sourceColumn.key]
            break
    }
}

const onFormatChanged = (selDecode, selFormat) => {
    emit('reload', selDecode, selFormat)
}

defineExpose({
    reset: () => {
        clearFilter()
        resetEdit()
    },
})
</script>

<template>
    <div class="content-wrapper flex-box-v">
        <content-toolbar
            v-show="!inFullEdit"
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
        <div v-show="!inFullEdit" class="tb2 value-item-part flex-box-h">
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
        <div id="content-table" class="value-wrapper value-item-part flex-box-h flex-item-expand">
            <!-- table -->
            <n-data-table
                v-show="!inFullEdit"
                ref="tableRef"
                :bordered="false"
                :bottom-bordered="false"
                :columns="columns"
                :data="props.value"
                :loading="props.loading"
                :row-key="(row) => row.k"
                :row-props="rowProps"
                :single-column="true"
                :single-line="false"
                class="flex-item-expand"
                flex-height
                size="small"
                striped
                virtual-scroll
                @update:filters="onUpdateFilter" />

            <!-- edit pane -->
            <content-entry-editor
                v-show="inEdit"
                v-model:fullscreen="fullEdit"
                :decode="currentEditRow.decode"
                :field="currentEditRow.key"
                :field-label="$t('common.field')"
                :format="currentEditRow.format"
                :value="currentEditRow.value"
                :value-label="$t('common.value')"
                class="flex-item-expand"
                style="width: 100%"
                @close="resetEdit"
                @save="saveEdit" />
        </div>
        <div class="value-footer flex-box-h">
            <n-text v-if="!isNaN(props.length)">{{ $t('interface.entries') }}: {{ entries }}</n-text>
            <n-divider v-if="!isNaN(props.length)" vertical />
            <n-text v-if="!isNaN(props.size)">{{ $t('interface.memory_usage') }}: {{ bytes(props.size) }}</n-text>
            <div class="flex-item-expand"></div>
            <format-selector
                v-show="!inEdit"
                :decode="props.decode"
                :disabled="inEdit"
                :format="props.format"
                @format-changed="onFormatChanged" />
        </div>
    </div>
</template>

<style lang="scss" scoped>
.value-footer {
    border-top: v-bind('themeVars.borderColor') 1px solid;
    background-color: v-bind('themeVars.tableHeaderColor');
}
</style>
