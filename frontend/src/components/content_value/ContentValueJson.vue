<script setup>
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import Copy from '@/components/icons/Copy.vue'
import Save from '@/components/icons/Save.vue'
import { useThemeVars } from 'naive-ui'
import { types as redisTypes } from '@/consts/support_redis_type.js'
import { isEmpty, toLower } from 'lodash'
import useBrowserStore from 'stores/browser.js'
import { decodeRedisKey } from '@/utils/key_convert.js'
import ContentEditor from '@/components/content_value/ContentEditor.vue'
import { decodeTypes, formatTypes } from '@/consts/value_view_type.js'
import { formatBytes } from '@/utils/byte_convert.js'
import copy from 'copy-text-to-clipboard'

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
    value: String,
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

const keyType = redisTypes.JSON

const editingContent = ref('')

const displayValue = computed(() => {
    return decodeRedisKey(props.value) || ''
})

const enableSave = computed(() => {
    return editingContent.value !== displayValue.value && !props.loading
})

const showMemoryUsage = computed(() => {
    return !isNaN(props.size) && props.size > 0
})

/**
 * Copy value
 */
const onCopyValue = () => {
    copy(displayValue.value)
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
            format: formatTypes.JSON,
            decode: decodeTypes.NONE,
        })
        if (success) {
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
                :loading="props.loading"
                :offset-key="props.keyPath"
                class="flex-item-expand"
                keep-offset
                language="json"
                style="height: 100%"
                @input="onInput"
                @reset="onInput"
                @save="onSave" />
            <n-spin v-show="props.loading" />
        </div>
        <div class="value-footer flex-box-h">
            <n-text v-if="showMemoryUsage">{{ $t('interface.memory_usage') }}: {{ formatBytes(props.size) }}</n-text>
            <div class="flex-item-expand" />
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
