<script setup>
import useDialogStore from 'stores/dialog.js'
import { h, markRaw, nextTick, reactive, ref } from 'vue'
import useConnectionStore from 'stores/connections.js'
import { NIcon, NSpace, NText, useThemeVars } from 'naive-ui'
import { ConnectionType } from '@/consts/connection_type.js'
import Folder from '@/components/icons/Folder.vue'
import Server from '@/components/icons/Server.vue'
import Cluster from '@/components/icons/Cluster.vue'
import { debounce, get, includes, indexOf, isEmpty, split } from 'lodash'
import Config from '@/components/icons/Config.vue'
import Delete from '@/components/icons/Delete.vue'
import Unlink from '@/components/icons/Unlink.vue'
import CopyLink from '@/components/icons/CopyLink.vue'
import Connect from '@/components/icons/Connect.vue'
import { useI18n } from 'vue-i18n'
import useTabStore from 'stores/tab.js'
import Edit from '@/components/icons/Edit.vue'
import { hexGammaCorrection, parseHexColor, toHexColor } from '@/utils/rgb.js'
import IconButton from '@/components/common/IconButton.vue'
import usePreferencesStore from 'stores/preferences.js'
import useBrowserStore from 'stores/browser.js'
import { useRender } from '@/utils/render.js'

const themeVars = useThemeVars()
const i18n = useI18n()
const render = useRender()
const connectingServer = ref('')
const connectionStore = useConnectionStore()
const browserStore = useBrowserStore()
const tabStore = useTabStore()
const prefStore = usePreferencesStore()
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

const menuOptions = {
    [ConnectionType.Group]: ({ opened }) => [
        {
            key: 'group_rename',
            label: 'interface.rename_conn_group',
            icon: Edit,
        },
        {
            key: 'group_delete',
            label: 'interface.remove_conn_group',
            icon: Delete,
        },
    ],
    [ConnectionType.Server]: ({ name }) => {
        const connected = browserStore.isConnected(name)
        if (connected) {
            return [
                {
                    key: 'server_close',
                    label: 'interface.disconnect',
                    icon: Unlink,
                },
                {
                    key: 'server_edit',
                    label: 'interface.edit_conn',
                    icon: Config,
                },
                {
                    key: 'server_dup',
                    label: 'interface.dup_conn',
                    icon: CopyLink,
                },
                {
                    type: 'divider',
                    key: 'd1',
                },
                {
                    key: 'server_remove',
                    label: 'interface.remove_conn',
                    icon: Delete,
                },
            ]
        } else {
            return [
                {
                    key: 'server_open',
                    label: 'interface.open_connection',
                    icon: Connect,
                },
                {
                    key: 'server_edit',
                    label: 'interface.edit_conn',
                    icon: Config,
                },
                {
                    key: 'server_dup',
                    label: 'interface.dup_conn',
                    icon: CopyLink,
                },
                {
                    type: 'divider',
                    key: 'd1',
                },
                {
                    key: 'server_remove',
                    label: 'interface.remove_conn',
                    icon: Delete,
                },
            ]
        }
    },
}

/**
 * get mark color of server saved in preferences
 * @param name
 * @return {null|string}
 */
const getServerMarkColor = (name) => {
    const { markColor = '' } = connectionStore.serverProfile[name] || {}
    if (!isEmpty(markColor)) {
        const rgb = parseHexColor(markColor)
        const rgb2 = hexGammaCorrection(rgb, 0.75)
        return toHexColor(rgb2)
    }
    return null
}

const renderLabel = ({ option }) => {
    if (option.type === ConnectionType.Server) {
        const color = getServerMarkColor(option.name)
        if (color != null) {
            return h(
                NText,
                {
                    style: {
                        color,
                        fontWeight: '450',
                    },
                },
                () => option.label,
            )
        }
    }
    return option.label
}

