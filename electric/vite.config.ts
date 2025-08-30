import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  // Configuración del servidor de desarrollo
  server: {
    port: 5173, // Puerto explícito
    host: true, // Escuchar en todas las interfaces
    open: false, // No abrir automáticamente el navegador
    hmr: {
      overlay: false, // Desactivar overlay de errores HMR
    },
  },
  // Configuración para build
  build: {
    reportCompressedSize: false,
    chunkSizeWarningLimit: 1000,
  },
})
