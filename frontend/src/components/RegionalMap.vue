<template>
  <v-chart v-if="mapReady" :option="chartOption" style="height: 500px" autoresize />
  <div v-else class="h-[500px] flex items-center justify-center text-gray-400">Loading map...</div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import VChart from 'vue-echarts'
import { use, registerMap } from 'echarts/core'
import { MapChart } from 'echarts/charts'
import { GeoComponent, TooltipComponent, VisualMapComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import usaGeoJSON from '@/assets/usa.json'

use([MapChart, GeoComponent, TooltipComponent, VisualMapComponent, CanvasRenderer])

const props = defineProps<{
  comparisons: Array<{ state_code: string; state: string; median: number }>
  userSalary: number
}>()

const mapReady = ref(false)

// Register map filtered to only the states in the data
watch(() => props.comparisons, (comps) => {
  if (comps.length === 0) return
  const stateNames = new Set(comps.map(c => c.state))
  const filtered = {
    ...usaGeoJSON,
    features: (usaGeoJSON as any).features.filter((f: any) =>
      stateNames.has(f.properties.name)
    ),
  }
  registerMap('filtered-states', filtered as any)
  mapReady.value = true
}, { immediate: true })

const chartOption = computed(() => ({
  tooltip: {
    trigger: 'item',
    formatter: (params: any) => {
      const comp = props.comparisons.find((c) => c.state === params.name)
      if (!comp) return params.name
      const diff = comp.median - props.userSalary
      const diffStr = diff >= 0 ? `+$${diff.toLocaleString()}` : `-$${Math.abs(diff).toLocaleString()}`
      return `${params.name}<br/>Median: $${comp.median.toLocaleString()}<br/>vs You: ${diffStr}`
    },
  },
  visualMap: {
    min: 40000,
    max: 150000,
    text: ['High', 'Low'],
    realtime: false,
    calculable: true,
    inRange: { color: ['#e0f2fe', '#1d4ed8'] },
  },
  series: [
    {
      name: 'Median Salary',
      type: 'map',
      map: 'filtered-states',
      roam: true,
      data: props.comparisons.map((c) => ({
        name: c.state,
        value: c.median,
      })),
      label: { show: true, fontSize: 10 },
    },
  ],
}))
</script>
