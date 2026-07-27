import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createI18n } from 'vue-i18n'
import PasswordField from '../PasswordField.vue'

// Mock the usePasswordRequirements composable
vi.mock('../../composables/usePasswordRequirements', () => ({
  usePasswordRequirements: () => ({
    checkRequirements: (password: string) => ({
      checks: { value: [] },
      strengthScore: { value: 0 },
      strengthColor: { value: 'text-base-content' },
      strengthLabel: { value: '' }
    })
  })
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  messages: {
    en: {
      register: {
        passwordStrength: 'Password Strength'
      }
    }
  }
})

describe('PasswordField', () => {
  it('renders with basic props', () => {
    const wrapper = mount(PasswordField, {
      props: {
        modelValue: '',
        label: 'Password',
        testId: 'test-password'
      },
      global: {
        plugins: [i18n]
      }
    })

    expect(wrapper.find('input[type="password"]').exists()).toBe(true)
    expect(wrapper.find('label .label-text').text()).toBe('Password')
    expect(wrapper.find('input[data-testid="test-password"]').exists()).toBe(true)
  })

  it('emits update:modelValue on input', async () => {
    const wrapper = mount(PasswordField, {
      props: {
        modelValue: '',
        label: 'Password'
      },
      global: {
        plugins: [i18n]
      }
    })

    const input = wrapper.find('input[type="password"]')
    await input.setValue('test123')

    expect(wrapper.emitted('update:modelValue')).toBeTruthy()
    expect(wrapper.emitted('update:modelValue')?.[0]).toEqual(['test123'])
  })

  it('shows error message when error prop is provided', () => {
    const wrapper = mount(PasswordField, {
      props: {
        modelValue: '',
        label: 'Password',
        error: 'Password is required'
      },
      global: {
        plugins: [i18n]
      }
    })

    expect(wrapper.find('.text-error').text()).toBe('Password is required')
    expect(wrapper.find('input').classes()).toContain('input-error')
  })

  it('hides requirements when showRequirements is false', () => {
    const wrapper = mount(PasswordField, {
      props: {
        modelValue: 'test123',
        label: 'Password',
        showRequirements: false
      },
      global: {
        plugins: [i18n]
      }
    })

    expect(wrapper.find('.space-y-2').exists()).toBe(false)
  })

  it('shows requirements when showRequirements is true and has value', () => {
    const wrapper = mount(PasswordField, {
      props: {
        modelValue: 'test123',
        label: 'Password',
        showRequirements: true
      },
      global: {
        plugins: [i18n]
      }
    })

    expect(wrapper.find('.space-y-2').exists()).toBe(true)
  })
})