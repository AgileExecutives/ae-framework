<script setup lang="ts">
import ToastContainer from '@@/components/ToastContainer.vue';

// TODO: Auth Restricted View should not be part of framework, but part of the testing 
// TODO: API Client should be part of framework
</script>

<template>
  <div class="flex bg-gradient-unburdy">
    <RouterView />
    <ToastContainer />
  </div>
</template>

