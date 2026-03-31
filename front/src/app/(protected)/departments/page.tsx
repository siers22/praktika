'use client'

import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import toast from 'react-hot-toast'
import { departments } from '@/lib/api'
import { Department } from '@/types'
import Header from '@/components/layout/Header'
import Table from '@/components/ui/Table'
import Button from '@/components/ui/Button'
import Modal from '@/components/ui/Modal'
import Input from '@/components/ui/Input'
import ConfirmDialog from '@/components/ui/ConfirmDialog'
import { formatDate } from '@/lib/utils'
import { Plus, Pencil, Trash2 } from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'

const schema = z.object({
  name: z.string().min(1, 'Введите название'),
  location: z.string().optional(),
})
type FormData = z.infer<typeof schema>

export default function DepartmentsPage() {
  const { hasRole } = useAuth()
  const isAdmin = hasRole('admin')
  const [data, setData] = useState<Department[]>([])
  const [loading, setLoading] = useState(true)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Department | null>(null)
  const [saving, setSaving] = useState(false)
  const [deleteTarget, setDeleteTarget] = useState<Department | null>(null)
  const [deleting, setDeleting] = useState(false)

  const { register, handleSubmit, reset, formState: { errors } } = useForm<FormData>({
    resolver: zodResolver(schema),
  })

  const load = () => {
    setLoading(true)
    departments.list()
      .then((res: { data: Department[] }) => setData(res.data || []))
      .catch((e: Error) => toast.error(e.message))
      .finally(() => setLoading(false))
  }

  useEffect(() => { load() }, [])

  const openCreate = () => {
    setEditing(null)
    reset({ name: '', location: '' })
    setModalOpen(true)
  }

  const openEdit = (dept: Department) => {
    setEditing(dept)
    reset({ name: dept.name, location: dept.location || '' })
    setModalOpen(true)
  }

  const onSubmit = async (data: FormData) => {
    setSaving(true)
    try {
      const payload = { name: data.name, location: data.location || null }
      if (editing) {
        await departments.update(editing.id, payload)
        toast.success('Подразделение обновлено')
      } else {
        await departments.create(payload)
        toast.success('Подразделение добавлено')
      }
      setModalOpen(false)
      load()
    } catch (e: unknown) {
      toast.error(e instanceof Error ? e.message : 'Ошибка')
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async () => {
    if (!deleteTarget) return
    setDeleting(true)
    try {
      await departments.delete(deleteTarget.id)
      toast.success('Подразделение удалено')
      setDeleteTarget(null)
      load()
    } catch (e: unknown) {
      toast.error(e instanceof Error ? e.message : 'Ошибка')
    } finally {
      setDeleting(false)
    }
  }

  const columns = [
    { key: 'name', header: 'Название', className: 'font-semibold' },
    { key: 'location', header: 'Местоположение' },
    { key: 'created_at', header: 'Создано', render: (row: Department) => formatDate(row.created_at) },
    ...(isAdmin ? [{
      key: 'actions',
      header: '',
      render: (row: Department) => (
        <div className="flex gap-1 justify-end">
          <button onClick={(e) => { e.stopPropagation(); openEdit(row) }}
            className="p-1.5 rounded-lg text-gray-400 hover:text-primary hover:shadow-neu-sm transition-all">
            <Pencil size={14} />
          </button>
          <button onClick={(e) => { e.stopPropagation(); setDeleteTarget(row) }}
            className="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:shadow-neu-sm transition-all">
            <Trash2 size={14} />
          </button>
        </div>
      ),
    }] : []),
  ]

  return (
    <div>
      <Header
        title="Подразделения"
        actions={
          isAdmin && (
            <Button size="sm" onClick={openCreate}>
              <Plus size={14} className="mr-1.5" />
              Добавить
            </Button>
          )
        }
      />

      <Table columns={columns} data={data} keyField="id" loading={loading} emptyText="Подразделений нет" />

      <Modal open={modalOpen} onClose={() => setModalOpen(false)} title={editing ? 'Редактировать подразделение' : 'Новое подразделение'}>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <Input label="Название *" error={errors.name?.message} {...register('name')} />
          <Input label="Местоположение" {...register('location')} />
          <div className="flex justify-end gap-3 pt-2">
            <Button variant="secondary" type="button" onClick={() => setModalOpen(false)}>Отмена</Button>
            <Button type="submit" loading={saving}>Сохранить</Button>
          </div>
        </form>
      </Modal>

      <ConfirmDialog
        open={!!deleteTarget}
        onClose={() => setDeleteTarget(null)}
        onConfirm={handleDelete}
        message={`Удалить подразделение "${deleteTarget?.name}"?`}
        loading={deleting}
      />
    </div>
  )
}
