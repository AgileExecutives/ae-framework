import { test, expect } from '@playwright/test';

const API_BASE_URL = 'http://localhost:8080/api/v1';

test.describe('Simple Form Tests', () => {
  test('should find registration form elements', async ({ page }) => {
    console.log('🔧 Testing registration form visibility');
    
    // Navigate to registration page
    await page.goto('/register');
    
    // Wait a moment for the page to load
    await page.waitForLoadState('networkidle');
    
    // Check if form elements are visible
    const firstnameField = page.locator('[data-testid="register-firstname"]');
    const lastnameField = page.locator('[data-testid="register-lastname"]');
    const emailField = page.locator('[data-testid="register-email"]');
    const passwordField = page.locator('[data-testid="register-password"]');
    const submitButton = page.locator('[data-testid="register-submit"]');
    
    await expect(firstnameField).toBeVisible();
    await expect(lastnameField).toBeVisible();
    await expect(emailField).toBeVisible();
    await expect(passwordField).toBeVisible();
    await expect(submitButton).toBeVisible();
    
    console.log('✅ All registration form elements are visible');
  });

  test('should find login form elements', async ({ page }) => {
    console.log('🔧 Testing login form visibility');
    
    // Navigate to login page
    await page.goto('/login');
    
    // Wait a moment for the page to load
    await page.waitForLoadState('networkidle');
    
    // Check if form elements are visible
    const emailField = page.locator('[data-testid="login-email"]');
    const passwordField = page.locator('[data-testid="login-password"]');
    const submitButton = page.locator('[data-testid="login-submit"]');
    
    await expect(emailField).toBeVisible();
    await expect(passwordField).toBeVisible();
    await expect(submitButton).toBeVisible();
    
    console.log('✅ All login form elements are visible');
  });

  test('should test basic form functionality', async ({ page, request }) => {
    console.log('🔧 Testing basic form functionality');
    
    // Create a unique test user first
    const timestamp = Date.now();
    const testUser = {
      email: `test_simple_${timestamp}@example.com`,
      password: 'TestPass123!',
      firstName: 'Test',
      lastName: 'User'
    };
    
    // Navigate to registration page
    await page.goto('/register');
    await page.waitForLoadState('networkidle');
    
    // Fill out the form
    await page.fill('[data-testid="register-firstname"]', testUser.firstName);
    await page.fill('[data-testid="register-lastname"]', testUser.lastName);
    await page.fill('[data-testid="register-email"]', testUser.email);
    await page.fill('[data-testid="register-password"]', testUser.password);
    
    console.log('✅ Form filled successfully');
    
    // Submit and see what happens
    await page.click('[data-testid="register-submit"]');
    
    // Wait a moment to see what happens
    await page.waitForTimeout(5000);
    
    // Check current URL
    const currentUrl = page.url();
    console.log('📍 Current URL after registration:', currentUrl);
    
    // Check if there are any error messages on the page
    const bodyText = await page.textContent('body');
    console.log('📄 Page content sample:', bodyText?.substring(0, 500));
    
    console.log('✅ Form submission test completed');
  });
});