import { createClient } from '@supabase/supabase-js'

// Configuración de Supabase
// Nota: El anon key es público y seguro de exponer en el cliente
const supabaseUrl = process.env.NEXT_PUBLIC_SUPABASE_URL || 'https://rbpdhapfcljecofrscnj.supabase.co'
const supabaseAnonKey = process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY || 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InJicGRoYXBmY2xqZWNvZnJzY25qIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NjQ5NjY3NDMsImV4cCI6MjA4MDU0Mjc0M30.Ba-W2KHJS8U6OYVAjU98Y7JDn87gYPuhFvg_0vhcFfI'

if (!supabaseUrl || !supabaseAnonKey) {
  console.error('⚠️ Variables de entorno de Supabase no configuradas')
}

// Crear cliente de Supabase
export const supabase = createClient(supabaseUrl, supabaseAnonKey, {
  auth: {
    autoRefreshToken: true,
    persistSession: true,
    detectSessionInUrl: true,
  },
})

