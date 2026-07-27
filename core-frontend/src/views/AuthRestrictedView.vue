<template>
  <div class="flex flex-col h-full w-full">
    <div class="flex flex-row w-full items-end justify-end p-4 gap-2">
      <LocaleSwitcher />
      <ThemeToggle />
      <LogoutButton />
    </div>
    <div class="flex h-full w-full items-center justify-center px-4 ">
      <h1 class="text-3xl font-bold">Welcome to a restricted Page</h1>
    </div>
  </div>
</template>

<script lang="ts" setup>
import LogoutButton from '@@/components/LogoutButton.vue';
import ThemeToggle from "@@/components/ThemeToggle.vue"
import LocaleSwitcher from "@@/components/LocaleSwitcher.vue"
</script>
