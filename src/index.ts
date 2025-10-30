// Re-export everything from the client
export * from './client';
export * from './types';
export * from './composables/index';

// Import for type definitions
import type { components } from './types';

// Additional types that may not be in the generated types yet
export interface CostProvider {
  id: number;
  tenant_id?: number;
  organization: string;
  department: string;
  street_address?: string;
  zip?: string;
  city?: string;
  created_at?: string;
  updated_at?: string;
}

// Export Client type from generated types
export type Client = components['schemas']['models.ClientResponse'];

// Vue plugin
import type { App } from 'vue';
import { createApiClient } from './composables/index';
import type { ApiClientConfig } from './client';

export interface AESaasApiClientOptions extends ApiClientConfig {}

export const AESaasApiClientPlugin = {
  install(app: App, options: AESaasApiClientOptions) {
    const client = createApiClient(options);
    
    // Make client available globally
    app.config.globalProperties['$apiClient'] = client;
    app.provide('apiClient', client);
  }
};

// Default export for plugin usage
export default AESaasApiClientPlugin;