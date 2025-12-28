# Configuración de Supabase para MenuID

Este documento explica cómo configurar la tabla necesaria en Supabase para almacenar la relación entre usuarios y sus menuIDs.

## Tabla `user_menus`

Necesitas crear una tabla en Supabase para almacenar la relación entre usuarios y sus menuIDs.

### SQL para crear la tabla

Ejecuta este SQL en el SQL Editor de Supabase:

```sql
-- Crear la tabla user_menus
-- IMPORTANTE: user_id es PRIMARY KEY, lo que garantiza que solo haya UN registro por usuario
CREATE TABLE IF NOT EXISTS user_menus (
  user_id UUID PRIMARY KEY REFERENCES auth.users(id) ON DELETE CASCADE,
  menu_id TEXT NOT NULL UNIQUE,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Si la tabla ya existe y tiene un campo 'id' separado, agregar restricción única en user_id
-- (Esto es solo por si acaso, normalmente user_id ya es PRIMARY KEY)
DO $$
BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns 
             WHERE table_name = 'user_menus' 
             AND column_name = 'id' 
             AND column_name != 'user_id') THEN
    -- Si tiene campo 'id' separado, agregar restricción única en user_id
    IF NOT EXISTS (
      SELECT 1 FROM pg_constraint 
      WHERE conrelid = 'user_menus'::regclass 
      AND contype = 'u' 
      AND conkey::int[] = ARRAY[(SELECT attnum FROM pg_attribute WHERE attrelid = 'user_menus'::regclass AND attname = 'user_id')]
    ) THEN
      ALTER TABLE user_menus ADD CONSTRAINT user_menus_user_id_key UNIQUE (user_id);
    END IF;
  END IF;
END $$;

-- Crear índice para búsquedas rápidas por menu_id
CREATE INDEX IF NOT EXISTS idx_user_menus_menu_id ON user_menus(menu_id);

-- Crear función para actualizar updated_at automáticamente
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Crear trigger para actualizar updated_at
CREATE TRIGGER update_user_menus_updated_at 
  BEFORE UPDATE ON user_menus 
  FOR EACH ROW 
  EXECUTE FUNCTION update_updated_at_column();

-- Habilitar Row Level Security (RLS)
ALTER TABLE user_menus ENABLE ROW LEVEL SECURITY;

-- Eliminar políticas existentes si existen (para evitar conflictos)
DROP POLICY IF EXISTS "Users can view their own menu" ON user_menus;
DROP POLICY IF EXISTS "Users can insert their own menu" ON user_menus;
DROP POLICY IF EXISTS "Users can update their own menu" ON user_menus;

-- Política: Los usuarios solo pueden ver su propio menuID
CREATE POLICY "Users can view their own menu"
  ON user_menus FOR SELECT
  USING (auth.uid() = user_id);

-- Política: Los usuarios solo pueden insertar su propio menuID
-- IMPORTANTE: Usa USING y WITH CHECK para INSERT
CREATE POLICY "Users can insert their own menu"
  ON user_menus FOR INSERT
  WITH CHECK (auth.uid() = user_id);

-- Política: Los usuarios solo pueden actualizar su propio menuID
CREATE POLICY "Users can update their own menu"
  ON user_menus FOR UPDATE
  USING (auth.uid() = user_id)
  WITH CHECK (auth.uid() = user_id);
```

## Verificación

Después de ejecutar el SQL, puedes verificar que la tabla se creó correctamente:

```sql
-- Verificar que la tabla existe
SELECT * FROM information_schema.tables 
WHERE table_name = 'user_menus';

-- Verificar las políticas RLS
SELECT * FROM pg_policies 
WHERE tablename = 'user_menus';
```

## Diagnóstico de Problemas

Si el menuID no se está creando automáticamente, puedes usar la función de diagnóstico:

1. Abre la consola del navegador (F12)
2. En la página de autenticación, ejecuta:
   ```javascript
   await diagnoseMenuIdIssue()
   ```

Esta función verificará:
- ✅ Si hay una sesión activa
- ✅ Si la tabla existe
- ✅ Si tienes permisos de lectura (RLS)
- ✅ Si tienes permisos de escritura (RLS)

### Problemas Comunes

**Error: "La tabla user_menus no existe"**
- Solución: Ejecuta el SQL completo en el SQL Editor de Supabase

**Error: "Permiso denegado por RLS"**
- Solución: Verifica que las políticas RLS estén creadas correctamente
- Ve a Supabase Dashboard > Authentication > Policies
- Asegúrate de que las 3 políticas estén activas:
  - "Users can view their own menu" (SELECT)
  - "Users can insert their own menu" (INSERT)
  - "Users can update their own menu" (UPDATE)

**Error: "RLS bloquea la inserción"**
- Verifica que `auth.uid()` esté funcionando correctamente
- Asegúrate de que la sesión esté completamente establecida antes de crear el menuID

## Notas

- El `menu_id` es un UUID v7 (36 caracteres con guiones) generado en el frontend
- UUID v7 es ordenable cronológicamente y único globalmente
- La tabla usa RLS (Row Level Security) para asegurar que los usuarios solo puedan acceder a su propio menuID
- El `user_id` hace referencia a `auth.users(id)` de Supabase, por lo que se eliminará automáticamente si se elimina el usuario
- El campo `updated_at` se actualiza automáticamente cuando se modifica un registro

