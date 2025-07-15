# Transport All-in-One

Este directorio contiene una solución completa "all-in-one" que integra todos los servicios necesarios para una aplicación de transporte:

- **OSRM**: Servicio de enrutamiento y geocodificación
- **VROOM Optimizer**: Optimización de flotas y rutas
- **VROOM Planner**: Planificación de entregas
- **Transport App**: Aplicación principal de transporte

## 🚀 Inicio Rápido

### Opción 1: Script automático (Recomendado)

```bash
cd dagger/allinone
chmod +x start.sh
./start.sh
```

### Opción 2: Manual

```bash
cd dagger/allinone

# 1. Generar los binarios y datos (si no existen)
go run main.go

# 2. Construir y ejecutar
docker-compose up -d

# 3. Verificar estado
docker-compose exec transport-all-in-one supervisorctl status
```

## 📋 Servicios Disponibles

| Servicio | Puerto | Descripción |
|----------|--------|-------------|
| OSRM API | 5000 | Enrutamiento y geocodificación |
| VROOM Optimizer | 3000 | Optimización de flotas |
| VROOM Planner | 3001 | Planificación de entregas |
| Transport App | 8080 | API principal de transporte |
| Supervisord UI | 9001 | Interfaz web de monitoreo |

## 🛠️ Comandos Útiles

### Gestión de servicios
```bash
# Ver estado de todos los servicios
docker-compose exec transport-all-in-one supervisorctl status

# Reiniciar un servicio específico
docker-compose exec transport-all-in-one supervisorctl restart transport-app

# Ver logs en tiempo real
docker-compose exec transport-all-in-one supervisorctl tail transport-app

# Parar todos los servicios
docker-compose down

# Reiniciar todos los servicios
docker-compose restart
```

### Logs y debugging
```bash
# Ver logs de todos los servicios
docker-compose logs -f

# Ver logs de un servicio específico
docker-compose logs -f transport-all-in-one

# Acceder al contenedor
docker-compose exec transport-all-in-one bash
```

## 🔧 Configuración

### Variables de entorno
Puedes modificar las variables en `docker-compose.yml`:

```yaml
environment:
  - TZ=America/Santiago
  - TRANSPORT_APP_PORT=8080
  - VROOM_OPTIMIZER_PORT=3000
  - VROOM_PLANNER_PORT=3001
  - OSRM_PORT=5000
```

### Volúmenes montados
- `./logs:/var/log/supervisor`: Logs de todos los servicios
- `./config:/app/config:ro`: Configuraciones externas (solo lectura)

## 📊 Monitoreo

### Health Check
El contenedor incluye health checks automáticos que verifican:
- Disponibilidad de OSRM API
- Estado de todos los servicios vía supervisord

### Supervisord Web UI
Accede a http://localhost:9001 para:
- Ver estado de todos los servicios
- Reiniciar servicios individuales
- Ver logs en tiempo real

## 🚨 Troubleshooting

### Si un servicio no inicia
```bash
# Ver logs detallados
docker-compose exec transport-all-in-one cat /var/log/supervisor/transport-app.log

# Verificar que los binarios existen
docker-compose exec transport-all-in-one ls -la /transport-app/

# Verificar permisos
docker-compose exec transport-all-in-one ls -la /transport-app/transport-app
```

### Si necesitas regenerar los datos
```bash
# Parar servicios
docker-compose down

# Regenerar con Dagger
go run main.go

# Reconstruir y ejecutar
docker-compose up -d --build
```

### Problemas de memoria
El contenedor está configurado con límites de memoria:
- **Límite**: 4GB
- **Reserva**: 2GB

Si tienes problemas de rendimiento, considera aumentar estos valores en `docker-compose.yml`.

## 🏗️ Arquitectura

```
┌─────────────────────────────────────────────────────────────┐
│                    Transport All-in-One                    │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐         │
│  │   OSRM      │ │   VROOM     │ │   VROOM     │         │
│  │   API       │ │ Optimizer   │ │  Planner    │         │
│  │  (5000)     │ │   (3000)    │ │   (3001)    │         │
│  └─────────────┘ └─────────────┘ └─────────────┘         │
│                                                           │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Transport App                         │   │
│  │              (8080)                               │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                           │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Supervisord                          │   │
│  │              (9001)                               │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## 📝 Características

- ✅ **Todo en uno**: Todos los servicios en un solo contenedor
- ✅ **Auto-restart**: Servicios se reinician automáticamente si fallan
- ✅ **Logs centralizados**: Todos los logs en `/var/log/supervisor`
- ✅ **Monitoreo**: Health checks y UI web de supervisord
- ✅ **Configuración flexible**: Variables de entorno y volúmenes
- ✅ **Binarios estáticos**: Sin dependencias del sistema
- ✅ **Docker Compose**: Fácil despliegue y gestión 