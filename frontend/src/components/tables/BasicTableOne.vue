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
            <th class="px-5 py-3 text-left w-3/12 sm:px-6">
              <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">Пользователь</p>
            </th>
            <th class="px-5 py-3 text-left w-2/12 sm:px-6">
              <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">Тип</p>
            </th>
            <th class="px-5 py-3 text-left w-1/12 sm:px-6">
              <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">Статус</p>
            </th>
            <th class="px-5 py-3 text-left w-2/12 sm:px-6">
              <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">Дата создания</p>
            </th>
            <th class="px-5 py-3 text-left w-2/12 sm:px-6">
              <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">Действия</p>
            </th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
          <tr
            v-for="user in users"
            :key="user.id"
            class="border-t border-gray-100 dark:border-gray-800"
          >

            <td class="px-5 py-4 sm:px-6">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 overflow-hidden rounded-full">
                  <div :class="getUserColor(user.id)" class="user-avatar" :alt="user.name"></div>
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
              <p class="text-gray-500 text-theme-sm dark:text-gray-400">{{ new Date(user.createdAt).toLocaleDateString() }}</p>
            </td>
            <td class="px-5 py-4 sm:px-6">
                            <div class="flex items-center gap-4">
                              <button @click="showQrCode(user)" class="text-blue-500 hover:text-blue-700">
                                QR
                              </button>
                              <button @click="downloadConfig(user)" class="text-green-500 hover:text-green-700">
                                Скачать
                              </button>
                              <button @click="deleteClient(user.id)" class="text-error-500 hover:text-error-700">
                                Удалить
                              </button>
                            </div>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                   <Pagination
                      v-if="!loading && totalClients > 0"
                      :current-page="currentPage"
                      :total-pages="totalPages"
                      :total-items="totalClients"
                      :per-page="perPage"
                      @page-changed="fetchClients"
                    />
                </div>

                <!-- Modal for Adding Client -->
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

                <!-- Modal for QR Code -->
                <Modal v-if="isQrModalOpen" @close="isQrModalOpen = false" :fullScreenBackdrop="true">
                  <template #body>
                      <div class="relative w-full max-w-sm p-6 bg-white rounded-lg shadow-lg dark:bg-gray-800 text-center">
                          <h2 class="text-xl font-bold mb-4">Конфигурация для {{ selectedClientName }}</h2>
                          <div v-if="qrError" class="p-3 mb-4 text-sm text-red-700 bg-red-100 rounded-lg dark:bg-red-200 dark:text-red-800">
                            Не удалось сгенерировать QR-код: {{ qrError }}
                          </div>
                          <div v-else-if="!qrCodeUrl" class="my-8">
                            <p>Генерация QR-кода...</p>
                          </div>
                          <div v-else class="flex justify-center p-4 bg-white rounded-lg">
                            <QrcodeVue :value="qrCodeUrl" :size="250" level="H" />
                          </div>
                          <p class="mt-4 text-sm text-gray-500 dark:text-gray-400">
                            Отсканируйте код камерой телефона. Ссылка действительна 5 минут и работает один раз.
                          </p>
                           <div class="mt-6 flex justify-center">
                              <button type="button" @click="isQrModalOpen = false" class="px-6 py-2 text-sm font-medium text-gray-700 bg-gray-200 rounded-md hover:bg-gray-300 dark:bg-gray-600 dark:text-gray-200 dark:hover:bg-gray-500">
                                Закрыть
                              </button>
                          </div>
                      </div>
                  </template>
                </Modal>

              </template>

              <script setup>
import { ref, onMounted } from 'vue';
import QrcodeVue from 'qrcode.vue';
import api from '../../api';
import Modal from '../ui/Modal.vue';
import Pagination from '../common/Pagination.vue';

const users = ref([]);
const loading = ref(true);
const error = ref(null);

// Pagination State
const currentPage = ref(1);
const totalPages = ref(1);
const totalClients = ref(0);
const perPage = 10; // or make it a ref if you want it to be dynamic

// State for Add Client Modal
const isModalOpen = ref(false);
const newClientName = ref('');
const newClientType = ref('openvpn');
const createError = ref(null);

