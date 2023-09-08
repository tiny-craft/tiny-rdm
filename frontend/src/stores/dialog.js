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

        selectTTL: -1,
        ttlDialogVisible: false,

        preferencesDialogVisible: false,
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
            this.connParam = profile || connStore.newDefaultConnection(name)
            this.connType = ConnDialogType.EDIT
            this.connDialogVisible = true
        },

        async openDuplicateDialog(name) {
            const connStore = useConnectionStore()
            const profile = await connStore.getConnectionProfile(name)
            if (profile != null) {
                profile.name += '2'
                this.connParam = profile
            } else {
                this.connParam = connStore.newDefaultConnection(name)
            }
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

        openDeleteKeyDialog(server, db, key) {
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
         * @param {string} type
         */
        openAddFieldsDialog(server, db, key, type) {
            this.addFieldParam.server = server
            this.addFieldParam.db = db
            this.addFieldParam.key = key
            this.addFieldParam.type = type
            this.addFieldsDialogVisible = true
        },
        closeAddFieldsDialog() {
            this.addFieldsDialogVisible = false
        },

        openTTLDialog(ttl) {
            this.selectTTL = ttl
            this.ttlDialogVisible = true
        },
        closeTTLDialog() {
            this.selectTTL = -1
            this.ttlDialogVisible = false
        },

        openPreferencesDialog() {
            this.preferencesDialogVisible = true
        },
        closePreferencesDialog() {
            this.preferencesDialogVisible = false
        },
    },
})

export default useDialogStore
