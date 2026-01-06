'use client'

import { useState, useRef } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Upload } from 'lucide-react'
import type { EditInstruction } from '@/types/edit'

interface ImageToolProps {
  onAdd: (instruction: EditInstruction) => void
  onCancel: () => void
}

export function ImageTool({ onAdd, onCancel }: ImageToolProps) {
  const [imageFile, setImageFile] = useState<File | null>(null)
  const [preview, setPreview] = useState<string | null>(null)
  const [x, setX] = useState(100)
  const [y, setY] = useState(100)
  const [width, setWidth] = useState(200)
  const [height, setHeight] = useState(200)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file && file.type.startsWith('image/')) {
      setImageFile(file)
      const reader = new FileReader()
      reader.onloadend = () => {
        setPreview(reader.result as string)
      }
      reader.readAsDataURL(file)
    }
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!imageFile) {
      return
    }

    // Converte imagem para base64 para enviar na instrução
    const reader = new FileReader()
    reader.onloadend = () => {
      const base64 = reader.result as string
      onAdd({
        type: 'image',
        page: 1, // Será atualizado pelo componente pai
        x,
        y,
        width,
        height,
        content: base64,
        metadata: {
          filename: imageFile.name,
          mimeType: imageFile.type,
        },
      })

      // Reset form
      setImageFile(null)
      setPreview(null)
      setX(100)
      setY(100)
      setWidth(200)
      setHeight(200)
      if (fileInputRef.current) {
        fileInputRef.current.value = ''
      }
    }
    reader.readAsDataURL(imageFile)
  }

  return (
    <Card className="w-80">
      <CardHeader>
        <CardTitle>Adicionar Imagem</CardTitle>
        <CardDescription>
          Selecione uma imagem para adicionar ao PDF
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="image">Imagem</Label>
            <Input
              ref={fileInputRef}
              id="image"
              type="file"
              accept="image/*"
              onChange={handleFileSelect}
              className="cursor-pointer"
            />
            {preview && (
              <div className="mt-2">
                <img
                  src={preview}
                  alt="Preview"
                  className="max-w-full h-32 object-contain border rounded"
                />
              </div>
            )}
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

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="width">Largura</Label>
              <Input
                id="width"
                type="number"
                value={width}
                onChange={(e) => setWidth(Number(e.target.value))}
                min={1}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="height">Altura</Label>
              <Input
                id="height"
                type="number"
                value={height}
                onChange={(e) => setHeight(Number(e.target.value))}
                min={1}
              />
            </div>
          </div>

          <div className="flex gap-2">
            <Button type="submit" className="flex-1" disabled={!imageFile}>
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
