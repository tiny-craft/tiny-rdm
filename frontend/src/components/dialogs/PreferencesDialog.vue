<script setup>
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import useDialog from '../../stores/dialog'
import usePreferencesStore from '../../stores/preferences.js'
import { useMessage } from 'naive-ui'

const prefStore = usePreferencesStore()

const prevPreferences = ref({})
const tab = ref('general')
const formLabelWidth = '80px'
const dialogStore = useDialog()
const i18n = useI18n()
const message = useMessage()
const loading = ref(false)

watch(
    () => dialogStore.preferencesDialogVisible,
    async (visible) => {
        if (visible) {
            try {
                loading.value = true
                tab.value = 'general'
                await prefStore.loadFontList()
                await prefStore.loadPreferences()
                prevPreferences.value = {
                    general: prefStore.general,
                    editor: prefStore.editor,
                }
            } finally {
                loading.value = false
            }
        }
    },
)

const onSavePreferences = async () => {
    const success = await prefStore.savePreferences()
    if (success) {
        message.success(i18n.t('handle_succ'))
        dialogStore.closePreferencesDialog()
    }
}

// Watch language and dynamically switch
watch(
    () => prefStore.general.language,
    (lang) => (i18n.locale.value = prefStore.currentLanguage),
)

const onClose = () => {
    // restore to old preferences
    prefStore.resetToLastPreferences()
    dialogStore.closePreferencesDialog()
}
</script>

<template>
    <n-modal
        v-model:show="dialogStore.preferencesDialogVisible"
        :closable="false"
        :close-on-esc="false"
        :mask-closable="false"
        :show-icon="false"
        :auto-focus="false"
        :title="$t('preferences')"
        preset="dialog"
        transform-origin="center"
    >
        <n-spin :show="loading">
            <n-tabs v-model:value="tab" type="line" animated>
                <n-tab-pane :tab="$t('general')" display-directive="show" name="general">
                    <n-form
                        :label-width="formLabelWidth"
                        :model="prefStore.general"
                        :show-require-mark="false"
                        label-align="right"
                        label-placement="left"
                    >
                        <n-form-item :label="$t('theme')" required>
                            <n-radio-group v-model:value="prefStore.general.theme" name="theme" size="medium">
                                <n-radio-button
                                    v-for="opt in prefStore.themeOption"
                                    :key="opt.value"
                                    :value="opt.value"
                                >
                                    {{ opt.label }}
                                </n-radio-button>
                            </n-radio-group>
                        </n-form-item>
                        <n-form-item :label="$t('language')" required>
                            <n-select
                                v-model:value="prefStore.general.language"
                                :options="prefStore.langOption"
                                filterable
                            />
                        </n-form-item>
                        <n-form-item :label="$t('font')" required>
                            <n-select
                                v-model:value="prefStore.general.font"
                                :options="prefStore.fontOption"
                                filterable
                            />
                        </n-form-item>
                        <n-form-item :label="$t('font_size')">
                            <n-input-number v-model:value="prefStore.general.fontSize" :max="65535" :min="1" />
                        </n-form-item>
                        <n-form-item :label="$t('proxy')">
                            <n-space>
                                <n-checkbox v-model:checked="prefStore.general.useSysProxy">
                                    {{ $t('use_system_proxy') }}
                                </n-checkbox>
                                <n-checkbox v-model:checked="prefStore.general.useSysProxyHttp">
                                    {{ $t('use_system_proxy_http') }}
                                </n-checkbox>
                            </n-space>
                        </n-form-item>
                        <n-form-item :label="$t('update')">
                            <n-checkbox v-model:checked="prefStore.general.checkUpdate">
                                {{ $t('auto_check_update') }}
                            </n-checkbox>
                        </n-form-item>
                    </n-form>
                </n-tab-pane>

                <n-tab-pane :tab="$t('editor')" display-directive="show" name="editor">
                    <n-form
                        :label-width="formLabelWidth"
                        :model="prefStore.editor"
                        :show-require-mark="false"
                        label-align="right"
                        label-placement="left"
                    >
                        <n-form-item :label="$t('font')" required>
                            <n-select
                                v-model:value="prefStore.editor.font"
                                :options="prefStore.fontOption"
                                filterable
                            />
                        </n-form-item>
                        <n-form-item :label="$t('font_size')">
                            <n-input-number v-model:value="prefStore.editor.fontSize" :max="65535" :min="1" />
                        </n-form-item>
                    </n-form>
                </n-tab-pane>
            </n-tabs>
        </n-spin>

        <template #action>
            <div class="flex-item-expand">
                <n-button :disabled="loading" @click="prefStore.restorePreferences">{{
                    $t('restore_defaults')
                }}</n-button>
            </div>
            <div class="flex-item n-dialog__action">
                <n-button :disabled="loading" @click="onClose">{{ $t('cancel') }}</n-button>
                <n-button type="primary" :disabled="loading" @click="onSavePreferences">{{ $t('save') }}</n-button>
            </div>
        </template>
    </n-modal>
</template>

<style lang="scss" scoped>
.inline-form-item {
    padding-right: 10px;
}
</style>
