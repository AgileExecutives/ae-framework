export default {
  auth: {
    signIn: 'Sign In',
    signInSubtitle: 'Sign in to your account',
    email: 'Email',
    password: 'Password',
    forgotPassword: 'Forgot your password?',
    or: 'or',
    noAccount: "Don't have an account?",
    signUp: 'Sign up'
  },

  forgot: {
    title: 'Forgot Password',
    subtitle: 'Enter your email to receive a password reset link',
    email: 'Email',
    sendLink: 'Send Reset Link',
    sending: 'Sending...',
    successMessage: 'If the email exists, a password reset link has been sent to your email.',
    errorMessage: 'Failed to send reset email. Please try again.'
  },

  register: {
    title: 'Create Account',
    subtitle: 'Sign up for a new account',
    firstname: 'First Name',
    lastname: 'Last Name',
    email: 'Email',
    password: 'Password',
    passwordRepeat: 'Repeat Password',
    passwordRepeatDescription: 'Please enter your password again',
    passwordRequirements: 'Must be at least 8 characters long',
    passwordStrength: 'Password Strength',
    acceptTerms: 'I accept the {termsOfService} and {privacyPolicy}',
    acceptTermsSimple: 'I accept the Terms of Service and Privacy Policy',
    newsletterOptIn: 'I would like to receive email updates (product news, learning therapy tips, new features). You can unsubscribe at any time.',
    createAccount: 'Create Account',
    creatingAccount: 'Creating account...',
    successMessage: 'Registration successful! Redirecting to dashboard...',
    errorMessage: 'Registration failed. Please try again.'
  },

  reset: {
    title: 'Reset Password',
    subtitle: 'Enter your new password',
    newPassword: 'New Password',
    newPasswordPlaceholder: 'Enter new password',
    confirmPassword: 'Repeat Password',
    confirmPasswordPlaceholder: 'Repeat new password',
    passwordRequirements: 'Password Requirements',
    resetButton: 'Reset Password',
    resetting: 'Resetting...',
    successMessage: 'Password has been reset successfully. Redirecting to login...',
    errorMessage: 'Failed to reset password. The link may be invalid or expired.'
  },

  change: {
    title: 'Change Password',
    subtitle: 'Update your account password',
    currentPassword: 'Current Password',
    currentPasswordPlaceholder: 'Enter current password',
    newPassword: 'New Password',
    newPasswordPlaceholder: 'Enter new password',
    confirmPassword: 'Repeat Password',
    confirmPasswordPlaceholder: 'Repeat new password',
    passwordRequirements: 'Must be at least 8 characters long',
    confirmPasswordDescription: 'Please enter your new password again',
    changeButton: 'Change Password',
    changing: 'Changing...',
    successMessage: 'Password has been changed successfully!',
    errorMessage: 'Failed to change password. Please check your current password.'
  },

  validation: {
    firstNameRequired: 'First name is required',
    firstNameMin: 'First name must be at least 2 characters',
    firstNameMax: 'First name must be at most 50 characters',
    lastNameRequired: 'Last name is required',
    lastNameMin: 'Last name must be at least 2 characters',
    lastNameMax: 'Last name must be at most 50 characters',
    emailRequired: 'Email is required',
    emailInvalid: 'Please enter a valid email address',
    passwordRequired: 'Password is required',
    passwordMin: 'Password must be at least 8 characters',
    passwordRepeatRequired: 'Please confirm your password',
    passwordsDontMatch: "Passwords don't match",
    currentPasswordRequired: 'Current password is required',
    passwordSameAsCurrent: 'New password must be different from current password',
    termsRequired: 'You must accept the Terms of Service and Privacy Policy'
  },

  passwordRequirements: {
    minLength: 'At least 8 characters',
    hasUppercase: 'At least one uppercase letter',
    hasLowercase: 'At least one lowercase letter',
    hasNumber: 'At least one number',
    hasSpecialChar: 'At least one special character',
    noCommonWords: 'No common words',
    notTooSimilar: 'Not too similar to personal information'
  },

  passwordStrength: {
    veryWeak: 'very weak',
    weak: 'weak',
    middle: 'middle',
    strong: 'strong'
  },

  notFound: {
    title: 'Page Not Found',
    message: 'Sorry, the page you are looking for does not exist or has been moved.',
    goBack: 'Go Back',
    goHome: 'Go to Home',
    helpfulLinks: 'You might find these links helpful:',
    mobileHint: 'Tap the back button or try the home page'
  }
}