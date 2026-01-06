import { z } from 'zod'

export const createUserSchema = z
  .object({
    name: z
      .string()
      .min(3, 'Nome deve ter no mínimo 3 caracteres')
      .max(100, 'Nome deve ter no máximo 100 caracteres'),
    email: z.string().email('Email inválido'),
    password: z.string().min(8, 'Senha deve ter no mínimo 8 caracteres'),
    confirmPassword: z.string(),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: 'Senhas não coincidem',
    path: ['confirmPassword'],
  })

export type CreateUserFormData = z.infer<typeof createUserSchema>

export const updateUserSchema = z.object({
  name: z
    .string()
    .min(3, 'Nome deve ter no mínimo 3 caracteres')
    .max(100, 'Nome deve ter no máximo 100 caracteres')
    .optional(),
  email: z.string().email('Email inválido').optional(),
})

export type UpdateUserFormData = z.infer<typeof updateUserSchema>
