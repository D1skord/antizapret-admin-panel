<template>
  <div
    v-if="isMobileOpen"
    class="fixed inset-0 bg-gray-900/50 z-40 lg:hidden"
    @click="toggleMobileSidebar"
  ></div>
</template>

<script setup>
import { useSidebar } from '@/composables/useSidebar'
const { toggleMobileSidebar, isMobileOpen } = useSidebar()
</script>