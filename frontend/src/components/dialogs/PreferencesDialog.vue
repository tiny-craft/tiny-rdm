<script setup>
import { computed, ref, watchEffect } from 'vue'
import { useI18n } from 'vue-i18n'
import useDialog from 'stores/dialog'
import usePreferencesStore from 'stores/preferences.js'
import { map, sortBy } from 'lodash'
import { typesIconStyle } from '@/consts/support_redis_type.js'

const prefStore = usePreferencesStore()

const prevPreferences = ref({})
const tab = ref('general')
const dialogStore = useDialog()
const i18n = useI18n()
const loading = ref(false)

const initPreferences = async () => {
    try {
        loading.value = true
        tab.value = 'general'
        await prefStore.loadPreferences()
        prevPreferences.value = {
            general: prefStore.general,
            editor: prefStore.editor,
        }
    } finally {
        loading.value = false
    }
}

watchEffect(() => {
    if (dialogStore.preferencesDialogVisible) {
        initPreferences()
    }
})

const keyOptions = computed(() => {
    const opts = map(typesIconStyle, (v) => ({
        value: v,
        label: i18n.t('preferences.general.key_icon_style' + v),
    }))
    return sortBy(opts, (o) => o.value)
})

const onSavePreferences = async () => {
    const success = await prefStore.savePreferences()
    if (success) {
        $message.success(i18n.t('dialogue.handle_succ'))
        dialogStore.closePreferencesDialog()
    }
}

const onClose = () => {
    // restore to old preferences
    prefStore.resetToLastPreferences()
    dialogStore.closePreferencesDialog()
}
</script>

<template>
    <n-modal
        v-model:show="dialogStore.preferencesDialogVisible"
        :auto-focus="false"
        :closable="false"
        :close-on-esc="false"
        :mask-closable="false"
        :show-icon="false"
        :title="$t('preferences.name')"
        preset="dialog"
        style="width: 500px"
        transform-origin="center">
        <!-- FIXME: set loading will slow down appear animation of dialog in linux -->
        <!-- <n-spin :show="loading"> -->
        <n-tabs v-model:value="tab" animated type="line">
            <n-tab-pane :tab="$t('preferences.general.name')" display-directive="show" name="general">
                <n-form :disabled="loading" :model="prefStore.general" :show-require-mark="false" label-placement="top">
                    <n-grid :x-gap="10">
                        <n-form-item-gi :label="$t('preferences.general.theme')" :span="24" required>
                            <n-radio-group v-model:value="prefStore.general.theme" name="theme" size="medium">
                                <n-radio-button
                                    v-for="opt in prefStore.themeOption"
                                    :key="opt.value"
                                    :value="opt.value">
                                    {{ opt.label }}
                                </n-radio-button>
                            </n-radio-group>
                        </n-form-item-gi>
                        <n-form-item-gi :label="$t('preferences.general.language')" :span="24" required>
                            <n-select
                                v-model:value="prefStore.general.language"
                                :options="prefStore.langOption"
                                filterable />
                        </n-form-item-gi>
                        <n-form-item-gi :label="$t('preferences.general.font')" :span="12" required>
                            <n-select
                                v-model:value="prefStore.general.font"
                                :options="prefStore.fontOption"
                                filterable />
                        </n-form-item-gi>
                        <n-form-item-gi :label="$t('preferences.general.font_size')" :span="12">
                            <n-input-number v-model:value="prefStore.general.fontSize" :max="65535" :min="1" />
                        </n-form-item-gi>
                        <n-form-item-gi :label="$t('preferences.general.scan_size')" :span="12">
                            <n-input-number v-model:value="prefStore.general.scanSize" :min="1" />
                        </n-form-item-gi>
                        <n-form-item-gi :label="$t('preferences.general.key_icon_style')" :span="12">
                            <n-select v-model:value="prefStore.general.keyIconStyle" :options="keyOptions" />
                        </n-form-item-gi>
                        <n-form-item-gi :label="$t('preferences.general.proxy')" :span="24">
                            <n-space>
                                <n-checkbox v-model:checked="prefStore.general.useSysProxy">
                                    {{ $t('preferences.general.use_system_proxy') }}
                                </n-checkbox>
                                <n-checkbox v-model:checked="prefStore.general.useSysProxyHttp">
                                    {{ $t('preferences.general.use_system_proxy_http') }}
                                </n-checkbox>
                            </n-space>
                        </n-form-item-gi>
                        <n-form-item-gi :label="$t('preferences.general.update')" :span="24">
                            <n-checkbox v-model:checked="prefStore.general.checkUpdate">
                                {{ $t('preferences.general.auto_check_update') }}
                            </n-checkbox>
                        </n-form-item-gi>
                    </n-grid>
                </n-form>
            </n-tab-pane>

            <n-tab-pane :tab="$t('preferences.editor.name')" display-directive="show" name="editor">
                <n-form :disabled="loading" :model="prefStore.editor" :show-require-mark="false" label-placement="top">
                    <n-form-item :label="$t('preferences.general.font')" required>
                        <n-select v-model:value="prefStore.editor.font" :options="prefStore.fontOption" filterable />
                    </n-form-item>
                    <n-form-item :label="$t('preferences.general.font_size')">
                        <n-input-number v-model:value="prefStore.editor.fontSize" :max="65535" :min="1" />
                    </n-form-item>
                    <n-checkbox v-model:checked="prefStore.editor.showLineNum">
                        {{ $t('preferences.editor.show_linenum') }}
                    </n-checkbox>
                </n-form>
            </n-tab-pane>
        </n-tabs>
        <!-- </n-spin> -->

        <template #action>
            <div class="flex-item-expand">
                <n-button :disabled="loading" @click="prefStore.restorePreferences">
                    {{ $t('preferences.restore_defaults') }}
                </n-button>
            </div>
            <div class="flex-item n-dialog__action">
                <n-button :disabled="loading" @click="onClose">{{ $t('common.cancel') }}</n-button>
                <n-button :disabled="loading" type="primary" @click="onSavePreferences">
                    {{ $t('common.save') }}
                </n-button>
            </div>
        </template>
    </n-modal>
</template>

<style lang="scss" scoped>
.inline-form-item {
    padding-right: 10px;
}
</style>
