// SPDX-License-Identifier: MIT

package gen

import (
	"strconv"
	"strings"
	"text/template"

	"github.com/otherview/solgen/internal/types"
)

// TemplateData holds data for template rendering
type TemplateData struct {
	Contract *types.Contract
	Imports  []string
}

// templateFuncs returns template helper functions
func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"formatGoType": formatGoType,
		"quote":        strconv.Quote,
		"lower":        strings.ToLower,
		"title":        titleCase,
		"join":         strings.Join,
		"add":          func(a, b int) int { return a + b },
		"default":      func(def, val string) string { if val == "" { return def }; return val },
		"hasPrefix":    strings.HasPrefix,
	}
}

// formatGoType formats a GoType for use in generated code
func formatGoType(goType types.GoType) string {
	return goType.TypeName
}

// titleCase provides a simple title case conversion
func titleCase(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}