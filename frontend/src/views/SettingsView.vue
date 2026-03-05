<template>
  <div class="max-w-xl mx-auto p-4">
    <header class="flex justify-between items-center mb-8">
      <h1 class="text-2xl font-bold text-gray-900">Settings</h1>
      <router-link to="/" class="text-indigo-600 hover:underline">Back to Dashboard</router-link>
    </header>

    <section class="bg-white p-6 rounded-lg shadow-sm mb-6">
      <h2 class="text-lg font-semibold text-gray-800 mb-4">Profile</h2>
      <form @submit.prevent="updateProfile">
        <div class="mb-4">
          <label class="block mb-1 font-semibold text-gray-700">First Name</label>
          <input v-model="profile.first_name" type="text"
            class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-indigo-500" />
        </div>
        <div class="mb-4">
          <label class="block mb-1 font-semibold text-gray-700">Last Name</label>
          <input v-model="profile.last_name" type="text"
            class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-indigo-500" />
        </div>
        <div class="mb-4">
          <label class="block mb-1 font-semibold text-gray-700">Email</label>
          <input v-model="profile.email" type="email"
            class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-indigo-500" />
        </div>
        <button type="submit"
          class="px-4 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700">
          Save Profile
        </button>
        <p v-if="profileMsg" class="text-green-600 mt-2 text-sm">{{ profileMsg }}</p>
      </form>
    </section>

    <section class="bg-white p-6 rounded-lg shadow-sm mb-6">
      <h2 class="text-lg font-semibold text-gray-800 mb-4">Current Salary</h2>
      <form @submit.prevent="saveSalary">
        <div class="mb-4">
          <label class="block mb-1 font-semibold text-gray-700">Base Salary ($)</label>
          <input v-model.number="salary" type="number" step="1000"
            class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-indigo-500" />
        </div>
        <button type="submit"
          class="px-4 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700">
          Save Salary
        </button>
        <p v-if="salaryMsg" class="text-green-600 mt-2 text-sm">{{ salaryMsg }}</p>
      </form>
    </section>

    <section class="bg-white p-6 rounded-lg shadow-sm">
      <h2 class="text-lg font-semibold text-gray-800 mb-4">Two-Factor Authentication</h2>

      <div v-if="!mfaEnabled && !qrUrl">
        <p class="text-gray-600 mb-3 text-sm">Protect your account with Google Authenticator.</p>
        <button @click="setupMFA"
          class="px-4 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700">
          Enable MFA
        </button>
      </div>

      <div v-if="qrUrl">
        <p class="text-gray-600 mb-3 text-sm">Scan this QR code with Google Authenticator:</p>
        <img :src="qrCodeImage" alt="QR Code" class="mb-4" />
        <p class="text-sm text-gray-600 mb-4">Or enter this secret manually: <code class="bg-gray-100 px-2 py-1 rounded text-sm">{{ mfaSecret }}</code></p>
        <h3 class="font-semibold text-gray-800 mb-2">Recovery Codes</h3>
        <p class="text-sm text-gray-600 mb-2">Save these codes somewhere safe. Each can only be used once:</p>
        <ul class="list-disc list-inside mb-4">
          <li v-for="code in recoveryCodes" :key="code">
            <code class="bg-gray-100 px-2 py-1 rounded text-sm">{{ code }}</code>
          </li>
        </ul>
      </div>

      <div v-if="mfaEnabled && !qrUrl">
        <p class="text-green-600 text-sm font-semibold">MFA is enabled.</p>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '@/api/client'

const profile = ref({ first_name: '', last_name: '', email: '' })
const salary = ref(0)
const compensationId = ref<number | null>(null)
const profileMsg = ref('')
const salaryMsg = ref('')
const mfaEnabled = ref(false)
const qrUrl = ref('')
const mfaSecret = ref('')
const recoveryCodes = ref<string[]>([])
const qrCodeImage = ref('')

onMounted(async () => {
  const [profileRes, compRes] = await Promise.all([
    api.get('/profile'),
    api.get('/compensation'),
  ])
  profile.value = profileRes.data
  mfaEnabled.value = profileRes.data.mfa_enabled ?? false

  if (compRes.data.length > 0) {
    salary.value = compRes.data[0].base_salary
    compensationId.value = compRes.data[0].id
  }
})

async function updateProfile() {
  await api.put('/profile', profile.value)
  profileMsg.value = 'Profile updated!'
  setTimeout(() => (profileMsg.value = ''), 3000)
}

async function saveSalary() {
  const payload = {
    base_salary: salary.value,
    effective_date: new Date().toISOString(),
    notes: 'Updated from settings',
  }

  if (compensationId.value) {
    await api.put(`/compensation/${compensationId.value}`, payload)
  } else {
    const { data } = await api.post('/compensation', payload)
    compensationId.value = data.id
  }

  salaryMsg.value = 'Salary saved!'
  setTimeout(() => (salaryMsg.value = ''), 3000)
}

async function setupMFA() {
  const { data } = await api.post('/auth/mfa/setup')
  qrUrl.value = data.qr_url
  mfaSecret.value = data.secret
  recoveryCodes.value = data.recovery_codes
  mfaEnabled.value = true
  qrCodeImage.value = `https://api.qrserver.com/v1/create-qr-code/?data=${encodeURIComponent(data.qr_url)}&size=200x200`
}
</script>