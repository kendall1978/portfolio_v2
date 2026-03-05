<template>
  <div class="max-w-4xl mx-auto p-4">
    <header class="flex justify-between items-center mb-8">
      <h1 class="text-2xl font-bold text-gray-900">Data Management</h1>
      <router-link to="/" class="text-indigo-600 hover:underline">Back to Dashboard</router-link>
    </header>

    <section class="bg-white p-6 rounded-lg shadow-sm mb-6">
      <h2 class="text-lg font-semibold text-gray-800 mb-2">Import OES Data</h2>
      <p class="text-gray-600 text-sm mb-4">Import salary data from BLS Occupational Employment Statistics files.</p>
      <div class="flex flex-col gap-3 mb-4">
        <label class="flex items-center gap-2 text-sm">
          <input type="checkbox" v-model="backfill" class="rounded" />
          <span>Backfill historical data (2018-present)</span>
        </label>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Local zip directory (optional)</label>
          <input v-model="localZipDir" type="text" placeholder="/home/user/Downloads"
            class="w-full px-3 py-2 border border-gray-300 rounded text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500" />
          <p class="text-xs text-gray-400 mt-1">Leave empty to download from BLS automatically.</p>
        </div>
      </div>
      <button @click="triggerScrape" :disabled="scraping"
        class="px-4 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700 disabled:opacity-60 disabled:cursor-not-allowed">
        {{ scraping ? 'Importing...' : 'Import Data' }}
      </button>
      <p v-if="scrapeMsg" class="mt-2 text-sm text-gray-600">{{ scrapeMsg }}</p>
    </section>

    <section class="bg-white p-6 rounded-lg shadow-sm">
      <div class="flex justify-between items-center mb-4">
        <h2 class="text-lg font-semibold text-gray-800">Scrape History</h2>
        <button @click="fetchLogs"
          class="px-3 py-1 bg-gray-200 text-gray-700 rounded hover:bg-gray-300 text-sm">
          Refresh Logs
        </button>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full border-collapse">
          <thead>
            <tr class="bg-gray-50">
              <th class="px-4 py-3 text-left font-semibold text-sm border-b border-gray-200">Source</th>
              <th class="px-4 py-3 text-left font-semibold text-sm border-b border-gray-200">Status</th>
              <th class="px-4 py-3 text-left font-semibold text-sm border-b border-gray-200">Records</th>
              <th class="px-4 py-3 text-left font-semibold text-sm border-b border-gray-200">Started</th>
              <th class="px-4 py-3 text-left font-semibold text-sm border-b border-gray-200">Completed</th>
              <th class="px-4 py-3 text-left font-semibold text-sm border-b border-gray-200">Error</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="log in logs" :key="log.id" class="border-b border-gray-100">
              <td class="px-4 py-3 text-sm">{{ log.source }}</td>
              <td class="px-4 py-3 text-sm">
                <span :class="{
                  'text-green-600': log.status === 'success',
                  'text-red-600': log.status === 'failed',
                  'text-yellow-600': log.status === 'running',
                }">{{ log.status }}</span>
              </td>
              <td class="px-4 py-3 text-sm">{{ log.records_found }}</td>
              <td class="px-4 py-3 text-sm">{{ formatDate(log.started_at) }}</td>
              <td class="px-4 py-3 text-sm">{{ log.completed_at ? formatDate(log.completed_at) : '-' }}</td>
              <td class="px-4 py-3 text-sm text-red-600">{{ log.error_message || '-' }}</td>
            </tr>
            <tr v-if="logs.length === 0">
              <td colspan="6" class="px-4 py-6 text-center text-gray-400 text-sm">No scrape history yet.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '@/api/client'

interface ScrapeLog {
  id: number
  source: string
  status: string
  records_found: number
  error_message: string
  started_at: string
  completed_at: string | null
}

const logs = ref<ScrapeLog[]>([])
const scraping = ref(false)
const scrapeMsg = ref('')
const backfill = ref(false)
const localZipDir = ref('')

onMounted(fetchLogs)

async function fetchLogs() {
  const { data } = await api.get('/scrape/logs')
  logs.value = data
}

async function triggerScrape() {
  scraping.value = true
  scrapeMsg.value = ''
  try {
    await api.post('/scrape/trigger', {
      backfill: backfill.value,
      local_zip_dir: localZipDir.value || undefined,
    })
    scrapeMsg.value = 'OES import started! Refresh logs to see progress.'
  } catch (e: any) {
    scrapeMsg.value = 'Failed to start import: ' + (e.response?.data?.error || e.message)
  } finally {
    scraping.value = false
  }
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleString()
}
</script>