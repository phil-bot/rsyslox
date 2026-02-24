<template>
  <div class="setup-page">
    <div class="setup-card">
      <img :src="logoSrc" alt="rsyslox" class="logo" />
      <div class="setup-header">
        <h1>Setup</h1>
        <p class="subtitle">Configure rsyslox to get started.</p>
      </div>

      <form @submit.prevent="submit">
        <!-- Database -->
        <fieldset>
          <legend>Database (MySQL / MariaDB)</legend>
          <div class="row-2">
            <div class="field">
              <label for="db_host">Host</label>
              <input id="db_host" v-model="form.db_host" required placeholder="localhost" />
            </div>
            <div class="field">
              <label for="db_port">Port</label>
              <input id="db_port" v-model.number="form.db_port" type="number" placeholder="3306" />
            </div>
          </div>
          <div class="field">
            <label for="db_name">Database name</label>
            <input id="db_name" v-model="form.db_name" required placeholder="Syslog" />
          </div>
          <div class="row-2">
            <div class="field">
              <label for="db_user">User</label>
              <input id="db_user" v-model="form.db_user" required />
            </div>
            <div class="field">
              <label for="db_password">Password</label>
              <input id="db_password" v-model="form.db_password" type="password" required />
            </div>
          </div>
        </fieldset>

        <!-- Admin -->
        <fieldset>
          <legend>Admin account</legend>
          <div class="field">
            <label for="admin_password">Password</label>
            <input id="admin_password" v-model="form.admin_password" type="password"
              required minlength="12" placeholder="At least 12 characters" />
          </div>
          <div class="field">
            <label for="confirm_password">Confirm password</label>
            <input id="confirm_password" v-model="confirmPassword" type="password" required />
          </div>
        </fieldset>

        <!-- Server -->
        <fieldset>
          <legend>Server</legend>
          <div class="row-2">
            <div class="field">
              <label for="server_host">Bind host</label>
              <input id="server_host" v-model="form.server_host" placeholder="0.0.0.0" />
            </div>
            <div class="field">
              <label for="server_port">Port</label>
              <input id="server_port" v-model.number="form.server_port" type="number" placeholder="8000" />
            </div>
          </div>
        </fieldset>

        <p v-if="error" class="msg error">{{ error }}</p>

        <button type="submit" class="btn btn-primary submit-btn" :disabled="loading">
          {{ loading ? 'Saving…' : 'Complete Setup' }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, inject, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '@/api/client'
import { useAuthStore } from '@/stores/auth'

const theme   = inject('theme')
const logoSrc = computed(() => theme?.value === 'dark' ? '/logo-dark.svg' : '/logo-light.svg')
const router  = useRouter()
const auth    = useAuthStore()

const form = ref({
  db_host: 'localhost',
  db_port: 3306,
  db_name: 'Syslog',
  db_user: '',
  db_password: '',
  admin_password: '',
  server_host: '0.0.0.0',
  server_port: 8000,
  use_ssl: false,
})
const confirmPassword = ref('')
const error   = ref('')
const loading = ref(false)

// Load prefill values from server (env vars set by Docker entrypoint)
onMounted(async () => {
  try {
    const prefill = await fetch('/api/setup').then(r => r.ok ? r.json() : null)
    if (prefill) {
      if (prefill.db_host)   form.value.db_host   = prefill.db_host
      if (prefill.db_port)   form.value.db_port   = prefill.db_port
      if (prefill.db_name)   form.value.db_name   = prefill.db_name
      if (prefill.db_user)   form.value.db_user   = prefill.db_user
      if (prefill.server_host) form.value.server_host = prefill.server_host
      if (prefill.server_port) form.value.server_port = prefill.server_port
    }
  } catch { /* prefill is optional */ }
})

async function submit() {
  error.value = ''
  if (form.value.admin_password !== confirmPassword.value) {
    error.value = 'Passwords do not match'
    return
  }
  loading.value = true
  try {
    await api.setup(form.value)
    // Config is now written on disk. The router re-checks /health on every
    // navigation — it will see setup_mode=false and allow /login.
    router.push('/login')
  } catch (e) {
    error.value = e.body?.message || 'Setup failed'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.setup-page {
  min-height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg);
  padding: 2rem 1rem;
}

.setup-card {
  width: 100%;
  max-width: 520px;
  padding: 2rem;
  background: var(--bg-surface);
  border: 1px solid var(--border);
  border-radius: calc(var(--radius) * 2);
  box-shadow: 0 4px 24px rgba(0,0,0,.07);
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.logo { height: 32px; width: auto; }

.setup-header { display: flex; flex-direction: column; gap: .25rem; }
h1 { font-size: 1.4rem; font-weight: 700; color: var(--color-primary); }
.subtitle { color: var(--text-muted); font-size: .875rem; }

fieldset {
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: .75rem;
}
legend {
  font-size: .8rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: .05em;
  color: var(--text-muted);
  padding: 0 .25rem;
}

.row-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: .75rem;
}

.field {
  display: flex;
  flex-direction: column;
  gap: .3rem;
}
.field label {
  font-size: .8rem;
  font-weight: 500;
  color: var(--text-muted);
}
.field input {
  padding: .5rem .625rem;
  border: 1px solid var(--border);
  border-radius: var(--radius);
  background: var(--bg);
  color: var(--text);
  font-size: .875rem;
  width: 100%;
  transition: border-color .15s;
}
.field input:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(2,132,199,.12);
}

.msg { font-size: .875rem; padding: .5rem .75rem; border-radius: var(--radius); border: 1px solid; }
.error { color: #dc2626; background: #fef2f2; border-color: #fca5a5; }
[data-theme="dark"] .error { background: #2d1212; border-color: #7f1d1d; color: #fca5a5; }

.submit-btn { width: 100%; justify-content: center; padding: .625rem; }
</style>
