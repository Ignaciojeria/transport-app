# Configuración de Supabase para MenuID

Este documento explica cómo configurar la tabla necesaria en Supabase para almacenar la relación entre usuarios y sus menuIDs.

## Tabla `user_menus`

Necesitas crear una tabla en Supabase para almacenar la relación entre usuarios y sus menuIDs.

### SQL para crear la tabla

Ejecuta este SQL en el SQL Editor de Supabase:

```sql
-- Crear la tabla user_menus
CREATE TABLE IF NOT EXISTS user_menus (
  user_id UUID PRIMARY KEY REFERENCES auth.users(id) ON DELETE CASCADE,
  menu_id TEXT NOT NULL UNIQUE,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

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

-- Política: Los usuarios solo pueden ver y modificar su propio menuID
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

## Notas

- El `menu_id` es un UUID v7 (36 caracteres con guiones) generado en el frontend
- UUID v7 es ordenable cronológicamente y único globalmente
- La tabla usa RLS (Row Level Security) para asegurar que los usuarios solo puedan acceder a su propio menuID
- El `user_id` hace referencia a `auth.users(id)` de Supabase, por lo que se eliminará automáticamente si se elimina el usuario
- El campo `updated_at` se actualiza automáticamente cuando se modifica un registro

