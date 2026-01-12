<template>
  <div class="min-h-screen bg-base-200 flex items-center justify-center p-4">
    <div class="card w-full max-w-md bg-base-100 shadow-xl">
      <div class="card-body">
        <h2 class="card-title justify-center text-2xl mb-4">
          {{ $t('verifyEmail.title') }}
        </h2>

        <!-- Loading State -->
        <div v-if="isVerifying" class="text-center py-8">
          <div class="loading loading-spinner loading-lg text-primary mb-4"></div>
          <p class="text-base-content/70">{{ $t('verifyEmail.verifying') }}</p>
        </div>

        <!-- Success State -->
        <div v-else-if="verified" class="text-center py-8">
          <svg class="w-16 h-16 text-success mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
          </svg>
          <p class="text-lg font-semibold text-success mb-2">{{ $t('verifyEmail.success') }}</p>
          <p class="text-base-content/70 mb-6">{{ $t('verifyEmail.successMessage') }}</p>
          <button @click="goToLogin" class="btn btn-primary">
            {{ $t('verifyEmail.goToLogin') }}
          </button>
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="text-center py-8">
          <svg class="w-16 h-16 text-error mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
          </svg>
          <p class="text-lg font-semibold text-error mb-2">{{ $t('verifyEmail.error') }}</p>
          <p class="text-base-content/70 mb-6">{{ errorMessage }}</p>
          <button @click="goToLogin" class="btn btn-outline">
            {{ $t('verifyEmail.goToLogin') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getApiClient } from '@/config/api-config'
import { useI18n } from 'vue-i18n'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const isVerifying = ref(true)
const verified = ref(false)
const error = ref(false)
const errorMessage = ref('')

const verifyEmail = async () => {
  const token = route.query.token as string

  if (!token) {
    error.value = true
    errorMessage.value = t('verifyEmail.missingToken')
    isVerifying.value = false
    return
  }

  try {
    const apiClient = getApiClient()
    await apiClient.verifyEmail(token)
    verified.value = true
  } catch (err: any) {
    error.value = true
    errorMessage.value = err.message || t('verifyEmail.defaultError')
  } finally {
    isVerifying.value = false
  }
}

const goToLogin = () => {
  router.push('/login')
}

onMounted(() => {
  verifyEmail()
})
</script>
