import axios from 'axios';
import type { AxiosInstance, AxiosResponse, AxiosRequestConfig } from 'axios';
import type { paths, components } from './types';

// API Response types from backend
export type ApiResponse<T = any> = {
  success: boolean;
  message?: string;
  data?: T;
  error?: string;
};

export type ListResponse<T = any> = {
  success: boolean;
  message?: string;
  data: T[];
  error?: string;
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

  // Dynamically generated methods from OpenAPI spec
  async changePassword(data: any) {
    const response = await this.request<ApiResponse<any>>('POST', `/auth/change-password`, data);
    return response;
  }

  async forgotPassword(data: any) {
    const response = await this.request<ApiResponse<any>>('POST', `/auth/forgot-password`, data);
    return response;
  }

  async login(data: any) {
    const response = await this.request<ApiResponse<any>>('POST', `/auth/login`, data);
    return response;
  }

  async logout() {
    const response = await this.request<ApiResponse<any>>('POST', `/auth/logout`, undefined);
    return response;
  }

  async getCurrentUser() {
    const response = await this.request<any>('GET', `/auth/me`, undefined);
    return response;
  }

  async resetPassword(token: string, data: any) {
    if (!token) throw new Error('token is required');
    const response = await this.request<ApiResponse<any>>('POST', `/auth/new-password/${token}`, data);
    return response;
  }

  async refreshToken() {
    const response = await this.request<ApiResponse<any>>('POST', `/auth/refresh`, undefined);
    return response;
  }

  async register(data: any) {
    const response = await this.request<ApiResponse<any>>('POST', `/auth/register`, data);
    return response;
  }

  async verifyEmail(token: string) {
    if (!token) throw new Error('token is required');
    const response = await this.request<ApiResponse<any>>('GET', `/auth/verify-email/${token}`, undefined);
    return response;
  }

  async getContacts(params?: Record<string, any>) {
    const response = await this.request<any>('GET', `/contacts`, undefined, params);
    return response;
  }

  async createContact(data: any) {
    const response = await this.request<ApiResponse<any>>('POST', `/contacts`, data);
    return response;
  }

  async submitContactForm(data: any) {
    const response = await this.request<ApiResponse<any>>('POST', `/contacts/form`, data);
    return response;
  }

  async getNewsletterSubscriptions() {
    const response = await this.request<any>('GET', `/contacts/newsletter`, undefined);
    return response;
  }

  async unsubscribeFromNewsletter(email: string) {
    if (!email) throw new Error('email is required');
    const response = await this.request<ApiResponse<any>>('DELETE', `/contacts/newsletter/${email}`, undefined);
    return response || { success: true };
  }

  async getContactById(id: string) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('GET', `/contacts/${id}`, undefined);
    return response;
  }

  async updateContact(id: string, data: any) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('PUT', `/contacts/${id}`, data);
    return response;
  }

  async deleteContact(id: string) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('DELETE', `/contacts/${id}`, undefined);
    return response || { success: true };
  }

  async getCustomers(params?: Record<string, any>) {
    const response = await this.request<any>('GET', `/customers`, undefined, params);
    return response;
  }

  async createCustomer(data: any) {
    const response = await this.request<ApiResponse<any>>('POST', `/customers`, data);
    return response;
  }

  async getCustomerById(id: number) {
    if (!id) throw new Error('id is required');
    const response = await this.request<any>('GET', `/customers/${id}`, undefined);
    return response;
  }

  async updateCustomer(id: number, data: any) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('PUT', `/customers/${id}`, data);
    return response;
  }

  async deleteCustomer(id: number) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('DELETE', `/customers/${id}`, undefined);
    return response || { success: true };
  }

  async getEmails(params?: Record<string, any>) {
    const response = await this.request<any>('GET', `/emails`, undefined, params);
    return response;
  }

  async sendEmail(data: any) {
    const response = await this.request<ApiResponse<any>>('POST', `/emails/send`, data);
    return response;
  }

  async getEmailStats() {
    const response = await this.request<any>('GET', `/emails/stats`, undefined);
    return response;
  }

  async getEmailById(id: number) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('GET', `/emails/${id}`, undefined);
    return response;
  }

  async healthCheck() {
    const response = await this.request<ApiResponse<any>>('GET', `/health`, undefined);
    return response;
  }

  async createPdf(data: any) {
    const response = await this.request<ApiResponse<any>>('POST', `/pdf/create`, data);
    return response;
  }

  async ping() {
    const response = await this.request<ApiResponse<any>>('GET', `/ping`, undefined);
    return response;
  }

  async getPlans() {
    const response = await this.request<any>('GET', `/plans`, undefined);
    return response;
  }

  async createPlan(data: any) {
    const response = await this.request<ApiResponse<any>>('POST', `/plans`, data);
    return response;
  }

  async getPlanById(id: number) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('GET', `/plans/${id}`, undefined);
    return response;
  }

  async updatePlan(id: number, data: any) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('PUT', `/plans/${id}`, data);
    return response;
  }

  async deletePlan(id: number) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('DELETE', `/plans/${id}`, undefined);
    return response || { success: true };
  }

  async listStaticFiles() {
    const response = await this.request<ApiResponse<any>>('GET', `/static`, undefined);
    return response;
  }

  async getStaticFile(filename: string) {
    if (!filename) throw new Error('filename is required');
    const response = await this.request<ApiResponse<any>>('GET', `/static/${filename}`, undefined);
    return response;
  }

  async getUserSettings() {
    const response = await this.request<any>('GET', `/user-settings`, undefined);
    return response;
  }

  async updateUserSettings(data: any) {
    const response = await this.request<ApiResponse<any>>('PUT', `/user-settings`, data);
    return response;
  }

  async resetUserSettings() {
    const response = await this.request<ApiResponse<any>>('POST', `/user-settings/reset`, undefined);
    return response;
  }

  async getCalendarEntries(params?: Record<string, any>) {
    const response = await this.request<any>('GET', `/calendar-entries`, undefined, params);
    return response;
  }

  async createCalendarEntry(data: any) {
    const response = await this.request<ApiResponse<any>>('POST', `/calendar-entries`, data);
    return response;
  }

  async getCalendarEntryById(id: number) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('GET', `/calendar-entries/${id}`, undefined);
    return response;
  }

  async updateCalendarEntry(id: number, data: any) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('PUT', `/calendar-entries/${id}`, data);
    return response;
  }

  async deleteCalendarEntry(id: number) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('DELETE', `/calendar-entries/${id}`, undefined);
    return response || { success: true };
  }

  async getCalendarSeries(params?: Record<string, any>) {
    const response = await this.request<any>('GET', `/calendar-series`, undefined, params);
    return response;
  }

  async createCalendarSeries(data: any) {
    const response = await this.request<ApiResponse<any>>('POST', `/calendar-series`, data);
    return response;
  }

  async getCalendarSeriesById(id: number) {
    if (!id) throw new Error('id is required');
    const response = await this.request<any>('GET', `/calendar-series/${id}`, undefined);
    return response;
  }

  async updateCalendarSeries(id: number, data: any) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('PUT', `/calendar-series/${id}`, data);
    return response;
  }

  async deleteCalendarSeries(id: number) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('DELETE', `/calendar-series/${id}`, undefined);
    return response || { success: true };
  }

  async getCalendars() {
    const response = await this.request<any>('GET', `/calendars`, undefined);
    return response;
  }

  async createCalendar(data: any) {
    const response = await this.request<ApiResponse<any>>('POST', `/calendars`, data);
    return response;
  }

  async getFreeSlots(params?: Record<string, any>) {
    const response = await this.request<any>('GET', `/calendars/free-slots`, undefined, params);
    return response;
  }

  async getCalendarWeek(params?: Record<string, any>) {
    const response = await this.request<ApiResponse<any>>('GET', `/calendars/week`, undefined, params);
    return response;
  }

  async getCalendarYear(params?: Record<string, any>) {
    const response = await this.request<ApiResponse<any>>('GET', `/calendars/year`, undefined, params);
    return response;
  }

  async getCalendarById(id: number) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('GET', `/calendars/${id}`, undefined);
    return response;
  }

  async updateCalendar(id: number, data: any) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('PUT', `/calendars/${id}`, data);
    return response;
  }

  async deleteCalendar(id: number) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('DELETE', `/calendars/${id}`, undefined);
    return response || { success: true };
  }

  async importHolidays(id: number, data: any) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('POST', `/calendars/${id}/import_holidays`, data);
    return response;
  }

  async getClients(params?: Record<string, any>) {
    const response = await this.request<any>('GET', `/clients`, undefined, params);
    return response;
  }

  async createClient(data: any) {
    const response = await this.request<ApiResponse<any>>('POST', `/clients`, data);
    return response;
  }

  async searchClients(params?: Record<string, any>) {
    const response = await this.request<ApiResponse<any>>('GET', `/clients/search`, undefined, params);
    return response;
  }

  async getClientById(id: number) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('GET', `/clients/${id}`, undefined);
    return response;
  }

  async updateClient(id: number, data: any) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('PUT', `/clients/${id}`, data);
    return response;
  }

  async deleteClient(id: number) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('DELETE', `/clients/${id}`, undefined);
    return response || { success: true };
  }

  async getCostProviders(params?: Record<string, any>) {
    const response = await this.request<any>('GET', `/cost-providers`, undefined, params);
    return response;
  }

  async createCostProvider(data: any) {
    const response = await this.request<ApiResponse<any>>('POST', `/cost-providers`, data);
    return response;
  }

  async searchCostProviders(params?: Record<string, any>) {
    const response = await this.request<ApiResponse<any>>('GET', `/cost-providers/search`, undefined, params);
    return response;
  }

  async getCostProviderById(id: number) {
    if (!id) throw new Error('id is required');
    const response = await this.request<any>('GET', `/cost-providers/${id}`, undefined);
    return response;
  }

  async updateCostProvider(id: number, data: any) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('PUT', `/cost-providers/${id}`, data);
    return response;
  }

  async deleteCostProvider(id: number) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('DELETE', `/cost-providers/${id}`, undefined);
    return response || { success: true };
  }

  async getExternalCalendars(params?: Record<string, any>) {
    const response = await this.request<any>('GET', `/external-calendars`, undefined, params);
    return response;
  }

  async createExternalCalendar(data: any) {
    const response = await this.request<ApiResponse<any>>('POST', `/external-calendars`, data);
    return response;
  }

  async getExternalCalendarById(id: number) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('GET', `/external-calendars/${id}`, undefined);
    return response;
  }

  async updateExternalCalendar(id: number, data: any) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('PUT', `/external-calendars/${id}`, data);
    return response;
  }

  async deleteExternalCalendar(id: number) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('DELETE', `/external-calendars/${id}`, undefined);
    return response || { success: true };
  }

}

export default AESaasApiClient;
