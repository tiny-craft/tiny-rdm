<script setup>
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { types } from '@/consts/support_redis_type.js'
import ContentValueHash from '@/components/content_value/ContentValueHash.vue'
import ContentValueList from '@/components/content_value/ContentValueList.vue'
import ContentValueString from '@/components/content_value/ContentValueString.vue'
import ContentValueSet from '@/components/content_value/ContentValueSet.vue'
import ContentValueZset from '@/components/content_value/ContentValueZSet.vue'
import { isEmpty, map, toUpper } from 'lodash'
import useTabStore from 'stores/tab.js'
import useConnectionStore from 'stores/connections.js'
import ContentServerStatus from '@/components/content_value/ContentServerStatus.vue'
import ContentValueStream from '@/components/content_value/ContentValueStream.vue'

const serverInfo = ref({})
const autoRefresh = ref(false)
const serverName = computed(() => {
    if (tabContent.value != null) {
        return tabContent.value.name
    }
    return ''
})
const loadingServerInfo = ref(false)

/**
 * refresh server status info
 * @param {boolean} [force] force refresh will show loading indicator
 * @returns {Promise<void>}
 */
const refreshInfo = async (force) => {
    if (force) {
        loadingServerInfo.value = true
    }
    if (!isEmpty(serverName.value) && connectionStore.isConnected(serverName.value)) {
        try {
            serverInfo.value = await connectionStore.getServerInfo(serverName.value)
        } finally {
            loadingServerInfo.value = false
        }
    }
}

let intervalId
onMounted(() => {
    refreshInfo(true)
    intervalId = setInterval(() => {
        if (autoRefresh.value) {
            refreshInfo()
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
            refreshInfo()
        }
    },
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
        ttl: tab.ttl,
        value: tab.value,
    }
})

const showServerStatus = computed(() => {
    return tabContent.value == null || isEmpty(tabContent.value.keyPath)
})

const showNonexists = computed(() => {
    return tabContent.value.value == null
})

const onUpdateValue = (tabIndex) => {
    tabStore.switchTab(tabIndex)
}

/**
 * reload current selection key
 * @returns {Promise<null>}
 */
const onReloadKey = async () => {
    const tab = tabStore.currentTab
    if (tab == null || isEmpty(tab.key)) {
        return null
    }
    await connectionStore.loadKeyValue(tab.name, tab.db, tab.key)
}
</script>

<template>
    <div class="content-container flex-box-v">
        <div v-if="showServerStatus" class="content-container flex-item-expand flex-box-v">
            <!-- select nothing or select server node, display server status -->
            <content-server-status
                v-model:auto-refresh="autoRefresh"
                :info="serverInfo"
                :loading="loadingServerInfo"
                :server="serverName"
                @refresh="refreshInfo(true)" />
        </div>
        <div v-else-if="showNonexists" class="content-container flex-item-expand flex-box-v">
            <n-empty :description="$t('nonexist_tab_content')" class="empty-content">
                <template #extra>
                    <n-button @click="onReloadKey">{{ $t('reload') }}</n-button>
                </template>
            </n-empty>
        </div>
        <component
            :is="valueComponents[tabContent.type]"
            v-else
            :db="tabContent.db"
            :key-path="tabContent.keyPath"
            :name="tabContent.name"
            :ttl="tabContent.ttl"
            :value="tabContent.value" />
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
