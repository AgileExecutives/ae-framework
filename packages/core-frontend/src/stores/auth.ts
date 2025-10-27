import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { createApiClient, ApiError } from '@agile-exec/api-client'
import type { AESaasApiClient } from '@agile-exec/api-client'
import type * as AuthTypes from './auth-types'

export const useAuthStore = defineStore('auth', () => {
  // API Client instance
  const apiClient = ref<AESaasApiClient | null>(null)
  
  // State
  const user = ref<AuthTypes.User | null>(null)
  const token = ref<string | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)
  const initialized = ref(false)

  // Initialize API client
  const initializeApiClient = () => {
    if (!apiClient.value) {
      console.log('🔍 Auth Store: Creating new API client')
      apiClient.value = createApiClient({
        baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'
      })
      
      // If we have a token, set it on the newly created client
      if (token.value) {
        console.log('🔍 Auth Store: Setting token on new client')
        apiClient.value.setToken(token.value)
      }
    }
    return apiClient.value
  }

  // Getters
  const isAuthenticated = computed(() => !!token.value && !!user.value)
  const userRole = computed(() => user.value?.role || null)
  const isAdmin = computed(() => userRole.value === 'admin' || userRole.value === 'super-admin')
  const isSuperAdmin = computed(() => userRole.value === 'super-admin')
  const userDisplayName = computed(() => {
    if (!user.value) return null
    const { first_name, last_name, username } = user.value
    if (first_name && last_name) {
      return `${first_name} ${last_name}`
    }
    return username
  })

  // Actions
  const setToken = (newToken: string | null) => {
    token.value = newToken
    const client = initializeApiClient()
    
    console.log('🔍 Auth Store: setToken called', {
      hasToken: !!newToken,
      tokenPreview: newToken ? `${newToken.substring(0, 20)}...` : 'null',
      hasClient: !!client
    })
    
    if (newToken) {
      localStorage.setItem('auth_token', newToken)
      client.setToken(newToken)
      console.log('🔍 Auth Store: Token set on client')
    } else {
      localStorage.removeItem('auth_token')
      client.clearToken()
      console.log('🔍 Auth Store: Token cleared from client')
    }
  }

  const setUser = (newUser: AuthTypes.User | null) => {
    user.value = newUser
  }

  const setError = (newError: string | null) => {
    error.value = newError
  }

  const clearAuth = () => {
    setUser(null)
    setToken(null)
    setError(null)
    initialized.value = false
  }

  const login = async (credentials: AuthTypes.LoginCredentials): Promise<void> => {
    try {
      isLoading.value = true
      setError(null)

      const client = initializeApiClient()
      const response = await client.login(credentials)

      if (response.token && response.user) {
        setToken(response.token)
        setUser(response.user as AuthTypes.User)
        console.log('✅ Auth Store: Login successful', { user: response.user.username })
      } else {
        throw new Error('Invalid response from server')
      }
    } catch (err: any) {
      console.error('❌ Auth Store: Login failed', err)
      
      if (err instanceof ApiError) {
        setError(err.data?.message || err.message || 'Login failed')
      } else {
        setError(err.message || 'Login failed')
      }
      
      clearAuth()
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const register = async (credentials: AuthTypes.RegisterCredentials): Promise<AuthTypes.User> => {
    try {
      isLoading.value = true
      setError(null)

      const client = initializeApiClient()
      const response = await client.register(credentials)
      console.log('✅ Auth Store: Registration response received', response)
      
      // Handle the response structure from the API client
      if (response && response.user) {
        // If token is provided, set it (auto-login after registration)
        if (response.token) {
          setToken(response.token)
          setUser(response.user as AuthTypes.User)
          console.log('✅ Auth Store: Registration successful with auto-login', { user: response.user.username })
        } else {
          // Registration successful but no auto-login
          console.log('✅ Auth Store: Registration successful, please login', { user: response.user.username })
        }
        return response.user as AuthTypes.User
      } else {
        throw new Error('Invalid response from server')
      }
    } catch (err: any) {
      console.error('❌ Auth Store: Registration failed', err)
      
      if (err instanceof ApiError) {
        setError(err.data?.message || err.message || 'Registration failed')
      } else {
        setError(err.message || 'Registration failed')
      }
      
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const changePassword = async (credentials: AuthTypes.ChangePasswordCredentials): Promise<void> => {
    try {
      isLoading.value = true
      setError(null)

      console.log('🔍 Auth Store: changePassword called', {
        hasToken: !!token.value,
        tokenValue: token.value ? `${token.value.substring(0, 20)}...` : 'null',
        isAuthenticated: isAuthenticated.value,
        hasUser: !!user.value
      })

      const client = initializeApiClient()
      await client.changePassword(credentials)

      console.log('✅ Auth Store: Password changed successfully')
    } catch (err: any) {
      console.error('❌ Auth Store: Password change failed', err)
      
      if (err instanceof ApiError) {
        setError(err.data?.message || err.message || 'Password change failed')
      } else {
        setError(err.message || 'Password change failed')
      }
      
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const forgotPassword = async (email: string): Promise<void> => {
    try {
      isLoading.value = true
      setError(null)

      const client = initializeApiClient()
      await client.forgotPassword(email)

      console.log('✅ Auth Store: Password reset email sent')
    } catch (err: any) {
      console.error('❌ Auth Store: Forgot password failed', err)
      
      if (err instanceof ApiError) {
        setError(err.data?.message || err.message || 'Failed to send reset email')
      } else {
        setError(err.message || 'Failed to send reset email')
      }
      
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const resetPassword = async (token: string, newPassword: string): Promise<void> => {
    try {
      isLoading.value = true
      setError(null)

      const client = initializeApiClient()
      await client.resetPassword(token, newPassword)

      console.log('✅ Auth Store: Password reset successful')
    } catch (err: any) {
      console.error('❌ Auth Store: Password reset failed', err)
      
      if (err instanceof ApiError) {
        setError(err.data?.message || err.message || 'Password reset failed')
      } else {
        setError(err.message || 'Password reset failed')
      }
      
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const logout = async (): Promise<void> => {
    try {
      isLoading.value = true
      setError(null)

      const client = initializeApiClient()
      
      // Try to logout from server, but don't fail if it doesn't work
      try {
        await client.logout()
      } catch (err) {
        console.warn('⚠️ Auth Store: Server logout failed, clearing local auth anyway', err)
      }

      console.log('✅ Auth Store: Logout successful')
    } catch (err: any) {
      console.error('❌ Auth Store: Logout error', err)
      setError(err.message || 'Logout failed')
    } finally {
      clearAuth()
      isLoading.value = false
    }
  }

  const getCurrentUser = async (): Promise<AuthTypes.User> => {
    try {
      isLoading.value = true
      setError(null)

      const client = initializeApiClient()
      const userData = await client.getCurrentUser()
      
      setUser(userData as AuthTypes.User)
      console.log('✅ Auth Store: User data refreshed', { user: userData.username })
      
      return userData as AuthTypes.User
    } catch (err: any) {
      console.error('❌ Auth Store: Failed to get current user', err)
      
      if (err instanceof ApiError && err.status === 401) {
        // Token is invalid, clear auth
        clearAuth()
      }
      
      setError(err.message || 'Failed to get user data')
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const initializeAuth = async (): Promise<void> => {
    if (initialized.value) return

    try {
      console.log('🔄 Auth Store: Initializing authentication...')
      
      // Initialize API client
      initializeApiClient()
      
      // Check for stored token
      const storedToken = localStorage.getItem('auth_token')
      
      if (!storedToken) {
        console.log('🔄 Auth Store: No stored token found')
        initialized.value = true
        return
      }

      console.log('🔄 Auth Store: Found stored token, validating...')
      
      // Set token and validate by getting current user
      setToken(storedToken)
      await getCurrentUser()
      
      console.log('✅ Auth Store: Authentication initialized successfully')
    } catch (err: any) {
      console.error('❌ Auth Store: Auth initialization failed', err)
      
      // Clear invalid token
      clearAuth()
    } finally {
      initialized.value = true
    }
  }

  const refreshAuth = async (): Promise<void> => {
    if (!isAuthenticated.value) {
      throw new Error('Not authenticated')
    }

    await getCurrentUser()
  }

  // Utility methods
  const hasRole = (role: string): boolean => {
    return userRole.value === role
  }

  const hasAnyRole = (roles: string[]): boolean => {
    return !!userRole.value && roles.includes(userRole.value)
  }

  const canAccessAdmin = (): boolean => {
    return isAdmin.value
  }

  const canAccessSuperAdmin = (): boolean => {
    return isSuperAdmin.value
  }

  return {
    // State
    user: computed(() => user.value),
    token: computed(() => token.value),
    isLoading: computed(() => isLoading.value),
    error: computed(() => error.value),
    initialized: computed(() => initialized.value),
    
    // Getters
    isAuthenticated,
    userRole,
    isAdmin,
    isSuperAdmin,
    userDisplayName,
    
    // Actions
    login,
    register,
    changePassword,
    forgotPassword,
    resetPassword,
    logout,
    getCurrentUser,
    initializeAuth,
    refreshAuth,
    clearAuth,
    setError,
    
    // API client access
    apiClient: computed(() => apiClient.value),
    
    // Utility methods
    hasRole,
    hasAnyRole,
    canAccessAdmin,
    canAccessSuperAdmin
  }
})
