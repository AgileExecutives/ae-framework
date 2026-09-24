// Minimal runtime style loader used by tests and host apps.
export default async function applyAppStyles(): Promise<void> {
  const styles = (globalThis as any).__CORE_FRONTEND_INLINE_STYLES__ || ''
  if (!styles) return
  const el = document.createElement('style')
  el.setAttribute('data-core-frontend', 'inline-styles')
  el.textContent = styles
  document.head.appendChild(el)
}
