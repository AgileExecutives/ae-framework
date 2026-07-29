import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import type { ApiError } from '@/config/api-config'
import { getApiClient as getGlobalApiClient } from '@/config/api-config'
import type * as AuthTypes from './auth-types'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref<AuthTypes.User | null>(null)
  const token = ref<string | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)
  const initialized = ref(false)

  // Get the global API client
  const getApiClient = () => {
    return getGlobalApiClient()
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
    const client = getApiClient()
    
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

      const client = getApiClient()
      const response = await client.login({
        email: credentials.username,
        password: credentials.password
      })

      console.log('🔍 Auth Store: Raw API response:', JSON.stringify(response, null, 2))

      // Handle successful response with flexible shapes
      // Example shapes:
      // { success: true, data: { token, user } }
      // { success: true, data: { access_token, user } }
      // { success: true, data: { data: { token, user } } }
      // { success: true, data: user } (no token)
      const respData = response?.data ?? response

      // Try common token fields
      const tokenCandidate = respData?.token || respData?.access_token || respData?.data?.token || respData?.data?.access_token
      const userCandidate = respData?.user || respData?.data?.user || (respData?.data && typeof respData.data === 'object' && !respData.data.token ? respData.data : undefined)

      if (tokenCandidate) {
        console.log('✅ Auth Store: Detected token in login response')
        setToken(tokenCandidate)
      }

      if (userCandidate) {
        console.log('✅ Auth Store: Detected user in login response')
        setUser(userCandidate as AuthTypes.User)
      }

      if (tokenCandidate && userCandidate) {
        console.log('✅ Auth Store: Login successful', { username: (userCandidate as any)?.username })
        // Mark auth as initialized to avoid immediate re-validation during navigation
        initialized.value = true
        return
      }

      // Fallback: some APIs return the user object directly without token (session-based)
      if (respData && typeof respData === 'object' && !tokenCandidate && respData?.username) {
        setUser(respData as AuthTypes.User)
        console.log('✅ Auth Store: Login succeeded (user-only response)')
        return
      }

      // If we reach here, the response was invalid
      console.error('❌ Auth Store: Invalid response structure')
      throw new Error('Invalid response from server')
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

      const client = getApiClient()
      const response = await client.register(credentials)
      console.log('✅ Auth Store: Registration response received', response)
      
      // Handle the response structure from the API client
      if (response.success && response.data?.user) {
        // If token is provided, set it (auto-login after registration)
        if (response.data.token) {
          setToken(response.data.token)
          setUser(response.data.user as AuthTypes.User)
          // Mark auth initialized after auto-login
          initialized.value = true
          console.log('✅ Auth Store: Registration successful with auto-login', { user: response.data.user.username })
        } else {
          // Registration successful but no auto-login
          console.log('✅ Auth Store: Registration successful, please login', { user: response.data.user.username })
        }
        return response.data.user as AuthTypes.User
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

      const client = getApiClient()
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

      const client = getApiClient()
      await client.forgotPassword({ email })

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

      console.log('🔑 Auth Store: resetPassword called with token:', token ? `${token.substring(0, 10)}...` : 'MISSING', 'password:', newPassword ? '***' : 'MISSING')
      
      const client = getApiClient()
      await client.resetPassword(token, { new_password: newPassword })

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

      const client = getApiClient()
      
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

      const client = getApiClient()
      const response = await client.getCurrentUser()

      // Response may have several shapes:
      // - { success: true, data: { ...user } }
      // - { user: { ... } }
      // - { ...user }
      let userData: AuthTypes.User | null = null

      if (!response) {
        userData = null
      } else if (response.data && typeof response.data === 'object') {
        userData = response.data as AuthTypes.User
      } else if (response.user && typeof response.user === 'object') {
        userData = response.user as AuthTypes.User
      } else if (typeof response === 'object' && response.username) {
        userData = response as AuthTypes.User
      }

      if (!userData) {
        // No usable user data returned
        console.warn('⚠️ Auth Store: getCurrentUser returned no user data', { rawResponse: response })
        throw new Error('No user data returned')
      }

      setUser(userData)
      console.log('✅ Auth Store: User data refreshed', { user: userData.username })

      return userData
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
      getApiClient()
      
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
    apiClient: computed(() => getApiClient()),
    
    // Utility methods
    hasRole,
    hasAnyRole,
    canAccessAdmin,
    canAccessSuperAdmin
  }
})
