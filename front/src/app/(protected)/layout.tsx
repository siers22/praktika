import AuthGuard from '@/components/layout/AuthGuard'
import MainLayout from '@/components/layout/MainLayout'

export default function ProtectedLayout({ children }: { children: React.ReactNode }) {
  return (
    <AuthGuard>
      <MainLayout>{children}</MainLayout>
    </AuthGuard>
  )
}
