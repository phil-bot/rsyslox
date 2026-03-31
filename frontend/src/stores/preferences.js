import { ref, watch } from 'vue'

const STORAGE_KEY = 'rsyslox_prefs'

function load() {
  try { return JSON.parse(localStorage.getItem(STORAGE_KEY) || '{}') } catch { return {} }
}
function save(prefs) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(prefs))
}

const stored = load()

export const language            = ref(stored.language            ?? 'en')
export const timeFormat          = ref(stored.timeFormat          ?? '24h')
export const fontSize            = ref(stored.fontSize            ?? 'medium')
export const autoRefreshInterval = ref(stored.autoRefreshInterval ?? 30)
export const defaultTimeRange    = ref(stored.defaultTimeRange    ?? '24h')

// Apply font-size immediately on load
applyFontSize(fontSize.value)

/**
 * Font size map — 14 / 16 / 18 px
 */
export function applyFontSize(size) {
  const map = { small: '14px', medium: '16px', large: '18px' }
  document.documentElement.style.setProperty('font-size', map[size] ?? '16px')
}

/**
 * Apply server-configured defaults for preferences that the user has not yet
 * explicitly set in their own localStorage.
 *
 * Called from the router guard before any route component is mounted.
 * Uses hasOwnProperty on the raw stored object so we can distinguish
 * "never set" from "set to a value that happens to equal the fallback".
 */
export function applyServerDefaults(defaults) {
  if (!defaults) return

  let rawStored = {}
  try {
    rawStored = JSON.parse(localStorage.getItem(STORAGE_KEY) || '{}')
  } catch {}

  const has = key => Object.prototype.hasOwnProperty.call(rawStored, key)

  if (!has('defaultTimeRange') && defaults.time_range)
    defaultTimeRange.value = defaults.time_range

  if (!has('autoRefreshInterval') && defaults.auto_refresh_interval)
    autoRefreshInterval.value = defaults.auto_refresh_interval

  if (!has('language') && defaults.language)
    language.value = defaults.language

  if (!has('fontSize') && defaults.font_size)
    fontSize.value = defaults.font_size

  if (!has('timeFormat') && defaults.time_format)
    timeFormat.value = defaults.time_format
}

// Persist + apply on every change
watch([language, timeFormat, fontSize, autoRefreshInterval, defaultTimeRange], () => {
  save({
    language:            language.value,
    timeFormat:          timeFormat.value,
    fontSize:            fontSize.value,
    autoRefreshInterval: autoRefreshInterval.value,
    defaultTimeRange:    defaultTimeRange.value,
  })
  applyFontSize(fontSize.value)
})

export function usePreferences() {
  return { language, timeFormat, fontSize, autoRefreshInterval, defaultTimeRange }
}
