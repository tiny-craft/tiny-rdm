<script setup>
import { useThemeVars } from 'naive-ui'
import BrowserTree from './BrowserTree.vue'
import IconButton from '@/components/common/IconButton.vue'
import useTabStore from 'stores/tab.js'
import { computed, onMounted, reactive, ref, unref, watch } from 'vue'
import { get, map } from 'lodash'
import Refresh from '@/components/icons/Refresh.vue'
import useDialogStore from 'stores/dialog.js'
import { useI18n } from 'vue-i18n'
import Search from '@/components/icons/Search.vue'
import Unlink from '@/components/icons/Unlink.vue'
import ContentSearchInput from '@/components/content_value/ContentSearchInput.vue'
import LoadAll from '@/components/icons/LoadAll.vue'
import LoadList from '@/components/icons/LoadList.vue'
import Delete from '@/components/icons/Delete.vue'
import useBrowserStore from 'stores/browser.js'
import { useRender } from '@/utils/render.js'
import RedisTypeSelector from '@/components/common/RedisTypeSelector.vue'
import { types } from '@/consts/support_redis_type.js'
import Plus from '@/components/icons/Plus.vue'

const themeVars = useThemeVars()
const i18n = useI18n()
const dialogStore = useDialogStore()
// const prefStore = usePreferencesStore()
const tabStore = useTabStore()
const browserStore = useBrowserStore()
const render = useRender()
const currentName = computed(() => get(tabStore.currentTab, 'name', ''))
const browserTreeRef = ref(null)
const loading = ref(false)
const fullyLoaded = ref(false)

const selectedDB = computed(() => {
    return browserStore.selectedDatabases[currentName.value] || 0
})

const dbSelectOptions = computed(() => {
    const dblist = browserStore.getDBList(currentName.value)
    return map(dblist, (db) => {
        if (selectedDB.value === db.db) {
            return {
                value: db.db,
                label: `db${db.db} (${db.keys}/${db.maxKeys})`,
            }
        }
        return {
            value: db.db,
            label: `db${db.db} (${db.maxKeys})`,
        }
    })
})

const loadProgress = computed(() => {
    const db = browserStore.getDatabase(currentName.value, selectedDB.value)
    if (db.maxKeys <= 0) {
        return 100
    }
    return (db.keys * 100) / db.maxKeys
})

const onReload = async () => {
    try {
        loading.value = true
        tabStore.setSelectedKeys(currentName.value)
        const db = selectedDB.value
        browserStore.closeDatabase(currentName.value, db)
        browserTreeRef.value?.resetExpandKey(currentName.value, db)

        let matchType = unref(filterForm.type)
        if (!types.hasOwnProperty(matchType)) {
            matchType = ''
        }
        browserStore.setKeyFilter(currentName.value, {
            type: matchType,
            pattern: unref(filterForm.pattern),
        })
        await browserStore.openDatabase(currentName.value, db)
        fullyLoaded.value = await browserStore.loadMoreKeys(currentName.value, db)
        // $message.success(i18n.t('dialogue.reload_succ'))
    } catch (e) {
        console.warn(e)
    } finally {
        loading.value = false
    }
}

const onAddKey = () => {
    dialogStore.openNewKeyDialog('', currentName.value, selectedDB.value)
}

const onLoadMore = async () => {
    try {
        loading.value = true
        fullyLoaded.value = await browserStore.loadMoreKeys(currentName.value, selectedDB.value)
    } catch (e) {
        $message.error(e.message)
    } finally {
        loading.value = false
    }
}

const onLoadAll = async () => {
    try {
        loading.value = true
        await browserStore.loadAllKeys(currentName.value, selectedDB.value)
        fullyLoaded.value = true
    } catch (e) {
        $message.error(e.message)
    } finally {
        loading.value = false
    }
}

const onFlush = () => {
    dialogStore.openFlushDBDialog(currentName.value, selectedDB.value)
}

const onDisconnect = () => {
    browserStore.closeConnection(currentName.value)
}

const handleSelectDB = async (db, prevDB) => {
    // watch 'browserStore.openedDB[currentName.value]' instead
}

const filterForm = reactive({
    type: '',
    pattern: '',
    filter: '',
})
const onSelectFilterType = (select) => {
    onReload()
}

const onFilterInput = (val) => {
    filterForm.filter = val
}

const onMatchInput = (matchVal, filterVal) => {
    filterForm.pattern = matchVal
    filterForm.filter = filterVal
    onReload()
}

watch(
    () => browserStore.openedDB[currentName.value],
    async (db, prevDB) => {
        if (db === undefined) {
            return
        }

        try {
            loading.value = true
            browserStore.closeDatabase(currentName.value, prevDB)
            browserStore.setKeyFilter(currentName.value, {})
            await browserStore.openDatabase(currentName.value, db)
            browserTreeRef.value?.resetExpandKey(currentName.value, db)
            fullyLoaded.value = await browserStore.loadMoreKeys(currentName.value, db)
            browserTreeRef.value?.refreshTree()
        } catch (e) {
            $message.error(e.message)
        } finally {
            loading.value = false
        }
    },
)

