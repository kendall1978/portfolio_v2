<template>
  <div class="relative" ref="container">
    <!-- Trigger button -->
    <button
      type="button"
      @click="open = !open"
      class="w-full flex items-center justify-between px-3 py-2 border rounded text-sm text-left transition-colors"
      :style="{
        backgroundColor: 'var(--bg-primary)',
        borderColor: open ? 'var(--accent-color)' : 'var(--border-color)',
        color: modelValue.length ? 'var(--text-primary)' : 'var(--text-muted)',
      }"
    >
      <span class="truncate">{{ triggerLabel }}</span>
      <svg class="w-4 h-4 shrink-0 ml-2 transition-transform" :class="{ 'rotate-180': open }" viewBox="0 0 20 20" fill="currentColor" style="color: var(--text-muted)">
        <path fill-rule="evenodd" d="M5.23 7.21a.75.75 0 011.06.02L10 11.168l3.71-3.938a.75.75 0 111.08 1.04l-4.25 4.5a.75.75 0 01-1.08 0l-4.25-4.5a.75.75 0 01.02-1.06z" clip-rule="evenodd" />
      </svg>
    </button>

    <!-- Dropdown panel -->
    <div
      v-if="open"
      class="absolute z-20 mt-1 w-full rounded-lg border shadow-lg overflow-hidden"
      :style="{ backgroundColor: 'var(--bg-card)', borderColor: 'var(--border-color)' }"
    >
      <!-- Search input -->
      <div v-if="searchable" class="p-2 border-b" :style="{ borderColor: 'var(--border-color)' }">
        <input
          ref="searchInput"
          v-model="search"
          type="text"
          placeholder="Search..."
          class="w-full px-2 py-1.5 border rounded text-sm"
          @keydown.escape="open = false"
        />
      </div>

      <!-- Options list -->
      <ul class="max-h-60 overflow-y-auto py-1">
        <li
          v-for="opt in filteredOptions"
          :key="opt.value"
          @click="toggle(opt.value)"
          class="flex items-center gap-2 px-3 py-2 text-sm cursor-pointer transition-colors"
          :style="{ color: 'var(--text-primary)' }"
          @mouseenter="($event.currentTarget as HTMLElement).style.backgroundColor = 'var(--bg-secondary)'"
          @mouseleave="($event.currentTarget as HTMLElement).style.backgroundColor = 'transparent'"
        >
          <span
            class="w-4 h-4 shrink-0 rounded border flex items-center justify-center transition-colors"
            :style="{
              borderColor: isSelected(opt.value) ? 'var(--accent-color)' : 'var(--border-color)',
              backgroundColor: isSelected(opt.value) ? 'var(--accent-color)' : 'transparent',
            }"
          >
            <svg v-if="isSelected(opt.value)" class="w-3 h-3 text-white" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
            </svg>
          </span>
          <span class="truncate">{{ opt.label }}</span>
        </li>
        <li v-if="filteredOptions.length === 0" class="px-3 py-2 text-sm" style="color: var(--text-muted)">
          No results found
        </li>
      </ul>
    </div>

    <!-- Selected chips -->
    <div v-if="selectedChips.length" class="flex flex-wrap gap-1 mt-1.5">
      <span
        v-for="chip in selectedChips"
        :key="chip.value"
        class="inline-flex items-center gap-1 px-2 py-0.5 rounded text-xs bg-accent-light text-accent"
      >
        <span class="truncate max-w-[150px]">{{ chip.label }}</span>
        <button type="button" @click.stop="toggle(chip.value)" class="hover:opacity-70">&times;</button>
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick, onMounted, onBeforeUnmount } from 'vue'

export interface SelectOption {
  value: string
  label: string
}

const props = defineProps<{
  options: SelectOption[]
  modelValue: string[]
  placeholder?: string
  searchable?: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string[]]
}>()

const container = ref<HTMLElement | null>(null)
const searchInput = ref<HTMLInputElement | null>(null)
const open = ref(false)
const search = ref('')

const triggerLabel = computed(() => {
  if (props.modelValue.length === 0) return props.placeholder ?? 'Select...'
  if (props.modelValue.length === 1) {
    const opt = props.options.find(o => o.value === props.modelValue[0])
    return opt?.label ?? props.modelValue[0]
  }
  return `${props.modelValue.length} selected`
})

const filteredOptions = computed(() => {
  if (!search.value) return props.options
  const q = search.value.toLowerCase()
  return props.options.filter(o => o.label.toLowerCase().includes(q))
})

const selectedChips = computed(() =>
  props.modelValue.map(v => {
    const opt = props.options.find(o => o.value === v)
    return { value: v, label: opt?.label ?? v }
  })
)

function isSelected(value: string): boolean {
  return props.modelValue.includes(value)
}

function toggle(value: string) {
  const next = isSelected(value)
    ? props.modelValue.filter(v => v !== value)
    : [...props.modelValue, value]
  emit('update:modelValue', next)
}

// Focus search input when dropdown opens
watch(open, (isOpen) => {
  if (isOpen && props.searchable) {
    nextTick(() => searchInput.value?.focus())
  }
  if (!isOpen) {
    search.value = ''
  }
})

// Close on click outside
function onClickOutside(e: MouseEvent) {
  if (container.value && !container.value.contains(e.target as Node)) {
    open.value = false
  }
}

// Close on Escape
function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') open.value = false
}

onMounted(() => {
  document.addEventListener('mousedown', onClickOutside)
  document.addEventListener('keydown', onKeydown)
})

onBeforeUnmount(() => {
  document.removeEventListener('mousedown', onClickOutside)
  document.removeEventListener('keydown', onKeydown)
})
</script>
