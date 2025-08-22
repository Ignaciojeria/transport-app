import { globalIgnores } from 'eslint/config'

export default [
  // Ignora todos los archivos: desactiva ESLint globalmente en este paquete
  globalIgnores(['**/*'])
]
