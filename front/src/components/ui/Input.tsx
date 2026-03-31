import { InputHTMLAttributes, forwardRef } from 'react'
import { cn } from '@/lib/utils'

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string
  error?: string
}

const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ className, label, error, id, ...props }, ref) => {
    return (
      <div className="flex flex-col gap-1.5">
        {label && (
          <label htmlFor={id} className="text-xs font-mono font-semibold text-gray-500 uppercase tracking-wider">
            {label}
          </label>
        )}
        <input
          ref={ref}
          id={id}
          className={cn(
            'w-full rounded-xl bg-surface px-4 py-2.5 font-mono text-sm text-gray-800',
            'shadow-neu-inset placeholder:text-gray-400',
            'focus:outline-none focus:ring-2 focus:ring-primary/30',
            'transition-all duration-150',
            error && 'ring-2 ring-red-400',
            className
          )}
          {...props}
        />
        {error && <p className="text-xs text-red-500 font-mono">{error}</p>}
      </div>
    )
  }
)

Input.displayName = 'Input'
export default Input
