# PasswordField Component

A reusable Vue 3 component for password input fields with optional visual password requirements validation.

## Features

- ✅ **Password Input Field**: Standard password input with customizable labels and placeholders
- ✅ **Visual Requirements**: Real-time password requirements checking with visual indicators
- ✅ **Strength Meter**: Progress bar showing password strength percentage
- ✅ **Color-coded Feedback**: Requirements and strength indicators with appropriate colors
- ✅ **Internationalization**: Full i18n support for all text labels
- ✅ **Accessibility**: Proper ARIA labels and semantic HTML structure
- ✅ **TypeScript Support**: Full TypeScript definitions and type safety

## Usage

### Basic Password Field

```vue
<template>
  <PasswordField
    v-model="password"
    label="Password"
    :show-requirements="false"
  />
</template>

<script setup>
import { PasswordField } from '@agile-exec/core-frontend/components'
import { ref } from 'vue'

const password = ref('')
</script>
```

### Password Field with Requirements

```vue
<template>
  <PasswordField
    v-model="password"
    label="Create Password"
    placeholder="Enter a secure password"
    :show-requirements="true"
    @requirements-changed="onRequirementsChanged"
  />
</template>

<script setup>
import { PasswordField } from '@agile-exec/core-frontend/components'
import { ref } from 'vue'

const password = ref('')

const onRequirementsChanged = (isValid: boolean, checks: any) => {
  console.log('Password requirements met:', isValid)
  console.log('Individual checks:', checks)
}
</script>
```

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `modelValue` | `string` | `""` | The password value (v-model) |
| `label` | `string` | `""` | Label text displayed above the input |
| `placeholder` | `string` | `"••••••••"` | Placeholder text in the input |
| `autocomplete` | `string` | `"new-password"` | HTML autocomplete attribute |
| `error` | `string` | `""` | Error message to display |
| `testId` | `string` | `"password-field"` | Test ID for the input element |
| `required` | `boolean` | `false` | Whether the field is required |
| `minLength` | `number` | `8` | Minimum password length |
| `showRequirements` | `boolean` | `true` | Show visual password requirements |

## Events

| Event | Payload | Description |
|-------|---------|-------------|
| `update:modelValue` | `string` | Emitted when password value changes |
| `blur` | `string` | Emitted when input loses focus |
| `requirements-changed` | `(isValid: boolean, checks: any)` | Emitted when requirements validation changes |

## Password Requirements

When `showRequirements` is `true`, the component displays:

1. **Strength Indicator**: Shows password strength as "Weak", "Fair", "Good", or "Strong"
2. **Progress Bar**: Visual representation of requirements completion (0-100%)
3. **Requirements List**: Individual requirements with checkmarks:
   - Minimum length (configurable)
   - Uppercase letter (if required by server)
   - Lowercase letter (always required)
   - Number (if required by server) 
   - Special character (if required by server)

## Styling

The component uses DaisyUI classes and follows the design system:

- **Input**: `input input-bordered` with error state support
- **Labels**: `label-text` for main label, `label-text-alt` for errors
- **Requirements**: Color-coded with `success`, `error`, `warning`, `info` colors
- **Progress Bar**: Smooth transitions with appropriate color mapping

## Server Integration

Password requirements are loaded from the server via the `usePasswordRequirements` composable:

```typescript
// Server response example
{
  "minLength": 8,
  "capital": true,    // Requires uppercase letter
  "numbers": true,    // Requires numbers
  "special": true     // Requires special characters
}
```

## Accessibility

- Proper label association with input field
- Error messages linked to input via ARIA
- Semantic HTML structure with `form-control` containers
- Color indicators supplemented with checkmark icons
- Screen reader friendly requirement descriptions

## Examples

### Registration Form
```vue
<PasswordField
  v-model="formData.password"
  :label="$t('register.password')"
  :error="formErrors.password"
  test-id="register-password"
  :required="true"
  :show-requirements="true"
  @requirements-changed="onPasswordRequirementsChanged"
/>
```

### Password Confirmation (No Requirements)
```vue
<PasswordField
  v-model="formData.passwordRepeat"
  :label="$t('register.passwordRepeat')"
  :error="formErrors.passwordRepeat"
  test-id="register-password-repeat"
  :required="true"
  :show-requirements="false"
/>
```

### Password Change Form
```vue
<PasswordField
  v-model="newPassword"
  label="New Password"
  autocomplete="new-password"
  :show-requirements="true"
  @requirements-changed="validateNewPassword"
/>
```