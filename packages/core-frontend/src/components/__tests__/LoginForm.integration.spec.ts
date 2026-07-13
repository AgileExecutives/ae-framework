import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createMemoryHistory } from 'vue-router'
import applyAppStyles from '../../plugins/appStyle'

// Mock i18n used by component
vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string) => {
      const map: Record<string, string> = {
        'auth.signIn': 'Sign In',
        'auth.signInSubtitle': 'Sign in to your account',
        'auth.email': 'Email',
        'auth.password': 'Password',
        'auth.forgotPassword': 'Forgot?',
        'messages.loginSuccessful': 'Login successful',
        'messages.loginSuccessfulWelcome': 'Welcome back',
        'messages.loginFailed': 'Login failed'
      }
      return map[key] || key
    }
  })
}))

describe('LoginForm integration', () => {
  let router: ReturnType<typeof createRouter>

  let pinia: ReturnType<typeof createPinia>

  beforeEach(() => {
    pinia = createPinia()
    setActivePinia(pinia)

    router = createRouter({
      history: createMemoryHistory(),
      routes: [
        { path: '/', name: 'home' },
        { path: '/dashboard', name: 'dashboard' }
      ]
    })

    // clear head
    document.head.innerHTML = ''
  })

  it('applies host styles and redirects to dashboard on successful login', async () => {
    // Provide inline styles from host
    ;(globalThis as any).__CORE_FRONTEND_INLINE_STYLES__ = ':root{--brand-primary:#112233}'
    await applyAppStyles()

    // Assert styles applied
    const styleEl = document.querySelector('style[data-core-frontend="inline-styles"]')
    expect(styleEl).toBeTruthy()
    expect(styleEl!.textContent).toContain('--brand-primary')

    // Import component after Pinia is active so its internal `useAuthStore()` binds to our test Pinia
    const { default: LoginForm } = await import('../LoginForm.vue')

    const wrapper = mount(LoginForm, {
      global: {
        plugins: [pinia, router],
        mocks: {
          $t: (k: string) => k
        }
      }
    })

    // Prepare store spy before submitting
    const { useAuthStore } = await import('../../stores/auth')
    const authStore = useAuthStore()
    const loginSpy = vi.spyOn(authStore, 'login').mockResolvedValue()

    // Spy on router.push
    const pushSpy = vi.spyOn(router, 'push')


    // Fill form
    const emailInput = wrapper.find('[data-testid="login-email"]')
    const passwordInput = wrapper.find('[data-testid="login-password"]')
    const submitButton = wrapper.find('[data-testid="login-submit"]')

    expect(submitButton.exists()).toBe(true)
    expect(submitButton.attributes('disabled')).toBeUndefined()

    await emailInput.setValue('test@example.com')
    await passwordInput.setValue('password123')
    await wrapper.vm.$nextTick()

    // Submit form by triggering submit event
    await wrapper.find('form').trigger('submit')

    // wait for promise microtasks
    await new Promise(resolve => setTimeout(resolve, 0))

    expect(loginSpy).toHaveBeenCalled()
    expect(pushSpy).toHaveBeenCalledWith('/dashboard')
  })
})
