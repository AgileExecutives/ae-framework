// Re-export everything from the client
export * from './client';
export * from './types';
export * from './composables/index';

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