'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { FileText, Home, LayoutDashboard } from 'lucide-react'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'

const navigation = [
  {
    name: 'Dashboard',
    href: '/',
    icon: LayoutDashboard,
  },
  {
    name: 'Documentos',
    href: '/documents',
    icon: FileText,
  },
]

export function Sidebar() {
  const pathname = usePathname()

  return (
    <div className="flex h-full w-64 flex-col border-r bg-background">
      <div className="flex h-16 items-center border-b px-6">
        <Link href="/" className="flex items-center gap-2">
          <FileText className="h-6 w-6 text-primary" />
          <span className="text-lg font-semibold">Editor PDF</span>
        </Link>
      </div>

      <nav className="flex-1 space-y-1 p-4">
        {navigation.map((item) => {
          const isActive = pathname === item.href
          return (
            <Link key={item.name} href={item.href}>
              <Button
                variant={isActive ? 'secondary' : 'ghost'}
                className={cn(
                  'w-full justify-start gap-3',
                  isActive && 'bg-secondary'
                )}
              >
                <item.icon className="h-5 w-5" />
                {item.name}
              </Button>
            </Link>
          )
        })}
      </nav>

      <div className="border-t p-4">
        <div className="text-xs text-muted-foreground">
          <p>Editor PDF v1.0</p>
        </div>
      </div>
    </div>
  )
}
