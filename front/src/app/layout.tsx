import type { Metadata } from 'next'
import { Space_Mono } from 'next/font/google'
import { Toaster } from 'react-hot-toast'
import { AuthProvider } from '@/contexts/AuthContext'
import './globals.css'

const spaceMono = Space_Mono({
  subsets: ['latin'],
  weight: ['400', '700'],
  variable: '--font-space-mono',
})

export const metadata: Metadata = {
  title: 'Инвентаризация оборудования',
  description: 'Система учёта и инвентаризации оборудования',
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="ru" className={spaceMono.variable}>
      <body className="bg-surface text-gray-800 antialiased">
        <AuthProvider>
          {children}
          <Toaster
            position="top-right"
            toastOptions={{
              style: {
                background: '#E7E5E4',
                color: '#374151',
                fontFamily: 'var(--font-space-mono)',
                fontSize: '13px',
                boxShadow: '6px 6px 12px #c8c6c4, -6px -6px 12px #ffffff',
                borderRadius: '12px',
              },
            }}
          />
        </AuthProvider>
      </body>
    </html>
  )
}
