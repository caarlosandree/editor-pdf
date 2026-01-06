'use client'

import { useDocument } from '@/hooks/queries/useDocument'
import { PDFEditor } from '@/components/pdf-editor/PDFEditor'
import { Skeleton } from '@/components/ui/skeleton'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { AlertCircle, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { use } from 'react'

interface PageProps {
  params: Promise<{ id: string }>
}

export default function EditDocumentPage({ params }: PageProps) {
  const { id } = use(params)
  const { data: document, isLoading, isError, error } = useDocument(id)
  const router = useRouter()

  if (isLoading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-16 w-full" />
        <Skeleton className="h-[600px] w-full" />
      </div>
    )
  }

  if (isError) {
    return (
      <div className="space-y-4">
        <Link href="/">
          <Button variant="ghost" size="sm">
            <ArrowLeft className="h-4 w-4 mr-2" />
            Voltar
          </Button>
        </Link>
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>Erro ao carregar documento</AlertTitle>
          <AlertDescription>
            {error instanceof Error ? error.message : 'Erro desconhecido'}
          </AlertDescription>
        </Alert>
      </div>
    )
  }

  if (!document) {
    return (
      <div className="space-y-4">
        <Link href="/">
          <Button variant="ghost" size="sm">
            <ArrowLeft className="h-4 w-4 mr-2" />
            Voltar
          </Button>
        </Link>
        <Alert>
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>Documento não encontrado</AlertTitle>
          <AlertDescription>
            O documento solicitado não foi encontrado.
          </AlertDescription>
        </Alert>
      </div>
    )
  }

  return (
    <div className="h-full">
      <PDFEditor documentId={document.id} pageCount={document.page_count} />
    </div>
  )
}
