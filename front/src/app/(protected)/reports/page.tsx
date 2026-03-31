'use client'

import { useEffect, useState } from 'react'
import toast from 'react-hot-toast'
import { reports } from '@/lib/api'
import Header from '@/components/layout/Header'
import Card from '@/components/ui/Card'
import { formatCurrency } from '@/lib/utils'
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

interface DeptReport {
  department_id: string
  department_name: string
  total_count: number
  total_value: number
  by_status: Record<string, number>
}

const COLORS = ['#006666', '#0284c7', '#059669', '#d97706', '#dc2626', '#7c3aed']

export default function ReportsPage() {
  const [deptData, setDeptData] = useState<DeptReport[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    reports.byDepartment()
      .then((res: { data: DeptReport[] }) => setDeptData(res.data || []))
      .catch((e: Error) => toast.error(e.message))
      .finally(() => setLoading(false))
  }, [])

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <span className="h-8 w-8 animate-spin rounded-full border-2 border-primary border-t-transparent" />
      </div>
    )
  }

  const countChartData = deptData.map((d) => ({
    name: d.department_name,
    count: d.total_count,
  }))

  const valueChartData = deptData.map((d) => ({
    name: d.department_name,
    value: d.total_value,
  }))

  const totalCount = deptData.reduce((sum, d) => sum + d.total_count, 0)
  const totalValue = deptData.reduce((sum, d) => sum + d.total_value, 0)

  return (
    <div>
      <Header title="Отчёты" />

      {/* Summary */}
      <div className="grid grid-cols-2 gap-4 mb-6">
        <Card>
          <p className="text-xs font-mono font-semibold text-gray-400 uppercase tracking-wider">Всего единиц</p>
          <p className="text-3xl font-mono font-bold text-gray-800 mt-1">{totalCount}</p>
        </Card>
        <Card>
          <p className="text-xs font-mono font-semibold text-gray-400 uppercase tracking-wider">Общая стоимость</p>
          <p className="text-3xl font-mono font-bold text-gray-800 mt-1">{formatCurrency(totalValue)}</p>
        </Card>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
        <Card title="Количество по подразделениям">
          <ResponsiveContainer width="100%" height={280}>
            <BarChart data={countChartData} layout="vertical" margin={{ left: 0, right: 20 }}>
              <XAxis type="number" tick={{ fontFamily: 'monospace', fontSize: 11 }} />
              <YAxis type="category" dataKey="name" width={120} tick={{ fontFamily: 'monospace', fontSize: 11 }} />
              <Tooltip />
              <Bar dataKey="count" fill="#006666" radius={[0, 6, 6, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </Card>

        <Card title="Стоимость по подразделениям">
          <ResponsiveContainer width="100%" height={280}>
            <PieChart>
              <Pie
                data={valueChartData}
                dataKey="value"
                nameKey="name"
                cx="50%"
                cy="50%"
                outerRadius={90}
                label={({ name }) => name}
              >
                {valueChartData.map((_, i) => (
                  <Cell key={i} fill={COLORS[i % COLORS.length]} />
                ))}
              </Pie>
              <Legend formatter={(v) => <span style={{ fontFamily: 'monospace', fontSize: 11 }}>{v}</span>} />
              <Tooltip formatter={(v: number) => [formatCurrency(v), 'Стоимость']} />
            </PieChart>
          </ResponsiveContainer>
        </Card>
      </div>

      {/* Table by department */}
      <Card title="Детализация по подразделениям">
        <div className="overflow-x-auto">
          <table className="w-full text-xs font-mono">
            <thead>
              <tr className="border-b border-gray-200">
                <th className="text-left py-2 px-3 text-gray-400 uppercase tracking-wider font-semibold">Подразделение</th>
                <th className="text-right py-2 px-3 text-gray-400 uppercase tracking-wider font-semibold">Кол-во</th>
                <th className="text-right py-2 px-3 text-gray-400 uppercase tracking-wider font-semibold">Стоимость</th>
              </tr>
            </thead>
            <tbody>
              {deptData.map((d) => (
                <tr key={d.department_id} className="border-b border-gray-100 last:border-0">
                  <td className="py-2 px-3 text-gray-700 font-semibold">{d.department_name}</td>
                  <td className="py-2 px-3 text-right text-gray-600">{d.total_count}</td>
                  <td className="py-2 px-3 text-right text-gray-600">{formatCurrency(d.total_value)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </Card>
    </div>
  )
}
