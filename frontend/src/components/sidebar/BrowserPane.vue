<script setup>
import { NIcon, useThemeVars } from 'naive-ui'
import BrowserTree from './BrowserTree.vue'
import IconButton from '@/components/common/IconButton.vue'
import useTabStore from 'stores/tab.js'
import { computed, reactive, ref } from 'vue'
import { get } from 'lodash'
import Refresh from '@/components/icons/Refresh.vue'
import useDialogStore from 'stores/dialog.js'
import { useI18n } from 'vue-i18n'
import useConnectionStore from 'stores/connections.js'
import { types } from '@/consts/support_redis_type.js'
import Search from '@/components/icons/Search.vue'
import Unlink from '@/components/icons/Unlink.vue'
import Status from '@/components/icons/Status.vue'
import SwitchButton from '@/components/common/SwitchButton.vue'
import ListView from '@/components/icons/ListView.vue'
import TreeView from '@/components/icons/TreeView.vue'

const themeVars = useThemeVars()
const dialogStore = useDialogStore()
const tabStore = useTabStore()
const currentName = computed(() => get(tabStore.currentTab, 'name', ''))
const browserTreeRef = ref(null)
/**
 *
 * @type {ComputedRef<{server: string, db: number, key: string}>}
 */
const currentSelect = computed(() => {
    const { server, db, key } = tabStore.currentTab || {}
    return { server, db, key }
})

const onInfo = () => {
    browserTreeRef.value?.handleSelectContextMenu('server_info')
}

const i18n = useI18n()
const connectionStore = useConnectionStore()
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

const viewType = ref(0)
const onSwitchView = (selectView) => {
    const { server } = tabStore.currentTab
    connectionStore.switchKeyView(server, selectView)
}
</script>

<template>
    <div class="nav-pane-container flex-box-v">
        <browser-tree ref="browserTreeRef" :server="currentName" />

        <div v-if="filterForm.showFilter" class="nav-pane-bottom flex-box-h">
            <n-input-group>
                <n-select
                    v-model:value="filterForm.type"
                    :consistent-menu-width="false"
                    :options="filterTypeOptions"
                    style="width: 120px" />
                <n-input clearable placeholder="" />
                <n-button :focusable="false" ghost>
                    <template #icon>
                        <n-icon :component="Search" />
                    </template>
                </n-button>
            </n-input-group>
        </div>
        <!-- bottom function bar -->
        <div class="nav-pane-bottom flex-box-h">
            <switch-button
                v-model:value="viewType"
                :icons="[TreeView, ListView]"
                :t-tooltips="['interface.tree_view', 'interface.list_view']"
                stroke-width="4"
                unselect-stroke-width="3"
                @update:value="onSwitchView" />
            <icon-button :icon="Status" size="20" stroke-width="4" t-tooltip="interface.status" @click="onInfo" />
            <icon-button :icon="Refresh" size="20" stroke-width="4" t-tooltip="interface.reload" @click="onRefresh" />
            <div class="flex-item-expand" />
            <icon-button
                :icon="Unlink"
                size="20"
                stroke-width="4"
                t-tooltip="interface.disconnect"
                @click="onDisconnect" />
        </div>
    </div>
</template>

<style lang="scss" scoped>
.nav-pane-bottom {
    color: v-bind('themeVars.iconColor');
    border-top: v-bind('themeVars.borderColor') 1px solid;
}
</style>
