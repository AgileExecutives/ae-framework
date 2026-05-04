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

// Template Contract System Types
export type TemplateChannel = 'EMAIL' | 'DOCUMENT';

export interface VariableSchema {
  type: string;
  description?: string;
  required?: boolean;
  properties?: Record<string, VariableSchema>;
  items?: VariableSchema;
  example?: any;
}

export interface TemplateContract {
  id: number;
  module: string;
  template_key: string;
  description: string;
  supported_channels: TemplateChannel[];
  variable_schema: Record<string, VariableSchema>;
  default_sample_data: Record<string, any>;
  created_at: string;
  updated_at: string;
}

export interface Template {
  id: number;
  tenant_id?: number;
  organization_id?: number;
  template_type: string; // Reference to contract's module/template_key
  name: string;
  description?: string;
  content: string;
  variables?: string[];
  sample_data?: Record<string, any>;
  is_active: boolean;
  is_default: boolean;
  version?: number;
  created_at: string;
  updated_at: string;
  storage_key?: string;
  preview_url?: string;
}

export interface CreateTemplateRequest {
  template_type: string;
  name: string;
  description?: string;
  content: string;
  variables?: string[];
  sample_data?: Record<string, any>;
  is_active?: boolean;
  is_default?: boolean;
}

export interface UpdateTemplateRequest {
  name?: string;
  description?: string;
  content?: string;
  variables?: string[];
  sample_data?: Record<string, any>;
  is_active?: boolean;
  is_default?: boolean;
}

export interface DuplicateTemplateRequest {
  name: string;
  description?: string;
}

export interface RenderTemplateRequest {
  data: Record<string, any>;
}

export interface RenderTemplateResponse {
  content: string;
  channel?: TemplateChannel;
  subject?: string; // For EMAIL channel
}

// Export Client type from generated types
export type Client = components['schemas']['entities.ClientResponse'];

// Additional request types that may not be in generated types yet
export interface CreateClientRequest {
  name: string;
  email?: string;
  phone?: string;
  address?: string;
  // Add other client creation fields as needed
}

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