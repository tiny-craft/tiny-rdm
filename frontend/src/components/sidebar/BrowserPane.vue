<script setup>
import { useMessage, useThemeVars } from 'naive-ui'
import AddLink from '@/components/icons/AddLink.vue'
import BrowserTree from './BrowserTree.vue'
import IconButton from '@/components/common/IconButton.vue'
import useTabStore from 'stores/tab.js'
import { computed, reactive, ref } from 'vue'
import { get } from 'lodash'
import Delete from '@/components/icons/Delete.vue'
import Refresh from '@/components/icons/Refresh.vue'
import useDialogStore from 'stores/dialog.js'
import { useConfirmDialog } from '@/utils/confirm_dialog.js'
import { useI18n } from 'vue-i18n'
import useConnectionStore from 'stores/connections.js'
import { types } from '@/consts/support_redis_type.js'
import Search from '@/components/icons/Search.vue'

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

const onNewKey = () => {
    const { server, key, db = 0 } = currentSelect.value
    dialogStore.openNewKeyDialog(key, server, db)
}

const i18n = useI18n()
const connectionStore = useConnectionStore()
const message = useMessage()
const onDeleteKey = () => {
    browserTreeRef.value?.handleSelectContextMenu('value_remove')
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
        label: i18n.t('all'),
    })
    return options
})
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
                    style="width: 120px"
                />
                <n-input clearable placeholder="" />
                <n-button ghost>
                    <template #icon>
                        <n-icon :component="Search" />
                    </template>
                </n-button>
            </n-input-group>
        </div>
        <!-- bottom function bar -->
        <div class="nav-pane-bottom flex-box-h">
            <icon-button :icon="AddLink" size="20" stroke-width="4" t-tooltip="new_key" @click="onNewKey" />
            <icon-button :icon="Refresh" size="20" stroke-width="4" t-tooltip="reload" @click="onRefresh" />
            <!--            <icon-button-->
            <!--                :icon="Filter"-->
            <!--                size="20"-->
            <!--                stroke-width="4"-->
            <!--                t-tooltip="filter_key"-->
            <!--                @click="filterForm.showFilter = !filterForm.showFilter"-->
            <!--            />-->
            <div class="flex-item-expand"></div>
            <icon-button
                :disabled="currentSelect.key == null"
                :icon="Delete"
                size="20"
                stroke-width="4"
                t-tooltip="remove_key"
                @click="onDeleteKey"
            />
        </div>
    </div>
</template>

<style lang="scss" scoped>
.nav-pane-bottom {
    color: v-bind('themeVars.iconColor');
    border-top: v-bind('themeVars.borderColor') 1px solid;
}
</style>
