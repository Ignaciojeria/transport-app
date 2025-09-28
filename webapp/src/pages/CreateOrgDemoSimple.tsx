import { CreateOrganizationSimple } from "@/components/CreateOrganizationSimple"

export default function CreateOrgDemoSimple() {
  return (
    <div>
      <CreateOrganizationSimple userEmail="demo@example.com" />
      
      {/* Botón para volver */}
      <div className="fixed top-4 left-4 z-50">
        <a 
          href="/"
          className="inline-flex items-center px-4 py-2 bg-black/50 hover:bg-black/70 text-white rounded-lg backdrop-blur-sm border border-white/20 transition-all duration-200"
        >
          ← Volver al inicio
        </a>
      </div>
    </div>
  )
}
