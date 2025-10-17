# Editor de Markdown para Prompts

## 🎯 Funcionalidad Implementada

Se ha implementado un **editor de Markdown por defecto** para todos los prompts de productos. Ahora todas las preguntas se escriben en Markdown, lo que permite crear contenido más rico y estructurado.

## ✨ Características

### **Editor de Markdown por Defecto**
- **Siempre activo**: No hay toggle, todas las preguntas son Markdown
- **Toolbar completa**: Bold, italic, headers, listas, links, etc.
- **Sintaxis highlighting**: Colores para mejor legibilidad
- **Altura optimizada**: 200px por defecto
- **Modo claro**: Optimizado para interfaces claras

### **Vista Previa Automática**
- **Preview en tiempo real**: Ve cómo se verá el Markdown
- **Estilos personalizados**: Diseño consistente con la app
- **Responsive**: Se adapta al contenido
- **Solo cuando hay contenido**: Se muestra automáticamente

### **Ayuda Integrada**
- **Ejemplos de sintaxis**: Guía rápida de Markdown
- **Casos de uso comunes**: Negrita, cursiva, listas, etc.
- **Siempre visible**: Ayuda contextual permanente

## 🚀 Casos de Uso

### **Pregunta Simple (Texto)**
```
¿Qué tamaño prefieres?
```

### **Pregunta Rica (Markdown)**
```markdown
## 🎨 Personalización de Color

¿Qué color prefieres para tu producto?

### Opciones disponibles:
- **Rojo**: Clásico y elegante
- **Azul**: Profesional y confiable  
- **Verde**: Natural y fresco
- **Negro**: Sofisticado y moderno

> 💡 *Puedes combinar hasta 2 colores si lo deseas*

### Información adicional:
- Los colores están disponibles en todos los tamaños
- Tiempo de producción: 2-3 días hábiles
- Garantía de color: 1 año
```

## 📋 Tipos de Preguntas Soportadas

1. **Texto libre** - Con soporte completo de Markdown
2. **Selección única** - Dropdown con opciones
3. **Selección múltiple** - Checkboxes con opciones
4. **Número** - Input numérico
5. **Sí/No** - Checkbox simple

## 🎨 Estilos Personalizados

### **Editor**
- Fuente monospace para código
- Bordes redondeados
- Colores consistentes con la app

### **Preview**
- Tipografía optimizada
- Espaciado adecuado
- Colores de acento (azul para links, etc.)

### **Ayuda**
- Fondo azul claro
- Ejemplos interactivos
- Iconos emoji para mejor UX

## 🔧 Configuración Técnica

### **Dependencias**
- `@uiw/react-md-editor`: Editor principal
- `react-markdown`: Renderizado del preview

### **Estructura de Datos Completamente Markdown**
```typescript
interface PromptItem {
  questionMarkdown: string      // Pregunta en Markdown
  type: 'text' | 'select' | 'multiselect' | 'number' | 'boolean'
  options?: string[]
  required?: boolean
  placeholderMarkdown?: string  // Placeholder también en Markdown
}
```

## 💡 Beneficios

### **Para Modelos de IA**
- ✅ **Mejor comprensión**: Los LLMs entienden perfectamente Markdown
- ✅ **Estructura clara**: Headers, listas, énfasis mejoran la interpretación
- ✅ **Contexto rico**: Instrucciones más detalladas y claras

### **Para Usuarios**
- ✅ **Flexibilidad**: Preguntas más ricas y descriptivas
- ✅ **Formato profesional**: Mejor presentación visual
- ✅ **Instrucciones claras**: Ejemplos y detalles integrados

### **Para Desarrolladores**
- ✅ **Código más simple**: Un solo campo en lugar de dos
- ✅ **Menos complejidad**: No hay toggle ni lógica condicional
- ✅ **Más mantenible**: Estructura más clara y directa
- ✅ **Consistente**: Diseño unificado
- ✅ **Completamente Markdown**: Tanto pregunta como placeholder

### **Ejemplos de Placeholders**

**Placeholder Simple:**
```json
{
  "placeholderMarkdown": "Describe tus preferencias..."
}
```

**Placeholder Rico:**
```json
{
  "placeholderMarkdown": "## Instrucciones\n\nDescribe tus preferencias:\n- **Color**: ¿Cuál prefieres?\n- **Tamaño**: ¿Pequeño o grande?\n\n> 💡 *Sé específico para mejor atención*"
}
```

## 🎯 Ejemplos Prácticos

### **Producto: Palta Hass**
```markdown
## 🥑 Palta Hass Premium

¿Qué tamaño prefieres?

### Opciones disponibles:
- **Pequeño** (200-300g): Ideal para 1-2 personas
- **Mediano** (300-400g): Perfecto para familias
- **Grande** (400-500g): Para eventos especiales

> 🌟 *Todas nuestras paltas son orgánicas y certificadas*

### Información adicional:
- **Madurez**: ¿Prefieres que esté lista para comer o verde?
- **Entrega**: Disponible en 24-48 horas
- **Garantía**: Si no estás satisfecho, te devolvemos el dinero
```

### **Producto: Servicio de Transporte**
```markdown
## 🚚 Servicio de Transporte Express

¿Qué tipo de servicio necesitas?

### Opciones:
- **🏠 Residencial**: Mudanzas y entregas a domicilio
- **🏢 Comercial**: Transporte de mercancías
- **📦 Paquetería**: Envíos pequeños y medianos

> ⚡ *Servicio disponible 24/7*

### Información importante:
- **Seguro incluido**: Todos los envíos están asegurados
- **Tracking**: Seguimiento en tiempo real
- **Soporte**: Atención al cliente 24/7
```

¡El editor de Markdown hace que los prompts sean mucho más profesionales y efectivos! 🚀
