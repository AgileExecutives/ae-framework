import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue'

type ThemePreference = 'dark' | 'light' | 'system'
const STORAGE_KEY = 'theme'

// Returns true if the OS/Browser prefers dark mode right now
function systemPrefersDark(): boolean {
  return typeof window !== 'undefined' && window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches
}

// Apply theme to the root element (used by DaisyUI)
function applyTheme(dark: boolean) {
  if (typeof document === 'undefined') return
  const root = document.documentElement
  
  // DaisyUI uses data-theme attribute
  const theme = dark ? 'dark' : 'light'
  root.setAttribute('data-theme', theme)
  
  // Also add/remove dark class for Tailwind dark: variants
  if (dark) root.classList.add('dark')
  else root.classList.remove('dark')
}

// Global state shared across all instances
const _isDark = ref<boolean>(false)
const pref = ref<ThemePreference>('system')
let mq: MediaQueryList | null = null
let mqListener: ((e: MediaQueryListEvent) => void) | null = null
let isInitialized = false

export function useDarkMode() {

  function readStoredPref(): ThemePreference | null {
    try {
      const v = localStorage.getItem(STORAGE_KEY)
      if (v === 'dark' || v === 'light' || v === 'system') return v
      return null
    } catch {
      return null
    }
  }

  function updateFromSources() {
    // stored preference wins if set, otherwise follow system
    const stored = readStoredPref()
    if (stored) pref.value = stored
    else pref.value = 'system'

    if (pref.value === 'system') _isDark.value = systemPrefersDark()
    else _isDark.value = pref.value === 'dark'

    applyTheme(_isDark.value)
  }

  function setPreference(p: ThemePreference) {
    try {
      if (p === 'system') localStorage.removeItem(STORAGE_KEY)
      else localStorage.setItem(STORAGE_KEY, p)
    } catch {}
    updateFromSources()
  }

  function toggle() {
    setPreference(_isDark.value ? 'light' : 'dark')
  }

  onMounted(() => {
    // Only initialize once globally
    if (!isInitialized) {
      isInitialized = true
      updateFromSources()

      if (typeof window !== 'undefined' && window.matchMedia) {
        mq = window.matchMedia('(prefers-color-scheme: dark)')
        mqListener = (e: MediaQueryListEvent) => {
          // Only update if user hasn't explicitly set a preference
          try {
            const stored = localStorage.getItem(STORAGE_KEY)
            if (!stored || stored === 'system') {
              _isDark.value = e.matches
              applyTheme(_isDark.value)
            }
          } catch {
            _isDark.value = e.matches
            applyTheme(_isDark.value)
          }
        }
        // Modern browsers support addEventListener on MediaQueryList
        if (mq.addEventListener) mq.addEventListener('change', mqListener)
        else mq.addListener(mqListener as any)
      }
    }
  })

  onBeforeUnmount(() => {
    if (mq && mqListener) {
      if (mq.removeEventListener) mq.removeEventListener('change', mqListener)
      else mq.removeListener(mqListener as any)
    }
  })

  // Keep theme in sync when _isDark changes
  watch(_isDark, (v) => applyTheme(v))

  const isDark = computed(() => _isDark.value)

  return {
    // reactive read-only
    isDark,
    // low-level control
    setPreference, // setPreference('dark'|'light'|'system')
    toggle,
    // convenience
    STORAGE_KEY,
  }
}
