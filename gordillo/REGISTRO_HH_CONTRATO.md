# Registro HH - Contrato Refactorizado

## Nuevo Contrato para Múltiples Trabajadores

El contrato ha sido refactorizado para permitir el registro de horas de múltiples trabajadores en una sola operación, lo que es más eficiente y práctico.

### Estructura del Contrato:

```json
{
  "fechaEjecucion": "yyyy-mm-dd",
  "orderDeTrabajo": "OT-001",
  "actividad": "trabajo de mantenimiento",
  "trabajadores": [
    {
      "trabajador": "alexander gutierrez",
      "horasNormales": 8,
      "horasExtras": 2
    },
    {
      "trabajador": "maria rodriguez", 
      "horasNormales": 6,
      "horasExtras": 1
    }
  ]
}
```

### Campos del Contrato:

#### Información General:
- **`fechaEjecucion`**: Fecha de ejecución del trabajo (requerido)
- **`orderDeTrabajo`**: Orden de trabajo asociada (opcional)
- **`actividad`**: Tipo de actividad realizada (requerido)

#### Array de Trabajadores:
- **`trabajador`**: Nombre del trabajador (requerido)
- **`horasNormales`**: Horas normales trabajadas (número)
- **`horasExtras`**: Horas extras trabajadas (número)

### Ventajas del Nuevo Contrato:

1. **Eficiencia**: Permite registrar múltiples trabajadores en una sola operación
2. **Consistencia**: Todos los trabajadores comparten la misma fecha, orden de trabajo y actividad
3. **Flexibilidad**: Se puede agregar o quitar trabajadores dinámicamente
4. **Totales Automáticos**: Calcula automáticamente los totales por trabajador y generales

### Funcionalidades del Formulario:

#### Gestión de Trabajadores:
- **Agregar Trabajador**: Botón verde para agregar nuevos trabajadores
- **Eliminar Trabajador**: Botón rojo para quitar trabajadores (mínimo 1)
- **Botón "+"**: Para agregar horas extras rápidamente

#### Validaciones:
- Fecha de ejecución requerida
- Actividad requerida
- Al menos un trabajador con horas > 0
- Trabajador seleccionado para cada entrada

#### Visualización:
- **Tarjetas separadas**: Cada registro se muestra en su propia tarjeta
- **Información general**: Fecha, orden de trabajo, actividad y fecha de creación
- **Tabla de trabajadores**: Con horas normales, extras y totales
- **Totales**: Suma automática de todas las horas por tipo

### Ejemplo de Uso:

1. **Seleccionar fecha** de ejecución
2. **Elegir orden de trabajo** (opcional)
3. **Seleccionar actividad**
4. **Agregar trabajadores**:
   - Seleccionar trabajador del dropdown
   - Ingresar horas normales
   - Ingresar horas extras (o usar botón "+")
5. **Crear registro** para todos los trabajadores

### Beneficios Operacionales:

- **Reducción de tiempo**: Un solo formulario para múltiples trabajadores
- **Menos errores**: Información compartida evita duplicación
- **Mejor organización**: Registros agrupados por actividad/fecha
- **Reportes más claros**: Totales automáticos por registro
