<script setup>
import { NIcon, useMessage } from 'naive-ui'
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
import { useConfirmDialog } from '../../utils/confirm_dialog.js'
import { useI18n } from 'vue-i18n'
import useConnectionStore from '../../stores/connections.js'

const dialogStore = useDialogStore()
const tabStore = useTabStore()
const currentName = computed(() => get(tabStore.currentTab, 'name', ''))
/**
 *
 * @type {ComputedRef<{server: string, db: number, key: string}>}
 */
const currentSelect = computed(() => {
    const { server, db, key } = tabStore.currentTab || {}
    return { server, db, key }
})

const onNewKey = () => {
    const { server, db, key } = currentSelect.value
    dialogStore.openNewKeyDialog(key, server, db)
}

const i18n = useI18n()
const connectionStore = useConnectionStore()
const confirmDialog = useConfirmDialog()
const message = useMessage()
const onDeleteKey = () => {
    const { server, db, key } = currentSelect.value
    confirmDialog.warning(i18n.t('remove_tip', { name: key }), () => {
        connectionStore.deleteKey(server, db, key).then((success) => {
            if (success) {
                message.success(i18n.t('delete_key_succ', { key }))
            }
        })
    })
}

const onRefresh = () => {
    connectionStore.openConnection(currentSelect.value.server, true).then(() => {
        message.success(i18n.t('reload_succ'))
    })
}
</script>

<template>
    <div class="nav-pane-container flex-box-v">
        <browser-tree :server="currentName" />

        <!-- bottom function bar -->
        <div class="nav-pane-bottom flex-box-h">
            <icon-button
                :icon="AddLink"
                color="#555"
                size="20"
                stroke-width="4"
                t-tooltip="new_key"
                @click="onNewKey"
            />
            <icon-button
                :icon="Refresh"
                color="#555"
                size="20"
                stroke-width="4"
                t-tooltip="reload"
                @click="onRefresh"
            />
            <div class="flex-item-expand"></div>
            <icon-button
                :disabled="currentSelect.key == null"
                :icon="Delete"
                color="#555"
                size="20"
                stroke-width="4"
                t-tooltip="remove_key"
                @click="onDeleteKey"
            />
            <!--            <n-input placeholder="">-->
            <!--                <template #prefix>-->
            <!--                    <n-icon :component="Filter" color="#aaa" size="20" />-->
            <!--                </template>-->
            <!--            </n-input>-->
        </div>
    </div>
</template>

<style lang="scss" scoped></style>
