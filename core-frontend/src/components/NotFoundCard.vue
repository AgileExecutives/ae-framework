<template>
      <div class="max-w-md w-full text-center space-y-6 sm:space-y-8 px-4 sm:px-6">
        <!-- 404 Icon/Number -->
        <div class="space-y-3 sm:space-y-4">
          <div class="text-6xl sm:text-8xl md:text-9xl font-bold text-primary opacity-20">
            404
          </div>
          <h1 class="text-2xl sm:text-3xl md:text-4xl font-bold tracking-tight px-2">
            {{ $t('notFound.title') }}
          </h1>
          <p class="text-base sm:text-lg text-base-content/70 px-2 leading-relaxed">
            {{ $t('notFound.message') }}
          </p>
        </div>

        <!-- Action Buttons -->
        <div class="flex flex-col sm:flex-row gap-3 justify-center px-2">
          <button 
            @click="goBack"
            class="btn btn-outline btn-md sm:btn-lg w-full sm:w-auto min-h-12"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
            </svg>
            {{ $t('notFound.goBack') }}
          </button>
          <button 
            @click="goHome"
            class="btn btn-primary btn-md sm:btn-lg w-full sm:w-auto min-h-12"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
            </svg>
            {{ $t('notFound.goHome') }}
          </button>
        </div>

        <!-- Helpful Links -->
        <div class="pt-6 sm:pt-8 border-t border-base-300 mx-2">
          <p class="text-sm text-base-content/70 mb-4">
            {{ $t('notFound.helpfulLinks') }}
          </p>
          <div class="flex flex-col sm:flex-row gap-3 sm:gap-6 justify-center text-sm">
            <router-link 
              to="/login" 
              class="text-primary hover:text-primary-focus hover:underline transition-colors duration-200 py-2 px-3 rounded-md hover:bg-base-200/50"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 inline mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1" />
              </svg>
              {{ $t('auth.signIn') }}
            </router-link>
            <router-link 
              to="/register" 
              class="text-primary hover:text-primary-focus hover:underline transition-colors duration-200 py-2 px-3 rounded-md hover:bg-base-200/50"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 inline mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
              </svg>
              {{ $t('auth.signUp') }}
            </router-link>
          </div>
        </div>

        <!-- Mobile Enhancement: Add a subtle illustration or additional help -->
        <div class="block sm:hidden pt-4 text-xs text-base-content/50">
          <p>{{ $t('notFound.mobileHint') }}</p>
        </div>
      </div>

</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'

const router = useRouter()
const { t } = useI18n()

const goHome = () => {
  router.push('/')
}

const goBack = () => {
  router.back()
}
</script>
