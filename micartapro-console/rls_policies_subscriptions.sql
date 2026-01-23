-- ============================================
-- POLÍTICAS RLS PARA LA TABLA subscriptions
-- ============================================
-- Este script crea las políticas RLS necesarias para la tabla subscriptions
-- Ejecuta este SQL en el SQL Editor de Supabase
-- ============================================

-- Habilitar Row Level Security (RLS) si no está habilitado
ALTER TABLE public.subscriptions ENABLE ROW LEVEL SECURITY;

-- ============================================
-- ELIMINAR POLÍTICAS EXISTENTES (si existen)
-- ============================================
-- Esto evita conflictos si ya tienes algunas políticas creadas

DROP POLICY IF EXISTS "Users can view their own subscriptions" ON public.subscriptions;
DROP POLICY IF EXISTS "Users can insert their own subscriptions" ON public.subscriptions;
DROP POLICY IF EXISTS "Users can update their own subscriptions" ON public.subscriptions;

-- ============================================
-- CREAR POLÍTICAS RLS
-- ============================================

-- SELECT: Los usuarios pueden ver sus propias suscripciones
-- Propósito: Permite que los usuarios lean su propia suscripción
-- Uso: Necesario para verificar si el usuario tiene una suscripción activa
CREATE POLICY "Users can view their own subscriptions"
  ON public.subscriptions FOR SELECT
  USING (auth.uid() = user_id);

-- INSERT: Los usuarios pueden insertar sus propias suscripciones
-- Propósito: Permite crear nuevos registros de suscripciones
-- Nota: Normalmente esto se hace desde el backend, pero por si acaso
CREATE POLICY "Users can insert their own subscriptions"
  ON public.subscriptions FOR INSERT
  WITH CHECK (auth.uid() = user_id);

-- UPDATE: Los usuarios pueden actualizar sus propias suscripciones
-- Propósito: Permite actualizar la suscripción de un usuario
-- Nota: Normalmente esto se hace desde el backend, pero por si acaso
CREATE POLICY "Users can update their own subscriptions"
  ON public.subscriptions FOR UPDATE
  USING (auth.uid() = user_id)
  WITH CHECK (auth.uid() = user_id);

-- ============================================
-- VERIFICACIÓN
-- ============================================
-- Ejecuta esta consulta para verificar que las políticas se crearon correctamente

SELECT 
  schemaname,
  tablename,
  policyname,
  permissive,
  roles,
  cmd as operation,  -- SELECT, UPDATE, INSERT, DELETE
  qual as using_clause,
  with_check as with_check_clause
FROM pg_policies 
WHERE tablename = 'subscriptions'
ORDER BY policyname;
