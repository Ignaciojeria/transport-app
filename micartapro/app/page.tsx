'use client'

import { useState } from 'react'
import Image from 'next/image'
import Link from 'next/link'
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { 
  Menu, 
  Palette, 
  QrCode, 
  Smartphone,
  Calculator,
  ShoppingCart,
  Truck,
  MessageCircle,
  ArrowRight,
  Star,
  CheckCircle,
  Sparkles
} from "lucide-react"
import { motion } from "framer-motion"
import { DemoEmbed } from "@/components/DemoEmbed"
import { openWhatsAppQuote } from "@/lib/whatsapp"

export default function LandingPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-indigo-50">
      {/* Navigation */}
      <nav className="border-b bg-white/80 backdrop-blur-sm sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center space-x-2">
              <Image 
                src="/logo.png" 
                alt="MiCartaPro Logo" 
                width={240} 
                height={72}
                className="h-16 md:h-20 w-auto"
              />
            </div>
            <div className="hidden md:flex items-center space-x-8">
              <a href="#servicio" className="text-gray-600 hover:text-blue-600 transition-colors">Servicio</a>
              <a href="#beneficios" className="text-gray-600 hover:text-blue-600 transition-colors">Beneficios</a>
              <a href="#demo" className="text-gray-600 hover:text-blue-600 transition-colors">Demo</a>
              <Button 
                className="bg-blue-600 hover:bg-blue-700"
                onClick={openWhatsAppQuote}
              >
                Cotizar
              </Button>
            </div>
          </div>
        </div>
      </nav>

      {/* Hero Section */}
      <section className="relative pt-0 pb-12 md:pt-2 md:pb-16 px-4 sm:px-6 lg:px-8 particles-container">
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
        <div className="max-w-7xl mx-auto relative z-10 -mt-8 md:-mt-10">
          <div className="text-center">
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6 }}
            >
              <div className="flex items-center justify-center -mb-8 md:-mb-10">
                <Image 
                  src="/logo.png" 
                  alt="MiCartaPro Logo" 
                  width={600} 
                  height={180}
                  className="h-40 md:h-52 lg:h-64 w-auto"
                />
              </div>
              <p className="text-3xl md:text-4xl lg:text-5xl text-gray-700 -mt-6 md:-mt-8 mb-4 font-semibold">
                Tu Menú Digital, Sin Complicaciones
              </p>
              <p className="text-xl md:text-2xl text-gray-600 mb-5 max-w-3xl mx-auto">
                Gestiona tu menú digital y deja que las ventas fluyan.
              </p>
              <div className="flex flex-col sm:flex-row gap-4 justify-center mt-3">
                <Button 
                  size="lg" 
                  className="bg-blue-600 hover:bg-blue-700 text-lg px-8 py-3"
                  onClick={openWhatsAppQuote}
                >
                  Cotizar Ahora
                  <ArrowRight className="ml-2 h-5 w-5" />
                </Button>
                <Button 
                  size="lg" 
                  variant="outline" 
                  className="text-lg px-8 py-3"
                  onClick={() => document.getElementById('demo')?.scrollIntoView({ behavior: 'smooth' })}
                >
                  Ver Demo
                </Button>
              </div>
            </motion.div>
          </div>
        </div>
      </section>

      {/* Servicio Section */}
      <section id="servicio" className="pt-8 pb-28 md:pt-12 md:pb-36 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
              🚀 Nuestro Servicio
            </h2>
            <p className="text-xl text-gray-600 max-w-2xl mx-auto">
              En MiCartaPro lo hacemos diferente
            </p>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8 mb-12">
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6 }}
              className="text-center"
            >
              <div className="bg-blue-100 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-6">
                <Palette className="h-8 w-8 text-blue-600" />
              </div>
              <h3 className="text-xl font-semibold text-gray-900 mb-4">
                🎨 Diseño 100% personalizado
              </h3>
            </motion.div>
            
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 0.1 }}
              className="text-center"
            >
              <div className="bg-purple-100 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-6">
                <Star className="h-8 w-8 text-purple-600" />
              </div>
              <h3 className="text-xl font-semibold text-gray-900 mb-4">
                ✨ Logo a medida
              </h3>
            </motion.div>
            
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 0.2 }}
              className="text-center"
            >
              <div className="bg-green-100 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-6">
                <QrCode className="h-8 w-8 text-green-600" />
              </div>
              <h3 className="text-xl font-semibold text-gray-900 mb-4">
                🔗 Código QR exclusivo
              </h3>
            </motion.div>
            
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 0.3 }}
              className="text-center"
            >
              <div className="bg-orange-100 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-6">
                <Smartphone className="h-8 w-8 text-orange-600" />
              </div>
              <h3 className="text-xl font-semibold text-gray-900 mb-4">
                📱 Diseño responsivo para todos los dispositivos
              </h3>
            </motion.div>
          </div>

          {/* Precio */}
          <div className="bg-gradient-to-br from-blue-600 to-indigo-700 rounded-2xl p-8 md:p-12 text-center text-white">
            <h3 className="text-2xl md:text-3xl font-bold mb-4">
              Oferta única con cupos limitados — desde $150 USD
            </h3>
            <div className="space-y-2 mb-6">
              <p className="text-xl">✅ Primer año gratis</p>
              <p className="text-lg opacity-90">Renovación desde el segundo año: $10 USD mensuales</p>
            </div>
            <Button 
              size="lg" 
              className="bg-white text-blue-600 hover:bg-gray-100 text-lg px-8 py-3"
              onClick={openWhatsAppQuote}
            >
              Cotizar Ahora
              <ArrowRight className="ml-2 h-5 w-5" />
            </Button>
          </div>
        </div>
      </section>

      {/* Beneficios Section */}
      <section id="beneficios" className="py-20 bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
              🎯 Beneficios incluidos
            </h2>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6 }}
              className="bg-white rounded-lg shadow-lg p-8"
            >
              <div className="flex items-start space-x-4">
                <div className="bg-green-100 w-12 h-12 rounded-full flex items-center justify-center flex-shrink-0">
                  <Calculator className="h-6 w-6 text-green-600" />
                </div>
                <div>
                  <h3 className="text-xl font-semibold text-gray-900 mb-2">
                    💰 Cálculo de costos automático
                  </h3>
                  <p className="text-gray-600">
                    Olvídate del cálculo manual. Tu carta procesa y muestra el costo total de cada plato.
                  </p>
                </div>
              </div>
            </motion.div>
            
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 0.1 }}
              className="bg-white rounded-lg shadow-lg p-8"
            >
              <div className="flex items-start space-x-4">
                <div className="bg-blue-100 w-12 h-12 rounded-full flex items-center justify-center flex-shrink-0">
                  <ShoppingCart className="h-6 w-6 text-blue-600" />
                </div>
                <div>
                  <h3 className="text-xl font-semibold text-gray-900 mb-2">
                    🛒 Carrito de compras integrado
                  </h3>
                  <p className="text-gray-600">
                    Permite que tus clientes armen su pedido de manera simple, ordenada y rápida.
                  </p>
                </div>
              </div>
            </motion.div>
            
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 0.2 }}
              className="bg-white rounded-lg shadow-lg p-8"
            >
              <div className="flex items-start space-x-4">
                <div className="bg-purple-100 w-12 h-12 rounded-full flex items-center justify-center flex-shrink-0">
                  <Truck className="h-6 w-6 text-purple-600" />
                </div>
                <div>
                  <h3 className="text-xl font-semibold text-gray-900 mb-2">
                    🚚 Envío o retiro en tienda
                  </h3>
                  <p className="text-gray-600">
                    Tu carta pregunta automáticamente por los detalles necesarios para completar el pedido.
                  </p>
                </div>
              </div>
            </motion.div>
            
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 0.3 }}
              className="bg-white rounded-lg shadow-lg p-8"
            >
              <div className="flex items-start space-x-4">
                <div className="bg-green-100 w-12 h-12 rounded-full flex items-center justify-center flex-shrink-0">
                  <MessageCircle className="h-6 w-6 text-green-600" />
                </div>
                <div>
                  <h3 className="text-xl font-semibold text-gray-900 mb-2">
                    📩 Recepción de pedidos por WhatsApp
                  </h3>
                  <p className="text-gray-600">
                    Recibe los pedidos de forma ordenada, clara y transparente tanto para la cocina como para tus clientes.
                  </p>
                </div>
              </div>
            </motion.div>
          </div>
        </div>
      </section>

      {/* Demo Section */}
      <div id="demo">
        <DemoEmbed />
      </div>

      {/* CTA Section */}
      <section className="py-20 bg-blue-600">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6 }}
          >
            <h2 className="text-3xl md:text-4xl font-bold text-white mb-6">
              ¿Listo para digitalizar tu menú?
            </h2>
            <p className="text-xl text-blue-100 mb-8 max-w-2xl mx-auto">
              Contáctanos ahora y obtén tu menú digital personalizado
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <Button 
                size="lg" 
                className="bg-white text-blue-600 hover:bg-gray-100 text-lg px-8 py-3"
                onClick={openWhatsAppQuote}
              >
                Cotizar Ahora
                <ArrowRight className="ml-2 h-5 w-5" />
              </Button>
            </div>
          </motion.div>
        </div>
      </section>

      {/* Footer */}
      <footer className="bg-gray-900 text-white py-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <div>
              <div className="flex items-center space-x-2 mb-4">
                <Image 
                  src="/logo.png" 
                  alt="MiCartaPro Logo" 
                  width={180} 
                  height={60}
                  className="h-12 w-auto"
                />
              </div>
              <p className="text-gray-400">
                Tu menú digital, sin complicaciones. Gestiona tu restaurante y deja que las ventas fluyan.
              </p>
            </div>
            
            <div>
              <h3 className="font-semibold mb-4">Servicio</h3>
              <ul className="space-y-2 text-gray-400">
                <li><a href="#servicio" className="hover:text-white transition-colors">Nuestro Servicio</a></li>
                <li><a href="#beneficios" className="hover:text-white transition-colors">Beneficios</a></li>
                <li><a href="#demo" className="hover:text-white transition-colors">Demo</a></li>
              </ul>
            </div>
            
            <div>
              <h3 className="font-semibold mb-4">Contacto</h3>
              <ul className="space-y-2 text-gray-400">
                <li>
                  <button 
                    onClick={openWhatsAppQuote}
                    className="hover:text-white transition-colors text-left"
                  >
                    Cotizar por WhatsApp
                  </button>
                </li>
              </ul>
            </div>
          </div>
          
          <div className="border-t border-gray-800 mt-8 pt-8 text-center text-gray-400">
            <div className="flex flex-col sm:flex-row justify-center items-center gap-4 mb-4">
              <Link href="/privacy" className="hover:text-white transition-colors">
                Política de Privacidad
              </Link>
              <span className="hidden sm:inline">•</span>
              <Link href="/terms" className="hover:text-white transition-colors">
                Términos y Condiciones
              </Link>
            </div>
            <p>&copy; {new Date().getFullYear()} MiCartaPro. Todos los derechos reservados.</p>
          </div>
        </div>
      </footer>
    </div>
  )
}

