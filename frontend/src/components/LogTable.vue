<template>
  <div class="log-table-wrap">
    <!-- Toolbar -->
    <div class="toolbar">
      <div class="toolbar-left">

        <!-- Filter open button — visible only when sidebar is collapsed -->
        <button
          v-if="sidebarCollapsed"
          class="filter-open-btn"
          @click="$emit('toggle-sidebar')"
          :title="t('filter.open')"
        >
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3"/>
          </svg>
          {{ t('filter.open') }}
        </button>

        <span v-if="selectedCount > 0" class="selection-info">
          {{ selectedCount }} {{ t('table.selected') }}
          <button class="toolbar-btn danger" @click="$emit('clear-selection')">{{ t('table.clear') }}</button>
          <button class="toolbar-btn" @click="$emit('export-csv')">{{ t('table.export_csv') }}</button>
          <button class="toolbar-btn" @click="$emit('export-json')">{{ t('table.export_json') }}</button>
        </span>
        <span v-else class="total-info">
          <template v-if="!loading">
            <span>{{ fmtNumber(total) }}</span>
            <span class="of-label">{{ t('table.entries') }}</span>
            <template v-if="dbTotal > 0 && dbTotal !== total">
              <span class="of-label">·</span>
              <span>{{ fmtNumber(dbTotal) }}</span>
              <span class="of-label">{{ t('table.db_total') }}</span>
            </template>
          </template>
          <span v-else class="loading-pulse">{{ t('table.loading') }}</span>
        </span>
      </div>

      <div class="toolbar-right">
        <!-- Live indicator — only shown when filter sidebar is collapsed -->
        <span v-if="autoRefresh && sidebarCollapsed" class="live-badge">
          <span class="live-badge-dot"></span>
          {{ t('table.live') }}
        </span>

        <!-- View mode: paginated vs. all -->
        <button
          class="ctrl-btn"
          :class="{ active: showAll }"
          @click="$emit('toggle-show-all')"
          :title="showAll ? t('table.view_multipage') : t('table.view_singlepage')"
        >
          <svg v-if="!showAll" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/></svg>
          <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="3" y1="6" x2="21" y2="6"/><line x1="3" y1="12" x2="21" y2="12"/><line x1="3" y1="18" x2="21" y2="18"/></svg>
          {{ showAll ? t('table.view_multipage') : t('table.view_singlepage') }}
        </button>
      </div>
    </div>

    <!-- Table -->
    <div class="table-scroll" @click.self="$emit('close-detail')">
      <table class="log-table">
        <thead>
          <tr>
            <th class="col-check">
              <input type="checkbox"
                :checked="allSelected"
                :indeterminate="someSelected"
                @change="$emit('toggle-select-all')" />
            </th>
            <th class="col-time">{{ t('table.col_time') }}</th>
            <th class="col-sev">{{ t('table.col_severity') }}</th>
            <th class="col-fac">{{ t('table.col_facility') }}</th>
            <th class="col-host">{{ t('table.col_host') }}</th>
            <th class="col-tag">{{ t('table.col_tag') }}</th>
            <th class="col-msg">{{ t('table.col_message') }}</th>
          </tr>
        </thead>
        <tbody>
          <template v-if="loading && firstLoad">
            <tr v-for="i in pageSize" :key="'sk-' + i" class="skeleton-row">
              <td></td>
              <td><span class="skel skel-time"></span></td>
              <td><span class="skel skel-sev"></span></td>
              <td><span class="skel skel-fac"></span></td>
              <td><span class="skel skel-host"></span></td>
              <td><span class="skel skel-tag"></span></td>
              <td><span class="skel skel-msg"></span></td>
            </tr>
          </template>

          <template v-else-if="!logs.length && !loading">
            <tr>
              <td colspan="7" class="empty-state">{{ t('table.no_entries') }}</td>
            </tr>
          </template>

          <template v-else>
            <tr
              v-for="entry in logs"
              :key="entry.ID"
              class="log-row"
              :class="{
                selected:   selectedIds.has(entry.ID),
                detail:     detailId === entry.ID,
                'row-new':  newIds.has(entry.ID),
                ['sev-row-' + severityVal(entry)]: true
              }"
              @click="$emit('open-detail', entry)"
            >
              <td class="col-check" @click.stop>
                <input type="checkbox"
                  :checked="selectedIds.has(entry.ID)"
                  @change="$emit('toggle-selection', entry.ID)" />
              </td>
              <td class="col-time mono">{{ formatTime(entry.ReceivedAt) }}</td>
              <td class="col-sev">
                <span :class="'sev-badge sev-' + severityVal(entry)">
                  {{ severityLabel(entry) }}
                </span>
              </td>
              <td class="col-fac mono">{{ entry.Facility_Label || entry.Facility }}</td>
              <td class="col-host mono">{{ entry.FromHost }}</td>
              <td class="col-tag mono">{{ trimTag(entry.SysLogTag) }}</td>
              <td class="col-msg mono" v-html="highlightMessage(entry.Message)"></td>
            </tr>
          </template>
        </tbody>
      </table>
    </div>

    <!-- Pagination -->
    <div v-if="!showAll" class="pagination">
      <button class="pag-btn" :disabled="page <= 1" @click="$emit('set-page', 1)">«</button>
      <button class="pag-btn" :disabled="page <= 1" @click="$emit('set-page', page - 1)">‹</button>
      <span class="pag-info">{{ t('table.page_of', { page, total: totalPages }) }}</span>
      <button class="pag-btn" :disabled="page >= totalPages" @click="$emit('set-page', page + 1)">›</button>
      <button class="pag-btn" :disabled="page >= totalPages" @click="$emit('set-page', totalPages)">»</button>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useLocale } from '@/composables/useLocale'
