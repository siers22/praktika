import { TextareaHTMLAttributes, forwardRef } from 'react'
import { cn } from '@/lib/utils'

interface TextareaProps extends TextareaHTMLAttributes<HTMLTextAreaElement> {
  label?: string
  error?: string
}

const Textarea = forwardRef<HTMLTextAreaElement, TextareaProps>(
  ({ className, label, error, id, ...props }, ref) => {
    return (
      <div className="flex flex-col gap-1.5">
        {label && (
          <label htmlFor={id} className="text-xs font-mono font-semibold text-gray-500 uppercase tracking-wider">
            {label}
          </label>
        )}
        <textarea
          ref={ref}
          id={id}
          rows={3}
          className={cn(
            'w-full rounded-xl bg-surface px-4 py-2.5 font-mono text-sm text-gray-800',
            'shadow-neu-inset placeholder:text-gray-400 resize-none',
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

Textarea.displayName = 'Textarea'
export default Textarea
