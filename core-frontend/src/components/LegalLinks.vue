<script setup lang="ts">
import { useAppConfig } from '../composables/useAppConfig'

const { 
  termsOfServiceUrl, 
  termsOfServiceText, 
  privacyPolicyUrl, 
  privacyPolicyText 
} = useAppConfig()

defineProps<{
  inline?: boolean
  separator?: string
}>()
</script>

<template>
  <div :class="inline ? 'inline-flex items-center gap-2' : 'space-y-2'">
    <a 
      :href="termsOfServiceUrl" 
      target="_blank"
      rel="noopener noreferrer"
      class="link link-primary text-sm hover:underline"
    >
      {{ termsOfServiceText }}
    </a>
    
    <span v-if="inline" class="text-base-content/50">{{ separator || '•' }}</span>
    
    <a 
      :href="privacyPolicyUrl" 
      target="_blank"
      rel="noopener noreferrer" 
      class="link link-primary text-sm hover:underline"
    >
      {{ privacyPolicyText }}
    </a>
  </div>
</template>