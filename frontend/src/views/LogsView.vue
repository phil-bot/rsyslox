<template>
  <div class="logs-layout">
    <AppHeader @toggle-sidebar="sidebarOpen = !sidebarOpen" />

    <div class="logs-body">
      <FilterPanel
        v-model="sidebarOpen"
        :time-mode="timeMode"
        :relative-dur="relativeDur"
        :start-date="startDate"
        :end-date="endDate"
        :severities="severities"
        :facilities="facilities"
        :hosts="hosts"
        :tags="tags"
        :message-search="messageSearch"
        :available-hosts="availableHosts"
        :available-tags="availableTags"
        :available-severities="availableSeverities"
        :available-facilities="availableFacilities"
        @update:timeMode="timeMode = $event"
        @update:relativeDur="relativeDur = $event"
        @update:startDate="startDate = $event"
        @update:endDate="endDate = $event"
        @update:severities="severities = $event"
        @update:facilities="facilities = $event"
        @update:hosts="hosts = $event"
        @update:tags="tags = $event"
        @update:messageSearch="messageSearch = $event"
        @shift="shiftTimeWindow"
        @reset="resetFilters"
      />

      <div v-if="sidebarOpen && isMobile" class="sidebar-backdrop" @click="sidebarOpen = false" />

      <div class="logs-main" :class="{ 'is-paginated': !showAll }" ref="logsMainRef">
        <div v-if="error" class="error-banner">
          ⚠ {{ error }}
          <button @click="fetchLogs">Retry</button>
        </div>

        <LogTable
          :logs="logs"
          :total="total"
          :loading="loading"
          :page="page"
          :page-size="pageSize"
          :total-pages="totalPages"
          :selected-ids="selectedIds"
          :selected-count="selectedCount"
          :detail-id="detailEntry ? detailEntry.ID : null"
          :auto-refresh="autoRefresh"
          :auto-refresh-interval="autoRefreshInterval"
          :countdown="countdown"
          :new-ids="newIds"
          :first-load="firstLoad"
          :show-all="showAll"
          @open-detail="openDetail"
          @close-detail="closeDetail"
          @toggle-selection="toggleSelection"
          @toggle-select-all="toggleSelectAll"
          @clear-selection="clearSelection"
          @export-csv="exportCSV(selectedLogs.length ? selectedLogs : logs)"
          @export-json="exportJSON(selectedLogs.length ? selectedLogs : logs)"
          @set-page="setPage"
          @toggle-refresh="toggleAutoRefresh"
          @toggle-show-all="showAll = !showAll; fetchLogs()"
        />
      </div>

      <LogDetail :entry="detailEntry" @close="closeDetail" />
    </div>
  </div>
</template>

<script setup>
import { ref, watch, nextTick, onMounted, onBeforeUnmount } from 'vue'
import AppHeader   from '@/components/AppHeader.vue'
import FilterPanel from '@/components/FilterPanel.vue'
import LogTable    from '@/components/LogTable.vue'
import LogDetail   from '@/components/LogDetail.vue'
import { useLogsStore } from '@/stores/logs'
import { autoRefreshInterval as prefAutoRefresh, fontSize as prefFontSize } from '@/stores/preferences'

const {
  logs, total, loading, error,
  page, pageSize, totalPages, showAll,
  timeMode, relativeDur, startDate, endDate,
  severities, facilities, hosts, tags, messageSearch,
  selectedIds, selectedCount, selectedLogs,
  detailEntry,
  availableHosts, availableTags, availableSeverities, availableFacilities,
  autoRefresh, autoRefreshInterval, newIds, countdown, firstLoad,
  fetchLogs, fetchFilterOptions,
  setPage, resetFilters,
  toggleSelection, toggleSelectAll, clearSelection,
  openDetail, closeDetail,
  toggleAutoRefresh, startAutoRefresh, stopAutoRefresh,
  exportCSV, exportJSON,
  setPageSize,
} = useLogsStore()

const sidebarOpen = ref(true)
const isMobile    = ref(window.innerWidth < 768)
const logsMainRef = ref(null)

// ── Dynamic page size ────────────────────────────────────────────────────────
// Computed in LogsView (owns the container), not in LogTable (avoids circular
// dependency: watch(logs)→emit pageSize→setPageSize→fetchLogs→watch(logs)).
//
let ro = null

