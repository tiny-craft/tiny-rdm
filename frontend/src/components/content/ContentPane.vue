<script setup>
import { computed, onMounted, onUnmounted, reactive, watch } from 'vue'
import { get, isEmpty, keyBy, map, size, toUpper } from 'lodash'
import useTabStore from 'stores/tab.js'
import useConnectionStore from 'stores/connections.js'
import ContentServerStatus from '@/components/content_value/ContentServerStatus.vue'
import Status from '@/components/icons/Status.vue'
import { useThemeVars } from 'naive-ui'
import { BrowserTabType } from '@/consts/browser_tab_type.js'
import Terminal from '@/components/icons/Terminal.vue'
import Log from '@/components/icons/Log.vue'
import Detail from '@/components/icons/Detail.vue'
import ContentValueWrapper from '@/components/content_value/ContentValueWrapper.vue'
import ContentCli from '@/components/content_value/ContentCli.vue'

const themeVars = useThemeVars()

/**
 * @typedef {Object} ServerStatusItem
 * @property {string} name
 * @property {Object} info
 * @property {boolean} autoRefresh
 * @property {boolean} loading loading status for refresh
 * @property {boolean} autoLoading loading status for auto refresh
 */

/**
 *
 * @type {UnwrapNestedRefs<Object.<string, ServerStatusItem>>}
 */
const serverStatusTab = reactive({})

/**
 *
 * @param {string} serverName
 * @return {UnwrapRef<ServerStatusItem>}
 */
const getServerInfo = (serverName) => {
    if (isEmpty(serverName)) {
        return {
            name: serverName,
            info: {},
            autoRefresh: false,
            autoLoading: false,
            loading: false,
        }
    }
    if (!serverStatusTab.hasOwnProperty(serverName)) {
        serverStatusTab[serverName] = {
            name: serverName,
            info: {},
            autoRefresh: false,
            autoLoading: false,
            loading: false,
        }
    }
    return serverStatusTab[serverName]
}
const serverName = computed(() => {
    if (tabContent.value != null) {
        return tabContent.value.name
    }
    return ''
})
/**
 *
 * @type {ComputedRef<ServerStatusItem>}
 */
const currentServer = computed(() => {
    return getServerInfo(serverName.value)
})

/**
 * refresh server status info
 * @param {string} serverName
 * @param {boolean} [force] force refresh will show loading indicator
 * @returns {Promise<void>}
 */
const refreshInfo = async (serverName, force) => {
    const info = getServerInfo(serverName)
    if (force) {
        info.loading = true
    } else {
        info.autoLoading = true
    }
    if (!isEmpty(serverName) && connectionStore.isConnected(serverName)) {
        try {
            info.info = await connectionStore.getServerInfo(serverName)
        } finally {
            info.loading = false
            info.autoLoading = false
        }
    }
}

const refreshAllInfo = async (force) => {
    for (const serverName in serverStatusTab) {
        await refreshInfo(serverName, force)
    }
}

let intervalId
onMounted(() => {
    refreshAllInfo(true)
    intervalId = setInterval(() => {
        for (const serverName in serverStatusTab) {
            if (get(serverStatusTab, [serverName, 'autoRefresh'])) {
                refreshInfo(serverName)
            }
        }
    }, 5000)
})

onUnmounted(() => {
    clearInterval(intervalId)
})

const connectionStore = useConnectionStore()
const tabStore = useTabStore()
const tab = computed(() =>
    map(tabStore.tabs, (item) => ({
        key: item.name,
        label: item.title,
    })),
)

watch(
    () => tabStore.nav,
    (nav) => {
        if (nav === 'browser') {
            refreshInfo(serverName.value)
        }
    },
)

watch(
    () => tabStore.tabList,
    (tabs) => {
        if (size(tabs) < size(serverStatusTab)) {
            const tabMap = keyBy(tabs, 'name')
            // remove unused server status tab
            for (const t in serverStatusTab) {
                if (!tabMap.hasOwnProperty(t)) {
                    delete serverStatusTab[t]
                }
            }
        }
    },
    { deep: true },
)

const tabContent = computed(() => {
    const tab = tabStore.currentTab
    if (tab == null) {
        return null
    }
    return {
        name: tab.name,
        subTab: tab.subTab,
        type: toUpper(tab.type),
        db: tab.db,
        keyPath: tab.key,
        keyCode: tab.keyCode,
        ttl: tab.ttl,
        value: tab.value,
        size: tab.size || 0,
        viewAs: tab.viewAs,
    }
})

const showServerStatus = computed(() => {
    return tabContent.value == null || isEmpty(tabContent.value.keyPath)
})

