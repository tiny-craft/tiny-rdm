import { assign, find, findIndex, get, indexOf, isEmpty, pullAt, remove, set, size } from 'lodash'
import { defineStore } from 'pinia'

const useTabStore = defineStore('tab', {
    /**
     * @typedef {Object} TabItem
     * @property {string} name connection name
     * @property {boolean} blank is blank tab
     * @property {string} subTab secondary tab value
     * @property {string} [title] tab title
     * @property {string} [icon] tab icon
     * @property {string[]} selectedKeys
     * @property {string} [type] key type
     * @property {*} [value] key value
     * @property {string} [server] server name
     * @property {int} [db] database index
     * @property {string} [key] current key name
     * @property {number[]|null|undefined} [keyCode] current key name as char array
     * @param {number} [size] memory usage
     * @param {number} [length] length of content or entries
     * @property {int} [ttl] ttl of current key
     * @param {string} [viewAs]
     * @param {string} [decode]
     * @param {boolean} [end]
     * @param {boolean} [loading]
     */

    /**
     *
     * @returns {{tabList: TabItem[], activatedTab: string, activatedIndex: number}}
     */
    state: () => ({
        nav: 'server',
        asideWidth: 300,
        tabList: [],
        activatedTab: '',
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
            // let current = find(this.tabs, {name: this.activatedTab})
            // if (current == null) {
            //     current = this.tabs[0]
            // }
            // return current
        },

        currentSelectedKeys() {
            const tab = this.currentTab()
            return get(tab, 'selectedKeys', [])
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
            this.upsertTab({ server, db: 0 })
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
         * @param {*} [value]
         */
        upsertTab({ subTab, server, db, type, ttl, key, keyCode, size, length }) {
            let tabIndex = findIndex(this.tabList, { name: server })
            if (tabIndex === -1) {
                this.tabList.push({
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
                    value: undefined,
                })
                tabIndex = this.tabList.length - 1
            } else {
                const tab = this.tabList[tabIndex]
                tab.blank = false
                tab.subTab = subTab
                // tab.title = db !== undefined ? `${server}/db${db}` : `${server}`
                tab.title = server
                tab.server = server
                tab.db = db
                tab.type = type
                tab.ttl = ttl
                tab.key = key
                tab.keyCode = keyCode
                tab.size = size
                tab.length = length
                tab.value = undefined
            }
            this._setActivatedIndex(tabIndex, true, subTab)
            // this.activatedTab = tab.name
        },

        /**
         * keep update value in tab
         * @param {string} server
         * @param {number} db
         * @param {string} key
         * @param {*} value
         * @param {string} [viewAs]
         * @param {string] [decode]
         * @param {boolean} reset
         * @param {boolean} [end] keep end status if not set
         */
        updateValue({ server, db, key, value, viewAs, decode, reset, end }) {
            const tab = find(this.tabList, { name: server, db, key })
            if (tab == null) {
                return
            }

            tab.viewAs = viewAs || tab.viewAs
            tab.decode = decode || tab.decode
            if (typeof end === 'boolean') {
                tab.end = end
            }
            if (!reset && typeof value === 'object') {
                if (value instanceof Array) {
                    tab.value = tab.value || []
                    tab.value.push(...value)
                } else {
                    tab.value = assign(value, tab.value || {})
                }
            } else {
                tab.value = value
            }
        },

        /**
         * update or insert value entries
         * @param {string} server
         * @param {number} db
         * @param {string} key
         * @param {string} type
         * @param {string[]|Object.<string, number>|Object.<number, string>} entries
         * @param {boolean} [prepend] for list only
         * @param {boolean} [reset]
         * @param {boolean} [nocheck] ignore conflict checking for hash/set/zset
         */
        upsertValueEntries({ server, db, key, type, entries, prepend, reset, nocheck }) {
            const tab = find(this.tabList, { name: server, db, key })
            if (tab == null) {
                return
            }

            switch (type.toLowerCase()) {
                case 'list': // string[] | Object.<number, string>
                    if (entries instanceof Array) {
                        // append or prepend items
                        if (reset === true) {
                            tab.value = entries
                        } else {
                            tab.value = tab.value || []
                            if (prepend === true) {
                                tab.value = [...entries, ...tab.value]
                            } else {
                                tab.value.push(...entries)
                            }
                            tab.length += size(entries)
                        }
                    } else {
                        // replace by index
                        tab.value = tab.value || []
                        for (const idx in entries) {
                            set(tab.value, idx, entries[idx])
                        }
                    }
                    break

                case 'hash': // Object.<string, string>
                    if (reset === true) {
                        tab.value = {}
                        tab.length = 0
                    } else {
                        tab.value = tab.value || {}
                    }
                    for (const k in entries) {
                        if (nocheck !== true && !tab.value.hasOwnProperty(k)) {
                            tab.length += 1
                        }
                        tab.value[k] = entries[k]
                    }
                    break

                case 'set': // string[] | Object.{string, string}
                    if (reset === true) {
                        tab.value = entries
                    } else {
                        tab.value = tab.value || []
                        if (entries instanceof Array) {
                            // add items
                            for (const elem of entries) {
                                if (nocheck !== true && indexOf(tab.value, elem) === -1) {
                                    tab.value.push(elem)
                                    tab.length += 1
                                }
                            }
                        } else {
                            // replace items
                            for (const k in entries) {
                                const idx = indexOf(tab.value, k)
                                if (idx !== -1) {
                                    tab.value[idx] = entries[k]
                                } else {
                                    tab.value.push(entries[k])
                                    tab.length += 1
                                }
                            }
                        }
                    }
                    break

                case 'zset': // {value: string, score: number}
                    if (reset === true) {
                        tab.value = Object.entries(entries).map(([value, score]) => ({ value, score }))
                    } else {
                        tab.value = tab.value || []
                        for (const val in entries) {
                            if (nocheck !== true) {
                                const ent = find(tab.value, (e) => e.value === val)
                                if (ent != null) {
                                    ent.score = entries[val]
                                } else {
                                    tab.value.push({ value: val, score: entries[val] })
                                    tab.length += 1
                                }
                            } else {
                                tab.value.push({ value: val, score: entries[val] })
                                tab.length += 1
                            }
                        }
                    }
                    break

                case 'stream': // [{id: string, value: []any}]
                    if (reset === true) {
                        tab.value = entries
                    } else {
                        tab.value = tab.value || []
                        tab.value = [...entries, ...tab.value]
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
                        // remove by index、
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
                        if (tab.value.hasOwnProperty(k)) {
                            delete tab.value[k]
                            tab.length -= 1
                        }
                    }
                    break

                case 'set': // []string
                    tab.value = tab.value || []
                    tab.length -= size(remove(tab.value, (v) => entries.indexOf(v) >= 0))
                    break

                case 'zset': // string[]
                    tab.value = tab.value || []
                    tab.length -= size(remove(tab.value, (v) => entries.indexOf(v.value) >= 0))
                    break

                case 'stream': // string[]
                    tab.value = tab.value || []
                    tab.length -= size(remove(tab.value, (v) => entries.indexOf(v.id) >= 0))
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
            let tab = find(this.tabList, { name: server, db, key })
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
         * set selected keys of current display browser tree
         * @param {string} server
         * @param {string|string[]} [keys]
         */
        setSelectedKeys(server, keys = null) {
            let tab = find(this.tabList, { name: server })
            if (tab != null) {
                if (keys == null) {
                    // select nothing
                    tab.selectedKeys = [server]
                } else if (typeof keys === 'string') {
                    tab.selectedKeys = [keys]
                } else {
                    tab.selectedKeys = keys
                }
            }
        },
    },
})

export default useTabStore
