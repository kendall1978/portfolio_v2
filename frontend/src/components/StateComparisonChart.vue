<template>
  <v-chart :option="chartOption" style="height: 400px" autoresize />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { BarChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, MarkLineComponent, LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

use([BarChart, GridComponent, TooltipComponent, MarkLineComponent, LegendComponent, CanvasRenderer])

const props = defineProps<{
  comparisons: Array<{ state_code: string; state: string; median: number; p25: number; p75: number }>
  userSalary: number
}>()

const chartOption = computed(() => ({
  tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
  legend: { data: ['25th Percentile', 'Median', '75th Percentile'] },
  grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
  yAxis: {
    type: 'category',
    data: props.comparisons.map((c) => c.state),
  },
  xAxis: {
    type: 'value',
    axisLabel: { formatter: '${value}' },
  },
  series: [
    {
      name: '25th Percentile',
      type: 'bar',
      data: props.comparisons.map((c) => c.p25),
      itemStyle: { color: '#93c5fd' },
      markLine: {
        data: [{ xAxis: props.userSalary, name: 'You' }],
        label: { formatter: 'You: ${c}' },
        lineStyle: { color: '#dc2626', type: 'dashed', width: 2 },
      },
    },
    {
      name: 'Median',
      type: 'bar',
      data: props.comparisons.map((c) => c.median),
      itemStyle: { color: '#3b82f6' },
    },
    {
      name: '75th Percentile',
      type: 'bar',
      data: props.comparisons.map((c) => c.p75),
      itemStyle: { color: '#1d4ed8' },
    },
  ],
}))
</script>