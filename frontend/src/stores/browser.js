import { defineStore } from 'pinia'
import {
    endsWith,
    find,
    get,
    isEmpty,
    join,
    last,
    map,
    remove,
    set,
    size,
    slice,
    sortedIndexBy,
    split,
    sumBy,
    toUpper,
} from 'lodash'
import {
    AddHashField,
    AddListItem,
    AddStreamValue,
    AddZSetValue,
    CleanCmdHistory,
    CloseConnection,
    ConvertValue,
    DeleteKey,
    FlushDB,
    GetCmdHistory,
    GetKeyDetail,
    GetKeySummary,
    GetSlowLogs,
    LoadAllKeys,
    LoadNextKeys,
    OpenConnection,
    OpenDatabase,
    RemoveStreamValues,
    RenameKey,
    ServerInfo,
    SetHashValue,
    SetKeyTTL,
    SetKeyValue,
    SetListItem,
    SetSetItem,
    UpdateSetItem,
    UpdateZSetValue,
} from 'wailsjs/go/services/browserService.js'
import useTabStore from 'stores/tab.js'
import { decodeRedisKey, nativeRedisKey } from '@/utils/key_convert.js'
import { BrowserTabType } from '@/consts/browser_tab_type.js'
import { KeyViewType } from '@/consts/key_view_type.js'
import { ConnectionType } from '@/consts/connection_type.js'
import useConnectionStore from 'stores/connections.js'
import { decodeTypes, formatTypes } from '@/consts/value_view_type.js'
import { isRedisGlob } from '@/utils/glob_pattern.js'

