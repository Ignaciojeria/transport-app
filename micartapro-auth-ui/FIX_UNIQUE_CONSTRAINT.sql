-- ============================================
-- SCRIPT PARA AGREGAR RESTRICCIÓN ÚNICA
-- ============================================
-- Este SQL asegura que solo haya UN registro por usuario
-- Ejecuta esto en el SQL Editor de Supabase
-- ============================================

-- 1. Verificar si ya existe una restricción única en user_id
-- Si la tabla tiene user_id como PRIMARY KEY, ya está protegida
-- Pero si tiene un campo 'id' separado, necesitamos agregar una restricción única

-- Opción A: Si user_id NO es PRIMARY KEY, agregar restricción única
DO $$
BEGIN
  -- Verificar si ya existe una restricción única en user_id
  IF NOT EXISTS (
    SELECT 1 
    FROM pg_constraint 
    WHERE conname = 'user_menus_user_id_key' 
    OR (conrelid = 'user_menus'::regclass AND contype = 'u' AND conkey::int[] = ARRAY[(SELECT attnum FROM pg_attribute WHERE attrelid = 'user_menus'::regclass AND attname = 'user_id')])
  ) THEN
    -- Agregar restricción única
    ALTER TABLE user_menus 
    ADD CONSTRAINT user_menus_user_id_key UNIQUE (user_id);
    
    RAISE NOTICE 'Restricción única agregada a user_id';
  ELSE
    RAISE NOTICE 'La restricción única en user_id ya existe';
  END IF;
END $$;

-- 2. Limpiar registros duplicados (mantener solo el más reciente)
-- IMPORTANTE: Esto eliminará registros duplicados, manteniendo solo el más reciente por usuario
WITH ranked_menus AS (
  SELECT 
    id,
    user_id,
    menu_id,
    created_at,
    ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY created_at DESC) as rn
  FROM user_menus
)
DELETE FROM user_menus
WHERE id IN (
  SELECT id FROM ranked_menus WHERE rn > 1
);

-- 3. Verificar que la restricción está activa
SELECT 
  conname as constraint_name,
  contype as constraint_type,
  pg_get_constraintdef(oid) as constraint_definition
FROM pg_constraint
WHERE conrelid = 'user_menus'::regclass
AND contype = 'u';

