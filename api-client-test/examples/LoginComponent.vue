<template>
  <div class="login-form">
    <h2>Login</h2>
    
    <form @submit.prevent="handleLogin">
      <div class="form-group">
        <label for="username">Username:</label>
        <input
          id="username"
          v-model="credentials.username"
          type="text"
          required
          :disabled="loading"
        />
      </div>
      
      <div class="form-group">
        <label for="password">Password:</label>
        <input
          id="password"
          v-model="credentials.password"
          type="password"
          required
          :disabled="loading"
        />
      </div>
      
      <button type="submit" :disabled="loading">
        {{ loading ? 'Logging in...' : 'Login' }}
      </button>
      
      <div v-if="error" class="error">
        {{ error }}
      </div>
    </form>
    
    <div v-if="isAuthenticated && user" class="user-info">
      <h3>Welcome, {{ user.username }}!</h3>
      <p>Email: {{ user.email }}</p>
      <button @click="handleLogout">Logout</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { useAuth } from '@unburdy/ae-saas-api-client';

// Use the auth composable
const { isAuthenticated, user, loading, error, login, logout } = useAuth();

// Form data
const credentials = reactive({
  username: '',
  password: ''
});

// Handle login
const handleLogin = async () => {
  try {
    await login(credentials);
    // Reset form on success
    credentials.username = '';
    credentials.password = '';
  } catch (err) {
    // Error is automatically handled by the composable
    console.error('Login failed:', err);
  }
};

// Handle logout
const handleLogout = async () => {
  try {
    await logout();
  } catch (err) {
    console.error('Logout failed:', err);
  }
};
</script>

<style scoped>
.login-form {
  max-width: 400px;
  margin: 0 auto;
  padding: 2rem;
  border: 1px solid #ddd;
  border-radius: 8px;
}

.form-group {
  margin-bottom: 1rem;
}

label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: bold;
}

input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 1rem;
}

input:disabled {
  background-color: #f5f5f5;
  cursor: not-allowed;
}

button {
  width: 100%;
  padding: 0.75rem;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
}

button:disabled {
  background-color: #6c757d;
  cursor: not-allowed;
}

button:hover:not(:disabled) {
  background-color: #0056b3;
}

.error {
  margin-top: 1rem;
  padding: 0.75rem;
  background-color: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
  border-radius: 4px;
}

.user-info {
  margin-top: 2rem;
  padding: 1rem;
  background-color: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
  border-radius: 4px;
}

.user-info button {
  margin-top: 1rem;
  background-color: #dc3545;
}

.user-info button:hover {
  background-color: #c82333;
}
</style>