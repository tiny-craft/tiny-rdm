<script setup>
import { computed, h, nextTick, reactive, ref, watchEffect } from 'vue'
import { ConnectionType } from '@/consts/connection_type.js'
import { NIcon, NSpace, useThemeVars } from 'naive-ui'
import Key from '@/components/icons/Key.vue'
import Binary from '@/components/icons/Binary.vue'
import Database from '@/components/icons/Database.vue'
import { filter, find, get, includes, indexOf, isEmpty, map, remove, size, startsWith, toUpper } from 'lodash'
import { useI18n } from 'vue-i18n'
import Refresh from '@/components/icons/Refresh.vue'
import CopyLink from '@/components/icons/CopyLink.vue'
import Add from '@/components/icons/Add.vue'
import Layer from '@/components/icons/Layer.vue'
import Delete from '@/components/icons/Delete.vue'
import useDialogStore from 'stores/dialog.js'
import { ClipboardSetText } from 'wailsjs/runtime/runtime.js'
import useConnectionStore from 'stores/connections.js'
import Filter from '@/components/icons/Filter.vue'
import useTabStore from 'stores/tab.js'
import IconButton from '@/components/common/IconButton.vue'
import { parseHexColor } from '@/utils/rgb.js'
import LoadList from '@/components/icons/LoadList.vue'
import LoadAll from '@/components/icons/LoadAll.vue'
import useBrowserStore from 'stores/browser.js'
import { useRender } from '@/utils/render.js'
import RedisTypeTag from '@/components/common/RedisTypeTag.vue'
import usePreferencesStore from 'stores/preferences.js'
import { typesIconStyle } from '@/consts/support_redis_type.js'

const props = defineProps({
    server: String,
    db: Number,
    keyView: String,
    loading: Boolean,
    pattern: String,
    fullLoaded: Boolean,
    checkMode: Boolean,
})

const themeVars = useThemeVars()
const render = useRender()
const i18n = useI18n()
/** @type {Ref<UnwrapRef<string[]>>} */
const expandedKeys = ref([])
const connectionStore = useConnectionStore()
const browserStore = useBrowserStore()
const prefStore = usePreferencesStore()
const tabStore = useTabStore()
const dialogStore = useDialogStore()

/**
 *
 * @type {ComputedRef<string[]>}
 */
const selectedKeys = computed(() => {
    const tab = find(tabStore.tabList, { name: props.server })
    if (tab != null) {
        return get(tab, 'selectedKeys', [])
    }
    return []
})

/**
 *
 * @type {ComputedRef<string[]>}
 */
const checkedKeys = computed(() => {
    const tab = find(tabStore.tabList, { name: props.server })
    if (tab != null) {
        const checkedKeys = get(tab, 'checkedKeys', [])
        return map(checkedKeys, 'key')
    }
    return []
})

