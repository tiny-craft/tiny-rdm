<script setup>
import { computed, h, ref, watchEffect } from 'vue'
import { useI18n } from 'vue-i18n'
import useDialog from 'stores/dialog'
import usePreferencesStore from 'stores/preferences.js'
import { find, map, sortBy } from 'lodash'
import { typesIconStyle } from '@/consts/support_redis_type.js'
import Help from '@/components/icons/Help.vue'
import Delete from '@/components/icons/Delete.vue'
import IconButton from '@/components/common/IconButton.vue'
import { NButton, NEllipsis, NIcon, NSpace, NTooltip } from 'naive-ui'
import Edit from '@/components/icons/Edit.vue'
import { joinCommand } from '@/utils/decoder_cmd.js'
import AddLink from '@/components/icons/AddLink.vue'
import Checked from '@/components/icons/Checked.vue'

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
            cli: prefStore.cli,
            decoder: prefStore.decoder,
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
        label: 'preferences.general.key_icon_style' + v,
    }))
    return sortBy(opts, (o) => o.value)
})

const decoderList = computed(() => {
    const decoder = prefStore.decoder || []
    const list = []
    for (const d of decoder) {
        // decode command
        list.push({
            name: d.name,
            auto: d.auto,
            decodeCmd: joinCommand(d.decodePath, d.decodeArgs),
            encodeCmd: joinCommand(d.encodePath, d.encodeArgs),
        })
    }
    return list
})

const decoderColumns = computed(() => {
    return [
        {
            key: 'name',
            title: () => i18n.t('preferences.decoder.decoder_name'),
            width: 120,
            align: 'center',
            titleAlign: 'center',
        },
        {
            key: 'cmd',
            title: () => i18n.t('preferences.decoder.cmd_preview'),
            titleAlign: 'center',
            render: ({ decodeCmd, encodeCmd }, index) => {
                return h(NSpace, { vertical: true, wrapItem: false, wrap: false, justify: 'center', size: 15 }, () => [
                    h(NEllipsis, {}, { default: () => decodeCmd, tooltip: () => decodeCmd + '\n\n' + encodeCmd }),
                    h(NEllipsis, {}, { default: () => encodeCmd, tooltip: () => decodeCmd + '\n\n' + encodeCmd }),
                ])
            },
        },
        {
            key: 'status',
            title: () => i18n.t('preferences.decoder.status'),
            width: 80,
            align: 'center',
            titleAlign: 'center',
            render: ({ auto }, index) => {
                if (auto) {
                    return h(
                        NTooltip,
                        { delay: 0, showArrow: false },
                        {
                            default: () => i18n.t('preferences.decoder.auto_enabled'),
                            trigger: () => h(NIcon, { component: Checked, size: 16 }),
                        },
                    )
                }
                return '-'
            },
        },
        {
            key: 'action',
            title: () => i18n.t('interface.action'),
            width: 80,
            align: 'center',
            titleAlign: 'center',
            render: ({ name, auto }, index) => {
                return h(NSpace, { wrapItem: false, wrap: false, justify: 'center', size: 'small' }, () => [
                    h(IconButton, {
                        icon: Delete,
                        tTooltip: 'interface.delete_row',
                        onClick: () => {
                            prefStore.removeCustomDecoder(name)
                        },
                    }),
                    h(IconButton, {
                        icon: Edit,
                        tTooltip: 'interface.edit_row',
                        onClick: () => {
                            const decoders = prefStore.decoder || []
                            const decoder = find(decoders, { name })
                            const { auto, decodePath, decodeArgs, encodePath, encodeArgs } = decoder
                            dialogStore.openDecoderDialog({
                                name,
                                auto,
                                decodePath,
                                decodeArgs,
                                encodePath,
                                encodeArgs,
                            })
                        },
                    }),
                ])
            },
        },
    ]
})

