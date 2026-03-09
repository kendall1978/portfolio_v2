import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)

app.mount('#app')

// Initialize theme after app is mounted
import { useThemeStore } from '@/stores/theme'
const themeStore = useThemeStore()
themeStore.init()
