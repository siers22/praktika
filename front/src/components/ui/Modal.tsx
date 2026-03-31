'use client'

import { useEffect } from 'react'
import { cn } from '@/lib/utils'
import { X } from 'lucide-react'

interface ModalProps {
  open: boolean
  onClose: () => void
  title?: string
  children: React.ReactNode
  size?: 'sm' | 'md' | 'lg' | 'xl'
}

const sizeClasses = {
  sm: 'max-w-sm',
  md: 'max-w-md',
  lg: 'max-w-lg',
  xl: 'max-w-2xl',
}

export default function Modal({ open, onClose, title, children, size = 'md' }: ModalProps) {
  useEffect(() => {
    if (!open) return
    const handler = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose()
    }
    document.addEventListener('keydown', handler)
    return () => document.removeEventListener('keydown', handler)
  }, [open, onClose])

  if (!open) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div
        className="absolute inset-0 bg-black/20 backdrop-blur-sm"
        onClick={onClose}
      />
      <div
        className={cn(
          'relative w-full rounded-2xl bg-surface shadow-[0_20px_60px_rgba(0,0,0,0.15)] p-6',
          sizeClasses[size]
        )}
      >
        {title && (
          <div className="mb-5 flex items-center justify-between">
            <h2 className="font-mono font-bold text-gray-800 text-base">{title}</h2>
            <button
              onClick={onClose}
              className="rounded-xl p-1.5 text-gray-400 hover:text-gray-600 shadow-neu-sm hover:shadow-neu-sm active:shadow-neu-inset transition-all"
            >
              <X size={16} />
            </button>
          </div>
        )}
        {children}
      </div>
    </div>
  )
}