// State for QR Code Modal
const isQrModalOpen = ref(false);
const qrCodeUrl = ref('');
const qrError = ref(null);
const selectedClientName = ref('');

const fetchClients = async (page = 1) => {
  try {
    loading.value = true;
    const response = await api.get('/api/clients', {
      params: { page: page, limit: perPage }
    });
    const data = response.data || { clients: [], total: 0 };
    users.value = data.clients || [];
    totalClients.value = data.total || 0;
    totalPages.value = Math.ceil(totalClients.value / perPage);
    currentPage.value = page;
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
    await api.post('/api/clients', {
      name: newClientName.value,
      type: newClientType.value,
    });
    isModalOpen.value = false;
    newClientName.value = '';
    newClientType.value = 'openvpn';
    await fetchClients(1); // Refresh and go to the first page to see the new client
  } catch (err) {
    console.error('Failed to create client:', err);
    createError.value = (err.response && err.response.data && err.response.data.details) || err.message || 'Не удалось создать клиента.';
  }
};

const deleteClient = async (id) => {
  if (confirm('Вы уверены, что хотите удалить этого клиента?')) {
    try {
      await api.delete(`/api/clients/${id}`);
      // If the user was on the last page and deleted the last item, go to the previous page.
      if (users.value.length === 1 && currentPage.value > 1) {
        await fetchClients(currentPage.value - 1);
      } else {
        await fetchClients(currentPage.value); // Refresh the current page
      }
    } catch (err) {
      console.error(`Failed to delete client ${id}:`, err);
      alert('Не удалось удалить клиента.');
    }
  }
};

const showQrCode = async (client) => {
  selectedClientName.value = client.name;
  isQrModalOpen.value = true;
  qrCodeUrl.value = '';
  qrError.value = null;

  try {
    const response = await api.get(`/api/clients/${client.id}/qr-token`);
    qrCodeUrl.value = window.location.origin + response.data.download_url;
  } catch (err) {
    console.error(`Failed to get QR token for client ${client.id}:`, err);
    qrError.value = (err.response && err.response.data && err.response.data.error) || 'Не удалось получить ссылку для QR-кода.';
  }
};

const downloadConfig = async (client) => {
  try {
    const response = await api.get(`/api/clients/${client.name}/config`, {
      responseType: 'blob',
    });

    let filename = `${client.name}.${client.type === 'wireguard' ? 'conf' : 'ovpn'}`;
    const contentDisposition = response.headers['content-disposition'];
    if (contentDisposition) {
      const filenameMatch = contentDisposition.match(/filename="([^"]+)"/);
      if (filenameMatch && filenameMatch[1]) {
        filename = filenameMatch[1];
      }
    }

    const url = window.URL.createObjectURL(new Blob([response.data]));
    const a = document.createElement('a');
    a.style.display = 'none';
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    window.URL.revokeObjectURL(url);
    document.body.removeChild(a);
  } catch (err) {
    console.error('Download failed:', err);
    let errorMessage = 'Произошла ошибка при скачивании файла.';
    if (err.response && err.response.data) {
        try {
            const errorBlob = new Blob([err.response.data], { type: 'application/json' });
            const reader = new FileReader();
            reader.onload = function() {
                const errorJson = JSON.parse(reader.result);
                errorMessage = errorJson.error || errorMessage;
                alert(`Не удалось скачать конфигурацию: ${errorMessage}`);
            };
            reader.readAsText(errorBlob);
            return;
        } catch (e) {
            // Игнорируем ошибку парсинга, если ответ не JSON
        }
    }
    alert(`Не удалось скачать конфигурацию: ${errorMessage}`);
  }
};
onMounted(() => fetchClients(1));



const bgColors = [
  'bg-red-500',
  'bg-blue-500',
  'bg-green-500',
  'bg-yellow-500',
  'bg-purple-500',
  'bg-pink-500'
];

// Функция для выбора цвета
const getUserColor = (id) => {
  const index = id % bgColors.length;
  return bgColors[index];
};
</script>
