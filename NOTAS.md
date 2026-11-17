# Notas de estado y pendientes

Estado del proyecto AIRecon a 17/11/2026 y cosas que podrían quedar por hacer.

## Hecho

- **Rebranding completo a AIRecon.** Nombre de app, usuario de login, módulo Go
  (`github.com/diegoperezg7/ai-recon`), binario (`cmd/airecon`), imagen
  (`airecon:local`), usuario y base de datos de Postgres (`airecon`), variables de
  entorno (`AIRECON_*`) y ficheros de detección (`airecon_*`). Cero rastros de
  `artex` ni del fork intermedio en el árbol versionado.
- **Logo y paleta Darc3** (mandala azul), favicon incluido.
- **Idioma ES/EN.** Selector de interfaz solo con esos dos; informe de pentest en
  español (antes salía en coreano); español por defecto.
- **Login rediseñado**: columna centrada, logo arriba, nombre AIRecon.
- **Menú lateral** reorganizado en tres secciones (Operación / Evidencia /
  Sistema) con grupos colapsables Motor IA e Interceptación.
- **Diversidad de ataque**: la semilla inicial rota entre 10 enfoques distintos
  por ejecución (reconocimiento, auth, control de acceso, inyección, pre-auth,
  lógica de negocio, API/JWT, lado cliente, fuga de información, parsing/caché).
- **Sin barras de scroll** en las pantallas principales (el scroll vive dentro de
  cada panel).

## Pendiente o mejorable

- **Strings CJK internos** sin traducir, heredados del fork, que solo afloran en
  casos de borde: mensajes de cola en coreano (`server/goals.go`), descripción por
  defecto y hecho de origen en chino (`server/server.go`). No son de cara al
  usuario en el flujo normal.
- **Prompt del agente en chino** (`agent/promptcatalog.go`): es el "cerebro" de
  razonamiento, probado en benchmark; `agent/outputlang.go` lo preserva a
  propósito. No traducir sin medir el impacto.
- **Más entropía entre ejecuciones**: el SDK `norma` (v0.4.4) no expone
  `temperature` ni `seed` en `llm.Config`, así que la variación queda en la
  semilla. Para ir más allá habría que forkear el SDK.
- **Pulido de front por pantalla**: se mejoró el menú, el login y el espaciado del
  dashboard. Las demás pantallas (Hallazgos, Detalle de tarea, Tráfico, etc.) y un
  repaso fino de cada gráfica quedan por revisar a fondo.
- **Deduplicación del planner**: no se tocó (no se localizó lógica de "contención"
  que suprima direcciones nuevas). Si la exploración se vuelve repetitiva, revisar
  `agent/planner.go`.

## Build y despliegue

El binario embebe el frontend estático (`-tags embedui`). Flujo completo:

```bash
# 1. Frontend estático
cd web && NEXT_EXPORT=1 npm run build:static && cd ..
rsync -a --delete web/out/ server/webui/dist/

# 2. Binario Go (arm64) con el frontend embebido
docker run --rm -v "$PWD":/src -w /src -e CGO_ENABLED=0 -e GOOS=linux -e GOARCH=arm64 \
  golang:1.26 go build -tags embedui -o dist/arm64/airecon ./cmd/airecon

# 3. Imagen + arranque
docker build --platform linux/arm64 -t airecon:local .
docker compose up -d
```

El frontend se construye con Node de Homebrew en el PATH.

## Configuración y datos

- El **perfil LLM** (proveedor, modelo, clave) y la **contraseña de acceso** viven
  en la base de datos (`llm_profiles` y `settings`), no en el código ni en el
  repo. Recrear el volumen de Postgres los borra: hay que volver a configurarlos
  desde la consola (Sistema → LLM) y el setup inicial de contraseña.
- Con Google AI Studio (Gemini) conviene `workers=1` en `settings` por el límite
  de la capa gratuita (15 peticiones/min).
- El `.env` local (con `POSTGRES_PASSWORD`) y `docker-compose.override.yml` están
  en `.gitignore`: no se publican.

## Notas

- El token de sesión se invalida al recrear el contenedor (la clave de firma del
  JWT se regenera): hay que volver a iniciar sesión, es esperado.
- Licencia AGPL-3.0: se conserva `LICENSE` y el crédito al proyecto original
  ARTEX (Autumn-27), como exige la licencia.
