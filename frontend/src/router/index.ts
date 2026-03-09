import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
    history: createWebHistory(),
    routes: [
        // Public routes
        {
            path: '/',
            name: 'home',
            component: () => import('@/views/HomeView.vue'),
        },
        {
            path: '/blog',
            name: 'blog',
            component: () => import('@/views/BlogView.vue'),
        },
        {
            path: '/login',
            name: 'login',
            component: () => import('@/views/LoginView.vue'),
        },
        // Protected routes
        {
            path: '/dashboard',
            name: 'dashboard',
            component: () => import('@/views/DashboardView.vue'),
            meta: { requiresAuth: true },
        },
        {
            path: '/settings',
            name: 'settings',
            component: () => import('@/views/SettingsView.vue'),
            meta: { requiresAuth: true },
        },
        {
            path: '/scrape',
            name: 'scrape',
            component: () => import('@/views/ScrapeView.vue'),
            meta: { requiresAuth: true },
        },
        {
            path: '/admin',
            name: 'admin',
            component: () => import('@/views/AdminView.vue'),
            meta: { requiresAuth: true },
        },
    ],
})

router.beforeEach((to) => {
    const token = localStorage.getItem('token')
    if (to.meta.requiresAuth && !token) {
        return { name: 'login' }
    }
})

export default router
