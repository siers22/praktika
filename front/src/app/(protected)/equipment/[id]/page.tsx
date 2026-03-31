'use client'

import { useEffect, useState } from 'react'
import { useParams, useRouter } from 'next/navigation'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import toast from 'react-hot-toast'
import { equipment, categories, departments, users } from '@/lib/api'
import { Equipment, STATUS_LABELS, STATUS_COLORS } from '@/types'
import type { BadgeVariant } from '@/components/ui/Badge'
import Header from '@/components/layout/Header'
import Card from '@/components/ui/Card'
import Badge from '@/components/ui/Badge'
import Input from '@/components/ui/Input'
import Select from '@/components/ui/Select'
import Textarea from '@/components/ui/Textarea'
import Button from '@/components/ui/Button'
import ConfirmDialog from '@/components/ui/ConfirmDialog'
import { formatDate, formatCurrency } from '@/lib/utils'
import { ArrowLeft, Archive, Upload, Trash2 } from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'

const schema = z.object({
  inventory_number: z.string().min(1),
  name: z.string().min(1),
  category_id: z.string().min(1),
  department_id: z.string().min(1),
  status: z.string().min(1),
  serial_number: z.string().optional(),
  model: z.string().optional(),
  manufacturer: z.string().optional(),
  purchase_date: z.string().optional(),
  purchase_price: z.string().optional(),
  warranty_expiry: z.string().optional(),
  responsible_person_id: z.string().optional(),
  description: z.string().optional(),
  notes: z.string().optional(),
})

type FormData = z.infer<typeof schema>

const STATUS_OPTIONS = [
  { value: 'in_use', label: 'В эксплуатации' },
  { value: 'in_storage', label: 'На складе' },
  { value: 'in_repair', label: 'В ремонте' },
  { value: 'written_off', label: 'Списано' },
  { value: 'reserved', label: 'Зарезервировано' },
]

const BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'

