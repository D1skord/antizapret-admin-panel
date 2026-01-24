import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import api from '../api';
import router from '../router';

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || null);

  const isAuthenticated = computed(() => !!token.value);

  function setToken(newToken) {
    token.value = newToken;
    if (newToken) {
      localStorage.setItem('token', newToken);
    } else {
      localStorage.removeItem('token');
    }
  }

  async function login(username, password) {
    try {
      const response = await api.post('/api/login', { username, password });
      setToken(response.data.token);
      await router.push('/');
    } catch (error) {
      console.error('Login failed:', error);
      setToken(null);
      throw error;
    }
  }

  function logout() {
    setToken(null);
    router.push('/login');
  }

  async function checkAuth() {
    if (!token.value) {
      return;
    }
    try {
      await api.get('/api/check-auth');
    } catch (error) {
      if (error.response && error.response.status === 401) {
        logout();
      }
    }
  }

  return { token, isAuthenticated, login, logout, checkAuth, setToken };
});
