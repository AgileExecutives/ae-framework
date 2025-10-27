import { describe, it, expect, beforeAll } from 'vitest';
import { AESaasApiClient } from '../src/client';

describe('API Health Check', () => {
  let client: AESaasApiClient;

  beforeAll(() => {
    client = new AESaasApiClient({
      baseURL: 'http://localhost:8080/api/v1',
      timeout: 5000,
    });
  });

  it('should successfully call the health endpoint', async () => {
    const health = await client.health();
    
    expect(health).toBeDefined();
    expect(health.status).toBe('healthy');
  });

  it('should return health status with expected structure', async () => {
    const health = await client.health();
    
    expect(health).toHaveProperty('status');
    expect(typeof health.status).toBe('string');
    
    // The health response should contain these fields based on the server implementation
    expect(health).toHaveProperty('timestamp');
    expect(health).toHaveProperty('version');
    expect(health).toHaveProperty('database');
  });
});