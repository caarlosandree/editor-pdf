'use client'

import { Type, Image, PenTool, Save, X } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import { cn } from '@/lib/utils'
import type { EditInstructionType } from '@/types/edit'

interface EditorToolbarProps {
  selectedTool: EditInstructionType | null
  onToolSelect: (tool: EditInstructionType | null) => void
  onSave: () => void
  onCancel: () => void
  isSaving?: boolean
  hasChanges?: boolean
}

export function EditorToolbar({
  selectedTool,
  onToolSelect,
  onSave,
  onCancel,
  isSaving = false,
  hasChanges = false,
}: EditorToolbarProps) {
  const tools: Array<{
    type: EditInstructionType
    label: string
    icon: React.ComponentType<{ className?: string }>
  }> = [
    { type: 'text', label: 'Texto', icon: Type },
    { type: 'image', label: 'Imagem', icon: Image },
    { type: 'drawing', label: 'Desenho', icon: PenTool },
  ]

  return (
    <div className="flex items-center justify-between p-4 border-b bg-background">
      <div className="flex items-center gap-2">
        <span className="text-sm font-medium text-muted-foreground mr-2">
          Ferramentas:
        </span>
        {tools.map((tool) => {
          const Icon = tool.icon
          const isActive = selectedTool === tool.type
          return (
            <Button
              key={tool.type}
              variant={isActive ? 'default' : 'outline'}
              size="sm"
              onClick={() => onToolSelect(isActive ? null : tool.type)}
              className={cn(isActive && 'bg-primary')}
            >
              <Icon className="h-4 w-4 mr-2" />
              {tool.label}
            </Button>
          )
        })}
        {selectedTool && (
          <>
            <Separator orientation="vertical" className="h-6 mx-2" />
            <Button
              variant="ghost"
              size="sm"
              onClick={() => onToolSelect(null)}
            >
              <X className="h-4 w-4 mr-2" />
              Cancelar
            </Button>
          </>
        )}
      </div>

      <div className="flex items-center gap-2">
        {hasChanges && (
          <span className="text-xs text-muted-foreground mr-2">
            Alterações não salvas
          </span>
        )}
        <Button variant="outline" size="sm" onClick={onCancel}>
          Cancelar
        </Button>
        <Button
          variant="default"
          size="sm"
          onClick={onSave}
          disabled={isSaving || !hasChanges}
        >
          <Save className="h-4 w-4 mr-2" />
          {isSaving ? 'Salvando...' : 'Salvar'}
        </Button>
      </div>
    </div>
  )
}
