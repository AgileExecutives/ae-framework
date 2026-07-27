import  AuthRestrictedView from '@/views/AuthRestrictedView.vue'

const routes = [
  {
    path: '/restricted',
    name: 'AuthRestrictedView',
    component: AuthRestrictedView,
    meta: { requiresAuth: true },
  }
]
export default routes