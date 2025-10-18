# 📾 Contract Definition — Product Schema v1

Este contrato define la estructura completa de un **producto vendible** dentro del ecosistema **Einar**.
Representa artículos **físicos**, **digitales** o **mixtos (bundles)**, incluyendo condiciones comerciales, atributos descriptivos, estructura de componentes y reglas logísticas.

---

## 🧱 Estructura general

| Campo                      | Tipo                | Descripción                                                                                     |
| -------------------------- | ------------------- | ----------------------------------------------------------------------------------------------- |
| **referenceID**            | `string`            | Identificador interno o externo del producto. Puede provenir de un ERP, catálogo o importación. |
| **name**                   | `string`            | Nombre visible para el cliente.                                                                 |
| **descriptionMarkdown**    | `string (Markdown)` | Descripción extendida del producto. Soporta formato Markdown.                                   |
| **status**                 | `object`            | Define la disponibilidad, visibilidad y comportamiento del producto.                            |
| **attachments**            | `array<object>`     | Archivos informativos o descargables (ej. fichas técnicas, manuales, certificados).             |
| **properties**             | `object`            | Propiedades comerciales como SKU, marca o código de barras.                                     |
| **purchaseConditions**     | `object`            | Define límites y múltiplos de compra por tipo de unidad (fija, peso, volumen).                  |
| **attributes**             | `array<object>`     | Características visibles o filtrables (color, variedad, origen, etc.).                          |
| **categories**             | `array<object>`     | Estructura jerárquica de clasificación del producto.                                            |
| **digitalBundle**          | `object`            | Representa un componente digital adicional (descarga, guía o acceso exclusivo).                 |
| **welcomeMessageMarkdown** | `string (Markdown)` | Mensaje mostrado por el bot o interfaz al presentar el producto.                                |
| **media**                  | `object`            | Contiene galería de imágenes y/o videos promocionales.                                          |
| **payment**                | `object`            | Define moneda, métodos y proveedor de pago disponibles.                                         |
| **stock**                  | `object`            | Cantidades disponibles (por unidades, peso o volumen).                                          |
| **price**                  | `object`            | Define el precio base o por unidad.                                                             |
| **cost**                   | `object`            | Costos de adquisición o fabricación, usados para cálculo de margen.                             |
| **components**             | `array<object>`     | Subproductos o agregados (ingredientes, addons, complementos).                                  |
| **logistics**              | `object`            | Contiene información física, horarios y costos de entrega.                                      |

---

## 🧩 Definición de bloques

### 🧠 `status`

Controla la visibilidad y el comportamiento del producto.

```json
"status": {
  "isAvailable": true,
  "isFeatured": false,
  "allowReviews": true
}
```

| Campo          | Tipo      | Descripción                                   |
| -------------- | --------- | --------------------------------------------- |
| `isAvailable`  | `boolean` | Si el producto está disponible para la venta. |
| `isFeatured`   | `boolean` | Si debe destacarse en catálogos o campañas.   |
| `allowReviews` | `boolean` | Si permite reseñas o calificaciones.          |

---

### 💎 `attachments`

Archivos asociados al producto.

```json
"attachments": [
  {
    "name": "Ficha técnica",
    "description": "Detalles sobre calidad y procedencia.",
    "url": "https://example.com/ficha.pdf",
    "type": "pdf",
    "sizeKb": 450
  }
]
```

---

### 🏷️ `properties`

Identificadores y metadatos comerciales.

```json
"properties": {
  "sku": "PALTA-HASS-200G",
  "brand": "Einar Produce",
  "barcode": "7801234567890"
}
```

---

### 💰 `purchaseConditions`

Define los límites mínimos, máximos y múltiplos de compra.

```json
"purchaseConditions": {
  "fixed": {
    "minUnits": 1,
    "maxUnits": 10,
    "multiplesOf": 1
  },
  "weight": {
    "minWeight": 1000,
    "maxWeight": 10000,
    "multiplesOf": 250,
    "notes": "La compra mínima es 1kg y máxima 10kg."
  }
}
```

---

### 🧲 `attributes`

Permite definir etiquetas personalizables visibles para el cliente.

```json
"attributes": [
  { "name": "color", "value": "verde" },
  { "name": "variedad", "value": "Hass" },
  { "name": "origen", "value": "Chile" }
]
```

---

### 🗂️ `categories`

Jerarquía que representa la posición del producto dentro de un catálogo.

```json
"categories": [
  { "id": "frutas", "name": "Frutas", "parent": null },
  { "id": "frutas/paltas", "name": "Paltas", "parent": "frutas" },
  { "id": "frutas/paltas/hass", "name": "Palta Hass", "parent": "frutas/paltas" }
]
```

