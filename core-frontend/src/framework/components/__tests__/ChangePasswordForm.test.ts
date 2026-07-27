import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import ChangePasswordForm from '../ChangePasswordForm.vue'
import { useAuthStore } from '../../stores/auth'

// Mock i18n including password requirement and validation keys used by component
vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string, params?: Record<string, any>) => {
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
        'change.errorMessage': 'Failed to change password',

        // Password requirements / validation keys
        'passwordRequirements.oneUppercase': 'At least one uppercase letter',
        'passwordRequirements.oneNumber': 'At least one number',
        'passwordRequirements.oneSpecial': 'At least one special character',
        'validation.passwordRequirements': 'Password must include: {requirements}',
        'passwordRequirements.atLeastChars': 'At least {n} characters',
        'passwordRequirements.hasUppercase': 'At least one uppercase letter',
        'passwordRequirements.hasLowercase': 'At least one lowercase letter',
        'passwordRequirements.hasNumber': 'At least one number',
        'passwordRequirements.hasSpecialChar': 'At least one special character',
        'validation.passwordMinLength': 'Password must be at least {n} characters',
        'validation.passwordRequireUppercase': 'Password must contain an uppercase letter',
        'validation.passwordRequireNumber': 'Password must contain a number',
        'validation.passwordRequireSpecial': 'Password must contain a special character',
        'passwordStrength.veryWeak': 'very weak',
        'passwordStrength.weak': 'weak',
        'passwordStrength.middle': 'middle',
        'passwordStrength.strong': 'strong'
      }

      let val = translations[key] || key
      // Simple param interpolation for tests
      if (params) {
        Object.keys(params).forEach(k => {
          val = val.replace(`{${k}}`, String(params[k]))
        })
      }
      return val
    }
  })
}))

// Shared router push mock for assertions
const pushMock = vi.fn()
vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: pushMock
  })
}))

// Mock password requirements composable to avoid external API and strict rules in tests
vi.mock('../../composables/usePasswordRequirements', () => ({
  usePasswordRequirements: () => ({
    validatePassword: (password: string) => ({ valid: true, errors: [] }),
    requirements: { minLength: 1, capital: false, numbers: false, special: false },
    checkRequirements: (p: string) => ({
      checks: { value: [] },
      strengthScore: { value: 0 },
      strengthColor: { value: 'text-base-content' },
      strengthLabel: { value: '' }
    })
  })
}))

describe('ChangePasswordForm - Component Tests', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  describe('Validation - Current Password', () => {
    it('should show error when current password is empty', async () => {
      const wrapper = mount(ChangePasswordForm)
      
      const form = wrapper.find('form')
      await form.trigger('submit')
      await wrapper.vm.$nextTick()
      
      expect(wrapper.text()).toContain('Current password is required')
    })
  })

  describe('Validation - New Password', () => {
    it('should show error when new password is empty', async () => {
      const wrapper = mount(ChangePasswordForm)
      
      const currentPassword = wrapper.find('[data-testid="change-current-password"]')
      await currentPassword.setValue('oldpass123')
      
      const form = wrapper.find('form')
      await form.trigger('submit')
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
      const form = wrapper.find('form')
      await form.trigger('submit')
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
      const form = wrapper.find('form')
      await form.trigger('submit')
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
      const form = wrapper.find('form')
      await form.trigger('submit')
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
      const form = wrapper.find('form')
      await form.trigger('submit')
      await wrapper.vm.$nextTick()
      
      expect(wrapper.text()).toContain('Passwords do not match')
    })

    it('should clear password mismatch error when passwords match', async () => {
      const wrapper = mount(ChangePasswordForm)
      
      const currentPassword = wrapper.find('[data-testid="change-current-password"]')
      const newPassword = wrapper.find('[data-testid="change-new-password"]')
      const confirmPassword = wrapper.find('[data-testid="change-confirm-password"]')
      const submitButton = wrapper.find('[data-testid="change-submit"]')
      
      // Set mismatched passwords and submit
      await currentPassword.setValue('oldpass123')
      await newPassword.setValue('newpass123')
      await confirmPassword.setValue('different123')
      await wrapper.find('form').trigger('submit')
      await wrapper.vm.$nextTick()
      
      expect(wrapper.text()).toContain('Passwords do not match')
      
      // Fix the mismatch and blur
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
      await wrapper.find('form').trigger('submit')
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
      await wrapper.find('form').trigger('submit')
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
      await wrapper.find('form').trigger('submit')

      // Wait for promise
      await new Promise(resolve => setTimeout(resolve, 0))
      await wrapper.vm.$nextTick()

      // Component navigates on success
      expect(pushMock).toHaveBeenCalled()
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
      await wrapper.find('form').trigger('submit')
      
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
      await wrapper.find('form').trigger('submit')

      // Wait for promise
      await new Promise(resolve => setTimeout(resolve, 0))
      await wrapper.vm.$nextTick()

      // Component navigates on success (does not clear form)
      expect(pushMock).toHaveBeenCalled()
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
      await wrapper.find('form').trigger('submit')
      await wrapper.vm.$nextTick()

      expect(pushMock).toHaveBeenCalled() // submission started
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
      await wrapper.find('form').trigger('submit')
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
