<script setup>
import { useThemeVars } from 'naive-ui'
import BrowserTree from './BrowserTree.vue'
import IconButton from '@/components/common/IconButton.vue'
import useTabStore from 'stores/tab.js'
import { computed, nextTick, onMounted, reactive, ref, unref } from 'vue'
import { find, map } from 'lodash'
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
import useConnectionStore from 'stores/connections.js'
import ListCheckbox from '@/components/icons/ListCheckbox.vue'
import Close from '@/components/icons/Close.vue'
import More from '@/components/icons/More.vue'

const props = defineProps({
    server: String,
    db: {
        type: Number,
        default: 0,
    },
})

const themeVars = useThemeVars()
const i18n = useI18n()
const dialogStore = useDialogStore()
const tabStore = useTabStore()
const browserStore = useBrowserStore()
const connectionStore = useConnectionStore()
const render = useRender()
const browserTreeRef = ref(null)
const loading = ref(false)
const fullyLoaded = ref(false)
const inCheckState = ref(false)
const checkedCount = ref(0)

const dbSelectOptions = computed(() => {
    const dblist = browserStore.getDBList(props.server)
    return map(dblist, (db) => {
        if (props.db === db.db) {
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

const moreOptions = computed(() => {
    return [
        { key: 'flush', label: i18n.t('interface.flush_db'), icon: render.renderIcon(Delete, { strokeWidth: 3.5 }) },
        {
            key: 'disconnect',
            label: i18n.t('interface.disconnect'),
            icon: render.renderIcon(Unlink, { strokeWidth: 3.5 }),
        },
    ]
})

const loadProgress = computed(() => {
    const db = browserStore.getDatabase(props.server, props.db)
    if (db.maxKeys <= 0) {
        return 100
    }
    return (db.keys * 100) / Math.max(db.keys, db.maxKeys)
})

const checkedTip = computed(() => {
    const dblist = browserStore.getDBList(props.server)
    const db = find(dblist, { db: props.db })
    return `${checkedCount.value} / ${db.maxKeys}`
})

const onReload = async () => {
    try {
        loading.value = true
        tabStore.setSelectedKeys(props.server)
        const db = props.db
        browserStore.closeDatabase(props.server, db)
        browserTreeRef.value?.resetExpandKey(props.server, db)

        let matchType = unref(filterForm.type)
        if (!types.hasOwnProperty(matchType)) {
            matchType = ''
        }
        browserStore.setKeyFilter(props.server, {
            type: matchType,
            pattern: unref(filterForm.pattern),
        })
        await browserStore.openDatabase(props.server, db)
        fullyLoaded.value = await browserStore.loadMoreKeys(props.server, db)
        // $message.success(i18n.t('dialogue.reload_succ'))
    } catch (e) {
        console.warn(e)
    } finally {
        loading.value = false
    }
}

const onAddKey = () => {
    dialogStore.openNewKeyDialog('', props.server, props.db)
}

const onLoadMore = async () => {
    try {
        loading.value = true
        fullyLoaded.value = await browserStore.loadMoreKeys(props.server, props.db)
    } catch (e) {
        $message.error(e.message)
    } finally {
        loading.value = false
    }
}

const onLoadAll = async () => {
    try {
        loading.value = true
        await browserStore.loadAllKeys(props.server, props.db)
        fullyLoaded.value = true
    } catch (e) {
        $message.error(e.message)
    } finally {
        loading.value = false
    }
}

const onDeleteChecked = () => {
    browserTreeRef.value?.deleteCheckedItems()
}

const onFlush = () => {
    dialogStore.openFlushDBDialog(props.server, props.db)
}

const onDisconnect = () => {
    browserStore.closeConnection(props.server)
}

const handleSelectDB = async (db) => {
    if (db === props.db) {
        return
    }

    try {
        loading.value = true
        browserStore.setKeyFilter(props.server, {})
        browserStore.closeDatabase(props.server, props.db)
        await browserStore.openDatabase(props.server, db)
        await nextTick()
        await connectionStore.saveLastDB(props.server, db)
        tabStore.upsertTab({ server: props.server, db })
        // browserTreeRef.value?.resetExpandKey(props.server, db)
        fullyLoaded.value = await browserStore.loadMoreKeys(props.server, db)
        browserTreeRef.value?.refreshTree()
        tabStore.setSelectedKeys(props.server)
    } catch (e) {
        $message.error(e.message)
    } finally {
        loading.value = false
    }
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

const onSelectOptions = (select) => {
    switch (select) {
        case 'flush':
            onFlush()
            break
        case 'disconnect':
            onDisconnect()
            break
    }
}

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
        <div class="flex-box-h nav-pane-func" style="height: 36px">
            <content-search-input
                :debounce-wait="1000"
                :full-search-icon="Search"
                small
                use-glob
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
            v-model:checked-count="checkedCount"
            :check-mode="inCheckState"
            :db="props.db"
            :full-loaded="fullyLoaded"
            :loading="loading && loadProgress <= 0"
            :pattern="filterForm.filter"
            :server="props.server" />
        <!-- bottom function bar -->
        <div class="nav-pane-bottom flex-box-v">
            <!--            <switch-button-->
            <!--                v-model:value="viewType"-->
            <!--                :icons="[TreeView, ListView]"-->
            <!--                :t-tooltips="['interface.tree_view', 'interface.list_view']"-->
            <!--                :stroke-width="3.5"-->
            <!--                unselect-stroke-width="3"-->
            <!--                @update:value="onSwitchView" />-->
            <transition mode="out-in" name="fade">
                <div v-if="!inCheckState" class="flex-box-h nav-pane-func">
                    <n-select
                        :consistent-menu-width="false"
                        :filter="(pattern, option) => option.value.toString() === pattern"
                        :input-props="{ spellcheck: 'false' }"
                        :options="dbSelectOptions"
                        :value="props.db"
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
                        :icon="ListCheckbox"
                        :stroke-width="3.5"
                        size="20"
                        t-tooltip="interface.check_mode"
                        @click="inCheckState = true" />
                    <n-dropdown :options="moreOptions" placement="top-end" @select="onSelectOptions">
                        <icon-button :button-class="['nav-pane-func-btn']" :icon="More" :stroke-width="3.5" size="20" />
                    </n-dropdown>
                </div>

                <!-- check mode function bar -->
                <div v-else class="flex-box-h nav-pane-func">
                    <icon-button
                        :button-class="['nav-pane-func-btn']"
                        :disabled="checkedCount <= 0"
                        :icon="Delete"
                        :stroke-width="3.5"
                        size="20"
                        t-tooltip="interface.delete_checked"
                        @click="onDeleteChecked" />
                    <div class="flex-item-expand ellipsis" style="text-align: center; margin: 0 5px">
                        <n-text>{{ checkedTip }}</n-text>
                    </div>
                    <icon-button
                        :button-class="['nav-pane-func-btn']"
                        :icon="Close"
                        :stroke-width="3.5"
                        size="20"
                        t-tooltip="interface.quit_check_mode"
                        @click="inCheckState = false" />
                </div>
            </transition>
        </div>
    </div>
</template>

<style lang="scss" scoped>
@import '@/styles/style';

:deep(.toggle-btn) {
    border-style: solid;
    border-width: 1px;
    border-radius: 3px;
    padding: 4px;
}

:deep(.toggle-on) {
    border-color: v-bind('themeVars.iconColorDisabled');
    background-color: v-bind('themeVars.iconColorDisabled');
}

:deep(.toggle-off) {
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
