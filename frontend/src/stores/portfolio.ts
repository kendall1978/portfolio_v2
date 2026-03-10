import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/api/client'

export interface BlogPost {
    id: number
    title: string
    content: string
    image_path: string
    date: string
    created_at: string
    updated_at: string
}

export interface HomeSection {
    id: number
    title: string
    content: string
    sort_order: number
}

export interface SocialLink {
    id: number
    platform: string
    url: string
    icon: string
    sort_order: number
}

export const usePortfolioStore = defineStore('portfolio', () => {
    const blogPosts = ref<BlogPost[]>([])
    const homeSections = ref<HomeSection[]>([])
    const socialLinks = ref<SocialLink[]>([])
    const siteSettings = ref<Record<string, string>>({})
    const loading = ref(false)

    async function fetchBlogPosts() {
        const { data } = await api.get('/blog')
        blogPosts.value = data
    }

    async function fetchHomeSections() {
        const { data } = await api.get('/home-sections')
        homeSections.value = data
    }

    async function fetchSocialLinks() {
        const { data } = await api.get('/social-links')
        socialLinks.value = data
    }

    async function fetchSiteSettings() {
        const { data } = await api.get('/site-settings')
        siteSettings.value = data
    }

    async function fetchHomeData() {
        loading.value = true
        await Promise.all([
            fetchHomeSections(),
            fetchSocialLinks(),
            fetchSiteSettings(),
        ])
        loading.value = false
    }

    // Admin CRUD methods
    async function createBlogPost(post: Omit<BlogPost, 'id' | 'created_at' | 'updated_at'>) {
        const { data } = await api.post('/blog', post)
        blogPosts.value.unshift(data)
        return data
    }

    async function updateBlogPost(id: number, post: Omit<BlogPost, 'id' | 'created_at' | 'updated_at'>) {
        const { data } = await api.put(`/blog/${id}`, post)
        const idx = blogPosts.value.findIndex(p => p.id === id)
        if (idx !== -1) blogPosts.value[idx] = data
        return data
    }

    async function deleteBlogPost(id: number) {
        await api.delete(`/blog/${id}`)
        blogPosts.value = blogPosts.value.filter(p => p.id !== id)
    }

    async function createHomeSection(section: Omit<HomeSection, 'id'>) {
        const { data } = await api.post('/home-sections', section)
        homeSections.value.push(data)
        homeSections.value.sort((a, b) => a.sort_order - b.sort_order)
        return data
    }

    async function updateHomeSection(id: number, section: Omit<HomeSection, 'id'>) {
        const { data } = await api.put(`/home-sections/${id}`, section)
        const idx = homeSections.value.findIndex(s => s.id === id)
        if (idx !== -1) homeSections.value[idx] = data
        return data
    }

    async function deleteHomeSection(id: number) {
        await api.delete(`/home-sections/${id}`)
        homeSections.value = homeSections.value.filter(s => s.id !== id)
    }

    async function createSocialLink(link: Omit<SocialLink, 'id'>) {
        const { data } = await api.post('/social-links', link)
        socialLinks.value.push(data)
        socialLinks.value.sort((a, b) => a.sort_order - b.sort_order)
        return data
    }

    async function updateSocialLink(id: number, link: Omit<SocialLink, 'id'>) {
        const { data } = await api.put(`/social-links/${id}`, link)
        const idx = socialLinks.value.findIndex(l => l.id === id)
        if (idx !== -1) socialLinks.value[idx] = data
        return data
    }

    async function deleteSocialLink(id: number) {
        await api.delete(`/social-links/${id}`)
        socialLinks.value = socialLinks.value.filter(l => l.id !== id)
    }

    async function updateSiteSettings(settings: Record<string, string>) {
        await api.put('/site-settings', settings)
        Object.assign(siteSettings.value, settings)
    }

    async function uploadImage(file: File): Promise<string> {
        const formData = new FormData()
        formData.append('image', file)
        const { data } = await api.post('/upload', formData, {
            headers: { 'Content-Type': 'multipart/form-data' },
        })
        return data.path
    }

    return {
        blogPosts, homeSections, socialLinks, siteSettings, loading,
        fetchBlogPosts, fetchHomeSections, fetchSocialLinks, fetchSiteSettings, fetchHomeData,
        createBlogPost, updateBlogPost, deleteBlogPost,
        createHomeSection, updateHomeSection, deleteHomeSection,
        createSocialLink, updateSocialLink, deleteSocialLink,
        updateSiteSettings, uploadImage,
    }
})
