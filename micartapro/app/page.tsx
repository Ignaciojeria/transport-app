'use client'

import { useState } from 'react'
import Image from 'next/image'
import Link from 'next/link'
import { Button } from "@/components/ui/button"
import { 
  Palette, 
  QrCode, 
  Smartphone,
  Calculator,
  ShoppingCart,
  Truck,
  MessageCircle,
  ArrowRight,
  Star
} from "lucide-react"
import { motion } from "framer-motion"
import { DemoEmbed } from "@/components/DemoEmbed"
import { openWhatsAppQuote } from "@/lib/whatsapp"
import { useLanguage } from "@/lib/useLanguage"
import { LanguageSelector } from "@/components/LanguageSelector"

export default function LandingPage() {
  const { language, changeLanguage, t, isLoading, availableLanguages, languageNames, languageFlags } = useLanguage()

  // Show loading while language is loading
  if (isLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-indigo-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading...</p>
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
            <Link href="/" className="flex items-center space-x-2">
              <Image 
                src="/logo.png" 
                alt="MiCartaPro Logo" 
                width={240} 
                height={72}
                className="h-16 md:h-20 w-auto"
              />
            </Link>
            <div className="hidden md:flex items-center space-x-8">
              <a href="#servicio" className="text-gray-600 hover:text-blue-600 transition-colors">{t.nav.service}</a>
              <a href="#beneficios" className="text-gray-600 hover:text-blue-600 transition-colors">{t.nav.benefits}</a>
              <a href="#demo" className="text-gray-600 hover:text-blue-600 transition-colors">{t.nav.demo}</a>
              <LanguageSelector
                currentLanguage={language}
                onLanguageChange={changeLanguage}
                availableLanguages={availableLanguages}
                languageNames={languageNames}
                languageFlags={languageFlags}
              />
              <Button 
                className="bg-blue-600 hover:bg-blue-700"
                onClick={() => window.location.href = `https://auth.micartapro.com?lang=${language}`}
              >
                {t.nav.signIn}
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
                {t.hero.subtitle}
              </p>
              <p className="text-xl md:text-2xl text-gray-600 mb-5 max-w-3xl mx-auto">
                {t.hero.description}
              </p>
              <div className="flex flex-col sm:flex-row gap-4 justify-center mt-3">
                <Button 
                  size="lg" 
                  className="bg-blue-600 hover:bg-blue-700 text-lg px-8 py-3"
                  onClick={() => window.location.href = `https://auth.micartapro.com?lang=${language}`}
                >
                  {t.hero.startFree}
                  <ArrowRight className="ml-2 h-5 w-5" />
                </Button>
                <Button 
                  size="lg" 
                  variant="outline" 
                  className="text-lg px-8 py-3"
                  onClick={() => document.getElementById('demo')?.scrollIntoView({ behavior: 'smooth' })}
                >
                  {t.hero.viewDemo}
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
              {t.service.title}
            </h2>
            <p className="text-xl text-gray-600 max-w-2xl mx-auto">
              {t.service.subtitle}
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
                {t.service.customDesign}
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
                {t.service.customLogo}
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
                {t.service.exclusiveQR}
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
                {t.service.responsiveDesign}
              </h3>
            </motion.div>
          </div>

          {/* SaaS Description */}
          <div className="mb-12 max-w-3xl mx-auto">
            <div className="bg-blue-50 rounded-lg p-6 border border-blue-100">
              <p className="text-gray-700 text-center leading-relaxed">
                {t.service.saasDescription}
              </p>
            </div>
          </div>

          {/* Precio */}
          <Link href="/pricing" className="block group">
            <motion.div 
              className="bg-gradient-to-br from-blue-600 to-indigo-700 rounded-2xl p-8 md:p-12 text-center text-white cursor-pointer transform transition-all duration-300 hover:scale-105 hover:shadow-2xl"
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
            >
              <h3 className="text-2xl md:text-3xl font-bold mb-4">
                {t.service.pricingTitle}
              </h3>
              <div className="space-y-2 mb-6">
                <p className="text-xl">{t.service.firstYearFree}</p>
                <p className="text-lg opacity-90">{t.service.renewalPrice}</p>
              </div>
              <div className="flex flex-col sm:flex-row gap-4 justify-center">
                <Button 
                  size="lg" 
                  className="bg-white text-blue-600 hover:bg-gray-100 text-lg px-8 py-3"
                  onClick={(e) => {
                    e.preventDefault()
                    window.location.href = `https://auth.micartapro.com?lang=${language}`
                  }}
                >
                  {t.service.startFreeButton}
                  <ArrowRight className="ml-2 h-5 w-5" />
                </Button>
                <Button 
                  size="lg" 
                  variant="outline"
                  className="bg-white border-2 border-white text-blue-600 hover:bg-gray-100 text-lg px-8 py-3"
                  onClick={(e) => {
                    e.preventDefault()
                    window.location.href = '/pricing'
                  }}
                >
                  {t.service.viewPricing}
                </Button>
              </div>
            </motion.div>
          </Link>
        </div>
      </section>

      {/* Beneficios Section */}
      <section id="beneficios" className="py-20 bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
              {t.benefits.title}
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
                    {t.benefits.autoCostCalculation.title}
                  </h3>
                  <p className="text-gray-600">
                    {t.benefits.autoCostCalculation.description}
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
                    {t.benefits.shoppingCart.title}
                  </h3>
                  <p className="text-gray-600">
                    {t.benefits.shoppingCart.description}
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
                    {t.benefits.deliveryOptions.title}
                  </h3>
                  <p className="text-gray-600">
                    {t.benefits.deliveryOptions.description}
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
                    {t.benefits.whatsappOrders.title}
                  </h3>
                  <p className="text-gray-600">
                    {t.benefits.whatsappOrders.description}
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
              {t.cta.title}
            </h2>
            <p className="text-xl text-blue-100 mb-8 max-w-2xl mx-auto">
              {t.cta.subtitle}
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <Button 
                size="lg" 
                className="bg-white text-blue-600 hover:bg-gray-100 text-lg px-8 py-3"
                onClick={() => window.location.href = `https://auth.micartapro.com?lang=${language}`}
              >
                {t.cta.startFreeButton}
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
                {t.footer.description}
              </p>
            </div>
            
            <div>
              <h3 className="font-semibold mb-4">{t.footer.service}</h3>
              <ul className="space-y-2 text-gray-400">
                <li><a href="#servicio" className="hover:text-white transition-colors">{t.footer.ourService}</a></li>
                <li><a href="#beneficios" className="hover:text-white transition-colors">{t.footer.benefits}</a></li>
                <li><a href="#demo" className="hover:text-white transition-colors">{t.footer.demo}</a></li>
              </ul>
            </div>
            
            <div>
              <h3 className="font-semibold mb-4">{t.footer.contact}</h3>
              <ul className="space-y-2 text-gray-400">
                <li>
                  <button 
                    onClick={openWhatsAppQuote}
                    className="hover:text-white transition-colors text-left"
                  >
                    {t.footer.quoteWhatsApp}
                  </button>
                </li>
              </ul>
            </div>
          </div>
          
          <div className="border-t border-gray-800 mt-8 pt-8 text-center text-gray-400">
            <div className="flex flex-col sm:flex-row justify-center items-center gap-4 mb-4">
              <Link href="/privacy" className="hover:text-white transition-colors">
                {t.footer.privacy}
              </Link>
              <span className="hidden sm:inline">•</span>
              <Link href="/terms" className="hover:text-white transition-colors">
                {t.footer.terms}
              </Link>
              <span className="hidden sm:inline">•</span>
              <Link href="/refund" className="hover:text-white transition-colors">
                {t.footer.refund}
              </Link>
            </div>
            <p>&copy; {new Date().getFullYear()} MiCartaPro. {t.footer.copyright}</p>
          </div>
        </div>
      </footer>
    </div>
  )
}
