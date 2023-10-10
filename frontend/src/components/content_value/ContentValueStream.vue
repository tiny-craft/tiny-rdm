<script setup>
import { computed, h, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import ContentToolbar from './ContentToolbar.vue'
import AddLink from '@/components/icons/AddLink.vue'
import { NButton, NCode, NIcon, NInput } from 'naive-ui'
import { types, types as redisTypes } from '@/consts/support_redis_type.js'
import EditableTableColumn from '@/components/common/EditableTableColumn.vue'
import useDialogStore from 'stores/dialog.js'
import useConnectionStore from 'stores/connections.js'
import { includes, isEmpty, keys, some, values } from 'lodash'

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
                    const { success, msg } = await connectionStore.removeStreamValues(
                        props.name,
                        props.db,
                        keyName.value,
                        row.id,
                    )
                    if (success) {
                        connectionStore.loadKeyValue(props.name, props.db, keyName.value).then((r) => {})
                        $message.success(i18n.t('dialogue.delete_key_succ', { key: row.id }))
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
        })
    },
}
const columns = reactive([idColumn, valueColumn, actionColumn])

const tableData = computed(() => {
    const data = []
    for (const elem of props.value) {
        data.push({
            id: elem.id,
            value: elem.value,
        })
    }
    return data
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
                    <n-input
                        v-model:value="filterValue"
                        :placeholder="$t('interface.search')"
                        clearable
                        @clear="clearFilter"
                        @update:value="onFilterInput" />
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
                :key="(row) => row.id"
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
