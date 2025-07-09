import tailwindcss from '@tailwindcss/vite';
import react from '@vitejs/plugin-react';
import path from 'path';
import { defineConfig } from 'vite';

export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    allowedHosts: ['whitelabel.abodemine.test', 'app.omega.test', 'clients-react-saas-whitelabel.abodemine.local'],
    hmr: {
      // Enable HMR
      overlay: true,
    },
    // Add watch options to improve file watching
    watch: {
      usePolling: true,
      interval: 1000,
    },
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
      '@am/commons': path.resolve(__dirname, '../../packages/commons'),
    },
  },
});
