import { QueryClient } from '@tanstack/react-query';

export class NavigateError extends Error {
  constructor(public to: string) {
    super(`Navigation to ${to}`);
    this.name = 'NavigateError';
  }
}

export function authLoader(queryClient: QueryClient) {
  return async () => {
    try {
      // Check if we have a user in local storage for development with mock auth
      const localUser = localStorage.getItem('@webapp:user');

      if (!localUser) {
        // Redirect to login if not authenticated
        throw new NavigateError('/login');
      }

      try {
        const user = JSON.parse(localUser);
        return { user };
      } catch (parseError) {
        // Handle JSON parse error by redirecting to login
        throw new NavigateError('/login');
      }
    } catch (error) {
      if (error instanceof NavigateError) {
        throw error;
      }

      // Handle other errors by redirecting to login
      throw new NavigateError('/login');
    }
  };
}
