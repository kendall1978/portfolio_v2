<template>
  <div class="min-h-screen" style="background-color: var(--bg-primary); color: var(--text-primary)">
    <nav class="border-b" style="border-color: var(--border-color); background-color: var(--bg-secondary)">
      <div class="max-w-6xl mx-auto px-4 h-16 flex items-center justify-between">
        <router-link to="/" class="text-lg font-semibold" style="color: var(--text-primary)">
          Kendall Roberts
        </router-link>

        <div class="flex items-center gap-6">
          <router-link to="/" class="text-sm hover:text-accent transition-colors" style="color: var(--text-secondary)">Home</router-link>
          <router-link to="/blog" class="text-sm hover:text-accent transition-colors" style="color: var(--text-secondary)">Blog</router-link>

          <template v-if="authStore.isAuthenticated">
            <router-link to="/dashboard" class="text-sm hover:text-accent transition-colors" style="color: var(--text-secondary)">Dashboard</router-link>
            <router-link to="/settings" class="text-sm hover:text-accent transition-colors" style="color: var(--text-secondary)">Settings</router-link>
            <router-link to="/admin" class="text-sm hover:text-accent transition-colors" style="color: var(--text-secondary)">Admin</router-link>
            <button @click="logout" class="text-sm hover:text-accent transition-colors" style="color: var(--text-secondary)">Logout</button>
          </template>
          <template v-else>
            <router-link to="/login" class="text-sm hover:text-accent transition-colors" style="color: var(--text-secondary)">Login</router-link>
          </template>

          <ThemeToggle />
        </div>
      </div>
    </nav>

    <main>
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import ThemeToggle from '@/components/ThemeToggle.vue'

const authStore = useAuthStore()
const router = useRouter()

function logout() {
    authStore.logout()
    router.push('/login')
}
</script>