// render horizontal item
const renderIconMenu = (items) => {
    return h(
        NSpace,
        {
            align: 'center',
            inline: true,
            size: 3,
            wrapItem: false,
            wrap: false,
            style: 'margin-right: 5px',
        },
        () => items,
    )
}

const renderPrefix = ({ option }) => {
    const iconTransparency = prefStore.isDark ? 0.75 : 1
    switch (option.type) {
        case ConnectionType.Group:
            const opened = indexOf(expandedKeys.value, option.key) !== -1
            return h(
                NIcon,
                { size: 20 },
                {
                    default: () =>
                        h(Folder, {
                            open: opened,
                            fillColor: `rgba(255,206,120,${iconTransparency})`,
                        }),
                },
            )
        case ConnectionType.Server:
            const connected = browserStore.isConnected(option.name)
            const color = getServerMarkColor(option.name)
            const icon = option.cluster === true ? Cluster : Server
            return h(
                NIcon,
                { size: 20, color: !!!connected ? color : '#dc423c' },
                {
                    default: () =>
                        h(icon, {
                            inverse: !!connected,
                            fillColor: `rgba(220,66,60,${iconTransparency})`,
                        }),
                },
            )
    }
}

const getServerMenu = (connected) => {
    const btns = []
    if (connected) {
        btns.push(
            h(IconButton, {
                tTooltip: 'interface.disconnect',
                icon: Unlink,
                onClick: () => handleSelectContextMenu('server_close'),
            }),
            h(IconButton, {
                tTooltip: 'interface.edit_conn',
                icon: Config,
                onClick: () => handleSelectContextMenu('server_edit'),
            }),
        )
    } else {
        btns.push(
            h(IconButton, {
                tTooltip: 'interface.open_connection',
                icon: Connect,
                onClick: () => handleSelectContextMenu('server_open'),
            }),
            h(IconButton, {
                tTooltip: 'interface.edit_conn',
                icon: Config,
                onClick: () => handleSelectContextMenu('server_edit'),
            }),
            h(IconButton, {
                tTooltip: 'interface.remove_conn',
                icon: Delete,
                onClick: () => handleSelectContextMenu('server_remove'),
            }),
        )
    }
    return btns
}

const getGroupMenu = () => {
    return [
        h(IconButton, {
            tTooltip: 'interface.rename_conn_group',
            icon: Config,
            onClick: () => handleSelectContextMenu('group_rename'),
        }),
        h(IconButton, {
            tTooltip: 'interface.remove_conn_group',
            icon: Delete,
            onClick: () => handleSelectContextMenu('group_delete'),
        }),
    ]
}

const renderSuffix = ({ option }) => {
    if (includes(selectedKeys.value, option.key)) {
        switch (option.type) {
            case ConnectionType.Server:
                const connected = browserStore.isConnected(option.name)
                return renderIconMenu(getServerMenu(connected))
            case ConnectionType.Group:
                return renderIconMenu(getGroupMenu())
        }
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
        connectingServer.value = name
        if (!browserStore.isConnected(name)) {
            await browserStore.openConnection(name)
        }
        // check if connection already canceled before finish open
        if (!isEmpty(connectingServer.value)) {
            tabStore.upsertTab({
                server: name,
                db: browserStore.getSelectedDB(name),
            })
        }
    } catch (e) {
        $message.error(e.message)
        // node.isLeaf = undefined
    } finally {
        connectingServer.value = ''
    }
}

const removeConnection = (name) => {
    $dialog.warning(
        i18n.t('dialogue.remove_tip', { type: i18n.t('dialogue.connection.conn_name'), name }),
        async () => {
            connectionStore.deleteConnection(name).then(({ success, msg }) => {
                if (!success) {
                    $message.error(msg)
                }
            })
        },
    )
}

