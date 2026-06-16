// SPDX-License-Identifier: MIT

package gen

// methodRegistryTemplate generates the method registry and method types
const methodRegistryTemplate = `{{- range .Contract.Methods}}
// {{.Name | title}}Method returns a packable method for {{.Name}}
func (mr MethodRegistry) {{.Name | title}}Method() *{{.Name | title}}Method {
	return &{{.Name | title}}Method{
		PackableMethod: PackableMethod{
			Name:      {{.Name | quote}},
			Signature: {{.Signature | quote}},
			Selector:  HexData({{.Selector.Hex | quote}}),
		},
	}
}
{{- end}}

// Methods returns the method registry
func Methods() MethodRegistry {
	return MethodRegistry{}
}

{{/* Generate specific method types */}}
{{- range .Contract.Methods}}

// {{.Name | title}}Method represents the {{.Name}} method with type-safe decode functionality
type {{.Name | title}}Method struct {
	PackableMethod
}
{{- end}}`

// methodDecodersTemplate generates method decode functions. Both the single- and
// multi-return decoders delegate to the shared "decodeValue" template.
const methodDecodersTemplate = `{{/* Generate type-specific decoders for methods */}}
{{- range .Contract.Methods}}
{{- if gt (len .Outputs) 0}}

// Decode decodes return values for {{.Name}} method
func (m *{{.Name | title}}Method) Decode(data []byte) ({{if eq (len .Outputs) 1}}{{$output := index .Outputs 0}}{{formatGoType $output.Type}}{{else}}{{.Name | title}}Result{{end}}, error) {
	return m.decodeImpl(data)
}

// MustDecode decodes return values for {{.Name}} method
func (m *{{.Name | title}}Method) MustDecode(data []byte) {{if eq (len .Outputs) 1}}{{$output := index .Outputs 0}}{{formatGoType $output.Type}}{{else}}{{.Name | title}}Result{{end}} {
	result, err := m.decodeImpl(data)
	if err != nil {
		panic(err)
	}
	return result
}

// decodeImpl contains the actual decode logic
func (m *{{.Name | title}}Method) decodeImpl(data []byte) ({{if eq (len .Outputs) 1}}{{$output := index .Outputs 0}}{{formatGoType $output.Type}}{{else}}{{.Name | title}}Result{{end}}, error) {
{{- if eq (len .Outputs) 1}}
	{{- $output := index .Outputs 0}}
	var result {{formatGoType $output.Type}}
	{{- if decoderSupports $output.Type}}
	offset := 0
	{{- template "decodeValue" (dict "Data" "data" "Field" "result" "T" $output.Type "Ctx" "return value" "Zero" "result")}}
	return result, nil
	{{- else}}
	return result, errors.New("unsupported return type: {{$output.Type.TypeName}}")
	{{- end}}
{{- else}}
	// Multiple return values - return as struct
	var result {{.Name | title}}Result
	{{- $allSupported := true}}
	{{- range .Outputs}}{{- if not (decoderSupports .Type)}}{{- $allSupported = false}}{{- end}}{{- end}}
	{{- if $allSupported}}
	offset := 0
	{{- range $i, $output := .Outputs}}
	{{- template "decodeValue" (dict "Data" "data" "Field" (printf "result.%s" ($output.Name | title)) "T" $output.Type "Ctx" (printf "return value %d" $i) "Zero" "result")}}
	{{- end}}
	return result, nil
	{{- else}}
	return result, errors.New("unsupported return type in {{.Name}}")
	{{- end}}
{{- end}}
}
{{- end}}
{{- end}}`
