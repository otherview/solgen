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
		"staticWordDecoder": staticWordDecoder,
		"decoderSupports":   decoderSupports,
		"dict":              dict,
	}
}

// dict builds a map from alternating key/value arguments, for passing multiple
// named values to a {{template}} invocation.
func dict(pairs ...any) map[string]any {
	m := make(map[string]any, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		key, _ := pairs[i].(string)
		m[key] = pairs[i+1]
	}
	return m
}

// decoderSupports reports whether the shared "decodeValue" template can decode a
// value of the given type from an ABI data section. Types it cannot handle get a
// single "unsupported" backstop instead.
func decoderSupports(t types.GoType) bool {
	switch t.TypeName {
	case "*big.Int", "uint8", "uint16", "uint32", "uint64",
		"int8", "int16", "int32", "int64",
		"bool", "Address", "Hash", "[1]byte", "[32]byte",
		"string", "[]byte":
		return true
	}
	if t.IsStruct {
		return true
	}
	if t.IsSlice && t.ElemType != nil && t.ElemType.IsStruct {
		return true
	}
	if t.IsSlice && staticWordDecoder(t.ElemType) != "" {
		return true
	}
	if t.IsArray && staticWordDecoder(t.ElemType) != "" {
		return true
	}
	return false
}

// staticWordDecoder returns the name of a generated decoder func of signature
// func([]byte) (T, error) that decodes a single static 32-byte ABI word into the
// given element type. It returns "" when the element type has no such decoder
// (e.g. dynamic types, structs, or int sizes below 64 bits), in which case the
// caller emits an "unsupported" backstop.
func staticWordDecoder(elem *types.GoType) string {
	if elem == nil {
		return ""
	}
	switch elem.TypeName {
	case "*big.Int":
		if elem.IsSigned {
			return "decodeInt256"
		}
		return "decodeUint256"
	case "uint8":
		return "decodeUint8"
	case "uint16":
		return "decodeUint16"
	case "uint32":
		return "decodeUint32"
	case "uint64":
		return "decodeUint64"
	case "bool":
		return "decodeBool"
	case "Address":
		return "decodeAddress"
	case "Hash":
		return "decodeHash"
	case "[1]byte":
		return "decodeBytes1"
	case "[32]byte":
		return "decodeBytes32"
	}
	return ""
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