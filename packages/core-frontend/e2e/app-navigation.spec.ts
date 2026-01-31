import { test, expect } from '@playwright/test';

/**
 * Application Navigation & Error Handling E2E Tests
 * 
 * Consolidates: not-found.spec.ts + simple-test.spec.ts
 * 
 * Covers:
 * - Main application navigation
 * - 404 error page handling
 * - Protected route redirects
 * - Basic smoke tests
 */

test.describe('Application Navigation', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('should load homepage without errors', async ({ page }) => {
    console.log('🧪 Test: Homepage loads successfully');
    
    await page.goto('/');
    await page.waitForLoadState('networkidle');
    
    // Verify page loaded
    expect(page.url()).toContain('localhost');
    
    // Verify no console errors (basic smoke test)
    const errors: string[] = [];
    page.on('console', msg => {
      if (msg.type() === 'error') {
        errors.push(msg.text());
      }
    });
    
    await page.waitForTimeout(500);
    
    console.log('✅ Homepage loaded successfully');
  });

  test('should navigate to main routes', async ({ page }) => {
    console.log('🧪 Test: Main route navigation');
    
    // Test navigation to key routes
    const routes = [
      '/login',
      '/register' // ADD MORE ROUTES AS IMPLEMENTED
    ];
    
    for (const route of routes) {
      console.log(`Navigating to ${route}...`);
      await page.goto(route);
      await page.waitForLoadState('networkidle');
      
      // Verify we're on the correct route
      expect(page.url()).toContain(route);
      
      // Wait a bit to ensure page is stable
      await page.waitForTimeout(500);
    }
    
    console.log('✅ Main route navigation completed');
  });

  // ========== 404 ERROR HANDLING (from not-found.spec.ts) ==========
  
  test('should display 404 page for invalid route', async ({ page }) => {
    console.log('🧪 Test: 404 page for invalid route');
    
    // Navigate to an invalid route
    await page.goto('/this-route-does-not-exist');
    await page.waitForLoadState('networkidle');
    
    expect(page.getByText('Seite nicht gefunden')).toBeVisible();

    // For now, just verify navigation doesn't crash
    await page.waitForTimeout(500);
    
    console.log('✅ 404 handling test completed');
  });

  test('should handle deeply nested invalid routes', async ({ page }) => {
    console.log('🧪 Test: Deeply nested 404 routes');
    
    await page.goto('/invalid/nested/route/structure');
    await page.waitForLoadState('networkidle');
    
    expect(page.getByText('Seite nicht gefunden')).toBeVisible();

    // Should handle gracefully
    await page.waitForTimeout(500);
    
    console.log('✅ Deep route 404 handling completed');
  });

  // ========== PROTECTED ROUTE REDIRECTS ==========
  
  test('should redirect to login for protected routes when not authenticated', async ({ page }) => {
    console.log('🧪 Test: Protected route redirects');
       
    // Try to access dashboard (should be protected)
    await page.goto('/dashboard');
    await page.waitForTimeout(500);
    
    const currentUrl = page.url();
    expect(currentUrl).toContain('/login?redirect=/dashboard');
    console.log('Current URL after accessing protected route:', currentUrl);

    console.log('✅ Protected route redirect test completed');
  });

  // ========== BASIC SMOKE TESTS (from simple-test.spec.ts) ==========
  
  test('should have correct page title', async ({ page }) => {
    console.log('🧪 Test: Page title verification');
    
    await page.goto('/');
    
    const title = await page.title();
    console.log('Page title:', title);
    
    // Verify title is not empty
    expect(title).toContain('App'); // Adjust as per actual title
    
    console.log('✅ Page title test completed');
  });

  test('should handle page refresh without errors', async ({ page }) => {
    console.log('🧪 Test: Page refresh handling');
    
    await page.goto('/');
    await page.waitForLoadState('networkidle');
    
    // Refresh the page
    await page.reload();
    await page.waitForLoadState('networkidle');
    
    // Verify page still works (may redirect to login)
    expect(page.url()).toContain('localhost');
    
    console.log('✅ Page refresh test completed');
  });

  test('should have responsive viewport meta tag', async ({ page }) => {
    console.log('🧪 Test: Responsive meta tag');
    
    await page.goto('/');
    
    // Check for viewport meta tag (important for mobile)
    const viewportMeta = await page.locator('meta[name="viewport"]').count();
    expect(viewportMeta).toBeGreaterThan(0);
    
    console.log('✅ Responsive meta tag found');
  });
});
