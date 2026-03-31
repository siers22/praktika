'use client'

import { useCallback, useEffect, useState } from 'react'
import toast from 'react-hot-toast'
import { audit } from '@/lib/api'
import { AuditLog } from '@/types'
import Header from '@/components/layout/Header'
import Table from '@/components/ui/Table'
import Pagination from '@/components/ui/Pagination'
import Input from '@/components/ui/Input'
import Select from '@/components/ui/Select'
import { formatDateTime } from '@/lib/utils'
import { Search } from 'lucide-react'

const ACTION_OPTIONS = [
  { value: '', label: 'Все действия' },
  { value: 'create', label: 'Создание' },
  { value: 'update', label: 'Обновление' },
  { value: 'delete', label: 'Удаление' },
  { value: 'login', label: 'Вход' },
  { value: 'logout', label: 'Выход' },
]

export default function AuditPage() {
  const [data, setData] = useState<AuditLog[]>([])
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [action, setAction] = useState('')
  const [entityType, setEntityType] = useState('')

  const load = useCallback(() => {
    setLoading(true)
    const params: Record<string, string | number> = { page, per_page: 20 }
    if (action) params.action = action
    if (entityType) params.entity_type = entityType
    audit.list(params)
      .then((res: { data: AuditLog[]; meta: { total_pages: number } }) => {
        setData(res.data || [])
        setTotalPages(res.meta?.total_pages || 1)
      })
      .catch((e: Error) => toast.error(e.message))
      .finally(() => setLoading(false))
  }, [page, action, entityType])

  useEffect(() => { load() }, [load])

  const columns = [
    { key: 'created_at', header: 'Время', render: (row: AuditLog) => formatDateTime(row.created_at) },
    { key: 'username', header: 'Пользователь' },
    { key: 'action', header: 'Действие', className: 'font-semibold' },
    { key: 'entity_type', header: 'Тип объекта' },
    { key: 'entity_id', header: 'ID объекта', render: (row: AuditLog) => row.entity_id?.slice(0, 8) + '…' || '—' },
  ]

  return (
    <div>
      <Header title="Журнал аудита" />

      <div className="grid grid-cols-2 lg:grid-cols-3 gap-3 mb-5">
        <Select
          options={ACTION_OPTIONS}
          value={action}
          onChange={(e) => { setAction(e.target.value); setPage(1) }}
        />
        <div className="relative">
          <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
          <input
            className="w-full rounded-xl bg-surface pl-8 pr-4 py-2.5 font-mono text-sm text-gray-800 shadow-neu-inset focus:outline-none focus:ring-2 focus:ring-primary/30"
            placeholder="Тип объекта..."
            value={entityType}
            onChange={(e) => { setEntityType(e.target.value); setPage(1) }}
          />
        </div>
      </div>

      <Table columns={columns} data={data} keyField="id" loading={loading} emptyText="Записей нет" />
      <Pagination page={page} totalPages={totalPages} onPageChange={setPage} />
    </div>
  )
}
