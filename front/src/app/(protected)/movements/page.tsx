'use client'

import { useCallback, useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import toast from 'react-hot-toast'
import { movements, equipment as equipmentApi, departments } from '@/lib/api'
import { Movement } from '@/types'
import Header from '@/components/layout/Header'
import Table from '@/components/ui/Table'
import Pagination from '@/components/ui/Pagination'
import Button from '@/components/ui/Button'
import Modal from '@/components/ui/Modal'
import Select from '@/components/ui/Select'
import Textarea from '@/components/ui/Textarea'
import Input from '@/components/ui/Input'
import { formatDateTime } from '@/lib/utils'
import { Plus } from 'lucide-react'

const schema = z.object({
  equipment_id: z.string().min(1, 'Выберите оборудование'),
  to_department_id: z.string().min(1, 'Выберите подразделение'),
  reason: z.string().optional(),
})
type FormData = z.infer<typeof schema>

export default function MovementsPage() {
  const [data, setData] = useState<Movement[]>([])
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [modalOpen, setModalOpen] = useState(false)
  const [saving, setSaving] = useState(false)
  const [equipOptions, setEquipOptions] = useState<{ value: string; label: string }[]>([])
  const [deptOptions, setDeptOptions] = useState<{ value: string; label: string }[]>([])

  const { register, handleSubmit, reset, formState: { errors } } = useForm<FormData>({
    resolver: zodResolver(schema),
  })

  const load = useCallback(() => {
    setLoading(true)
    movements.list({ page, per_page: 20 })
      .then((res: { data: Movement[]; meta: { total_pages: number } }) => {
        setData(res.data || [])
        setTotalPages(res.meta?.total_pages || 1)
      })
      .catch((e: Error) => toast.error(e.message))
      .finally(() => setLoading(false))
  }, [page])

  useEffect(() => { load() }, [load])

  const openModal = async () => {
    const [eqRes, deptRes] = await Promise.all([
      equipmentApi.list({ per_page: 200 }),
      departments.list(),
    ])
    setEquipOptions((eqRes as { data: { id: string; name: string; inventory_number: string }[] }).data.map((e) => ({
      value: e.id,
      label: `${e.inventory_number} — ${e.name}`,
    })))
    setDeptOptions((deptRes as { data: { id: string; name: string }[] }).data.map((d) => ({
      value: d.id,
      label: d.name,
    })))
    reset()
    setModalOpen(true)
  }

  const onSubmit = async (data: FormData) => {
    setSaving(true)
    try {
      await movements.create({ ...data, reason: data.reason || null })
      toast.success('Перемещение зафиксировано')
      setModalOpen(false)
      load()
    } catch (e: unknown) {
      toast.error(e instanceof Error ? e.message : 'Ошибка')
    } finally {
      setSaving(false)
    }
  }

  const columns = [
    { key: 'inventory_number', header: 'Инв. номер' },
    { key: 'equipment_name', header: 'Оборудование' },
    { key: 'from_department_name', header: 'Откуда' },
    { key: 'to_department_name', header: 'Куда' },
    { key: 'moved_by_name', header: 'Кто переместил' },
    {
      key: 'moved_at',
      header: 'Дата',
      render: (row: Movement) => formatDateTime(row.moved_at),
    },
    { key: 'reason', header: 'Причина' },
  ]

  return (
    <div>
      <Header
        title="Перемещения"
        actions={
          <Button size="sm" onClick={openModal}>
            <Plus size={14} className="mr-1.5" />
            Зафиксировать
          </Button>
        }
      />

      <Table
        columns={columns}
        data={data}
        keyField="id"
        loading={loading}
        emptyText="Перемещений нет"
      />
      <Pagination page={page} totalPages={totalPages} onPageChange={setPage} />

      <Modal open={modalOpen} onClose={() => setModalOpen(false)} title="Новое перемещение">
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <Select
            label="Оборудование *"
            options={equipOptions}
            placeholder="Выберите..."
            error={errors.equipment_id?.message}
            {...register('equipment_id')}
          />
          <Select
            label="Подразделение назначения *"
            options={deptOptions}
            placeholder="Выберите..."
            error={errors.to_department_id?.message}
            {...register('to_department_id')}
          />
          <Textarea label="Причина" {...register('reason')} />
          <div className="flex justify-end gap-3 pt-2">
            <Button variant="secondary" type="button" onClick={() => setModalOpen(false)}>Отмена</Button>
            <Button type="submit" loading={saving}>Сохранить</Button>
          </div>
        </form>
      </Modal>
    </div>
  )
}
