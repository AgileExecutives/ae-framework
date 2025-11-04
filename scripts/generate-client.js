#!/usr/bin/env node

import { execSync } from 'child_process';
import { writeFileSync, readFileSync, mkdirSync, existsSync } from 'fs';
import { join, dirname } from 'path';
import { fileURLToPath } from 'url';

const __dirname = dirname(fileURLToPath(import.meta.url));
const projectRoot = join(__dirname, '..');

// Configuration
const config = {
  swaggerUrl: process.env.SWAGGER_URL || 'http://localhost:8080/swagger/doc.json',
  outputDir: join(projectRoot, 'src'),
  typesFile: 'types.ts',
  clientFile: 'client.ts',
  composablesFile: 'composables/index.ts'
};

console.log('🚀 Generating Vue API Client from Swagger...\n');

async function fetchSwaggerSpec() {
  console.log(`📡 Fetching Swagger spec from: ${config.swaggerUrl}`);
  
  try {
    const response = await fetch(config.swaggerUrl);
    if (!response.ok) {
      throw new Error(`Failed to fetch: ${response.status} ${response.statusText}`);
    }
    
    const spec = await response.json();
    console.log('✅ Swagger spec fetched successfully');
    return spec;
  } catch (error) {
    console.error('❌ Failed to fetch Swagger spec:', error.message);
    process.exit(1);
  }
}

function convertSwaggerToOpenAPI(swaggerSpec) {
  console.log('🔄 Converting Swagger 2.0 to OpenAPI 3.0...');
  
  try {
    // Import swagger2openapi dynamically for ES modules
    return import('swagger2openapi').then(({ convert }) => {
      return new Promise((resolve, reject) => {
        convert(swaggerSpec, { patch: true, warnOnly: true }, (err, options) => {
          if (err) {
            reject(err);
          } else {
            console.log('✅ Successfully converted to OpenAPI 3.0');
            resolve(options.openapi);
          }
        });
      });
    });
  } catch (error) {
    console.error('❌ Failed to convert Swagger spec:', error.message);
    process.exit(1);
  }
}

async function generateTypes(openApiSpec) {
  console.log('⚙️  Generating TypeScript types...');
  
  // Save the OpenAPI spec temporarily
  const tempSpecPath = join(projectRoot, 'temp-openapi.json');
  writeFileSync(tempSpecPath, JSON.stringify(openApiSpec, null, 2));
  
  try {
    // Generate types using openapi-typescript
    const typesOutput = execSync(`npx openapi-typescript ${tempSpecPath}`, { 
      encoding: 'utf8',
      cwd: projectRoot 
    });
    
    // Clean up temp file
    execSync(`rm ${tempSpecPath}`);
    
    console.log('✅ TypeScript types generated');
    return typesOutput;
  } catch (error) {
    console.error('❌ Failed to generate types:', error.message);
    process.exit(1);
  }
}

