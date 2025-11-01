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

  // Calendar methods
  async getCalendar() {
    try {
      console.log('📅 Fetching calendar list from /calendar endpoint')
      const response = await this.request<any>('GET', '/calendar');
      
      return {
        success: true,
        data: response.data || response || []
      };
    } catch (error) {
      console.error('📅 Error fetching calendars:', error);
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Failed to fetch calendars',
        data: []
      };
    }
  }

  async getCalendarEvents(params?: { limit?: number; startDate?: string; date?: string }) {
    try {
      console.log('📅 Fetching calendars from /calendar endpoint')
      
      // First, get the list of calendars
      const calendarsResponse = await this.request<any>('GET', '/calendar');
      
      if (!calendarsResponse || !calendarsResponse.data || !Array.isArray(calendarsResponse.data) || calendarsResponse.data.length === 0) {
        console.warn('📅 No calendars found or invalid response');
        return {
          success: false,
          error: 'No calendars available',
          data: []
        };
      }
      
      // Use the first calendar in the list
      const firstCalendar = calendarsResponse.data[0];
      console.log('📅 Using first calendar:', firstCalendar);
      
      // If the calendar has events directly, return them
      if (firstCalendar.events && Array.isArray(firstCalendar.events)) {
        console.log('📅 Found events in calendar:', firstCalendar.events.length);
        return {
          success: true,
          data: firstCalendar.events
        };
      }
      
      // If calendar has an ID, try to fetch events for that calendar
      if (firstCalendar.id) {
        console.log('📅 Fetching events for calendar ID:', firstCalendar.id);
        try {
          const eventsResponse = await this.request<any>('GET', `/calendar/${firstCalendar.id}/events`, undefined, params);
          return {
            success: true,
            data: eventsResponse.data || eventsResponse || []
          };
        } catch (eventsError) {
          console.warn('📅 Failed to fetch events for calendar, using calendar data as events');
          // Fallback: treat the calendar list as events if events endpoint fails
          return {
            success: true,
            data: calendarsResponse.data
          };
        }
      }
      
      // Fallback: return the calendars data as events
      return {
        success: true,
        data: calendarsResponse.data
      };
      
    } catch (error) {
      console.error('📅 Error fetching calendar events:', error);
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Failed to fetch calendar events',
        data: []
      };
    }
  }

  async getCalendarSeries(params?: { limit?: number }) {
    try {
      const response = await this.request<any>('GET', '/calendar-series', undefined, params);
      return {
        success: true,
        data: response.data || response || []
      };
    } catch (error) {
      console.error('Error fetching calendar series:', error);
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Failed to fetch calendar series',
        data: []
      };
    }
  }

  async createCalendarEvent(data: any) {
    try {
      const response = await this.request<any>('POST', '/calendar-entries', data);
      return {
        success: true,
        data: response.data || response
      };
    } catch (error) {
      console.error('Error creating calendar event:', error);
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Failed to create calendar event'
      };
    }
  }

  // Client management methods
  async getClients(params?: { page?: number; limit?: number }) {
    try {
      const response = await this.request<any>('GET', '/clients', undefined, params);
      console.log('🔍 API Client getClients - Raw response:', response);
      console.log('🔍 API Client getClients - Response type:', typeof response);
      console.log('🔍 API Client getClients - Response.data:', response.data);
      console.log('🔍 API Client getClients - Response.data type:', typeof response.data);
      
      const rawData = response.data || response;
      
      // Extract clients array from the response structure
      const clientsArray = rawData.clients || rawData;
      console.log('🔍 API Client getClients - Extracted clients array:', clientsArray);
      console.log('🔍 API Client getClients - Clients array is array:', Array.isArray(clientsArray));
      console.log('🔍 API Client getClients - Clients array length:', clientsArray?.length);
      
      return {
        success: true,
        data: clientsArray,
        pagination: rawData.pagination // Include pagination info if available
      };
    } catch (error) {
      console.error('Error fetching clients:', error);
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Failed to fetch clients',
        data: []
      };
    }
  }

  async getClient(id: number | string) {
    try {
      const response = await this.request<any>('GET', `/clients/${id}`);
      return {
        success: true,
        data: response.data || response
      };
    } catch (error) {
      console.error('Error fetching client:', error);
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Failed to fetch client'
      };
    }
  }

  async searchClients(params: { q: string }) {
    try {
      const response = await this.request<any>('GET', '/clients/search', undefined, params);
      return {
        success: true,
        data: response.data || response
      };
    } catch (error) {
      console.error('Error searching clients:', error);
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Failed to search clients',
        data: []
      };
    }
  }

  async createClient(data: any) {
    try {
      const response = await this.request<any>('POST', '/clients', data);
      return {
        success: true,
        data: response.data || response
      };
    } catch (error) {
      console.error('Error creating client:', error);
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Failed to create client'
      };
    }
  }

  async updateClient(id: number | string, data: any) {
    try {
      const response = await this.request<any>('PUT', `/clients/${id}`, data);
      return {
        success: true,
        data: response.data || response
      };
    } catch (error) {
      console.error('Error updating client:', error);
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Failed to update client'
      };
    }
  }

  async deleteClient(id: number | string) {
    try {
      const response = await this.request<any>('DELETE', `/clients/${id}`);
      return {
        success: true,
        data: response.data || response
      };
    } catch (error) {
      console.error('Error deleting client:', error);
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Failed to delete client'
      };
    }
  }

  async getCostProviders() {
    try {
      const response = await this.request<any>('GET', '/cost-providers');
      console.log('🔍 API Client getCostProviders - Raw response:', response);
      
      const rawData = response.data || response;
      
      // Extract cost_providers array from the response structure
      const costProvidersArray = rawData.cost_providers || rawData;
      console.log('🔍 API Client getCostProviders - Extracted array:', costProvidersArray);
      console.log('🔍 API Client getCostProviders - Array length:', costProvidersArray?.length);
      
      return {
        success: true,
        data: costProvidersArray,
        pagination: rawData.pagination // Include pagination info if available
      };
    } catch (error) {
      console.error('Error fetching cost providers:', error);
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Failed to fetch cost providers',
        data: []
      };
    }
  }

  async getStatic(type: string) {
    try {
      const response = await this.request<any>('GET', `/static/${type}`);
      return {
        success: true,
        data: response.data || response
      };
    } catch (error) {
      console.error('Error fetching static data:', error);
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Failed to fetch static data',
        data: []
      };
    }
  }
}

export default AESaasApiClient;
