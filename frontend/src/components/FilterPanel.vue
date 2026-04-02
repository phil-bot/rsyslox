<template>
  <aside class="filter-panel" :class="{ open: modelValue }">

    <!-- ── Header ────────────────────────────────── -->
    <div class="panel-header">
      <span class="panel-title">{{ t('filter.title') }}</span>
      <div class="panel-header-actions">
        <button class="reset-btn" @click="$emit('reset')">{{ t('filter.reset') }}</button>
        <!-- Close button lives inside the sidebar when it is open -->
        <button class="close-btn" @click="$emit('close')" :title="t('filter.close')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
          <span class="close-label">{{ t('filter.close') }}</span>
        </button>
      </div>
    </div>

    <!-- ── Time Range ─────────────────────────────── -->
    <section class="filter-section fixed-section">
      <button class="section-toggle" @click="toggle('timeRange')">
        <span class="section-label">{{ t('filter.time_range') }}</span>
        <svg class="chevron" :class="{ rotated: collapsed.timeRange }" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
          <polyline points="6 9 12 15 18 9"/>
        </svg>
      </button>
      <div v-show="!collapsed.timeRange" class="section-body">
        <div class="dur-seg">
          <button
            v-for="d in durations" :key="d.value"
            class="dur-btn"
            :class="{ active: activeDur === d.value }"
            @click="selectDuration(d.value)"
          >{{ d.label }}</button>
        </div>
        <div class="date-row">
          <span class="date-lbl">{{ t('filter.from') }}</span>
          <input type="datetime-local" class="date-input" :value="startDate"
            @input="onDateInput('start', $event.target.value)" />
        </div>
        <div class="date-row">
          <span class="date-lbl">{{ t('filter.to') }}</span>
          <input
            type="datetime-local" class="date-input"
            :class="{ 'live-input': autoRefresh }"
            :value="autoRefresh ? liveNow : endDate"
            :disabled="autoRefresh"
            @input="onDateInput('end', $event.target.value)"
          />
        </div>
        <div class="shift-row">
          <button class="shift-btn" @click="onEarlier">{{ t('filter.earlier') }}</button>
          <button
            class="shift-btn live-btn" :class="{ active: autoRefresh }"
            @click="$emit('toggle-live')"
            :title="autoRefresh ? t('filter.live_title_on') : t('filter.live_title_off')"
          >
            <span class="live-dot" :class="{ pulse: autoRefresh }"></span>
            {{ t('filter.live') }}
            <span v-if="autoRefresh && countdown > 0" class="live-countdown">{{ countdown }}s</span>
          </button>
          <button class="shift-btn" :disabled="laterDisabled" @click="$emit('shift', 1)">{{ t('filter.later') }}</button>
        </div>
      </div>
    </section>

    <!-- ── Severity ───────────────────────────────── -->
    <section class="filter-section fixed-section">
      <button class="section-toggle" @click="toggle('severity')">
        <span class="section-label">{{ t('filter.severity') }}</span>
        <svg class="chevron" :class="{ rotated: collapsed.severity }" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
          <polyline points="6 9 12 15 18 9"/>
        </svg>
      </button>
      <div v-show="!collapsed.severity" class="section-body">
        <div v-if="availableSeverities.length" class="sev-row">
          <button
            v-for="item in availableSeverities" :key="item.val"
            class="sev-badge-btn"
            :class="sevBadgeClass(item.val)"
            @click="cycle('severities', 'excludeSeverities', item.val)"
            :title="pillTitle('sev', item.val)"
          >{{ SEV_LABEL_FULL[item.val] ?? item.label }}</button>
        </div>
        <p v-else class="empty-hint">{{ t('filter.loading') }}</p>
      </div>
    </section>

    <!-- ── Facility ───────────────────────────────── -->
    <section class="filter-section fixed-section">
      <button class="section-toggle" @click="toggle('facility')">
        <span class="section-label">{{ t('filter.facility') }}</span>
        <svg class="chevron" :class="{ rotated: collapsed.facility }" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
          <polyline points="6 9 12 15 18 9"/>
        </svg>
      </button>
      <div v-show="!collapsed.facility" class="section-body">
        <div v-if="availableFacilities.length" class="sev-row">
          <button
            v-for="item in availableFacilities" :key="item.val"
            class="fac-badge-btn"
            :class="facBadgeClass(item.val)"
            @click="cycle('facilities', 'excludeFacilities', item.val)"
            :title="pillTitle('fac', item.val)"
          >{{ item.label ?? item.val }}</button>
        </div>
        <p v-else class="empty-hint">{{ t('filter.loading') }}</p>
      </div>
    </section>

    <!-- ── Tag ────────────────────────────────────── -->
    <section class="filter-section fixed-section">
      <button class="section-toggle" @click="toggle('tag')">
        <span class="section-label">{{ t('filter.tag') }}</span>
        <svg class="chevron" :class="{ rotated: collapsed.tag }" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
          <polyline points="6 9 12 15 18 9"/>
        </svg>
      </button>
      <div v-show="!collapsed.tag" class="section-body">
        <div v-if="availableTags.length" class="sev-row">
          <button
            v-for="tag in availableTags" :key="tag"
            class="fac-badge-btn"
            :class="tagBadgeClass(tag)"
            @click="cycle('tags', 'excludeTags', tag)"
            :title="pillTitle('tag', tag)"
          >{{ tag }}</button>
        </div>
        <p v-else class="empty-hint">{{ t('filter.no_tags') }}</p>
      </div>
    </section>

    <!-- ── Host ──────────────────────────────────── -->
    <section class="filter-section flex-section" :style="collapsed.host ? {} : flexStyle(availableHosts.length)">
      <button class="section-toggle" @click="toggle('host')">
        <span class="section-label">{{ t('filter.host') }}</span>
        <svg class="chevron" :class="{ rotated: collapsed.host }" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
          <polyline points="6 9 12 15 18 9"/>
        </svg>
      </button>
      <div v-show="!collapsed.host" class="section-body host-body">
        <template v-if="availableHosts.length">
          <input class="list-search" type="text"
            :placeholder="t('filter.search_hosts')" v-model="hostSearch" />
          <div class="filter-list">
            <button
              v-for="host in filteredHosts" :key="host"
              class="list-item"
              :class="listClass('host', host)"
              @click="cycle('hosts', 'excludeHosts', host)"
            >
              <span class="list-state" :class="listStateClass('host', host)">
                {{ listIcon('host', host) }}
              </span>
              <span class="list-label mono-text">{{ host }}</span>
            </button>
            <p v-if="!filteredHosts.length" class="empty-hint pad">{{ t('filter.no_match') }}</p>
          </div>
        </template>
        <p v-else class="empty-hint">{{ t('filter.no_hosts') }}</p>
      </div>
    </section>

    <!-- ── Message Search ─────────────────────────── -->
    <section class="filter-section fixed-section message-section">
      <button class="section-toggle" @click="toggle('message')">
        <span class="section-label">{{ t('filter.message_search') }}</span>
        <svg class="chevron" :class="{ rotated: collapsed.message }" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
          <polyline points="6 9 12 15 18 9"/>
        </svg>
      </button>
      <div v-show="!collapsed.message" class="section-body">
        <div class="search-wrap">
          <input
            class="search-input" type="text"
            :placeholder="t('filter.search_placeholder')"
            :value="messageSearch"
            @input="$emit('update:messageSearch', $event.target.value)"
          />
          <button v-if="messageSearch" class="search-clear"
            @click="$emit('update:messageSearch', '')" title="Clear search" type="button">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
              <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
      </div>
    </section>

  </aside>
