import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import { Suspense } from 'react'
import './globals.css'
import { AuthProvider } from '@/components/AuthProvider'
import { RemotionDemoListener } from '@/components/RemotionDemoListener'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'MiCartaPro - Autenticación',
  description: 'Sistema de autenticación para MiCartaPro',
  icons: {
    icon: '/logo.png',
    apple: '/logo.png',
  },
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="es">
      <head>
        {/* Supabase maneja Google OAuth internamente */}
      </head>
      <body className={inter.className}>
        <RemotionDemoListener />
        <AuthProvider>
          <Suspense fallback={
            <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-indigo-50 flex items-center justify-center">
              <div className="text-center">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
                <p className="text-gray-600">Loading...</p>
              </div>
            </div>
          }>
            {children}
          </Suspense>
        </AuthProvider>
      </body>
    </html>
  )
}

