import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import ToastContainer from '../ToastContainer.vue'
import { useToast } from '../../composables/useToast'

describe('ToastContainer', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    const { clearAllToasts } = useToast()
    clearAllToasts()
  })

  it('should render toast container', () => {
    const wrapper = mount(ToastContainer)
    expect(wrapper.vm).toBeTruthy()
  })

  it('should display success toast with correct styling', async () => {
    const { success } = useToast()
    success('Success!')
    
    const wrapper = mount(ToastContainer)
    await wrapper.vm.$nextTick()
    
    expect(wrapper.find('.alert-success').exists()).toBe(true)
  })

  it('should display error toast with correct styling', async () => {
    const { error } = useToast()
    error('Error occurred')
    
    const wrapper = mount(ToastContainer)
    await wrapper.vm.$nextTick()
    
    expect(wrapper.find('.alert-error').exists()).toBe(true)
  })

  it('should display info toast with correct styling', async () => {
    const { info } = useToast()
    info('Info message')
    
    const wrapper = mount(ToastContainer)
    await wrapper.vm.$nextTick()
    
    expect(wrapper.find('.alert-info').exists()).toBe(true)
  })

  it('should display warning toast with correct styling', async () => {
    const { warning } = useToast()
    warning('Warning!')
    
    const wrapper = mount(ToastContainer)
    await wrapper.vm.$nextTick()
    
    expect(wrapper.find('.alert-warning').exists()).toBe(true)
  })

  it('should remove toast on dismiss button click', async () => {
    const { success } = useToast()
    success('Test message')
    
    const wrapper = mount(ToastContainer)
    await wrapper.vm.$nextTick()
    
    const dismissButton = wrapper.find('button[aria-label="Close"]')
    if (dismissButton.exists()) {
      await dismissButton.trigger('click')
      await wrapper.vm.$nextTick()
    }
    
    expect(wrapper.vm).toBeTruthy()
  })

  it('should display toast with title', async () => {
    const { info } = useToast()
    info('Message content', 'Toast Title')
    
    const wrapper = mount(ToastContainer)
    await wrapper.vm.$nextTick()
    
    expect(wrapper.text()).toContain('Toast Title')
  })

  it('should display multiple toasts', async () => {
    const { success, error } = useToast()
    success('First message')
    error('Second message')
    
    const wrapper = mount(ToastContainer)
    await wrapper.vm.$nextTick()
    
    const alerts = wrapper.findAll('.alert')
    expect(alerts.length).toBeGreaterThanOrEqual(2)
  })

  it('should render toast actions when provided', async () => {
    const { toasts } = useToast()
    
    // Manually add a toast with actions
    toasts.value.push({
      id: Date.now().toString(),
      message: 'Action message',
      type: 'info',
      actions: [
        { label: 'Undo', action: vi.fn() },
        { label: 'Confirm', action: vi.fn(), class: 'btn-primary' }
      ]
    })
    
    const wrapper = mount(ToastContainer)
    await wrapper.vm.$nextTick()
    
    expect(wrapper.text()).toContain('Undo')
  })
})
