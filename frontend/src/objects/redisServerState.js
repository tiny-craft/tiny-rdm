import { initial, isEmpty, join, last, mapValues, size, slice, sortBy, split, toUpper } from 'lodash'
import useConnectionStore from 'stores/connections.js'
import { ConnectionType } from '@/consts/connection_type.js'
import { RedisDatabaseItem } from '@/objects/redisDatabaseItem.js'
import { KeyViewType } from '@/consts/key_view_type.js'
import { RedisNodeItem } from '@/objects/redisNodeItem.js'
import { decodeRedisKey, nativeRedisKey } from '@/utils/key_convert.js'

/**
 * server connection state
 */
export class RedisServerState {
    /**
     * @typedef {Object} LoadingState
     * @property {boolean} loading indicated that is loading children now
     * @property {boolean} fullLoaded indicated that all children already loaded
     */

    /**
     * @param {string} name server name
     * @param {number} db current opened database
     * @param {{}} stats current server status info
     * @param {Object.<number, RedisDatabaseItem>} databases database list
     * @param {string|null} patternFilter pattern filter
     * @param {string|null} typeFilter redis type filter
     * @param {LoadingState} loadingState all loading state in opened connections map by server and LoadingState
     * @param {KeyViewType} viewType view type selection for all opened connections group by 'server'
     * @param {Map<string, RedisNodeItem>} nodeMap map nodes by "type#key"
     */
    constructor({
        name,
        db = 0,
        stats = {},
        databases = {},
        patternFilter = null,
        typeFilter = null,
        loadingState = {},
        viewType = KeyViewType.Tree,
        nodeMap = new Map(),
    }) {
        this.name = name
        this.db = db
        this.stats = stats
        this.databases = databases
        this.patternFilter = patternFilter
        this.typeFilter = typeFilter
        this.loadingState = loadingState
        this.viewType = viewType
        this.nodeMap = nodeMap
        this.getRoot()

        const connStore = useConnectionStore()
        const { keySeparator } = connStore.getDefaultSeparator(name)
        this.separator = isEmpty(keySeparator) ? ':' : keySeparator
    }

    dispose() {
        this.stats = {}
        this.patternFilter = null
        this.typeFilter = null
        this.nodeMap.clear()
    }

    closeDatabase() {
        this.patternFilter = null
        this.typeFilter = null
        this.nodeMap.clear()
    }

    setDatabaseKeyCount(db, maxKeys) {
        const dbInst = this.databases[db]
        if (dbInst == null) {
            this.databases[db] = new RedisDatabaseItem({ db, maxKeys })
        } else {
            dbInst.maxKeys = maxKeys
        }
        return dbInst
    }

    /**
     * update max key by increase/decrease value
     * @param {number} db
     * @param {number} updateVal
     */
    updateDBKeyCount(db, updateVal) {
        const dbInst = this.databases[this.db]
        if (dbInst != null) {
            dbInst.maxKeys = Math.max(0, dbInst.maxKeys + updateVal)
        }
    }

    /**
     * set db max keys value
     * @param {number} db
     * @param {number} count
     */
    setDBKeyCount(db, count) {
        const dbInst = this.databases[db]
        if (dbInst != null) {
            dbInst.maxKeys = Math.max(0, count)
        }
    }

    /**
     * get tree root item
     * @returns {RedisNodeItem}
     */
    getRoot() {
        const rootKey = `${ConnectionType.RedisDB}`
        let root = this.nodeMap.get(rootKey)
        if (root == null) {
            // create root node
            root = new RedisNodeItem({
                key: rootKey,
                label: this.separator,
                type: ConnectionType.RedisDB,
                children: [],
            })
            this.nodeMap.set(rootKey, root)
        }
        return root
    }

    /**
     * get database list sort by db asc
     * @return {RedisDatabaseItem[]}
     */
    getDatabase() {
        return sortBy(mapValues(this.databases), 'db')
    }

    /**
     *
     * @param {ConnectionType} type
     * @param {string} keyPath
     * @param {RedisNodeItem} node
     */
    addNode(type, keyPath, node) {
        this.nodeMap.set(`${type}/${keyPath}`, node)
    }

