# AE SaaS API Client for Vue

A TypeScript-first, Vue 3 compatible API client that auto-generates from Swagger/OpenAPI specifications. Built specifically for Vue applications with composables, reactive state management, and full type safety.

## ✨ Features

- 🚀 **Auto-generated from Swagger** - Keeps in sync with your API
- 📱 **Vue 3 Composables** - Reactive state management out of the box
- 🔒 **TypeScript First** - Full type safety and IntelliSense
- 🎯 **Authentication Handling** - JWT token management
- 📊 **CRUD Operations** - Resource management made easy
- 🔍 **Search Functionality** - Built-in search with pagination
- ⚡ **ESM/CJS Support** - Works in any modern bundler
- 🧪 **Vue Plugin** - Easy integration with Vue apps

## 📦 Installation

```bash
npm install @unburdy/ae-saas-api-client
```

## 🚀 Quick Start

### 1. Install the Plugin

```typescript
// main.ts
import { createApp } from 'vue';
import App from './App.vue';
import { AESaasApiClientPlugin } from '@unburdy/ae-saas-api-client';

const app = createApp(App);

app.use(AESaasApiClientPlugin, {
  baseURL: 'http://localhost:8080/api/v1',
  timeout: 10000
});

app.mount('#app');
```

### 2. Use in Components

```vue
<template>
  <div>
    <!-- Authentication -->
    <div v-if="!isAuthenticated">
      <input v-model="credentials.username" placeholder="Username" />
      <input v-model="credentials.password" type="password" placeholder="Password" />
      <button @click="login(credentials)" :disabled="loading">
        {{ loading ? 'Logging in...' : 'Login' }}
      </button>
    </div>
    
    <!-- Customer Management -->
    <div v-else>
      <h2>Welcome, {{ user?.username }}!</h2>
      
      <!-- Customer List -->
      <div v-for="customer in customers" :key="customer.id">
        {{ customer.name }} - {{ customer.email }}
      </div>
      
      <!-- Pagination -->
      <button @click="prevPage" :disabled="pagination.page <= 1">Previous</button>
      <span>Page {{ pagination.page }} of {{ pagination.totalPages }}</span>
      <button @click="nextPage" :disabled="pagination.page >= pagination.totalPages">Next</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, onMounted } from 'vue';
import { useAuth, useCustomers } from '@unburdy/ae-saas-api-client';

// Authentication
const { isAuthenticated, user, loading, login } = useAuth();
const credentials = reactive({ username: '', password: '' });

// Customer management
const { 
  items: customers, 
  pagination, 
  fetchItems, 
  prevPage, 
  nextPage 
} = useCustomers();

onMounted(() => {
  if (isAuthenticated.value) {
    fetchItems();
  }
});
</script>
```

## 📚 Available Composables

### Authentication
- `useAuth()` - Login, logout, current user management
- `useHealth()` - API health checks and ping

### Resources
- `useCustomers()` - Customer CRUD operations
- `useContacts()` - Contact management
- `usePlans()` - Plan management
- `useEmails()` - Email operations

### Utilities
- `useSearch()` - Search functionality with quick search
- `useApiClient()` - Direct access to the API client

## 🔧 Auto-Generation

This package can be regenerated from your Swagger/OpenAPI specification:

```bash
# Generate from local server
npm run generate

# Generate from custom URL
SWAGGER_URL=https://your-api.com/swagger/doc.json npm run generate

# Build the package
npm run build
```

### Custom Generation Script

```javascript
// In your CI/CD or development scripts
import { execSync } from 'child_process';

// Regenerate when API changes
execSync('cd path/to/vue-api-client && npm run generate && npm run build');
```

## 📖 API Reference

### useAuth()

```typescript
const {
  isAuthenticated,    // Ref<boolean> - Current auth state
  user,              // Ref<User | null> - Current user data
  loading,           // Ref<boolean> - Loading state
  error,             // Ref<string | null> - Error message
  login,             // (credentials) => Promise<LoginResponse>
  logout,            // () => Promise<void>
  getCurrentUser     // () => Promise<User>
} = useAuth();
```

### useResource() (useCustomers, useContacts, etc.)

```typescript
const {
  items,             // Ref<T[]> - List of items
  loading,           // Ref<boolean> - Loading state
  error,             // Ref<string | null> - Error message
  pagination,        // ComputedRef<Pagination> - Pagination info
  fetchItems,        // (params?) => Promise<ListResponse<T>>
  createItem,        // (data) => Promise<T>
  updateItem,        // (id, data) => Promise<T>
  deleteItem,        // (id) => Promise<void>
  nextPage,          // () => void
  prevPage,          // () => void
  setPage            // (page: number) => void
} = useCustomers();
```

### useSearch()

```typescript
const {
  results,           // Ref<any[]> - Search results
  loading,           // Ref<boolean> - Loading state
  error,             // Ref<string | null> - Error message
  query,             // Ref<string> - Current query
  search,            // (query, options?) => Promise<SearchResults>
  quickSearch,       // (query, limit?) => Promise<SearchResults>
  clearResults       // () => void
} = useSearch();
```

## 🏗️ Direct API Client Usage

```typescript
import { AESaasApiClient } from '@unburdy/ae-saas-api-client';

const client = new AESaasApiClient({
  baseURL: 'http://localhost:8080/api/v1'
});

// Direct usage
const loginResponse = await client.login({ username, password });
const customers = await client.getCustomers({ page: 1, limit: 10 });
const customer = await client.createCustomer({ name: 'Acme Corp', email: 'contact@acme.com' });
```

## 🔒 Authentication

The client automatically manages JWT tokens:

```typescript
// Login sets the token automatically
await client.login({ username: 'user', password: 'pass' });

// All subsequent requests include the Bearer token
const user = await client.getCurrentUser(); // ✅ Authenticated

// Clear token on logout
await client.logout();
```

## 🛠️ TypeScript Support

Full type safety with auto-generated types:

```typescript
import type { 
  Customer, 
  CustomerCreateRequest, 
  ApiResponse,
  ListResponse 
} from '@unburdy/ae-saas-api-client';

// Type-safe operations
const newCustomer: CustomerCreateRequest = {
  name: 'Acme Corp',
  email: 'contact@acme.com'
};

const response: ApiResponse<Customer> = await client.createCustomer(newCustomer);
```

## 📊 Error Handling

```typescript
const { error } = useCustomers();

// Reactive error state
watch(error, (err) => {
  if (err) {
    toast.error(err);
  }
});

// Manual error handling
try {
  await createItem(newCustomer);
} catch (error) {
  if (error instanceof ApiError) {
    console.log('Status:', error.status);
    console.log('Data:', error.data);
  }
}
```

## 🎯 Examples

See the `examples/` directory for complete component examples:

- `LoginComponent.vue` - Authentication flow
- `CustomerList.vue` - CRUD operations with pagination
- `SearchComponent.vue` - Search functionality
- `setup.ts` - App configuration

## 📝 License

MIT

## 🤝 Contributing

1. Fork the repository
2. Make your changes
3. Run `npm run generate` to update the client
4. Run `npm run build` to build the package
5. Submit a pull request

## 🔄 Regeneration Workflow

For regular updates from your API:

1. **Manual**: Run `npm run generate` when your API changes
2. **CI/CD**: Set up automated regeneration on API deployments
3. **Development**: Use with `npm run dev` for watch mode during development

The generated client stays in sync with your Swagger specification, ensuring type safety and API compatibility.