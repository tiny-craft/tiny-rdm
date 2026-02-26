import { defineStore } from 'pinia'
import { lang } from '@/langs/index.js'
import { cloneDeep, findIndex, get, isEmpty, join, map, pick, set, some, split } from 'lodash'
import {
    CheckForUpdate,
    GetBuildInDecoder,
    GetFontList,
    GetPreferences,
    RestorePreferences,
    SetPreferences
} from 'wailsjs/go/services/preferencesService.js'
import { BrowserOpenURL } from 'wailsjs/runtime/runtime.js'
import { i18nGlobal } from '@/utils/i18n.js'
import { enUS, NButton, NSpace, useOsTheme, zhCN } from 'naive-ui'
import { h, nextTick } from 'vue'
import { compareVersion } from '@/utils/version.js'
import { typesIconStyle } from '@/consts/support_redis_type.js'
import { TextAlignType } from '@/consts/text_align_type.js'

const osTheme = useOsTheme()
const usePreferencesStore = defineStore('preferences', {
    /**
     * @typedef {Object} FontItem
     * @property {string} name
     * @property {string} path
     */
    /**
     * @typedef {Object} Preferences
     * @property {Object} general
     * @property {Object} editor
     * @property {FontItem[]} fontList
     */
    /**
     *
     * @returns {Preferences}
     */
    state: () => ({
        behavior: {
            welcomed: false,
            asideWidth: 300,
            windowWidth: 0,
            windowHeight: 0,
            windowMaximised: false,
        },
        general: {
            theme: 'auto',
            language: 'auto',
            font: '',
            fontFamily: [],
            fontSize: 14,
            scanSize: 3000,
            keyIconStyle: 0,
            useSysProxy: false,
            useSysProxyHttp: false,
            checkUpdate: true,
            skipVersion: '',
            allowTrack: true,
        },
        editor: {
            font: '',
            fontFamily: [],
            fontSize: 14,
            showLineNum: true,
            showFolding: true,
            dropText: true,
            links: true,
            entryTextAlign: TextAlignType.Center,
        },
        cli: {
            fontFamily: [],
            fontSize: 14,
            cursorStyle: 'block',
        },
        buildInDecoder: [],
        decoder: [],
        lastPref: {},
        fontList: [],
    }),
    getters: {
        getSeparator() {
            return ':'
        },

        themeOption() {
            return [
                {
                    value: 'light',
                    label: 'preferences.general.theme_light',
                },
                {
                    value: 'dark',
                    label: 'preferences.general.theme_dark',
                },
                {
                    value: 'auto',
                    label: 'preferences.general.theme_auto',
                },
            ]
        },

        /**
         * all available language
         * @returns {{label: string, value: string}[]}
         */
        langOption() {
            const options = Object.entries(lang).map(([key, value]) => ({
                value: key,
                label: value['name'],
            }))
            options.splice(0, 0, {
                value: 'auto',
                label: 'preferences.general.system_lang',
            })
            return options
        },

        /**
         * all system font list
         * @returns {{path: string, label: string, value: string}[]}
         */
        fontOption() {
            return map(this.fontList, (font) => ({
                value: font.name,
                label: font.name,
                path: font.path,
            }))
        },

        /**
         * current font selection
         * @returns {{fontSize: string, fontFamily?: string}}
         */
        generalFont() {
            const fontStyle = {
                fontSize: this.general.fontSize + 'px',
            }
            if (!isEmpty(this.general.fontFamily)) {
                fontStyle['fontFamily'] = join(
                    map(this.general.fontFamily, (f) => `"${f}"`),
                    ',',
                )
            }
            // compatible with old preferences
            // if (isEmpty(fontStyle['fontFamily'])) {
            //     if (!isEmpty(this.general.font) && this.general.font !== 'none') {
            //         const font = find(this.fontList, { name: this.general.font })
            //         if (font != null) {
            //             fontStyle['fontFamily'] = `${font.name}`
            //         }
            //     }
            // }
            return fontStyle
        },

        /**
         * current editor font
         * @return {{fontSize: string, fontFamily?: string}}
         */
        editorFont() {
            const fontStyle = {
                fontSize: (this.editor.fontSize || 14) + 'px',
            }
            if (!isEmpty(this.editor.fontFamily)) {
                fontStyle['fontFamily'] = join(
                    map(this.editor.fontFamily, (f) => `"${f}"`),
                    ',',
                )
            }
            // compatible with old preferences
            // if (isEmpty(fontStyle['fontFamily'])) {
            //     if (!isEmpty(this.editor.font) && this.editor.font !== 'none') {
            //         const font = find(this.fontList, { name: this.editor.font })
            //         if (font != null) {
            //             fontStyle['fontFamily'] = `${font.name}`
            //         }
            //     }
            // }
            if (isEmpty(fontStyle['fontFamily'])) {
                fontStyle['fontFamily'] = ['monaco']
            }
            return fontStyle
        },

        /**
         * current cli font
         * @return {{fontSize: string, fontFamily?: string}}
         */
        cliFont() {
            const fontStyle = {
                fontSize: this.cli.fontSize || 14,
            }
            if (!isEmpty(this.cli.fontFamily)) {
                fontStyle['fontFamily'] = join(
                    map(this.cli.fontFamily, (f) => `"${f}"`),
                    ',',
                )
            }
            if (isEmpty(fontStyle['fontFamily'])) {
                fontStyle['fontFamily'] = ['Courier New']
            }
            return fontStyle
        },

        cliCursorStyleOption() {
            return [
                {
                    value: 'block',
                    label: 'preferences.cli.cursor_style_block',
                },
                {
                    value: 'underline',
                    label: 'preferences.cli.cursor_style_underline',
                },
                {
                    value: 'bar',
                    label: 'preferences.cli.cursor_style_bar',
                },
            ]
        },

        /**
         * get current language setting
         * @return {string}
         */
        currentLanguage() {
            let lang = get(this.general, 'language', 'auto')
            if (lang === 'auto') {
                const systemLang = navigator.language || navigator.userLanguage
                lang = split(systemLang, '-')[0]
            }
            return lang || 'en'
        },

        isDark() {
            const th = get(this.general, 'theme', 'auto')
            if (th !== 'auto') {
                return th === 'dark'
            } else {
                return osTheme.value === 'dark'
            }
        },

        themeLocale() {
            const lang = this.currentLanguage
            switch (lang) {
                case 'zh':
                    return zhCN
                default:
                    return enUS
            }
        },

        autoCheckUpdate() {
            return get(this.general, 'checkUpdate', false)
        },

        showLineNum() {
            return get(this.editor, 'showLineNum', true)
        },

        showFolding() {
            return get(this.editor, 'showFolding', true)
        },

        dropText() {
            return get(this.editor, 'dropText', true)
        },

        editorLinks() {
            return get(this.editor, 'links', true)
        },

        keyIconType() {
            return get(this.general, 'keyIconStyle', typesIconStyle.SHORT)
        },

        entryTextAlign() {
            return get(this.editor, 'entryTextAlign', TextAlignType.Center)
        },
    },
    actions: {
        _applyPreferences(data) {
            for (const key in data) {
                set(this, key, data[key])
            }
        },

        /**
         * load preferences from local
         * @returns {Promise<boolean>}
         */
        async loadPreferences() {
            const { success, data } = await GetPreferences()
            if (success) {
                this.lastPref = cloneDeep(data)
                this._applyPreferences(data)
                // default value
                const showLineNum = get(data, 'editor.showLineNum')
                if (showLineNum === undefined) {
                    set(data, 'editor.showLineNum', true)
                }
                const showFolding = get(data, 'editor.showFolding')
                if (showFolding === undefined) {
                    set(data, 'editor.showFolding', true)
                }
                const dropText = get(data, 'editor.dropText')
                if (dropText === undefined) {
                    set(data, 'editor.dropText', true)
                }
                const links = get(data, 'editor.links')
                if (links === undefined) {
                    set(data, 'editor.links', true)
                }
                i18nGlobal.locale.value = this.currentLanguage
            }
            return success
        },

        /**
         * load system font list
         * @returns {Promise<string[]>}
         */
        async loadFontList() {
            const { success, data } = await GetFontList()
            if (success) {
                const { fonts = [] } = data
                this.fontList = fonts
            } else {
                this.fontList = []
            }
            return this.fontList
        },

        /**
         * get all available build-in decoder
         * @return {Promise<void>}
         */
        async loadBuildInDecoder() {
            const { success, data } = await GetBuildInDecoder()
            if (success) {
                const { decoder = [] } = data
                this.buildInDecoder = decoder
            } else {
                this.buildInDecoder = []
            }
        },

        /**
         * save preferences to local
         * @returns {Promise<boolean>}
         */
        async savePreferences() {
            const pf = pick(this, ['behavior', 'general', 'editor', 'cli', 'decoder'])
            const { success } = await SetPreferences(pf)
            return success === true
        },

        /**
         * reset to last-loaded preferences
         * @returns {Promise<void>}
         */
        async resetToLastPreferences() {
            if (!isEmpty(this.lastPref)) {
                this._applyPreferences(this.lastPref)
            }
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

        /**
         * add a new custom decoder
         * @param {string} name
         * @param {boolean} enable
         * @param {boolean} auto
         * @param {string} encodePath
         * @param {string[]} encodeArgs
         * @param {string} decodePath
         * @param {string[]} decodeArgs
         */
        addCustomDecoder({ name, enable = true, auto = true, encodePath, encodeArgs, decodePath, decodeArgs }) {
            if (some(this.decoder, { name })) {
                return false
            }
            this.decoder = this.decoder || []
            this.decoder.push({ name, enable, auto, encodePath, encodeArgs, decodePath, decodeArgs })
            return true
        },

        /**
         * update an existing custom decoder
         * @param {string} newName
         * @param {boolean} enable
         * @param {boolean} auto
         * @param {string} name
         * @param {string} encodePath
         * @param {string[]} encodeArgs
         * @param {string} decodePath
         * @param {string[]} decodeArgs
         */
        updateCustomDecoder({
            newName,
            enable = true,
            auto = true,
            name,
            encodePath,
            encodeArgs,
            decodePath,
            decodeArgs,
        }) {
            const idx = findIndex(this.decoder, { name })
            if (idx === -1) {
                return false
            }
            // conflicted
            if (newName !== name && some(this.decoder, { name: newName })) {
                return false
            }

            let selDecoder = this.decoder[idx]
            selDecoder.name = newName || name
            selDecoder.enable = enable
            selDecoder.auto = auto
            selDecoder.encodePath = encodePath
            selDecoder.encodeArgs = encodeArgs
            selDecoder.decodePath = decodePath
            selDecoder.decodeArgs = decodeArgs
            this.decoder[idx] = selDecoder
            return true
        },

        /**
         * remove an existing custom decoder
         * @param {string} name
         * @return {boolean}
         */
        removeCustomDecoder(name) {
            const idx = findIndex(this.decoder, { name })
            if (idx === -1) {
                return false
            }
            this.decoder.splice(idx, 1)
            return true
        },

        setAsWelcomed(acceptTrack) {
            this.behavior.welcomed = true
            this.general.allowTrack = acceptTrack
            this.savePreferences()
        },

        async checkForUpdate(manual = false) {
            let msgRef = null
            if (manual) {
                msgRef = $message.loading(i18nGlobal.t('interface.retrieving_version'), { duration: 0 })
            }
            try {
                const { success, data = {} } = await CheckForUpdate()
                if (success) {
                    const {
                        version = 'v1.0.0',
                        latest,
                        download_page: pageUrl = {},
                        description = {},
                        sponsor = [],
                        banner = [],
                    } = data
                    const downUrl = pageUrl[this.currentLanguage] || pageUrl['en']
                    const descStr = description[this.currentLanguage] || description['en']
                    // save sponsor ad
                    if (!isEmpty(sponsor)) {
                        localStorage.setItem('sponsor_ad', JSON.stringify(sponsor))
                    }
                    if (!isEmpty(banner)) {
                        localStorage.setItem('banner', JSON.stringify(banner))
                    }
                    if (
                        (manual || compareVersion(latest, this.general.skipVersion) !== 0) &&
                        compareVersion(latest, version) > 0 &&
                        !isEmpty(downUrl)
                    ) {
                        const notiRef = $notification.show({
                            title: `${i18nGlobal.t('dialogue.upgrade.title')} - ${latest}`,
                            content: descStr || i18nGlobal.t('dialogue.upgrade.new_version_tip', { ver: latest }),
                            action: () =>
                                h('div', { class: 'flex-box-h flex-item-expand' }, [
                                    h(NSpace, { wrapItem: false }, () => [
                                        h(
                                            NButton,
                                            {
                                                size: 'small',
                                                secondary: true,
                                                onClick: () => {
                                                    // skip this update
                                                    this.general.skipVersion = latest
                                                    this.savePreferences()
                                                    notiRef.destroy()
                                                },
                                            },
                                            () => i18nGlobal.t('dialogue.upgrade.skip'),
                                        ),
                                        h(
                                            NButton,
                                            {
                                                size: 'small',
                                                secondary: true,
                                                onClick: notiRef.destroy,
                                            },
                                            () => i18nGlobal.t('dialogue.upgrade.later'),
                                        ),
                                        h(
                                            NButton,
                                            {
                                                type: 'primary',
                                                size: 'small',
                                                secondary: true,
                                                onClick: () => BrowserOpenURL(downUrl),
                                            },
                                            () => i18nGlobal.t('dialogue.upgrade.download_now'),
                                        ),
                                    ]),
                                ]),
                            onPositiveClick: () => BrowserOpenURL(downUrl),
                        })
                        return
                    }
                }

                if (manual) {
                    $message.info(i18nGlobal.t('dialogue.upgrade.no_update'))
                }
            } finally {
                nextTick().then(() => {
                    if (msgRef != null) {
                        msgRef.destroy()
                        msgRef = null
                    }
                })
            }
        },
    },
})

export default usePreferencesStore
