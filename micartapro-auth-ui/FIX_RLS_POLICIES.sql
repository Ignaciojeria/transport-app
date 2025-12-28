-- ============================================
-- SCRIPT PARA CORREGIR POLÍTICAS RLS
-- ============================================
-- Si estás teniendo problemas con "Permiso denegado por RLS",
-- ejecuta este SQL en el SQL Editor de Supabase
-- ============================================

-- 1. Eliminar todas las políticas existentes
DROP POLICY IF EXISTS "Users can view their own menu" ON user_menus;
DROP POLICY IF EXISTS "Users can insert their own menu" ON user_menus;
DROP POLICY IF EXISTS "Users can update their own menu" ON user_menus;

-- 2. Verificar que RLS esté habilitado
ALTER TABLE user_menus ENABLE ROW LEVEL SECURITY;

-- 3. Crear las políticas nuevamente con sintaxis correcta
CREATE POLICY "Users can view their own menu"
  ON user_menus FOR SELECT
  USING (auth.uid() = user_id);

CREATE POLICY "Users can insert their own menu"
  ON user_menus FOR INSERT
  WITH CHECK (auth.uid() = user_id);

CREATE POLICY "Users can update their own menu"
  ON user_menus FOR UPDATE
  USING (auth.uid() = user_id)
  WITH CHECK (auth.uid() = user_id);

-- 4. Verificar que las políticas se crearon correctamente
SELECT 
  schemaname,
  tablename,
  policyname,
  permissive,
  roles,
  cmd,
  qual,
  with_check
FROM pg_policies 
WHERE tablename = 'user_menus';

-- Si después de ejecutar esto sigues teniendo problemas,
-- prueba esta versión alternativa que usa auth.jwt():

-- ALTERNATIVA (solo si la anterior no funciona):
-- DROP POLICY IF EXISTS "Users can insert their own menu" ON user_menus;
-- CREATE POLICY "Users can insert their own menu"
--   ON user_menus FOR INSERT
--   WITH CHECK ((auth.jwt() ->> 'sub')::uuid = user_id);

