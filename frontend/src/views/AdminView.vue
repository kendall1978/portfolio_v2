<template>
  <div class="max-w-4xl mx-auto px-4 py-8">
    <h1 class="text-2xl font-bold mb-6" style="color: var(--text-primary)">Admin</h1>

    <!-- Tabs -->
    <div class="flex gap-1 mb-8 border-b" style="border-color: var(--border-color)">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        @click="activeTab = tab.id"
        class="px-4 py-2 text-sm font-medium -mb-px border-b-2 transition-colors"
        :style="{
          borderColor: activeTab === tab.id ? 'var(--accent-color)' : 'transparent',
          color: activeTab === tab.id ? 'var(--accent-color)' : 'var(--text-muted)',
        }"
      >
        {{ tab.label }}
      </button>
    </div>

    <!-- Blog Posts Tab -->
    <div v-if="activeTab === 'blog'">
      <div class="flex justify-between items-center mb-4">
        <h2 class="text-lg font-semibold" style="color: var(--text-primary)">Blog Posts</h2>
        <button @click="showBlogForm = !showBlogForm" class="px-3 py-1.5 rounded text-sm text-white bg-accent hover:bg-accent-hover transition-colors">
          {{ showBlogForm ? 'Cancel' : 'New Post' }}
        </button>
      </div>

      <!-- Blog form -->
      <form v-if="showBlogForm" @submit.prevent="saveBlogPost" class="p-4 rounded-lg border mb-6" :style="{ backgroundColor: 'var(--bg-card)', borderColor: 'var(--border-color)' }">
        <div class="grid gap-4">
          <input v-model="blogForm.title" placeholder="Title" required class="w-full px-3 py-2 rounded border text-sm" :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)', color: 'var(--text-primary)' }" />
          <textarea v-model="blogForm.content" placeholder="Content" rows="4" required class="w-full px-3 py-2 rounded border text-sm" :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)', color: 'var(--text-primary)' }"></textarea>
          <input v-model="blogForm.date" type="date" required class="w-full px-3 py-2 rounded border text-sm" :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)', color: 'var(--text-primary)' }" />
          <div>
            <input type="file" accept="image/jpeg,image/png,image/webp,image/gif" @change="onBlogImageSelect" class="text-sm" style="color: var(--text-secondary)" />
          </div>
          <button type="submit" class="px-4 py-2 rounded text-sm text-white bg-accent hover:bg-accent-hover transition-colors w-fit">
            {{ editingBlogId ? 'Update' : 'Create' }}
          </button>
        </div>
      </form>

      <!-- Blog list -->
      <div class="grid gap-3">
        <div
          v-for="post in store.blogPosts"
          :key="post.id"
          class="flex items-center justify-between p-3 rounded-lg border"
          :style="{ backgroundColor: 'var(--bg-card)', borderColor: 'var(--border-color)' }"
        >
          <div>
            <p class="font-medium text-sm" style="color: var(--text-primary)">{{ post.title }}</p>
            <p class="text-xs" style="color: var(--text-muted)">{{ formatDate(post.date) }}</p>
          </div>
          <div class="flex gap-2">
            <button @click="editBlogPost(post)" class="text-xs px-2 py-1 rounded" style="color: var(--text-secondary)">Edit</button>
            <button @click="deleteBlogPost(post.id)" class="text-xs px-2 py-1 rounded text-red-500 hover:text-red-400">Delete</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Home Sections Tab -->
    <div v-if="activeTab === 'sections'">
      <div class="flex justify-between items-center mb-4">
        <h2 class="text-lg font-semibold" style="color: var(--text-primary)">Home Sections</h2>
        <button @click="showSectionForm = !showSectionForm" class="px-3 py-1.5 rounded text-sm text-white bg-accent hover:bg-accent-hover transition-colors">
          {{ showSectionForm ? 'Cancel' : 'New Section' }}
        </button>
      </div>

      <form v-if="showSectionForm" @submit.prevent="saveSection" class="p-4 rounded-lg border mb-6" :style="{ backgroundColor: 'var(--bg-card)', borderColor: 'var(--border-color)' }">
        <div class="grid gap-4">
          <input v-model="sectionForm.title" placeholder="Title" required class="w-full px-3 py-2 rounded border text-sm" :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)', color: 'var(--text-primary)' }" />
          <textarea v-model="sectionForm.content" placeholder="Content" rows="3" required class="w-full px-3 py-2 rounded border text-sm" :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)', color: 'var(--text-primary)' }"></textarea>
          <input v-model.number="sectionForm.sort_order" type="number" placeholder="Sort order" class="w-full px-3 py-2 rounded border text-sm" :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)', color: 'var(--text-primary)' }" />
          <button type="submit" class="px-4 py-2 rounded text-sm text-white bg-accent hover:bg-accent-hover transition-colors w-fit">
            {{ editingSectionId ? 'Update' : 'Create' }}
          </button>
        </div>
      </form>

      <div class="grid gap-3">
        <div
          v-for="section in store.homeSections"
          :key="section.id"
          class="flex items-center justify-between p-3 rounded-lg border"
          :style="{ backgroundColor: 'var(--bg-card)', borderColor: 'var(--border-color)' }"
        >
          <div>
            <p class="font-medium text-sm" style="color: var(--text-primary)">{{ section.title }}</p>
            <p class="text-xs" style="color: var(--text-muted)">Order: {{ section.sort_order }}</p>
          </div>
          <div class="flex gap-2">
            <button @click="editSection(section)" class="text-xs px-2 py-1 rounded" style="color: var(--text-secondary)">Edit</button>
            <button @click="deleteSection(section.id)" class="text-xs px-2 py-1 rounded text-red-500 hover:text-red-400">Delete</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Social Links Tab -->
    <div v-if="activeTab === 'links'">
      <div class="flex justify-between items-center mb-4">
        <h2 class="text-lg font-semibold" style="color: var(--text-primary)">Social Links</h2>
        <button @click="showLinkForm = !showLinkForm" class="px-3 py-1.5 rounded text-sm text-white bg-accent hover:bg-accent-hover transition-colors">
          {{ showLinkForm ? 'Cancel' : 'New Link' }}
        </button>
      </div>

      <form v-if="showLinkForm" @submit.prevent="saveLink" class="p-4 rounded-lg border mb-6" :style="{ backgroundColor: 'var(--bg-card)', borderColor: 'var(--border-color)' }">
        <div class="grid gap-4">
          <input v-model="linkForm.platform" placeholder="Platform (e.g. GitHub)" required class="w-full px-3 py-2 rounded border text-sm" :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)', color: 'var(--text-primary)' }" />
          <input v-model="linkForm.url" placeholder="URL" required class="w-full px-3 py-2 rounded border text-sm" :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)', color: 'var(--text-primary)' }" />
          <input v-model="linkForm.icon" placeholder="Icon key (github, linkedin, twitter, instagram)" required class="w-full px-3 py-2 rounded border text-sm" :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)', color: 'var(--text-primary)' }" />
          <input v-model.number="linkForm.sort_order" type="number" placeholder="Sort order" class="w-full px-3 py-2 rounded border text-sm" :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)', color: 'var(--text-primary)' }" />
          <button type="submit" class="px-4 py-2 rounded text-sm text-white bg-accent hover:bg-accent-hover transition-colors w-fit">
            {{ editingLinkId ? 'Update' : 'Create' }}
          </button>
        </div>
      </form>

      <div class="grid gap-3">
        <div
          v-for="link in store.socialLinks"
          :key="link.id"
          class="flex items-center justify-between p-3 rounded-lg border"
          :style="{ backgroundColor: 'var(--bg-card)', borderColor: 'var(--border-color)' }"
        >
          <div>
            <p class="font-medium text-sm" style="color: var(--text-primary)">{{ link.platform }}</p>
            <p class="text-xs" style="color: var(--text-muted)">{{ link.url }}</p>
          </div>
          <div class="flex gap-2">
            <button @click="editLink(link)" class="text-xs px-2 py-1 rounded" style="color: var(--text-secondary)">Edit</button>
            <button @click="deleteLink(link.id)" class="text-xs px-2 py-1 rounded text-red-500 hover:text-red-400">Delete</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Site Settings Tab -->
    <div v-if="activeTab === 'settings'">
      <h2 class="text-lg font-semibold mb-4" style="color: var(--text-primary)">Site Settings</h2>
      <form @submit.prevent="saveSiteSettings" class="p-4 rounded-lg border" :style="{ backgroundColor: 'var(--bg-card)', borderColor: 'var(--border-color)' }">
        <div class="grid gap-4">
          <div>
            <label class="block text-xs font-medium mb-1" style="color: var(--text-muted)">Site Title</label>
            <input v-model="settingsForm.site_title" class="w-full px-3 py-2 rounded border text-sm" :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)', color: 'var(--text-primary)' }" />
          </div>
          <div>
            <label class="block text-xs font-medium mb-1" style="color: var(--text-muted)">Subtitle</label>
            <input v-model="settingsForm.site_subtitle" class="w-full px-3 py-2 rounded border text-sm" :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)', color: 'var(--text-primary)' }" />
          </div>
          <div>
            <label class="block text-xs font-medium mb-1" style="color: var(--text-muted)">Profile Image URL</label>
            <input v-model="settingsForm.profile_image_url" class="w-full px-3 py-2 rounded border text-sm" :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)', color: 'var(--text-primary)' }" />
          </div>
          <button type="submit" class="px-4 py-2 rounded text-sm text-white bg-accent hover:bg-accent-hover transition-colors w-fit">
            Save Settings
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { usePortfolioStore, type BlogPost, type HomeSection, type SocialLink } from '@/stores/portfolio'