    /**
     * add keys to current opened database
     * @param {Array<string|number[]>|Set<string|number[]>} keys
     * @param {boolean} [sortInsert]
     * @return {{newKey: number, newLayer: number, success: boolean, replaceKey: number}}
     */
    addKeyNodes(keys, sortInsert) {
        const result = {
            success: false,
            newLayer: 0,
            newKey: 0,
            replaceKey: 0,
        }
        const root = this.getRoot()

        if (this.viewType === KeyViewType.List) {
            // construct list view data
            for (const key of keys) {
                const k = decodeRedisKey(key)
                const isBinaryKey = k !== key
                const nodeKey = `${ConnectionType.RedisValue}/${nativeRedisKey(key)}`
                const replaceKey = this.nodeMap.has(nodeKey)
                const selectedNode = new RedisNodeItem({
                    key: `${this.name}/db${this.db}#${nodeKey}`,
                    label: k,
                    db: this.db,
                    keyCount: 0,
                    redisKey: k,
                    redisKeyCode: isBinaryKey ? key : undefined,
                    redisKeyType: undefined,
                    type: ConnectionType.RedisValue,
                    isLeaf: true,
                })
                this.nodeMap.set(nodeKey, selectedNode)
                if (!replaceKey) {
                    root.addChild(selectedNode, sortInsert)
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
                const keyParts = isBinaryKey ? [nativeRedisKey(key)] : split(k, this.separator)
                const len = size(keyParts)
                const lastIdx = len - 1
                let handlePath = ''
                let node = root
                for (let i = 0; i < len; i++) {
                    handlePath += keyParts[i]
                    if (i !== lastIdx) {
                        // layer
                        const nodeKey = `${ConnectionType.RedisKey}/${handlePath}`
                        let selectedNode = this.nodeMap.get(nodeKey)
                        if (selectedNode == null) {
                            selectedNode = new RedisNodeItem({
                                key: `${this.name}/db${this.db}#${nodeKey}`,
                                label: keyParts[i],
                                db: this.db,
                                keyCount: 0,
                                redisKey: handlePath,
                                type: ConnectionType.RedisKey,
                                isLeaf: false,
                                children: [],
                            })
                            this.nodeMap.set(nodeKey, selectedNode)
                            node.addChild(selectedNode, sortInsert)
                            result.newLayer += 1
                        }
                        node = selectedNode
                        handlePath += this.separator
                    } else {
                        // key
                        const nodeKey = `${ConnectionType.RedisValue}/${handlePath}`
                        const replaceKey = this.nodeMap.has(nodeKey)
                        const selectedNode = new RedisNodeItem({
                            key: `${this.name}/db${this.db}#${nodeKey}`,
                            label: isBinaryKey ? k : keyParts[i],
                            db: this.db,
                            keyCount: 0,
                            redisKey: handlePath,
                            redisKeyCode: isBinaryKey ? key : undefined,
                            redisKeyType: undefined,
                            type: ConnectionType.RedisValue,
                            isLeaf: true,
                        })
                        this.nodeMap.set(nodeKey, selectedNode)
                        if (!replaceKey) {
                            node.addChild(selectedNode, sortInsert)
                            result.newKey += 1
                        } else {
                            result.replaceKey += 1
                        }
                    }
                }
            }
        }
        return result
    }

    /**
     * rename key to a new name
     * @param key
     * @param newKey
     */
    renameKey(key, newKey) {
        const oldLayer = initial(key.split(this.separator)).join(this.separator)
        const newLayer = initial(newKey.split(this.separator)).join(this.separator)
        if (oldLayer !== newLayer) {
            // also change layer
            this.removeKeyNode(key, false)
            const { success } = this.addKeyNodes([newKey], true)
            if (success) {
                this.tidyNode(newLayer)
            }
        } else {
            // change key name only
            const oldNodeKeyName = `${ConnectionType.RedisValue}/${key}`
            const newNodeKeyName = `${ConnectionType.RedisValue}/${newKey}`
            const keyNode = this.nodeMap.get(oldNodeKeyName)
            keyNode.key = `${this.name}/db${this.db}#${newNodeKeyName}`
            keyNode.label = last(split(newKey, this.separator))
            keyNode.redisKey = newKey
            // not support rename binary key name yet
            // keyNode.redisKeyCode = []
            this.nodeMap.set(newNodeKeyName, keyNode)
            this.nodeMap.delete(oldNodeKeyName)
        }
    }

    /**
     * remove key node by key name
     * @param {string} [key]
     * @param {boolean} [isLayer]
     * @return {boolean}
     */
    removeKeyNode(key, isLayer) {
        if (isLayer === true) {
            this.deleteChildrenKeyNodes(key)
        }

        const dbRoot = this.getRoot()
        if (isEmpty(key)) {
            // clear all key nodes
            this.nodeMap.clear()
            this.getRoot()
        } else {
            const keyParts = split(key, this.separator)
            const totalParts = size(keyParts)
            // remove from parent in tree node
            const parentKey = slice(keyParts, 0, totalParts - 1)
            let parentNode
            if (isEmpty(parentKey)) {
                parentNode = dbRoot
            } else {
                parentNode = this.nodeMap.get(`${ConnectionType.RedisKey}/${join(parentKey, this.separator)}`)
            }

            // not found parent node
            if (parentNode == null) {
                return false
            }
            parentNode.removeChild({
                type: isLayer ? ConnectionType.RedisKey : ConnectionType.RedisValue,
                redisKey: key,
            })

            // // check and remove empty layer node
            // let i = totalParts - 1
            // for (; i >= 0; i--) {
            //     const anceKey = join(slice(keyParts, 0, i), this.separator)
            //     if (i > 0) {
            //         const anceNode = this.nodeMap.get(`${ConnectionType.RedisKey}/${anceKey}`)
            //         const redisKey = join(slice(keyParts, 0, i + 1), this.separator)
            //         anceNode.removeChild({ type: ConnectionType.RedisKey, redisKey })
            //
            //         if (isEmpty(anceNode.children)) {
            //             this.nodeMap.delete(`${ConnectionType.RedisKey}/${anceKey}`)
            //         } else {
            //             break
            //         }
            //     } else {
            //         // last one, remove from db node
            //         dbRoot.removeChild({ type: ConnectionType.RedisKey, redisKey: keyParts[0] })
            //         this.nodeMap.delete(`${ConnectionType.RedisValue}/${keyParts[0]}`)
            //     }
            // }
        }

        return true
    }

    /**
     * tidy node by key
     * @param {string} [key]
     * @param {boolean} [skipResort]
     * @return
     */
    tidyNode(key, skipResort) {
        const rootNode = this.getRoot()
        const keyParts = split(key, this.separator)
        const totalParts = size(keyParts)
        let node
        // find last exists ancestor key
        let i = totalParts - 1
        for (; i > 0; i--) {
            const parentKey = join(slice(keyParts, 0, i), this.separator)
            node = this.nodeMap.get(`${ConnectionType.RedisKey}/${parentKey}`)
            if (node != null) {
                break
            }
        }
        if (node == null) {
            node = rootNode
        }
        const keyCountUpdated = node.tidy(skipResort)
        if (keyCountUpdated) {
            // update key count of parent and above
            for (; i > 0; i--) {
                const parentKey = join(slice(keyParts, 0, i), this.separator)
                const parentNode = this.nodeMap.get(`${ConnectionType.RedisKey}/${parentKey}`)
                if (parentNode == null) {
                    break
                }
                const count = parentNode.reCalcKeyCount()
                if (count <= 0) {
                    let anceKeyNode = rootNode
                    // remove from ancestor node
                    if (i > 1) {
                        const anceKey = join(slice(keyParts, 0, i - 1), this.separator)
                        anceKeyNode = this.nodeMap.get(`${ConnectionType.RedisKey}/${anceKey}`)
                    }
                    if (anceKeyNode != null) {
                        anceKeyNode.removeChild({ type: ConnectionType.RedisKey, redisKey: parentKey })
                    }
                }
            }
            // update key count of db
            const dbInst = this.databases[this.db]
            if (dbInst != null) {
                dbInst.keyCount = rootNode.reCalcKeyCount()
            }
        }
    }

    /**
     * add keys to current opened database
     * @param {ConnectionType} type
     * @param {string} keyPath
     * @return {RedisNodeItem|null}
     */
    getNode(type, keyPath) {
        return this.nodeMap.get(`${type}/${keyPath}`) || null
    }

    /**
     * delete node and all it's children from nodeMap
     * @param {string} [key] clean nodeMap if key is empty
     * @private
     */
    deleteChildrenKeyNodes(key) {
        if (isEmpty(key)) {
            this.nodeMap.clear()
            this.getRoot()
        } else {
            const nodeKey = `${ConnectionType.RedisKey}/${key}`
            const node = this.nodeMap.get(nodeKey)
            const children = node.children || []
            for (const child of children) {
                if (child.type === ConnectionType.RedisValue) {
                    if (!this.nodeMap.delete(`${ConnectionType.RedisValue}/${child.redisKey}`)) {
                        console.warn('delete:', `${ConnectionType.RedisValue}/${child.redisKey}`)
                    }
                } else if (child.type === ConnectionType.RedisKey) {
                    this.deleteChildrenKeyNodes(child.redisKey)
                }
            }
            if (!this.nodeMap.delete(nodeKey)) {
                console.warn('delete map key', nodeKey)
            }
        }
    }

    getFilter() {
        let pattern = this.patternFilter
        if (isEmpty(pattern)) {
            const conn = useConnectionStore()
            pattern = conn.getDefaultKeyFilter(this.name)
        }
        return {
            match: pattern,
            type: toUpper(this.typeFilter),
        }
    }

    /**
     * set key filter
     * @param {string} [pattern]
     * @param {string} [type]
     */
    setFilter({ pattern, type }) {
        this.patternFilter = pattern === null ? this.patternFilter : pattern
        this.typeFilter = type === null ? this.typeFilter : type
    }
}
