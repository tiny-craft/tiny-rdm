<script setup>
import { computed, defineEmits, defineProps, nextTick, reactive, ref, watch } from 'vue'
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
const pin = ref(false)
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

const btnStyle = computed(() => ({
    padding: '3px',
    border: 'solid 1px #0000',
    borderRadius: '3px',
}))

const pinBtnStyle = computed(() => ({
    padding: '3px',
    border: `solid 1px ${themeVars.value.borderColor}`,
    borderRadius: '3px',
    backgroundColor: themeVars.value.borderColor,
}))

/**
 *
 * @param {decodeTypes} decode
 * @param {formatTypes} format
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

const onToggleFullscreen = () => {
    emit('update:fullscreen', !!!props.fullscreen)
}

const onClose = () => {
    pin.value = false
    emit('close')
}

const onSave = () => {
    emit('save', viewAs.field, viewAs.value, viewAs.decode, viewAs.format)
    if (!pin.value) {
        nextTick().then(onClose)
    }
}
</script>

<template>
    <div class="entry-editor flex-box-v">
        <n-card :title="$t('interface.edit_row')" autofocus size="small" style="height: 100%">
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
                        :resizable="false"
                        @update:value="onUpdateValue" />
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
                        :button-style="pin ? pinBtnStyle : btnStyle"
                        :size="19"
                        :t-tooltip="pin ? 'interface.unpin_edit' : 'interface.pin_edit'"
                        stroke-width="4"
                        @click="pin = !pin">
                        <Pin :inverse="pin" stroke-width="4" />
                    </icon-button>
                    <icon-button
                        :button-style="btnStyle"
                        :icon="props.fullscreen ? OffScreen : FullScreen"
                        :size="18"
                        stroke-width="5"
                        t-tooltip="interface.fullscreen"
                        @click="onToggleFullscreen" />
                    <icon-button
                        :button-style="btnStyle"
                        :icon="WindowClose"
                        :size="18"
                        stroke-width="5"
                        t-tooltip="menu.close"
                        @click="onClose" />
                </n-space>
            </template>
            <template #action>
                <n-space :wrap="false" :wrap-item="false" justify="end">
                    <n-button ghost @click="onSave">
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

//:deep(.n-card--bordered) {
//    border-radius: 0;
//}
</style>
