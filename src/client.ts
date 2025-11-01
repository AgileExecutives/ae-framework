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
}

export default AESaasApiClient;
