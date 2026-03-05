import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/api/client'

interface Region {
  id: number
  area_code: string
  area_title: string
  area_type: number
  state: string
  state_code: string
}

interface Occupation {
  occ_code: string
  occ_title: string
}

interface StateComparison {
  state_code: string
  state: string
  median: number
  p25: number
  p75: number
}

interface CompareData {
  user_salary: number
  comparisons: StateComparison[]
}

interface TrendPoint {
  state_code: string
  state: string
  year: number
  avg_median: number
  avg_p25: number
  avg_p75: number
}

export const useSalaryStore = defineStore('salary', () => {
  const regions = ref<Region[]>([])
  const occupations = ref<Occupation[]>([])
  const compareData = ref<CompareData | null>(null)
  const trends = ref<TrendPoint[]>([])
  const loading = ref(false)

  // Filter state
  const selectedOccCodes = ref<string[]>(['15-1254'])
  const selectedStates = ref<string[]>([])
  const areaType = ref<'state' | 'metro'>('state')

  function buildFilterParams() {
    const params: Record<string, string> = {}
    if (selectedOccCodes.value.length > 0) {
      params.occ_codes = selectedOccCodes.value.join(',')
    }
    if (selectedStates.value.length > 0) {
      params.states = selectedStates.value.join(',')
    }
    params.area_type = areaType.value
    return params
  }

  async function fetchRegions() {
    const { data } = await api.get('/regions')
    regions.value = data
  }

  async function fetchOccupations() {
    const { data } = await api.get('/occupations')
    occupations.value = data
  }

  async function fetchCompare() {
    loading.value = true
    try {
      const { data } = await api.get('/salaries/compare', { params: buildFilterParams() })
      compareData.value = data
    } finally {
      loading.value = false
    }
  }

  async function fetchTrends() {
    const { data } = await api.get('/salaries/trends', { params: buildFilterParams() })
    trends.value = data
  }

  async function fetchAll() {
    await Promise.all([fetchCompare(), fetchTrends()])
  }

  return {
    regions, occupations, compareData, trends, loading,
    selectedOccCodes, selectedStates, areaType,
    fetchRegions, fetchOccupations, fetchCompare, fetchTrends, fetchAll,
  }
})