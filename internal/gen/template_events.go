// SPDX-License-Identifier: MIT

package gen

// eventDecodersTemplate generates event decoder functions
const eventDecodersTemplate = `{{/* Generate type-specific decoders for events */}}
{{- range .Contract.Events}}

// Decode decodes log data for {{.Name}} event
func (e *{{.Name}}EventDecoder) Decode(data []byte) ({{.Struct.Name}}, error) {
	return e.decodeImpl(data)
}

// MustDecode decodes log data for {{.Name}} event
func (e *{{.Name}}EventDecoder) MustDecode(data []byte) {{.Struct.Name}} {
	result, err := e.decodeImpl(data)
	if err != nil {
		panic(err)
	}
	return result
}

// decodeImpl contains the actual decode logic
func (e *{{.Name}}EventDecoder) decodeImpl(data []byte) ({{.Struct.Name}}, error) {
	// Decode event parameters (only non-indexed parameters are in data)
	var result {{.Struct.Name}}
	{{- $hasNonIndexedParams := false}}
	{{- range $i, $input := .Inputs}}
	{{- if not $input.Indexed}}
	{{- $hasNonIndexedParams = true}}
	{{- end}}
	{{- end}}
	{{- if $hasNonIndexedParams}}
	{{- $needsVal := false}}
	{{- $needsValUint64 := false}}
	{{- $needsValInt64 := false}}
	{{- $needsValAddr := false}}
	{{- $needsValBool := false}}
	{{- $needsValString := false}}
	{{- $needsValBytes := false}}
	{{- range .Inputs}}
		{{- if not .Indexed}}
			{{- if eq .Type.TypeName "*big.Int"}}
				{{- $needsVal = true}}
			{{- end}}
			{{- if eq .Type.TypeName "uint64"}}
				{{- $needsValUint64 = true}}
			{{- end}}
			{{- if eq .Type.TypeName "int64"}}
				{{- $needsValInt64 = true}}
			{{- end}}
			{{- if eq .Type.TypeName "Address"}}
				{{- $needsValAddr = true}}
			{{- end}}
			{{- if eq .Type.TypeName "bool"}}
				{{- $needsValBool = true}}
			{{- end}}
			{{- if eq .Type.TypeName "string"}}
				{{- $needsValString = true}}
			{{- end}}
			{{- if eq .Type.TypeName "[]byte"}}
				{{- $needsValBytes = true}}
			{{- end}}
		{{- end}}
	{{- end}}
	{{- if $needsVal}}
	var val *big.Int
	{{- end}}
	{{- if $needsValUint64}}
	var valUint64 uint64
	{{- end}}
	{{- if $needsValInt64}}
	var valInt64 int64
	{{- end}}
	{{- if $needsValAddr}}
	var valAddr Address
	{{- end}}
	{{- if $needsValBool}}
	var valBool bool
	{{- end}}
	{{- if $needsValString}}
	var valString string
	{{- end}}
	{{- if $needsValBytes}}
	var valBytes []byte
	{{- end}}
	var err error
	offset := 0
	{{- range $i, $input := .Inputs}}
	{{- if not $input.Indexed}}
	{{- if eq $input.Type.TypeName "*big.Int"}}
	if len(data) < offset+32 {
		return result, errors.New("insufficient data for event parameter {{$input.Name}}")
	}
	{{- if $input.Type.IsSigned}}
	val, err = decodeInt256(data[offset:offset+32])
	if err != nil {
		return result, fmt.Errorf("decoding event parameter {{$input.Name}}: %w", err)
	}
	result.{{$input.Name | title}} = val
	{{- else}}
	val, err = decodeUint256(data[offset:offset+32])
	if err != nil {
		return result, fmt.Errorf("decoding event parameter {{$input.Name}}: %w", err)
	}
	result.{{$input.Name | title}} = val
	{{- end}}
	offset += 32
	{{- else if eq $input.Type.TypeName "uint64"}}
	if len(data) < offset+32 {
		return result, errors.New("insufficient data for event parameter {{$input.Name}}")
	}
	valUint64, err = decodeUint64(data[offset:offset+32])
	if err != nil {
		return result, fmt.Errorf("decoding event parameter {{$input.Name}}: %w", err)
	}
	result.{{$input.Name | title}} = valUint64
	offset += 32
	{{- else if eq $input.Type.TypeName "int64"}}
	if len(data) < offset+32 {
		return result, errors.New("insufficient data for event parameter {{$input.Name}}")
	}
	valInt64, err = decodeInt64(data[offset:offset+32])
	if err != nil {
		return result, fmt.Errorf("decoding event parameter {{$input.Name}}: %w", err)
	}
	result.{{$input.Name | title}} = valInt64
	offset += 32
	{{- else if eq $input.Type.TypeName "bool"}}
	if len(data) < offset+32 {
		return result, errors.New("insufficient data for event parameter {{$input.Name}}")
	}
	valBool, err = decodeBool(data[offset:offset+32])
	if err != nil {
		return result, fmt.Errorf("decoding event parameter {{$input.Name}}: %w", err)
	}
	result.{{$input.Name | title}} = valBool
	offset += 32
	{{- else if eq $input.Type.TypeName "Address"}}
	if len(data) < offset+32 {
		return result, errors.New("insufficient data for event parameter {{$input.Name}}")
	}
	valAddr, err = decodeAddress(data[offset:offset+32])
	if err != nil {
		return result, fmt.Errorf("decoding event parameter {{$input.Name}}: %w", err)
	}
	result.{{$input.Name | title}} = valAddr
	offset += 32
	{{- else if eq $input.Type.TypeName "string"}}
	var nextOffset int
	valString, nextOffset, err = decodeString(data, offset)
	if err != nil {
		return result, fmt.Errorf("decoding event parameter {{$input.Name}}: %w", err)
	}
	result.{{$input.Name | title}} = valString
	offset = nextOffset
	{{- else if eq $input.Type.TypeName "[]byte"}}
	var nextOffset int
	valBytes, nextOffset, err = decodeBytes(data, offset)
	if err != nil {
		return result, fmt.Errorf("decoding event parameter {{$input.Name}}: %w", err)
	}
	result.{{$input.Name | title}} = valBytes
	offset = nextOffset
	{{- else}}
	return result, errors.New("unsupported event parameter type: {{$input.Type.TypeName}}")
	{{- end}}
	{{- end}}
	{{- end}}
	{{- else}}
	// Event has no non-indexed parameters, return empty struct
	{{- end}}
	return result, nil
}
{{- end}}`

// eventRegistryTemplate generates the event registry and event types
const eventRegistryTemplate = `{{- range .Contract.Events}}
// {{.Name | title}}EventDecoder returns a decoder for {{.Name}} events
func (er EventRegistry) {{.Name | title}}EventDecoder() *{{.Name}}EventDecoder {
	return &{{.Name}}EventDecoder{
		PackableEvent: PackableEvent{
			Name:  {{.Name | quote}},
			Topic: HashFromHex({{printf "0x%x" .Topic.Bytes | quote}}),
		},
	}
}
{{- end}}

// Events returns the event registry
func Events() EventRegistry {
	return EventRegistry{}
}

{{/* Generate specific event decoder types */}}
{{- range .Contract.Events}}

// {{.Name | title}}EventDecoder represents the {{.Name}} event with type-safe decode functionality
type {{.Name | title}}EventDecoder struct {
	PackableEvent
}
{{- end}}`