<script setup>
import useDialogStore from 'stores/dialog.js'
import { NIcon, useThemeVars } from 'naive-ui'
import AddGroup from '@/components/icons/AddGroup.vue'
import AddLink from '@/components/icons/AddLink.vue'
import IconButton from '@/components/common/IconButton.vue'
import Filter from '@/components/icons/Filter.vue'
import ConnectionTree from './ConnectionTree.vue'
import { ref } from 'vue'

const themeVars = useThemeVars()
const dialogStore = useDialogStore()
const filterPattern = ref('')
</script>

<template>
    <div class="nav-pane-container flex-box-v">
        <connection-tree :filter-pattern="filterPattern" />

        <!-- bottom function bar -->
        <div class="nav-pane-bottom nav-pane-func flex-box-h">
            <icon-button
                :button-class="['nav-pane-func-btn']"
                :icon="AddLink"
                size="20"
                stroke-width="4"
                t-tooltip="interface.new_conn"
                @click="dialogStore.openNewDialog()" />
            <icon-button
                :button-class="['nav-pane-func-btn']"
                :icon="AddGroup"
                size="20"
                stroke-width="4"
                t-tooltip="interface.new_group"
                @click="dialogStore.openNewGroupDialog()" />
            <n-input
                v-model:value="filterPattern"
                :input-props="{ spellcheck: 'false' }"
                :placeholder="$t('interface.filter')"
                clearable>
                <template #prefix>
                    <n-icon :component="Filter" size="20" />
                </template>
            </n-input>
        </div>
    </div>
</template>

<style lang="scss" scoped>
.nav-pane-bottom {
    color: v-bind('themeVars.iconColor');
    border-top: v-bind('themeVars.borderColor') 1px solid;
}
</style>