function generateApiClient(openApiSpec) {
  console.log('🏗️  Generating API client...');
  
  const { info, paths, components } = openApiSpec;
  const baseUrl = openApiSpec.servers?.[0]?.url || '/api/v1';
  
  // Extract endpoints and generate client methods
  const endpoints = Object.entries(paths).map(([path, methods]) => {
    return Object.entries(methods).map(([method, spec]) => ({
      path,
      method: method.toUpperCase(),
      operationId: spec.operationId,
      summary: spec.summary,
      parameters: spec.parameters || [],
      requestBody: spec.requestBody,
      responses: spec.responses
    }));
  }).flat();
  
  // Generate methods dynamically from endpoints
  const generateMethods = () => {
    const methods = [];
    
    endpoints.forEach((endpoint, index) => {
      if (!endpoint.operationId) {
        console.log(`❌ Skipping endpoint ${index + 1}: ${endpoint.method} ${endpoint.path} - missing operationId`);
        return;
      }
      
      console.log(`✅ Processing endpoint ${index + 1}: ${endpoint.method} ${endpoint.path} - operationId: ${endpoint.operationId}`);
      
      const methodName = endpoint.operationId;
      // Strip /api/v1 prefix from path since it's already in baseURL
      const cleanPath = endpoint.path.replace(/^\/api\/v1/, '') || '/';
      const pathTemplate = cleanPath.replace(/\{([^}]+)\}/g, '${$1}');
      const hasPathParams = endpoint.path.includes('{');
      const hasQueryParams = endpoint.parameters.some(p => p.in === 'query');
      const hasBody = endpoint.requestBody;
      
      let methodCode = `  async ${methodName}(`;
      
      // Add parameters
      const params = [];
      if (hasPathParams) {
        const pathParams = endpoint.parameters.filter(p => p.in === 'path');
        pathParams.forEach(param => {
          params.push(`${param.name}: ${param.schema?.type === 'integer' ? 'number' : 'string'}`);
        });
      }
      
      if (hasBody) {
        params.push('data: any');
      }
      
      if (hasQueryParams) {
        params.push('params?: Record<string, any>');
      }
      
      methodCode += params.join(', ');
      methodCode += ') {\n';
      
      // Determine response type
      const isListEndpoint = methodName.includes('get') && (methodName.includes('s') || methodName.includes('search'));
      const responseType = isListEndpoint ? 'any' : 'ApiResponse<any>';
      
      // Generate method body
      if (hasPathParams) {
        const pathParams = endpoint.parameters.filter(p => p.in === 'path');
        pathParams.forEach(param => {
          methodCode += `    if (!${param.name}) throw new Error('${param.name} is required');\n`;
        });
      }
      
      const requestArgs = [`'${endpoint.method}'`, `\`${pathTemplate}\``];
      
      if (hasBody) {
        requestArgs.push('data');
      } else {
        requestArgs.push('undefined');
      }
      
      if (hasQueryParams) {
        requestArgs.push('params');
      }
      
      if (endpoint.method === 'DELETE') {
        methodCode += `    await this.request<ApiResponse<any>>(${requestArgs.join(', ')});\n`;
        methodCode += `    return { success: true };\n`;
      } else {
        methodCode += `    const response = await this.request<${responseType}>(${requestArgs.join(', ')});\n`;
        if (isListEndpoint) {
          methodCode += `    return response;\n`;
        } else {
          methodCode += `    return response;\n`;
        }
      }
      
      methodCode += '  }\n';
      methods.push(methodCode);
    });
    
    return methods.join('\n');
  };

  const dynamicMethods = generateMethods();
  
  return `import axios, { AxiosInstance, AxiosResponse, AxiosRequestConfig } from 'axios';
import type { paths, components } from './types';

// API Response types
export type ApiResponse<T = any> = {
  success: boolean;
  message?: string;
  data?: T;
  error?: string;
};

export type ListResponse<T = any> = {
  success: boolean;
  data: T[];
  pagination?: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
  };
};

// Configuration interface
export interface ApiClientConfig {
  baseURL: string;
  timeout?: number;
  headers?: Record<string, string>;
}

// Error class
export class ApiError extends Error {
  constructor(
    public status: number,
    public data: any,
    message: string
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

// Main API Client class
export class AESaasApiClient {
  private client: AxiosInstance;
  private token: string | null = null;

  constructor(config: ApiClientConfig) {
    this.client = axios.create({
      baseURL: config.baseURL,
      timeout: config.timeout || 10000,
      headers: {
        'Content-Type': 'application/json',
        ...config.headers,
      },
    });

    // Request interceptor to add authorization
    this.client.interceptors.request.use((config) => {
      if (this.token) {
        config.headers.Authorization = \`Bearer \${this.token}\`;
      }
      return config;
    });

    // Response interceptor for error handling
    this.client.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response) {
          throw new ApiError(
            error.response.status,
            error.response.data,
            error.message
          );
        }
        throw error;
      }
    );
  }

  setToken(token: string) {
    this.token = token;
  }

  clearToken() {
    this.token = null;
  }

  // Helper method for making requests
  private async request<T>(
    method: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH',
    path: string,
    data?: any,
    params?: Record<string, any>
  ): Promise<T> {
    const config: AxiosRequestConfig = {
      method,
      url: path,
      data,
      params,
    };

    const response: AxiosResponse<T> = await this.client.request(config);
    return response.data;
  }

  // Dynamically generated methods from OpenAPI spec
${dynamicMethods}
}

export default AESaasApiClient;
`;
}

