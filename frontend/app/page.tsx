'use client'

import { DocumentList } from '@/components/documents/DocumentList'
import { DocumentUpload } from '@/components/documents/DocumentUpload'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'

export default function Home() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Dashboard</h1>
        <p className="text-muted-foreground">
          Gerencie seus documentos PDF
        </p>
      </div>

      <Tabs defaultValue="documents" className="space-y-4">
        <TabsList>
          <TabsTrigger value="documents">Documentos</TabsTrigger>
          <TabsTrigger value="upload">Enviar Documento</TabsTrigger>
        </TabsList>
        <TabsContent value="documents" className="space-y-4">
          <DocumentList />
        </TabsContent>
        <TabsContent value="upload" className="space-y-4">
          <DocumentUpload />
        </TabsContent>
      </Tabs>
    </div>
  )
}
