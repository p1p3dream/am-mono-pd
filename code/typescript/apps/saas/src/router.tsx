import { createRootRoute, createRoute, createRouter } from '@tanstack/react-router';
import { QueryClient } from '@tanstack/react-query';
import { HomePage } from './pages/home';
import { DashboardPage } from './pages/dashboard';
import { DebugPage } from './pages/debug';
import { MapPage } from './pages/property/map';
import { PropertyDetailPage } from './pages/property/detail';
import { PropertyCMAPage } from './pages/property/cma';
import { PropertyPhotosPage } from './pages/property/photos';
import { UnauthorizedPage } from './pages/unauthorized';
import { LoginPage } from './pages/auth/login';
import { ProtectedLayout } from './components/layout/protected-layout';
import { PropertyLayoutWithContext } from './components/layout/property-layout-with-context';
import { PortfolioLayoutWithContext } from './components/layout/portfolio-layout';
import { PortfolioDetail } from './components/portfolio/portfolio-detail';

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

// Debug route
export const debugRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/debug',
  component: DebugPage,
});

// Protected layout route
export const protectedRoute = createRoute({
  getParentRoute: () => rootRoute,
  id: 'protected',
  component: ProtectedLayout,
});

// Portfolio layout route
export const portfolioLayoutRoute = createRoute({
  getParentRoute: () => protectedRoute,
  id: 'portfolio-layout',
  component: PortfolioLayoutWithContext,
});

// Property layout route
export const propertyLayoutRoute = createRoute({
  getParentRoute: () => protectedRoute,
  id: 'property',
  component: PropertyLayoutWithContext,
});

// Portfolio routes
export const homeRoute = createRoute({
  getParentRoute: () => portfolioLayoutRoute,
  path: '/',
  component: HomePage,
});

export const dashboardRoute = createRoute({
  getParentRoute: () => portfolioLayoutRoute,
  path: '/dashboard',
  component: DashboardPage,
});

export const portfolioDetailRoute = createRoute({
  getParentRoute: () => portfolioLayoutRoute,
  path: '/portfolio/$portfolioId',
  component: PortfolioDetail,
});

// Property routes
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
  debugRoute,
  protectedRoute.addChildren([
    portfolioLayoutRoute.addChildren([homeRoute, dashboardRoute, portfolioDetailRoute]),
    propertyLayoutRoute.addChildren([mapRoute, detailRoute, cmaRoute, photosRoute]),
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
