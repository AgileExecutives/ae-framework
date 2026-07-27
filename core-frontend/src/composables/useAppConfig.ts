import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { getAppConfig, getLogos, getLegal, getAppInfo } from '../config/app-config'

/**
 * Composable for accessing app configuration in Vue components
 */
export const useAppConfig = () => {
  const { locale } = useI18n()
  
  const config = computed(() => getAppConfig())
  const logos = computed(() => getLogos())
  const legal = computed(() => getLegal())
  const appInfo = computed(() => getAppInfo())
  
  // Get localized text with fallback
  const getLocalizedText = (textObj: { de?: string; en?: string } | undefined, fallback: string) => {
    if (!textObj) return fallback
    return textObj[locale.value as 'de' | 'en'] || textObj.en || textObj.de || fallback
  }
  
  // Computed properties for easy access
  const currentLogo = computed(() => {
    // For now, just return main logo - dark mode detection can be added later
    return logos.value.main
  })
  
  const currentLogoDark = computed(() => {
    return logos.value.mainDark
  })
  
  const currentSmallLogo = computed(() => {
    return logos.value.small || logos.value.main
  })
  
  const currentSmallLogoDark = computed(() => {
    return logos.value.smallDark || logos.value.mainDark
  })
  
  const termsOfServiceText = computed(() => 
    getLocalizedText(legal.value.termsOfService.text, 'Terms of Service')
  )
  
  const privacyPolicyText = computed(() => 
    getLocalizedText(legal.value.privacyPolicy.text, 'Privacy Policy')
  )
  
  const appName = computed(() => appInfo.value.name)
  
  const appDescription = computed(() => 
    getLocalizedText(appInfo.value.description, appInfo.value.name)
  )
  
  return {
    // Raw config access
    config,
    logos,
    legal,
    appInfo,
    
    // Convenient computed properties
    currentLogo,
    currentLogoDark,
    currentSmallLogo,
    currentSmallLogoDark,
    
    // Legal links
    termsOfServiceUrl: computed(() => legal.value.termsOfService.url),
    termsOfServiceText,
    privacyPolicyUrl: computed(() => legal.value.privacyPolicy.url), 
    privacyPolicyText,
    
    // App info
    appName,
    appDescription,
    
    // Utility function
    getLocalizedText
  }
}