export default {
  auth: {
    signIn: 'Anmelden',
    signInSubtitle: 'Melden Sie sich an',
    email: 'E-Mail',
    password: 'Passwort',
    forgotPassword: 'Passwort vergessen?',
    or: 'oder',
    noAccount: 'Noch kein Konto?',
    signUp: 'Registrieren'
  },

  forgot: {
    title: 'Passwort vergessen',
    subtitle: 'Geben Sie Ihre E-Mail-Adresse ein, um einen Link zum Zurücksetzen zu erhalten',
    email: 'E-Mail',
    sendLink: 'Link senden',
    sending: 'Wird gesendet...',
    successMessage: 'Wenn die E-Mail existiert, wurde ein Link zum Zurücksetzen an Ihre E-Mail gesendet.',
    errorMessage: 'Fehler beim Senden der E-Mail. Bitte versuchen Sie es erneut.'
  },

  register: {
    title: 'Konto erstellen',
    subtitle: 'Registrieren Sie sich für ein neues Konto',
    firstname: 'Vorname',
    lastname: 'Nachname',
    email: 'E-Mail',
    password: 'Passwort',
    passwordRepeat: 'Passwort wiederholen',
    passwordRepeatDescription: 'Bitte geben Sie Ihr Passwort erneut ein',
    passwordRequirements: 'Muss mindestens 8 Zeichen lang sein',
    passwordStrength: 'Passwort-Stärke',
    acceptTerms: 'Ich akzeptiere die {termsOfService} und {privacyPolicy}',
    acceptTermsSimple: 'Ich akzeptiere die AGB und Datenschutzerklärung',
    newsletterOptIn: 'Ich möchte per E-Mail informiert werden (Produktneuigkeiten, Lerntherapie-Tipps, neue Funktionen). Du kannst dich jederzeit abmelden.',
    createAccount: 'Konto erstellen',
    creatingAccount: 'Konto wird erstellt...',
    successMessage: 'Registrierung erfolgreich! Weiterleitung zum Dashboard...',
    errorMessage: 'Registrierung fehlgeschlagen. Bitte versuchen Sie es erneut.'
  },

  reset: {
    title: 'Passwort zurücksetzen',
    subtitle: 'Geben Sie Ihr neues Passwort ein',
    newPassword: 'Neues Passwort',
    newPasswordPlaceholder: 'Neues Passwort eingeben',
    confirmPassword: 'Passwort wiederholen',
    confirmPasswordPlaceholder: 'Neues Passwort wiederholen',
    passwordRequirements: 'Passwort-Anforderungen',
    resetButton: 'Passwort zurücksetzen',
    resetting: 'Wird zurückgesetzt...',
    successMessage: 'Passwort wurde erfolgreich zurückgesetzt. Weiterleitung zur Anmeldung...',
    errorMessage: 'Fehler beim Zurücksetzen des Passworts. Der Link ist möglicherweise ungültig oder abgelaufen.'
  },

  change: {
    title: 'Passwort ändern',
    subtitle: 'Aktualisieren Sie Ihr Kontopasswort',
    currentPassword: 'Aktuelles Passwort',
    currentPasswordPlaceholder: 'Aktuelles Passwort eingeben',
    newPassword: 'Neues Passwort',
    newPasswordPlaceholder: 'Neues Passwort eingeben',
    confirmPassword: 'Passwort wiederholen',
    confirmPasswordPlaceholder: 'Neues Passwort bestätigen',
    passwordRequirements: 'Muss mindestens 8 Zeichen lang sein',
    confirmPasswordDescription: 'Bitte geben Sie Ihr neues Passwort erneut ein',
    changeButton: 'Passwort ändern',
    changing: 'Wird geändert...',
    successMessage: 'Passwort wurde erfolgreich geändert!',
    errorMessage: 'Fehler beim Ändern des Passworts. Bitte überprüfen Sie Ihr aktuelles Passwort.'
  },

  validation: {
    firstNameRequired: 'Vorname ist erforderlich',
    firstNameMin: 'Vorname muss mindestens 2 Zeichen lang sein',
    firstNameMax: 'Vorname darf maximal 50 Zeichen lang sein',
    lastNameRequired: 'Nachname ist erforderlich',
    lastNameMin: 'Nachname muss mindestens 2 Zeichen lang sein',
    lastNameMax: 'Nachname darf maximal 50 Zeichen lang sein',
    emailRequired: 'E-Mail ist erforderlich',
    emailInvalid: 'Bitte geben Sie eine gültige E-Mail-Adresse ein',
    passwordRequired: 'Passwort ist erforderlich',
    passwordMin: 'Passwort muss mindestens 8 Zeichen lang sein',
    passwordRepeatRequired: 'Bitte bestätigen Sie Ihr Passwort',
    passwordsDontMatch: 'Die Passwörter stimmen nicht überein',
    currentPasswordRequired: 'Aktuelles Passwort ist erforderlich',
    passwordSameAsCurrent: 'Neues Passwort muss sich vom aktuellen unterscheiden',
    termsRequired: 'Sie müssen die AGB und Datenschutzerklärung akzeptieren'
  },

  passwordRequirements: {
    minLength: 'Mindestens 8 Zeichen',
    hasUppercase: 'Mindestens ein Großbuchstabe',
    hasLowercase: 'Mindestens ein Kleinbuchstabe', 
    hasNumber: 'Mindestens eine Zahl',
    hasSpecialChar: 'Mindestens ein Sonderzeichen',
    noCommonWords: 'Keine häufigen Wörter',
    notTooSimilar: 'Nicht zu ähnlich zu persönlichen Informationen'
  },

  passwordStrength: {
    veryWeak: 'sehr schwach',
    weak: 'schwach',
    middle: 'mittelmäßig', 
    strong: 'stark'
  },



  notFound: {
    title: 'Seite nicht gefunden',
    message: 'Entschuldigung, die gesuchte Seite existiert nicht oder wurde verschoben.',
    goBack: 'Zurück',
    goHome: 'Zur Startseite',
    helpfulLinks: 'Diese Links könnten hilfreich sein:',
    mobileHint: 'Tippe auf den Zurück-Button oder versuche die Startseite'
  }
}

