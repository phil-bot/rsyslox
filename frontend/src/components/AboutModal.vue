<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="modelValue" class="about-backdrop" @click.self="close">
        <div class="about-dialog" role="dialog" aria-modal="true" aria-label="About rsyslox">

          <div class="about-header">
            <img :src="logoSrc" alt="rsyslox" class="about-logo" />
            <button class="close-btn" @click="close" title="Close">✕</button>
          </div>

          <div class="about-body">
            <p class="about-tagline">
              A self-hosted syslog viewer for rsyslog data stored in MySQL&thinsp;/&thinsp;MariaDB.
            </p>

            <table class="about-table">
              <tbody>
                <tr>
                  <td class="about-key">Version</td>
                  <td class="about-val mono">{{ appState?.version || '—' }}</td>
                </tr>
                <tr>
                  <td class="about-key">License</td>
                  <td class="about-val">MIT</td>
                </tr>
                <tr>
                  <td class="about-key">Author</td>
                  <td class="about-val">Phillip Grothues</td>
                </tr>
                <tr>
                  <td class="about-key">Source</td>
                  <td class="about-val">
                    <a
                      href="https://github.com/phil-bot/rsyslox"
                      target="_blank"
                      rel="noopener"
                      class="about-link"
                    >github.com/phil-bot/rsyslox ↗</a>
                  </td>
                </tr>
                <tr>
                  <td class="about-key">Docs</td>
                  <td class="about-val">
                    <a
                      href="https://rsyslox.grothu.net"
                      target="_blank"
                      rel="noopener"
                      class="about-link"
                    >rsyslox.grothu.net ↗</a>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <div class="about-footer">
            <button class="btn btn-ghost" @click="close">Close</button>
          </div>

        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { inject, computed } from 'vue'

defineProps({
  modelValue: { type: Boolean, default: false },
})
const emit = defineEmits(['update:modelValue'])

const theme    = inject('theme')
const appState = inject('appVersion', null) // reactive appState { version, defaults }

const logoSrc = computed(() =>
  theme?.value === 'dark' ? '/logo-dark.svg' : '/logo-light.svg'
)

function close() {
  emit('update:modelValue', false)
}
</script>

<style scoped>
.about-backdrop {
  position: fixed; inset: 0;
  background: rgba(0, 0, 0, .45);
  z-index: 400;
  display: flex; align-items: center; justify-content: center;
  padding: 1rem;
}
.about-dialog {
  background: var(--bg-surface);
  border: 1px solid var(--border);
  border-radius: calc(var(--radius) * 2);
  box-shadow: 0 20px 60px rgba(0, 0, 0, .2);
  width: 100%; max-width: 420px;
  display: flex; flex-direction: column; overflow: hidden;
}
.about-header {
  display: flex; align-items: center; justify-content: space-between;
  padding: 1.25rem 1.25rem .875rem;
  border-bottom: 1px solid var(--border);
}
.about-logo { height: 36px; width: auto; }
.close-btn {
  background: none; border: none; cursor: pointer;
  color: var(--text-muted); font-size: 1rem; padding: .25rem .4rem;
  border-radius: var(--radius); transition: background .15s, color .15s;
}
.close-btn:hover { background: var(--bg-hover); color: var(--text); }
.about-body {
  padding: 1rem 1.25rem;
  display: flex; flex-direction: column; gap: .875rem;
}
.about-tagline { font-size: .825rem; color: var(--text-muted); line-height: 1.5; }
.about-table { width: 100%; border-collapse: collapse; }
.about-table tr { border-bottom: 1px solid var(--border); }
.about-table tr:last-child { border-bottom: none; }
.about-key {
  padding: .5rem 0; font-size: .72rem; font-weight: 700;
  text-transform: uppercase; letter-spacing: .06em;
  color: var(--text-muted); width: 32%; vertical-align: middle;
}
.about-val {
  padding: .5rem 0 .5rem .75rem;
  font-size: .875rem; color: var(--text); vertical-align: middle;
}
.mono { font-family: ui-monospace, 'SF Mono', Menlo, monospace; }
.about-link { color: var(--color-primary); text-decoration: none; }
.about-link:hover { text-decoration: underline; }
.about-footer {
  padding: .875rem 1.25rem;
  border-top: 1px solid var(--border);
  display: flex; justify-content: flex-end;
}
.modal-enter-active, .modal-leave-active { transition: opacity .2s; }
.modal-enter-from, .modal-leave-to { opacity: 0; }
</style>