const store = usePortfolioStore()

const tabs = [
    { id: 'blog', label: 'Blog Posts' },
    { id: 'sections', label: 'Home Sections' },
    { id: 'links', label: 'Social Links' },
    { id: 'settings', label: 'Site Settings' },
]
const activeTab = ref('blog')

// Blog
const showBlogForm = ref(false)
const editingBlogId = ref<number | null>(null)
const blogForm = reactive({ title: '', content: '', date: '', image_path: '' })
let blogImageFile: File | null = null

function editBlogPost(post: BlogPost) {
    editingBlogId.value = post.id
    blogForm.title = post.title
    blogForm.content = post.content
    blogForm.date = post.date.split('T')[0]
    blogForm.image_path = post.image_path
    showBlogForm.value = true
}

function onBlogImageSelect(e: Event) {
    const input = e.target as HTMLInputElement
    if (input.files?.length) blogImageFile = input.files[0]
}

async function saveBlogPost() {
    if (blogImageFile) {
        blogForm.image_path = await store.uploadImage(blogImageFile)
        blogImageFile = null
    }
    if (editingBlogId.value) {
        await store.updateBlogPost(editingBlogId.value, { ...blogForm })
    } else {
        await store.createBlogPost({ ...blogForm })
    }
    resetBlogForm()
}

