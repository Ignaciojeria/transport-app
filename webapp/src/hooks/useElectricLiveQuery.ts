import { useState, useEffect, useRef } from 'react'

/**
 * Hook para manejar consultas en tiempo real con Electric SQL usando LiveQuery
 */
export const useElectricLiveQuery = <T>(
  queryFn: () => Promise<T>,
  dependencies: any[] = [],
  options: {
    enabled?: boolean
    refetchInterval?: number
    onSuccess?: (data: T) => void
    onError?: (error: Error) => void
  } = {}
) => {
  const [data, setData] = useState<T | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)
  const intervalRef = useRef<number | null>(null)

  const {
    enabled = true,
    refetchInterval,
    onSuccess,
    onError
  } = options

  const executeQuery = async () => {
    if (!enabled) return

    try {
      setIsLoading(true)
      setError(null)
      
      console.log('🔄 Ejecutando consulta LiveQuery...')
      const result = await queryFn()
      
      setData(result)
      onSuccess?.(result)
      
      console.log('✅ Consulta LiveQuery exitosa:', result)
    } catch (err) {
      const error = err instanceof Error ? err : new Error('Error desconocido')
      setError(error)
      onError?.(error)
      
      console.error('❌ Error en consulta LiveQuery:', error)
    } finally {
      setIsLoading(false)
    }
  }

  // Ejecutar consulta cuando cambien las dependencias
  useEffect(() => {
    executeQuery()
  }, dependencies)

  // Configurar refetch automático si se especifica
  useEffect(() => {
    if (refetchInterval && enabled) {
      intervalRef.current = setInterval(executeQuery, refetchInterval)
      
      return () => {
        if (intervalRef.current) {
          clearInterval(intervalRef.current)
        }
      }
    }
  }, [refetchInterval, enabled])

  // Limpiar interval al desmontar
  useEffect(() => {
    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current)
      }
    }
  }, [])

  return {
    data,
    isLoading,
    error,
    refetch: executeQuery
  }
}

/**
 * Hook específico para consultas de cuentas con Electric SQL
 */
export const useAccountLiveQuery = (token: string, email: string) => {
  return useElectricLiveQuery(
    async () => {
      const { findAccountByEmail } = await import('../services/electricService')
      return findAccountByEmail(token, email)
    },
    [token, email],
    {
      enabled: !!token && !!email,
      refetchInterval: 5000, // Refetch cada 5 segundos
      onSuccess: (data) => {
        console.log('🔄 Datos de cuenta actualizados:', data)
      },
      onError: (error) => {
        console.error('❌ Error en consulta de cuenta:', error)
      }
    }
  )
}

/**
 * Hook específico para consultas de tenants con Electric SQL
 */
export const useTenantsLiveQuery = (token: string, accountId: string) => {
  return useElectricLiveQuery(
    async () => {
      const { findTenantsByAccountId } = await import('../services/electricService')
      return findTenantsByAccountId(token, accountId)
    },
    [token, accountId],
    {
      enabled: !!token && !!accountId,
      refetchInterval: 5000, // Refetch cada 5 segundos
      onSuccess: (data) => {
        console.log('🔄 Datos de tenants actualizados:', data)
      },
      onError: (error) => {
        console.error('❌ Error en consulta de tenants:', error)
      }
    }
  )
}
