const BASE = ''

function toQueryString(params) {
  const parts = []
  for (const [key, val] of Object.entries(params)) {
    if (val === undefined || val === null) continue
    if (Array.isArray(val)) {
      val.forEach(v => parts.push(`${encodeURIComponent(key)}=${encodeURIComponent(v)}`))
    } else {
      parts.push(`${encodeURIComponent(key)}=${encodeURIComponent(val)}`)
    }
  }
  return parts.join('&')
}

async function request(path, options = {}) {
  const token   = sessionStorage.getItem('rsyslox_token')
  const headers = { 'Content-Type': 'application/json' }
  if (token) headers['X-Session-Token'] = token

  const res = await fetch(BASE + path, { ...options, headers })

  if (res.status === 401 && !path.startsWith('/api/admin/login')) {
    sessionStorage.removeItem('rsyslox_token')
    sessionStorage.removeItem('rsyslox_role')
    const redirect = encodeURIComponent(window.location.pathname)
    window.location.href = `/login?redirect=${redirect}`
    throw new Error('Session expired')
  }

  const text = await res.text()
  let data
  try { data = text ? JSON.parse(text) : null } catch {}

  if (!res.ok) {
    const msg = data?.message || data?.error || `HTTP ${res.status}` + (text ? ': ' + text.slice(0, 120) : '')
    throw new Error(msg)
  }

  return data
}

export const api = {
  login:  (password) =>
    request('/api/admin/login', { method: 'POST', body: JSON.stringify({ password }) }),

  logout: () =>
    request('/api/admin/logout', { method: 'POST' }),

  setup: (payload) =>
    request('/api/setup', { method: 'POST', body: JSON.stringify(payload) }),

  getLogs: (params) =>
    request('/api/logs?' + toQueryString(params)),

  getMeta: () =>
    request('/api/meta'),

  getMetaColumn: (column, params = {}) =>
    request('/api/meta/' + column + (Object.keys(params).length ? '?' + toQueryString(params) : '')),

  getConfig: () =>
    request('/api/admin/config'),

  updateConfig: (patch) =>
    request('/api/admin/config', { method: 'PATCH', body: JSON.stringify(patch) }),

  generateSSL: () =>
    request('/api/admin/ssl/generate', { method: 'POST' }),

  getDiskUsage: () =>
    request('/api/admin/disk'),

  getCleanupStatus: () =>
    request('/api/admin/cleanup/status'),

  restart: () =>
    request('/api/admin/restart', { method: 'POST' }),

  uploadSSL: (certFile, keyFile) => {
    const form = new FormData()
    form.append('cert', certFile)
    form.append('key',  keyFile)
    const token = sessionStorage.getItem('rsyslox_token')
    const headers = {}
    if (token) headers['X-Session-Token'] = token
    return fetch('/api/admin/ssl/upload', { method: 'POST', headers, body: form })
      .then(async res => {
        if (!res.ok) {
          let msg = `HTTP ${res.status}`
          try { const e = await res.json(); msg = e.message || e.error || msg } catch {}
          throw new Error(msg)
        }
        return res.json()
      })
  },

  getKeys: () =>
    request('/api/admin/keys'),

  createKey: (name) =>
    request('/api/admin/keys', { method: 'POST', body: JSON.stringify({ name }) }),

  deleteKey: (name) =>
    request('/api/admin/keys/' + encodeURIComponent(name), { method: 'DELETE' }),

  health: () =>
    fetch(BASE + '/health').then(r => r.json()),
}
