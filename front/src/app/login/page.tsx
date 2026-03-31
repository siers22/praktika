'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import toast from 'react-hot-toast'
import { useAuth } from '@/contexts/AuthContext'
import { auth } from '@/lib/api'
import Input from '@/components/ui/Input'
import Button from '@/components/ui/Button'

const schema = z.object({
  username: z.string().min(1, 'Введите логин'),
  password: z.string().min(1, 'Введите пароль'),
})

type FormData = z.infer<typeof schema>

export default function LoginPage() {
  const router = useRouter()
  const { login } = useAuth()
  const [loading, setLoading] = useState(false)

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<FormData>({ resolver: zodResolver(schema) })

  const onSubmit = async (data: FormData) => {
    setLoading(true)
    try {
      const res = await auth.login(data.username, data.password)
      login(res.data.access_token, res.data.refresh_token)
      router.replace('/dashboard')
    } catch (e: unknown) {
      toast.error(e instanceof Error ? e.message : 'Ошибка входа')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-surface flex items-center justify-center p-4">
      <div className="w-full max-w-sm">
        {/* Logo block */}
        <div className="mb-8 text-center">
          <div className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-surface shadow-neu mb-4">
            <span className="text-2xl font-mono font-bold text-primary">И</span>
          </div>
          <h1 className="font-mono font-bold text-gray-800 text-xl">Инвентаризация</h1>
          <p className="font-mono text-sm text-gray-400 mt-1">оборудования</p>
        </div>

        {/* Form card */}
        <div className="rounded-2xl bg-surface shadow-neu p-8">
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <Input
              id="username"
              label="Логин"
              placeholder="admin"
              autoComplete="username"
              error={errors.username?.message}
              {...register('username')}
            />
            <Input
              id="password"
              type="password"
              label="Пароль"
              placeholder="••••••••"
              autoComplete="current-password"
              error={errors.password?.message}
              {...register('password')}
            />
            <Button type="submit" className="w-full mt-2" loading={loading}>
              Войти
            </Button>
          </form>
        </div>
      </div>
    </div>
  )
}
