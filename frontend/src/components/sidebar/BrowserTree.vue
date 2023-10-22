<script setup>
import { computed, h, nextTick, onMounted, reactive, ref } from 'vue'
import { ConnectionType } from '@/consts/connection_type.js'
import { NIcon, NSpace, NTag } from 'naive-ui'
import Key from '@/components/icons/Key.vue'
import Binary from '@/components/icons/Binary.vue'
import ToggleDb from '@/components/icons/ToggleDb.vue'
import { find, get, includes, indexOf, isEmpty, remove, size, startsWith } from 'lodash'
import { useI18n } from 'vue-i18n'
import Refresh from '@/components/icons/Refresh.vue'
import CopyLink from '@/components/icons/CopyLink.vue'
import Add from '@/components/icons/Add.vue'
import Layer from '@/components/icons/Layer.vue'
import Delete from '@/components/icons/Delete.vue'
import Connect from '@/components/icons/Connect.vue'
import useDialogStore from 'stores/dialog.js'
import { ClipboardSetText } from 'wailsjs/runtime/runtime.js'
import useConnectionStore from 'stores/connections.js'
import Unlink from '@/components/icons/Unlink.vue'
import Filter from '@/components/icons/Filter.vue'
import Close from '@/components/icons/Close.vue'
import { typesBgColor, typesColor } from '@/consts/support_redis_type.js'
import useTabStore from 'stores/tab.js'
import IconButton from '@/components/common/IconButton.vue'
import { parseHexColor } from '@/utils/rgb.js'
import LoadList from '@/components/icons/LoadList.vue'
import LoadAll from '@/components/icons/LoadAll.vue'

const props = defineProps({
    server: String,
    keyView: String,
})

const i18n = useI18n()
const loading = ref(false)
const loadingConnections = ref(false)
const expandedKeys = ref([props.server])
const connectionStore = useConnectionStore()
const tabStore = useTabStore()
const dialogStore = useDialogStore()

/**
 *
 * @type {ComputedRef<string[]>}
 */
const selectedKeys = computed(() => {
    const tab = find(tabStore.tabList, { name: props.server })
    if (tab != null) {
        return get(tab, 'selectedKeys', [props.server])
    }
    return [props.server]
})

const data = computed(() => {
    const dbs = get(connectionStore.databases, props.server, [])
    return dbs
})

const backgroundColor = computed(() => {
    const { markColor: hex = '' } = connectionStore.serverProfile[props.server] || {}
    if (isEmpty(hex)) {
        return ''
    }
    const { r, g, b } = parseHexColor(hex)
    return `rgba(${r}, ${g}, ${b}, 0.2)`
})

