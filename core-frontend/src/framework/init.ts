// =============================================================================
// BASE-APP MINIMAL WORKING EXPORTS
// =============================================================================
// Only exports that work without dependency issues

import { createApp } from 'vue'
import { createI18n } from 'vue-i18n'
import type { RouteRecordRaw } from 'vue-router'

// Import language files
import en from './locales/en.ts'
import de from './locales/de.ts'
import routerInstance from './router/index.ts'
import { createPinia } from 'pinia'

// Export router for app-level access
export const router = routerInstance

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
        routerInstance.removeRoute('AuthRestrictedView')
    }
    
    // Add additional routes if provided
    if (options?.routes) {
        options.routes.forEach(route => {
            routerInstance.addRoute(route)
        })
    }
    
    // Add catch-all route
    addCatchAllRoute(routerInstance)

    app.use(routerInstance)
    app.use(i18n)
    
    // Add middleware AFTER Pinia is registered
    addMiddleware(routerInstance)

    return { app, router: routerInstance }
}

// App Configuration
export { useAppConfig } from './composables/useAppConfig.ts'
export { useDarkMode } from './composables/useDarkMode.ts'
export { useToast } from './composables/useToast.ts'

// Core Components
export * from './components/index.ts'

// Additional component exports for direct access
export { default as SingleFormCard } from './components/SingleFormCard.vue'
export { default as AuthLayout } from './components/AuthLayout.vue'
export { default as ThemeToggle } from './components/ThemeToggle.vue'
export { default as LocaleSwitcher } from './components/LocaleSwitcher.vue'
export { default as LegalLinks } from './components/LegalLinks.vue'

