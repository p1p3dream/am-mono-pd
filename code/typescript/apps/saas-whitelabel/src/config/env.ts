export const config = {
  // Use a default value of false for production, and true can be set in development
  useMockAuth: import.meta.env.VITE_USE_MOCK_AUTH === 'true' || false,

  apiUrl: import.meta.env.VITE_API_URL || '/api',
  customerWebsiteUrl: import.meta.env.VITE_CUSTOMER_WEBSITE_URL || 'https://customer-website.com',
  useMockData: import.meta.env.VITE_USE_MOCK_DATA === 'true',
};
