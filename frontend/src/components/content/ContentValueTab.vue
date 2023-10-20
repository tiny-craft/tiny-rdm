<script setup>
import ToggleServer from '@/components/icons/ToggleServer.vue'
import useTabStore from 'stores/tab.js'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { get, map } from 'lodash'
import { useThemeVars } from 'naive-ui'
import useConnectionStore from 'stores/connections.js'

const themeVars = useThemeVars()
const i18n = useI18n()
const tabStore = useTabStore()
const connectionStore = useConnectionStore()

const onCloseTab = (tabIndex) => {
    $dialog.warning(i18n.t('dialogue.close_confirm'), () => {
        const tab = get(tabStore.tabs, tabIndex)
        if (tab != null) {
            connectionStore.closeConnection(tab.name)
        }
    })
}

const activeTabStyle = computed(() => {
    const { name } = tabStore.currentTab
    const { markColor = '' } = connectionStore.serverProfile[name] || {}
    return {
        backgroundColor: themeVars.value.bodyColor,
        borderTopWidth: markColor ? '3px' : '1px',
        borderTopColor: markColor || themeVars.value.borderColor,
        borderBottomColor: themeVars.value.bodyColor,
        borderTopLeftRadius: themeVars.value.borderRadius,
        borderTopRightRadius: themeVars.value.borderRadius,
    }
})
const inactiveTabStyle = computed(() => ({
    borderWidth: '0 0 1px',
    // borderBottomColor: themeVars.value.borderColor,
    borderTopLeftRadius: themeVars.value.borderRadius,
    borderTopRightRadius: themeVars.value.borderRadius,
}))

const tab = computed(() =>
    map(tabStore.tabs, (item) => ({
        key: item.name,
        label: item.title,
    })),
)
</script>

<template>
    <n-tabs
        v-model:value="tabStore.activatedIndex"
        :closable="true"
        :tab-style="{
            borderStyle: 'solid',
            borderWidth: '1px',
            borderLeftColor: themeVars.borderColor,
            borderRightColor: themeVars.borderColor,
        }"
        :theme-overrides="{
            tabFontWeightActive: 800,
            tabBorderRadius: 0,
            tabGapSmallCard: 0,
            tabGapMediumCard: 0,
            tabGapLargeCard: 0,
            tabColor: '#0000',
            // tabBorderColor: themeVars.borderColor,
            tabBorderColor: '#0000',
            tabTextColorCard: themeVars.closeIconColor,
        }"
        size="small"
        type="card"
        @close="onCloseTab"
        @update:value="(tabIndex) => tabStore.switchTab(tabIndex)">
        <n-tab
            v-for="(t, i) in tab"
            :key="i"
            :closable="tabStore.activatedIndex === i"
            :name="i"
            :style="tabStore.activatedIndex === i ? activeTabStyle : inactiveTabStyle"
            style="--wails-draggable: none"
            @dblclick.stop="() => {}">
            <n-space :size="5" :wrap-item="false" align="center" inline justify="center">
                <n-icon :component="ToggleServer" size="18" />
                <n-ellipsis style="max-width: 150px">{{ t.label }}</n-ellipsis>
            </n-space>
        </n-tab>
    </n-tabs>
</template>

<style lang="scss" scoped></style>
