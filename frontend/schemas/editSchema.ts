// Schemas Zod para validação de edições

import { z } from 'zod'

export const editInstructionSchema = z.object({
  type: z.enum(['text', 'image', 'drawing'], {
    errorMap: () => ({ message: 'Tipo deve ser text, image ou drawing' }),
  }),
  page: z.number().int().min(1, 'Página deve ser maior que 0'),
  x: z.number().min(0, 'Coordenada X deve ser maior ou igual a 0'),
  y: z.number().min(0, 'Coordenada Y deve ser maior ou igual a 0'),
  width: z.number().positive().optional(),
  height: z.number().positive().optional(),
  content: z.string().optional(),
  fontSize: z.number().positive().optional(),
  metadata: z.record(z.unknown()).optional(),
})

export const processDocumentRequestSchema = z.object({
  instructions: z
    .array(editInstructionSchema)
    .min(1, 'Deve haver pelo menos uma instrução de edição'),
})

export type EditInstructionFormData = z.infer<typeof editInstructionSchema>
export type ProcessDocumentRequestFormData = z.infer<
  typeof processDocumentRequestSchema
>
