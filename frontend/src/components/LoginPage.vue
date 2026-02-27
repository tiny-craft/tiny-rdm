<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useThemeVars } from 'naive-ui'
import iconUrl from '@/assets/images/icon.png'
import usePreferencesStore from '@/stores/preferences.js'
import LangIcon from '@/components/icons/Lang.vue'
import Sun from '@/components/icons/Sun.vue'
import Moon from '@/components/icons/Moon.vue'
import ThemeAuto from '@/components/icons/ThemeAuto.vue'

import { Login } from '@/utils/api.js'
import { lang } from '@/langs/index.js'
import { useI18n } from 'vue-i18n'
import { useRender } from '@/utils/render.js'

const themeVars = useThemeVars()
const prefStore = usePreferencesStore()
const i18n = useI18n()
const emit = defineEmits(['login'])

// --- Theme ---
const THEME_KEY = 'rdm_login_theme'
const themeMode = ref(localStorage.getItem(THEME_KEY) || 'auto')

onMounted(() => {
    prefStore.general.theme = themeMode.value
})

const getThemeLabels = (langKey) => {
    const l = lang[langKey] || lang.en
    const g = l.preferences?.general || {}
    return {
        auto: g.theme_auto || 'Auto',
        light: g.theme_light || 'Light',
        dark: g.theme_dark || 'Dark',
    }
}

const themeOptions = computed(() => {
    const labels = getThemeLabels(currentLang.value)
    return [
        { label: labels.light, key: 'light', icon: Sun },
        { label: labels.dark, key: 'dark', icon: Moon },
        { label: labels.auto, key: 'auto', icon: ThemeAuto },
    ]
})

const currentThemeLabel = computed(() => {
    const labels = getThemeLabels(currentLang.value)
    return labels[themeMode.value]
})

const onThemeSelect = (key) => {
    if (!['auto', 'light', 'dark'].includes(key)) {
        return
    }
    themeMode.value = key
    prefStore.general.theme = key
    localStorage.setItem(THEME_KEY, key)
}

// --- Language ---
const LANG_KEY = 'rdm_login_lang'
const langNames = Object.fromEntries(Object.entries(lang).map(([k, v]) => [k, v.name]))
const autoLabel = Object.fromEntries(
    Object.entries(lang).map(([k, v]) => [k, v.preferences?.general?.theme_auto || 'Auto']),
)

const detectSystemLang = () => {
    const sysLang = (navigator.language || '').toLowerCase()
    if (sysLang.startsWith('zh-tw') || sysLang.startsWith('zh-hant')) {
        return 'tw'
    }
    const prefix = sysLang.split('-')[0]
    return langNames[prefix] ? prefix : 'en'
}

const langSetting = ref(localStorage.getItem(LANG_KEY) || 'auto')
const currentLang = computed(() => (langSetting.value === 'auto' ? detectSystemLang() : langSetting.value))

const langOptions = computed(() => [
    { label: autoLabel[currentLang.value] || 'Auto', key: 'auto' },
    { type: 'divider' },
    ...Object.entries(langNames).map(([k, v]) => ({ label: v, key: k })),
])

const currentLangLabel = computed(() => {
    if (langSetting.value === 'auto') return autoLabel[currentLang.value] || 'Auto'
    return langNames[langSetting.value] || langSetting.value
})

const onLangSelect = (key) => {
    const valid = ['auto', ...Object.keys(langNames)]
    if (!valid.includes(key)) {
        return
    }
    langSetting.value = key
    localStorage.setItem(LANG_KEY, key)
}

const render = useRender()

// --- i18n ---
watch(
    currentLang,
    (val) => {
        i18n.locale.value = val
    },
    { immediate: true },
)

// --- Form ---
const username = ref('')
const password = ref('')
const loading = ref(false)
const errorMsg = ref('')

const canSubmit = computed(() => username.value.length > 0 && password.value.length > 0)

const handleLogin = async () => {
    if (!canSubmit.value || loading.value) return
    loading.value = true
    errorMsg.value = ''

    try {
        const { msg, success = false } = await Login(username.value, password.value)
        if (msg === 'too_many_attempts') {
            errorMsg.value = $t('login.too_many_attempts')
            return
        }
        if (!success) {
            errorMsg.value = $t('login.invalid_credentials')
            return
        }
        emit('login')
    } catch (e) {
        errorMsg.value = $t('login.network_error')
    } finally {
        loading.value = false
    }
}
</script>

