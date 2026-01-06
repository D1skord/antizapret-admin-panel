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
            <th class="px-5 py-3 text-left sm:px-6">
              <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">Пользователь</p>
            </th>
            <th class="px-5 py-3 text-left sm:px-6">
              <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">Статус</p>
            </th>
            <th class="px-5 py-3 text-left sm:px-6">
              <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">Дата создания</p>
            </th>
            <th class="px-5 py-3 text-center sm:px-6">
              <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">OpenVPN</p>
            </th>
            <th class="px-5 py-3 text-center sm:px-6">
              <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">Antizapret</p>
            </th>
            <th class="px-5 py-3 text-center sm:px-6 w-16">
              <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400"></p>
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
            <!-- Обычный VPN -->
            <td class="px-5 py-4 sm:px-6">
              <div class="flex items-center justify-center gap-2">
                <button @click="downloadConfig(user, 'vpn')" class="p-1.5 text-green-500 hover:text-green-700 rounded-md hover:bg-green-50 dark:hover:bg-green-500/10" title="Скачать VPN конфиг">
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-5 h-5">
                    <path d="M10.75 2.75a.75.75 0 0 0-1.5 0v8.614L6.295 8.235a.75.75 0 1 0-1.09 1.03l4.25 4.5a.75.75 0 0 0 1.09 0l4.25-4.5a.75.75 0 0 0-1.09-1.03l-2.955 3.129V2.75Z" />
                    <path d="M3.5 12.75a.75.75 0 0 0-1.5 0v2.5A2.75 2.75 0 0 0 4.75 18h10.5A2.75 2.75 0 0 0 18 15.25v-2.5a.75.75 0 0 0-1.5 0v2.5c0 .69-.56 1.25-1.25 1.25H4.75c-.69 0-1.25-.56-1.25-1.25v-2.5Z" />
                  </svg>
                </button>
                <button @click="showQrCode(user, 'vpn')" class="p-1.5 text-blue-500 hover:text-blue-700 rounded-md hover:bg-blue-50 dark:hover:bg-blue-500/10" title="QR-код VPN конфига">
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-5 h-5">
                    <path fill-rule="evenodd" d="M3.75 2A1.75 1.75 0 0 0 2 3.75v3.5C2 8.216 2.784 9 3.75 9h3.5A1.75 1.75 0 0 0 9 7.25v-3.5A1.75 1.75 0 0 0 7.25 2h-3.5ZM3.5 3.75a.25.25 0 0 1 .25-.25h3.5a.25.25 0 0 1 .25.25v3.5a.25.25 0 0 1-.25.25h-3.5a.25.25 0 0 1-.25-.25v-3.5ZM3.75 11A1.75 1.75 0 0 0 2 12.75v3.5c0 .966.784 1.75 1.75 1.75h3.5A1.75 1.75 0 0 0 9 16.25v-3.5A1.75 1.75 0 0 0 7.25 11h-3.5Zm-.25 1.75a.25.25 0 0 1 .25-.25h3.5a.25.25 0 0 1 .25.25v3.5a.25.25 0 0 1-.25.25h-3.5a.25.25 0 0 1-.25-.25v-3.5ZM12.75 2A1.75 1.75 0 0 0 11 3.75v3.5c0 .966.784 1.75 1.75 1.75h3.5A1.75 1.75 0 0 0 18 7.25v-3.5A1.75 1.75 0 0 0 16.25 2h-3.5Zm-.25 1.75a.25.25 0 0 1 .25-.25h3.5a.25.25 0 0 1 .25.25v3.5a.25.25 0 0 1-.25.25h-3.5a.25.25 0 0 1-.25-.25v-3.5ZM11 12.75a.75.75 0 0 1 1.5 0v1.5h1.5a.75.75 0 0 1 0 1.5H12v3.5a.75.75 0 0 1-1.5 0v-3.5h-.75a.75.75 0 0 1 0-1.5h1.25v-1.5ZM15.5 13.5a.75.75 0 0 1 .75-.75h1a.75.75 0 0 1 .75.75v3a.75.75 0 0 1-.75.75h-1a.75.75 0 0 1-.75-.75v-3Zm-2 4.25a.75.75 0 0 1 .75-.75h3.5a.75.75 0 0 1 0 1.5h-3.5a.75.75 0 0 1-.75-.75Z" clip-rule="evenodd" />
                  </svg>
                </button>
              </div>
            </td>
            <!-- Антизапрет -->
            <td class="px-5 py-4 sm:px-6">
              <div class="flex items-center justify-center gap-2">
                <button @click="downloadConfig(user, 'antizapret')" class="p-1.5 text-green-500 hover:text-green-700 rounded-md hover:bg-green-50 dark:hover:bg-green-500/10" title="Скачать Антизапрет конфиг">
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-5 h-5">
                    <path d="M10.75 2.75a.75.75 0 0 0-1.5 0v8.614L6.295 8.235a.75.75 0 1 0-1.09 1.03l4.25 4.5a.75.75 0 0 0 1.09 0l4.25-4.5a.75.75 0 0 0-1.09-1.03l-2.955 3.129V2.75Z" />
                    <path d="M3.5 12.75a.75.75 0 0 0-1.5 0v2.5A2.75 2.75 0 0 0 4.75 18h10.5A2.75 2.75 0 0 0 18 15.25v-2.5a.75.75 0 0 0-1.5 0v2.5c0 .69-.56 1.25-1.25 1.25H4.75c-.69 0-1.25-.56-1.25-1.25v-2.5Z" />
                  </svg>
                </button>
                <button @click="showQrCode(user, 'antizapret')" class="p-1.5 text-blue-500 hover:text-blue-700 rounded-md hover:bg-blue-50 dark:hover:bg-blue-500/10" title="QR-код Антизапрет конфига">
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-5 h-5">
                    <path fill-rule="evenodd" d="M3.75 2A1.75 1.75 0 0 0 2 3.75v3.5C2 8.216 2.784 9 3.75 9h3.5A1.75 1.75 0 0 0 9 7.25v-3.5A1.75 1.75 0 0 0 7.25 2h-3.5ZM3.5 3.75a.25.25 0 0 1 .25-.25h3.5a.25.25 0 0 1 .25.25v3.5a.25.25 0 0 1-.25.25h-3.5a.25.25 0 0 1-.25-.25v-3.5ZM3.75 11A1.75 1.75 0 0 0 2 12.75v3.5c0 .966.784 1.75 1.75 1.75h3.5A1.75 1.75 0 0 0 9 16.25v-3.5A1.75 1.75 0 0 0 7.25 11h-3.5Zm-.25 1.75a.25.25 0 0 1 .25-.25h3.5a.25.25 0 0 1 .25.25v3.5a.25.25 0 0 1-.25.25h-3.5a.25.25 0 0 1-.25-.25v-3.5ZM12.75 2A1.75 1.75 0 0 0 11 3.75v3.5c0 .966.784 1.75 1.75 1.75h3.5A1.75 1.75 0 0 0 18 7.25v-3.5A1.75 1.75 0 0 0 16.25 2h-3.5Zm-.25 1.75a.25.25 0 0 1 .25-.25h3.5a.25.25 0 0 1 .25.25v3.5a.25.25 0 0 1-.25.25h-3.5a.25.25 0 0 1-.25-.25v-3.5ZM11 12.75a.75.75 0 0 1 1.5 0v1.5h1.5a.75.75 0 0 1 0 1.5H12v3.5a.75.75 0 0 1-1.5 0v-3.5h-.75a.75.75 0 0 1 0-1.5h1.25v-1.5ZM15.5 13.5a.75.75 0 0 1 .75-.75h1a.75.75 0 0 1 .75.75v3a.75.75 0 0 1-.75.75h-1a.75.75 0 0 1-.75-.75v-3Zm-2 4.25a.75.75 0 0 1 .75-.75h3.5a.75.75 0 0 1 0 1.5h-3.5a.75.75 0 0 1-.75-.75Z" clip-rule="evenodd" />
                  </svg>
                </button>
              </div>
            </td>
            <!-- Удалить -->
            <td class="px-5 py-4 sm:px-6">
              <div class="flex items-center justify-center">
                <button @click="deleteClient(user.id)" class="p-1.5 text-error-500 hover:text-error-700 rounded-md hover:bg-error-50 dark:hover:bg-error-500/10" title="Удалить клиента">
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-5 h-5">
                    <path fill-rule="evenodd" d="M8.75 1A2.75 2.75 0 0 0 6 3.75v.443c-.795.077-1.584.176-2.365.298a.75.75 0 1 0 .23 1.482l.149-.022.841 10.518A2.75 2.75 0 0 0 7.596 19h4.807a2.75 2.75 0 0 0 2.742-2.53l.841-10.519.149.023a.75.75 0 0 0 .23-1.482A41.03 41.03 0 0 0 14 4.193V3.75A2.75 2.75 0 0 0 11.25 1h-2.5ZM10 4c.84 0 1.673.025 2.5.075V3.75c0-.69-.56-1.25-1.25-1.25h-2.5c-.69 0-1.25.56-1.25 1.25v.325C8.327 4.025 9.16 4 10 4ZM8.58 7.72a.75.75 0 0 0-1.5.06l.3 7.5a.75.75 0 1 0 1.5-.06l-.3-7.5Zm4.34.06a.75.75 0 1 0-1.5-.06l-.3 7.5a.75.75 0 1 0 1.5.06l.3-7.5Z" clip-rule="evenodd" />
                  </svg>
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
             <label for="clientExpiration" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Срок действия (дней)</label>
             <input type="number" id="clientExpiration" v-model="newClientExpiration" class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white" placeholder="3650" min="1" max="3650">
             <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">От 1 до 3650. По умолчанию: 3650.</p>
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
            <h2 class="text-xl font-bold mb-4">{{ qrModalTitle }}</h2>
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
import { ref, onMounted, computed } from 'vue';
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
const perPage = 10;

