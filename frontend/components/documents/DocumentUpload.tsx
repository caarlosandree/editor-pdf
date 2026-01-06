'use client'

import { useRef, useState } from 'react'
import { Upload, FileText, X } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { useUploadDocument } from '@/hooks/mutations/useUploadDocument'
import { Progress } from '@/components/ui/progress'

export function DocumentUpload() {
  const fileInputRef = useRef<HTMLInputElement>(null)
  const [selectedFile, setSelectedFile] = useState<File | null>(null)
  const uploadMutation = useUploadDocument()

  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (file) {
      if (file.type !== 'application/pdf') {
        alert('Por favor, selecione um arquivo PDF')
        return
      }
      setSelectedFile(file)
    }
  }

  const handleUpload = () => {
    if (selectedFile) {
      uploadMutation.mutate(selectedFile, {
        onSuccess: () => {
          setSelectedFile(null)
          if (fileInputRef.current) {
            fileInputRef.current.value = ''
          }
        },
      })
    }
  }

  const handleRemoveFile = () => {
    setSelectedFile(null)
    if (fileInputRef.current) {
      fileInputRef.current.value = ''
    }
  }

  return (
    <Card>
      <CardContent className="pt-6">
        <div className="space-y-4">
          <input
            ref={fileInputRef}
            type="file"
            accept=".pdf,application/pdf"
            onChange={handleFileSelect}
            className="hidden"
          />

          {!selectedFile ? (
            <Button
              variant="outline"
              className="w-full h-32 border-dashed"
              onClick={() => fileInputRef.current?.click()}
            >
              <div className="flex flex-col items-center gap-2">
                <Upload className="h-8 w-8 text-muted-foreground" />
                <span className="text-sm text-muted-foreground">
                  Clique para selecionar um arquivo PDF
                </span>
              </div>
            </Button>
          ) : (
            <div className="space-y-4">
              <div className="flex items-center gap-3 p-4 border rounded-lg">
                <FileText className="h-8 w-8 text-primary" />
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-medium truncate">
                    {selectedFile.name}
                  </p>
                  <p className="text-xs text-muted-foreground">
                    {(selectedFile.size / 1024 / 1024).toFixed(2)} MB
                  </p>
                </div>
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={handleRemoveFile}
                  disabled={uploadMutation.isPending}
                >
                  <X className="h-4 w-4" />
                </Button>
              </div>

              {uploadMutation.isPending && (
                <div className="space-y-2">
                  <Progress value={undefined} className="w-full" />
                  <p className="text-xs text-muted-foreground text-center">
                    Enviando documento...
                  </p>
                </div>
              )}

              <div className="flex gap-2">
                <Button
                  onClick={handleUpload}
                  disabled={uploadMutation.isPending}
                  className="flex-1"
                >
                  {uploadMutation.isPending ? 'Enviando...' : 'Enviar'}
                </Button>
                <Button
                  variant="outline"
                  onClick={handleRemoveFile}
                  disabled={uploadMutation.isPending}
                >
                  Cancelar
                </Button>
              </div>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  )
}
