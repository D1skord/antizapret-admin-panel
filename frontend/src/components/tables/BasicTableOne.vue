<template>
  <div
    class="overflow-hidden rounded-xl border border-gray-200 bg-white dark:border-gray-800 dark:bg-white/[0.03]"
  >
    <div class="p-5 flex justify-end">
      <button @click="isModalOpen = true" class="px-4 py-2 text-white bg-blue-500 rounded-md hover:bg-blue-600">
        Добавить клиента
      </button>
    </div>
    <div v-if="loading" class="p-5 text-center text-gray-500 dark:text-gray-400">Загрузка данных...</div>
    <div v-else-if="error" class="p-5 text-center text-error-500">Ошибка при загрузке данных: {{ error }}</div>
    <div v-else-if="users.length === 0" class="p-5 text-center text-gray-500 dark:text-gray-400">
      Нет данных о клиентах.
    </div>
    <div v-else class="max-w-full overflow-x-auto custom-scrollbar">
      <table class="min-w-full">
        <thead>
          <tr class="border-b border-gray-200 dark:border-gray-700">
            <th class="px-5 py-3 text-left w-3/11 sm:px-6">
              <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">Пользователь</p>
            </th>
            <th class="px-5 py-3 text-left w-2/11 sm:px-6">
              <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">Тип</p>
            </th>
            <th class="px-5 py-3 text-left w-2/11 sm:px-6">
              <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">Проект</p>
            </th>
            <th class="px-5 py-3 text-left w-2/11 sm:px-6">
              <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">Команда</p>
            </th>
            <th class="px-5 py-3 text-left w-1/11 sm:px-6">
              <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">Статус</p>
            </th>
            <th class="px-5 py-3 text-left w-1/11 sm:px-6">
              <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">Бюджет</p>
            </th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
          <tr
            v-for="(user, index) in users"
            :key="index"
            class="border-t border-gray-100 dark:border-gray-800"
          >
            <td class="px-5 py-4 sm:px-6">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 overflow-hidden rounded-full">
                  <img :src="user.avatar" :alt="user.name" />
                </div>
                <div>
                  <span class="block font-medium text-gray-800 text-theme-sm dark:text-white/90">
                    {{ user.name }}
                  </span>
                  <span class="block text-gray-500 text-theme-xs dark:text-gray-400">
                    {{ user.role }}
                  </span>
                </div>
              </div>
            </td>
            <td class="px-5 py-4 sm:px-6">
              <p class="text-gray-500 text-theme-sm dark:text-gray-400">{{ user.type }}</p>
            </td>
            <td class="px-5 py-4 sm:px-6">
              <p class="text-gray-500 text-theme-sm dark:text-gray-400">{{ user.project }}</p>
            </td>
            <td class="px-5 py-4 sm:px-6">
              <div class="flex -space-x-2">
                <div
                  v-for="(member, memberIndex) in user.team"
                  :key="memberIndex"
                  class="w-6 h-6 overflow-hidden border-2 border-white rounded-full dark:border-gray-900"
                >
                  <img :src="member" alt="team member" />
                </div>
              </div>
            </td>
            <td class="px-5 py-4 sm:px-6">
              <span
                :class="[
                  'rounded-full px-2 py-0.5 text-theme-xs font-medium',
                  {
                    'bg-success-50 text-success-700 dark:bg-success-500/15 dark:text-success-500':
                      user.status === 'Active',
                    'bg-warning-50 text-warning-700 dark:bg-warning-500/15 dark:text-warning-400':
                      user.status === 'Pending',
                    'bg-error-50 text-error-700 dark:bg-error-500/15 dark:text-error-500':
                      user.status === 'Cancel',
                  },
                ]"
              >
                {{ user.status }}
              </span>
            </td>
            <td class="px-5 py-4 sm:px-6">
              <p class="text-gray-500 text-theme-sm dark:text-gray-400">{{ user.budget }}</p>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>

  <Modal v-if="isModalOpen" @close="isModalOpen = false" :fullScreenBackdrop="true">
    <template #body>
      <div class="relative w-full max-w-md p-6 bg-white rounded-lg shadow-lg dark:bg-gray-800">
        <h2 class="text-xl font-bold mb-4">Добавить нового клиента</h2>
        <form @submit.prevent="createClient">
          <div class="mb-4">
            <label for="clientName" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Имя клиента</label>
            <input type="text" id="clientName" v-model="newClientName" class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white" placeholder="Например, my-client" required>
          </div>

          <div class="mb-4">
            <label class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Тип клиента</label>
            <div class="flex items-center">
              <input id="openvpn" type="radio" value="openvpn" v-model="newClientType" class="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600">
              <label for="openvpn" class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300">OpenVPN</label>
            </div>
            <div class="flex items-center mt-2">
              <input id="wireguard" type="radio" value="wireguard" v-model="newClientType" class="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600">
              <label for="wireguard" class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300">WireGuard</label>
            </div>
          </div>
          
          <div v-if="createError" class="p-3 mb-4 text-sm text-red-700 bg-red-100 rounded-lg dark:bg-red-200 dark:text-red-800">
            {{ createError }}
          </div>

          <div class="flex justify-end gap-3">
            <button type="button" @click="isModalOpen = false" class="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-200 rounded-md hover:bg-gray-300 dark:bg-gray-600 dark:text-gray-200 dark:hover:bg-gray-500">
              Отмена
            </button>
            <button type="submit" class="px-4 py-2 text-sm font-medium text-white bg-blue-500 rounded-md hover:bg-blue-600">
              Создать
            </button>
          </div>
        </form>
      </div>
    </template>
  </Modal>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import api from '../../api';
import Modal from '../ui/Modal.vue';

const users = ref([]);
const loading = ref(true);
const error = ref(null);

const isModalOpen = ref(false);
const newClientName = ref('');
const newClientType = ref('openvpn');
const createError = ref(null);

const fetchClients = async () => {
  try {
    loading.value = true;
    const response = await api.get('/api/clients');
    users.value = response.data;
  } catch (err) {
    console.error('Failed to fetch clients:', err);
    error.value = err.message || 'Неизвестная ошибка';
  } finally {
    loading.value = false;
  }
};

const createClient = async () => {
  createError.value = null;
  try {
    const response = await api.post('/api/clients', {
      name: newClientName.value,
      type: newClientType.value,
    });
    users.value.push(response.data);
    isModalOpen.value = false;
    newClientName.value = '';
    newClientType.value = 'openvpn';
  } catch (err) {
    console.error('Failed to create client:', err);
    createError.value = (err.response && err.response.data && err.response.data.details) || err.message || 'Не удалось создать клиента.';
  }
};

onMounted(fetchClients);
</script>
