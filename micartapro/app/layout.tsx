import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import './globals.css'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'MiCartaPro — Tu Menú Digital, Sin Complicaciones',
  description: 'Self-service SaaS platform for managing your digital menu with our AI agent. Integrated shopping cart, and WhatsApp order reception.',
  keywords: 'menú digital, carta digital, restaurante, QR code, pedidos online, WhatsApp',
  authors: [{ name: 'MiCartaPro' }],
  icons: {
    icon: '/logo.png',
    shortcut: '/logo.png',
    apple: '/logo.png',
  },
  openGraph: {
    title: 'MiCartaPro — Tu Menú Digital, Sin Complicaciones',
    description: 'Gestiona tu menú digital con nuestro agente de IA y deja que las ventas fluyan.',
    type: 'website',
    locale: 'es_ES',
  },
  twitter: {
    card: 'summary_large_image',
    title: 'MiCartaPro — Tu Menú Digital, Sin Complicaciones',
    description: 'Gestiona tu menú digital con nuestro agente de IA y deja que las ventas fluyan.',
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

