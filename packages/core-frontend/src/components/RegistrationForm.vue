<script setup lang="ts">
import SingleFormCard from "./SingleFormCard.vue"
import LegalLinks from "./LegalLinks.vue"
import PasswordField from "./PasswordField.vue"
import { useAuthStore } from "../stores/auth"
import { useRouter } from "vue-router"
import { ref, reactive, computed } from "vue"
import { useI18n } from 'vue-i18n'
import { useAppConfig } from '../composables/useAppConfig'
import { usePasswordRequirements } from '../composables/usePasswordRequirements'

const { t } = useI18n()
const authStore = useAuthStore()
const router = useRouter()
const { termsOfServiceText, privacyPolicyText } = useAppConfig()
const { validatePassword } = usePasswordRequirements()

const successMessage = ref<string | null>(null)
const errorMessage = ref<string | null>(null)
const isSubmitting = ref(false)
const formErrors = reactive<Record<string, string>>({})
const hasAttemptedSubmit = ref(false)

// Create localized accept terms text with links
const acceptTermsText = computed(() => {
  return t('register.acceptTerms', {
    termsOfService: termsOfServiceText.value,
    privacyPolicy: privacyPolicyText.value
  })
})

// Form data
const formData = reactive({
  firstname: '',
  lastname: '',
  email: '',
  password: '',
  passwordRepeat: '',
  acceptTerms: false,
  newsletterOptIn: false
})

// Password requirements state
const passwordRequirementsMet = ref(false)

// Password requirements handler
const onPasswordRequirementsChanged = (isValid: boolean, checks: any) => {
  passwordRequirementsMet.value = isValid
}

