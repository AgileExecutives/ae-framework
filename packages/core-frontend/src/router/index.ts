import { createRouter, createWebHistory } from 'vue-router'
import { createBaseRoutes } from './shared-routes.js'
import  AuthRestrictedView from '../views/AuthRestrictedView.vue'

// Create default routes using base components
const routes = [
  {
    path: '/restricted',
    name: 'AuthRestrictedView',
    component: AuthRestrictedView,
    meta: { requiresAuth: true },
  },
  ...createBaseRoutes()
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// Don't add middleware here - it will be added in createBaseApp after Pinia is ready
// addMiddleware(router)

export { routes }
export default router
