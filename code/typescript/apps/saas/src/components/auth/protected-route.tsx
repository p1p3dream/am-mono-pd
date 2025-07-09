// src/components/auth/protected-route.tsx
import { useEffect } from 'react';
import { useRouter, useRouterState } from '@tanstack/react-router';
import { useAuth } from '@/contexts/auth-context';
import { UnauthorizedPage } from '@/pages/unauthorized';

const LOG_PREFIX = '[PROTECTED_ROUTE_DEBUG]';

type ProtectedRouteProps = {
  children: React.ReactNode;
};

export function ProtectedRoute({ children }: ProtectedRouteProps) {
  const { user, authState, checkAuth } = useAuth();
  const router = useRouter();
  const { location } = useRouterState();

  // Determine if this is the hello page
  const isHelloPage = location.pathname === '/hello';

  useEffect(() => {
    // For /hello, always check auth via cookie method
    if (isHelloPage && authState !== 'authenticated') {
      console.log(`${LOG_PREFIX} Checking cookie-based auth for /hello page`);
      checkAuth();
    }
  }, [isHelloPage, authState, checkAuth]);

  // If still loading, show spinner
  if (authState === 'loading') {
    return (
      <div className="flex h-screen w-screen items-center justify-center">
        <div className="h-16 w-16 animate-spin rounded-full border-4 border-primary border-t-transparent"></div>
      </div>
    );
  }

  // For /hello path, handle differently - show unauthorized page if not authenticated
  if (isHelloPage) {
    console.log(`${LOG_PREFIX} Handling /hello page protection`);
    return user ? children : <UnauthorizedPage />;
  }

  // For all other paths, including "/", redirect to login if not authenticated
  if (!user) {
    console.log(`${LOG_PREFIX} User not authenticated, redirecting to login`);
    // Use setTimeout to avoid React state update during render
    setTimeout(() => {
      router.navigate({ to: '/login', replace: true });
    }, 0);
    return null;
  }

  // User is authenticated, render protected content
  return <>{children}</>;
}
