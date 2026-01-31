<template>
  <AuthLayout>
    <div class="card w-full max-w-md bg-base-100 shadow-xl">
      <div class="card-body items-center text-center">
        <!-- Error Icon -->
        <div class="mb-6">
          <svg 
            class="w-20 h-20 text-error" 
            fill="none" 
            stroke="currentColor" 
            viewBox="0 0 24 24"
          >
            <path 
              stroke-linecap="round" 
              stroke-linejoin="round" 
              stroke-width="2" 
              d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
        </div>

        <!-- Title -->
        <h1 class="card-title text-2xl mb-2" data-testid="error-title">
          {{ title }}
        </h1>

        <!-- Message -->
        <p class="text-base-content/70 mb-6" data-testid="error-message">
          {{ message }}
        </p>

        <!-- Actions -->
        <div class="card-actions flex-col w-full gap-2">
          <!-- Back Button -->
          <button
            @click="goBack"
            class="btn btn-primary w-full"
            data-testid="error-back-button"
          >
            {{ $t('error.goBack') }}
          </button>

          <!-- Home Link -->
          <RouterLink
            to="/"
            class="btn btn-ghost w-full"
            data-testid="error-home-link"
          >
            {{ $t('error.goHome') }}
          </RouterLink>
        </div>
      </div>
    </div>
  </AuthLayout>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter, RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'
import AuthLayout from '@@/components/AuthLayout.vue'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

// Get parameters from route query
const customTitle = computed(() => route.query.title as string)
const customMessage = computed(() => route.query.message as string)

const title = computed(() => customTitle.value || t('error.default.title'))
const message = computed(() => customMessage.value || t('error.default.message'))

const goBack = () => {
  router.back()
}
</script>
