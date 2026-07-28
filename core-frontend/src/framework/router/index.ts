import { createRouter, createWebHistory } from 'vue-router'
import { createBaseRoutes } from './shared-routes.js'
import DashboardView from '@@/views/DashboardView.vue'

// Create default routes using base components
const routes = [
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