</template>

<script setup>
import { computed, onBeforeUnmount, reactive, ref, watch } from 'vue'
import { useLocale } from '@/composables/useLocale'

const props = defineProps({
  modelValue:           { type: Boolean, default: true },
  timeMode:             { type: String,  default: 'relative' },
  relativeDur:          { type: String,  default: '1h' },
  startDate:            { type: String,  default: '' },
  endDate:              { type: String,  default: '' },
  severities:           { type: Array,   default: () => [] },
  excludeSeverities:    { type: Array,   default: () => [] },
  facilities:           { type: Array,   default: () => [] },
  excludeFacilities:    { type: Array,   default: () => [] },
  hosts:                { type: Array,   default: () => [] },
  excludeHosts:         { type: Array,   default: () => [] },
  tags:                 { type: Array,   default: () => [] },
  excludeTags:          { type: Array,   default: () => [] },
  messageSearch:        { type: String,  default: '' },
  availableHosts:       { type: Array,   default: () => [] },
  availableTags:        { type: Array,   default: () => [] },
  availableSeverities:  { type: Array,   default: () => [] },
  availableFacilities:  { type: Array,   default: () => [] },
  autoRefresh:          { type: Boolean, default: false },
  countdown:            { type: Number,  default: 0 },
})

const { t } = useLocale()
const emit = defineEmits([
  'update:timeMode','update:relativeDur','update:startDate','update:endDate',
  'update:severities','update:excludeSeverities',
  'update:facilities','update:excludeFacilities',
  'update:hosts','update:excludeHosts',
  'update:tags','update:excludeTags',
  'update:messageSearch','shift','reset','toggle-live','exit-live','close',
])

