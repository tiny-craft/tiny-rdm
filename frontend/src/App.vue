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

hljs.registerLanguage('json', json)
hljs.registerLanguage('plaintext', plaintext)

const themeOverrides = {
    common: {
        // primaryColor: '#409EFF',
        borderRadius: '4px',
        borderRadiusSmall: '3px',
        fontFamily: `"Nunito", -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto",
  "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue"`,
        lineHeight: 1.5,
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
</script>

<template>
    <n-config-provider :hljs="hljs" :inline-theme-disabled="true" :theme-overrides="themeOverrides" class="fill-height">
        <n-message-provider>
            <n-dialog-provider>
                <app-content />

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

<style lang="scss"></style>
