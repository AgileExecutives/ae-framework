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
function readonly<T>(source: Ref<T> | T) {
  // Check if it's a Ref by looking for the 'value' property
  if (source && typeof source === 'object' && 'value' in source) {
    // It's a Ref, return computed with .value
    return computed(() => (source as Ref<T>).value);
  } else {
    // It's a reactive object, return computed without .value
    return computed(() => source as T);
  }
}
