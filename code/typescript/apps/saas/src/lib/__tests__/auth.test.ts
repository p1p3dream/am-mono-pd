import { jest, describe, it, expect, beforeEach } from '@jest/globals';
import { authLoader, NavigateError } from '../auth';

describe('authLoader', () => {
  const mockQueryClient = {
    getQueryData: jest.fn(),
    setQueryData: jest.fn(),
  };

  beforeEach(() => {
    // Clear localStorage before each test
    localStorage.clear();
    // Clear all mocks
    jest.clearAllMocks();
  });

  it('should return user data when user is in localStorage', async () => {
    const mockUser = { id: 1, name: 'Test User' };
    localStorage.setItem('@webapp:user', JSON.stringify(mockUser));

    const loader = authLoader(mockQueryClient as any);
    const result = await loader();

    expect(result).toEqual({ user: mockUser });
  });

  it('should throw NavigateError to /login when no user is in localStorage', async () => {
    const loader = authLoader(mockQueryClient as any);

    await expect(loader()).rejects.toThrow(NavigateError);
    await expect(loader()).rejects.toMatchObject({
      to: '/login',
    });
  });

  it('should handle invalid JSON in localStorage', async () => {
    localStorage.setItem('@webapp:user', 'invalid-json');

    const loader = authLoader(mockQueryClient as any);

    await expect(loader()).rejects.toThrow(NavigateError);
    await expect(loader()).rejects.toMatchObject({
      to: '/login',
    });
  });

  it('should handle empty user object in localStorage', async () => {
    localStorage.setItem('@webapp:user', '{}');

    const loader = authLoader(mockQueryClient as any);
    const result = await loader();

    expect(result).toEqual({ user: {} });
  });
});
