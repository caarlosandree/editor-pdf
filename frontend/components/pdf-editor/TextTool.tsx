'use client'

import { useState } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Slider } from '@/components/ui/slider'
import type { EditInstruction } from '@/types/edit'

interface TextToolProps {
  onAdd: (instruction: EditInstruction) => void
  onCancel: () => void
}

export function TextTool({ onAdd, onCancel }: TextToolProps) {
  const [text, setText] = useState('')
  const [fontSize, setFontSize] = useState(12)
  const [x, setX] = useState(100)
  const [y, setY] = useState(100)

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!text.trim()) {
      return
    }

    onAdd({
      type: 'text',
      page: 1, // Será atualizado pelo componente pai
      x,
      y,
      content: text,
      fontSize,
    })

    // Reset form
    setText('')
    setFontSize(12)
    setX(100)
    setY(100)
  }

  return (
    <Card className="w-80">
      <CardHeader>
        <CardTitle>Adicionar Texto</CardTitle>
        <CardDescription>
          Configure o texto que será adicionado ao PDF
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="text">Texto</Label>
            <Input
              id="text"
              value={text}
              onChange={(e) => setText(e.target.value)}
              placeholder="Digite o texto..."
              required
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="fontSize">
              Tamanho da Fonte: {fontSize}pt
            </Label>
            <Slider
              id="fontSize"
              min={8}
              max={72}
              step={1}
              value={[fontSize]}
              onValueChange={([value]) => setFontSize(value)}
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="x">Posição X</Label>
              <Input
                id="x"
                type="number"
                value={x}
                onChange={(e) => setX(Number(e.target.value))}
                min={0}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="y">Posição Y</Label>
              <Input
                id="y"
                type="number"
                value={y}
                onChange={(e) => setY(Number(e.target.value))}
                min={0}
              />
            </div>
          </div>

          <div className="flex gap-2">
            <Button type="submit" className="flex-1" disabled={!text.trim()}>
              Adicionar
            </Button>
            <Button type="button" variant="outline" onClick={onCancel}>
              Cancelar
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  )
}
