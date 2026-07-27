<script setup lang="ts">
import { useToast } from '../composables/useToast'

const { toasts, removeToast } = useToast()

const getAlertClass = (type: string) => {
  switch (type) {
    case 'success':
      return 'alert-success'
    case 'error':
      return 'alert-error'
    case 'warning':
      return 'alert-warning'
    case 'info':
      return 'alert-info'
    default:
      return 'alert-info'
  }
}

const getIcon = (type: string) => {
  switch (type) {
    case 'success':
      return '✓'
    case 'error':
      return '✕'
    case 'warning':
      return '⚠'
    case 'info':
      return 'ℹ'
    default:
      return 'ℹ'
  }
}
</script>

<template>
  <div class="toast toast-top toast-end z-50">
    <div
      v-for="toast in toasts"
      :key="toast.id"
      :class="['alert', getAlertClass(toast.type), 'shadow-lg min-w-80 max-w-md']"
    >
      <div class="flex-1">
        <div class="flex items-start gap-2">
          <span class="text-lg flex-shrink-0">{{ getIcon(toast.type) }}</span>
          <div class="flex-1 min-w-0">
            <div v-if="toast.title" class="font-bold text-sm">{{ toast.title }}</div>
            <div class="text-sm break-words">{{ toast.message }}</div>
          </div>
        </div>
        
        <div v-if="toast.actions && toast.actions.length > 0" class="flex gap-2 mt-2">
          <button
            v-for="action in toast.actions"
            :key="action.label"
            @click="action.action"
            :class="['btn btn-sm', action.class || 'btn-ghost']"
          >
            {{ action.label }}
          </button>
        </div>
      </div>
      
      <button
        @click="removeToast(toast.id)"
        class="btn btn-sm btn-ghost btn-square"
        aria-label="Close"
      >
        ✕
      </button>
    </div>
  </div>
</template>