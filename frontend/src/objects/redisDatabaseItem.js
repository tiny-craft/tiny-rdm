/**
 * redis database item
 */
export class RedisDatabaseItem {
    constructor({ db = 0, keyCount = 0, maxKeys = 0 }) {
        this.db = db
        this.keyCount = keyCount
        this.maxKeys = maxKeys
    }
}
