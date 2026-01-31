import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import ChangePasswordForm from '../ChangePasswordForm.vue'
import { useAuthStore } from '../../stores/auth'

// Mock i18n
vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string) => {
      const translations: Record<string, string> = {
        'validation.currentPasswordRequired': 'Current password is required',
        'validation.passwordRequired': 'Password is required',
        'validation.passwordRepeatRequired': 'Please repeat your password',
        'validation.passwordsDontMatch': 'Passwords do not match',
        'validation.passwordSameAsCurrent': 'New password must be different from current password',
        'change.title': 'Change Password',
        'change.subtitle': 'Update your password',
        'change.currentPassword': 'Current Password',
        'change.newPassword': 'New Password',
        'change.confirmPassword': 'Confirm New Password',
        'change.changeButton': 'Change Password',
        'change.changing': 'Changing...',
        'change.successMessage': 'Password changed successfully',
        'change.errorMessage': 'Failed to change password'
      }
      return translations[key] || key
    }
  })
}))

// Mock router
vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: vi.fn()
  })
}))

describe('ChangePasswordForm - Component Tests', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  describe('Validation - Current Password', () => {
    it('should show error when current password is empty', async () => {
      const wrapper = mount(ChangePasswordForm)
      
      const submitButton = wrapper.find('[data-testid="change-submit"]')
      await submitButton.trigger('click')
      await wrapper.vm.$nextTick()
      
      expect(wrapper.text()).toContain('Current password is required')
    })
  })

  describe('Validation - New Password', () => {
    it('should show error when new password is empty', async () => {
      const wrapper = mount(ChangePasswordForm)
      
      const currentPassword = wrapper.find('[data-testid="change-current-password"]')
      await currentPassword.setValue('oldpass123')
      
      const submitButton = wrapper.find('[data-testid="change-submit"]')
      await submitButton.trigger('click')
      await wrapper.vm.$nextTick()
      
      expect(wrapper.text()).toContain('Password is required')
    })

    it('should show error when new password is same as current password', async () => {
      const wrapper = mount(ChangePasswordForm)
      
      const currentPassword = wrapper.find('[data-testid="change-current-password"]')
      const newPassword = wrapper.find('[data-testid="change-new-password"]')
      const submitButton = wrapper.find('[data-testid="change-submit"]')
      
      await currentPassword.setValue('samepass123')
      await newPassword.setValue('samepass123')
      await submitButton.trigger('click')
      await wrapper.vm.$nextTick()
      
      expect(wrapper.text()).toContain('New password must be different from current password')
    })

    it('should require password to meet minimum requirements', async () => {
      const wrapper = mount(ChangePasswordForm)
      
      const currentPassword = wrapper.find('[data-testid="change-current-password"]')
      const newPassword = wrapper.find('[data-testid="change-new-password"]')
      const submitButton = wrapper.find('[data-testid="change-submit"]')
      
      await currentPassword.setValue('oldpass123')
      await newPassword.setValue('123') // Too short
      await submitButton.trigger('click')
      await wrapper.vm.$nextTick()
      
      // Should show password requirements error
      const errorText = wrapper.text()
      expect(errorText).toMatch(/password|requirements|characters/i)
    })
  })

  describe('Validation - Confirm Password', () => {
    it('should show error when confirm password is empty', async () => {
      const wrapper = mount(ChangePasswordForm)
      
      const currentPassword = wrapper.find('[data-testid="change-current-password"]')
      const newPassword = wrapper.find('[data-testid="change-new-password"]')
      const submitButton = wrapper.find('[data-testid="change-submit"]')
      
      await currentPassword.setValue('oldpass123')
      await newPassword.setValue('newpass123')
      await submitButton.trigger('click')
      await wrapper.vm.$nextTick()
      
      expect(wrapper.text()).toContain('Please repeat your password')
    })

    it('should show error when passwords do not match', async () => {
      const wrapper = mount(ChangePasswordForm)
      
      const currentPassword = wrapper.find('[data-testid="change-current-password"]')
      const newPassword = wrapper.find('[data-testid="change-new-password"]')
      const confirmPassword = wrapper.find('[data-testid="change-confirm-password"]')
      const submitButton = wrapper.find('[data-testid="change-submit"]')
      
      await currentPassword.setValue('oldpass123')
      await newPassword.setValue('newpass123')
      await confirmPassword.setValue('different123')
      await submitButton.trigger('click')
      await wrapper.vm.$nextTick()
      
      expect(wrapper.text()).toContain('Passwords do not match')
    })

    it('should clear password mismatch error when passwords match', async () => {
      const wrapper = mount(ChangePasswordForm)
      
      const currentPassword = wrapper.find('[data-testid="change-current-password"]')
      const newPassword = wrapper.find('[data-testid="change-new-password"]')
      const confirmPassword = wrapper.find('[data-testid="change-confirm-password"]')
      const submitButton = wrapper.find('[data-testid="change-submit"]')
      
      // Set mismatched passwords
      await currentPassword.setValue('oldpass123')
      await newPassword.setValue('newpass123')
      await confirmPassword.setValue('different123')
      await submitButton.trigger('click')
      await wrapper.vm.$nextTick()
      
      expect(wrapper.text()).toContain('Passwords do not match')
      
      // Fix the mismatch
      await confirmPassword.setValue('newpass123')
      await confirmPassword.trigger('blur')
      await wrapper.vm.$nextTick()
      
      expect(wrapper.text()).not.toContain('Passwords do not match')
    })
  })

  describe('Form Submission', () => {
    it('should not submit when form is invalid', async () => {
      const wrapper = mount(ChangePasswordForm)
      const authStore = useAuthStore()
      
      const changePasswordSpy = vi.spyOn(authStore, 'changePassword').mockResolvedValue()
      
      const submitButton = wrapper.find('[data-testid="change-submit"]')
      await submitButton.trigger('click')
      await wrapper.vm.$nextTick()
      
      expect(changePasswordSpy).not.toHaveBeenCalled()
    })

    it('should call changePassword with correct data when form is valid', async () => {
      const wrapper = mount(ChangePasswordForm)
      const authStore = useAuthStore()
      
      const changePasswordSpy = vi.spyOn(authStore, 'changePassword').mockResolvedValue()
      
      const currentPassword = wrapper.find('[data-testid="change-current-password"]')
      const newPassword = wrapper.find('[data-testid="change-new-password"]')
      const confirmPassword = wrapper.find('[data-testid="change-confirm-password"]')
      const submitButton = wrapper.find('[data-testid="change-submit"]')
      
      await currentPassword.setValue('oldpass123')
      await newPassword.setValue('newpass123')
      await confirmPassword.setValue('newpass123')
      await submitButton.trigger('click')
      await wrapper.vm.$nextTick()
      
      expect(changePasswordSpy).toHaveBeenCalledWith({
        current_password: 'oldpass123',
        new_password: 'newpass123'
      })
    })

    it('should display success message after successful password change', async () => {
      const wrapper = mount(ChangePasswordForm)
      const authStore = useAuthStore()
      
      vi.spyOn(authStore, 'changePassword').mockResolvedValue()
      
      const currentPassword = wrapper.find('[data-testid="change-current-password"]')
      const newPassword = wrapper.find('[data-testid="change-new-password"]')
      const confirmPassword = wrapper.find('[data-testid="change-confirm-password"]')
      const submitButton = wrapper.find('[data-testid="change-submit"]')
      
      await currentPassword.setValue('oldpass123')
      await newPassword.setValue('newpass123')
      await confirmPassword.setValue('newpass123')
      await submitButton.trigger('click')
      
      // Wait for promise
      await new Promise(resolve => setTimeout(resolve, 0))
      await wrapper.vm.$nextTick()
      
      expect(wrapper.find('[data-testid="change-success-message"]').exists()).toBe(true)
      expect(wrapper.text()).toContain('Password changed successfully')
    })

    it('should display error message on password change failure', async () => {
      const wrapper = mount(ChangePasswordForm)
      const authStore = useAuthStore()
      
      vi.spyOn(authStore, 'changePassword').mockRejectedValue(new Error('Invalid current password'))
      
      const currentPassword = wrapper.find('[data-testid="change-current-password"]')
      const newPassword = wrapper.find('[data-testid="change-new-password"]')
      const confirmPassword = wrapper.find('[data-testid="change-confirm-password"]')
      const submitButton = wrapper.find('[data-testid="change-submit"]')
      
      await currentPassword.setValue('wrongpass123')
      await newPassword.setValue('newpass123')
      await confirmPassword.setValue('newpass123')
      await submitButton.trigger('click')
      
      // Wait for promise
      await new Promise(resolve => setTimeout(resolve, 0))
      await wrapper.vm.$nextTick()
      
      expect(wrapper.find('[data-testid="change-error-message"]').exists()).toBe(true)
      expect(wrapper.text()).toContain('Invalid current password')
    })

    it('should clear form after successful password change', async () => {
      const wrapper = mount(ChangePasswordForm)
      const authStore = useAuthStore()
      
      vi.spyOn(authStore, 'changePassword').mockResolvedValue()
      
      const currentPassword = wrapper.find('[data-testid="change-current-password"]')
      const newPassword = wrapper.find('[data-testid="change-new-password"]')
      const confirmPassword = wrapper.find('[data-testid="change-confirm-password"]')
      const submitButton = wrapper.find('[data-testid="change-submit"]')
      
      await currentPassword.setValue('oldpass123')
      await newPassword.setValue('newpass123')
      await confirmPassword.setValue('newpass123')
      await submitButton.trigger('click')
      
      // Wait for promise
      await new Promise(resolve => setTimeout(resolve, 0))
      await wrapper.vm.$nextTick()
      
      // Form should be cleared
      expect((currentPassword.element as HTMLInputElement).value).toBe('')
      expect((newPassword.element as HTMLInputElement).value).toBe('')
      expect((confirmPassword.element as HTMLInputElement).value).toBe('')
    })

    it('should disable submit button while submitting', async () => {
      const wrapper = mount(ChangePasswordForm)
      const authStore = useAuthStore()
      
      vi.spyOn(authStore, 'changePassword').mockImplementation(() =>
        new Promise(resolve => setTimeout(resolve, 100))
      )
      
      const currentPassword = wrapper.find('[data-testid="change-current-password"]')
      const newPassword = wrapper.find('[data-testid="change-new-password"]')
      const confirmPassword = wrapper.find('[data-testid="change-confirm-password"]')
      const submitButton = wrapper.find('[data-testid="change-submit"]')
      
      await currentPassword.setValue('oldpass123')
      await newPassword.setValue('newpass123')
      await confirmPassword.setValue('newpass123')
      await submitButton.trigger('click')
      await wrapper.vm.$nextTick()
      
      expect(submitButton.attributes('disabled')).toBeDefined()
      expect(wrapper.text()).toContain('Changing...')
    })
  })

  describe('Validation Behavior', () => {
    it('should not show validation errors before first submit attempt', async () => {
      const wrapper = mount(ChangePasswordForm)
      
      const currentPassword = wrapper.find('[data-testid="change-current-password"]')
      const newPassword = wrapper.find('[data-testid="change-new-password"]')
      
      await currentPassword.setValue('test')
      await currentPassword.trigger('blur')
      await newPassword.setValue('123')
      await newPassword.trigger('blur')
      await wrapper.vm.$nextTick()
      
      // Should not show errors before submit
      expect(wrapper.text()).not.toContain('required')
      expect(wrapper.text()).not.toContain('characters')
    })

    it('should show validation errors on blur after submit attempt', async () => {
      const wrapper = mount(ChangePasswordForm)
      
      const submitButton = wrapper.find('[data-testid="change-submit"]')
      const newPassword = wrapper.find('[data-testid="change-new-password"]')
      
      // Trigger submit with empty form
      await submitButton.trigger('click')
      await wrapper.vm.$nextTick()
      
      // Now errors should show
      expect(wrapper.text()).toContain('required')
      
      // Change field and blur
      await newPassword.setValue('123')
      await newPassword.trigger('blur')
      await wrapper.vm.$nextTick()
      
      // Should show validation error for short password
      const errorText = wrapper.text()
      expect(errorText).toMatch(/password|requirements|characters/i)
    })
  })
})
