# Deployment Plan

## Estado

No se ejecuto un despliegue real desde este entorno. El repositorio quedo preparado para desplegarse en Render con Blueprint.

## Plataforma elegida

Render.

Motivo:

- acepta el repo actual sin reestructurarlo
- soporta Docker para `api-go` y `api-node`
- permite publicar el frontend como static site
- deja todo definido en un solo archivo `render.yaml`

## Servicios a crear

El archivo [render.yaml](C:/Users/theda/Documents/talsory_interseguro/repo/render.yaml) define:

1. `talsory-interseguro-api-node`
2. `talsory-interseguro-api-go`
3. `talsory-interseguro-frontend`

## Como desplegarlo en Render

1. Entra a Render.
2. Conecta tu cuenta de GitHub si aun no lo hiciste.
3. Haz clic en `New > Blueprint`.
4. Selecciona el repositorio `DanielRCor/talsory_interseguro`.
5. Usa la rama `main`.
6. Confirma que Render detecte `render.yaml`.
7. Revisa los 3 servicios que se van a crear.
8. Haz clic en `Deploy Blueprint`.

## Configuracion clave

### api-node

- runtime: Docker
- plan: free
- health check: `/health`
- auth: activada
- `JWT_SECRET`: generado por Render

### api-go

- runtime: Docker
- plan: free
- health check: `/health`
- auth: activada
- `NODE_API_URL`: toma automaticamente la URL publica de `api-node`
- `JWT_SECRET`: reutiliza el mismo secreto de `api-node`

### frontend

- runtime: static site
- build: `npm ci && npm run build`
- salida: `dist`
- `VITE_GO_API_BASE`: toma automaticamente la URL publica de `api-go`
- `VITE_NODE_API_BASE`: toma automaticamente la URL publica de `api-node`

## Verificacion despues del deploy

1. Abrir la URL del frontend en Render.
2. Verificar que las tarjetas de salud respondan.
3. Generar un JWT demo.
4. Ejecutar el flujo `QR + Estadisticas`.
5. Probar la rotacion.

## Limitaciones del free tier

- Los web services free pueden entrar en reposo despues de 15 minutos sin trafico.
- El primer request despues de ese reposo puede tardar cerca de 1 minuto.
- Render da 750 horas gratis al mes por workspace para web services free.

## Si algo falla

- Revisar logs de `api-node` primero.
- Luego revisar logs de `api-go`.
- Si el frontend carga pero falla el flujo, normalmente el problema sera:
  - servicio dormido
  - variable de entorno mal resuelta
  - primer deploy aun no terminado

## Pendiente despues del despliegue

Despues de que Render lo publique, lo siguiente seria:

1. validar el flujo con la URL real
2. agregar dominio custom si quieres una entrega mas pulida
3. dejar una demo grabada corta como respaldo para entrevista
