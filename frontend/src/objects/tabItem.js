/**
 * tab item
 */
export class TabItem {
    /**
     * @typedef {Object} CheckedKey
     * @property {string} key
     * @property {string} [redisKey]
     */

    /**
     *
     * @param {string} name connection name
     * @param {string} title tab title
     * @param {boolean} blank is blank tab
     * @param {string} subTab secondary tab value
     * @param {string} [title] tab title
     * @param {string} [icon] tab icon
     * @param {string[]} selectedKeys
     * @param {CheckedKey[]} checkedKeys
     * @param {string} [type] key type
     * @param {*} [value] key value
     * @param {string} [server] server name
     * @param {int} [db] database index
     * @param {string} [key] current key name
     * @param {number[]|null|undefined} [keyCode] current key name as char array
     * @param {number} [size] memory usage
     * @param {number} [length] length of content or entries
     * @param {int} [ttl] ttl of current key
     * @param {string} [decode]
     * @param {string} [format]
     * @param {string} [matchPattern]
     * @param {boolean} [end]
     * @param {boolean} [loading]
     */
    constructor({
        name,
        title,
        blank,
        subTab,
        icon,
        selectedKeys,
        checkedKeys,
        type,
        value,
        server,
        db = 0,
        key,
        keyCode,
        size = 0,
        length = 0,
        ttl = 0,
        decode = '',
        format = '',
        matchPattern = '',
        end = false,
        loading = false,
    }) {
        this.name = name
        this.title = title
        this.blank = blank
        this.subTab = subTab
        this.icon = icon
        this.selectedKeys = selectedKeys
        this.checkedKeys = checkedKeys
        this.type = type
        this.value = value
        this.server = server
        this.db = db
        this.key = key
        this.keyCode = keyCode
        this.size = size
        this.length = length
        this.ttl = ttl
        this.decode = decode
        this.format = format
        this.matchPattern = matchPattern
        this.end = end
        this.loading = loading
    }
}
