<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { computed, watchEffect, onMounted } from 'vue'

const { locale, availableLocales } = useI18n()

// Computed property to check if current locale is German
const isGerman = computed(() => locale.value === 'de')

// Initialize locale on mount from localStorage
onMounted(() => {
  const savedLocale = localStorage.getItem('locale')
  if (savedLocale && availableLocales.includes(savedLocale)) {
    locale.value = savedLocale
  }
})

// Function to toggle between English and German
const toggleLocale = () => {
  locale.value = isGerman.value ? 'en' : 'de'
}

// Save locale preference to localStorage
watchEffect(() => {
  localStorage.setItem('locale', locale.value)
})
</script>

<template>
  <!-- DaisyUI Swap Component for Locale Toggle -->
  <label class="swap btn btn-ghost btn-circle" title="Switch Language">
    <!-- Locale controller checkbox -->
    <input
      type="checkbox"
      :checked="isGerman"
      @change="toggleLocale"
      class="sr-only"
    />
    
    <!-- German Flag (when German is active) -->
    <svg v-if="isGerman"
      class="swap-on size-5 opacity-90"
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 5 3"
    >
      <!-- German flag: More visible muted colors -->
      <rect width="5" height="1" fill="#1a1a1a"/>
      <rect width="5" height="1" y="1" fill="#b85450"/>
      <rect width="5" height="1" y="2" fill="#d4af37"/>
    </svg>
    
    <!-- US/UK Flag (when English is active) -->
    <svg v-else
      class="swap-off size-5 opacity-85"
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 12 6"
    >
      <!-- Simplified US flag with more colorful but still muted colors -->
      <rect width="12" height="6" fill="#ffffff"/>
      <rect width="12" height="1" y="1" fill="#c85a54"/>
      <rect width="12" height="1" y="3" fill="#c85a54"/>
      <rect width="12" height="1" y="5" fill="#c85a54"/>
      <rect width="5" height="3" fill="#4a5568"/>
      <!-- Simplified stars -->
      <circle cx="1.5" cy="1" r="0.3" fill="#ffffff"/>
      <circle cx="3.5" cy="1" r="0.3" fill="#ffffff"/>
      <circle cx="2.5" cy="2" r="0.3" fill="#ffffff"/>
    </svg>
    
    <span class="sr-only">Switch Language</span>
  </label>
</template>
