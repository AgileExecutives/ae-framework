<script setup lang="ts">
// Reset Password Form - handles password reset with token validation
import SingleFormCard from "./SingleFormCard.vue"
import { useAuthStore } from "../stores/auth"
import { useRouter } from "vue-router"
import { ref, reactive } from "vue"
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const authStore = useAuthStore()
const router = useRouter()

const token = router.currentRoute.value.params.token as string
const successMessage = ref<string | null>(null)
const errorMessage = ref<string | null>(null)
const isSubmitting = ref(false)
const formErrors = reactive<Record<string, string>>({})
const hasAttemptedSubmit = ref(false)

// Form data
const formData = reactive({
  newPassword: '',
  confirmPassword: ''
})

// Validation functions
const validateField = (field: string, value: string) => {
  switch (field) {
    case 'newPassword':
      if (!value) return t('validation.passwordRequired')
      if (value.length < 8) return t('validation.passwordMin')
      break
    case 'confirmPassword':
      if (!value) return t('validation.passwordRepeatRequired')
      if (value !== formData.newPassword) return t('validation.passwordsDontMatch')
      break
  }
  return ''
}

const onFieldBlur = (field: string, value: string) => {
  // Only show validation errors after first submit attempt
  if (!hasAttemptedSubmit.value) return
  
  const error = validateField(field, value)
  if (error) {
    formErrors[field] = error
  } else {
    delete formErrors[field]
  }
}

const onFieldInput = (field: string, value: string) => {
  // Only show validation errors after first submit attempt
  if (!hasAttemptedSubmit.value) return
  
  const error = validateField(field, value)
  if (error) {
    formErrors[field] = error
  } else {
    delete formErrors[field]
  }
}

const validateForm = () => {
  const errors: Record<string, string> = {}
  
  Object.keys(formData).forEach(field => {
    const error = validateField(field, (formData as any)[field])
    if (error) errors[field] = error
  })
  
  Object.assign(formErrors, errors)
  return Object.keys(errors).length === 0
}

const onSubmit = async (event: Event) => {
  event.preventDefault()
  
  // Mark that user has attempted to submit
  hasAttemptedSubmit.value = true
  
  if (!validateForm()) {
    return
  }

  try {
    isSubmitting.value = true
    successMessage.value = null
    errorMessage.value = null
    
    await authStore.resetPassword(token, formData.newPassword)
    
    successMessage.value = t('reset.successMessage')
    
    // Reset form
    Object.assign(formData, {
      newPassword: '',
      confirmPassword: ''
    })
    Object.keys(formErrors).forEach(key => delete formErrors[key])
    
    // Redirect to login after 2 seconds
    setTimeout(() => {
      router.push('/login')
    }, 2000)
  } catch (error: any) {
    errorMessage.value = error.message || t('reset.errorMessage')
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <SingleFormCard :title="$t('reset.title')" :subtitle="$t('reset.subtitle')">
    <div 
      v-if="successMessage" 
      class="alert alert-success mb-4"
      data-testid="reset-success-message"
    >
      <span class="break-words">{{ successMessage }}</span>
    </div>
    
    <div 
      v-if="errorMessage" 
      class="alert alert-error mb-4"
      data-testid="reset-error-message"
    >
      <span class="break-words">{{ errorMessage }}</span>
    </div>
    
    <form @submit="onSubmit" autocomplete="off" class="space-y-6">
      <!-- New Password Setup -->
      <fieldset class="space-y-4">
        <legend class="text-base font-semibold text-base-content mb-4 break-words text-wrap">New Password Setup</legend>
        
        <div class="form-control">
          <label class="label">
            <span class="label-text">{{ $t('reset.newPassword') }}</span>
          </label>
          <input 
            type="password" 
            autocomplete="new-password" 
            placeholder="••••••••"
            v-model="formData.newPassword"
            @blur="onFieldBlur('newPassword', formData.newPassword)"
            @input="onFieldInput('confirmPassword', formData.confirmPassword)"
            :class="['input input-bordered w-full', formErrors.newPassword ? 'input-error' : '']"
            data-testid="reset-new-password"
            required
            minlength="8"
          />
          <label class="label">
            <span class="label-text-alt break-words">{{ $t('reset.passwordRequirements') }}</span>
          </label>
          <label v-if="formErrors.newPassword" class="label">
            <span class="label-text-alt text-error break-words">{{ formErrors.newPassword }}</span>
          </label>
        </div>
        
        <div class="form-control">
          <label class="label">
            <span class="label-text">{{ $t('reset.confirmPassword') }}</span>
          </label>
          <input 
            type="password" 
            autocomplete="new-password" 
            placeholder="••••••••" 
            v-model="formData.confirmPassword"
            @blur="onFieldBlur('confirmPassword', formData.confirmPassword)"
            :class="['input input-bordered w-full', formErrors.confirmPassword ? 'input-error' : '']"
            data-testid="reset-confirm-password"
            required
          />
          <label v-if="formErrors.confirmPassword" class="label">
            <span class="label-text-alt text-error break-words">{{ formErrors.confirmPassword }}</span>
          </label>
        </div>
      </fieldset>
      
      <!-- Submit Button -->
      <button 
        type="submit" 
        class="btn btn-primary w-full"
        :disabled="isSubmitting"
        data-testid="reset-submit"
      >
        <span v-if="isSubmitting" class="loading loading-spinner loading-sm"></span>
        <span v-if="isSubmitting">{{ $t('reset.resetting') }}</span>
        <span v-else>{{ $t('reset.resetButton') }}</span>
      </button>
    </form>
  </SingleFormCard>
</template>
