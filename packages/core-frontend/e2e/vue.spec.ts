import { test, expect } from '@playwright/test';

// Test configuration
const API_BASE_URL = 'http://localhost:8080/api/v1';

test.describe('Vue App E2E Tests', () => {
  test('visits the app root url without authentication', async ({ page }) => {
    // Root URL requires authentication, so should redirect to login
    await page.goto('/');
    await page.waitForLoadState('networkidle');
    
    // Should redirect to login page
    expect(page.url()).toContain('/login');
  });

  test('visits the app root url with authentication', async ({ page }) => {
    // Register user to get onboarding token
    const uniqueId = `${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
    const uniqueUsername = `testuser_${uniqueId}`;
    const uniqueEmail = `test_${uniqueId}@example.com`;
    
    const registerResponse = await page.request.post(`${API_BASE_URL}/auth/register`, {
      data: {
        username: uniqueUsername,
        email: uniqueEmail,
        password: 'TestPass123!',
        first_name: 'Test',
        last_name: 'User',
        company_name: `Test Company ${uniqueId}`,
        accept_terms: true
      }
    });
    
    if (registerResponse.status() === 429) {
      console.warn('⚠️ Rate limited, skipping test');
      test.skip();
      return;
    }
    
    expect(registerResponse.status()).toBe(201);
    const registerData = await registerResponse.json();
    const token = registerData.data.token;
    
    // Set token in localStorage
    await page.goto('/');
    await page.evaluate((tkn) => {
      localStorage.setItem('auth_token', tkn);
    }, token);
    
    // Now visit root URL with authentication
    await page.goto('/');
    await page.waitForLoadState('networkidle');
    
    // Should see the restricted page content
    await expect(page.locator('h1')).toContainText('Welcome to a restricted Page');
  });
});
