<template>
  <v-chart :option="chartOption" style="width: 100%; height: 100%" autoresize />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

use([LineChart, GridComponent, TooltipComponent, LegendComponent, CanvasRenderer])

const props = defineProps<{
  trends: Array<{
    state_code: string
    state: string
    year: number
    avg_median: number
    avg_p25: number
    avg_p75: number
  }>
  userSalary: number
}>()

const chartOption = computed(() => {
  const states = [...new Set(props.trends.map(t => t.state))]
  const years = [...new Set(props.trends.map(t => t.year))].sort()

  const series: any[] = states.map(state => ({
    name: state,
    type: 'line',
    smooth: true,
    data: years.map(year => {
      const point = props.trends.find(t => t.state === state && t.year === year)
      return point?.avg_median ?? null
    }),
  }))

  // Add "Your Salary" reference line
  series.push({
    name: 'Your Salary',
    type: 'line',
    smooth: false,
    data: years.map(() => props.userSalary),
    lineStyle: { type: 'dashed', color: '#dc2626', width: 2 },
    itemStyle: { color: '#dc2626' },
    symbol: 'none',
  })

  return {
    tooltip: { trigger: 'axis' },
    legend: { type: 'scroll', bottom: 0 },
    grid: { left: '3%', right: '4%', bottom: '10%', containLabel: true },
    xAxis: { type: 'category', data: years.map(String) },
    yAxis: { type: 'value', axisLabel: { formatter: '${value}' } },
    series,
  }
})
</script>
