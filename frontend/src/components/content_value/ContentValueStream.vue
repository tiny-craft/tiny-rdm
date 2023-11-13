<script setup>
import { computed, h, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import ContentToolbar from './ContentToolbar.vue'
import AddLink from '@/components/icons/AddLink.vue'
import { NButton, NCode, NIcon, NInput, useThemeVars } from 'naive-ui'
import { types, types as redisTypes } from '@/consts/support_redis_type.js'
import EditableTableColumn from '@/components/common/EditableTableColumn.vue'
import useDialogStore from 'stores/dialog.js'
import { includes, isEmpty, keys, size, some, values } from 'lodash'
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
    value: {
        type: Array,
        default: () => [],
    },
    size: Number,
    length: Number,
    viewAs: {
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
const keyType = redisTypes.STREAM
const idColumn = reactive({
    key: 'id',
    title: 'ID',
    align: 'center',
    titleAlign: 'center',
    resizable: true,
})
const valueColumn = reactive({
    key: 'value',
    title: i18n.t('common.value'),
    align: 'center',
    titleAlign: 'center',
    resizable: true,
    filterOptionValue: null,
    filter(value, row) {
        const v = value.toString()
        if (filterType.value === 1) {
            // filter key
            return some(keys(row.value), (key) => includes(key, v))
        } else {
            // filter value
            return some(values(row.value), (val) => includes(val, v))
        }
    },
    // sorter: (row1, row2) => row1.value - row2.value,
    // ellipsis: {
    //     tooltip: true
    // },
    render: (row) => {
        return h(NCode, { language: 'json', wordWrap: true }, { default: () => JSON.stringify(row.value) })
    },
})
const actionColumn = {
    key: 'action',
    title: i18n.t('interface.action'),
    width: 60,
    align: 'center',
    titleAlign: 'center',
    fixed: 'right',
    render: (row) => {
        return h(EditableTableColumn, {
            bindKey: row.id,
            readonly: true,
            onDelete: async () => {
                try {
                    const { success, msg } = await browserStore.removeStreamValues(
                        props.name,
                        props.db,
                        keyName.value,
                        row.id,
                    )
                    if (success) {
                        $message.success(i18n.t('dialogue.delete_key_succ', { key: row.id }))
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
const columns = reactive([idColumn, valueColumn, actionColumn])

const tableData = computed(() => {
    const data = []
    if (!isEmpty(props.value)) {
        for (const elem of props.value) {
            data.push({
                id: elem.id,
                value: elem.value,
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
    dialogStore.openAddFieldsDialog(props.name, props.db, props.keyPath, props.keyCode, types.STREAM)
}

const filterValue = ref('')
const onFilterInput = (val) => {
    valueColumn.filterOptionValue = val
}

const onChangeFilterType = (type) => {
    onFilterInput(filterValue.value)
}

const clearFilter = () => {
    idColumn.filterOptionValue = null
    valueColumn.filterOptionValue = null
}

const onUpdateFilter = (filters, sourceColumn) => {
    switch (filterType.value) {
        case filterOption[0].value:
            idColumn.filterOptionValue = filters[sourceColumn.key]
            break
        case filterOption[1].value:
            valueColumn.filterOptionValue = filters[sourceColumn.key]
            break
    }
}

defineExpose({
    reset: () => {
        clearFilter()
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
        <div class="value-wrapper value-item-part flex-box-v flex-item-expand">
            <n-data-table
                :key="(row) => row.id"
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
