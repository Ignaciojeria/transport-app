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
  
  // Generar UUID solo en el cliente para evitar errores de hidratación
  useEffect(() => {
    setRouteId(generateDemoUUID())
  }, [])
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
              <a href="#como-funciona" className="text-gray-600 hover:text-blue-600 transition-colors">Cómo funciona</a>
              <a href="#beneficios" className="text-gray-600 hover:text-blue-600 transition-colors">Beneficios</a>
              <Button 
                className="bg-blue-600 hover:bg-blue-700"
                onClick={() => window.open('https://calendly.com/ignaciovl-j/30min', '_blank')}
              >
                Evaluación Gratuita
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
                Optimiza tu flota en{" "}
                <span className="text-blue-600">minutos</span>
              </h1>
              <p className="text-xl text-gray-600 mb-8 max-w-3xl mx-auto">
                Optimiza rutas, genera enlaces para conductores y monitorea entregas en tiempo real. 
                <span className="font-semibold text-blue-600">Todo desde una sola plataforma</span>.
              </p>
              <div className="flex flex-col sm:flex-row gap-4 justify-center">
                <Button 
                  size="lg" 
                  className="bg-blue-600 hover:bg-blue-700 text-lg px-8 py-3"
                  onClick={() => window.open('https://calendly.com/ignaciovl-j/30min', '_blank')}
                >
                  Evaluación Gratuita
                  <ArrowRight className="ml-2 h-5 w-5" />
                </Button>
                <Button 
                  size="lg" 
                  variant="outline" 
                  className="text-lg px-8 py-3"
                  onClick={() => window.open('https://calendly.com/ignaciovl-j/30min', '_blank')}
                >
                  Consulta Personalizada
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
                        <CardTitle className="text-lg">Vehículos</CardTitle>
                      </div>
                    </CardHeader>
                                         <CardContent>
                                               <div className="space-y-3">
                          <div className="space-y-1">
                            <div className="text-sm font-medium text-gray-900">Patente: ABC-123</div>
                            <div className="text-xs text-gray-600">Peso: 1000kg, Vol: 5000cm³</div>
                            <div className="text-xs text-gray-600">Seguro: $50,000</div>
                            <div className="text-xs text-gray-400">Variables adicionales...</div>
                          </div>
                          <div className="space-y-1">
                            <div className="text-sm font-medium text-gray-900">Patente: XYZ-789</div>
                            <div className="text-xs text-gray-600">Peso: 800kg, Vol: 4000cm³</div>
                            <div className="text-xs text-gray-600">Seguro: $40,000</div>
                            <div className="text-xs text-gray-400">Variables adicionales...</div>
                          </div>
                          <div className="text-xs text-gray-400 text-center">
                            + 3 vehículos más...
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
                        <CardTitle className="text-lg">Entregas</CardTitle>
                      </div>
                    </CardHeader>
                                         <CardContent>
                                               <div className="space-y-3">
                          <div className="space-y-1">
                            <div className="text-sm font-medium text-gray-900">Cliente A</div>
                            <div className="text-xs text-gray-600">Dirección: Las Condes</div>
                            <div className="text-xs text-gray-600">Peso: 15kg, Vol: 500cm³</div>
                            <div className="text-xs text-gray-600">Precio: $1,200</div>
                            <div className="text-xs text-gray-400">Variables adicionales...</div>
                          </div>
                          <div className="space-y-1">
                            <div className="text-sm font-medium text-gray-900">Cliente B</div>
                            <div className="text-xs text-gray-600">Dirección: Providencia</div>
                            <div className="text-xs text-gray-600">Peso: 25kg, Vol: 800cm³</div>
                            <div className="text-xs text-gray-600">Precio: $1,800</div>
                            <div className="text-xs text-gray-400">Variables adicionales...</div>
                          </div>
                          <div className="text-xs text-gray-400 text-center">
                            + 8 entregas más...
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
                        <CardTitle className="text-lg">Rutas</CardTitle>
                      </div>
                    </CardHeader>
                                         <CardContent>
                                               <div className="space-y-3">
                          <div className="space-y-1">
                            <div className="text-sm font-medium text-gray-900">Ruta 1</div>
                            <div className="text-xs text-gray-600">Patente asignada: ABC-123</div>
                            <div className="text-xs text-gray-600">5 entregas</div>
                            <div className="text-xs text-gray-400">Variables adicionales...</div>
                          </div>
                          <div className="space-y-1">
                            <div className="text-sm font-medium text-gray-900">Ruta 2</div>
                            <div className="text-xs text-gray-600">Patente asignada: XYZ-789</div>
                            <div className="text-xs text-gray-600">3 entregas</div>
                            <div className="text-xs text-gray-400">Variables adicionales...</div>
                          </div>
                          <div className="text-xs text-gray-400 text-center">
                            + 2 rutas más...
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
              ¿Cómo funciona?
            </h2>
            <p className="text-xl text-gray-600 max-w-2xl mx-auto">
              En solo 3 pasos simples, optimiza tu logística y reduce costos
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
                 Crea tus vehículos
               </h3>
               <p className="text-gray-600">
                 Configura y registra tus vehículos con sus capacidades, dimensiones y características. 
                 Define la flota que utilizarás para tus entregas.
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
                 Carga las entregas
               </h3>
               <p className="text-gray-600">
                 Crea y configura tus entregas con destinos, productos y restricciones. 
                 Define todos los detalles necesarios para la planificación de rutas.
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
                 Optimiza y ejecuta
               </h3>
               <p className="text-gray-600">
                 Llama a nuestro agente de optimización y recibe los enlaces para el conductor. 
                 Ejecuta las rutas optimizadas con trazabilidad en tiempo real.
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
              Beneficios para empresas
            </h2>
            <p className="text-xl text-gray-600 max-w-2xl mx-auto">
              Optimiza tu operación logística y reduce costos significativamente
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
                Reduce costos logísticos
              </h3>
              <p className="text-gray-600">
                Optimiza rutas y reduce combustible, tiempo y recursos.
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
                 Balancea carga automáticamente
               </h3>
               <p className="text-gray-600">
                 Distribuye la carga de manera óptima entre vehículos, maximizando eficiencia y minimizando viajes vacíos
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
                Ahorra tiempo planificando
              </h3>
              <p className="text-gray-600">
                De 2 horas a 5 minutos. Planifica rutas automáticamente
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
                 Trazabilidad en tiempo real
               </h3>
               <p className="text-gray-600">
                 Monitorea el progreso de entregas en tiempo real desde la web mobile, sin necesidad de software adicional
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
              ¿Listo para optimizar tu flota?
            </h2>
            <p className="text-xl text-blue-100 mb-8 max-w-2xl mx-auto">
              Agenda una evaluación gratuita de tus procesos logísticos
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <Button 
                size="lg" 
                className="bg-white text-blue-600 hover:bg-gray-100 text-lg px-8 py-3"
                onClick={() => window.open('https://calendly.com/ignaciovl-j/30min', '_blank')}
              >
                Evaluación Gratuita
                <ArrowRight className="ml-2 h-5 w-5" />
              </Button>
              <button 
                className="border-2 border-white text-white bg-transparent hover:bg-white hover:text-blue-600 text-lg px-8 py-3 rounded-lg font-medium transition-colors duration-200"
                onClick={() => window.open('https://calendly.com/ignaciovl-j/30min', '_blank')}
              >
                Consulta Personalizada
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
                  Optimiza rutas, genera enlaces para conductores y monitorea entregas en tiempo real. Todo desde una sola plataforma.
                </p>
            </div>
            
            <div>
              <h3 className="font-semibold mb-4">Producto</h3>
              <ul className="space-y-2 text-gray-400">
                <li><a href="#" className="hover:text-white transition-colors">Cómo funciona</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Precios</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Demo</a></li>
              </ul>
            </div>
            
            <div>
              <h3 className="font-semibold mb-4">Soporte</h3>
              <ul className="space-y-2 text-gray-400">
                <li><a href="#" className="hover:text-white transition-colors">Ayuda</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Contacto</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Documentación</a></li>
              </ul>
            </div>
            
            <div>
              <h3 className="font-semibold mb-4">Empresa</h3>
              <ul className="space-y-2 text-gray-400">
                <li><a href="#" className="hover:text-white transition-colors">Acerca de</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Blog</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Carreras</a></li>
              </ul>
            </div>
          </div>
          
          <div className="border-t border-gray-800 mt-8 pt-8 text-center text-gray-400">
            <p>&copy; 2025 TransportApp. Todos los derechos reservados.</p>
          </div>
        </div>
      </footer>
    </div>
  )
}
