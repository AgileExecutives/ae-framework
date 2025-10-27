import { test, expect } from './fixtures/auth-fixtures';

// Example test showing how to use the auth fixtures
// These are examples only - skip them in actual test runs

test.describe.skip('Example: Using Auth Fixtures', () => {
  test('accessing authenticated route with fixture user', async ({ page, authenticatedUser }) => {
    console.log(`🧪 Using pre-created user: ${authenticatedUser.username}`);
    
    // Set auth token from the pool user
    await page.goto('/');
    await page.evaluate((token) => {
      localStorage.setItem('auth_token', token);
    }, authenticatedUser.token);
    
    // Navigate to change password page (requires authentication)
    await page.goto('/change-password');
    await page.waitForLoadState('networkidle');
    
    // Should see the authenticated page, not be redirected to login
    await expect(page).toHaveURL('/change-password');
    await expect(page.getByRole('heading', { name: /passwort ändern/i })).toBeVisible();
    
    console.log('✅ Successfully accessed authenticated route with fixture user');
  });
  
  test('multiple users from pool', async ({ page, authUserPool }) => {
    console.log(`🧪 User pool has ${authUserPool.length} users available`);
    
    // You can access any user from the pool
    const user1 = authUserPool[0];
    const user2 = authUserPool[1];
    
    console.log(`User 1: ${user1.username}`);
    console.log(`User 2: ${user2.username}`);
    
    // Use whichever user you need
    await page.goto('/');
    await page.evaluate((token) => {
      localStorage.setItem('auth_token', token);
    }, user1.token);
    
    await page.goto('/change-password');
    await expect(page).toHaveURL('/change-password');
    
    console.log('✅ Successfully used user from pool');
  });
  
  test('example without authentication fixture', async ({ page }) => {
    // Tests that don't need auth don't have to use the fixtures
    await page.goto('/login');
    await expect(page.getByRole('heading', { name: /anmelden/i })).toBeVisible();
    
    console.log('✅ Non-authenticated test works normally');
  });
});
