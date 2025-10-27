import { ref, computed, onMounted } from 'vue'
import { createApiClient } from '@agile-exec/api-client'
import * as z from 'zod'
import { useI18n } from 'vue-i18n'

export interface PasswordRequirements {
  minLength: number
  capital: boolean
  numbers: boolean
  special: boolean
}

export function usePasswordRequirements() {
  const { t } = useI18n()
  
  const requirements = ref<PasswordRequirements>({
    minLength: 8,
    capital: true,
    numbers: true,
    special: true
  })
  
  const isLoading = ref(true)
  const error = ref<string | null>(null)

  // Load requirements from server
  const loadRequirements = async () => {
    try {
      isLoading.value = true
      error.value = null
      
      const apiClient = createApiClient({
        baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'
      })
      
      const data = await apiClient.getPasswordSecurity()
      requirements.value = data
    } catch (err: any) {
      console.error('Failed to load password requirements:', err)
      error.value = err.message
      // Keep using default requirements if loading fails
    } finally {
      isLoading.value = false
    }
  }

  // Auto-load on mount
  onMounted(() => {
    loadRequirements()
  })

  // Build error message for validation
  const buildErrorMessage = (req: PasswordRequirements): string => {
    const parts: string[] = []
    if (req.capital) parts.push(t('passwordRequirements.oneUppercase'))
    if (req.numbers) parts.push(t('passwordRequirements.oneNumber'))
    if (req.special) parts.push(t('passwordRequirements.oneSpecial'))
    
    if (parts.length === 0) return ''
    return t('validation.passwordRequirements', { requirements: parts.join(', ') })
  }

  // Build description for UI
  const description = computed(() => {
    const req = requirements.value
    const parts: string[] = [t('passwordRequirements.atLeastChars', { n: req.minLength })]
    
    if (req.capital) parts.push(t('passwordRequirements.oneUppercase'))
    if (req.numbers) parts.push(t('passwordRequirements.oneNumber'))
    if (req.special) parts.push(t('passwordRequirements.oneSpecial'))
    
    return parts.join(', ')
  })

  // Create Zod schema for password validation
  const createPasswordSchema = () => {
    const req = requirements.value
    
    return z.string()
      .min(req.minLength, t('validation.passwordMinLength', { n: req.minLength }))
      .refine(
        (password) => {
          if (req.capital && !/[A-Z]/.test(password)) return false
          if (req.numbers && !/[0-9]/.test(password)) return false
          if (req.special && !/[^a-zA-Z0-9]/.test(password)) return false
          return true
        },
        {
          message: buildErrorMessage(req)
        }
      )
  }

  // Validate password manually
  const validatePassword = (password: string): { valid: boolean; errors: string[] } => {
    const req = requirements.value
    const errors: string[] = []

    if (password.length < req.minLength) {
      errors.push(t('validation.passwordMinLength', { n: req.minLength }))
    }

    if (req.capital && !/[A-Z]/.test(password)) {
      errors.push(t('validation.passwordRequireUppercase'))
    }

    if (req.numbers && !/[0-9]/.test(password)) {
      errors.push(t('validation.passwordRequireNumber'))
    }

    if (req.special && !/[^a-zA-Z0-9]/.test(password)) {
      errors.push(t('validation.passwordRequireSpecial'))
    }

    return {
      valid: errors.length === 0,
      errors
    }
  }

  // Check individual requirements for visual indicators
  const checkRequirements = (password: string) => {
    const req = requirements.value
    
    const checks = computed(() => [
      {
        id: 'length',
        label: t('passwordRequirements.minLength'),
        met: password.length >= req.minLength,
        required: true
      },
      {
        id: 'uppercase',
        label: t('passwordRequirements.hasUppercase'),
        met: /[A-Z]/.test(password),
        required: req.capital
      },
      {
        id: 'lowercase',
        label: t('passwordRequirements.hasLowercase'),
        met: /[a-z]/.test(password),
        required: true
      },
      {
        id: 'number',
        label: t('passwordRequirements.hasNumber'),
        met: /[0-9]/.test(password),
        required: req.numbers
      },
      {
        id: 'special',
        label: t('passwordRequirements.hasSpecialChar'),
        met: /[^a-zA-Z0-9]/.test(password),
        required: req.special
      }
    ].filter(check => check.required))

    const allRequiredMet = computed(() => 
      checks.value.every(check => check.met)
    )

    const metCount = computed(() => 
      checks.value.filter(check => check.met).length
    )

    const strengthScore = computed(() => {
      // If password is empty, return 0
      if (!password || password.length === 0) return 0
      
      const requiredCount = checks.value.length
      if (requiredCount === 0) return 0
      
      return Math.round((metCount.value / requiredCount) * 100)
    })

    const strengthColor = computed(() => {
      const score = strengthScore.value
      // Only show success when ALL requirements are met (100%)
      if (score === 0) return 'text-base-content'
      if (score < 50) return 'text-error'
      if (score < 75) return 'text-warning'  
      if (score < 100) return 'text-info'
      return 'text-success'
    })

    const strengthLabel = computed(() => {
      const score = strengthScore.value
      if (score === 0) return ''
      if (score < 50) return t('passwordStrength.veryWeak')
      if (score < 75) return t('passwordStrength.weak')
      if (score < 100) return t('passwordStrength.middle')
      return t('passwordStrength.strong')
    })

    return {
      checks,
      allRequiredMet,
      metCount,
      strengthScore,
      strengthColor,
      strengthLabel
    }
  }

  return {
    requirements,
    isLoading,
    error,
    description,
    loadRequirements,
    createPasswordSchema,
    validatePassword,
    checkRequirements
  }
}
