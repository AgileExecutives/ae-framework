<template>
  <div class="form-control">
    <label class="label">
      <span class="label-text">{{ label }}</span>
    </label>
    <input 
      type="password" 
      :placeholder="placeholder"
      :autocomplete="autocomplete"
      :value="modelValue"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
      @blur="$emit('blur', ($event.target as HTMLInputElement).value)"
      :class="['input input-bordered w-full', error ? 'input-error' : '']"
      :data-testid="testId"
      :required="required"
      :minlength="minLength"
    />
    
    <!-- Password Requirements Visual Indicators -->
    <div v-if="showRequirements && modelValue" class="mt-2 space-y-2">
      <!-- Password Strength Indicator -->
      <div class="flex items-center justify-between">
        <span class="text-xs font-medium">{{ $t('register.passwordStrength') }}</span>
        <span :class="['text-xs font-medium', passwordChecks.strengthColor.value]">
          {{ passwordChecks.strengthLabel.value }}
        </span>
      </div>
      
      <!-- Strength Progress Bar -->
      <div class="w-full bg-base-200 rounded-full h-2">
        <div 
          class="h-2 rounded-full transition-all duration-300"
          :class="[
            passwordChecks.strengthScore.value === 0 ? 'bg-base-200' :
            passwordChecks.strengthScore.value < 50 ? 'bg-error' :
            passwordChecks.strengthScore.value < 75 ? 'bg-warning' :
            passwordChecks.strengthScore.value < 100 ? 'bg-info' : 'bg-success'
          ]"
          :style="`width: ${passwordChecks.strengthScore.value}%`">
        </div>
      </div>
      
      <!-- Requirements List -->
      <div class="space-y-1">
        <div 
          v-for="requirement in passwordChecks.checks.value" 
          :key="requirement.id"
          class="flex items-center gap-2 text-xs"
        >
          <div :class="[
            'w-3 h-3 rounded-full flex items-center justify-center',
            requirement.met ? 'bg-success text-success-content' : 'bg-base-300 text-base-content'
          ]">
            <svg v-if="requirement.met" class="w-2 h-2" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
            </svg>
          </div>
          <span :class="[
            'break-words',
            requirement.met ? 'text-success' : 'text-base-content/70'
          ]">
            {{ requirement.label }}
          </span>
        </div>
      </div>
    </div>
    
    <label v-if="error" class="label">
      <span class="label-text-alt text-error break-words">{{ error }}</span>
    </label>
  </div>
</template>

<script setup lang="ts">
import { computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { usePasswordRequirements } from '../composables/usePasswordRequirements'

export interface PasswordFieldProps {
  modelValue: string
  label?: string
  placeholder?: string
  autocomplete?: string
  error?: string
  testId?: string
  required?: boolean
  minLength?: number
  showRequirements?: boolean
}

const props = withDefaults(defineProps<PasswordFieldProps>(), {
  label: '',
  placeholder: '••••••••',
  autocomplete: 'new-password',
  error: '',
  testId: 'password-field',
  required: false,
  minLength: 8,
  showRequirements: true
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
  'blur': [value: string]
  'requirements-changed': [isValid: boolean, checks: any]
}>()

const { t } = useI18n()
const { checkRequirements } = usePasswordRequirements()

// Get password checks when showRequirements is enabled
const passwordChecks = computed(() => {
  if (!props.showRequirements) {
    return {
      checks: { value: [] },
      strengthScore: { value: 0 },
      strengthColor: { value: 'text-base-content' },
      strengthLabel: { value: '' }
    }
  }
  return checkRequirements(props.modelValue || '')
})

// Emit requirements changes for parent component validation
watch(() => passwordChecks.value, (newChecks) => {
  if (props.showRequirements) {
    const allMet = newChecks.checks.value.every((check: any) => check.met)
    emit('requirements-changed', allMet, newChecks)
  }
}, { deep: true })
</script>