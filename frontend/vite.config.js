import { defineConfig } from 'vite';

export default defineConfig({
  server: {
    proxy: {
      // Go backend; override with API_URL if it's running somewhere else.
      '/api': process.env.API_URL ?? 'http://localhost:8080',
    },
  },
});
