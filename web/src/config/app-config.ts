import { DEFAULT_LOCALE, type Locale, resolveLocale } from "@/i18n/config";

import packageJson from "../../package.json";

const currentYear = new Date().getFullYear();

// Título de la pestaña del navegador y meta descripción para buscadores. Entra
// como metadata estática en layout.tsx y se renderiza en el <title> y
// <meta name="description"> de todas las páginas, según el locale por defecto
// del build. El título usa dos puntos en vez de raya ("Producto: descripción").
// Esta metadata es un valor estático de build, así que no reacciona al cambio
// de idioma en ejecución (es solo el título de la pestaña, se admite).
const META_BY_LOCALE: Record<Locale, { title: string; description: string }> = {
  en: {
    title: "AIRecon: Consola de Pentesting Autónomo",
    description: "Consola de sistema impulsada por LLM para pruebas de penetración autónomas",
  },
  es: {
    title: "AIRecon: Consola de Pentesting Autónomo",
    description: "Una consola de sistema impulsada por LLM para pruebas de penetración autónomas",
  },
};

export const APP_CONFIG = {
  name: "AIRecon",
  version: packageJson.version,
  copyright: `© ${currentYear}, AIRecon.`,
  meta: META_BY_LOCALE[resolveLocale()] ?? META_BY_LOCALE[DEFAULT_LOCALE],
};
