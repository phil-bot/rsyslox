<template>
  <div class="logs-layout">
    <AppHeader />

    <div class="logs-body">
      <FilterPanel
        v-model="sidebarOpen"
        :time-mode="timeMode"
        :relative-dur="relativeDur"
        :start-date="startDate"
        :end-date="endDate"
        :severities="severities"
        :exclude-severities="excludeSeverities"
        :facilities="facilities"
        :exclude-facilities="excludeFacilities"
        :hosts="hosts"
        :exclude-hosts="excludeHosts"
        :tags="tags"
        :exclude-tags="excludeTags"
        :message-search="messageSearch"
        :available-hosts="availableHosts"
        :available-tags="availableTags"
        :available-severities="availableSeverities"
        :available-facilities="availableFacilities"
        :auto-refresh="autoRefresh"
        :countdown="countdown"
        @update:timeMode="timeMode = $event"
        @update:relativeDur="relativeDur = $event"
        @update:startDate="startDate = $event"
        @update:endDate="endDate = $event"
        @update:severities="severities = $event"
        @update:excludeSeverities="excludeSeverities = $event"
        @update:facilities="facilities = $event"
        @update:excludeFacilities="excludeFacilities = $event"
        @update:hosts="hosts = $event"
        @update:excludeHosts="excludeHosts = $event"
        @update:tags="tags = $event"
        @update:excludeTags="excludeTags = $event"
        @update:messageSearch="messageSearch = $event"
        @shift="shiftTimeWindow"
        @toggle-live="toggleAutoRefresh"
        @exit-live="stopAutoRefresh"
        @reset="resetFilters"
        @close="sidebarOpen = false"
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
          :db-total="dbTotal"
          :loading="loading"
          :page="page"
          :page-size="pageSize"
          :total-pages="totalPages"
          :selected-ids="selectedIds"
          :selected-count="selectedCount"
          :detail-id="detailEntry ? detailEntry.ID : null"
          :auto-refresh="autoRefresh"
          :sidebar-collapsed="!sidebarOpen"
          :message-search="messageSearch"
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
          @toggle-show-all="showAll = !showAll; fetchLogs()"
          @toggle-sidebar="sidebarOpen = !sidebarOpen"
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
  logs, total, dbTotal, loading, error,
  page, pageSize, totalPages, showAll,
  timeMode, relativeDur, startDate, endDate,
  severities, excludeSeverities, facilities, excludeFacilities,
  hosts, excludeHosts, tags, excludeTags, messageSearch,
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

let ro = null

function computePageSize() {
  const wrap = logsMainRef.value
  if (!wrap) return
  const toolbar    = wrap.querySelector('.toolbar')
  const thead      = wrap.querySelector('thead')
  const pagination = wrap.querySelector('.pagination')
  const chromeH = (toolbar    ? toolbar.offsetHeight    : 40)
               + (thead       ? thead.offsetHeight      : 28)
               + (pagination  ? pagination.offsetHeight : 36)
  const available = wrap.clientHeight - chromeH
  if (available <= 0) return
  const naturalRowH = 31
  const n = Math.max(5, Math.floor(available / naturalRowH))
  const exactRowH = ((available / n)).toFixed(2)
  const tableScroll = wrap.querySelector('.table-scroll')
  if (tableScroll) tableScroll.style.setProperty('--row-h', `${exactRowH}px`)
  if (n !== pageSize.value) { setPageSize(n); fetchLogs() }
}

function onResize() {
  isMobile.value = window.innerWidth < 768
  if (!isMobile.value) sidebarOpen.value = true
}

watch(prefAutoRefresh, (val) => { startAutoRefresh(val) }, { immediate: true })
watch(prefFontSize, () => { nextTick(computePageSize) })

onMounted(() => {
  window.addEventListener('resize', onResize)
  if (isMobile.value) sidebarOpen.value = false
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
.logs-layout { display: flex; flex-direction: column; height: 100%; overflow: hidden; }
.logs-body   { display: flex; flex: 1; overflow: hidden; position: relative; }
.logs-main   { flex: 1; display: flex; flex-direction: column; overflow: hidden; min-width: 0; }
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
.sidebar-backdrop { position: fixed; inset: 0; background: rgba(0,0,0,.3); z-index: 40; }
.is-paginated :deep(.table-scroll) { overflow: hidden !important; }
</style>
