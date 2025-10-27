import { useAuthStore } from '../stores/auth'

export function addMiddleware(router: any) {
    // Authentication guard - temporarily disabled for testing
    router.beforeEach(async (to: any, _from: any, next: any) => {
        console.log('🛡️ Route guard: navigating to', to.path, 'from', _from?.path)
        console.log('🛡️ Route guard: route meta:', to.meta)
        console.log('🛡️ Route guard: matched records:', to.matched.map((r: any) => ({ name: r.name, meta: r.meta })))

        const authStore = useAuthStore()
        console.log('🛡️ Route guard: current auth state:', {
            isAuthenticated: authStore.isAuthenticated,
            hasUser: !!authStore.user,
            hasToken: !!authStore.token
        })

        // Check if route explicitly does not require auth
        const requiresAuth = to.matched.some((record: any) => record.meta?.requiresAuth === true)
        const explicitlyPublic = to.matched.some((record: any) => record.meta?.requiresAuth === false)

        console.log('🛡️ Route guard: requiresAuth =', requiresAuth, ', explicitlyPublic =', explicitlyPublic)

        // If route is explicitly public, skip all auth checks
        if (explicitlyPublic) {
            console.log('🛡️ Route guard: route is explicitly public, allowing access')
            next()
            return
        }

        // Initialize auth if not already done
        if (!authStore.initialized) {
            console.log('🛡️ Route guard: auth not initialized, initializing...')
            try {
                await authStore.initializeAuth()
                console.log('🛡️ Route guard: auth initialized successfully')
            } catch (error) {
                console.error('🛡️ Route guard: auth initialization failed:', error)
            }
        }

        // Check if route requires authentication
        if (requiresAuth) {
            console.log('🛡️ Route guard: route requires auth')
            if (!authStore.isAuthenticated) {
                console.log('🛡️ Route guard: not authenticated, redirecting to login')
                next({ name: 'Login', query: { redirect: to.fullPath } })
                return
            }

            // Check if route requires super admin role
            if (to.matched.some((record: any) => record.meta?.requiresSuperAdmin)) {
                console.log('🛡️ Route guard: route requires super admin')
                const userRole = authStore.user?.role as string
                if (userRole !== 'super-admin') {
                    console.log('🛡️ Route guard: insufficient permissions, redirecting to dashboard')
                    next({ name: 'Dashboard' })
                    return
                }
            }

            console.log('🛡️ Route guard: authenticated, allowing access')
        }

        // Check if route is for guests only (like login page)
        if (to.matched.some((record: any) => record.meta?.requiresGuest)) {
            console.log('🛡️ Route guard: guest route')
            if (authStore.isAuthenticated) {
                console.log('🛡️ Route guard: already authenticated, redirecting to dashboard')
                next({ name: 'Dashboard' })
                return
            }
        }

        console.log('🛡️ Route guard: allowing navigation')
        next()
    })

}

