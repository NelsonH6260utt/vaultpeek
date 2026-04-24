package vault

import (
	"fmt"
	"regexp"
	"strings"
)

// TemplateResult holds the outcome of a secret template rendering operation.
type TemplateResult struct {
	Path     string
	Rendered map[string]string
	Missing  []string
	Success  bool
	Err      error
}

// templateVarRe matches placeholders like {{secret.key}}
var templateVarRe = regexp.MustCompile(`\{\{([^}]+)\}\}`)

// RenderTemplate replaces {{key}} placeholders in a template string with
// values from the provided secret data map. Keys not found in data are
// collected in TemplateResult.Missing and left unreplaced.
func RenderTemplate(path, tmpl string, data map[string]interface{}) TemplateResult {
	if data == nil {
		return TemplateResult{
			Path:    path,
			Success: false,
			Err:     fmt.Errorf("secret data is nil for path %q", path),
		}
	}

	rendered := make(map[string]string)
	missing := []string{}

	result := templateVarRe.ReplaceAllStringFunc(tmpl, func(match string) string {
		key := strings.TrimSpace(match[2 : len(match)-2])
		if val, ok := data[key]; ok {
			s := fmt.Sprintf("%v", val)
			rendered[key] = s
			return s
		}
		missing = append(missing, key)
		return match
	})

	_ = result
	rendered["__output"] = result

	return TemplateResult{
		Path:     path,
		Rendered: rendered,
		Missing:  missing,
		Success:  len(missing) == 0,
	}
}
