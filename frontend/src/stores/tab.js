import { assign, find, findIndex, get, indexOf, isEmpty, pullAt, remove, set, size } from 'lodash'
import { defineStore } from 'pinia'
import { TabItem } from '@/objects/tabItem.js'
import { decodeRedisKey } from '@/utils/key_convert.js'

const useTabStore = defineStore('tab', {
    /**
     * @typedef {Object} ListEntryItem
     * @property {string|number[]} v value
     * @property {string} [dv] display value
     */

    /**
     * @typedef {Object} ListReplaceItem
     * @property {number} index
     * @property {string|number[]} v value
     * @property {string} [dv] display value
     */

    /**
     * @typedef {Object} HashEntryItem
     * @property {string} k field name
     * @property {string|number[]} v value
     * @property {string} [dv] display value
     */

    /**
     * @typedef {Object} HashReplaceItem
     * @property {string|number[]} k field name
     * @property {string|number[]} nk new field name
     * @property {string|number[]} v value
     * @property {string} [dv] display value
     */

    /**
     * @typedef {Object} SetEntryItem
     * @property {string|number[]} v value
     * @property {string} [dv] display value
     */

    /**
     * @typedef {Object} ZSetEntryItem
     * @property {number} s score
     * @property {string|number[]} v value
     * @property {string} [dv] display value
     */

    /**
     * @typedef {Object} ZSetReplaceItem
     * @property {number} s score
     * @property {string|number[]} v value
     * @property {string|number[]} nv new value
     * @property {string} [dv] display value
     */

    /**
     * @typedef {Object} StreamEntryItem
     * @property {string} id
     * @property {Object.<string, *>} v value
     * @property {string} [dv] display value
     */

    /**
     * @typedef {Object} TabState
     * @property {string} nav
     * @property {number} asideWidth
     * @property {TabItem[]} tabList
     * @property {number} activatedIndex
     */

    /**
     *
     * @returns {TabState}
     */
    state: () => ({
        nav: 'server',
        asideWidth: 300,
        tabList: [],
        activatedIndex: 0, // current activated tab index
    }),
    getters: {
        /**
         * get current tab list item
         * @returns {TabItem[]}
         */
        tabs() {
            // if (isEmpty(this.tabList)) {
            //     this.newBlankTab()
            // }
            return this.tabList
        },

        /**
         * get current activated tab item
         * @returns {TabItem|null}
         */
        currentTab() {
            return get(this.tabs, this.activatedIndex)
        },

        currentTabName() {
            return get(this.tabs, [this.activatedIndex, 'name'])
        },

        currentSelectedKeys() {
            const tab = this.currentTab
            return get(tab, 'selectedKeys', [])
        },

        currentCheckedKeys() {
            const tab = this.currentTab
            return get(tab, 'checkedKeys', [])
        },
    },
    actions: {
        /**
         *
         * @param idx
         * @param {boolean} [switchNav]
         * @param {string} [subTab]
         * @private
         */
        _setActivatedIndex(idx, switchNav, subTab) {
            this.activatedIndex = idx
            if (switchNav === true) {
                this.nav = idx >= 0 ? 'browser' : 'server'
                if (!isEmpty(subTab)) {
                    set(this.tabList, [idx, 'subTab'], subTab)
                }
            } else {
                if (idx < 0) {
                    this.nav = 'server'
                }
            }
        },

        openBlank(server) {
            this.upsertTab({ server, clearValue: true })
        },

        /**
         * update or insert a new tab if not exists with the same name
         * @param {string} subTab
         * @param {string} server
         * @param {number} [db]
         * @param {number} [type]
         * @param {number} [ttl]
         * @param {string} [key]
         * @param {string} [keyCode]
         * @param {number} [size]
         * @param {number} [length]
         * @param {string} [matchPattern]
         * @param {boolean} [clearValue]
         * @param {*} [value]
         */
        upsertTab({ subTab, server, db, type, ttl, key, keyCode, size, length, matchPattern = '', clearValue }) {
            let tabIndex = findIndex(this.tabList, { name: server })
            if (tabIndex === -1) {
                const tabItem = new TabItem({
                    name: server,
                    title: server,
                    subTab,
                    server,
                    db,
                    type,
                    ttl,
                    key,
                    keyCode,
                    size,
                    length,
                    matchPattern,
                    value: undefined,
                })
                this.tabList.push(tabItem)
                tabIndex = this.tabList.length - 1
            } else {
                const tab = this.tabList[tabIndex]
                tab.blank = false
                tab.subTab = subTab
                // tab.title = db !== undefined ? `${server}/db${db}` : `${server}`
                tab.title = server
                tab.server = server
                tab.db = db == null ? tab.db : db
                tab.type = type
                tab.ttl = ttl
                tab.key = key
                tab.keyCode = keyCode
                tab.size = size
                tab.length = length
                tab.matchPattern = matchPattern
                if (clearValue === true) {
                    tab.value = undefined
                }
            }
            this._setActivatedIndex(tabIndex, true, subTab)
        },

        /**
         * keep update value in tab
         * @param {string} server
         * @param {number} db
         * @param {string} key
         * @param {*} [value]
         * @param {string} [format]
         * @param {string] [decode]
         * @param {string} [matchPattern]
         * @param {boolean} [reset]
         * @param {boolean} [end] keep end status if not set
         * @param {number} [size]
         * @param {number} [length]
         */
        updateValue({ server, db, key, value, format, decode, matchPattern, reset, end, size = -1, length = -1 }) {
            const tabData = find(this.tabList, { name: server, db, key })
            if (tabData == null) {
                return
            }

            tabData.format = format || tabData.format
            tabData.decode = decode || tabData.decode
            tabData.matchPattern = matchPattern || ''
            if (size >= 0) {
                tabData.size = size
            }
            if (length >= 0) {
                tabData.length = length
            }
            if (typeof end === 'boolean') {
                tabData.end = end
            }
            if (!!!reset && typeof value === 'object') {
                if (value instanceof Array) {
                    tabData.value = tabData.value || []
                    tabData.value.push(...value)
                } else {
                    tabData.value = assign(value, tabData.value || {})
                }
            } else {
                tabData.value = value
            }
        },

        /**
         * insert entries
         * @param {string} server
         * @param {number} db
         * @param {string|number[]} key
         * @param {string} type
         * @param {ListEntryItem[]|HashEntryItem[]|SetEntryItem[]|ZSetEntryItem[]|StreamEntryItem[]} entries
         * @param {boolean} [prepend] for list only
         */
        insertValueEntries({ server, db, key, type, entries, prepend }) {
            const tab = find(this.tabList, { name: server, db, key })
            if (tab == null) {
                return
            }

            switch (type.toLowerCase()) {
                case 'list': // {v:string, dv:[string]}[]
                    tab.value = tab.value || []
                    if (prepend === true) {
                        tab.value = [...entries, ...tab.value]
                    } else {
                        tab.value.push(...entries)
                    }
                    tab.length += size(entries)
                    break

                case 'hash': // {k:string, v:string, dv:[string]}[]
                case 'set': // {v: string, s: number}[]
                case 'zset': // {v: string, s: number}[]
                    tab.value = tab.value || []
                    tab.value.push(...entries)
                    tab.length += size(entries)
                    break

                case 'stream': // {id: string, v: {}}[]
                    tab.value = tab.value || []
                    tab.value = [...entries, ...tab.value]
                    tab.length += size(entries)
                    break
            }
        },

        /**
         * update entries' value
         * @param {string} server
         * @param {number} db
         * @param {string|number[]} key
         * @param {string} type
         * @param {ListEntryItem[]|HashEntryItem[]|SetEntryItem[]|ZSetEntryItem[]|StreamEntryItem[]} entries
         */
        updateValueEntries({ server, db, key, type, entries }) {
            const tab = find(this.tabList, { name: server, db, key })
            if (tab == null) {
                return
            }

            switch (type.toLowerCase()) {
                case 'hash': // {k:string, v:string, dv:string}[]
                    tab.value = tab.value || []
                    for (const entry of entries) {
                        let updated = false
                        for (const val of tab.value) {
                            if (val.k === entry.k) {
                                val.v = entry.v
                                val.dv = entry.dv
                                updated = true
                                break
                            }
                        }
                        if (!updated) {
                            // no match element, append
                            tab.value.push(entry)
                            tab.length += 1
                        }
                    }
                    break

                case 'zset': // {s:number, v:string, dv:string}[]
                    tab.value = tab.value || []
                    for (const entry of entries) {
                        let updated = false
                        for (const val of tab.value) {
                            if (val.v === entry.v) {
                                val.s = entry.s
                                val.dv = entry.dv
                                updated = true
                                break
                            }
                        }
                        if (!updated) {
                            // no match element, append
                            tab.value.push(entry)
                            tab.length += 1
                        }
                    }
                    break
            }
        },

        /**
         * replace entry item key or field in value(modify the index key)
         * @param {string} server
         * @param {number} db
         * @param {string|number[]} key
         * @param {string} type
         * @param {ListReplaceItem[]|HashReplaceItem[]|ZSetReplaceItem[]} entries
         * @param {number[]} [index] indexes for replacement, can improve search efficiency if configured
         */
        replaceValueEntries({ server, db, key, type, entries, index }) {
            const tab = find(this.tabList, { name: server, db, key })
            if (tab == null) {
                return
            }

            switch (type.toLowerCase()) {
                case 'list': // ListReplaceItem[]
                    tab.value = tab.value || []
                    for (const entry of entries) {
                        if (size(tab.value) > entry.index) {
                            tab.value[entry.index] = {
                                v: entry.v,
                                dv: entry.dv,
                            }
                        } else {
                            // out of range, append
                            tab.value.push(entry)
                            tab.length += 1
                        }
                    }
                    break

                case 'hash': // HashReplaceItem[]
                    tab.value = tab.value || []
                    for (const idx of index) {
                        const entry = get(tab.value, idx)
                        if (entry != null) {
                            /** @type HashReplaceItem[] **/
                            const replaceEntry = remove(entries, (e) => e.k === entry.k)
                            if (!isEmpty(replaceEntry)) {
                                entry.k = replaceEntry[0].nk
                                entry.v = replaceEntry[0].v
                                entry.dv = replaceEntry[0].dv
                            }
                        }
                    }

                    // the left entries do not included in index list, try to retrieve the whole list
                    for (const entry of entries) {
                        let updated = false
                        for (const val of tab.value) {
                            if (val.k === entry.k) {
                                val.k = entry.nk
                                val.v = entry.v
                                val.dv = entry.dv
                                updated = true
                                break
                            }
                        }
                        if (!updated) {
                            // no match element, append
                            tab.value.push({
                                k: entry.nk,
                                v: entry.v,
                                dv: entry.dv,
                            })
                            tab.length += 1
                        }
                    }
                    break

                case 'zset': // ZSetReplaceItem[]
                    tab.value = tab.value || []
                    for (const idx of index) {
                        const entry = get(tab.value, idx)
                        if (entry != null) {
                            /** @type ZSetReplaceItem[] **/
                            const replaceEntry = remove(entries, ({ v }) => v === entry.k)
                            if (!isEmpty(replaceEntry)) {
                                entry.s = replaceEntry[0].s
                                entry.v = replaceEntry[0].nv
                                entry.dv = replaceEntry[0].dv
                            }
                        }
                    }

                    // the left entries do not included in index list, try to retrieve the whole list
                    for (const entry of entries) {
                        let updated = false
                        for (const val of tab.value) {
                            if (val.v === entry.v) {
                                val.s = entry.s
                                val.v = entry.nv
                                val.dv = entry.dv
                                updated = true
                                break
                            }
                        }
                        if (!updated) {
                            // no match element, append
                            tab.value.push({
                                s: entry.s,
                                v: entry.nv,
                                dv: entry.dv,
                            })
                            tab.length += 1
                        }
                    }
                    break
            }
        },

        /**
         * remove value entries
         * @param {string} server
         * @param {number} db
         * @param {string} key
         * @param {string} type
         * @param {string[] | number[]} entries
         */
        removeValueEntries({ server, db, key, type, entries }) {
            const tab = find(this.tabList, { name: server, db, key })
            if (tab == null) {
                return
            }

            switch (type.toLowerCase()) {
                case 'list': // string[] | number[]
                    tab.value = tab.value || []
                    if (typeof entries[0] === 'number') {
                        // remove by index, sort by desc first
                        entries.sort((a, b) => b - a)
                        const removed = pullAt(tab.value, ...entries)
                        tab.length -= size(removed)
                    } else {
                        // append or prepend items
                        for (const elem of entries) {
                            if (!isEmpty(remove(tab.value, elem))) {
                                tab.length -= 1
                            }
                        }
                    }
                    break

                case 'hash': // string[]
                    tab.value = tab.value || {}
                    for (const k of entries) {
                        for (let i = 0; i < tab.value.length; i++) {
                            if (tab.value[i].k === k) {
                                tab.value.splice(i, 1)
                                tab.length -= 1
                                break
                            }
                        }
                    }
                    break

                case 'set': // string[]
                case 'zset': // string[]
                    tab.value = tab.value || []
                    for (const v of entries) {
                        for (let i = 0; i < tab.value.length; i++) {
                            if (tab.value[i].v === v) {
                                tab.value.splice(i, 1)
                                tab.length -= 1
                                break
                            }
                        }
                    }
                    break

                case 'stream': // string[]
                    tab.value = tab.value || []
                    for (const id of entries) {
                        for (let i = 0; i < tab.value.length; i++) {
                            if (tab.value[i].id === id) {
                                tab.value.splice(i, 1)
                                tab.length -= 1
                                break
                            }
                        }
                    }
                    break
            }
        },

        /**
         * update loading status of content in tab
         * @param {string} server
         * @param {number} db
         * @param {boolean} loading
         */
        updateLoading({ server, db, loading }) {
            const tab = find(this.tabList, { name: server, db })
            if (tab == null) {
                return
            }

            tab.loading = loading
        },

        /**
         * update ttl in tab
         * @param {string} server
         * @param {number} db
         * @param {string|number[]} key
         * @param {number} ttl
         */
        updateTTL({ server, db, key, ttl }) {
            let tab = find(this.tabList, { name: server, db, key: decodeRedisKey(key) })
            if (tab == null) {
                return
            }
            tab.ttl = ttl
        },

        /**
         * set tab's content to empty
         * @param {string} name
         */
        emptyTab(name) {
            const tab = find(this.tabList, { name })
            if (tab != null) {
                tab.key = null
                tab.value = null
            }
        },
        switchTab(tabIndex) {
            // const len = size(this.tabList)
            // if (tabIndex < 0 || tabIndex >= len) {
            //     tabIndex = 0
            // }
            // this.activatedIndex = tabIndex
            // const tabIndex = findIndex(this.tabList, {name})
            // if (tabIndex === -1) {
            //     return
            // }
            // this.activatedIndex = tabIndex
        },

        switchSubTab(name) {
            const tab = this.currentTab
            if (tab == null) {
                return
            }
            tab.subTab = name
        },

        /**
         *
         * @param {number} tabIndex
         * @returns {*|null}
         */
        removeTab(tabIndex) {
            const len = size(this.tabs)
            // ignore remove last blank tab
            if (len === 1 && this.tabs[0].blank) {
                return null
            }

            if (tabIndex < 0 || tabIndex >= len) {
                return null
            }
            const removed = this.tabList.splice(tabIndex, 1)

            // update select index if removed index equal current selected
            this.activatedIndex -= 1
            if (this.activatedIndex < 0) {
                if (this.tabList.length > 0) {
                    this._setActivatedIndex(0, false)
                } else {
                    this._setActivatedIndex(-1, false)
                }
            } else {
                this._setActivatedIndex(this.activatedIndex, false)
            }

            return size(removed) > 0 ? removed[0] : null
        },

        /**
         *
         * @param {string} tabName
         */
        removeTabByName(tabName) {
            const idx = findIndex(this.tabs, { name: tabName })
            if (idx !== -1) {
                this.removeTab(idx)
            }
        },

        /**
         *
         */
        removeAllTab() {
            this.tabList = []
            this._setActivatedIndex(-1, false)
        },

        /**
         * set expanded keys for server
         * @param {string} server
         * @param {string[]} keys
         */
        setExpandedKeys(server, keys = []) {
            /** @type TabItem**/
            let tab = find(this.tabList, { name: server })
            if (tab != null) {
                if (isEmpty(keys)) {
                    tab.expandedKeys = []
                } else {
                    tab.expandedKeys = keys
                }
            }
        },

        /**
         *
         * @param {string} server
         * @param {string} key
         */
        addExpandedKey(server, key) {
            /** @type TabItem**/
            let tab = find(this.tabList, { name: server })
            if (tab != null) {
                tab.expandedKeys.push(key)
            }
        },

        /**
         *
         * @param {string} server
         * @param {string} key
         */
        toggleExpandKey(server, key) {
            /** @type TabItem**/
            let tab = find(this.tabList, { name: server })
            if (tab != null) {
                const idx = indexOf(tab.expandedKeys, key)
                if (idx === -1) {
                    tab.expandedKeys.push(key)
                } else {
                    tab.expandedKeys.splice(idx, 1)
                }
            }
        },

        /**
         *
         * @param {string} server
         * @param {string} key
         */
        removeExpandedKey(server, key) {
            /** @type TabItem**/
            let tab = find(this.tabList, { name: server })
            if (tab != null) {
                remove(tab.expandedKeys, (v) => v === key)
            }
        },

        /**
         * set selected keys for server
         * @param {string} server
         * @param {string|string[]} [keys]
         */
        setSelectedKeys(server, keys = null) {
            /** @type TabItem**/
            let tab = find(this.tabList, { name: server })
            if (tab != null) {
                if (keys == null) {
                    // select nothing
                    tab.selectedKeys = []
                } else if (typeof keys === 'string') {
                    tab.selectedKeys = [keys]
                } else {
                    tab.selectedKeys = keys
                }
            }
        },

        /**
         * get checked keys
         * @param server
         * @returns {CheckedKey[]}
         */
        getCheckedKeys(server) {
            /** @type TabItem**/
            let tab = find(this.tabList, { name: server })
            if (tab != null) {
                return tab.checkedKeys || []
            }
            return []
        },

        /**
         * set checked keys for server
         * @param {string} server
         * @param {CheckedKey[]} [keys]
         */
        setCheckedKeys(server, keys = null) {
            /** @type TabItem**/
            let tab = find(this.tabList, { name: server })
            if (tab != null) {
                if (isEmpty(keys)) {
                    // select nothing
                    tab.checkedKeys = []
                } else {
                    tab.checkedKeys = keys
                }
            }
        },
    },
})

export default useTabStore
