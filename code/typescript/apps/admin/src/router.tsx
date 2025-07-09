import { createRootRoute, createRoute, createRouter, redirect } from '@tanstack/react-router';
import { QueryClient } from '@tanstack/react-query';
import { Login } from '@/pages/login';
import { Dashboard } from '@/pages/dashboard';
import { Users } from '@/pages/users';
import { AdminLayout } from '@/layouts/AdminLayout';

// Create the query client
const queryClient = new QueryClient();

// Define the root route
export const rootRoute = createRootRoute();

// Helper function to check if user is authenticated
const checkAuth = () => {
  // Check for user in localStorage
  const user = localStorage.getItem('@admin:user');

  if (!user) {
    // Redirect to login if no user is found
    throw redirect({
      to: '/login',
      replace: true,
    });
  }

  return { user: JSON.parse(user) };
};

// Authentication-related routes
export const loginRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/login',
  beforeLoad: () => {
    // Check if user is already logged in, redirect to admin if they are
    const user = localStorage.getItem('@admin:user');
    if (user) {
      throw redirect({
        to: '/admin',
        replace: true,
      });
    }
    return {};
  },
  component: Login,
});

// Admin layout route
export const adminLayoutRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/admin',
  // Add authentication check before loading any admin routes
  beforeLoad: () => {
    return checkAuth();
  },
  component: AdminLayout,
});

// Admin routes
export const dashboardRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/',
  component: Dashboard,
});

export const usersRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/users',
  component: Users,
});

// Settings route - we'll add this even if the component doesn't exist yet
export const settingsRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/settings',
  component: () => <div>Settings Page</div>,
});

// Root index route to redirect to login
export const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/',
  beforeLoad: () => {
    // Using the TanStack redirect function instead of direct DOM manipulation
    throw redirect({
      to: '/login',
      replace: true,
    });
  },
  component: () => null,
});

// Create and configure the router
export const routeTree = rootRoute.addChildren([
  loginRoute,
  adminLayoutRoute.addChildren([dashboardRoute, usersRoute, settingsRoute]),
  indexRoute,
]);

export const router = createRouter({
  routeTree,
  defaultPreload: 'intent',
  context: {
    queryClient,
  },
});

// Register the router for type safety
declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router;
  }
}
