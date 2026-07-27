import { ref, computed, onMounted } from 'vue'
import { getApiClient } from '@/config/api-config'
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
      
      // Try to fetch password requirements from the password security endpoint
      try {
        const apiClient = getApiClient()
        const response = await apiClient.getPasswordSecurity()
        console.log('🔐 Password security API response:', response)
        
        if (response && typeof response === 'object') {
          // The API returns the PasswordRequirements object directly
          // Type: { capital?: boolean, minLength?: number, numbers?: boolean, special?: boolean }
          const passwordSettings = response as any
          
          requirements.value = {
            minLength: passwordSettings.minLength ?? 8,
            capital: passwordSettings.capital ?? true,
            numbers: passwordSettings.numbers ?? true,
            special: passwordSettings.special ?? true
          }
          console.log('✅ Loaded password requirements from backend:', requirements.value)
          return
        }
        
        console.warn('⚠️ Password security response format unexpected:', response)
      } catch (settingsErr) {
        console.error('❌ Password security endpoint error:', settingsErr)
        // Settings endpoint may not exist or may not have password requirements
      }
      
      // Fall back to default requirements
      console.log('ℹ️ Using default password requirements:', requirements.value)
    } catch (err: any) {
      console.error('Failed to load password requirements:', err)
      error.value = err.message
      // Keep using default requirements if loading fails
    } finally {
      isLoading.value = false
    }
  }


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
    
    console.log('🔍 Validating password with requirements:', req)

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
    
    console.log('🔍 Validation result:', { valid: errors.length === 0, errors })

    return {
      valid: errors.length === 0,
      errors
    }
  }

  // Check individual requirements for visual indicators
  const checkRequirements = (password: string) => {
    const req = requirements.value
    
    // Only require lowercase if at least one other character requirement is enabled
    const requireLowercase = req.capital || req.numbers || req.special
    
    const checks = computed(() => [
      {
        id: 'length',
        label: t('passwordRequirements.atLeastChars', { n: req.minLength }),
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
        required: requireLowercase
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

  // Auto-load on mount
  onMounted(() => {
    loadRequirements()
  })

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
