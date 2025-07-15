# transport-app

[![codecov](https://codecov.io/gh/Ignaciojeria/transport-app/branch/main/graph/badge.svg)](https://codecov.io/gh/Ignaciojeria/transport-app)
[![Go Report Card](https://goreportcard.com/badge/github.com/Ignaciojeria/transport-app)](https://goreportcard.com/report/github.com/Ignaciojeria/transport-app)

---

# 🗺️ Optimización de rutas con OSRM + VROOM

Este entorno Docker permite levantar una instancia local de **OSRM** para ruteo y **VROOM** para optimización de rutas de última milla. Ideal para SaaS logísticos, simulaciones o soluciones descentralizadas.

## 📦 Requisitos

- Docker
- Docker Compose
- Conexión a Internet para descargar los mapas

## 📁 Estructura del proyecto

```
.
├── docker-compose.yml
└── osrm/
    └── chile-latest.osm.pbf      # Archivo del mapa descargado manualmente
```

## 🌍 Paso 1: Descargar el mapa de Chile

```bash
mkdir -p ./osrm
curl -L https://download.geofabrik.de/south-america/chile-latest.osm.pbf -o ./osrm/chile-latest.osm.pbf
```

Puedes cambiar el país o región desde: https://download.geofabrik.de/

## 🐳 Paso 2: Levantar los servicios

```bash

✅ Comandos para levantar cada entorno

🔧 Solo backend (desarrollo, sin transport-app)
docker-compose -f docker-compose.yml up

🚀 Producción (con transport-app)
docker-compose -f docker-compose.allinone.yml pull
docker-compose -f docker-compose.allinone.yml up 

Aquí docker-compose.allinone.yml agrega transport-app al conjunto de servicios, sin necesidad de duplicar el resto.
```


# JWT Setup Guide - Transport App

## 🔑 Generación de Claves RSA

### **Requisitos previos:**
- OpenSSL instalado en el sistema
- Acceso a terminal/PowerShell

### **Comandos para generar claves JWT y exponer .well-known/jwks.json:**
Estos comandos generan una clave privada y su clave pública correspondiente, necesarias para firmar tokens JWT y exponer tu endpoint .well-known/jwks.json como Authorization Server.
```bash
# Generar clave privada PKCS#8
openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:2048

# Extraer clave pública
openssl rsa -pubout -in private_key.pem -out public_key.pem
```

## 🔐 Comando para generar clave de encriptación de Client Credentials (AES-256-GCM)
Esta clave se utiliza para encriptar y desencriptar los client_secret antes de almacenarlos en la base de datos.
```bash
# Generar 32 bytes aleatorios en Base64 (clave de 256 bits)
openssl rand -base64 32
```

Esto ejecutará:

- **OSRM**: procesará `chile-latest.osm.pbf` (extract → partition → customize) y servirá en `http://localhost:5000`
- **VROOM**: se conectará a OSRM y servirá en `http://localhost:3000`

## ✅ Verificación de estado

```bash
# Verifica OSRM
curl "http://localhost:5000/route/v1/driving/-70.65,-33.45;-70.66,-33.46"

# Verifica VROOM Powershell
$body = '{"jobs":[{"id":1,"location":[-70.6483,-33.4372]}],"vehicles":[{"id":1,"start":[-70.6483,-33.4372]}]}'
Invoke-RestMethod -Uri "http://localhost:3000/" -Method Post -Body $body -ContentType "application/json"
```

## 📤 Ejemplo de solicitud a VROOM

```bash
curl -X POST http://localhost:3000 \
  -H "Content-Type: application/json" \
  -d '{
    "vehicles": [{ "id": 1, "start": [ -70.653, -33.45 ] }],
    "jobs": [
      { "id": 1, "location": [ -70.65, -33.44 ] },
      { "id": 2, "location": [ -70.66, -33.46 ] }
    ]
  }'
```

## ⚙️ Personalización

- Cambia el mapa reemplazando `chile-latest.osm.pbf` en `./osrm/`.
- Puedes modificar `car.lua` para ajustar reglas de tráfico.
- VROOM permite una configuración avanzada vía `config.yml`.

## 🧩 Extensiones posibles

- Convertir en sidecar para microservicio de optimización.
- Usar `zrok` para habilitar nodos descentralizados P2P.
- Integrar sistema de staking y recompensas basado en rendimiento.
- Registrar nodos vía healthcheck en un SaaS logístico.


## 🔄 Propagación de Contexto en Transport App

Esta aplicación utiliza OpenTelemetry Baggage para propagar contexto entre componentes (API → eventos Pub/Sub → consumidores) de forma trazable, desacoplada y estandarizada.

### 🧾 Headers del API

El consumidor debe enviar los siguientes headers:

El consumidor debe enviar los siguientes headers:

| Header         | Ejemplo            | Descripción                                                                            |
| -------------- | ------------------ | -------------------------------------------------------------------------------------- |
| `tenant`       | `53f0cea9-ee4c-4c0f-adc8-dabc850e4d7b-CL`             | ID del tenant + país (`tenantId-country`). Siempre se utiliza para operar el recurso.  |
| `consumer`     | `CONSUMER_EXAMPLE` | Identificador técnico de quien inicializó la creación del recurso vía API. (origen)             |
| `commerce`     | `COMMERCE_EXAMPLE` | Identificador del comercio o tienda que genera la necesidad inicial de operar el recurso.       |
| `channel`      | `WEB` / `APP` / `YOUR_API_SERVICE`      | Canal de operación previo a creación del recurso (planificado : Web, Entregado: App Mobile, etc.). Puede cambiar. |

Estos headers son procesados por un middleware que los desglosa e inyecta como baggage en el contexto de OpenTelemetry.

### 🧠 OpenTelemetry Baggage Keys

Estas claves se propagan en el `context.Context`:

| Clave Baggage       | Ejemplo                  | Fuente                     |
| ------------------- | ------------------------ | -------------------------- |
| `tenant.id`         | `2`                      | Derivado de `tenant` |
| `tenant.country`    | `CL`                     | Derivado de `tenant` |
| `consumer`          | `CONSUMER_EXAMPLE`       | Header `consumer`          |
| `commerce`          | `COMMERCE_EXAMPLE`       | Header `commerce`          |
| `channel`           | `WEB/MOBILE/INTEGRATION` | Header `channel`           |
| `event.type`        | `orderSubmitted`         | Inyectado por código       |
| `entity.type`       | `order`                  | Inyectado por código       |

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
| `state`                 | Región         | `region metropolitana de santiago` |
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

## Filtros

### CoordinatesConfidenceLevel

🎯 **Filtros y resultados esperados**

| Filtro Min     | Filtro Max     | Resultado esperado (IDs) | Lógica SQL aplicada                                      |
|----------------|----------------|---------------------------|-----------------------------------------------------------|
| *(sin filtro)* | *(sin filtro)* | A, B, C, D, E, F          | No se aplica condición (`WHERE 1=1`)                      |
| 0.5            | *(sin filtro)* | C, D, E, F                | `level >= 0.5`                                            |
| *(sin filtro)* | 0.5            | A, B                      | `level <= 0.5`                                            |
| 0.5            | 0.8            | C, D                      | `level >= 0.5 AND level <= 0.8`                           |
| 0.7            | 0.7            | *(vacío)*                 | `level >= 0.7 AND level <= 0.7` (match exacto)            |
| 1.1            | *(sin filtro)* | *(vacío)*                 | `level >= 1.1` (ningún registro cumple)                   |
| *(sin filtro)* | 0.0            | *(vacío)*                 | `level <= 0.0` (ningún registro cumple)                   |
| 0.0            | 1.0            | A, B, C, D, E, F          | `level >= 0.0 AND level <= 1.0` (todos los registros)     |


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


# 📦 Consulta GraphQL para Delivery Units (Relay Style)

Este sistema expone una API GraphQL basada en la especificación **Relay** para gestionar consultas de unidades entregables (Delivery Units).

## 🚀 Query: Obtener Delivery Units

```graphql
query GetDeliveryUnitsReports(
  $filter: DeliveryUnitsReportFilterInput, 
  $first: Int, 
  $after: String
) {
  deliveryUnitsReports(
    filter: $filter,
    first: $first,
    after: $after
  ) {
    edges {
      cursor
      node {
        referenceID
        collectAvailabilityDate {
          date
          timeRange {
            startTime
            endTime
          }
        }
        destination {
          addressInfo {
            addressLine1
            addressLine2
            contact {
              additionalContactMethods {
                type
                value
              }
              documents {
                type
                value
              }
              email
              fullName
              nationalID
              phone
            }
            district
            latitude
            longitude
            province
            state
            timeZone
            zipCode
          }
          deliveryInstructions
          nodeInfo {
            referenceId
            name
          }
        }
        origin {
          addressInfo {
            addressLine1
            addressLine2
            contact {
              additionalContactMethods {
                type
                value
              }
              documents {
                type
                value
              }
              email
              fullName
              nationalID
              phone
            }
            district
            latitude
            longitude
            province
            state
            timeZone
            zipCode
          }
          deliveryInstructions
          nodeInfo {
            referenceId
            name
          }
        }
        orderType {
          type
          description
        }
        package {
          dimensions {
            length
            height
            width
            unit
          }
          insurance {
            currency
            unitValue
          }
          items {
            sku
            description
            dimensions {
              length
              height
              width
              unit
            }
            insurance {
              currency
              unitValue
            }
            skills {
              type
              value
              description
            }
            quantity {
              quantityNumber
              quantityDelivered
              quantityUnit
            }
            weight {
              unit
              value
            }
          }
          labels {
            type
            value
          }
          lpn
          weight {
            unit
            value
          }
        }
        promisedDate {
          dateRange {
            startDate
            endDate
          }
          serviceCategory
          timeRange {
            startTime
            endTime
          }
        }
        references {
          type
          value
        }
        extraFields {
          key
          value
        }
        carrier {
          nationalID
          name
        }
        vehicle {
          plate
        }
        driver {
          nationalID
          name
          email
        }
        route {
          routeID
          routePosition
        }
        delivery {
          status
          handledAt
          failure {
            detail
            reason
            referenceID
          }
          location {
            latitude
            longitude
          }
          recipient {
            fullName
            nationalID
          }
        }
      }
    }
    pageInfo {
      hasNextPage
      endCursor
    }
  }
}
```

## 📖 Consideraciones importantes

* **Relay Style:** La API utiliza el patrón Relay para manejar paginación (`edges`, `cursor`, `pageInfo`).
* **PageInfo:** Permite saber si hay más resultados (`hasNextPage`) y avanzar (`endCursor`).

Con esta estructura puedes construir reportes logísticos detallados, dashboards de entregas, o integraciones hacia otros sistemas de monitoreo.