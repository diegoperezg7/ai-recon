<div align="center">

<img src="web/public/logo.png" width="130" alt="AIRecon" />

# AIRecon

**Consola de pentesting autónomo impulsada por LLM**

Un agente multinivel que explora, razona y verifica vulnerabilidades reales sobre un objetivo web, de principio a fin, y redacta el informe por ti.

[![Licencia: AGPL-3.0](https://img.shields.io/badge/licencia-AGPL--3.0-2563eb.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.26-2563eb.svg)](go.mod)
[![Next.js](https://img.shields.io/badge/Next.js-16-111827.svg)](web/package.json)
[![Docker](https://img.shields.io/badge/Docker-Compose-2563eb.svg)](docker-compose.yml)

</div>

---

## Qué es

AIRecon orquesta varios agentes LLM (planificador y workers) que colaboran para auditar una
aplicación web: reconocen la superficie, deciden qué atacar, ejecutan herramientas, confirman
cada hallazgo contra el objetivo y lo documentan. Todo queda guardado (sesiones, tráfico,
activos, evidencias) y se consulta desde una consola web.

Pensado para **entornos autorizados**: laboratorios, CTF, aplicaciones propias y objetivos con
permiso explícito.

## Características

- **Autónomo de punta a punta.** Del reconocimiento al informe, sin guiar cada paso.
- **Selección de ataque dinámica.** Cada ejecución arranca por una superficie distinta
  (reconocimiento, autenticación, control de acceso, inyección, pre-auth y carga de ficheros,
  lógica de negocio), así dos runs del mismo objetivo no siguen el mismo camino.
- **Verificación real.** Un hallazgo solo se marca como confirmado si se reproduce contra el
  objetivo, con su PoC y su tráfico asociado.
- **Informe listo para entregar.** Markdown y CSV, con severidad, evidencia y activos.
- **Multiproveedor.** Compatible con la API estilo OpenAI: OpenRouter, Google AI Studio
  (Gemini), Groq, Cerebras, Ollama local, y más.
- **Interfaz bilingüe (ES/EN)** con modo claro y oscuro.

## Stack

| Capa | Tecnología |
|------|------------|
| Backend | Go 1.26, binario único con frontend embebido (`-tags embedui`) |
| Frontend | Next.js 16, TypeScript, Tailwind v4, Radix |
| Datos | PostgreSQL 16 |
| Despliegue | Docker Compose |

## Inicio rápido

Requiere Docker.

```bash
# 1. Construir la imagen (frontend export + binario Go)
docker build --platform linux/arm64 -t darc3-artex:local .

# 2. Levantar PostgreSQL + AIRecon
docker compose up -d

# 3. Abrir la consola
open http://127.0.0.1:8787
```

En el primer arranque, la consola pide crear la contraseña de acceso (el usuario es fijo).
La configuración del proveedor LLM se hace desde **Sistema → LLM**.

### Compilación desde el código

El binario embebe el frontend estático. El flujo es: exportar el frontend, copiarlo a
`server/webui/dist/` y compilar el Go con la etiqueta `embedui`.

```bash
cd web && NEXT_EXPORT=1 npm run build:static && cd ..
rsync -a --delete web/out/ server/webui/dist/
CGO_ENABLED=0 go build -tags embedui -o dist/artex ./cmd/artex
```

## Aviso de uso

AIRecon ejecuta acciones ofensivas contra el objetivo que se le indica. Úsalo **solo** sobre
sistemas de tu propiedad o para los que tengas autorización escrita. El uso no autorizado puede
ser ilegal. La responsabilidad es de quien lo opera.

## Licencia y origen

Distribuido bajo **AGPL-3.0** (ver [LICENSE](LICENSE)). AIRecon es un derivado del proyecto de
código abierto [ARTEX](https://github.com/Autumn-27/artex) de Autumn-27, con una capa propia de
interfaz, idioma y selección de ataque. Se conservan los avisos de copyright originales según
exige la licencia.
