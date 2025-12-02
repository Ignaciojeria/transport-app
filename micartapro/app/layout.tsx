import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import './globals.css'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'MiCartaPro — Tu Menú Digital, Sin Complicaciones',
  description: 'Gestiona tu menú digital y deja que las ventas fluyan. Diseño 100% personalizado, carrito de compras integrado y recepción de pedidos por WhatsApp.',
  keywords: 'menú digital, carta digital, restaurante, QR code, pedidos online, WhatsApp',
  authors: [{ name: 'MiCartaPro' }],
  openGraph: {
    title: 'MiCartaPro — Tu Menú Digital, Sin Complicaciones',
    description: 'Gestiona tu menú digital y deja que las ventas fluyan.',
    type: 'website',
    locale: 'es_ES',
  },
  twitter: {
    card: 'summary_large_image',
    title: 'MiCartaPro — Tu Menú Digital, Sin Complicaciones',
    description: 'Gestiona tu menú digital y deja que las ventas fluyan.',
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

