<script setup>
import { NIcon } from 'naive-ui'
import AddLink from '../icons/AddLink.vue'
import BrowserTree from './BrowserTree.vue'
import IconButton from '../common/IconButton.vue'
import Filter from '../icons/Filter.vue'
import useTabStore from '../../stores/tab.js'
import { computed } from 'vue'
import { get } from 'lodash'
import Delete from '../icons/Delete.vue'
import Refresh from '../icons/Refresh.vue'
import useDialogStore from '../../stores/dialog.js'

const dialogStore = useDialogStore()
const tabStore = useTabStore()
const currentName = computed(() => get(tabStore.currentTab, 'name', ''))
</script>

<template>
    <div class="nav-pane-container flex-box-v">
        <browser-tree :server="currentName" />

        <!-- bottom function bar -->
        <div class="nav-pane-bottom flex-box-h">
            <icon-button :icon="AddLink" color="#555" size="20" stroke-width="4" t-tooltip="new_key" />
            <icon-button :icon="Delete" color="#555" size="20" stroke-width="4" t-tooltip="remove_key" />
            <icon-button :icon="Refresh" color="#555" size="20" stroke-width="4" t-tooltip="reload" />
            <n-input placeholder="">
                <template #prefix>
                    <n-icon :component="Filter" color="#aaa" size="20" />
                </template>
            </n-input>
        </div>
    </div>
</template>

<style lang="scss" scoped>
.nav-pane-container {
    overflow: hidden;
    background-color: var(--bg-color);

    .nav-pane-bottom {
        align-items: center;
        gap: 5px;
        padding: 3px 3px 5px 5px;
    }
}
</style>
