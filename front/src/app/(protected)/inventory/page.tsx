'use client'

import { useCallback, useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import toast from 'react-hot-toast'
import { inventory, departments } from '@/lib/api'
import { InventorySession } from '@/types'
import Header from '@/components/layout/Header'
import Table from '@/components/ui/Table'
import Badge from '@/components/ui/Badge'
import Button from '@/components/ui/Button'
import Modal from '@/components/ui/Modal'
import Select from '@/components/ui/Select'
import Pagination from '@/components/ui/Pagination'
import { formatDateTime } from '@/lib/utils'
import { Plus } from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'

const schema = z.object({
  department_id: z.string().min(1, 'Выберите подразделение'),
})
type FormData = z.infer<typeof schema>

export default function InventoryPage() {
  const router = useRouter()
  const { hasRole } = useAuth()
  const canCreate = hasRole('admin', 'inventory')
  const [data, setData] = useState<InventorySession[]>([])
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [modalOpen, setModalOpen] = useState(false)
  const [saving, setSaving] = useState(false)
  const [deptOptions, setDeptOptions] = useState<{ value: string; label: string }[]>([])

  const { register, handleSubmit, reset, formState: { errors } } = useForm<FormData>({
    resolver: zodResolver(schema),
  })

  const load = useCallback(() => {
    setLoading(true)
    inventory.listSessions(page)
      .then((res: { data: InventorySession[]; meta: { total_pages: number } }) => {
        setData(res.data || [])
        setTotalPages(res.meta?.total_pages || 1)
      })
      .catch((e: Error) => toast.error(e.message))
      .finally(() => setLoading(false))
  }, [page])

  useEffect(() => { load() }, [load])

  const openModal = async () => {
    const res = await departments.list()
    setDeptOptions((res as { data: { id: string; name: string }[] }).data.map((d) => ({ value: d.id, label: d.name })))
    reset()
    setModalOpen(true)
  }

  const onSubmit = async (data: FormData) => {
    setSaving(true)
    try {
      const res = await inventory.createSession(data.department_id)
      toast.success('Сессия создана')
      setModalOpen(false)
      router.push(`/inventory/${(res as { data: { id: string } }).data.id}`)
    } catch (e: unknown) {
      toast.error(e instanceof Error ? e.message : 'Ошибка')
    } finally {
      setSaving(false)
    }
  }

  const columns = [
    { key: 'department_name', header: 'Подразделение', className: 'font-semibold' },
    {
      key: 'status',
      header: 'Статус',
      render: (row: InventorySession) => (
        <Badge variant={row.status === 'completed' ? 'success' : 'warning'}>
          {row.status === 'completed' ? 'Завершена' : 'В процессе'}
        </Badge>
      ),
    },
    { key: 'created_by_name', header: 'Создал' },
    { key: 'started_at', header: 'Начата', render: (row: InventorySession) => formatDateTime(row.started_at) },
    {
      key: 'finished_at',
      header: 'Завершена',
      render: (row: InventorySession) => formatDateTime(row.finished_at),
    },
  ]

  return (
    <div>
      <Header
        title="Инвентаризация"
        actions={
          canCreate && (
            <Button size="sm" onClick={openModal}>
              <Plus size={14} className="mr-1.5" />
              Новая сессия
            </Button>
          )
        }
      />

      <Table
        columns={columns}
        data={data}
        keyField="id"
        loading={loading}
        onRowClick={(row) => router.push(`/inventory/${row.id}`)}
        emptyText="Сессий инвентаризации нет"
      />
      <Pagination page={page} totalPages={totalPages} onPageChange={setPage} />

      <Modal open={modalOpen} onClose={() => setModalOpen(false)} title="Новая инвентаризация">
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <Select
            label="Подразделение *"
            options={deptOptions}
            placeholder="Выберите..."
            error={errors.department_id?.message}
            {...register('department_id')}
          />
          <div className="flex justify-end gap-3 pt-2">
            <Button variant="secondary" type="button" onClick={() => setModalOpen(false)}>Отмена</Button>
            <Button type="submit" loading={saving}>Создать</Button>
          </div>
        </form>
      </Modal>
    </div>
  )
}
