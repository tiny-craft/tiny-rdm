<script setup>
import useDialogStore from '../stores/dialog.js'
import { NIcon } from 'naive-ui'
import AddGroup from './icons/AddGroup.vue'
import AddLink from './icons/AddLink.vue'
import Sort from './icons/Sort.vue'
import ConnectionsTree from './ConnectionsTree.vue'
import IconButton from './IconButton.vue'
import Filter from './icons/Filter.vue'

const dialogStore = useDialogStore()

const onSort = () => {
    dialogStore.openPreferencesDialog()
}
</script>

<template>
    <div class="nav-pane-container flex-box-v">
        <ConnectionsTree />

        <!-- bottom function bar -->
        <div class="nav-pane-bottom flex-box-h">
            <icon-button
                :icon="AddLink"
                color="#555"
                size="20"
                stroke-width="4"
                t-tooltip="new_conn"
                @click="dialogStore.openNewDialog()"
            />
            <icon-button
                :icon="AddGroup"
                color="#555"
                size="20"
                stroke-width="4"
                t-tooltip="new_group"
                @click="dialogStore.openNewKeyDialog('aa:bb')"
            />
            <n-divider style="margin: 0 4px; --n-color: #aaa; width: 2px" vertical />
            <icon-button :icon="Sort" color="#555" size="20" stroke-width="4" t-tooltip="sort_conn" @click="onSort" />
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
