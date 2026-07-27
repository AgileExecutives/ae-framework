<script setup lang="ts">
// Reset Password Form - handles password reset with token validation
import SingleFormCard from "./SingleFormCard.vue"
import PasswordField from "./PasswordField.vue"
import { useAuthStore } from "../stores/auth"
import { useRouter } from "vue-router"
import { ref, reactive, onMounted } from "vue"
import { useI18n } from 'vue-i18n'
import { usePasswordRequirements } from '../composables/usePasswordRequirements'
import { getApiClient } from '@/config/api-config'

const { t } = useI18n()
const authStore = useAuthStore()
const router = useRouter()
const { validatePassword, requirements } = usePasswordRequirements()

const token = ref<string>('')
const successMessage = ref<string | null>(null)
const errorMessage = ref<string | null>(null)
const isValidatingToken = ref(false)

// Validate token on mount
onMounted(async () => {
  // Get token from query params
  const queryToken = router.currentRoute.value.query.token as string
  if (!queryToken) {
    router.push('/forgot-password')
    return
  }
  
  token.value = queryToken
  console.log('🔑 Validating reset token:', token)

  // Validate the token with the backend
  try {
    isValidatingToken.value = true
    const apiClient = getApiClient()
    const data = await apiClient.checkResetToken(token.value)
    console.log('✅ Reset token is valid' + JSON.stringify(data) )
  } catch (error: any) {
    console.error('❌ Reset token validation failed:', error)
    // Redirect to error page with message
    router.push({
      path: '/error',
      query: {
        title: t('reset.invalidTokenTitle'),
        message: t('reset.invalidTokenMessage')
      }
    })
  } finally {
    isValidatingToken.value = false
  }
})
const isSubmitting = ref(false)
const formErrors = reactive<Record<string, string>>({})
const hasAttemptedSubmit = ref(false)

// Password requirements state
const newPasswordRequirementsMet = ref(false)

// Password requirements handler
const onNewPasswordRequirementsChanged = (isValid: boolean, checks: any) => {
  newPasswordRequirementsMet.value = isValid
}

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
      const passwordValidation = validatePassword(value)
      if (!passwordValidation.valid) return passwordValidation.errors[0]
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
    
    console.log('🔐 Resetting password with token:', token.value.substring(0, 10) + '...')
    await authStore.resetPassword(token.value, formData.newPassword)
    
    // Redirect to success page
    router.push({ 
      path: '/success', 
      query: { type: 'password-reset' } 
    })
  } catch (error: any) {
    errorMessage.value = error.message || t('reset.errorMessage')
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <SingleFormCard :title="$t('reset.title')" :subtitle="$t('reset.subtitle')">
    <!-- Token validation loading state -->
    <div 
      v-if="isValidatingToken" 
      class="alert alert-info mb-4"
      data-testid="reset-validating-token"
    >
      <span class="loading loading-spinner"></span>
      <span>{{ $t('reset.validatingToken') }}</span>
    </div>
    
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
    
    <form @submit="onSubmit" autocomplete="off" class="space-y-6" v-if="!isValidatingToken">
      <!-- New Password Setup -->
      <fieldset class="space-y-4">
        <legend class="text-base font-semibold text-base-content mb-4 break-words text-wrap">New Password Setup</legend>
        
        <PasswordField
          v-model="formData.newPassword"
          :label="$t('reset.newPassword')"
          name="newPassword"
          autocomplete="new-password"
          :error="formErrors.newPassword"
          :min-length="requirements.minLength"
          :show-requirements="true"
          @blur="onFieldBlur('newPassword', formData.newPassword)"
          @input="onFieldInput('confirmPassword', formData.confirmPassword)"
          @requirements-changed="onNewPasswordRequirementsChanged"
          test-id="reset-password"
          required
        />
        
        <PasswordField
          v-model="formData.confirmPassword"
          :label="$t('reset.confirmPassword')"
          name="confirmPassword"
          autocomplete="new-password"
          :error="formErrors.confirmPassword"
          :min-length="requirements.minLength"
          @blur="onFieldBlur('confirmPassword', formData.confirmPassword)"
          test-id="reset-password-repeat"
          required
        />
      </fieldset>
      
      <!-- Submit Button -->
      <button 
        type="submit" 
        class="btn btn-primary w-full"
        :disabled="isSubmitting || !newPasswordRequirementsMet"
        data-testid="reset-submit"
      >
        <span v-if="isSubmitting" class="loading loading-spinner loading-sm"></span>
        <span v-if="isSubmitting">{{ $t('reset.resetting') }}</span>
        <span v-else>{{ $t('reset.resetButton') }}</span>
      </button>
    </form>
  </SingleFormCard>
</template>
