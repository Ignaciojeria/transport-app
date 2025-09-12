'use client'

import { useState, useEffect } from 'react'
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { 
  Truck, 
  MapPin, 
  Users, 
  BarChart3, 
  CheckCircle, 
  ArrowRight, 
  Smartphone, 
  Navigation,
  Clock,
  DollarSign,
  Zap,
  Shield,
  Star,
  Quote
} from "lucide-react"
import { motion } from "framer-motion"
import { DemoEmbed } from "@/components/DemoEmbed"
import { useLanguage } from "@/lib/useLanguage"
import { LanguageSelector } from "@/components/LanguageSelector"

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

export default function LandingPage() {
  const [routeId, setRouteId] = useState<string>('')
  const { language, changeLanguage, t, isLoading, availableLanguages, languageNames, languageFlags } = useLanguage()
  
  // Generar UUID solo en el cliente para evitar errores de hidratación
  useEffect(() => {
    setRouteId(generateDemoUUID())
  }, [])

  // Mostrar loading mientras se carga el idioma
  if (isLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-indigo-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Cargando...</p>
        </div>
      </div>
    )
  }
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-indigo-50">
      {/* Navigation */}
      <nav className="border-b bg-white/80 backdrop-blur-sm sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center space-x-2">
              <Truck className="h-8 w-8 text-blue-600" />
              <span className="text-xl font-bold text-gray-900">TransportApp</span>
            </div>
            <div className="hidden md:flex items-center space-x-8">
              <a href="#como-funciona" className="text-gray-600 hover:text-blue-600 transition-colors">{t.nav.howItWorks}</a>
              <a href="#beneficios" className="text-gray-600 hover:text-blue-600 transition-colors">{t.nav.benefits}</a>
              <LanguageSelector
                currentLanguage={language}
                onLanguageChange={changeLanguage}
                availableLanguages={availableLanguages}
                languageNames={languageNames}
                languageFlags={languageFlags}
              />
              <Button 
                className="bg-blue-600 hover:bg-blue-700"
                onClick={() => window.open('https://calendly.com/ignaciovl-j/30min', '_blank')}
              >
                {t.nav.freeEvaluation}
              </Button>
            </div>
          </div>
        </div>
      </nav>

      {/* Hero Section */}
      <section className="relative py-20 px-4 sm:px-6 lg:px-8 particles-container">
        {/* Partículas de fondo */}
        <div className="hero-particles">
          <div className="particle"></div>
          <div className="particle"></div>
          <div className="particle"></div>
          <div className="particle"></div>
          <div className="particle"></div>
          <div className="particle"></div>
          <div className="particle"></div>
          <div className="particle"></div>
          <div className="particle"></div>
          <div className="particle"></div>
        </div>
        <div className="max-w-7xl mx-auto relative z-10">
          <div className="text-center">
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6 }}
            >
              <h1 className="text-4xl md:text-6xl font-bold text-gray-900 mb-6">
                {t.hero.title}{" "}
                <span className="text-blue-600">
                  {language === 'CL' ? 'minutos' : language === 'BR' ? 'minutos' : 'minutes'}
                </span>
              </h1>
              <p className="text-xl text-gray-600 mb-8 max-w-3xl mx-auto">
                {t.hero.subtitle}{" "}
                <span className="font-semibold text-blue-600">
                  {language === 'CL' ? 'Todo desde una sola plataforma' : 
                   language === 'BR' ? 'Tudo em uma única plataforma' : 
                   'Everything from a single platform'}
                </span>.
              </p>
              <div className="flex flex-col sm:flex-row gap-4 justify-center">
                <Button 
                  size="lg" 
                  className="bg-blue-600 hover:bg-blue-700 text-lg px-8 py-3"
                  onClick={() => window.open('https://calendly.com/ignaciovl-j/30min', '_blank')}
                >
                  {t.hero.freeEvaluation}
                  <ArrowRight className="ml-2 h-5 w-5" />
                </Button>
                <Button 
                  size="lg" 
                  variant="outline" 
                  className="text-lg px-8 py-3"
                  onClick={() => window.open('https://calendly.com/ignaciovl-j/30min', '_blank')}
                >
                  {t.hero.personalizedConsultation}
                </Button>
              </div>
            </motion.div>
            
            {/* Hero Visual */}
            <motion.div
              initial={{ opacity: 0, y: 40 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.8, delay: 0.2 }}
              className="mt-16 relative"
            >
              <div className="bg-white rounded-2xl shadow-2xl p-8 max-w-4xl mx-auto">
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                   <Card className="border-2 border-blue-100 card-glow">
                    <CardHeader>
                      <div className="flex items-center space-x-2">
                        <div className="p-2 bg-blue-100 rounded-lg">
                          <Truck className="h-6 w-6 text-blue-600" />
                        </div>
                        <CardTitle className="text-lg">{t.hero.vehicles}</CardTitle>
                      </div>
                    </CardHeader>
                                         <CardContent>
                                               <div className="space-y-3">
                          <div className="space-y-1">
                            <div className="text-sm font-medium text-gray-900">{t.hero.vehicleDetails.license}: ABC-123</div>
                            <div className="text-xs text-gray-600">{t.hero.vehicleDetails.weight}: 1000kg, {t.hero.vehicleDetails.volume}: 5000cm³</div>
                            <div className="text-xs text-gray-600">{t.hero.vehicleDetails.insurance}: $50,000</div>
                            <div className="text-xs text-gray-400">{t.hero.vehicleDetails.additionalVars}</div>
                          </div>
                          <div className="space-y-1">
                            <div className="text-sm font-medium text-gray-900">{t.hero.vehicleDetails.license}: XYZ-789</div>
                            <div className="text-xs text-gray-600">{t.hero.vehicleDetails.weight}: 800kg, {t.hero.vehicleDetails.volume}: 4000cm³</div>
                            <div className="text-xs text-gray-600">{t.hero.vehicleDetails.insurance}: $40,000</div>
                            <div className="text-xs text-gray-400">{t.hero.vehicleDetails.additionalVars}</div>
                          </div>
                          <div className="text-xs text-gray-400 text-center">
                            {t.hero.vehicleDetails.moreVehicles}
                          </div>
                        </div>
                     </CardContent>
                  </Card>
                  
                                     <Card className="border-2 border-green-100 card-glow">
                    <CardHeader>
                      <div className="flex items-center space-x-2">
                        <div className="p-2 bg-green-100 rounded-lg">
                          <MapPin className="h-6 w-6 text-green-600" />
                        </div>
                        <CardTitle className="text-lg">{t.hero.deliveries}</CardTitle>
                      </div>
                    </CardHeader>
                                         <CardContent>
                                               <div className="space-y-3">
                          <div className="space-y-1">
                            <div className="text-sm font-medium text-gray-900">{t.hero.deliveryDetails.client} A</div>
                            <div className="text-xs text-gray-600">{t.hero.deliveryDetails.address}: Las Condes</div>
                            <div className="text-xs text-gray-600">{t.hero.deliveryDetails.weight}: 15kg, {t.hero.deliveryDetails.volume}: 500cm³</div>
                            <div className="text-xs text-gray-600">{t.hero.deliveryDetails.price}: $1,200</div>
                            <div className="text-xs text-gray-400">{t.hero.deliveryDetails.additionalVars}</div>
                          </div>
                          <div className="space-y-1">
                            <div className="text-sm font-medium text-gray-900">{t.hero.deliveryDetails.client} B</div>
                            <div className="text-xs text-gray-600">{t.hero.deliveryDetails.address}: Providencia</div>
                            <div className="text-xs text-gray-600">{t.hero.deliveryDetails.weight}: 25kg, {t.hero.deliveryDetails.volume}: 800cm³</div>
                            <div className="text-xs text-gray-600">{t.hero.deliveryDetails.price}: $1,800</div>
                            <div className="text-xs text-gray-400">{t.hero.deliveryDetails.additionalVars}</div>
                          </div>
                          <div className="text-xs text-gray-400 text-center">
                            {t.hero.deliveryDetails.moreDeliveries}
                          </div>
                        </div>
                     </CardContent>
                  </Card>
                  
                                     <Card className="border-2 border-purple-100 card-glow">
                    <CardHeader>
                      <div className="flex items-center space-x-2">
                        <div className="p-2 bg-blue-100 rounded-lg">
                          <Navigation className="h-6 w-6 text-blue-600" />
                        </div>
                        <CardTitle className="text-lg">{t.hero.routes}</CardTitle>
                      </div>
                    </CardHeader>
                                         <CardContent>
                                               <div className="space-y-3">
                          <div className="space-y-1">
                            <div className="text-sm font-medium text-gray-900">{t.hero.routeDetails.route} 1</div>
                            <div className="text-xs text-gray-600">{t.hero.routeDetails.assignedLicense}: ABC-123</div>
                            <div className="text-xs text-gray-600">5 {t.hero.routeDetails.deliveries}</div>
                            <div className="text-xs text-gray-400">{t.hero.routeDetails.additionalVars}</div>
                          </div>
                          <div className="space-y-1">
                            <div className="text-sm font-medium text-gray-900">{t.hero.routeDetails.route} 2</div>
                            <div className="text-xs text-gray-600">{t.hero.routeDetails.assignedLicense}: XYZ-789</div>
                            <div className="text-xs text-gray-600">3 {t.hero.routeDetails.deliveries}</div>
                            <div className="text-xs text-gray-400">{t.hero.routeDetails.additionalVars}</div>
                          </div>
                          <div className="text-xs text-gray-400 text-center">
                            {t.hero.routeDetails.moreRoutes}
                          </div>
                        </div>
                     </CardContent>
                  </Card>
                </div>
              </div>
            </motion.div>
          </div>
        </div>
      </section>

      {/* How it Works Section */}
      <section id="como-funciona" className="py-20 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
              {t.howItWorks.title}
            </h2>
            <p className="text-xl text-gray-600 max-w-2xl mx-auto">
              {t.howItWorks.subtitle}
            </p>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6 }}
              className="text-center"
            >
              <div className="bg-blue-100 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-6">
                <span className="text-2xl font-bold text-blue-600">1</span>
              </div>
                             <h3 className="text-xl font-semibold text-gray-900 mb-4">
                 {t.howItWorks.step1.title}
               </h3>
               <p className="text-gray-600">
                 {t.howItWorks.step1.description}
               </p>
            </motion.div>
            
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 0.2 }}
              className="text-center"
            >
              <div className="bg-green-100 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-6">
                <span className="text-2xl font-bold text-green-600">2</span>
              </div>
                             <h3 className="text-xl font-semibold text-gray-900 mb-4">
                 {t.howItWorks.step2.title}
               </h3>
               <p className="text-gray-600">
                 {t.howItWorks.step2.description}
               </p>
            </motion.div>
            
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 0.4 }}
              className="text-center"
            >
              <div className="bg-purple-100 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-6">
                <span className="text-2xl font-bold text-purple-600">3</span>
              </div>
                             <h3 className="text-xl font-semibold text-gray-900 mb-4">
                 {t.howItWorks.step3.title}
               </h3>
               <p className="text-gray-600">
                 {t.howItWorks.step3.description}
               </p>
            </motion.div>
          </div>
        </div>
      </section>

            {/* Demo Section */}
      <DemoEmbed />

      {/* Benefits Section */}
      <section id="beneficios" className="py-20 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
              {t.benefits.title}
            </h2>
            <p className="text-xl text-gray-600 max-w-2xl mx-auto">
              {t.benefits.subtitle}
            </p>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6 }}
              className="text-center"
            >
              <div className="bg-red-100 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-6">
                <DollarSign className="h-8 w-8 text-red-600" />
              </div>
              <h3 className="text-xl font-semibold text-gray-900 mb-4">
                {t.benefits.reduceCosts.title}
              </h3>
              <p className="text-gray-600">
                {t.benefits.reduceCosts.description}
              </p>
            </motion.div>
            
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 0.1 }}
              className="text-center"
            >
              <div className="bg-blue-100 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-6">
                <Truck className="h-8 w-8 text-blue-600" />
              </div>
                             <h3 className="text-xl font-semibold text-gray-900 mb-4">
                 {t.benefits.balanceLoad.title}
               </h3>
               <p className="text-gray-600">
                 {t.benefits.balanceLoad.description}
               </p>
            </motion.div>
            
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 0.2 }}
              className="text-center"
            >
              <div className="bg-green-100 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-6">
                <Clock className="h-8 w-8 text-green-600" />
              </div>
              <h3 className="text-xl font-semibold text-gray-900 mb-4">
                {t.benefits.saveTime.title}
              </h3>
              <p className="text-gray-600">
                {t.benefits.saveTime.description}
              </p>
            </motion.div>
            
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 0.3 }}
              className="text-center"
            >
              <div className="bg-blue-100 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-6">
                <Zap className="h-8 w-8 text-blue-600" />
              </div>
                             <h3 className="text-xl font-semibold text-gray-900 mb-4">
                 {t.benefits.realTimeTracking.title}
               </h3>
               <p className="text-gray-600">
                 {t.benefits.realTimeTracking.description}
               </p>
            </motion.div>
          </div>
        </div>
      </section>



      {/* CTA Section */}
      <section className="py-20 bg-blue-600">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6 }}
          >
            <h2 className="text-3xl md:text-4xl font-bold text-white mb-6">
              {t.cta.title}
            </h2>
            <p className="text-xl text-blue-100 mb-8 max-w-2xl mx-auto">
              {t.cta.subtitle}
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <Button 
                size="lg" 
                className="bg-white text-blue-600 hover:bg-gray-100 text-lg px-8 py-3"
                onClick={() => window.open('https://calendly.com/ignaciovl-j/30min', '_blank')}
              >
                {t.cta.freeEvaluation}
                <ArrowRight className="ml-2 h-5 w-5" />
              </Button>
              <button 
                className="border-2 border-white text-white bg-transparent hover:bg-white hover:text-blue-600 text-lg px-8 py-3 rounded-lg font-medium transition-colors duration-200"
                onClick={() => window.open('https://calendly.com/ignaciovl-j/30min', '_blank')}
              >
                {t.cta.personalizedConsultation}
              </button>
            </div>
          </motion.div>
        </div>
      </section>

      {/* Footer */}
      <footer className="bg-gray-900 text-white py-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
            <div>
              <div className="flex items-center space-x-2 mb-4">
                <Truck className="h-8 w-8 text-blue-400" />
                <span className="text-xl font-bold">TransportApp</span>
              </div>
                              <p className="text-gray-400">
                  {t.footer.description}
                </p>
            </div>
            
            <div>
              <h3 className="font-semibold mb-4">{t.footer.product}</h3>
              <ul className="space-y-2 text-gray-400">
                <li><a href="#" className="hover:text-white transition-colors">{t.footer.howItWorks}</a></li>
                <li><a href="#" className="hover:text-white transition-colors">{t.footer.prices}</a></li>
                <li><a href="#" className="hover:text-white transition-colors">{t.footer.demo}</a></li>
              </ul>
            </div>
            
            <div>
              <h3 className="font-semibold mb-4">{t.footer.support}</h3>
              <ul className="space-y-2 text-gray-400">
                <li><a href="#" className="hover:text-white transition-colors">{t.footer.help}</a></li>
                <li><a href="#" className="hover:text-white transition-colors">{t.footer.contact}</a></li>
                <li><a href="#" className="hover:text-white transition-colors">{t.footer.documentation}</a></li>
              </ul>
            </div>
            
            <div>
              <h3 className="font-semibold mb-4">{t.footer.company}</h3>
              <ul className="space-y-2 text-gray-400">
                <li><a href="#" className="hover:text-white transition-colors">{t.footer.about}</a></li>
                <li><a href="#" className="hover:text-white transition-colors">{t.footer.blog}</a></li>
                <li><a href="#" className="hover:text-white transition-colors">{t.footer.careers}</a></li>
              </ul>
            </div>
          </div>
          
          <div className="border-t border-gray-800 mt-8 pt-8 text-center text-gray-400">
            <p>{t.footer.copyright}</p>
          </div>
        </div>
      </footer>
    </div>
  )
}
