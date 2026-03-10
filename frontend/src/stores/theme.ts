import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export const useThemeStore = defineStore('theme', () => {
    const isDark = ref(true)

    function init() {
        const saved = localStorage.getItem('theme')
        if (saved) {
            isDark.value = saved === 'dark'
        } else {
            // Default to dark, but respect system preference if set
            isDark.value = !window.matchMedia('(prefers-color-scheme: light)').matches
        }
        applyTheme()
    }

    function toggle() {
        isDark.value = !isDark.value
        localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
        applyTheme()
    }

    function applyTheme() {
        document.documentElement.classList.toggle('dark', isDark.value)
    }

    watch(isDark, applyTheme)

    return { isDark, init, toggle }
})
