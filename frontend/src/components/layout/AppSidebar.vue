<template>
  <aside
    :class="[
      'fixed mt-16 flex flex-col lg:mt-0 top-0 px-5 left-0 bg-white dark:bg-gray-900 dark:border-gray-800 text-gray-900 h-screen transition-all duration-300 ease-in-out z-40 border-r border-gray-200',
      {
        'lg:w-[290px]': isExpanded || isMobileOpen || isHovered,
        'lg:w-[90px]': !isExpanded && !isHovered,
        'translate-x-0 w-[290px]': isMobileOpen,
        '-translate-x-full': !isMobileOpen,
        'lg:translate-x-0': true,
      },
    ]"
    @mouseenter="!isExpanded && (isHovered = true)"
    @mouseleave="isHovered = false"
  >
    <div
      :class="[
        'py-8 flex',
        !isExpanded && !isHovered ? 'lg:justify-center' : 'justify-start',
      ]"
    >
      <router-link to="/">
        <img
          v-if="isExpanded || isHovered || isMobileOpen"
          class="dark:hidden"
          src="/images/logo/logo.svg"
          alt="Logo"
          width="150"
          height="40"
        />
        <img
          v-if="isExpanded || isHovered || isMobileOpen"
          class="hidden dark:block"
          src="/images/logo/logo-dark.svg"
          alt="Logo"
          width="150"
          height="40"
        />
        <img
          v-else
          src="/images/logo/logo-icon.svg"
          alt="Logo"
          width="32"
          height="32"
        />
      </router-link>
    </div>
    <div
      class="flex flex-col overflow-y-auto duration-300 ease-linear no-scrollbar"
    >
      <nav class="mb-6">
        <div class="flex flex-col gap-4">
          <div>
            <h2
              :class="[
                'mb-4 text-xs uppercase flex leading-[20px] text-gray-400',
                !isExpanded && !isHovered
                  ? 'lg:justify-center'
                  : 'justify-start',
              ]"
            >
              <template v-if="isExpanded || isHovered || isMobileOpen">
                Меню
              </template>
              <HorizontalDots v-else />
            </h2>
            <ul class="flex flex-col gap-4">
              <li>
                <router-link
                  to="/"
                  :class="[
                    'menu-item group',
                    {
                      'menu-item-active': isActive('/'),
                      'menu-item-inactive': !isActive('/'),
                    },
                  ]"
                >
                  <span
                    :class="[
                      isActive('/')
                        ? 'menu-item-icon-active'
                        : 'menu-item-icon-inactive',
                    ]"
                  >
                    <GridIcon />
                  </span>
                  <span
                    v-if="isExpanded || isHovered || isMobileOpen"
                    class="menu-item-text"
                    >Клиенты</span
                  >
                </router-link>
              </li>
              <li>
                <router-link
                  to="/profile"
                  :class="[
                    'menu-item group',
                    {
                      'menu-item-active': isActive('/profile'),
                      'menu-item-inactive': !isActive('/profile'),
                    },
                  ]"
                >
                  <span
                    :class="[
                      isActive('/profile')
                        ? 'menu-item-icon-active'
                        : 'menu-item-icon-inactive',
                    ]"
                  >
                    <UserCircleIcon />
                  </span>
                  <span
                    v-if="isExpanded || isHovered || isMobileOpen"
                    class="menu-item-text"
                    >Профиль</span
                  >
                </router-link>
              </li>
              <li>
                <router-link
                  to="/settings"
                  :class="[
                    'menu-item group',
                    {
                      'menu-item-active': isActive('/settings'),
                      'menu-item-inactive': !isActive('/settings'),
                    },
                  ]"
                >
                  <span
                    :class="[
                      isActive('/settings')
                        ? 'menu-item-icon-active'
                        : 'menu-item-icon-inactive',
                    ]"
                  >
                    <SettingsIcon />
                  </span>
                  <span
                    v-if="isExpanded || isHovered || isMobileOpen"
                    class="menu-item-text"
                    >Настройки</span
                  >
                </router-link>
              </li>
            </ul>
          </div>
        </div>
      </nav>
      <SidebarWidget v-if="isExpanded || isHovered || isMobileOpen" />
    </div>
  </aside>
</template>

<script setup>
import { useRoute } from "vue-router";
import { GridIcon, HorizontalDots, UserCircleIcon, SettingsIcon } from "@/icons";
import SidebarWidget from "./SidebarWidget.vue";
import { useSidebar } from "@/composables/useSidebar";

const route = useRoute();
const { isExpanded, isMobileOpen, isHovered } = useSidebar();
const isActive = (path) => route.path === path;
</script>

<style>
/* Basic styles for menu items, you might need to copy more from tailAdmin's css if something is missing */
.menu-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0.5rem 1rem;
  border-radius: 0.5rem;
  transition: all 0.2s;
}
.menu-item-text {
  font-size: 0.875rem;
  line-height: 1.25rem;
  font-weight: 500;
}
.menu-item-active {
  background-color: #4f46e5;
  color: white;
}
.menu-item-icon-active {
  color: white;
}
.menu-item-inactive {
  color: #6b7280;
}
.menu-item-inactive:hover {
  background-color: #f3f4f6;
}
.dark .menu-item-inactive:hover {
    background-color: #1f2937;
}
</style>
