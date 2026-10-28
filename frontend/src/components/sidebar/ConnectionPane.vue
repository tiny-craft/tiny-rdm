<script setup>
import useDialogStore from 'stores/dialog.js'
import { NIcon, useThemeVars } from 'naive-ui'
import AddGroup from '@/components/icons/AddGroup.vue'
import AddLink from '@/components/icons/AddLink.vue'
import IconButton from '@/components/common/IconButton.vue'
import Filter from '@/components/icons/Filter.vue'
import ConnectionTree from './ConnectionTree.vue'
import { computed, ref } from 'vue'
import More from '@/components/icons/More.vue'
import Import from '@/components/icons/Import.vue'
import { useRender } from '@/utils/render.js'
import Export from '@/components/icons/Export.vue'
import Sort from '@/components/icons/Sort.vue'
import Checked from '@/components/icons/Checked.vue'
import useConnectionStore from 'stores/connections.js'
import usePreferencesStore from 'stores/preferences.js'

const themeVars = useThemeVars()
const dialogStore = useDialogStore()
const connectionStore = useConnectionStore()
const prefStore = usePreferencesStore()
const render = useRender()
const filterPattern = ref('')
const sorting = ref(false)

const moreOptions = computed(() => [
    {
        key: 'sort',
        label: 'interface.sort_conn',
        icon: Sort,
        children: [
            { key: 'original', label: 'interface.sort_conn_original', fallback: '保持原有' },
            { key: 'asc', label: 'A → Z', literal: true },
            { key: 'desc', label: 'Z → A', literal: true },
        ].map((option) => ({
            ...option,
            icon: connectionStore.connectionSortOrder === option.key ? Checked : undefined,
        })),
    },
    { key: 'import', label: 'interface.import_conn', icon: Import },
    { key: 'export', label: 'interface.export_conn', icon: Export },
])

const onSelectOptions = async (select) => {
    if (sorting.value) {
        return
    }
    switch (select) {
        case 'original':
        case 'asc':
        case 'desc':
            sorting.value = true
            try {
                const { success, msg } = await prefStore.setConnectionSortOrder(select)
                if (!success) {
                    $message.error(msg)
                }
            } catch (e) {
                $message.error(e.message)
            } finally {
                sorting.value = false
            }
            break
        case 'import':
            await connectionStore.importConnections()
            await connectionStore.initConnections(true)
            break
        case 'export':
            await connectionStore.exportConnections()
            break
    }
}
</script>

<template>
    <div class="nav-pane-container flex-box-v">
        <connection-tree :filter-pattern="filterPattern" />

        <!-- bottom function bar -->
        <div class="nav-pane-bottom nav-pane-func flex-box-h">
            <icon-button
                :button-class="['nav-pane-func-btn']"
                :icon="AddLink"
                :stroke-width="3.5"
                size="20"
                t-tooltip="interface.new_conn"
                @click="dialogStore.openNewDialog()" />
            <icon-button
                :button-class="['nav-pane-func-btn']"
                :icon="AddGroup"
                :stroke-width="3.5"
                size="20"
                t-tooltip="interface.new_group"
                @click="dialogStore.openNewGroupDialog()" />
            <n-divider vertical />
            <n-input v-model:value="filterPattern" :autofocus="false" :placeholder="$t('interface.filter')" clearable>
                <template #prefix>
                    <n-icon :component="Filter" size="20" />
                </template>
            </n-input>
            <n-dropdown
                :disabled="sorting"
                :options="moreOptions"
                :render-icon="({ icon }) => render.renderIcon(icon, { strokeWidth: 3.5 })"
                :render-label="
                    ({ label, literal, fallback }) => (literal ? label : fallback ? $t(label, fallback) : $t(label))
                "
                placement="top-end"
                style="min-width: 130px"
                trigger="click"
                @select="onSelectOptions">
                <icon-button :button-class="['nav-pane-func-btn']" :icon="More" :stroke-width="3.5" size="20" />
            </n-dropdown>
        </div>
    </div>
</template>

<style lang="scss" scoped>
.nav-pane-bottom {
    color: v-bind('themeVars.iconColor');
    border-top: v-bind('themeVars.borderColor') 1px solid;
}
</style>