const onSavePreferences = async () => {
    const success = await prefStore.savePreferences()
    if (success) {
        // $message.success(i18n.t('dialogue.handle_succ'))
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
        style="width: 640px"
        transform-origin="center">
        <!-- FIXME: set loading will slow down appear animation of dialog in linux -->
        <!-- <n-spin :show="loading"> -->
        <n-tabs
            v-model:value="tab"
            animated
            pane-style="min-height: 300px"
            placement="left"
            tab-style="justify-content: right; font-weight: 420;"
            type="line">
            <!-- general pane -->
            <n-tab-pane :tab="$t('preferences.general.name')" display-directive="show" name="general">
                <n-form :disabled="loading" :model="prefStore.general" :show-require-mark="false" label-placement="top">
                    <n-grid :x-gap="10">
                        <n-form-item-gi :label="$t('preferences.general.theme')" :span="24" required>
                            <n-radio-group v-model:value="prefStore.general.theme" name="theme" size="medium">
                                <n-radio-button
                                    v-for="opt in prefStore.themeOption"
                                    :key="opt.value"
                                    :value="opt.value">
                                    {{ $t(opt.label) }}
                                </n-radio-button>
                            </n-radio-group>
                        </n-form-item-gi>
                        <n-form-item-gi :label="$t('preferences.general.language')" :span="24" required>
                            <n-select
                                v-model:value="prefStore.general.language"
                                :options="prefStore.langOption"
                                :render-label="({ label, value }) => (value === 'auto' ? $t(label) : label)"
                                filterable />
                        </n-form-item-gi>
                        <n-form-item-gi :span="24" required>
                            <template #label>
                                {{ $t('preferences.general.font') }}
                                <n-tooltip trigger="hover">
                                    <template #trigger>
                                        <n-icon :component="Help" />
                                    </template>
                                    <div class="text-block">
                                        {{ $t('preferences.font_tip') }}
                                    </div>
                                </n-tooltip>
                            </template>
                            <n-select
                                v-model:value="prefStore.general.fontFamily"
                                :options="prefStore.fontOption"
                                :placeholder="$t('preferences.general.font_tip')"
                                :render-label="({ label, value }) => (value === '' ? $t(label) : label)"
                                filterable
                                multiple
                                tag />
                        </n-form-item-gi>
                        <n-form-item-gi :label="$t('preferences.general.font_size')" :span="24">
                            <n-input-number v-model:value="prefStore.general.fontSize" :max="65535" :min="1" />
                        </n-form-item-gi>
                        <n-form-item-gi :label="$t('preferences.general.scan_size')" :span="12">
                            <n-input-number
                                v-model:value="prefStore.general.scanSize"
                                :min="1"
                                :show-button="false"
                                style="width: 100%" />
                        </n-form-item-gi>
                        <n-form-item-gi :label="$t('preferences.general.key_icon_style')" :span="12">
                            <n-select
                                v-model:value="prefStore.general.keyIconStyle"
                                :options="keyOptions"
                                :render-label="({ label }) => $t(label)" />
                        </n-form-item-gi>
                        <!--                        <n-form-item-gi :label="$t('preferences.general.proxy')" :span="24">-->
                        <!--                            <n-space>-->
                        <!--                                <n-checkbox v-model:checked="prefStore.general.useSysProxy">-->
                        <!--                                    {{ $t('preferences.general.use_system_proxy') }}-->
                        <!--                                </n-checkbox>-->
                        <!--                                <n-checkbox v-model:checked="prefStore.general.useSysProxyHttp">-->
                        <!--                                    {{ $t('preferences.general.use_system_proxy_http') }}-->
                        <!--                                </n-checkbox>-->
                        <!--                            </n-space>-->
                        <!--                        </n-form-item-gi>-->
                        <n-form-item-gi :label="$t('preferences.general.update')" :span="24">
                            <n-checkbox v-model:checked="prefStore.general.checkUpdate">
                                {{ $t('preferences.general.auto_check_update') }}
                            </n-checkbox>
                        </n-form-item-gi>
                    </n-grid>
                </n-form>
            </n-tab-pane>

            <!-- editor pane -->
            <n-tab-pane :tab="$t('preferences.editor.name')" display-directive="show" name="editor">
                <n-form :disabled="loading" :model="prefStore.editor" :show-require-mark="false" label-placement="top">
                    <n-grid :x-gap="10">
                        <n-form-item-gi :span="24" required>
                            <template #label>
                                {{ $t('preferences.general.font') }}
                                <n-tooltip trigger="hover">
                                    <template #trigger>
                                        <n-icon :component="Help" />
                                    </template>
                                    <div class="text-block">
                                        {{ $t('preferences.font_tip') }}
                                    </div>
                                </n-tooltip>
                            </template>
                            <n-select
                                v-model:value="prefStore.editor.fontFamily"
                                :options="prefStore.fontOption"
                                :placeholder="$t('preferences.general.font_tip')"
                                :render-label="({ label, value }) => value || $t(label)"
                                filterable
                                multiple
                                tag />
                        </n-form-item-gi>
                        <n-form-item-gi :label="$t('preferences.general.font_size')" :span="24">
                            <n-input-number v-model:value="prefStore.editor.fontSize" :max="65535" :min="1" />
                        </n-form-item-gi>
                        <n-form-item-gi :show-feedback="false" :show-label="false" :span="24">
                            <n-checkbox v-model:checked="prefStore.editor.showLineNum">
                                {{ $t('preferences.editor.show_linenum') }}
                            </n-checkbox>
                        </n-form-item-gi>
                        <n-form-item-gi :show-feedback="false" :show-label="false" :span="24">
                            <n-checkbox v-model:checked="prefStore.editor.showFolding">
                                {{ $t('preferences.editor.show_folding') }}
                            </n-checkbox>
                        </n-form-item-gi>
                        <n-form-item-gi :show-feedback="false" :show-label="false" :span="24">
                            <n-checkbox v-model:checked="prefStore.editor.dropText">
                                {{ $t('preferences.editor.drop_text') }}
                            </n-checkbox>
                        </n-form-item-gi>
                        <n-form-item-gi :show-feedback="false" :show-label="false" :span="24">
                            <n-checkbox v-model:checked="prefStore.editor.links">
                                {{ $t('preferences.editor.links') }}
                            </n-checkbox>
                        </n-form-item-gi>
                    </n-grid>
                </n-form>
            </n-tab-pane>

            <!-- cli pane -->
            <n-tab-pane :tab="$t('preferences.cli.name')" display-directive="show" name="cli">
                <n-form :disabled="loading" :model="prefStore.cli" :show-require-mark="false" label-placement="top">
                    <n-grid :x-gap="10">
                        <n-form-item-gi :span="24" required>
                            <template #label>
                                {{ $t('preferences.general.font') }}
                                <n-tooltip trigger="hover">
                                    <template #trigger>
                                        <n-icon :component="Help" />
                                    </template>
                                    <div class="text-block">
                                        {{ $t('preferences.font_tip') }}
                                    </div>
                                </n-tooltip>
                            </template>
                            <n-select
                                v-model:value="prefStore.cli.fontFamily"
                                :options="prefStore.fontOption"
                                :placeholder="$t('preferences.general.font_tip')"
                                :render-label="({ label, value }) => value || $t(label)"
                                filterable
                                multiple
                                tag />
                        </n-form-item-gi>
                        <n-form-item-gi :label="$t('preferences.general.font_size')" :span="24">
                            <n-input-number v-model:value="prefStore.cli.fontSize" :max="65535" :min="1" />
                        </n-form-item-gi>
                        <n-form-item-gi :label="$t('preferences.cli.cursor_style')" :span="24">
                            <n-radio-group v-model:value="prefStore.cli.cursorStyle" name="theme" size="medium">
                                <n-radio-button
                                    v-for="opt in prefStore.cliCursorStyleOption"
                                    :key="opt.value"
                                    :value="opt.value">
                                    {{ $t(opt.label) }}
                                </n-radio-button>
                            </n-radio-group>
                        </n-form-item-gi>
                    </n-grid>
                </n-form>
            </n-tab-pane>

            <!-- custom decoder pane -->
            <n-tab-pane :tab="$t('preferences.decoder.name')" display-directive="show:lazy" name="decoder">
                <n-space>
                    <n-button @click="dialogStore.openDecoderDialog()">
                        <template #icon>
                            <n-icon :component="AddLink" size="18" />
                        </template>
                        {{ $t('preferences.decoder.new') }}
                    </n-button>
                    <n-data-table
                        :columns="decoderColumns"
                        :data="decoderList"
                        :single-line="false"
                        max-height="350px" />
                </n-space>
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
