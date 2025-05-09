# transport-app

[![codecov](https://codecov.io/gh/Ignaciojeria/transport-app/branch/main/graph/badge.svg)](https://codecov.io/gh/Ignaciojeria/transport-app)
[![Go Report Card](https://goreportcard.com/badge/github.com/Ignaciojeria/transport-app)](https://goreportcard.com/report/github.com/Ignaciojeria/transport-app)

---

## 🔄 Propagación de Contexto en Transport App

Esta aplicación utiliza OpenTelemetry Baggage para propagar contexto entre componentes (API → eventos Pub/Sub → consumidores) de forma trazable, desacoplada y estandarizada.

### 🧾 Headers del API

El consumidor debe enviar los siguientes headers:

| Header         | Ejemplo            | Descripción                               |
| -------------- | ------------------ | ----------------------------------------- |
| `organization` | `2-CL`             | ID del tenant + país (`tenantId-country`) |
| `consumer`     | `CONSUMER_EXAMPLE` | Identificador del consumidor del API      |
| `commerce`     | `COMMERCE_EXAMPLE` | Identificador del comercio                |

Estos headers son procesados por un middleware que los desglosa e inyecta como baggage en el contexto de OpenTelemetry.

### 🧠 OpenTelemetry Baggage Keys

Estas claves se propagan en el `context.Context`:

| Clave Baggage       | Ejemplo            | Fuente                     |
| ------------------- | ------------------ | -------------------------- |
| `tenant.id`         | `2`                | Derivado de `organization` |
| `tenant.country`    | `CL`               | Derivado de `organization` |
| `business.consumer` | `CONSUMER_EXAMPLE` | Header `consumer`          |
| `business.commerce` | `COMMERCE_EXAMPLE` | Header `commerce`          |
| `event.type`        | `orderSubmitted`   | Inyectado por código       |
| `entity.type`       | `order`            | Inyectado por código       |

### 📦 Pub/Sub Message Attributes

Al publicar un evento como `orderSubmitted`, los atributos se propagan como:

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

---

## 📌 Convención de Direcciones

### Chile

| Campo JSON              | Descripción    | Ejemplo                            |
| ----------------------- | -------------- | ---------------------------------- |
| `addressLine1`          | Calle y número | `inglaterra 59`                    |
| `addressLine2`          | Info adicional | `dpto 2214`                        |
| `district`              | Comuna         | `la florida`                       |
| `province`              | Provincia      | `santiago`                         |
| `state`                 | Región         | `región metropolitana de santiago` |
| `latitude`, `longitude` | Coordenadas    | `-33.5204181`, `-70.6006178`       |

### México

| Campo JSON              | Descripción        | Ejemplo               |
| ----------------------- | ------------------ | --------------------- |
| `addressLine1`          | Calle y número     | `av insurgentes 1234` |
| `addressLine2`          | Info adicional     | `interior 5`          |
| `district`              | Localidad          | `condesa`             |
| `province`              | Alcaldía/Municipio | `cuauhtémoc`          |
| `state`                 | Estado             | `ciudad de méxico`    |
| `latitude`, `longitude` | Coordenadas        | `19.4326`, `-99.1332` |

---

## 📦 Sistema de Gestión de Paquetes

### Tipos de Paquetes

* **Con LPN:** El campo `lpn` es un identificador único físico.
* **Sin LPN:** Los ítems se tratan como agrupaciones lógicas.

### Generación de DocID

* **Con LPN:** Se usa directamente `HashByTenant(ctx, lpn)`.
* **Sin LPN:** Se usa:

  * La referencia de la orden (`referenceID`)
  * Lista de SKUs de los ítems (ordenados)
  * Un índice (`Index`) si hay paquetes con los mismos SKUs

Esto garantiza que:

* Un conjunto idéntico de ítems genera el mismo DocID
* Cambios en ítems o índice generan un nuevo DocID

### Separación de Unidades Entregables

Dividir los ítems en distintos paquetes sin LPN indica que cada paquete será una **unidad entregable**.

---

## 🧪 Ejemplos

### Paquete con LPN

```json
{
  "lpn": "PKG-12345",
  "dimensions": {"length": 10, "height": 10, "width": 10, "unit": "cm"},
  "weight": {"value": 1, "unit": "kg"},
  "insurance": {"currency": "USD", "unitValue": 10},
  "items": [
    {
      "sku": "ITEM-001",
      "description": "Producto A",
      "quantity": {"quantityNumber": 1, "quantityUnit": "unit"},
      "weight": {"value": 0.5, "unit": "kg"},
      "dimensions": {"length": 1, "height": 1, "width": 1, "unit": "cm"},
      "insurance": {"currency": "USD", "unitValue": 100}
    }
  ]
}
```

### Paquetes sin LPN – Ítems separados como unidades entregables

```json
{
  "packages": [
    {"items": [{"sku": "PROD-A", "quantity": {"quantityNumber": 1, "quantityUnit": "unit"}}]},
    {"items": [{"sku": "PROD-B", "quantity": {"quantityNumber": 1, "quantityUnit": "unit"}}]}
  ]
}
```

### Paquete sin LPN – Ítems agrupados como una unidad entregable

```json
{
  "packages": [
    {
      "items": [
        {"sku": "PROD-A", "quantity": {"quantityNumber": 1, "quantityUnit": "unit"}},
        {"sku": "PROD-B", "quantity": {"quantityNumber": 1, "quantityUnit": "unit"}}
      ]
    }
  ]
}
```

---

## 🧠 Conclusión

* Los paquetes con LPN tienen trazabilidad directa.
* Los paquetes sin LPN se identifican por los SKUs ordenados, la referencia y un índice si corresponde.
* La estructura del JSON refleja directamente la granularidad de entrega esperada.
* Separar ítems en distintos paquetes sin LPN implica que serán tratados como entregas independientes.

Esto otorga gran flexibilidad al sistema para representar diversos escenarios logísticos.
