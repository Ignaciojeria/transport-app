import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  // Configuración para reducir logs en desarrollo
  logLevel: 'warn', // Solo muestra warnings y errores
  clearScreen: false, // No limpia la pantalla al recargar
  // Configuración para reducir logs de consola
  define: {
    __DEV__: JSON.stringify(process.env.NODE_ENV === 'development'),
  },
  // Configuración para el servidor de desarrollo
  server: {
    // Reducir logs del servidor
    hmr: {
      overlay: false, // Desactivar overlay de errores HMR
    },
  },
  // Configuración para build
  build: {
    // Reducir logs durante el build
    reportCompressedSize: false,
    chunkSizeWarningLimit: 1000,
  },
})
