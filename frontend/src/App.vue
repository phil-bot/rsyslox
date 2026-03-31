<template>
  <div :data-theme="theme" class="app-root">
    <!--
      Show a blank screen until router.isReady() resolves.

      With the smart redirect on '/' the initial navigation goes directly
      to the correct destination (/login or /logs) without an intermediate
      step — so router.isReady() resolves only once, on the final route.
      No flash of protected views is possible.
    -->
    <div v-if="!ready" class="app-loading" aria-hidden="true"><div class="loader"></div></div>
    <RouterView v-else />
  </div>
</template>

<script setup>
import { ref, provide } from 'vue'
import { useRouter } from 'vue-router'
import { appState } from '@/stores/appState'

const router = useRouter()
const ready  = ref(false)

// ── Theme ─────────────────────────────────────────────────────────────────────
// Restore synchronously before first paint — avoids a colour-mode flash.
const theme = ref('light')
const savedTheme = localStorage.getItem('rsyslox_theme')
if (savedTheme) {
  theme.value = savedTheme
} else if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
  theme.value = 'dark'
}

provide('theme', theme)
provide('toggleTheme', () => {
  theme.value = theme.value === 'dark' ? 'light' : 'dark'
  localStorage.setItem('rsyslox_theme', theme.value)
})
provide('appVersion', appState)

// ── Render gate ───────────────────────────────────────────────────────────────
// Call at setup time (before mount) so the promise resolves only after the
// single initial navigation — directly to /login or /logs — has finished.
router.isReady().then(() => {
  ready.value = true
})
</script>

<style>
.app-root {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.app-loading {
  height: 100%;
  background: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center
}

[data-theme="dark"] .app-loading {
  background: #0b1220;
}

/* HTML: <div class="loader"></div> */
.loader {
  width: 50px;
  aspect-ratio: 1;
  border-radius: 50%;
  border: 8px solid #0000;
  border-right-color: #f59e0b97;
  position: relative;
  animation: l24 1s infinite linear;
}
.loader:before,
.loader:after {
  content: "";
  position: absolute;
  inset: -8px;
  border-radius: 50%;
  border: inherit;
  animation: inherit;
  animation-duration: 2s;
}
.loader:after {
  animation-duration: 4s;
}
@keyframes l24 {
  100% {transform: rotate(1turn)}
}


</style>