<template>
    <div class="login-wrapper">
        <div class="login-card">
            <div class="login-header">
                <n-avatar :size="64" :src="iconUrl" color="#0000" />
                <div class="login-title">Tiny RDM</div>
                <!--                <n-text depth="3" style="font-size: 13px">Redis Web Manager</n-text>-->
            </div>

            <n-form class="login-form" @submit.prevent="handleLogin">
                <n-form-item :label="$t('dialogue.connection.usr')">
                    <n-input
                        v-model:value="username"
                        :placeholder="$t('login.username_placeholder')"
                        autofocus
                        size="large"
                        @keydown.enter="handleLogin" />
                </n-form-item>
                <n-form-item :label="$t('dialogue.connection.pwd')">
                    <n-input
                        v-model:value="password"
                        :placeholder="$t('login.password_placeholder')"
                        show-password-on="click"
                        size="large"
                        type="password"
                        @keydown.enter="handleLogin" />
                </n-form-item>

                <n-text v-if="errorMsg" style="font-size: 13px; margin-bottom: 12px; display: block" type="error">
                    {{ errorMsg }}
                </n-text>

                <n-button
                    :disabled="!canSubmit"
                    :loading="loading"
                    attr-type="submit"
                    block
                    size="large"
                    style="margin-top: 8px"
                    type="primary"
                    @click="handleLogin">
                    {{ $t('login.submit') }}
                </n-button>
            </n-form>

            <div class="login-toolbar">
                <n-dropdown :options="langOptions" size="small" trigger="hover" @select="onLangSelect">
                    <span class="toolbar-btn">
                        <n-icon :component="LangIcon" :size="14" />
                        <span>{{ currentLangLabel }}</span>
                    </span>
                </n-dropdown>
                <n-divider style="margin: 0 4px" vertical />
                <n-dropdown
                    :options="themeOptions"
                    :render-icon="({ icon }) => render.renderIcon(icon)"
                    size="small"
                    trigger="hover"
                    @select="onThemeSelect">
                    <span class="toolbar-btn">
                        <n-icon
                            :component="themeMode === 'dark' ? Moon : themeMode === 'light' ? Sun : ThemeAuto"
                            :size="14" />
                        <span>{{ currentThemeLabel }}</span>
                    </span>
                </n-dropdown>
                <template v-if="prefStore.appVersion">
                    <n-divider style="margin: 0 4px" vertical />
                    <a
                        class="toolbar-btn toolbar-link"
                        href="https://github.com/tiny-craft/tiny-rdm"
                        rel="noopener noreferrer"
                        target="_blank">
                        {{ prefStore.appVersion }}
                    </a>
                </template>
            </div>
        </div>
    </div>
</template>

<style lang="scss" scoped>
.login-wrapper {
    width: 100vw;
    height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: v-bind('themeVars.bodyColor');
    padding: 16px;
    box-sizing: border-box;
}

.login-card {
    width: 420px;
    max-width: 100%;
    padding: 48px 40px 36px;
    border-radius: 10px;
    border: 1px solid v-bind('themeVars.borderColor');
    background-color: v-bind('themeVars.cardColor');
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
    box-sizing: border-box;
}

.login-header {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 6px;
    margin-bottom: 40px;
}

.login-title {
    font-size: 24px;
    font-weight: 800;
    margin-top: 8px;
    color: v-bind('themeVars.textColor1');
}

.login-form {
    :deep(.n-form-item) {
        margin-bottom: 18px;
    }

    :deep(.n-form-item-label) {
        color: v-bind('themeVars.textColor1');
        font-weight: 500;
    }
}

.login-toolbar {
    display: flex;
    align-items: center;
    justify-content: center;
    margin-top: 28px;
    flex-wrap: wrap;
    gap: 2px;
}

.toolbar-btn {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    font-size: 13px;
    color: v-bind('themeVars.textColor3');
    cursor: pointer;
    padding: 2px 6px;
    border-radius: 4px;
    transition:
        color 0.2s,
        background-color 0.2s;
    user-select: none;
    white-space: nowrap;

    &:hover {
        color: v-bind('themeVars.textColor2');
        background-color: v-bind('themeVars.buttonColor2Hover');
    }
}

.toolbar-link {
    text-decoration: none;
    color: v-bind('themeVars.textColor3');

    &:hover {
        color: v-bind('themeVars.textColor2');
    }
}

@media (max-width: 480px) {
    .login-wrapper {
        align-items: flex-start;
        padding-top: 12vh;
    }

    .login-card {
        padding: 32px 24px 28px;
        border: none;
        border-radius: 12px;
    }

    .login-header {
        margin-bottom: 28px;
    }

    .login-toolbar {
        margin-top: 20px;
    }
}
</style>