// ── Collapsible sections ──────────────────────────────────────────────────────
const STORAGE_KEY = 'rsyslox_filter_collapsed'

function loadCollapsed() {
  try { return JSON.parse(localStorage.getItem(STORAGE_KEY) || '{}') } catch { return {} }
}

const stored = loadCollapsed()
const collapsed = reactive({
  timeRange: stored.timeRange ?? false,
  severity:  stored.severity  ?? false,
  facility:  stored.facility  ?? true,
  tag:       stored.tag       ?? true,
  host:      stored.host      ?? false,
  message:   stored.message   ?? false,
})

function toggle(key) {
  collapsed[key] = !collapsed[key]
  localStorage.setItem(STORAGE_KEY, JSON.stringify({ ...collapsed }))
}

// ── Search inputs ─────────────────────────────────────────────────────────────
const hostSearch = ref('')

const filteredHosts = computed(() =>
  hostSearch.value
    ? props.availableHosts.filter(h => h.toLowerCase().includes(hostSearch.value.toLowerCase()))
    : props.availableHosts
)

// ── 3-state cycle ─────────────────────────────────────────────────────────────
function getState(includeArr, excludeArr, val) {
  if (includeArr.includes(val)) return 'include'
  if (excludeArr.includes(val)) return 'exclude'
  return 'neutral'
}

function cycle(includeKey, excludeKey, val) {
  const state   = getState(props[includeKey], props[excludeKey], val)
  const incCopy = [...props[includeKey]]
  const excCopy = [...props[excludeKey]]
  if (state === 'neutral') {
    incCopy.push(val)
  } else if (state === 'include') {
    incCopy.splice(incCopy.indexOf(val), 1)
    excCopy.push(val)
  } else {
    excCopy.splice(excCopy.indexOf(val), 1)
  }
  emit('update:' + includeKey, incCopy)
  emit('update:' + excludeKey, excCopy)
}

// ── Severity / Facility / Tag badge display ───────────────────────────────────
const SEV_LABEL_FULL = {
  0: 'Emergency', 1: 'Alert', 2: 'Critical', 3: 'Error',
  4: 'Warning', 5: 'Notice', 6: 'Info', 7: 'Debug',
}

