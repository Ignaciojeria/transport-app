'use client'

import { useState, useEffect } from 'react'
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Play, ExternalLink, Smartphone } from "lucide-react"

// Generar short UUID único para esta sesión de demo
const generateDemoUUID = () => {
  // Generar un short UUID de 8 caracteres alfanuméricos
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
  let result = '';
  for (let i = 0; i < 8; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length));
  }
  return 'DEMO-' + result;
};

export function DemoEmbed() {
  const [showDemo, setShowDemo] = useState(false)
  const [routeId, setRouteId] = useState<string>('')
  
  // Generar UUID solo en el cliente para evitar errores de hidratación
  useEffect(() => {
    setRouteId(generateDemoUUID())
  }, [])

  return (
    <section className="py-20 bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center mb-16">
          <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
            Prueba la App en Vivo
          </h2>
          <p className="text-xl text-gray-600 mb-4">
            Experimenta cómo funciona la app desde la perspectiva del conductor con datos de prueba
          </p>
          {routeId && (
            <div className="bg-blue-50 border border-blue-200 rounded-lg p-3 inline-block">
              <p className="text-sm text-blue-800">
                <span className="font-semibold">ID de Ruta:</span> <span className="font-mono text-blue-900">{routeId}</span>
              </p>
            </div>
          )}
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 items-center">
          {/* Lado izquierdo - Información */}
          <div className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center">
                  <Smartphone className="h-6 w-6 text-blue-600 mr-2" />
                  Demo Interactiva
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="space-y-3">
                  <div className="flex items-start space-x-3">
                    <div className="w-2 h-2 bg-blue-600 rounded-full mt-2"></div>
                    <div>
                      <h4 className="font-semibold text-gray-900">Inicia una ruta</h4>
                      <p className="text-gray-600 text-sm">Simula el inicio de una ruta con patente</p>
                    </div>
                  </div>
                  <div className="flex items-start space-x-3">
                    <div className="w-2 h-2 bg-green-600 rounded-full mt-2"></div>
                    <div>
                      <h4 className="font-semibold text-gray-900">Marca entregas</h4>
                      <p className="text-gray-600 text-sm">Simula el proceso de entrega con evidencia</p>
                    </div>
                  </div>
                  <div className="flex items-start space-x-3">
                    <div className="w-2 h-2 bg-purple-600 rounded-full mt-2"></div>
                    <div>
                      <h4 className="font-semibold text-gray-900">Navega con mapas</h4>
                      <p className="text-gray-600 text-sm">Integración con Google Maps y Waze</p>
                    </div>
                  </div>
                  <div className="flex items-start space-x-3">
                    <div className="w-2 h-2 bg-orange-600 rounded-full mt-2"></div>
                    <div>
                      <h4 className="font-semibold text-gray-900">Genera reportes</h4>
                      <p className="text-gray-600 text-sm">Descarga reportes en CSV y Excel</p>
                    </div>
                  </div>
                </div>
                
                <div className="pt-4 border-t">
                  <Button 
                    onClick={() => setShowDemo(!showDemo)}
                    className="w-full bg-blue-600 hover:bg-blue-700"
                  >
                    <Play className="h-4 w-4 mr-2" />
                    {showDemo ? 'Ocultar Demo' : 'Ver Demo'}
                  </Button>
                  
                  <Button 
                    variant="outline" 
                    className="w-full mt-2"
                    onClick={() => routeId && window.open(`http://localhost:5173/demo?routeId=${routeId}`, '_blank')}
                    disabled={!routeId}
                  >
                    <ExternalLink className="h-4 w-4 mr-2" />
                    Abrir en Nueva Pestaña
                  </Button>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Lado derecho - Demo embebida */}
          <div className="relative">
            {showDemo ? (
              <div className="relative">
                <div className="bg-white rounded-lg shadow-lg overflow-hidden">
                  {routeId && (
                    <iframe 
                      src={`https://einar-404623.web.app/demo?routeId=${routeId}`} 
                      width="100%" 
                      height="600"
                      className="border-0"
                      title="Demo TransportApp"
                      loading="lazy"
                      allow="camera; microphone; geolocation; autoplay; clipboard-write"
                      sandbox="allow-same-origin allow-scripts allow-forms allow-popups allow-popups-to-escape-sandbox allow-presentation allow-camera allow-microphone"
                    />
                  )}
                </div>
              </div>
            ) : (
              <div className="bg-gradient-to-br from-blue-50 to-indigo-100 rounded-lg p-8 text-center">
                <div className="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-4">
                  <Smartphone className="h-8 w-8 text-blue-600" />
                </div>
                <h3 className="text-xl font-semibold text-gray-900 mb-2">
                  Demo Interactiva
                </h3>
                <p className="text-gray-600 mb-4">
                  Haz clic en &quot;Ver Demo&quot; para experimentar la aplicación con datos simulados
                </p>
                <div className="text-sm text-gray-500">
                  • 9 entregas simuladas<br/>
                  • Funcionalidad completa<br/>
                  • Datos realistas de Santiago
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </section>
  )
}
