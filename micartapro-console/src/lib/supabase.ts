import { createClient } from '@supabase/supabase-js'
import type { SupabaseClient } from '@supabase/supabase-js'

// Configuración de Supabase
const supabaseUrl = 'https://rbpdhapfcljecofrscnj.supabase.co'
const supabaseAnonKey = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InJicGRoYXBmY2xqZWNvZnJzY25qIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NjQ5NjY3NDMsImV4cCI6MjA4MDU0Mjc0M30.Ba-W2KHJS8U6OYVAjU98Y7JDn87gYPuhFvg_0vhcFfI'

// Crear cliente de Supabase
export const supabase: SupabaseClient = createClient(supabaseUrl, supabaseAnonKey, {
  auth: {
    autoRefreshToken: true,
    persistSession: true,
    detectSessionInUrl: true,
  },
})

