package report

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"path"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/diegoperezg7/ai-recon/db"
)

// Renderizado para el "exportar" de la página de vulnerabilidades: convierte
// las filas de la tabla de findings en Markdown resumen, Markdown individual o
// CSV. El JSON lo serializa la capa server directamente como DTO, aquí no.

// sortFindingsForExport ordena por severidad descendente y luego por fecha descendente; igual que la agrupación del informe resumen.
func sortFindingsForExport(fs []*db.DBFinding) {
	sort.SliceStable(fs, func(i, j int) bool {
		ri, rj := sevRank[fs[i].Severity], sevRank[fs[j].Severity]
		if ri != rj {
			return ri < rj // en sevRank, cuanto menor el valor, más grave
		}
		return fs[i].CreatedAt.After(fs[j].CreatedAt)
	})
}

// findingTitle elige el título legible de la vulnerabilidad: nombre -> tipo -> "Sin clasificar".
func findingTitle(f *db.DBFinding) string {
	return nz(f.Name, nz(f.VulnClass, "Sin clasificar"))
}

// FindingsMarkdown integra varias findings en un único informe resumen (resumen
// + agrupación por severidad; cada una incluye tipo/estado/tarea/evidencia/informe).
func FindingsMarkdown(fs []*db.DBFinding, generatedAt time.Time) string {
	items := append([]*db.DBFinding(nil), fs...)
	sortFindingsForExport(items)

	var b strings.Builder
	b.WriteString("# Informe resumen de vulnerabilidades\n\n")
	fmt.Fprintf(&b, "- **Fecha de generación**: %s\n", generatedAt.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(&b, "- **Total de vulnerabilidades**: %d\n\n", len(items))

	// Resumen: recuento por severidad.
	counts := map[string]int{}
	for _, f := range items {
		counts[f.Severity]++
	}
	b.WriteString("## Resumen\n\n")
	b.WriteString("| Severidad | Nº |\n| --- | --- |\n")
	for _, s := range []struct{ key, label string }{
		{"critical", "Crítica"}, {"high", "Alta"}, {"medium", "Media"}, {"low", "Baja"},
	} {
		fmt.Fprintf(&b, "| %s | %d |\n", s.label, counts[s.key])
	}
	b.WriteString("\n")

	if len(items) == 0 {
		b.WriteString("_No hay vulnerabilidades coincidentes._\n")
		return b.String()
	}

	b.WriteString("## Detalle de vulnerabilidades\n\n")
	for i, f := range items {
		fmt.Fprintf(&b, "### %d. [%s] %s\n\n", i+1, strings.ToUpper(nz(f.Severity, "info")), findingTitle(f))
		if f.VulnClass != "" {
			fmt.Fprintf(&b, "- **Tipo**: %s\n", f.VulnClass)
		}
		fmt.Fprintf(&b, "- **Estado**: %s\n", nz(f.Status, "pending"))
		if desc := strings.TrimSpace(f.TaskDescription); desc != "" {
			fmt.Fprintf(&b, "- **Tarea**: %s\n", desc)
		}
		fmt.Fprintf(&b, "- **Fecha de detección**: %s\n\n", f.CreatedAt.Format("2006-01-02 15:04:05"))
		if s := strings.TrimSpace(f.Summary); s != "" {
			fmt.Fprintf(&b, "%s\n\n", s)
		}
		if e := strings.TrimSpace(f.Evidence); e != "" {
			fmt.Fprintf(&b, "**Evidencia:**\n\n```\n%s\n```\n\n", e)
		}
		if rep := strings.TrimSpace(f.Report); rep != "" {
			b.WriteString("**Informe detallado:**\n\n")
			b.WriteString(rep)
			b.WriteString("\n\n")
		}
		b.WriteString(findingTrafficMarkdown(f, false))
		b.WriteString("---\n\n")
	}
	return b.String()
}

// SingleFindingMarkdown renderiza una vulnerabilidad como un Markdown independiente (para el paquete "un fichero por vulnerabilidad").
func SingleFindingMarkdown(f *db.DBFinding, generatedAt time.Time) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# [%s] %s\n\n", strings.ToUpper(nz(f.Severity, "info")), findingTitle(f))
	if f.VulnClass != "" {
		fmt.Fprintf(&b, "- **Tipo**: %s\n", f.VulnClass)
	}
	fmt.Fprintf(&b, "- **Severidad**: %s\n", nz(f.Severity, "info"))
	fmt.Fprintf(&b, "- **Estado**: %s\n", nz(f.Status, "pending"))
	if desc := strings.TrimSpace(f.TaskDescription); desc != "" {
		fmt.Fprintf(&b, "- **Tarea**: %s\n", desc)
	}
	fmt.Fprintf(&b, "- **Fecha de detección**: %s\n", f.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(&b, "- **Fecha de generación**: %s\n\n", generatedAt.Format("2006-01-02 15:04:05"))
	if s := strings.TrimSpace(f.Summary); s != "" {
		fmt.Fprintf(&b, "## Resumen\n\n%s\n\n", s)
	}
	if e := strings.TrimSpace(f.Evidence); e != "" {
		fmt.Fprintf(&b, "## Evidencia\n\n```\n%s\n```\n\n", e)
	}
	if rep := strings.TrimSpace(f.Report); rep != "" {
		b.WriteString("## Informe detallado\n\n")
		b.WriteString(rep)
		b.WriteString("\n")
	}
	b.WriteString(findingTrafficMarkdown(f, true))
	return b.String()
}

var unsafeFilenameChars = regexp.MustCompile(`[^\p{Han}\p{L}\p{N}._-]+`)

// FindingFilename genera un nombre de fichero .md seguro para "un fichero por
// vulnerabilidad". Quita separadores de ruta y caracteres de control para evitar rutas anómalas dentro del zip.
func FindingFilename(f *db.DBFinding) string {
	sev := nz(f.Severity, "info")
	title := findingTitle(f)
	name := fmt.Sprintf("%s_%s_#%d", sev, title, f.ID)
	name = unsafeFilenameChars.ReplaceAllString(name, "_")
	name = strings.Trim(name, "._")
	if name == "" {
		name = fmt.Sprintf("finding_%d", f.ID)
	}
	// Defensa: se vuelve a reducir la ruta para bloquear zip slip.
	name = path.Base(name)
	if len(name) > 120 {
		name = name[:120]
	}
	return name + ".md"
}

// FindingsCSV renderiza varias findings como CSV (con BOM UTF-8 para que Excel
// reconozca bien los acentos). No incluye el report/evidence completo, solo
// campos de tipo resumen; para el texto completo se exporta a Markdown/JSON.
func FindingsCSV(fs []*db.DBFinding) []byte {
	items := append([]*db.DBFinding(nil), fs...)
	sortFindingsForExport(items)

	var buf bytes.Buffer
	buf.WriteString("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(&buf)
	_ = w.Write([]string{"ID", "Nombre", "Tipo", "Severidad", "Estado", "Tarea", "Fecha de detección", "Resumen", "Nº evidencias de tráfico", "IDs evidencias de tráfico"})
	for _, f := range items {
		_ = w.Write([]string{
			fmt.Sprintf("%d", f.ID),
			findingTitle(f),
			f.VulnClass,
			nz(f.Severity, "info"),
			nz(f.Status, "pending"),
			f.TaskDescription,
			f.CreatedAt.Format("2006-01-02 15:04:05"),
			strings.TrimSpace(f.Summary),
			fmt.Sprint(len(f.TrafficBindings)), findingTrafficIDs(f),
		})
	}
	w.Flush()
	return buf.Bytes()
}

func findingTrafficIDs(f *db.DBFinding) string {
	ids := make([]string, 0, len(f.TrafficBindings))
	for _, b := range f.TrafficBindings {
		ids = append(ids, fmt.Sprint(b.ID))
	}
	return strings.Join(ids, ",")
}

func findingTrafficMarkdown(f *db.DBFinding, attachments bool) string {
	stale := f.Report != "" && f.EvidenceVersion != f.ReportEvidenceVersion
	if len(f.TrafficBindings) == 0 && !stale {
		return ""
	}
	var out strings.Builder
	out.WriteString("\n## Evidencia de tráfico relacionada\n\n")
	fmt.Fprintf(&out, "Versión de evidencia: %d, nº de vínculos: %d.\n\n", f.EvidenceVersion, len(f.TrafficBindings))
	if stale {
		out.WriteString("La evidencia cambió; hay que actualizar el informe detallado.\n\n")
	}
	for i, b := range f.TrafficBindings {
		fmt.Fprintf(&out, "%d. **Evidencia #%d · %s**: `%s %s`, código de estado %d\n", i+1, b.ID, b.Role, b.Snapshot.Method, strings.ReplaceAll(b.Snapshot.URL, "`", "%60"), b.Snapshot.Status)
		if b.Note != "" {
			fmt.Fprintf(&out, "   %s\n", strings.ReplaceAll(b.Note, "\n", "\n   "))
		}
		if attachments {
			fmt.Fprintf(&out, "   [Mensaje de petición](evidence/%d/%d/request.http) · [Mensaje de respuesta](evidence/%d/%d/response.http)\n", f.ID, b.ID, f.ID, b.ID)
		}
	}
	out.WriteString("\n")
	return out.String()
}
