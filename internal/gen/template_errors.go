// SPDX-License-Identifier: MIT

package gen

// errorDecodersTemplate generates error decoder functions
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
	errorData := data[4:]
	// Decode error parameters
	var result {{.Struct.Name}}
{{- if gt (len .Inputs) 0}}
	var err error
	offset := 0
	{{- range $i, $input := .Inputs}}
	{{- if eq $input.Type.TypeName "*big.Int"}}
	if len(errorData) < offset+32 {
		return result, errors.New("insufficient data for error parameter {{$input.Name}}")
	}
	{{- if $input.Type.IsSigned}}
	val{{$i}}, err := decodeInt256(errorData[offset:offset+32])
	if err != nil {
		return result, fmt.Errorf("decoding error parameter {{$input.Name}}: %w", err)
	}
	result.{{$input.Name | title}} = val{{$i}}
	{{- else}}
	val{{$i}}, err := decodeUint256(errorData[offset:offset+32])
	if err != nil {
		return result, fmt.Errorf("decoding error parameter {{$input.Name}}: %w", err)
	}
	result.{{$input.Name | title}} = val{{$i}}
	{{- end}}
	offset += 32
	{{- else if eq $input.Type.TypeName "uint64"}}
	if len(errorData) < offset+32 {
		return result, errors.New("insufficient data for error parameter {{$input.Name}}")
	}
	val{{$i}}, err := decodeUint64(errorData[offset:offset+32])
	if err != nil {
		return result, fmt.Errorf("decoding error parameter {{$input.Name}}: %w", err)
	}
	result.{{$input.Name | title}} = val{{$i}}
	offset += 32
	{{- else if eq $input.Type.TypeName "int64"}}
	if len(errorData) < offset+32 {
		return result, errors.New("insufficient data for error parameter {{$input.Name}}")
	}
	val{{$i}}, err := decodeInt64(errorData[offset:offset+32])
	if err != nil {
		return result, fmt.Errorf("decoding error parameter {{$input.Name}}: %w", err)
	}
	result.{{$input.Name | title}} = val{{$i}}
	offset += 32
	{{- else if eq $input.Type.TypeName "bool"}}
	if len(errorData) < offset+32 {
		return result, errors.New("insufficient data for error parameter {{$input.Name}}")
	}
	val{{$i}}, err := decodeBool(errorData[offset:offset+32])
	if err != nil {
		return result, fmt.Errorf("decoding error parameter {{$input.Name}}: %w", err)
	}
	result.{{$input.Name | title}} = val{{$i}}
	offset += 32
	{{- else if eq $input.Type.TypeName "Address"}}
	if len(errorData) < offset+32 {
		return result, errors.New("insufficient data for error parameter {{$input.Name}}")
	}
	val{{$i}}, err := decodeAddress(errorData[offset:offset+32])
	if err != nil {
		return result, fmt.Errorf("decoding error parameter {{$input.Name}}: %w", err)
	}
	result.{{$input.Name | title}} = val{{$i}}
	offset += 32
	{{- else if eq $input.Type.TypeName "string"}}
	val{{$i}}, nextOffset, err := decodeString(errorData, offset)
	if err != nil {
		return result, fmt.Errorf("decoding error parameter {{$input.Name}}: %w", err)
	}
	result.{{$input.Name | title}} = val{{$i}}
	offset = nextOffset
	{{- else if eq $input.Type.TypeName "[]byte"}}
	val{{$i}}, nextOffset, err := decodeBytes(errorData, offset)
	if err != nil {
		return result, fmt.Errorf("decoding error parameter {{$input.Name}}: %w", err)
	}
	result.{{$input.Name | title}} = val{{$i}}
	offset = nextOffset
	{{- else}}
	return result, errors.New("unsupported error parameter type: {{$input.Type.TypeName}}")
	{{- end}}
	{{- end}}
{{- end}}
	return result, nil
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