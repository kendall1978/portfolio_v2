<template>
  <div class="max-w-6xl mx-auto p-4">
    <h1 class="text-2xl font-bold mb-8 pb-4 border-b" style="color: var(--text-primary); border-color: var(--border-color)">Salary Dashboard</h1>

    <!-- Filters -->
    <section class="bg-white p-4 rounded-lg shadow-sm mb-6">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <!-- Occupation filter -->
        <div>
          <label class="block mb-1 text-sm font-semibold text-gray-700">Occupations</label>
          <input v-model="occSearch" type="text" placeholder="Search occupations..."
            class="w-full px-3 py-2 border border-gray-300 rounded text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 mb-1" />
          <select v-model="salaryStore.selectedOccCodes" multiple
            class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-indigo-500 text-sm"
            style="height: 160px">
            <option v-for="occ in filteredOccupations" :key="occ.occ_code" :value="occ.occ_code">
              {{ occ.occ_title }}
            </option>
          </select>
          <div v-if="selectedOccLabels.length" class="flex flex-wrap gap-1 mt-1">
            <span v-for="sel in selectedOccLabels" :key="sel.code"
              class="inline-flex items-center gap-1 px-2 py-0.5 bg-indigo-100 text-indigo-700 rounded text-xs">
              {{ sel.title }}
              <button @click="removeOcc(sel.code)" class="hover:text-indigo-900">&times;</button>
            </span>
          </div>
          <p class="text-xs text-gray-400 mt-1">Hold Ctrl/Cmd to select multiple</p>
        </div>

        <!-- State filter -->
        <div>
          <label class="block mb-1 text-sm font-semibold text-gray-700">States</label>
          <select v-model="salaryStore.selectedStates" multiple
            class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-indigo-500 text-sm"
            style="min-height: 80px">
            <option v-for="r in stateList" :key="r.state_code" :value="r.state_code">
              {{ r.state }}
            </option>
          </select>
          <p class="text-xs text-gray-400 mt-1">Leave empty for all states</p>
        </div>

        <!-- Area type + Apply -->
        <div class="flex flex-col justify-between">
          <div></div>
          <button @click="applyFilters"
            class="mt-2 px-4 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700 text-sm">
            Apply Filters
          </button>
        </div>
      </div>
    </section>

    <div v-if="salaryStore.loading" class="text-center py-12 text-gray-500">
      Loading salary data...
    </div>

    <div v-else-if="comparisons.length > 0" class="grid gap-8">
      <section class="bg-white p-6 rounded-lg shadow-sm">
        <h2 class="text-lg font-semibold text-gray-800 mb-4">Regional Overview</h2>
        <RegionalMap
          :comparisons="comparisons"
          :user-salary="userSalary"
        />
      </section>

      <section class="bg-white p-6 rounded-lg shadow-sm">
        <h2 class="text-lg font-semibold text-gray-800 mb-4">State Comparison</h2>
        <StateComparisonChart
          :comparisons="comparisons"
          :user-salary="userSalary"
        />
      </section>

      <section class="bg-white p-6 rounded-lg shadow-sm">
        <h2 class="text-lg font-semibold text-gray-800 mb-4">Salary Trends</h2>
        <SalaryTrendsChart
          :trends="salaryStore.trends"
          :user-salary="userSalary"
        />
      </section>

      <section class="bg-white p-6 rounded-lg shadow-sm">
        <h2 class="text-lg font-semibold text-gray-800 mb-4">Raw Data</h2>
        <div class="overflow-x-auto">
          <table class="w-full border-collapse">
            <thead>
              <tr class="bg-gray-50">
                <th class="px-4 py-3 text-left font-semibold border-b border-gray-200">State</th>
                <th class="px-4 py-3 text-left font-semibold border-b border-gray-200">25th Percentile</th>
                <th class="px-4 py-3 text-left font-semibold border-b border-gray-200">Median</th>
                <th class="px-4 py-3 text-left font-semibold border-b border-gray-200">75th Percentile</th>
                <th class="px-4 py-3 text-left font-semibold border-b border-gray-200">vs You</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="row in comparisons" :key="row.state_code"
                class="border-b border-gray-100">
                <td class="px-4 py-3">{{ row.state }}</td>
                <td class="px-4 py-3">${{ row.p25.toLocaleString() }}</td>
                <td class="px-4 py-3">${{ row.median.toLocaleString() }}</td>
                <td class="px-4 py-3">${{ row.p75.toLocaleString() }}</td>
                <td class="px-4 py-3" :class="row.median > userSalary ? 'text-green-600' : 'text-red-600'">
                  {{ row.median > userSalary ? '+' : '' }}${{
                    (row.median - userSalary).toLocaleString()
                  }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </div>

    <div v-else class="text-center py-12 text-gray-500">
      <p class="mb-2">No salary data yet. Go to <router-link to="/scrape" class="text-indigo-600 hover:underline">Data Management</router-link> to trigger a data pull.</p>
      <p>Make sure to set your salary in <router-link to="/settings" class="text-indigo-600 hover:underline">Settings</router-link>.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useSalaryStore } from '@/stores/salary'
import RegionalMap from '@/components/RegionalMap.vue'
import StateComparisonChart from '@/components/StateComparisonChart.vue'
import SalaryTrendsChart from '@/components/SalaryTrendsChart.vue'

const salaryStore = useSalaryStore()

const comparisons = computed(() => salaryStore.compareData?.comparisons ?? [])
const userSalary = computed(() => salaryStore.compareData?.user_salary ?? 0)

// Unique states from regions for the filter dropdown
const stateList = computed(() => {
  const seen = new Set<string>()
  return salaryStore.regions
    .filter(r => r.area_type === 2)
    .filter(r => {
      if (seen.has(r.state_code)) return false
      seen.add(r.state_code)
      return true
    })
    .sort((a, b) => a.state.localeCompare(b.state))
})

onMounted(async () => {
  await Promise.all([
    salaryStore.fetchOccupations(),
    salaryStore.fetchRegions(),
  ])
  await salaryStore.fetchAll()
})


async function applyFilters() {
  await salaryStore.fetchAll()
}

const occSearch = ref('')

const filteredOccupations = computed(() => {
  const q = occSearch.value.toLowerCase()
  if (!q) return salaryStore.occupations
  return salaryStore.occupations.filter(o => o.occ_title.toLowerCase().includes(q))
})

const selectedOccLabels = computed(() =>
  salaryStore.selectedOccCodes.map(code => {
    const occ = salaryStore.occupations.find(o => o.occ_code === code)
    return { code, title: occ?.occ_title ?? code }
  })
)

function removeOcc(code: string) {
  salaryStore.selectedOccCodes = salaryStore.selectedOccCodes.filter(c => c !== code)
}

</script>