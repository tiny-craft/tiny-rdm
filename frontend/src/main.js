import { createPinia } from 'pinia'
import { createApp } from 'vue'
import App from './App.vue'
import './styles/style.scss'
import dayjs from 'dayjs'
import duration from 'dayjs/plugin/duration'
import relativeTime from 'dayjs/plugin/relativeTime'
import { i18n } from '@/utils/i18n.js'

dayjs.extend(duration)
dayjs.extend(relativeTime)

const app = createApp(App)
app.use(i18n)
app.use(createPinia())
app.mount('#app')
