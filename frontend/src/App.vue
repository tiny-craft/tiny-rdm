<script setup>
import ConnectionDialog from './components/dialogs/ConnectionDialog.vue'
import NewKeyDialog from './components/dialogs/NewKeyDialog.vue'
import PreferencesDialog from './components/dialogs/PreferencesDialog.vue'
import RenameKeyDialog from './components/dialogs/RenameKeyDialog.vue'
import SetTtlDialog from './components/dialogs/SetTtlDialog.vue'
import hljs from 'highlight.js/lib/core'
import json from 'highlight.js/lib/languages/json'
import plaintext from 'highlight.js/lib/languages/plaintext'
import AddFieldsDialog from './components/dialogs/AddFieldsDialog.vue'
import AppContent from './AppContent.vue'
import GroupDialog from './components/dialogs/GroupDialog.vue'
import DeleteKeyDialog from './components/dialogs/DeleteKeyDialog.vue'
import { onMounted, ref, watch } from 'vue'
import usePreferencesStore from './stores/preferences.js'
import useConnectionStore from './stores/connections.js'
import { useI18n } from 'vue-i18n'
import { darkTheme } from 'naive-ui'
import KeyFilterDialog from './components/dialogs/KeyFilterDialog.vue'
import { WindowSetDarkTheme, WindowSetLightTheme } from 'wailsjs/runtime/runtime.js'
import { themeOverrides } from '@/utils/theme.js'

hljs.registerLanguage('json', json)
hljs.registerLanguage('plaintext', plaintext)

const prefStore = usePreferencesStore()
const connectionStore = useConnectionStore()
const i18n = useI18n()
const initializing = ref(false)
onMounted(async () => {
    try {
        initializing.value = true
        await prefStore.loadFontList()
        await connectionStore.initConnections()
        if (prefStore.autoCheckUpdate) {
            prefStore.checkForUpdate()
        }
    } finally {
        initializing.value = false
    }
})

// watch theme and dynamically switch
watch(
    () => prefStore.isDark,
    (isDark) => (isDark ? WindowSetDarkTheme() : WindowSetLightTheme()),
)

// watch language and dynamically switch
watch(
    () => prefStore.general.language,
    (lang) => (i18n.locale.value = prefStore.currentLanguage),
)
</script>

<template>
    <n-config-provider
        :hljs="hljs"
        :inline-theme-disabled="true"
        :theme="prefStore.isDark ? darkTheme : undefined"
        :theme-overrides="themeOverrides"
        :locale="prefStore.themeLocale"
        class="fill-height">
        <n-global-style />
        <n-dialog-provider>
            <n-spin v-show="initializing" :theme-overrides="{ opacitySpinning: 0 }" style="--wails-draggable: drag">
                <div id="launch-container" />
            </n-spin>
            <app-content v-if="!initializing" class="flex-item-expand" />

            <!-- top modal dialogs -->
            <connection-dialog />
            <group-dialog />
            <new-key-dialog />
            <key-filter-dialog />
            <add-fields-dialog />
            <rename-key-dialog />
            <delete-key-dialog />
            <set-ttl-dialog />
            <preferences-dialog />
        </n-dialog-provider>
    </n-config-provider>
</template>

<style lang="scss">
#launch-container {
    width: 100vw;
    height: 100vh;
}

#app-title {
    text-align: center;
    width: 100vw;
    height: 30px;
    line-height: 30px;
}
</style>