export default function EquipmentDetailPage() {
  const { id } = useParams<{ id: string }>()
  const router = useRouter()
  const { hasRole } = useAuth()
  const isAdmin = hasRole('admin')

  const [eq, setEq] = useState<Equipment | null>(null)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [archiving, setArchiving] = useState(false)
  const [confirmArchive, setConfirmArchive] = useState(false)
  const [categoryOptions, setCategoryOptions] = useState<{ value: string; label: string }[]>([])
  const [deptOptions, setDeptOptions] = useState<{ value: string; label: string }[]>([])
  const [userOptions, setUserOptions] = useState<{ value: string; label: string }[]>([])

  const { register, handleSubmit, reset, formState: { errors } } = useForm<FormData>({
    resolver: zodResolver(schema),
  })

  useEffect(() => {
    Promise.all([
      categories.list(),
      departments.list(),
      users.list(1, 100),
      equipment.getById(id),
    ]).then(([cats, depts, usrs, eqRes]) => {
      setCategoryOptions((cats as { data: { id: string; name: string }[] }).data.map((c) => ({ value: c.id, label: c.name })))
      setDeptOptions((depts as { data: { id: string; name: string }[] }).data.map((d) => ({ value: d.id, label: d.name })))
      setUserOptions([
        { value: '', label: 'Не назначен' },
        ...(usrs as { data: { id: string; full_name: string }[] }).data.map((u) => ({ value: u.id, label: u.full_name })),
      ])
      const item = (eqRes as { data: Equipment }).data
      setEq(item)
      reset({
        inventory_number: item.inventory_number,
        name: item.name,
        category_id: item.category_id,
        department_id: item.department_id,
        status: item.status,
        serial_number: item.serial_number || '',
        model: item.model || '',
        manufacturer: item.manufacturer || '',
        purchase_date: item.purchase_date?.slice(0, 10) || '',
        purchase_price: item.purchase_price?.toString() || '',
        warranty_expiry: item.warranty_expiry?.slice(0, 10) || '',
        responsible_person_id: item.responsible_person_id || '',
        description: item.description || '',
        notes: item.notes || '',
      })
    }).catch((e: Error) => toast.error(e.message))
      .finally(() => setLoading(false))
  }, [id, reset])

  const onSubmit = async (data: FormData) => {
    setSaving(true)
    try {
      const payload: Record<string, unknown> = {
        ...data,
        purchase_price: data.purchase_price ? Number(data.purchase_price) : null,
        purchase_date: data.purchase_date || null,
        warranty_expiry: data.warranty_expiry || null,
        serial_number: data.serial_number || null,
        model: data.model || null,
        manufacturer: data.manufacturer || null,
        responsible_person_id: data.responsible_person_id || null,
        description: data.description || null,
        notes: data.notes || null,
      }
      await equipment.update(id, payload)
      toast.success('Сохранено')
    } catch (e: unknown) {
      toast.error(e instanceof Error ? e.message : 'Ошибка')
    } finally {
      setSaving(false)
    }
  }

  const handleArchive = async () => {
    setArchiving(true)
    try {
      await equipment.archive(id)
      toast.success('Оборудование архивировано')
      router.push('/equipment')
    } catch (e: unknown) {
      toast.error(e instanceof Error ? e.message : 'Ошибка')
    } finally {
      setArchiving(false)
      setConfirmArchive(false)
    }
  }

  const handlePhotoUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return
    try {
      const res = await equipment.uploadPhoto(id, file)
      setEq((prev) => prev ? { ...prev, photos: [...(prev.photos || []), (res as { data: Equipment['photos'][0] }).data] } : prev)
      toast.success('Фото загружено')
    } catch {
      toast.error('Ошибка загрузки фото')
    }
  }

  const handlePhotoDelete = async (photoId: string) => {
    try {
      await equipment.deletePhoto(id, photoId)
      setEq((prev) => prev ? { ...prev, photos: prev.photos.filter((p) => p.id !== photoId) } : prev)
    } catch {
      toast.error('Ошибка удаления фото')
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <span className="h-8 w-8 animate-spin rounded-full border-2 border-primary border-t-transparent" />
      </div>
    )
  }

  if (!eq) return null

  return (
    <div>
      <Header
        title={eq.name}
        actions={
          <div className="flex gap-2">
            <Button variant="ghost" size="sm" onClick={() => router.back()}>
              <ArrowLeft size={14} className="mr-1" />
              Назад
            </Button>
            {isAdmin && !eq.is_archived && (
              <Button variant="danger" size="sm" onClick={() => setConfirmArchive(true)}>
                <Archive size={14} className="mr-1" />
                Архив
              </Button>
            )}
          </div>
        }
      />

      {eq.is_archived && (
        <div className="mb-4 rounded-xl bg-amber-50 border border-amber-200 px-4 py-2.5 font-mono text-xs text-amber-700">
          Оборудование архивировано
        </div>
      )}

      <form onSubmit={handleSubmit(onSubmit)}>
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
          <Card title="Основная информация">
            <div className="space-y-4">
              <Input label="Инвентарный номер" error={errors.inventory_number?.message} {...register('inventory_number')} />
              <Input label="Наименование" error={errors.name?.message} {...register('name')} />
              <Select label="Категория" options={categoryOptions} {...register('category_id')} />
              <Select label="Подразделение" options={deptOptions} {...register('department_id')} />
              <Select label="Статус" options={STATUS_OPTIONS} {...register('status')} />
              <Select label="Ответственный" options={userOptions} {...register('responsible_person_id')} />
            </div>
          </Card>

          <Card title="Технические данные">
            <div className="space-y-4">
              <Input label="Серийный номер" {...register('serial_number')} />
              <Input label="Модель" {...register('model')} />
              <Input label="Производитель" {...register('manufacturer')} />
              <Input label="Дата покупки" type="date" {...register('purchase_date')} />
              <Input label="Стоимость (₽)" type="number" {...register('purchase_price')} />
              <Input label="Гарантия до" type="date" {...register('warranty_expiry')} />
            </div>
          </Card>
        </div>

        <Card title="Описание и примечания" className="mb-6">
          <div className="space-y-4">
            <Textarea label="Описание" rows={3} {...register('description')} />
            <Textarea label="Примечания" rows={2} {...register('notes')} />
          </div>
        </Card>

        {!eq.is_archived && (
          <div className="flex justify-end mb-6">
            <Button type="submit" loading={saving}>
              Сохранить изменения
            </Button>
          </div>
        )}
      </form>

      {/* Photos */}
      <Card title="Фотографии">
        <div className="flex flex-wrap gap-3 mb-4">
          {eq.photos?.map((photo) => (
            <div key={photo.id} className="relative group">
              <img
                src={`${BASE_URL.replace('/api/v1', '')}/uploads/${photo.file_path}`}
                alt="фото"
                className="w-28 h-28 object-cover rounded-xl shadow-neu-sm"
              />
              {isAdmin && (
                <button
                  onClick={() => handlePhotoDelete(photo.id)}
                  className="absolute top-1 right-1 hidden group-hover:flex items-center justify-center w-6 h-6 rounded-lg bg-red-500 text-white"
                >
                  <Trash2 size={12} />
                </button>
              )}
            </div>
          ))}
        </div>
        {isAdmin && !eq.is_archived && (
          <label className="inline-flex items-center gap-2 cursor-pointer rounded-xl bg-surface shadow-neu px-4 py-2 font-mono text-sm text-gray-700 hover:shadow-neu-sm active:shadow-neu-inset transition-all">
            <Upload size={14} className="mr-1" />
            Загрузить фото
            <input type="file" accept="image/*" className="hidden" onChange={handlePhotoUpload} />
          </label>
        )}
      </Card>

      <ConfirmDialog
        open={confirmArchive}
        onClose={() => setConfirmArchive(false)}
        onConfirm={handleArchive}
        title="Архивировать оборудование?"
        message={`Оборудование "${eq.name}" будет перемещено в архив.`}
        confirmLabel="Архивировать"
        loading={archiving}
      />
    </div>
  )
}
