// Locales soportados y valor por defecto. AIRecon usa español por defecto y
// permite cambiar a inglés desde los ajustes de usuario.
export const LOCALES = ["es", "en"] as const;
export type Locale = (typeof LOCALES)[number];

// El idioma por defecto es el español. El inglés se selecciona desde ajustes.
// En build se puede sobrescribir con NEXT_PUBLIC_LOCALE; en ejecución manda el
// valor que elige el usuario (cookie / localStorage).
export const DEFAULT_LOCALE: Locale = "es";

// Nombre de cada locale en su propia lengua para el selector de idioma. No se
// traduce: cada idioma se lee con su propio nombre, se vea en el idioma que se vea.
export const LOCALE_LABELS: Record<Locale, string> = {
  es: "Español",
  en: "English",
};

// Clave donde se guarda el idioma elegido. La pantalla lee el valor de
// localStorage; también se escribe en cookie por si el servidor lo necesita.
export const LOCALE_COOKIE = "airecon_locale";
export const LOCALE_STORAGE_KEY = "airecon_locale";

// Guarda que estrecha un string cualquiera a un locale soportado.
export function isLocale(raw: unknown): raw is Locale {
  return typeof raw === "string" && (LOCALES as readonly string[]).includes(raw);
}

// Resuelve el locale activo en tiempo de build. Debe ser compatible con la
// exportación estática (output: "export"), así que no usa cookies()/headers(),
// solo variables de entorno. Si NEXT_PUBLIC_LOCALE está vacío o fuera de la
// lista soportada, cae al valor por defecto. Este valor solo afecta al HTML
// inicial del prerender; el idioma real lo fija el cliente vía cookie.
export function resolveLocale(): Locale {
  const raw = process.env.NEXT_PUBLIC_LOCALE;
  return isLocale(raw) ? raw : DEFAULT_LOCALE;
}
