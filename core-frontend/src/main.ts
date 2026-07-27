// main.ts
import { createBaseApp } from '@@/init.ts'
import routes from './router/routes.ts'
import { createApiClient } from './config/api-config.ts'
import { useAuthStore } from '@@/stores/auth'
import MainApp from './App.vue'
import './style.css'

// Initialize the global API client (default to relative path so Vite proxy applies in dev)
createApiClient()

const {app} = createBaseApp(MainApp, {
    routes
} )

// Initialize auth store after app is mounted
const authStore = useAuthStore()
authStore.initializeAuth()

// Mount the app
app.mount('#app')
