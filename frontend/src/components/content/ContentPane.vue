<script setup>
import { computed, onMounted, onUnmounted, reactive, watch } from 'vue'
import { types } from '@/consts/support_redis_type.js'
import ContentValueHash from '@/components/content_value/ContentValueHash.vue'
import ContentValueList from '@/components/content_value/ContentValueList.vue'
import ContentValueString from '@/components/content_value/ContentValueString.vue'
import ContentValueSet from '@/components/content_value/ContentValueSet.vue'
import ContentValueZset from '@/components/content_value/ContentValueZSet.vue'
import { get, isEmpty, keyBy, map, size, toUpper } from 'lodash'
import useTabStore from 'stores/tab.js'
import useConnectionStore from 'stores/connections.js'
import ContentServerStatus from '@/components/content_value/ContentServerStatus.vue'
import ContentValueStream from '@/components/content_value/ContentValueStream.vue'

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

const valueComponents = {
    [types.STRING]: ContentValueString,
    [types.HASH]: ContentValueHash,
    [types.LIST]: ContentValueList,
    [types.SET]: ContentValueSet,
    [types.ZSET]: ContentValueZset,
    [types.STREAM]: ContentValueStream,
}

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

const showNonexists = computed(() => {
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
</script>

<template>
    <div class="content-container flex-box-v">
        <div v-if="showServerStatus" class="content-container flex-item-expand flex-box-v">
            <!-- select nothing or select server node, display server status -->
            <content-server-status
                v-model:auto-refresh="currentServer.autoRefresh"
                :info="currentServer.info"
                :loading="currentServer.loading"
                :auto-loading="currentServer.autoLoading"
                :server="currentServer.name"
                @refresh="refreshInfo(currentServer.name, true)" />
        </div>
        <div v-else-if="showNonexists" class="content-container flex-item-expand flex-box-v">
            <n-empty :description="$t('interface.nonexist_tab_content')" class="empty-content">
                <template #extra>
                    <n-button :focusable="false" @click="onReloadKey">{{ $t('interface.reload') }}</n-button>
                </template>
            </n-empty>
        </div>
        <component
            :is="valueComponents[tabContent.type]"
            v-else
            :db="tabContent.db"
            :key-path="tabContent.keyPath"
            :key-code="tabContent.keyCode"
            :name="tabContent.name"
            :ttl="tabContent.ttl"
            :value="tabContent.value"
            :size="tabContent.size"
            :view-as="tabContent.viewAs" />
    </div>
</template>

<style lang="scss" scoped>
@import '@/styles/content';

.content-container {
    padding: 5px;
    box-sizing: border-box;
}

//.tab-item {
//    gap: 5px;
//    padding: 0 5px 0 10px;
//    align-items: center;
//    max-width: 150px;
//
//    transition: all var(--transition-duration-fast) var(--transition-function-ease-in-out-bezier);
//
//    &-label {
//        font-size: 15px;
//        text-align: center;
//    }
//
//    &-close {
//        &:hover {
//            background-color: rgb(176, 177, 182, 0.4);
//        }
//    }
//}
</style>
