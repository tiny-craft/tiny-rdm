import { defineStore } from 'pinia'
import { get, isEmpty, isObject, uniq } from 'lodash'
import {
    CreateGroup,
    DeleteConnection,
    DeleteGroup,
    ExportConnections,
    GetConnection,
    ImportConnections,
    ListConnection,
    ParseConnectURL,
    RenameGroup,
    SaveConnection,
    SaveLastDB,
    SaveRefreshInterval,
    SaveSortedConnection,
} from 'wailsjs/go/services/connectionService.js'
import { ConnectionType } from '@/consts/connection_type.js'
import { KeyViewType } from '@/consts/key_view_type.js'
import useBrowserStore from 'stores/browser.js'
import { i18nGlobal } from '@/utils/i18n.js'
import { ClipboardGetText } from 'wailsjs/runtime/runtime.js'

const useConnectionStore = defineStore('connections', {
    /**
     * @typedef {Object} ConnectionItem
     * @property {string} key
     * @property {string} label display label
     * @property {string} name database name
     * @property {number} type
     * @property {boolean} cluster is cluster node
     * @property {ConnectionItem[]} children
     */

    /**
     * @typedef {Object} ConnectionProfile
     * @property {string} defaultFilter
     * @property {string} keySeparator
     * @property {string} markColor
     * @property {number} refreshInterval
     */

    /**
     * @typedef {Object} ConnectionState
     * @property {string[]} groups
     * @property {ConnectionItem[]} connections
     * @property {Object.<string, ConnectionProfile>} serverProfile
     */

    /**
     *
     * @returns {ConnectionState}
     */
    state: () => ({
        groups: [], // all group name set
        connections: [], // all connections
        serverProfile: {}, // all server profile in flat list
    }),
    getters: {},
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
            const { data = [{ groupName: '', connections: [], refreshInterval: 5 }] } = await ListConnection()
            for (const conn of data) {
                if (conn.type !== 'group') {
                    // top level
                    conns.push({
                        key: '/' + conn.name,
                        label: conn.name,
                        name: conn.name,
                        type: ConnectionType.Server,
                        cluster: get(conn, 'cluster.enable', false),
                        // isLeaf: false,
                    })
                    profiles[conn.name] = {
                        defaultFilter: conn.defaultFilter,
                        keySeparator: conn.keySeparator,
                        markColor: conn.markColor,
                        refreshInterval: conn.refreshInterval,
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
                            cluster: get(item, 'cluster.enable', false),
                            // isLeaf: false,
                        })
                        profiles[item.name] = {
                            defaultFilter: item.defaultFilter,
                            keySeparator: item.keySeparator,
                            markColor: item.markColor,
                            refreshInterval: item.refreshInterval,
                        }
                    }
                    conns.push({
                        key: conn.name + '/',
                        label: conn.name,
                        type: ConnectionType.Group,
                        children,
                    })
                }
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
                network: 'tcp',
                sock: '/tmp/redis.sock',
                addr: '127.0.0.1',
                port: 6379,
                username: '',
                password: '',
                defaultFilter: '*',
                keySeparator: ':',
                connTimeout: 60,
                execTimeout: 60,
                dbFilterType: 'none',
                dbFilterList: [],
                keyView: KeyViewType.Tree,
                loadSize: 10000,
                markColor: '',
                alias: {},
                ssl: {
                    enable: false,
                    allowInsecure: true,
                    sni: '',
                    certFile: '',
                    keyFile: '',
                    caFile: '',
                },
                ssh: {
                    enable: false,
                    addr: '',
                    port: 22,
                    loginType: 'pwd',
                    username: '',
                    password: '',
                    pkFile: '',
                    passphrase: '',
                },
                sentinel: {
                    enable: false,
                    master: 'mymaster',
                    username: '',
                    password: '',
                },
                cluster: {
                    enable: false,
                },
                proxy: {
                    type: 0,
                    schema: 'http',
                    addr: '',
                    port: 0,
                    auth: false,
                    username: '',
                    password: '',
                },
            }
        },

        mergeConnectionProfile(dest, src) {
            const mergeObj = (destObj, srcObj) => {
                for (const k in srcObj) {
                    const t = typeof srcObj[k]
                    if (t === 'string') {
                        destObj[k] = srcObj[k] || destObj[k] || ''
                    } else if (t === 'number') {
                        destObj[k] = srcObj[k] || destObj[k] || 0
                    } else if (t === 'object') {
                        mergeObj(destObj[k], srcObj[k] || {})
                    } else {
                        destObj[k] = srcObj[k]
                    }
                }
                return destObj
            }
            return mergeObj(dest, src)
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
         * remove connection
         * @param name
         * @returns {Promise<{success: boolean, [msg]: string}>}
         */
        async deleteConnection(name) {
            // close connection first
            const browser = useBrowserStore()
            await browser.closeConnection(name)
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
         * save last selected database
         * @param {string} name
         * @param {number} db
         * @return {Promise<{success: boolean, [msg]: string}>}
         */
        async saveLastDB(name, db) {
            const { success, msg } = await SaveLastDB(name, db)
            if (!success) {
                return { success: false, msg }
            }
            return { success: true }
        },

        /**
         * get default key filter pattern by server name
         * @param name
         * @return {string}
         */
        getDefaultKeyFilter(name) {
            const { defaultFilter = '*' } = this.serverProfile[name] || {}
            return defaultFilter
        },

        /**
         * get default key separator by server name
         * @param name
         * @return {string}
         */
        getDefaultSeparator(name) {
            const { keySeparator = ':' } = this.serverProfile[name] || {}
            return keySeparator
        },

        /**
         * get default status refresh interval by server name
         * @param {string} name
         * @return {number}
         */
        getRefreshInterval(name) {
            const { refreshInterval = 5 } = this.serverProfile[name] || {}
            return refreshInterval
        },

        /**
         * set and save default refresh interval
         * @param {string} name
         * @param {number} interval
         * @return {Promise<{success: boolean}|{msg: undefined, success: boolean}>}
         */
        async saveRefreshInterval(name, interval) {
            const profile = this.serverProfile[name] || {}
            profile.refreshInterval = interval
            const { success, msg } = await SaveRefreshInterval(name, interval)
            if (!success) {
                return { success: false, msg }
            }
            return { success: true }
        },

        /**
         * export connections to zip
         * @return {Promise<void>}
         */
        async exportConnections() {
            const {
                success,
                msg,
                data: { path = '' },
            } = await ExportConnections()
            if (!success) {
                if (!isEmpty(msg)) {
                    $message.error(msg)
                    return
                }
            }

            $message.success(i18nGlobal.t('dialogue.handle_succ'))
        },

        /**
         * import connections from zip
         * @return {Promise<void>}
         */
        async importConnections() {
            const { success, msg } = await ImportConnections()
            if (!success) {
                if (!isEmpty(msg)) {
                    $message.error(msg)
                    return
                }
            }

            $message.success(i18nGlobal.t('dialogue.handle_succ'))
        },

        /**
         * parse redis url from text in clipboard
         * @return {Promise<{}>}
         */
        async parseUrlFromClipboard() {
            const urlString = await ClipboardGetText()
            if (isEmpty(urlString)) {
                throw new Error('no text in clipboard')
            }

            const { success, msg, data } = await ParseConnectURL(urlString)
            if (!success || !isObject(data)) {
                throw new Error(msg || 'unknown')
            }

            data.url = urlString
            return data
        },
    },
})

export default useConnectionStore