function sevBadgeClass(val) {
  const state = getState(props.severities, props.excludeSeverities, val)
  return ['sev-' + val, 'sev-' + state]
}
function facBadgeClass(val) {
  return ['fac-' + getState(props.facilities, props.excludeFacilities, val)]
}
function tagBadgeClass(val) {
  return ['fac-' + getState(props.tags, props.excludeTags, val)]
}

function pillTitle(type, val) {
  const incArr = type === 'sev' ? props.severities : type === 'fac' ? props.facilities : props.tags
  const excArr = type === 'sev' ? props.excludeSeverities : type === 'fac' ? props.excludeFacilities : props.excludeTags
  const state  = getState(incArr, excArr, val)
  if (state === 'neutral') return t('filter.click_include')
  if (state === 'include') return t('filter.click_exclude')
  return t('filter.click_reset')
}

function getListArrays(type) {
  if (type === 'fac')  return [props.facilities,  props.excludeFacilities]
  if (type === 'host') return [props.hosts,        props.excludeHosts]
                       return [props.tags,          props.excludeTags]
}
function listClass(type, val) {
  const [inc, exc] = getListArrays(type)
  const s = getState(inc, exc, val)
  return s === 'include' ? 'list-include' : s === 'exclude' ? 'list-exclude' : ''
}
function listStateClass(type, val) {
  const [inc, exc] = getListArrays(type)
  const s = getState(inc, exc, val)
  return s === 'include' ? 'state-inc' : s === 'exclude' ? 'state-exc' : 'state-neu'
}
function listIcon(type, val) {
  const [inc, exc] = getListArrays(type)
  const s = getState(inc, exc, val)
  return s === 'include' ? '+' : s === 'exclude' ? '−' : '·'
}

// ── Duration control ──────────────────────────────────────────────────────────
const durations = [
  { value: '15m', label: '15m' }, { value: '1h',  label: '1h'  },
  { value: '6h',  label: '6h'  }, { value: '24h', label: '24h' },
  { value: '7d',  label: '7d'  }, { value: '30d', label: '30d' },
]
const DURATION_MS = {
  '15m': 15*60*1000, '1h': 60*60*1000, '6h': 6*60*60*1000,
  '24h': 24*60*60*1000, '7d': 7*24*60*60*1000, '30d': 30*24*60*60*1000,
}

const activeDur = computed(() => {
  if (props.autoRefresh) return props.relativeDur
  if (!props.startDate || !props.endDate) return props.relativeDur
  // Die Strings kommen von toISOString() und sind UTC ("YYYY-MM-DDTHH:mm").
  // Chrome parst diese ohne 'Z' als Lokalzeit. Überspannt das Fenster eine
  // DST-Grenze, ist der Offset zwischen Start und End unterschiedlich → diff
  // ist um 1h falsch → kein Preset-Match. Das 'Z' erzwingt UTC in allen Browsern.
  const diff  = new Date(props.endDate + 'Z') - new Date(props.startDate + 'Z')
  const match = Object.entries(DURATION_MS).find(([, ms]) => Math.abs(diff - ms) < 60000)
  return match ? match[0] : null
})

function selectDuration(val) {
  const now  = new Date()
  const from = new Date(now - DURATION_MS[val])
  const fmt  = d => d.toISOString().slice(0, 16)
  emit('update:startDate', fmt(from))
  if (!props.autoRefresh) emit('update:endDate', fmt(now))
  emit('update:timeMode',    'absolute')
  emit('update:relativeDur', val)
}

const laterDisabled = computed(() => {
  if (props.autoRefresh) return true
  if (!props.endDate) return false
  return new Date(props.endDate) >= new Date()
})

function onEarlier() {
  if (props.autoRefresh) emit('exit-live')
  emit('shift', -1)
}
function onDateInput(which, value) {
  if (which === 'end' && props.autoRefresh) emit('exit-live')
  emit(which === 'start' ? 'update:startDate' : 'update:endDate', value)
  emit('update:timeMode', 'absolute')
}

