/**
 * redis database item
 */
export class RedisDatabaseItem {
    constructor({ db = 0, alias = '', keyCount = 0, maxKeys = 0 }) {
        this.db = db
        this.alias = alias
        this.keyCount = keyCount
        this.maxKeys = maxKeys
    }
}
