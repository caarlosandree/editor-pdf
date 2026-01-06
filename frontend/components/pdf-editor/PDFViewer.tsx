'use client'

import { useState, useEffect, useRef } from 'react'
import { ChevronLeft, ChevronRight, ZoomIn, ZoomOut, RotateCw } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { documentService } from '@/services/documentService'
import { cn } from '@/lib/utils'

interface PDFViewerProps {
  documentId: string
  pageCount: number
  currentPage: number
  onPageChange: (page: number) => void
  zoom?: number
  onZoomChange?: (zoom: number) => void
  className?: string
}

export function PDFViewer({
  documentId,
  pageCount,
  currentPage,
  onPageChange,
  zoom = 1,
  onZoomChange,
  className,
}: PDFViewerProps) {
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [imageUrl, setImageUrl] = useState<string | null>(null)
  const containerRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (!documentId || currentPage < 1 || currentPage > pageCount) {
      return
    }

    setLoading(true)
    setError(null)

    const previewUrl = documentService.getPreviewUrl(documentId, currentPage)
    
    // Adiciona timestamp para evitar cache
    const urlWithTimestamp = `${previewUrl}?t=${Date.now()}`

    // Carrega a imagem diretamente (sem autenticação)
    const img = new Image()
    img.onload = () => {
      setImageUrl(urlWithTimestamp)
      setLoading(false)
    }
    img.onerror = () => {
      setError('Erro ao carregar página')
      setLoading(false)
    }
    img.src = urlWithTimestamp
  }, [documentId, currentPage, pageCount])

  const handlePreviousPage = () => {
    if (currentPage > 1) {
      onPageChange(currentPage - 1)
    }
  }

  const handleNextPage = () => {
    if (currentPage < pageCount) {
      onPageChange(currentPage + 1)
    }
  }

  const handleZoomIn = () => {
    if (onZoomChange) {
      onZoomChange(Math.min(zoom + 0.25, 3))
    }
  }

  const handleZoomOut = () => {
    if (onZoomChange) {
      onZoomChange(Math.max(zoom - 0.25, 0.5))
    }
  }

  const handleResetZoom = () => {
    if (onZoomChange) {
      onZoomChange(1)
    }
  }

  return (
    <div className={cn('flex flex-col h-full', className)}>
      {/* Controles */}
      <div className="flex items-center justify-between p-4 border-b bg-background">
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="icon"
            onClick={handlePreviousPage}
            disabled={currentPage <= 1}
          >
            <ChevronLeft className="h-4 w-4" />
          </Button>
          <span className="text-sm font-medium min-w-[120px] text-center">
            Página {currentPage} de {pageCount}
          </span>
          <Button
            variant="outline"
            size="icon"
            onClick={handleNextPage}
            disabled={currentPage >= pageCount}
          >
            <ChevronRight className="h-4 w-4" />
          </Button>
        </div>

        <div className="flex items-center gap-2">
          <Button variant="outline" size="icon" onClick={handleZoomOut}>
            <ZoomOut className="h-4 w-4" />
          </Button>
          <span className="text-sm font-medium min-w-[60px] text-center">
            {Math.round(zoom * 100)}%
          </span>
          <Button variant="outline" size="icon" onClick={handleZoomIn}>
            <ZoomIn className="h-4 w-4" />
          </Button>
          <Button variant="outline" size="icon" onClick={handleResetZoom}>
            <RotateCw className="h-4 w-4" />
          </Button>
        </div>
      </div>

      {/* Área de visualização */}
      <div
        ref={containerRef}
        className="flex-1 overflow-auto bg-gray-100 dark:bg-gray-900 p-4"
      >
        <div className="flex items-center justify-center min-h-full">
          {loading && (
            <div className="space-y-4">
              <Skeleton className="w-[800px] h-[1000px]" />
            </div>
          )}

          {error && (
            <div className="text-center text-destructive">
              <p>{error}</p>
            </div>
          )}

          {!loading && !error && imageUrl && (
            <div
              className="relative"
              style={{
                transform: `scale(${zoom})`,
                transformOrigin: 'top center',
                transition: 'transform 0.2s',
              }}
            >
              <img
                src={imageUrl}
                alt={`Página ${currentPage}`}
                className="max-w-full h-auto shadow-lg"
                style={{ maxHeight: 'none' }}
                onError={() => setError('Erro ao carregar imagem')}
              />
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
