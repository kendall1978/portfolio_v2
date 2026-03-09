<template>
  <div class="max-w-4xl mx-auto px-4 py-16">
    <h1 class="text-3xl font-bold mb-10" style="color: var(--text-primary)">Blog</h1>

    <div v-if="!store.blogPosts.length" class="text-center py-12" style="color: var(--text-muted)">
      No posts yet.
    </div>

    <div class="grid gap-8">
      <article
        v-for="post in store.blogPosts"
        :key="post.id"
        class="rounded-xl border overflow-hidden transition-shadow hover:shadow-md"
        :style="{ backgroundColor: 'var(--bg-card)', borderColor: 'var(--border-color)' }"
      >
        <img
          v-if="post.image_path"
          :src="imageUrl(post.image_path)"
          :alt="post.title"
          class="w-full h-48 object-cover"
        />
        <div class="p-6">
          <time class="text-xs uppercase tracking-wide" style="color: var(--text-muted)">
            {{ formatDate(post.date) }}
          </time>
          <h2 class="text-xl font-semibold mt-1 mb-3" style="color: var(--text-primary)">
            {{ post.title }}
          </h2>
          <p class="text-sm leading-relaxed" style="color: var(--text-secondary)">
            {{ post.content }}
          </p>
        </div>
      </article>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { usePortfolioStore } from '@/stores/portfolio'

const store = usePortfolioStore()

onMounted(() => {
    store.fetchBlogPosts()
})

function formatDate(dateStr: string): string {
    return new Date(dateStr).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
    })
}

function imageUrl(path: string): string {
    // If it's already a full URL (e.g., Firebase URL from migration), use as-is
    if (path.startsWith('http')) return path
    // Otherwise it's a local upload path
    return `http://localhost:8080${path}`
}
</script>
