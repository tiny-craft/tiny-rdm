<script setup>
import { computed, defineEmits, defineProps, reactive, ref, watch } from 'vue'
import { useThemeVars } from 'naive-ui'
import Save from '@/components/icons/Save.vue'
import { decodeTypes, formatTypes } from '@/consts/value_view_type.js'
import { decodeRedisKey } from '@/utils/key_convert.js'
import useBrowserStore from 'stores/browser.js'
import FormatSelector from '@/components/content_value/FormatSelector.vue'

const props = defineProps({
    field: {
        type: [String, Number],
    },
    value: {
        type: String,
    },
    fieldLabel: {
        type: String,
    },
    valueLabel: {
        type: String,
    },
    decode: {
        type: String,
    },
    format: {
        type: String,
    },
    fieldReadonly: {
        type: Boolean,
    },
})

const themeVars = useThemeVars()
const browserStore = useBrowserStore()
const emit = defineEmits(['update:field', 'update:value', 'update:decode', 'update:format', 'save', 'cancel'])
const model = reactive({
    field: '',
    value: '',
})

watch(
    () => props.value,
    (val) => {
        if (val != null) {
            onFormatChanged()
        } else {
            viewAs.value = ''
        }
    },
)

const loading = ref(false)
const viewAs = reactive({
    field: '',
    value: '',
    format: formatTypes.RAW,
    decode: decodeTypes.NONE,
})
const displayValue = computed(() => {
    if (loading.value) {
        return ''
    }
    if (viewAs.value == null) {
        return decodeRedisKey(props.value)
    }
    return viewAs.value
})

/**
 *
 * @param {string} decode
 * @param {string} format
 * @return {Promise<void>}
 */
const onFormatChanged = async (decode = '', format = '') => {
    try {
        loading.value = true
        const {
            value,
            decode: retDecode,
            format: retFormat,
        } = await browserStore.convertValue({
            value: props.value,
            decode,
            format,
        })
        viewAs.field = props.field + ''
        viewAs.value = value
        viewAs.decode = decode || retDecode
        viewAs.format = format || retFormat
    } finally {
        loading.value = false
    }
}

const onUpdateValue = (value) => {
    // emit('update:value', value)
    viewAs.value = value
}

const onSave = () => {
    emit('save', viewAs.field, viewAs.value, viewAs.decode, viewAs.format)
}
</script>

<template>
    <div class="entry-editor flex-box-v">
        <n-card
            :title="$t('interface.edit_row')"
            autofocus
            closable
            size="small"
            style="height: 100%"
            @close="emit('cancel')">
            <div class="editor-content flex-box-v" style="height: 100%">
                <!-- field -->
                <div class="editor-content-item flex-box-v">
                    <div class="editor-content-item-label">{{ props.fieldLabel }}</div>
                    <n-input
                        v-model:value="viewAs.field"
                        :readonly="props.fieldReadonly"
                        class="editor-content-item-input"
                        type="text" />
                </div>

                <!-- value -->
                <div class="editor-content-item flex-box-v flex-item-expand">
                    <div class="editor-content-item-label">{{ props.valueLabel }}</div>
                    <n-input
                        :value="displayValue"
                        autofocus
                        class="flex-item-expand"
                        type="textarea"
                        @update:value="onUpdateValue" />
                    <format-selector
                        :decode="viewAs.decode"
                        :format="viewAs.format"
                        @format-changed="(d, f) => onFormatChanged(d, f)" />
                </div>
            </div>
            <template #action>
                <n-space :wrap="false" :wrap-item="false" justify="end">
                    <n-button ghost type="primary" @click="onSave">
                        <template #icon>
                            <n-icon :component="Save"></n-icon>
                        </template>
                        {{ $t('common.update') }}
                    </n-button>
                </n-space>
            </template>
        </n-card>
    </div>
</template>

<style lang="scss" scoped>
.entry-editor {
    padding-left: 2px;

    .editor-content {
        &-item {
            &:not(:last-child) {
                margin-bottom: 18px;
            }

            &-label {
                line-height: 1.25;
                color: v-bind('themeVars.textColor3');
                font-size: 13px;
                padding: 5px 0;
            }

            &-input {
            }
        }
    }
}

:deep(.n-card__action) {
    padding: 5px 10px;
    background-color: unset;
}

:deep(.n-card--bordered) {
    border-radius: 0;
}
</style>
