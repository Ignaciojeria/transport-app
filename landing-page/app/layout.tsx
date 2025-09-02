import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import './globals.css'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'TransportApp - Optimiza tu flota en minutos',
  description: 'Asigna rutas automáticamente, comparte enlaces con tus conductores y recibe reportes de entrega. Todo desde Google Sheets.',
  keywords: 'logística, optimización de rutas, flota, transporte, Google Sheets, automatización',
  authors: [{ name: 'TransportApp' }],
  openGraph: {
    title: 'TransportApp - Optimiza tu flota en minutos',
    description: 'Asigna rutas automáticamente, comparte enlaces con tus conductores y recibe reportes de entrega. Todo desde Google Sheets.',
    type: 'website',
    locale: 'es_ES',
  },
  twitter: {
    card: 'summary_large_image',
    title: 'TransportApp - Optimiza tu flota en minutos',
    description: 'Asigna rutas automáticamente, comparte enlaces con tus conductores y recibe reportes de entrega. Todo desde Google Sheets.',
  },
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="es">
      <body className={inter.className}>
        {children}
      </body>
    </html>
  )
}
