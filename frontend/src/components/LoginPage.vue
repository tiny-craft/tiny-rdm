<script setup>
// --- Render icon helper for dropdown items ---
import { computed, h, onMounted, ref } from 'vue'
import { NIcon, useThemeVars } from 'naive-ui'
import iconUrl from '@/assets/images/icon.png'
import usePreferencesStore from '@/stores/preferences.js'
import LangIcon from '@/components/icons/Lang.vue'

import { Login } from '@/utils/api.js'
import { lang } from '@/langs/index.js'

const themeVars = useThemeVars()
const prefStore = usePreferencesStore()
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
        { label: labels.light, key: 'light', icon: renderIcon('sun') },
        { label: labels.dark, key: 'dark', icon: renderIcon('moon') },
        { label: labels.auto, key: 'auto', icon: renderIcon('auto') },
    ]
})

const currentThemeLabel = computed(() => {
    const labels = getThemeLabels(currentLang.value)
    return labels[themeMode.value]
})

const onThemeSelect = (key) => {
    if (!['auto', 'light', 'dark'].includes(key)) return
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
    if (!valid.includes(key)) return
    langSetting.value = key
    localStorage.setItem(LANG_KEY, key)
}

const SunSvg = {
    render: () =>
        h(
            'svg',
            {
                xmlns: 'http://www.w3.org/2000/svg',
                viewBox: '0 0 24 24',
                width: '1em',
                height: '1em',
                fill: 'none',
                stroke: 'currentColor',
                'stroke-width': '2',
                'stroke-linecap': 'round',
                'stroke-linejoin': 'round',
            },
            [
                h('circle', { cx: '12', cy: '12', r: '5' }),
                h('path', {
                    d: 'M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42',
                }),
            ],
        ),
}
const MoonSvg = {
    render: () =>
        h(
            'svg',
            {
                xmlns: 'http://www.w3.org/2000/svg',
                viewBox: '0 0 24 24',
                width: '1em',
                height: '1em',
                fill: 'none',
                stroke: 'currentColor',
                'stroke-width': '2',
                'stroke-linecap': 'round',
                'stroke-linejoin': 'round',
            },
            [h('path', { d: 'M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z' })],
        ),
}
const AutoSvg = {
    render: () =>
        h(
            'svg',
            {
                xmlns: 'http://www.w3.org/2000/svg',
                viewBox: '0 0 24 24',
                width: '1em',
                height: '1em',
                fill: 'none',
                stroke: 'currentColor',
                'stroke-width': '2',
                'stroke-linecap': 'round',
                'stroke-linejoin': 'round',
            },
            [h('circle', { cx: '12', cy: '12', r: '10' }), h('path', { d: 'M12 2a10 10 0 0 1 0 20V2' })],
        ),
}
const LangSvg = {
    render: () =>
        h(
            'svg',
            {
                xmlns: 'http://www.w3.org/2000/svg',
                viewBox: '0 0 24 24',
                width: '1em',
                height: '1em',
                fill: 'none',
                stroke: 'currentColor',
                'stroke-width': '2',
                'stroke-linecap': 'round',
                'stroke-linejoin': 'round',
            },
            [
                h('path', { d: 'M5 8l6 6' }),
                h('path', { d: 'M4 14l6-6 2-3' }),
                h('path', { d: 'M2 5h12' }),
                h('path', { d: 'M7 2h1' }),
                h('path', { d: 'M22 22l-5-10-5 10' }),
                h('path', { d: 'M14 18h6' }),
            ],
        ),
}

const iconMap = { sun: SunSvg, moon: MoonSvg, auto: AutoSvg, lang: LangSvg }
const renderIcon = (name) => () => h(NIcon, null, { default: () => h(iconMap[name]) })

