'use client'

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

export default function LandingPage() {
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
              <a href="#testimonios" className="text-gray-600 hover:text-blue-600 transition-colors">Testimonios</a>
              <Button className="bg-blue-600 hover:bg-blue-700">Probar Gratis</Button>
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
                <Button size="lg" className="bg-blue-600 hover:bg-blue-700 text-lg px-8 py-3">
                  Probar Gratis
                  <ArrowRight className="ml-2 h-5 w-5" />
                </Button>
                <Button size="lg" variant="outline" className="text-lg px-8 py-3">
                  Solicitar Demo
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
                           <div className="text-sm font-medium text-gray-900">ABC-123</div>
                           <div className="text-xs text-gray-600">Peso: 1000g, Vol: 5000cm³</div>
                           <div className="text-xs text-gray-600">Seguro: $50,000</div>
                         </div>
                         <div className="space-y-1">
                           <div className="text-sm font-medium text-gray-900">XYZ-789</div>
                           <div className="text-xs text-gray-600">Peso: 800g, Vol: 4000cm³</div>
                           <div className="text-xs text-gray-600">Seguro: $40,000</div>
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
                           <div className="text-xs text-gray-600">Peso: 15g, Vol: 500cm³</div>
                           <div className="text-xs text-gray-600">Precio: $1,200</div>
                         </div>
                         <div className="space-y-1">
                           <div className="text-sm font-medium text-gray-900">Cliente B</div>
                           <div className="text-xs text-gray-600">Peso: 25g, Vol: 800cm³</div>
                           <div className="text-xs text-gray-600">Precio: $1,800</div>
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
                         </div>
                         <div className="space-y-1">
                           <div className="text-sm font-medium text-gray-900">Ruta 2</div>
                           <div className="text-xs text-gray-600">Patente asignada: XYZ-789</div>
                           <div className="text-xs text-gray-600">3 entregas</div>
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
                 Planifica desde Google Sheets
               </h3>
               <p className="text-gray-600">
                 Carga tus vehículos y entregas directamente en Google Sheets para la planificación. 
                 No necesitas instalar software ni aprender nuevas herramientas.
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
                 Balancea carga y optimiza rutas
               </h3>
               <p className="text-gray-600">
                 Balancea automáticamente la carga entre vehículos y genera el orden óptimo 
                 de entregas considerando peso, volumen, seguro y restricciones de tiempo.
               </p>
            </motion.div>
            
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 0.4 }}
              className="text-center"
            >
              <div className="bg-blue-100 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-6">
                <span className="text-2xl font-bold text-blue-600">3</span>
              </div>
                             <h3 className="text-xl font-semibold text-gray-900 mb-4">
                 Ejecuta desde web mobile
               </h3>
               <p className="text-gray-600">
                 Comparte enlaces con conductores para acceder a la web mobile, obtén trazabilidad en tiempo real 
                 y genera reportes de entrega automáticamente.
               </p>
            </motion.div>
          </div>
        </div>
      </section>

      {/* Driver View Section */}
      <section className="py-20 bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
            <motion.div
              initial={{ opacity: 0, x: -20 }}
              whileInView={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.6 }}
            >
              <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-6">
                Vista del conductor
              </h2>
              <p className="text-xl text-gray-600 mb-8">
                Tus conductores tienen una experiencia simple y funcional
              </p>
              
              <div className="space-y-6">
                <div className="flex items-start space-x-4">
                  <div className="bg-blue-100 p-3 rounded-lg">
                    <Smartphone className="h-6 w-6 text-blue-600" />
                  </div>
                  <div>
                    <h3 className="font-semibold text-gray-900 mb-2">Abre el enlace</h3>
                    <p className="text-gray-600">Recibe un enlace directo en su teléfono</p>
                  </div>
                </div>
                
                <div className="flex items-start space-x-4">
                  <div className="bg-green-100 p-3 rounded-lg">
                    <Navigation className="h-6 w-6 text-green-600" />
                  </div>
                  <div>
                    <h3 className="font-semibold text-gray-900 mb-2">Ve la ruta optimizada</h3>
                    <p className="text-gray-600">Visualiza todas las paradas en orden óptimo</p>
                  </div>
                </div>
                
                <div className="flex items-start space-x-4">
                  <div className="bg-blue-100 p-3 rounded-lg">
                    <MapPin className="h-6 w-6 text-blue-600" />
                  </div>
                  <div>
                    <h3 className="font-semibold text-gray-900 mb-2">Usa Google Maps</h3>
                    <p className="text-gray-600">Navega directamente con Google Maps integrado</p>
                  </div>
                </div>
                
                                 <div className="flex items-start space-x-4">
                   <div className="bg-orange-100 p-3 rounded-lg">
                     <CheckCircle className="h-6 w-6 text-orange-600" />
                   </div>
                   <div>
                                        <h3 className="font-semibold text-gray-900 mb-2">Trazabilidad en tiempo real</h3>
                   <p className="text-gray-600">Monitorea y confirma entregas con actualizaciones instantáneas desde web mobile</p>
                   </div>
                 </div>
              </div>
            </motion.div>
            
            <motion.div
              initial={{ opacity: 0, x: 20 }}
              whileInView={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.6, delay: 0.2 }}
              className="relative"
            >
              <div className="bg-white rounded-2xl shadow-2xl p-6">
                <div className="bg-gray-900 rounded-lg p-4 mb-4">
                  <div className="flex items-center space-x-2 mb-3">
                    <div className="w-3 h-3 bg-red-500 rounded-full"></div>
                    <div className="w-3 h-3 bg-yellow-500 rounded-full"></div>
                    <div className="w-3 h-3 bg-green-500 rounded-full"></div>
                  </div>
                  <div className="text-white text-sm">TransportApp - Conductor</div>
                </div>
                
                <div className="space-y-4">
                  <div className="border rounded-lg p-4">
                    <div className="flex items-center justify-between mb-2">
                      <span className="font-medium">Ruta del día</span>
                      <span className="text-sm text-gray-500">8 entregas</span>
                    </div>
                    <div className="space-y-2">
                      <div className="flex items-center space-x-3">
                        <div className="w-6 h-6 bg-blue-100 rounded-full flex items-center justify-center">
                          <span className="text-xs font-bold text-blue-600">1</span>
                        </div>
                                                 <div className="flex-1">
                           <div className="font-medium text-sm">Cliente A - Centro</div>
                           <div className="text-xs text-gray-500">15g, 500cm³ - 10:00 AM</div>
                         </div>
                        <CheckCircle className="h-5 w-5 text-green-500" />
                      </div>
                      
                      <div className="flex items-center space-x-3">
                        <div className="w-6 h-6 bg-blue-100 rounded-full flex items-center justify-center">
                          <span className="text-xs font-bold text-blue-600">2</span>
                        </div>
                                                 <div className="flex-1">
                           <div className="font-medium text-sm">Cliente B - Norte</div>
                           <div className="text-xs text-gray-500">25g, 800cm³ - 10:30 AM</div>
                         </div>
                        <div className="w-5 h-5 border-2 border-gray-300 rounded-full"></div>
                      </div>
                      
                      <div className="flex items-center space-x-3">
                        <div className="w-6 h-6 bg-gray-100 rounded-full flex items-center justify-center">
                          <span className="text-xs font-bold text-gray-600">3</span>
                        </div>
                                                 <div className="flex-1">
                           <div className="font-medium text-sm text-gray-400">Cliente C - Sur</div>
                           <div className="text-xs text-gray-400">20g, 600cm³ - 11:00 AM</div>
                         </div>
                        <div className="w-5 h-5 border-2 border-gray-300 rounded-full"></div>
                      </div>
                    </div>
                  </div>
                  
                  <Button className="w-full bg-blue-600 hover:bg-blue-700">
                    <Navigation className="mr-2 h-4 w-4" />
                    Abrir en Google Maps
                  </Button>
                </div>
              </div>
            </motion.div>
          </div>
        </div>
      </section>

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

      {/* Testimonials Section */}
      <section id="testimonios" className="py-20 bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
              Lo que dicen nuestros clientes
            </h2>
            <p className="text-xl text-gray-600 max-w-2xl mx-auto">
              Empresas que ya optimizaron su logística con TransportApp
            </p>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6 }}
            >
              <Card className="h-full">
                <CardContent className="p-6">
                  <div className="flex items-center mb-4">
                    {[...Array(5)].map((_, i) => (
                      <Star key={i} className="h-5 w-5 text-yellow-400 fill-current" />
                    ))}
                  </div>
                  <Quote className="h-8 w-8 text-blue-600 mb-4" />
                                     <p className="text-gray-600 mb-6">
                     "El balanceo automático de carga y la trazabilidad en tiempo real nos han revolucionado la logística. Planificamos desde Sheets."
                   </p>
                  <div className="flex items-center">
                    <div className="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center mr-4">
                      <Users className="h-6 w-6 text-blue-600" />
                    </div>
                    <div>
                      <div className="font-semibold text-gray-900">Carlos Mendoza</div>
                      <div className="text-sm text-gray-500">Operador Logístico</div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </motion.div>
            
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 0.1 }}
            >
              <Card className="h-full">
                <CardContent className="p-6">
                  <div className="flex items-center mb-4">
                    {[...Array(5)].map((_, i) => (
                      <Star key={i} className="h-5 w-5 text-yellow-400 fill-current" />
                    ))}
                  </div>
                  <Quote className="h-8 w-8 text-green-600 mb-4" />
                                     <p className="text-gray-600 mb-6">
                     "La optimización de rutas y el balanceo de carga nos redujeron costos 25% y mejoraron la satisfacción de conductores."
                   </p>
                  <div className="flex items-center">
                    <div className="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center mr-4">
                      <Truck className="h-6 w-6 text-green-600" />
                    </div>
                    <div>
                      <div className="font-semibold text-gray-900">María González</div>
                      <div className="text-sm text-gray-500">Gerente de Flota</div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </motion.div>
            
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 0.2 }}
            >
              <Card className="h-full">
                <CardContent className="p-6">
                  <div className="flex items-center mb-4">
                    {[...Array(5)].map((_, i) => (
                      <Star key={i} className="h-5 w-5 text-yellow-400 fill-current" />
                    ))}
                  </div>
                  <Quote className="h-8 w-8 text-blue-600 mb-4" />
                                     <p className="text-gray-600 mb-6">
                     "La planificación desde Google Sheets es perfecta. La web mobile para conductores es muy intuitiva."
                   </p>
                  <div className="flex items-center">
                    <div className="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center mr-4">
                      <BarChart3 className="h-6 w-6 text-blue-600" />
                    </div>
                    <div>
                      <div className="font-semibold text-gray-900">Roberto Silva</div>
                      <div className="text-sm text-gray-500">Director de Operaciones</div>
                    </div>
                  </div>
                </CardContent>
              </Card>
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
              Únete a las empresas que ya están ahorrando tiempo y dinero con TransportApp
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <Button size="lg" className="bg-white text-blue-600 hover:bg-gray-100 text-lg px-8 py-3">
                Solicitar Acceso
                <ArrowRight className="ml-2 h-5 w-5" />
              </Button>
              <Button size="lg" variant="outline" className="border-white text-white hover:bg-white hover:text-blue-600 text-lg px-8 py-3">
                Agendar Demo
              </Button>
              <Button size="lg" variant="outline" className="border-white text-white hover:bg-white hover:text-blue-600 text-lg px-8 py-3">
                Ver Precios
              </Button>
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
                Optimiza tu flota en minutos. Todo desde Google Sheets.
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
            <p>&copy; 2024 TransportApp. Todos los derechos reservados.</p>
          </div>
        </div>
      </footer>
    </div>
  )
}
