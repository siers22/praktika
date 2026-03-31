import { cn } from '@/lib/utils'

export type BadgeVariant = 'success' | 'warning' | 'danger' | 'info' | 'primary' | 'default'

interface BadgeProps {
  variant?: BadgeVariant
  children: React.ReactNode
  className?: string
}

const variantStyles: Record<BadgeVariant, string> = {
  success: 'text-emerald-700 bg-emerald-100',
  warning: 'text-amber-700 bg-amber-100',
  danger: 'text-red-700 bg-red-100',
  info: 'text-sky-700 bg-sky-100',
  primary: 'text-teal-700 bg-teal-100',
  default: 'text-gray-600 bg-gray-200',
}

export default function Badge({ variant = 'default', children, className }: BadgeProps) {
  return (
    <span
      className={cn(
        'inline-flex items-center rounded-lg px-2.5 py-0.5 text-xs font-mono font-semibold',
        variantStyles[variant],
        className
      )}
    >
      {children}
    </span>
  )
}
