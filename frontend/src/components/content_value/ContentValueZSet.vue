<script setup>
import { computed, h, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import AddLink from '@/components/icons/AddLink.vue'
import { NButton, NIcon, useThemeVars } from 'naive-ui'
import { types, types as redisTypes } from '@/consts/support_redis_type.js'
import EditableTableColumn from '@/components/common/EditableTableColumn.vue'
import { isEmpty, size, truncate } from 'lodash'
import useDialogStore from 'stores/dialog.js'
import { decodeTypes, formatTypes } from '@/consts/value_view_type.js'
import useBrowserStore from 'stores/browser.js'
import LoadList from '@/components/icons/LoadList.vue'
import LoadAll from '@/components/icons/LoadAll.vue'
import IconButton from '@/components/common/IconButton.vue'
import ContentEntryEditor from '@/components/content_value/ContentEntryEditor.vue'
import FormatSelector from '@/components/content_value/FormatSelector.vue'
import Edit from '@/components/icons/Edit.vue'
import ContentSearchInput from '@/components/content_value/ContentSearchInput.vue'
import { formatBytes } from '@/utils/byte_convert.js'
import { ClipboardSetText } from 'wailsjs/runtime/runtime.js'
import { TextAlignType } from '@/consts/text_align_type.js'
import AlignLeft from '@/components/icons/AlignLeft.vue'
import AlignCenter from '@/components/icons/AlignCenter.vue'
import SwitchButton from '@/components/common/SwitchButton.vue'
import { nativeRedisKey } from '@/utils/key_convert.js'

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
    textAlign: Number,
})

const emit = defineEmits(['loadmore', 'loadall', 'reload', 'match', 'update:textAlign'])

/**
 *
 * @type {ComputedRef<string|number[]>}
 */
const keyName = computed(() => {
    return !isEmpty(props.keyCode) ? props.keyCode : props.keyPath
})

const browserStore = useBrowserStore()
const dialogStore = useDialogStore()
const keyType = redisTypes.ZSET
const currentEditRow = reactive({
    no: 0,
    score: 0,
    value: null,
    format: formatTypes.RAW,
    decode: decodeTypes.NONE,
})

const inEdit = computed(() => {
    return currentEditRow.no > 0
})
const fullEdit = ref(false)

// const scoreFilterOption = ref(null)
const scoreColumn = computed(() => ({
    key: 'score',
    title: () => i18n.t('common.score'),
    align: props.textAlign !== TextAlignType.Left ? 'center' : 'left',
    titleAlign: 'center',
    resizable: true,
    sorter: (row1, row2) => row1.s - row2.s,
    // filterOptionValue: scoreFilterOption.value,
    // filter(value, row) {
    //     const score = parseFloat(row.s)
    //     if (isNaN(score)) {
    //         return true
    //     }
    //
    //     const regex = /^(>=|<=|>|<|=|!=)?(\d+(\.\d*)?)?$/
    //     const matches = value.match(regex)
    //     if (matches) {
    //         const operator = matches[1] || ''
    //         const filterScore = parseFloat(matches[2] || '')
    //         if (!isNaN(filterScore)) {
    //             switch (operator) {
    //                 case '>=':
    //                     return score >= filterScore
    //                 case '<=':
    //                     return score <= filterScore
    //                 case '>':
    //                     return score > filterScore
    //                 case '<':
    //                     return score < filterScore
    //                 case '=':
    //                     return score === filterScore
    //                 case '!=':
    //                     return score !== filterScore
    //             }
    //         }
    //     } else {
    //         return !!~row.v.indexOf(value.toString())
    //     }
    //     return true
    // },
    render: (row) => {
        return row.ss || row.s
    },
}))

const isCode = computed(() => {
    return props.format === formatTypes.JSON || props.format === formatTypes.UNICODE_JSON
})
const valueFilterOption = ref(null)
const valueColumn = computed(() => ({
    key: 'value',
    title: () => i18n.t('common.value'),
    align: isCode.value ? 'left' : props.textAlign !== TextAlignType.Left ? 'center' : 'left',
    titleAlign: 'center',
    resizable: true,
    ellipsis: isCode.value
        ? false
        : {
              tooltip: {
                  style: {
                      maxWidth: '50vw',
                      maxHeight: '50vh',
                  },
                  scrollable: true,
              },
              lineClamp: 1,
          },
    filterOptionValue: valueFilterOption.value,
    className: inEdit.value ? 'clickable' : '',
    filter(filterValue, row) {
        const val = row.dv || nativeRedisKey(row.v)
        return !!~val.indexOf(filterValue.toString())
    },
    // sorter: (row1, row2) => row1.value - row2.value,
    render: (row) => {
        if (isCode.value) {
            const val = row.dv || nativeRedisKey(row.v)
            return h('pre', { class: 'pre-wrap' }, val)
        } else {
            const val = row.dv || nativeRedisKey(row.v, 500)
            return truncate(val, { length: 500 })
        }
    },
}))

const startEdit = async (no, score, value) => {
    currentEditRow.no = no
    currentEditRow.score = score
    currentEditRow.value = value
    currentEditRow.decode = props.decode
    currentEditRow.format = props.format
}

