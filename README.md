# StreamX - Backend Server

StreamX es una plataforma de streaming multimedia que permite a los usuarios buscar,
subir y reproducir contenido multimedia. El backend de StreamX está diseñado para
funcionar en entornos de red escalables y distribuidos, utilizando un balanceador de
carga con Nginx para gestionar el tráfico entre tres instancias del servidor. Este
sistema es ideal para manejar grandes volúmenes de contenido, almacenando archivos
multimedia en PostgreSQL utilizando su capacidad de almacenamiento en BLOB (Binary
Large Object).

## Características clave

- Búsqueda de contenido multimedia: API pública que permite a los usuarios buscar
  videos y audios almacenados en la plataforma.
- Subida de contenido multimedia: Soporte para subir archivos de audio y video a
  través de la API.
- Streaming: Capacidad de transmitir el contenido almacenado directamente a los
  usuarios.
- Escalabilidad: El backend está distribuido entre tres instancias, con Nginx
  gestionando el balanceo de carga.
- Este sistema está diseñado para funcionar eficientemente en entornos de red tanto
  locales como en la nube.

# Instrucciones de Instalación
## Requisitos Previos
PostgreSQL configurado para soportar almacenamiento BLOB.
Nginx para balancear la carga entre instancias del servidor.
Herramientas de compilación de Go instaladas para ejecutar el servidor.

## Instalación:

1. **Clonar el repositorio**

```bash
git clone https://github.com/zSnails/StreamX.git
cd streamx-backend
```

2. **Configurar las variables de entorno**: Crea un archivo `.env` con las siguientes variables:

```bash
Copy code
DB_USER=
DB_PASSWORD=
DB_HOST=
# El formato del puerto es el siguiente :<número> (con el : antes del número y sin los <>)
PORT=
```
3. **Configurar Nginx para balanceo de carga**: Asegúrate de que tu archivo de
   configuración de Nginx incluya las tres instancias del servidor backend. Aquí
   tienes un ejemplo de configuración:

```nginx
events {}
http {
    client_max_body_size 1G;
    upstream backend {
            server localhost:8081;
            server localhost:8082;
            server localhost:8083;
    }

    resolver localhost;

    server {
        listen 8080;
        location / {
            proxy_pass http://backend;
        }
    }
}
```

4. **Compilar el Servidor de Backend**: Desde la carpeta principal del proyecto
   ejecuta el siguiente comando.

```bash
go build ./cmd/streamx
```

5. **Lanzar las instancias del servidor con mprocs**: Usa mprocs para lanzar las tres
   instancias del servidor en puertos diferentes. Crea un archivo de configuración
   para mprocs y ejecuta el siguiente comando `mprocs`:

`mprocs.yaml`
```yaml
procs:
  server-1:
    cmd: ["./streamx", "-port", ":8081"]
  server-2:
    cmd: ["./streamx", "-port", ":8082"]
  server-3:
    cmd: ["./streamx", "-port", ":8083"]
```

# Configuraciones Recomendadas

Para entornos locales y de prueba, se recomienda ejecutar las tres instancias
utilizando mprocs con Nginx balanceando el tráfico. En entornos de producción, puedes
distribuir las instancias entre diferentes servidores para mejorar la redundancia y
escalabilidad.


# Autores:

- Aaron González
    - Desarrollador principal del backend.
- Daniel Porras
    - Desarrollador secundario del backend.
