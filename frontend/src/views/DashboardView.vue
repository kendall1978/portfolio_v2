<template>
  <div class="max-w-6xl mx-auto px-4 py-8">
    <!-- Sticky Filter Bar -->
    <div
      class="sticky top-0 z-10 -mx-4 px-4 py-3 mb-6 border-b"
      :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)' }"
    >
      <div class="flex items-end gap-4">
        <div class="flex-1 min-w-0">
          <MultiSelect
            v-model="salaryStore.selectedOccCodes"
            :options="occupationOptions"
            placeholder="Select occupations..."
            :searchable="true"
          />
        </div>
        <div class="flex-1 min-w-0">
          <MultiSelect
            v-model="salaryStore.selectedStates"
            :options="stateOptions"
            placeholder="Select states..."
            :searchable="false"
          />
        </div>
        <button
          @click="applyFilters"
          class="px-4 py-2 rounded text-sm text-white bg-accent hover:bg-accent-hover transition-colors shrink-0"
        >
          Apply Filters
        </button>
      </div>
    </div>

    <div v-if="salaryStore.loading" class="text-center py-12" style="color: var(--text-muted)">
      Loading salary data...
    </div>

    <div v-else-if="comparisons.length > 0" class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <!-- Regional Map Widget -->
      <div class="rounded-xl border overflow-hidden" :style="{ backgroundColor: 'var(--bg-card)', borderColor: 'var(--border-color)' }">
        <div class="flex items-center justify-between px-4 py-3 border-b" :style="{ borderColor: 'var(--border-color)' }">
          <h2 class="text-sm font-semibold" style="color: var(--text-primary)">Regional Overview</h2>
          <button @click="expandedWidget = 'map'" class="p-1 rounded transition-colors" style="color: var(--text-muted)" title="Expand">
            <svg class="w-4 h-4" viewBox="0 0 20 20" fill="currentColor"><path d="M13.28 7.78l3.22-3.22v2.69a.75.75 0 001.5 0v-4.5a.75.75 0 00-.75-.75h-4.5a.75.75 0 000 1.5h2.69l-3.22 3.22a.75.75 0 001.06 1.06zM2 17.25v-4.5a.75.75 0 011.5 0v2.69l3.22-3.22a.75.75 0 011.06 1.06L4.56 16.5h2.69a.75.75 0 010 1.5h-4.5a.75.75 0 01-.75-.75z" /></svg>
          </button>
        </div>
        <div class="p-4">
          <RegionalMap :comparisons="comparisons" :user-salary="userSalary" style="height: 280px" />
        </div>
      </div>

      <!-- State Comparison Widget -->
      <div class="rounded-xl border overflow-hidden" :style="{ backgroundColor: 'var(--bg-card)', borderColor: 'var(--border-color)' }">
        <div class="flex items-center justify-between px-4 py-3 border-b" :style="{ borderColor: 'var(--border-color)' }">
          <h2 class="text-sm font-semibold" style="color: var(--text-primary)">State Comparison</h2>
          <button @click="expandedWidget = 'comparison'" class="p-1 rounded transition-colors" style="color: var(--text-muted)" title="Expand">
            <svg class="w-4 h-4" viewBox="0 0 20 20" fill="currentColor"><path d="M13.28 7.78l3.22-3.22v2.69a.75.75 0 001.5 0v-4.5a.75.75 0 00-.75-.75h-4.5a.75.75 0 000 1.5h2.69l-3.22 3.22a.75.75 0 001.06 1.06zM2 17.25v-4.5a.75.75 0 011.5 0v2.69l3.22-3.22a.75.75 0 011.06 1.06L4.56 16.5h2.69a.75.75 0 010 1.5h-4.5a.75.75 0 01-.75-.75z" /></svg>
          </button>
        </div>
        <div class="p-4">
          <StateComparisonChart :comparisons="comparisons" :user-salary="userSalary" style="height: 280px" />
        </div>
      </div>

      <!-- Salary Trends Widget (full width) -->
      <div class="md:col-span-2 rounded-xl border overflow-hidden" :style="{ backgroundColor: 'var(--bg-card)', borderColor: 'var(--border-color)' }">
        <div class="flex items-center justify-between px-4 py-3 border-b" :style="{ borderColor: 'var(--border-color)' }">
          <h2 class="text-sm font-semibold" style="color: var(--text-primary)">Salary Trends</h2>
          <button @click="expandedWidget = 'trends'" class="p-1 rounded transition-colors" style="color: var(--text-muted)" title="Expand">
            <svg class="w-4 h-4" viewBox="0 0 20 20" fill="currentColor"><path d="M13.28 7.78l3.22-3.22v2.69a.75.75 0 001.5 0v-4.5a.75.75 0 00-.75-.75h-4.5a.75.75 0 000 1.5h2.69l-3.22 3.22a.75.75 0 001.06 1.06zM2 17.25v-4.5a.75.75 0 011.5 0v2.69l3.22-3.22a.75.75 0 011.06 1.06L4.56 16.5h2.69a.75.75 0 010 1.5h-4.5a.75.75 0 01-.75-.75z" /></svg>
          </button>
        </div>
        <div class="p-4">
          <SalaryTrendsChart :trends="salaryStore.trends" :user-salary="userSalary" style="height: 280px" />
        </div>
      </div>

      <!-- Raw Data Widget (full width) -->
      <div class="md:col-span-2 rounded-xl border overflow-hidden" :style="{ backgroundColor: 'var(--bg-card)', borderColor: 'var(--border-color)' }">
        <div class="flex items-center justify-between px-4 py-3 border-b" :style="{ borderColor: 'var(--border-color)' }">
          <h2 class="text-sm font-semibold" style="color: var(--text-primary)">Raw Data</h2>
          <button v-if="comparisons.length > 5" @click="expandedWidget = 'data'" class="text-xs text-accent hover:opacity-80 transition-opacity">
            Show all {{ comparisons.length }} rows
          </button>
        </div>
        <div class="overflow-x-auto">
          <table class="w-full border-collapse">
            <thead>
              <tr :style="{ backgroundColor: 'var(--bg-secondary)' }">
                <th class="px-4 py-3 text-left font-semibold text-sm" style="color: var(--text-secondary)">State</th>
                <th class="px-4 py-3 text-left font-semibold text-sm" style="color: var(--text-secondary)">25th Pctl</th>
                <th class="px-4 py-3 text-left font-semibold text-sm" style="color: var(--text-secondary)">Median</th>
                <th class="px-4 py-3 text-left font-semibold text-sm" style="color: var(--text-secondary)">75th Pctl</th>
                <th class="px-4 py-3 text-left font-semibold text-sm" style="color: var(--text-secondary)">vs You</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="row in previewRows"
                :key="row.state_code"
                class="border-t transition-colors"
                :style="{ borderColor: 'var(--border-color)' }"
                @mouseenter="($event.currentTarget as HTMLElement).style.backgroundColor = 'var(--bg-secondary)'"
                @mouseleave="($event.currentTarget as HTMLElement).style.backgroundColor = 'transparent'"
              >
                <td class="px-4 py-3 text-sm" style="color: var(--text-primary)">{{ row.state }}</td>
                <td class="px-4 py-3 text-sm" style="color: var(--text-secondary)">${{ row.p25.toLocaleString() }}</td>
                <td class="px-4 py-3 text-sm" style="color: var(--text-primary)">${{ row.median.toLocaleString() }}</td>
                <td class="px-4 py-3 text-sm" style="color: var(--text-secondary)">${{ row.p75.toLocaleString() }}</td>
                <td class="px-4 py-3 text-sm" :class="row.median > userSalary ? 'text-green-500' : 'text-red-500'">
                  {{ row.median > userSalary ? '+' : '' }}${{ (row.median - userSalary).toLocaleString() }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <div v-else class="text-center py-12" style="color: var(--text-muted)">
      <p class="mb-2">No salary data yet. Go to <router-link to="/scrape" class="text-accent hover:opacity-80">Data Management</router-link> to trigger a data pull.</p>
      <p>Make sure to set your salary in <router-link to="/settings" class="text-accent hover:opacity-80">Settings</router-link>.</p>
    </div>

    <!-- Expand Modals -->
    <WidgetModal :open="expandedWidget === 'map'" title="Regional Overview" @close="expandedWidget = null">
      <RegionalMap :comparisons="comparisons" :user-salary="userSalary" style="height: 70vh" />
    </WidgetModal>

    <WidgetModal :open="expandedWidget === 'comparison'" title="State Comparison" @close="expandedWidget = null">
      <StateComparisonChart :comparisons="comparisons" :user-salary="userSalary" style="height: 70vh" />
    </WidgetModal>

    <WidgetModal :open="expandedWidget === 'trends'" title="Salary Trends" @close="expandedWidget = null">
      <SalaryTrendsChart :trends="salaryStore.trends" :user-salary="userSalary" style="height: 70vh" />
    </WidgetModal>

    <WidgetModal :open="expandedWidget === 'data'" title="Raw Data" @close="expandedWidget = null">
      <div class="overflow-x-auto">
        <table class="w-full border-collapse">
          <thead>
            <tr :style="{ backgroundColor: 'var(--bg-secondary)' }">
              <th class="px-4 py-3 text-left font-semibold text-sm" style="color: var(--text-secondary)">State</th>
              <th class="px-4 py-3 text-left font-semibold text-sm" style="color: var(--text-secondary)">25th Pctl</th>
              <th class="px-4 py-3 text-left font-semibold text-sm" style="color: var(--text-secondary)">Median</th>
              <th class="px-4 py-3 text-left font-semibold text-sm" style="color: var(--text-secondary)">75th Pctl</th>
              <th class="px-4 py-3 text-left font-semibold text-sm" style="color: var(--text-secondary)">vs You</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="row in comparisons"
              :key="row.state_code"
              class="border-t transition-colors"
              :style="{ borderColor: 'var(--border-color)' }"
              @mouseenter="($event.currentTarget as HTMLElement).style.backgroundColor = 'var(--bg-secondary)'"
              @mouseleave="($event.currentTarget as HTMLElement).style.backgroundColor = 'transparent'"
            >
              <td class="px-4 py-3 text-sm" style="color: var(--text-primary)">{{ row.state }}</td>
              <td class="px-4 py-3 text-sm" style="color: var(--text-secondary)">${{ row.p25.toLocaleString() }}</td>
              <td class="px-4 py-3 text-sm" style="color: var(--text-primary)">${{ row.median.toLocaleString() }}</td>
              <td class="px-4 py-3 text-sm" style="color: var(--text-secondary)">${{ row.p75.toLocaleString() }}</td>
              <td class="px-4 py-3 text-sm" :class="row.median > userSalary ? 'text-green-500' : 'text-red-500'">
                {{ row.median > userSalary ? '+' : '' }}${{ (row.median - userSalary).toLocaleString() }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </WidgetModal>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useSalaryStore } from '@/stores/salary'
import RegionalMap from '@/components/RegionalMap.vue'
import StateComparisonChart from '@/components/StateComparisonChart.vue'
import SalaryTrendsChart from '@/components/SalaryTrendsChart.vue'
import MultiSelect from '@/components/MultiSelect.vue'
import WidgetModal from '@/components/WidgetModal.vue'

const salaryStore = useSalaryStore()

const comparisons = computed(() => salaryStore.compareData?.comparisons ?? [])
const userSalary = computed(() => salaryStore.compareData?.user_salary ?? 0)
const previewRows = computed(() => comparisons.value.slice(0, 5))

const expandedWidget = ref<string | null>(null)

// Build options arrays for MultiSelect
const occupationOptions = computed(() =>
  salaryStore.occupations.map(o => ({ value: o.occ_code, label: o.occ_title }))
)

const stateOptions = computed(() => {
  const seen = new Set<string>()
  return salaryStore.regions
    .filter(r => r.area_type === 2)
    .filter(r => {
      if (seen.has(r.state_code)) return false
      seen.add(r.state_code)
      return true
    })
    .sort((a, b) => a.state.localeCompare(b.state))
    .map(r => ({ value: r.state_code, label: r.state }))
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
</script>
