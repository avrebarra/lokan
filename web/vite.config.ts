import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    // forward API calls to the mock server during `runtask web dev`
    proxy: {
      '/api': {
        target: `http://localhost:${process.env.MOCK_API_PORT ?? 8787}`,
        changeOrigin: true,
      },
    },
  },
})