function computePageSize() {
  const wrap = logsMainRef.value
  if (!wrap) return

  // Measure chrome elements from the real DOM — accurate regardless of theme or zoom.
  const toolbar    = wrap.querySelector('.toolbar')
  const thead      = wrap.querySelector('thead')
  const pagination = wrap.querySelector('.pagination')
  const chromeH = (toolbar    ? toolbar.offsetHeight    : 40)
               + (thead       ? thead.offsetHeight      : 28)
               + (pagination  ? pagination.offsetHeight : 36)

  const available = wrap.clientHeight - chromeH
  if (available <= 0) return

  // Estimate natural row height from the current base font size.
  // td has .38rem top + .38rem bottom padding; text is .8rem at 1.5 line-height.
  const basePx = parseFloat(getComputedStyle(document.documentElement).fontSize) || 14

  //  const naturalRowH = Math.ceil(basePx * (0.76 + 0.8 * 1.5))
  // set naturalRowH to fixed value
  const naturalRowH = 31

  // Number of rows that fit — subtract 1 so the last row is never clipped by pagination.
  const n = Math.max(5, Math.floor(available / naturalRowH))

  // Exact row height so n rows fill the container with zero leftover space.
  //const exactRowH = Math.floor((available / n)).toFixed(2)
  const exactRowH = ((available / n)).toFixed(2)

  // Push --row-h into the table-scroll element; LogTable uses it for td min-height.
  const tableScroll = wrap.querySelector('.table-scroll')
  if (tableScroll) tableScroll.style.setProperty('--row-h', `${exactRowH}px`)

  if (n !== pageSize.value) {
    setPageSize(n)
    fetchLogs()
  }
}

function onResize() {
  isMobile.value = window.innerWidth < 768
  if (!isMobile.value) sidebarOpen.value = true
}

// Sync auto-refresh interval from preferences into the logs store
watch(prefAutoRefresh, (val) => {
  startAutoRefresh(val)
}, { immediate: true })

// Recompute row count when font size changes (naturalRowH depends on base font size)
watch(prefFontSize, () => { nextTick(computePageSize) })

onMounted(() => {
  window.addEventListener('resize', onResize)
  if (isMobile.value) sidebarOpen.value = false

  // Observe .logs-main height changes (sidebar open/close, window resize)
  nextTick(() => {
    if (logsMainRef.value) {
      ro = new ResizeObserver(computePageSize)
      ro.observe(logsMainRef.value)
      computePageSize()
    }
  })

  fetchLogs()
  fetchFilterOptions()
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', onResize)
  ro?.disconnect()
  stopAutoRefresh()
})

function shiftTimeWindow(direction) {
  const durations = {
    '15m': 15*60*1000, '1h': 60*60*1000, '6h': 6*60*60*1000,
    '24h': 24*60*60*1000, '7d': 7*24*60*60*1000, '30d': 30*24*60*60*1000,
  }
  const durMs = durations[relativeDur.value] ?? durations['1h']
  let end, start
  if (timeMode.value === 'absolute' && endDate.value && startDate.value) {
    end   = new Date(endDate.value).getTime()
    start = new Date(startDate.value).getTime()
    const winMs = end - start
    start += direction * winMs
    end   += direction * winMs
  } else {
    end   = Date.now()
    start = end - durMs
    start += direction * durMs
    end   += direction * durMs
  }
  const fmt = ts => new Date(ts).toISOString().slice(0, 16)
  timeMode.value  = 'absolute'
  startDate.value = fmt(start)
  endDate.value   = fmt(end)
}
</script>

<style scoped>
.logs-layout {
  display: flex; flex-direction: column;
  height: 100%; overflow: hidden;
}
.logs-body {
  display: flex; flex: 1; overflow: hidden; position: relative;
}
.logs-main {
  flex: 1; display: flex; flex-direction: column; overflow: hidden; min-width: 0;
}
.error-banner {
  background: #fef2f2; color: #dc2626;
  border-bottom: 1px solid #fca5a5;
  padding: .5rem .75rem; font-size: .875rem;
  display: flex; align-items: center; gap: .75rem;
}
[data-theme="dark"] .error-banner { background: #2d1212; border-color: #7f1d1d; }
.error-banner button {
  background: none; border: 1px solid currentColor;
  border-radius: var(--radius); cursor: pointer;
  padding: .2rem .5rem; font-size: .8rem; color: inherit;
}
.sidebar-backdrop {
  position: fixed; inset: 0;
  background: rgba(0,0,0,.3); z-index: 40;
}
.is-paginated :deep(.table-scroll) {
  overflow: hidden !important;
}
</style>
