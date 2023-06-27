<script setup>
import { computed } from 'vue'
import { types } from '../consts/support_redis_type.js'
import ContentValueHash from './content_value/ContentValueHash.vue'
import ContentValueList from './content_value/ContentValueList.vue'
import ContentValueString from './content_value/ContentValueString.vue'
import ContentValueSet from './content_value/ContentValueSet.vue'
import ContentValueZset from './content_value/ContentValueZset.vue'
import { isEmpty, map, toUpper } from 'lodash'
import useTabStore from '../stores/tab.js'

const valueComponents = {
    [types.STRING]: ContentValueString,
    [types.HASH]: ContentValueHash,
    [types.LIST]: ContentValueList,
    [types.SET]: ContentValueSet,
    [types.ZSET]: ContentValueZset,
}

const tabStore = useTabStore()
const tab = computed(() =>
    map(tabStore.tabs, (item) => ({
        key: item.name,
        label: item.title,
    }))
)

/**
 *
 * @type {ComputedRef<TabItem>}
 */
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

const onUpdateValue = (tabIndex) => {
    tabStore.switchTab(tabIndex)
}

const onAddTab = () => {
    tabStore.newBlankTab()
}

const onCloseTab = (tabIndex) => {
    tabStore.removeTab(tabIndex)
    console.log('TODO: close connection also')
}
</script>

<template>
    <div class="content-container flex-box-v">
        <!--    <content-tab :model-value="tab"></content-tab>-->
        <n-tabs
            v-model:value="tabStore.activatedIndex"
            :closable="tab.length > 1"
            addable
            size="small"
            type="card"
            @add="onAddTab"
            @close="onCloseTab"
            @update:value="onUpdateValue"
        >
            <n-tab v-for="(t, i) in tab" :key="i" :name="i">
                <n-ellipsis style="max-width: 150px">{{ t.label }}</n-ellipsis>
            </n-tab>
        </n-tabs>
        <!-- add loading status -->
        <component
            :is="valueComponents[tabContent.type]"
            v-if="tabContent != null && !isEmpty(tabContent.keyPath)"
            :db="tabContent.db"
            :key-path="tabContent.keyPath"
            :name="tabContent.name"
            :ttl="tabContent.ttl"
            :value="tabContent.value"
        />
        <div v-else class="flex-item-expand flex-box-v">
            <n-empty :description="$t('empty_tab_content')" class="empty-content" />
        </div>
    </div>
</template>

<style lang="scss" scoped>
.content-container {
    height: 100%;
    overflow: hidden;
    background-color: var(--bg-color);
    padding-top: 2px;
    padding-bottom: 5px;
    box-sizing: border-box;
}

.empty-content {
    height: 100%;
    justify-content: center;
}

.tab-content {
}
</style>
