import  AuthRestrictedView from '@/views/AuthRestrictedView.vue'
import DashboardView from '@@/views/DashboardView.vue'

const routes = [
    {
    path: '/',
    name: 'Root',
    redirect: '/dashboard'
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: DashboardView,
    meta: { requiresAuth: true },
  },
  {
    path: '/restricted',
    name: 'AuthRestrictedView',
    component: AuthRestrictedView,
    meta: { requiresAuth: true },
  }
]
export default routes