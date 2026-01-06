'use client'

import { ReactNode } from 'react'
import { usePathname } from 'next/navigation'
import { Sidebar } from './Sidebar'
import { AppBar } from './AppBar'
import { cn } from '@/lib/utils'

interface AppLayoutProps {
  children: ReactNode
  title?: string
  onUploadClick?: () => void
}

export function AppLayout({
  children,
  title,
  onUploadClick,
}: AppLayoutProps) {
  const pathname = usePathname()
  const isEditorPage = pathname?.includes('/documents/') && pathname?.includes('/edit')

  return (
    <div className="flex h-screen overflow-hidden">
      <Sidebar />
      <div className="flex flex-1 flex-col overflow-hidden">
        {!isEditorPage && <AppBar title={title} onUploadClick={onUploadClick} />}
        <main
          className={cn(
            'flex-1 overflow-y-auto bg-background',
            !isEditorPage && 'p-6'
          )}
        >
          {children}
        </main>
      </div>
    </div>
  )
}
