'use client'

import { useState } from 'react'
import Image from 'next/image'
import Link from 'next/link'
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { 
  Check, 
  ArrowRight, 
  CheckCircle2,
  ChevronDown,
  ChevronUp
} from "lucide-react"
import { motion } from "framer-motion"
import { useLanguage } from "@/lib/useLanguage"
import { LanguageSelector } from "@/components/LanguageSelector"
import { openWhatsAppQuote } from "@/lib/whatsapp"

export default function PricingPage() {
  const { language, changeLanguage, t, isLoading, availableLanguages, languageNames, languageFlags } = useLanguage()
  const [openFaqIndex, setOpenFaqIndex] = useState<number | null>(null)

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
              <Link href="/#servicio" className="text-gray-600 hover:text-blue-600 transition-colors">{t.nav.service}</Link>
              <Link href="/#beneficios" className="text-gray-600 hover:text-blue-600 transition-colors">{t.nav.benefits}</Link>
              <Link href="/#demo" className="text-gray-600 hover:text-blue-600 transition-colors">{t.nav.demo}</Link>
              <LanguageSelector
                currentLanguage={language}
                onLanguageChange={changeLanguage}
                availableLanguages={availableLanguages}
                languageNames={languageNames}
                languageFlags={languageFlags}
              />
              <Button 
                className="bg-blue-600 hover:bg-blue-700"
                onClick={openWhatsAppQuote}
              >
                {t.nav.quote}
              </Button>
            </div>
          </div>
        </div>
      </nav>

      {/* Hero Section */}
      <section className="py-16 md:py-24">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6 }}
            >
              <h1 className="text-4xl md:text-5xl font-bold text-gray-900 mb-4">
                {t.pricing.title}
              </h1>
              <p className="text-xl text-gray-600 max-w-2xl mx-auto">
                {t.pricing.subtitle}
              </p>
            </motion.div>
          </div>

          {/* Pricing Card */}
          <div className="max-w-4xl mx-auto">
            <motion.div
              initial={{ opacity: 0, y: 30 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 0.2 }}
            >
              <Card className="border-2 border-blue-200 shadow-2xl">
                <CardHeader className="bg-gradient-to-br from-blue-600 to-indigo-700 text-white rounded-t-lg p-8 md:p-12">
                  <div className="text-center">
                    <div className="inline-block bg-white/20 backdrop-blur-sm rounded-full px-4 py-2 mb-4">
                      <span className="text-sm font-semibold">{t.pricing.startingFrom}</span>
                    </div>
                    <div className="mb-4">
                      <span className="text-5xl md:text-6xl font-bold">$150</span>
                      <span className="text-2xl md:text-3xl text-blue-100 ml-2">USD</span>
                    </div>
                    <p className="text-blue-100 text-lg mb-2">{t.pricing.oneTimePayment}</p>
                    <div className="flex items-center justify-center gap-2 mt-4">
                      <span className="bg-green-500 text-white px-4 py-2 rounded-full text-sm font-semibold">
                        {t.pricing.firstYearFree}
                      </span>
                      <span className="text-blue-100 text-sm">{t.pricing.firstYear}</span>
                    </div>
                  </div>
                </CardHeader>
                <CardContent className="p-8 md:p-12">
                  {/* Renewal Info */}
                  <div className="bg-gray-50 rounded-lg p-6 mb-8 text-center">
                    <p className="text-gray-600 mb-2">
                      <span className="font-semibold text-gray-900">{t.pricing.renewal}:</span> {t.pricing.fromSecondYear}
                    </p>
                    <p className="text-2xl font-bold text-gray-900">
                      $10 USD <span className="text-lg font-normal text-gray-600">{t.pricing.monthly}</span>
                    </p>
                  </div>

                  {/* Features */}
                  <div className="mb-8">
                    <h3 className="text-xl font-semibold text-gray-900 mb-6">Everything Included:</h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div className="flex items-start space-x-3">
                        <CheckCircle2 className="h-6 w-6 text-green-600 flex-shrink-0 mt-0.5" />
                        <span className="text-gray-700">{t.pricing.features.customDesign}</span>
                      </div>
                      <div className="flex items-start space-x-3">
                        <CheckCircle2 className="h-6 w-6 text-green-600 flex-shrink-0 mt-0.5" />
                        <span className="text-gray-700">{t.pricing.features.customLogo}</span>
                      </div>
                      <div className="flex items-start space-x-3">
                        <CheckCircle2 className="h-6 w-6 text-green-600 flex-shrink-0 mt-0.5" />
                        <span className="text-gray-700">{t.pricing.features.exclusiveQR}</span>
                      </div>
                      <div className="flex items-start space-x-3">
                        <CheckCircle2 className="h-6 w-6 text-green-600 flex-shrink-0 mt-0.5" />
                        <span className="text-gray-700">{t.pricing.features.responsiveDesign}</span>
                      </div>
                      <div className="flex items-start space-x-3">
                        <CheckCircle2 className="h-6 w-6 text-green-600 flex-shrink-0 mt-0.5" />
                        <span className="text-gray-700">{t.pricing.features.shoppingCart}</span>
                      </div>
                      <div className="flex items-start space-x-3">
                        <CheckCircle2 className="h-6 w-6 text-green-600 flex-shrink-0 mt-0.5" />
                        <span className="text-gray-700">{t.pricing.features.whatsappIntegration}</span>
                      </div>
                      <div className="flex items-start space-x-3">
                        <CheckCircle2 className="h-6 w-6 text-green-600 flex-shrink-0 mt-0.5" />
                        <span className="text-gray-700">{t.pricing.features.costCalculation}</span>
                      </div>
                      <div className="flex items-start space-x-3">
                        <CheckCircle2 className="h-6 w-6 text-green-600 flex-shrink-0 mt-0.5" />
                        <span className="text-gray-700">{t.pricing.features.deliveryOptions}</span>
                      </div>
                      <div className="flex items-start space-x-3 md:col-span-2">
                        <CheckCircle2 className="h-6 w-6 text-green-600 flex-shrink-0 mt-0.5" />
                        <span className="text-gray-700">{t.pricing.features.support}</span>
                      </div>
                    </div>
                  </div>

                  {/* CTA Button */}
                  <Button 
                    size="lg" 
                    className="w-full bg-blue-600 hover:bg-blue-700 text-lg py-6"
                    onClick={openWhatsAppQuote}
                  >
                    {t.pricing.cta.quoteButton}
                    <ArrowRight className="ml-2 h-5 w-5" />
                  </Button>
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
            viewport={{ once: true }}
          >
            <h2 className="text-3xl md:text-4xl font-bold text-white mb-6">
              {t.pricing.cta.title}
            </h2>
            <p className="text-xl text-blue-100 mb-8 max-w-2xl mx-auto">
              {t.pricing.cta.subtitle}
            </p>
            <Button 
              size="lg" 
              className="bg-white text-blue-600 hover:bg-gray-100 text-lg px-8 py-3"
              onClick={openWhatsAppQuote}
            >
              {t.pricing.cta.quoteButton}
              <ArrowRight className="ml-2 h-5 w-5" />
            </Button>
          </motion.div>
        </div>
      </section>

      {/* FAQ Section */}
      <section className="py-20 bg-white">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-12 text-center">
            {t.pricing.faq.title}
          </h2>
          
          <div className="space-y-4">
            {t.pricing.faq.questions.map((faq, index) => (
              <motion.div
                key={index}
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.6, delay: index * 0.1 }}
                viewport={{ once: true }}
              >
                <Card className="border border-gray-200 hover:border-blue-300 transition-colors">
                  <CardHeader 
                    className="cursor-pointer"
                    onClick={() => setOpenFaqIndex(openFaqIndex === index ? null : index)}
                  >
                    <div className="flex items-center justify-between">
                      <CardTitle className="text-lg font-semibold text-gray-900 pr-4">
                        {faq.question}
                      </CardTitle>
                      {openFaqIndex === index ? (
                        <ChevronUp className="h-5 w-5 text-gray-500 flex-shrink-0" />
                      ) : (
                        <ChevronDown className="h-5 w-5 text-gray-500 flex-shrink-0" />
                      )}
                    </div>
                  </CardHeader>
                  {openFaqIndex === index && (
                    <CardContent className="pt-0">
                      <p className="text-gray-600 leading-relaxed">
                        {faq.answer}
                      </p>
                    </CardContent>
                  )}
                </Card>
              </motion.div>
            ))}
          </div>
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
                <li><Link href="/#servicio" className="hover:text-white transition-colors">{t.footer.ourService}</Link></li>
                <li><Link href="/#beneficios" className="hover:text-white transition-colors">{t.footer.benefits}</Link></li>
                <li><Link href="/#demo" className="hover:text-white transition-colors">{t.footer.demo}</Link></li>
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

