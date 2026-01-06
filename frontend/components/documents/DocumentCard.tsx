'use client'

import Link from 'next/link'
import { FileText, Edit, Trash2, Calendar } from 'lucide-react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog'
import { Badge } from '@/components/ui/badge'
import { useDeleteDocument } from '@/hooks/mutations/useDeleteDocument'
import type { Document } from '@/types/document'
import { format } from 'date-fns'
import { ptBR } from 'date-fns/locale'

interface DocumentCardProps {
  document: Document
}

export function DocumentCard({ document }: DocumentCardProps) {
  const deleteMutation = useDeleteDocument()

  const handleDelete = () => {
    deleteMutation.mutate(document.id)
  }

  const statusColors = {
    uploaded: 'bg-blue-500',
    processing: 'bg-yellow-500',
    processed: 'bg-green-500',
    error: 'bg-red-500',
  }

  const statusLabels = {
    uploaded: 'Enviado',
    processing: 'Processando',
    processed: 'Processado',
    error: 'Erro',
  }

  return (
    <Card className="hover:shadow-md transition-shadow">
      <CardHeader>
        <div className="flex items-start justify-between">
          <div className="flex items-center gap-3">
            <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10">
              <FileText className="h-6 w-6 text-primary" />
            </div>
            <div>
              <CardTitle className="line-clamp-1">
                {document.file_path.split('/').pop() || 'Documento'}
              </CardTitle>
              <CardDescription className="flex items-center gap-2 mt-1">
                <Calendar className="h-3 w-3" />
                {format(new Date(document.created_at), "dd 'de' MMM 'de' yyyy", {
                  locale: ptBR,
                })}
              </CardDescription>
            </div>
          </div>
          <Badge
            variant="secondary"
            className={`${statusColors[document.status]} text-white`}
          >
            {statusLabels[document.status]}
          </Badge>
        </div>
      </CardHeader>
      <CardContent>
        <div className="flex flex-col gap-2">
          <div className="text-sm text-muted-foreground">
            <p>{document.page_count} página{document.page_count !== 1 ? 's' : ''}</p>
            <p>Versão {document.version}</p>
          </div>
          <div className="flex gap-2 mt-2">
            <Link href={`/documents/${document.id}/edit`} className="flex-1">
              <Button variant="default" className="w-full" size="sm">
                <Edit className="h-4 w-4 mr-2" />
                Editar
              </Button>
            </Link>
            <AlertDialog>
              <AlertDialogTrigger asChild>
                <Button
                  variant="destructive"
                  size="sm"
                  disabled={deleteMutation.isPending}
                >
                  <Trash2 className="h-4 w-4" />
                </Button>
              </AlertDialogTrigger>
              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle>Confirmar exclusão</AlertDialogTitle>
                  <AlertDialogDescription>
                    Tem certeza que deseja excluir este documento? Esta ação
                    não pode ser desfeita.
                  </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                  <AlertDialogCancel>Cancelar</AlertDialogCancel>
                  <AlertDialogAction onClick={handleDelete}>
                    Excluir
                  </AlertDialogAction>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
