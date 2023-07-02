<script setup>
import { computed } from 'vue'
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
        <!--    <content-tab :model-value="tab"></content-tab>-->
        <n-tabs
            v-model:value="tabStore.activatedIndex"
            :closable="true"
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
        <!-- TODO: add loading status -->
        <component
            v-if="tabContent != null && !isEmpty(tabContent.keyPath)"
            :is="valueComponents[tabContent.type]"
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
@import 'content';
</style>