const data = computed(() => {
    return browserStore.getKeyStruct(props.server, props.checkMode)
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

const menuOptions = {
    [ConnectionType.RedisDB]: ({ opened }) => {
        if (opened) {
            return [
                {
                    key: 'db_filter',
                    label: i18n.t('interface.filter_key'),
                    icon: render.renderIcon(Filter),
                },
                {
                    type: 'divider',
                    key: 'd1',
                },
                {
                    type: 'divider',
                    key: 'd2',
                },
            ]
        } else {
            return []
        }
    },
    [ConnectionType.RedisKey]: () => [
        // {
        //     key: 'key_reload',
        //     label: i18n.t('interface.reload'),
        //     icon: render.renderIcon(Refresh),
        // },
        {
            key: 'key_newkey',
            label: i18n.t('interface.new_key'),
            icon: render.renderIcon(Add),
        },
        {
            key: 'key_copy',
            label: i18n.t('interface.copy_path'),
            icon: render.renderIcon(CopyLink),
        },
        {
            type: 'divider',
            key: 'd1',
        },
        {
            key: 'key_remove',
            label: i18n.t('interface.batch_delete_key'),
            icon: render.renderIcon(Delete),
        },
    ],
    [ConnectionType.RedisValue]: () => [
        {
            key: 'value_reload',
            label: i18n.t('interface.reload'),
            icon: render.renderIcon(Refresh),
        },
        {
            key: 'value_copy',
            label: i18n.t('interface.copy_key'),
            icon: render.renderIcon(CopyLink),
        },
        {
            type: 'divider',
            key: 'd1',
        },
        {
            key: 'value_remove',
            label: i18n.t('interface.remove_key'),
            icon: render.renderIcon(Delete),
        },
    ],
}

const renderContextLabel = (option) => {
    return h('div', { class: 'context-menu-item' }, option.label)
}

/**
 *
 * @param {string} key
 * @param {boolean} [toggle]
 */
const expandKey = (key, toggle) => {
    const idx = indexOf(expandedKeys.value, key)
    if (idx === -1) {
        expandedKeys.value.push(key)
    } else if (toggle !== false) {
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
    const node = browserStore.getNode(selectedKey)
    const { db = 0, key: nodeKey, redisKey: rk = '', redisKeyCode: rkc, label } = node || {}
    const redisKey = rkc || rk
    const redisKeyName = !!rkc ? label : redisKey
    switch (key) {
        case 'key_newkey':
            dialogStore.openNewKeyDialog(redisKey, props.server, db)
            break
        case 'db_filter':
            const { match: pattern, type } = browserStore.getKeyFilter(props.server)
            dialogStore.openKeyFilterDialog(props.server, db, pattern, type)
            break
        case 'key_reload':
            if (node != null && !!!node.loading) {
                node.loading = true
                browserStore.reloadLayer(props.server, db, redisKey).finally(() => {
                    delete node.loading
                })
            }
            break
        case 'value_reload':
            browserStore.reloadKey({
                server: props.server,
                db,
                key: redisKey,
            })
            break
        case 'key_remove':
            dialogStore.openDeleteKeyDialog(props.server, db, isEmpty(redisKey) ? '*' : redisKey + ':*')
            break
        case 'value_remove':
            $dialog.warning(i18n.t('dialogue.remove_tip', { name: redisKeyName }), () => {
                browserStore.deleteKey(props.server, db, redisKey).then((success) => {
                    if (success) {
                        $message.success(i18n.t('dialogue.delete.success', { key: redisKeyName }))
                    }
                })
            })
            break
        case 'key_copy':
        case 'value_copy':
            ClipboardSetText(redisKey)
                .then((succ) => {
                    if (succ) {
                        $message.success(i18n.t('interface.copy_succ'))
                    }
                })
                .catch((e) => {
                    $message.error(e.message)
                })
            break
        case 'db_loadall':
            if (node != null && !!!node.loading) {
                node.loading = true
                browserStore
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

const onUpdateSelectedKeys = (keys, options) => {
    try {
        if (!isEmpty(options)) {
            // prevent load duplicate key
            for (const node of options) {
                if (node.type === ConnectionType.RedisValue) {
                    const { key, db } = node
                    const redisKey = node.redisKeyCode || node.redisKey
                    if (!includes(selectedKeys.value, key)) {
                        browserStore.loadKeySummary({
                            server: props.server,
                            db,
                            key: redisKey,
                            clearValue: true,
                        })
                    }
                    return
                }
            }
        }
        // default is load blank key to display server status
        tabStore.openBlank(props.server, props.db)
    } finally {
        tabStore.setSelectedKeys(props.server, keys)
    }
}

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

/**
 *
 * @param {string[]} keys
 * @param {TreeOption[]} options
 */
const onUpdateCheckedKeys = (keys, options) => {
    const checkedKeys = map(
        filter(options, (o) => o.type === ConnectionType.RedisValue),
        (o) => ({ key: o.key, redisKey: o.redisKeyCode || o.redisKey }),
    )
    tabStore.setCheckedKeys(props.server, checkedKeys)
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
                { size: 20, color: option.opened === true ? '#dc423c' : undefined },
                {
                    default: () => h(Database, { inverse: option.opened === true }),
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
            if (prefStore.keyIconType === typesIconStyle.ICON) {
                return h(NIcon, { size: 20 }, () => h(Key))
            }
            const loading = isEmpty(option.redisType) || option.redisType === 'loading'
            if (loading) {
                browserStore.loadKeyType({
                    server: props.server,
                    db: option.db,
                    key: option.redisKey,
                    keyCode: option.redisKeyCode,
                })
            }
            switch (prefStore.keyIconType) {
                case typesIconStyle.FULL:
                    return h(RedisTypeTag, {
                        type: toUpper(option.redisType),
                        short: false,
                        size: 'small',
                        inverse: includes(selectedKeys.value, option.key),
                        loading,
                    })

                case typesIconStyle.POINT:
                    return h(RedisTypeTag, {
                        type: toUpper(option.redisType),
                        point: true,
                        loading,
                    })

                case typesIconStyle.SHORT:
                default:
                    return h(RedisTypeTag, {
                        type: toUpper(option.redisType),
                        short: true,
                        size: 'small',
                        loading,
                        inverse: includes(selectedKeys.value, option.key),
                    })
            }
    }
}

// render tree item label
const renderLabel = ({ option }) => {
    switch (option.type) {
        case ConnectionType.RedisKey:
            return `${option.label} (${option.keyCount || 0})`
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

const calcDBMenu = (opened, loading, end) => {
    const btns = []
    if (opened) {
        btns.push(
            h(IconButton, {
                tTooltip: 'interface.load_more',
                icon: LoadList,
                disabled: end === true,
                loading: loading === true,
                color: loading === true ? themeVars.value.primaryColor : '',
                onClick: () => handleSelectContextMenu('db_loadmore'),
            }),
            h(IconButton, {
                tTooltip: 'interface.load_all',
                icon: LoadAll,
                disabled: end === true,
                loading: loading === true,
                color: loading === true ? themeVars.value.primaryColor : '',
                onClick: () => handleSelectContextMenu('db_loadall'),
            }),
            // h(IconButton, {
            //     tTooltip: 'interface.more_action',
            //     icon: More,
            //     onClick: () => handleSelectContextMenu('more_action'),
            // }),
        )
    }
    return btns
}

const calcLayerMenu = (loading) => {
    return [
        // reload layer enable only full loaded
        h(IconButton, {
            tTooltip: props.fullLoaded ? 'interface.reload' : 'interface.reload_disable',
            icon: Refresh,
            loading: loading === true,
            disabled: !props.fullLoaded,
            onClick: () => handleSelectContextMenu('key_reload'),
        }),
        h(IconButton, {
            tTooltip: 'interface.new_key',
            icon: Add,
            onClick: () => handleSelectContextMenu('key_newkey'),
        }),
        h(IconButton, {
            tTooltip: 'interface.batch_delete_key',
            icon: Delete,
            onClick: () => handleSelectContextMenu('key_remove'),
        }),
    ]
}

const calcValueMenu = () => {
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
    const selected = includes(selectedKeys.value, option.key)
    if (selected && !props.checkMode) {
        switch (option.type) {
            case ConnectionType.RedisDB:
                return renderIconMenu(calcDBMenu(option.opened, option.loading, option.fullLoaded))
            case ConnectionType.RedisKey:
                return renderIconMenu(calcLayerMenu(option.loading))
            case ConnectionType.RedisValue:
                return renderIconMenu(calcValueMenu())
        }
    } else if (!selected && !!option.redisKeyCode && option.type === ConnectionType.RedisValue) {
        // render binary icon
        return renderIconMenu(h(NIcon, { size: 20 }, () => h(Binary)))
    }
    return null
}

const nodeProps = ({ option }) => {
    return {
        onDblclick: () => {
            if (props.loading) {
                console.warn('TODO: alert to ignore double click when loading')
                return
            }
            if (!props.checkMode) {
                // default handle is expand current node
                nextTick().then(() => expandKey(option.key))
            }
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
                onUpdateSelectedKeys([option.key], [option])
            })
        },
        // onMouseover() {
        //   console.log('mouse over')
        // }
    }
}

const handleOutsideContextMenu = () => {
    contextMenuParam.show = false
}

watchEffect(
    () => {
        if (!props.checkMode) {
            tabStore.setCheckedKeys(props.server)
        } else {
            expandKey(`${ConnectionType.RedisDB}`, false)
        }
    },
    { flush: 'post' },
)

// the NTree node may get incorrect height after change data
// add key property to force refresh the component and then everything back to normal
const treeKey = ref(0)
defineExpose({
    handleSelectContextMenu,
    resetExpandKey,
    refreshTree: () => {
        treeKey.value = Date.now()
        expandedKeys.value = []
        tabStore.setCheckedKeys(props.server)
    },
    deleteCheckedItems: () => {
        const checkedKeys = tabStore.currentCheckedKeys
        const redisKeys = map(checkedKeys, 'redisKey')
        if (!isEmpty(redisKeys)) {
            dialogStore.openDeleteKeyDialog(props.server, props.db, redisKeys)
        }
    },
    exportCheckedItems: () => {
        const checkedKeys = tabStore.currentCheckedKeys
        const redisKeys = map(checkedKeys, 'redisKey')
        if (!isEmpty(redisKeys)) {
            dialogStore.openExportKeyDialog(props.server, props.db, redisKeys)
        }
    },
    updateTTLCheckedItems: () => {
        const checkedKeys = tabStore.currentCheckedKeys
        const redisKeys = map(checkedKeys, 'redisKey')
        if (!isEmpty(redisKeys)) {
            dialogStore.openTTLDialog({
                server: props.server,
                db: props.db,
                keys: redisKeys,
            })
        }
    },
    getSelectedKey: () => {
        return selectedKeys.value || []
    },
})
</script>

<template>
    <div :style="{ backgroundColor }" class="flex-box-v browser-tree-wrapper">
        <n-spin v-if="props.loading" class="fill-height" />
        <n-empty v-else-if="!props.loading && isEmpty(data)" class="empty-content" />
        <n-tree
            v-show="!props.loading && !isEmpty(data)"
            :key="treeKey"
            :animated="false"
            :block-line="true"
            :block-node="true"
            :cancelable="false"
            :cascade="true"
            :checkable="props.checkMode"
            :checked-keys="checkedKeys"
            :data="data"
            :expand-on-click="false"
            :expanded-keys="expandedKeys"
            :filter="(pattern, node) => includes(node.redisKey, pattern)"
            :node-props="nodeProps"
            :pattern="props.pattern"
            :render-label="renderLabel"
            :render-prefix="renderPrefix"
            :render-suffix="renderSuffix"
            :selected-keys="selectedKeys"
            :show-irrelevant-nodes="false"
            check-strategy="child"
            class="fill-height"
            virtual-scroll
            @update:selected-keys="onUpdateSelectedKeys"
            @update:expanded-keys="onUpdateExpanded"
            @update:checked-keys="onUpdateCheckedKeys" />
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
@import '@/styles/content';

.browser-tree-wrapper {
    height: 100%;
    overflow: hidden;
}
</style>
