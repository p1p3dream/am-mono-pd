import { AuthProvider } from './contexts/auth-context';
import { PortfolioProvider } from './contexts/portfolio-context';
import { RouterProvider } from '@tanstack/react-router';
import { router } from './router';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <PortfolioProvider>
          <RouterProvider router={router} />
          <ReactQueryDevtools initialIsOpen={false} />
        </PortfolioProvider>
      </AuthProvider>
    </QueryClientProvider>
  );
}

export default App;
