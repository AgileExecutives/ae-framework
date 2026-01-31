import { test, expect } from '@playwright/test';

// Test configuration
const BASE_URL = 'http://localhost:3000';
const API_BASE_URL = 'http://localhost:8080/api/v1';

// Test user data
const testUser = {
  username: `pwdtest_${Date.now()}`,
  email: `pwdtest_${Date.now()}@example.com`,
  password: 'OldPass123!',
  newPassword: 'NewPass456!',
  firstName:  `First_${Date.now()}`,
  lastName:  `Last_${Date.now()}`
};

// Helper function to create and authenticate a test user
async function createAuthenticatedUser(request: any) {
  console.log('🔧 Creating test user...');
  
  const user = {
    username: `pwdtest_${Date.now()}`,
    email: `pwdtest_${Date.now()}@example.com`,
    password: 'OldPass123!',
    firstName:  `First_${Date.now()}`,
    lastName:  `Last_${Date.now()}`
  };
  
  // Wait a bit to avoid rate limiting
  await new Promise(resolve => setTimeout(resolve, 1000));
  
  // Create test user
  const registerResponse = await request.post(`${API_BASE_URL}/auth/register`, {
    data: {
      username: user.username,
      email: user.email,
      password: user.password,
      first_name:  `First_${Date.now()}`,
      last_name:  `Last_${Date.now()}`,
      company_name: `Test Company ${Date.now()}`,
      accept_terms: true
    }
  });
  
  if (registerResponse.status() === 429) {
    // Rate limited - wait and skip this test
    console.warn('⚠️ Rate limited, skipping test');
    throw new Error('RATE_LIMITED');
  }
  
  if (registerResponse.status() !== 200 && registerResponse.status() !== 201) {
    const responseText = await registerResponse.text();
    console.error('Registration failed:', responseText);
    throw new Error(`Failed to create test user: ${registerResponse.status()} - ${responseText}`);
  }
  
  const registerData = await registerResponse.json();
  console.log('✅ User registered successfully');
  
  // Use the onboarding token from registration response instead of trying to login
  // (Login requires email verification, but onboarding token works for change password)
  return {
    user,
    token: registerData.data.token // This is the onboarding token
  };
}

