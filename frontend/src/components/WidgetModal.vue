<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="open" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-black/50" @click="$emit('close')" />

        <!-- Modal content -->
        <div
          class="relative w-full rounded-xl border shadow-2xl flex flex-col"
          style="max-width: 90vw; max-height: 80vh;"
          :style="{ backgroundColor: 'var(--bg-card)', borderColor: 'var(--border-color)' }"
        >
          <!-- Header -->
          <div class="flex items-center justify-between px-6 py-4 border-b" :style="{ borderColor: 'var(--border-color)' }">
            <h2 class="text-lg font-semibold" style="color: var(--text-primary)">{{ title }}</h2>
            <button
              @click="$emit('close')"
              class="p-1 rounded-lg transition-colors"
              :style="{ color: 'var(--text-muted)' }"
              @mouseenter="($event.currentTarget as HTMLElement).style.backgroundColor = 'var(--bg-secondary)'"
              @mouseleave="($event.currentTarget as HTMLElement).style.backgroundColor = 'transparent'"
            >
              <svg class="w-5 h-5" viewBox="0 0 20 20" fill="currentColor">
                <path d="M6.28 5.22a.75.75 0 00-1.06 1.06L8.94 10l-3.72 3.72a.75.75 0 101.06 1.06L10 11.06l3.72 3.72a.75.75 0 101.06-1.06L11.06 10l3.72-3.72a.75.75 0 00-1.06-1.06L10 8.94 6.28 5.22z" />
              </svg>
            </button>
          </div>

          <!-- Body -->
          <div class="flex-1 overflow-auto p-6">
            <slot />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { watch, onMounted, onBeforeUnmount } from 'vue'

const props = defineProps<{
  open: boolean
  title: string
}>()

const emit = defineEmits<{
  close: []
}>()

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape' && props.open) {
    emit('close')
  }
}

// Prevent body scroll when modal is open
watch(() => props.open, (isOpen) => {
  document.body.style.overflow = isOpen ? 'hidden' : ''
})

onMounted(() => {
  document.addEventListener('keydown', onKeydown)
})

onBeforeUnmount(() => {
  document.removeEventListener('keydown', onKeydown)
  document.body.style.overflow = ''
})
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}
.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
</style>
