// Archivo de prueba para verificar que GunJS funciona correctamente
// Este archivo puede ser eliminado después de las pruebas

import { 
  setRouteStarted, 
  setDeliveryStatus, 
  setDeliveryEvidence,
  gun,
  driverData
} from './driver-gun-state'

// Función de prueba para verificar funcionalidad básica
export function testGunBasicFunctionality() {
  console.log('🧪 Probando GunJS básico...')
  
  // Probar escritura y lectura básica
  const testKey = 'test-key'
  driverData.get(testKey).put('test-value')
  
  // Verificar que se puede leer
  driverData.get(testKey).once((data) => {
    console.log('✅ Lectura básica funciona:', data)
  })
  
  // Probar funciones del driver
  console.log('🧪 Probando funciones del driver...')
  
  const testRouteId = 'test-route-123'
  
  // Probar start route
  setRouteStarted(testRouteId, true)
  console.log('✅ setRouteStarted ejecutado')
  
  // Probar delivery status
  setDeliveryStatus(testRouteId, 0, 0, 0, 'delivered')
  console.log('✅ setDeliveryStatus ejecutado')
  
  // Probar evidence
  setDeliveryEvidence(testRouteId, 0, 0, 0, {
    recipientName: 'Juan Pérez',
    recipientRut: '12345678-9',
    photoDataUrl: 'data:image/jpeg;base64,test',
    takenAt: Date.now()
  })
  console.log('✅ setDeliveryEvidence ejecutado')
  
  console.log('🎉 Todas las pruebas básicas completadas!')
}

// Función para escuchar cambios en tiempo real
export function testGunReactivity() {
  console.log('🧪 Probando reactividad de GunJS...')
  
  const testKey = 'reactive-test'
  let changeCount = 0
  
  // Escuchar cambios
  driverData.get(testKey).on((data) => {
    changeCount++
    console.log(`🔄 Cambio #${changeCount} detectado:`, data)
  })
  
  // Hacer algunos cambios
  setTimeout(() => {
    driverData.get(testKey).put('valor-1')
  }, 100)
  
  setTimeout(() => {
    driverData.get(testKey).put('valor-2')
  }, 200)
  
  setTimeout(() => {
    driverData.get(testKey).put('valor-3')
  }, 300)
  
  setTimeout(() => {
    console.log('🎉 Prueba de reactividad completada!')
  }, 500)
}

// Función para verificar persistencia
export function testGunPersistence() {
  console.log('🧪 Probando persistencia de GunJS...')
  
  const persistenceKey = 'persistence-test'
  const testValue = `test-${Date.now()}`
  
  // Escribir valor
  driverData.get(persistenceKey).put(testValue)
  
  // Verificar inmediatamente
  driverData.get(persistenceKey).once((data) => {
    if (data === testValue) {
      console.log('✅ Persistencia inmediata funciona:', data)
    } else {
      console.log('❌ Error en persistencia inmediata:', data)
    }
  })
  
  console.log('🎉 Prueba de persistencia completada!')
  console.log('💡 Recarga la página para verificar persistencia a largo plazo')
}

// Ejecutar todas las pruebas
export function runAllTests() {
  console.log('🚀 Iniciando pruebas completas de GunJS...')
  
  testGunBasicFunctionality()
  setTimeout(() => testGunReactivity(), 1000)
  setTimeout(() => testGunPersistence(), 2000)
  
  console.log('📝 Revisa la consola para ver los resultados')
}
