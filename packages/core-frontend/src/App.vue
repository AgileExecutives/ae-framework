<script setup lang="ts">
import ToastContainer from './components/ToastContainer.vue';
</script>

<template>
  <div class="flex bg-gradient-unburdy">
    <RouterView />
    <ToastContainer />
  </div>
</template>