// State for Add Client Modal
const isModalOpen = ref(false);
const newClientName = ref('');
const newClientExpiration = ref(3650);
const createError = ref(null);

// State for QR Code Modal
const isQrModalOpen = ref(false);
const qrCodeUrl = ref('');
const qrError = ref(null);
const selectedClientName = ref('');
const selectedConfigType = ref('vpn');

const configTypeLabels = {
  vpn: 'VPN',
  antizapret: 'Антизапрет',
};

const qrModalTitle = computed(() => {
  const typeLabel = configTypeLabels[selectedConfigType.value] || selectedConfigType.value;
  return `${typeLabel} — ${selectedClientName.value}`;
});

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
    const payload = {
      name: newClientName.value,
      type: 'openvpn',
      expires_in: parseInt(newClientExpiration.value, 10) || 3650,
    };

    await api.post('/api/clients', payload);
    isModalOpen.value = false;
    newClientName.value = '';
    newClientExpiration.value = 3650;
    await fetchClients(1);
  } catch (err) {
    console.error('Failed to create client:', err);
    createError.value = (err.response && err.response.data && err.response.data.details) || err.message || 'Не удалось создать клиента.';
  }
};

const deleteClient = async (id) => {
  if (confirm('Вы уверены, что хотите удалить этого клиента?')) {
    try {
      await api.delete(`/api/clients/${id}`);
      if (users.value.length === 1 && currentPage.value > 1) {
        await fetchClients(currentPage.value - 1);
      } else {
        await fetchClients(currentPage.value);
      }
    } catch (err) {
      console.error(`Failed to delete client ${id}:`, err);
      alert('Не удалось удалить клиента.');
    }
  }
};

const showQrCode = async (client, configType) => {
  selectedClientName.value = client.name;
  selectedConfigType.value = configType;
  isQrModalOpen.value = true;
  qrCodeUrl.value = '';
  qrError.value = null;

  try {
    const response = await api.get(`/api/clients/${client.id}/qr-token?type=${configType}`);
    qrCodeUrl.value = window.location.origin + response.data.download_url;
  } catch (err) {
    console.error(`Failed to get QR token for client ${client.id}:`, err);
    qrError.value = (err.response && err.response.data && err.response.data.error) || 'Не удалось получить ссылку для QR-кода.';
  }
};

const downloadConfig = async (client, configType) => {
  try {
    const response = await api.get(`/api/clients/${client.name}/config?type=${configType}`, {
      responseType: 'blob',
    });

    let filename = `${client.name}.ovpn`;
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

const getUserColor = (id) => {
  const index = id % bgColors.length;
  return bgColors[index];
};
</script>
