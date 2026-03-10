<template>
  <div class="max-w-4xl mx-auto px-4 py-8">
    <h1 class="text-2xl font-bold mb-8 pb-4 border-b" style="color: var(--text-primary); border-color: var(--border-color)">Data Management</h1>

    <section class="p-6 rounded-xl border mb-6" :style="{ backgroundColor: 'var(--bg-card)', borderColor: 'var(--border-color)' }">
      <h2 class="text-lg font-semibold mb-2" style="color: var(--text-primary)">Import OES Data</h2>
      <p class="text-sm mb-4" style="color: var(--text-secondary)">Import salary data from BLS Occupational Employment Statistics files.</p>
      <div class="flex flex-col gap-3 mb-4">
        <label class="flex items-center gap-2 text-sm" style="color: var(--text-secondary)">
          <input type="checkbox" v-model="backfill" class="rounded accent-[var(--accent-color)]" />
          <span>Backfill historical data (2018-present)</span>
        </label>
        <div>
          <label class="block text-sm font-medium mb-1" style="color: var(--text-secondary)">Local zip directory (optional)</label>
          <input v-model="localZipDir" type="text" placeholder="/home/user/Downloads"
            class="w-full px-3 py-2 border rounded text-sm" />
          <p class="text-xs mt-1" style="color: var(--text-muted)">Leave empty to download from BLS automatically.</p>
        </div>
      </div>
      <button @click="triggerScrape" :disabled="scraping"
        class="px-4 py-2 rounded text-white text-sm bg-accent hover:bg-accent-hover disabled:opacity-60 disabled:cursor-not-allowed transition-colors">
        {{ scraping ? 'Importing...' : 'Import Data' }}
      </button>
      <p v-if="scrapeMsg" class="mt-2 text-sm" style="color: var(--text-secondary)">{{ scrapeMsg }}</p>
    </section>

    <section class="p-6 rounded-xl border" :style="{ backgroundColor: 'var(--bg-card)', borderColor: 'var(--border-color)' }">
      <div class="flex justify-between items-center mb-4">
        <h2 class="text-lg font-semibold" style="color: var(--text-primary)">Scrape History</h2>
        <button @click="fetchLogs"
          class="px-3 py-1.5 rounded text-sm transition-colors"
          :style="{ backgroundColor: 'var(--bg-secondary)', color: 'var(--text-secondary)' }">
          Refresh Logs
        </button>
      </div>
      <div class="overflow-x-auto rounded-lg border" :style="{ borderColor: 'var(--border-color)' }">
        <table class="w-full border-collapse">
          <thead>
            <tr :style="{ backgroundColor: 'var(--bg-secondary)' }">
              <th class="px-4 py-3 text-left font-semibold text-sm" style="color: var(--text-secondary)">Source</th>
              <th class="px-4 py-3 text-left font-semibold text-sm" style="color: var(--text-secondary)">Status</th>
              <th class="px-4 py-3 text-left font-semibold text-sm" style="color: var(--text-secondary)">Records</th>
              <th class="px-4 py-3 text-left font-semibold text-sm" style="color: var(--text-secondary)">Started</th>
              <th class="px-4 py-3 text-left font-semibold text-sm" style="color: var(--text-secondary)">Completed</th>
              <th class="px-4 py-3 text-left font-semibold text-sm" style="color: var(--text-secondary)">Error</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="log in logs" :key="log.id"
              class="border-t transition-colors"
              :style="{ borderColor: 'var(--border-color)' }"
              style="cursor: default;"
              @mouseenter="($event.currentTarget as HTMLElement).style.backgroundColor = 'var(--bg-secondary)'"
              @mouseleave="($event.currentTarget as HTMLElement).style.backgroundColor = 'transparent'">
              <td class="px-4 py-3 text-sm" style="color: var(--text-primary)">{{ log.source }}</td>
              <td class="px-4 py-3 text-sm">
                <span :class="{
                  'text-green-500': log.status === 'success',
                  'text-red-500': log.status === 'failed',
                  'text-yellow-500': log.status === 'running',
                }">{{ log.status }}</span>
              </td>
              <td class="px-4 py-3 text-sm" style="color: var(--text-primary)">{{ log.records_found }}</td>
              <td class="px-4 py-3 text-sm" style="color: var(--text-secondary)">{{ formatDate(log.started_at) }}</td>
              <td class="px-4 py-3 text-sm" style="color: var(--text-secondary)">{{ log.completed_at ? formatDate(log.completed_at) : '-' }}</td>
              <td class="px-4 py-3 text-sm text-red-500">{{ log.error_message || '-' }}</td>
            </tr>
            <tr v-if="logs.length === 0">
              <td colspan="6" class="px-4 py-6 text-center text-sm" style="color: var(--text-muted)">No scrape history yet.</td>
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