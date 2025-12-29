'use client'

import { useState } from 'react'
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Play, ExternalLink, Smartphone } from "lucide-react"
import { useLanguage } from "@/lib/useLanguage"

export function DemoEmbed() {
  const [showDemo, setShowDemo] = useState(false)
  const { t, isLoading, language } = useLanguage()

  if (isLoading) {
    return (
      <section className="py-20 bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
            <p className="text-gray-600">Loading...</p>
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
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 items-center">
          {/* Left side - Information */}
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
                      <h4 className="font-semibold text-gray-900">{t.demo.selectDish.title}</h4>
                      <p className="text-gray-600 text-sm">{t.demo.selectDish.description}</p>
                    </div>
                  </div>
                  <div className="flex items-start space-x-3">
                    <div className="w-2 h-2 bg-green-600 rounded-full mt-2"></div>
                    <div>
                      <h4 className="font-semibold text-gray-900">{t.demo.chooseSide.title}</h4>
                      <p className="text-gray-600 text-sm">{t.demo.chooseSide.description}</p>
                    </div>
                  </div>
                  <div className="flex items-start space-x-3">
                    <div className="w-2 h-2 bg-purple-600 rounded-full mt-2"></div>
                    <div>
                      <h4 className="font-semibold text-gray-900">{t.demo.enterName.title}</h4>
                      <p className="text-gray-600 text-sm">{t.demo.enterName.description}</p>
                    </div>
                  </div>
                  <div className="flex items-start space-x-3">
                    <div className="w-2 h-2 bg-orange-600 rounded-full mt-2"></div>
                    <div>
                      <h4 className="font-semibold text-gray-900">{t.demo.selectTime.title}</h4>
                      <p className="text-gray-600 text-sm">{t.demo.selectTime.description}</p>
                    </div>
                  </div>
                  <div className="flex items-start space-x-3">
                    <div className="w-2 h-2 bg-green-500 rounded-full mt-2"></div>
                    <div>
                      <h4 className="font-semibold text-gray-900">{t.demo.sendWhatsApp.title}</h4>
                      <p className="text-gray-600 text-sm">{t.demo.sendWhatsApp.description}</p>
                    </div>
                  </div>
                </div>
                
                <div className="pt-4 border-t">
                  <Button 
                    onClick={() => setShowDemo(!showDemo)}
                    className="w-full bg-blue-600 hover:bg-blue-700"
                  >
                    <Play className="h-4 w-4 mr-2" />
                    {showDemo ? t.demo.hideDemo : t.demo.viewDemo}
                  </Button>
                  
                  <Button 
                    variant="outline" 
                    className="w-full mt-2"
                    onClick={() => window.open('https://cadorago.web.app/', '_blank')}
                  >
                    <ExternalLink className="h-4 w-4 mr-2" />
                    {t.demo.openNewTab}
                  </Button>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Right side - Embedded demo */}
          <div className="relative">
            {showDemo ? (
              <div className="relative">
                <div className="bg-white rounded-lg shadow-lg overflow-hidden iframe-container">
                  <iframe 
                    src="https://cadorago.web.app/" 
                    width="100%" 
                    height="600"
                    className="border-0"
                    title="Demo MiCartaPro"
                    loading="lazy"
                    allow="camera; microphone; geolocation; autoplay; clipboard-write"
                    sandbox="allow-same-origin allow-scripts allow-forms allow-popups allow-popups-to-escape-sandbox allow-presentation allow-camera allow-microphone"
                  />
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
                  • {t.demo.preview.features.digitalMenu}<br/>
                  • {t.demo.preview.features.shoppingCart}<br/>
                  • {t.demo.preview.features.whatsapp}
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </section>
  )
}
