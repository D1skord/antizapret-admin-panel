<template>
  <div>
    <h1 class="text-2xl font-bold mb-5">Главная</h1>

    <div v-if="error" class="mb-4 p-4 bg-red-100 dark:bg-red-900/20 text-red-700 dark:text-red-400 rounded-lg text-sm">
      Не удалось загрузить данные: {{ error }}
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
      <!-- CPU -->
      <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-5">
        <div class="flex items-center justify-between mb-3">
          <span class="text-sm font-medium text-gray-500 dark:text-gray-400">Процессор (CPU)</span>
          <span class="text-2xl font-bold text-gray-900 dark:text-white">
            {{ stats ? stats.cpu_percent.toFixed(1) : '—' }}%
          </span>
        </div>
        <div class="w-full bg-gray-100 dark:bg-gray-800 rounded-full h-2">
          <div
            class="h-2 rounded-full transition-all duration-500"
            :class="cpuBarColor"
            :style="{ width: stats ? stats.cpu_percent + '%' : '0%' }"
          />
        </div>
        <div class="mt-2 text-xs text-gray-400">
          Load avg: {{ stats ? `${stats.load_avg_1} / ${stats.load_avg_5} / ${stats.load_avg_15}` : '—' }}
        </div>
      </div>

      <!-- RAM -->
      <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-5">
        <div class="flex items-center justify-between mb-3">
          <span class="text-sm font-medium text-gray-500 dark:text-gray-400">Память (RAM)</span>
          <span class="text-2xl font-bold text-gray-900 dark:text-white">
            {{ stats ? stats.ram_percent.toFixed(1) : '—' }}%
          </span>
        </div>
        <div class="w-full bg-gray-100 dark:bg-gray-800 rounded-full h-2">
          <div
            class="h-2 rounded-full transition-all duration-500"
            :class="ramBarColor"
            :style="{ width: stats ? stats.ram_percent + '%' : '0%' }"
          />
        </div>
        <div class="mt-2 text-xs text-gray-400">
          {{ stats ? `${formatMB(stats.ram_used_mb)} / ${formatMB(stats.ram_total_mb)}` : '—' }}
        </div>
      </div>

      <!-- Network -->
      <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-5">
        <div class="flex items-center justify-between mb-3">
          <span class="text-sm font-medium text-gray-500 dark:text-gray-400">Сеть</span>
          <span class="text-xs text-gray-400">за 5 сек</span>
        </div>
        <div class="flex gap-6">
          <div>
            <div class="text-xs text-gray-400 mb-1">↓ Входящий</div>
            <div class="text-lg font-semibold text-gray-900 dark:text-white">{{ formatBytes(netRxSpeed) }}/с</div>
            <div class="text-xs text-gray-400 mt-1">всего: {{ formatBytes(stats?.net_rx_bytes ?? 0) }}</div>
          </div>
          <div>
            <div class="text-xs text-gray-400 mb-1">↑ Исходящий</div>
            <div class="text-lg font-semibold text-gray-900 dark:text-white">{{ formatBytes(netTxSpeed) }}/с</div>
            <div class="text-xs text-gray-400 mt-1">всего: {{ formatBytes(stats?.net_tx_bytes ?? 0) }}</div>
          </div>
        </div>
      </div>

      <!-- Uptime -->
      <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-5">
        <div class="flex items-center justify-between mb-3">
          <span class="text-sm font-medium text-gray-500 dark:text-gray-400">Uptime</span>
        </div>
        <div class="text-2xl font-bold text-gray-900 dark:text-white">
          {{ stats ? formatUptime(stats.uptime_seconds) : '—' }}
        </div>
        <div class="mt-2 text-xs text-gray-400">
          Сервер работает без перезагрузки
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import api from '@/api/index.js';

const stats = ref(null);
const error = ref(null);
const prevStats = ref(null);
const prevTimestamp = ref(null);
const netRxSpeed = ref(0);
const netTxSpeed = ref(0);

let timer = null;

async function fetchStats() {
  try {
    const now = Date.now();
    const res = await api.get('/api/system/stats');
    const data = res.data;

    if (prevStats.value && prevTimestamp.value) {
      const elapsed = (now - prevTimestamp.value) / 1000;
      if (elapsed > 0) {
        netRxSpeed.value = Math.max(0, (data.net_rx_bytes - prevStats.value.net_rx_bytes) / elapsed);
        netTxSpeed.value = Math.max(0, (data.net_tx_bytes - prevStats.value.net_tx_bytes) / elapsed);
      }
    }

    prevStats.value = data;
    prevTimestamp.value = now;
    stats.value = data;
    error.value = null;
  } catch (err) {
    error.value = err.message;
  }
}

onMounted(() => {
  fetchStats();
  timer = setInterval(fetchStats, 5000);
});

onUnmounted(() => {
  clearInterval(timer);
});

const cpuBarColor = computed(() => {
  if (!stats.value) return 'bg-indigo-500';
  if (stats.value.cpu_percent >= 85) return 'bg-red-500';
  if (stats.value.cpu_percent >= 60) return 'bg-yellow-500';
  return 'bg-indigo-500';
});

const ramBarColor = computed(() => {
  if (!stats.value) return 'bg-indigo-500';
  if (stats.value.ram_percent >= 90) return 'bg-red-500';
  if (stats.value.ram_percent >= 70) return 'bg-yellow-500';
  return 'bg-indigo-500';
});

function formatMB(mb) {
  if (mb >= 1024) return (mb / 1024).toFixed(1) + ' GB';
  return mb + ' MB';
}

function formatBytes(bytes) {
  if (bytes >= 1073741824) return parseFloat((bytes / 1073741824).toFixed(2)) + ' GB';
  if (bytes >= 1048576) return parseFloat((bytes / 1048576).toFixed(1)) + ' MB';
  if (bytes >= 1024) return parseFloat((bytes / 1024).toFixed(1)) + ' KB';
  return bytes + ' B';
}

function formatUptime(seconds) {
  const d = Math.floor(seconds / 86400);
  const h = Math.floor((seconds % 86400) / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  if (d > 0) return `${d}д ${h}ч ${m}м`;
  if (h > 0) return `${h}ч ${m}м`;
  return `${m}м`;
}
</script>
