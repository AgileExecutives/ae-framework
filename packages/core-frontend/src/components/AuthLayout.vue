<template>
  <div class="relative h-full w-full bg-base-100/40 sm:bg-transparent">
    <!-- Main content area -->
    <div class="h-full w-full overflow-auto sm:flex sm:items-center sm:justify-center sm:px-4 py-16 sm:0">
      <slot />
    </div>
    
    <!-- Floating controls - positioned fixed on mobile, absolute on desktop -->
    <div class="fixed md:absolute top-1 md:top-4 right-1 md:right-4 flex flex-row items-center md:gap-2 z-10 
                bg-accent-200/80 md:bg-transparent backdrop-blur-sm md:backdrop-blur-none 
                rounded-lg md:rounded-none shadow-lg shadow-none">
      <LocaleSwitcher />
      <ThemeToggle />
      <LogoutButton v-if="showLogoutButton" />
    </div>
  </div>
</template>

<script lang="ts" setup>
import LogoutButton from '@/components/LogoutButton.vue';
import ThemeToggle from "@/components/ThemeToggle.vue"
import LocaleSwitcher from "@/components/LocaleSwitcher.vue"
defineProps<{ showLogoutButton?: boolean }>();
</script>
