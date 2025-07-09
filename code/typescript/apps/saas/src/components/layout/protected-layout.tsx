import { Outlet, useRouter } from '@tanstack/react-router';
import { useAuth } from '@/contexts/auth-context';

export function ProtectedLayout() {
  const { user, authState } = useAuth();
  const router = useRouter();

  // If loading, show spinner
  if (authState === 'loading') {
    return (
      <div className="flex h-screen w-screen items-center justify-center">
        <div className="h-16 w-16 animate-spin rounded-full border-4 border-primary border-t-transparent"></div>
      </div>
    );
  }

  // If not authenticated, redirect to login
  if (!user) {
    // Use setTimeout to avoid React state update during render
    setTimeout(() => {
      router.navigate({ to: '/login' });
    }, 0);

    return null;
  }

  // If authenticated, render children
  return (
    <div>
      <Outlet />
    </div>
  );
}
