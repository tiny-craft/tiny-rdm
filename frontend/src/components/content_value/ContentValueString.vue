<script setup>
import { computed, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import ContentToolbar from './ContentToolbar.vue'
import Copy from '@/components/icons/Copy.vue'
import Save from '@/components/icons/Save.vue'
import { useThemeVars } from 'naive-ui'
import { decodeTypes, formatTypes } from '@/consts/value_view_type.js'
import Close from '@/components/icons/Close.vue'
import { types as redisTypes } from '@/consts/support_redis_type.js'
import { ClipboardSetText } from 'wailsjs/runtime/runtime.js'
import { isEmpty, toLower } from 'lodash'
import EditFile from '@/components/icons/EditFile.vue'
import bytes from 'bytes'
import useBrowserStore from 'stores/browser.js'
import { decodeRedisKey } from '@/utils/key_convert.js'
import FormatSelector from '@/components/content_value/FormatSelector.vue'

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
    value: [String, Array],
    size: Number,
    length: Number,
    loading: Boolean,
})

const emit = defineEmits(['reload', 'rename', 'delete'])

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
            return 'json'
        default:
            return 'plaintext'
    }
})

const viewAs = reactive({
    value: '',
    format: formatTypes.RAW,
    decode: decodeTypes.NONE,
})

const displayValue = computed(() => {
    if (props.loading) {
        return ''
    }
    return viewAs.value || decodeRedisKey(props.value)
})

watch(
    () => props.value,
    (val, oldVal) => {
        if (val !== undefined && oldVal !== undefined) {
            onFormatChanged(viewAs.decode, viewAs.format)
        }
    },
)

const onFormatChanged = async (decode = '', format = '') => {
    const {
        value,
        decode: retDecode,
        format: retFormat,
    } = await browserStore.convertValue({
        value: props.value,
        decode,
        format,
    })
    viewAs.value = value
    viewAs.decode = decode || retDecode
    viewAs.format = format || retFormat
}

/**
 * Copy value
 */
const onCopyValue = () => {
    ClipboardSetText(displayValue.value)
        .then((succ) => {
            if (succ) {
                $message.success(i18n.t('dialogue.copy_succ'))
            }
        })
        .catch((e) => {
            $message.error(e.message)
        })
}

const editValue = ref('')
const inEdit = ref(false)
const onEditValue = () => {
    editValue.value = displayValue.value
    inEdit.value = true
}

const onCancelEdit = () => {
    inEdit.value = false
}

/**
 * Save value
 */
const browserStore = useBrowserStore()
const saving = ref(false)
const onSaveValue = async () => {
    saving.value = true
    try {
        const { success, msg } = await browserStore.setKey({
            server: props.name,
            db: props.db,
            key: keyName.value,
            keyType: toLower(keyType),
            value: editValue.value,
            ttl: -1,
            format: viewAs.format,
            decode: viewAs.decode,
        })
        if (success) {
            await browserStore.loadKeyDetail({ server: props.name, db: props.db, key: keyName.value })
            $message.success(i18n.t('dialogue.save_value_succ'))
        } else {
            $message.error(msg)
        }
    } catch (e) {
        $message.error(e.message)
    } finally {
        inEdit.value = false
        saving.value = false
    }
}

defineExpose({
    reset: () => {
        viewAs.value = ''
        inEdit.value = false
    },
    beforeShow: () => onFormatChanged(),
})
</script>

<template>
    <div class="content-wrapper flex-box-v">
        <content-toolbar
            :db="props.db"
            :key-code="keyCode"
            :key-path="keyPath"
            :key-type="keyType"
            :loading="loading"
            :server="props.name"
            :ttl="ttl"
            class="value-item-part"
            @delete="emit('delete')"
            @reload="emit('reload')"
            @rename="emit('rename')" />
        <div class="tb2 value-item-part flex-box-h">
            <div class="flex-item-expand"></div>
            <n-button-group v-if="!inEdit">
                <n-button :focusable="false" @click="onCopyValue">
                    <template #icon>
                        <n-icon :component="Copy" size="18" />
                    </template>
                    {{ $t('interface.copy_value') }}
                </n-button>
                <n-button :focusable="false" plain @click="onEditValue">
                    <template #icon>
                        <n-icon :component="EditFile" size="18" />
                    </template>
                    {{ $t('interface.edit_value') }}
                </n-button>
            </n-button-group>
            <n-button-group v-else>
                <n-button :focusable="false" :loading="saving" plain @click="onSaveValue">
                    <template #icon>
                        <n-icon :component="Save" size="18" />
                    </template>
                    {{ $t('interface.save_update') }}
                </n-button>
                <n-button :focusable="false" :loading="saving" plain @click="onCancelEdit">
                    <template #icon>
                        <n-icon :component="Close" size="18" />
                    </template>
                    {{ $t('common.cancel') }}
                </n-button>
            </n-button-group>
        </div>
        <div class="value-wrapper value-item-part flex-item-expand flex-box-v">
            <n-scrollbar v-if="!inEdit" class="flex-item-expand">
                <n-code :code="displayValue" :language="viewLanguage" style="cursor: text" word-wrap />
            </n-scrollbar>
            <n-input
                v-else
                v-model:value="editValue"
                :disabled="saving"
                :resizable="false"
                class="flex-item-expand"
                type="textarea" />
        </div>
        <div class="value-footer flex-box-h">
            <n-text v-if="!isNaN(props.length)">{{ $t('interface.length') }}: {{ props.length }}</n-text>
            <n-divider v-if="!isNaN(props.length)" vertical />
            <n-text v-if="!isNaN(props.size)">{{ $t('interface.memory_usage') }}: {{ bytes(props.size) }}</n-text>
            <div class="flex-item-expand"></div>
            <format-selector
                :decode="viewAs.decode"
                :disabled="inEdit"
                :format="viewAs.format"
                @format-changed="onFormatChanged" />
        </div>
    </div>
</template>

<style lang="scss" scoped>
.value-wrapper {
    overflow: hidden;
    border-top: v-bind('themeVars.borderColor') 1px solid;
    padding: 5px;
}

.value-footer {
    border-top: v-bind('themeVars.borderColor') 1px solid;
    background-color: v-bind('themeVars.tableHeaderColor');
}
</style>
