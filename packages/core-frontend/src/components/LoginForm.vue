<script setup lang="ts">
import { useRouter } from 'vue-router';
import SingleFormCard from './SingleFormCard.vue';
import { nextTick, ref, watch } from "vue";
import { useAuthStore } from "../stores/auth";
import { useToast } from "../composables/useToast";
import { useI18n } from 'vue-i18n';

// Initialize router at top level
const router = useRouter();
const { t } = useI18n();

const authStore = useAuthStore();
const { error: showErrorToast, success: showSuccessToast } = useToast();

const credentials = ref({
  username: '',
  password: ''
});

// Helper function to translate API error messages
const translateMessage = (message: string): string => {
  if (!message) return message;
  
  // Define the message mappings based on our locale structure
  const messageMap: Record<string, string> = {
    'Invalid credentials': t('messages.Invalid credentials'),
    'Invalid request': t('messages.Invalid request'),
    'Terms not accepted': t('messages.Terms not accepted'),
    'User already exists': t('messages.User already exists'),
    'Password mismatch': t('messages.Password mismatch'),
    'Company name required': t('messages.Company name required'),
    'User created successfully. Please check your email to verify your account.': t('messages.User created successfully. Please check your email to verify your account.'),
    'Login successful': t('messages.Login successful'),
    'User not found': t('messages.User not found'),
    'Username or email already taken': t('messages.Username or email already taken'),
    'You must accept the terms and conditions to register': t('messages.You must accept the terms and conditions to register'),
    'Company name is required to create a new tenant': t('messages.Company name is required to create a new tenant'),
    'Login failed': t('messages.loginFailed'),
    'Invalid response from server': t('messages.invalidResponseFromServer')
  };
  
  // Return translated message if available, otherwise return original
  return messageMap[message] || message;
};

// Clear errors when user starts typing
watch(credentials, () => {
  if (authStore.error) {
    authStore.setError(null);
  }
}, { deep: true });

async function handleLogin() {
  try {
    console.log('🔐 Login attempt started')
    // Clear any previous errors
    authStore.setError(null);
    
    console.log('🔐 Calling authStore.login...')
    await authStore.login(credentials.value)
    
    console.log('🔐 Login completed, auth state:', {
      isAuthenticated: authStore.isAuthenticated,
      hasUser: !!authStore.user,
      hasToken: !!authStore.token
    })
    
    // Show success toast
    showSuccessToast(t('messages.loginSuccessful'), t('messages.loginSuccessfulWelcome'));
    
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
    
  } catch (error: any) {
    console.error('🔐 Login failed:', error)
    
    // Show error toast in addition to the inline alert
    const errorMessage = authStore.error || error.message || t('messages.loginFailed');
    const translatedErrorMessage = translateMessage(errorMessage);
    showErrorToast(translatedErrorMessage, t('messages.loginFailed'));
  }
}
</script>

<template>
  <SingleFormCard :title="$t('auth.signIn')" :subtitle="$t('auth.signInSubtitle')">
    <div v-if="authStore.error" class="alert alert-error mb-4 shadow-lg">
      <div class="flex items-start gap-2">
        <span class="text-lg flex-shrink-0">⚠️</span>
        <div class="flex-1">
          <div class="font-bold text-sm mb-1">{{ $t('messages.loginFailed') }}</div>
          <div class="text-sm break-words">{{ translateMessage(authStore.error) }}</div>
        </div>
        <button
          @click="authStore.setError(null)"
          class="btn btn-sm btn-ghost btn-square"
          aria-label="Close error"
        >
          ✕
        </button>
      </div>
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
