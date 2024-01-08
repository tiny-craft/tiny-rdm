<script setup>
import { computed, reactive, ref, watchEffect } from 'vue'
import useDialog from 'stores/dialog'
import useBrowserStore from 'stores/browser.js'
import { useI18n } from 'vue-i18n'
import { isEmpty, size } from 'lodash'
import TtlInput from '@/components/common/TtlInput.vue'

const ttlForm = reactive({
    server: '',
    db: 0,
    key: '',
    keys: [],
    ttl: -1,
    unit: 1,
})

const dialogStore = useDialog()
const browserStore = useBrowserStore()

watchEffect(() => {
    if (dialogStore.ttlDialogVisible) {
        // get ttl from current tab
        const { server, db, key, keys, ttl } = dialogStore.ttlParam
        ttlForm.server = server
        ttlForm.db = db
        ttlForm.key = key
        ttlForm.keys = keys
        ttlForm.unit = 1
        if (ttl < 0) {
            // forever
            ttlForm.ttl = -1
        } else {
            ttlForm.ttl = ttl
        }
        procssing.value = false
    }
})

const procssing = ref(false)
const isBatchAction = computed(() => {
    return !isEmpty(ttlForm.keys)
})

const title = computed(() => {
    if (isBatchAction.value) {
        return i18n.t('dialogue.ttl.title_batch', { count: size(ttlForm.keys) })
    } else {
        return i18n.t('dialogue.ttl.title')
    }
})

const i18n = useI18n()
const quickOption = computed(() => [
    { value: -1, unit: 1, label: i18n.t('interface.forever') },
    { value: 10, unit: 1, label: `10 ${i18n.t('common.second')}` },
    { value: 1, unit: 60, label: `1 ${i18n.t('common.minute')}` },
    { value: 1, unit: 3600, label: `1 ${i18n.t('common.hour')}` },
    { value: 1, unit: 86400, label: `1 ${i18n.t('common.day')}` },
])

const onQuickSet = (opt) => {
    ttlForm.ttl = opt.value
    ttlForm.unit = opt.unit
}

const onClose = () => {
    dialogStore.closeTTLDialog()
}

const onConfirm = async () => {
    try {
        procssing.value = true
        const ttl = ttlForm.ttl * (ttlForm.unit || 1)
        let success = false
        if (isBatchAction.value) {
            success = await browserStore.setTTLs(ttlForm.server, ttlForm.db, ttlForm.keys, ttl)
        } else {
            success = await browserStore.setTTL(ttlForm.server, ttlForm.db, ttlForm.key, ttl)
        }
        if (success) {
        }
    } catch (e) {
        $message.error(e.message || 'set ttl fail')
    } finally {
        procssing.value = false
        dialogStore.closeTTLDialog()
    }
}
</script>

<template>
    <n-modal
        v-model:show="dialogStore.ttlDialogVisible"
        :closable="false"
        :close-on-esc="false"
        :mask-closable="false"
        :negative-button-props="{ focusable: false, size: 'medium' }"
        :negative-text="$t('common.cancel')"
        :on-negative-click="onClose"
        :on-positive-click="onConfirm"
        :positive-button-props="{ focusable: false, size: 'medium', loading: procssing }"
        :positive-text="$t('common.save')"
        :show-icon="false"
        :title="title"
        preset="dialog"
        transform-origin="center">
        <n-form :model="ttlForm" :show-require-mark="false" label-placement="top">
            <n-form-item v-if="!isBatchAction" :label="$t('common.key')">
                <n-input :value="ttlForm.key" readonly />
            </n-form-item>
            <n-form-item :label="$t('interface.ttl')" required>
                <ttl-input v-model:unit="ttlForm.unit" v-model:value="ttlForm.ttl" />
            </n-form-item>
            <n-form-item :label="$t('dialogue.ttl.quick_set')" :show-feedback="false">
                <n-space :wrap="true" :wrap-item="false">
                    <n-button
                        v-for="(opt, i) in quickOption"
                        :key="i"
                        round
                        secondary
                        size="small"
                        @click="onQuickSet(opt)">
                        {{ opt.label }}
                    </n-button>
                </n-space>
            </n-form-item>
        </n-form>
    </n-modal>
</template>

<style lang="scss" scoped></style>
