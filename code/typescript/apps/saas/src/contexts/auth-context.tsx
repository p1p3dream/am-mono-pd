// src/contexts/auth-context.tsx
import { config } from '@/config/env';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { createContext, ReactNode, useCallback, useContext, useEffect, useState } from 'react';

const LOG_PREFIX = '[AUTH_DEBUG]';

type User = {
  id: string;
  email: string;
  name: string;
};

type AuthState = 'loading' | 'authenticated' | 'unauthenticated' | 'error';

type AuthContextType = {
  user: User | null;
  authState: AuthState;
  checkAuth: () => Promise<boolean>;
  signIn: (email: string, password: string) => Promise<boolean>;
  signOut: () => void;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

type AuthProviderProps = {
  children: ReactNode;
};

// Query keys for TanStack Query
const authKeys = {
  me: ['auth', 'me'],
  login: () => ['auth', 'login'],
  loginWithCredentials: (email: string, password: string) => ['auth', 'login', email, password],
  logout: ['auth', 'logout'],
};

export function AuthProvider({ children }: AuthProviderProps) {
  const [user, setUser] = useState<User | null>(null);
  const [authState, setAuthState] = useState<AuthState>('loading');
  const [authChecked, setAuthChecked] = useState(false);
  const queryClient = useQueryClient();

  // Check authentication status using TanStack Query
  const checkAuth = useCallback(async (): Promise<boolean> => {
    if (authChecked && authState !== 'loading') {
      console.log(`${LOG_PREFIX} Authentication already checked, state: ${authState}`);
      return authState === 'authenticated';
    }

    console.log(`${LOG_PREFIX} Checking authentication`);
    setAuthState('loading');

    try {
      // For development, use mock
      if (config.useMockAuth) {
        console.log(`${LOG_PREFIX} Using mock authentication for development`);
        await new Promise((resolve) => setTimeout(resolve, 500));

        const mockUser = {
          id: 'dev-1',
          email: 'dev-user@example.com',
          name: 'Dev User',
        };

        setUser(mockUser);
        setAuthState('authenticated');
        setAuthChecked(true);
        return true;
      }

      // Use TanStack Query to fetch auth status
      const result = await queryClient.fetchQuery({
        queryKey: authKeys.me,
        queryFn: async () => {
          const response = await fetch(`${config.apiUrl}/auth/load`, {
            method: 'POST',
            credentials: 'include', // Important to send cookies
          });

          if (!response.ok) {
            throw new Error(`Authentication failed: ${response.status}`);
          }

          return response.json();
        },
      });

      console.log(`${LOG_PREFIX} User authenticated:`, result);
      setUser(result);
      setAuthState('authenticated');
      setAuthChecked(true);
      return true;
    } catch (error) {
      console.error(`${LOG_PREFIX} Error checking authentication:`, error);
      setUser(null);
      setAuthState('unauthenticated');
      setAuthChecked(true);
      return false;
    }
  }, [authChecked, authState, queryClient]);

  // Login implementation using useQuery
  const { refetch: loginRefetch } = useQuery({
    queryKey: authKeys.login(),
    queryFn: async () => {
      // This function will be called with actual credentials during refetch
      return null;
    },
    enabled: false, // Don't run on component mount
  });

  // Traditional sign in with email/password
  const signIn = useCallback(
    async (email: string, password: string): Promise<boolean> => {
      console.log(`${LOG_PREFIX} Signing in with email: ${email}`);
      setAuthState('loading');

      try {
        // Create a custom login query function with the credentials
        const customLoginFn = async () => {
          if (config.useMockAuth) {
            // Simulate network delay
            await new Promise((resolve) => setTimeout(resolve, 800));

            // Mock credentials validation
            if (email === 'client@abodemine.com' && password === '123') {
              return {
                id: 'login-1',
                email: email,
                name: 'client User',
              };
            } else {
              throw new Error('Invalid credentials');
            }
          }

          // Real implementation calling API
          const response = await fetch(`${config.apiUrl}/auth/login`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            credentials: 'include', // This will handle cookies if backend sets them
            body: JSON.stringify({ email, password }),
          });

          if (!response.ok) {
            throw new Error('Login failed');
          }

          return response.json();
        };

        // Use the refetch functionality with our custom query function
        const userData = await queryClient.fetchQuery({
          queryKey: authKeys.loginWithCredentials(email, password),
          queryFn: customLoginFn,
        });

        if (!userData) {
          throw new Error('No user data returned from login');
        }

        if (config.useMockAuth) {
          localStorage.setItem('@webapp:user', JSON.stringify(userData));
        }

        setUser(userData);
        setAuthState('authenticated');
        // Invalidate the auth status query to refresh it
        queryClient.invalidateQueries({ queryKey: authKeys.me });
        return true;
      } catch (error) {
        console.error(`${LOG_PREFIX} Login failed:`, error);
        setAuthState('error');
        return false;
      }
    },
    [queryClient]
  );

  // Logout implementation using useQuery
  const { refetch: logoutRefetch } = useQuery({
    queryKey: authKeys.logout,
    queryFn: async () => null, // Placeholder, actual implementation happens in signOut
    enabled: false, // Don't run on component mount
  });

  // Sign out functionality
  const signOut = useCallback(() => {
    console.log(`${LOG_PREFIX} Signing out`);

    // Custom logout function
    const performLogout = async () => {
      try {
        if (!config.useMockAuth) {
          const response = await fetch(`${config.apiUrl}/auth/logout`, {
            method: 'POST',
            credentials: 'include',
          });

          if (!response.ok) {
            throw new Error('Logout failed');
          }
        }

        // Clear any local storage items
        localStorage.removeItem('@webapp:user');

        // Update state
        setUser(null);
        setAuthState('unauthenticated');

        // Invalidate auth queries
        queryClient.invalidateQueries({ queryKey: authKeys.me });

        // Redirect to login page
        window.location.href = '/login';
      } catch (error) {
        console.error(`${LOG_PREFIX} Error logging out:`, error);
        // Still try to clean up client-side state
        setUser(null);
        setAuthState('unauthenticated');
      }
    };

    // Execute the logout query
    queryClient.fetchQuery({
      queryKey: authKeys.logout,
      queryFn: performLogout,
    });
  }, [queryClient]);

  // Initial effect to check authentication
  useEffect(() => {
    if (!authChecked) {
      checkAuth();
    }
  }, [authChecked, checkAuth]);

  return (
    <AuthContext.Provider
      value={{
        user,
        authState,
        checkAuth,
        signIn,
        signOut,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
