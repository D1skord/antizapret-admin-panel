<template>
  <div class="flex items-center justify-between p-4">
    <div class="text-sm text-gray-700 dark:text-gray-400">
      Показано с <span class="font-medium">{{ startItem }}</span> по <span class="font-medium">{{ endItem }}</span> из <span class="font-medium">{{ totalItems }}</span>
    </div>
    <div class="flex items-center">
      <button
        @click="changePage(currentPage - 1)"
        :disabled="currentPage === 1"
        class="px-3 py-1 mx-1 rounded-md bg-gray-200 dark:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        &laquo;
      </button>
      <button
        v-for="page in pages"
        :key="page"
        @click="changePage(page)"
        :class="[
          'px-3 py-1 mx-1 rounded-md',
          currentPage === page ? 'bg-blue-500 text-white' : 'bg-gray-200 dark:bg-gray-700'
        ]"
      >
        {{ page }}
      </button>
      <button
        @click="changePage(currentPage + 1)"
        :disabled="currentPage === totalPages"
        class="px-3 py-1 mx-1 rounded-md bg-gray-200 dark:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        &raquo;
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  currentPage: {
    type: Number,
    required: true,
  },
  totalPages: {
    type: Number,
    required: true,
  },
  totalItems: {
    type: Number,
    required: true,
  },
  perPage: {
    type: Number,
    required: true,
  }
});

const emit = defineEmits(['page-changed']);

const startItem = computed(() => (props.currentPage - 1) * props.perPage + 1);
const endItem = computed(() => {
  const end = props.currentPage * props.perPage;
  return end > props.totalItems ? props.totalItems : end;
});

const pages = computed(() => {
    const range = [];
    for (let i = 1; i <= props.totalPages; i++) {
        range.push(i);
    }
    return range;
});


function changePage(page) {
  if (page > 0 && page <= props.totalPages) {
    emit('page-changed', page);
  }
}
</script>
