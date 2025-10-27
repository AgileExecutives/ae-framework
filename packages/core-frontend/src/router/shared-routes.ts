// Shared routes configuration from base-app
// These routes can be imported and extended for different apps
import type { Router } from 'vue-router'

export function createBaseRoutes() {
  return [
    {
      path: '/login',
      name: 'Login',
      component: () => import('../views/LoginView.vue'),
      meta: { requiresAuth: false, requiresGuest: true },
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('../views/RegisterView.vue'),
      meta: { requiresAuth: false, requiresGuest: true },
    },
    {
      path: '/forgot-password',
      name: 'ForgotPassword',
      component: () => import('../views/ForgotPasswordView.vue'),
      meta: { requiresAuth: false },
    },
    {
      path: '/new-password/:token',
      name: 'NewPassword',
      component: () => import('../views/ResetPasswordView.vue'),
      meta: { requiresAuth: false },
    },
    {
      path: '/change-password',
      name: 'ChangePassword',
      component: () => import('../views/ChangePasswordView.vue'),
      meta: { requiresAuth: true },
    }
  ]
}

export function addCatchAllRoute( router: Router ) {
  return router.addRoute({
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('../views/NotFoundView.vue'),
    meta: { requiresAuth: false },
  })
}
