'use client'

import { useEffect, useState } from 'react'
import { useParams, useRouter } from 'next/navigation'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import toast from 'react-hot-toast'
import { inventory } from '@/lib/api'
import { InventorySession, InventoryItem, ACTUAL_STATUS_LABELS } from '@/types'
import Header from '@/components/layout/Header'
import Card from '@/components/ui/Card'
import Badge from '@/components/ui/Badge'
import Table from '@/components/ui/Table'
import Button from '@/components/ui/Button'
import Modal from '@/components/ui/Modal'
import Input from '@/components/ui/Input'
import Select from '@/components/ui/Select'
import Textarea from '@/components/ui/Textarea'
import ConfirmDialog from '@/components/ui/ConfirmDialog'
import { formatDateTime, downloadBlob } from '@/lib/utils'
import { ArrowLeft, CheckCircle, Download } from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'

const schema = z.object({
  equipment_id: z.string().min(1, 'Обязательное поле'),
  actual_status: z.enum(['found', 'not_found', 'damaged']),
  comment: z.string().optional(),
})
type FormData = z.infer<typeof schema>

const ACTUAL_STATUS_OPTIONS = Object.entries(ACTUAL_STATUS_LABELS).map(([v, l]) => ({
  value: v,
  label: l,
}))

export default function InventorySessionPage() {
  const { id } = useParams<{ id: string }>()
  const router = useRouter()
  const { hasRole } = useAuth()
  const canEdit = hasRole('admin', 'inventory')

  const [session, setSession] = useState<InventorySession | null>(null)
  const [items, setItems] = useState<InventoryItem[]>([])
  const [loading, setLoading] = useState(true)
  const [modalOpen, setModalOpen] = useState(false)
  const [editingItem, setEditingItem] = useState<InventoryItem | null>(null)
  const [saving, setSaving] = useState(false)
  const [completing, setCompleting] = useState(false)
  const [confirmComplete, setConfirmComplete] = useState(false)
  const [exportLoading, setExportLoading] = useState(false)

  const { register, handleSubmit, reset, formState: { errors } } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: { actual_status: 'found' },
  })

  const load = () => {
    setLoading(true)
    inventory.getSession(id)
      .then((res: { data: { session: InventorySession; items: InventoryItem[] } }) => {
        setSession(res.data.session)
        setItems(res.data.items || [])
      })
      .catch((e: Error) => toast.error(e.message))
      .finally(() => setLoading(false))
  }

  useEffect(() => { load() }, [id])

  const openCheckModal = (item?: InventoryItem) => {
    if (item) {
      setEditingItem(item)
      reset({
        equipment_id: item.equipment_id,
        actual_status: item.actual_status,
        comment: item.comment || '',
      })
    } else {
      setEditingItem(null)
      reset({ equipment_id: '', actual_status: 'found', comment: '' })
    }
    setModalOpen(true)
  }

  const onSubmit = async (data: FormData) => {
    setSaving(true)
    try {
      if (editingItem) {
        await inventory.updateItem(id, editingItem.id, data)
        toast.success('Обновлено')
      } else {
        await inventory.checkItem(id, data)
        toast.success('Позиция отмечена')
      }
      setModalOpen(false)
      load()
    } catch (e: unknown) {
      toast.error(e instanceof Error ? e.message : 'Ошибка')
    } finally {
      setSaving(false)
    }
  }

  const handleComplete = async () => {
    setCompleting(true)
    try {
      await inventory.complete(id)
      toast.success('Инвентаризация завершена')
      setConfirmComplete(false)
      load()
    } catch (e: unknown) {
      toast.error(e instanceof Error ? e.message : 'Ошибка')
    } finally {
      setCompleting(false)
    }
  }

  const handleExport = async () => {
    setExportLoading(true)
    try {
      const res = await inventory.exportCSV(id)
      const blob = await res.blob()
      downloadBlob(blob, `inventory_${id.slice(0, 8)}.csv`)
    } catch {
      toast.error('Ошибка экспорта')
    } finally {
      setExportLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <span className="h-8 w-8 animate-spin rounded-full border-2 border-primary border-t-transparent" />
      </div>
    )
  }

  if (!session) return null

  const isCompleted = session.status === 'completed'

  const statusVariantMap: Record<string, 'success' | 'danger' | 'warning'> = {
    found: 'success',
    not_found: 'danger',
    damaged: 'warning',
  }

  const columns = [
    { key: 'inventory_number', header: 'Инв. номер' },
    { key: 'equipment_name', header: 'Оборудование' },
    { key: 'expected_status', header: 'Ожидаемый статус' },
    {
      key: 'actual_status',
      header: 'Фактический статус',
      render: (row: InventoryItem) => (
        <Badge variant={statusVariantMap[row.actual_status] || 'default'}>
          {ACTUAL_STATUS_LABELS[row.actual_status]}
        </Badge>
      ),
    },
    { key: 'comment', header: 'Комментарий' },
    { key: 'checked_at', header: 'Проверено', render: (row: InventoryItem) => formatDateTime(row.checked_at) },
  ]

  return (
    <div>
      <Header
        title={`Инвентаризация: ${session.department_name}`}
        actions={
          <div className="flex gap-2">
            <Button variant="ghost" size="sm" onClick={() => router.back()}>
              <ArrowLeft size={14} className="mr-1" />
              Назад
            </Button>
            <Button variant="secondary" size="sm" onClick={handleExport} loading={exportLoading}>
              <Download size={14} className="mr-1" />
              CSV
            </Button>
            {canEdit && !isCompleted && (
              <>
                <Button size="sm" onClick={() => openCheckModal()}>
                  Добавить позицию
                </Button>
                <Button variant="primary" size="sm" onClick={() => setConfirmComplete(true)}>
                  <CheckCircle size={14} className="mr-1" />
                  Завершить
                </Button>
              </>
            )}
          </div>
        }
      />

      <Card className="mb-6">
        <div className="flex flex-wrap gap-6 text-xs font-mono">
          <div>
            <span className="text-gray-400 uppercase tracking-wider">Статус: </span>
            <Badge variant={isCompleted ? 'success' : 'warning'}>
              {isCompleted ? 'Завершена' : 'В процессе'}
            </Badge>
          </div>
          <div>
            <span className="text-gray-400 uppercase tracking-wider">Создал: </span>
            <span className="text-gray-700">{session.created_by_name}</span>
          </div>
          <div>
            <span className="text-gray-400 uppercase tracking-wider">Начата: </span>
            <span className="text-gray-700">{formatDateTime(session.started_at)}</span>
          </div>
          {session.finished_at && (
            <div>
              <span className="text-gray-400 uppercase tracking-wider">Завершена: </span>
              <span className="text-gray-700">{formatDateTime(session.finished_at)}</span>
            </div>
          )}
          <div>
            <span className="text-gray-400 uppercase tracking-wider">Позиций: </span>
            <span className="text-gray-700">{items.length}</span>
          </div>
        </div>
      </Card>

      <Table
        columns={columns}
        data={items}
        keyField="id"
        onRowClick={canEdit && !isCompleted ? (row: InventoryItem) => openCheckModal(row) : undefined}
        emptyText="Позиции не добавлены"
      />

      <Modal open={modalOpen} onClose={() => setModalOpen(false)} title={editingItem ? 'Редактировать позицию' : 'Добавить позицию'}>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          {!editingItem && (
            <Input
              label="ID оборудования *"
              placeholder="UUID оборудования"
              error={errors.equipment_id?.message}
              {...register('equipment_id')}
            />
          )}
          <Select
            label="Фактический статус *"
            options={ACTUAL_STATUS_OPTIONS}
            error={errors.actual_status?.message}
            {...register('actual_status')}
          />
          <Textarea label="Комментарий" {...register('comment')} />
          <div className="flex justify-end gap-3 pt-2">
            <Button variant="secondary" type="button" onClick={() => setModalOpen(false)}>Отмена</Button>
            <Button type="submit" loading={saving}>Сохранить</Button>
          </div>
        </form>
      </Modal>

      <ConfirmDialog
        open={confirmComplete}
        onClose={() => setConfirmComplete(false)}
        onConfirm={handleComplete}
        title="Завершить инвентаризацию?"
        message="После завершения редактирование позиций будет недоступно."
        confirmLabel="Завершить"
        loading={completing}
      />
    </div>
  )
}
