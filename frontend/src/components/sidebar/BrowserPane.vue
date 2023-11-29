<script setup>
import { useThemeVars } from 'naive-ui'
import BrowserTree from './BrowserTree.vue'
import IconButton from '@/components/common/IconButton.vue'
import useTabStore from 'stores/tab.js'
import { computed, reactive, ref } from 'vue'
import { get } from 'lodash'
import Refresh from '@/components/icons/Refresh.vue'
import useDialogStore from 'stores/dialog.js'
import { useI18n } from 'vue-i18n'
import { types } from '@/consts/support_redis_type.js'
import Search from '@/components/icons/Search.vue'
import Unlink from '@/components/icons/Unlink.vue'
import Filter from '@/components/icons/Filter.vue'
import ContentSearchInput from '@/components/content_value/ContentSearchInput.vue'

const themeVars = useThemeVars()
const dialogStore = useDialogStore()
const tabStore = useTabStore()
const currentName = computed(() => get(tabStore.currentTab, 'name', ''))
const browserTreeRef = ref(null)

const onInfo = () => {
    browserTreeRef.value?.handleSelectContextMenu('server_info')
}

const i18n = useI18n()
const onDisconnect = () => {
    browserTreeRef.value?.handleSelectContextMenu('server_close')
}

const onRefresh = () => {
    browserTreeRef.value?.handleSelectContextMenu('server_reload')
}

const filterForm = reactive({
    showFilter: false,
    type: '',
    pattern: '',
})

const filterTypeOptions = computed(() => {
    const options = Object.keys(types).map((t) => ({
        value: t,
        label: t,
    }))
    options.splice(0, 0, {
        value: '',
        label: i18n.t('common.all'),
    })
    return options
})

// forbid dynamic switch key view due to performance issues
// const viewType = ref(0)
// const onSwitchView = (selectView) => {
//     const { server } = tabStore.currentTab
//     browserStore.switchKeyView(server, selectView)
// }
</script>

<template>
    <div class="nav-pane-container flex-box-v">
        <browser-tree ref="browserTreeRef" :server="currentName" />
        s
        <!-- bottom function bar -->
        <div class="nav-pane-bottom flex-box-v">
            <!--            <switch-button-->
            <!--                v-model:value="viewType"-->
            <!--                :icons="[TreeView, ListView]"-->
            <!--                :t-tooltips="['interface.tree_view', 'interface.list_view']"-->
            <!--                stroke-width="4"-->
            <!--                unselect-stroke-width="3"-->
            <!--                @update:value="onSwitchView" />-->
            <!--            <icon-button :icon="Status" size="20" stroke-width="4" t-tooltip="interface.status" @click="onInfo" />-->
            <div
                v-show="filterForm.showFilter"
                class="flex-box-h nav-pane-func"
                style="padding-left: 3px; padding-right: 3px">
                <!--                <n-input-group v-show="filterForm.showFilter">-->
                <!--                    <n-select-->
                <!--                        v-model:value="filterForm.type"-->
                <!--                        :consistent-menu-width="false"-->
                <!--                        :options="filterTypeOptions"-->
                <!--                        style="width: 120px" />-->
                <!--                    <n-input clearable placeholder="">-->
                <!--                        <template #prefix></template>-->
                <!--                    </n-input>-->
                <!--                    <n-button :focusable="false" ghost>-->
                <!--                        <template #icon>-->
                <!--                            <n-icon :component="Search" />-->
                <!--                        </template>-->
                <!--                    </n-button>-->
                <!--                </n-input-group>-->
                <content-search-input :full-search-icon="Search">
                    <template #prepend>
                        <n-select
                            v-model:value="filterForm.type"
                            :consistent-menu-width="false"
                            :options="filterTypeOptions"
                            style="width: 120px" />
                    </template>
                </content-search-input>
            </div>
            <div class="flex-box-h nav-pane-func">
                <icon-button
                    :button-class="{
                        'filter-on': filterForm.showFilter,
                        'filter-off': !filterForm.showFilter,
                        'toggle-btn': true,
                        'nav-pane-func-btn': true,
                    }"
                    :icon="Filter"
                    size="20"
                    stroke-width="4"
                    t-tooltip="interface.filter_key"
                    @click="filterForm.showFilter = !filterForm.showFilter" />
                <icon-button
                    :button-class="['nav-pane-func-btn']"
                    :icon="Refresh"
                    size="20"
                    stroke-width="4"
                    t-tooltip="interface.reload"
                    @click="onRefresh" />
                <div class="flex-item-expand" />
                <icon-button
                    :button-class="['nav-pane-func-btn']"
                    :icon="Unlink"
                    size="20"
                    stroke-width="4"
                    t-tooltip="interface.disconnect"
                    @click="onDisconnect" />
            </div>
        </div>
    </div>
</template>

<style lang="scss" scoped>
:deep(.toggle-btn) {
    border-style: solid;
    border-width: 1px;
}

:deep(.filter-on) {
    border-color: v-bind('themeVars.borderColor');
    background-color: v-bind('themeVars.borderColor');
}

:deep(.filter-off) {
    border-color: #0000;
}

.nav-pane-bottom {
    color: v-bind('themeVars.iconColor');
    border-top: v-bind('themeVars.borderColor') 1px solid;
}
</style>
