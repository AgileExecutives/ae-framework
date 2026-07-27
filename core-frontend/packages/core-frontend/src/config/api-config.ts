/**
 * API Configuration
 * 
 * Centralized API client configuration using environment variables
 */

import { AESaasApiClient } from '@agile-exec/api-client'

export interface ApiConfig {
  baseURL: string
}

// Get API configuration from environment variables
export const getApiConfig = (): ApiConfig => {
  const baseURL = import.meta.env.VITE_API_BASE_URL

  if (!baseURL) {
    throw new Error('VITE_API_BASE_URL environment variable is required')
  }

  return {
    baseURL
  }
}

// Create and configure API client instance
export const createApiClient = (): AESaasApiClient => {
  const config = getApiConfig()
  
  return new AESaasApiClient({
    baseURL: config.baseURL
  })
}

// Global API client instance - initialized lazily
let globalApiClient: AESaasApiClient | null = null

export const getApiClient = (): AESaasApiClient => {
  if (!globalApiClient) {
    globalApiClient = createApiClient()
  }
  return globalApiClient
}

// Reset the global client (useful for testing or reconfiguration)
export const resetApiClient = () => {
  globalApiClient = null
}