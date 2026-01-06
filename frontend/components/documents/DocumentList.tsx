'use client'

import { useDocuments } from '@/hooks/queries/useDocuments'
import { DocumentCard } from './DocumentCard'
import { Skeleton } from '@/components/ui/skeleton'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { AlertCircle, FileText } from 'lucide-react'

export function DocumentList() {
  const { data, isLoading, isError, error } = useDocuments({ limit: 20, offset: 0 })

  if (isLoading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {Array.from({ length: 6 }).map((_, i) => (
          <div key={i} className="space-y-4">
            <Skeleton className="h-48 w-full" />
            <Skeleton className="h-4 w-3/4" />
            <Skeleton className="h-4 w-1/2" />
          </div>
        ))}
      </div>
    )
  }

  if (isError) {
    return (
      <Alert variant="destructive">
        <AlertCircle className="h-4 w-4" />
        <AlertTitle>Erro ao carregar documentos</AlertTitle>
        <AlertDescription>
          {error instanceof Error ? error.message : 'Erro desconhecido'}
        </AlertDescription>
      </Alert>
    )
  }

  if (!data || data.documents.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-12 text-center">
        <FileText className="h-12 w-12 text-muted-foreground mb-4" />
        <h3 className="text-lg font-semibold mb-2">Nenhum documento encontrado</h3>
        <p className="text-sm text-muted-foreground">
          Comece enviando seu primeiro documento PDF
        </p>
      </div>
    )
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {data.documents.map((document) => (
        <DocumentCard key={document.id} document={document} />
      ))}
    </div>
  )
}
