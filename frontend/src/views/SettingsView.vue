<template>
  <div class="max-w-xl mx-auto px-4 py-8">
    <h1 class="text-2xl font-bold mb-8 pb-4 border-b" style="color: var(--text-primary); border-color: var(--border-color)">Settings</h1>

    <section class="p-6 rounded-xl border mb-6" :style="{ backgroundColor: 'var(--bg-card)', borderColor: 'var(--border-color)' }">
      <h2 class="text-lg font-semibold mb-4" style="color: var(--text-primary)">Profile</h2>
      <form @submit.prevent="updateProfile">
        <div class="mb-4">
          <label class="block mb-1 font-semibold text-sm" style="color: var(--text-secondary)">First Name</label>
          <input v-model="profile.first_name" type="text"
            class="w-full px-3 py-2 border rounded text-sm" />
        </div>
        <div class="mb-4">
          <label class="block mb-1 font-semibold text-sm" style="color: var(--text-secondary)">Last Name</label>
          <input v-model="profile.last_name" type="text"
            class="w-full px-3 py-2 border rounded text-sm" />
        </div>
        <div class="mb-4">
          <label class="block mb-1 font-semibold text-sm" style="color: var(--text-secondary)">Email</label>
          <input v-model="profile.email" type="email"
            class="w-full px-3 py-2 border rounded text-sm" />
        </div>
        <button type="submit"
          class="px-4 py-2 rounded text-white text-sm bg-accent hover:bg-accent-hover transition-colors">
          Save Profile
        </button>
        <p v-if="profileMsg" class="text-green-500 mt-2 text-sm">{{ profileMsg }}</p>
      </form>
    </section>

    <section class="p-6 rounded-xl border mb-6" :style="{ backgroundColor: 'var(--bg-card)', borderColor: 'var(--border-color)' }">
      <h2 class="text-lg font-semibold mb-4" style="color: var(--text-primary)">Current Salary</h2>
      <form @submit.prevent="saveSalary">
        <div class="mb-4">
          <label class="block mb-1 font-semibold text-sm" style="color: var(--text-secondary)">Base Salary ($)</label>
          <input v-model.number="salary" type="number" step="1000"
            class="w-full px-3 py-2 border rounded text-sm" />
        </div>
        <button type="submit"
          class="px-4 py-2 rounded text-white text-sm bg-accent hover:bg-accent-hover transition-colors">
          Save Salary
        </button>
        <p v-if="salaryMsg" class="text-green-500 mt-2 text-sm">{{ salaryMsg }}</p>
      </form>
    </section>

    <section class="p-6 rounded-xl border" :style="{ backgroundColor: 'var(--bg-card)', borderColor: 'var(--border-color)' }">
      <h2 class="text-lg font-semibold mb-4" style="color: var(--text-primary)">Two-Factor Authentication</h2>

      <div v-if="!mfaEnabled && !qrUrl">
        <p class="text-sm mb-3" style="color: var(--text-secondary)">Protect your account with Google Authenticator.</p>
        <button @click="setupMFA"
          class="px-4 py-2 rounded text-white text-sm bg-accent hover:bg-accent-hover transition-colors">
          Enable MFA
        </button>
      </div>

      <div v-if="qrUrl">
        <p class="text-sm mb-3" style="color: var(--text-secondary)">Scan this QR code with Google Authenticator:</p>
        <img :src="qrCodeImage" alt="QR Code" class="mb-4" />
        <p class="text-sm mb-4" style="color: var(--text-secondary)">Or enter this secret manually:
          <code class="px-2 py-1 rounded text-sm" :style="{ backgroundColor: 'var(--bg-secondary)', color: 'var(--text-primary)' }">{{ mfaSecret }}</code>
        </p>
        <h3 class="font-semibold mb-2" style="color: var(--text-primary)">Recovery Codes</h3>
        <p class="text-sm mb-2" style="color: var(--text-secondary)">Save these codes somewhere safe. Each can only be used once:</p>
        <ul class="list-disc list-inside mb-4">
          <li v-for="code in recoveryCodes" :key="code">
            <code class="px-2 py-1 rounded text-sm" :style="{ backgroundColor: 'var(--bg-secondary)', color: 'var(--text-primary)' }">{{ code }}</code>
          </li>
        </ul>
      </div>

      <div v-if="mfaEnabled && !qrUrl">
        <p class="text-green-500 text-sm font-semibold">MFA is enabled.</p>
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