async function deleteBlogPost(id: number) {
    await store.deleteBlogPost(id)
}

function resetBlogForm() {
    showBlogForm.value = false
    editingBlogId.value = null
    blogForm.title = ''
    blogForm.content = ''
    blogForm.date = ''
    blogForm.image_path = ''
}

// Sections
const showSectionForm = ref(false)
const editingSectionId = ref<number | null>(null)
const sectionForm = reactive({ title: '', content: '', sort_order: 0 })

function editSection(section: HomeSection) {
    editingSectionId.value = section.id
    sectionForm.title = section.title
    sectionForm.content = section.content
    sectionForm.sort_order = section.sort_order
    showSectionForm.value = true
}

async function saveSection() {
    if (editingSectionId.value) {
        await store.updateHomeSection(editingSectionId.value, { ...sectionForm })
    } else {
        await store.createHomeSection({ ...sectionForm })
    }
    showSectionForm.value = false
    editingSectionId.value = null
    sectionForm.title = ''
    sectionForm.content = ''
    sectionForm.sort_order = 0
}

async function deleteSection(id: number) {
    await store.deleteHomeSection(id)
}

// Links
const showLinkForm = ref(false)
const editingLinkId = ref<number | null>(null)
const linkForm = reactive({ platform: '', url: '', icon: '', sort_order: 0 })

function editLink(link: SocialLink) {
    editingLinkId.value = link.id
    linkForm.platform = link.platform
    linkForm.url = link.url
    linkForm.icon = link.icon
    linkForm.sort_order = link.sort_order
    showLinkForm.value = true
}

async function saveLink() {
    if (editingLinkId.value) {
        await store.updateSocialLink(editingLinkId.value, { ...linkForm })
    } else {
        await store.createSocialLink({ ...linkForm })
    }
    showLinkForm.value = false
    editingLinkId.value = null
    linkForm.platform = ''
    linkForm.url = ''
    linkForm.icon = ''
    linkForm.sort_order = 0
}

async function deleteLink(id: number) {
    await store.deleteSocialLink(id)
}

// Site Settings
const settingsForm = reactive({ site_title: '', site_subtitle: '', profile_image_url: '' })

async function saveSiteSettings() {
    await store.updateSiteSettings({ ...settingsForm })
}

function formatDate(dateStr: string): string {
    return new Date(dateStr).toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })
}

onMounted(async () => {
    await Promise.all([
        store.fetchBlogPosts(),
        store.fetchHomeSections(),
        store.fetchSocialLinks(),
        store.fetchSiteSettings(),
    ])
    settingsForm.site_title = store.siteSettings.site_title || ''
    settingsForm.site_subtitle = store.siteSettings.site_subtitle || ''
    settingsForm.profile_image_url = store.siteSettings.profile_image_url || ''
})
</script>
