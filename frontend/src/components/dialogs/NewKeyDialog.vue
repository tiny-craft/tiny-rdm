<script setup>
import { computed, h, nextTick, reactive, ref, watchEffect } from 'vue'
import { types, typesColor } from '@/consts/support_redis_type.js'
import useDialog from 'stores/dialog'
import { endsWith, get, isEmpty, keys, map, trim } from 'lodash'
import NewStringValue from '@/components/new_value/NewStringValue.vue'
import NewHashValue from '@/components/new_value/NewHashValue.vue'
import NewListValue from '@/components/new_value/NewListValue.vue'
import NewZSetValue from '@/components/new_value/NewZSetValue.vue'
import NewSetValue from '@/components/new_value/NewSetValue.vue'
import { useI18n } from 'vue-i18n'
import { NSpace } from 'naive-ui'
import useTabStore from 'stores/tab.js'
import NewStreamValue from '@/components/new_value/NewStreamValue.vue'
import useBrowserStore from 'stores/browser.js'
import Import from '@/components/icons/Import.vue'
import NewJsonValue from '@/components/new_value/NewJsonValue.vue'

const i18n = useI18n()
const newForm = reactive({
    server: '',
    db: 0,
    key: '',
    type: '',
    ttl: -1,
    value: null,
})
const formRules = computed(() => {
    const requiredMsg = i18n.t('dialogue.field_required')
    return {
        key: { required: true, message: requiredMsg, trigger: 'input' },
        type: { required: true, message: requiredMsg, trigger: 'input' },
        ttl: { required: true, message: requiredMsg, trigger: 'input' },
    }
})
const dbOptions = computed(() =>
    map(keys(browserStore.getDBList(newForm.server)), (key) => ({
        label: key,
        value: parseInt(key),
    })),
)
const newFormRef = ref(null)
const subFormRef = ref(null)

const options = computed(() => {
    return Object.keys(types).map((t) => ({
        value: t,
        label: t,
    }))
})
const newValueComponent = {
    [types.STRING]: NewStringValue,
    [types.HASH]: NewHashValue,
    [types.LIST]: NewListValue,
    [types.SET]: NewSetValue,
    [types.ZSET]: NewZSetValue,
    [types.STREAM]: NewStreamValue,
    [types.JSON]: NewJsonValue,
}
const defaultValue = {
    [types.STRING]: '',
    [types.HASH]: [],
    [types.LIST]: [],
    [types.SET]: [],
    [types.ZSET]: [],
    [types.STREAM]: [],
    [types.JSON]: '{}',
}

const dialogStore = useDialog()
const scrollRef = ref(null)
watchEffect(() => {
    if (dialogStore.newKeyDialogVisible) {
        const { prefix, server, db } = dialogStore.newKeyParam
        const separator = browserStore.getSeparator(server)
        newForm.server = server
        if (isEmpty(prefix)) {
            newForm.key = ''
        } else {
            if (!endsWith(prefix, separator)) {
                newForm.key = prefix + separator
            } else {
                newForm.key = prefix
            }
        }
        newForm.db = db
        newForm.type = options.value[0].value
        newForm.ttl = -1
        newForm.value = null
    }
})

const renderTypeLabel = (option) => {
    return h(
        NSpace,
        {
            align: 'center',
            inline: true,
            size: 3,
            itemStyle: {
                lineHeight: 'var(--n-blank-height)',
            },
        },
        {
            default: () => [
                h('div', {
                    style: {
                        borderRadius: '9999px',
                        backgroundColor: typesColor[option.value],
                        width: '13px',
                        height: '13px',
                        border: '2px solid white',
                    },
                }),
                option.value,
            ],
        },
    )
}

const onAppend = () => {
    nextTick(() => {
        scrollRef.value?.scrollTo({ position: 'bottom' })
    })
}

const onChangeType = () => {
    newForm.value = null
}

