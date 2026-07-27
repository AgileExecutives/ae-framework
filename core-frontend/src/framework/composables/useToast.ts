import { ref } from 'vue'

export interface ToastMessage {
  id: string
  type: 'success' | 'error' | 'warning' | 'info'
  title?: string
  message: string
  duration?: number
  actions?: Array<{
    label: string
    action: () => void
    class?: string
  }>
}

const toasts = ref<ToastMessage[]>([])

export function useToast() {
  const addToast = (toast: Omit<ToastMessage, 'id'>) => {
    const id = `toast-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
    const newToast: ToastMessage = {
      id,
      duration: 5000,
      ...toast
    }
    
    toasts.value.push(newToast)
    
    if (newToast.duration && newToast.duration > 0) {
      setTimeout(() => {
        removeToast(id)
      }, newToast.duration)
    }
    
    return id
  }

  const removeToast = (id: string) => {
    const index = toasts.value.findIndex(toast => toast.id === id)
    if (index > -1) {
      toasts.value.splice(index, 1)
    }
  }

  const clearAllToasts = () => {
    toasts.value = []
  }

  // Convenience methods
  const success = (message: string, title?: string, options?: Partial<ToastMessage>) => {
    return addToast({ type: 'success', message, title, ...options })
  }

  const error = (message: string, title?: string, options?: Partial<ToastMessage>) => {
    return addToast({ type: 'error', message, title, duration: 8000, ...options })
  }

  const warning = (message: string, title?: string, options?: Partial<ToastMessage>) => {
    return addToast({ type: 'warning', message, title, ...options })
  }

  const info = (message: string, title?: string, options?: Partial<ToastMessage>) => {
    return addToast({ type: 'info', message, title, ...options })
  }

  return {
    toasts: toasts,
    addToast,
    removeToast,
    clearAllToasts,
    success,
    error,
    warning,
    info
  }
}