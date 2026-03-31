import { reactive } from 'vue'

/**
 * Shared reactive state populated by the router's health check before any
 * route component is mounted. Used by App.vue (provide version) and
 * preferences.js (apply server defaults).
 */
export const appState = reactive({
  version: '',
  defaults: null, // { time_range: string, auto_refresh_interval: number }
})
