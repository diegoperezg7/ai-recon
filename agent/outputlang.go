package agent

// Idioma de salida de cara al usuario. AIRecon fuerza al agente a escribir el
// lenguaje natural que ve el usuario en español o inglés. El "cerebro" del
// agente (el cuerpo de razonamiento en chino, probado en benchmark) se deja
// igual: aquí solo se decide el idioma de la superficie que lee el usuario.

// OutputLanguage es el hook que el servidor cablea para devolver el código de
// idioma de salida ("es"/"en") elegido en la configuración (AIRECON_LANG, KV de
// settings en runtime). Si es nil o vacío, cae al valor por defecto de abajo.
var OutputLanguage func() string

// Idioma de salida por defecto cuando el servidor no cablea el hook (p. ej. al
// testear el paquete agent en solitario). Por defecto, español; el inglés se
// cambia desde la configuración.
const defaultOutputLang = "es"

// SupportedOutputLangs es la lista de códigos de idioma de salida disponibles.
var SupportedOutputLangs = []string{"es", "en"}

// normalizeLang 은 지원 목록 밖의 코드를 기본값으로 좁힌다.
func normalizeLang(code string) string {
	for _, c := range SupportedOutputLangs {
		if c == code {
			return code
		}
	}
	return defaultOutputLang
}

// resolveOutputLang 은 현재 적용할 출력 언어 코드를 돌려준다. 서버가 배선한 훅을 먼저
// 보고, 없으면 기본값을 쓴다.
func resolveOutputLang() string {
	if OutputLanguage != nil {
		if code := OutputLanguage(); code != "" {
			return normalizeLang(code)
		}
	}
	return defaultOutputLang
}

// langPhrase es el nombre con el que se señala "escribe en este idioma" dentro
// de la directiva (redactada en chino). El nombre nativo entre paréntesis evita
// que el modelo confunda de qué idioma se trata.
var langPhrase = map[string]string{
	"en": "英语（English）",
	"es": "español (Español)",
}

// langBare es el nombre corto del idioma para señalarlo de forma repetida en la directiva.
var langBare = map[string]string{
	"en": "英语",
	"es": "español",
}
