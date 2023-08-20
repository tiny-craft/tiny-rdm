import { defineStore } from 'pinia'
import { endsWith, get, isEmpty, join, remove, size, slice, sortedIndexBy, split, sumBy, toUpper, uniq } from 'lodash'
import {
    AddHashField,
    AddListItem,
    AddStreamValue,
    AddZSetValue,
    CloseConnection,
    CreateGroup,
    DeleteConnection,
    DeleteGroup,
    DeleteKey,
    GetCmdHistory,
    GetConnection,
    GetKeyValue,
    ListConnection,
    OpenConnection,
    OpenDatabase,
    RemoveStreamValues,
    RenameGroup,
    RenameKey,
    SaveConnection,
    SaveSortedConnection,
    ScanKeys,
    ServerInfo,
    SetHashValue,
    SetKeyTTL,
    SetKeyValue,
    SetListItem,
    SetSetItem,
    UpdateSetItem,
    UpdateZSetValue,
} from 'wailsjs/go/services/connectionService.js'
import { ConnectionType } from '@/consts/connection_type.js'
import useTabStore from './tab.js'
import { types } from '@/consts/support_redis_type.js'

const useConnectionStore = defineStore('connections', {
    /**
     * @typedef {Object} ConnectionItem
     * @property {string} key
     * @property {string} label display label
     * @property {string} name database name
     * @property {number} type
     * @property {ConnectionItem[]} children
     */

    /**
     * @typedef {Object} DatabaseItem
     * @property {string} key - tree node unique key
     * @property {string} label
     * @property {string} [name] - server name, type != ConnectionType.Group only
     * @property {number} type
     * @property {number} [db] - database index, type == ConnectionType.RedisDB only
     * @property {string} [redisKey] - redis key, type == ConnectionType.RedisKey || type == ConnectionType.RedisValue only
     * @property {number} [keys] - children key count
     * @property {boolean} [isLeaf]
     * @property {boolean} [opened] - redis db is opened, type == ConnectionType.RedisDB only
     * @property {boolean} [expanded] - current node is expanded
     * @property {DatabaseItem[]} [children]
     */

    /**
     * @typedef {Object} ConnectionState
     * @property {string[]} groups
     * @property {ConnectionItem[]} connections
     * @property {Object} serverStats
     * @property {Object.<string, ConnectionProfile>} serverProfile
     * @property {Object.<string, string>} keyFilter key is 'server#db', 'server#-1' stores default filter pattern
     * @property {Object.<string, string>} typeFilter key is 'server#db'
     * @property {Object.<string, DatabaseItem[]>} databases
     * @property {Object.<string, Map<string, DatabaseItem>>} nodeMap key format likes 'server#db', children key format likes 'key#type'
     */

    /**
     * @typedef {Object} HistoryItem
     * @property {string} time
     * @property {string} server
     * @property {string} cmd
     * @property {number} cost
     */

    /**
     * @typedef {Object} ConnectionProfile
     * @property {string} defaultFilter
     * @property {string} keySeparator
     * @property {string} markColor
     */

    /**
     *
     * @returns {ConnectionState}
     */
    state: () => ({
        groups: [], // all group name set
        connections: [], // all connections
        serverStats: {}, // current server status info
        serverProfile: {}, // all server profile
        keyFilter: {}, // all key filters in opened connections group by server+db
        typeFilter: {}, // all key type filters in opened connections group by server+db
        databases: {}, // all databases in opened connections group by server name
        nodeMap: {}, // all nodes in opened connections group by server#db and type/key
    }),
    getters: {
        anyConnectionOpened() {
            return !isEmpty(this.databases)
        },
    },
    actions: {
        /**
         * load all store connections struct from local profile
         * @param {boolean} [force]
         * @returns {Promise<void>}
         */
        async initConnections(force) {
            if (!force && !isEmpty(this.connections)) {
                return
            }
            const conns = []
            const groups = []
            const profiles = {}
            const { data = [{ groupName: '', connections: [] }] } = await ListConnection()
            for (const conn of data) {
                if (conn.type !== 'group') {
                    // top level
                    conns.push({
                        key: conn.name,
                        label: conn.name,
                        name: conn.name,
                        type: ConnectionType.Server,
                        // isLeaf: false,
                    })
                    profiles[conn.name] = {
                        defaultFilter: conn.defaultFilter,
                        keySeparator: conn.keySeparator,
                        markColor: conn.markColor,
                    }
                } else {
                    // custom group
                    groups.push(conn.name)
                    const subConns = get(conn, 'connections', [])
                    const children = []
                    for (const item of subConns) {
                        const value = conn.name + '/' + item.name
                        children.push({
                            key: value,
                            label: item.name,
                            name: item.name,
                            type: ConnectionType.Server,
                            // isLeaf: false,
                        })
                    }
                    conns.push({
                        key: conn.name,
                        label: conn.name,
                        type: ConnectionType.Group,
                        children,
                    })
                    profiles[conn.name] = {
                        defaultFilter: conn.defaultFilter,
                        keySeparator: conn.keySeparator,
                        markColor: conn.markColor,
                    }
                }
                this.setKeyFilter(conn.name, -1, conn.defaultFilter)
            }
            this.connections = conns
            this.serverProfile = profiles
            this.groups = uniq(groups)
        },

        /**
         * get connection by name from local profile
         * @param name
         * @returns {Promise<ConnectionProfile|null>}
         */
        async getConnectionProfile(name) {
            try {
                const { data, success } = await GetConnection(name)
                if (success) {
                    this.serverProfile[name] = {
                        defaultFilter: data.defaultFilter,
                        keySeparator: data.keySeparator,
                        markColor: data.markColor,
                    }
                    return data
                }
            } finally {
            }
            return null
        },

        /**
         * create a new default connection
         * @param {string} [name]
         * @returns {{}}
         */
        newDefaultConnection(name) {
            return {
                group: '',
                name: name || '',
                addr: '127.0.0.1',
                port: 6379,
                username: '',
                password: '',
                defaultFilter: '*',
                keySeparator: ':',
                connTimeout: 60,
                execTimeout: 60,
                markColor: '',
            }
        },

        /**
         * get database server by name
         * @param name
         * @returns {ConnectionItem|null}
         */
        getConnection(name) {
            const conns = this.connections
            for (let i = 0; i < conns.length; i++) {
                if (conns[i].type === ConnectionType.Server && conns[i].key === name) {
                    return conns[i]
                } else if (conns[i].type === ConnectionType.Group) {
                    const children = conns[i].children
                    for (let j = 0; j < children.length; j++) {
                        if (children[j].type === ConnectionType.Server && conns[i].key === name) {
                            return children[j]
                        }
                    }
                }
            }
            return null
        },

        /**
         * create a new connection or update current connection profile
         * @param {string} name set null if create a new connection
         * @param {{}} param
         * @returns {Promise<{success: boolean, [msg]: string}>}
         */
        async saveConnection(name, param) {
            const { success, msg } = await SaveConnection(name, param)
            if (!success) {
                return { success: false, msg }
            }

            // reload connection list
            await this.initConnections(true)
            return { success: true }
        },

        /**
         * save connection after sort
         * @returns {Promise<void>}
         */
        async saveConnectionSorted() {
            const mapToList = (conns) => {
                const list = []
                for (const conn of conns) {
                    if (conn.type === ConnectionType.Group) {
                        const children = mapToList(conn.children)
                        list.push({
                            name: conn.label,
                            type: 'group',
                            connections: children,
                        })
                    } else if (conn.type === ConnectionType.Server) {
                        list.push({
                            name: conn.name,
                        })
                    }
                }
                return list
            }
            const s = mapToList(this.connections)
            SaveSortedConnection(s)
        },

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
            const { db } = data
            if (isEmpty(db)) {
                throw new Error('no db loaded')
            }
            const dbs = []
            for (let i = 0; i < db.length; i++) {
                this._getNodeMap(name, i).clear()
                dbs.push({
                    key: `${name}/${db[i].name}`,
                    label: db[i].name,
                    name: name,
                    keys: db[i].keys,
                    db: i,
                    type: ConnectionType.RedisDB,
                    isLeaf: false,
                    children: undefined,
                })
            }
            this.databases[name] = dbs
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
            for (const db of dbs) {
                this.removeKeyFilter(name, db.db)
                this._getNodeMap(name, db.db).clear()
            }
            this.removeKeyFilter(name, -1)
            delete this.databases[name]
            delete this.serverStats[name]

            const tabStore = useTabStore()
            tabStore.removeTabByName(name)
            return true
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
            this.serverStats = {}
            const tabStore = useTabStore()
            tabStore.removeAllTab()
        },

        /**
         * remove connection
         * @param name
         * @returns {Promise<{success: boolean, [msg]: string}>}
         */
        async deleteConnection(name) {
            // close connection first
            await this.closeConnection(name)
            const { success, msg } = await DeleteConnection(name)
            if (!success) {
                return { success: false, msg }
            }
            await this.initConnections(true)
            return { success: true }
        },

        /**
         * create a connection group
         * @param name
         * @returns {Promise<{success: boolean, [msg]: string}>}
         */
        async createGroup(name) {
            const { success, msg } = await CreateGroup(name)
            if (!success) {
                return { success: false, msg }
            }
            await this.initConnections(true)
            return { success: true }
        },

        /**
         * rename connection group
         * @param name
         * @param newName
         * @returns {Promise<{success: boolean, [msg]: string}>}
         */
        async renameGroup(name, newName) {
            if (name === newName) {
                return { success: true }
            }
            const { success, msg } = await RenameGroup(name, newName)
            if (!success) {
                return { success: false, msg }
            }
            await this.initConnections(true)
            return { success: true }
        },

        /**
         * delete group by name
         * @param {string} name
         * @param {boolean} [includeConn]
         * @returns {Promise<{success: boolean, [msg]: string}>}
         */
        async deleteGroup(name, includeConn) {
            const { success, msg } = await DeleteGroup(name, includeConn === true)
            if (!success) {
                return { success: false, msg }
            }
            await this.initConnections(true)
            return { success: true }
        },

        /**
         * open database and load all keys
         * @param connName
         * @param db
         * @returns {Promise<void>}
         */
        async openDatabase(connName, db) {
            const { match: filterPattern, type: keyType } = this.getKeyFilter(connName, db)
            const { data, success, msg } = await OpenDatabase(connName, db, filterPattern, keyType)
            if (!success) {
                throw new Error(msg)
            }
            const { keys = [] } = data
            const dbs = this.databases[connName]
            dbs[db].opened = true
            if (isEmpty(keys)) {
                dbs[db].children = []
                return
            }

            // append db node to current connection's children
            this._addKeyNodes(connName, db, keys)
            this._tidyNode(connName, db)
        },

        /**
         * reopen database
         * @param connName
         * @param db
         * @returns {Promise<void>}
         */
        async reopenDatabase(connName, db) {
            const dbs = this.databases[connName]
            dbs[db].children = undefined
            dbs[db].isLeaf = false

            this._getNodeMap(connName, db).clear()
        },

        /**
         * close database
         * @param connName
         * @param db
         */
        closeDatabase(connName, db) {
            const dbs = this.databases[connName]
            delete dbs[db].children
            dbs[db].isLeaf = false
            dbs[db].opened = false

            this._getNodeMap(connName, db).clear()
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
         * load redis key
         * @param {string} server
         * @param {number} db
         * @param {string} [key] when key is null or blank, update tab to display normal content (blank content or server status)
         */
        async loadKeyValue(server, db, key) {
            try {
                const tab = useTabStore()
                if (!isEmpty(key)) {
                    const { data, success, msg } = await GetKeyValue(server, db, key)
                    if (success) {
                        const { type, ttl, value } = data
                        tab.upsertTab({
                            server,
                            db,
                            type,
                            ttl,
                            key,
                            value,
                        })
                        return
                    } else {
                        // key not exists, remove this key
                        await this.deleteKey(server, db, key)
                    }
                }

                tab.upsertTab({
                    server,
                    db,
                    type: 'none',
                    ttl: -1,
                    key: null,
                    value: null,
                })
            } finally {
            }
        },

        /**
         * scan keys with prefix
         * @param {string} connName
         * @param {number} db
         * @param {string} [prefix] full reload database if prefix is null
         * @returns {Promise<{keys: string[]}>}
         */
        async scanKeys(connName, db, prefix) {
            const { data, success, msg } = await ScanKeys(connName, db, prefix || '*')
            if (!success) {
                throw new Error(msg)
            }
            const { keys = [] } = data
            return { keys, success }
        },

        /**
         * load keys with prefix
         * @param {string} connName
         * @param {number} db
         * @param {string} [prefix]
         * @returns {Promise<void>}
         */
        async loadKeys(connName, db, prefix) {
            let scanPrefix = prefix
            if (isEmpty(scanPrefix)) {
                scanPrefix = '*'
            } else {
                const separator = this._getSeparator(connName)
                if (!endsWith(prefix, separator + '*')) {
                    scanPrefix = prefix + separator + '*'
                }
            }
            const { keys, success } = await this.scanKeys(connName, db, scanPrefix)
            if (!success) {
                return
            }

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
            const { keySeparator } = this.serverProfile[server] || { keySeparator: ':' }
            if (isEmpty(keySeparator)) {
                return ':'
            }
            return keySeparator
        },

        /**
         * get node map
         * @param connName
         * @param db
         * @returns {Map<string, DatabaseItem>}
         * @private
         */
        _getNodeMap(connName, db) {
            if (this.nodeMap[`${connName}#${db}`] == null) {
                this.nodeMap[`${connName}#${db}`] = new Map()
            }
            // construct a tree node list, the format of item key likes 'server/db#type/key'
            return this.nodeMap[`${connName}#${db}`]
        },

        /**
         * remove keys in db
         * @param {string} connName
         * @param {number} db
         * @param {string[]} keys
         * @param {boolean} [sortInsert]
         * @return {{success: boolean, newKey: number, newLayer: number, replaceKey: number}}
         * @private
         */
        _addKeyNodes(connName, db, keys, sortInsert) {
            const result = { success: false, newLayer: 0, newKey: 0, replaceKey: 0 }
            if (isEmpty(keys)) {
                return result
            }
            const separator = this._getSeparator(connName)
            const dbs = this.databases[connName]
            if (dbs[db].children == null) {
                dbs[db].children = []
            }
            const nodeMap = this._getNodeMap(connName, db)
            const rootChildren = dbs[db].children
            for (const key of keys) {
                const keyParts = split(key, separator)
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
                                const index = sortedIndexBy(children, selectedNode, (elem) => {
                                    return elem.key
                                })
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
                            label: keyParts[i],
                            db,
                            keys: 0,
                            redisKey: handlePath,
                            type: ConnectionType.RedisValue,
                            isLeaf: true,
                        }
                        nodeMap.set(nodeKey, selectedNode)
                        if (!replaceKey) {
                            if (sortInsert) {
                                const index = sortedIndexBy(children, selectedNode, (elem) => {
                                    return elem.key > selectedNode.key
                                })
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
            const separator = this._getSeparator(connName)
            const keyParts = split(key, separator)
            const totalParts = size(keyParts)
            const dbNode = get(this.databases, [connName, db], {})
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
         * @param node
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
                count += 1
            }
            if (node.keys !== count) {
                node.keys = count
                return true
            }
            return false
        },

        /**
         * get tree node by key name
         * @param key
         * @return {DatabaseItem|null}
         */
        getNode(key) {
            const idx = key.indexOf('#')
            if (idx < 0) {
                return null
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
                return this.databases[server][db]
            }
        },

        /**
         * set redis key
         * @param {string} connName
         * @param {number} db
         * @param {string} key
         * @param {string} keyType
         * @param {any} value
         * @param {number} ttl
         * @returns {Promise<{[msg]: string, success: boolean, [nodeKey]: {string}}>}
         */
        async setKey(connName, db, key, keyType, value, ttl) {
            try {
                const { data, success, msg } = await SetKeyValue(connName, db, key, keyType, value, ttl)
                if (success) {
                    // update tree view data
                    const { newKey = 0 } = this._addKeyNodes(connName, db, [key], true)
                    if (newKey > 0) {
                        this._tidyNode(connName, db, key)
                    }
                    return { success, nodeKey: `${connName}/db${db}#${ConnectionType.RedisValue}/${key}` }
                } else {
                    return { success, msg }
                }
            } catch (e) {
                return { success: false, msg: e.message }
            }
        },

        /**
         * update hash field
         * when field is set, newField is null, delete field
         * when field is null, newField is set, add new field
         * when both field and newField are set, and field === newField, update field
         * when both field and newField are set, and field !== newField, delete field and add newField
         * @param {string} connName
         * @param {number} db
         * @param {string} key
         * @param {string} field
         * @param {string} newField
         * @param {string} value
         * @returns {Promise<{[msg]: string, success: boolean, [updated]: {}}>}
         */
        async setHash(connName, db, key, field, newField, value) {
            try {
                const { data, success, msg } = await SetHashValue(connName, db, key, field, newField || '', value || '')
                if (success) {
                    const { updated = {} } = data
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
         * @param {string} connName
         * @param {number} db
         * @param {string} key
         * @param {number }action 0:ignore duplicated fields 1:overwrite duplicated fields
         * @param {string[]} fieldItems field1, value1, filed2, value2...
         * @returns {Promise<{[msg]: string, success: boolean, [updated]: {}}>}
         */
        async addHashField(connName, db, key, action, fieldItems) {
            try {
                const { data, success, msg } = await AddHashField(connName, db, key, action, fieldItems)
                if (success) {
                    const { updated = {} } = data
                    return { success, updated }
                } else {
                    return { success: false, msg }
                }
            } catch (e) {
                return { success: false, msg: e.message }
            }
        },

        /**
         * remove hash field
         * @param {string} connName
         * @param {number} db
         * @param {string} key
         * @param {string} field
         * @returns {Promise<{[msg]: {}, success: boolean, [removed]: string[]}>}
         */
        async removeHashField(connName, db, key, field) {
            try {
                const { data, success, msg } = await SetHashValue(connName, db, key, field, '', '')
                if (success) {
                    const { removed = [] } = data
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
         * @param {string} key
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
         * @param connName
         * @param db
         * @param key
         * @param values
         * @returns {Promise<[msg]: string, success: boolean, [item]: []>}
         */
        async prependListItem(connName, db, key, values) {
            try {
                const { data, success, msg } = await AddListItem(connName, db, key, 0, values)
                if (success) {
                    const { left = [] } = data
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
         * @param connName
         * @param db
         * @param key
         * @param values
         * @returns {Promise<[msg]: string, success: boolean, [item]: any[]>}
         */
        async appendListItem(connName, db, key, values) {
            try {
                const { data, success, msg } = await AddListItem(connName, db, key, 1, values)
                if (success) {
                    const { right = [] } = data
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
         * @param {string} connName
         * @param {number} db
         * @param {string} key
         * @param {number} index
         * @param {string} value
         * @returns {Promise<{[msg]: string, success: boolean, [updated]: {}}>}
         */
        async updateListItem(connName, db, key, index, value) {
            try {
                const { data, success, msg } = await SetListItem(connName, db, key, index, value)
                if (success) {
                    const { updated = {} } = data
                    return { success, updated }
                } else {
                    return { success, msg }
                }
            } catch (e) {
                return { success: false, msg: e.message }
            }
        },

        /**
         * remove list item
         * @param {string} connName
         * @param {number} db
         * @param {string} key
         * @param {number} index
         * @returns {Promise<{[msg]: string, success: boolean, [removed]: string[]}>}
         */
        async removeListItem(connName, db, key, index) {
            try {
                const { data, success, msg } = await SetListItem(connName, db, key, index, '')
                if (success) {
                    const { removed = [] } = data
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
         * @param {string} connName
         * @param {number} db
         * @param {string} key
         * @param {string} value
         * @returns {Promise<{[msg]: string, success: boolean}>}
         */
        async addSetItem(connName, db, key, value) {
            try {
                const { success, msg } = await SetSetItem(connName, db, key, false, [value])
                if (success) {
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
         * @param {string} connName
         * @param {number} db
         * @param {string} key
         * @param {string} value
         * @param {string} newValue
         * @returns {Promise<{[msg]: string, success: boolean}>}
         */
        async updateSetItem(connName, db, key, value, newValue) {
            try {
                const { success, msg } = await UpdateSetItem(connName, db, key, value, newValue)
                if (success) {
                    return { success: true }
                } else {
                    return { success, msg }
                }
            } catch (e) {
                return { success: false, msg: e.message }
            }
        },

        /**
         * remove item from set
         * @param connName
         * @param db
         * @param key
         * @param value
         * @returns {Promise<{[msg]: string, success: boolean}>}
         */
        async removeSetItem(connName, db, key, value) {
            try {
                const { success, msg } = await SetSetItem(connName, db, key, true, [value])
                if (success) {
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
         * @param {string} connName
         * @param {number} db
         * @param {string} key
         * @param {number} action
         * @param {Object.<string, number>} vs value: score
         * @returns {Promise<{[msg]: string, success: boolean}>}
         */
        async addZSetItem(connName, db, key, action, vs) {
            try {
                const { success, msg } = await AddZSetValue(connName, db, key, action, vs)
                if (success) {
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
         * @param {string} connName
         * @param {number} db
         * @param {string} key
         * @param {string} value
         * @param {string} newValue
         * @param {number} score
         * @returns {Promise<{[msg]: string, success: boolean, [updated]: {}, [removed]: []}>}
         */
        async updateZSetItem(connName, db, key, value, newValue, score) {
            try {
                const { data, success, msg } = await UpdateZSetValue(connName, db, key, value, newValue, score)
                if (success) {
                    const { updated, removed } = data
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
         * @param {string} connName
         * @param {number} db
         * @param key
         * @param {string} value
         * @returns {Promise<{[msg]: string, success: boolean, [removed]: []}>}
         */
        async removeZSetItem(connName, db, key, value) {
            try {
                const { data, success, msg } = await UpdateZSetValue(connName, db, key, value, '', 0)
                if (success) {
                    const { removed } = data
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
         * @param {string} connName
         * @param {number} db
         * @param {string} key
         * @param {string} id
         * @param {string[]} values field1, value1, filed2, value2...
         * @returns {Promise<{[msg]: string, success: boolean, [updated]: {}}>}
         */
        async addStreamValue(connName, db, key, id, values) {
            try {
                const { data = {}, success, msg } = await AddStreamValue(connName, db, key, id, values)
                if (success) {
                    const { updated = {} } = data
                    return { success, updated }
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
         * @param {string} key
         * @param {string[]|string} ids
         * @returns {Promise<{[msg]: {}, success: boolean, [removed]: string[]}>}
         */
        async removeStreamValues(connName, db, key, ids) {
            if (typeof ids === 'string') {
                ids = [ids]
            }
            try {
                const { data = {}, success, msg } = await RemoveStreamValues(connName, db, key, ids)
                if (success) {
                    const { removed = [] } = data
                    return { success, removed }
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
         * @param {boolean} [isLayer]
         * @private
         */
        _deleteKeyNode(connName, db, key, isLayer) {
            const dbRoot = get(this.databases, [connName, db], {})
            const separator = this._getSeparator(connName)

            if (dbRoot == null) {
                return false
            }

            const nodeMap = this._getNodeMap(connName, db)
            const keyParts = split(key, separator)
            const totalParts = size(keyParts)
            if (isLayer === true) {
                this._deleteChildrenKeyNodes(nodeMap, key)
            }
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
                    } else {
                        break
                    }
                } else {
                    // last one, remove from db node
                    remove(dbRoot.children, { type: ConnectionType.RedisKey, redisKey: keyParts[0] })
                }
            }

            return true
        },

        /**
         * delete node and all it's children from nodeMap
         * @param nodeMap
         * @param key
         * @private
         */
        _deleteChildrenKeyNodes(nodeMap, key) {
            const mapKey = `${ConnectionType.RedisKey}/${key}`
            const node = nodeMap.get(mapKey)
            for (const child of node.children || []) {
                if (child.type === ConnectionType.RedisValue) {
                    if (!nodeMap.delete(`${ConnectionType.RedisValue}/${child.redisKey}`)) {
                        console.warn('delete:', `${ConnectionType.RedisValue}/${child.redisKey}`)
                    }
                } else if (child.type === ConnectionType.RedisKey) {
                    this._deleteChildrenKeyNodes(nodeMap, child.redisKey)
                }
            }
            if (!nodeMap.delete(mapKey)) {
                console.warn('delete map key', mapKey)
            }
        },

        /**
         * delete redis key
         * @param {string} connName
         * @param {number} db
         * @param {string} key
         * @returns {Promise<boolean>}
         */
        async deleteKey(connName, db, key) {
            try {
                const { data, success, msg } = await DeleteKey(connName, db, key)
                if (success) {
                    // update tree view data
                    this._deleteKeyNode(connName, db, key)
                    this._tidyNode(connName, db, key, true)

                    // set tab content empty
                    const tab = useTabStore()
                    tab.emptyTab(connName)
                    return true
                }
            } finally {
            }
            return false
        },

        /**
         * delete keys with prefix
         * @param connName
         * @param db
         * @param prefix
         * @returns {Promise<boolean>}
         */
        async deleteKeyPrefix(connName, db, prefix) {
            if (isEmpty(prefix)) {
                return false
            }
            try {
                if (!endsWith(prefix, '*')) {
                    prefix += '*'
                }
                const { data, success, msg } = await DeleteKey(connName, db, prefix)
                if (success) {
                    // const { deleted: keys = [] } = data
                    // for (const key of keys) {
                    //     await this._deleteKeyNode(connName, db, key)
                    // }
                    const separator = this._getSeparator(connName)
                    if (endsWith(prefix, '*')) {
                        prefix = prefix.substring(0, prefix.length - 1)
                    }
                    if (endsWith(prefix, separator)) {
                        prefix = prefix.substring(0, prefix.length - 1)
                    }
                    this._deleteKeyNode(connName, db, prefix, true)
                    this._tidyNode(connName, db, prefix, true)
                    return true
                }
            } finally {
            }
            return false
        },

        /**
         * rename key
         * @param {string} connName
         * @param {number} db
         * @param {string} key
         * @param {string} newKey
         * @returns {Promise<{[msg]: string, success: boolean}>}
         */
        async renameKey(connName, db, key, newKey) {
            const { success = false, msg } = await RenameKey(connName, db, key, newKey)
            if (success) {
                // delete old key and add new key struct
                this._deleteKeyNode(connName, db, key)
                this._addKeyNodes(connName, db, [newKey])
                return { success: true }
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
         * get key filter pattern and filter type
         * @param {string} server
         * @param {number} db
         * @returns {{match: string, type: string}}
         */
        getKeyFilter(server, db) {
            let match, type
            const key = `${server}#${db}`
            if (!this.keyFilter.hasOwnProperty(key)) {
                match = this.keyFilter[`${server}#-1`] || '*'
            } else {
                match = this.keyFilter[key] || '*'
            }
            type = this.typeFilter[`${server}#${db}`] || ''
            return {
                match,
                type: toUpper(type),
            }
        },

        /**
         * set key filter
         * @param {string} server
         * @param {number} db
         * @param {string} pattern
         * @param {string} [type]
         */
        setKeyFilter(server, db, pattern, type) {
            this.keyFilter[`${server}#${db}`] = pattern || '*'
            this.typeFilter[`${server}#${db}`] = types[toUpper(type)] || ''
        },

        removeKeyFilter(server, db) {
            this.keyFilter[`${server}#${db}`] = '*'
            delete this.typeFilter[`${server}#${db}`]
        },
    },
})

export default useConnectionStore
