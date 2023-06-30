import { get, isEmpty, last, remove, size, sortedIndexBy, split } from 'lodash'
import { defineStore } from 'pinia'
import {
    AddHashField,
    AddListItem,
    AddZSetValue,
    CloseConnection,
    GetKeyValue,
    OpenConnection,
    OpenDatabase,
    RemoveKey,
    RenameKey,
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
import useConnectionStore from './connections.js'

const separator = ':'

const useDatabaseStore = defineStore('database', {
    /**
     * @typedef {Object} DatabaseItem
     * @property {string} key
     * @property {string} label
     * @property {string} name - server name, type != ConnectionType.Group only
     * @property {number} type
     * @property {number} [db] - database index, type == ConnectionType.RedisDB only
     * @property {number} keys
     * @property {boolean} [connected] - redis server is connected, type == ConnectionType.Server only
     * @property {boolean} [opened] - redis db is opened, type == ConnectionType.RedisDB only
     * @property {boolean} [expanded] - current node is expanded
     */

    /**
     *
     * @returns {{connections: DatabaseItem[]}}
     */
    state: () => ({
        connections: [], // all connections list
        databases: {}, // all database group by opened connections
    }),
    getters: {},
    actions: {
        /**
         * Open connection
         * @param {string} connName
         * @returns {Promise<void>}
         */
        async openConnection(connName) {
            const { data, success, msg } = await OpenConnection(connName)
            if (!success) {
                throw new Error(msg)
            }
            // append to db node to current connection
            const connStore = useConnectionStore()
            const connNode = connStore.getConnection(connName)
            if (connNode == null) {
                throw new Error('no such connection')
            }
            const { db } = data
            if (isEmpty(db)) {
                throw new Error('no db loaded')
            }
            const children = []
            for (let i = 0; i < db.length; i++) {
                children.push({
                    key: `${connName}/${db[i].name}`,
                    label: db[i].name,
                    name: connName,
                    keys: db[i].keys,
                    db: i,
                    type: ConnectionType.RedisDB,
                    // isLeaf: false,
                })
            }
            connNode.children = children
            connNode.connected = true
        },

        /**
         * Close connection
         * @param {string} connName
         * @returns {Promise<boolean>}
         */
        async closeConnection(connName) {
            const { success, msg } = await CloseConnection(connName)
            if (!success) {
                // throw new Error(msg)
                return false
            }

            // successful close connection, remove all children
            const connNode = this.getConnection(connName)
            if (connNode == null) {
                // throw new Error('no such connection')
                return false
            }
            connNode.children = undefined
            connNode.isLeaf = undefined
            connNode.connected = false
            connNode.expanded = false
            return true
        },

        /**
         * Get Connection by path name
         * @param {string} connName
         * @returns
         */
        getConnection(connName) {
            const conn = this.connections
            for (let i = 0; i < conn.length; i++) {
                if (conn[i].type === ConnectionType.Server && conn[i].key === connName) {
                    return conn[i]
                } else if (conn[i].type === ConnectionType.Group) {
                    const children = conn[i].children
                    for (let j = 0; j < children.length; j++) {
                        if (children[j].type === ConnectionType.Server && conn[i].key === connName) {
                            return children[j]
                        }
                    }
                }
            }
            return null
        },

        /**
         * Open database and load all keys
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
            if (isEmpty(keys)) {
                const connNode = this.getConnection(connName)
                const { children = [] } = connNode
                children[db].children = []
                children[db].opened = true
                return
            }

            // insert child to children list by order
            const sortedInsertChild = (childrenList, item) => {
                const insertIdx = sortedIndexBy(childrenList, item, 'key')
                childrenList.splice(insertIdx, 0, item)
                // childrenList.push(item)
            }
            // update all node item's children num
            const updateChildrenNum = (node) => {
                let count = 0
                const totalChildren = size(node.children)
                if (totalChildren > 0) {
                    for (const elem of node.children) {
                        updateChildrenNum(elem)
                        count += elem.keys
                    }
                } else {
                    count += 1
                }
                node.keys = count
                // node.children = sortBy(node.children, 'label')
            }

            const keyStruct = []
            const mark = {}
            for (const key in keys) {
                const keyPart = split(key, separator)
                // const prefixLen = size(keyPart) - 1
                const len = size(keyPart)
                let handlePath = ''
                let ks = keyStruct
                for (let i = 0; i < len; i++) {
                    handlePath += keyPart[i]
                    if (i !== len - 1) {
                        // layer
                        const treeKey = `${handlePath}@${ConnectionType.RedisKey}`
                        if (!mark.hasOwnProperty(treeKey)) {
                            mark[treeKey] = {
                                key: `${connName}/db${db}/${treeKey}`,
                                label: keyPart[i],
                                name: connName,
                                db,
                                keys: 0,
                                redisKey: handlePath,
                                type: ConnectionType.RedisKey,
                                children: [],
                            }
                            sortedInsertChild(ks, mark[treeKey])
                        }
                        ks = mark[treeKey].children
                        handlePath += separator
                    } else {
                        // key
                        const treeKey = `${handlePath}@${ConnectionType.RedisValue}`
                        mark[treeKey] = {
                            key: `${connName}/db${db}/${treeKey}`,
                            label: keyPart[i],
                            name: connName,
                            db,
                            keys: 0,
                            redisKey: handlePath,
                            type: ConnectionType.RedisValue,
                        }
                        sortedInsertChild(ks, mark[treeKey])
                    }
                }
            }

            // append db node to current connection's children
            const connNode = this.getConnection(connName)
            const { children = [] } = connNode
            children[db].children = keyStruct
            children[db].opened = true
            updateChildrenNum(children[db])
        },

        /**
         * select node
         * @param key
         * @param name
         * @param db
         * @param type
         * @param redisKey
         */
        select({ key, name, db, type, redisKey }) {
            if (type === ConnectionType.RedisValue) {
                console.log(`[click]key:${key} db: ${db} redis key: ${redisKey}`)

                // async get value for key
                this.loadKeyValue(name, db, redisKey).then(() => {})
            }
        },

        /**
         * load redis key
         * @param server
         * @param db
         * @param key
         */
        async loadKeyValue(server, db, key) {
            try {
                const { data, success, msg } = await GetKeyValue(server, db, key)
                if (success) {
                    const { type, ttl, value } = data
                    const tab = useTabStore()
                    tab.upsertTab({
                        server,
                        db,
                        type,
                        ttl,
                        key,
                        value,
                    })
                } else {
                    console.warn('TODO: handle get key fail')
                }
            } finally {
            }
        },

        /**
         *
         * @param {string} connName
         * @param {number} db
         * @param {string} key
         * @private
         */
        _addKey(connName, db, key) {
            const connNode = this.getConnection(connName)
            const { children: dbs = [] } = connNode
            const dbDetail = get(dbs, db, {})

            if (dbDetail == null) {
                return
            }

            const descendantChain = [dbDetail]

            const keyPart = split(key, separator)
            let redisKey = ''
            const keyLen = size(keyPart)
            let added = false
            for (let i = 0; i < keyLen; i++) {
                redisKey += keyPart[i]

                const node = last(descendantChain)
                const nodeList = get(node, 'children', [])
                const len = size(nodeList)
                const isLastKeyPart = i === keyLen - 1
                for (let j = 0; j < len + 1; j++) {
                    const treeKey = get(nodeList[j], 'key')
                    const isLast = j >= len - 1
                    const currentKey = `${connName}/db${db}/${redisKey}@${
                        isLastKeyPart ? ConnectionType.RedisValue : ConnectionType.RedisKey
                    }`
                    if (treeKey > currentKey || isLast) {
                        // out of search range, add new item
                        if (isLastKeyPart) {
                            // key not exists, add new one
                            const item = {
                                key: currentKey,
                                label: keyPart[i],
                                name: connName,
                                db,
                                keys: 1,
                                redisKey,
                                type: ConnectionType.RedisValue,
                            }
                            if (isLast) {
                                nodeList.push(item)
                            } else {
                                nodeList.splice(j, 0, item)
                            }
                            added = true
                        } else {
                            // layer not exists, add new one
                            const item = {
                                key: currentKey,
                                label: keyPart[i],
                                name: connName,
                                db,
                                keys: 0,
                                redisKey,
                                type: ConnectionType.RedisKey,
                                children: [],
                            }
                            if (isLast) {
                                nodeList.push(item)
                                descendantChain.push(last(nodeList))
                            } else {
                                nodeList.splice(j, 0, item)
                                descendantChain.push(nodeList[j])
                            }
                            redisKey += separator
                            added = true
                        }
                        break
                    } else if (treeKey === currentKey) {
                        if (isLastKeyPart) {
                            // same key exists, do nothing
                            console.log('TODO: same key exist, do nothing now, should replace value later')
                        } else {
                            // same group exists, find into it's children
                            descendantChain.push(nodeList[j])
                            redisKey += separator
                        }
                        break
                    }
                }
            }

            // update ancestor node's info
            if (added) {
                const desLen = size(descendantChain)
                for (let i = 0; i < desLen; i++) {
                    const children = get(descendantChain[i], 'children', [])
                    let keys = 0
                    for (const child of children) {
                        if (child.type === ConnectionType.RedisKey) {
                            keys += get(child, 'keys', 1)
                        } else if (child.type === ConnectionType.RedisValue) {
                            keys += get(child, 'keys', 0)
                        }
                    }
                    descendantChain[i].keys = keys
                }
            }
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
                    this._addKey(connName, db, key)
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
        _removeKey(connName, db, key) {
            const connNode = this.getConnection(connName)
            const { children: dbs = [] } = connNode
            const dbDetail = get(dbs, db, {})

            if (dbDetail == null) {
                return
            }

            const descendantChain = [dbDetail]
            const keyPart = split(key, separator)
            let redisKey = ''
            const keyLen = size(keyPart)
            let deleted = false
            let forceBreak = false
            for (let i = 0; i < keyLen && !forceBreak; i++) {
                redisKey += keyPart[i]

                const node = last(descendantChain)
                const nodeList = get(node, 'children', [])
                const len = size(nodeList)
                const isLastKeyPart = i === keyLen - 1
                for (let j = 0; j < len; j++) {
                    const treeKey = get(nodeList[j], 'key')
                    const currentKey = `${connName}/db${db}/${redisKey}@${
                        isLastKeyPart ? ConnectionType.RedisValue : ConnectionType.RedisKey
                    }`
                    if (treeKey > currentKey) {
                        // out of search range, target not exists
                        forceBreak = true
                        break
                    } else if (treeKey === currentKey) {
                        if (isLastKeyPart) {
                            // find target
                            nodeList.splice(j, 1)
                            node.keys -= 1
                            deleted = true
                            forceBreak = true
                        } else {
                            // find into it's children
                            descendantChain.push(nodeList[j])
                            redisKey += separator
                        }
                        break
                    }
                }

                if (forceBreak) {
                    break
                }
            }
            // console.log(JSON.stringify(descendantChain))

            // update ancestor node's info
            if (deleted) {
                const desLen = size(descendantChain)
                for (let i = desLen - 1; i > 0; i--) {
                    const children = get(descendantChain[i], 'children', [])
                    const parent = descendantChain[i - 1]
                    if (isEmpty(children)) {
                        const parentChildren = get(parent, 'children', [])
                        const k = get(descendantChain[i], 'key')
                        remove(parentChildren, (item) => item.key === k)
                    }
                    parent.keys -= 1
                }
            }
        },

        /**
         * remove redis key
         * @param {string} connName
         * @param {number} db
         * @param {string} key
         * @returns {Promise<boolean>}
         */
        async removeKey(connName, db, key) {
            try {
                const { data, success, msg } = await RemoveKey(connName, db, key)
                if (success) {
                    // update tree view data
                    this._removeKey(connName, db, key)

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
                this._removeKey(connName, db, key)
                this._addKey(connName, db, newKey)
                return { success: true }
            } else {
                return { success: false, msg }
            }
        },
    },
})

export default useDatabaseStore
