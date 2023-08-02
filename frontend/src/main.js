import { createPinia } from 'pinia'
import { createApp } from 'vue'
import { createI18n } from 'vue-i18n'
import App from './App.vue'
import { lang } from './langs'
import './styles/style.scss'
import dayjs from 'dayjs'
import duration from 'dayjs/plugin/duration'
import relativeTime from 'dayjs/plugin/relativeTime'

dayjs.extend(duration)
dayjs.extend(relativeTime)

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
    }),
)
app.use(createPinia())
app.mount('#app')
