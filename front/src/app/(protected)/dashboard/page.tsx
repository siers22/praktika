'use client'

import { useEffect, useState } from 'react'
import { reports } from '@/lib/api'
import { DashboardData, STATUS_LABELS } from '@/types'
import Header from '@/components/layout/Header'
import StatCard from '@/components/ui/StatCard'
import Card from '@/components/ui/Card'
import Badge from '@/components/ui/Badge'
import { formatDate, formatCurrency } from '@/lib/utils'
import { Package, AlertTriangle, ArrowLeftRight, ClipboardCheck } from 'lucide-react'
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
  Legend,
} from 'recharts'
import type { EquipmentStatus } from '@/types'

const STATUS_COLORS_HEX: Record<string, string> = {
  in_use: '#059669',
  in_storage: '#0284c7',
  in_repair: '#d97706',
  written_off: '#dc2626',
  reserved: '#006666',
}

export default function DashboardPage() {
  const [data, setData] = useState<DashboardData | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    reports.dashboard()
      .then((res: { data: DashboardData }) => setData(res.data))
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [])

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <span className="h-8 w-8 animate-spin rounded-full border-2 border-primary border-t-transparent" />
      </div>
    )
  }

  if (!data) return null

  const statusChartData = Object.entries(data.by_status || {}).map(([key, count]) => ({
    name: STATUS_LABELS[key as EquipmentStatus] || key,
    value: count,
    color: STATUS_COLORS_HEX[key] || '#9ca3af',
  }))

  const categoryChartData = (data.by_category || []).slice(0, 8)
  const recentMovements = data.recent_movements || []
  const warrantyExpiring = data.warranty_expiring_soon || []

  const inUse = data.by_status['in_use'] || 0
  const inRepair = data.by_status['in_repair'] || 0

  return (
    <div>
      <Header title="Дашборд" />

      {/* Stats */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
        <StatCard
          label="Всего единиц"
          value={data.total_equipment}
          icon={<Package size={20} />}
        />
        <StatCard
          label="В эксплуатации"
          value={inUse}
          icon={<ClipboardCheck size={20} />}
        />
        <StatCard
          label="В ремонте"
          value={inRepair}
          icon={<AlertTriangle size={20} />}
        />
        <StatCard
          label="Перемещений"
          value={recentMovements.length}
          icon={<ArrowLeftRight size={20} />}
        />
      </div>

      {/* Charts */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
        <Card title="По статусу">
          <ResponsiveContainer width="100%" height={220}>
            <PieChart>
              <Pie
                data={statusChartData}
                dataKey="value"
                nameKey="name"
                cx="50%"
                cy="50%"
                outerRadius={80}
                label={({ name, value }) => `${value}`}
              >
                {statusChartData.map((entry, i) => (
                  <Cell key={i} fill={entry.color} />
                ))}
              </Pie>
              <Legend
                formatter={(value) => (
                  <span style={{ fontFamily: 'monospace', fontSize: '11px' }}>{value}</span>
                )}
              />
              <Tooltip formatter={(v: number) => [v, 'Кол-во']} />
            </PieChart>
          </ResponsiveContainer>
        </Card>

        <Card title="По категориям">
          <ResponsiveContainer width="100%" height={220}>
            <BarChart data={categoryChartData} layout="vertical" margin={{ left: 0, right: 20 }}>
              <XAxis type="number" tick={{ fontFamily: 'monospace', fontSize: 11 }} />
              <YAxis
                type="category"
                dataKey="category"
                width={100}
                tick={{ fontFamily: 'monospace', fontSize: 11 }}
              />
              <Tooltip />
              <Bar dataKey="count" fill="#006666" radius={[0, 6, 6, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </Card>
      </div>

      {/* Bottom row */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Warranty expiring */}
        {warrantyExpiring.length > 0 && (
          <Card title="Гарантия истекает (30 дней)">
            <ul className="space-y-2">
              {warrantyExpiring.map((eq) => (
                <li key={eq.id} className="flex items-center justify-between text-xs font-mono">
                  <span className="text-gray-700 truncate max-w-[60%]">{eq.name}</span>
                  <span className="text-amber-600">{formatDate(eq.warranty_expiry)}</span>
                </li>
              ))}
            </ul>
          </Card>
        )}

        {/* Recent movements */}
        {recentMovements.length > 0 && (
          <Card title="Последние перемещения">
            <ul className="space-y-2">
              {recentMovements.slice(0, 6).map((m) => (
                <li key={m.id} className="flex items-center justify-between text-xs font-mono">
                  <span className="text-gray-700 truncate max-w-[55%]">{m.equipment_name || m.inventory_number}</span>
                  <span className="text-gray-400">{formatDate(m.moved_at)}</span>
                </li>
              ))}
            </ul>
          </Card>
        )}
      </div>
    </div>
  )
}
