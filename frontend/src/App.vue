<script setup>
import ConnectionDialog from './components/dialogs/ConnectionDialog.vue'
import NewKeyDialog from './components/dialogs/NewKeyDialog.vue'
import PreferencesDialog from './components/dialogs/PreferencesDialog.vue'
import RenameKeyDialog from './components/dialogs/RenameKeyDialog.vue'
import SetTtlDialog from './components/dialogs/SetTtlDialog.vue'
import AddFieldsDialog from './components/dialogs/AddFieldsDialog.vue'
import AppContent from './AppContent.vue'
import GroupDialog from './components/dialogs/GroupDialog.vue'
import DeleteKeyDialog from './components/dialogs/DeleteKeyDialog.vue'
import { defineAsyncComponent, h, onMounted, onUnmounted, ref, watch } from 'vue'
import usePreferencesStore from './stores/preferences.js'
import useConnectionStore from './stores/connections.js'
import { useI18n } from 'vue-i18n'
import { darkTheme, NButton, NSpace } from 'naive-ui'
import KeyFilterDialog from './components/dialogs/KeyFilterDialog.vue'
import { Environment, WindowSetDarkTheme, WindowSetLightTheme } from 'wailsjs/runtime/runtime.js'
import { darkThemeOverrides, themeOverrides } from '@/utils/theme.js'
import AboutDialog from '@/components/dialogs/AboutDialog.vue'
import FlushDbDialog from '@/components/dialogs/FlushDbDialog.vue'
import ExportKeyDialog from '@/components/dialogs/ExportKeyDialog.vue'
import ImportKeyDialog from '@/components/dialogs/ImportKeyDialog.vue'
import { Info } from 'wailsjs/go/services/systemService.js'
import DecoderDialog from '@/components/dialogs/DecoderDialog.vue'
import { loadModule, trackEvent } from '@/utils/analytics.js'

const prefStore = usePreferencesStore()
const connectionStore = useConnectionStore()
const i18n = useI18n()
const initializing = ref(true)

// Detect if running in web mode (VITE_WEB=true at build time)
const isWebMode = import.meta.env.VITE_WEB === 'true'

// Web-only: lazy load LoginPage to avoid importing websocket.js in desktop mode
const LoginPage = isWebMode ? defineAsyncComponent(() => import('@/components/LoginPage.vue')) : null

// Viewport management for mobile
const setViewport = (mode) => {
    const meta = document.querySelector('meta[name="viewport"]')
    if (!meta) return
    if (mode === 'mobile') {
        meta.setAttribute('content', 'width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no')
    } else {
        const ratio = (window.innerWidth || screen.width) / (window.innerHeight || screen.height)
        const sh = window.innerHeight || screen.height
        let vw
        if (ratio < 1) {
            vw = 680
        } else if (sh < 500) {
            vw = 1280
        } else {
            vw = 1024
        }
        meta.setAttribute('content', `width=${vw}, user-scalable=yes`)
    }
}

// Auth state (web mode only)
const authChecking = ref(isWebMode) // desktop: false (skip), web: true (checking)
const authenticated = ref(false)
const authEnabled = ref(false)

const checkAuth = async () => {
    try {
        const resp = await fetch('/api/auth/status', { credentials: 'same-origin' })
        const result = await resp.json()
        if (result.success) {
            authEnabled.value = result.data.enabled
            authenticated.value = result.data.authenticated
        }
    } catch {
        authenticated.value = false
    } finally {
        authChecking.value = false
    }
}

const onLogin = async () => {
    authenticated.value = true
    setViewport('desktop')
    // Reconnect WebSocket with auth cookie (dynamic import to avoid desktop issues)
    try {
        const runtime = await import('wailsjs/runtime/runtime.js')
        if (runtime.ReconnectWebSocket) runtime.ReconnectWebSocket()
    } catch {}
    // Capture login page choices before loadPreferences overwrites them
    const loginTheme = localStorage.getItem('rdm_login_theme')
    const loginLang = localStorage.getItem('rdm_login_lang')
    await initApp()
    // Sync login page choices to preferences
    let needSave = false
    if (loginTheme && ['auto', 'light', 'dark'].includes(loginTheme)) {
        prefStore.general.theme = loginTheme
        needSave = true
    }
    if (loginLang) {
        const validLangs = ['auto', 'zh', 'tw', 'en', 'ja', 'ko', 'es', 'fr', 'ru', 'pt', 'tr']
        if (validLangs.includes(loginLang)) {
            prefStore.general.language = loginLang
            i18n.locale.value = prefStore.currentLanguage
            needSave = true
        }
    }
    if (needSave) prefStore.savePreferences()
}

