import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import ForgotPasswordForm from '../ForgotPasswordForm.vue'
import { useAuthStore } from '../../stores/auth'

// Mock i18n
vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string) => {
      const translations: Record<string, string> = {
        'validation.emailRequired': 'Email is required',
        'validation.emailInvalid': 'Invalid email format',
        'forgot.title': 'Reset Password',
        'forgot.subtitle': 'Enter your email',
        'forgot.email': 'Email Address',
        'forgot.emailDescription': 'We will send you a reset link',
        'forgot.sendLink': 'Send Reset Link',
        'forgot.sending': 'Sending...',
        'forgot.successMessage': 'Password reset link sent to your email',
        'forgot.errorMessage': 'Failed to send reset link'
      }
      return translations[key] || key
    }
  })
}))

describe('ForgotPasswordForm - Email Validation', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('should display email validation error for empty email', async () => {
    const wrapper = mount(ForgotPasswordForm)
    
    const emailInput = wrapper.find('[data-testid="forgot-email"]')
    const submitButton = wrapper.find('[data-testid="forgot-submit"]')
    
    // Try to submit without entering email
    await submitButton.trigger('click')
    await wrapper.vm.$nextTick()
    
    // Should show validation error
    expect(wrapper.text()).toContain('Email is required')
  })

  it('should display email validation error for invalid email format', async () => {
    const wrapper = mount(ForgotPasswordForm)
    
    const emailInput = wrapper.find('[data-testid="forgot-email"]')
    const submitButton = wrapper.find('[data-testid="forgot-submit"]')
    
    // Enter invalid email
    await emailInput.setValue('invalid-email')
    await submitButton.trigger('click')
    await wrapper.vm.$nextTick()
    
    // Should show validation error
    expect(wrapper.text()).toContain('Invalid email format')
  })

  it('should display email validation error for email without @', async () => {
    const wrapper = mount(ForgotPasswordForm)
    
    const emailInput = wrapper.find('[data-testid="forgot-email"]')
    const submitButton = wrapper.find('[data-testid="forgot-submit"]')
    
    // Enter email without @
    await emailInput.setValue('invalidemail.com')
    await submitButton.trigger('click')
    await wrapper.vm.$nextTick()
    
    // Should show validation error
    expect(wrapper.text()).toContain('Invalid email format')
  })

  it('should display email validation error for email without domain', async () => {
    const wrapper = mount(ForgotPasswordForm)
    
    const emailInput = wrapper.find('[data-testid="forgot-email"]')
    const submitButton = wrapper.find('[data-testid="forgot-submit"]')
    
    // Enter email without domain
    await emailInput.setValue('user@')
    await submitButton.trigger('click')
    await wrapper.vm.$nextTick()
    
    // Should show validation error
    expect(wrapper.text()).toContain('Invalid email format')
  })

  it('should not show validation errors before submit attempt', async () => {
    const wrapper = mount(ForgotPasswordForm)
    
    const emailInput = wrapper.find('[data-testid="forgot-email"]')
    
    // Enter invalid email
    await emailInput.setValue('invalid-email')
    await emailInput.trigger('blur')
    await wrapper.vm.$nextTick()
    
    // Should NOT show error before submit attempt
    expect(wrapper.text()).not.toContain('Invalid email format')
  })

  it('should show validation errors on blur after submit attempt', async () => {
    const wrapper = mount(ForgotPasswordForm)
    
    const emailInput = wrapper.find('[data-testid="forgot-email"]')
    const submitButton = wrapper.find('[data-testid="forgot-submit"]')
    
    // Try to submit with invalid email
    await emailInput.setValue('invalid')
    await submitButton.trigger('click')
    await wrapper.vm.$nextTick()
    
    // Should show error
    expect(wrapper.text()).toContain('Invalid email format')
    
    // Change to valid email
    await emailInput.setValue('valid@example.com')
    await emailInput.trigger('blur')
    await wrapper.vm.$nextTick()
    
    // Error should clear
    expect(wrapper.text()).not.toContain('Invalid email format')
  })

  it('should accept valid email formats', async () => {
    const wrapper = mount(ForgotPasswordForm)
    const authStore = useAuthStore()
    
    // Mock the forgotPassword method
    vi.spyOn(authStore, 'forgotPassword').mockResolvedValue()
    
    const emailInput = wrapper.find('[data-testid="forgot-email"]')
    const submitButton = wrapper.find('[data-testid="forgot-submit"]')
    
    const validEmails = [
      'user@example.com',
      'test.user@example.co.uk',
      'user+tag@example.com',
      'user_name@example-domain.com'
    ]
    
    for (const email of validEmails) {
      await emailInput.setValue(email)
      await submitButton.trigger('click')
      await wrapper.vm.$nextTick()
      
      // Should not show validation error
      expect(wrapper.text()).not.toContain('Invalid email format')
      expect(wrapper.text()).not.toContain('Email is required')
    }
  })

  it('should call forgotPassword with valid email', async () => {
    const wrapper = mount(ForgotPasswordForm)
    const authStore = useAuthStore()
    
    // Mock the forgotPassword method
    const forgotPasswordSpy = vi.spyOn(authStore, 'forgotPassword').mockResolvedValue()
    
    const emailInput = wrapper.find('[data-testid="forgot-email"]')
    const submitButton = wrapper.find('[data-testid="forgot-submit"]')
    
    // Enter valid email and submit
    await emailInput.setValue('test@example.com')
    await submitButton.trigger('click')
    await wrapper.vm.$nextTick()
    
    // Should call forgotPassword
    expect(forgotPasswordSpy).toHaveBeenCalledWith('test@example.com')
  })

  it('should not call forgotPassword with invalid email', async () => {
    const wrapper = mount(ForgotPasswordForm)
    const authStore = useAuthStore()
    
    // Mock the forgotPassword method
    const forgotPasswordSpy = vi.spyOn(authStore, 'forgotPassword').mockResolvedValue()
    
    const emailInput = wrapper.find('[data-testid="forgot-email"]')
    const submitButton = wrapper.find('[data-testid="forgot-submit"]')
    
    // Enter invalid email and submit
    await emailInput.setValue('invalid-email')
    await submitButton.trigger('click')
    await wrapper.vm.$nextTick()
    
    // Should NOT call forgotPassword
    expect(forgotPasswordSpy).not.toHaveBeenCalled()
  })

  it('should display success message after successful submission', async () => {
    const wrapper = mount(ForgotPasswordForm)
    const authStore = useAuthStore()
    
    // Mock successful response
    vi.spyOn(authStore, 'forgotPassword').mockResolvedValue()
    
    const emailInput = wrapper.find('[data-testid="forgot-email"]')
    const submitButton = wrapper.find('[data-testid="forgot-submit"]')
    
    // Submit valid email
    await emailInput.setValue('test@example.com')
    await submitButton.trigger('click')
    
    // Wait for promise to resolve
    await new Promise(resolve => setTimeout(resolve, 0))
    await wrapper.vm.$nextTick()
    
    // Should show success message
    expect(wrapper.find('[data-testid="forgot-success-message"]').exists()).toBe(true)
    expect(wrapper.text()).toContain('Password reset link sent to your email')
  })

  it('should display error message on submission failure', async () => {
    const wrapper = mount(ForgotPasswordForm)
    const authStore = useAuthStore()
    
    // Mock error response
    vi.spyOn(authStore, 'forgotPassword').mockRejectedValue(new Error('Network error'))
    
    const emailInput = wrapper.find('[data-testid="forgot-email"]')
    const submitButton = wrapper.find('[data-testid="forgot-submit"]')
    
    // Submit valid email
    await emailInput.setValue('test@example.com')
    await submitButton.trigger('click')
    
    // Wait for promise to reject
    await new Promise(resolve => setTimeout(resolve, 0))
    await wrapper.vm.$nextTick()
    
    // Should show error message
    expect(wrapper.find('[data-testid="forgot-error-message"]').exists()).toBe(true)
    expect(wrapper.text()).toContain('Network error')
  })

  it('should disable submit button while submitting', async () => {
    const wrapper = mount(ForgotPasswordForm)
    const authStore = useAuthStore()
    
    // Mock slow response
    vi.spyOn(authStore, 'forgotPassword').mockImplementation(() => 
      new Promise(resolve => setTimeout(resolve, 100))
    )
    
    const emailInput = wrapper.find('[data-testid="forgot-email"]')
    const submitButton = wrapper.find('[data-testid="forgot-submit"]')
    
    // Submit valid email
    await emailInput.setValue('test@example.com')
    await submitButton.trigger('click')
    await wrapper.vm.$nextTick()
    
    // Button should be disabled
    expect(submitButton.attributes('disabled')).toBeDefined()
    expect(wrapper.text()).toContain('Sending...')
  })

  it('should clear form after successful submission', async () => {
    const wrapper = mount(ForgotPasswordForm)
    const authStore = useAuthStore()
    
    // Mock successful response
    vi.spyOn(authStore, 'forgotPassword').mockResolvedValue()
    
    const emailInput = wrapper.find('[data-testid="forgot-email"]')
    const submitButton = wrapper.find('[data-testid="forgot-submit"]')
    
    // Submit valid email
    await emailInput.setValue('test@example.com')
    await submitButton.trigger('click')
    
    // Wait for promise to resolve
    await new Promise(resolve => setTimeout(resolve, 0))
    await wrapper.vm.$nextTick()
    
    // Form should be cleared
    expect((emailInput.element as HTMLInputElement).value).toBe('')
  })
})
