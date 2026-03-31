'use client'

import { useCallback, useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import toast from 'react-hot-toast'
import { equipment, categories, departments } from '@/lib/api'
import { Equipment, EquipmentStatus, STATUS_LABELS, STATUS_COLORS } from '@/types'
import type { BadgeVariant } from '@/components/ui/Badge'
import Header from '@/components/layout/Header'
import Table from '@/components/ui/Table'
import Badge from '@/components/ui/Badge'
import Button from '@/components/ui/Button'
import Pagination from '@/components/ui/Pagination'
import Input from '@/components/ui/Input'
import Select from '@/components/ui/Select'
import { formatCurrency, formatDate, downloadBlob } from '@/lib/utils'
import { Plus, Download, Search } from 'lucide-react'

const STATUS_OPTIONS = [
  { value: '', label: 'Все статусы' },
  ...Object.entries(STATUS_LABELS).map(([v, l]) => ({ value: v, label: l })),
]

export default function EquipmentPage() {
  const router = useRouter()
  const [data, setData] = useState<Equipment[]>([])
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [search, setSearch] = useState('')
  const [status, setStatus] = useState('')
  const [categoryId, setCategoryId] = useState('')
  const [departmentId, setDepartmentId] = useState('')
  const [categoryOptions, setCategoryOptions] = useState<{ value: string; label: string }[]>([])
  const [deptOptions, setDeptOptions] = useState<{ value: string; label: string }[]>([])
  const [exportLoading, setExportLoading] = useState(false)

  useEffect(() => {
    categories.list().then((res: { data: { id: string; name: string }[] }) => {
      setCategoryOptions([{ value: '', label: 'Все категории' }, ...res.data.map((c) => ({ value: c.id, label: c.name }))])
    })
    departments.list().then((res: { data: { id: string; name: string }[] }) => {
      setDeptOptions([{ value: '', label: 'Все подразделения' }, ...res.data.map((d) => ({ value: d.id, label: d.name }))])
    })
  }, [])

  const load = useCallback(() => {
    setLoading(true)
    const params: Record<string, string | number> = { page, per_page: 20 }
    if (search) params.search = search
    if (status) params.status = status
    if (categoryId) params.category_id = categoryId
    if (departmentId) params.department_id = departmentId

    equipment.list(params)
      .then((res: { data: Equipment[]; meta: { total_pages: number } }) => {
        setData(res.data || [])
        setTotalPages(res.meta?.total_pages || 1)
      })
      .catch((e: Error) => toast.error(e.message))
      .finally(() => setLoading(false))
  }, [page, search, status, categoryId, departmentId])

  useEffect(() => { load() }, [load])

  const handleExport = async () => {
    setExportLoading(true)
    try {
      const res = await equipment.exportCSV()
      const blob = await res.blob()
      downloadBlob(blob, `equipment_${new Date().toISOString().slice(0, 10)}.csv`)
    } catch {
      toast.error('Ошибка экспорта')
    } finally {
      setExportLoading(false)
    }
  }

  const columns = [
    { key: 'inventory_number', header: 'Инв. номер', className: 'font-semibold' },
    { key: 'name', header: 'Наименование' },
    {
      key: 'status',
      header: 'Статус',
      render: (row: Equipment) => (
        <Badge variant={STATUS_COLORS[row.status] as BadgeVariant}>
          {STATUS_LABELS[row.status]}
        </Badge>
      ),
    },
    { key: 'category_name', header: 'Категория' },
    { key: 'department_name', header: 'Подразделение' },
    {
      key: 'purchase_price',
      header: 'Стоимость',
      render: (row: Equipment) => formatCurrency(row.purchase_price),
    },
    {
      key: 'warranty_expiry',
      header: 'Гарантия до',
      render: (row: Equipment) => formatDate(row.warranty_expiry),
    },
  ]

  return (
    <div>
      <Header
        title="Оборудование"
        actions={
          <>
            <Button variant="secondary" size="sm" onClick={handleExport} loading={exportLoading}>
              <Download size={14} className="mr-1.5" />
              CSV
            </Button>
            <Button size="sm" onClick={() => router.push('/equipment/new')}>
              <Plus size={14} className="mr-1.5" />
              Добавить
            </Button>
          </>
        }
      />

      {/* Filters */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-3 mb-5">
        <div className="relative">
          <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
          <input
            className="w-full rounded-xl bg-surface pl-8 pr-4 py-2.5 font-mono text-sm text-gray-800 shadow-neu-inset focus:outline-none focus:ring-2 focus:ring-primary/30"
            placeholder="Поиск..."
            value={search}
            onChange={(e) => { setSearch(e.target.value); setPage(1) }}
          />
        </div>
        <Select
          options={STATUS_OPTIONS}
          value={status}
          onChange={(e) => { setStatus(e.target.value); setPage(1) }}
        />
        <Select
          options={categoryOptions}
          value={categoryId}
          onChange={(e) => { setCategoryId(e.target.value); setPage(1) }}
        />
        <Select
          options={deptOptions}
          value={departmentId}
          onChange={(e) => { setDepartmentId(e.target.value); setPage(1) }}
        />
      </div>

      <Table
        columns={columns}
        data={data}
        keyField="id"
        loading={loading}
        onRowClick={(row) => router.push(`/equipment/${row.id}`)}
        emptyText="Оборудование не найдено"
      />
      <Pagination page={page} totalPages={totalPages} onPageChange={setPage} />
    </div>
  )
}
