## transport-app

[![codecov](https://codecov.io/gh/Ignaciojeria/transport-app/branch/main/graph/badge.svg)](https://codecov.io/gh/Ignaciojeria/transport-app) [![Go Report Card](https://goreportcard.com/badge/github.com/Ignaciojeria/transport-app)](https://goreportcard.com/report/github.com/Ignaciojeria/transport-app)


# 🔄 Propagación de Contexto en Transport App

En esta aplicación utilizamos OpenTelemetry Baggage para propagar contexto entre componentes (API → eventos Pub/Sub → consumidores) de forma trazable, desacoplada y estandarizada.

---

## 🧾 Headers del API

El consumidor de la API debe enviar los siguientes headers en cada request:

| Header         | Ejemplo              | Descripción                            |
|----------------|----------------------|----------------------------------------|
| `organization` | `2-CL`               | ID del tenant + país (`tenantId-country`) |
| `consumer`     | `CONSUMER_EXAMPLE`   | Identificador del consumidor del API   |
| `commerce`     | `COMMERCE_EXAMPLE`   | Identificador del comercio             |

Estos headers son procesados por un middleware que:

- Desglosa el `organization` en `tenant.id` y `tenant.country`.
- Carga todos los valores como `baggage` en el contexto de OpenTelemetry.

---

## 🧠 OpenTelemetry Baggage Keys

Estas claves viven dentro del contexto propagado (`context.Context`) y son visibles desde trazas (Jaeger, OpenObserve, etc).

| Clave Baggage       | Ejemplo              | Fuente               |
|---------------------|----------------------|----------------------|
| `tenant.id`         | `2`                  | Derivado de `organization` |
| `tenant.country`    | `CL`                 | Derivado de `organization` |
| `business.consumer` | `CONSUMER_EXAMPLE`   | Header `consumer`    |
| `business.commerce` | `COMMERCE_EXAMPLE`   | Header `commerce`    |
| `event.type`        | `orderSubmitted`     | Inyectado por código |
| `entity.type`       | `order`              | Inyectado por código |

---

## 📦 Pub/Sub Message Attributes

Cuando se publica un evento (como `orderSubmitted`), los atributos del contexto son exportados al mensaje en formato `camelCase`:

```json
{
  "tenantId": "2",
  "tenantCountry": "CL",
  "consumer": "CONSUMER_EXAMPLE",
  "commerce": "COMMERCE_EXAMPLE",
  "eventType": "orderSubmitted",
  "entityType": "order"
}
```


## 📌 Convención de Campos para Direcciones en Chile

| Campo en JSON                | Descripción                      | Ejemplo                       |
| ---------------------------- | --------------------------------- | ----------------------------- |
| `addressInfo.addressLine1` | Calle y número                    | `inglaterra 59`                    |
| `addressInfo.addressLine2` | Información adicional (opcional)  | `dpto 2214`                        |
| `addressInfo.district`     | Comuna                            | `la florida`                       |
| `addressInfo.province`     | Provincia                         | `santiago`                         |
| `addressInfo.state`        | Región                            | `región metropolitana de santiago` |
| `addressInfo.latitude`     | Latitud                           | `-33.5204181`                      |
| `addressInfo.longitude`    | Longitud                          | `-70.6006178`               |      
---

## 🌐 Ejemplos de Direcciones en Chile

### 🏢 Ejemplo en Santiago (con departamento)

```json
{
  "addressInfo": {
    "addressLine1": "av providencia 1234",
    "addressLine2": "depto 1202",
    "district": "providencia",
    "province": "santiago",
    "state": "región metropolitana de santiago",
    "latitude": -33.4263,
    "longitude": -70.6200
  }
}
```

### 🏙️ Ejemplo en Santiago (sin departamento)

```json
{
  "addressInfo": {
    "addressLine1": "avenida las condes 7890",
    "addressLine2": "",
    "district": "las condes",
    "province": "santiago",
    "state": "región metropolitana de santiago",
    "latitude": -33.4085,
    "longitude": -70.5666
  }
}
```

## 📌 Convención de Campos para Direcciones en México

| Campo en JSON                | Descripción                      | Ejemplo                       |
| ---------------------------- | --------------------------------- | ----------------------------- |
| `addressInfo.addressLine1` | Calle y número                    | `av insurgentes 1234`              |
| `addressInfo.addressLine2` | Información adicional (opcional)  | `interior 5` o `depto 202`         |
| `addressInfo.district`     | Colonia                           | `condesa`                          |
| `addressInfo.province`     | Alcaldía/Municipio                | `cuauhtémoc`                       |
| `addressInfo.state`        | Estado                            | `ciudad de méxico`                 |
| `addressInfo.latitude`     | Latitud                           | `19.4326`                          |
| `addressInfo.longitude`    | Longitud                          | `-99.1332`                         |
---
> ℹ️ **Nota:** El campo `province` representa:
> - Una **alcaldía** si el estado es `ciudad de méxico`.
> - Un **municipio** en el resto del país.

# 📦 Sistema de Gestión de Paquetes en Órdenes de Transporte

## Estructura de Paquetes

La API de órdenes de transporte permite definir paquetes de dos maneras diferentes:

1. **Paquetes con LPN (License Plate Number)**: Utilizados cuando existe una identificación única para la unidad física de entrega.
2. **Paquetes sin LPN**: Utilizados cuando los items deben ser tratados como entidades individuales o cuando no se dispone de una identificación de embalaje.

### Ejemplos de Estructuras

#### Paquete con LPN Especificado

```json
{
    "dimensions": {"length": 10, "height": 10, "unit": "cm", "width": 10},
    "insurance": {"currency": "USD", "unitValue": 10},
    "items": [
        {
            "sku": "ITEM-001",
            "description": "Producto A",
            "dimensions": {"length": 1, "height": 1, "unit": "cm", "width": 1},
            "insurance": {"currency": "USD", "unitValue": 100},
            "logisticCondition": "normal",
            "quantity": {"quantityNumber": 1, "quantityUnit": "unit"},
            "weight": {"unit": "kg", "value": 0.5}
        }
    ],
    "lpn": "PKG-12345",
    "weight": {"unit": "kg", "value": 1}
}
```

#### Paquete sin LPN

```json
{
    "items": [
        {
            "sku": "ITEM-001",
            "description": "Producto A",
            "dimensions": {"length": 1, "height": 1, "unit": "cm", "width": 1},
            "insurance": {"currency": "USD", "unitValue": 100},
            "logisticCondition": "normal",
            "quantity": {"quantityNumber": 1, "quantityUnit": "unit"},
            "weight": {"unit": "kg", "value": 0.5}
        }
    ]
}
```

## Comportamiento del Sistema

### Manejo de Paquetes con LPN

Cuando un paquete incluye un valor LPN, este valor se utiliza como identificador principal para el paquete. El sistema:

- Utiliza el LPN para generar un identificador único (DocID) para el paquete
- Permite buscar, actualizar y seguir el paquete específicamente por su LPN
- Considera el paquete como una unidad física única que contiene todos los items listados

### Manejo de Paquetes sin LPN

Cuando un paquete se envía sin LPN, el sistema:

1. **Generación de identificador**: Genera un identificador único basado en:
   - La referencia de la orden (`referenceID`)
   - Los SKUs de los items contenidos en el paquete

2. **Tratamiento lógico**: Considera todos los items como parte de la misma orden, donde:
   - La orden completa es la unidad principal de entrega
   - Las agrupaciones de items sin LPN representan conjuntos lógicos sin asumir empaquetado físico específico

3. **Trazabilidad**: Mantiene la trazabilidad a través del identificador generado, asegurando que:
   - Actualizaciones posteriores afecten al mismo conjunto de items
   - Se puedan realizar seguimientos coherentes aunque no exista un LPN explícito

4. **Flexibilidad**: Permite adaptarse a diferentes flujos logísticos sin imponer restricciones en el empaquetado físico.

## Consideraciones para el Uso

### Cuándo Usar Paquetes con LPN

- Cuando sus sistemas logísticos asignan identificadores únicos a los paquetes físicos
- Cuando necesita mantener consistencia con sistemas externos que utilizan LPNs
- Para seguimiento de alto nivel de unidades de envío específicas

### Cuándo Usar Paquetes sin LPN

- Cuando trabaja con proveedores que no utilizan sistemas de codificación de paquetes
- Cuando los items se transportan individualmente sin empaquetado conjunto
- Cuando solo necesita identificar los productos a entregar sin preocuparse por su agrupación física

## Ejemplos de Escenarios

### Escenario 1: Múltiples Productos en un Paquete

Un pedido contiene tres productos diferentes que se embalan juntos:

```json
{
    "packages": [
        {
            "lpn": "BOX-12345",
            "items": [
                {"sku": "PROD-A", "quantity": {"quantityNumber": 2, "quantityUnit": "unit"}},
                {"sku": "PROD-B", "quantity": {"quantityNumber": 1, "quantityUnit": "unit"}},
                {"sku": "PROD-C", "quantity": {"quantityNumber": 3, "quantityUnit": "unit"}}
            ]
        }
    ]
}
```

### Escenario 2: Productos Individuales sin Embalaje Común

Un pedido contiene dos productos que se entregan de forma independiente:

```json
{
    "packages": [
        {
            "items": [
                {"sku": "PROD-X", "quantity": {"quantityNumber": 1, "quantityUnit": "unit"}}
            ]
        },
        {
            "items": [
                {"sku": "PROD-Y", "quantity": {"quantityNumber": 1, "quantityUnit": "unit"}}
            ]
        }
    ]
}
```

## 🎯 Ejemplo de Equivalencia entre Agrupaciones de Items sin LPN

En el sistema, la ausencia de LPN implica que **no existe una unidad de embalaje físico identificable**. Esto permite representar ítems como agrupaciones lógicas, independientemente de cómo estén organizados en el cuerpo de la solicitud.

A continuación se muestran dos formas distintas de estructurar los mismos ítems en una orden, que **son tratadas de forma equivalente** por el sistema.

---

#### 📦 Escenario A: Agrupación múltiple sin LPN

```json
{
    "packages": [
        {
            "items": [
                {"sku": "PROD-A", "quantity": {"quantityNumber": 1, "quantityUnit": "unit"}},
                {"sku": "PROD-B", "quantity": {"quantityNumber": 1, "quantityUnit": "unit"}}
            ]
        },
        {
            "items": [
                {"sku": "PROD-C", "quantity": {"quantityNumber": 1, "quantityUnit": "unit"}}
            ]
        }
    ]
}
```

#### 📦 Escenario B: Paquetes individuales sin LPN

```json
{
    "packages": [
        {
            "items": [
                {"sku": "PROD-A", "quantity": {"quantityNumber": 1, "quantityUnit": "unit"}}
            ]
        },
        {
            "items": [
                {"sku": "PROD-B", "quantity": {"quantityNumber": 1, "quantityUnit": "unit"}}
            ]
        },
        {
            "items": [
                {"sku": "PROD-C", "quantity": {"quantityNumber": 1, "quantityUnit": "unit"}}
            ]
        }
    ]
}
```

### ✅ Resultado en el sistema

- Ambos escenarios representan los mismos ítems con las mismas cantidades.
- Ambos omiten el LPN, por lo que no existe una referencia a un paquete físico específico.
- Ambos son interpretados como una serie de ítems asociados directamente a la orden, no como unidades físicas independientes.
- El sistema genera identificadores únicos por agrupación, pero todos están ligados al mismo orderReferenceId.

### 🧠 Conclusión

- Aunque la estructura visual difiere, el tratamiento lógico y funcional dentro del sistema es el mismo.
- Lo relevante no es cómo se organizan los items en el payload, sino la ausencia de LPN, que indica que estas agrupaciones no representan una caja o contenedor físico formal.
- Esto otorga al sistema una gran flexibilidad operativa, permitiendo representar tanto productos individuales como agrupaciones lógicas sin imponer restricciones rígidas de empaquetado.


### Escenario 3: Mezcla de Paquetes con y sin LPN

Un pedido incluye tanto productos empaquetados como individuales:

```json
{
    "packages": [
        {
            "lpn": "BOX-A",
            "items": [
                {"sku": "PROD-1", "quantity": {"quantityNumber": 5, "quantityUnit": "unit"}},
                {"sku": "PROD-2", "quantity": {"quantityNumber": 2, "quantityUnit": "unit"}}
            ]
        },
        {
            "items": [
                {"sku": "PROD-3", "quantity": {"quantityNumber": 1, "quantityUnit": "unit"}}
            ]
        }
    ]
}
```

## Nota Técnica

El sistema automáticamente genera identificadores de documento (DocID) para los paquetes utilizando una función hash que garantiza la consistencia y unicidad de los identificadores. Para paquetes sin LPN, este identificador se basa en la combinación de la referencia de orden y los SKUs de los items, lo que asegura que:

1. El mismo conjunto de items en la misma orden siempre genere el mismo identificador
2. Diferentes conjuntos de items o diferentes órdenes generen identificadores distintos
3. El sistema pueda rastrear y actualizar correctamente los paquetes incluso sin un LPN explícito

Este enfoque proporciona flexibilidad para trabajar con diversos proveedores y escenarios logísticos, manteniendo la integridad y trazabilidad del sistema.