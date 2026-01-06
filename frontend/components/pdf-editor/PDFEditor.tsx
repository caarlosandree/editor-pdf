'use client'

import { useState } from 'react'
import { PDFViewer } from './PDFViewer'
import { EditorToolbar } from './EditorToolbar'
import { TextTool } from './TextTool'
import { ImageTool } from './ImageTool'
import { DrawingTool } from './DrawingTool'
import type { EditInstruction, EditInstructionType } from '@/types/edit'
import { useProcessDocument } from '@/hooks/mutations/useProcessDocument'
import { useRouter } from 'next/navigation'

interface PDFEditorProps {
  documentId: string
  pageCount: number
}

export function PDFEditor({ documentId, pageCount }: PDFEditorProps) {
  const [currentPage, setCurrentPage] = useState(1)
  const [zoom, setZoom] = useState(1)
  const [selectedTool, setSelectedTool] = useState<EditInstructionType | null>(null)
  const [instructions, setInstructions] = useState<EditInstruction[]>([])
  const processMutation = useProcessDocument()
  const router = useRouter()

  const handleToolSelect = (tool: EditInstructionType | null) => {
    setSelectedTool(tool)
  }

  const handleAddInstruction = (instruction: EditInstruction) => {
    const newInstruction: EditInstruction = {
      ...instruction,
      page: currentPage,
    }
    setInstructions([...instructions, newInstruction])
    setSelectedTool(null) // Fecha a ferramenta após adicionar
  }

  const handleSave = () => {
    if (instructions.length === 0) {
      return
    }

    processMutation.mutate(
      {
        id: documentId,
        request: {
          instructions,
        },
      },
      {
        onSuccess: () => {
          setInstructions([])
          router.refresh()
        },
      }
    )
  }

  const handleCancel = () => {
    if (instructions.length > 0) {
      if (confirm('Tem certeza que deseja cancelar? As alterações não salvas serão perdidas.')) {
        setInstructions([])
        setSelectedTool(null)
      }
    } else {
      router.back()
    }
  }

  const renderTool = () => {
    if (!selectedTool) {
      return null
    }

    switch (selectedTool) {
      case 'text':
        return (
          <TextTool
            onAdd={handleAddInstruction}
            onCancel={() => setSelectedTool(null)}
          />
        )
      case 'image':
        return (
          <ImageTool
            onAdd={handleAddInstruction}
            onCancel={() => setSelectedTool(null)}
          />
        )
      case 'drawing':
        return (
          <DrawingTool
            onAdd={handleAddInstruction}
            onCancel={() => setSelectedTool(null)}
          />
        )
      default:
        return null
    }
  }

  return (
    <div className="flex h-full flex-col">
      <EditorToolbar
        selectedTool={selectedTool}
        onToolSelect={handleToolSelect}
        onSave={handleSave}
        onCancel={handleCancel}
        isSaving={processMutation.isPending}
        hasChanges={instructions.length > 0}
      />

      <div className="flex flex-1 overflow-hidden">
        <div className="flex-1">
          <PDFViewer
            documentId={documentId}
            pageCount={pageCount}
            currentPage={currentPage}
            onPageChange={setCurrentPage}
            zoom={zoom}
            onZoomChange={setZoom}
          />
        </div>

        {selectedTool && (
          <div className="w-80 border-l bg-background p-4 overflow-y-auto">
            {renderTool()}
          </div>
        )}

        {instructions.length > 0 && !selectedTool && (
          <div className="w-80 border-l bg-background p-4 overflow-y-auto">
            <div className="space-y-2">
              <h3 className="font-semibold">Instruções Pendentes</h3>
              <p className="text-sm text-muted-foreground">
                {instructions.length} instrução(ões) aguardando salvamento
              </p>
              <div className="space-y-2 mt-4">
                {instructions.map((instruction, index) => (
                  <div
                    key={index}
                    className="p-3 border rounded-lg text-sm"
                  >
                    <div className="font-medium capitalize">{instruction.type}</div>
                    <div className="text-xs text-muted-foreground mt-1">
                      Página {instruction.page} - ({instruction.x}, {instruction.y})
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