const contextMenuParam = reactive({
    show: false,
    x: 0,
    y: 0,
    options: null,
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
        return [
            {
                key: 'server_reload',
                label: i18n.t('interface.reload'),
                icon: renderIcon(Refresh),
            },
            {
                key: 'server_close',
                label: i18n.t('interface.disconnect'),
                icon: renderIcon(Unlink),
            },
        ]
    },
    [ConnectionType.RedisDB]: ({ opened }) => {
        if (opened) {
            return [
                {
                    key: 'db_reload',
                    label: i18n.t('interface.reload'),
                    icon: renderIcon(Refresh),
                },
                {
                    key: 'db_newkey',
                    label: i18n.t('interface.new_key'),
                    icon: renderIcon(Add),
                },
                {
                    key: 'db_filter',
                    label: i18n.t('interface.filter_key'),
                    icon: renderIcon(Filter),
                },
                {
                    type: 'divider',
                    key: 'd1',
                },
                {
                    key: 'key_remove',
                    label: i18n.t('interface.batch_delete'),
                    icon: renderIcon(Delete),
                },
                {
                    type: 'divider',
                    key: 'd2',
                },
                {
                    key: 'db_close',
                    label: i18n.t('interface.close_db'),
                    icon: renderIcon(Close),
                },
            ]
        } else {
            return [
                {
                    key: 'db_open',
                    label: i18n.t('interface.open_db'),
                    icon: renderIcon(Connect),
                },
            ]
        }
    },
    [ConnectionType.RedisKey]: () => [
        // {
        //     key: 'key_reload',
        //     label: i18n.t('interface.reload'),
        //     icon: renderIcon(Refresh),
        // },
        {
            key: 'key_newkey',
            label: i18n.t('interface.new_key'),
            icon: renderIcon(Add),
        },
        {
            key: 'key_copy',
            label: i18n.t('interface.copy_path'),
            icon: renderIcon(CopyLink),
        },
        {
            type: 'divider',
            key: 'd1',
        },
        {
            key: 'key_remove',
            label: i18n.t('interface.batch_delete'),
            icon: renderIcon(Delete),
        },
    ],
    [ConnectionType.RedisValue]: () => [
        {
            key: 'value_reload',
            label: i18n.t('interface.reload'),
            icon: renderIcon(Refresh),
        },
        {
            key: 'value_copy',
            label: i18n.t('interface.copy_key'),
            icon: renderIcon(CopyLink),
        },
        {
            type: 'divider',
            key: 'd1',
        },
        {
            key: 'value_remove',
            label: i18n.t('interface.remove_key'),
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

const resetExpandKey = (server, db, includeDB) => {
    const prefix = `${server}/db${db}`
    remove(expandedKeys.value, (k) => {
        if (!!!includeDB) {
            return k !== prefix && startsWith(k, prefix)
        } else {
            return startsWith(k, prefix)
        }
    })
}

const handleSelectContextMenu = (key) => {
    contextMenuParam.show = false
    const selectedKey = get(selectedKeys.value, 0)
    if (selectedKey == null) {
        return
    }
    const node = connectionStore.getNode(selectedKey)
    const { db = 0, key: nodeKey, redisKey: rk = '', redisKeyCode: rkc, label } = node || {}
    const redisKey = rkc || rk
    const redisKeyName = !!rkc ? label : redisKey
    switch (key) {
        case 'server_info':
            tabStore.setSelectedKeys(props.server)
            onUpdateSelectedKeys()
            break
        case 'server_reload':
            expandedKeys.value = [props.server]
            tabStore.setSelectedKeys(props.server)
            connectionStore.openConnection(props.server, true).then(() => {
                $message.success(i18n.t('dialogue.reload_succ'))
            })
            break
        case 'server_close':
            connectionStore.closeConnection(props.server)
            break
        case 'db_open':
            nextTick().then(() => expandKey(nodeKey))
            break
        case 'db_reload':
            resetExpandKey(props.server, db)
            connectionStore.reopenDatabase(props.server, db)
            break
        case 'db_close':
            resetExpandKey(props.server, db, true)
            connectionStore.closeDatabase(props.server, db)
            break
        case 'db_newkey':
        case 'key_newkey':
            dialogStore.openNewKeyDialog(redisKey, props.server, db)
            break
        case 'db_filter':
            const { match: pattern, type } = connectionStore.getKeyFilter(props.server, db)
            dialogStore.openKeyFilterDialog(props.server, db, pattern, type)
            break
        // case 'key_reload':
        //     connectionStore.loadKeys(props.server, db, redisKey)
        //     break
        case 'value_reload':
            connectionStore.loadKeyValue(props.server, db, redisKey)
            break
        case 'key_remove':
            dialogStore.openDeleteKeyDialog(props.server, db, isEmpty(redisKey) ? '*' : redisKey + ':*')
            break
        case 'value_remove':
            $dialog.warning(i18n.t('dialogue.remove_tip', { name: redisKeyName }), () => {
                connectionStore.deleteKey(props.server, db, redisKey).then((success) => {
                    if (success) {
                        $message.success(i18n.t('dialogue.delete_key_succ', { key: redisKeyName }))
                    }
                })
            })
            break
        case 'key_copy':
        case 'value_copy':
            ClipboardSetText(redisKey)
                .then((succ) => {
                    if (succ) {
                        $message.success(i18n.t('dialogue.copy_succ'))
                    }
                })
                .catch((e) => {
                    $message.error(e.message)
                })
            break
        case 'db_loadmore':
            if (node != null && !!!node.loading && !!!node.fullLoaded) {
                node.loading = true
                connectionStore
                    .loadMoreKeys(props.server, db)
                    .then((end) => {
                        // fully loaded
                        node.fullLoaded = end === true
                    })
                    .catch((e) => {
                        $message.error(e.message)
                    })
                    .finally(() => {
                        delete node.loading
                    })
            }
            break
        case 'db_loadall':
            if (node != null && !!!node.loading) {
                node.loading = true
                connectionStore
                    .loadAllKeys(props.server, db)
                    .catch((e) => {
                        $message.error(e.message)
                    })
                    .finally(() => {
                        delete node.loading
                        node.fullLoaded = true
                    })
            }
            break
        case 'more_action':
        default:
            console.warn('TODO: handle context menu:' + key)
    }
}

defineExpose({
    handleSelectContextMenu,
})

const onUpdateExpanded = (value, option, meta) => {
    expandedKeys.value = value
    if (!meta.node) {
        return
    }

    // keep expand or collapse children while they own more than 1 child
    let node = meta.node
    while (node != null && size(node.children) === 1) {
        const key = node.children[0].key
        switch (meta.action) {
            case 'expand':
                node.expanded = true
                if (!includes(expandedKeys.value, key)) {
                    expandedKeys.value.push(key)
                }
                break
            case 'collapse':
                node.expanded = false
                remove(expandedKeys.value, (v) => v === key)
                break
        }
        node = node.children[0]
    }
}

const onUpdateSelectedKeys = (keys, options) => {
    try {
        if (!isEmpty(options)) {
            // prevent load duplicate key
            for (const node of options) {
                if (node.type === ConnectionType.RedisValue) {
                    const { key, db } = node
                    const redisKey = node.redisKeyCode || node.redisKey
                    if (!includes(selectedKeys.value, key)) {
                        connectionStore.loadKeyValue(props.server, db, redisKey)
                    }
                    return
                }
            }
        }
        // default is load blank key to display server status
        connectionStore.loadKeyValue(props.server, 0)
    } finally {
        tabStore.setSelectedKeys(props.server, keys)
    }
}

const renderPrefix = ({ option }) => {
    switch (option.type) {
        // case ConnectionType.Server:
        //     const icon = option.cluster === true ? ToggleCluster : ToggleServer
        //     return h(
        //         NIcon,
        //         { size: 20 },
        //         {
        //             default: () => h(icon, { modelValue: false }),
        //         },
        //     )
        case ConnectionType.RedisDB:
            return h(
                NIcon,
                { size: 20 },
                {
                    default: () => h(ToggleDb, { modelValue: option.opened === true }),
                },
            )
        case ConnectionType.RedisKey:
            return h(
                NIcon,
                { size: 20 },
                {
                    default: () => h(Layer),
                },
            )
        case ConnectionType.RedisValue:
            return h(
                NIcon,
                { size: 20 },
                {
                    default: () => h(!!option.redisKeyCode ? Binary : Key),
                },
            )
    }
}

// render tree item label
const renderLabel = ({ option }) => {
    switch (option.type) {
        case ConnectionType.Server:
            return h('b', {}, { default: () => option.label })
        case ConnectionType.RedisDB:
            const { name: server, db, opened = false } = option
            let { match: matchPattern, type: typeFilter } = connectionStore.getKeyFilter(server, db)
            const items = []
            if (opened) {
                items.push(`${option.label} (${option.keys || 0}/${Math.max(option.maxKeys || 0, option.keys || 0)})`)
            } else {
                items.push(`${option.label} (${Math.max(option.maxKeys || 0, option.keys || 0)})`)
            }
            // show filter tag after label
            // type filter tag
            if (!isEmpty(typeFilter)) {
                items.push(
                    h(
                        NTag,
                        {
                            size: 'small',
                            closable: true,
                            bordered: false,
                            color: {
                                color: typesBgColor[typeFilter],
                                textColor: typesColor[typeFilter],
                            },
                            onClose: () => {
                                // remove type filter
                                connectionStore.setKeyFilter(server, db, matchPattern)
                                connectionStore.reopenDatabase(server, db)
                            },
                        },
                        { default: () => typeFilter },
                    ),
                )
            }
            // match pattern tag
            if (!isEmpty(matchPattern) && matchPattern !== '*') {
                items.push(
                    h(
                        NTag,
                        {
                            bordered: false,
                            closable: true,
                            size: 'small',
                            onClose: () => {
                                // remove key match pattern
                                connectionStore.setKeyFilter(server, db, '*', typeFilter)
                                connectionStore.reopenDatabase(server, db)
                            },
                        },
                        { default: () => matchPattern },
                    ),
                )
            }
            return renderIconMenu(items)
        case ConnectionType.RedisKey:
            return `${option.label} (${option.keys || 0})`
        // case ConnectionType.RedisValue:
        //   return `[${option.keyType}]${option.label}`
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

const getDatabaseMenu = (opened, loading, end) => {
    const btns = []
    if (opened) {
        btns.push(
            h(IconButton, {
                tTooltip: 'interface.filter_key',
                icon: Filter,
                disabled: loading === true,
                onClick: () => handleSelectContextMenu('db_filter'),
            }),
            h(IconButton, {
                tTooltip: 'interface.reload',
                icon: Refresh,
                disabled: loading === true,
                onClick: () => handleSelectContextMenu('db_reload'),
            }),
            h(IconButton, {
                tTooltip: 'interface.new_key',
                icon: Add,
                disabled: loading === true,
                onClick: () => handleSelectContextMenu('db_newkey'),
            }),
            h(IconButton, {
                tTooltip: 'interface.load_more',
                icon: LoadList,
                disabled: end === true,
                loading: loading === true,
                onClick: () => handleSelectContextMenu('db_loadmore'),
            }),
            h(IconButton, {
                tTooltip: 'interface.load_all',
                icon: LoadAll,
                disabled: end === true,
                loading: loading === true,
                onClick: () => handleSelectContextMenu('db_loadall'),
            }),
            h(IconButton, {
                tTooltip: 'interface.batch_delete',
                icon: Delete,
                disabled: loading === true,
                onClick: () => handleSelectContextMenu('key_remove'),
            }),
            // h(IconButton, {
            //     tTooltip: 'interface.more_action',
            //     icon: More,
            //     onClick: () => handleSelectContextMenu('more_action'),
            // }),
        )
    } else {
        btns.push(
            h(IconButton, {
                tTooltip: 'interface.open_db',
                icon: Connect,
                onClick: () => handleSelectContextMenu('db_open'),
            }),
        )
    }
    return btns
}

const getLayerMenu = () => {
    return [
        // disable reload by layer, due to conflict with partial loading keys
        // h(IconButton, {
        //     tTooltip: 'interface.reload',
        //     icon: Refresh,
        //     onClick: () => handleSelectContextMenu('key_reload'),
        // }),
        h(IconButton, {
            tTooltip: 'interface.new_key',
            icon: Add,
            onClick: () => handleSelectContextMenu('key_newkey'),
        }),
        h(IconButton, {
            tTooltip: 'interface.batch_delete',
            icon: Delete,
            onClick: () => handleSelectContextMenu('key_remove'),
        }),
    ]
}

const getValueMenu = () => {
    return [
        h(IconButton, {
            tTooltip: 'interface.remove_key',
            icon: Delete,
            onClick: () => handleSelectContextMenu('value_remove'),
        }),
    ]
}

// render menu function icon
const renderSuffix = ({ option }) => {
    if ((option.type === ConnectionType.RedisDB && option.opened) || includes(selectedKeys.value, option.key)) {
        switch (option.type) {
            case ConnectionType.RedisDB:
                return renderIconMenu(getDatabaseMenu(option.opened, option.loading, option.fullLoaded))
            case ConnectionType.RedisKey:
                return renderIconMenu(getLayerMenu())
            case ConnectionType.RedisValue:
                return renderIconMenu(getValueMenu())
        }
    }
    return null
}

const nodeProps = ({ option }) => {
    return {
        onDblclick: () => {
            if (loading.value) {
                console.warn('TODO: alert to ignore double click when loading')
                return
            }
            // default handle is expand current node
            nextTick().then(() => expandKey(option.key))
        },
        onContextmenu(e) {
            e.preventDefault()
            if (!menuOptions.hasOwnProperty(option.type)) {
                return
            }
            contextMenuParam.show = false
            contextMenuParam.options = menuOptions[option.type](option)
            nextTick().then(() => {
                contextMenuParam.x = e.clientX
                contextMenuParam.y = e.clientY
                contextMenuParam.show = true
                tabStore.setSelectedKeys(props.server, option.key)
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
                $message.error(e.message)
                node.isLeaf = undefined
            } finally {
                loading.value = false
            }
            break
        // case ConnectionType.RedisKey:
        //     console.warn('load redis key', node.redisKey)
        //     node.keys = sumBy(node.children, 'keys')
        //     break
        // case ConnectionType.RedisValue:
        //     node.keys = 1
        //     break
    }
}

const handleOutsideContextMenu = () => {
    contextMenuParam.show = false
}
</script>

<template>
    <div :style="{ backgroundColor }" class="browser-tree-wrapper">
        <n-tree
            :animated="false"
            :block-line="true"
            :block-node="true"
            :cancelable="false"
            :data="data"
            :expand-on-click="false"
            :expanded-keys="expandedKeys"
            :node-props="nodeProps"
            :render-label="renderLabel"
            :render-prefix="renderPrefix"
            :render-suffix="renderSuffix"
            :selected-keys="selectedKeys"
            class="fill-height"
            virtual-scroll
            @load="onLoadTree"
            @update:selected-keys="onUpdateSelectedKeys"
            @update:expanded-keys="onUpdateExpanded" />
        <n-dropdown
            :options="contextMenuParam.options"
            :render-label="renderContextLabel"
            :show="contextMenuParam.show"
            :x="contextMenuParam.x"
            :y="contextMenuParam.y"
            placement="bottom-start"
            trigger="manual"
            @clickoutside="handleOutsideContextMenu"
            @select="handleSelectContextMenu" />
    </div>
</template>

<style lang="scss" scoped>
.browser-tree-wrapper {
    height: 100%;
    overflow: hidden;
}
</style>
