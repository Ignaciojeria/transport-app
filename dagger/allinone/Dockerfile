# Dockerfile que usa los archivos extraídos por Dagger
FROM debian:bookworm-slim

# Instalar supervisord y herramientas necesarias
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        supervisor \
        curl \
        git-core \
        build-essential \
        g++ \
        libssl-dev \
        libasio-dev \
        libglpk-dev \
        pkg-config \
        libssl3 \
        libglpk40 \
        nodejs \
        npm && \
    mkdir -p /var/log/supervisor /var/run /conf/optimizer /conf/planner /vroom-optimizer-app /vroom-planner-app /transport-app && \
    rm -rf /var/lib/apt/lists/*

# Copiar archivos extraídos por Dagger
COPY ./osrm-data/ /data/
COPY ./osrm-static/ /usr/local/bin/
COPY ./vroom-optimizer/vroom-optimizer-static /usr/local/bin/vroom-optimizer
COPY ./vroom-planner/vroom-planner-static /usr/local/bin/vroom-planner
COPY ./vroom-express/ /vroom-express/
COPY ./transport-app/transport-app /transport-app/transport-app

# Copiar configuraciones específicas
COPY ./vroom-optimizer/config.yml /conf/optimizer/config.yml
COPY ./vroom-planner/config.yml /conf/planner/config.yml
COPY supervisord.conf /etc/supervisor/conf.d/osrm.conf

# Hacer ejecutables los binarios
RUN chmod +x /usr/local/bin/osrm-routed-static && \
    chmod +x /usr/local/bin/osrm-extract-static && \
    chmod +x /usr/local/bin/osrm-partition-static && \
    chmod +x /usr/local/bin/osrm-customize-static && \
    chmod +x /usr/local/bin/vroom-optimizer && \
    chmod +x /usr/local/bin/vroom-planner && \
    chmod +x /transport-app/transport-app

# Crear directorios separados para cada servicio y configurar enlaces
RUN cp -r /vroom-express/* /vroom-optimizer-app/ && \
    cp -r /vroom-express/* /vroom-planner-app/ && \
    ln -sf /usr/local/bin/vroom-optimizer /vroom-optimizer-app/vroom && \
    ln -sf /usr/local/bin/vroom-planner /vroom-planner-app/vroom && \
    # Copiar configuraciones específicas a cada directorio de aplicación
    cp /conf/optimizer/config.yml /vroom-optimizer-app/config.yml && \
    cp /conf/planner/config.yml /vroom-planner-app/config.yml && \
    # Instalar dependencias de Node.js para cada aplicación
    cd /vroom-optimizer-app && npm install && \
    cd /vroom-planner-app && npm install

# Exponer puertos (OSRM: 5000, VROOM Optimizer: 3000, VROOM Planner: 3001, Supervisor: 9001)
EXPOSE 5000 3000 3001 9001

# Comando de inicio
CMD ["supervisord", "-c", "/etc/supervisor/conf.d/osrm.conf"] 