function generateComposables() {
  console.log('🔧 Generating Vue composables...');
  
  return `import { ref, reactive, computed, type Ref } from 'vue';
import { AESaasApiClient, type ApiClientConfig, type ApiResponse, type ListResponse } from '../client';

// Global client instance
let globalClient: AESaasApiClient | null = null;

// Plugin installation
export function createApiClient(config: ApiClientConfig) {
  globalClient = new AESaasApiClient(config);
  return globalClient;
}

export function useApiClient() {
  if (!globalClient) {
    throw new Error('API client not initialized. Call createApiClient() first.');
  }
  return globalClient;
}

// Authentication composable
export function useAuth() {
  const client = useApiClient();
  const isAuthenticated = ref(false);
  const user = ref(null);
  const loading = ref(false);
  const error = ref<string | null>(null);

  const login = async (credentials: { username: string; password: string }) => {
    try {
      loading.value = true;
      error.value = null;
      
      const response = await client.login(credentials);
      user.value = response.data?.user;
      isAuthenticated.value = true;
      
      return response;
    } catch (err: any) {
      error.value = err.message || 'Login failed';
      isAuthenticated.value = false;
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const logout = async () => {
    try {
      loading.value = true;
      await client.logout();
    } catch (err: any) {
      error.value = err.message || 'Logout failed';
    } finally {
      user.value = null;
      isAuthenticated.value = false;
      loading.value = false;
    }
  };

  const getCurrentUser = async () => {
    try {
      loading.value = true;
      error.value = null;
      
      const userData = await client.getCurrentUser();
      user.value = userData;
      isAuthenticated.value = true;
      
      return userData;
    } catch (err: any) {
      error.value = err.message || 'Failed to get user';
      isAuthenticated.value = false;
      throw err;
    } finally {
      loading.value = false;
    }
  };

  return {
    isAuthenticated: readonly(isAuthenticated),
    user: readonly(user),
    loading: readonly(loading),
    error: readonly(error),
    login,
    logout,
    getCurrentUser
  };
}

// Generic resource composable
export function useResource<T = any>(resourceName: string) {
  const client = useApiClient();
  const items = ref<T[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);
  const pagination = reactive({
    page: 1,
    limit: 10,
    total: 0,
    totalPages: 0
  });

  const fetchItems = async (params?: Record<string, any>) => {
    try {
      loading.value = true;
      error.value = null;
      
      const method = \`get\${resourceName.charAt(0).toUpperCase() + resourceName.slice(1)}\` as keyof AESaasApiClient;
      const response = await (client[method] as Function)({ 
        page: pagination.page, 
        limit: pagination.limit, 
        ...params 
      });
      
      items.value = response.data || [];
      if (response.pagination) {
        Object.assign(pagination, response.pagination);
      }
      
      return response;
    } catch (err: any) {
      error.value = err.message || \`Failed to fetch \${resourceName}\`;
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const createItem = async (data: Partial<T>) => {
    try {
      loading.value = true;
      error.value = null;
      
      const method = \`create\${resourceName.charAt(0).toUpperCase() + resourceName.slice(1, -1)}\` as keyof AESaasApiClient;
      const item = await (client[method] as Function)(data);
      
      items.value.unshift(item);
      return item;
    } catch (err: any) {
      error.value = err.message || \`Failed to create \${resourceName.slice(0, -1)}\`;
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const updateItem = async (id: number, data: Partial<T>) => {
    try {
      loading.value = true;
      error.value = null;
      
      const method = \`update\${resourceName.charAt(0).toUpperCase() + resourceName.slice(1, -1)}\` as keyof AESaasApiClient;
      const updatedItem = await (client[method] as Function)(id, data);
      
      const index = items.value.findIndex((item: any) => item.id === id);
      if (index !== -1) {
        items.value[index] = updatedItem;
      }
      
      return updatedItem;
    } catch (err: any) {
      error.value = err.message || \`Failed to update \${resourceName.slice(0, -1)}\`;
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const deleteItem = async (id: number) => {
    try {
      loading.value = true;
      error.value = null;
      
      const method = \`delete\${resourceName.charAt(0).toUpperCase() + resourceName.slice(1, -1)}\` as keyof AESaasApiClient;
      await (client[method] as Function)(id);
      
      items.value = items.value.filter((item: any) => item.id !== id);
    } catch (err: any) {
      error.value = err.message || \`Failed to delete \${resourceName.slice(0, -1)}\`;
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const nextPage = () => {
    if (pagination.page < pagination.totalPages) {
      pagination.page++;
      fetchItems();
    }
  };

  const prevPage = () => {
    if (pagination.page > 1) {
      pagination.page--;
      fetchItems();
    }
  };

  const setPage = (page: number) => {
    pagination.page = page;
    fetchItems();
  };

  return {
    items: readonly(items),
    loading: readonly(loading),
    error: readonly(error),
    pagination: readonly(pagination),
    fetchItems,
    createItem,
    updateItem,
    deleteItem,
    nextPage,
    prevPage,
    setPage
  };
}

// Specific resource composables
export const usePlans = () => useResource('plans');
export const useCustomers = () => useResource('customers');
export const useContacts = () => useResource('contacts');
export const useEmails = () => useResource('emails');

// Search composable
// Note: Commented out until search/quickSearch/health/ping endpoints have operationIds
/*
export function useSearch() {
  // Implementation disabled - requires search/quickSearch methods
}

export function useHealth() {
  // Implementation disabled - requires health/ping methods  
}
*/

// Utility composables
function readonly<T>(ref: Ref<T>) {
  return computed(() => ref.value);
}
`;
}

