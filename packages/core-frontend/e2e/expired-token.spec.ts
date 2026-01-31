import { test, expect } from '@playwright/test';

// an old token
// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InJlc2V0X3Rlc3RfMTc2OTg4ODQyNDIwMEBleGFtcGxlLmNvbSIsImV4cCI6MTc2OTg5NTYyNiwibmJmIjoxNzY5ODg4NDI2LCJpYXQiOjE3Njk4ODg0MjYsImp0aSI6InJlc2V0X3Jlc2V0X3Rlc3RfMTc2OTg4ODQyNDIwMEBleGFtcGxlLmNvbV8xNzY5ODg4NDI2In0.YKhu90ywG_d_8bMQ2d_2nQoK4Ug019Qf4WyawGw1QrY

test.describe('Expired Token Handling', () => {
  test('should properly redirect to login when token expires', async ({ page }) => {
    console.log('🧪 Testing expired token redirect behavior');
    
    // Set an invalid/expired token in localStorage
    await page.goto('/');
    await page.evaluate(() => {
      localStorage.setItem('auth_token', 'invalid_or_expired_token_xyz123');
    });
    
    console.log('✅ Set invalid token in localStorage');
    
    // Try to navigate to protected route
    await page.goto('/');
    
    // Wait for navigation to complete
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(1000); // Give time for auth check
    
    console.log('📍 Current URL:', page.url());
    
    // Should redirect to login page
    expect(page.url()).toContain('/login');
    
    // Login form should be visible
    const loginEmailInput = page.locator('[data-testid="login-email"]');
    const loginPasswordInput = page.locator('[data-testid="login-password"]');
    const loginSubmitButton = page.locator('[data-testid="login-submit"]');
    
    await expect(loginEmailInput).toBeVisible({ timeout: 5000 });
    await expect(loginPasswordInput).toBeVisible({ timeout: 5000 });
    await expect(loginSubmitButton).toBeVisible({ timeout: 5000 });
    
    console.log('✅ Login form is properly displayed after expired token redirect');
  });
  
  test('should show login form immediately on direct navigation', async ({ page }) => {
    console.log('🧪 Testing direct login page navigation');
    
    // Navigate directly to login
    await page.goto('/login');
    await page.waitForLoadState('networkidle');
    
    console.log('📍 Current URL:', page.url());
    
    // Login form should be immediately visible
    const loginEmailInput = page.locator('[data-testid="login-email"]');
    const loginPasswordInput = page.locator('[data-testid="login-password"]');
    const loginSubmitButton = page.locator('[data-testid="login-submit"]');
    
    await expect(loginEmailInput).toBeVisible({ timeout: 5000 });
    await expect(loginPasswordInput).toBeVisible({ timeout: 5000 });
    await expect(loginSubmitButton).toBeVisible({ timeout: 5000 });
    
    console.log('✅ Login form is properly displayed on direct navigation');
  });
});
