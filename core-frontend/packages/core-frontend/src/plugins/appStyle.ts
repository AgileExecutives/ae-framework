export async function applyAppStyles(): Promise<void> {
  try {
    const win = window as any;
    if (typeof document === 'undefined') return;

    // Inline CSS string provided by host application
    if (win.__CORE_FRONTEND_INLINE_STYLES__) {
      const existing = document.querySelector('style[data-core-frontend="inline-styles"]');
      if (!existing) {
        const styleEl = document.createElement('style');
        styleEl.setAttribute('data-core-frontend', 'inline-styles');
        styleEl.textContent = win.__CORE_FRONTEND_INLINE_STYLES__;
        document.head.appendChild(styleEl);
      }
      return;
    }

    // URL to a stylesheet provided by host, or fallback to /app-styles.css
    const rawUrl = win.__CORE_FRONTEND_STYLE_URL__;

    // If the host explicitly set `null` or `false`, skip automatic insertion
    if (rawUrl === null || rawUrl === false) return;

    // If host provided a value (including empty string), use it. Otherwise fallback.
    const url = rawUrl !== undefined ? rawUrl : '/app-styles.css';

    // Don't duplicate insertion
    if (!document.querySelector(`link[data-core-frontend="app-style"][href="${url}"]`)) {
      const link = document.createElement('link');
      link.rel = 'stylesheet';
      link.href = url;
      link.setAttribute('data-core-frontend', 'app-style');
      document.head.appendChild(link);
    }
  } catch (err) {
    // eslint-disable-next-line no-console
    console.warn('core-frontend: applyAppStyles failed', err);
  }
}

export default applyAppStyles;
