import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

export default defineConfig({
  plugins: [react()],
  server: {
    // The browser only talks to Vite (port 5173). Vite forwards every request starting with
    // /v1 to the Go API, so the browser never makes a cross-origin (CORS) request.
    proxy: {
      '/v1': {
        target: process.env.API_URL || 'http://localhost:4000',
        changeOrigin: true,
      },
    },
  },
})
