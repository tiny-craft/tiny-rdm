<script setup>
import { h, nextTick, onMounted, reactive, ref } from 'vue'
import { ConnectionType } from '../consts/connection_type.js'
import useConnection from '../stores/connection.js'
import { NIcon, useDialog, useMessage } from 'naive-ui'
import ToggleFolder from './icons/ToggleFolder.vue'
import Key from './icons/Key.vue'
import ToggleDb from './icons/ToggleDb.vue'
import ToggleServer from './icons/ToggleServer.vue'
import { indexOf, remove, startsWith } from 'lodash'
import { useI18n } from 'vue-i18n'
import Refresh from './icons/Refresh.vue'
import Config from './icons/Config.vue'
import CopyLink from './icons/CopyLink.vue'
import Unlink from './icons/Unlink.vue'
import Add from './icons/Add.vue'
import Layer from './icons/Layer.vue'
import Delete from './icons/Delete.vue'
import Connect from './icons/Connect.vue'
import useDialogStore from '../stores/dialog.js'
import { ClipboardSetText } from '../../wailsjs/runtime/runtime.js'
import useTabStore from '../stores/tab.js'

const i18n = useI18n()
const loading = ref(false)
const loadingConnections = ref(false)
const expandedKeys = ref([])
const connectionStore = useConnection()
const tabStore = useTabStore()
const dialogStore = useDialogStore()

const showContextMenu = ref(false)
const contextPos = reactive({ x: 0, y: 0 })
const contextMenuOptions = ref(null)
const currentContextNode = ref(null)
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
            key: 'group_reload',
            label: i18n.t('config_conn_group'),
            icon: renderIcon(Config),
        },
        {
            key: 'group_delete',
            label: i18n.t('remove_conn_group'),
            icon: renderIcon(Delete),
        },
    ],
    [ConnectionType.Server]: ({ connected }) => {
        if (connected) {
            return [
                {
                    key: 'server_reload',
                    label: i18n.t('reload'),
                    icon: renderIcon(Refresh),
                },
                {
                    key: 'server_disconnect',
                    label: i18n.t('disconnect'),
                    icon: renderIcon(Unlink),
                },
                {
                    key: 'server_dup',
                    label: i18n.t('dup_conn'),
                    icon: renderIcon(CopyLink),
                },
                {
                    key: 'server_config',
                    label: i18n.t('config_conn'),
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
            ]
        }
    },
    [ConnectionType.RedisDB]: ({ opened }) => {
        if (opened) {
            return [
                {
                    key: 'db_reload',
                    label: i18n.t('reload'),
                    icon: renderIcon(Refresh),
                },
                {
                    key: 'db_newkey',
                    label: i18n.t('new_key'),
                    icon: renderIcon(Add),
                },
            ]
        } else {
            return [
                {
                    key: 'db_open',
                    label: i18n.t('open_db'),
                    icon: renderIcon(Connect),
                },
            ]
        }
    },
    [ConnectionType.RedisKey]: () => [
        {
            key: 'key_reload',
            label: i18n.t('reload'),
            icon: renderIcon(Refresh),
        },
        {
            key: 'key_newkey',
            label: i18n.t('new_key'),
            icon: renderIcon(Add),
        },
        {
            key: 'key_copy',
            label: i18n.t('copy_path'),
            icon: renderIcon(CopyLink),
        },
        {
            type: 'divider',
            key: 'd1',
        },
        {
            key: 'key_remove',
            label: i18n.t('remove_path'),
            icon: renderIcon(Delete),
        },
    ],
    [ConnectionType.RedisValue]: () => [
        {
            key: 'value_reload',
            label: i18n.t('reload'),
            icon: renderIcon(Refresh),
        },
        {
            key: 'value_copy',
            label: i18n.t('copy_key'),
            icon: renderIcon(CopyLink),
        },
        {
            type: 'divider',
            key: 'd1',
        },
        {
            key: 'value_remove',
            label: i18n.t('remove_key'),
            icon: renderIcon(Delete),
        },
    ],
}

const renderContextLabel = (option) => {
    return h('div', { class: 'context-menu-item' }, option.label)
}

onMounted(async () => {
    try {
        // TODO: Show loading list status
        loadingConnections.value = true
        nextTick(connectionStore.initConnection)
    } finally {
        loadingConnections.value = false
    }
})

const expandKey = (key) => {
    const idx = indexOf(expandedKeys.value, key)
    if (idx === -1) {
        expandedKeys.value.push(key)
    } else {
        expandedKeys.value.splice(idx, 1)
    }
}

const collapseKeyAndChildren = (key) => {
    remove(expandedKeys.value, (k) => startsWith(k, key))
    // console.log(key)
    // const idx = indexOf(expandedKeys.value, key)
    // console.log(JSON.stringify(expandedKeys.value))
    // if (idx !== -1) {
    //     expandedKeys.value.splice(idx, 1)
    //     return true
    // }
    // return false
}

const message = useMessage()
const dialog = useDialog()
const onUpdateExpanded = (value, option, meta) => {
    expandedKeys.value = value
    if (!meta.node) {
        return
    }
    // console.log(JSON.stringify(meta))
    switch (meta.action) {
        case 'expand':
            meta.node.expanded = true
            break
        case 'collapse':
            meta.node.expanded = false
            break
    }
}