const useBrowserStore = defineStore('browser', {
    /**
     * @typedef {Object} DatabaseItem
     * @property {string} key - tree node unique key
     * @property {string} label
     * @property {string} [name] - server name, type != ConnectionType.Group only
     * @property {number} type
     * @property {number} [db] - database index, type == ConnectionType.RedisDB only
     * @property {string} [redisKey] - redis key, type == ConnectionType.RedisKey || type == ConnectionType.RedisValue only
     * @property {number[]} [redisKeyCode] - redis key char code array, optional for redis key which contains binary data
     * @property {number} [keys] - children key count
     * @property {number} [maxKeys] - max key count for database
     * @property {boolean} [isLeaf]
     * @property {boolean} [opened] - redis db is opened, type == ConnectionType.RedisDB only
     * @property {boolean} [expanded] - current node is expanded
     * @property {DatabaseItem[]} [children]
     */

    /**
     * @typedef {Object} FilterItem
     * @property {string} pattern key pattern filter
     * @property {string} type type filter
     */

    /**
     * @typedef {Object} LoadingState
     * @property {string} loading indicated that is loading children now
     * @property {string} fullLoaded indicated that all children already loaded
     */

    /**
     * @typedef {Object} HistoryItem
     * @property {string} time
     * @property {string} server
     * @property {string} cmd
     * @property {number} cost
     */

    /**
     * @typedef {Object} BrowserState
     * @property {Object} serverStats
     * @property {Object.<string, FilterItem>} filter
     * @property {Object.<string, KeyViewType>} viewType
     * @property {Object.<string, DatabaseItem[]>} databases
     * @property {Object.<string, Map<string, DatabaseItem>>} nodeMap key format likes 'server#db', children key format likes 'key#type'
     */

    /**
     *
     * @returns {BrowserState}
     */
    state: () => ({
        serverStats: {}, // current server status info
        filter: {}, // all filters in opened connections map by server and FilterItem
        loadingState: {}, // all loading state in opened connections map by server and LoadingState
        viewType: {}, // view type selection for all opened connections group by 'server'
        databases: {}, // all databases in opened connections group by 'server name'
        nodeMap: {}, // all nodes in opened connections group by 'server#db' and 'type/key'
        keySet: {}, // all keys set in opened connections group by 'server#db
        openedDB: {}, // opened database map by server and database index
    }),
    getters: {
        anyConnectionOpened() {
            return !isEmpty(this.databases)
        },
        selectedDatabases() {
            return this.openedDB || {}
        },
    },
    actions: {
        /**
         * check if connection is connected
         * @param name
         * @returns {boolean}
         */
        isConnected(name) {
            let dbs = get(this.databases, name, [])
            return !isEmpty(dbs)
        },

        /**
         * close all connections
         * @returns {Promise<void>}
         */
        async closeAllConnection() {
            for (const name in this.databases) {
                await CloseConnection(name)
            }

            this.databases = {}
            this.nodeMap.clear()
            this.keySet.clear()
            this.serverStats = {}
            const tabStore = useTabStore()
            tabStore.removeAllTab()
        },

        /**
         * get database info list
         * @param server
         * @return {DatabaseItem[]}
         */
        getDBList(server) {
            return this.databases[server] || []
        },

        /**
         * get database by server name and database index
         * @param {string} server
         * @param {number} db
         * @return {DatabaseItem|null}
         */
        getDatabase(server, db) {
            const dbs = this.databases[server]
            if (dbs != null) {
                const selDB = find(dbs, (item) => item.db === db)
                if (selDB != null) {
                    return selDB
                }
            }
            return null
        },

        /**
         * get current selection database by server
         * @param server
         * @return {number}
         */
        getSelectedDB(server) {
            return this.selectedDatabases[server] || 0
        },

        /**
         * get key list in current database
         * @param server
         * @return {DatabaseItem[]}
         */
        getKeyList(server) {
            const db = this.getSelectedDB(server)
            const dbNodes = this.databases[server]
            const node = find(dbNodes, (n) => n.db === db)
            return node.children
        },

        /**
         * switch key view
         * @param {string} connName
         * @param {number} viewType
         */
        async switchKeyView(connName, viewType) {
            if (viewType !== KeyViewType.Tree && viewType !== KeyViewType.List) {
                return
            }

            const t = get(this.viewType, connName, KeyViewType.Tree)
            if (t === viewType) {
                return
            }

            this.viewType[connName] = viewType
            const dbs = get(this.databases, connName, [])
            for (const dbItem of dbs) {
                if (!dbItem.opened) {
                    continue
                }

                dbItem.children = undefined
                dbItem.keys = 0
                const { db = 0 } = dbItem
                this._getNodeMap(connName, db).clear()
                const keys = this._getKeySet(connName, db)
                this._addKeyNodes(connName, db, keys)
                this._tidyNode(connName, db, '')
            }
        },

        /**
         * open connection
         * @param {string} name
         * @param {boolean} [reload]
         * @returns {Promise<void>}
         */
        async openConnection(name, reload) {
            if (this.isConnected(name)) {
                if (reload !== true) {
                    return
                } else {
                    // reload mode, try close connection first
                    await CloseConnection(name)
                }
            }

            const { data, success, msg } = await OpenConnection(name)
            if (!success) {
                throw new Error(msg)
            }
            // append to db node to current connection
            // const connNode = this.getConnection(name)
            // if (connNode == null) {
            //     throw new Error('no such connection')
            // }
            const { db, view = KeyViewType.Tree } = data
            if (isEmpty(db)) {
                throw new Error('no db loaded')
            }
            const dbs = []
            for (let i = 0; i < db.length; i++) {
                this._getNodeMap(name, i).clear()
                this._getKeySet(name, i).clear()
                dbs.push({
                    key: `${name}/${db[i].name}`,
                    label: db[i].name,
                    name: name,
                    keys: 0,
                    maxKeys: db[i].keys,
                    db: db[i].index,
                    type: ConnectionType.RedisDB,
                    isLeaf: false,
                    children: undefined,
                })
            }
            this.databases[name] = dbs
            this.viewType[name] = view
            this.openedDB[name] = get(dbs, '0.db', 0)
        },

        /**
         * close connection
         * @param {string} name
         * @returns {Promise<boolean>}
         */
        async closeConnection(name) {
            const { success, msg } = await CloseConnection(name)
            if (!success) {
                // throw new Error(msg)
                return false
            }

            const dbs = this.databases[name]
            if (!isEmpty(dbs)) {
                for (const db of dbs) {
                    this._getNodeMap(name, db.db).clear()
                    this._getKeySet(name, db.db).clear()
                }
            }
            delete this.filter[name]
            delete this.databases[name]
            delete this.serverStats[name]

            const tabStore = useTabStore()
            tabStore.removeTabByName(name)
            return true
        },

        /**
         * open database and load all keys
         * @param server
         * @param db
         * @returns {Promise<void>}
         */
        async openDatabase(server, db) {
            const { match: filterPattern, type: filterType } = this.getKeyFilter(server)
            const { data, success, msg } = await OpenDatabase(server, db, filterPattern, filterType)
            if (!success) {
                throw new Error(msg)
            }
            const { keys = [], end = false, maxKeys = 0 } = data
            const selDB = this.getDatabase(server, db)
            if (selDB == null) {
                return
            }

            selDB.opened = true
            selDB.maxKeys = maxKeys
            set(this.loadingState, 'fullLoaded', end)
            if (isEmpty(keys)) {
                selDB.children = []
            } else {
                // append db node to current connection's children
                this._addKeyNodes(server, db, keys)
            }
            this._tidyNode(server, db)
        },

        /**
         * reopen database
         * @param connName
         * @param db
         * @returns {Promise<void>}
         */
        async reopenDatabase(connName, db) {
            const selDB = this.getDatabase(connName, db)
            if (selDB == null) {
                return
            }
            selDB.children = undefined
            selDB.isLeaf = false

            this._getNodeMap(connName, db).clear()
            this._getKeySet(connName, db).clear()
            delete this.filter[connName]
        },

        /**
         * close database
         * @param server
         * @param db
         */
        closeDatabase(server, db) {
            const selDB = this.getDatabase(server, db)
            if (selDB == null) {
                return
            }
            delete selDB.children
            selDB.isLeaf = false
            selDB.opened = false
            selDB.keys = 0

            this._getNodeMap(server, db).clear()
            this._getKeySet(server, db).clear()
            delete this.filter[server]
        },

        /**
         *
         * @param server
         * @returns {Promise<{}>}
         */
        async getServerInfo(server) {
            try {
                const { success, data } = await ServerInfo(server)
                if (success) {
                    this.serverStats[server] = data
                    return data
                }
            } finally {
            }
            return {}
        },

        /**
         * load key summary info
         * @param {string} server
         * @param {number} db
         * @param {string|number[]} [key] null or blank indicate that update tab to display normal content (blank content or server status)
         * @return {Promise<void>}
         */
        async loadKeySummary({ server, db, key }) {
            try {
                const tab = useTabStore()
                if (!isEmpty(key)) {
                    const { data, success, msg } = await GetKeySummary({
                        server,
                        db,
                        key,
                    })
                    if (success) {
                        const { type, ttl, size, length } = data
                        const k = decodeRedisKey(key)
                        const binaryKey = k !== key
                        tab.upsertTab({
                            subTab: BrowserTabType.KeyDetail,
                            server,
                            db,
                            type,
                            ttl,
                            keyCode: binaryKey ? key : undefined,
                            key: k,
                            size,
                            length,
                        })
                        return
                    } else {
                        if (!isEmpty(msg)) {
                            $message.error('load key summary fail: ' + msg)
                        }
                        // its danger to delete "non-exists" key, just remove from tree view
                        await this.deleteKey(server, db, key, true)
                        // TODO: show key not found page or check exists on server first?
                    }
                }

                tab.upsertTab({
                    subTab: BrowserTabType.Status,
                    server,
                    db,
                    type: 'none',
                    ttl: -1,
                    key: null,
                    keyCode: null,
                    value: null,
                    size: 0,
                    length: 0,
                })
            } catch (e) {
                $message.error('')
            } finally {
            }
        },

        /**
         * reload key
         * @param {string} server
         * @param {number} db
         * @param {string|number[]} key
         * @param {string} [decode]
         * @param {string} [format]
         * @param {string} [matchPattern]
         * @return {Promise<void>}
         */
        async reloadKey({ server, db, key, decode, format, matchPattern }) {
            const tab = useTabStore()
            try {
                tab.updateLoading({ server, db, loading: true })
                await this.loadKeySummary({ server, db, key })
                await this.loadKeyDetail({ server, db, key, decode, format, matchPattern, reset: true })
            } finally {
                tab.updateLoading({ server, db, loading: false })
            }
        },

        /**
         * load key content
         * @param {string} server
         * @param {number} db
         * @param {string|number[]} key
         * @param {string} [format]
         * @param {string} [decode]
         * @param {string} [matchPattern]
         * @param {boolean} [reset]
         * @param {boolean} [full]
         * @return {Promise<void>}
         */
        async loadKeyDetail({ server, db, key, format, decode, matchPattern, reset, full }) {
            const tab = useTabStore()
            try {
                tab.updateLoading({ server, db, loading: true })
                const { data, success, msg } = await GetKeyDetail({
                    server,
                    db,
                    key,
                    format,
                    decode,
                    matchPattern,
                    full: full === true,
                    reset,
                    lite: true,
                })
                if (success) {
                    const { value, decode: retDecode, format: retFormat, match: retMatch, reset: retReset, end } = data
                    tab.updateValue({
                        server,
                        db,
                        key: decodeRedisKey(key),
                        value,
                        decode: retDecode,
                        format: retFormat,
                        reset: retReset,
                        matchPattern: retMatch || '',
                        end,
                    })
                } else {
                    $message.error('load key detail fail:' + msg)
                }
            } finally {
                tab.updateLoading({ server, db, loading: false })
            }
        },

        /**
         * convert value by decode type or format
         * @param {string|number[]} value
         * @param {string} [decode]
         * @param {string} [format]
         * @return {Promise<{[format]: string, [decode]: string, value: string}>}
         */
        async convertValue({ value, decode, format }) {
            try {
                const { data, success } = await ConvertValue(value, decode, format)
                if (success) {
                    const { value: retVal, decode: retDecode, format: retFormat } = data
                    return { value: retVal, decode: retDecode, format: retFormat }
                }
            } catch (e) {}
            return { value, decode, format }
        },

        /**
         * scan keys with prefix
         * @param {string} connName
         * @param {number} db
         * @param {string} match
         * @param {string} [matchType]
         * @param {boolean} [full]
         * @returns {Promise<{keys: string[], maxKeys: number, end: boolean}>}
         */
        async scanKeys(connName, db, match, matchType, full) {
            let resp
            if (full) {
                resp = await LoadAllKeys(connName, db, match || '*', matchType)
            } else {
                resp = await LoadNextKeys(connName, db, match || '*', matchType)
            }
            const { data, success, msg } = resp || {}
            if (!success) {
                throw new Error(msg)
            }
            const { keys = [], maxKeys, end } = data
            return { keys, end, maxKeys, success }
        },

        /**
         *
         * @param {string} connName
         * @param {number} db
         * @param {string|null} prefix
         * @param {string|null} matchType
         * @param {boolean} [all]
         * @return {Promise<{keys: Array<string|number[]>, maxKeys: number, end: boolean}>}
         * @private
         */
        async _loadKeys(connName, db, prefix, matchType, all) {
            let match = prefix
            if (isEmpty(match)) {
                match = '*'
            } else if (!isRedisGlob(match)) {
                const separator = this._getSeparator(connName)
                if (!endsWith(prefix, separator + '*')) {
                    match = prefix + separator + '*'
                }
            }
            return this.scanKeys(connName, db, match, matchType, all)
        },

        /**
         * load more keys within the database
         * @param {string} connName
         * @param {number} db
         * @return {Promise<boolean>}
         */
        async loadMoreKeys(connName, db) {
            const { match, type: keyType } = this.getKeyFilter(connName)
            const { keys, maxKeys, end } = await this._loadKeys(connName, db, match, keyType, false)
            this._setDBMaxKeys(connName, db, maxKeys)
            // remove current keys below prefix
            this._addKeyNodes(connName, db, keys)
            this._tidyNode(connName, db, '')
            return end
        },

        /**
         * load all left keys within the database
         * @param {string} connName
         * @param {number} db
         * @return {Promise<void>}
         */
        async loadAllKeys(connName, db) {
            const { match, type: keyType } = this.getKeyFilter(connName, db)
            const { keys, maxKeys } = await this._loadKeys(connName, db, match, keyType, true)
            this._setDBMaxKeys(connName, db, maxKeys)
            this._addKeyNodes(connName, db, keys)
            this._tidyNode(connName, db, '')
        },

        /**
         * reload keys under layer
         * @param {string} connName
         * @param {number} db
         * @param {string} prefix
         * @return {Promise<void>}
         */
        async reloadLayer(connName, db, prefix) {
            if (isEmpty(prefix)) {
                return
            }
            let match = prefix
            const separator = this._getSeparator(connName)
            if (!endsWith(match, separator)) {
                match += separator + '*'
            } else {
                match += '*'
            }
            // FIXME: ignore original match pattern due to redis not support combination matching
            const { match: originMatch, type: keyType } = this.getKeyFilter(connName, db)
            const { keys, maxKeys, success } = await this._loadKeys(connName, db, match, keyType, true)
            if (!success) {
                return
            }

            this._setDBMaxKeys(connName, db, maxKeys)
            // remove current keys below prefix
            this._deleteKeyNode(connName, db, prefix, true)
            this._addKeyNodes(connName, db, keys)
            this._tidyNode(connName, db, prefix)
        },

        /**
         * get custom separator of connection
         * @param server
         * @returns {string}
         * @private
         */
        _getSeparator(server) {
            const connStore = useConnectionStore()
            const { keySeparator } = connStore.getDefaultSeparator(server)
            if (isEmpty(keySeparator)) {
                return ':'
            }
            return keySeparator
        },

        /**
         * get node map
         * @param {string} connName
         * @param {number} db
         * @returns {Map<string, DatabaseItem>}
         * @private
         */
        _getNodeMap(connName, db) {
            if (!this.nodeMap.hasOwnProperty(`${connName}#${db}`)) {
                this.nodeMap[`${connName}#${db}`] = new Map()
            }
            // construct a tree node list, the format of item key likes 'server/db#type/key'
            return this.nodeMap[`${connName}#${db}`]
        },

        /**
         * get all keys in a database
         * @param {string} connName
         * @param {number} db
         * @return {Set<string|number[]>}
         * @private
         */
        _getKeySet(connName, db) {
            if (!this.keySet.hasOwnProperty(`${connName}#${db}`)) {
                this.keySet[`${connName}#${db}`] = new Set()
            }
            // construct a key set
            return this.keySet[`${connName}#${db}`]
        },

        /**
         * remove keys in db
         * @param {string} connName
         * @param {number} db
         * @param {Array<string|number[]>|Set<string|number[]>} keys
         * @param {boolean} [sortInsert]
         * @return {{success: boolean, newKey: number, newLayer: number, replaceKey: number}}
         * @private
         */
        _addKeyNodes(connName, db, keys, sortInsert) {
            const result = {
                success: false,
                newLayer: 0,
                newKey: 0,
                replaceKey: 0,
            }
            if (isEmpty(keys)) {
                return result
            }
            const separator = this._getSeparator(connName)
            const selDB = this.getDatabase(connName, db)
            if (selDB == null) {
                return result
            }

            if (selDB.children == null) {
                selDB.children = []
            }
            const nodeMap = this._getNodeMap(connName, db)
            const keySet = this._getKeySet(connName, db)
            const rootChildren = selDB.children
            const viewType = get(this.viewType, connName, KeyViewType.Tree)
            if (viewType === KeyViewType.List) {
                // construct list view data
                for (const key of keys) {
                    const k = decodeRedisKey(key)
                    const isBinaryKey = k !== key
                    const nodeKey = `${ConnectionType.RedisValue}/${nativeRedisKey(key)}`
                    const replaceKey = nodeMap.has(nodeKey)
                    const selectedNode = {
                        key: `${connName}/db${db}#${nodeKey}`,
                        label: k,
                        db,
                        keys: 0,
                        redisKey: k,
                        redisKeyCode: isBinaryKey ? key : undefined,
                        type: ConnectionType.RedisValue,
                        isLeaf: true,
                    }
                    nodeMap.set(nodeKey, selectedNode)
                    keySet.add(key)
                    if (!replaceKey) {
                        if (sortInsert) {
                            const index = sortedIndexBy(rootChildren, selectedNode, 'key')
                            rootChildren.splice(index, 0, selectedNode)
                        } else {
                            rootChildren.push(selectedNode)
                        }
                        result.newKey += 1
                    } else {
                        result.replaceKey += 1
                    }
                }
            } else {
                // construct tree view data
                for (const key of keys) {
                    const k = decodeRedisKey(key)
                    const isBinaryKey = k !== key
                    const keyParts = isBinaryKey ? [nativeRedisKey(key)] : split(k, separator)
                    const len = size(keyParts)
                    const lastIdx = len - 1
                    let handlePath = ''
                    let children = rootChildren
                    for (let i = 0; i < len; i++) {
                        handlePath += keyParts[i]
                        if (i !== lastIdx) {
                            // layer
                            const nodeKey = `${ConnectionType.RedisKey}/${handlePath}`
                            let selectedNode = nodeMap.get(nodeKey)
                            if (selectedNode == null) {
                                selectedNode = {
                                    key: `${connName}/db${db}#${nodeKey}`,
                                    label: keyParts[i],
                                    db,
                                    keys: 0,
                                    redisKey: handlePath,
                                    type: ConnectionType.RedisKey,
                                    isLeaf: false,
                                    children: [],
                                }
                                nodeMap.set(nodeKey, selectedNode)
                                if (sortInsert) {
                                    const index = sortedIndexBy(children, selectedNode, 'key')
                                    children.splice(index, 0, selectedNode)
                                } else {
                                    children.push(selectedNode)
                                }
                                result.newLayer += 1
                            }
                            children = selectedNode.children
                            handlePath += separator
                        } else {
                            // key
                            const nodeKey = `${ConnectionType.RedisValue}/${handlePath}`
                            const replaceKey = nodeMap.has(nodeKey)
                            const selectedNode = {
                                key: `${connName}/db${db}#${nodeKey}`,
                                label: isBinaryKey ? k : keyParts[i],
                                db,
                                keys: 0,
                                redisKey: handlePath,
                                redisKeyCode: isBinaryKey ? key : undefined,
                                type: ConnectionType.RedisValue,
                                isLeaf: true,
                            }
                            nodeMap.set(nodeKey, selectedNode)
                            keySet.add(key)
                            if (!replaceKey) {
                                if (sortInsert) {
                                    const index = sortedIndexBy(children, selectedNode, 'key')
                                    children.splice(index, 0, selectedNode)
                                } else {
                                    children.push(selectedNode)
                                }
                                result.newKey += 1
                            } else {
                                result.replaceKey += 1
                            }
                        }
                    }
                }
            }
            return result
        },

        /**
         *
         * @param {DatabaseItem[]} nodeList
         * @private
         */
        _sortNodes(nodeList) {
            if (nodeList == null) {
                return
            }
            nodeList.sort((a, b) => {
                return a.key > b.key ? 1 : -1
            })
        },

        /**
         * tidy node by key
         * @param {string} connName
         * @param {number} db
         * @param {string} [key]
         * @param {boolean} [skipResort]
         * @private
         */
        _tidyNode(connName, db, key, skipResort) {
            const nodeMap = this._getNodeMap(connName, db)
            const dbNode = this.getDatabase(connName, db) || {}
            const separator = this._getSeparator(connName)
            const keyParts = split(key, separator)
            const totalParts = size(keyParts)
            let node
            // find last exists ancestor key
            let i = totalParts - 1
            for (; i > 0; i--) {
                const parentKey = join(slice(keyParts, 0, i), separator)
                node = nodeMap.get(`${ConnectionType.RedisKey}/${parentKey}`)
                if (node != null) {
                    break
                }
            }
            if (node == null) {
                node = dbNode
            }
            const keyCountUpdated = this._tidyNodeChildren(node, skipResort)

            if (keyCountUpdated) {
                // update key count of parent and above
                for (; i > 0; i--) {
                    const parentKey = join(slice(keyParts, 0, i), separator)
                    const parentNode = nodeMap.get(`${ConnectionType.RedisKey}/${parentKey}`)
                    if (parentNode == null) {
                        break
                    }
                    parentNode.keys = sumBy(parentNode.children, 'keys')
                }
                // update key count of db
                dbNode.keys = sumBy(dbNode.children, 'keys')
            }
            return true
        },

        /**
         * sort all node item's children and calculate keys count
         * @param {DatabaseItem} node
         * @param {boolean} skipSort skip sorting children
         * @returns {boolean} return whether key count changed
         * @private
         */
        _tidyNodeChildren(node, skipSort) {
            let count = 0
            if (!isEmpty(node.children)) {
                if (skipSort !== true) {
                    this._sortNodes(node.children)
                }

                for (const elem of node.children) {
                    this._tidyNodeChildren(elem, skipSort)
                    count += elem.keys
                }
            } else {
                if (node.type === ConnectionType.RedisValue) {
                    count += 1
                } else {
                    // no children in db node or layer node, set count to 0
                    count = 0
                }
            }
            if (node.keys !== count) {
                node.keys = count
                return true
            }
            return false
        },

        /**
         * update max key by increase/decrease value
         * @param {string} connName
         * @param {number} db
         * @param {number} [updateValue]
         * @private
         */
        _updateDBMaxKeys(connName, db, updateValue) {
            if (updateValue === undefined) {
                return
            }
            const database = this.getDatabase(connName, db)
            if (database != null) {
                const maxKeys = get(database, 'maxKeys', 0)
                database.maxKeys = Math.max(0, maxKeys + updateValue)
            }
        },

        /**
         * set db max keys value
         * @param {string} connName
         * @param {number} db
         * @param {number} maxKeys
         * @private
         */
        _setDBMaxKeys(connName, db, maxKeys) {
            const database = this.getDatabase(connName, db)
            if (database != null) {
                set(database, 'maxKeys', maxKeys)
            }
        },

        /**
         * get tree node by key name
         * @param key
         * @return {DatabaseItem|null}
         */
        getNode(key) {
            let idx = key.indexOf('#')
            if (idx < 0) {
                idx = size(key)
            }
            const dbPart = key.substring(0, idx)
            // parse server and db index
            const idx2 = dbPart.lastIndexOf('/db')
            if (idx2 < 0) {
                return null
            }
            const server = dbPart.substring(0, idx2)
            const db = parseInt(dbPart.substring(idx2 + 3))
            if (isNaN(db)) {
                return null
            }

            if (size(key) > idx + 1) {
                const keyPart = key.substring(idx + 1)
                // contains redis key
                const nodeMap = this._getNodeMap(server, db)
                return nodeMap.get(keyPart)
            } else {
                return this.getDatabase(server, db)
            }
        },

        /**
         * set redis key
         * @param {string} server
         * @param {number} db
         * @param {string|number[]} key
         * @param {string} keyType
         * @param {any} value
         * @param {number} ttl
         * @param {string} [format]
         * @param {string} [decode]
         * @returns {Promise<{[msg]: string, success: boolean, [nodeKey]: {string}}>}
         */
        async setKey({ server, db, key, keyType, value, ttl, format = formatTypes.RAW, decode = decodeTypes.NONE }) {
            try {
                const { data, success, msg } = await SetKeyValue({
                    server,
                    db,
                    key,
                    keyType,
                    value,
                    ttl,
                    format,
                    decode,
                })
                if (success) {
                    // const { value } = data
                    // update tree view data
                    const { newKey = 0 } = this._addKeyNodes(server, db, [key], true)
                    if (newKey > 0) {
                        this._tidyNode(server, db, key)
                        this._updateDBMaxKeys(server, db, newKey)
                    }
                    const tab = useTabStore()
                    tab.updateValue({ server, db, key, value })
                    return {
                        success,
                        nodeKey: `${server}/db${db}#${ConnectionType.RedisValue}/${key}`,
                        updatedValue: value,
                    }
                } else {
                    return { success, msg }
                }
            } catch (e) {
                return { success: false, msg: e.message }
            }
        },

        /**
         * update hash entry
         * when field is set, newField is null, delete field
         * when field is null, newField is set, add new field
         * when both field and newField are set, and field === newField, update field
         * when both field and newField are set, and field !== newField, delete field and add newField
         * @param {string} server
         * @param {number} db
         * @param {string|number[]} key
         * @param {string} field
         * @param {string} [newField]
         * @param {string} [value]
         * @param {decodeTypes} [decode]
         * @param {formatTypes} [format]
         * @param {decodeTypes} [retDecode]
         * @param {formatTypes} [retFormat]
         * @param {boolean} [refresh]
         * @param {number} [index] index for retrieve affect entries quickly
         * @returns {Promise<{[msg]: string, success: boolean, [updated]: {}}>}
         */
        async setHash({
            server,
            db,
            key,
            field,
            newField = '',
            value = '',
            decode = decodeTypes.NONE,
            format = formatTypes.RAW,
            retDecode,
            retFormat,
            index,
        }) {
            try {
                const { data, success, msg } = await SetHashValue({
                    server,
                    db,
                    key,
                    field,
                    newField,
                    value,
                    decode,
                    format,
                    retDecode,
                    retFormat,
                })
                if (success) {
                    /**
                     * @type {{updated: HashEntryItem[], removed: HashEntryItem[], updated: HashEntryItem[], replaced: HashReplaceItem[]}}
                     */
                    const { updated = [], removed = [], added = [], replaced = [] } = data
                    const tab = useTabStore()
                    if (!isEmpty(removed)) {
                        const removedKeys = map(removed, 'k')
                        tab.removeValueEntries({ server, db, key, type: 'hash', entries: removedKeys })
                    }
                    if (!isEmpty(updated)) {
                        tab.updateValueEntries({ server, db, key, type: 'hash', entries: updated })
                    }
                    if (!isEmpty(added)) {
                        tab.insertValueEntries({ server, db, key, type: 'hash', entries: added })
                    }
                    if (!isEmpty(replaced)) {
                        tab.replaceValueEntries({
                            server,
                            db,
                            key,
                            type: 'hash',
                            entries: replaced,
                            index: [index],
                        })
                    }
                    return { success, updated }
                } else {
                    return { success, msg }
                }
            } catch (e) {
                return { success: false, msg: e.message }
            }
        },

        /**
         * insert or update hash field item
         * @param {string} server
         * @param {number} db
         * @param {string|number[]} key
         * @param {number }action 0:ignore duplicated fields 1:overwrite duplicated fields
         * @param {string[]} fieldItems field1, value1, filed2, value2...
         * @returns {Promise<{[msg]: string, success: boolean, [updated]: [], [added]: []}>}
         */
        async addHashField(server, db, key, action, fieldItems) {
            try {
                const { data, success, msg } = await AddHashField(server, db, key, action, fieldItems)
                if (success) {
                    const { updated = [], added = [] } = data
                    const tab = useTabStore()
                    if (!isEmpty(updated)) {
                        tab.updateValueEntries({ server, db, key, type: 'hash', entries: updated })
                    }
                    if (!isEmpty(added)) {
                        tab.insertValueEntries({ server, db, key, type: 'hash', entries: added })
                    }
                    return { success, updated, added }
                } else {
                    return { success: false, msg }
                }
            } catch (e) {
                return { success: false, msg: e.message }
            }
        },

        /**
         * remove hash field
         * @param {string} server
         * @param {number} db
         * @param {string|number[]} key
         * @param {string} field
         * @returns {Promise<{[msg]: {}, success: boolean, [removed]: string[]}>}
         */
        async removeHashField(server, db, key, field) {
            try {
                const { data, success, msg } = await SetHashValue({ server, db, key, field, newField: '' })
                if (success) {
                    const { removed = [] } = data
                    // if (!isEmpty(removed)) {
                    //     const tab = useTabStore()
                    //     tab.removeValueEntries({ server, db, key, type: 'hash', entries: removed })
                    // }
                    return { success, removed }
                } else {
                    return { success, msg }
                }
            } catch (e) {
                return { success: false, msg: e.message }
            }
        },

        /**
         * insert list item
         * @param {string} connName
         * @param {number} db
         * @param {string|number[]} key
         * @param {int} action 0: push to head, 1: push to tail
         * @param {string[]}values
         * @returns {Promise<*|{msg, success: boolean}>}
         */
        async addListItem(connName, db, key, action, values) {
            try {
                return AddListItem(connName, db, key, action, values)
            } catch (e) {
                return { success: false, msg: e.message }
            }
        },

        /**
         * prepend item to head of list
         * @param {string} server
         * @param {number} db
         * @param {string|number[]} key
         * @param {string[]} values
         * @returns {Promise<{[msg]: string, success: boolean, [item]: []}>}
         */
        async prependListItem({ server, db, key, values }) {
            try {
                const { data, success, msg } = await AddListItem(server, db, key, 0, values)
                if (success) {
                    const { left = [] } = data
                    if (!isEmpty(left)) {
                        const tab = useTabStore()
                        tab.insertValueEntries({
                            server: server,
                            db,
                            key,
                            type: 'list',
                            entries: left,
                            prepend: true,
                        })
                    }
                    return { success, item: left }
                } else {
                    return { success: false, msg }
                }
            } catch (e) {
                return { success: false, msg: e.message }
            }
        },

        /**
         * append item to tail of list
         * @param {string} server
         * @param {number} db
         * @param {string|number[]} key
         * @param {string[]} values
         * @returns {Promise<{[msg]: string, success: boolean, [item]: any[]}>}
         */
        async appendListItem({ server, db, key, values }) {
            try {
                const { data, success, msg } = await AddListItem(server, db, key, 1, values)
                if (success) {
                    const { right = [] } = data
                    // FIXME: do not append items if not all items loaded
                    if (!isEmpty(right)) {
                        const tab = useTabStore()
                        tab.insertValueEntries({
                            server: server,
                            db,
                            key,
                            type: 'list',
                            entries: right,
                            prepend: false,
                        })
                    }
                    return { success, item: right }
                } else {
                    return { success: false, msg }
                }
            } catch (e) {
                return { success: false, msg: e.message }
            }
        },

        /**
         * update value of list item by index
         * @param {string} server
         * @param {number} db
         * @param {string|number[]} key
         * @param {number} index
         * @param {string|number[]} value
         * @param {decodeTypes} decode
         * @param {formatTypes} format
         * @param {decodeTypes} [retDecode]
         * @param {formatTypes} [retFormat]
         * @returns {Promise<{[msg]: string, success: boolean}>}
         */
        async updateListItem({
            server,
            db,
            key,
            index,
            value,
            decode = decodeTypes.NONE,
            format = formatTypes.RAW,
            retDecode,
            retFormat,
        }) {
            try {
                const { data, success, msg } = await SetListItem({
                    server,
                    db,
                    key,
                    index,
                    value,
                    decode,
                    format,
                    retDecode,
                    retFormat,
                })
                if (success) {
                    /** @type {{replaced: ListReplaceItem[]}} **/
                    const { replaced = [], removed = [] } = data
                    const tab = useTabStore()
                    if (!isEmpty(replaced)) {
                        tab.replaceValueEntries({
                            server,
                            db,
                            key,
                            type: 'list',
                            entries: replaced,
                        })
                    }
                    if (!isEmpty(removed)) {
                        const removedIndex = map(removed, 'index')
                        tab.removeValueEntries({
                            server,
                            db,
                            key,
                            type: 'list',
                            entries: removedIndex,
                        })
                    }
                    return { success }
                } else {
                    return { success, msg }
                }
            } catch (e) {
                return { success: false, msg: e.message }
            }
        },

        /**
         * remove list item
         * @param {string} server
         * @param {number} db
         * @param {string|number[]} key
         * @param {number} index
         * @returns {Promise<{[msg]: string, success: boolean, [removed]: string[]}>}
         */
        async removeListItem(server, db, key, index) {
            try {
                const { data, success, msg } = await SetListItem({ server, db, key, index })
                if (success) {
                    const { removed = [] } = data
                    const tab = useTabStore()
                    if (!isEmpty(removed)) {
                        const removedIndexes = map(removed, 'index')
                        tab.removeValueEntries({
                            server,
                            db,
                            key,
                            type: 'list',
                            entries: removedIndexes,
                        })
                    }
                    return { success, removed }
                } else {
                    return { success, msg }
                }
            } catch (e) {
                return { success: false, msg: e.message }
            }
        },

        /**
         * add item to set
         * @param {string} server
         * @param {number} db
         * @param {string|number} key
         * @param {string|string[]} value
         * @returns {Promise<{[msg]: string, success: boolean}>}
         */
        async addSetItem(server, db, key, value) {
            try {
                if ((!value) instanceof Array) {
                    value = [value]
                }
                const { data, success, msg } = await SetSetItem(server, db, key, false, value)
                if (success) {
                    const { added } = data
                    if (!isEmpty(added)) {
                        const tab = useTabStore()
                        tab.insertValueEntries({ server, db, key, type: 'set', entries: added })
                    }
                    return { success }
                } else {
                    return { success, msg }
                }
            } catch (e) {
                return { success: false, msg: e.message }
            }
        },

        /**
         * update value of set item
         * @param {string} server
         * @param {number} db
         * @param {string|number[]} key
         * @param {string|number[]} value
         * @param {string|number[]} newValue
         * @param {decodeTypes} [decode]
         * @param {formatTypes} [format]
         * @param {decodeTypes} [retDecode]
         * @param {formatTypes} [retFormat]
         * @returns {Promise<{[msg]: string, success: boolean}>}
         */
        async updateSetItem({
            server,
            db,
            key,
            value,
            newValue,
            decode = decodeTypes.NONE,
            format = formatTypes.RAW,
            retDecode,
            retFormat,
        }) {
            try {
                const { data, success, msg } = await UpdateSetItem({
                    server,
                    db,
                    key,
                    value,
                    newValue,
                    decode,
                    format,
                    retDecode,
                    retFormat,
                })
                if (success) {
                    const { added, removed } = data
                    const tab = useTabStore()
                    if (!isEmpty(removed)) {
                        const removedValues = map(removed, 'v')
                        tab.removeValueEntries({ server, db, key, type: 'set', entries: removedValues })
                    }
                    if (!isEmpty(added)) {
                        tab.insertValueEntries({ server, db, key, type: 'set', entries: added })
                    }
                    return { success }
                } else {
                    return { success: false, msg }
                }
            } catch (e) {
                return { success: false, msg: e.message }
            }
        },

        /**
         * remove item from set
         * @param {string} server
         * @param {number} db
         * @param {string|number[]} key
         * @param {string} value
         * @returns {Promise<{[msg]: string, success: boolean}>}
         */
        async removeSetItem(server, db, key, value) {
            try {
                const { data, success, msg } = await SetSetItem(server, db, key, true, [value])
                if (success) {
                    const { removed } = data
                    const tab = useTabStore()
                    if (!isEmpty(removed)) {
                        const removedValues = map(removed, 'v')
                        tab.removeValueEntries({ server, db, key, type: 'set', entries: removedValues })
                    }
                    return { success }
                } else {
                    return { success, msg }
                }
            } catch (e) {
                return { success: false, msg: e.message }
            }
        },

        /**
         * add item to sorted set
         * @param {string} server
         * @param {number} db
         * @param {string|number[]} key
         * @param {number} action
         * @param {Object.<string, number>} vs value: score
         * @returns {Promise<{[msg]: string, success: boolean}>}
         */
        async addZSetItem(server, db, key, action, vs) {
            try {
                const { data, success, msg } = await AddZSetValue(server, db, key, action, vs)
                if (success) {
                    const { added, updated } = data
                    const tab = useTabStore()
                    if (!isEmpty(added)) {
                        tab.insertValueEntries({ server, db, key, type: 'zset', entries: added })
                    }
                    if (!isEmpty(updated)) {
                        tab.updateValueEntries({ server, db, key, type: 'zset', entries: updated })
                    }
                    return { success }
                } else {
                    return { success, msg }
                }
            } catch (e) {
                return { success: false, msg: e.message }
            }
        },

        /**
         * update item of sorted set
         * @param {string} server
         * @param {number} db
         * @param {string|number[]} key
         * @param {string} value
         * @param {string} newValue
         * @param {number} score
         * @param {decodeTypes} decode
         * @param {formatTypes} format
         * @param {decodeTypes} [retDecode]
         * @param {formatTypes} [retFormat]
         * @param {number} [index] index for retrieve affect entries quickly
         * @returns {Promise<{[msg]: string, success: boolean}>}
         */
        async updateZSetItem({
            server,
            db,
            key,
            value = '',
            newValue,
            score,
            decode = decodeTypes.NONE,
            format = formatTypes.RAW,
            retDecode,
            retFormat,
            index,
        }) {
            try {
                const { data, success, msg } = await UpdateZSetValue({
                    server,
                    db,
                    key,
                    value,
                    newValue,
                    score,
                    decode,
                    format,
                    retDecode,
                    retFormat,
                })
                if (success) {
                    const { updated = [], added = [], removed = [], replaced = [] } = data
                    const tab = useTabStore()
                    if (!isEmpty(removed)) {
                        const removedValues = map(removed, 'v')
                        tab.removeValueEntries({ server, db, key, type: 'zset', entries: removedValues })
                    }
                    if (!isEmpty(updated)) {
                        tab.updateValueEntries({ server, db, key, type: 'zset', entries: updated })
                    }
                    if (!isEmpty(added)) {
                        tab.insertValueEntries({ server, db, key, type: 'zset', entries: added })
                    }
                    if (!isEmpty(replaced)) {
                        tab.replaceValueEntries({ server, db, key, type: 'zset', entries: replaced, index: [index] })
                    }
                    return { success, updated, removed }
                } else {
                    return { success, msg }
                }
            } catch (e) {
                return { success: false, msg: e.message }
            }
        },

        /**
         * remove item from sorted set
         * @param {string} server
         * @param {number} db
         * @param {string|number[]} key
         * @param {string} value
         * @returns {Promise<{[msg]: string, success: boolean, [removed]: []}>}
         */
        async removeZSetItem(server, db, key, value) {
            try {
                const { data, success, msg } = await UpdateZSetValue({ server, db, key, value, newValue: '', score: 0 })
                if (success) {
                    const { removed } = data
                    const tab = useTabStore()
                    if (!isEmpty(removed)) {
                        const removeValues = map(removed, 'v')
                        tab.removeValueEntries({ server, db, key, type: 'zset', entries: removeValues })
                    }
                    return { success, removed }
                } else {
                    return { success, msg }
                }
            } catch (e) {
                return { success: false, msg: e.message }
            }
        },

        /**
         * insert new stream field item
         * @param {string} server
         * @param {number} db
         * @param {string|number[]} key
         * @param {string} id
         * @param {string[]} values field1, value1, filed2, value2...
         * @returns {Promise<{[msg]: string, success: boolean}>}
         */
        async addStreamValue(server, db, key, id, values) {
            try {
                const { data = {}, success, msg } = await AddStreamValue(server, db, key, id, values)
                if (success) {
                    const { added = [] } = data
                    if (!isEmpty(added)) {
                        const tab = useTabStore()
                        tab.insertValueEntries({
                            server,
                            db,
                            key,
                            type: 'stream',
                            entries: added,
                        })
                    }
                    return { success }
                } else {
                    return { success: false, msg }
                }
            } catch (e) {
                return { success: false, msg: e.message }
            }
        },

        /**
         * remove stream field
         * @param {string} connName
         * @param {number} db
         * @param {string|number[]} key
         * @param {string[]|string} ids
         * @returns {Promise<{[msg]: {}, success: boolean}>}
         */
        async removeStreamValues(connName, db, key, ids) {
            if (typeof ids === 'string') {
                ids = [ids]
            }
            try {
                const { data = {}, success, msg } = await RemoveStreamValues(connName, db, key, ids)
                if (success) {
                    const tab = useTabStore()
                    tab.removeValueEntries({ server: connName, db, key, type: 'stream', entries: ids })
                    return { success }
                } else {
                    return { success, msg }
                }
            } catch (e) {
                return { success: false, msg: e.message }
            }
        },

        /**
         * reset key's ttl
         * @param {string} connName
         * @param {number} db
         * @param {string} key
         * @param {number} ttl
         * @returns {Promise<boolean>}
         */
        async setTTL(connName, db, key, ttl) {
            try {
                const { success, msg } = await SetKeyTTL(connName, db, key, ttl)
                return success === true
            } catch (e) {
                return false
            }
        },

        /**
         *
         * @param {string} connName
         * @param {number} db
         * @param {string} key
         * @param {string} newKey
         * @private
         */
        _renameKeyNode(connName, db, key, newKey) {
            const nodeMap = this._getNodeMap(connName, db)
            const nodeKey = `${ConnectionType.RedisValue}/${key}`
            const newNodeKey = `${ConnectionType.RedisValue}/${newKey}`
            const node = nodeMap.get(nodeKey)
            if (node != null) {
                // replace node map item
                const separator = this._getSeparator(connName)
                node.label = last(split(newKey, separator))
                node.key = `${connName}/db${db}#${newNodeKey}`
                node.redisKey = newKey
                nodeMap[newNodeKey] = node
                nodeMap.delete(nodeKey)
                // replace key set item
                const keySet = this._getKeySet(connName, db)
                keySet.delete(key)
                keySet.add(newKey)
            }
        },

        /**
         *
         * @param {string} connName
         * @param {number} db
         * @param {string} [key]
         * @param {boolean} [isLayer]
         * @private
         */
        _deleteKeyNode(connName, db, key, isLayer) {
            const dbRoot = this.getDatabase(connName, db) || {}
            const separator = this._getSeparator(connName)

            if (dbRoot == null) {
                return false
            }

            const nodeMap = this._getNodeMap(connName, db)
            const keySet = this._getKeySet(connName, db)
            if (isLayer === true) {
                this._deleteChildrenKeyNodes(nodeMap, keySet, key)
            }
            if (isEmpty(key)) {
                // clear all key nodes
                dbRoot.children = []
                dbRoot.keys = 0
            } else {
                const keyParts = split(key, separator)
                const totalParts = size(keyParts)
                // remove from parent in tree node
                const parentKey = slice(keyParts, 0, totalParts - 1)
                let parentNode
                if (isEmpty(parentKey)) {
                    parentNode = dbRoot
                } else {
                    parentNode = nodeMap.get(`${ConnectionType.RedisKey}/${join(parentKey, separator)}`)
                }

                // not found parent node
                if (parentNode == null) {
                    return false
                }
                remove(parentNode.children, {
                    type: isLayer ? ConnectionType.RedisKey : ConnectionType.RedisValue,
                    redisKey: key,
                })

                // check and remove empty layer node
                let i = totalParts - 1
                for (; i >= 0; i--) {
                    const anceKey = join(slice(keyParts, 0, i), separator)
                    if (i > 0) {
                        const anceNode = nodeMap.get(`${ConnectionType.RedisKey}/${anceKey}`)
                        const redisKey = join(slice(keyParts, 0, i + 1), separator)
                        remove(anceNode.children, { type: ConnectionType.RedisKey, redisKey })

                        if (isEmpty(anceNode.children)) {
                            nodeMap.delete(`${ConnectionType.RedisKey}/${anceKey}`)
                            keySet.delete(anceNode.redisKeyCode || anceNode.redisKey)
                        } else {
                            break
                        }
                    } else {
                        // last one, remove from db node
                        remove(dbRoot.children, { type: ConnectionType.RedisKey, redisKey: keyParts[0] })
                        const node = nodeMap.get(`${ConnectionType.RedisValue}/${keyParts[0]}`)
                        if (node != null) {
                            nodeMap.delete(`${ConnectionType.RedisValue}/${keyParts[0]}`)
                            keySet.delete(node.redisKeyCode || node.redisKey)
                        }
                    }
                }
            }

            return true
        },

        /**
         * delete node and all it's children from nodeMap
         * @param {Map<string, DatabaseItem>} nodeMap
         * @param {Set<string|number[]>} keySet
         * @param {string} [key] clean nodeMap if key is empty
         * @private
         */
        _deleteChildrenKeyNodes(nodeMap, keySet, key) {
            if (isEmpty(key)) {
                nodeMap.clear()
                keySet.clear()
            } else {
                const mapKey = `${ConnectionType.RedisKey}/${key}`
                const node = nodeMap.get(mapKey)
                for (const child of node.children || []) {
                    if (child.type === ConnectionType.RedisValue) {
                        if (!nodeMap.delete(`${ConnectionType.RedisValue}/${child.redisKey}`)) {
                            console.warn('delete:', `${ConnectionType.RedisValue}/${child.redisKey}`)
                        }
                        keySet.delete(child.redisKeyCode || child.redisKey)
                    } else if (child.type === ConnectionType.RedisKey) {
                        this._deleteChildrenKeyNodes(nodeMap, keySet, child.redisKey)
                    }
                }
                if (!nodeMap.delete(mapKey)) {
                    console.warn('delete map key', mapKey)
                }
                keySet.delete(node.redisKeyCode || node.redisKey)
            }
        },

        /**
         * delete redis key
         * @param {string} connName
         * @param {number} db
         * @param {string|number[]} key
         * @param {boolean} [soft] do not try to remove from redis if true, just remove from tree data
         * @returns {Promise<boolean>}
         */
        async deleteKey(connName, db, key, soft) {
            try {
                let deleteCount = 0
                if (soft !== true) {
                    const { data } = await DeleteKey(connName, db, key)
                    deleteCount = get(data, 'deleteCount', 0)
                }

                const k = nativeRedisKey(key)
                // update tree view data
                this._deleteKeyNode(connName, db, k)
                this._tidyNode(connName, db, k, true)
                this._updateDBMaxKeys(connName, db, -deleteCount)

                // set tab content empty
                const tab = useTabStore()
                tab.emptyTab(connName)
                return true
            } finally {
            }
            return false
        },

        /**
         * delete keys with prefix
         * @param {string} connName
         * @param {number} db
         * @param {string} prefix
         * @param {boolean} async
         * @returns {Promise<boolean>}
         */
        async deleteKeyPrefix(connName, db, prefix, async) {
            if (isEmpty(prefix)) {
                return false
            }
            try {
                if (!endsWith(prefix, '*')) {
                    prefix += '*'
                }
                const { data, success, msg } = await DeleteKey(connName, db, prefix, async)
                if (success) {
                    // const { deleted: keys = [] } = data
                    // for (const key of keys) {
                    //     await this._deleteKeyNode(connName, db, key)
                    // }
                    const deleteCount = get(data, 'deleteCount', 0)
                    const separator = this._getSeparator(connName)
                    if (endsWith(prefix, '*')) {
                        prefix = prefix.substring(0, prefix.length - 1)
                    }
                    if (endsWith(prefix, separator)) {
                        prefix = prefix.substring(0, prefix.length - 1)
                    }
                    this._deleteKeyNode(connName, db, prefix, true)
                    this._tidyNode(connName, db, prefix, true)
                    this._updateDBMaxKeys(connName, db, -deleteCount)
                    return true
                }
            } finally {
            }
            return false
        },

        /**
         * flush database
         * @param connName
         * @param db
         * @param async
         * @return {Promise<boolean>}
         */
        async flushDatabase(connName, db, async) {
            try {
                const { success = false } = await FlushDB(connName, db, async)

                if (success === true) {
                    // update tree view data
                    this._deleteKeyNode(connName, db)
                    this._setDBMaxKeys(connName, db, 0)
                    // set tab content empty
                    const tab = useTabStore()
                    tab.emptyTab(connName)
                    return true
                }
            } finally {
            }
            return true
        },

        /**
         * rename key
         * @param {string} connName
         * @param {number} db
         * @param {string} key
         * @param {string} newKey
         * @returns {Promise<{[msg]: string, success: boolean, [nodeKey]: string}>}
         */
        async renameKey(connName, db, key, newKey) {
            const { success = false, msg } = await RenameKey(connName, db, key, newKey)
            if (success) {
                // delete old key and add new key struct
                this._renameKeyNode(connName, db, key, newKey)
                return { success: true, nodeKey: `${connName}/db${db}#${ConnectionType.RedisValue}/${newKey}` }
            } else {
                return { success: false, msg }
            }
        },

        /**
         * get command history
         * @param {number} [pageNo]
         * @param {number} [pageSize]
         * @returns {Promise<HistoryItem[]>}
         */
        async getCmdHistory(pageNo, pageSize) {
            if (pageNo === undefined || pageSize === undefined) {
                pageNo = -1
                pageSize = -1
            }
            try {
                const { success, data = { list: [] } } = await GetCmdHistory(pageNo, pageSize)
                const { list } = data
                return list
            } catch {
                return []
            }
        },

        /**
         * clean cmd history
         * @return {Promise<boolean>}
         */
        async cleanCmdHistory() {
            try {
                const { success } = await CleanCmdHistory()
                return success === true
            } catch {
                return false
            }
        },

        /**
         * get slow log list
         * @param {string} server
         * @param {number} db
         * @param {number} num
         * @return {Promise<[]>}
         */
        async getSlowLog(server, db, num) {
            try {
                const { success, data = { list: [] } } = await GetSlowLogs(server, db, num)
                const { list } = data
                return list
            } catch {
                return []
            }
        },

        /**
         * get key filter pattern and filter type
         * @param {string} server
         * @returns {{match: string, type: string}}
         */
        getKeyFilter(server) {
            let { pattern = '', type = '' } = this.filter[server] || {}
            if (isEmpty(pattern)) {
                // no custom match pattern, use default
                const conn = useConnectionStore()
                pattern = conn.getDefaultKeyFilter(server)
            }
            return {
                match: pattern,
                type: toUpper(type),
            }
        },

        /**
         *
         * @param {string} server
         * @param {string} [type]
         * @param {string} [pattern]
         */
        setKeyFilter(server, { type, pattern }) {
            const filter = this.filter[server] || {}
            filter.type = type === null ? filter.type : type
            filter.pattern = type === null ? filter.pattern : pattern
            this.filter[server] = filter
        },
    },
})

export default useBrowserStore
