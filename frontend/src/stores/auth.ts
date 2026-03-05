import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/api/client'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const isAuthenticated = ref(!!token.value)
  const mfaRequired = ref(false)
  const challengeToken = ref<string | null>(null)

  const registrationOpen = ref(false)

  async function checkStatus() {
    const { data } = await api.get('/auth/status')
    registrationOpen.value = data.registration_open
  }

  async function register(email: string, password: string, firstName: string, lastName: string) {
    await api.post('/auth/register', {
      email,
      password,
      first_name: firstName,
      last_name: lastName,
    })
    // Auto-login after registration
    return login(email, password)
  }

  async function login(email: string, password: string) {
    const { data } = await api.post('/auth/login', { email, password })
    if (data.mfa_required) {
      mfaRequired.value = true
      challengeToken.value = data.challenge_token
      return { mfaRequired: true }
    }
    setToken(data.token)
    return { mfaRequired: false }
  }

  async function verifyMFA(code: string) {
    const { data } = await api.post('/auth/mfa/verify', {
      challenge_token: challengeToken.value,
      code,
    })
    setToken(data.token)
    mfaRequired.value = false
    challengeToken.value = null
  }

  async function recoverMFA(recoveryCode: string) {
    const { data } = await api.post('/auth/mfa/recover', {
      challenge_token: challengeToken.value,
      recovery_code: recoveryCode,
    })
    setToken(data.token)
    mfaRequired.value = false
    challengeToken.value = null
  }

  function setToken(newToken: string) {
    token.value = newToken
    localStorage.setItem('token', newToken)
    isAuthenticated.value = true
  }

  function logout() {
    token.value = ''
    localStorage.removeItem('token')
    isAuthenticated.value = false
  }

  return { token, isAuthenticated, mfaRequired, challengeToken, registrationOpen, checkStatus, register, login, verifyMFA, recoverMFA, logout }
})