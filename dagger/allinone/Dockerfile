FROM alpine:latest

# Instalar supervisord y herramientas necesarias
RUN apk update && \
    apk add --no-cache supervisor bash curl && \
    mkdir -p /var/log/supervisor /var/run

# Copiar los binarios estáticos de OSRM
COPY ./osrm-static/ /usr/local/bin/

# Copiar los datos de OSRM
COPY ./osrm-data/ /data/

# Copiar la configuración de supervisord
COPY supervisord.conf /etc/supervisor/conf.d/osrm.conf

# Hacer ejecutables los binarios
RUN chmod +x /usr/local/bin/osrm-routed-static && \
    chmod +x /usr/local/bin/osrm-extract-static && \
    chmod +x /usr/local/bin/osrm-partition-static && \
    chmod +x /usr/local/bin/osrm-customize-static

# Exponer puertos
EXPOSE 5000 9001

# Comando de inicio
CMD ["supervisord", "-c", "/etc/supervisor/conf.d/osrm.conf"] 