const isBlankValue = computed(() => {
    return tabContent.value.value == null
})

/**
 * reload current selection key
 * @returns {Promise<null>}
 */
const onReloadKey = async () => {
    const tab = tabStore.currentTab
    if (tab == null || isEmpty(tab.key)) {
        return null
    }
    await connectionStore.loadKeyValue(tab.name, tab.db, tab.key, tab.viewAs)
}

const selectedSubTab = computed(() => {
    const { subTab = 'status' } = tabStore.currentTab || {}
    return subTab
})

const onSwitchSubTab = (name) => {
    tabStore.switchSubTab(name)
}
</script>

<template>
    <div class="content-container flex-box-v">
        <n-tabs
            :tabs-padding="5"
            :theme-overrides="{
                tabFontWeightActive: 'normal',
                tabGapSmallLine: '10px',
                tabGapMediumLine: '10px',
                tabGapLargeLine: '10px',
            }"
            :value="selectedSubTab"
            class="content-sub-tab"
            default-value="status"
            pane-class="content-sub-tab-pane"
            placement="top"
            tab-style="padding-left: 10px; padding-right: 10px;"
            type="line"
            @update:value="onSwitchSubTab">
            <!-- server status pane -->
            <n-tab-pane :name="BrowserTabType.Status.toString()">
                <template #tab>
                    <n-space :size="5" :wrap-item="false" align="center" inline justify="center">
                        <n-icon size="16">
                            <status :inverse="selectedSubTab === BrowserTabType.Status.toString()" stroke-width="4" />
                        </n-icon>
                        <span>{{ $t('interface.sub_tab.status') }}</span>
                    </n-space>
                </template>
                <content-server-status
                    v-model:auto-refresh="currentServer.autoRefresh"
                    :auto-loading="currentServer.autoLoading"
                    :info="currentServer.info"
                    :loading="currentServer.loading"
                    :server="currentServer.name"
                    @refresh="refreshInfo(currentServer.name, true)" />
            </n-tab-pane>

            <!-- key detail pane -->
            <n-tab-pane :name="BrowserTabType.KeyDetail.toString()">
                <template #tab>
                    <n-space :size="5" :wrap-item="false" align="center" inline justify="center">
                        <n-icon size="16">
                            <detail
                                :inverse="selectedSubTab === BrowserTabType.KeyDetail.toString()"
                                fill-color="none" />
                        </n-icon>
                        <span>{{ $t('interface.sub_tab.key_detail') }}</span>
                    </n-space>
                </template>
                <content-value-wrapper
                    :blank="isBlankValue"
                    :type="tabContent.type"
                    :db="tabContent.db"
                    :key-code="tabContent.keyCode"
                    :key-path="tabContent.keyPath"
                    :name="tabContent.name"
                    :size="tabContent.size"
                    :ttl="tabContent.ttl"
                    :value="tabContent.value"
                    :view-as="tabContent.viewAs"
                    @reload="onReloadKey" />
            </n-tab-pane>

            <!-- cli pane -->
            <n-tab-pane :name="BrowserTabType.Cli.toString()">
                <template #tab>
                    <n-space :size="5" :wrap-item="false" align="center" inline justify="center">
                        <n-icon size="16">
                            <terminal :inverse="selectedSubTab === BrowserTabType.Cli.toString()" />
                        </n-icon>
                        <span>{{ $t('interface.sub_tab.cli') }}</span>
                    </n-space>
                </template>
                <content-cli :name="currentServer.name" />
            </n-tab-pane>

            <!-- slow log pane -->
            <n-tab-pane :name="BrowserTabType.SlowLog.toString()">
                <template #tab>
                    <n-space :size="5" :wrap-item="false" align="center" inline justify="center">
                        <n-icon size="16">
                            <log :inverse="selectedSubTab === BrowserTabType.SlowLog.toString()" />
                        </n-icon>
                        <span>{{ $t('interface.sub_tab.slow_log') }}</span>
                    </n-space>
                </template>
            </n-tab-pane>
        </n-tabs>
    </div>
</template>

<style lang="scss" scoped>
@import '@/styles/content';

.content-container {
    //padding: 5px 5px 0;
    //padding-top: 0;
    box-sizing: border-box;
    background-color: v-bind('themeVars.tabColor');
}
</style>

<style lang="scss">
.content-sub-tab {
    background-color: v-bind('themeVars.bodyColor');
    height: 100%;
}

.content-sub-tab-pane {
    padding: 0 !important;
    height: 100%;
    background-color: v-bind('themeVars.tabColor');
    overflow: hidden;
}
</style>
