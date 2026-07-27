import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createI18n } from 'vue-i18n'
import { createPinia } from 'pinia'
import ChangePasswordForm from '../ChangePasswordForm.vue'

// Mock the auth store
vi.mock('../../stores/auth', () => ({
  useAuthStore: () => ({
    changePassword: vi.fn()
  })
}))

// Mock the router
vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: vi.fn()
  })
}))

// Mock the PasswordField component
vi.mock('../PasswordField.vue', () => ({
  default: {
    name: 'PasswordField',
    props: ['modelValue', 'label', 'placeholder', 'autocomplete', 'error', 'testId', 'required', 'showRequirements'],
    emits: ['update:modelValue', 'blur', 'requirements-changed'],
    template: `
      <div>
        <input 
          :value="modelValue" 
          @input="$emit('update:modelValue', $event.target.value)"
          :data-testid="testId"
        />
      </div>
    `
  }
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  messages: {
    en: {
      change: {
        title: 'Change Password',
        subtitle: 'Update your account password',
        currentPassword: 'Current Password',
        currentPasswordPlaceholder: 'Enter current password',
        newPassword: 'New Password',
        newPasswordPlaceholder: 'Enter new password',
        confirmPassword: 'Confirm Password',
        confirmPasswordPlaceholder: 'Confirm new password',
        confirmPasswordDescription: 'Please enter your new password again',
        changeButton: 'Change Password',
        changing: 'Changing...'
      }
    }
  }
})

const pinia = createPinia()

describe('ChangePasswordForm', () => {
  it('renders all password fields using PasswordField component', () => {
    const wrapper = mount(ChangePasswordForm, {
      global: {
        plugins: [i18n, pinia]
      }
    })

    // Should have three password fields
    const passwordFields = wrapper.findAllComponents({ name: 'PasswordField' })
    expect(passwordFields).toHaveLength(3)

    // Check current password field
    const currentPasswordField = wrapper.find('[data-testid="change-current-password"]')
    expect(currentPasswordField.exists()).toBe(true)

    // Check new password field
    const newPasswordField = wrapper.find('[data-testid="change-new-password"]')
    expect(newPasswordField.exists()).toBe(true)

    // Check confirm password field
    const confirmPasswordField = wrapper.find('[data-testid="change-confirm-password"]')
    expect(confirmPasswordField.exists()).toBe(true)
  })

  it('shows form title and subtitle', () => {
    const wrapper = mount(ChangePasswordForm, {
      global: {
        plugins: [i18n, pinia]
      }
    })

    expect(wrapper.text()).toContain('Change Password')
    expect(wrapper.text()).toContain('Update your account password')
  })

  it('has submit button', () => {
    const wrapper = mount(ChangePasswordForm, {
      global: {
        plugins: [i18n, pinia]
      }
    })

    const submitButton = wrapper.find('[data-testid="change-submit"]')
    expect(submitButton.exists()).toBe(true)
    expect(submitButton.text()).toBe('Change Password')
  })

  it('shows fieldset legends', () => {
    const wrapper = mount(ChangePasswordForm, {
      global: {
        plugins: [i18n, pinia]
      }
    })

    expect(wrapper.text()).toContain('Current Authentication')
    expect(wrapper.text()).toContain('New Password Setup')
  })
})