// ── Host section height ───────────────────────────────────────────────────────
const ROW_H    = 29
const CHROME_H = 80
function flexStyle(count) {
  const maxH = Math.max(count, 1) * ROW_H + CHROME_H
  return { flex: '0 1 auto', maxHeight: maxH + 'px', minHeight: 0 }
}

// ── Live "now" ticker ─────────────────────────────────────────────────────────
const liveNow = ref('')
let liveTimer = null
function updateLiveNow() { liveNow.value = new Date().toISOString().slice(0, 16) }
watch(() => props.autoRefresh, (val) => {
  if (val) { updateLiveNow(); liveTimer = setInterval(updateLiveNow, 1000) }
  else     { clearInterval(liveTimer); liveTimer = null }
}, { immediate: true })
onBeforeUnmount(() => clearInterval(liveTimer))
</script>

<style scoped>
/* ── Panel shell ─────────────────────────────────── */
.filter-panel {
  width: var(--sidebar-width);
  flex-shrink: 0;
  background: var(--bg-surface);
  border-right: 1px solid var(--border);
  display: flex; flex-direction: column;
  height: 100%; overflow: hidden;
  transition: width .2s;
}
@media (max-width: 768px) {
  .filter-panel {
    position: fixed; top: var(--header-height); left: 0; bottom: 0;
    z-index: 50; transform: translateX(-100%);
    transition: transform .25s; box-shadow: 4px 0 16px rgba(0,0,0,.15);
  }
  .filter-panel.open { transform: translateX(0); }
}
@media (min-width: 769px) {
  .filter-panel:not(.open) { width: 0; border: none; }
}

/* ── Panel header ────────────────────────────────── */
.panel-header {
  display: flex; align-items: center; justify-content: space-between;
  padding: .4rem .5rem .4rem .75rem;
  border-bottom: 1px solid var(--border);
  min-height: 40px; flex-shrink: 0;
}
.panel-title { font-weight: 600; font-size: .875rem; color: var(--text); }
.panel-header-actions { display: flex; align-items: center; gap: .25rem; }

.reset-btn {
  background: none; border: none; cursor: pointer;
  color: var(--color-primary); font-size: .78rem;
  padding: .2rem .45rem; border-radius: var(--radius);
}
.reset-btn:hover { background: var(--bg-hover); }

.close-btn {
  display: flex; align-items: center; gap: .3rem;
  background: none; border: 1px solid var(--border);
  border-radius: var(--radius); cursor: pointer;
  color: var(--text-muted); font-size: .78rem;
  padding: .2rem .5rem;
  transition: background .15s, color .15s, border-color .15s;
}
.close-btn:hover { background: var(--bg-hover); color: var(--text); border-color: var(--text-muted); }
.close-label { white-space: nowrap; }

/* ── Section toggle header ───────────────────────── */
.filter-section {
  border-bottom: 1px solid var(--border);
  display: flex; flex-direction: column;
}
.fixed-section  { flex-shrink: 0; }
.flex-section   { flex: 1; min-height: 0; }
.message-section { flex-shrink: 0; border-bottom: none; }

.section-toggle {
  display: flex; align-items: center; justify-content: space-between;
  width: 100%; padding: .45rem .875rem;
  background: none; border: none; cursor: pointer;
  text-align: left; gap: .5rem;
  transition: background .12s;
}
.section-toggle:hover { background: var(--bg-hover); }

.section-label {
  font-size: .68rem; font-weight: 700;
  text-transform: uppercase; letter-spacing: .07em;
  color: var(--text-muted);
}

.chevron {
  color: var(--text-muted); flex-shrink: 0;
  transition: transform .2s ease;
}
.chevron.rotated { transform: rotate(-90deg); }

/* Section body holds the actual content with padding */
.section-body {
  display: flex; flex-direction: column; gap: .4rem;
  padding: 0 .875rem .6rem;
}
/* Host section body needs to participate in flex scroll */
.host-body {
  flex: 1; min-height: 0; overflow: hidden;
}

