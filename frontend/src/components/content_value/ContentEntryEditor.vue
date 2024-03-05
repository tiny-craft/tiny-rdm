<script setup>
import { computed, defineEmits, defineProps, nextTick, reactive, ref, watchEffect } from 'vue'
import { useThemeVars } from 'naive-ui'
import Save from '@/components/icons/Save.vue'
import { decodeTypes, formatTypes } from '@/consts/value_view_type.js'
import { decodeRedisKey } from '@/utils/key_convert.js'
import useBrowserStore from 'stores/browser.js'
import FormatSelector from '@/components/content_value/FormatSelector.vue'
import IconButton from '@/components/common/IconButton.vue'
import FullScreen from '@/components/icons/FullScreen.vue'
import WindowClose from '@/components/icons/WindowClose.vue'
import Pin from '@/components/icons/Pin.vue'
import OffScreen from '@/components/icons/OffScreen.vue'
import ContentEditor from '@/components/content_value/ContentEditor.vue'
import { toString } from 'lodash'

const props = defineProps({
    show: {
        type: Boolean,
    },
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
    fullscreen: {
        type: Boolean,
    },
})

const themeVars = useThemeVars()
const browserStore = useBrowserStore()
const emit = defineEmits([
    'update:field',
    'update:value',
    'update:decode',
    'update:format',
    'update:fullscreen',
    'save',
    'close',
])

watchEffect(
    () => {
        if (props.show && props.value != null) {
            onFormatChanged()
        } else {
            viewAs.value = ''
        }
    },
    {
        flush: 'post',
    },
)

const loading = ref(false)
const isPin = ref(false)
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
const editingContent = ref('')
const enableSave = computed(() => {
    return toString(props.field) !== viewAs.field || editingContent.value !== viewAs.value
})

const viewLanguage = computed(() => {
    switch (viewAs.format) {
        case formatTypes.JSON:
        case formatTypes.UNICODE_JSON:
            return 'json'
        case formatTypes.YAML:
            return 'yaml'
        case formatTypes.XML:
            return 'xml'
        default:
            return 'plaintext'
    }
})

/**
 *
 * @param {decodeTypes|null} decode
 * @param {formatTypes|null} format
 * @return {Promise<void>}
 */
const onFormatChanged = async (decode = null, format = null) => {
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
        editingContent.value = viewAs.value = value
        viewAs.decode = decode || retDecode
        viewAs.format = format || retFormat
    } finally {
        loading.value = false
    }
}

const onInput = (content) => {
    editingContent.value = content
}

const onToggleFullscreen = () => {
    emit('update:fullscreen', !!!props.fullscreen)
}

const onClose = () => {
    isPin.value = false
    emit('close')
}

const onSave = () => {
    emit('save', viewAs.field, editingContent.value, viewAs.decode, viewAs.format)
    if (!isPin.value) {
        nextTick().then(onClose)
    }
}
</script>

<template>
    <div v-show="show" class="entry-editor flex-box-v">
        <n-card :title="$t('interface.edit_row')" autofocus class="flex-item-expand" size="small">
            <div class="editor-content flex-box-v flex-item-expand">
                <!-- field -->
                <div class="editor-content-item flex-box-v">
                    <div class="editor-content-item-label">{{ props.fieldLabel }}</div>
                    <n-input
                        v-model:value="viewAs.field"
                        :placeholder="props.field + ''"
                        :readonly="props.fieldReadonly"
                        class="editor-content-item-input"
                        type="text" />
                </div>

                <!-- value -->
                <div class="editor-content-item flex-box-v flex-item-expand">
                    <div class="editor-content-item-label">{{ props.valueLabel }}</div>
                    <content-editor
                        :border="true"
                        :content="displayValue"
                        :key-path="viewAs.field"
                        :language="viewLanguage"
                        class="flex-item-expand"
                        @input="onInput"
                        @reset="onInput"
                        @save="onSave" />
                    <format-selector
                        :decode="viewAs.decode"
                        :format="viewAs.format"
                        style="margin-top: 5px"
                        @format-changed="(d, f) => onFormatChanged(d, f)" />
                </div>
            </div>
            <template #header-extra>
                <n-space :size="5">
                    <icon-button
                        :button-class="{ 'pinable-btn': true, 'unpin-btn': !isPin, 'pin-btn': isPin }"
                        :icon="Pin"
                        :size="19"
                        :t-tooltip="isPin ? 'interface.unpin_edit' : 'interface.pin_edit'"
                        stroke-width="4"
                        @click="isPin = !isPin" />
                    <icon-button
                        :button-class="['pinable-btn', 'unpin-btn']"
                        :icon="props.fullscreen ? OffScreen : FullScreen"
                        :size="18"
                        stroke-width="5"
                        t-tooltip="interface.fullscreen"
                        @click="onToggleFullscreen" />
                    <icon-button
                        :button-class="['pinable-btn', 'unpin-btn']"
                        :icon="WindowClose"
                        :size="18"
                        stroke-width="5"
                        t-tooltip="menu.close"
                        @click="onClose" />
                </n-space>
            </template>
            <template #action>
                <n-space :wrap="false" :wrap-item="false" justify="end">
                    <n-button :disabled="!enableSave" :secondary="enableSave" type="primary" @click="onSave">
                        <template #icon>
                            <n-icon :component="Save" />
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
    box-sizing: border-box;
    position: absolute;
    left: 0;
    right: 0;
    top: 0;
    bottom: 0;
    z-index: 100;

    .editor-content {
        &-item {
            &:not(:last-child) {
                margin-bottom: 16px;
            }

            &-label {
                height: 18px;
                color: v-bind('themeVars.textColor3');
                font-size: 13px;
                padding: 5px 0;
            }

            &-input {
            }
        }
    }
}

:deep(.n-card__content) {
    display: flex;
    flex-direction: column;
    flex-grow: 1;
}

:deep(.n-card__action) {
    padding: 5px 10px;
    background-color: unset;
}

:deep(.pinable-btn) {
    padding: 3px;
    border-style: solid;
    border-width: 1px;
    border-radius: 3px;
}

:deep(.unpin-btn) {
    border-color: #0000;
}

:deep(.pin-btn) {
    border-color: v-bind('themeVars.iconColorDisabled');
    background-color: v-bind('themeVars.iconColorDisabled');
}

//:deep(.n-card--bordered) {
//    border-radius: 0;
//}
</style>
