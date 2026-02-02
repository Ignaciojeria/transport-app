# Templates de Cadorago

Esta carpeta contiene los diseños y recursos para diferentes templates de la UI de Cadorago.

## Uso de Templates

Puedes seleccionar diferentes templates usando el query parameter `template` en la URL:

- **Template original (Hero)**: `/?template=hero` o sin parámetro (por defecto)
- **Template moderno**: `/?template=modern`

## Estructura

```
template/
├── README.md          # Este archivo
├── designs/          # Diseños, mockups, imágenes de referencia
└── assets/           # Recursos específicos del template (si aplica)
```

## Agregar un nuevo Template

1. Crea un nuevo componente en `src/components/templates/` (ej: `MiTemplate.svelte`)
2. El componente debe aceptar un slot para el contenido
3. Actualiza `Home.svelte` para incluir el nuevo template en el switch
4. Agrega el diseño de referencia en esta carpeta `template/designs/`

### Ejemplo de Template

```svelte
<script>
  export let bgColor = 'bg-white';
  export let className = '';
</script>

<div class={`min-h-screen ${bgColor} ${className}`}>
  <!-- Tu diseño aquí -->
  <slot />
</div>
```

## Templates Disponibles

- **hero** (por defecto): Template original con BrandHero y diseño clásico
- **modern**: Template moderno con gradientes y espaciado mejorado
