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
import { computed, onBeforeMount, ref } from 'vue'
import { get } from 'lodash'
import usePreferencesStore from './stores/preferences.js'
import useConnectionStore from './stores/connections.js'
import { useI18n } from 'vue-i18n'
import { darkTheme, lightTheme, useOsTheme } from 'naive-ui'

hljs.registerLanguage('json', json)
hljs.registerLanguage('plaintext', plaintext)

/**
 *
 * @type import('naive-ui').GlobalThemeOverrides
 */
const themeOverrides = {
    common: {
        primaryColor: '#D33A31',
        primaryColorHover: '#FF6B6B',
        primaryColorPressed: '#D5271C',
        primaryColorSuppl: '#FF6B6B',
        borderRadius: '4px',
        borderRadiusSmall: '3px',
        lineHeight: 1.5,
        scrollbarWidth: '8px',
    },
    Tag: {
        // borderRadius: '3px'
    },
    Tabs: {
        tabGapSmallCard: '1px',
        tabGapMediumCard: '1px',
        tabGapLargeCard: '1px',
    },
}

const prefStore = usePreferencesStore()
const connectionStore = useConnectionStore()
const i18n = useI18n()
const initializing = ref(false)
onBeforeMount(async () => {
    try {
        initializing.value = true
        await prefStore.loadPreferences()
        i18n.locale.value = get(prefStore.general, 'language', 'en')
        await prefStore.loadFontList()
        await connectionStore.initConnections()
    } finally {
        initializing.value = false
    }
})

const osTheme = useOsTheme()
const theme = computed(() => {
    if (prefStore.general.theme === 'auto') {
        if (osTheme.value === 'dark') {
            return darkTheme
        }
    } else if (prefStore.general.theme === 'dark') {
        return darkTheme
    }
    return lightTheme
})
</script>

<template>
    <n-config-provider
        :hljs="hljs"
        :inline-theme-disabled="true"
        :theme="theme"
        :theme-overrides="themeOverrides"
        class="fill-height"
    >
        <n-global-style />
        <n-message-provider>
            <n-dialog-provider>
                <n-spin v-show="initializing" :theme-overrides="{ opacitySpinning: 0 }">
                    <div id="launch-container" />
                </n-spin>
                <app-content v-if="!initializing" class="flex-item-expand" />

                <!-- top modal dialogs -->
                <connection-dialog />
                <group-dialog />
                <new-key-dialog />
                <add-fields-dialog />
                <rename-key-dialog />
                <delete-key-dialog />
                <set-ttl-dialog />
                <preferences-dialog />
            </n-dialog-provider>
        </n-message-provider>
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
