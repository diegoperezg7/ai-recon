import { getRequestConfig } from "next-intl/server";

import { type Locale, resolveLocale } from "./config";

// Cargadores de mensajes por locale. Mapa de imports explícito para que el
// bundler no arrastre todo el directorio messages/ como módulo de contexto.
const loaders: Record<Locale, () => Promise<{ default: Record<string, unknown> }>> = {
  es: () => import("../../messages/es.json"),
  en: () => import("../../messages/en.json"),
};

// Configuración de petición de next-intl. Como no se usa enrutado por locale
// (segmentos/middleware), aquí se fija el locale directamente. No se consulta
// requestLocale, así que no fuerza render dinámico. Esto solo fija el idioma
// por defecto de los componentes de servidor y el prerender; el cambio en
// ejecución lo gestiona IntlRuntimeProvider en el cliente.
export default getRequestConfig(async () => {
  const locale = resolveLocale();
  const messages = (await loaders[locale]()).default;
  return { locale, messages };
});
