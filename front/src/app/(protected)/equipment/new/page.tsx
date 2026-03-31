'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import toast from 'react-hot-toast'
import { equipment, categories, departments, users } from '@/lib/api'
import Header from '@/components/layout/Header'
import Card from '@/components/ui/Card'
import Input from '@/components/ui/Input'
import Select from '@/components/ui/Select'
import Textarea from '@/components/ui/Textarea'
import Button from '@/components/ui/Button'

const schema = z.object({
  inventory_number: z.string().min(1, 'Обязательное поле'),
  name: z.string().min(1, 'Обязательное поле'),
  category_id: z.string().min(1, 'Выберите категорию'),
  department_id: z.string().min(1, 'Выберите подразделение'),
  status: z.string().min(1, 'Выберите статус'),
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
  { value: 'reserved', label: 'Зарезервировано' },
]

export default function NewEquipmentPage() {
  const router = useRouter()
  const [loading, setLoading] = useState(false)
  const [categoryOptions, setCategoryOptions] = useState<{ value: string; label: string }[]>([])
  const [deptOptions, setDeptOptions] = useState<{ value: string; label: string }[]>([])
  const [userOptions, setUserOptions] = useState<{ value: string; label: string }[]>([])

  useEffect(() => {
    categories.list().then((res: { data: { id: string; name: string }[] }) =>
      setCategoryOptions(res.data.map((c) => ({ value: c.id, label: c.name })))
    )
    departments.list().then((res: { data: { id: string; name: string }[] }) =>
      setDeptOptions(res.data.map((d) => ({ value: d.id, label: d.name })))
    )
    users.list(1, 100).then((res: { data: { id: string; full_name: string }[] }) =>
      setUserOptions([{ value: '', label: 'Не назначен' }, ...res.data.map((u) => ({ value: u.id, label: u.full_name }))])
    )
  }, [])

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<FormData>({ resolver: zodResolver(schema) })

  const onSubmit = async (data: FormData) => {
    setLoading(true)
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
      await equipment.create(payload)
      toast.success('Оборудование добавлено')
      router.push('/equipment')
    } catch (e: unknown) {
      toast.error(e instanceof Error ? e.message : 'Ошибка')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div>
      <Header title="Новое оборудование" />
      <form onSubmit={handleSubmit(onSubmit)}>
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
          <Card title="Основная информация">
            <div className="space-y-4">
              <Input label="Инвентарный номер *" error={errors.inventory_number?.message} {...register('inventory_number')} />
              <Input label="Наименование *" error={errors.name?.message} {...register('name')} />
              <Select
                label="Категория *"
                options={categoryOptions}
                placeholder="Выберите..."
                error={errors.category_id?.message}
                {...register('category_id')}
              />
              <Select
                label="Подразделение *"
                options={deptOptions}
                placeholder="Выберите..."
                error={errors.department_id?.message}
                {...register('department_id')}
              />
              <Select
                label="Статус *"
                options={STATUS_OPTIONS}
                placeholder="Выберите..."
                error={errors.status?.message}
                {...register('status')}
              />
              <Select
                label="Ответственный"
                options={userOptions}
                {...register('responsible_person_id')}
              />
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

        <div className="flex gap-3 justify-end">
          <Button variant="secondary" type="button" onClick={() => router.back()}>
            Отмена
          </Button>
          <Button type="submit" loading={loading}>
            Сохранить
          </Button>
        </div>
      </form>
    </div>
  )
}
