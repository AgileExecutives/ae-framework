// =============================================================================
// BASE-APP MINIMAL WORKING EXPORTS
// =============================================================================
// Only exports that work without dependency issues

import { createApp } from 'vue'
import { createI18n } from 'vue-i18n'
import type { RouteRecordRaw } from 'vue-router'

// Import language files
import en from './locales/en'
import de from './locales/de'
import router from './router/index.ts'
import { createPinia } from 'pinia'

// Create base routes function (re-exported for compatibility)
export { createBaseRoutes } from './router/shared-routes.ts'
import { addCatchAllRoute } from './router/shared-routes.ts'

// auth route protection middleware
export { addMiddleware } from './router/middleware.ts'
import { addMiddleware } from './router/middleware.ts'

// Auth store
export * from './stores/auth-types.ts'
export { useAuthStore } from './stores/auth.ts'

export interface CreateBaseAppOptions {
    routes?: RouteRecordRaw[]
    overrideRoot?: boolean  // If true, removes the default '/' route
}

export function createBaseApp(AppComponent: any, options?: CreateBaseAppOptions) {

    const messages = {
        en,
        de,
    }
    const i18n = createI18n({
        legacy: false,
        locale: 'de',
        messages,
    })
    
    const app = createApp(AppComponent)
    const pinia = createPinia()
    app.use(pinia)
    
    // Remove the default root route if overrideRoot is true
    if (options?.overrideRoot) {
        router.removeRoute('AuthRestrictedView')
    }
    
    // Add additional routes if provided
    if (options?.routes) {
        options.routes.forEach(route => {
            router.addRoute(route)
        })
    }
    
    // Add catch-all route
    addCatchAllRoute(router)

    app.use(router)
    app.use(i18n)
    
    // Add middleware AFTER Pinia is registered
    addMiddleware(router)

    return app
}

// Utilities
export { cn } from './lib/utils'

// App Configuration
export * from './config/app-config'
export * from './config/api-config'
export { useAppConfig } from './composables/useAppConfig'
export { useDarkMode } from './composables/useDarkMode'

// Core Components
export * from './components'

// Additional component exports for direct access
export { default as SingleFormCard } from './components/SingleFormCard.vue'
export { default as AuthLayout } from './components/AuthLayout.vue'
export { default as ThemeToggle } from './components/ThemeToggle.vue'
export { default as LocaleSwitcher } from './components/LocaleSwitcher.vue'
export { default as LegalLinks } from './components/LegalLinks.vue'

