import axios from 'axios';
import type { AxiosInstance, AxiosResponse, AxiosRequestConfig } from 'axios';
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
        config.headers.Authorization = `Bearer ${this.token}`;
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

    // Debug logging for authentication
    if (this.token) {
      console.log(`🔐 Making authenticated request: ${method} ${path} with Bearer token`);
    }

    const response: AxiosResponse<T> = await this.client.request(config);
    return response.data;
  }

  // Authentication methods
  async login(credentials: { username: string; password: string }) {
    const response = await this.request<ApiResponse<{ token: string; user: any }>>('POST', '/auth/login', credentials);
    if (response.data?.token) {
      this.setToken(response.data.token);
    }
    return response.data!;
  }

  async logout() {
    await this.request<ApiResponse>('POST', '/auth/logout');
    this.clearToken();
  }

  async getCurrentUser() {
    const response = await this.request<ApiResponse<any>>('GET', '/auth/me');
    return response.data!;
  }

  async register(credentials: { username: string; email: string; password: string }) {
    const response = await this.request<ApiResponse<{ token: string; user: any }>>('POST', '/auth/register', credentials);
    if (response.data?.token) {
      this.setToken(response.data.token);
    }
    return response.data!;
  }

  async changePassword(credentials: { current_password: string; new_password: string }) {
    const response = await this.request<ApiResponse>('POST', '/auth/change-password', credentials);
    return response.data!;
  }

  async forgotPassword(email: string) {
    const response = await this.request<ApiResponse>('POST', '/auth/forgot-password', { email });
    return response.data!;
  }

  async resetPassword(token: string, newPassword: string) {
    const response = await this.request<ApiResponse>('POST', '/auth/reset-password', { token, password: newPassword });
    return response.data!;
  }

  async getPasswordSecurity() {
    const response = await this.request<ApiResponse<any>>('GET', '/auth/password-security');
    return response.data!;
  }

  // Health check methods
  async health() {
    return this.request<{ status: string }>('GET', '/health');
  }

  async ping() {
    return this.request<any>('GET', '/ping');
  }

  // Resource methods
  async getPlans(params?: { page?: number; limit?: number }) {
    const response = await this.request<ListResponse<any>>('GET', '/plans', undefined, params);
    return response;
  }

  async getPlan(id: number) {
    const response = await this.request<ApiResponse<any>>('GET', `/plans/${id}`);
    return response.data!;
  }

  async getCustomers(params?: { page?: number; limit?: number }) {
    const response = await this.request<ListResponse<any>>('GET', '/customers', undefined, params);
    return response;
  }

  async getCustomer(id: number) {
    const response = await this.request<ApiResponse<any>>('GET', `/customers/${id}`);
    return response.data!;
  }

  async createCustomer(data: any) {
    const response = await this.request<ApiResponse<any>>('POST', '/customers', data);
    return response.data!;
  }

  async updateCustomer(id: number, data: any) {
    const response = await this.request<ApiResponse<any>>('PUT', `/customers/${id}`, data);
    return response.data!;
  }

  async deleteCustomer(id: number) {
    await this.request<ApiResponse>('DELETE', `/customers/${id}`);
  }

  async getContacts(params?: { page?: number; limit?: number }) {
    const response = await this.request<ListResponse<any>>('GET', '/contacts', undefined, params);
    return response;
  }

  async getContact(id: number) {
    const response = await this.request<ApiResponse<any>>('GET', `/contacts/${id}`);
    return response.data!;
  }

  async createContact(data: any) {
    const response = await this.request<ApiResponse<any>>('POST', '/contacts', data);
    return response.data!;
  }

  async updateContact(id: number, data: any) {
    const response = await this.request<ApiResponse<any>>('PUT', `/contacts/${id}`, data);
    return response.data!;
  }

  async deleteContact(id: number) {
    await this.request<ApiResponse>('DELETE', `/contacts/${id}`);
  }

  async getEmails(params?: { page?: number; limit?: number; status?: string }) {
    const response = await this.request<ListResponse<any>>('GET', '/emails', undefined, params);
    return response;
  }

  async sendEmail(data: any) {
    const response = await this.request<ApiResponse<any>>('POST', '/emails/send', data);
    return response.data!;
  }

  async quickSearch(params: { q: string; limit?: number }) {
    return this.request<any>('GET', '/search/quick', undefined, params);
  }

  async search(data: any) {
    return this.request<any>('POST', '/search', data);
  }

  async getUserSettings() {
    const response = await this.request<ApiResponse<any>>('GET', '/user-settings');
    return response.data!;
  }

  async updateUserSettings(data: any) {
    const response = await this.request<ApiResponse<any>>('PUT', '/user-settings', data);
    return response.data!;
  }

  // Client methods
  async getClients(params?: { page?: number; limit?: number }) {
    try {
      const response = await this.request<any>('GET', '/clients', undefined, params);
      
      // Handle the actual API response format: {"clients": [...], "total": 46, "page": 1, "limit": 10}
      if (response && response.clients) {
        return {
          success: true,
          data: response.clients,
          pagination: {
            total: response.total,
            page: response.page,
            limit: response.limit
          }
        };
      }
      
      // Fallback for other formats
      return {
        success: true,
        data: Array.isArray(response) ? response : []
      };
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Failed to fetch clients',
        data: []
      };
    }
  }

  async getClient(id: number) {
    try {
      const response = await this.request<any>('GET', `/clients/${id}`);
      return {
        success: true,
        data: response
      };
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Failed to fetch client'
      };
    }
  }

  async createClient(data: any) {
    try {
      const response = await this.request<any>('POST', '/clients', data);
      return {
        success: true,
        data: response
      };
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Failed to create client'
      };
    }
  }

  async updateClient(id: number, data: any) {
    try {
      const response = await this.request<any>('PUT', `/clients/${id}`, data);
      return {
        success: true,
        data: response
      };
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Failed to update client'
      };
    }
  }

  async deleteClient(id: number) {
    try {
      await this.request<any>('DELETE', `/clients/${id}`);
      return {
        success: true
      };
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Failed to delete client'
      };
    }
  }

  // Cost Provider methods
  async getCostProviders(params?: { page?: number; limit?: number }) {
    try {
      const response = await this.request<any>('GET', '/cost-providers', undefined, params);
      
      // Handle the actual API response format
      if (response && response.cost_providers) {
        return {
          success: true,
          data: response.cost_providers,
          pagination: {
            total: response.total,
            page: response.page,
            limit: response.limit
          }
        };
      }
      
      // Handle array response
      if (Array.isArray(response)) {
        return {
          success: true,
          data: response
        };
      }
      
      // Fallback
      return {
        success: true,
        data: []
      };
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Failed to fetch cost providers',
        data: []
      };
    }
  }

  async getCostProvider(id: number) {
    try {
      const response = await this.request<any>('GET', `/cost-providers/${id}`);
      return {
        success: true,
        data: response
      };
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Failed to fetch cost provider'
      };
    }
  }

  async searchClients(params: { q: string; page?: number; limit?: number }) {
    try {
      const response = await this.request<any>('GET', '/clients/search', undefined, params);
      
      // Handle the actual API response format
      if (response && response.clients) {
        return {
          success: true,
          data: response.clients,
          pagination: {
            total: response.total,
            page: response.page,
            limit: response.limit
          }
        };
      }
      
      // Fallback for other formats
      return {
        success: true,
        data: Array.isArray(response) ? response : []
      };
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Failed to search clients',
        data: []
      };
    }
  }
}

export default AESaasApiClient;
