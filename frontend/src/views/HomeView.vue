<template>
  <div class="max-w-4xl mx-auto px-4 py-16">
    <!-- Loading state -->
    <div v-if="store.loading" class="text-center py-24" style="color: var(--text-muted)">
      Loading...
    </div>

    <template v-else>
      <!-- Hero Section -->
      <section class="text-center mb-20">
        <img
          v-if="store.siteSettings.profile_image_url"
          :src="store.siteSettings.profile_image_url"
          :alt="store.siteSettings.site_title || 'Profile'"
          class="w-32 h-32 rounded-full mx-auto mb-6 object-cover shadow-lg"
        />
        <h1 class="text-4xl font-bold mb-3" style="color: var(--text-primary)">
          {{ store.siteSettings.site_title || 'Welcome' }}
        </h1>
        <p class="text-lg" style="color: var(--text-secondary)">
          {{ store.siteSettings.site_subtitle || '' }}
        </p>

        <!-- Social Links -->
        <div v-if="store.socialLinks.length" class="flex justify-center gap-5 mt-6">
          <a
            v-for="link in store.socialLinks"
            :key="link.id"
            :href="link.url"
            target="_blank"
            rel="noopener noreferrer"
            class="text-accent hover:opacity-80 transition-opacity"
            :aria-label="link.platform"
          >
            <SocialIcon :icon="link.icon" class="w-6 h-6" />
          </a>
        </div>
      </section>

      <!-- About Sections -->
      <section v-if="store.homeSections.length" class="grid gap-8 md:grid-cols-3">
        <div
          v-for="section in store.homeSections"
          :key="section.id"
          class="p-6 rounded-xl border transition-shadow hover:shadow-md"
          :style="{ backgroundColor: 'var(--bg-card)', borderColor: 'var(--border-color)' }"
        >
          <h3 class="text-lg font-semibold mb-3" style="color: var(--text-primary)">
            {{ section.title }}
          </h3>
          <p class="text-sm leading-relaxed" style="color: var(--text-secondary)">
            {{ section.content }}
          </p>
        </div>
      </section>
    </template>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { usePortfolioStore } from '@/stores/portfolio'
import SocialIcon from '@/components/SocialIcon.vue'

const store = usePortfolioStore()

onMounted(() => {
    store.fetchHomeData()
})
</script>
