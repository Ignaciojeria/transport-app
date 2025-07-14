# OSRM con Supervisord

Este directorio contiene los archivos necesarios para ejecutar OSRM con supervisord usando los binarios estáticos generados por Dagger.

## Archivos generados por Dagger

- `./osrm-data/` - Datos procesados de OSRM para Chile
- `./osrm-static/` - Binarios estáticos de OSRM

## Archivos de configuración

- `supervisord.conf` - Configuración de supervisord para ejecutar OSRM
- `Dockerfile` - Dockerfile para crear la imagen con supervisord

## Cómo usar

### 1. Generar los binarios y datos (si no existen)

```bash
cd dagger/allinone
go run main.go
```

### 2. Construir la imagen Docker

```bash
docker build -t osrm-supervisord .
```

### 3. Ejecutar el contenedor

```bash
docker run -d \
  --name osrm-test \
  -p 5000:5000 \
  -p 9001:9001 \
  osrm-supervisord
```

### 4. Verificar que OSRM esté funcionando

```bash
# Verificar el estado de supervisord
docker exec osrm-test supervisorctl status

# Ver logs de OSRM
docker exec osrm-test tail -f /var/log/supervisor/osrm-routed.log

# Probar OSRM
curl "http://localhost:5000/route/v1/driving/-70.6483,-33.4372;-70.6500,-33.4400?overview=false"
```

### 5. Comandos útiles de supervisord

```bash
# Ver estado de todos los procesos
docker exec osrm-test supervisorctl status

# Reiniciar OSRM
docker exec osrm-test supervisorctl restart osrm-routed

# Parar OSRM
docker exec osrm-test supervisorctl stop osrm-routed

# Iniciar OSRM
docker exec osrm-test supervisorctl start osrm-routed

# Ver logs en tiempo real
docker exec osrm-test supervisorctl tail osrm-routed
```

## Puertos

- **5000**: OSRM API
- **9001**: Interfaz web de supervisord (opcional)

## Características

- **Reinicio automático**: Si OSRM se cae, supervisord lo reinicia automáticamente
- **Logs rotativos**: Los logs se guardan con rotación automática
- **Monitoreo**: Puedes monitorear el estado a través de supervisord
- **Binarios estáticos**: No depende de librerías del sistema

## Troubleshooting

### Si OSRM no inicia

```bash
# Ver logs detallados
docker exec osrm-test cat /var/log/supervisor/osrm-routed.log

# Verificar que los binarios existen
docker exec osrm-test ls -la /usr/local/bin/

# Verificar que los datos existen
docker exec osrm-test ls -la /data/
```

### Si necesitas regenerar los datos

```bash
# Eliminar contenedor
docker stop osrm-test && docker rm osrm-test

# Regenerar con Dagger
go run main.go

# Reconstruir imagen
docker build -t osrm-supervisord .

# Ejecutar nuevamente
docker run -d --name osrm-test -p 5000:5000 -p 9001:9001 osrm-supervisord
``` 