---

### 📻 `digitalBundle`

Contenido digital adicional asociado al producto.

```json
"digitalBundle": {
  "hasDigitalContent": true,
  "type": "downloadable",
  "title": "Recetario de Desayunos Saludables",
  "description": "Incluye 10 recetas para acompañar tu Palta Hass.",
  "access": {
    "method": "link",
    "url": "https://example.com/recetario.pdf",
    "expiresInDays": 30
  }
}
```

---

### 💬 `welcomeMessageMarkdown`

Mensaje que el bot o la UI muestran al presentar el producto.

```json
"welcomeMessageMarkdown": "🥑 **¡Hola! Soy Einar Bot.** Hoy tenemos Palta Hass fresca lista para tu desayuno."
```

---

### 🎥 `media`

Imágenes o videos relacionados con el producto.

```json
"media": {
  "videos": [
    {
      "title": "Conoce la Palta Hass",
      "platform": "YouTube",
      "url": "https://www.youtube.com/watch?v=abcd1234",
      "thumbnail": "https://img.youtube.com/vi/abcd1234/hqdefault.jpg"
    }
  ],
  "gallery": [
    "https://example.com/images/palta-1.jpg",
    "https://example.com/images/palta-2.jpg"
  ]
}
```

---

### 💳 `payment`

Define las reglas de pago.

```json
"payment": {
  "currency": "CLP",
  "methods": ["credit_card", "debit_card", "transfer"],
  "provider": "Transbank"
}
```

---

### 📦 `stock`

Controla la cantidad disponible en distintos formatos.

```json
"stock": {
  "fixed": { "availableUnits": 40 },
  "weight": { "availableWeight": 250 },
  "volume": { "availableVolume": 0 }
}
```

---

### 💵 `price`

Define el precio base o unitario.

```json
"price": {
  "fixedPrice": 250000,
  "weight": { "unitSize": 1, "pricePerUnit": 1000 },
  "volume": { "unitSize": 1, "pricePerUnit": 0 }
}
```

---

### 📉 `cost`

Define los costos internos o de adquisición.

```json
"cost": {
  "fixedCost": 180000,
  "weight": { "unitSize": 1, "costPerUnit": 700 },
  "volume": { "unitSize": 1, "costPerUnit": 0 }
}
```

---

### 🧩 `components`

Subproductos o ingredientes del producto principal.

```json
"components": [
  {
    "type": "base",
    "name": "Palta Hass",
    "quantity": "200g",
    "required": true
  },
  {
    "type": "addon",
    "name": "Pan tostado",
    "description": "Pan artesanal tostado al momento",
    "price": 500,
    "required": false
  }
]
```

---

### 🚚 `logistics`

Información logística y restricciones de entrega.

```json
"logistics": {
  "dimensions": { "height": 100, "length": 100, "width": 100 },
  "weight": 250,
  "availabilityTime": [
    {
      "timeRange": { "from": "09:00", "to": "22:00" },
      "daysOfWeek": ["mon", "tue", "wed", "thu", "fri"]
    },
    {
      "timeRange": { "from": "10:00", "to": "18:00" },
      "daysOfWeek": ["sat"]
    }
  ],
  "deliveryFees": [
    {
      "condition": "prime",
      "type": "fixed",
      "value": 50,
      "timeRange": { "from": "09:00", "to": "18:00" }
    },
    {
      "condition": "night",
      "type": "fixed",
      "value": 100,
      "timeRange": { "from": "18:01", "to": "22:00" }
    }
  ]
}
```

---

## 📘 Convenciones generales

* Todos los textos con formato deben indicarse como `Markdown` (por ejemplo, `descriptionMarkdown` o `welcomeMessageMarkdown`).
* Todos los montos deben expresarse en **moneda local** (por ejemplo, CLP o USD).
* Las medidas (`weight`, `volume`, `dimensions`) deben expresarse en **gramos, mililitros y milímetros** por defecto.
* Los campos `id` de categorías deben seguir una **convención jerárquica** (`frutas/paltas/hass`).
* Se recomienda incluir un campo `meta` en el nivel raíz:

```json
"meta": {
  "version": "1.0.0",
  "contractType": "product",
  "createdAt": "2025-10-15T15:00:00Z"
}
```

---

> ✅ **Este contrato está optimizado para interoperar con catálogos Einar, asistentes LLM y motores de órdenes.**
> Puede ser serializado en JSON, YAML o integrado en bases de datos documentales sin pérdida semántica.