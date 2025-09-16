'use client'

import { useState, useEffect } from 'react'
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Play, ExternalLink, Smartphone } from "lucide-react"
import { useLanguage } from "@/lib/useLanguage"

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
  const { t, isLoading, language } = useLanguage()
  
  // Generar UUID solo en el cliente para evitar errores de hidratación
  useEffect(() => {
    setRouteId(generateDemoUUID())
  }, [])

  // Forzar re-render del iframe cuando cambie el idioma
  const iframeKey = `${routeId}-${language}`

  // Mostrar loading mientras se carga el idioma
  if (isLoading) {
    return (
      <section className="py-20 bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
            <p className="text-gray-600">Cargando...</p>
          </div>
        </div>
      </section>
    )
  }

  return (
    <section className="py-20 bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center mb-16">
          <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
            {t.demo.title}
          </h2>
          <p className="text-xl text-gray-600 mb-4">
            {t.demo.subtitle}
          </p>
          {routeId && (
            <div className="bg-blue-50 border border-blue-200 rounded-lg p-3 inline-block">
              <p className="text-sm text-blue-800">
                <span className="font-semibold">{t.demo.routeId}:</span> <span className="font-mono text-blue-900">{routeId}</span>
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
                  {t.demo.interactiveDemo}
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="space-y-3">
                  <div className="flex items-start space-x-3">
                    <div className="w-2 h-2 bg-blue-600 rounded-full mt-2"></div>
                    <div>
                      <h4 className="font-semibold text-gray-900">{t.demo.features.startRoute.title}</h4>
                      <p className="text-gray-600 text-sm">{t.demo.features.startRoute.description}</p>
                    </div>
                  </div>
                  <div className="flex items-start space-x-3">
                    <div className="w-2 h-2 bg-green-600 rounded-full mt-2"></div>
                    <div>
                      <h4 className="font-semibold text-gray-900">{t.demo.features.markDeliveries.title}</h4>
                      <p className="text-gray-600 text-sm">{t.demo.features.markDeliveries.description}</p>
                    </div>
                  </div>
                  <div className="flex items-start space-x-3">
                    <div className="w-2 h-2 bg-purple-600 rounded-full mt-2"></div>
                    <div>
                      <h4 className="font-semibold text-gray-900">{t.demo.features.navigateMaps.title}</h4>
                      <p className="text-gray-600 text-sm">{t.demo.features.navigateMaps.description}</p>
                    </div>
                  </div>
                  <div className="flex items-start space-x-3">
                    <div className="w-2 h-2 bg-orange-600 rounded-full mt-2"></div>
                    <div>
                      <h4 className="font-semibold text-gray-900">{t.demo.features.generateReports.title}</h4>
                      <p className="text-gray-600 text-sm">{t.demo.features.generateReports.description}</p>
                    </div>
                  </div>
                </div>
                
                <div className="pt-4 border-t">
                  <Button 
                    onClick={() => setShowDemo(!showDemo)}
                    className="w-full bg-blue-600 hover:bg-blue-700"
                  >
                    <Play className="h-4 w-4 mr-2" />
                    {showDemo ? t.demo.buttons.hideDemo : t.demo.buttons.viewDemo}
                  </Button>
                  
                  <Button 
                    variant="outline" 
                    className="w-full mt-2"
                    onClick={() => routeId && window.open(`https://einar-404623.web.app/demo?routeId=${routeId}&lang=${language}`, '_blank')}
                    disabled={!routeId}
                  >
                    <ExternalLink className="h-4 w-4 mr-2" />
                    {t.demo.buttons.openNewTab}
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
                      key={iframeKey}
                      src={`https://einar-404623.web.app/demo?routeId=${routeId}&lang=${language}`} 
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
                  {t.demo.preview.title}
                </h3>
                <p className="text-gray-600 mb-4">
                  {t.demo.preview.description}
                </p>
                <div className="text-sm text-gray-500">
                  • {t.demo.preview.features.simulatedDeliveries}<br/>
                  • {t.demo.preview.features.fullFunctionality}<br/>
                  • {t.demo.preview.features.realisticData}
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </section>
  )
}
