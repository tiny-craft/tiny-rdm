import { defineStore } from 'pinia'
import { endsWith, get, isEmpty, join, map, now, size, slice, split } from 'lodash'
import {
    AddHashField,
    AddListItem,
    AddStreamValue,
    AddZSetValue,
    BatchSetTTL,
    CleanCmdHistory,
    CloseConnection,
    ConvertValue,
    DeleteKey,
    DeleteKeys,
    DeleteKeysByPattern,
    ExportKey,
    FlushDB,
    GetClientList,
    GetCmdHistory,
    GetHashValue,
    GetKeyDetail,
    GetKeySummary,
    GetKeyType,
    GetSlowLogs,
    ImportCSV,
    LoadAllKeys,
    LoadNextAllKeys,
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
import { nativeRedisKey } from '@/utils/key_convert.js'
import { BrowserTabType } from '@/consts/browser_tab_type.js'
import { KeyViewType } from '@/consts/key_view_type.js'
import { ConnectionType } from '@/consts/connection_type.js'
import useConnectionStore from 'stores/connections.js'
import { decodeTypes, formatTypes } from '@/consts/value_view_type.js'
import { isRedisGlob } from '@/utils/glob_pattern.js'
import { i18nGlobal } from '@/utils/i18n.js'
import { EventsEmit, EventsOn } from 'wailsjs/runtime/runtime.js'
import { RedisNodeItem } from '@/objects/redisNodeItem.js'
import { RedisServerState } from '@/objects/redisServerState.js'
import { RedisDatabaseItem } from '@/objects/redisDatabaseItem.js'
import { timeout } from '@/utils/promise.js'

const useBrowserStore = defineStore('browser', {
    /**
     * @typedef {Object} FilterItem
     * @property {string} pattern key pattern filter
     * @property {string} type type filter
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
     * @property {Object.<string, RedisServerState>} servers
     */

    /**
     *
     * @returns {BrowserState}
     */
    state: () => ({
        servers: {},
    }),
    getters: {
        anyConnectionOpened() {
            return !isEmpty(this.servers)
        },
    },
    actions: {
        /**
         * check if connection is connected
         * @param name
         * @returns {boolean}
         */
        isConnected(name) {
            return this.servers.hasOwnProperty(name)
        },

        /**
         * close all connections
         * @returns {Promise<void>}
         */
        async closeAllConnection() {
            for (const serverName in this.servers) {
                await CloseConnection(serverName)
                this.servers[serverName].dispose()
            }

            const tabStore = useTabStore()
            tabStore.removeAllTab()
        },

        /**
         * get database info list
         * @param {string} server
         * @return {RedisDatabaseItem[]}
         */
        getDBList(server) {
            const serverInst = this.servers[server]
            if (serverInst != null) {
                return serverInst.getDatabase()
            }
            return []
        },

        /**
         * get server version
         * @param {string} server
         */
        getServerVersion(server) {
            const serverInst = this.servers[server]
            if (serverInst != null) {
                return serverInst.version
            }
            return '1.0.0'
        },

        /**
         * get database by server name and database index
         * @param {string} server
         * @param {number} db
         * @return {RedisDatabaseItem|null}
         */
        getDatabase(server, db) {
            /** @type {RedisServerState} **/
            const serverInst = this.servers[server]
            if (serverInst != null) {
                return serverInst.databases[db] || null
            }
            return null
        },

        /**
         * get current selection database by server
         * @param server
         * @return {number}
         */
        getSelectedDB(server) {
            /** @type {RedisServerState} **/
            const serverInst = this.servers[server]
            if (serverInst != null) {
                return serverInst.db
            }
            return 0
        },

        /**
         * get key struct in current database
         * @param {string} server
         * @param {boolean} [includeRoot]
         * @return {RedisNodeItem[]}
         */
        getKeyStruct(server, includeRoot) {
            /** @type {RedisServerState} **/
            const serverInst = this.servers[server]
            let rootNode = null
            if (serverInst != null) {
                rootNode = serverInst.getRoot()
            }
            if (includeRoot === true) {
                return [rootNode]
            }
            return get(rootNode, 'children', [])
        },

        getReloadKey(server) {
            /** @type {RedisServerState} **/
            const serverInst = this.servers[server]
            return serverInst != null ? serverInst.reloadKey : 0
        },

        reloadServer(server) {
            /** @type {RedisServerState} **/
            const serverInst = this.servers[server]
            if (serverInst != null) {
                serverInst.reloadKey = Date.now()
            }
        },

        /**
         * switch key view
         * @param {string} connName
         * @param {number} viewType
         */
        // async switchKeyView(connName, viewType) {
        //     if (viewType !== KeyViewType.Tree && viewType !== KeyViewType.List) {
        //         return
        //     }
        //
        //     const t = get(this.viewType, connName, KeyViewType.Tree)
        //     if (t === viewType) {
        //         return
        //     }
        //
        //     this.viewType[connName] = viewType
        //     const dbs = get(this.databases, connName, [])
        //     for (const dbItem of dbs) {
        //         if (!dbItem.opened) {
        //             continue
        //         }
        //
        //         dbItem.children = undefined
        //         dbItem.keyCount = 0
        //         const { db = 0 } = dbItem
        //         this._getNodeMap(connName, db).clear()
        //         this._addKeyNodes(connName, db, keys)
        //         this._tidyNode(connName, db, '')
        //     }
        // },

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
            const { db, view = KeyViewType.Tree, lastDB, version } = data
            if (isEmpty(db)) {
                throw new Error('no db loaded')
            }
            const serverInst = new RedisServerState({
                name,
                separator: this.getSeparator(name),
                db: -1,
                viewType: view,
                version,
            })
            /** @type {Object.<number,RedisDatabaseItem>} **/
            const databases = {}
            for (const dbItem of db) {
                databases[dbItem.index] = new RedisDatabaseItem({
                    db: dbItem.index,
                    alias: dbItem.alias,
                    maxKeys: dbItem.maxKeys,
                })
                if (dbItem.index === lastDB) {
                    // set last opened database as default
                    serverInst.db = dbItem.index
                } else if (serverInst.db === -1) {
                    // set the first database as default
                    serverInst.db = dbItem.index
                }
            }
            serverInst.databases = databases
            this.servers[name] = serverInst
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
            delete this.servers[name]

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
            const { data, success, msg } = await OpenDatabase(server, db)
            if (!success) {
                throw new Error(msg)
            }
            const { keys = [], end = false, maxKeys = 0 } = data

            /** @type {RedisServerState} **/
            const serverInst = this.servers[server]
            if (serverInst == null) {
                return
            }
            serverInst.db = db
            serverInst.setDatabaseKeyCount(db, maxKeys)
            serverInst.loadingState.fullLoaded = end

            if (isEmpty(keys)) {
                serverInst.nodeMap.clear()
            } else {
                // append db node to current connection's children
                serverInst.addKeyNodes(keys)
            }
            serverInst.tidyNode('', false)
        },

        /**
         * close database
         * @param server
         * @param db
         */
        closeDatabase(server, db) {
            /** @type {RedisServerState} **/
            const serverInst = this.servers[server]
            if (serverInst == null) {
                return
            }
            if (serverInst.db !== db) {
                return
            }
            serverInst.closeDatabase()

            /** @type {RedisDatabaseItem} **/
            const selDB = this.getDatabase(server, db)
            if (selDB == null) {
                return
            }
            selDB.keyCount = 0
        },

        /**
         *
         * @param {string} server
         * @param {boolean} mute
         * @returns {Promise<{}>}
         */
        async getServerInfo(server, mute) {
            try {
                const { success, data, msg } = await ServerInfo(server)
                if (success) {
                    /** @type {RedisServerState} **/
                    const serverInst = this.servers[server]
                    if (serverInst != null) {
                        serverInst.stats = data
                    }
                    return data
                } else if (!isEmpty(msg) && mute !== true) {
                    $message.warning(msg)
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
         * @param {boolean} [clearValue]
         * @param {boolean} [redirect] redirect to key detail tab
         * @return {Promise<void>}
         */
        async loadKeySummary({ server, db, key, clearValue, redirect = true }) {
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
                        const k = nativeRedisKey(key)
                        const binaryKey = k !== key
                        tab.upsertTab({
                            subTab: redirect === false ? null : BrowserTabType.KeyDetail,
                            server,
                            db,
                            type,
                            ttl,
                            keyCode: binaryKey ? key : undefined,
                            key: k,
                            size,
                            length,
                            clearValue,
                        })
                        return
                    } else {
                        if (!isEmpty(msg)) {
                            $message.error('load key summary fail: ' + msg)
                        }
                        // its danger to delete "non-exists" key, just remove from tree view
                        // await this.deleteKey(server, db, key, true)
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
                    size: 0,
                    length: 0,
                    clearValue,
                })
            } catch (e) {
                $message.error(e.message || 'unknown error')
            } finally {
            }
        },

        /**
         * load key type
         * @param {string} server
         * @param {number} db
         * @param {string|number[]} key
         * @param {number[]} keyCode
         * @return {Promise<void>}
         */
        async loadKeyType({ server, db, key }) {
            /** @type {RedisServerState} **/
            const serverInst = this.servers[server]
            if (serverInst == null) {
                return
            }
            const node = serverInst.getNode(ConnectionType.RedisValue, nativeRedisKey(key))
            if (node == null || !isEmpty(node.redisType)) {
                return
            }
            try {
                node.redisType = 'loading'
                const { data, success, msg } = await GetKeyType({ server, db, key })
                if (success) {
                    const { type } = data || {}
                    node.redisType = type
                } else {
                    node.redisType = 'NONE'
                }
            } catch (e) {
                node.redisType = 'NONE'
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
         * @param {boolean} [showLoading]
         * @return {Promise<void>}
         */
        async reloadKey({ server, db, key, decode, format, matchPattern, showLoading = true }) {
            const tab = useTabStore()
            try {
                if (showLoading) {
                    tab.updateLoading({ server, db, loading: true })
                }
                await this.loadKeySummary({ server, db, key, clearValue: true, redirect: false })
                await this.loadKeyDetail({
                    server,
                    db,
                    key,
                    decode,
                    format,
                    matchPattern,
                    reset: true,
                    showLoading: false,
                })
            } finally {
                if (showLoading) {
                    tab.updateLoading({ server, db, loading: false })
                }
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
         * @param {boolean} [showLoading]
         * @return {Promise<void>}
         */
        async loadKeyDetail({ server, db, key, format, decode, matchPattern, reset, full, showLoading = true }) {
            const tab = useTabStore()
            const serverInst = this.servers[server]
            if (serverInst == null) {
                return
            }
            try {
                if (showLoading) {
                    tab.updateLoading({ server, db, loading: true })
                }
                const [storeFormat, storeDecode] = serverInst.getDecodeHistory(key, db)
                const { data, success, msg } = await GetKeyDetail({
                    server,
                    db,
                    key,
                    format: isEmpty(format) ? storeFormat : format,
                    decode: isEmpty(decode) ? storeDecode : decode,
                    matchPattern,
                    full: full === true,
                    reset,
                    lite: true,
                })
                if (success) {
                    const {
                        value,
                        keyType,
                        decode: retDecode,
                        format: retFormat,
                        match: retMatch,
                        reset: retReset,
                        end,
                    } = data
                    tab.updateValue({
                        server,
                        db,
                        key: nativeRedisKey(key),
                        value,
                        decode: retDecode || storeDecode,
                        format: retFormat || storeFormat,
                        reset: retReset,
                        matchPattern: retMatch || '',
                        end,
                    })
                } else {
                    $message.error('load key detail fail:' + msg)
                }
            } finally {
                if (showLoading) {
                    tab.updateLoading({ server, db, loading: false })
                }
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
         * @param {string} server
         * @param {number} db
         * @param {string} match
         * @param {boolean} exact
         * @param {string} [matchType]
         * @param {number} [loadType] 0.load next; 1.load next full; 2.reload load all
         * @returns {Promise<{keys: string[], maxKeys: number, end: boolean}>}
         */
        async scanKeys({ server, db, match = '*', exact = false, matchType = '', loadType = 0 }) {
            let resp
            switch (loadType) {
                case 0:
                default:
                    resp = await LoadNextKeys(server, db, match, matchType, exact)
                    break
                case 1:
                    resp = await LoadNextAllKeys(server, db, match, matchType, exact)
                    break
                case 2:
                    resp = await LoadAllKeys(server, db, match, matchType, exact)
                    break
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
         * @param {string} server
         * @param {number} db
         * @param {string|null} match
         * @param {boolean} exact
         * @param {string|null} matchType
         * @param {boolean} [all]
         * @return {Promise<{keys: Array<string|number[]>, maxKeys: number, end: boolean}>}
         * @private
         */
        async _loadKeys({ server, db, match, exact, matchType, all }) {
            if (isEmpty(match)) {
                match = '*'
            }

            if (!isRedisGlob(match) && !exact) {
                if (!endsWith(match, '*')) {
                    match = match + '*'
                }
            }
            return this.scanKeys({ server, db, match, exact, matchType, loadType: all ? 1 : 0 })
        },

        /**
         * load more keys within the database
         * @param {string} server
         * @param {number} db
         * @return {Promise<boolean>}
         */
        async loadMoreKeys(server, db) {
            const { match, type: keyType, exact } = this.getKeyFilter(server)
            const { keys, maxKeys, end } = await this._loadKeys({
                server,
                db,
                match,
                exact,
                matchType: keyType,
                all: false,
            })
            /** @type RedisServerState **/
            const serverInst = this.servers[server]
            if (serverInst != null) {
                serverInst.setDBKeyCount(db, maxKeys)
                // remove current keys below prefix
                serverInst.addKeyNodes(keys)
                serverInst.tidyNode('')
            }
            return end
        },

        /**
         * load all left keys within the database
         * @param {string} server
         * @param {number} db
         * @return {Promise<void>}
         */
        async loadAllKeys(server, db) {
            const { match, type: keyType, exact } = this.getKeyFilter(server)
            const { keys, maxKeys } = await this._loadKeys({ server, db, match, exact, matchType: keyType, all: true })
            /** @type RedisServerState **/
            const serverInst = this.servers[server]
            if (serverInst != null) {
                serverInst.setDBKeyCount(db, maxKeys)
                serverInst.addKeyNodes(keys)
                serverInst.tidyNode('')
            }
        },

        /**
         * reload keys under layer
         * @param {string} server
         * @param {number} db
         * @param {string} prefix
         * @return {Promise<void>}
         */
        async reloadLayer(server, db, prefix) {
            if (isEmpty(prefix)) {
                return
            }
            let match = prefix
            const separator = this.getSeparator(server)
            if (!isEmpty(match)) {
                if (!endsWith(match, separator)) {
                    match += separator + '*'
                } else {
                    match += '*'
                }
            }
            // FIXME: ignore original match pattern due to redis not support combination matching
            const { match: originMatch, type: keyType, exact } = this.getKeyFilter(server)
            const { keys, maxKeys, success } = await this._loadKeys({
                server,
                db,
                match: match || originMatch,
                exact: false,
                matchType: keyType,
                all: true,
            })
            if (!success) {
                return
            }

            /** @type RedisServerState **/
            const serverInst = this.servers[server]
            if (serverInst != null) {
                serverInst.setDBKeyCount(db, maxKeys)
                // remove current keys below prefix
                serverInst.removeKeyNode(prefix, true)
                serverInst.addKeyNodes(keys)
                serverInst.tidyNode(prefix)
            }
        },

        /**
         * get custom separator of connection
         * @param server
         * @returns {string}
         * @private
         */
        getSeparator(server) {
            const connStore = useConnectionStore()
            const { keySeparator } = connStore.getDefaultSeparator(server)
            if (isEmpty(keySeparator)) {
                return ':'
            }
            return keySeparator
        },

        /**
         * get tree node by key name
         * @param key
         * @return {RedisNodeItem|null}
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
            /** @type {RedisServerState} **/
            const serverInst = this.servers[server]
            if (serverInst == null) {
                return null
            }

            const db = parseInt(dbPart.substring(idx2 + 3))
            if (isNaN(db)) {
                return null
            }

            if (size(key) <= idx + 1) {
                return null
            }
            // contains redis key
            const keyPart = key.substring(idx + 1)
            return serverInst.nodeMap.get(keyPart)
        },

        /**
         * get parent tree node by key name
         * @param key
         * @return {RedisNodeItem|null}
         */
        getParentNode(key) {
            const i = key.indexOf('#')
            if (i < 0) {
                return null
            }
            const [server, db] = split(key.substring(0, i), '/')
            if (isEmpty(server) || isEmpty(db)) {
                return null
            }
            /** @type {RedisServerState} **/
            const serverInst = this.servers[server]
            if (serverInst == null) {
                return null
            }
            const separator = this.getSeparator(server)
            const keyPart = key.substring(i)
            const keyStartIdx = keyPart.indexOf('/')
            const redisKey = keyPart.substring(keyStartIdx + 1)
            const redisKeyParts = split(redisKey, separator)
            const parentKey = slice(redisKeyParts, 0, size(redisKeyParts) - 1)
            if (isEmpty(parentKey)) {
                return serverInst.getRoot()
            }
            return serverInst.nodeMap.get(`${ConnectionType.RedisKey}/${join(parentKey, separator)}`)
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
                    /** @type RedisServerState **/
                    const serverInst = this.servers[server]
                    if (serverInst != null && serverInst.db === db) {
                        // const { value } = data
                        // update tree view data
                        const { newKey = 0 } = serverInst.addKeyNodes([key], true)
                        if (newKey > 0) {
                            serverInst.tidyNode(key)
                            serverInst.updateDBKeyCount(db, newKey)
                        }

                        const { value: updatedValue } = data
                        if (updatedValue != null) {
                            const tab = useTabStore()
                            tab.updateValue({ server, db, key, value: updatedValue })
                        }
                    }
                    // this.loadKeySummary({ server, db, key })
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
         * @param {boolean} [reload]
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
            reload,
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
                    if (reload === true) {
                        this.reloadKey({ server, db, key })
                    } else {
                        // reload summary only
                        this.loadKeySummary({ server, db, key })
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
         * @param {boolean} [reload]
         * @returns {Promise<{[msg]: string, success: boolean, [updated]: [], [added]: []}>}
         */
        async addHashField({ server, db, key, action, fieldItems, reload }) {
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
                    if (reload === true) {
                        this.reloadKey({ server, db, key })
                    } else {
                        // reload summary only
                        this.loadKeySummary({ server, db, key })
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
         * get hash field
         * @param {string} server
         * @param {number} db
         * @param {string} key
         * @param {string} field
         * @param {decodeTypes} [decode]
         * @param {formatTypes} [format]
         * @return {Promise<{{msg: string, success: boolean, updated: HashEntryItem[]}>}
         */
        async getHashField({ server, db, key, field, decode = decodeTypes.NONE, format = formatTypes.RAW }) {
            try {
                const { data, success, msg } = await GetHashValue({ server, db, key, field, decode, format })
                if (success && !isEmpty(data)) {
                    const tab = useTabStore()
                    tab.updateValueEntries({ server, db, key, type: 'hash', entries: [data] })
                    return { success, updated: data }
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
         * @param {boolean} reload
         * @returns {Promise<{[msg]: {}, success: boolean, [removed]: string[]}>}
         */
        async removeHashField({ server, db, key, field, reload }) {
            try {
                const { data, success, msg } = await SetHashValue({ server, db, key, field, newField: '' })
                if (success) {
                    const { removed = [] } = data
                    // if (!isEmpty(removed)) {
                    //     const tab = useTabStore()
                    //     tab.removeValueEntries({ server, db, key, type: 'hash', entries: removed })
                    // }
                    if (reload === true) {
                        this.reloadKey({ server, db, key })
                    } else {
                        // reload summary only
                        this.loadKeySummary({ server, db, key })
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
         * prepend item to head of list
         * @param {string} server
         * @param {number} db
         * @param {string|number[]} key
         * @param {string[]} values
         * @param {boolean} reload
         * @returns {Promise<{[msg]: string, success: boolean, [item]: []}>}
         */
        async prependListItem({ server, db, key, values, reload }) {
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
                        if (reload === true) {
                            this.reloadKey({ server, db, key })
                        } else {
                            // reload summary only
                            this.loadKeySummary({ server, db, key })
                        }
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
         * @param {boolean} [reload]
         * @returns {Promise<{[msg]: string, success: boolean, [item]: any[]}>}
         */
        async appendListItem({ server, db, key, values, reload }) {
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
                        if (reload === true) {
                            this.reloadKey({ server, db, key })
                        } else {
                            // reload summary only
                            this.loadKeySummary({ server, db, key })
                        }
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
         * @param {boolean} [reload]
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
            reload,
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
                    if (reload === true) {
                        this.reloadKey({ server, db, key })
                    } else {
                        // reload summary only
                        this.loadKeySummary({ server, db, key })
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
         * @param {boolean} [reload]
         * @returns {Promise<{[msg]: string, success: boolean, [removed]: string[]}>}
         */
        async removeListItem({ server, db, key, index, reload }) {
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
                        if (reload === true) {
                            this.reloadKey({ server, db, key })
                        } else {
                            // reload summary only
                            this.loadKeySummary({ server, db, key })
                        }
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
         * @param {boolean} [reload]
         * @returns {Promise<{[msg]: string, success: boolean}>}
         */
        async addSetItem({ server, db, key, value, reload }) {
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
                    if (reload === true) {
                        this.reloadKey({ server, db, key })
                    } else {
                        // reload summary only
                        this.loadKeySummary({ server, db, key })
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
         * @param {boolean} [reload]
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
            reload,
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
                    if (reload === true) {
                        this.reloadKey({ server, db, key })
                    } else {
                        // reload summary only
                        this.loadKeySummary({ server, db, key })
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
         * @param {boolean} [reload]
         * @returns {Promise<{[msg]: string, success: boolean}>}
         */
        async removeSetItem({ server, db, key, value, reload }) {
            try {
                const { data, success, msg } = await SetSetItem(server, db, key, true, [value])
                if (success) {
                    const { removed } = data
                    const tab = useTabStore()
                    if (!isEmpty(removed)) {
                        const removedValues = map(removed, 'v')
                        tab.removeValueEntries({ server, db, key, type: 'set', entries: removedValues })
                    }
                    if (reload === true) {
                        this.reloadKey({ server, db, key })
                    } else {
                        // reload summary only
                        this.loadKeySummary({ server, db, key })
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
         * @param {boolean} [reload]
         * @returns {Promise<{[msg]: string, success: boolean}>}
         */
        async addZSetItem({ server, db, key, action, vs, reload }) {
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
                    if (reload === true) {
                        this.reloadKey({ server, db, key })
                    } else {
                        // reload summary only
                        this.loadKeySummary({ server, db, key })
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
         * @param {boolean} [reload]
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
            reload,
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
                    if (reload === true) {
                        this.reloadKey({ server, db, key })
                    } else {
                        // reload summary only
                        this.loadKeySummary({ server, db, key })
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
         * @param {boolean} [reload]
         * @returns {Promise<{[msg]: string, success: boolean, [removed]: []}>}
         */
        async removeZSetItem({ server, db, key, value, reload }) {
            try {
                const { data, success, msg } = await UpdateZSetValue({ server, db, key, value, newValue: '', score: 0 })
                if (success) {
                    const { removed } = data
                    const tab = useTabStore()
                    if (!isEmpty(removed)) {
                        const removeValues = map(removed, 'v')
                        tab.removeValueEntries({ server, db, key, type: 'zset', entries: removeValues })
                    }
                    if (reload === true) {
                        this.reloadKey({ server, db, key })
                    } else {
                        // reload summary only
                        this.loadKeySummary({ server, db, key })
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
         * @param {boolean} [reload]
         * @returns {Promise<{[msg]: string, success: boolean}>}
         */
        async addStreamValue({ server, db, key, id, values, reload }) {
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
                        if (reload === true) {
                            this.reloadKey({ server, db, key })
                        } else {
                            // reload summary only
                            this.loadKeySummary({ server, db, key })
                        }
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
         * @param {string} server
         * @param {number} db
         * @param {string|number[]} key
         * @param {string[]|string} ids
         * @param {boolean} [reload]
         * @returns {Promise<{[msg]: {}, success: boolean}>}
         */
        async removeStreamValues({ server, db, key, ids, reload }) {
            if (typeof ids === 'string') {
                ids = [ids]
            }
            try {
                const { data = {}, success, msg } = await RemoveStreamValues(server, db, key, ids)
                if (success) {
                    const tab = useTabStore()
                    tab.removeValueEntries({ server, db, key, type: 'stream', entries: ids })
                    if (reload === true) {
                        this.reloadKey({ server, db, key })
                    } else {
                        // reload summary only
                        this.loadKeySummary({ server, db, key })
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
         * reset key's ttl
         * @param {string} server
         * @param {number} db
         * @param {string|number[]} key
         * @param {number} ttl
         * @returns {Promise<boolean>}
         */
        async setTTL(server, db, key, ttl) {
            try {
                const { success, msg } = await SetKeyTTL(server, db, key, ttl)
                if (success) {
                    const tabStore = useTabStore()
                    tabStore.updateTTL({
                        server,
                        db,
                        key: nativeRedisKey(key),
                        ttl,
                    })
                }
                return success === true
            } catch (e) {
                return false
            }
        },

        async setTTLs(server, db, keys, ttl) {
            // const msgRef = $message.loading('', { duration: 0, closable: true })
            // let updated = []
            // let failCount = 0
            // let canceled = false
            const serialNo = Date.now().valueOf().toString()
            // const eventName = 'ttling:' + serialNo
            // const cancelEvent = 'ttling:stop:' + serialNo
            try {
                // let maxProgress = 0
                // EventsOn(eventName, ({ total, progress, processing }) => {
                //     // update delete progress
                //     if (progress > maxProgress) {
                //         maxProgress = progress
                //     }
                //     const k = decodeRedisKey(processing)
                //     msgRef.content = i18nGlobal.t('dialogue.delete.doing', {
                //         key: k,
                //         index: maxProgress,
                //         count: total,
                //     })
                // })
                // msgRef.onClose = () => {
                //     EventsEmit(cancelEvent)
                // }
                const { data, success, msg } = await BatchSetTTL(server, db, keys, ttl, serialNo)
                if (success) {
                    // canceled = get(data, 'canceled', false)
                    // updated = get(data, 'updated', [])
                    // failCount = get(data, 'failed', 0)
                } else {
                    $message.error(msg)
                }
            } finally {
                // msgRef.destroy()
                // EventsOff(eventName)
            }
            $message.success(i18nGlobal.t('dialogue.ttl.success'))
            // const deletedCount = size(updated)
            // if (canceled) {
            //     $message.info(i18nGlobal.t('dialogue.handle_cancel'))
            // } else if (failCount <= 0) {
            //     // no fail
            //     $message.success(i18nGlobal.t('dialogue.delete.completed', { success: deletedCount, fail: failCount }))
            // } else if (failCount >= deletedCount) {
            //     // all fail
            //     $message.error(i18nGlobal.t('dialogue.delete.completed', { success: deletedCount, fail: failCount }))
            // } else {
            //     // some fail
            //     $message.warning(i18nGlobal.t('dialogue.delete.completed', { success: deletedCount, fail: failCount }))
            // }
        },

        /**
         * delete redis key
         * @param {string} server
         * @param {number} db
         * @param {string|number[]} key
         * @param {boolean} [soft] do not try to remove from redis if true, just remove from tree data
         * @returns {Promise<boolean>}
         */
        async deleteKey(server, db, key, soft) {
            try {
                let deleteCount = 1
                if (soft !== true) {
                    const { data } = await DeleteKey(server, db, key)
                    deleteCount = get(data, 'deleteCount', 0)
                }

                const k = nativeRedisKey(key)
                // update tree view data
                /** @type RedisServerState **/
                const serverInst = this.servers[server]
                if (serverInst != null) {
                    serverInst.removeKeyNode(k)
                    serverInst.tidyNode(k, true)
                    serverInst.updateDBKeyCount(db, -deleteCount)
                }

                // set tab content empty
                const tab = useTabStore()
                tab.emptyTab(server)
                tab.setSelectedKeys(server)
                tab.setCheckedKeys(server)
                return true
            } finally {
            }
            return false
        },

        /**
         * delete multiple keys
         * @param {string} server
         * @param {number} db
         * @param {string[]|number[][]} keys
         * @return {Promise<void>}
         */
        async deleteKeys(server, db, keys) {
            const msgRef = $message.loading(i18nGlobal.t('dialogue.delete.deleting'), { duration: 0, closable: true })
            let deleted = []
            let failCount = 0
            let canceled = false
            const serialNo = Date.now().valueOf().toString()
            msgRef.onClose = () => {
                EventsEmit('delete:stop:' + serialNo)
            }
            try {
                const { success, msg, data } = await DeleteKeys(server, db, keys, serialNo)
                if (success) {
                    canceled = get(data, 'canceled', false)
                    deleted = get(data, 'deleted', [])
                    failCount = get(data, 'failed', 0)
                } else {
                    $message.error(msg)
                }
            } finally {
                msgRef.destroy()
                // clear checked keys
                const tab = useTabStore()
                tab.setCheckedKeys(server)
            }
            // refresh model data
            const deletedCount = size(deleted)
            if (canceled) {
                $message.info(i18nGlobal.t('dialogue.handle_cancel'))
            } else if (failCount <= 0) {
                // no fail
                $message.success(i18nGlobal.t('dialogue.delete.completed', { success: deletedCount, fail: failCount }))
            } else if (failCount >= deletedCount) {
                // all fail
                $message.error(i18nGlobal.t('dialogue.delete.completed', { success: deletedCount, fail: failCount }))
            } else {
                // some fail
                $message.warning(i18nGlobal.t('dialogue.delete.completed', { success: deletedCount, fail: failCount }))
            }
            // update ui
            timeout(100).then(async () => {
                /** @type RedisServerState **/
                const serverInst = this.servers[server]
                if (serverInst != null) {
                    let start = now()
                    for (let i = 0; i < deleted.length; i++) {
                        serverInst.removeKeyNode(deleted[i], false)
                        if (now() - start > 300) {
                            await timeout(100)
                            start = now()
                        }
                    }
                    serverInst.tidyNode('', true)
                    serverInst.updateDBKeyCount(db, -deletedCount)
                }
            })
        },

        /**
         * delete multiple keys by pattern
         * @param server
         * @param db
         * @param pattern
         * @return {Promise<void>}
         */
        async deleteByPattern(server, db, pattern) {
            const msgRef = $message.loading(i18nGlobal.t('dialogue.delete.deleting'), { duration: 0, closable: true })
            let deleted = []
            let failCount = 0
            let canceled = false
            try {
                const { success, msg, data } = await DeleteKeysByPattern(server, db, pattern)
                if (success) {
                    canceled = get(data, 'canceled', false)
                    deleted = get(data, 'deleted', [])
                    failCount = get(data, 'failed', 0)
                } else {
                    $message.error(msg)
                }
            } finally {
                msgRef.destroy()
                // clear checked keys
                const tab = useTabStore()
                tab.setCheckedKeys(server)
            }
            // refresh model data
            const deletedCount = size(deleted)
            if (canceled) {
                $message.info(i18nGlobal.t('dialogue.handle_cancel'))
            } else if (failCount <= 0) {
                // no fail
                $message.success(i18nGlobal.t('dialogue.delete.completed', { success: deletedCount, fail: failCount }))
            } else if (failCount >= deletedCount) {
                // all fail
                $message.error(i18nGlobal.t('dialogue.delete.completed', { success: deletedCount, fail: failCount }))
            } else {
                // some fail
                $message.warning(i18nGlobal.t('dialogue.delete.completed', { success: deletedCount, fail: failCount }))
            }
            // update ui
            timeout(100).then(async () => {
                /** @type RedisServerState **/
                const serverInst = this.servers[server]
                if (serverInst != null) {
                    let start = now()
                    for (let i = 0; i < deleted.length; i++) {
                        serverInst.removeKeyNode(deleted[i], false)
                        if (now() - start > 300) {
                            await timeout(100)
                            start = now()
                        }
                    }
                    serverInst.tidyNode('', true)
                    serverInst.updateDBKeyCount(db, -deletedCount)
                }
            })
        },

        /**
         * export multiple keys
         * @param {string} server
         * @param {number} db
         * @param {string[]|number[][]} keys
         * @param {string} path
         * @param {boolean} [expire]
         * @returns {Promise<void>}
         */
        async exportKeys(server, db, keys, path, expire) {
            const msgRef = $message.loading('', { duration: 0, closable: true })
            let exported = 0
            let failCount = 0
            let canceled = false
            const cancelEventFn = EventsOn('exporting:' + path, ({ total, progress, processing }) => {
                // update export progress
                msgRef.content = i18nGlobal.t('dialogue.export.exporting', {
                    // key: decodeRedisKey(processing),
                    index: progress,
                    count: total,
                })
            })
            msgRef.onClose = () => {
                EventsEmit('export:stop:' + path)
            }
            try {
                const { data, success, msg } = await ExportKey(server, db, keys, path, expire)
                if (success) {
                    canceled = get(data, 'canceled', false)
                    exported = get(data, 'exported', 0)
                    failCount = get(data, 'failed', 0)
                } else {
                    $message.error(msg)
                }
            } finally {
                msgRef.destroy()
                cancelEventFn()
            }
            if (canceled) {
                $message.info(i18nGlobal.t('dialogue.handle_cancel'))
            } else if (failCount <= 0) {
                // no fail
                $message.success(
                    i18nGlobal.t('dialogue.export.export_completed', { success: exported, fail: failCount }),
                )
            } else if (failCount >= exported) {
                // all fail
                $message.error(i18nGlobal.t('dialogue.export.export_completed', { success: exported, fail: failCount }))
            } else {
                // some fail
                $message.warning(
                    i18nGlobal.t('dialogue.export.export_completed', {
                        success: exported,
                        fail: failCount,
                    }),
                )
            }
        },

        /**
         * import multiple keys from csv file
         * @param {string} server
         * @param {number} db
         * @param {string} path
         * @param {number} conflict
         * @param {number} [ttl] <0:use previous; ==0: persist; >0: custom ttl
         * @param {boolean} [reload]
         * @return {Promise<void>}
         */
        async importKeysFromCSVFile(server, db, path, conflict, ttl, reload) {
            const msgRef = $message.loading('', { duration: 0, closable: true })
            let imported = 0
            let ignored = 0
            let canceled = false
            const cancelEventFn = EventsOn('importing:' + path, ({ imported = 0, ignored = 0 }) => {
                // update export progress
                msgRef.content = i18nGlobal.t('dialogue.import.importing', {
                    // key: decodeRedisKey(processing),
                    imported,
                    conflict: ignored,
                })
            })
            msgRef.onClose = () => {
                EventsEmit('import:stop:' + path)
            }
            try {
                const { data, success, msg } = await ImportCSV(server, db, path, conflict, ttl)
                if (success) {
                    canceled = get(data, 'canceled', false)
                    imported = get(data, 'imported', 0)
                    ignored = get(data, 'ignored', 0)
                } else {
                    $message.error(msg)
                }
            } finally {
                cancelEventFn()
                msgRef.destroy()
            }
            if (canceled) {
                $message.info(i18nGlobal.t('dialogue.handle_cancel'))
            } else {
                // finish
                $message.success(i18nGlobal.t('dialogue.import.import_completed', { success: imported, ignored }))
                if (reload) {
                    this.reloadServer(server)
                }
            }
        },

        /**
         * flush database
         * @param {string} server
         * @param {number} db
         * @param {boolean} async
         * @return {Promise<boolean>}
         */
        async flushDatabase(server, db, async) {
            try {
                const { success = false } = await FlushDB(server, db, async)

                if (success === true) {
                    /** @type RedisServerState **/
                    const serverInst = this.servers[server]
                    if (serverInst != null) {
                        // update tree view data
                        serverInst.removeKeyNode()
                    }
                    // set tab content empty
                    const tab = useTabStore()
                    tab.emptyTab(server)
                    tab.setSelectedKeys(server)
                    tab.setCheckedKeys(server)
                    tab.setExpandedKeys(server)
                    return true
                }
            } finally {
            }
            return true
        },

        /**
         * rename key
         * @param {string} server
         * @param {number} db
         * @param {string} key
         * @param {string} newKey
         * @returns {Promise<{[msg]: string, success: boolean, [nodeKey]: string}>}
         */
        async renameKey(server, db, key, newKey) {
            const { success = false, msg } = await RenameKey(server, db, key, newKey)
            if (success) {
                // delete old key and add new key struct
                /** @type RedisServerState **/
                const serverInst = this.servers[server]
                if (serverInst != null) {
                    serverInst.renameKey(key, newKey)
                }
                return { success: true, nodeKey: `${server}/db${db}#${ConnectionType.RedisValue}/${newKey}` }
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
         * get client list info
         * @param {string} server
         * @return {Promise<{idle: number, name: string, addr: string, age: number, db: number}[]>}
         */
        async getClientList(server) {
            const { success, msg, data } = await GetClientList(server)
            if (success) {
                const { list = [] } = data
                return map(list, (item) => ({
                    addr: item['addr'],
                    name: item['name'],
                    age: item['age'] || 0,
                    idle: item['idle'] || 0,
                    db: item['db'] || 0,
                }))
            }
            return []
        },

        /**
         * get slow log list
         * @param {string} server
         * @param {number} num
         * @return {Promise<[]>}
         */
        async getSlowLog(server, num) {
            try {
                const { success, data = { list: [] } } = await GetSlowLogs(server, num)
                const { list } = data
                return list
            } catch {
                return []
            }
        },

        /**
         * get key filter pattern and filter type
         * @param {string} server
         * @returns {{match: string, type: string, exact: boolean}}
         */
        getKeyFilter(server) {
            let serverInst = this.servers[server]
            if (serverInst == null) {
                serverInst = new RedisServerState({
                    name: server,
                    separator: this.getSeparator(server),
                })
            }
            return serverInst.getFilter()
        },

        /**
         *
         * @param {string} server
         * @param {string} [pattern]
         * @param {string} [type]
         * @param {boolean} [exact]
         */
        setKeyFilter(server, { pattern, type, exact = false }) {
            const serverInst = this.servers[server]
            if (serverInst != null) {
                serverInst.setFilter({ pattern, type, exact })
            }
        },

        /**
         *
         * @param {string} server
         * @param {string} key
         * @param {number} db
         * @param {string} format
         * @param {string} decode
         */
        setSelectedFormat(server, key, db, format, decode) {
            const serverInst = this.servers[server]
            if (serverInst == null) {
                return
            }
            serverInst.addDecodeHistory(key, db, format, decode)
        },
    },
})

export default useBrowserStore
