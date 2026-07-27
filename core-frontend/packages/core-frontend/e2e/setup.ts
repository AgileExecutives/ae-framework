import { request } from '@playwright/test';

const API_BASE_URL = 'http://localhost:8080/api/v1';

export default async function globalSetup() {
  console.log('🔧 Verifying server setup...');
  
  // Ensure RATE_LIMIT_ENABLED is set to false for testing
  if (process.env.RATE_LIMIT_ENABLED !== 'false') {
    console.warn('⚠️  RATE_LIMIT_ENABLED is not set to false. Tests may be rate limited.');
    console.warn('   Set RATE_LIMIT_ENABLED=false in your server environment for testing.');
  }
  
  // Create a request context
  const apiContext = await request.newContext();
  
  // Check if server is running and properly configured
  try {
    // First try the direct health endpoint for configuration details
    const directHealthResponse = await fetch('http://localhost:8080/health');
    if (directHealthResponse.ok) {
      const healthData = await directHealthResponse.json();
      console.log('🔍 Server configuration check:');
      console.log(`   Status: ${healthData.status}`);
      
      if (healthData.settings) {
        console.log(`   Mock Email: ${healthData.settings.mock_email}`);
        console.log(`   Rate Limit: ${healthData.settings.rate_limit_enabled}`);
        
        if (healthData.settings.mock_email !== true) {
          throw new Error('Server is running but MOCK_EMAIL is not set to true');
        }
        if (healthData.settings.rate_limit_enabled !== false) {
          throw new Error('Server is running but RATE_LIMIT_ENABLED is not set to false');
        }
        
        console.log('✅ Server is running and properly configured for E2E testing');
      } else {
        console.log('⚠️  Could not verify server configuration, proceeding anyway...');
      }
    } else {
      // Fallback to API health endpoint
      const healthResponse = await apiContext.get(`${API_BASE_URL}/health`);
      
      if (healthResponse.status() !== 200) {
        throw new Error(`Server health check failed with status: ${healthResponse.status()}`);
      }
      
      const healthData = await healthResponse.json();
      if (healthData.status !== 'healthy') {
        throw new Error(`Server is not healthy: ${JSON.stringify(healthData)}`);
      }
      
      console.log('✅ Server is running and healthy');
    }
  } catch (error) {
    console.error('❌ Server is not running. Please start the server with:');
    console.error('');
    console.error('   cd ../server-api');
    console.error('   RATE_LIMIT_ENABLED=false MOCK_EMAIL=true air');
    console.error('');
    console.error('💡 Or use the helper script:');
    console.error('   ./run-e2e-tests.sh');
    console.error('');
    console.error('📖 See E2E_TEST_GUIDE.md for more details');
    console.error('');
    throw new Error(`Server not available: ${error instanceof Error ? error.message : String(error)}`);
  }
  
  // Verify email service is working (should not fail even with mock)
  try {
    // Try a simple registration to test email service indirectly
    const testResponse = await apiContext.post(`${API_BASE_URL}/auth/register`, {
      data: {
        email: `setup_test_${Date.now()}@example.com`,
        password: 'TestPass123!',
        first_name: 'Setup',
        last_name: 'Test',
        company_name: 'Test Company',
        accept_terms: true,
        newsletter_opt_in: false
      }
    });
    
    // Should either succeed or fail with 409 (user exists), but not fail due to email service
    if ([200, 201, 409].includes(testResponse.status())) {
      console.log('✅ Email service is properly configured (likely mocked)');
    } else {
      // Don't fail setup on registration errors - just log and continue
      console.log(`ℹ️  Setup registration test returned status ${testResponse.status()}`);
    }
    
  } catch (error) {
    // Don't fail setup on email service verification errors
    console.log('ℹ️  Email service verification skipped');
  }
  
  // Dispose of the request context
  await apiContext.dispose();
  
  console.log('🔧 Global setup completed successfully');
}