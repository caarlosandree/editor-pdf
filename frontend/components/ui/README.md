# Componentes shadcn/ui

Este diretório contém todos os componentes shadcn/ui instalados no projeto.

## Componentes Disponíveis

### Formulários
- **Input** (`input.tsx`) - Campo de entrada de texto
- **Label** (`label.tsx`) - Rótulo para campos de formulário
- **Textarea** (`textarea.tsx`) - Área de texto multilinha
- **Select** (`select.tsx`) - Seleção dropdown
- **Checkbox** (`checkbox.tsx`) - Caixa de seleção
- **Radio Group** (`radio-group.tsx`) - Grupo de botões de opção
- **Switch** (`switch.tsx`) - Interruptor toggle
- **Form** (`form.tsx`) - Componente de formulário integrado com React Hook Form

### Feedback e Notificações
- **Alert** (`alert.tsx`) - Alerta de mensagem
- **Alert Dialog** (`alert-dialog.tsx`) - Diálogo de confirmação
- **Toast** (`sonner.tsx`) - Notificações toast (usando Sonner)

### Overlays e Modais
- **Dialog** (`dialog.tsx`) - Diálogo modal
- **Sheet** (`sheet.tsx`) - Painel lateral deslizante
- **Popover** (`popover.tsx`) - Popover flutuante
- **Tooltip** (`tooltip.tsx`) - Dica de ferramenta
- **Dropdown Menu** (`dropdown-menu.tsx`) - Menu dropdown
- **Menubar** (`menubar.tsx`) - Barra de menu

### Navegação
- **Tabs** (`tabs.tsx`) - Abas de navegação
- **Navigation Menu** (`navigation-menu.tsx`) - Menu de navegação
- **Accordion** (`accordion.tsx`) - Acordeão expansível

### Dados e Conteúdo
- **Table** (`table.tsx`) - Tabela de dados
- **Card** (`card.tsx`) - Card de conteúdo
- **Badge** (`badge.tsx`) - Badge/etiqueta
- **Avatar** (`avatar.tsx`) - Avatar de usuário
- **Separator** (`separator.tsx`) - Separador visual
- **Skeleton** (`skeleton.tsx`) - Placeholder de carregamento

### Interação
- **Button** (`button.tsx`) - Botão
- **Slider** (`slider.tsx`) - Controle deslizante
- **Progress** (`progress.tsx`) - Barra de progresso
- **Scroll Area** (`scroll-area.tsx`) - Área com scroll customizado
- **Command** (`command.tsx`) - Comando/paleta de comandos
- **Calendar** (`calendar.tsx`) - Calendário

## Uso

### Importação

```typescript
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
```

### Exemplo com React Hook Form

```typescript
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { createUserSchema } from '@/schemas/userSchema'

const MyForm = () => {
  const form = useForm({
    resolver: zodResolver(createUserSchema),
  })

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)}>
        <FormField
          control={form.control}
          name="email"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Email</FormLabel>
              <FormControl>
                <Input {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button type="submit">Enviar</Button>
      </form>
    </Form>
  )
}
```

### Exemplo com Toast (Sonner)

```typescript
import { toast } from 'sonner'

// No seu componente
const handleSuccess = () => {
  toast.success('Operação realizada com sucesso!')
}

const handleError = () => {
  toast.error('Ocorreu um erro')
}
```

**Importante**: Adicione o `<Toaster />` no seu layout:

```typescript
import { Toaster } from '@/components/ui/sonner'

export default function RootLayout({ children }) {
  return (
    <html>
      <body>
        {children}
        <Toaster />
      </body>
    </html>
  )
}
```

## Documentação

Para mais informações sobre cada componente, consulte a [documentação oficial do shadcn/ui](https://ui.shadcn.com/docs/components).
