<script setup>
import useDialogStore from '../../stores/dialog.js'
import { NIcon } from 'naive-ui'
import AddGroup from '../icons/AddGroup.vue'
import AddLink from '../icons/AddLink.vue'
import IconButton from '../common/IconButton.vue'
import Filter from '../icons/Filter.vue'
import ConnectionTree from './ConnectionTree.vue'
import Unlink from '../icons/Unlink.vue'
import useConnectionStore from '../../stores/connections.js'
import { ref } from 'vue'

const dialogStore = useDialogStore()
const connectionStore = useConnectionStore()
const filterPattern = ref('')

const onDisconnectAll = () => {
    connectionStore.closeAllConnection()
}
</script>

<template>
    <div class="nav-pane-container flex-box-v">
        <connection-tree :filter-pattern="filterPattern" />

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
                @click="dialogStore.openNewGroupDialog()"
            />
            <icon-button
                :disabled="!connectionStore.anyConnectionOpened"
                :icon="Unlink"
                color="#555"
                size="20"
                stroke-width="4"
                t-tooltip="disconnect_all"
                @click="onDisconnectAll"
            />
            <n-input v-model:value="filterPattern" :placeholder="$t('filter')" clearable>
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
