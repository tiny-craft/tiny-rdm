<script setup>
import { computed, h, nextTick, onMounted, reactive, ref } from 'vue'
import { ConnectionType } from '../../consts/connection_type.js'
import { NIcon, NTag, useDialog, useMessage } from 'naive-ui'
import Key from '../icons/Key.vue'
import ToggleDb from '../icons/ToggleDb.vue'
import { get, indexOf, isEmpty, remove } from 'lodash'
import { useI18n } from 'vue-i18n'
import Refresh from '../icons/Refresh.vue'
import CopyLink from '../icons/CopyLink.vue'
import Add from '../icons/Add.vue'
import Layer from '../icons/Layer.vue'
import Delete from '../icons/Delete.vue'
import Connect from '../icons/Connect.vue'
import useDialogStore from '../../stores/dialog.js'
import { ClipboardSetText } from '../../../wailsjs/runtime/runtime.js'
import useConnectionStore from '../../stores/connections.js'
import { useConfirmDialog } from '../../utils/confirm_dialog.js'
import ToggleServer from '../icons/ToggleServer.vue'
import Unlink from '../icons/Unlink.vue'
import Filter from '../icons/Filter.vue'
import Close from '../icons/Close.vue'

const props = defineProps({
    server: String,
})

const i18n = useI18n()
const loading = ref(false)
const loadingConnections = ref(false)
const expandedKeys = ref([props.server])
const selectedKeys = ref([props.server])
const connectionStore = useConnectionStore()
const dialogStore = useDialogStore()

const data = computed(() => {
    const dbs = get(connectionStore.databases, props.server, [])
    return [
        {
            key: props.server,
            label: props.server,
            type: ConnectionType.Server,
            children: dbs,
        },
    ]
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
    [ConnectionType.Server]: () => {
        console.log('open server context')
        return [
            {
                key: 'server_reload',
                label: i18n.t('reload'),
                icon: renderIcon(Refresh),
            },
            {
                key: 'server_close',
                label: i18n.t('disconnect'),
                icon: renderIcon(Unlink),
            },
        ]
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
                {
                    key: 'db_filter',
                    label: i18n.t('filter_key'),
                    icon: renderIcon(Filter),
                },
                {
                    type: 'divider',
                    key: 'd1',
                },
                {
                    key: 'db_close',
                    label: i18n.t('close_db'),
                    icon: renderIcon(Close),
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
        // nextTick(connectionStore.initConnection)
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

const onUpdateSelectedKeys = (keys, options) => {
    try {
        if (!isEmpty(options)) {
            // prevent load duplicate key
            for (const node of options) {
                if (node.type === ConnectionType.RedisValue) {
                    const { key, db, redisKey } = node
                    if (indexOf(selectedKeys.value, key) === -1) {
                        connectionStore.loadKeyValue(props.server, db, redisKey)
                    }
                    return
                }
            }

            // default is load blank key to display server status
            connectionStore.loadKeyValue(props.server, 0)
        }
    } finally {
        selectedKeys.value = keys
    }
}

const renderPrefix = ({ option }) => {
    switch (option.type) {
        case ConnectionType.Server:
            return h(
                NIcon,
                { size: 20 },
                {
                    default: () => h(ToggleServer, { modelValue: false }),
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
    if (option.type === ConnectionType.RedisDB) {
        const { name: server, db } = option
        const filterPattern = connectionStore.getKeyFilter(server, db)
        if (!isEmpty(filterPattern) && filterPattern !== '*') {
            return h(
                NTag,
                {
                    bordered: false,
                    closable: true,
                    size: 'small',
                    onClose: () => {
                        connectionStore.removeKeyFilter(server, db)
                        connectionStore.reopenDatabase(server, db)
                    },
                },
                { default: () => filterPattern }
            )
        }
    }
    return null
}

const nodeProps = ({ option }) => {
    return {
        onDblclick: async () => {
            if (loading.value) {
                console.warn('TODO: alert to ignore double click when loading')
                return
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
        // onMouseover() {
        //   console.log('mouse over')
        // }
    }
}

const onLoadTree = async (node) => {
    switch (node.type) {
        case ConnectionType.RedisDB:
            loading.value = true
            try {
                await connectionStore.openDatabase(props.server, node.db)
            } catch (e) {
                message.error(e.message)
                node.isLeaf = undefined
            } finally {
                loading.value = false
            }
            break
        case ConnectionType.RedisKey:
            // load all children
            // node.children = []
            break
    }
}

const confirmDialog = useConfirmDialog()
const handleSelectContextMenu = (key) => {
    contextMenuParam.show = false
    const { db, key: nodeKey, redisKey } = contextMenuParam.currentNode
    switch (key) {
        case 'server_reload':
            connectionStore.openConnection(props.server, true).then(() => {
                message.success(i18n.t('reload_succ'))
            })
            break
        case 'server_close':
            connectionStore.closeConnection(props.server)
            break
        case 'db_open':
            nextTick().then(() => expandKey(nodeKey))
            break
        case 'db_reload':
            connectionStore.reopenDatabase(props.server, db)
            break
        case 'db_close':
            remove(expandedKeys.value, (k) => k === `${props.server}/db${db}`)
            connectionStore.closeDatabase(props.server, db)
            break
        case 'db_newkey':
        case 'key_newkey':
            dialogStore.openNewKeyDialog(redisKey, props.server, db)
            break
        case 'db_filter':
            const pattern = connectionStore.getKeyFilter(props.server, db)
            dialogStore.openKeyFilterDialog(props.server, db, pattern)
            break
        case 'key_reload':
            connectionStore.loadKeys(props.server, db, redisKey)
            break
        case 'value_reload':
            connectionStore.loadKeyValue(props.server, db, redisKey)
            break
        case 'key_remove':
            dialogStore.openDeleteKeyDialog(props.server, db, redisKey + ':*')
            break
        case 'value_remove':
            confirmDialog.warning(i18n.t('remove_tip', { name: redisKey }), () => {
                connectionStore.deleteKey(props.server, db, redisKey).then((success) => {
                    if (success) {
                        message.success(i18n.t('delete_key_succ', { key: redisKey }))
                    }
                })
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
    contextMenuParam.show = false
}
</script>

<template>
    <n-tree
        :block-line="true"
        :block-node="true"
        :animated="false"
        :cancelable="false"
        :data="data"
        :expand-on-click="false"
        :expanded-keys="expandedKeys"
        :selected-keys="selectedKeys"
        @update:selected-keys="onUpdateSelectedKeys"
        :node-props="nodeProps"
        @load="onLoadTree"
        @update:expanded-keys="onUpdateExpanded"
        :render-label="renderLabel"
        :render-prefix="renderPrefix"
        :render-suffix="renderSuffix"
        class="fill-height"
        virtual-scroll
    />
    <n-dropdown
        :animated="false"
        :options="contextMenuParam.options"
        :render-label="renderContextLabel"
        :show="contextMenuParam.show"
        :x="contextMenuParam.x"
        :y="contextMenuParam.y"
        placement="bottom-start"
        trigger="manual"
        @clickoutside="handleOutsideContextMenu"
        @select="handleSelectContextMenu"
    />
</template>

<style lang="scss" scoped></style>
