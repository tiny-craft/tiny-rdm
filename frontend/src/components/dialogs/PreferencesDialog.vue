<script setup>
import { reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { GetPreferences, RestoreDefault, SetPreferencesN } from '../../../wailsjs/go/storage/PreferencesStorage.js'
import { lang } from '../../langs/index'
import useDialog from '../../stores/dialog'

const langOption = Object.entries(lang).map(([key, value]) => ({
    value: key,
    label: `${value['lang_name']}`,
}))

const fontOption = [
    {
        label: 'JetBrains Mono',
        value: 'JetBrains Mono',
    },
]

const generalForm = reactive({
    language: langOption[0].value,
    font: '',
    fontSize: 14,
    useSystemProxy: false,
    useSystemProxyHttp: false,
    checkUpdate: false,
})

const editorForm = reactive({
    font: '',
    fontSize: 14,
})
const prevPreferences = ref({})
const tab = ref('general')
const formLabelWidth = '80px'
const dialogStore = useDialog()
const i18n = useI18n()

const applyPreferences = (pf) => {
    const { general = {}, editor = {} } = pf
    generalForm.language = general['language']
    generalForm.font = general['font']
    generalForm.fontSize = general['font_size'] || 14
    generalForm.useSystemProxy = general['use_system_proxy'] === true
    generalForm.useSystemProxyHttp = general['use_system_proxy_http'] === true
    generalForm.checkUpdate = general['check_update'] === true

    editorForm.font = editor['font']
    editorForm.fontSize = editor['font_size'] || 14
}

watch(
    () => dialogStore.preferencesDialogVisible,
    (visible) => {
        if (visible) {
            GetPreferences()
                .then((pf) => {
                    // load preferences from local
                    applyPreferences(pf)
                    prevPreferences.value = pf
                })
                .catch((e) => {
                    console.log(e)
                })
        }
    }
)

const onSavePreferences = async () => {
    const pf = {
        'general.language': generalForm.language,
        'general.font': generalForm.font,
        'general.font_size': generalForm.fontSize,
        'general.use_system_proxy': generalForm.useSystemProxy,
        'general.use_system_proxy_http': generalForm.useSystemProxyHttp,
        'general.check_update': generalForm.checkUpdate,

        'editor.font': editorForm.font,
        'editor.font_size': editorForm.fontSize,
    }
    await SetPreferencesN(pf)
    dialogStore.closePreferencesDialog()
}

// Watch language and dynamically switch
watch(
    () => generalForm.language,
    (lang) => (i18n.locale.value = lang)
)

watch(
    () => generalForm.font,
    (font) => {}
)

const onRestoreDefaults = async () => {
    const pf = await RestoreDefault()
    applyPreferences(pf)
}

const onClose = () => {
    dialogStore.closePreferencesDialog()
    // restore to old preferences
    applyPreferences(prevPreferences.value)
}
</script>

<template>
    <n-modal
        v-model:show="dialogStore.preferencesDialogVisible"
        :closable="false"
        :close-on-esc="false"
        :mask-closable="false"
        :show-icon="false"
        :title="$t('preferences')"
        preset="dialog"
        transform-origin="center"
    >
        <n-tabs v-model:value="tab" type="line" animated>
            <n-tab-pane :tab="$t('general')" display-directive="show" name="general">
                <n-form
                    :label-width="formLabelWidth"
                    :model="generalForm"
                    :show-require-mark="false"
                    label-align="right"
                    label-placement="left"
                >
                    <n-form-item :label="$t('language')" required>
                        <n-select v-model:value="generalForm.language" :options="langOption" filterable />
                    </n-form-item>
                    <n-form-item :label="$t('font')" required>
                        <n-select v-model:value="generalForm.font" :options="fontOption" filterable />
                    </n-form-item>
                    <n-form-item :label="$t('font_size')">
                        <n-input-number v-model:value="generalForm.fontSize" :max="65535" :min="1" />
                    </n-form-item>
                    <n-form-item :label="$t('proxy')">
                        <n-space>
                            <n-checkbox v-model:checked="generalForm.useSystemProxy">
                                {{ $t('use_system_proxy') }}
                            </n-checkbox>
                            <n-checkbox v-model:checked="generalForm.useSystemProxyHttp">
                                {{ $t('use_system_proxy_http') }}
                            </n-checkbox>
                        </n-space>
                    </n-form-item>
                    <n-form-item :label="$t('update')">
                        <n-checkbox v-model:checked="generalForm.checkUpdate"
                            >{{ $t('auto_check_update') }}
                        </n-checkbox>
                    </n-form-item>
                </n-form>
            </n-tab-pane>

            <n-tab-pane :tab="$t('editor')" display-directive="show" name="editor">
                <n-form
                    :label-width="formLabelWidth"
                    :model="editorForm"
                    :show-require-mark="false"
                    label-align="right"
                    label-placement="left"
                >
                    <n-form-item :label="$t('font')" :label-width="formLabelWidth" required>
                        <n-select v-model="editorForm.font" :options="fontOption" filterable />
                    </n-form-item>
                    <n-form-item :label="$t('font_size')" :label-width="formLabelWidth">
                        <n-input-number v-model="editorForm.fontSize" :max="65535" :min="1" />
                    </n-form-item>
                </n-form>
            </n-tab-pane>
        </n-tabs>

        <template #action>
            <div class="flex-item-expand">
                <n-button @click="onRestoreDefaults">{{ $t('restore_defaults') }}</n-button>
            </div>
            <div class="flex-item n-dialog__action">
                <n-button @click="onClose">{{ $t('cancel') }}</n-button>
                <n-button type="primary" @click="onSavePreferences">{{ $t('save') }}</n-button>
            </div>
        </template>
    </n-modal>
</template>

<style lang="scss" scoped>
.inline-form-item {
    padding-right: 10px;
}
</style>
