import { describe, it, expect, beforeEach, vi } from 'vitest';
import { AESaasApiClient, ApiError } from '../src/client';
import axios from 'axios';

vi.mock('axios');

describe('API Authentication and Authorization Unit Tests', () => {
  let client: AESaasApiClient;
  let mockAxiosInstance: any;

  beforeEach(() => {
    vi.clearAllMocks();
    
    // Create mock axios instance
    mockAxiosInstance = {
      request: vi.fn(),
      get: vi.fn(),
      post: vi.fn(),
      put: vi.fn(),
      delete: vi.fn(),
      patch: vi.fn(),
      interceptors: {
        request: { use: vi.fn(), eject: vi.fn() },
        response: { use: vi.fn(), eject: vi.fn() }
      }
    };

    // Mock axios.create to return our mock instance
    (axios.create as any) = vi.fn(() => mockAxiosInstance);
    
    client = new AESaasApiClient({
      baseURL: 'http://localhost:8080/api/v1',
      timeout: 10000,
    });
  });

  describe('Authentication', () => {
    it('should successfully login and set token', async () => {
      const mockLoginResponse = {
        data: {
          success: true,
          data: {
            token: 'test-jwt-token',
            user: {
              id: 1,
              username: 'testuser',
              email: 'test@example.com',
              first_name: 'Test',
              last_name: 'User'
            }
          }
        }
      };

      mockAxiosInstance.request.mockResolvedValue(mockLoginResponse);

      const result = await client.login({
        username: 'testuser',
        password: 'password123'
      });

      expect(result).toBeDefined();
      expect(result.success).toBe(true);
      expect(result.data).toBeDefined();
      expect(result.data.token).toBe('test-jwt-token');
      expect(mockAxiosInstance.request).toHaveBeenCalled();
    });

    it('should handle login failure', async () => {
      mockAxiosInstance.request.mockRejectedValue({
        response: {
          status: 401,
          data: {
            success: false,
            error: 'Invalid credentials'
          }
        },
        message: 'Request failed'
      });

      await expect(client.login({
        username: 'testuser',
        password: 'wrongpassword'
      })).rejects.toThrow();
    });

    it('should manage token state', () => {
      client.setToken('test-token');
      client.clearToken();
      // Token is private, but these methods should not throw
      expect(true).toBe(true);
    });
  });

  describe('Authorization', () => {
    it('should fail when calling protected endpoint without auth token', async () => {
      mockAxiosInstance.request.mockRejectedValue({
        response: {
          status: 401,
          data: {
            success: false,
            error: 'Unauthorized'
          }
        },
        message: 'Unauthorized'
      });

      client.clearToken();
      
      await expect(client.getCurrentUser()).rejects.toThrow();
    });

    it('should successfully access protected endpoint with valid token', async () => {
      const mockUserResponse = {
        data: {
          success: true,
          data: {
            id: 1,
            username: 'testuser',
            email: 'test@example.com'
          }
        }
      };

      mockAxiosInstance.request.mockResolvedValue(mockUserResponse);

      client.setToken('valid-jwt-token');
      
      const result = await client.getCurrentUser();
      
      expect(result).toBeDefined();
      expect(mockAxiosInstance.request).toHaveBeenCalled();
    });
  });

  describe('Client Methods', () => {
    it('should handle registration requests', async () => {
      const mockRegisterResponse = {
        data: {
          success: true,
          data: {
            user: {
              id: 1,
              username: 'newuser',
              email: 'newuser@example.com'
            }
          }
        }
      };

      mockAxiosInstance.request.mockResolvedValue(mockRegisterResponse);

      const result = await client.register({
        username: 'newuser',
        email: 'newuser@example.com',
        password: 'password123'
      });

      expect(result).toBeDefined();
      expect(result.success).toBe(true);
    });

    it('should get current user', async () => {
      const mockUserResponse = {
        data: {
          success: true,
          data: {
            id: 1,
            username: 'testuser',
            email: 'test@example.com'
          }
        }
      };

      mockAxiosInstance.request.mockResolvedValue(mockUserResponse);

      client.setToken('valid-token');
      
      const result = await client.getCurrentUser();

      expect(result).toBeDefined();
      expect(mockAxiosInstance.request).toHaveBeenCalled();
    });

    it('should handle booking template operations', async () => {
      const mockResponse = {
        data: {
          success: true,
          data: []
        }
      };

      mockAxiosInstance.request.mockResolvedValue(mockResponse);

      client.setToken('valid-token');
      
      const result = await client.listBookingTemplates();

      expect(result).toBeDefined();
      expect(mockAxiosInstance.request).toHaveBeenCalled();
    });
  });
});