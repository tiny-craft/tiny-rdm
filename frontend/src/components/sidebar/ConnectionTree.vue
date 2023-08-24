<script setup>
import useDialogStore from 'stores/dialog.js'
import { h, nextTick, reactive, ref } from 'vue'
import useConnectionStore from 'stores/connections.js'
import { NIcon, useDialog, useThemeVars } from 'naive-ui'
import { ConnectionType } from '@/consts/connection_type.js'
import ToggleFolder from '@/components/icons/ToggleFolder.vue'
import ToggleServer from '@/components/icons/ToggleServer.vue'
import { debounce, indexOf, isEmpty } from 'lodash'
import Config from '@/components/icons/Config.vue'
import Delete from '@/components/icons/Delete.vue'
import Unlink from '@/components/icons/Unlink.vue'
import CopyLink from '@/components/icons/CopyLink.vue'
import Connect from '@/components/icons/Connect.vue'
import { useI18n } from 'vue-i18n'
import useTabStore from 'stores/tab.js'
import Edit from '@/components/icons/Edit.vue'

const themeVars = useThemeVars()
const i18n = useI18n()
const openingConnection = ref(false)
const connectionStore = useConnectionStore()
const tabStore = useTabStore()
const dialogStore = useDialogStore()

const expandedKeys = ref([])
const selectedKeys = ref([])

const props = defineProps({
    filterPattern: {
        type: String,
    },
})

const contextMenuParam = reactive({
    show: false,
    x: 0,
    y: 0,
    options: null,
    currentNode: null,
})

const renderIcon = (icon) => {
    return () => {
        return h(NIcon, null, {
            default: () => h(icon),
        })
    }
}
const menuOptions = {
    [ConnectionType.Group]: ({ opened }) => [
        {
            key: 'group_rename',
            label: i18n.t('rename_conn_group'),
            icon: renderIcon(Edit),
        },
        {
            key: 'group_delete',
            label: i18n.t('remove_conn_group'),
            icon: renderIcon(Delete),
        },
    ],
    [ConnectionType.Server]: ({ name }) => {
        const connected = connectionStore.isConnected(name)
        if (connected) {
            return [
                {
                    key: 'server_close',
                    label: i18n.t('disconnect'),
                    icon: renderIcon(Unlink),
                },
                {
                    key: 'server_dup',
                    label: i18n.t('dup_conn'),
                    icon: renderIcon(CopyLink),
                },
                {
                    key: 'server_edit',
                    label: i18n.t('edit_conn'),
                    icon: renderIcon(Config),
                },
                {
                    type: 'divider',
                    key: 'd1',
                },
                {
                    key: 'server_remove',
                    label: i18n.t('remove_conn'),
                    icon: renderIcon(Delete),
                },
            ]
        } else {
            return [
                {
                    key: 'server_open',
                    label: i18n.t('open_connection'),
                    icon: renderIcon(Connect),
                },
                {
                    key: 'server_edit',
                    label: i18n.t('edit_conn'),
                    icon: renderIcon(Edit),
                },
                {
                    type: 'divider',
                    key: 'd1',
                },
                {
                    key: 'server_remove',
                    label: i18n.t('remove_conn'),
                    icon: renderIcon(Delete),
                },
            ]
        }
    },
}

const renderLabel = ({ option }) => {
    return option.label
}

const renderPrefix = ({ option }) => {
    switch (option.type) {
        case ConnectionType.Group:
            const opened = indexOf(expandedKeys.value, option.key) !== -1
            return h(
                NIcon,
                { size: 20 },
                {
                    default: () => h(ToggleFolder, { modelValue: opened }),
                },
            )
        case ConnectionType.Server:
            const connected = connectionStore.isConnected(option.name)
            return h(
                NIcon,
                { size: 20 },
                {
                    default: () => h(ToggleServer, { modelValue: !!connected }),
                },
            )
    }
}

const renderSuffix = ({ option }) => {
    if (option.type === ConnectionType.Server) {
        const { markColor = '' } = connectionStore.serverProfile[option.name] || {}
        if (isEmpty(markColor)) {
            return ''
        }
        return h('div', {
            style: {
                borderRadius: '50%',
                backgroundColor: markColor,
                width: '13px',
                height: '13px',
                border: '2px solid white',
            },
        })
    }
    return null
}

const onUpdateExpandedKeys = (keys, option) => {
    expandedKeys.value = keys
}

const onUpdateSelectedKeys = (keys, option) => {
    selectedKeys.value = keys
}

/**
 * Open connection
 * @param name
 * @returns {Promise<void>}
 */
const openConnection = async (name) => {
    try {
        if (!connectionStore.isConnected(name)) {
            openingConnection.value = true
            await connectionStore.openConnection(name)
        }
        tabStore.upsertTab({
            server: name,
        })
    } catch (e) {
        $message.error(e.message)
        // node.isLeaf = undefined
    } finally {
        openingConnection.value = false
    }
}

const dialog = useDialog()
const removeConnection = (name) => {
    $dialog.warning(i18n.t('remove_tip', { type: i18n.t('conn_name'), name }), async () => {
        connectionStore.deleteConnection(name).then(({ success, msg }) => {
            if (!success) {
                $message.error(msg)
            }
        })
    })
}

const removeGroup = async (name) => {
    $dialog.warning(i18n.t('remove_tip', { type: i18n.t('conn_group'), name }), async () => {
        connectionStore.deleteGroup(name).then(({ success, msg }) => {
            if (!success) {
                $message.error(msg)
            }
        })
    })
}

