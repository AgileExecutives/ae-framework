<script setup lang="ts">
import SingleFormCard from "./SingleFormCard.vue"
import { useAuthStore } from "../stores/auth"
import { ref, reactive } from "vue"
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const authStore = useAuthStore()
const successMessage = ref<string | null>(null)
const errorMessage = ref<string | null>(null)
const isSubmitting = ref(false)
const formErrors = reactive<Record<string, string>>({})
const hasAttemptedSubmit = ref(false)

// Form data
const formData = reactive({
  email: ''
})

// Validation functions
const validateField = (field: string, value: string) => {
  switch (field) {
    case 'email':
      if (!value) return t('validation.emailRequired')
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
      if (!emailRegex.test(value)) return t('validation.emailInvalid')
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
    
    await authStore.forgotPassword(formData.email)
    
    successMessage.value = t('forgot.successMessage')
    
    // Reset form
    Object.assign(formData, { email: '' })
    Object.keys(formErrors).forEach(key => delete formErrors[key])
    
  } catch (error: any) {
    errorMessage.value = error.message || t('forgot.errorMessage')
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <SingleFormCard :title="$t('forgot.title')" :subtitle="$t('forgot.subtitle')">
    <div 
      v-if="successMessage" 
      class="alert alert-success mb-4"
      data-testid="forgot-success-message"
    >
      <span class="break-words">{{ successMessage }}</span>
    </div>
    
    <div 
      v-if="errorMessage" 
      class="alert alert-error mb-4"
      data-testid="forgot-error-message"
    >
      <span class="break-words">{{ errorMessage }}</span>
    </div>
    
    <form @submit="onSubmit" autocomplete="off" class="space-y-6">
      <!-- Account Recovery Information -->
      <fieldset class="space-y-4">
        <legend class="text-base font-semibold text-base-content mb-4 break-words text-wrap">Account Recovery Information</legend>
        
        <div class="form-control">
          <label class="label">
            <span class="label-text">{{ $t('forgot.email') }}</span>
          </label>
          <input 
            type="email" 
            autocomplete="email" 
            placeholder="john.doe@example.com" 
            v-model="formData.email"
            @blur="onFieldBlur('email', formData.email)"
            :class="['input input-bordered w-full', formErrors.email ? 'input-error' : '']"
            data-testid="forgot-email"
            required
          />
          <label class="label">
            <span class="label-text-alt break-words">{{ $t('forgot.emailDescription') }}</span>
          </label>
          <label v-if="formErrors.email" class="label">
            <span class="label-text-alt text-error break-words">{{ formErrors.email }}</span>
          </label>
        </div>
      </fieldset>
      
      <!-- Submit Button -->
      <button 
        type="submit" 
        class="btn btn-primary w-full"
        :disabled="isSubmitting"
        data-testid="forgot-submit"
      >
        <span v-if="isSubmitting" class="loading loading-spinner loading-sm"></span>
        <span v-if="isSubmitting">{{ $t('forgot.sending') }}</span>
        <span v-else>{{ $t('forgot.sendLink') }}</span>
      </button>
    </form>
  </SingleFormCard>
</template>