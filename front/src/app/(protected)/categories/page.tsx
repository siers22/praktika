'use client'

import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import toast from 'react-hot-toast'
import { categories } from '@/lib/api'
import { Category } from '@/types'
import Header from '@/components/layout/Header'
import Table from '@/components/ui/Table'
import Button from '@/components/ui/Button'
import Modal from '@/components/ui/Modal'
import Input from '@/components/ui/Input'
import Textarea from '@/components/ui/Textarea'
import ConfirmDialog from '@/components/ui/ConfirmDialog'
import { formatDate } from '@/lib/utils'
import { Plus, Pencil, Trash2 } from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'

const schema = z.object({
  name: z.string().min(1, 'Введите название'),
  description: z.string().optional(),
})
type FormData = z.infer<typeof schema>

export default function CategoriesPage() {
  const { hasRole } = useAuth()
  const isAdmin = hasRole('admin')
  const [data, setData] = useState<Category[]>([])
  const [loading, setLoading] = useState(true)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Category | null>(null)
  const [saving, setSaving] = useState(false)
  const [deleteTarget, setDeleteTarget] = useState<Category | null>(null)
  const [deleting, setDeleting] = useState(false)

  const { register, handleSubmit, reset, formState: { errors } } = useForm<FormData>({
    resolver: zodResolver(schema),
  })

  const load = () => {
    setLoading(true)
    categories.list()
      .then((res: { data: Category[] }) => setData(res.data || []))
      .catch((e: Error) => toast.error(e.message))
      .finally(() => setLoading(false))
  }

  useEffect(() => { load() }, [])

  const openCreate = () => {
    setEditing(null)
    reset({ name: '', description: '' })
    setModalOpen(true)
  }

  const openEdit = (cat: Category) => {
    setEditing(cat)
    reset({ name: cat.name, description: cat.description || '' })
    setModalOpen(true)
  }

  const onSubmit = async (data: FormData) => {
    setSaving(true)
    try {
      const payload = { name: data.name, description: data.description || null }
      if (editing) {
        await categories.update(editing.id, payload)
        toast.success('Категория обновлена')
      } else {
        await categories.create(payload)
        toast.success('Категория добавлена')
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
      await categories.delete(deleteTarget.id)
      toast.success('Категория удалена')
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
    { key: 'description', header: 'Описание' },
    { key: 'created_at', header: 'Создана', render: (row: Category) => formatDate(row.created_at) },
    ...(isAdmin ? [{
      key: 'actions',
      header: '',
      render: (row: Category) => (
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
        title="Категории"
        actions={
          isAdmin && (
            <Button size="sm" onClick={openCreate}>
              <Plus size={14} className="mr-1.5" />
              Добавить
            </Button>
          )
        }
      />

      <Table columns={columns} data={data} keyField="id" loading={loading} emptyText="Категорий нет" />

      <Modal open={modalOpen} onClose={() => setModalOpen(false)} title={editing ? 'Редактировать категорию' : 'Новая категория'}>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <Input label="Название *" error={errors.name?.message} {...register('name')} />
          <Textarea label="Описание" {...register('description')} />
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
        message={`Удалить категорию "${deleteTarget?.name}"?`}
        loading={deleting}
      />
    </div>
  )
}
