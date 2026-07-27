<template>
  <AuthLayout>
    <div class="card w-full max-w-md bg-base-100 shadow-xl">
      <div class="card-body items-center text-center">
        <!-- Success Icon -->
        <div class="mb-6">
          <svg 
            class="w-20 h-20 text-success" 
            fill="none" 
            stroke="currentColor" 
            viewBox="0 0 24 24"
          >
            <path 
              stroke-linecap="round" 
              stroke-linejoin="round" 
              stroke-width="2" 
              d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
        </div>

        <!-- Title -->
        <h1 class="card-title text-2xl mb-2" data-testid="success-title">
          {{ title }}
        </h1>

        <!-- Message -->
        <p class="text-base-content/70 mb-6" data-testid="success-message">
          {{ message }}
        </p>

        <!-- Actions -->
        <div class="card-actions flex-col w-full gap-2">
          <!-- Primary Action (if provided) -->
          <button
            v-if="primaryAction"
            @click="handlePrimaryAction"
            class="btn btn-primary w-full"
            data-testid="success-primary-action"
          >
            {{ primaryAction.label }}
          </button>

          <!-- Login Link -->
          <RouterLink
            to="/login"
            class="btn btn-outline w-full"
            data-testid="success-login-link"
          >
            {{ $t('success.goToLogin') }}
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
const type = computed(() => route.query.type as string || 'default')
const customTitle = computed(() => route.query.title as string)
const customMessage = computed(() => route.query.message as string)
const actionUrl = computed(() => route.query.actionUrl as string)
const actionLabel = computed(() => route.query.actionLabel as string)

// Determine title and message based on type
const title = computed(() => {
  if (customTitle.value) return customTitle.value
  
  switch (type.value) {
    case 'registration':
      return t('success.registration.title')
    case 'password-reset':
      return t('success.passwordReset.title')
    case 'password-change':
      return t('success.passwordChange.title')
    case 'email-verified':
      return t('success.emailVerified.title')
    default:
      return t('success.default.title')
  }
})

const message = computed(() => {
  if (customMessage.value) return customMessage.value
  
  switch (type.value) {
    case 'registration':
      return t('success.registration.message')
    case 'password-reset':
      return t('success.passwordReset.message')
    case 'password-change':
      return t('success.passwordChange.message')
    case 'email-verified':
      return t('success.emailVerified.message')
    default:
      return t('success.default.message')
  }
})

const primaryAction = computed(() => {
  if (actionUrl.value && actionLabel.value) {
    return {
      label: actionLabel.value,
      url: actionUrl.value
    }
  }
  return null
})

const handlePrimaryAction = () => {
  if (primaryAction.value?.url) {
    router.push(primaryAction.value.url)
  }
}
</script>
