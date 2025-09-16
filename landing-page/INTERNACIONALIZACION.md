# Sistema de Internacionalización - TransportApp Landing Page

## Descripción

Se ha implementado un sistema completo de internacionalización (i18n) que permite mostrar la landing page en tres idiomas:

- **CL (Chile)**: Español 🇨🇱
- **BR (Brasil)**: Portugués 🇧🇷  
- **EU (Europa)**: Inglés 🇪🇺

## Características

### 1. Query Parameters
- Soporte para `?lang=CL`, `?lang=BR`, `?lang=EU`
- Ejemplo: `https://tu-dominio.com/?lang=BR`

### 2. Persistencia
- El idioma seleccionado se guarda en `localStorage`
- Al recargar la página, mantiene el idioma seleccionado

### 3. Detección Automática
- Si no hay query param ni idioma guardado, detecta automáticamente el idioma del navegador
- Portugués para `pt-*`, Inglés para `en-*`, Español por defecto

### 4. Selector de Idioma
- Dropdown en la navegación con banderas y nombres de idiomas
- Cambio instantáneo de idioma
- Actualiza la URL automáticamente

## Archivos Creados/Modificados

### Nuevos Archivos
- `lib/translations.ts` - Definiciones de tipos y traducciones
- `lib/useLanguage.ts` - Hook personalizado para manejo de idioma
- `components/LanguageSelector.tsx` - Componente selector de idioma

### Archivos Modificados
- `app/page.tsx` - Landing page principal con traducciones dinámicas

## Uso

### Para Desarrolladores

1. **Agregar nuevas traducciones**:
   ```typescript
   // En lib/translations.ts
   export interface Translations {
     // ... propiedades existentes
     nuevaSeccion: {
       titulo: string
       descripcion: string
     }
   }
   
   // Agregar traducciones para cada idioma
   export const translations: Record<Language, Translations> = {
     CL: {
       // ... traducciones existentes
       nuevaSeccion: {
         titulo: "Nuevo Título",
         descripcion: "Nueva descripción"
       }
     },
     // ... otros idiomas
   }
   ```

2. **Usar traducciones en componentes**:
   ```tsx
   import { useLanguage } from "@/lib/useLanguage"
   
   function MiComponente() {
     const { t, language } = useLanguage()
     
     return (
       <div>
         <h1>{t.nuevaSeccion.titulo}</h1>
         <p>{t.nuevaSeccion.descripcion}</p>
       </div>
     )
   }
   ```

### Para Usuarios

1. **Cambiar idioma**:
   - Usar el selector en la navegación (globo 🌐)
   - Agregar `?lang=XX` a la URL

2. **URLs de ejemplo**:
   - Español: `https://tu-dominio.com/?lang=CL`
   - Portugués: `https://tu-dominio.com/?lang=BR`
   - Inglés: `https://tu-dominio.com/?lang=EU`

## Estructura de Traducciones

Las traducciones están organizadas por secciones:

- `nav` - Navegación
- `hero` - Sección principal
- `howItWorks` - Cómo funciona
- `benefits` - Beneficios
- `cta` - Call to Action
- `footer` - Pie de página

Cada sección contiene las traducciones específicas para todos los elementos de esa sección.

## Consideraciones Técnicas

- **SSR/SSG**: El hook `useLanguage` maneja la hidratación correctamente
- **Performance**: Las traducciones se cargan una sola vez
- **SEO**: Los query parameters permiten URLs específicas por idioma
- **Accesibilidad**: El selector de idioma es accesible por teclado

## Próximos Pasos

1. Agregar más idiomas si es necesario
2. Implementar traducciones para el componente `DemoEmbed`
3. Agregar meta tags específicos por idioma
4. Implementar rutas específicas por idioma (`/es`, `/pt`, `/en`)