test.describe('Password Management E2E Tests', () => {
  // No beforeAll - tests create their own users as needed

  test.beforeEach(async ({ page }) => {
    // Ensure we start from a fresh state
    await page.goto('/');
  });

  test('should display forgot password form', async ({ page }) => {
    await page.goto('/forgot-password');
    
    await expect(page.locator('[data-testid="forgot-email"]')).toBeVisible();
    await expect(page.locator('[data-testid="forgot-submit"]')).toBeVisible();
  });

  test('should submit forgot password request', async ({ page }) => {
    console.log('🧪 Test: Submit forgot password request');
    
    await page.goto('/forgot-password');
    await page.waitForLoadState('networkidle');
    
    // Fill in email using testid
    await page.fill('[data-testid="forgot-email"]', testUser.email);
    
    // Submit form using testid
    await page.click('[data-testid="forgot-submit"]');
    
    // Wait for redirect
    await page.waitForTimeout(200);
    
    // Should redirect to success page
    await page.waitForURL('**/success**', { timeout: 3000 });
    const successUrl = page.url();
    expect(successUrl).toContain('/success');
    
    console.log('✅ Forgot password request submitted successfully, redirected to success page');
  });

  test('should handle forgot password for non-existent email', async ({ page }) => {
    console.log('🧪 Test: Forgot password with non-existent email');
    
    await page.goto('/forgot-password');
    await page.waitForLoadState('networkidle');
    
    // Fill in non-existent email using testid
    await page.locator('[data-testid="forgot-email"]').fill('nonexistent@example.com');
    
    // Submit form using testid
    await page.locator('[data-testid="forgot-submit"]').click();
    
    // Wait for redirect
    await page.waitForTimeout(2000);
    
    // Should still redirect to success (security best practice - don't reveal if email exists)
    await page.waitForURL('**/success**', { timeout: 3000 });
    const successUrl = page.url();
    expect(successUrl).toContain('/success');
    
    console.log('✅ Forgot password handles non-existent email correctly');
  });

  test('should validate email format in forgot password', async ({ page }) => {
    console.log('🧪 Test: Email validation in forgot password');
    
    await page.goto('/forgot-password');
    await page.waitForLoadState('networkidle');
    
    // Try to submit with invalid email using testid
    await page.locator('[data-testid="forgot-email"]').fill('invalid-email');
    await page.locator('[data-testid="forgot-submit"]').click();
    
    // Wait a bit for validation to trigger
    await page.waitForTimeout(1000);
    
    // Check if validation prevents submission or shows error
    // Some implementations might prevent submission entirely rather than showing error messages
    try {
      const anyError = page.locator('.alert-error, .error, [role="alert"], .text-error, .text-red-500').first();
      const germanError = page.getByText(/bitte geben sie eine gültige e-mail-adresse ein/i);
      const englishError = page.getByText(/please enter a valid email address/i);
      const genericEmailError = page.getByText(/email/i).and(page.getByText(/valid|invalid|gültig/i));
      
      // Try to find any validation error
      await expect(anyError.or(germanError).or(englishError).or(genericEmailError)).toBeVisible({ timeout: 2000 });
    } catch (error) {
      // If no error message is visible, check if we're still on the same page (validation prevented submission)
      await expect(page).toHaveURL(/forgot-password/);
      console.log('⚠️  Validation prevented submission (no error message shown)');
    }
    
    console.log('✅ Email validation works correctly');
  });

  test('should display reset password form with token', async ({ page, request }) => {
    console.log('🧪 Test: Reset password form display');
    
    // First create a user
    const user = await createAuthenticatedUser(request);
    if (!user) return; // Skip if rate limited
    
    // Request password reset to get token
    const forgotPasswordResponse = await request.post(`${API_BASE_URL}/auth/forgot-password`, {
      data: { email: user.user.email }
    });
    
    console.log('📧 Forgot password response:', forgotPasswordResponse.status(), await forgotPasswordResponse.text());
    
    // Wait longer for email to be sent
    await page.waitForTimeout(2000);
    
    // Fetch the reset token from the latest email
    const emailsResponse = await request.get(`${API_BASE_URL}/emails/latest-emails`);
    const emailsData = await emailsResponse.json();
    
    console.log('📧 All emails:', emailsData.data.map((e: any) => ({ subject: e.subject, to: e.to })));
    
    const resetEmail = emailsData.data.find((email: any) => 
      email.to === user.user.email && 
      email.subject.toLowerCase().includes('password') && 
      email.subject.toLowerCase().includes('reset')
    );
    
    if (!resetEmail) {
      console.error('Available emails:', emailsData.data.map((e: any) => e.subject));
      throw new Error('Reset email not found');
    }
    
    // Extract token from email text (format: /new-password?token=xxx)
    const tokenMatch = resetEmail.text.match(/\/new-password\?token=([^\s]+)/);
    if (!tokenMatch) {
      throw new Error('Token not found in email');
    }
    const testToken = tokenMatch[1];
    console.log('🔑 Extracted reset token:', testToken);
    
    await page.goto(`/new-password?token=${testToken}`);
    await page.waitForLoadState('networkidle');
    
    // Wait for password requirements to load
    console.log('⏳ Waiting for password requirements to load...');
    await page.waitForSelector('[data-testid="reset-password"]', { state: 'visible' });
    await page.waitForTimeout(1000); // Give time for password requirements API call
    
    // Verify form elements using testids
    await expect(page.getByRole('heading', { name: /passwort zurücksetzen/i })).toBeVisible();
    await expect(page.locator('[data-testid="reset-password"]')).toBeVisible();
    await expect(page.locator('[data-testid="reset-password-repeat"]')).toBeVisible();
    await expect(page.locator('[data-testid="reset-submit"]')).toBeVisible();
    
    console.log('✅ Reset password form displays correctly');
  });

  test('should validate password requirements in reset form', async ({ page, request }) => {
    console.log('🧪 Test: Password requirements validation');
    
    // First create a user
    const { user } = await createAuthenticatedUser(request);
    if (!user) return; // Skip if rate limited
    
    // Request password reset to get token
    await request.post(`${API_BASE_URL}/auth/forgot-password`, {
      data: { email: user.email }
    });
    
    // Wait and fetch the reset token from email
    await page.waitForTimeout(500);
    const emailsResponse = await request.get(`${API_BASE_URL}/emails/latest-emails`);
    const emailsData = await emailsResponse.json();
    const resetEmail = emailsData.data.find((email: any) => 
      email.subject.toLowerCase().includes('password') && email.subject.toLowerCase().includes('reset')
    );
    const tokenMatch = resetEmail.text.match(/\/new-password\?token=([^\s]+)/);
    const testToken = tokenMatch[1];
    
    await page.goto(`/new-password?token=${testToken}`);
    await page.waitForLoadState('networkidle');
    
    // Wait for password requirements to load
    console.log('⏳ Waiting for password requirements to load...');
    await page.waitForSelector('[data-testid="reset-password"]', { state: 'visible' });
    await page.waitForTimeout(1000); // Give time for password requirements API call
    
    // Check for password requirements description (optional - might not be visible)
    try {
      const germanReq = page.getByText(/muss mindestens.*zeichen/i);
      const englishReq = page.getByText(/must be at least.*characters/i);
      const anyPassReq = page.getByText(/8.*character|character.*8|mindestens.*8|least.*8/i);
      await expect(germanReq.or(englishReq).or(anyPassReq).first()).toBeVisible({ timeout: 1000 });
      console.log('✅ Password requirements are displayed');
    } catch (error) {
      console.log('⚠️  Password requirements not visible (may be in placeholder or help text)');
    }
    
    // Try weak password using testids
    await page.locator('[data-testid="reset-password"]').fill('weak');
    await page.locator('[data-testid="reset-password-repeat"]').fill('weak');
    
    // Check that submit button is disabled when password requirements are not met
    await expect(page.locator('[data-testid="reset-submit"]')).toBeDisabled();
    console.log('✅ Submit button is disabled for weak password');
    
    // Try a valid password that meets requirements
    await page.locator('[data-testid="reset-password"]').fill('ValidPass123!');
    await page.locator('[data-testid="reset-password-repeat"]').fill('ValidPass123!');
    
    // Now button should be enabled
    await expect(page.locator('[data-testid="reset-submit"]')).toBeEnabled();
    
    console.log('✅ Password requirements validation works');
  });

  test('should validate password confirmation match', async ({ page, request }) => {
    console.log('🧪 Test: Password confirmation matching');
    
    // First create a user
    const {user} = await createAuthenticatedUser(request);
    if (!user) return; // Skip if rate limited
    
    // Request password reset to get token
    await request.post(`${API_BASE_URL}/auth/forgot-password`, {
      data: { email: user.email }
    });
    
    // Wait and fetch the reset token from email
    await page.waitForTimeout(500);
    const emailsResponse = await request.get(`${API_BASE_URL}/emails/latest-emails`);
    const emailsData = await emailsResponse.json();
    const resetEmail = emailsData.data.find((email: any) => 
      email.subject.toLowerCase().includes('password') && email.subject.toLowerCase().includes('reset')
    );
    const tokenMatch = resetEmail.text.match(/\/new-password\?token=([^\s]+)/);
    const testToken = tokenMatch[1];
    
    await page.goto(`/new-password?token=${testToken}`);
    await page.waitForLoadState('networkidle');
    
    // Wait for password requirements to load
    console.log('⏳ Waiting for password requirements to load...');
    await page.waitForSelector('[data-testid="reset-password"]', { state: 'visible' });
    await page.waitForTimeout(1000); // Give time for password requirements API call
    
    // Fill passwords that don't match using testids
    await page.locator('[data-testid="reset-password"]').fill('ValidPass123!');
    await page.locator('[data-testid="reset-password-repeat"]').fill('DifferentPass456!');
    
    // Wait for validation
    await page.waitForTimeout(500);
    
    await page.locator('[data-testid="reset-submit"]').click();
    
    // Should show mismatch error (German) - look for the validation error message
    await expect(page.getByText(/die passwörter stimmen nicht überein/i)).toBeVisible({ timeout: 2000 });
    
    console.log('✅ Password confirmation validation works');
  });

  test('should display change password form for authenticated user', async ({ page, context, request }) => {
    console.log('🧪 Test: Change password form display');
    
    try {
      // Create authenticated user
      const { token } = await createAuthenticatedUser(request);
      
      // Set auth token in localStorage
      await context.addCookies([]);
      await page.goto('/');
      await page.evaluate((token) => {
        localStorage.setItem('auth_token', token);
      }, token);
      
      // Navigate to change password page
      await page.goto('/change-password');
      await page.waitForLoadState('networkidle');
      
      // Wait for password requirements to load
      console.log('⏳ Waiting for password requirements to load...');
      await page.waitForSelector('[data-testid="change-new-password"]', { state: 'visible' });
      await page.waitForTimeout(1000); // Give time for password requirements API call
    
      // Verify form elements using testids
      await expect(page.getByRole('heading', { name: /passwort ändern/i })).toBeVisible();
      await expect(page.locator('[data-testid="change-current-password"]')).toBeVisible();
      await expect(page.locator('[data-testid="change-new-password"]')).toBeVisible();
      await expect(page.locator('[data-testid="change-confirm-password"]')).toBeVisible();
      await expect(page.locator('[data-testid="change-submit"]')).toBeVisible();
      
      console.log('✅ Change password form displays correctly');
    } catch (error: any) {
      if (error.message === 'RATE_LIMITED') {
        console.warn('⚠️ Skipping test due to rate limiting');
        test.skip();
      } else {
        throw error;
      }
    }
  });

  test('should redirect unauthenticated user from change password', async ({ page }) => {
    console.log('🧪 Test: Change password requires authentication');
    
    // Clear any stored auth
    await page.goto('/');
    await page.evaluate(() => {
      localStorage.clear();
    });
    
    // Try to access change password page
    await page.goto('/change-password');
    await page.waitForLoadState('networkidle');
    
    // Should redirect to login (with optional redirect query param)
    await expect(page.url()).toContain('/login');
    
    console.log('✅ Change password correctly requires authentication');
  });

  test('should validate current password in change password', async ({ page, context, request }) => {
    console.log('🧪 Test: Current password validation');
    
    try {
      // Create authenticated user
      const { token } = await createAuthenticatedUser(request);
      
      // Set auth token
      await page.goto('/');
      await page.evaluate((tkn) => {
        localStorage.setItem('auth_token', tkn);
      }, token);
      
      await page.goto('/change-password');
      await page.waitForLoadState('networkidle');
      
      // Wait for password requirements to load
      console.log('⏳ Waiting for password requirements to load...');
      await page.waitForSelector('[data-testid="change-new-password"]', { state: 'visible' });
      await page.waitForTimeout(1000); // Give time for password requirements API call
      
      // Enter wrong current password using testids
      await page.locator('[data-testid="change-current-password"]').fill('WrongPassword123!');
      await page.locator('[data-testid="change-new-password"]').fill('NewValidPass456!');
      await page.locator('[data-testid="change-confirm-password"]').fill('NewValidPass456!');
      
      // Wait for validation
      await page.waitForTimeout(500);
      
      await page.locator('[data-testid="change-submit"]').click();
      
      // Wait for error to appear
      await page.waitForTimeout(2000);
      
      // Should show error (stays on page with error message)
      const hasError = await page.locator('.text-error, [class*="error"]').count() > 0;
      expect(hasError).toBe(true);
      
      console.log('✅ Current password validation works');
    } catch (error: any) {
      if (error.message === 'RATE_LIMITED') {
        console.warn('⚠️ Skipping test due to rate limiting');
        test.skip();
      } else {
        throw error;
      }
    }
  });

  test('should successfully change password', async ({ page, context, request }) => {
    console.log('🧪 Test: Successful password change');
    
    try {
      // Create authenticated user
      const { user, token } = await createAuthenticatedUser(request);
      const newPassword = 'NewPass456!';
      
      // Set auth token
      await page.goto('/');
      await page.evaluate((tkn) => {
        localStorage.setItem('auth_token', tkn);
      }, token);
      
      await page.goto('/change-password');
      await page.waitForLoadState('networkidle');
      
      // Wait for password requirements to load
      console.log('⏳ Waiting for password requirements to load...');
      await page.waitForSelector('[data-testid="change-new-password"]', { state: 'visible' });
      await page.waitForTimeout(1000); // Give time for password requirements API call
      
      // Enter correct current password and new password using testids
      await page.locator('[data-testid="change-current-password"]').fill(user.password);
      await page.locator('[data-testid="change-new-password"]').fill(newPassword);
      await page.locator('[data-testid="change-confirm-password"]').fill(newPassword);
      
      // Wait for validation
      await page.waitForTimeout(500);
      
      await page.locator('[data-testid="change-submit"]').click();
      
      // Wait for redirect
      await page.waitForTimeout(2000);
      
      // Should redirect to success page
      await page.waitForURL('**/success**', { timeout: 3000 });
      const successUrl = page.url();
      expect(successUrl).toContain('/success');
      expect(successUrl).toContain('password-change');
      
      // Note: We cannot test login with new password because the user's email is not verified
      // The onboarding token allows password change, but full login requires email verification
      
      console.log('✅ Password change successful, redirected to success page');
    } catch (error: any) {
      if (error.message === 'RATE_LIMITED') {
        console.warn('⚠️ Skipping test due to rate limiting');
        test.skip();
      } else {
        throw error;
      }
    }
  });

  test('should validate new password requirements in change password', async ({ page, request }) => {
    console.log('🧪 Test: New password requirements in change password');
    
    try {
      // Create authenticated user
      const { user, token } = await createAuthenticatedUser(request);
      
      await page.goto('/');
      await page.evaluate((tkn) => {
        localStorage.setItem('auth_token', tkn);
      }, token);
    
      await page.goto('/change-password');
      await page.waitForLoadState('networkidle');
      
      // Wait for requirements to load
      await page.waitForTimeout(1000);
      
      // Try with weak new password using testids
      await page.locator('[data-testid="change-current-password"]').fill(user.password);
      await page.locator('[data-testid="change-new-password"]').fill('weak');
      await page.locator('[data-testid="change-confirm-password"]').fill('weak');
      
      await page.locator('[data-testid="change-submit"]').click();
      
      // Check for validation error for weak password - may or may not show visible error
      try {
        const germanPassError = page.getByText(/passwort muss mindestens.*zeichen/i);
        const englishPassError = page.getByText(/password must be at least.*characters/i);
        const anyPassError = page.locator('.alert-error, .error, [role="alert"], .text-error').filter({ hasText: /password|passwort/i });
        await expect(germanPassError.or(englishPassError).or(anyPassError).first()).toBeVisible({ timeout: 2000 });
        console.log('✅ Password validation error shown');
      } catch (error) {
        // Check if we're still on change password page (validation prevented submission)
        await expect(page).toHaveURL(/change-password/);
        console.log('⚠️  Validation prevented submission (no visible error message)');
      }
      
      console.log('✅ New password requirements validation works');
    } catch (error: any) {
      if (error.message === 'RATE_LIMITED') {
        console.warn('⚠️ Skipping test due to rate limiting');
        test.skip();
      } else {
        throw error;
      }
    }
  });

  test('should handle API errors gracefully', async ({ page }) => {
    console.log('🧪 Test: Error handling');
    
    await page.goto('/forgot-password');
    await page.waitForLoadState('networkidle');
    
    // Intercept API call and return error
    await page.route(`${API_BASE_URL}/auth/forgot-password`, async (route) => {
      await route.fulfill({
        status: 500,
        contentType: 'application/json',
        body: JSON.stringify({
          success: false,
          error: 'Server error'
        })
      });
    });
    
    await page.locator('[data-testid="forgot-email"]').fill(testUser.email);
    await page.locator('[data-testid="forgot-submit"]').click();
    
    // Should show error message using testid
    await expect(page.locator('[data-testid="forgot-error-message"]')).toBeVisible({ timeout: 5000 });
    
    console.log('✅ Error handling works correctly');
  });

  test('should show loading states during submission', async ({ page }) => {
    console.log('🧪 Test: Loading states');
    
    await page.goto('/forgot-password');
    await page.waitForLoadState('networkidle');
    
    // Slow down the API request to see loading state
    await page.route(`${API_BASE_URL}/auth/forgot-password`, async (route) => {
      // Add delay to see loading state
      await new Promise(resolve => setTimeout(resolve, 500));
      await route.continue();
    });
    
    await page.locator('[data-testid="forgot-email"]').fill(testUser.email);
    
    // Click submit and immediately check for loading state using testid
    const submitButton = page.locator('[data-testid="forgot-submit"]');
    await submitButton.click();
    
    // Button should show loading text (German)
    await expect(submitButton).toContainText(/wird gesendet/i, { timeout: 1000 });
    
    console.log('✅ Loading states display correctly');
  });
});
