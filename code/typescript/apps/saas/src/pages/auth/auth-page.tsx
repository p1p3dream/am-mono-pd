// src/pages/auth/auth-page.tsx
import { useEffect, useState } from 'react';
import { useRouter } from '@tanstack/react-router';
import { Card, CardContent } from '@am/commons/components/ui/card';
import { usePageTitle } from '@/hooks/use-page-title';
import { useAuth } from '@/contexts/auth-context';

const LOG_PREFIX = '[AUTH_PAGE_DEBUG]';

export function AuthPage() {
  usePageTitle('Authenticating');
  const router = useRouter();
  const { user, authState, checkAuth } = useAuth();
  const [processingAuth, setProcessingAuth] = useState(false);

  // Effect to process authentication
  useEffect(() => {
    const processAuth = async () => {
      // If already processing or user is authenticated, do nothing
      if (processingAuth || (user && authState === 'authenticated')) {
        return;
      }

      console.log(`${LOG_PREFIX} Processing authentication`);
      setProcessingAuth(true);

      try {
        // Try to authenticate
        const success = await checkAuth();
        console.log(`${LOG_PREFIX} Authentication result: ${success ? 'success' : 'failure'}`);

        // If successfully authenticated, redirect to dashboard after a small delay
        if (success) {
          setTimeout(() => {
            router.navigate({ to: '/', replace: true });
          }, 1000);
        }
      } catch (error) {
        console.error(`${LOG_PREFIX} Error during authentication:`, error);
      }
    };

    processAuth();
  }, [user, authState, checkAuth, router, processingAuth]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-linear-to-br from-primary/5 via-secondary/5 to-primary/5">
      <Card className="w-[400px]">
        <CardContent className="pt-6 flex flex-col items-center">
          {/* <Logo size={60} /> */}
          <h2 className="mt-4 text-xl font-semibold text-center">Authenticating</h2>
          <p className="mt-2 text-sm text-center text-muted-foreground">
            Please wait while we set up your session
          </p>

          <div className="mt-6 w-full flex flex-col items-center">
            {authState === 'loading' ? (
              <div className="h-8 w-8 animate-spin rounded-full border-2 border-primary border-t-transparent"></div>
            ) : authState === 'authenticated' ? (
              <div className="flex flex-col items-center">
                <svg
                  className="h-8 w-8 text-green-500"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M5 13l4 4L19 7"
                  />
                </svg>
                <p className="mt-2 text-sm text-green-600">Authentication completed!</p>
              </div>
            ) : authState === 'error' || authState === 'unauthenticated' ? (
              <div className="flex flex-col items-center">
                <svg
                  className="h-8 w-8 text-red-500"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M6 18L18 6M6 6l12 12"
                  />
                </svg>
                <p className="mt-2 text-sm text-red-600">Authentication failed</p>
              </div>
            ) : null}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