const renderPrefix = ({ option }) => {
    switch (option.type) {
        case ConnectionType.Group:
            return h(
                NIcon,
                { size: 20 },
                {
                    default: () => h(ToggleFolder, { modelValue: option.expanded === true }),
                }
            )
        case ConnectionType.Server:
            return h(
                NIcon,
                { size: 20 },
                {
                    default: () => h(ToggleServer, { modelValue: option.connected === true }),
                }
            )
        case ConnectionType.RedisDB:
            return h(
                NIcon,
                { size: 20 },
                {
                    default: () => h(ToggleDb, { modelValue: option.opened === true }),
                }
            )
        case ConnectionType.RedisKey:
            return h(
                NIcon,
                { size: 20 },
                {
                    default: () => h(Layer),
                }
            )
        case ConnectionType.RedisValue:
            return h(
                NIcon,
                { size: 20 },
                {
                    default: () => h(Key),
                }
            )
    }
}

const renderLabel = ({ option }) => {
    switch (option.type) {
        case ConnectionType.RedisDB:
        case ConnectionType.RedisKey:
            return `${option.label} (${option.keys || 0})`
        // case ConnectionType.RedisValue:
        //   return `[${option.keyType}]${option.label}`
    }
    return option.label
}

const renderSuffix = ({ option }) => {
    // return h(NButton,
    //     { text: true, type: 'primary' },
    //     { default: () => h(Key) })
}

const nodeProps = ({ option }) => {
    return {
        onClick() {
            connectionStore.select(option)
            // console.log('[click]:' + JSON.stringify(option))
        },
        onDblclick: async () => {
            if (loading.value) {
                console.warn('TODO: alert to ignore double click when loading')
                return
            }
            switch (option.type) {
                case ConnectionType.Server:
                    option.isLeaf = false
                    break

                case ConnectionType.RedisDB:
                    option.isLeaf = false
            }

            // default handle is expand current node
            nextTick().then(() => expandKey(option.key))
        },
        onContextmenu(e) {
            e.preventDefault()
            const mop = menuOptions[option.type]
            if (mop == null) {
                return
            }
            showContextMenu.value = false
            nextTick().then(() => {
                contextMenuOptions.value = mop(option)
                currentContextNode.value = option
                contextPos.x = e.clientX
                contextPos.y = e.clientY
                showContextMenu.value = true
            })
        },
        // onMouseover() {
        //   console.log('mouse over')
        // }
    }
}

const onLoadTree = async (node) => {
    switch (node.type) {
        case ConnectionType.Server:
            loading.value = true
            try {
                await connectionStore.openConnection(node.name)
            } catch (e) {
                message.error(e.message)
                node.isLeaf = undefined
            } finally {
                loading.value = false
            }
            break
        case ConnectionType.RedisDB:
            loading.value = true
            try {
                await connectionStore.openDatabase(node.name, node.db)
            } catch (e) {
                message.error(e.message)
                node.isLeaf = undefined
            } finally {
                loading.value = false
            }
            break
    }
}

const handleSelectContextMenu = (key) => {
    showContextMenu.value = false
    const { name, db, key: nodeKey, redisKey } = currentContextNode.value
    switch (key) {
        case 'server_disconnect':
            connectionStore.closeConnection(nodeKey).then((success) => {
                if (success) {
                    collapseKeyAndChildren(nodeKey)
                    tabStore.removeTabByName(name)
                }
            })
            break
        // case 'server_reload':
        // case 'db_reload':
        //     connectionStore.loadKeyValue()
        //     break
        case 'db_newkey':
        case 'key_newkey':
            dialogStore.openNewKeyDialog(redisKey, name, db)
            break
        case 'key_remove':
        case 'value_remove':
            dialog.warning({
                title: i18n.t('warning'),
                content: i18n.t('delete_key_tip', { key: redisKey }),
                closable: false,
                autoFocus: false,
                transformOrigin: 'center',
                positiveText: i18n.t('confirm'),
                negativeText: i18n.t('cancel'),
                onPositiveClick: () => {
                    connectionStore.removeKey(name, db, redisKey).then((success) => {
                        if (success) {
                            message.success(i18n.t('delete_key_succ', { key: redisKey }))
                        }
                    })
                },
            })
            break
        case 'key_copy':
        case 'value_copy':
            ClipboardSetText(redisKey)
                .then((succ) => {
                    if (succ) {
                        message.success(i18n.t('copy_succ'))
                    }
                })
                .catch((e) => {
                    message.error(e.message)
                })
            break
        default:
            console.warn('TODO: handle context menu:' + key)
    }
}

const handleOutsideContextMenu = () => {
    showContextMenu.value = false
}
</script>

<template>
    <n-tree
        :block-line="true"
        :block-node="true"
        :data="connectionStore.connections"
        :expand-on-click="false"
        :expanded-keys="expandedKeys"
        :node-props="nodeProps"
        :on-load="onLoadTree"
        :on-update:expanded-keys="onUpdateExpanded"
        :render-label="renderLabel"
        :render-prefix="renderPrefix"
        :render-suffix="renderSuffix"
        block-line
        class="fill-height"
        virtual-scroll
    />
    <n-dropdown
        :animated="false"
        :options="contextMenuOptions"
        :render-label="renderContextLabel"
        :show="showContextMenu"
        :x="contextPos.x"
        :y="contextPos.y"
        placement="bottom-start"
        trigger="manual"
        @clickoutside="handleOutsideContextMenu"
        @select="handleSelectContextMenu"
    />
</template>

<style lang="scss" scoped></style>
