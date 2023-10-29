<script setup>
import { computed, ref } from 'vue'
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
import useConnectionStore from 'stores/connections.js'
import DropdownSelector from '@/components/content_value/DropdownSelector.vue'
import Code from '@/components/icons/Code.vue'
import Conversion from '@/components/icons/Conversion.vue'
import EditFile from '@/components/icons/EditFile.vue'

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
    value: String,
    size: Number,
    viewAs: {
        type: String,
        default: formatTypes.PLAIN_TEXT,
    },
    decode: {
        type: String,
        default: decodeTypes.NONE,
    },
})

/**
 *
 * @type {ComputedRef<string|number[]>}
 */
const keyName = computed(() => {
    return !isEmpty(props.keyCode) ? props.keyCode : props.keyPath
})

// const viewOption = computed(() =>
//     map(types, (t) => {
//         return {
//             value: t,
//             label: t,
//             key: t,
//         }
//     }),
// )

const keyType = redisTypes.STRING
const viewLanguage = computed(() => {
    switch (props.viewAs) {
        case formatTypes.JSON:
            return 'json'
        default:
            return 'plaintext'
    }
})

const onViewTypeUpdate = (viewType) => {
    connectionStore.loadKeyValue(props.name, props.db, keyName.value, viewType, props.decode)
}

const onDecodeTypeUpdate = (decodeType) => {
    connectionStore.loadKeyValue(props.name, props.db, keyName.value, props.viewAs, decodeType)
}

/**
 * Copy value
 */
const onCopyValue = () => {
    ClipboardSetText(props.value)
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
    editValue.value = props.value
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
            keyName.value,
            toLower(keyType),
            editValue.value,
            -1,
            props.viewAs,
            props.decode,
        )
        if (success) {
            await connectionStore.loadKeyValue(props.name, props.db, keyName.value)
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
        <content-toolbar
            :db="props.db"
            :key-code="keyCode"
            :key-path="keyPath"
            :key-type="keyType"
            :server="props.name"
            :ttl="ttl"
            class="value-item-part" />
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
                <n-code :code="props.value" :language="viewLanguage" show-line-numbers style="cursor: text" word-wrap />
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
            <div class="flex-item-expand"></div>
            <dropdown-selector
                :icon="Code"
                :options="formatTypes"
                :tooltip="$t('interface.view_as')"
                :value="props.viewAs"
                @update:value="onViewTypeUpdate" />

            <n-divider vertical />

            <dropdown-selector
                :icon="Conversion"
                :options="decodeTypes"
                :tooltip="$t('interface.decode_with')"
                :value="props.decode"
                @update:value="onDecodeTypeUpdate" />
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
    background-color: v-bind('themeVars.bodyColor');
}
</style>
