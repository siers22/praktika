import { cn } from '@/lib/utils'
import { HTMLAttributes } from 'react'

interface CardProps extends HTMLAttributes<HTMLDivElement> {
  title?: string
  actions?: React.ReactNode
}

export default function Card({ className, title, actions, children, ...props }: CardProps) {
  return (
    <div
      className={cn('rounded-2xl bg-surface p-6 shadow-neu', className)}
      {...props}
    >
      {(title || actions) && (
        <div className="mb-4 flex items-center justify-between">
          {title && <h3 className="font-mono font-semibold text-gray-700 text-sm uppercase tracking-wider">{title}</h3>}
          {actions && <div className="flex gap-2">{actions}</div>}
        </div>
      )}
      {children}
    </div>
  )
}
