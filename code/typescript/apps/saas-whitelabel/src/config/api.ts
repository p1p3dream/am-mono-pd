import { config } from './env';

export const API_CONFIG = {
  useMockData: config.useMockData || false,
  apiBaseUrl: config.apiUrl || 'http://localhost:3000/api',
} as const;