const saveEdit = async (field, value, decode, format) => {
    try {
        const score = parseFloat(field)
        const row = props.value[currentEditRow.no - 1]
        if (row == null) {
            throw new Error('row not exists')
        }

        if (isEmpty(value)) {
            value = currentEditRow.value
        }

        const { success, msg } = await browserStore.updateZSetItem({
            server: props.name,
            db: props.db,
            key: keyName.value,
            value: row.v,
            newValue: value,
            score,
            decode,
            format,
        })
        if (success) {
            $message.success(i18n.t('interface.save_value_succ'))
        } else {
            $message.error(msg)
        }
    } catch (e) {
        $message.error(e.message)
    }
}

const resetEdit = () => {
    currentEditRow.no = 0
    currentEditRow.score = 0
    currentEditRow.value = null
    // if (currentEditRow.format !== props.format || currentEditRow.decode !== props.decode) {
    //     nextTick(() => onFormatChanged(currentEditRow.decode, currentEditRow.format))
    // }
}

const actionColumn = {
    key: 'action',
    title: () => i18n.t('interface.action'),
    width: 120,
    align: 'center',
    titleAlign: 'center',
    fixed: 'right',
    render: (row, index) => {
        return h(EditableTableColumn, {
            editing: false,
            bindKey: row.v,
            onCopy: async () => {
                await ClipboardSetText(row.v)
                $message.success(i18n.t('interface.copy_succ'))
            },
            onEdit: () => startEdit(index + 1, row.s, row.v),
            onDelete: async () => {
                try {
                    const { success, msg } = await browserStore.removeZSetItem({
                        server: props.name,
                        db: props.db,
                        key: keyName.value,
                        value: row.v,
                    })
                    if (success) {
                        $message.success(i18n.t('dialogue.delete.success', { key: row.v }))
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
            valueColumn.value,
            scoreColumn.value,
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
            valueColumn.value,
        ]
    }
})

const entries = computed(() => {
    const len = size(props.value)
    return `${len} / ${Math.max(len, props.length)}`
})

const loadProgress = computed(() => {
    const len = size(props.value)
    return (len * 100) / Math.max(len, props.length)
})

const showMemoryUsage = computed(() => {
    return !isNaN(props.size) && props.size > 0
})

const onAddRow = () => {
    dialogStore.openAddFieldsDialog(props.name, props.db, props.keyPath, props.keyCode, types.ZSET)
}

const onFilterInput = (val) => {
    valueFilterOption.value = val
}

const onMatchInput = (matchVal, filterVal) => {
    valueFilterOption.value = filterVal
    emit('match', matchVal)
}

const onUpdateFilter = (filters, sourceColumn) => {
    valueFilterOption.value = filters[sourceColumn.key]
}

const onFormatChanged = (selDecode, selFormat) => {
    emit('reload', selDecode, selFormat)
}

const searchInputRef = ref(null)
defineExpose({
    reset: () => {
        resetEdit()
        searchInputRef.value?.reset()
    },
})
</script>

<template>
    <div class="content-wrapper flex-box-v">
        <slot name="toolbar" />
        <div class="tb2 value-item-part flex-box-h">
            <div class="flex-box-h" style="max-width: 50%">
                <content-search-input
                    ref="searchInputRef"
                    @filter-changed="onFilterInput"
                    @match-changed="onMatchInput" />
            </div>
            <div class="flex-item-expand"></div>
            <switch-button
                :icons="[AlignCenter, AlignLeft]"
                :stroke-width="3.5"
                :t-tooltips="['interface.text_align_center', 'interface.text_align_left']"
                :value="props.textAlign"
                size="medium"
                unselect-stroke-width="3"
                @update:value="(val) => emit('update:textAlign', val)" />
            <n-divider vertical />
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
        <!-- loaded progress -->
        <n-progress
            :border-radius="0"
            :color="props.end ? '#0000' : themeVars.primaryColor"
            :height="2"
            :percentage="loadProgress"
            :processing="props.loading"
            :show-indicator="false"
            status="success"
            type="line" />
        <div class="value-wrapper value-item-part flex-box-h flex-item-expand">
            <!-- table -->
            <n-data-table
                :bordered="false"
                :bottom-bordered="false"
                :columns="columns"
                :data="props.value"
                :loading="props.loading"
                :row-key="(row) => row.v"
                :single-column="true"
                :single-line="false"
                class="flex-item-expand"
                flex-height
                size="small"
                striped
                virtual-scroll
                @update:filters="onUpdateFilter" />

            <!-- edit pane -->
            <div
                v-show="inEdit"
                :style="{ position: fullEdit ? 'static' : 'relative' }"
                class="entry-editor-container flex-item-expand"
                style="width: 100%">
                <content-entry-editor
                    v-model:decode="currentEditRow.decode"
                    v-model:format="currentEditRow.format"
                    v-model:fullscreen="fullEdit"
                    :field="currentEditRow.score"
                    :field-label="$t('common.score')"
                    :key-path="props.keyPath"
                    :show="inEdit"
                    :value="currentEditRow.value"
                    :value-label="$t('common.value')"
                    class="flex-item-expand"
                    style="width: 100%"
                    @close="resetEdit"
                    @save="saveEdit" />
            </div>
        </div>
        <div class="value-footer flex-box-h">
            <n-text v-if="!isNaN(props.length)">{{ $t('interface.entries') }}: {{ entries }}</n-text>
            <n-divider v-if="showMemoryUsage" vertical />
            <n-text v-if="showMemoryUsage">{{ $t('interface.memory_usage') }}: {{ formatBytes(props.size) }}</n-text>
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
