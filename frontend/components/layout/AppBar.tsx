'use client'

import { Upload } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useUploadDocument } from '@/hooks/mutations/useUploadDocument'
import { useRef } from 'react'

interface AppBarProps {
  title?: string
  onUploadClick?: () => void
}

export function AppBar({ title = 'Editor PDF', onUploadClick }: AppBarProps) {
  const uploadMutation = useUploadDocument()
  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (file) {
      uploadMutation.mutate(file)
    }
    // Reset input para permitir selecionar o mesmo arquivo novamente
    if (fileInputRef.current) {
      fileInputRef.current.value = ''
    }
  }

  const handleUploadClick = () => {
    if (onUploadClick) {
      onUploadClick()
    } else {
      fileInputRef.current?.click()
    }
  }

  return (
    <div className="flex h-16 items-center justify-between border-b bg-background px-6">
      <h1 className="text-xl font-semibold">{title}</h1>

      <div className="flex items-center gap-4">
        <input
          ref={fileInputRef}
          type="file"
          accept=".pdf,application/pdf"
          onChange={handleFileSelect}
          className="hidden"
        />
        <Button
          onClick={handleUploadClick}
          disabled={uploadMutation.isPending}
          size="sm"
        >
          <Upload className="h-4 w-4 mr-2" />
          {uploadMutation.isPending ? 'Enviando...' : 'Enviar PDF'}
        </Button>
      </div>
    </div>
  )
}
