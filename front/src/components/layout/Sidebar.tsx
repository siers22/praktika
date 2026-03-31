'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { useAuth } from '@/contexts/AuthContext'
import { cn } from '@/lib/utils'
import {
  LayoutDashboard,
  Package,
  Tags,
  Building2,
  ArrowLeftRight,
  ClipboardList,
  BarChart3,
  Users,
  ScrollText,
  LogOut,
} from 'lucide-react'

interface NavItem {
  href: string
  label: string
  icon: React.ReactNode
  roles?: ('admin' | 'inventory' | 'viewer')[]
}

const navItems: NavItem[] = [
  { href: '/dashboard', label: 'Дашборд', icon: <LayoutDashboard size={18} /> },
  { href: '/equipment', label: 'Оборудование', icon: <Package size={18} /> },
  { href: '/movements', label: 'Перемещения', icon: <ArrowLeftRight size={18} /> },
  { href: '/inventory', label: 'Инвентаризация', icon: <ClipboardList size={18} />, roles: ['admin', 'inventory'] },
  { href: '/categories', label: 'Категории', icon: <Tags size={18} />, roles: ['admin'] },
  { href: '/departments', label: 'Подразделения', icon: <Building2 size={18} />, roles: ['admin'] },
  { href: '/reports', label: 'Отчёты', icon: <BarChart3 size={18} /> },
  { href: '/users', label: 'Пользователи', icon: <Users size={18} />, roles: ['admin'] },
  { href: '/audit', label: 'Журнал аудита', icon: <ScrollText size={18} />, roles: ['admin'] },
]

export default function Sidebar() {
  const pathname = usePathname()
  const { user, logout, hasRole } = useAuth()

  const visible = navItems.filter(
    (item) => !item.roles || item.roles.some((r) => hasRole(r))
  )

  return (
    <aside className="flex h-screen w-56 flex-col bg-surface shadow-[6px_0_20px_rgba(0,0,0,0.06)] z-10">
      {/* Logo */}
      <div className="px-6 py-5 border-b border-gray-200">
        <h1 className="font-mono font-bold text-primary text-base tracking-tight leading-tight">
          Инвент.<br />
          <span className="text-gray-400 font-normal text-xs">оборудования</span>
        </h1>
      </div>

      {/* Nav */}
      <nav className="flex-1 overflow-y-auto py-4 px-3">
        <ul className="space-y-1">
          {visible.map((item) => {
            const active = pathname === item.href || pathname.startsWith(item.href + '/')
            return (
              <li key={item.href}>
                <Link
                  href={item.href}
                  className={cn(
                    'flex items-center gap-3 px-3 py-2.5 rounded-xl font-mono text-xs font-medium transition-all duration-150',
                    active
                      ? 'text-primary shadow-neu-inset bg-surface'
                      : 'text-gray-500 hover:text-gray-700 hover:shadow-neu-sm hover:bg-surface'
                  )}
                >
                  <span className={active ? 'text-primary' : 'text-gray-400'}>{item.icon}</span>
                  {item.label}
                </Link>
              </li>
            )
          })}
        </ul>
      </nav>

      {/* User / logout */}
      <div className="border-t border-gray-200 px-3 py-3">
        <div className="rounded-xl bg-surface shadow-neu-inset px-3 py-2.5 mb-2">
          <p className="text-xs font-mono font-semibold text-gray-700 truncate">{user?.id?.slice(0, 8)}…</p>
          <p className="text-[10px] font-mono text-gray-400 capitalize">{user?.role}</p>
        </div>
        <button
          onClick={logout}
          className="flex w-full items-center gap-2 rounded-xl px-3 py-2 text-xs font-mono text-gray-500 hover:text-red-500 hover:shadow-neu-sm transition-all"
        >
          <LogOut size={14} />
          Выйти
        </button>
      </div>
    </aside>
  )
}
