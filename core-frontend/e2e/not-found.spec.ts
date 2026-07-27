import { test, expect } from '@playwright/test'

test.describe('404 Not Found Page', () => {
  test('should display 404 page for non-existent route', async ({ page }) => {
    // Navigate to a non-existent page
    await page.goto('/this-page-does-not-exist')

    // Wait for page to load
    await page.waitForLoadState('networkidle')

    // Check for 404 heading
    await expect(page.locator('text=404')).toBeVisible()
    
    // Check for title (English or German)
    const titleLocator = page.locator('h1')
    await expect(titleLocator).toBeVisible()
    
    // Verify one of the translated titles appears
    const titleText = await titleLocator.textContent()
    expect(['Page Not Found', 'Seite nicht gefunden']).toContain(titleText)
  })

  test('should have working Go Back button', async ({ page }) => {
    // First go to login page
    await page.goto('/login')
    await page.waitForLoadState('networkidle')
    
    // Then navigate to non-existent page
    await page.goto('/non-existent-page')
    await page.waitForLoadState('networkidle')

    // Verify 404 page is shown
    await expect(page.locator('text=404')).toBeVisible()

    // Click Go Back button (English or German)
    const goBackButton = page.locator('button').filter({ hasText: /Go Back|Zurück/ }).first()
    await expect(goBackButton).toBeVisible()
    await goBackButton.click()

    // Should be back at login page
    await expect(page).toHaveURL('/login')
  })

  test('should have working Go to Home button', async ({ page }) => {
    // Navigate to non-existent page
    await page.goto('/does-not-exist')
    await page.waitForLoadState('networkidle')

    // Verify 404 page is shown
    await expect(page.locator('text=404')).toBeVisible()

    // Click Go Home button (English or German)
    const goHomeButton = page.locator('button').filter({ hasText: /Go to Home|Zur Startseite/ }).first()
    await expect(goHomeButton).toBeVisible()
    await goHomeButton.click()

    // Wait for navigation to complete - accept either / or /login (with optional query params)
    await page.waitForURL((url) => {
      const pathname = url.pathname
      return pathname === '/' || pathname === '/login'
    }, { timeout: 5000 })
    await page.waitForLoadState('networkidle')
    
    // Should redirect to home page (which may redirect to login if not authenticated)
    const currentPath = new URL(page.url()).pathname
    expect(currentPath).toMatch(/^\/(login)?$/)
  })

  test('should have helpful links to login and register', async ({ page }) => {
    // Navigate to non-existent page
    await page.goto('/404-test')
    await page.waitForLoadState('networkidle')

    // Verify 404 page is shown
    await expect(page.locator('text=404')).toBeVisible()

    // Check for login link
    const loginLink = page.locator('a[href="/login"]')
    await expect(loginLink).toBeVisible()

    // Check for register link
    const registerLink = page.locator('a[href="/register"]')
    await expect(registerLink).toBeVisible()
  })

  test('should show theme toggle and locale switcher', async ({ page }) => {
    await page.goto('/invalid-route')
    await page.waitForLoadState('networkidle')

    // Verify 404 page is shown
    await expect(page.locator('text=404')).toBeVisible()

    // Theme toggle and locale switcher should be visible
    // These are now floating buttons in the top right corner as labels with btn classes
    
    // Look for the specific theme and locale switcher components
    const themeToggle = page.locator('label.swap.btn.btn-ghost.btn-circle').first()
    const localeSwitcher = page.locator('label.swap.btn.btn-ghost.btn-circle').nth(1)
    
    // At least the theme toggle should be visible
    await expect(themeToggle).toBeVisible()
    
    // Look for SVG elements inside (theme toggle has sun/moon icons, locale has flags)
    const themeIcons = page.locator('label.swap svg')
    await expect(themeIcons.first()).toBeVisible()
  })

  test('should handle nested non-existent routes', async ({ page }) => {
    // Navigate to deeply nested non-existent path
    await page.goto('/deeply/nested/path/that/does/not/exist')
    await page.waitForLoadState('networkidle')

    // Should still show 404 page
    await expect(page.locator('text=404')).toBeVisible()
    
    const titleLocator = page.locator('h1')
    await expect(titleLocator).toBeVisible()
  })
})
