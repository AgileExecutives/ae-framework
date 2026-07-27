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
  messages: {
    'Invalid credentials': 'Ungültige Anmeldedaten',
    'Invalid request': 'Ungültige Anfrage',
    'Terms not accepted': 'AGB nicht akzeptiert',
    'User already exists': 'Benutzer existiert bereits',
    'Password mismatch': 'Passwort stimmt nicht überein',
    'Company name required': 'Firmenname erforderlich',
    'User created successfully. Please check your email to verify your account.': 'Benutzer erfolgreich erstellt. Bitte überprüfen Sie Ihre E-Mail, um Ihr Konto zu verifizieren.',
    'Login successful': 'Anmeldung erfolgreich',
    'User not found': 'Benutzer nicht gefunden',
    'Username or email already taken': 'Benutzername oder E-Mail bereits vergeben',
    'You must accept the terms and conditions to register': 'Sie müssen die AGB akzeptieren, um sich zu registrieren',
    'Company name is required to create a new tenant': 'Firmenname ist erforderlich, um einen neuen Mandanten zu erstellen',
    'Key: \'UserCreateRequest.FirstName\' Error:Field validation for \'FirstName\' failed on the \'required\' tag': 'Vorname ist erforderlich',
    'Key: \'UserCreateRequest.LastName\' Error:Field validation for \'LastName\' failed on the \'required\' tag': 'Nachname ist erforderlich',
    loginSuccessful: 'Anmeldung erfolgreich! Weiterleitung...',
    loginSuccessfulWelcome: 'Willkommen zurück',
    loginFailed: 'Anmeldung fehlgeschlagen',
    invalidResponseFromServer: 'Ungültige Antwort vom Server',
    registrationSuccessful: 'Registrierung erfolgreich! Weiterleitung zum Dashboard...',
    registrationFailed: 'Registrierung fehlgeschlagen. Bitte versuchen Sie es erneut.',
    passwordChangedSuccessfully: 'Passwort wurde erfolgreich geändert!',
    passwordResetSuccessfully: 'Passwort wurde erfolgreich zurückgesetzt. Weiterleitung zur Anmeldung...',
    forgotPasswordEmailSent: 'Wenn die E-Mail existiert, wurde ein Link zum Zurücksetzen an Ihre E-Mail gesendet.',
    forgotPasswordFailed: 'Fehler beim Senden der E-Mail. Bitte versuchen Sie es erneut.',
    usingDefaultPasswordRequirements: 'Verwende Standard-Passwort-Anforderungen (getPasswordSecurity-Endpunkt nicht verfügbar)',
    failedToFetchClients: 'Fehler beim Laden der Kunden',
    failedToDeleteClient: 'Fehler beim Löschen des Kunden',
    failedToLoadCalendarData: 'Fehler beim Laden der Kalenderdaten',
    failedToCreateEvent: 'Fehler beim Erstellen des Termins',
    failedToDeleteEvent: 'Fehler beim Löschen des Termins',
    calendarDataCachedSuccessfully: 'Kalenderdaten erfolgreich zwischengespeichert und geladen'
  },
  forgot: {
    title: 'Passwort vergessen',
    subtitle: 'Geben Sie Ihre E-Mail-Adresse ein, um einen Link zum Zurücksetzen zu erhalten',
    emailDescription: 'Geben Sie die E-Mail-Adresse ein, die mit Ihrem Konto verknüpft ist; wir senden einen Link zum Zurücksetzen, falls diese existiert.',
    email: 'E-Mail',
    sendLink: 'Link senden',
    sending: 'Wird gesendet...',
    successTitle: 'Prüfen Sie Ihre E-Mail',
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
    validatingToken: 'Link wird überprüft...',
    successMessage: 'Passwort wurde erfolgreich zurückgesetzt. Weiterleitung zur Anmeldung...',
    errorMessage: 'Fehler beim Zurücksetzen des Passworts. Der Link ist möglicherweise ungültig oder abgelaufen.',
    invalidTokenTitle: 'Ungültiger Link zum Zurücksetzen',
    invalidTokenMessage: 'Dieser Link zum Zurücksetzen des Passworts ist ungültig oder abgelaufen. Bitte fordere einen neuen Link an.'
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
    atLeastChars: 'Mindestens {n} Zeichen',
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
  },
  verifyEmail: {
    title: 'E-Mail-Verifizierung',
    verifying: 'Verifiziere Ihre E-Mail-Adresse...',
    success: 'E-Mail erfolgreich verifiziert!',
    successMessage: 'Ihre E-Mail-Adresse wurde erfolgreich verifiziert. Sie können sich jetzt anmelden.',
    error: 'Verifizierung fehlgeschlagen',
    missingToken: 'Verifizierungstoken fehlt',
    defaultError: 'Verifizierung fehlgeschlagen. Der Link ist möglicherweise abgelaufen oder ungültig.',
    goToLogin: 'Zur Anmeldung'
  },
  dashboard: {
    welcome: 'Dashboard',
    welcomeMessage: 'Willkommen im Dashboard',
    description: 'Hier finden Sie eine Übersicht über Ihr Konto und können verschiedene Aktionen durchführen.',
    plans: 'Pläne',
    customers: 'Kunden',
    users: 'Benutzer',
    emails: 'E-Mails',
    newsletters: 'Newsletter',
    activePlans: 'Aktive Pläne',
    totalCustomers: 'Gesamte Kunden',
    registeredUsers: 'Registrierte Benutzer',
    emailsSent: 'Gesendete E-Mails',
    subscribers: 'Abonnenten',
    recentPlans: 'Neueste Pläne',
    recentCustomers: 'Neueste Kunden',
    name: 'Name',
    status: 'Status',
    created: 'Erstellt',
    joined: 'Beigetreten',
    active: 'Aktiv',
    inactive: 'Inaktiv',
    email: 'E-Mail',
    loadError: 'Fehler beim Laden der Dashboard-Daten'
  },
  success: {
    goToLogin: 'Zur Anmeldung',
    registration: {
      title: 'Registrierung erfolgreich',
      message: 'Ihr Konto wurde erstellt. Bitte überprüfen Sie Ihre E-Mails, um Ihre E-Mail-Adresse zu bestätigen.'
    },
    passwordReset: {
      title: 'Passwort zurückgesetzt!',
      message: 'Ihr Passwort wurde zurückgesetzt. Sie können sich jetzt mit Ihrem neuen Passwort anmelden.'
    },
    passwordChange: {
      title: 'Passwort geändert!',
      message: 'Ihr Passwort wurde erfolgreich geändert.'
    },
    emailVerified: {
      title: 'E-Mail bestätigt!',
      message: 'Ihre E-Mail-Adresse wurde erfolgreich bestätigt. Sie können sich jetzt anmelden.'
    },
    default: {
      title: 'Erfolg!',
      message: 'Die Aktion wurde erfolgreich ausgeführt.'
    }
  },
  error: {
    goBack: 'Zurück',
    goHome: 'Zur Startseite',
    default: {
      title: 'Etwas ist schiefgelaufen',
      message: 'Es ist ein Fehler aufgetreten. Bitte versuchen Sie es erneut oder wenden Sie sich an den Support, falls das Problem weiterhin besteht.'
    }
  }
}
