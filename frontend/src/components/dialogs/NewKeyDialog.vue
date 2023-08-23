<script setup>
import { computed, h, reactive, ref, watch } from 'vue'
import { types, typesColor } from '@/consts/support_redis_type.js'
import useDialog from 'stores/dialog'
import { isEmpty, keys, map } from 'lodash'
import NewStringValue from '@/components/new_value/NewStringValue.vue'
import NewHashValue from '@/components/new_value/NewHashValue.vue'
import NewListValue from '@/components/new_value/NewListValue.vue'
import NewZSetValue from '@/components/new_value/NewZSetValue.vue'
import NewSetValue from '@/components/new_value/NewSetValue.vue'
import { useI18n } from 'vue-i18n'
import useConnectionStore from 'stores/connections.js'
import { NSpace } from 'naive-ui'
import useTabStore from 'stores/tab.js'
import NewStreamValue from '@/components/new_value/NewStreamValue.vue'

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
    const requiredMsg = i18n.t('field_required')
    return {
        key: { required: true, message: requiredMsg, trigger: 'input' },
        type: { required: true, message: requiredMsg, trigger: 'input' },
        ttl: { required: true, message: requiredMsg, trigger: 'input' },
    }
})
const dbOptions = computed(() =>
    map(keys(connectionStore.databases[newForm.server]), (key) => ({
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
const addValueComponent = {
    [types.STRING]: NewStringValue,
    [types.HASH]: NewHashValue,
    [types.LIST]: NewListValue,
    [types.SET]: NewSetValue,
    [types.ZSET]: NewZSetValue,
    [types.STREAM]: NewStreamValue,
}
const defaultValue = {
    [types.STRING]: '',
    [types.HASH]: [],
    [types.LIST]: [],
    [types.SET]: [],
    [types.ZSET]: [],
    [types.STREAM]: [],
}

const dialogStore = useDialog()
watch(
    () => dialogStore.newKeyDialogVisible,
    (visible) => {
        if (visible) {
            const { prefix, server, db } = dialogStore.newKeyParam
            newForm.server = server
            newForm.key = isEmpty(prefix) ? '' : prefix
            newForm.db = db
            newForm.type = options.value[0].value
            newForm.ttl = -1
            newForm.value = null
        }
    },
)

const renderTypeLabel = (option) => {
    return h(
        NSpace,
        { align: 'center', inline: true, size: 3 },
        {
            default: () => [
                h('div', {
                    style: {
                        borderRadius: '50%',
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

const connectionStore = useConnectionStore()
const tabStore = useTabStore()
const onAdd = async () => {
    await newFormRef.value?.validate().catch((err) => {
        $message.error(err.message)
    })
    if (subFormRef.value?.validate && !subFormRef.value?.validate()) {
        $message.error(i18n.t('spec_field_required', { key: i18n.t('element') }))
        return false
    }
    try {
        const { server, db, key, type, ttl } = newForm
        let { value } = newForm
        if (value == null) {
            value = defaultValue[type]
        }
        const { success, msg, nodeKey } = await connectionStore.setKey(server, db, key, type, value, ttl)
        if (success) {
            // select current key
            tabStore.setSelectedKeys(server, nodeKey)
            connectionStore.loadKeyValue(server, db, key).then(() => {})
        } else if (!isEmpty(msg)) {
            $message.error(msg)
        }
        dialogStore.closeNewKeyDialog()
    } catch (e) {}
}

const onClose = () => {
    dialogStore.closeNewKeyDialog()
}
</script>

<template>
    <n-modal
        v-model:show="dialogStore.newKeyDialogVisible"
        :closable="false"
        :close-on-esc="false"
        :mask-closable="false"
        :negative-button-props="{ size: 'medium' }"
        :negative-text="$t('cancel')"
        :positive-button-props="{ size: 'medium' }"
        :positive-text="$t('confirm')"
        :show-icon="false"
        :title="$t('new_key')"
        preset="dialog"
        style="width: 600px"
        transform-origin="center"
        @positive-click="onAdd"
        @negative-click="onClose"
    >
        <n-scrollbar style="max-height: 500px">
            <n-form
                ref="newFormRef"
                :model="newForm"
                :rules="formRules"
                :show-require-mark="false"
                label-placement="top"
                style="padding-right: 15px"
            >
                <n-form-item :label="$t('key')" path="key" required>
                    <n-input v-model:value="newForm.key" placeholder="" />
                </n-form-item>
                <n-form-item :label="$t('db_index')" path="db" required>
                    <n-select v-model:value="newForm.db" :options="dbOptions" filterable />
                </n-form-item>
                <n-form-item :label="$t('type')" path="type" required>
                    <n-select v-model:value="newForm.type" :options="options" :render-label="renderTypeLabel" />
                </n-form-item>
                <n-form-item :label="$t('ttl')" required>
                    <n-input-group>
                        <n-input-number
                            v-model:value="newForm.ttl"
                            :max="Number.MAX_SAFE_INTEGER"
                            :min="-1"
                            placeholder="TTL"
                        >
                            <template #suffix>
                                {{ $t('second') }}
                            </template>
                        </n-input-number>
                        <n-button secondary type="primary" @click="newForm.ttl = -1">{{ $t('persist_key') }}</n-button>
                    </n-input-group>
                </n-form-item>
                <component :is="addValueComponent[newForm.type]" ref="subFormRef" v-model:value="newForm.value" />
                <!--  TODO: Add import from txt file option -->
            </n-form>
        </n-scrollbar>
    </n-modal>
</template>

<style lang="scss" scoped></style>
