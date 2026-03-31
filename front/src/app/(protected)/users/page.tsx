'use client'

import { useCallback, useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import toast from 'react-hot-toast'
import { users } from '@/lib/api'
import { User, ROLE_LABELS } from '@/types'
import Header from '@/components/layout/Header'
import Table from '@/components/ui/Table'
import Badge from '@/components/ui/Badge'
import Button from '@/components/ui/Button'
import Modal from '@/components/ui/Modal'
import Input from '@/components/ui/Input'
import Select from '@/components/ui/Select'
import Pagination from '@/components/ui/Pagination'
import { formatDate } from '@/lib/utils'
import { Plus, UserCheck, UserX } from 'lucide-react'

const schema = z.object({
  username: z.string().min(3, 'Минимум 3 символа'),
  password: z.string().min(6, 'Минимум 6 символов'),
  full_name: z.string().min(1, 'Обязательное поле'),
  email: z.string().email('Некорректный email'),
  role: z.enum(['admin', 'inventory', 'viewer']),
})
type FormData = z.infer<typeof schema>

const ROLE_OPTIONS = Object.entries(ROLE_LABELS).map(([v, l]) => ({ value: v, label: l }))

export default function UsersPage() {
  const [data, setData] = useState<User[]>([])
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [modalOpen, setModalOpen] = useState(false)
  const [saving, setSaving] = useState(false)

  const { register, handleSubmit, reset, formState: { errors } } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: { role: 'viewer' },
  })

  const load = useCallback(() => {
    setLoading(true)
    users.list(page, 20)
      .then((res: { data: User[]; meta: { total_pages: number } }) => {
        setData(res.data || [])
        setTotalPages(res.meta?.total_pages || 1)
      })
      .catch((e: Error) => toast.error(e.message))
      .finally(() => setLoading(false))
  }, [page])

  useEffect(() => { load() }, [load])

  const toggleStatus = async (user: User) => {
    try {
      await users.updateStatus(user.id, !user.is_active)
      toast.success(user.is_active ? 'Пользователь деактивирован' : 'Пользователь активирован')
      load()
    } catch (e: unknown) {
      toast.error(e instanceof Error ? e.message : 'Ошибка')
    }
  }

  const onSubmit = async (data: FormData) => {
    setSaving(true)
    try {
      await users.create(data)
      toast.success('Пользователь создан')
      setModalOpen(false)
      load()
    } catch (e: unknown) {
      toast.error(e instanceof Error ? e.message : 'Ошибка')
    } finally {
      setSaving(false)
    }
  }

  const columns = [
    { key: 'username', header: 'Логин', className: 'font-semibold' },
    { key: 'full_name', header: 'Полное имя' },
    { key: 'email', header: 'Email' },
    {
      key: 'role',
      header: 'Роль',
      render: (row: User) => (
        <Badge variant={row.role === 'admin' ? 'primary' : row.role === 'inventory' ? 'info' : 'default'}>
          {ROLE_LABELS[row.role]}
        </Badge>
      ),
    },
    {
      key: 'is_active',
      header: 'Статус',
      render: (row: User) => (
        <Badge variant={row.is_active ? 'success' : 'danger'}>
          {row.is_active ? 'Активен' : 'Отключён'}
        </Badge>
      ),
    },
    { key: 'created_at', header: 'Создан', render: (row: User) => formatDate(row.created_at) },
    {
      key: 'actions',
      header: '',
      render: (row: User) => (
        <button
          onClick={(e) => { e.stopPropagation(); toggleStatus(row) }}
          className={`p-1.5 rounded-lg transition-all hover:shadow-neu-sm ${row.is_active ? 'text-gray-400 hover:text-red-500' : 'text-gray-400 hover:text-emerald-600'}`}
          title={row.is_active ? 'Деактивировать' : 'Активировать'}
        >
          {row.is_active ? <UserX size={14} /> : <UserCheck size={14} />}
        </button>
      ),
    },
  ]

  return (
    <div>
      <Header
        title="Пользователи"
        actions={
          <Button size="sm" onClick={() => { reset(); setModalOpen(true) }}>
            <Plus size={14} className="mr-1.5" />
            Добавить
          </Button>
        }
      />

      <Table columns={columns} data={data} keyField="id" loading={loading} emptyText="Пользователей нет" />
      <Pagination page={page} totalPages={totalPages} onPageChange={setPage} />

      <Modal open={modalOpen} onClose={() => setModalOpen(false)} title="Новый пользователь">
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <Input label="Логин *" error={errors.username?.message} {...register('username')} />
          <Input label="Пароль *" type="password" error={errors.password?.message} {...register('password')} />
          <Input label="Полное имя *" error={errors.full_name?.message} {...register('full_name')} />
          <Input label="Email *" type="email" error={errors.email?.message} {...register('email')} />
          <Select label="Роль *" options={ROLE_OPTIONS} error={errors.role?.message} {...register('role')} />
          <div className="flex justify-end gap-3 pt-2">
            <Button variant="secondary" type="button" onClick={() => setModalOpen(false)}>Отмена</Button>
            <Button type="submit" loading={saving}>Создать</Button>
          </div>
        </form>
      </Modal>
    </div>
  )
}