// Validation functions
const validateField = (field: string, value: string) => {
  switch (field) {
    case 'firstname':
      if (!value) return t('validation.firstNameRequired')
      if (value.length < 2) return t('validation.firstNameMin')
      if (value.length > 50) return t('validation.firstNameMax')
      break
    case 'lastname':
      if (!value) return t('validation.lastNameRequired')
      if (value.length < 2) return t('validation.lastNameMin')
      if (value.length > 50) return t('validation.lastNameMax')
      break
    case 'email':
      if (!value) return t('validation.emailRequired')
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
      if (!emailRegex.test(value)) return t('validation.emailInvalid')
      break
    case 'password':
      if (!value) return t('validation.passwordRequired')
      const passwordValidation = validatePassword(value)
      if (!passwordValidation.valid) return passwordValidation.errors[0]
      break
    case 'passwordRepeat':
      if (!value) return t('validation.passwordRepeatRequired')
      if (value !== formData.password) return t('validation.passwordsDontMatch')
      break
    case 'acceptTerms':
      if (!value) return t('validation.termsRequired')
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
  
  // Validate text fields
  const textFields = ['firstname', 'lastname', 'email', 'password', 'passwordRepeat'] as const
  textFields.forEach((field) => {
    const error = validateField(field, formData[field])
    if (error) errors[field] = error
  })
  
  // Validate required checkboxes
  if (!formData.acceptTerms) {
    errors.acceptTerms = t('validation.termsRequired')
  }
  
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
    errorMessage.value = null
    successMessage.value = null

    console.log('📝 Registration attempt with values:', {
      email: formData.email,
      hasPassword: !!formData.password,
      acceptTerms: formData.acceptTerms,
      newsletterOptIn: formData.newsletterOptIn
    })

    await authStore.register({
      username: formData.email, // Use email as username
      email: formData.email,
      password: formData.password,
      first_name: formData.firstname,
      last_name: formData.lastname,
      company_name: 'Lerntherapie '+formData.lastname,
      accept_terms: formData.acceptTerms,
      newsletter_opt_in: formData.newsletterOptIn
    })

    successMessage.value = t('register.successMessage')
    
    // Check if user was auto-logged in after registration
    if (authStore.isAuthenticated) {
      // Auto-login successful, redirect to home page
      setTimeout(() => {
        router.push('/')
      }, 1500)
    } else {
      // Registration successful but no auto-login, redirect to login page
      setTimeout(() => {
        router.push('/login')
      }, 2000)
    }

  } catch (error: any) {
    console.error('❌ Registration failed:', error)
    errorMessage.value = error.message || authStore.error || t('register.errorMessage')
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <SingleFormCard :title="$t('register.title')" :subtitle="$t('register.subtitle')">
    <!-- Success Message -->
    <div v-if="successMessage" class="alert alert-success mb-4">
      <span>{{ successMessage }}</span>
    </div>

    <!-- Error Message -->
    <div v-if="errorMessage" class="alert alert-error mb-4">
      <span>{{ errorMessage }}</span>
    </div>

    <form @submit="onSubmit" autocomplete="off" class="space-y-6">
      <!-- Personal Information Fieldset -->
      <fieldset class="space-y-4">
        <legend class="text-base font-semibold text-base-content mb-4 break-words text-wrap">Personal Information</legend>
        
        <!-- First Name -->
        <div class="form-control">
          <label class="label">
            <span class="label-text">{{ $t('register.firstname') }}</span>
          </label>
          <input 
            type="text" 
            placeholder="John" 
            autocomplete="given-name"
            v-model="formData.firstname"
            @blur="onFieldBlur('firstname', formData.firstname)"
            :class="['input input-bordered w-full', formErrors.firstname ? 'input-error' : '']"
            data-testid="register-firstname"
            required
          />
          <label v-if="formErrors.firstname" class="label">
            <span class="label-text-alt text-error">{{ formErrors.firstname }}</span>
          </label>
        </div>

        <!-- Last Name -->
        <div class="form-control flex-1">
          <label class="label">
            <span class="label-text">{{ $t('register.lastname') }}</span>
          </label>
          <input 
            type="text" 
            placeholder="Doe" 
            autocomplete="family-name"
            v-model="formData.lastname"
            @blur="onFieldBlur('lastname', formData.lastname)"
            :class="['input input-bordered w-full', formErrors.lastname ? 'input-error' : '']"
            data-testid="register-lastname"
            required
            minlength="2"
            maxlength="50"
          />
          <label v-if="formErrors.lastname" class="label">
            <span class="label-text-alt text-error">{{ formErrors.lastname }}</span>
          </label>
        </div>

        <!-- Email -->
        <div class="form-control">
          <label class="label">
            <span class="label-text">{{ $t('register.email') }}</span>
          </label>
          <input 
            type="email" 
            placeholder="john.doe@example.com" 
            autocomplete="email"
            v-model="formData.email"
            @blur="onFieldBlur('email', formData.email)"
            :class="['input input-bordered w-full', formErrors.email ? 'input-error' : '']"
            data-testid="register-email"
            required
          />
          <label v-if="formErrors.email" class="label">
            <span class="label-text-alt text-error">{{ formErrors.email }}</span>
          </label>
        </div>
      </fieldset>

      <!-- Security Information Fieldset -->
      <fieldset class="space-y-4">
        <legend class="text-base font-semibold text-base-content mb-4 break-words text-wrap">Security Information</legend>
        
        <!-- Password -->
        <PasswordField
          v-model="formData.password"
          :label="$t('register.password')"
          placeholder="••••••••"
          autocomplete="new-password"
          :error="formErrors.password"
          test-id="register-password"
          :required="true"
          :min-length="8"
          :show-requirements="true"
          @blur="onFieldBlur('password', formData.password)"
          @update:model-value="onFieldInput('passwordRepeat', formData.passwordRepeat)"
          @requirements-changed="onPasswordRequirementsChanged"
        />

        <!-- Password Repeat -->
        <PasswordField
          v-model="formData.passwordRepeat"
          :label="$t('register.passwordRepeat')"
          placeholder="••••••••"
          autocomplete="new-password"
          :error="formErrors.passwordRepeat"
          test-id="register-password-repeat"
          :required="true"
          :show-requirements="false"
          @blur="onFieldBlur('passwordRepeat', formData.passwordRepeat)"
        />
      </fieldset>

      <!-- Terms and Preferences Fieldset -->
      <fieldset class="space-y-4">
        <legend class="text-base font-semibold text-base-content mb-4 break-words text-wrap">Terms and Preferences</legend>
        
        <!-- Terms and Conditions Checkbox -->
        <div class="form-control">
          <label class="label cursor-pointer justify-start gap-3 items-start p-0">
            <input 
              type="checkbox"
              v-model="formData.acceptTerms"
              :class="['checkbox checkbox-primary flex-shrink-0 mt-1', formErrors.acceptTerms ? 'checkbox-error' : '']"
              data-testid="register-accept-terms"
              required
            />
            <span class="label-text text-sm  break-words hyphens-auto leading-relaxed text-wrap min-w-0 flex-1">
              {{ acceptTermsText }}
            </span>
          </label>
          <label v-if="formErrors.acceptTerms" class="label">
            <span class="label-text-alt text-error break-words text-wrap">{{ formErrors.acceptTerms }}</span>
          </label>
          
          <!-- Legal Links -->
          <div class="mt-2 pl-8">
            <LegalLinks :inline="true" separator="|" />
          </div>
        </div>

        <!-- Newsletter Opt-in Checkbox -->
        <div class="form-control">
          <label class="label cursor-pointer justify-start gap-3 items-start p-0">
            <input 
              type="checkbox"
              v-model="formData.newsletterOptIn"
              class="checkbox checkbox-primary flex-shrink-0 mt-1"
              data-testid="register-newsletter-opt-in"
            />
            <span class="label-text text-sm break-words hyphens-auto leading-relaxed text-wrap min-w-0 flex-1">{{ $t('register.newsletterOptIn') }}</span>
          </label>
        </div>
      </fieldset>

      <!-- Submit Button -->
      <button 
        type="submit" 
        class="btn btn-primary w-full"
        :disabled="isSubmitting" 
        data-testid="register-submit"
      >
        <span v-if="isSubmitting" class="loading loading-spinner loading-sm"></span>
        <span v-if="isSubmitting">{{ $t('register.creatingAccount') }}</span>
        <span v-else>{{ $t('register.createAccount') }}</span>
      </button>
    </form>
  </SingleFormCard>
</template>