<template>
  <div class="min-h-screen bg-base-100 flex flex-col">
    <!-- Header -->
    <header class="bg-base-200 shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <h1 class="text-xl font-semibold text-base-content">{{ $t('dashboard.welcome') }}</h1>
          <div class="flex items-center space-x-4">
            <span class="text-sm text-base-content/70">{{ user?.email || user?.username }}</span>
            <LogoutButton />
          </div>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="flex-1 bg-base-100">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <!-- Statistics Cards -->
        <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-5 gap-6 mb-8">
          <div class="stat bg-base-200 rounded-lg shadow-sm">
            <div class="stat-figure text-primary">
              <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
              </svg>
            </div>
            <div class="stat-title">{{ $t('dashboard.plans') }}</div>
            <div class="stat-value text-primary">{{ dashboardData.plans.count }}</div>
            <div class="stat-desc">{{ $t('dashboard.activePlans') }}</div>
          </div>

          <div class="stat bg-base-200 rounded-lg shadow-sm">
            <div class="stat-figure text-secondary">
              <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path>
              </svg>
            </div>
            <div class="stat-title">{{ $t('dashboard.customers') }}</div>
            <div class="stat-value text-secondary">{{ dashboardData.customers.count }}</div>
            <div class="stat-desc">{{ $t('dashboard.totalCustomers') }}</div>
          </div>

          <div class="stat bg-base-200 rounded-lg shadow-sm">
            <div class="stat-figure text-accent">
              <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z"></path>
              </svg>
            </div>
            <div class="stat-title">{{ $t('dashboard.users') }}</div>
            <div class="stat-value text-accent">{{ dashboardData.users.count }}</div>
            <div class="stat-desc">{{ $t('dashboard.registeredUsers') }}</div>
          </div>

          <div class="stat bg-base-200 rounded-lg shadow-sm">
            <div class="stat-figure text-info">
              <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 4.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path>
              </svg>
            </div>
            <div class="stat-title">{{ $t('dashboard.emails') }}</div>
            <div class="stat-value text-info">{{ dashboardData.emails.sent }}</div>
            <div class="stat-desc">{{ $t('dashboard.emailsSent') }}</div>
          </div>

          <div class="stat bg-base-200 rounded-lg shadow-sm">
            <div class="stat-figure text-success">
              <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"></path>
              </svg>
            </div>
            <div class="stat-title">{{ $t('dashboard.newsletters') }}</div>
            <div class="stat-value text-success">{{ dashboardData.newsletters.subscribers }}</div>
            <div class="stat-desc">{{ $t('dashboard.subscribers') }}</div>
          </div>
        </div>

        <!-- Data Tables -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <!-- Recent Plans -->
          <div class="card bg-base-200 shadow-sm">
            <div class="card-body">
              <h3 class="card-title mb-4">{{ $t('dashboard.recentPlans') }}</h3>
              <div class="overflow-x-auto">
                <table class="table table-sm">
                  <thead>
                    <tr>
                      <th>{{ $t('dashboard.name') }}</th>
                      <th>{{ $t('dashboard.status') }}</th>
                      <th>{{ $t('dashboard.created') }}</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="plan in dashboardData.plans.recent" :key="plan.id">
                      <td>{{ plan.name }}</td>
                      <td>
                        <div class="badge badge-sm" :class="plan.active ? 'badge-success' : 'badge-ghost'">
                          {{ plan.active ? $t('dashboard.active') : $t('dashboard.inactive') }}
                        </div>
                      </td>
                      <td>{{ formatDate(plan.created_at) }}</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>

          <!-- Recent Customers -->
          <div class="card bg-base-200 shadow-sm">
            <div class="card-body">
              <h3 class="card-title mb-4">{{ $t('dashboard.recentCustomers') }}</h3>
              <div class="overflow-x-auto">
                <table class="table table-sm">
                  <thead>
                    <tr>
                      <th>{{ $t('dashboard.name') }}</th>
                      <th>{{ $t('dashboard.email') }}</th>
                      <th>{{ $t('dashboard.joined') }}</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="customer in dashboardData.customers.recent" :key="customer.id">
                      <td>{{ customer.name }}</td>
                      <td>{{ customer.email }}</td>
                      <td>{{ formatDate(customer.created_at) }}</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>

        <!-- Loading Indicator -->
        <div v-if="loading" class="mt-6 flex justify-center">
          <div class="loading loading-spinner loading-lg"></div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useToast } from '../composables/useToast'
import { useI18n } from 'vue-i18n'
import LogoutButton from '../components/LogoutButton.vue'
import { AESaasApiClient } from '@agile-exec/api-client'

// Create API client instance
const apiClient = new AESaasApiClient({ baseURL: 'http://localhost:8080' })

const authStore = useAuthStore()
const toast = useToast()
const { t } = useI18n()

const user = computed(() => authStore.user)
const loading = ref(false)

// Dashboard data structure
const dashboardData = ref({
  plans: {
    count: 0,
    recent: [] as any[]
  },
  customers: {
    count: 0,
    recent: [] as any[]
  },
  users: {
    count: 0
  },
  emails: {
    sent: 0
  },
  newsletters: {
    subscribers: 0
  }
})

// Format date helper
const formatDate = (dateString: string) => {
  if (!dateString) return 'N/A'
  try {
    return new Date(dateString).toLocaleDateString('de-DE', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit'
    })
  } catch (error) {
    return 'N/A'
  }
}

// Load dashboard data
const loadDashboardData = async () => {
  loading.value = true
  try {
    // Load plans data
    try {
      const plansResponse = await apiClient.getPlans()
      dashboardData.value.plans.count = plansResponse.data?.length || 0
      dashboardData.value.plans.recent = (plansResponse.data || []).slice(0, 5)
    } catch (error) {
      console.warn('Could not load plans data:', error)
    }

    // Load customers data
    try {
      const customersResponse = await apiClient.getCustomers()
      dashboardData.value.customers.count = customersResponse.data?.length || 0
      dashboardData.value.customers.recent = (customersResponse.data || []).slice(0, 5)
    } catch (error) {
      console.warn('Could not load customers data:', error)
    }

    // Load email stats
    try {
      const emailStatsResponse = await apiClient.getEmailStats()
      dashboardData.value.emails.sent = emailStatsResponse.data?.total_sent || 0
    } catch (error) {
      console.warn('Could not load email stats:', error)
    }

    // Load newsletter subscribers
    try {
      const newsletterResponse = await apiClient.getNewsletterSubscriptions()
      dashboardData.value.newsletters.subscribers = newsletterResponse.data?.length || 0
    } catch (error) {
      console.warn('Could not load newsletter data:', error)
    }

    // For users count, we'll use a placeholder since there might not be a direct endpoint
    dashboardData.value.users.count = 1 // At least the current user

  } catch (error) {
    console.error('Error loading dashboard data:', error)
    toast.error(t('dashboard.loadError'))
  } finally {
    loading.value = false
  }
}

// Load data on component mount
onMounted(() => {
  loadDashboardData()
})
</script>