/* ── Facility badge ──────────────────────────────── */
.fac-badge-btn {
  display: inline-flex; align-items: center;
  padding: .15rem .45rem; border-radius: 3px;
  font-size: .7rem; font-weight: 600; letter-spacing: .02em;
  cursor: pointer; background: var(--bg); color: var(--text-muted);
  border: 1px solid var(--border);
  transition: background .12s, color .12s, border-color .12s;
  user-select: none;
}
.fac-badge-btn:active { transform: scale(.95); }
.fac-badge-btn.fac-neutral:hover { background: var(--bg-hover); color: var(--text); }
.fac-badge-btn.fac-include { background: var(--color-primary); color: #fff; border-color: var(--color-primary); }
.fac-badge-btn.fac-exclude { color: #dc2626; border-color: #dc2626; text-decoration: line-through; }

/* ── Duration ────────────────────────────────────── */
.dur-seg {
  display: flex; width: 100%;
  background: var(--bg); border: 1px solid var(--border);
  border-radius: var(--radius); overflow: hidden; flex-shrink: 0;
}
.dur-btn {
  flex: 1; padding: .28rem .2rem;
  background: transparent; border: none;
  border-right: 1px solid var(--border);
  cursor: pointer; font-size: .78rem; color: var(--text-muted);
  transition: background .12s, color .12s; white-space: nowrap;
}
.dur-btn:last-child { border-right: none; }
.dur-btn:hover { background: var(--bg-hover); color: var(--text); }
.dur-btn.active { background: var(--color-primary); color: #fff; }

/* ── Date fields ─────────────────────────────────── */
.date-row { display: flex; align-items: center; gap: .5rem; flex-shrink: 0; }
.date-lbl { font-size: .78rem; color: var(--text-muted); flex-shrink: 0; width: 2rem; text-align: right; }
.date-input {
  flex: 1; padding: .28rem .4rem; border: 1px solid var(--border);
  border-radius: var(--radius); background: var(--bg); color: var(--text);
  font-size: .76rem; min-width: 0;
}
.date-input.live-input { color: var(--color-primary); border-color: var(--color-primary); opacity: .85; }

/* ── Shift / Live row ────────────────────────────── */
.shift-row { display: flex; gap: .375rem; flex-shrink: 0; }
.shift-btn {
  flex: 1; background: var(--bg); border: 1px solid var(--border);
  border-radius: var(--radius); cursor: pointer;
  padding: .25rem .4rem; font-size: .76rem; color: var(--text-muted);
  transition: background .15s, color .15s, border-color .15s; white-space: nowrap;
}
.shift-btn:hover:not(:disabled) { background: var(--bg-hover); color: var(--text); }
.shift-btn:disabled { opacity: .35; cursor: default; }
.live-btn { display: inline-flex; align-items: center; justify-content: center; gap: .3rem; }
.live-btn.active { border-color: var(--color-primary); color: var(--color-primary); background: var(--bg-selected); }
.live-btn.active:hover { background: var(--bg-selected); }
.live-dot { width: 6px; height: 6px; border-radius: 50%; background: var(--text-muted); flex-shrink: 0; transition: background .15s; }
.live-btn.active .live-dot { background: var(--color-primary); }
@keyframes live-pulse { 0%,100%{opacity:1;transform:scale(1)} 50%{opacity:.5;transform:scale(.75)} }
.live-dot.pulse { animation: live-pulse 1.5s ease-in-out infinite; }
.live-countdown { font-variant-numeric: tabular-nums; font-size: .72rem; opacity: .75; }

/* ── Severity badges ─────────────────────────────── */
.sev-row { display: flex; flex-wrap: wrap; gap: .3rem; flex-shrink: 0; }
.sev-badge-btn {
  display: inline-flex; align-items: center;
  padding: .15rem .4rem; border-radius: 3px;
  font-size: .7rem; font-weight: 700; letter-spacing: .02em;
  cursor: pointer; border: 2px solid transparent;
  transition: opacity .15s, border-color .15s, transform .1s;
  color: #fff; user-select: none;
}
.sev-badge-btn:active { transform: scale(.95); }
.sev-badge-btn.sev-neutral { opacity: .35; }
.sev-badge-btn.sev-neutral:hover { opacity: .7; }
.sev-badge-btn.sev-include { opacity: 1; border-color: rgba(255,255,255,.6); box-shadow: 0 0 0 1px rgba(255,255,255,.25); }
.sev-badge-btn.sev-exclude { opacity: .55; border-color: #dc2626; text-decoration: line-through; }
.sev-badge-btn.sev-0 { background: var(--sev-0); }
.sev-badge-btn.sev-1 { background: var(--sev-1); }
.sev-badge-btn.sev-2 { background: var(--sev-2); }
.sev-badge-btn.sev-3 { background: var(--sev-3); }
.sev-badge-btn.sev-4 { background: var(--sev-4); }
.sev-badge-btn.sev-5 { background: var(--sev-5); }
.sev-badge-btn.sev-6 { background: var(--sev-6); }
.sev-badge-btn.sev-7 { background: var(--sev-7); color: var(--text); }

/* ── Host list ───────────────────────────────────── */
.list-search {
  width: 100%; padding: .28rem .45rem; box-sizing: border-box;
  border: 1px solid var(--border); border-radius: var(--radius);
  background: var(--bg); color: var(--text); font-size: .76rem; flex-shrink: 0;
}
.list-search:focus { outline: 2px solid var(--color-primary); outline-offset: -1px; }
.filter-list {
  flex: 1; overflow-y: auto; min-height: 0;
  border: 1px solid var(--border); border-radius: var(--radius);
  background: var(--bg);
}
.list-item {
  display: flex; align-items: center; gap: .4rem;
  width: 100%; padding: .26rem .5rem;
  background: none; border: none; border-bottom: 1px solid var(--border);
  cursor: pointer; text-align: left; color: var(--text-muted);
  font-size: .78rem; transition: background .1s;
}
.list-item:last-child { border-bottom: none; }
.list-item:hover { background: var(--bg-hover); color: var(--text); }
.list-item.list-include { background: color-mix(in srgb, var(--color-primary) 10%, transparent); color: var(--text); }
.list-item.list-exclude { background: color-mix(in srgb, #dc2626 10%, transparent); color: var(--text); }
.list-state {
  font-size: .72rem; font-weight: 700;
  width: 1rem; text-align: center; flex-shrink: 0; line-height: 1;
  border-radius: 3px; padding: .05rem .15rem;
}
.state-neu { color: var(--text-muted); }
.state-inc { color: var(--color-primary); background: color-mix(in srgb, var(--color-primary) 15%, transparent); }
.state-exc { color: #dc2626; background: color-mix(in srgb, #dc2626 15%, transparent); }
.list-label { flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.mono-text  { font-family: ui-monospace,'SF Mono',Menlo,monospace; font-size: .74rem; }

/* ── Message search ──────────────────────────────── */
.search-wrap { position: relative; }
.search-input {
  width: 100%; padding: .35rem 1.75rem .35rem .55rem;
  border: 1px solid var(--border); border-radius: var(--radius);
  background: var(--bg); color: var(--text); font-size: .875rem;
}
.search-input:focus { outline: 2px solid var(--color-primary); outline-offset: -1px; }
.search-clear {
  position: absolute; right: .35rem; top: 50%; transform: translateY(-50%);
  background: none; border: none; cursor: pointer; color: var(--text-muted);
  padding: .2rem; border-radius: 3px; display: flex; align-items: center;
  transition: color .15s, background .15s;
}
.search-clear:hover { color: var(--text); background: var(--bg-hover); }

.empty-hint { font-size: .78rem; color: var(--text-muted); }
.empty-hint.pad { padding: .4rem .5rem; }
</style>
