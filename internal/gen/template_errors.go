// SPDX-License-Identifier: MIT

package gen

// errorDecodersTemplate generates error decoder functions. Parameters are
// decoded via the shared "decodeValue" template.
const errorDecodersTemplate = `{{/* Generate type-specific decoders for errors */}}
{{- range .Contract.Errors}}

// Decode decodes error data for {{.Name}} error
func (e *{{.Name}}ErrorDecoder) Decode(data []byte) ({{.Struct.Name}}, error) {
	return e.decodeImpl(data)
}

// MustDecode decodes error data for {{.Name}} error
func (e *{{.Name}}ErrorDecoder) MustDecode(data []byte) {{.Struct.Name}} {
	result, err := e.decodeImpl(data)
	if err != nil {
		panic(err)
	}
	return result
}

// decodeImpl contains the actual decode logic
func (e *{{.Name}}ErrorDecoder) decodeImpl(data []byte) ({{.Struct.Name}}, error) {
	// Skip the 4-byte selector
	if len(data) < 4 {
		return {{.Struct.Name}}{}, errors.New("insufficient data for error selector")
	}
	var result {{.Struct.Name}}
	{{- $allSupported := true}}
	{{- range .Inputs}}{{- if not (decoderSupports .Type)}}{{- $allSupported = false}}{{- end}}{{- end}}
	{{- if not $allSupported}}
	return result, errors.New("unsupported parameter type in {{.Name}} error")
	{{- else}}
	{{- if gt (len .Inputs) 0}}
	errorData := data[4:]
	offset := 0
	{{- range $i, $input := .Inputs}}
	{{- template "decodeValue" (dict "Data" "errorData" "Field" (printf "result.%s" ($input.Name | title)) "T" $input.Type "Ctx" (printf "error parameter %s" $input.Name) "Zero" "result")}}
	{{- end}}
	{{- end}}
	return result, nil
	{{- end}}
}
{{- end}}`

// errorRegistryTemplate generates the error registry and error types
const errorRegistryTemplate = `{{- range .Contract.Errors}}
// {{.Name}}Error returns a packable error for {{.Name}}
func (er ErrorRegistry) {{.Name}}Error() *{{.Name}}ErrorDecoder {
	return &{{.Name}}ErrorDecoder{
		PackableError: PackableError{
			Name:      {{.Name | quote}},
			Signature: {{.Signature | quote}},
			Selector:  HexData({{.Selector.Hex | quote}}),
		},
	}
}
{{- end}}

// Errors returns the error registry
func Errors() ErrorRegistry {
	return ErrorRegistry{}
}

{{/* Generate specific error decoder types */}}
{{- range .Contract.Errors}}

// {{.Name}}ErrorDecoder represents the {{.Name}} error with type-safe decode functionality
type {{.Name}}ErrorDecoder struct {
	PackableError
}
{{- end}}`