import { timeFormat } from '@/stores/preferences'

const props = defineProps({
  logs:              { type: Array,          default: () => [] },
  total:             { type: Number,         default: 0 },
  dbTotal:           { type: Number,         default: 0 },
  loading:           { type: Boolean,        default: false },
  page:              { type: Number,         default: 1 },
  pageSize:          { type: Number,         default: 20 },
  totalPages:        { type: Number,         default: 1 },
  selectedIds:       { type: Object,         default: () => new Set() },
  selectedCount:     { type: Number,         default: 0 },
  detailId:          { type: [Number, null], default: null },
  autoRefresh:       { type: Boolean,        default: false },
  sidebarCollapsed:  { type: Boolean,        default: false },
  newIds:            { type: Object,         default: () => new Set() },
  showAll:           { type: Boolean,        default: false },
  firstLoad:         { type: Boolean,        default: true },
  messageSearch:     { type: String,         default: '' },
})

defineEmits([
  'open-detail', 'close-detail',
  'toggle-selection', 'toggle-select-all', 'clear-selection',
  'export-csv', 'export-json',
  'set-page',
  'toggle-show-all',
  'toggle-sidebar',
])

const { t, fmtNumber } = useLocale()
const allSelected  = computed(() => props.logs.length > 0 && props.selectedCount === props.logs.length)
const someSelected = computed(() => props.selectedCount > 0 && props.selectedCount < props.logs.length)

