import { deliveriesData } from './gun'
import type { RouteStart } from '../domain/route-start'
import { useState, useEffect } from 'react'

// Clave para almacenar el estado de inicio de ruta
export const routeStartKey = (routeId: string) => `route_start:${routeId}`

// Clave para almacenar la patente del vehículo
export const vehiclePlateKey = (routeId: string) => `vehicle_plate:${routeId}`

// Clave para almacenar información del conductor
export const driverInfoKey = (routeId: string) => `driver_info:${routeId}`

// Clave para almacenar información del carrier
export const carrierInfoKey = (routeId: string) => `carrier_info:${routeId}`

// Clave para almacenar el estado de inicio de ruta (para el nuevo sistema)
export const routeStartedKey = (routeId: string) => `routeStarted:${routeId}`

// Clave para almacenar la licencia del vehículo (para el nuevo sistema)
export const routeLicenseKey = (routeId: string) => `routeLicense:${routeId}`

// Función para establecer el inicio de ruta
export const setRouteStart = async (routeId: string, routeStart: RouteStart): Promise<void> => {
  try {
    // console.log('🔍 Iniciando setRouteStart con:', { routeId, routeStart }) // Comentado para reducir logs
    
    const key = routeStartKey(routeId)
    // console.log('🔑 Clave generada:', key) // Comentado para reducir logs
    
    // Verificar si hay BigInt en los datos
    // console.log('🔍 Verificando campos de routeStart:', { // Comentado para reducir logs
    //   routeId: routeStart.route?.id,
    //   routeIdType: typeof routeStart.route?.id,
    //   documentID: routeStart.route?.documentID,
    //   documentIDType: typeof routeStart.route?.documentID,
    //   referenceID: routeStart.route?.referenceID,
    //   referenceIDType: typeof routeStart.route?.referenceID
    // })
    
    // Crear una copia segura para JSON
    const safeRouteStart = {
      ...routeStart,
      route: {
        ...routeStart.route,
        // Convertir BigInt a string si existe
        id: routeStart.route?.id ? String(routeStart.route.id) : undefined
      }
    }
    
    const data = {
      ...safeRouteStart,
      timestamp: Date.now(),
      deviceId: `device_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
    }
    
    // console.log('🚀 Guardando inicio de ruta:', data) // Comentado para reducir logs
    // console.log('📝 JSON stringificado:', JSON.stringify(data)) // Comentado para reducir logs
    
    // Verificar que deliveriesData esté disponible
    if (!deliveriesData) {
      throw new Error('deliveriesData no está disponible')
    }
    
    // console.log('💾 Guardando en GunJS...') // Comentado para reducir logs
    await deliveriesData.get(key).put(JSON.stringify(data))
    // console.log('✅ Guardado principal exitoso') // Comentado para reducir logs
    
    // También guardar por separado para acceso rápido
    // console.log('💾 Guardando datos separados...') // Comentado para reducir logs
    await deliveriesData.get(vehiclePlateKey(routeId)).put(routeStart.vehicle.plate)
    await deliveriesData.get(driverInfoKey(routeId)).put(JSON.stringify(routeStart.driver))
    await deliveriesData.get(carrierInfoKey(routeId)).put(JSON.stringify(routeStart.carrier))
    // console.log('✅ Todos los datos guardados exitosamente') // Comentado para reducir logs
    
  } catch (error) {
    console.error('❌ Error guardando inicio de ruta:', error)
    console.error('🔍 Detalles del error:', {
      message: error instanceof Error ? error.message : 'Error desconocido',
      stack: error instanceof Error ? error.stack : 'No disponible',
      routeId,
      routeStart
    })
    throw error
  }
}

// Función para obtener el estado de inicio de ruta
export const getRouteStart = async (routeId: string): Promise<RouteStart | null> => {
  try {
    const key = routeStartKey(routeId)
    const data = await deliveriesData.get(key).once()
    
    if (data && typeof data === 'string') {
      const parsed = JSON.parse(data)
      // Remover campos internos antes de retornar
      const { timestamp, deviceId, ...routeStart } = parsed
      return routeStart
    }
    
    return null
  } catch (error) {
    console.error('Error obteniendo inicio de ruta:', error)
    return null
  }
}

// Función para verificar si una ruta está iniciada
export const isRouteStarted = async (routeId: string): Promise<boolean> => {
  try {
    const routeStart = await getRouteStart(routeId)
    return routeStart !== null
  } catch (error) {
    console.error('Error verificando estado de ruta:', error)
    return false
  }
}

// Función para obtener solo la patente del vehículo
export const getVehiclePlate = async (routeId: string): Promise<string | null> => {
  try {
    const key = vehiclePlateKey(routeId)
    const data = await deliveriesData.get(key).once()
    return typeof data === 'string' ? data : null
  } catch (error) {
    console.error('Error obteniendo patente:', error)
    return null
  }
}

// Función para obtener solo la información del conductor
export const getDriverInfo = async (routeId: string): Promise<RouteStart['driver'] | null> => {
  try {
    const key = driverInfoKey(routeId)
    const data = await deliveriesData.get(key).once()
    
    if (data && typeof data === 'string') {
      return JSON.parse(data)
    }
    
    return null
  } catch (error) {
    console.error('Error obteniendo información del conductor:', error)
    return null
  }
}

// Función para obtener solo la información del carrier
export const getCarrierInfo = async (routeId: string): Promise<RouteStart['carrier'] | null> => {
  try {
    const key = carrierInfoKey(routeId)
    const data = await deliveriesData.get(key).once()
    
    if (data && typeof data === 'string') {
      return JSON.parse(data)
    }
    
    return null
  } catch (error) {
    console.error('Error obteniendo información del carrier:', error)
    return null
  }
}

// Función para limpiar el estado de inicio de ruta (útil para testing o reset)
export const clearRouteStart = async (routeId: string): Promise<void> => {
  try {
    const keys = [
      routeStartKey(routeId),
      vehiclePlateKey(routeId),
      driverInfoKey(routeId),
      carrierInfoKey(routeId)
    ]
    
    for (const key of keys) {
      await deliveriesData.get(key).put(null)
    }
    
    // console.log('🧹 Estado de inicio de ruta limpiado para:', routeId) // Comentado para reducir logs
  } catch (error) {
    console.error('Error limpiando estado de inicio de ruta:', error)
    throw error
  }
}

// ===== FUNCIONES PARA MANEJO DE ESTADO DE RUTAS =====

// Tipos para el estado de rutas
interface RouteStarted {
  status: 'true' | 'false'
  timestamp: number
  deviceId: string
  action: 'route_started' | 'route_stopped'
  routeId: string
}

interface RouteLicense {
  license: string
  timestamp: number
  deviceId: string
  action: 'license_set'
  routeId: string
}

// Helper para generar un ID único del dispositivo
function getRouteDeviceId(): string {
  let deviceId = localStorage.getItem('gun-device-id')
  if (!deviceId) {
    deviceId = `device-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
    localStorage.setItem('gun-device-id', deviceId)
  }
  return deviceId
}

// Funciones para manejar el estado de rutas
export function setRouteStarted(routeId: string, started: boolean): void {
  const key = routeStartedKey(routeId)
  const timestamp = Date.now()
  const deviceId = getRouteDeviceId()
  
  const value: RouteStarted = {
    status: started ? 'true' : 'false',
    timestamp,
    deviceId,
    action: started ? 'route_started' : 'route_stopped',
    routeId
  }
  
  deliveriesData.get(key).put(value)
  // Mantener versión simple para compatibilidad
  deliveriesData.get(`${key}_simple`).put(started ? 'true' : 'false')
}

export function setRouteLicense(routeId: string, license: string): void {
  const key = routeLicenseKey(routeId)
  const timestamp = Date.now()
  const deviceId = getRouteDeviceId()
  
  const value: RouteLicense = {
    license,
    timestamp,
    deviceId,
    action: 'license_set',
    routeId
  }
  
  deliveriesData.get(key).put(value)
}

export function getRouteLicenseFromState(
  state: Record<string, any>,
  routeId: string
): string | undefined {
  const key = routeLicenseKey(routeId)
  const licenseData = state[key]
  return licenseData?.license ?? undefined
}

// Hook específico para monitorear sincronización de ruta iniciada
export function useRouteStartedSync(routeId: string) {
  const [syncData, setSyncData] = useState<{
    isStarted: boolean
    lastAction: string
    deviceId: string
    timestamp: number
    syncedDevices: string[]
  } | null>(null)

  useEffect(() => {
    const key = routeStartedKey(routeId)
    
    const unsubscribe = deliveriesData.get(key).on((data) => {
      if (data && typeof data === 'object') {
        setSyncData({
          isStarted: data.status === 'true',
          lastAction: data.action || 'unknown',
          deviceId: data.deviceId || 'unknown',
          timestamp: data.timestamp || 0,
          syncedDevices: [data.deviceId || 'unknown']
        })
      }
    })

    return () => {
      if (unsubscribe && typeof unsubscribe.off === 'function') {
        unsubscribe.off()
      }
    }
  }, [routeId])

  return syncData
}

// Función para obtener información de sincronización de todas las rutas
export function getAllRoutesSyncInfo(): Promise<Record<string, RouteStarted>> {
  return new Promise((resolve) => {
    const routes: Record<string, RouteStarted> = {}
    
    deliveriesData.map().once((data, key) => {
      if (key && key.includes('routeStarted:') && !key.includes('_simple')) {
        const routeId = key.replace('routeStarted:', '')
        if (data && typeof data === 'object') {
          routes[routeId] = data as RouteStarted
        }
      }
    })
    
    setTimeout(() => resolve(routes), 500)
  })
}

// Exportar tipos
export type { RouteStarted, RouteLicense }


