
# Transport App

## **Parte Técnica**

* 🔹 **Binario ligero** escrito en  **Go** , diseñado para alto rendimiento, baja latencia y bajo consumo de memoria.
* 🔹 **Capacidad de ejecución en entornos serverless eg: Cloud Run** para escalabilidad sin gestión manual de infraestructura.
* 🔹 **Arquitectura orientada a eventos** para mayor flexibilidad y desacoplamiento.
* 🔹 **Identificadores de recursos definidos por los consumidores**, permitiendo personalización y control.
* 🔹 **Los consumidores se agrupan por organización**, asegurando aislamiento y gobernanza.
* 🔹 **Identificadores de recursos únicos por organización**, evitando colisiones entre entidades.
* 🔹 **Propagación de eventos homogéneos**: cada creación de recurso genera un evento basado en la especificación del contrato.
* 🔹 **Contratos basados en código**: la especificación del contrato se genera automáticamente desde los modelos en el código, eliminando la necesidad de mantener un `openapi.json`.

---

## **Integración con Proveedores/APIs Externos**

* 📍 **Geocodificación**: Soporte para conversión de direcciones a coordenadas mediante **LocationIQ**.
* 🚀 **Optimización de rutas**: Integración con **LocationIQ** para cálculos eficientes de rutas.
* 🏙️ **Autocompletado de direcciones**: Predicciones de direcciones en base a entrada del usuario (**LocationIQ**).
* 🌎 **Normalización de direcciones (Chile únicamente)**:
  * Basado en **Gemini modelo Flash 2.0 + Structured Outputs**.
  * Estandariza nombres de **comunas, provincias y regiones** para mayor precisión en datos locales y disminución en la redundancia de datos.

📌 **Pendientes en integración**:

* 🔄 **Algoritmos de clusterización**: Implementación de **Capacity K-Means** para optimización de agrupación de puntos de entrega.

---

## **Bases de Datos Soportadas**

* 🐘 **PostgreSQL**: Base de datos relacional robusta y ampliamente utilizada.
* ⚡ **TiDB (Serverless MySQL)**: Escalabilidad sin intervención manual, ideal para entornos dinámicos.

---

## **API Management (Opcional)**

* 🚀 **Comparación entre Apigee y Zuplo**:
  * Zuplo como alternativa moderna y flexible: [Apigee vs Zuplo](https://zuplo.com/api-gateways/apigee-alternative-zuplo).
  * Rate limit.
  * Firebase ID Token Validation.
  * Cloud Run invocation authentication.

---



## **🔐 Seguridad con Firebase Authentication**

El sistema utiliza **Firebase Authentication** para la gestión de identidad y seguridad de los usuarios.

### **🛠️ Características**

* ✅  **Autenticación basada en tokens** : Validación de `ID tokens` generados por Firebase.
* ✅  **Compatibilidad con múltiples proveedores** : Soporte para autenticación con  **Google, Apple, Facebook, email/password, y más** .
* ✅  **Validación automática de sesiones** : Prevención de accesos no autorizados mediante expiración de tokens y revocación de sesiones.
* ✅  **Integración con API Management** :
  * Uso de **Firebase ID Token Validation** en **Zuplo**.

---

## **🔹 Diseño (En progreso)**

📌 **Figma Prototype**:

🔗 [Transport App - Diseño](https://www.figma.com/design/dwZrBHZmlhe1lmOWAVqV35/TRANSPORT-APP?node-id=0-1&p=f&t=YqdEnDMSu69X568M-0)

---

## **📌 Documentación y Releases**

* 📖 **API Management & Docs**: Markdown + documentación en progreso.
  * 🔗 [Documentación Transport App](https://ignaciojeria.github.io/transport-app-docs/)
* 📝 **Release Notes Site**: **Pendiente de implementación**.

---

## **📌 Casos de Uso**

El sistema está diseñado para cubrir diversas operaciones dentro de la gestión logística:

### **👤 Gestión de Usuarios y Organizaciones**

- ✅ **Creación de cuentas (registro)**.
- ✅ **Inicio de sesión**.
- ✅ **Creación de organizaciones**.

### **🏢 Gestión de Empresas de Transporte**

- ✅ **Creación de vehículos y empresas de transporte**.
- ✅ **Búsqueda de vehículos por empresa de transporte**.

### **📍 Gestión de Nodos**

- ✅ **Creación de nodos**.
- ✅ **Búsqueda de nodos con paginación**.

### **📦 Gestión de Órdenes y Rutas**

- ✅ **Creación de órdenes**.
- ✅ **Creación de planes + rutas**.
- ✅ **Inicio de ruta** (_Out for delivery_).
- ✅ **Confirmación de entregas y no entregas**.

### **🌎 Servicios de Localización e Inteligencia Geográfica**

- ✅ **Geocodificación de direcciones** mediante **LocationIQ**.
- ✅ **Autocompletado de direcciones** basado en **LocationIQ**.
- ✅ **Normalización de direcciones** _(comunas, provincias y regiones en Chile)_ mediante **Gemini Model Flash 2.0 + Structured Outputs**.
