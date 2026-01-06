<template>
  <slot></slot>
</template>

<script setup>
import { useSidebarProvider } from '@/composables/useSidebar'

useSidebarProvider()
</script>