onMounted(() => onReload())
// forbid dynamic switch key view due to performance issues
// const viewType = ref(0)
// const onSwitchView = (selectView) => {
//     const { server } = tabStore.currentTab
//     browserStore.switchKeyView(server, selectView)
// }
</script>

<template>
    <div class="nav-pane-container flex-box-v">
        <!-- top function bar -->
        <div class="flex-box-h nav-pane-func">
            <content-search-input
                :debounce-wait="1000"
                :full-search-icon="Search"
                small
                @filter-changed="onFilterInput"
                @match-changed="onMatchInput">
                <template #prepend>
                    <redis-type-selector v-model:value="filterForm.type" @update:value="onSelectFilterType" />
                </template>
            </content-search-input>
            <n-button-group>
                <n-tooltip :show-arrow="false">
                    <template #trigger>
                        <n-button :disabled="loading" :focusable="false" bordered size="small" @click="onReload">
                            <template #icon>
                                <n-icon :component="Refresh" size="18" />
                            </template>
                        </n-button>
                    </template>
                    {{ $t('interface.reload') }}
                </n-tooltip>
                <n-tooltip :show-arrow="false">
                    <template #trigger>
                        <n-button :disabled="loading" :focusable="false" bordered size="small" @click="onAddKey">
                            <template #icon>
                                <n-icon :component="Plus" size="18" />
                            </template>
                        </n-button>
                    </template>
                    {{ $t('interface.new_key') }}
                </n-tooltip>
            </n-button-group>
        </div>

        <!-- loaded progress -->
        <n-progress
            :border-radius="0"
            :color="fullyLoaded ? '#0000' : themeVars.primaryColor"
            :height="2"
            :percentage="loadProgress"
            :processing="loading"
            :show-indicator="false"
            status="success"
            type="line" />

        <!-- tree view -->
        <browser-tree
            ref="browserTreeRef"
            :full-loaded="fullyLoaded"
            :loading="loading && loadProgress <= 0"
            :pattern="filterForm.filter"
            :server="currentName" />
        <!-- bottom function bar -->
        <div class="nav-pane-bottom flex-box-v">
            <!--            <switch-button-->
            <!--                v-model:value="viewType"-->
            <!--                :icons="[TreeView, ListView]"-->
            <!--                :t-tooltips="['interface.tree_view', 'interface.list_view']"-->
            <!--                :stroke-width="3.5"-->
            <!--                unselect-stroke-width="3"-->
            <!--                @update:value="onSwitchView" />-->
            <div class="flex-box-h nav-pane-func">
                <n-select
                    v-model:value="browserStore.openedDB[currentName]"
                    :consistent-menu-width="false"
                    :filter="(pattern, option) => option.value.toString() === pattern"
                    :options="dbSelectOptions"
                    filterable
                    size="small"
                    style="min-width: 100px; max-width: 200px"
                    @update:value="handleSelectDB" />
                <icon-button
                    :button-class="['nav-pane-func-btn']"
                    :disabled="fullyLoaded"
                    :icon="LoadList"
                    :loading="loading"
                    :stroke-width="3.5"
                    size="20"
                    t-tooltip="interface.load_more"
                    @click="onLoadMore" />
                <icon-button
                    :button-class="['nav-pane-func-btn']"
                    :disabled="fullyLoaded"
                    :icon="LoadAll"
                    :loading="loading"
                    :stroke-width="3.5"
                    size="20"
                    t-tooltip="interface.load_all"
                    @click="onLoadAll" />
                <div class="flex-item-expand" style="min-width: 10px" />
                <icon-button
                    :button-class="['nav-pane-func-btn']"
                    :icon="Delete"
                    :stroke-width="3.5"
                    size="20"
                    t-tooltip="interface.flush_db"
                    @click="onFlush" />
                <icon-button
                    :button-class="['nav-pane-func-btn']"
                    :icon="Unlink"
                    :stroke-width="3.5"
                    size="20"
                    t-tooltip="interface.disconnect"
                    @click="onDisconnect" />
            </div>
        </div>
    </div>
</template>

<style lang="scss" scoped>
@import '@/styles/style';

:deep(.toggle-btn) {
    border-style: solid;
    border-width: 1px;
}

:deep(.filter-on) {
    border-color: v-bind('themeVars.iconColorDisabled');
    background-color: v-bind('themeVars.iconColorDisabled');
}

:deep(.filter-off) {
    border-color: #0000;
}

.nav-pane-top {
    //@include bottom-shadow(0.1);
    color: v-bind('themeVars.iconColor');
    border-bottom: v-bind('themeVars.borderColor') 1px solid;
}

.nav-pane-bottom {
    @include top-shadow(0.1);
    color: v-bind('themeVars.iconColor');
    border-top: v-bind('themeVars.borderColor') 1px solid;
}
</style>
