import { defineStore } from 'pinia'
import { endsWith, findIndex, get, isEmpty, size, split, uniq } from 'lodash'
import {
    AddHashField,
    AddListItem,
    AddZSetValue,
    CloseConnection,
    CreateGroup,
    DeleteConnection,
    DeleteGroup,
    DeleteKey,
    GetConnection,
    GetKeyValue,
    ListConnection,
    OpenConnection,
    OpenDatabase,
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
} from '../../wailsjs/go/services/connectionService.js'
import { ConnectionType } from '../consts/connection_type.js'
import useTabStore from './tab.js'

const separator = ':'

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
     * @property {string} key
     * @property {string} label
     * @property {string} name - server name, type != ConnectionType.Group only
     * @property {number} type
     * @property {number} [db] - database index, type == ConnectionType.RedisDB only
     * @property {number} keys
     * @property {boolean} [isLeaf]
     * @property {boolean} [opened] - redis db is opened, type == ConnectionType.RedisDB only
     * @property {boolean} [expanded] - current node is expanded
     */

    /**
     * @typedef {Object} ConnectionState
     * @property {string[]} groups
     * @property {ConnectionItem[]} connections
     * @property {Object.<string, DatabaseItem[]>} databases
     * @property {Object.<string, Map<string, DatabaseItem>>} nodeMap key format likes 'server#db', children key format likes 'key#type'
     */

    /**
     *
     * @returns {ConnectionState}
     */
    state: () => ({
        groups: [], // all group name set
        connections: [], // all connections
        serverStats: {}, // current server status info
        databases: {}, // all databases in opened connections group by server name
        nodeMap: {}, // all node in opened connections group by server+db and key+type
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
                }
            }
            this.connections = conns
            this.groups = uniq(groups)
        },

        /**
         * get connection by name from local profile
         * @param name
         * @returns {Promise<{}|null>}
         */
        async getConnectionProfile(name) {
            try {
                const { data, success } = await GetConnection(name)
                if (success) {
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
                dbs.push({
                    key: `${name}/${db[i].name}`,
                    label: db[i].name,
                    name: name,
                    keys: db[i].keys,
                    db: i,
                    type: ConnectionType.RedisDB,
                    isLeaf: false,
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
            const { data, success, msg } = await OpenDatabase(connName, db)
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
            this._tidyNodeChildren(dbs[db])
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

            delete this.nodeMap[`${connName}#${db}`]
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
                if (!endsWith(prefix, separator + '*')) {
                    scanPrefix = prefix + separator + '*'
                }
            }
            const { keys, success } = await this.scanKeys(connName, db, scanPrefix)
            if (!success) {
                return
            }

            // remove current keys below prefix
            this._deleteKeyNodes(connName, db, prefix)
            this._addKeyNodes(connName, db, keys)
            this._tidyNodeChildren(this.databases[connName][db])
        },

        /**
         * remove key with prefix
         * @param {string} connName
         * @param {number} db
         * @param {string} prefix
         * @returns {boolean}
         * @private
         */
        _deleteKeyNodes(connName, db, prefix) {
            const dbs = this.databases[connName]
            let node = dbs[db]
            const prefixPart = split(prefix, separator)
            const partLen = size(prefixPart)
            for (let i = 0; i < partLen; i++) {
                let idx = findIndex(node.children, { label: prefixPart[i] })
                if (idx === -1) {
                    node = null
                    break
                }
                if (i === partLen - 1) {
                    // remove last part from parent
                    node.children.splice(idx, 1)
                    return true
                } else {
                    node = node.children[idx]
                }
            }
            return false
        },

        /**
         * remove keys in db
         * @param {string} connName
         * @param {number} db
         * @param {string[]} keys
         * @private
         */
        _addKeyNodes(connName, db, keys) {
            const dbs = this.databases[connName]
            if (dbs[db].children == null) {
                dbs[db].children = []
            }
            if (this.nodeMap[`${connName}#${db}`] == null) {
                this.nodeMap[`${connName}#${db}`] = new Map()
            }
            // construct tree node list, the format of item key likes 'server/db#type/key'
            const nodeMap = this.nodeMap[`${connName}#${db}`]
            const rootChildren = dbs[db].children
            let count = 0
            for (const key of keys) {
                const keyPart = split(key, separator)
                // const prefixLen = size(keyPart) - 1
                const len = size(keyPart)
                const lastIdx = len - 1
                let handlePath = ''
                let children = rootChildren
                for (let i = 0; i < len; i++) {
                    handlePath += keyPart[i]
                    if (i !== lastIdx) {
                        // layer
                        const nodeKey = `#${ConnectionType.RedisKey}/${handlePath}`
                        let selectedNode = nodeMap.get(nodeKey)
                        if (selectedNode == null) {
                            selectedNode = {
                                key: `${connName}/db${db}${nodeKey}`,
                                label: keyPart[i],
                                db,
                                keys: 0,
                                redisKey: handlePath,
                                type: ConnectionType.RedisKey,
                                isLeaf: false,
                                children: [],
                            }
                            nodeMap.set(nodeKey, selectedNode)
                            children.push(selectedNode)
                        }
                        children = selectedNode.children
                        handlePath += separator
                    } else {
                        // key
                        const nodeKey = `#${ConnectionType.RedisValue}/${handlePath}`
                        const selectedNode = {
                            key: `${connName}/db${db}${nodeKey}`,
                            label: keyPart[i],
                            db,
                            keys: 0,
                            redisKey: handlePath,
                            type: ConnectionType.RedisValue,
                            isLeaf: true,
                        }
                        nodeMap.set(nodeKey, selectedNode)
                        children.push(selectedNode)
                    }
                }
            }
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
         * sort all node item's children and calculate keys count
         * @param node
         * @private
         */
        _tidyNodeChildren(node) {
            let count = 0
            if (!isEmpty(node.children)) {
                this._sortNodes(node.children)

                for (const elem of node.children) {
                    this._tidyNodeChildren(elem)
                    count += elem.keys
                }
            } else {
                count += 1
            }
            node.keys = count
        },

        /**
         * set redis key
         * @param {string} connName
         * @param {number} db
         * @param {string} key
         * @param {number} keyType
         * @param {any} value
         * @param {number} ttl
         * @returns {Promise<{[msg]: string, success: boolean}>}
         */
        async setKey(connName, db, key, keyType, value, ttl) {
            try {
                const { data, success, msg } = await SetKeyValue(connName, db, key, keyType, value, ttl)
                if (success) {
                    // update tree view data
                    // this._addKey(connName, db, key)
                    this._addKeyNodes(connName, db, [key])
                    this._tidyNodeChildren(this.databases[connName][db])
                    return { success }
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
         * @private
         */
        _deleteKeyNode(connName, db, key) {
            const dbs = this.databases[connName]
            const dbDetail = get(dbs, db, {})

            if (dbDetail == null) {
                return
            }

            const nodeMap = this.nodeMap[`${connName}#${db}`]
            if (nodeMap == null) {
                return
            }
            const idx = key.lastIndexOf(separator)
            let parentNode = null
            let parentKey = ''
            if (idx === -1) {
                // root
                parentNode = dbDetail
            } else {
                parentKey = key.substring(0, idx)
                parentNode = nodeMap.get(`#${ConnectionType.RedisKey}/${parentKey}`)
            }

            if (parentNode == null || parentNode.children == null) {
                return
            }

            // remove children
            const delIdx = findIndex(parentNode.children, { redisKey: key })
            if (delIdx !== -1) {
                const childKeys = parentNode.children[delIdx].keys || 1
                parentNode.children.splice(delIdx, 1)
                parentNode.keys = Math.max(parentNode.keys - childKeys, 0)
            }

            // also remove parent node if no more children
            while (isEmpty(parentNode.children)) {
                const idx = parentKey.lastIndexOf(separator)
                if (idx !== -1) {
                    parentKey = parentKey.substring(0, idx)
                    parentNode = nodeMap.get(`#${ConnectionType.RedisKey}/${parentKey}`)
                    if (parentNode != null) {
                        parentNode.keys = (parentNode.keys || 1) - 1
                        parentNode.children = []
                    }
                } else {
                    // reach root, remove from db
                    const delIdx = findIndex(dbDetail.children, { redisKey: parentKey })
                    dbDetail.keys = (dbDetail.keys || 1) - 1
                    dbDetail.children.splice(delIdx, 1)
                    break
                }
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
                    if (endsWith(prefix, '*')) {
                        prefix = prefix.substring(0, prefix.length - 1)
                    }
                    if (endsWith(prefix, separator)) {
                        prefix = prefix.substring(0, prefix.length - 1)
                    }
                    await this._deleteKeyNode(connName, db, prefix)
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
    },
})

export default useConnectionStore
