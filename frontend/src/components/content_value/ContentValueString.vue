<script setup>
import { computed, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import Copy from '@/components/icons/Copy.vue'
import Save from '@/components/icons/Save.vue'
import { useThemeVars } from 'naive-ui'
import { formatTypes } from '@/consts/value_view_type.js'
import { types as redisTypes } from '@/consts/support_redis_type.js'
import { isEmpty, toLower } from 'lodash'
import useBrowserStore from 'stores/browser.js'
import { decodeRedisKey } from '@/utils/key_convert.js'
import FormatSelector from '@/components/content_value/FormatSelector.vue'
import ContentEditor from '@/components/content_value/ContentEditor.vue'
import { formatBytes } from '@/utils/byte_convert.js'
import { ClipboardSetText } from 'wailsjs/runtime/runtime.js'

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
    value: [String, Array],
    format: {
        type: String,
    },
    decode: {
        type: String,
    },
    size: Number,
    length: Number,
    loading: Boolean,
})

const i18n = useI18n()
const themeVars = useThemeVars()

/**
 *
 * @type {ComputedRef<string|number[]>}
 */
const keyName = computed(() => {
    return !isEmpty(props.keyCode) ? props.keyCode : props.keyPath
})

const keyType = redisTypes.STRING
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

const viewAs = reactive({
    value: '',
    format: '',
    decode: '',
})

const editingContent = ref('')
const resetKey = ref('')

const enableSave = computed(() => {
    return editingContent.value !== viewAs.value && !props.loading
})

const displayValue = computed(() => {
    return viewAs.value || decodeRedisKey(props.value) || ''
})

const showMemoryUsage = computed(() => {
    return !isNaN(props.size) && props.size > 0
})

watch(
    () => props.value,
    (val) => {
        if (!isEmpty(val)) {
            onFormatChanged(viewAs.decode, viewAs.format)
        }
    },
)

const converting = ref(false)
const onFormatChanged = async (decode = '', format = '') => {
    try {
        converting.value = true
        const {
            value,
            decode: retDecode,
            format: retFormat,
        } = await browserStore.convertValue({
            value: props.value,
            decode: decode || props.decode,
            format: format || props.format,
        })
        editingContent.value = viewAs.value = value
        viewAs.decode = decode || retDecode
        viewAs.format = format || retFormat
        browserStore.setSelectedFormat(props.name, props.keyPath, props.db, viewAs.format, viewAs.decode)
        resetKey.value = Date.now().toString()
    } finally {
        converting.value = false
    }
}

/**
 * Copy value
 */
const onCopyValue = async () => {
    await ClipboardSetText(displayValue.value)
    $message.success(i18n.t('interface.copy_succ'))
}

/**
 * Save value
 */
const browserStore = useBrowserStore()
const saving = ref(false)

const onInput = (content) => {
    editingContent.value = content
}

const onSave = async () => {
    saving.value = true
    try {
        const { success, msg } = await browserStore.setKey({
            server: props.name,
            db: props.db,
            key: keyName.value,
            keyType: toLower(keyType),
            value: editingContent.value,
            ttl: -1,
            format: viewAs.format,
            decode: viewAs.decode,
        })
        if (success) {
            viewAs.value = editingContent.value
            $message.success(i18n.t('interface.save_value_succ'))
        } else {
            $message.error(msg)
        }
    } catch (e) {
        $message.error(e.message)
    } finally {
        saving.value = false
    }
}

defineExpose({
    reset: () => {
        viewAs.value = ''
        viewAs.decode = ''
        viewAs.format = ''
        editingContent.value = ''
    },
})
</script>

<template>
    <div class="content-wrapper flex-box-v">
        <slot name="toolbar" />
        <div class="tb2 value-item-part flex-box-h">
            <div class="flex-item-expand"></div>
            <n-button-group>
                <n-button :disabled="saving" :focusable="false" @click="onCopyValue">
                    <template #icon>
                        <n-icon :component="Copy" size="18" />
                    </template>
                    {{ $t('interface.copy_value') }}
                </n-button>
                <n-button
                    :disabled="!enableSave"
                    :loading="saving"
                    :secondary="enableSave"
                    :type="enableSave ? 'primary' : ''"
                    @click="onSave">
                    <template #icon>
                        <n-icon :component="Save" size="18" />
                    </template>
                    {{ $t('common.save') }}
                </n-button>
            </n-button-group>
        </div>
        <div class="value-wrapper value-item-part flex-item-expand flex-box-v">
            <content-editor
                :content="displayValue"
                :language="viewLanguage"
                :loading="props.loading"
                :offset-key="props.keyPath"
                :reset-key="resetKey"
                class="flex-item-expand"
                keep-offset
                style="height: 100%"
                @input="onInput"
                @reset="onInput"
                @save="onSave" />
            <n-spin v-show="props.loading || converting" />
        </div>
        <div class="value-footer flex-box-h">
            <n-text v-if="!isNaN(props.length)">{{ $t('interface.length') }}: {{ props.length }}</n-text>
            <n-divider v-if="showMemoryUsage" vertical />
            <n-text v-if="showMemoryUsage">{{ $t('interface.memory_usage') }}: {{ formatBytes(props.size) }}</n-text>
            <div class="flex-item-expand" />
            <format-selector
                :decode="viewAs.decode"
                :disabled="enableSave"
                :format="viewAs.format"
                @format-changed="onFormatChanged" />
        </div>
    </div>
</template>

<style lang="scss" scoped>
.value-wrapper {
    //overflow: hidden;
    border-top: v-bind('themeVars.borderColor') 1px solid;
    padding: 5px;
}

.value-footer {
    border-top: v-bind('themeVars.borderColor') 1px solid;
    background-color: v-bind('themeVars.tableHeaderColor');
}
</style>
