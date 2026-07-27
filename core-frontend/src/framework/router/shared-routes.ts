// Shared routes configuration from base-app
// These routes can be imported and extended for different apps
import type { Router } from 'vue-router'
import LoginView from '@@/views/LoginView.vue'
import RegisterView from '@@/views/RegisterView.vue'
import VerifyEmailView from '@@/views/VerifyEmailView.vue'
import ForgotPasswordView from '@@/views/ForgotPasswordView.vue'
import ResetPasswordView from '@@/views/ResetPasswordView.vue'
import ChangePasswordView from '@@/views/ChangePasswordView.vue'
import SuccessPage from '@@/views/SuccessPage.vue'
import ErrorPage from '@@/views/ErrorPage.vue'
import NotFoundView from '@@/views/NotFoundView.vue'

export function createBaseRoutes() {
  return [
    {
      path: '/login',
      name: 'Login',
      component: LoginView,
      meta: { requiresAuth: false, requiresGuest: true },
    },
    {
      path: '/register',
      name: 'Register',
      component: RegisterView,
      meta: { requiresAuth: false, requiresGuest: true },
    },
    {
      path: '/verify-email',
      name: 'VerifyEmail',
      component: VerifyEmailView,
      meta: { requiresAuth: false },
    },
    {
      path: '/forgot-password',
      name: 'ForgotPassword',
      component: ForgotPasswordView,
      meta: { requiresAuth: false },
    },
    {
      path: '/new-password',
      name: 'NewPassword',
      component: ResetPasswordView,
      meta: { requiresAuth: false },
    },
    {
      path: '/change-password',
      name: 'ChangePassword',
      component: ChangePasswordView,
      meta: { requiresAuth: true },
    },
    {
      path: '/success',
      name: 'Success',
      component: SuccessPage,
      meta: { requiresAuth: false },
    },
    {
      path: '/error',
      name: 'Error',
      component: ErrorPage,
      meta: { requiresAuth: false },
    }
  ]
}

export function addCatchAllRoute( router: Router ) {
  return router.addRoute({
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: NotFoundView,
    meta: { requiresAuth: false },
  })
}
