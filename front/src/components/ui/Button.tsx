import { ButtonHTMLAttributes, forwardRef } from 'react'
import { cn } from '@/lib/utils'

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'danger' | 'ghost'
  size?: 'sm' | 'md' | 'lg'
  loading?: boolean
}

const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant = 'primary', size = 'md', loading, disabled, children, ...props }, ref) => {
    const base =
      'inline-flex items-center justify-center font-mono font-semibold tracking-wide transition-all duration-150 focus:outline-none disabled:opacity-50 disabled:cursor-not-allowed rounded-xl'

    const variants = {
      primary:
        'bg-surface text-primary shadow-neu hover:shadow-neu-sm active:shadow-neu-inset',
      secondary:
        'bg-surface text-gray-700 shadow-neu hover:shadow-neu-sm active:shadow-neu-inset',
      danger:
        'bg-surface text-red-600 shadow-neu hover:shadow-neu-sm active:shadow-neu-inset',
      ghost:
        'bg-transparent text-gray-600 hover:bg-surface hover:shadow-neu-sm',
    }

    const sizes = {
      sm: 'px-3 py-1.5 text-xs',
      md: 'px-5 py-2.5 text-sm',
      lg: 'px-7 py-3.5 text-base',
    }

    return (
      <button
        ref={ref}
        className={cn(base, variants[variant], sizes[size], className)}
        disabled={disabled || loading}
        {...props}
      >
        {loading ? (
          <span className="mr-2 inline-block h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
        ) : null}
        {children}
      </button>
    )
  }
)

Button.displayName = 'Button'
export default Button