const browserStore = useBrowserStore()
const tabStore = useTabStore()
const onAdd = async () => {
    await newFormRef.value?.validate((errs) => {
        const err = get(errs, '0.0.message')
        if (err != null) {
            $message.error(err)
        }
    })
    if (subFormRef.value?.validate) {
        await subFormRef.value?.validate((errs) => {
            const err = get(errs, '0.0.message')
            if (err != null) {
                $message.error(err)
            } else {
                $message.error(i18n.t('dialogue.spec_field_required', { key: i18n.t('dialogue.field.element') }))
            }
        })
    }
    try {
        const { server, db, key, type, ttl } = newForm
        let { value } = newForm
        if (value == null) {
            value = defaultValue[type]
        }
        // await browserStore.reloadKey({server, db, key: trim(key)})
        const { success, msg, nodeKey } = await browserStore.setKey({
            server,
            db,
            key: trim(key),
            keyType: type,
            value,
            ttl,
        })
        if (success) {
            // select current key
            await nextTick()
            tabStore.setSelectedKeys(server, nodeKey)
            browserStore.reloadKey({ server, db, key })
        } else if (!isEmpty(msg)) {
            $message.error(msg)
        }
        dialogStore.closeNewKeyDialog()
    } catch (e) {
        return false
    }
    return true
}

const onClose = () => {
    dialogStore.closeNewKeyDialog()
}

const onImport = () => {
    dialogStore.closeNewKeyDialog()
    dialogStore.openImportKeyDialog(newForm.server, newForm.db)
}
</script>

<template>
    <n-modal
        v-model:show="dialogStore.newKeyDialogVisible"
        :closable="false"
        :mask-closable="false"
        :show-icon="false"
        :title="$t('dialogue.key.new')"
        close-on-esc
        preset="dialog"
        style="width: 600px"
        transform-origin="center"
        @esc="onClose">
        <n-scrollbar ref="scrollRef" style="max-height: 500px">
            <n-form
                ref="newFormRef"
                :model="newForm"
                :rules="formRules"
                :show-require-mark="false"
                label-placement="top"
                style="padding-right: 15px">
                <n-form-item :label="$t('common.key')" path="key" required>
                    <n-input v-model:value="newForm.key" placeholder="" />
                </n-form-item>
                <n-form-item :label="$t('dialogue.key.db_index')" path="db" required>
                    <n-select v-model:value="newForm.db" :options="dbOptions" filterable />
                </n-form-item>
                <n-form-item :label="$t('interface.type')" path="type" required>
                    <n-select
                        v-model:value="newForm.type"
                        :options="options"
                        :render-label="renderTypeLabel"
                        @update:value="onChangeType" />
                </n-form-item>
                <n-form-item :label="$t('interface.ttl')" required>
                    <n-input-group>
                        <n-input-number
                            v-model:value="newForm.ttl"
                            :max="Number.MAX_SAFE_INTEGER"
                            :min="-1"
                            :show-button="false"
                            placeholder="TTL">
                            <template #suffix>
                                {{ $t('common.second') }}
                            </template>
                        </n-input-number>
                        <n-button :focusable="false" secondary type="primary" @click="() => (newForm.ttl = -1)">
                            {{ $t('interface.forever') }}
                        </n-button>
                    </n-input-group>
                </n-form-item>
                <component
                    :is="newValueComponent[newForm.type]"
                    ref="subFormRef"
                    v-model:value="newForm.value"
                    @append="onAppend" />
                <!--  TODO: Add import from txt file option -->
            </n-form>
        </n-scrollbar>

        <template #action>
            <div class="flex-item-expand">
                <n-button :focusable="false" size="medium" @click="onImport">
                    <template #icon>
                        <n-icon :component="Import" />
                    </template>
                    {{ $t('interface.import_key') }}
                </n-button>
            </div>
            <div class="flex-item n-dialog__action">
                <n-button :focusable="false" size="medium" @click="onClose">
                    {{ $t('common.cancel') }}
                </n-button>
                <n-button :focusable="false" size="medium" type="primary" @click="onAdd">
                    {{ $t('common.confirm') }}
                </n-button>
            </div>
        </template>
    </n-modal>
</template>

<style lang="scss" scoped></style>
