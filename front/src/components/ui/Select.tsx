import { SelectHTMLAttributes, forwardRef } from 'react'
import { cn } from '@/lib/utils'

interface SelectProps extends SelectHTMLAttributes<HTMLSelectElement> {
  label?: string
  error?: string
  options: { value: string; label: string }[]
  placeholder?: string
}

const Select = forwardRef<HTMLSelectElement, SelectProps>(
  ({ className, label, error, id, options, placeholder, ...props }, ref) => {
    return (
      <div className="flex flex-col gap-1.5">
        {label && (
          <label htmlFor={id} className="text-xs font-mono font-semibold text-gray-500 uppercase tracking-wider">
            {label}
          </label>
        )}
        <select
          ref={ref}
          id={id}
          className={cn(
            'w-full rounded-xl bg-surface px-4 py-2.5 font-mono text-sm text-gray-800',
            'shadow-neu-inset appearance-none cursor-pointer',
            'focus:outline-none focus:ring-2 focus:ring-primary/30',
            'transition-all duration-150',
            error && 'ring-2 ring-red-400',
            className
          )}
          {...props}
        >
          {placeholder && <option value="">{placeholder}</option>}
          {options.map((o) => (
            <option key={o.value} value={o.value}>
              {o.label}
            </option>
          ))}
        </select>
        {error && <p className="text-xs text-red-500 font-mono">{error}</p>}
      </div>
    )
  }
)

Select.displayName = 'Select'
export default Select