const initApp = async () => {
    try {
        initializing.value = true
        if (isWebMode) {
            const prefResult = await prefStore.loadPreferences()
            // If loadPreferences failed (e.g. 401 from expired session),
            // rdm:unauthorized event already fired → silently abort init
            if (prefResult === false || !authenticated.value) return
            i18n.locale.value = prefStore.currentLanguage
        }
        await prefStore.loadFontList()
        await prefStore.loadBuildInDecoder()
        await connectionStore.initConnections()
        if (prefStore.autoCheckUpdate) {
            prefStore.checkForUpdate()
        }
        const env = await Environment()
        loadModule(env.buildType !== 'dev' && prefStore.general.allowTrack !== false).then(() => {
            Info().then(({ data }) => {
                trackEvent('startup', data, true)
            })
        })

        // show greetings and user behavior tracking statements
        if (!!!prefStore.behavior.welcomed) {
            const n = $notification.show({
                title: () => i18n.t('dialogue.welcome.title'),
                content: () => i18n.t('dialogue.welcome.content'),
                // duration: 5000,
                keepAliveOnHover: true,
                closable: false,
                meta: ' ',
                action: () =>
                    h(
                        NSpace,
                        {},
                        {
                            default: () => [
                                h(
                                    NButton,
                                    {
                                        secondary: true,
                                        type: 'tertiary',
                                        onClick: () => {
                                            prefStore.setAsWelcomed(false)
                                            n.destroy()
                                        },
                                    },
                                    {
                                        default: () => i18n.t('dialogue.welcome.reject'),
                                    },
                                ),
                                h(
                                    NButton,
                                    {
                                        secondary: true,
                                        type: 'primary',
                                        onClick: () => {
                                            prefStore.setAsWelcomed(true)
                                            n.destroy()
                                        },
                                    },
                                    {
                                        default: () => i18n.t('dialogue.welcome.accept'),
                                    },
                                ),
                            ],
                        },
                    ),
            })
        }
    } finally {
        initializing.value = false
    }
}

const onUnauthorized = () => {
    if (authEnabled.value) {
        authenticated.value = false
        setViewport('mobile')
    }
}

let resizeTimer = null
const onOrientationChange = () => {
    if (!authenticated.value) return
    clearTimeout(resizeTimer)
    resizeTimer = setTimeout(() => setViewport('desktop'), 200)
}

onMounted(async () => {
    if (isWebMode) {
        // Apply saved login theme before auth check to prevent flash
        const savedTheme = localStorage.getItem('rdm_login_theme')
        if (savedTheme && ['auto', 'light', 'dark'].includes(savedTheme)) {
            prefStore.general.theme = savedTheme
        }
        window.addEventListener('rdm:unauthorized', onUnauthorized)
        window.addEventListener('orientationchange', onOrientationChange)
        window.addEventListener('resize', onOrientationChange)
        await checkAuth()
        if (authEnabled.value && !authenticated.value) {
            // Not authenticated — show login page, do NOT call any API
            setViewport('mobile')
        } else {
            setViewport('desktop')
            // Connect WebSocket before initApp
            try {
                const runtime = await import('wailsjs/runtime/runtime.js')
                if (runtime.ReconnectWebSocket) await runtime.ReconnectWebSocket()
            } catch {}
            await initApp()
        }
    } else {
        // Desktop mode: original Wails flow, no auth needed
        await initApp()
    }
})

onUnmounted(() => {
    if (isWebMode) {
        window.removeEventListener('rdm:unauthorized', onUnauthorized)
        window.removeEventListener('orientationchange', onOrientationChange)
        window.removeEventListener('resize', onOrientationChange)
    }
})

// watch theme and dynamically switch
watch(
    () => prefStore.isDark,
    (isDark) => (isDark ? WindowSetDarkTheme() : WindowSetLightTheme()),
)

// watch language and dynamically switch
watch(
    () => prefStore.general.language,
    (lang) => (i18n.locale.value = prefStore.currentLanguage),
)
</script>

<template>
    <n-config-provider
        :inline-theme-disabled="true"
        :locale="prefStore.themeLocale"
        :theme="prefStore.isDark ? darkTheme : undefined"
        :theme-overrides="prefStore.isDark ? darkThemeOverrides : themeOverrides"
        class="fill-height">
        <!-- Web mode: auth gate -->
        <template v-if="isWebMode && authChecking">
            <div style="width: 100vw; height: 100vh"></div>
        </template>
        <template v-else-if="isWebMode && authEnabled && !authenticated">
            <component :is="LoginPage" @login="onLogin" />
        </template>
        <template v-else>
            <n-dialog-provider>
                <app-content :loading="initializing" />

                <!-- top modal dialogs -->
                <connection-dialog />
                <group-dialog />
                <new-key-dialog />
                <key-filter-dialog />
                <add-fields-dialog />
                <rename-key-dialog />
                <delete-key-dialog />
                <export-key-dialog />
                <import-key-dialog />
                <flush-db-dialog />
                <set-ttl-dialog />
                <preferences-dialog />
                <decoder-dialog />
                <about-dialog />
            </n-dialog-provider>
        </template>
    </n-config-provider>
</template>

<style lang="scss"></style>
