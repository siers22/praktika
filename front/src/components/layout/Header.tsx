'use client'

interface HeaderProps {
  title: string
  actions?: React.ReactNode
}

export default function Header({ title, actions }: HeaderProps) {
  return (
    <div className="flex items-center justify-between mb-6">
      <h2 className="font-mono font-bold text-gray-800 text-xl tracking-tight">{title}</h2>
      {actions && <div className="flex items-center gap-3">{actions}</div>}
    </div>
  )
}
