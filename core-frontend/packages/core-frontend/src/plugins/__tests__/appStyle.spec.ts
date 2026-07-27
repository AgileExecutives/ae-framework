import { beforeEach, it, expect } from 'vitest'
import applyAppStyles from '../appStyle'

beforeEach(() => {
  document.head.innerHTML = ''
  ;(globalThis as any).__CORE_FRONTEND_INLINE_STYLES__ = undefined
  ;(globalThis as any).__CORE_FRONTEND_STYLE_URL__ = undefined
})

it('injects inline style when window var set', async () => {
  ;(globalThis as any).__CORE_FRONTEND_INLINE_STYLES__ = '.brand{color:#123456}'
  await applyAppStyles()
  const s = document.querySelector('style[data-core-frontend="inline-styles"]')
  expect(s).toBeTruthy()
  expect(s!.textContent).toContain('.brand')
})

it('injects link when url var set', async () => {
  ;(globalThis as any).__CORE_FRONTEND_STYLE_URL__ = '/app-styles.css'
  await applyAppStyles()
  const l = document.querySelector('link[data-core-frontend="app-style"]') as HTMLLinkElement
  expect(l).toBeTruthy()
  // href may be resolved in jsdom; check substring
  expect(l.href).toContain('/app-styles.css')
})

it('skips insertion when url is null', async () => {
  ;(globalThis as any).__CORE_FRONTEND_STYLE_URL__ = null
  await applyAppStyles()
  expect(document.querySelectorAll('link[data-core-frontend="app-style"]').length).toBe(0)
  expect(document.querySelectorAll('style[data-core-frontend="inline-styles"]').length).toBe(0)
})
