import { DEFAULT_LOCALE, isLocale, LOCALE_STORAGE_KEY, type Locale, resolveLocale } from "@/i18n/config";
import { getLocalStorageValue } from "@/lib/local-storage.client";

// Etiqueta de locale BCP47 para fechas y horas. Sigue el idioma de pantalla
// (locale-store), así que se usa esta función en vez de un valor fijo. En el
// render de servidor (sin window) usa el locale por defecto del build; en el
// cliente, el valor que eligió el usuario. Es seguro para los formateadores
// Intl que se leen una vez al cargar el módulo.
const DATE_LOCALE: Record<Locale, string> = {
  es: "es-ES",
  en: "en-US",
};

export function getDateLocale(): string {
  if (typeof window === "undefined") return DATE_LOCALE[DEFAULT_LOCALE];
  const raw = getLocalStorageValue(LOCALE_STORAGE_KEY);
  const loc: Locale = isLocale(raw) ? raw : resolveLocale();
  return DATE_LOCALE[loc];
}
