import { defineStore } from 'pinia'

const useDialogStore = defineStore('dialog', {
    state: () => ({
        newDialogVisible: false,

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

        selectTTL: -1,
        ttlDialogVisible: false,

        preferencesDialogVisible: false,
    }),
    actions: {
        openNewDialog() {
            this.newDialogVisible = true
        },
        closeNewDialog() {
            this.newDialogVisible = false
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
         * @param {string} prefix
         * @param {number} server
         * @param {string} db
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
