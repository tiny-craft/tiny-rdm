<script setup>
import { computed, reactive, watchEffect } from 'vue'
import useDialog from 'stores/dialog'
import useTabStore from 'stores/tab.js'
import Binary from '@/components/icons/Binary.vue'
import { isEmpty } from 'lodash'
import useBrowserStore from 'stores/browser.js'
import { useI18n } from 'vue-i18n'

const ttlForm = reactive({
    server: '',
    db: 0,
    key: '',
    keyCode: null,
    ttl: -1,
    unit: 1,
})

const dialogStore = useDialog()
const browserStore = useBrowserStore()
const tabStore = useTabStore()

watchEffect(() => {
    if (dialogStore.ttlDialogVisible) {
        // get ttl from current tab
        const tab = tabStore.currentTab
        if (tab != null) {
            ttlForm.server = tab.name
            ttlForm.db = tab.db
            ttlForm.key = tab.key
            ttlForm.keyCode = tab.keyCode
            ttlForm.unit = 1
            if (tab.ttl < 0) {
                // forever
                ttlForm.ttl = -1
            } else {
                ttlForm.ttl = tab.ttl
            }
        }
    }
})

const i18n = useI18n()
const unit = computed(() => [
    { value: 1, label: i18n.t('common.second') },
    {
        value: 60,
        label: i18n.t('common.minute'),
    },
    {
        value: 3600,
        label: i18n.t('common.hour'),
    },
    {
        value: 86400,
        label: i18n.t('common.day'),
    },
])

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
        const tab = tabStore.currentTab
        if (tab == null) {
            return
        }
        const key = isEmpty(ttlForm.keyCode) ? ttlForm.key : ttlForm.keyCode
        const ttl = ttlForm.ttl * (ttlForm.unit || 1)
        const success = await browserStore.setTTL(tab.name, tab.db, key, ttl)
        if (success) {
            tabStore.updateTTL({
                server: ttlForm.server,
                db: ttlForm.db,
                key: ttlForm.key,
                ttl: ttl,
            })
        }
    } catch (e) {
        $message.error(e.message || 'set ttl fail')
    } finally {
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
        :positive-button-props="{ focusable: false, size: 'medium' }"
        :positive-text="$t('common.save')"
        :show-icon="false"
        :title="$t('dialogue.ttl.title')"
        preset="dialog"
        transform-origin="center">
        <n-form :model="ttlForm" :show-require-mark="false" label-placement="top">
            <n-form-item :label="$t('common.key')">
                <n-input :value="ttlForm.key" readonly>
                    <template #prefix>
                        <n-icon v-if="!!ttlForm.keyCode" :component="Binary" size="20" />
                    </template>
                </n-input>
            </n-form-item>
            <n-form-item :label="$t('interface.ttl')" required>
                <n-input-group>
                    <n-input-number
                        v-model:value="ttlForm.ttl"
                        :max="Number.MAX_SAFE_INTEGER"
                        :min="-1"
                        :show-button="false"
                        class="flex-item-expand" />
                    <n-select v-model:value="ttlForm.unit" :options="unit" style="max-width: 150px" />
                </n-input-group>
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
