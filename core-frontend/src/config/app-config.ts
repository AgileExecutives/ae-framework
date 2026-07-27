/**
 * App Configuration
 * 
 * This file contains configurable elements that differ between 
 * core-frontend and unburdy-app implementations.
 */

// Logo imports - these should be overridden in each app
import defaultLogo from '../assets/logo-light.svg'
import defaultLogoDark from '../assets/logo-dark.svg'

export interface AppConfig {
  // Logos
  logos: {
    // Main logos (large size - used in forms, headers)
    main: string
    mainDark: string
    // Small logos (used in compact spaces, mobile)
    small?: string
    smallDark?: string
  }
  
  // Legal links
  legal: {
    termsOfService: {
      url: string
      text?: {
        de?: string
        en?: string
      }
    }
    privacyPolicy: {
      url: string
      text?: {
        de?: string
        en?: string
      }
    }
  }
  
  // App specific settings
  app: {
    name: string
    description?: {
      de?: string
      en?: string
    }
  }
}

// Default configuration for core-frontend
export const defaultAppConfig: AppConfig = {
  logos: {
    main: defaultLogo,
    mainDark: defaultLogoDark,
    // If no small logos provided, use main logos
    small: defaultLogo,
    smallDark: defaultLogoDark
  },
  
  legal: {
    termsOfService: {
      url: '/terms',
      text: {
        de: 'AGB',
        en: 'Terms of Service'
      }
    },
    privacyPolicy: {
      url: '/privacy',
      text: {
        de: 'Datenschutzerklärung', 
        en: 'Privacy Policy'
      }
    }
  },
  
  app: {
    name: 'Core Frontend',
    description: {
      de: 'Kern Frontend Anwendung',
      en: 'Core Frontend Application'
    }
  }
}

// Configuration override mechanism
let currentConfig = defaultAppConfig

export const setAppConfig = (config: Partial<AppConfig>) => {
  currentConfig = {
    ...defaultAppConfig,
    ...config,
    logos: {
      ...defaultAppConfig.logos,
      ...config.logos
    },
    legal: {
      termsOfService: {
        ...defaultAppConfig.legal.termsOfService,
        ...config.legal?.termsOfService
      },
      privacyPolicy: {
        ...defaultAppConfig.legal.privacyPolicy,
        ...config.legal?.privacyPolicy
      }
    },
    app: {
      ...defaultAppConfig.app,
      ...config.app
    }
  }
}

export const getAppConfig = (): AppConfig => currentConfig

// Convenience functions
export const getLogos = () => currentConfig.logos
export const getLegal = () => currentConfig.legal
export const getAppInfo = () => currentConfig.app