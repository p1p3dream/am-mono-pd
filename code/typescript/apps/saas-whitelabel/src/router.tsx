import { createRootRoute, createRoute, createRouter } from '@tanstack/react-router';
import { QueryClient } from '@tanstack/react-query';
import { HomePage } from './pages/selector';
import { MapPage } from './pages/property/map';
import { PropertyDetailPage } from './pages/property/detail';
import { PropertyCMAPage } from './pages/property/cma';
import { PropertyPhotosPage } from './pages/property/photos';
import { UnauthorizedPage } from './pages/unauthorized';
import { LoginPage } from './pages/auth/login';
import { ProtectedLayout } from './components/layout/protected-layout';
import { PropertyLayoutWithContext } from './components/layout/property-layout-with-context';

// Create the query client
const queryClient = new QueryClient();

// Define the root route
export const rootRoute = createRootRoute();

// Authentication-related routes
export const loginRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/login',
  component: LoginPage,
});

export const unauthorizedRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/unauthorized',
  component: UnauthorizedPage,
});

// Protected layout route
export const protectedRoute = createRoute({
  getParentRoute: () => rootRoute,
  id: 'protected',
  component: ProtectedLayout,
});

// Property layout route
export const propertyLayoutRoute = createRoute({
  getParentRoute: () => protectedRoute,
  id: 'property',
  component: PropertyLayoutWithContext,
});

// Property routes
export const homeRoute = createRoute({
  getParentRoute: () => protectedRoute,
  path: '/',
  component: HomePage,
});

export const mapRoute = createRoute({
  getParentRoute: () => propertyLayoutRoute,
  path: '/property/$propertyId/map',
  component: MapPage,
});

export const detailRoute = createRoute({
  getParentRoute: () => propertyLayoutRoute,
  path: '/property/$propertyId/detail',
  component: PropertyDetailPage,
});

export const cmaRoute = createRoute({
  getParentRoute: () => propertyLayoutRoute,
  path: '/property/$propertyId/cma',
  component: PropertyCMAPage,
});

export const photosRoute = createRoute({
  getParentRoute: () => propertyLayoutRoute,
  path: '/property/$propertyId/photos',
  component: PropertyPhotosPage,
});

// Create and configure the router
export const routeTree = rootRoute.addChildren([
  loginRoute,
  unauthorizedRoute,
  protectedRoute.addChildren([
    propertyLayoutRoute.addChildren([homeRoute, mapRoute, detailRoute, cmaRoute, photosRoute]),
  ]),
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
