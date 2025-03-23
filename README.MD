## transport-app

[![codecov](https://codecov.io/gh/Ignaciojeria/transport-app/branch/main/graph/badge.svg)](https://codecov.io/gh/Ignaciojeria/transport-app) [![Go Report Card](https://goreportcard.com/badge/github.com/Ignaciojeria/transport-app)](https://goreportcard.com/report/github.com/Ignaciojeria/transport-app)

## 📌 Convención de Campos para Direcciones en Chile

| Campo en JSON                | Descripción                      | Ejemplo                       |
| ---------------------------- | --------------------------------- | ----------------------------- |
| `addressInfo.addressLine1` | Calle y número                    | `inglaterra 59`             |
| `addressInfo.addressLine2` | Información adicional (opcional)  | `dpto 2214`                 |
| `addressInfo.district`     | Comuna                            | `la-florida`                |
| `addressInfo.province`     | Provincia                         | `santiago`                  |
| `addressInfo.state`        | Región                            | `metropolitana-de-santiago` |
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
    "state": "metropolitana-de-santiago",
    "latitude": "-33.4263",
    "longitude": "-70.6200"
  }
}
```

### 🏙️ Ejemplo en Santiago (sin departamento)

```json
{
  "addressInfo": {
    "addressLine1": "avenida las condes 7890",
    "addressLine2": "",
    "district": "las-condes",
    "province": "santiago",
    "state": "metropolitana-de-santiago",
    "latitude": "-33.4085",
    "longitude": "-70.5666"
  }
}
```

### 🏠 Ejemplo en el Sur (zona urbana)

```json
{
  "addressInfo": {
    "addressLine1": "avenida diego portales 1010",
    "addressLine2": "edificio costanera piso 5",
    "district": "puerto-montt",
    "province": "llanquihue",
    "state": "los-lagos",
    "latitude": "-41.4717",
    "longitude": "-72.9361"
  }
}
```

### 🌿 Ejemplo en el Sur (zona rural)

```json
{
  "addressInfo": {
    "addressLine1": "ruta 5 sur km 789",
    "addressLine2": "sector rucapequen",
    "district": "san-carlos",
    "province": "punilla",
    "state": "nuble",
    "latitude": "-36.4241",
    "longitude": "-71.9598"
  }
}
```
