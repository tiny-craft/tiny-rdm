import { isEmpty, remove, size, sumBy } from 'lodash'
import { ConnectionType } from '@/consts/connection_type.js'

/**
 * redis node item in tree view
 */
export class RedisNodeItem {
    /**
     *
     * @param {string} key - tree node unique key
     * @param {string} label
     * @param {string} [name] - server name, type != ConnectionType.Group only
     * @param {ConnectionType} type
     * @param {number} [db] - database index, type == ConnectionType.RedisDB only
     * @param {string} [redisKey] - redis key, type == ConnectionType.RedisKey || type == ConnectionType.RedisValue only
     * @param {number[]} [redisKeyCode] - redis key char code array, optional for redis key which contains binary data
     * @param {number} [keyCount] - children key count
     * @param {number} [maxKeys] - max key count for database
     * @param {boolean} [isLeaf]
     * @param {boolean} [opened] - redis db is opened, type == ConnectionType.RedisDB only
     * @param {boolean} [expanded] - current node is expanded
     * @param {RedisNodeItem[]} [children]
     * @param {string} [redisType] - redis type name, 'loading' indicate that is in loading progress
     */
    constructor({
        key,
        label,
        name,
        type,
        db = 0,
        redisKey,
        redisKeyCode,
        keyCount = 0,
        maxKeys = 0,
        isLeaf = false,
        opened = false,
        expanded = false,
        children,
        redisType,
    }) {
        this.key = key
        this.label = label
        this.name = name
        this.type = type
        this.db = db
        this.redisKey = redisKey
        this.redisKeyCode = redisKeyCode
        this.keyCount = keyCount
        this.maxKeys = maxKeys
        this.isLeaf = isLeaf
        this.opened = opened
        this.expanded = expanded
        this.children = children
        this.redisType = redisType
    }

    /**
     * sort node list
     * @param {RedisNodeItem[]} nodeList
     * @private
     */
    _sortNodes(nodeList) {
        if (nodeList == null) {
            return
        }
        nodeList.sort((a, b) => {
            return a.key > b.key ? 1 : -1
        })
    }

    /**
     * compare two items to determine the sort order
     * @param {*} a
     * @param {*} b
     * @return {number}
     */
    _sortingCompare(a, b) {
        if (a.type !== b.type) {
            return a.type - b.type
        }
        const isANum = isNaN(a.label)
        const isBNum = isNaN(b.label)
        if (!isANum && !isBNum) {
            return parseInt(a.label, 10) - parseInt(b.label, 10)
        } else if (!isANum) {
            return -1
        } else if (!isBNum) {
            return 1
        }
        return a.label.localeCompare(b.label)
    }

    /**
     * calculate insert sorted index
     * @param {[]} arr
     * @param {*} item
     * @return {number}
     */
    _sortedIndex(arr, item) {
        for (let i = 0; i < arr.length; i++) {
            const cmpResult = this._sortingCompare(arr[i], item)
            if (cmpResult > 0) {
                return i
            } else if (cmpResult === 0) {
                return i + 1
            }
        }
        return arr.length
    }

    /**
     * sort all node item's children and calculate keys count
     * @param skipSort skip sorting children
     * @returns {boolean} return whether key count changed
     */
    tidy(skipSort) {
        if (this.type === ConnectionType.RedisValue) {
            this.keyCount = 1
        } else if (this.type === ConnectionType.RedisKey || this.type === ConnectionType.RedisDB) {
            let keyCount = 0
            if (!isEmpty(this.children)) {
                if (!!!skipSort) {
                    this.sortChildren()
                }
                for (const child of this.children) {
                    child.tidy(skipSort)
                    keyCount += child.keyCount
                }
            } else {
                keyCount = 0
            }
            if (this.keyCount !== keyCount) {
                this.keyCount = keyCount
                return true
            }
        }
        return false
    }

    sortChildren() {
        this.children.sort(this._sortingCompare)
    }

    /**
     *
     * @param {RedisNodeItem} child
     * @param {boolean} [sorted]
     */
    addChild(child, sorted) {
        if (!!!sorted) {
            this.children.push(child)
        } else {
            const idx = this._sortedIndex(this.children, child)
            this.children.splice(idx, 0, child)
        }
    }

    /**
     *
     * @param {{}} predicate
     * @return {number}
     */
    removeChild(predicate) {
        if (this.type !== ConnectionType.RedisKey && this.type !== ConnectionType.RedisDB) {
            return 0
        }
        const removed = remove(this.children, predicate)
        return size(removed)
    }

    getChildren() {
        return this.children
    }

    reCalcKeyCount() {
        if (this.type === ConnectionType.RedisValue) {
            this.keyCount = 1
        }
        this.keyCount = sumBy(this.children, (c) => c.keyCount)
        return this.keyCount
    }
}
