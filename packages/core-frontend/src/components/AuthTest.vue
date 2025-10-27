<template>
  <div class="max-w-4xl mx-auto p-6 space-y-6">
    <div class="text-center">
      <h2 class="text-3xl font-bold">Authentication Test</h2>
      <p class="text-base-content/70">Test the auth store functionality</p>
    </div>

    <div class="card bg-base-100 shadow-xl">
      <div class="card-body">
        <h3 class="card-title">Auth Status</h3>
      <ul>
        <li><strong>Authenticated:</strong> {{ authStore.isAuthenticated ? 'Yes' : 'No' }}</li>
        <li><strong>Loading:</strong> {{ authStore.isLoading ? 'Yes' : 'No' }}</li>
        <li><strong>Initialized:</strong> {{ authStore.initialized ? 'Yes' : 'No' }}</li>
        <li><strong>User:</strong> {{ authStore.user?.username || 'None' }}</li>
        <li><strong>Role:</strong> {{ authStore.userRole || 'None' }}</li>
        <li><strong>Display Name:</strong> {{ authStore.userDisplayName || 'None' }}</li>
      </ul>
      
              
        <div v-if="authStore.error" class="alert alert-error">
          <strong>Error:</strong> {{ authStore.error }}
        </div>
      </div>
    </div>

    <div class="card bg-base-100 shadow-xl">
      <div class="card-body">
    </div>

    <div class="auth-actions">
      <div v-if="!authStore.isAuthenticated" class="login-form">
        <h3>Login</h3>
        <form @submit.prevent="handleLogin" class="space-y-4">
          <div class="form-control">
            <label class="label" for="username">
              <span class="label-text">Username:</span>
            </label>
            <input 
              id="username"
              v-model="loginForm.username" 
              type="text" 
              required 
              :disabled="authStore.isLoading"
              class="input input-bordered w-full"
            />
          </div>
          
          <div class="form-control">
            <label class="label" for="password">
              <span class="label-text">Password:</span>
            </label>
            <input 
              id="password"
              v-model="loginForm.password" 
              type="password" 
              required 
              :disabled="authStore.isLoading"
              class="input input-bordered w-full"
            />
          </div>
          
          <button 
            type="submit" 
            :disabled="authStore.isLoading"
            class="btn btn-primary w-full"
          >
            <span v-if="authStore.isLoading" class="loading loading-spinner loading-sm"></span>
            {{ authStore.isLoading ? 'Logging in...' : 'Login' }}
          </button>
        </form>
      </div>

      <div v-else class="authenticated-actions">
        <h3>Authenticated Actions</h3>
        <div class="flex gap-2 flex-wrap">
          <button 
            @click="handleRefreshUser"
            :disabled="authStore.isLoading"
            class="btn btn-secondary"
          >
            <span v-if="authStore.isLoading" class="loading loading-spinner loading-sm"></span>
            {{ authStore.isLoading ? 'Refreshing...' : 'Refresh User' }}
          </button>
          
          <button 
            @click="handleLogout"
            :disabled="authStore.isLoading"
            class="btn btn-error"
          >
            <span v-if="authStore.isLoading" class="loading loading-spinner loading-sm"></span>
            {{ authStore.isLoading ? 'Logging out...' : 'Logout' }}
          </button>
        </div>
      </div>
    </div>

    <div class="card bg-base-100 shadow-xl">
      <div class="card-body">
        <h3 class="card-title">Utility Methods</h3>
      <ul>
        <li><strong>Is Admin:</strong> {{ authStore.isAdmin ? 'Yes' : 'No' }}</li>
        <li><strong>Is Super Admin:</strong> {{ authStore.isSuperAdmin ? 'Yes' : 'No' }}</li>
        <li><strong>Can Access Admin:</strong> {{ authStore.canAccessAdmin() ? 'Yes' : 'No' }}</li>
        <li><strong>Has User Role:</strong> {{ authStore.hasRole('user') ? 'Yes' : 'No' }}</li>
        <li><strong>Has Any Admin Role:</strong> {{ authStore.hasAnyRole(['admin', 'super-admin']) ? 'Yes' : 'No' }}</li>
      </ul>
      </div>
    </div>

    <div class="card bg-base-100 shadow-xl">
      <div class="card-body">
        <h3 class="card-title">Debug Info</h3>
        <pre class="mockup-code text-sm">{{ debugInfo }}</pre>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()

const loginForm = ref({
  username: 'testuser', // Default credentials for ae-saas-basic
  password: 'newpass123'
})

const debugInfo = computed(() => ({
  token: authStore.token ? `${authStore.token.substring(0, 20)}...` : null,
  user: authStore.user,
  initialized: authStore.initialized,
  isAuthenticated: authStore.isAuthenticated,
  userRole: authStore.userRole,
  error: authStore.error
}))

const handleLogin = async () => {
  try {
    await authStore.login({
      username: loginForm.value.username,
      password: loginForm.value.password
    })
    console.log('✅ Login successful from component')
  } catch (error) {
    console.error('❌ Login failed in component:', error)
  }
}

const handleLogout = async () => {
  try {
    await authStore.logout()
    console.log('✅ Logout successful from component')
  } catch (error) {
    console.error('❌ Logout failed in component:', error)
  }
}

const handleRefreshUser = async () => {
  try {
    await authStore.refreshAuth()
    console.log('✅ User refresh successful from component')
  } catch (error) {
    console.error('❌ User refresh failed in component:', error)
  }
}

// Initialize auth on component mount
onMounted(async () => {
  try {
    console.log('🔄 Initializing auth from component...')
    await authStore.initializeAuth()
    console.log('✅ Auth initialization complete from component')
  } catch (error) {
    console.error('❌ Auth initialization failed in component:', error)
  }
})
</script>

<style scoped>
/* Custom overrides for specific styling if needed */
.login-form {
  max-width: 400px;
}

.auth-status ul, .auth-utils ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.auth-status li, .auth-utils li {
  padding: 8px 0;
  border-bottom: 1px solid hsl(var(--bc) / 0.1);
}

.auth-status li:last-child, .auth-utils li:last-child {
  border-bottom: none;
}
</style>