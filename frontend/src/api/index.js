import axios from 'axios';
import { useAuthStore } from '../stores/auth';

const api = axios.create({
  baseURL: '/',
});

api.interceptors.request.use(config => {
  const authStore = useAuthStore();
  if (authStore.token) {
    config.headers['X-Auth-Token'] = authStore.token;
  }
  return config;
});

api.interceptors.response.use(
  response => response,
  error => {
    if (error.response && error.response.status === 401) {
      const authStore = useAuthStore();
      authStore.logout();
    }
    return Promise.reject(error);
  }
);

export default api;
