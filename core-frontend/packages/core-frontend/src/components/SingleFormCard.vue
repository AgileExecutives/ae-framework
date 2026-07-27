<script setup lang="ts">
import { computed, ref } from "vue";
import { useAppConfig } from '../composables/useAppConfig';
import { useDarkMode } from '../composables/useDarkMode';

defineProps<{
  title: string,
  subtitle: string,
}>()

const { isDark } = useDarkMode()
const { currentLogo, currentLogoDark } = useAppConfig()

// Compute the appropriate logo based on theme
const logoUrl = computed(() => {
  return isDark.value ? currentLogoDark.value : currentLogo.value
})

</script>

<template>
  <!-- DaisyUI Card Component -->
  <div class="card w-full max-w-sm md:max-w-md lg:max-w-lg sm:bg-base-100 sm:shadow-xl mx-auto">
    <div class="card-body">
      <div class="text-center mb-6 md:mb-8">
        <div class="mx-auto h-10 md:h-12 lg:h-14 w-auto flex justify-center mb-6 md:mb-8">
          <img class="h-10 md:h-12 lg:h-14 w-auto" :src="logoUrl" alt="Unburdy" />
        </div>
        <h2 class="card-title text-xl md:text-2xl lg:text-3xl mb-3 md:mb-4 font-bold text-center break-words justify-center">
          {{ title }}
        </h2>
        <p class="text-base-content/70 break-words text-sm md:text-base leading-relaxed text-center">
          {{ subtitle }}
        </p>
      </div>
      <slot />
    </div>
  </div>
</template>