// --- i18n texts ---
const langTexts = {
    zh: {
        title: '登录',
        username: '用户名',
        password: '密码',
        usernamePh: '请输入用户名',
        passwordPh: '请输入密码',
        submit: '登 录',
        tooMany: '尝试次数过多，请稍后再试',
        failed: '用户名或密码错误',
        network: '网络错误',
    },
    tw: {
        title: '登入',
        username: '使用者名稱',
        password: '密碼',
        usernamePh: '請輸入使用者名稱',
        passwordPh: '請輸入密碼',
        submit: '登 入',
        tooMany: '嘗試次數過多，請稍後再試',
        failed: '使用者名稱或密碼錯誤',
        network: '網路錯誤',
    },
    ja: {
        title: 'ログイン',
        username: 'ユーザー名',
        password: 'パスワード',
        usernamePh: 'ユーザー名を入力',
        passwordPh: 'パスワードを入力',
        submit: 'ログイン',
        tooMany: '試行回数が多すぎます',
        failed: 'ユーザー名またはパスワードが正しくありません',
        network: 'ネットワークエラー',
    },
    ko: {
        title: '로그인',
        username: '사용자 이름',
        password: '비밀번호',
        usernamePh: '사용자 이름 입력',
        passwordPh: '비밀번호 입력',
        submit: '로그인',
        tooMany: '시도 횟수 초과, 잠시 후 다시 시도하세요',
        failed: '사용자 이름 또는 비밀번호가 올바르지 않습니다',
        network: '네트워크 오류',
    },
    es: {
        title: 'Iniciar sesión',
        username: 'Usuario',
        password: 'Contraseña',
        usernamePh: 'Ingrese usuario',
        passwordPh: 'Ingrese contraseña',
        submit: 'Entrar',
        tooMany: 'Demasiados intentos, intente más tarde',
        failed: 'Credenciales inválidas',
        network: 'Error de red',
    },
    fr: {
        title: 'Connexion',
        username: "Nom d'utilisateur",
        password: 'Mot de passe',
        usernamePh: "Entrez le nom d'utilisateur",
        passwordPh: 'Entrez le mot de passe',
        submit: 'Se connecter',
        tooMany: 'Trop de tentatives, réessayez plus tard',
        failed: 'Identifiants invalides',
        network: 'Erreur réseau',
    },
    ru: {
        title: 'Вход',
        username: 'Имя пользователя',
        password: 'Пароль',
        usernamePh: 'Введите имя пользователя',
        passwordPh: 'Введите пароль',
        submit: 'Войти',
        tooMany: 'Слишком много попыток, попробуйте позже',
        failed: 'Неверные учётные данные',
        network: 'Ошибка сети',
    },
    pt: {
        title: 'Entrar',
        username: 'Usuário',
        password: 'Senha',
        usernamePh: 'Digite o usuário',
        passwordPh: 'Digite a senha',
        submit: 'Entrar',
        tooMany: 'Muitas tentativas, tente novamente mais tarde',
        failed: 'Credenciais inválidas',
        network: 'Erro de rede',
    },
    tr: {
        title: 'Giriş',
        username: 'Kullanıcı adı',
        password: 'Şifre',
        usernamePh: 'Kullanıcı adını girin',
        passwordPh: 'Şifreyi girin',
        submit: 'Giriş Yap',
        tooMany: 'Çok fazla deneme, lütfen daha sonra tekrar deneyin',
        failed: 'Geçersiz kimlik bilgileri',
        network: 'Ağ hatası',
    },
    en: {
        title: 'Sign In',
        username: 'Username',
        password: 'Password',
        usernamePh: 'Enter username',
        passwordPh: 'Enter password',
        submit: 'Sign In',
        tooMany: 'Too many attempts, please try later',
        failed: 'Invalid credentials',
        network: 'Network error',
    },
}

const t = computed(() => langTexts[currentLang.value] || langTexts.en)

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
            errorMsg.value = t.value.tooMany
            return
        }
        if (!success) {
            errorMsg.value = t.value.failed
            return
        }
        emit('login')
    } catch (e) {
        errorMsg.value = t.value.network
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
                <n-text depth="3" style="font-size: 13px">Redis Web Manager</n-text>
            </div>

            <n-form class="login-form" @submit.prevent="handleLogin">
                <n-form-item :label="t.username">
                    <n-input
                        v-model:value="username"
                        :placeholder="t.usernamePh"
                        autofocus
                        size="large"
                        @keydown.enter="handleLogin" />
                </n-form-item>
                <n-form-item :label="t.password">
                    <n-input
                        v-model:value="password"
                        :placeholder="t.passwordPh"
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
                    {{ t.submit }}
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
                <n-dropdown :options="themeOptions" size="small" trigger="hover" @select="onThemeSelect">
                    <span class="toolbar-btn">
                        <n-icon
                            :component="themeMode === 'dark' ? MoonSvg : themeMode === 'light' ? SunSvg : AutoSvg"
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
