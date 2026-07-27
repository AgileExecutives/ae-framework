import { test as base, request as playwrightRequest } from '@playwright/test';

const API_BASE_URL = 'http://localhost:8080/api/v1';

export interface AuthUser {
  username: string;
  email: string;
  password: string;
  firstName: string;
  lastName: string;
  companyName: string;
  token?: string;
}

// Extend the base test with authentication fixtures
export const test = base.extend<{
  authenticatedUser: AuthUser;
}, {
  authUserPool: AuthUser[];
}>({
  // Create a pool of authenticated users (runs once per worker)
  authUserPool: [async ({}, use, workerInfo) => {
    const poolSize = 5; // Create 5 users per worker
    const users: AuthUser[] = [];
    
    console.log(`🔧 Worker ${workerInfo.workerIndex}: Creating user pool...`);
    
    const apiContext = await playwrightRequest.newContext();
    
    for (let i = 0; i < poolSize; i++) {
      const timestamp = Date.now();
      const workerIndex = workerInfo.workerIndex;
      const user: AuthUser = {
        username: `testuser_w${workerIndex}_${i}_${timestamp}`,
        email: `testuser_w${workerIndex}_${i}_${timestamp}@example.com`,
        password: 'TestPass123!',
        firstName: 'Test',
        lastName: `User${i}`,
        companyName: `Test Company W${workerIndex} ${i}`,
      };
      
      try {
        const response = await apiContext.post(`${API_BASE_URL}/auth/register`, {
          data: {
            username: user.username,
            email: user.email,
            password: user.password,
            first_name: user.firstName,
            last_name: user.lastName,
            company_name: user.companyName,
          },
        });
        
        if (response.status() === 201 || response.status() === 200) {
          const data = await response.json();
          user.token = data.data?.token;
          users.push(user);
          console.log(`✅ Created user: ${user.username}`);
        } else if (response.status() === 429) {
          console.warn(`⚠️  Rate limited while creating user pool. Ensure RATE_LIMIT_ENABLED=false is set on server.`);
          // Still push the user without token - tests can handle registration if needed
          users.push(user);
        } else {
          console.warn(`⚠️  Failed to create user ${user.username}: ${response.status()}`);
          users.push(user); // Add anyway for tests that might create it themselves
        }
      } catch (error) {
        console.error(`❌ Error creating user ${user.username}:`, error);
        users.push(user); // Add anyway for fallback
      }
      
      // Small delay between registrations
      if (i < poolSize - 1) {
        await new Promise(resolve => setTimeout(resolve, 100));
      }
    }
    
    await apiContext.dispose();
    console.log(`✅ Worker ${workerInfo.workerIndex}: User pool created (${users.length} users)`);
    
    // Provide the pool to tests
    await use(users);
    
    // Cleanup is handled by test database reset or similar
    console.log(`🧹 Worker ${workerInfo.workerIndex}: User pool cleanup (if needed)`);
  }, { scope: 'worker' }],
  
  // Provide a single authenticated user from the pool (runs per test)
  authenticatedUser: async ({ authUserPool }, use, testInfo) => {
    // Get a user from the pool using test index
    const userIndex = testInfo.testId.charCodeAt(0) % authUserPool.length;
    const user = authUserPool[userIndex];
    
    console.log(`🧪 Test "${testInfo.title}" using user: ${user.username}`);
    
    await use(user);
  },
});

export { expect } from '@playwright/test';
