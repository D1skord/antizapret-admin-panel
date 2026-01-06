import { createRouter, createWebHistory } from 'vue-router';
import AdminLayout from '../components/layout/AdminLayout.vue';
import { useAuthStore } from '../stores/auth';

// Ленивая загрузка (Lazy Loading) - хорошая привычка даже для малых проектов.
// Компоненты загружаются только когда нужны, ускоряя первую отрисовку.
const Clients = () => import('../views/Clients.vue');
const Login = () => import('../views/Login.vue');
const Profile = () => import('../views/Profile.vue');
const Settings = () => import('../views/Settings.vue');

const routes = [
  {
    path: '/',
    component: AdminLayout,
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Home',
        component: Clients,
        meta: { title: 'Клиенты' },
      },
      {
        path: '/profile',
        name: 'Profile',
        component: Profile,
        meta: { title: 'Профиль' },
      },
      {
        path: '/settings',
        name: 'Settings',
        component: Settings,
        meta: { title: 'Настройки' },
      },
    ]
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { title: 'Вход', guestOnly: true }, // Добавили метку "только для гостей"
  },
  // Обработка 404 (опционально, но полезно)
  {
    path: '/:pathMatch(.*)*',
    redirect: '/'
  }
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  scrollBehavior(to, from, savedPosition) {
    return savedPosition || { left: 0, top: 0 }
  },
  routes,
});

router.beforeEach((to, from, next) => {
  // Pinia уже активна здесь, если в main.js порядок верный.
  // try-catch не нужен.
  const authStore = useAuthStore();

  // Установка заголовка
  document.title = `Antizapret | ${to.meta.title || 'Panel'}`;

  // 1. Проверка авторизации
  // Проверяем, требует ли маршрут (или его родитель) авторизации
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth);

  if (requiresAuth && !authStore.isAuthenticated) {
    // Если страница закрыта, а юзер не вошел -> на логин
    return next('/login');
  }

  // 2. Проверка страницы логина
  // Если юзер уже вошел и пытается открыть логин -> на главную
  if (to.meta.guestOnly && authStore.isAuthenticated) {
    return next('/');
  }

  // Все ок
  next();
});

export default router;