const removeGroup = async (name) => {
    $dialog.warning(i18n.t('dialogue.remove_group_tip', { name }), async () => {
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
                contextMenuParam.options = markRaw(mop(option))
                contextMenuParam.currentNode = option
                contextMenuParam.x = e.clientX
                contextMenuParam.y = e.clientY
                contextMenuParam.show = true
                selectedKeys.value = [option.key]
            })
        },
    }
}

const handleSelectContextMenu = (key) => {
    contextMenuParam.show = false
    const selectedKey = get(selectedKeys.value, 0)
    if (selectedKey == null) {
        return
    }
    const [group, name] = split(selectedKey, '/')
    if (isEmpty(group) && isEmpty(name)) {
        return
    }
    switch (key) {
        case 'server_open':
            openConnection(name).then(() => {})
            break
        case 'server_edit':
            // ask for close relevant connections before edit
            if (browserStore.isConnected(name)) {
                $dialog.warning(i18n.t('dialogue.edit_close_confirm'), () => {
                    browserStore.closeConnection(name)
                    dialogStore.openEditDialog(name)
                })
            } else {
                dialogStore.openEditDialog(name)
            }
            break
        case 'server_dup':
            dialogStore.openDuplicateDialog(name)
            break
        case 'server_remove':
            removeConnection(name)
            break
        case 'server_close':
            browserStore.closeConnection(name).then((closed) => {
                if (closed) {
                    $message.success(i18n.t('dialogue.handle_succ'))
                }
            })
            break
        case 'group_rename':
            if (!isEmpty(group)) {
                dialogStore.openRenameGroupDialog(group)
            }
            break
        case 'group_delete':
            if (!isEmpty(group)) {
                removeGroup(group)
            }
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
const saveSort = debounce(connectionStore.saveConnectionSorted, 1500, { trailing: true })
const handleDrop = ({ node, dragNode, dropPosition }) => {
    const [dragNodeSiblings, dragNodeIndex] = findSiblingsAndIndex(dragNode, connectionStore.connections)
    if (dragNodeSiblings === null || dragNodeIndex === null) {
        return
    }
    if (node.type === ConnectionType.Group && dragNode.type === ConnectionType.Group) {
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

const onCancelOpen = () => {
    if (!isEmpty(connectingServer.value)) {
        browserStore.closeConnection(connectingServer.value)
        connectingServer.value = ''
    }
}
</script>

<template>
    <div class="connection-tree-wrapper" @keydown.esc="contextMenuParam.show = false">
        <n-tree
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
            @update:expanded-keys="onUpdateExpandedKeys">
            <template #empty>
                <n-empty :description="$t('interface.empty_server_list')" class="empty-content" />
            </template>
        </n-tree>

        <!-- status display modal -->
        <n-modal :show="connectingServer !== ''" transform-origin="center">
            <n-card
                :bordered="false"
                :content-style="{ textAlign: 'center' }"
                aria-model="true"
                role="dialog"
                style="width: 400px">
                <n-spin>
                    <template #description>
                        <n-space vertical>
                            <n-text strong>{{ $t('dialogue.opening_connection') }}</n-text>
                            <n-button :focusable="false" secondary size="small" @click="onCancelOpen">
                                {{ $t('dialogue.interrupt_connection') }}
                            </n-button>
                        </n-space>
                    </template>
                </n-spin>
            </n-card>
        </n-modal>

        <!-- context menu -->
        <n-dropdown
            :keyboard="true"
            :options="contextMenuParam.options"
            :render-icon="({ icon }) => render.renderIcon(icon)"
            :render-label="({ label }) => render.renderLabel($t(label), { class: 'context-menu-item' })"
            :show="contextMenuParam.show"
            :x="contextMenuParam.x"
            :y="contextMenuParam.y"
            placement="bottom-start"
            trigger="manual"
            @clickoutside="contextMenuParam.show = false"
            @select="handleSelectContextMenu" />
    </div>
</template>

<style lang="scss" scoped>
@use '@/styles/content';

.connection-tree-wrapper {
    height: 100%;
    overflow: hidden;
}
</style>
