import { defineStore } from 'pinia'
import useConnectionStore from './connections.js'

/**
 * connection dialog type
 * @enum {number}
 */
export const ConnDialogType = {
    NEW: 0,
    EDIT: 1,
}

const useDialogStore = defineStore('dialog', {
    state: () => ({
        connDialogVisible: false,
        /** @type {ConnDialogType} **/
        connType: ConnDialogType.NEW,
        connParam: null,

        groupDialogVisible: false,
        editGroup: '',

        /**
         * @property {string} prefix
         * @property {string} server
         * @property {int} db
         */
        newKeyParam: {
            prefix: '',
            server: '',
            db: 0,
        },
        newKeyDialogVisible: false,

        keyFilterParam: {
            server: '',
            db: 0,
            type: '',
            pattern: '*',
        },
        keyFilterDialogVisible: false,

        addFieldParam: {
            server: '',
            db: 0,
            key: '',
            keyCode: null,
            type: null,
        },
        addFieldsDialogVisible: false,

        renameKeyParam: {
            server: '',
            db: 0,
            key: '',
        },
        renameDialogVisible: false,

        deleteKeyParam: {
            server: '',
            db: 0,
            key: '',
        },
        deleteKeyDialogVisible: false,

        exportKeyParam: {
            server: '',
            db: 0,
            keys: [],
        },
        exportKeyDialogVisible: false,

        importKeyParam: {
            server: '',
            db: 0,
        },
        importKeyDialogVisible: false,

        flushDBParam: {
            server: '',
            db: 0,
        },
        flushDBDialogVisible: false,

        ttlDialogVisible: false,
        ttlParam: {
            server: '',
            db: 0,
            key: '',
            keys: [],
            ttl: 0,
        },

        decodeDialogVisible: false,
        decodeParam: {
            name: '',
            auto: true,
            decodePath: '',
            decodeArgs: [],
            encodePath: '',
            encodeArgs: [],
        },

        preferencesDialogVisible: false,
        aboutDialogVisible: false,
    }),
    actions: {
        openNewDialog() {
            this.connParam = null
            this.connType = ConnDialogType.NEW
            this.connDialogVisible = true
        },
        closeConnDialog() {
            this.connDialogVisible = false
        },

        async openEditDialog(name) {
            const connStore = useConnectionStore()
            const profile = await connStore.getConnectionProfile(name)
            this.connParam = connStore.mergeConnectionProfile(connStore.newDefaultConnection(name), profile)
            this.connType = ConnDialogType.EDIT
            this.connDialogVisible = true
        },

        async openDuplicateDialog(name) {
            const connStore = useConnectionStore()
            this.connParam = {}
            let profile
            let suffix = 1
            do {
                let profileName = name
                if (suffix > 1) {
                    profileName += suffix
                }
                profile = await connStore.getConnectionProfile(profileName)
                if (profile != null) {
                    suffix += 1
                    if (profileName === name) {
                        this.connParam = profile
                    }
                } else {
                    this.connParam = connStore.mergeConnectionProfile(
                        connStore.newDefaultConnection(profileName),
                        this.connParam,
                    )
                    this.connParam.name = profileName
                    break
                }
            } while (true)
            this.connType = ConnDialogType.NEW
            this.connDialogVisible = true
        },

        openNewGroupDialog() {
            this.editGroup = ''
            this.groupDialogVisible = true
        },
        closeNewGroupDialog() {
            this.groupDialogVisible = false
        },

        /**
         *
         * @param {string} server
         * @param {number} db
         * @param {string} [pattern]
         * @param {string} [type]
         */
        openKeyFilterDialog(server, db, pattern, type) {
            this.keyFilterParam.server = server
            this.keyFilterParam.db = db
            this.keyFilterParam.type = type || ''
            this.keyFilterParam.pattern = pattern || '*'
            this.keyFilterDialogVisible = true
        },
        closeKeyFilterDialog() {
            this.keyFilterDialogVisible = false
        },

        /**
         *
         * @param {string} name
         */
        openRenameGroupDialog(name) {
            this.editGroup = name
            this.groupDialogVisible = true
        },
        closeRenameGroupDialog() {
            this.groupDialogVisible = false
        },

        /**
         *
         * @param {string} server
         * @param {number} db
         * @param {string} key
         */
        openRenameKeyDialog(server, db, key) {
            this.renameKeyParam.server = server
            this.renameKeyParam.db = db
            this.renameKeyParam.key = key
            this.renameDialogVisible = true
        },
        closeRenameKeyDialog() {
            this.renameDialogVisible = false
        },

        /**
         *
         * @param {string} server
         * @param {number} db
         * @param {string|string[]} [key]
         */
        openDeleteKeyDialog(server, db, key = '*') {
            this.deleteKeyParam.server = server
            this.deleteKeyParam.db = db
            this.deleteKeyParam.key = key
            this.deleteKeyDialogVisible = true
        },
        closeDeleteKeyDialog() {
            this.deleteKeyDialogVisible = false
        },

        /**
         *
         * @param {string} server
         * @param {number} db
         * @param {string|string[]} keys
         */
        openExportKeyDialog(server, db, keys) {
            this.exportKeyParam.server = server
            this.exportKeyParam.db = db
            this.exportKeyParam.keys = keys
            this.exportKeyDialogVisible = true
        },
        closeExportKeyDialog() {
            this.exportKeyDialogVisible = false
        },

        /**
         *
         * @param {string} server
         * @param {number} db
         */
        openImportKeyDialog(server, db) {
            this.importKeyParam.server = server
            this.importKeyParam.db = db
            this.importKeyDialogVisible = true
        },
        closeImportKeyDialog() {
            this.importKeyDialogVisible = false
        },

        openFlushDBDialog(server, db) {
            this.flushDBParam.server = server
            this.flushDBParam.db = db
            this.flushDBDialogVisible = true
        },
        closeFlushDBDialog() {
            this.flushDBDialogVisible = false
        },

        /**
         *
         * @param {string} prefix
         * @param {string} server
         * @param {number} db
         */
        openNewKeyDialog(prefix, server, db) {
            this.newKeyParam.prefix = prefix
            this.newKeyParam.server = server
            this.newKeyParam.db = db
            this.newKeyDialogVisible = true
        },
        closeNewKeyDialog() {
            this.newKeyDialogVisible = false
        },

        /**
         *
         * @param {string} server
         * @param {number} db
         * @param {string} key
         * @param {number[]|null} keyCode
         * @param {string} type
         */
        openAddFieldsDialog(server, db, key, keyCode, type) {
            this.addFieldParam.server = server
            this.addFieldParam.db = db
            this.addFieldParam.key = key
            this.addFieldParam.keyCode = keyCode
            this.addFieldParam.type = type
            this.addFieldsDialogVisible = true
        },
        closeAddFieldsDialog() {
            this.addFieldsDialogVisible = false
        },

        /**
         *
         * @param {string} server
         * @param {number} db
         * @param {string|number[]} [key]
         * @param {string[]|number[][]} [keys]
         * @param {number} [ttl]
         */
        openTTLDialog({ server, db, key, keys, ttl = -1 }) {
            this.ttlDialogVisible = true
            this.ttlParam.server = server
            this.ttlParam.db = db
            this.ttlParam.key = key
            this.ttlParam.keys = keys
            this.ttlParam.ttl = ttl
        },
        closeTTLDialog() {
            this.ttlDialogVisible = false
        },

        /**
         *
         * @param {string} name
         * @param {boolean} auto
         * @param {string} decodePath
         * @param {string[]} decodeArgs
         * @param {string} encodePath
         * @param {string[]} encodeArgs
         */
        openDecoderDialog({
            name = '',
            auto = true,
            decodePath = '',
            decodeArgs = [],
            encodePath = '',
            encodeArgs = [],
        } = {}) {
            this.decodeDialogVisible = true
            this.decodeParam.name = name
            this.decodeParam.auto = auto !== false
            this.decodeParam.decodePath = decodePath
            this.decodeParam.decodeArgs = decodeArgs || []
            this.decodeParam.encodePath = encodePath
            this.decodeParam.encodeArgs = encodeArgs || []
        },

        closeDecoderDialog() {
            this.decodeDialogVisible = false
        },

        openPreferencesDialog() {
            this.preferencesDialogVisible = true
        },
        closePreferencesDialog() {
            this.preferencesDialogVisible = false
        },

        openAboutDialog() {
            this.aboutDialogVisible = true
        },
        closeAboutDialog() {
            this.aboutDialogVisible = false
        },
    },
})

export default useDialogStore
