import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { appState } from '@/stores/appState'
import { applyServerDefaults } from '@/stores/preferences'

const router = createRouter({
  history: createWebHistory(),
  routes: [{
      path: '/',
      name: 'home',
      meta: {
        public: true
      },
    },
    {
      path: '/setup',
      name: 'setup',
      component: () => import('@/views/SetupView.vue'),
      meta: {
        public: true
      },
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
      meta: {
        public: true
      },
    },
    {
      path: '/logs',
      name: 'logs',
      component: () => import('@/views/LogsView.vue'),
      meta: {
        requiresAuth: true
      },
    },
    {
      path: '/admin',
      name: 'admin',
      component: () => import('@/views/AdminView.vue'),
      meta: {
        requiresAdmin: true
      },
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: () => {
        const auth = useAuthStore()
        return auth.isAuthenticated.value ? '/logs' : '/login'
      },
    },
  ],
})

async function fetchHealthAndApplyDefaults() {
  try {
    const res = await fetch('/health')
    const data = await res.json()
    if (data.version)  appState.version  = data.version
      if (data.defaults) {
        appState.defaults = data.defaults
        applyServerDefaults(data.defaults)
      }
      return data
  } catch {
    return {}
  }
}

router.beforeEach(async (to) => {
  const data = await fetchHealthAndApplyDefaults()
  const auth = useAuthStore()

  // Setup-Modus aktiv → nur /setup erlaubt
  if (data.setup_mode === true) {
    if (to.name !== 'setup') return { name: 'setup' }
    return true
  }

  // Setup-Modus inaktiv → /setup sperren
  if (to.name === 'setup') {
    return auth.isAuthenticated.value ? { name: 'logs' } : { name: 'login' }
  }

  // "/" → direkt weiterleiten
  if (to.name === 'home') {
    return auth.isAuthenticated.value ? { name: 'logs' } : { name: 'login' }
  }

  if (to.meta.public) return true

    if (!auth.isAuthenticated.value) {
      return { name: 'login' }
    }

    if (to.meta.requiresAdmin && !auth.isAdmin.value) {
      return { name: 'logs' }
    }

    return true
})

export default router
