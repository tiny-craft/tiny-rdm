<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { types } from '../../consts/support_redis_type.js'
import ContentValueHash from '../content_value/ContentValueHash.vue'
import ContentValueList from '../content_value/ContentValueList.vue'
import ContentValueString from '../content_value/ContentValueString.vue'
import ContentValueSet from '../content_value/ContentValueSet.vue'
import ContentValueZset from '../content_value/ContentValueZSet.vue'
import { get, isEmpty, map, toUpper } from 'lodash'
import useTabStore from '../../stores/tab.js'
import { useDialog } from 'naive-ui'
import useConnectionStore from '../../stores/connections.js'
import { useI18n } from 'vue-i18n'
import { useConfirmDialog } from '../../utils/confirm_dialog.js'
import ContentServerStatus from '../content_value/ContentServerStatus.vue'

const serverInfo = ref({})
const autoRefresh = ref(false)
const serverName = computed(() => {
    if (tabContent.value != null) {
        return tabContent.value.name
    }
    return ''
})
const refreshInfo = async () => {
    if (!isEmpty(serverName.value) && connectionStore.isConnected(serverName.value)) {
        serverInfo.value = await connectionStore.getServerInfo(serverName.value)
    }
}

let intervalId
onMounted(() => {
    refreshInfo()
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
}

const dialog = useDialog()
const connectionStore = useConnectionStore()
const tabStore = useTabStore()
const tab = computed(() =>
    map(tabStore.tabs, (item) => ({
        key: item.name,
        label: item.title,
    }))
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

const i18n = useI18n()
const confirmDialog = useConfirmDialog()
const onCloseTab = (tabIndex) => {
    confirmDialog.warning(i18n.t('close_confirm'), () => {
        const tab = get(tabStore.tabs, tabIndex)
        if (tab != null) {
            connectionStore.closeConnection(tab.name)
        }
    })
}
</script>

<template>
    <div class="content-container flex-box-v">
        <!--        <content-tab :model-value="tab"></content-tab>-->
        <n-tabs
            v-model:value="tabStore.activatedIndex"
            :closable="true"
            size="small"
            type="card"
            @close="onCloseTab"
            @update:value="onUpdateValue"
        >
            <n-tab v-for="(t, i) in tab" :key="i" :name="i">
                <n-ellipsis style="max-width: 150px">{{ t.label }}</n-ellipsis>
            </n-tab>
        </n-tabs>
        <!-- TODO: add loading status -->

        <div v-if="showServerStatus" class="content-container flex-item-expand flex-box-v">
            <!-- select nothing or select server node, display server status -->
            <content-server-status v-model:auto-refresh="autoRefresh" :server="serverName" :info="serverInfo" />
        </div>
        <div v-else-if="showNonexists" class="content-container flex-item-expand flex-box-v">
            <n-empty :description="$t('nonexist_tab_content')" class="empty-content">
                <template #extra>
                    <n-button @click="onReloadKey">{{ $t('reload') }}</n-button>
                </template>
            </n-empty>
        </div>
        <component
            v-else
            :is="valueComponents[tabContent.type]"
            :db="tabContent.db"
            :key-path="tabContent.keyPath"
            :name="tabContent.name"
            :ttl="tabContent.ttl"
            :value="tabContent.value"
        />
    </div>
</template>

<style lang="scss" scoped>
@import 'content';

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