function severityVal(entry) {
  if (typeof entry.Severity === 'number') return entry.Severity
  if (typeof entry.Priority === 'number') return entry.Priority % 8
  return 6
}
function severityLabel(entry) {
  const short = { 0:'EMERG',1:'ALERT',2:'CRIT',3:'ERR',4:'WARN',5:'NOTICE',6:'INFO',7:'DEBUG' }
  return short[severityVal(entry)] ?? String(severityVal(entry))
}
function trimTag(tag) {
  if (!tag) return ''
  return tag.replace(/[:\[]+\d*\]?$/, '').trim()
}
function highlightMessage(text) {
  if (!props.messageSearch || !text) return escapeHtml(text ?? '')
  const escaped = escapeHtml(text)
  const term    = escapeHtml(props.messageSearch)
  const re      = new RegExp(term.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), 'gi')
  return escaped.replace(re, m => `<mark class="msg-hl">${m}</mark>`)
}
function escapeHtml(s) {
  return s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;')
}
function formatTime(ts) {
  const d = new Date(ts)
  if (isNaN(d)) return ts
  const pad = n => String(n).padStart(2, '0')
  const date = `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())}`
  if (timeFormat.value === '12h') {
    const h = d.getHours(); const ampm = h >= 12 ? 'PM' : 'AM'; const h12 = h % 12 || 12
    return `${date} ${pad(h12)}:${pad(d.getMinutes())}:${pad(d.getSeconds())} ${ampm}`
  }
  return `${date} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}
</script>

<style scoped>
.log-table-wrap {
  display: flex; flex-direction: column;
  flex: 1; overflow: hidden; min-height: 0;
}

/* ── Toolbar ──────────────────────────────────────── */
.toolbar {
  display: flex; align-items: center; justify-content: space-between;
  padding: .4rem .75rem;
  border-bottom: 1px solid var(--border);
  background: var(--bg-surface);
  flex-shrink: 0; gap: .5rem; min-height: 40px;
}
.toolbar-left, .toolbar-right { display: flex; align-items: center; gap: .5rem; }

/* Filter open button */
.filter-open-btn {
  display: inline-flex; align-items: center; gap: .35rem;
  padding: .28rem .6rem; font-size: .8rem;
  background: var(--bg); border: 1px solid var(--border);
  border-radius: var(--radius); cursor: pointer; color: var(--text-muted);
  transition: color .15s, background .15s, border-color .15s;
  white-space: nowrap; flex-shrink: 0;
}
.filter-open-btn:hover {
  background: var(--bg-hover); color: var(--text);
  border-color: var(--color-primary);
}

.total-info     { font-size: .8rem; color: var(--text-muted); display: flex; align-items: center; gap: .35rem; }
.of-label       { color: var(--text-muted); }
.selection-info { display: flex; align-items: center; gap: .375rem; font-size: .8rem; font-weight: 500; }

.toolbar-btn {
  background: var(--bg); border: 1px solid var(--border);
  border-radius: var(--radius); cursor: pointer;
  padding: .2rem .5rem; font-size: .78rem; color: var(--text); transition: all .15s;
}
.toolbar-btn:hover { background: var(--bg-hover); }
.toolbar-btn.danger { color: #dc2626; border-color: #fca5a5; }
.toolbar-btn.danger:hover { background: #fef2f2; }
[data-theme="dark"] .toolbar-btn.danger:hover { background: #2d1212; }

.live-badge {
  display: inline-flex; align-items: center; gap: .3rem;
  font-size: .75rem; font-weight: 600; color: var(--color-primary);
  padding: .15rem .45rem; border: 1px solid var(--color-primary);
  border-radius: var(--radius); background: var(--bg-selected); user-select: none;
}
.live-badge-dot {
  width: 6px; height: 6px; border-radius: 50%; background: var(--color-primary);
  animation: live-pulse 1.5s ease-in-out infinite;
}
@keyframes live-pulse { 0%, 100% { opacity: 1; transform: scale(1); } 50% { opacity: .45; transform: scale(.7); } }

.ctrl-btn {
  display: inline-flex; align-items: center; gap: .35rem;
  padding: .28rem .6rem; font-size: .8rem;
  background: var(--bg); border: 1px solid var(--border);
  border-radius: var(--radius); cursor: pointer; color: var(--text-muted);
  transition: color .15s, background .15s, border-color .15s;
  font-variant-numeric: tabular-nums; white-space: nowrap;
}
.ctrl-btn:hover { background: var(--bg-hover); color: var(--text); }
.ctrl-btn.active { color: var(--color-primary); border-color: var(--color-primary); background: var(--bg-selected); }

@keyframes pulse { 0%,100%{opacity:1} 50%{opacity:.4} }
.loading-pulse { animation: pulse 1.5s infinite; font-size: .8rem; }

/* ── Table ────────────────────────────────────────── */
.table-scroll { flex: 1; overflow-y: auto; overflow-x: auto; min-height: 0; background: var(--bg); }
.log-table {
  width: 100%; border-collapse: collapse;
  font-size: .8rem; table-layout: fixed; background: var(--bg);
}
.log-table thead tr { position: sticky; top: 0; z-index: 2; }
.log-table th {
  padding: .45rem 0; text-align: left;
  font-size: .68rem; font-weight: 700;
  text-transform: uppercase; letter-spacing: .05em;
  color: var(--text-muted); background: var(--bg-surface);
  border-bottom: 2px solid var(--border); white-space: nowrap;
}

.col-check { width: 36px; text-align: center; }
.col-time  { width: 152px; }
.col-sev   { width: 72px; }
.col-fac   { width: 88px; }
.col-host  { width: 140px; }
.col-tag   { width: 120px; }
.col-msg   { /* fill remaining */ }

.log-row {
  border-bottom: 1px solid var(--border); background: var(--bg);
  cursor: pointer; transition: background .1s;
  height: var(--row-h, 31px);
}
@keyframes row-flash {
  0%   { background: color-mix(in srgb, var(--color-primary) 22%, transparent); }
  100% { background: transparent; }
}
.row-new { animation: row-flash 1.4s ease-out forwards; }
.log-row:hover    { background: var(--bg-hover); }
.log-row.selected { background: var(--bg-selected); }
.log-row.detail   { outline: 2px solid var(--color-primary); outline-offset: -1px; }

.log-row.sev-row-0 { border-left: 3px solid var(--sev-0); }
.log-row.sev-row-1 { border-left: 3px solid var(--sev-1); }
.log-row.sev-row-2 { border-left: 3px solid var(--sev-2); }
.log-row.sev-row-3 { border-left: 3px solid var(--sev-3); }
.log-row.sev-row-4 { border-left: 3px solid var(--sev-4); }
.log-row.sev-row-5 { border-left: 3px solid var(--sev-5); }
.log-row.sev-row-6 { border-left: 3px solid transparent; }
.log-row.sev-row-7 { border-left: 3px solid transparent; }

.log-table td {
  padding: 0; overflow: hidden;
  text-overflow: ellipsis; white-space: nowrap;
  vertical-align: middle; color: var(--text);
}
.mono { font-family: ui-monospace,'Cascadia Code','SF Mono',Menlo,monospace; font-size: .77rem; }

.sev-badge {
  display: inline-block; padding: .1rem .3rem; border-radius: 3px;
  font-size: .68rem; font-weight: 700; letter-spacing: .02em; color: #fff;
}
.sev-0 { background: var(--sev-0); } .sev-1 { background: var(--sev-1); }
.sev-2 { background: var(--sev-2); } .sev-3 { background: var(--sev-3); }
.sev-4 { background: var(--sev-4); } .sev-5 { background: var(--sev-5); }
.sev-6 { background: var(--sev-6); } .sev-7 { background: var(--sev-7); color: var(--text); }

.empty-state { text-align: center; padding: 3rem; color: var(--text-muted); font-size: .875rem; }

.skeleton-row td { padding: .5rem .625rem; }
.skel { display: inline-block; height: .72rem; border-radius: 3px; background: var(--border); animation: pulse 1.5s ease-in-out infinite; }
.skel-time { width: 128px; } .skel-sev { width: 44px; } .skel-fac { width: 60px; }
.skel-host { width: 80px; }  .skel-tag { width: 64px; } .skel-msg { width: 60%; }

/* ── Pagination ───────────────────────────────────── */
.pagination {
  display: flex; align-items: center; justify-content: center;
  gap: .375rem; padding: .5rem;
  border-top: 1px solid var(--border); background: var(--bg-surface); flex-shrink: 0;
}
.pag-btn {
  background: var(--bg); border: 1px solid var(--border);
  border-radius: var(--radius); cursor: pointer;
  padding: .28rem .6rem; font-size: .875rem; color: var(--text);
  transition: all .15s; min-width: 34px;
}
.pag-btn:hover:not(:disabled) { background: var(--bg-hover); border-color: var(--color-primary); }
.pag-btn:disabled { opacity: .35; cursor: default; }
.pag-info { font-size: .8rem; color: var(--text-muted); padding: 0 .375rem; }
</style>

<style>
.msg-hl { background: #fef08a; color: #713f12; border-radius: 2px; padding: 0 .1rem; }
[data-theme="dark"] .msg-hl { background: #854d0e; color: #fef9c3; }
</style>
