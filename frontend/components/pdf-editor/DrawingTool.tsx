'use client'

import { useState } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Slider } from '@/components/ui/slider'
import type { EditInstruction } from '@/types/edit'

interface DrawingToolProps {
  onAdd: (instruction: EditInstruction) => void
  onCancel: () => void
}

export function DrawingTool({ onAdd, onCancel }: DrawingToolProps) {
  const [x1, setX1] = useState(100)
  const [y1, setY1] = useState(100)
  const [x2, setX2] = useState(200)
  const [y2, setY2] = useState(200)
  const [strokeWidth, setStrokeWidth] = useState(2)

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()

    onAdd({
      type: 'drawing',
      page: 1, // Será atualizado pelo componente pai
      x: x1,
      y: y1,
      width: Math.abs(x2 - x1),
      height: Math.abs(y2 - y1),
      metadata: {
        x2,
        y2,
        strokeWidth,
        drawingType: 'line',
      },
    })

    // Reset form
    setX1(100)
    setY1(100)
    setX2(200)
    setY2(200)
    setStrokeWidth(2)
  }

  return (
    <Card className="w-80">
      <CardHeader>
        <CardTitle>Adicionar Desenho</CardTitle>
        <CardDescription>
          Configure o desenho que será adicionado ao PDF
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label>Ponto Inicial</Label>
            <div className="grid grid-cols-2 gap-2">
              <div className="space-y-1">
                <Label htmlFor="x1" className="text-xs">X</Label>
                <Input
                  id="x1"
                  type="number"
                  value={x1}
                  onChange={(e) => setX1(Number(e.target.value))}
                  min={0}
                />
              </div>
              <div className="space-y-1">
                <Label htmlFor="y1" className="text-xs">Y</Label>
                <Input
                  id="y1"
                  type="number"
                  value={y1}
                  onChange={(e) => setY1(Number(e.target.value))}
                  min={0}
                />
              </div>
            </div>
          </div>

          <div className="space-y-2">
            <Label>Ponto Final</Label>
            <div className="grid grid-cols-2 gap-2">
              <div className="space-y-1">
                <Label htmlFor="x2" className="text-xs">X</Label>
                <Input
                  id="x2"
                  type="number"
                  value={x2}
                  onChange={(e) => setX2(Number(e.target.value))}
                  min={0}
                />
              </div>
              <div className="space-y-1">
                <Label htmlFor="y2" className="text-xs">Y</Label>
                <Input
                  id="y2"
                  type="number"
                  value={y2}
                  onChange={(e) => setY2(Number(e.target.value))}
                  min={0}
                />
              </div>
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="strokeWidth">
              Espessura da Linha: {strokeWidth}px
            </Label>
            <Slider
              id="strokeWidth"
              min={1}
              max={10}
              step={1}
              value={[strokeWidth]}
              onValueChange={([value]) => setStrokeWidth(value)}
            />
          </div>

          <div className="flex gap-2">
            <Button type="submit" className="flex-1">
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
