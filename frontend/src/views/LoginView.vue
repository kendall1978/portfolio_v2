<template>
  <div class="flex justify-center items-center min-h-screen bg-gray-100">
    <div class="bg-white p-8 rounded-lg shadow-md w-full max-w-md">
      <h1 class="text-2xl font-bold mb-6 text-center text-gray-900">Kendall Roberts Portfolio</h1>

      <div v-if="loading" class="text-center text-gray-500">Loading...</div>

      <!-- Registration form (first-time setup) -->
      <form v-else-if="authStore.registrationOpen" @submit.prevent="handleRegister">
        <p class="text-sm text-gray-500 mb-4 text-center">Create your account to get started.</p>
        <div class="mb-4">
          <label for="firstName" class="block mb-1 font-semibold text-gray-700">First Name</label>
          <input id="firstName" v-model="firstName" type="text" required
            class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-indigo-500" />
        </div>
        <div class="mb-4">
          <label for="lastName" class="block mb-1 font-semibold text-gray-700">Last Name</label>
          <input id="lastName" v-model="lastName" type="text" required
            class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-indigo-500" />
        </div>
        <div class="mb-4">
          <label for="regEmail" class="block mb-1 font-semibold text-gray-700">Email</label>
          <input id="regEmail" v-model="email" type="email" required
            class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-indigo-500" />
        </div>
        <div class="mb-4">
          <label for="regPassword" class="block mb-1 font-semibold text-gray-700">Password</label>
          <input id="regPassword" v-model="password" type="password" required minlength="8"
            class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-indigo-500" />
        </div>
        <button type="submit" :disabled="submitting"
          class="w-full py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700 disabled:opacity-60 disabled:cursor-not-allowed">
          {{ submitting ? 'Creating account...' : 'Create Account' }}
        </button>
        <p v-if="error" class="text-red-600 mt-2 text-sm">{{ error }}</p>
      </form>

      <!-- Login form -->
      <form v-else-if="!authStore.mfaRequired" @submit.prevent="handleLogin">
        <div class="mb-4">
          <label for="email" class="block mb-1 font-semibold text-gray-700">Email</label>
          <input id="email" v-model="email" type="email" required
            class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-indigo-500" />
        </div>
        <div class="mb-4">
          <label for="password" class="block mb-1 font-semibold text-gray-700">Password</label>
          <input id="password" v-model="password" type="password" required
            class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-indigo-500" />
        </div>
        <button type="submit" :disabled="submitting"
          class="w-full py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700 disabled:opacity-60 disabled:cursor-not-allowed">
          {{ submitting ? 'Signing in...' : 'Sign In' }}
        </button>
        <p v-if="error" class="text-red-600 mt-2 text-sm">{{ error }}</p>
      </form>

      <!-- MFA form -->
      <div v-else>
        <form v-if="!showRecovery" @submit.prevent="handleMFA">
          <div class="mb-4">
            <label for="code" class="block mb-1 font-semibold text-gray-700">Authenticator Code</label>
            <input id="code" v-model="mfaCode" type="text" inputmode="numeric" maxlength="6" required
              class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-indigo-500" />
          </div>
          <button type="submit" :disabled="submitting"
            class="w-full py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700 disabled:opacity-60 disabled:cursor-not-allowed">
            Verify
          </button>
          <button type="button" @click="showRecovery = true"
            class="w-full py-2 mt-2 bg-gray-200 text-gray-700 rounded hover:bg-gray-300">
            Use Recovery Code
          </button>
          <p v-if="error" class="text-red-600 mt-2 text-sm">{{ error }}</p>
        </form>

        <form v-else @submit.prevent="handleRecovery">
          <div class="mb-4">
            <label for="recovery" class="block mb-1 font-semibold text-gray-700">Recovery Code</label>
            <input id="recovery" v-model="recoveryCode" type="text" required
              class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-indigo-500" />
          </div>
          <button type="submit" :disabled="submitting"
            class="w-full py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700 disabled:opacity-60 disabled:cursor-not-allowed">
            Recover
          </button>
          <button type="button" @click="showRecovery = false"
            class="w-full py-2 mt-2 bg-gray-200 text-gray-700 rounded hover:bg-gray-300">
            Back to Code
          </button>
          <p v-if="error" class="text-red-600 mt-2 text-sm">{{ error }}</p>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const router = useRouter()

const email = ref('')
const password = ref('')
const firstName = ref('')
const lastName = ref('')
const mfaCode = ref('')
const recoveryCode = ref('')
const loading = ref(true)
const submitting = ref(false)
const error = ref('')
const showRecovery = ref(false)

onMounted(async () => {
  try {
    await authStore.checkStatus()
  } finally {
    loading.value = false
  }
})

async function handleRegister() {
  submitting.value = true
  error.value = ''
  try {
    const result = await authStore.register(email.value, password.value, firstName.value, lastName.value)
    if (!result.mfaRequired) {
      router.push('/')
    }
  } catch (e: any) {
    error.value = e.response?.data?.error || 'Registration failed'
  } finally {
    submitting.value = false
  }
}

async function handleLogin() {
  submitting.value = true
  error.value = ''
  try {
    const result = await authStore.login(email.value, password.value)
    if (!result.mfaRequired) {
      router.push('/')
    }
  } catch (e: any) {
    error.value = e.response?.data?.error || 'Login failed'
  } finally {
    submitting.value = false
  }
}

async function handleMFA() {
  submitting.value = true
  error.value = ''
  try {
    await authStore.verifyMFA(mfaCode.value)
    router.push('/')
  } catch (e: any) {
    error.value = e.response?.data?.error || 'Invalid code'
  } finally {
    submitting.value = false
  }
}

async function handleRecovery() {
  submitting.value = true
  error.value = ''
  try {
    await authStore.recoverMFA(recoveryCode.value)
    router.push('/')
  } catch (e: any) {
    error.value = e.response?.data?.error || 'Invalid recovery code'
  } finally {
    submitting.value = false
  }
}
</script>
