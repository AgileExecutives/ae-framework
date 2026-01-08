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

  // Auth methods (not auto-generated, manually maintained)
  async login(credentials: { email: string; password: string }) {
    const response = await this.request<ApiResponse<{ token: string; user: any }>>('POST', '/auth/login', credentials);
    return response;
  }

  async register(credentials: { email: string; password: string; name?: string }) {
    const response = await this.request<ApiResponse<{ token: string; user: any }>>('POST', '/auth/register', credentials);
    return response;
  }

  async logout() {
    const response = await this.request<ApiResponse<any>>('POST', '/auth/logout', undefined);
    return response;
  }

  async getCurrentUser() {
    const response = await this.request<ApiResponse<any>>('GET', '/auth/me', undefined);
    return response;
  }

  async changePassword(credentials: { current_password: string; new_password: string }) {
    const response = await this.request<ApiResponse<any>>('POST', '/auth/change-password', credentials);
    return response;
  }

  async forgotPassword(email: string) {
    const response = await this.request<ApiResponse<any>>('POST', '/auth/forgot-password', { email });
    return response;
  }

  async resetPassword(token: string, newPassword: string) {
    const response = await this.request<ApiResponse<any>>('POST', '/auth/reset-password', { token, new_password: newPassword });
    return response;
  }

  async getPlans() {
    const response = await this.request<ApiResponse<any>>('GET', '/plans', undefined);
    return response;
  }

  async getCustomers(params?: Record<string, any>) {
    const response = await this.request<ApiResponse<any>>('GET', '/customers', undefined, params);
    return response;
  }

  async getEmailStats() {
    const response = await this.request<ApiResponse<any>>('GET', '/email-stats', undefined);
    return response;
  }

  async getNewsletterSubscriptions(params?: Record<string, any>) {
    const response = await this.request<ApiResponse<any>>('GET', '/newsletter-subscriptions', undefined, params);
    return response;
  }

  // Calendar methods (not in Swagger spec)
  async getCalendars(params?: Record<string, any>) {
    const response = await this.request<any>('GET', '/calendars', undefined, params);
    return response;
  }

  async createCalendarEntry(data: any) {
    const response = await this.request<ApiResponse<any>>('POST', '/calendar/entries', data);
    return response;
  }

  async updateCalendarEntry(id: number, data: any) {
    const response = await this.request<ApiResponse<any>>('PUT', `/calendar/entries/${id}`, data);
    return response;
  }

  async deleteCalendarEntry(id: number) {
    const response = await this.request<ApiResponse<any>>('DELETE', `/calendar/entries/${id}`, undefined);
    return response;
  }

  async getCalendarEntryById(id: number) {
    const response = await this.request<ApiResponse<any>>('GET', `/calendar/entries/${id}`, undefined);
    return response;
  }

  async createCalendarSeries(data: any) {
    const response = await this.request<ApiResponse<any>>('POST', '/calendar/series', data);
    return response;
  }

  async updateCalendarSeries(id: number, data: any) {
    const response = await this.request<ApiResponse<any>>('PUT', `/calendar/series/${id}`, data);
    return response;
  }

  async deleteCalendarSeries(id: number, options?: any) {
    const response = await this.request<ApiResponse<any>>('DELETE', `/calendar/series/${id}`, options);
    return response;
  }

  async updateCalendar(id: number, data: any) {
    const response = await this.request<ApiResponse<any>>('PUT', `/calendars/${id}`, data);
    return response;
  }

  // Additional booking methods
  async listBookingTemplatesByUser(params?: Record<string, any>) {
    const response = await this.request<ApiResponse<any>>('GET', '/booking/templates/by-user', undefined, params);
    return response;
  }

  async getStaticFile(filename: string) {
    const response = await this.request<any>('GET', `/static/${filename}`, undefined);
    return response;
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
  async getEntityAuditLogs(entity_type: string, entity_id: number) {
    if (!entity_type) throw new Error('entity_type is required');
    if (!entity_id) throw new Error('entity_id is required');
    const response = await this.request<any>('GET', `/audit/entity/${entity_type}/${entity_id}`, undefined);
    return response;
  }

  async exportAuditLogs(params?: Record<string, any>) {
    const response = await this.request<ApiResponse<any>>('GET', `/audit/export`, undefined, params);
    return response;
  }

  async getAuditLogs(params?: Record<string, any>) {
    const response = await this.request<any>('GET', `/audit/logs`, undefined, params);
    return response;
  }

  async getAuditStatistics(params?: Record<string, any>) {
    const response = await this.request<any>('GET', `/audit/statistics`, undefined, params);
    return response;
  }

  async getInvoices(params?: Record<string, any>) {
    const response = await this.request<any>('GET', `/client-invoices`, undefined, params);
    return response;
  }

  async createInvoice(data: any) {
    const response = await this.request<ApiResponse<any>>('POST', `/client-invoices`, data);
    return response;
  }

  async getClientsWithUnbilledSessions() {
    const response = await this.request<any>('GET', `/client-invoices/unbilled-sessions`, undefined);
    return response;
  }

  async getInvoiceById(id: number) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('GET', `/client-invoices/${id}`, undefined);
    return response;
  }

  async updateInvoice(id: number, data: any) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('PUT', `/client-invoices/${id}`, data);
    return response;
  }

  async deleteInvoice(id: number) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('DELETE', `/client-invoices/${id}`, undefined);
    return response || { success: true };
  }

  async getClientByToken(token: string) {
    if (!token) throw new Error('token is required');
    const response = await this.request<ApiResponse<any>>('GET', `/client/${token}`, undefined);
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

  async getSessionsByClient(id: number, params?: Record<string, any>) {
    if (!id) throw new Error('id is required');
    const response = await this.request<any>('GET', `/clients/${id}/sessions`, undefined, params);
    return response;
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

  async listExtraEfforts(params?: Record<string, any>) {
    const response = await this.request<ApiResponse<any>>('GET', `/extra-efforts`, undefined, params);
    return response;
  }

  async createExtraEffort(data: any) {
    const response = await this.request<ApiResponse<any>>('POST', `/extra-efforts`, data);
    return response;
  }

  async getExtraEffort(id: number) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('GET', `/extra-efforts/${id}`, undefined);
    return response;
  }

  async updateExtraEffort(id: number, data: any) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('PUT', `/extra-efforts/${id}`, data);
    return response;
  }

  async deleteExtraEffort(id: number) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('DELETE', `/extra-efforts/${id}`, undefined);
    return response || { success: true };
  }

  async createDraftInvoice(data: any) {
    const response = await this.request<ApiResponse<any>>('POST', `/invoices/draft`, data);
    return response;
  }

  async getVATCategories() {
    const response = await this.request<any>('GET', `/invoices/vat-categories`, undefined);
    return response;
  }

  async updateDraftInvoice(id: number, data: any) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('PUT', `/invoices/${id}`, data);
    return response;
  }

  async cancelDraftInvoice(id: number) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('DELETE', `/invoices/${id}`, undefined);
    return response || { success: true };
  }

  async createCreditNote(id: number, data: any) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('POST', `/invoices/${id}/credit-note`, data);
    return response;
  }

  async finalizeInvoice(id: number) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('POST', `/invoices/${id}/finalize`, undefined);
    return response;
  }

  async markInvoiceAsOverdue(id: number) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('POST', `/invoices/${id}/mark-overdue`, undefined);
    return response;
  }

  async markInvoiceAsPaid(id: number, data: any) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('POST', `/invoices/${id}/mark-paid`, data);
    return response;
  }

  async sendReminder(id: number) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('POST', `/invoices/${id}/reminder`, undefined);
    return response;
  }

  async sendInvoiceEmail(id: number) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('POST', `/invoices/${id}/send-email`, undefined);
    return response;
  }

  async exportXRechnung(id: number) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('GET', `/invoices/${id}/xrechnung`, undefined);
    return response;
  }

  async getAllSessions(params?: Record<string, any>) {
    const response = await this.request<any>('GET', `/sessions`, undefined, params);
    return response;
  }

  async createSession(data: any) {
    const response = await this.request<ApiResponse<any>>('POST', `/sessions`, data);
    return response;
  }

  async bookSessions(data: any) {
    const response = await this.request<ApiResponse<any>>('POST', `/sessions/book`, data);
    return response;
  }

  async bookSessionsWithToken(token: string, data: any) {
    if (!token) throw new Error('token is required');
    const response = await this.request<ApiResponse<any>>('POST', `/sessions/book/${token}`, data);
    return response;
  }

  async getSessionByCalendarEntry(id: number) {
    if (!id) throw new Error('id is required');
    const response = await this.request<any>('GET', `/sessions/by_entry/${id}`, undefined);
    return response;
  }

  async getDetailedSessionsUpcoming(params?: Record<string, any>) {
    const response = await this.request<any>('GET', `/sessions/detail`, undefined, params);
    return response;
  }

  async getSessionById(id: number) {
    if (!id) throw new Error('id is required');
    const response = await this.request<any>('GET', `/sessions/${id}`, undefined);
    return response;
  }

  async updateSession(id: number, data: any) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('PUT', `/sessions/${id}`, data);
    return response;
  }

  async deleteSession(id: number) {
    if (!id) throw new Error('id is required');
    const response = await this.request<ApiResponse<any>>('DELETE', `/sessions/${id}`, undefined);
    return response || { success: true };
  }

}

export default AESaasApiClient;
