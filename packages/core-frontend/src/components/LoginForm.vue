<script setup lang="ts">
import { useRouter } from 'vue-router';
import SingleFormCard from './SingleFormCard.vue';
import { nextTick, ref } from "vue";
import { useAuthStore } from "../stores/auth";

// Initialize router at top level
const router = useRouter();

const authStore = useAuthStore(); // Assuming you have a Pinia store for auth
const credentials = ref({
  username: '',
  password: ''
});

async function handleLogin() {
  try {
    console.log('🔐 Login attempt started')
    //authStore.clearError()
    
    console.log('🔐 Calling authStore.login...')
    await authStore.login(credentials.value)
    
    console.log('🔐 Login completed, auth state:', {
      isAuthenticated: authStore.isAuthenticated,
      hasUser: !!authStore.user,
      hasToken: !!authStore.token
    })
    
    // Ensure Vue reactivity system processes all updates
    await nextTick()
    
    console.log('🔐 After nextTick, auth state:', {
      isAuthenticated: authStore.isAuthenticated,
      hasUser: !!authStore.user,
      hasToken: !!authStore.token
    })
    
    console.log('🔐 Attempting navigation to dashboard or redirect...')
    const redirectPath = router.currentRoute.value.query.redirect || '/';
    const result = await router.push(redirectPath as string)
    console.log('🔐 Navigation result:', result)
    
  } catch (error) {
    console.error('🔐 Login failed:', error)
    // Error is handled by the auth store
  }
}
</script>

<template>
  <SingleFormCard :title="$t('auth.signIn')" :subtitle="$t('auth.signInSubtitle')">
    <div v-if="authStore.error" class="alert alert-error mb-4">
      <span class="break-words">{{ authStore.error }}</span>
    </div>

    <form @submit.prevent="handleLogin" autocomplete="on" class="space-y-6">
      <!-- Authentication Credentials -->
      <fieldset class="space-y-4">
        <legend class="text-base font-semibold text-base-content mb-4 break-words text-wrap">Authentication Credentials</legend>
        
        <div class="form-control">
          <label class="label" for="email">
            <span class="label-text">{{ $t('auth.email') }}</span>
          </label>
          <input
            id="email"
            name="email"
            type="email"
            placeholder="m@example.com"
            autocomplete="email"
            required
            data-testid="login-email"
            v-model="credentials.username"
            class="input input-bordered w-full"
          />
        </div>

        <div class="form-control">
          <div class="flex items-center justify-between mb-1">
            <label class="label" for="password">
              <span class="label-text">{{ $t('auth.password') }}</span>
            </label>
            <router-link to="/forgot-password" class="link link-primary link-hover text-sm">
              {{ $t('auth.forgotPassword') }}
            </router-link>
          </div>
          <input 
            id="password"
            name="password"
            type="password"
            autocomplete="current-password"
            required 
            data-testid="login-password"
            v-model="credentials.password"
            class="input input-bordered w-full"
          />
        </div>
      </fieldset>
      <button 
        type="submit" 
        class="btn btn-primary w-full mt-8"
        :disabled="authStore.isLoading"
        data-testid="login-submit"
      >
        <span v-if="authStore.isLoading" class="loading loading-spinner loading-sm"></span>
        {{ $t('auth.signIn') }}
      </button>
    </form>

    <div class="divider">{{ $t('auth.or') }}</div>

    <div class="text-center text-sm">
      {{ $t('auth.noAccount') }}
      <router-link to="/register" class="link link-primary">
        {{ $t('auth.signUp') }}
      </router-link>
    </div>
  </SingleFormCard>
</template>
