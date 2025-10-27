<script setup lang="ts">
import SingleFormCard from "./SingleFormCard.vue"
import PasswordField from "./PasswordField.vue"
import { useAuthStore } from "../stores/auth"
import { useRouter } from "vue-router"
import { ref, reactive } from "vue"
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const authStore = useAuthStore()
const router = useRouter()

const successMessage = ref<string | null>(null)
const errorMessage = ref<string | null>(null)
const isSubmitting = ref(false)
const formErrors = reactive<Record<string, string>>({})
const hasAttemptedSubmit = ref(false)

// Form data
const formData = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// Password requirements state
const newPasswordRequirementsMet = ref(false)

// Password requirements handler
const onNewPasswordRequirementsChanged = (isValid: boolean, checks: any) => {
  newPasswordRequirementsMet.value = isValid
}

// Validation functions
const validateField = (field: string, value: string) => {
  switch (field) {
    case 'currentPassword':
      if (!value) return t('validation.currentPasswordRequired')
      break
    case 'newPassword':
      if (!value) return t('validation.passwordRequired')
      if (value.length < 8) return t('validation.passwordMin')
      if (value === formData.currentPassword) return t('validation.passwordSameAsCurrent')
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

    await authStore.changePassword({
      current_password: formData.currentPassword,
      new_password: formData.newPassword
    })

    successMessage.value = t('change.successMessage')
    
    // Reset form
    Object.assign(formData, {
      currentPassword: '',
      newPassword: '',
      confirmPassword: ''
    })
    Object.keys(formErrors).forEach(key => delete formErrors[key])
    
  } catch (error: any) {
    errorMessage.value = error.message || t('change.errorMessage')
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <SingleFormCard :title="$t('change.title')" :subtitle="$t('change.subtitle')">
    <div 
      v-if="successMessage" 
      class="alert alert-success mb-4"
      data-testid="change-success-message"
    >
      <span class="break-words">{{ successMessage }}</span>
    </div>
    
    <div 
      v-if="errorMessage" 
      class="alert alert-error mb-4"
      data-testid="change-error-message"
    >
      <span class="break-words">{{ errorMessage }}</span>
    </div>
    
    <form @submit="onSubmit" autocomplete="off" class="space-y-6">
      <!-- Current Authentication -->
      <fieldset class="space-y-4">
        <legend class="text-base font-semibold text-base-content mb-4 break-words text-wrap">Current Authentication</legend>
        
        <PasswordField
          v-model="formData.currentPassword"
          :label="$t('change.currentPassword')"
          placeholder="••••••••" 
          autocomplete="current-password"
          :error="formErrors.currentPassword"
          test-id="change-current-password"
          :required="true"
          :show-requirements="false"
          @blur="onFieldBlur('currentPassword', formData.currentPassword)"
        />
      </fieldset>

      <!-- New Password Setup -->
      <fieldset class="space-y-4">
        <legend class="text-base font-semibold text-base-content mb-4 break-words text-wrap">New Password Setup</legend>
        
        <PasswordField
          v-model="formData.newPassword"
          :label="$t('change.newPassword')"
          placeholder="••••••••" 
          autocomplete="new-password"
          :error="formErrors.newPassword"
          test-id="change-new-password"
          :required="true"
          :min-length="8"
          :show-requirements="true"
          @blur="onFieldBlur('newPassword', formData.newPassword)"
          @update:model-value="onFieldInput('confirmPassword', formData.confirmPassword)"
          @requirements-changed="onNewPasswordRequirementsChanged"
        />

        <PasswordField
          v-model="formData.confirmPassword"
          :label="$t('change.confirmPassword')"
          placeholder="••••••••" 
          autocomplete="new-password"
          :error="formErrors.confirmPassword"
          test-id="change-confirm-password"
          :required="true"
          :show-requirements="false"
          @blur="onFieldBlur('confirmPassword', formData.confirmPassword)"
        />

      </fieldset>

      <!-- Submit Button -->
      <button 
        type="submit" 
        class="btn btn-primary w-full mt-8"
        :disabled="isSubmitting"
        data-testid="change-submit"
      >
        <span v-if="isSubmitting" class="loading loading-spinner loading-sm"></span>
        <span v-if="isSubmitting">{{ $t('change.changing') }}</span>
        <span v-else>{{ $t('change.changeButton') }}</span>
      </button>
    </form>
  </SingleFormCard>
</template>
