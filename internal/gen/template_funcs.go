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
		"isDynamic":    isDynamicGoType,
	}
}

// isDynamicGoType reports whether a Go type corresponds to a dynamically-sized
// ABI type (i.e., one that uses a 32-byte offset pointer in head-tail encoding).
// Dynamic types: string, []byte, and any slice (IsSlice == true).
// Static types: bool, uint*, int*, address, hash, bytes1-32, and fixed arrays.
// NOTE: nested struct types are treated as static here; if a nested struct
// itself contains dynamic fields, the nested decode function handles the
// internal head-tail layout.
func isDynamicGoType(goType types.GoType) bool {
	if goType.IsSlice {
		return true
	}
	switch goType.TypeName {
	case "string", "[]byte":
		return true
	}
	return false
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