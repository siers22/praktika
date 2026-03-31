import { cn } from '@/lib/utils'

interface StatCardProps {
  label: string
  value: string | number
  icon?: React.ReactNode
  className?: string
}

export default function StatCard({ label, value, icon, className }: StatCardProps) {
  return (
    <div className={cn('rounded-2xl bg-surface p-5 shadow-neu flex items-center gap-4', className)}>
      {icon && (
        <div className="flex-shrink-0 rounded-xl bg-surface p-3 shadow-neu-sm text-primary">
          {icon}
        </div>
      )}
      <div>
        <p className="text-xs font-mono font-semibold text-gray-400 uppercase tracking-wider">{label}</p>
        <p className="text-2xl font-mono font-bold text-gray-800 mt-0.5">{value}</p>
      </div>
    </div>
  )
}
