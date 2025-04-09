## transport-app

[![codecov](https://codecov.io/gh/Ignaciojeria/transport-app/branch/main/graph/badge.svg)](https://codecov.io/gh/Ignaciojeria/transport-app) [![Go Report Card](https://goreportcard.com/badge/github.com/Ignaciojeria/transport-app)](https://goreportcard.com/report/github.com/Ignaciojeria/transport-app)

## 📌 Convención de Campos para Direcciones en Chile

| Campo en JSON                | Descripción                      | Ejemplo                       |
| ---------------------------- | --------------------------------- | ----------------------------- |
| `addressInfo.addressLine1` | Calle y número                    | `inglaterra 59`             |
| `addressInfo.addressLine2` | Información adicional (opcional)  | `dpto 2214`                 |
| `addressInfo.district`     | Comuna                            | `la florida`                |
| `addressInfo.province`     | Provincia                         | `santiago`                  |
| `addressInfo.state`        | Región                            | `metropolitana de santiago` |
| `addressInfo.latitude`     | Latitud                           | `-33.5204181`               |
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
    "latitude": "-33.4085",
    "longitude": "-70.5666"
  }
}
```

# 🔄 Context Propagation in Transport App

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
