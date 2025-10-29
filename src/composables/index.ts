import { ref, reactive, computed, type Ref } from 'vue';
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
      user.value = response.user;
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
      
      const method = `get${resourceName.charAt(0).toUpperCase() + resourceName.slice(1)}` as keyof AESaasApiClient;
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
      error.value = err.message || `Failed to fetch ${resourceName}`;
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const createItem = async (data: Partial<T>) => {
    try {
      loading.value = true;
      error.value = null;
      
      const method = `create${resourceName.charAt(0).toUpperCase() + resourceName.slice(1, -1)}` as keyof AESaasApiClient;
      const item = await (client[method] as Function)(data);
      
      items.value.unshift(item);
      return item;
    } catch (err: any) {
      error.value = err.message || `Failed to create ${resourceName.slice(0, -1)}`;
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const updateItem = async (id: number, data: Partial<T>) => {
    try {
      loading.value = true;
      error.value = null;
      
      const method = `update${resourceName.charAt(0).toUpperCase() + resourceName.slice(1, -1)}` as keyof AESaasApiClient;
      const updatedItem = await (client[method] as Function)(id, data);
      
      const index = items.value.findIndex((item: any) => item.id === id);
      if (index !== -1) {
        items.value[index] = updatedItem;
      }
      
      return updatedItem;
    } catch (err: any) {
      error.value = err.message || `Failed to update ${resourceName.slice(0, -1)}`;
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const deleteItem = async (id: number) => {
    try {
      loading.value = true;
      error.value = null;
      
      const method = `delete${resourceName.charAt(0).toUpperCase() + resourceName.slice(1, -1)}` as keyof AESaasApiClient;
      await (client[method] as Function)(id);
      
      items.value = items.value.filter((item: any) => item.id !== id);
    } catch (err: any) {
      error.value = err.message || `Failed to delete ${resourceName.slice(0, -1)}`;
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
export function useSearch() {
  const client = useApiClient();
  const results = ref([]);
  const loading = ref(false);
  const error = ref<string | null>(null);
  const query = ref('');

  const search = async (searchQuery: string, options?: any) => {
    try {
      loading.value = true;
      error.value = null;
      query.value = searchQuery;
      
      const searchResults = await client.search({ q: searchQuery, ...options });
      results.value = searchResults.data || [];
      
      return searchResults;
    } catch (err: any) {
      error.value = err.message || 'Search failed';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const quickSearch = async (searchQuery: string, limit = 5) => {
    try {
      loading.value = true;
      error.value = null;
      query.value = searchQuery;
      
      const searchResults = await client.quickSearch({ q: searchQuery, limit });
      results.value = searchResults.data || searchResults || [];
      
      return searchResults;
    } catch (err: any) {
      error.value = err.message || 'Quick search failed';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const clearResults = () => {
    results.value = [];
    query.value = '';
    error.value = null;
  };

  return {
    results: readonly(results),
    loading: readonly(loading),
    error: readonly(error),
    query: readonly(query),
    search,
    quickSearch,
    clearResults
  };
}

// Health check composable
export function useHealth() {
  const client = useApiClient();
  const status = ref<string | null>(null);
  const loading = ref(false);
  const error = ref<string | null>(null);
  const lastChecked = ref<Date | null>(null);

  const checkHealth = async () => {
    try {
      loading.value = true;
      error.value = null;
      
      const health = await client.health();
      status.value = health.status;
      lastChecked.value = new Date();
      
      return health;
    } catch (err: any) {
      error.value = err.message || 'Health check failed';
      status.value = 'unhealthy';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const ping = async () => {
    try {
      const response = await client.ping();
      return response;
    } catch (err: any) {
      error.value = err.message || 'Ping failed';
      throw err;
    }
  };

  const isHealthy = computed(() => status.value === 'healthy');

  return {
    status: readonly(status),
    loading: readonly(loading),
    error: readonly(error),
    lastChecked: readonly(lastChecked),
    isHealthy,
    checkHealth,
    ping
  };
}

// Utility composables
function readonly<T>(ref: Ref<T>) {
  return computed(() => ref.value);
}
