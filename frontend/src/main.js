import { createPinia } from 'pinia'
import { createApp } from 'vue'
import { createI18n } from 'vue-i18n'
import App from './App.vue'
import { lang } from './langs'
import './style.scss'

const app = createApp(App)
app.use(
    createI18n({
        locale: 'en',
        fallbackLocale: 'en',
        globalInjection: true,
        legacy: false,
        messages: {
            ...lang,
        },
    })
)
app.use(createPinia())
app.mount('#app')