const expandKey = (key) => {
    const idx = indexOf(expandedKeys.value, key)
    if (idx === -1) {
        expandedKeys.value.push(key)
    } else {
        expandedKeys.value.splice(idx, 1)
    }
}

const nodeProps = ({ option }) => {
    return {
        onDblclick: async () => {
            if (option.type === ConnectionType.Server) {
                openConnection(option.name).then(() => {})
            } else if (option.type === ConnectionType.Group) {
                // toggle expand
                nextTick().then(() => expandKey(option.key))
            }
        },
        onContextmenu(e) {
            e.preventDefault()
            const mop = menuOptions[option.type]
            if (mop == null) {
                return
            }
            contextMenuParam.show = false
            nextTick().then(() => {
                contextMenuParam.options = mop(option)
                contextMenuParam.currentNode = option
                contextMenuParam.x = e.clientX
                contextMenuParam.y = e.clientY
                contextMenuParam.show = true
                selectedKeys.value = [option.key]
            })
        },
    }
}

const renderContextLabel = (option) => {
    return h('div', { class: 'context-menu-item' }, option.label)
}

const handleSelectContextMenu = (key) => {
    contextMenuParam.show = false
    const { name, label, db, key: nodeKey, redisKey } = contextMenuParam.currentNode
    switch (key) {
        case 'server_open':
            openConnection(name).then(() => {})
            break
        case 'server_edit':
            // ask for close relevant connections before edit
            if (connectionStore.isConnected(name)) {
                $dialog.warning(i18n.t('edit_close_confirm'), () => {
                    connectionStore.closeConnection(name)
                    dialogStore.openEditDialog(name)
                })
            } else {
                dialogStore.openEditDialog(name)
            }
            break
        case 'server_remove':
            removeConnection(name)
            break
        case 'server_close':
            connectionStore.closeConnection(name)
            break
        case 'group_rename':
            dialogStore.openRenameGroupDialog(label)
            break
        case 'group_delete':
            removeGroup(label)
            break
        default:
            console.warn('TODO: handle context menu:' + key)
    }
}

const findSiblingsAndIndex = (node, nodes) => {
    if (!nodes) {
        return [null, null]
    }
    for (let i = 0; i < nodes.length; ++i) {
        const siblingNode = nodes[i]
        if (siblingNode.key === node.key) {
            return [nodes, i]
        }
        const [siblings, index] = findSiblingsAndIndex(node, siblingNode.children)
        if (siblings && index !== null) {
            return [siblings, index]
        }
    }
    return [null, null]
}

// delay save until drop stopped after 2 seconds
const saveSort = debounce(connectionStore.saveConnectionSorted, 2000, { trailing: true })
const handleDrop = ({ node, dragNode, dropPosition }) => {
    const [dragNodeSiblings, dragNodeIndex] = findSiblingsAndIndex(dragNode, connectionStore.connections)
    if (dragNodeSiblings === null || dragNodeIndex === null) {
        return
    }
    dragNodeSiblings.splice(dragNodeIndex, 1)
    if (dropPosition === 'inside') {
        if (node.children) {
            node.children.unshift(dragNode)
        } else {
            node.children = [dragNode]
        }
    } else if (dropPosition === 'before') {
        const [nodeSiblings, nodeIndex] = findSiblingsAndIndex(node, connectionStore.connections)
        if (nodeSiblings === null || nodeIndex === null) {
            return
        }
        nodeSiblings.splice(nodeIndex, 0, dragNode)
    } else if (dropPosition === 'after') {
        const [nodeSiblings, nodeIndex] = findSiblingsAndIndex(node, connectionStore.connections)
        if (nodeSiblings === null || nodeIndex === null) {
            return
        }
        nodeSiblings.splice(nodeIndex + 1, 0, dragNode)
    }
    connectionStore.connections = Array.from(connectionStore.connections)
    saveSort()
}
</script>

<template>
    <n-empty v-if="isEmpty(connectionStore.connections)" :description="$t('empty_server_list')" class="empty-content" />
    <n-tree
        v-else
        :animated="false"
        :block-line="true"
        :block-node="true"
        :cancelable="false"
        :data="connectionStore.connections"
        :draggable="true"
        :expanded-keys="expandedKeys"
        :node-props="nodeProps"
        :pattern="props.filterPattern"
        :render-label="renderLabel"
        :render-prefix="renderPrefix"
        :render-suffix="renderSuffix"
        :selected-keys="selectedKeys"
        class="fill-height"
        virtual-scroll
        @drop="handleDrop"
        @update:selected-keys="onUpdateSelectedKeys"
        @update:expanded-keys="onUpdateExpandedKeys" />

    <!-- status display modal -->
    <n-modal :show="openingConnection" transform-origin="center">
        <n-card
            :bordered="false"
            :content-style="{ textAlign: 'center' }"
            aria-model="true"
            role="dialog"
            style="width: 400px">
            <n-spin>
                <template #description>
                    {{ $t('opening_connection') }}
                </template>
            </n-spin>
        </n-card>
    </n-modal>

    <!-- context menu -->
    <n-dropdown
        :animated="false"
        :options="contextMenuParam.options"
        :render-label="renderContextLabel"
        :show="contextMenuParam.show"
        :x="contextMenuParam.x"
        :y="contextMenuParam.y"
        placement="bottom-start"
        trigger="manual"
        @clickoutside="contextMenuParam.show = false"
        @select="handleSelectContextMenu" />
</template>

<style lang="scss" scoped>
@import '@/styles/content';
</style>
