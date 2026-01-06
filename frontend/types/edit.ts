// Tipos relacionados a edições de PDF

export type EditInstructionType = 'text' | 'image' | 'drawing'

export interface EditInstruction {
  type: EditInstructionType
  page: number
  x: number
  y: number
  width?: number
  height?: number
  content?: string
  fontSize?: number
  metadata?: Record<string, unknown>
}

export interface ProcessDocumentRequest {
  instructions: EditInstruction[]
}
