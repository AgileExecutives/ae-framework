import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import { addMiddleware } from './router/middleware'
import AuthRestrictedView from './views/AuthRestrictedView.vue'

import { createI18n } from 'vue-i18n'
import './style.css'
import applyAppStyles from './plugins/appStyle'

// Import language files
import en from './locales/en'
import de from './locales/de'
import { addCatchAllRoute } from './router/shared-routes'
const messages = {
  en,
  de,
}
const i18n = createI18n({
  legacy: false,
  locale: 'de',
  messages,
})

const app = createApp(App)

app.use(createPinia())
addCatchAllRoute(router)
addMiddleware(router)
app.use(router)
app.use(i18n)

// Apply host-provided styles (inline or stylesheet) before mounting.
// Host apps can set `window.__CORE_FRONTEND_INLINE_STYLES__` or
// `window.__CORE_FRONTEND_STYLE_URL__` (or provide `/app-styles.css`).
applyAppStyles().finally(() => {
  app.mount('#app')
})