async function writeGeneratedFiles() {
  console.log('📝 Writing generated files...');
  
  // Ensure output directory exists
  if (!existsSync(config.outputDir)) {
    mkdirSync(config.outputDir, { recursive: true });
  }
  
  const composablesDir = join(config.outputDir, 'composables');
  if (!existsSync(composablesDir)) {
    mkdirSync(composablesDir, { recursive: true });
  }
  
  try {
    // Fetch and process the Swagger spec
    const swaggerSpec = await fetchSwaggerSpec();
    const openApiSpec = await convertSwaggerToOpenAPI(swaggerSpec);
    
    // Generate types
    const types = await generateTypes(openApiSpec);
    writeFileSync(join(config.outputDir, config.typesFile), types);
    console.log(`✅ Types written to ${config.typesFile}`);
    
    // Generate API client
    const clientCode = generateApiClient(openApiSpec);
    writeFileSync(join(config.outputDir, config.clientFile), clientCode);
    console.log(`✅ Client written to ${config.clientFile}`);
    
    // Generate composables
    const composablesCode = generateComposables();
    writeFileSync(join(config.outputDir, config.composablesFile), composablesCode);
    console.log(`✅ Composables written to ${config.composablesFile}`);
    
  } catch (error) {
    console.error('❌ Failed to write files:', error.message);
    process.exit(1);
  }
}

async function main() {
  try {
    await writeGeneratedFiles();
    console.log('\n🎉 Vue API Client generated successfully!\n');
    
    console.log('📦 Generated files:');
    console.log(`  - ${config.typesFile} (TypeScript types)`);
    console.log(`  - ${config.clientFile} (API client class)`);
    console.log(`  - ${config.composablesFile} (Vue composables)`);
    
    console.log('\n🚀 Next steps:');
    console.log('  1. Run "npm install" to install dependencies');
    console.log('  2. Run "npm run build" to build the package');
    console.log('  3. Import the client in your Vue app');
    
  } catch (error) {
    console.error('❌ Generation failed:', error.message);
    process.exit(1);
  }
}

// Run if called directly
if (import.meta.url === `file://${process.argv[1]}`) {
  main();
}