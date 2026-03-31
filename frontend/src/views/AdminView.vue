<template>
  <div class="admin-layout">
    <AppHeader />

    <div class="admin-body">
      <!-- Restart banner -->
      <Transition name="banner">
        <div v-if="restartNeeded" class="restart-banner">
          <span class="banner-icon">⚠</span>
          <span class="banner-text">{{ t('admin.restart_banner') }}</span>
          <div class="banner-actions">
            <button class="btn btn-ghost btn-sm" @click="restartNeeded = false">{{ t('admin.restart_dismiss') }}</button>
            <button class="btn btn-primary btn-sm" :disabled="restarting" @click="confirmRestart = true">{{ t('admin.restart_now') }}</button>
          </div>
        </div>
      </Transition>

      <div class="admin-content">
        <!-- Tab nav -->
        <nav class="admin-nav">
          <div>
            <button v-for="tab in tabs" :key="tab.id" class="nav-tab"
              :class="{ active: activeTab === tab.id }" @click="activeTab = tab.id">
              <span class="tab-icon" v-html="tab.svg"></span>
              {{ tab.label }}
            </button>
          </div>
          <!-- About — info/tag icon -->
          <button class="nav-tab" @click="showAbout = true" :title="t('nav.about')">
            <svg width="16" height="16" viewBox="0 0 1024 1024" fill="currentColor">
              <path fill="currentColor" d="M512 64a448 448 0 1 1 0 896.064A448 448 0 0 1 512 64zm67.2 275.072c33.28 0 60.288-23.104 60.288-57.344s-27.072-57.344-60.288-57.344c-33.28 0-60.16 23.104-60.16 57.344s26.88 57.344 60.16 57.344zM590.912 699.2c0-6.848 2.368-24.64 1.024-34.752l-52.608 60.544c-10.88 11.456-24.512 19.392-30.912 17.28a12.992 12.992 0 0 1-8.256-14.72l87.68-276.992c7.168-35.136-12.544-67.2-54.336-71.296-44.096 0-108.992 44.736-148.48 101.504 0 6.784-1.28 23.68.064 33.792l52.544-60.608c10.88-11.328 23.552-19.328 29.952-17.152a12.8 12.8 0 0 1 7.808 16.128L388.48 728.576c-10.048 32.256 8.96 63.872 55.04 71.04 67.84 0 107.904-43.648 147.456-100.416z"/>
            </svg>
            {{ t('nav.about') }}
          </button>
        </nav>



        <!-- Main content area -->
        <div class="admin-main-wrapper">
          <main class="admin-main">

            <!-- ── Server ──────────────────────────────────── -->
            <section v-if="activeTab === 'server'" class="admin-section">
              <div class="section-header">
                <h2>{{ t('admin.server_title') }}</h2>
                <p class="section-desc">{{ t('admin.server_desc') }}</p>
              </div>
              <div v-if="configLoading" class="loading-msg">{{ t('admin.loading') }}</div>
              <template v-else>
                <div class="config-form">
                  <div class="field-row">
                    <label class="field-label">{{ t('admin.field_host') }}
                      <input v-model="serverForm.host" class="field-input" placeholder="0.0.0.0" />
                      <span class="field-hint">{{ t('admin.field_host_hint') }}</span>
                      <span class="field-hint restart-hint">{{ t('admin.hint_restart') }}</span>
                    </label>
                    <label class="field-label">{{ t('admin.field_port') }}
                      <input v-model.number="serverForm.port" type="number" min="1" max="65535" class="field-input" />
                      <span class="field-hint restart-hint">{{ t('admin.hint_restart') }}</span>
                    </label>
                  </div>
                  <label class="field-label">{{ t('admin.field_origins') }}
                    <input v-model="serverForm.originsStr" class="field-input" placeholder="*" />
                    <span class="field-hint">{{ t('admin.field_origins_hint') }}</span>
                  </label>
                  <label class="toggle-label">
                    <input type="checkbox" v-model="serverForm.useSSL" />
                    {{ t('admin.field_ssl') }}
                  </label>
                </div>

                <!-- Default values -->
                <div class="subsection-header">
                  <h3>{{ t('admin.defaults_title') }}</h3>
                  <p class="section-desc">{{ t('admin.defaults_desc') }}</p>
                </div>
                <div class="config-form">
                  <div class="field-row">
                    <label class="field-label">{{ t('admin.default_time_range') }}
                      <select v-model="serverForm.defaultTimeRange" class="field-input">
                        <option value="15m">15m</option>
                        <option value="1h">1h</option>
                        <option value="6h">6h</option>
                        <option value="24h">24h</option>
                        <option value="7d">7d</option>
                        <option value="30d">30d</option>
                      </select>
                      <span class="field-hint">{{ t('admin.default_time_range_hint') }}</span>
                    </label>
                    <label class="field-label">{{ t('admin.default_auto_refresh') }}
                      <div class="inline-field">
                        <input v-model.number="serverForm.autoRefreshInterval" type="number" min="5" max="300"
                          class="field-input" style="max-width:100px" />
                        <span class="field-hint">s</span>
                      </div>
                    </label>
                  </div>
                  <div class="field-row">
                    <label class="field-label">{{ t('admin.default_language') }}
                      <select v-model="serverForm.defaultLanguage" class="field-input">
                        <option value="en">English</option>
                        <option value="de">Deutsch</option>
                      </select>
                    </label>
                    <label class="field-label">{{ t('admin.default_font_size') }}
                      <select v-model="serverForm.defaultFontSize" class="field-input">
                        <option value="small">{{ t('prefs.font_small') }}</option>
                        <option value="medium">{{ t('prefs.font_medium') }}</option>
                        <option value="large">{{ t('prefs.font_large') }}</option>
                      </select>
                    </label>
                  </div>
                  <label class="field-label">{{ t('admin.default_time_format') }}
                    <select v-model="serverForm.defaultTimeFormat" class="field-input" style="max-width:260px">
                      <option value="24h">{{ t('prefs.time_format_24h') }}</option>
                      <option value="12h">{{ t('prefs.time_format_12h') }}</option>
                    </select>
                  </label>
                </div>

                <!-- SSL section -->
                <template v-if="serverForm.useSSL">
                  <div class="subsection-header">
                    <h3>{{ t('admin.ssl_title') }}</h3>
                    <p class="section-desc">{{ t('admin.ssl_desc') }}</p>
                  </div>
                  <div class="config-form">
                    <div class="ssl-block">
                      <p class="field-label">{{ t('admin.ssl_generate_title') }}</p>
                      <p class="field-hint">{{ t('admin.ssl_selfsigned_desc') }}</p>
                      <div class="inline-actions">
                        <button type="button" class="btn btn-ghost" :disabled="sslGenerating" @click="generateSSL">
                          {{ sslGenerating ? t('admin.ssl_generating') : t('admin.ssl_generate') }}
                        </button>
                        <span v-if="sslMsg" class="action-msg" :class="sslMsgOk ? 'ok' : 'err'">{{ sslMsg }}</span>
                      </div>
                    </div>
                    <div class="ssl-block">
                      <p class="field-label">{{ t('admin.ssl_upload_title') }}</p>
                      <div class="field-row">
                        <label class="field-label">{{ t('admin.ssl_cert_file') }}
                          <input type="file" accept=".pem,.crt,.cer" class="field-input file-input"
                            @change="sslCertFile = $event.target.files[0]" />
                        </label>
                        <label class="field-label">{{ t('admin.ssl_key_file') }}
                          <input type="file" accept=".pem,.key" class="field-input file-input"
                            @change="sslKeyFile = $event.target.files[0]" />
                        </label>
                      </div>
                      <div class="inline-actions">
                        <button type="button" class="btn btn-ghost"
                          :disabled="sslUploading || !sslCertFile || !sslKeyFile" @click="uploadSSL">
                          {{ sslUploading ? t('admin.ssl_uploading') : t('admin.ssl_upload') }}
                        </button>
                        <span v-if="sslUploadMsg" class="action-msg" :class="sslUploadMsgOk ? 'ok' : 'err'">{{ sslUploadMsg }}</span>
                      </div>
                    </div>
                  </div>
                </template>
              </template>
            </section>

            <!-- ── Database + Cleanup ──────────────────────── -->
            <section v-if="activeTab === 'database'" class="admin-section">
              <div class="section-header">
                <h2>{{ t('admin.db_title') }}</h2>
                <p class="section-desc">{{ t('admin.db_desc') }}</p>
              </div>
              <div v-if="configLoading" class="loading-msg">{{ t('admin.loading') }}</div>
              <template v-else>
                <div class="config-form">
                  <div class="field-row">
                    <label class="field-label">{{ t('admin.field_host') }}
                      <input v-model="dbForm.host" class="field-input" placeholder="localhost" />
                      <span class="field-hint restart-hint">{{ t('admin.hint_restart') }}</span>
                    </label>
                    <label class="field-label">{{ t('admin.field_port') }}
                      <input v-model.number="dbForm.port" type="number" min="1" max="65535"
                        class="field-input" style="max-width:120px" />
                    </label>
                  </div>
                  <div class="field-row">
                    <label class="field-label">{{ t('admin.db_name') }}
                      <input v-model="dbForm.name" class="field-input" placeholder="Syslog" />
                    </label>
                    <label class="field-label">{{ t('admin.db_user') }}
                      <input v-model="dbForm.user" class="field-input" />
                    </label>
                  </div>
                  <label class="field-label">{{ t('admin.db_password') }}
                    <input v-model="dbForm.password" type="password" class="field-input"
                      :placeholder="t('admin.db_password_placeholder')" autocomplete="new-password" />
                    <span class="field-hint">{{ t('admin.db_password_hint') }}</span>
                  </label>
                </div>

                <div class="subsection-header">
                  <h3>{{ t('admin.cleanup_section_title') }}</h3>
                  <p class="section-desc">{{ t('admin.cleanup_desc') }}</p>
                </div>
                <div class="info-callout">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="flex-shrink:0;margin-top:1px"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
                  {{ t('admin.cleanup_local_hint') }}
                </div>
                <div class="config-form">
                  <label class="toggle-label">
                    <input type="checkbox" v-model="cleanupForm.enabled" />
                    {{ t('admin.cleanup_enable') }}
                  </label>
                  <div class="field-row" :class="{ disabled: !cleanupForm.enabled }">
                    <label class="field-label">{{ t('admin.cleanup_disk_path') }}
                      <input v-model="cleanupForm.diskPath" class="field-input"
                        :disabled="!cleanupForm.enabled" placeholder="/var/lib/mysql"
                        @change="loadDiskUsage" />
                    </label>
                    <label class="field-label">{{ t('admin.cleanup_threshold') }}
                      <div class="inline-field">
                        <input v-model.number="cleanupForm.thresholdPercent" type="number" min="1" max="100"
                          class="field-input" style="max-width:90px" :disabled="!cleanupForm.enabled" />
                        <span class="field-hint">%</span>
                      </div>
                    </label>
                  </div>
                  <div class="field-row" :class="{ disabled: !cleanupForm.enabled }">
                    <label class="field-label">{{ t('admin.cleanup_batch') }}
                      <input v-model.number="cleanupForm.batchSize" type="number" min="1"
                        class="field-input" style="max-width:120px" :disabled="!cleanupForm.enabled" />
                    </label>
                    <label class="field-label">{{ t('admin.cleanup_interval') }}
                      <div class="inline-field">
                        <input v-model.number="cleanupForm.intervalSeconds" type="number" min="60"
                          class="field-input" style="max-width:120px" :disabled="!cleanupForm.enabled" />
                        <span class="field-hint">s</span>
                      </div>
                    </label>
                  </div>
                  <div class="disk-widget">
                    <div class="disk-widget-header">
                      <span class="field-label" style="margin:0">{{ t('admin.disk_usage') }}</span>
                      <span class="disk-path-label">{{ cleanupForm.diskPath || '/' }}</span>
                      <button type="button" class="btn btn-ghost btn-sm" @click="loadDiskUsage" :disabled="diskLoading">↻</button>
                    </div>
                    <div v-if="diskLoading" class="disk-loading">{{ t('admin.disk_loading') }}</div>
                    <div v-else-if="diskError" class="disk-error">{{ t('admin.disk_error') }}: {{ diskError }}</div>
                    <template v-else-if="diskInfo">
                      <div class="disk-bar-track">
                        <div class="disk-bar-fill" :class="diskBarClass" :style="{ width: diskInfo.used_percent.toFixed(1) + '%' }"></div>
                      </div>
                      <div class="disk-bar-labels">
                        <span>{{ diskInfo.used_percent.toFixed(1) }}% used</span>
                        <span>{{ formatBytes(diskInfo.free_bytes) }} {{ t('admin.disk_free') }} {{ t('admin.disk_of') }} {{ formatBytes(diskInfo.total_bytes) }}</span>
                      </div>
                    </template>
                  </div>
                </div>
              </template>
            </section>

            <!-- ── API Keys ────────────────────────────────── -->
            <section v-if="activeTab === 'keys'" class="admin-section">
              <div class="section-header">
                <h2>{{ t('admin.keys_title') }}</h2>
                <p class="section-desc">{{ t('admin.keys_desc') }}</p>
              </div>

              <div class="info-callout">
                <!-- Docs — external link indicator -->
                <a href="/docs/" target="_blank" rel="noopener" class="no-gap" :title="t('nav.docs')">
                  <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 5px">
                    <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"/><path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"/>
                  </svg>
                  {{ t('nav.docs') }} ↗
                </a>
              </div>





              <div class="new-key-form">
                <input v-model="newKeyName" class="field-input"
                  :placeholder="t('admin.keys_placeholder')" @keydown.enter.prevent="createKey"
                  style="max-width:280px" />
                <button class="btn btn-primary" :disabled="!newKeyName.trim() || keyCreating" @click="createKey">
                  {{ keyCreating ? t('admin.keys_creating') : t('admin.keys_create') }}
                </button>
              </div>
              <div v-if="newKeyPlaintext" class="key-reveal">
                <div class="key-reveal-header">
                  <span>🔑</span>
                  <strong>{{ t('admin.keys_created_for') }} "{{ newKeyRevealName }}"</strong>
                  <span class="key-reveal-warn">{{ t('admin.keys_copy_note') }}</span>
                </div>
                <div class="key-reveal-value">
                  <code class="mono">{{ newKeyPlaintext }}</code>
                  <button class="btn btn-primary btn-sm" @click="copyKey">
                    {{ keyCopied ? t('admin.keys_copied') : t('admin.keys_copy') }}
                  </button>
                </div>
                <button class="btn btn-ghost btn-sm" @click="newKeyPlaintext = null">{{ t('admin.keys_dismiss') }}</button>
              </div>
              <div v-if="keysLoading" class="loading-msg">{{ t('admin.loading') }}</div>
              <div v-else-if="!keys.length" class="empty-keys">{{ t('admin.keys_none') }}</div>
              <ul v-else class="keys-list">
                <li v-for="key in keys" :key="key.name" class="key-item">
                  <div class="key-info">
                    <span class="key-name mono">{{ key.name }}</span>
                    <span class="key-badge">{{ t('admin.keys_readonly') }}</span>
                  </div>
                  <button class="btn btn-danger btn-sm" @click="confirmDelete(key.name)">{{ t('admin.keys_revoke') }}</button>
                </li>
              </ul>
            </section>

            <!-- ── Preferences ─────────────────────────────── -->
            <section v-if="activeTab === 'prefs'" class="admin-section">
              <div class="section-header">
                <h2>{{ t('prefs.title') }}</h2>
                <p class="section-desc">{{ t('prefs.desc') }}</p>
              </div>
              <div class="config-form">
                <label class="field-label">{{ t('prefs.language') }}
                  <select v-model="language" class="field-input" style="max-width:220px">
                    <option value="en">English</option>
                    <option value="de">Deutsch</option>
                  </select>
                </label>
                <label class="field-label">{{ t('prefs.time_format') }}
                  <select v-model="timeFormat" class="field-input" style="max-width:260px">
                    <option value="24h">{{ t('prefs.time_format_24h') }}</option>
                    <option value="12h">{{ t('prefs.time_format_12h') }}</option>
                  </select>
                </label>
                <label class="field-label">{{ t('prefs.font_size') }}
                  <div class="radio-group">
                    <label class="radio-opt"><input type="radio" v-model="fontSize" value="small" /> {{ t('prefs.font_small') }}</label>
                    <label class="radio-opt"><input type="radio" v-model="fontSize" value="medium" /> {{ t('prefs.font_medium') }}</label>
                    <label class="radio-opt"><input type="radio" v-model="fontSize" value="large" /> {{ t('prefs.font_large') }}</label>
                  </div>
                </label>
                <label class="field-label">{{ t('prefs.auto_refresh') }}
                  <div class="inline-field">
                    <input v-model.number="autoRefreshIntervalPref" type="number" min="5" max="300"
                      class="field-input" style="max-width:100px" />
                    <span class="field-hint">{{ t('prefs.auto_refresh_unit') }}</span>
                  </div>
                </label>
                <label class="field-label">{{ t('prefs.default_time_range') }}
                  <select v-model="defaultTimeRange" class="field-input" style="max-width:120px">
                    <option value="15m">15m</option>
                    <option value="1h">1h</option>
                    <option value="6h">6h</option>
                    <option value="24h">24h</option>
                    <option value="7d">7d</option>
                    <option value="30d">30d</option>
                  </select>
                  <span class="field-hint">{{ t('prefs.default_time_range_hint') }}</span>
                </label>
                <p class="prefs-autosave-note">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 6L9 17l-5-5"/></svg>
                  {{ t('prefs.autosave_note') }}
                </p>
              </div>
            </section>

          </main>

          <!-- ── Bottom save bar (server + database tabs only) ── -->
          <div v-if="activeTab === 'server' || activeTab === 'database'" class="save-bar">
            <button
              class="btn btn-primary save-btn"
              :disabled="saving || configLoading"
              @click="saveCurrentTab"
            >
              {{ saving ? t('admin.saving') : t('admin.save') }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Delete key confirmation -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="deleteTarget" class="modal-backdrop" @click.self="deleteTarget = null">
          <div class="confirm-dialog">
            <h3>{{ t('admin.keys_revoke_title', { name: deleteTarget }) }}</h3>
            <p>{{ t('admin.keys_revoke_desc') }}</p>
            <div class="confirm-actions">
              <button class="btn btn-ghost" @click="deleteTarget = null">{{ t('admin.cancel') }}</button>
              <button class="btn btn-danger" :disabled="deleting" @click="deleteKey">
                {{ deleting ? t('admin.keys_revoking') : t('admin.keys_revoke') }}
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Restart confirmation -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="confirmRestart" class="modal-backdrop" @click.self="confirmRestart = false">
          <div class="confirm-dialog">
            <h3>{{ t('admin.restart_confirm_title') }}</h3>
            <p>{{ t('admin.restart_confirm_desc') }}</p>
            <div class="confirm-actions">
              <button class="btn btn-ghost" @click="confirmRestart = false" :disabled="restarting">
                {{ t('admin.cancel') }}
              </button>
              <button class="btn btn-primary" :disabled="restarting" @click="doRestart">
                {{ restarting ? t('admin.restarting') : t('admin.restart_btn') }}
              </button>
            </div>
            <p v-if="restarting" class="restart-status">{{ restartStatus }}</p>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>


  <AboutModal v-model="showAbout" />

</template>



<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { language, timeFormat, fontSize, autoRefreshInterval as autoRefreshIntervalPref, defaultTimeRange } from '@/stores/preferences'
import { useLocale } from '@/composables/useLocale'
import AppHeader from '@/components/AppHeader.vue'
import { api } from '@/api/client'
import AboutModal from '@/components/AboutModal.vue'

const { t } = useLocale()

const showAbout = ref(false)

const PREFS_SVG = '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>'

const tabs = computed(() => [
  { id: 'prefs',    label: t('admin.tab_preferences'), svg: PREFS_SVG },
  { id: 'server',   label: t('admin.tab_server'),   svg: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="3" width="20" height="5" rx="1"/><rect x="2" y="11" width="20" height="5" rx="1"/><rect x="2" y="19" width="20" height="5" rx="1"/></svg>' },
  { id: 'database', label: t('admin.tab_database'), svg: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><ellipse cx="12" cy="5" rx="9" ry="3"/><path d="M3 5v4c0 1.66 4.03 3 9 3s9-1.34 9-3V5"/><path d="M3 9v4c0 1.66 4.03 3 9 3s9-1.34 9-3V9"/><path d="M3 13v4c0 1.66 4.03 3 9 3s9-1.34 9-3v-4"/></svg>' },
  { id: 'keys',     label: t('admin.tab_keys'),     svg: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="7.5" cy="15.5" r="5.5"/><path d="M21 2l-9.6 9.6"/><path d="M15.5 7.5l3 3L22 7l-3-3"/></svg>' },
])
const activeTab = ref('prefs')

const configLoading = ref(false)
const saving        = ref(false)

const serverForm = reactive({
  host: '', port: 8000, originsStr: '*', autoRefreshInterval: 30, useSSL: false,
  defaultTimeRange: '24h', defaultLanguage: 'en', defaultFontSize: 'medium', defaultTimeFormat: '24h',
})
const dbForm      = reactive({ host: 'localhost', port: 3306, name: '', user: '', password: '' })
const cleanupForm = reactive({ enabled: false, diskPath: '/var/lib/mysql', thresholdPercent: 85, batchSize: 1000, intervalSeconds: 900 })

// ── SSL ───────────────────────────────────────────────────────────────────────
const sslGenerating  = ref(false)
const sslUploading   = ref(false)
const sslCertFile    = ref(null)
const sslKeyFile     = ref(null)
const sslMsg         = ref('')
const sslMsgOk       = ref(true)
const sslUploadMsg   = ref('')
const sslUploadMsgOk = ref(true)

// ── Disk ──────────────────────────────────────────────────────────────────────
const diskInfo    = ref(null)
const diskLoading = ref(false)
const diskError   = ref('')
const diskBarClass = computed(() => {
  if (!diskInfo.value) return ''
  const p = diskInfo.value.used_percent
  if (p >= 90) return 'critical'
  if (p >= 75) return 'warning'
  return 'ok'
})
function formatBytes(bytes) {
  if (bytes >= 1e12) return (bytes / 1e12).toFixed(1) + ' TB'
  if (bytes >= 1e9)  return (bytes / 1e9).toFixed(1) + ' GB'
  if (bytes >= 1e6)  return (bytes / 1e6).toFixed(1) + ' MB'
  return (bytes / 1e3).toFixed(0) + ' KB'
}
async function loadDiskUsage() {
  diskLoading.value = true; diskError.value = ''
  try { diskInfo.value = await api.getDiskUsage() }
  catch (e) { diskError.value = e.message || 'unknown'; diskInfo.value = null }
  finally { diskLoading.value = false }
}

// ── Load config ───────────────────────────────────────────────────────────────
async function loadConfig() {
  configLoading.value = true
  try {
    const cfg = await api.getConfig()
    const s = cfg.server ?? {}
    serverForm.host                = s.host ?? '0.0.0.0'
    serverForm.port                = s.port ?? 8000
    serverForm.originsStr          = (s.allowed_origins ?? ['*']).join(', ')
    serverForm.autoRefreshInterval = s.auto_refresh_interval ?? 30
    serverForm.useSSL              = s.use_ssl ?? false
    serverForm.defaultTimeRange    = s.default_time_range ?? '24h'
    serverForm.defaultLanguage     = s.default_language ?? 'en'
    serverForm.defaultFontSize     = s.default_font_size ?? 'medium'
    serverForm.defaultTimeFormat   = s.default_time_format ?? '24h'

    const d = cfg.database ?? {}
    dbForm.host     = d.host ?? 'localhost'
    dbForm.port     = d.port ?? 3306
    dbForm.name     = d.name ?? ''
    dbForm.user     = d.user ?? ''
    dbForm.password = ''

    const c = cfg.cleanup ?? {}
    cleanupForm.enabled          = c.enabled ?? false
    cleanupForm.diskPath         = c.disk_path ?? '/var/lib/mysql'
    cleanupForm.thresholdPercent = c.threshold_percent ?? 85
    cleanupForm.batchSize        = c.batch_size ?? 1000
    cleanupForm.intervalSeconds  = c.interval_seconds ?? 900
  } catch (e) {
    console.error('Failed to load config:', e)
  } finally { configLoading.value = false }
}

// ── Save ──────────────────────────────────────────────────────────────────────
async function saveCurrentTab() {
  if (activeTab.value === 'server')   return saveServer()
  if (activeTab.value === 'database') return saveDatabaseAndCleanup()
}

async function saveServer() {
  saving.value = true
  try {
    await api.updateConfig({ server: {
      host:                 serverForm.host,
      port:                 serverForm.port,
      allowed_origins:      serverForm.originsStr.split(',').map(s => s.trim()).filter(Boolean),
      auto_refresh_interval: serverForm.autoRefreshInterval,
      use_ssl:              serverForm.useSSL,
      default_time_range:   serverForm.defaultTimeRange,
      default_language:     serverForm.defaultLanguage,
      default_font_size:    serverForm.defaultFontSize,
      default_time_format:  serverForm.defaultTimeFormat,
    }})
    restartNeeded.value = true
  } catch (e) { alert(e.message || t('admin.save_failed')) }
  finally { saving.value = false }
}

async function saveDatabaseAndCleanup() {
  saving.value = true
  try {
    const dbPatch = { host: dbForm.host, port: dbForm.port, name: dbForm.name, user: dbForm.user }
    if (dbForm.password) dbPatch.password = dbForm.password
    await api.updateConfig({
      database: dbPatch,
      cleanup: {
        enabled:           cleanupForm.enabled,
        disk_path:         cleanupForm.diskPath,
        threshold_percent: cleanupForm.thresholdPercent,
        batch_size:        cleanupForm.batchSize,
        interval_seconds:  cleanupForm.intervalSeconds,
      },
    })
    dbForm.password = ''
    restartNeeded.value = true
  } catch (e) { alert(e.message || t('admin.save_failed')) }
  finally { saving.value = false }
}

// ── SSL actions ───────────────────────────────────────────────────────────────
async function generateSSL() {
  sslGenerating.value = true; sslMsg.value = ''
  try {
    await api.generateSSL()
    sslMsg.value = t('admin.ssl_generated'); sslMsgOk.value = true
    setTimeout(() => { sslMsg.value = '' }, 4000)
  } catch (e) { sslMsg.value = e.message || t('admin.ssl_generate_failed'); sslMsgOk.value = false }
  finally { sslGenerating.value = false }
}
async function uploadSSL() {
  if (!sslCertFile.value || !sslKeyFile.value) return
  sslUploading.value = true; sslUploadMsg.value = ''
  try {
    await api.uploadSSL(sslCertFile.value, sslKeyFile.value)
    sslCertFile.value = null; sslKeyFile.value = null
    sslUploadMsg.value = t('admin.ssl_uploaded'); sslUploadMsgOk.value = true
    setTimeout(() => { sslUploadMsg.value = '' }, 4000)
  } catch (e) { sslUploadMsg.value = e.message || t('admin.ssl_upload_failed'); sslUploadMsgOk.value = false }
  finally { sslUploading.value = false }
}

// ── Restart ───────────────────────────────────────────────────────────────────
const confirmRestart = ref(false)
const restarting     = ref(false)
const restartStatus  = ref('')
const restartNeeded  = ref(false)

async function doRestart() {
  restarting.value = true
  restartStatus.value = t('admin.restarting_wait')
  const token = sessionStorage.getItem('rsyslox_token')
  try {
    const ctrl = new AbortController()
    setTimeout(() => ctrl.abort(), 2000)
    await fetch('/api/admin/restart', {
      method: 'POST', headers: { 'X-Session-Token': token ?? '' }, signal: ctrl.signal,
    })
  } catch {}
  restartStatus.value = t('admin.restarting_poll')
  const start = Date.now()
  const poll = async () => {
    if (Date.now() - start > 30000) { restartStatus.value = t('admin.restart_timeout'); restarting.value = false; return }
    try { const res = await fetch('/health'); if (res.ok) { restartNeeded.value = false; window.location.reload(); return } } catch {}
    setTimeout(poll, 1000)
  }
  setTimeout(poll, 1500)
}

// ── Keys ──────────────────────────────────────────────────────────────────────
const keys = ref([]); const keysLoading = ref(false)
const newKeyName = ref(''); const keyCreating = ref(false)
const newKeyPlaintext = ref(null); const newKeyRevealName = ref('')
const keyCopied = ref(false)
const deleteTarget = ref(null); const deleting = ref(false)

async function loadKeys() {
  keysLoading.value = true
  try { keys.value = await api.getKeys() } catch {}
  finally { keysLoading.value = false }
}
async function createKey() {
  if (!newKeyName.value.trim()) return
  keyCreating.value = true
  try {
    const res = await api.createKey(newKeyName.value.trim())
    newKeyPlaintext.value = res.key; newKeyRevealName.value = res.name; newKeyName.value = ''
    await loadKeys()
  } catch (e) { alert(e.message || 'Failed') }
  finally { keyCreating.value = false }
}
async function copyKey() {
  try { await navigator.clipboard.writeText(newKeyPlaintext.value); keyCopied.value = true; setTimeout(() => { keyCopied.value = false }, 2000) } catch {}
}
function confirmDelete(name) { deleteTarget.value = name }
async function deleteKey() {
  deleting.value = true
  try { await api.deleteKey(deleteTarget.value); deleteTarget.value = null; await loadKeys() }
  catch (e) { alert(e.message || 'Failed') }
  finally { deleting.value = false }
}

onMounted(() => { loadConfig(); loadKeys(); loadDiskUsage() })
</script>

<style scoped>
.admin-layout { display: flex; flex-direction: column; height: 100%; overflow: hidden; }
.admin-body   { display: flex; flex-direction: column; flex: 1; overflow: hidden; }
.admin-content { display: flex; flex: 1; overflow: hidden; }

.restart-banner {
  display: flex; align-items: center; gap: .75rem;
  padding: .6rem 1.25rem;
  background: #fefce8; border-bottom: 1px solid #fde68a;
  color: #92400e; flex-shrink: 0; font-size: .875rem;
}
[data-theme="dark"] .restart-banner { background: #1c1404; border-color: #78350f; color: #fcd34d; }
.banner-icon { font-size: 1rem; flex-shrink: 0; }
.banner-text { flex: 1; }
.banner-actions { display: flex; gap: .5rem; flex-shrink: 0; }
.banner-enter-active, .banner-leave-active { transition: max-height .25s ease, opacity .2s; max-height: 60px; overflow: hidden; }
.banner-enter-from, .banner-leave-to { max-height: 0; opacity: 0; }

.admin-nav {
  width: 180px; flex-shrink: 0;
  background: var(--bg-surface); border-right: 1px solid var(--border);
  padding: .75rem .5rem;
  display: flex; flex-direction: column; gap: .25rem; overflow-y: auto;
  justify-content: space-between;
}
.nav-tab {
  display: flex; align-items: center; gap: .5rem;
  padding: .5rem .75rem; border: none; border-radius: var(--radius);
  background: none; cursor: pointer;
  font-size: .875rem; color: var(--text-muted); text-align: left;
  transition: background .15s, color .15s; width: 100%;
  font-weight: 600;
}
.nav-tab:hover  { background: var(--bg-hover); color: var(--text); }
.nav-tab.active { background: var(--bg-selected); color: var(--color-primary); }
.tab-icon { display: flex; flex-shrink: 0; }

.admin-main-wrapper { flex: 1; display: flex; flex-direction: column; overflow: hidden; }
.admin-main { flex: 1; overflow-y: auto; padding: 1.5rem; background: var(--bg); color: var(--text); }
.admin-section { max-width: 640px; display: flex; flex-direction: column; gap: 1.25rem; }

/* ── Bottom save bar ──────────────────── */
.save-bar {
  display: flex; align-items: center; justify-content: flex-end;
  padding: .625rem 1.5rem;
  background: var(--bg-surface); border-top: 1px solid var(--border);
  flex-shrink: 0;
}
/* No min-width on save button — it sizes to its text only */
.save-btn { min-width: unset; }

.section-header { display: flex; flex-direction: column; gap: .25rem; }
.section-header h2 { font-size: 1.125rem; font-weight: 700; }
.section-desc { font-size: .875rem; color: var(--text-muted); }
.subsection-header { display: flex; flex-direction: column; gap: .2rem; }
.subsection-header h3 { font-size: 1rem; font-weight: 600; }

.info-callout {
  display: flex; align-items: flex-start; gap: .5rem;
  padding: .625rem .875rem;
  background: var(--bg-surface); border: 1px solid var(--border);
  border-left: 3px solid var(--color-primary);
  border-radius: var(--radius);
  font-size: .8rem; color: var(--text-muted); line-height: 1.5;
}
.loading-msg { color: var(--text-muted); font-size: .875rem; }

.config-form {
  display: flex; flex-direction: column; gap: 1rem;
  background: var(--bg-surface); border: 1px solid var(--border);
  border-radius: var(--radius); padding: 1.25rem;
}
.field-row { display: grid; grid-template-columns: 1fr 1fr; gap: .75rem; }
.field-label {
  display: flex; flex-direction: column; gap: .375rem;
  font-size: .875rem; font-weight: 500; color: var(--text);
}
.field-input {
  padding: .4rem .625rem; border: 1px solid var(--border); border-radius: var(--radius);
  background: var(--bg); color: var(--text); font-size: .875rem;
  width: 100%; box-sizing: border-box; transition: border-color .15s;
}
.field-input:focus { outline: 2px solid var(--color-primary); outline-offset: -1px; border-color: var(--color-primary); }
.field-input:disabled { opacity: .5; cursor: not-allowed; background: var(--bg-hover); }
.file-input { padding: .28rem .4rem; cursor: pointer; }
.field-hint { font-size: .75rem; color: var(--text-muted); font-weight: 400; }
.restart-hint { color: var(--color-accent, #f59e0b); }
.toggle-label {
  display: flex; align-items: center; gap: .5rem;
  font-size: .875rem; font-weight: 500; cursor: pointer; color: var(--text);
}
.toggle-label input { accent-color: var(--color-primary); width: 16px; height: 16px; cursor: pointer; }
.inline-field { display: flex; align-items: center; gap: .5rem; }
.disabled { opacity: .5; pointer-events: none; }
.inline-actions { display: flex; align-items: center; gap: .75rem; }
.action-msg { font-size: .875rem; }
.action-msg.ok  { color: #16a34a; }
.action-msg.err { color: #dc2626; }

.ssl-block {
  display: flex; flex-direction: column; gap: .625rem;
  padding-top: .75rem; border-top: 1px solid var(--border);
}

.disk-widget {
  display: flex; flex-direction: column; gap: .5rem;
  padding: .875rem 1rem; background: var(--bg);
  border: 1px solid var(--border); border-radius: var(--radius);
}
.disk-widget-header { display: flex; align-items: center; gap: .5rem; }
.disk-path-label {
  flex: 1; font-size: .78rem; color: var(--text-muted);
  font-family: ui-monospace, monospace; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
}
.disk-loading { font-size: .8rem; color: var(--text-muted); }
.disk-error   { font-size: .8rem; color: #dc2626; }
.disk-bar-track { height: 8px; border-radius: 999px; background: var(--bg-hover); overflow: hidden; }
.disk-bar-fill  { height: 100%; border-radius: 999px; transition: width .4s ease; }
.disk-bar-fill.ok       { background: #16a34a; }
.disk-bar-fill.warning  { background: #d97706; }
.disk-bar-fill.critical { background: #dc2626; }
.disk-bar-labels { display: flex; justify-content: space-between; font-size: .75rem; color: var(--text-muted); }

.btn-sm { font-size: .8rem; padding: .3rem .625rem; }
.new-key-form { display: flex; gap: .625rem; align-items: center; flex-wrap: wrap; }
.key-reveal {
  background: #f0fdf4; border: 1px solid #bbf7d0;
  border-radius: var(--radius); padding: 1rem;
  display: flex; flex-direction: column; gap: .75rem;
}
[data-theme="dark"] .key-reveal { background: #052e16; border-color: #166534; }
.key-reveal-header { display: flex; align-items: center; gap: .5rem; flex-wrap: wrap; font-size: .875rem; }
.key-reveal-warn   { font-size: .8rem; color: #16a34a; margin-left: auto; }
[data-theme="dark"] .key-reveal-warn { color: #4ade80; }
.key-reveal-value {
  display: flex; align-items: center; gap: .625rem;
  background: var(--bg); border: 1px solid var(--border);
  border-radius: var(--radius); padding: .5rem .75rem; overflow: hidden;
}
.key-reveal-value code { flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; font-size: .8rem; }
.empty-keys { font-size: .875rem; color: var(--text-muted); }
.keys-list { list-style: none; background: var(--bg-surface); border: 1px solid var(--border); border-radius: var(--radius); overflow: hidden; }
.key-item  { display: flex; align-items: center; justify-content: space-between; padding: .75rem 1rem; border-bottom: 1px solid var(--border); gap: .75rem; }
.key-item:last-child { border-bottom: none; }
.key-info  { display: flex; align-items: center; gap: .75rem; }
.key-name  { font-size: .875rem; font-family: ui-monospace, monospace; }
.key-badge { font-size: .7rem; padding: .1rem .375rem; background: var(--bg-hover); border: 1px solid var(--border); border-radius: 999px; color: var(--text-muted); }

.radio-group { display: flex; gap: 1rem; margin-top: .25rem; flex-wrap: wrap; }
.radio-opt   { display: flex; align-items: center; gap: .35rem; font-size: .875rem; cursor: pointer; }
.prefs-autosave-note { display: flex; align-items: center; gap: .4rem; font-size: .78rem; color: var(--text-muted); padding-top: .25rem; }

.modal-backdrop {
  position: fixed; inset: 0; background: rgba(0,0,0,.4);
  z-index: 300; display: flex; align-items: center; justify-content: center; padding: 1rem;
}
.confirm-dialog {
  background: var(--bg-surface); border: 1px solid var(--border);
  border-radius: var(--radius); padding: 1.5rem;
  max-width: 360px; width: 100%;
  display: flex; flex-direction: column; gap: .75rem;
  box-shadow: 0 20px 40px rgba(0,0,0,.2);
}
.confirm-dialog h3 { font-size: 1rem; font-weight: 700; }
.confirm-dialog p  { font-size: .875rem; color: var(--text-muted); }
.confirm-actions   { display: flex; gap: .625rem; justify-content: flex-end; }
.restart-status    { font-size: .8rem; color: var(--text-muted); text-align: center; }
.modal-enter-active, .modal-leave-active { transition: opacity .2s; }
.modal-enter-from,  .modal-leave-to      { opacity: 0; }

@media (max-width: 600px) {
  .admin-nav { width: 100%; flex-direction: row; border-right: none; border-bottom: 1px solid var(--border); overflow-x: auto; padding: .5rem; }
  .nav-tab   { flex-shrink: 0; }
  .admin-content { flex-direction: column; }
  .field-row  { grid-template-columns: 1fr; }
}
.no-gap { display: inline-flex; align-items: center; flex-wrap: wrap; }
.no-gap svg { flex: 0 0 auto; }
.no-gap { gap: 0; } /* optional, sorgt für keinen Abstand */
</style>
