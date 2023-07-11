import { defineStore } from 'pinia'
import { lang } from '../langs/index.js'
import { camelCase, isObject, map, set, snakeCase, split } from 'lodash'
import { GetPreferences, RestorePreferences, SetPreferences } from '../../wailsjs/go/services/preferencesService.js'

const usePreferencesStore = defineStore('preferences', {
    state: () => ({
        general: {
            language: 'en',
            font: '',
            fontSize: 14,
            useSysProxy: false,
            useSysProxyHttp: false,
            checkUpdate: false,
        },
        editor: {
            font: '',
            fontSize: 14,
        },
    }),
    getters: {
        getSeparator() {
            return ':'
        },

        /**
         * all available language
         * @returns {{label: string, value: string}[]}
         */
        langOption() {
            return Object.entries(lang).map(([key, value]) => ({
                value: key,
                label: `${value['lang_name']}`,
            }))
        },
    },
    actions: {
        _applyPreferences(data) {
            for (const key in data) {
                const keys = map(split(key, '.'), camelCase)
                set(this, keys, data[key])
            }
        },

        /**
         * load preferences from local
         * @returns {Promise<void>}
         */
        async loadPreferences() {
            const { success, data } = await GetPreferences()
            if (success) {
                this._applyPreferences(data)
            }
        },

        /**
         * save preferences to local
         * @returns {Promise<boolean>}
         */
        async savePreferences() {
            const obj2Map = (prefix, obj) => {
                const result = {}
                for (const key in obj) {
                    if (isObject(obj[key])) {
                        // TODO: perform sub object
                    } else {
                        result[`${prefix}.${snakeCase(key)}`] = obj[key]
                    }
                }
                return result
            }
            const pf = Object.assign({}, obj2Map('general', this.general), obj2Map('editor', this.editor))
            const { success, msg } = await SetPreferences(pf)
            return success === true
        },

        /**
         * restore preferences to default
         * @returns {Promise<boolean>}
         */
        async restorePreferences() {
            const { success, data } = await RestorePreferences()
            if (success === true) {
                const { pref } = data
                this._applyPreferences(pref)
                return true
            }
            return false
        },
    },
})

export default usePreferencesStore
