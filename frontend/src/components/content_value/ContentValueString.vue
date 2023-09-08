<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import ContentToolbar from './ContentToolbar.vue'
import Copy from '@/components/icons/Copy.vue'
import Save from '@/components/icons/Save.vue'
import { useThemeVars } from 'naive-ui'
import { types } from '@/consts/value_view_type.js'
import Close from '@/components/icons/Close.vue'
import Edit from '@/components/icons/Edit.vue'
import { IsJson } from '@/utils/check_string_format.js'
import { types as redisTypes } from '@/consts/support_redis_type.js'
import { ClipboardSetText } from 'wailsjs/runtime/runtime.js'
import { toLower } from 'lodash'
import useConnectionStore from 'stores/connections.js'

const i18n = useI18n()
const themeVars = useThemeVars()

const props = defineProps({
    name: String,
    db: Number,
    keyPath: String,
    ttl: {
        type: Number,
        default: -1,
    },
    value: String,
})

const viewOption = [
    {
        value: types.PLAIN_TEXT,
        label: types.PLAIN_TEXT,
    },
    {
        value: types.JSON,
        label: types.JSON,
    },
    {
        value: types.BASE64_TO_TEXT,
        label: types.BASE64_TO_TEXT,
    },
    {
        value: types.BASE64_TO_JSON,
        label: types.BASE64_TO_JSON,
    },
]
const viewAs = ref(types.PLAIN_TEXT)

const getJsonValue = () => {
    try {
        const jsonObj = JSON.parse(props.value)
        return JSON.stringify(jsonObj, null, 2)
    } catch (e) {
        return props.value
    }
}

const getBase64Value = () => {
    try {
        return atob(props.value)
    } catch (e) {
        return props.value
    }
}

const getBase64JsonValue = () => {
    try {
        const text = atob(props.value)
        const jsonObj = JSON.parse(text)
        return JSON.stringify(jsonObj, null, 2)
    } catch (e) {
        return props.value
    }
}

const autoDetectFormat = () => {
    // auto check format when loaded
    if (IsJson(props.value)) {
        viewAs.value = types.JSON
    } else {
        viewAs.value = types.PLAIN_TEXT
    }
}

onMounted(() => {
    autoDetectFormat()
})
watch(
    () => props.value,
    (value) => {
        autoDetectFormat()
    },
)

const keyType = redisTypes.STRING
/**
 * view value
 * @type {ComputedRef<string>}
 */
const viewValue = computed(() => {
    switch (viewAs.value) {
        case types.PLAIN_TEXT:
            return props.value
        case types.JSON:
            return getJsonValue()
        case types.BASE64_TO_TEXT:
            return getBase64Value()
        case types.BASE64_TO_JSON:
            return getBase64JsonValue()
        default:
            return props.value
    }
})

const viewLanguage = computed(() => {
    switch (viewAs.value) {
        case types.PLAIN_TEXT:
        case types.BASE64_TO_TEXT:
            return 'plaintext'
        case types.JSON:
        case types.BASE64_TO_JSON:
            return 'json'
        default:
            return 'plaintext'
    }
})

/**
 * Copy value
 */
const onCopyValue = () => {
    ClipboardSetText(viewValue.value)
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
    editValue.value = viewValue.value
    inEdit.value = true
}

const onCancelEdit = () => {
    inEdit.value = false
}

/**
 * Save value
 */
const connectionStore = useConnectionStore()
const saving = ref(false)
const onSaveValue = async () => {
    saving.value = true
    try {
        const { success, msg } = await connectionStore.setKey(
            props.name,
            props.db,
            props.keyPath,
            toLower(keyType),
            editValue.value,
            -1,
        )
        if (success) {
            await connectionStore.loadKeyValue(props.name, props.db, props.keyPath)
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
</script>

<template>
    <div class="content-wrapper flex-box-v">
        <content-toolbar :db="props.db" :key-path="keyPath" :key-type="keyType" :server="props.name" :ttl="ttl" />
        <div class="tb2 flex-box-h">
            <n-text>{{ $t('interface.view_as') }}</n-text>
            <n-select v-model:value="viewAs" :options="viewOption" style="width: 200px" />
            <div class="flex-item-expand"></div>
            <n-button-group v-if="!inEdit">
                <n-button :focusable="false" @click="onCopyValue">
                    <template #icon>
                        <n-icon :component="Copy" size="18" />
                    </template>
                    {{ $t('interface.copy_value') }}
                </n-button>
                <n-button plain :focusable="false" @click="onEditValue">
                    <template #icon>
                        <n-icon :component="Edit" size="18" />
                    </template>
                    {{ $t('interface.edit_value') }}
                </n-button>
            </n-button-group>
            <n-button-group v-else>
                <n-button :loading="saving" :focusable="false" plain @click="onSaveValue">
                    <template #icon>
                        <n-icon :component="Save" size="18" />
                    </template>
                    {{ $t('interface.save_update') }}
                </n-button>
                <n-button :loading="saving" :focusable="false" plain @click="onCancelEdit">
                    <template #icon>
                        <n-icon :component="Close" size="18" />
                    </template>
                    {{ $t('common.cancel') }}
                </n-button>
            </n-button-group>
        </div>
        <div class="value-wrapper flex-item-expand flex-box-v">
            <n-scrollbar v-if="!inEdit" class="flex-item-expand">
                <n-code :code="viewValue" :language="viewLanguage" show-line-numbers style="cursor: text" word-wrap />
            </n-scrollbar>
            <n-input
                v-else
                v-model:value="editValue"
                :disabled="saving"
                :resizable="false"
                class="flex-item-expand"
                type="textarea" />
        </div>
    </div>
</template>

<style lang="scss" scoped>
.value-wrapper {
    overflow: hidden;
    border-top: v-bind('themeVars.borderColor') 1px solid;
    padding: 5px;
}
</style>
