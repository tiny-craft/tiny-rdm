import { createI18n } from 'vue-i18n'
import { lang } from '@/langs/index.js'

export const i18n = createI18n({
    locale: 'en',
    fallbackLocale: 'en',
    globalInjection: true,
    legacy: false,
    messages: {
        ...lang,
    },
})

export const i18nGlobal = i18n.global
