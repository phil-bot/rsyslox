<template>
  <header class="app-header">
    <!-- Left: logo + text nav links -->
    <div class="header-left">
      <a href="/logs" class="logo-link">
        <img :src="logoSrc" alt="rsyslox" class="logo-img" />
      </a>

      <nav class="header-nav">
      <RouterLink to="/logs" class="nav-item" :class="{ active: route.path === '/logs' }">
          {{ t('nav.logs') }}
        </RouterLink>
        <span class="nav-item nav-item--soon" :title="t('nav.statistics_soon')">
          {{ t('nav.statistics') }}
        </span>
      </nav>
    </div>

    <!-- Right: icon buttons -->
    <div class="header-right">

      <!-- Dark / light toggle switch -->
      <button
        class="theme-toggle"
        :class="{ dark: theme === 'dark' }"
        @click="toggleTheme()"
        :title="theme === 'dark' ? t('nav.toggle_theme_dark') : t('nav.toggle_theme_light')"
        role="switch"
        :aria-checked="theme === 'dark'"
      >
        <!-- Sun icon -->
        <svg class="theme-icon theme-icon--sun" width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
          <circle cx="12" cy="12" r="5"/>
          <line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/>
          <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/>
          <line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/>
          <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/>
        </svg>
        <span class="toggle-track">
          <span class="toggle-thumb"></span>
        </span>
        <!-- Moon icon -->
        <svg class="theme-icon theme-icon--moon" width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
          <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
        </svg>
      </button>

      <div class="header-divider"></div>

      <!-- Settings (admin only) -->
      <RouterLink
        v-if="auth.isAdmin"
        to="/admin"
        class="icon-btn"
        :class="{ active: route.path.startsWith('/admin') }"
        :title="t('nav.settings')"
      >
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="3"/>
          <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/>
        </svg>
      </RouterLink>

      <!-- Logout — door/exit icon -->
      <button v-if="auth.isAuthenticated" class="icon-btn" @click="logout" :title="t('nav.logout')">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
          <polyline points="16 17 21 12 16 7"/>
          <line x1="21" y1="12" x2="9" y2="12"/>
        </svg>
      </button>
    </div>
  </header>

</template>

<script setup>
import { ref, inject, computed } from 'vue'
import { useRoute, useRouter, RouterLink } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { api } from '@/api/client'
import { useLocale } from '@/composables/useLocale'

const { t } = useLocale()

const theme       = inject('theme')
const toggleTheme = inject('toggleTheme')
const logoSrc     = computed(() => theme.value === 'dark' ? '/logo-dark.svg' : '/logo-light.svg')

const route  = useRoute()
const router = useRouter()
const auth   = useAuthStore()


async function logout() {
  try { await api.logout() } catch {}
  auth.clearSession()
  router.push('/login')
}
</script>

<style scoped>
.app-header {
  height: var(--header-height);
  display: flex; align-items: center;
  padding: 0 .75rem;
  background: var(--bg-surface);
  border-bottom: 1px solid var(--border);
  position: sticky; top: 0; z-index: 100;
  flex-shrink: 0; gap: .5rem;
}

.header-left {
  display: flex; align-items: center; gap: .5rem; flex: 1;
}
.header-right {
  display: flex; align-items: center; gap: .125rem; flex-shrink: 0;
}

.logo-link { display: flex; align-items: center; flex-shrink: 0; }
.logo-img  { height: 40px; width: auto; }

.header-nav {
  display: flex; align-items: center; gap: .125rem; margin-left: .25rem;
}

.nav-item {
  padding: .375rem .7rem; border-radius: var(--radius);
  font-size: .9rem; font-weight: 500; color: var(--text-muted);
  text-decoration: none; white-space: nowrap;
  transition: background .15s, color .15s;
}
.nav-item:hover  { background: var(--bg-hover); color: var(--text); }
.nav-item.active { color: var(--color-primary); background: var(--bg-selected); }
.nav-item--soon  { opacity: .4; cursor: default; pointer-events: none; font-style: italic; }

/* ── Generic icon button ────────────────── */
.icon-btn {
  display: flex; align-items: center; justify-content: center;
  width: 34px; height: 34px; position: relative;
  background: none; border: none; border-radius: var(--radius);
  cursor: pointer; color: var(--text-muted); text-decoration: none;
  transition: background .15s, color .15s; flex-shrink: 0;
}
.icon-btn:hover  { background: var(--bg-hover); color: var(--text); }
.icon-btn.active { color: var(--color-primary); background: var(--bg-selected); }

/* External-link badge on Docs button */
.ext-badge {
  position: absolute; top: 4px; right: 4px;
  color: var(--text-muted); opacity: .7;
}
.icon-btn:hover .ext-badge { opacity: 1; }

/* ── Theme toggle switch ────────────────── */
.theme-toggle {
  display: flex; align-items: center; gap: 5px;
  padding: 0 6px; height: 28px;
  background: none; border: 1px solid var(--border);
  border-radius: 999px; cursor: pointer;
  color: var(--text-muted);
  transition: border-color .2s, background .2s;
  flex-shrink: 0;
}
.theme-toggle:hover { border-color: var(--text-muted); background: var(--bg-hover); }

.theme-icon { flex-shrink: 0; transition: opacity .2s; }
.theme-icon--sun  { opacity: 1; }
.theme-icon--moon { opacity: .4; }
.theme-toggle.dark .theme-icon--sun  { opacity: .4; }
.theme-toggle.dark .theme-icon--moon { opacity: 1; }

.toggle-track {
  width: 28px; height: 16px; border-radius: 999px;
  background: var(--border);
  position: relative; flex-shrink: 0;
  transition: background .2s;
}
.theme-toggle.dark .toggle-track { background: var(--color-primary); }

.toggle-thumb {
  position: absolute; top: 2px; left: 2px;
  width: 12px; height: 12px; border-radius: 50%;
  background: #fff;
  box-shadow: 0 1px 3px rgba(0,0,0,.25);
  transition: transform .2s ease;
}
.theme-toggle.dark .toggle-thumb { transform: translateX(12px); }

/* ── Divider ────────────────────────────── */
.header-divider {
  /* USE ONLY MARGIN !
  /* width: 1px; height: 20px;
  background: var(--border);*/
  margin: 0 .25rem; flex-shrink: 0;
}

@media (max-width: 380px) { .header-nav { display: none; } }
</style>
