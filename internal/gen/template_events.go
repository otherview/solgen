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

// DecodeLog decodes a full log into {{.Struct.Name}}: non-indexed parameters are
// read from data, while indexed parameters are read from topics (topics[0] is the
// event signature, so indexed params start at topics[1]) in ABI order.
//
// Indexed dynamic types (string, bytes, arrays) are stored as the keccak hash of
// their value in the topic and cannot be recovered to their original value; such
// fields are left at their zero value.
func (e *{{.Name}}EventDecoder) DecodeLog(topics [][32]byte, data []byte) ({{.Struct.Name}}, error) {
	// Non-indexed parameters live in the data section.
	result, err := e.decodeImpl(data)
	if err != nil {
		return result, err
	}
	{{- $hasIndexed := false}}
	{{- range .Inputs}}{{- if .Indexed}}{{- $hasIndexed = true}}{{- end}}{{- end}}
	{{- if $hasIndexed}}
	// Indexed parameters live in the log topics, after the signature topic.
	topicIndex := 1
	{{- range $i, $input := .Inputs}}
	{{- if $input.Indexed}}
	if len(topics) <= topicIndex {
		return result, errors.New("missing topic for indexed parameter {{$input.Name}}")
	}
	{
		word := topics[topicIndex][:]
		{{- if eq $input.Type.TypeName "*big.Int"}}
		{{- if $input.Type.IsSigned}}
		v, e := decodeInt256(word)
		{{- else}}
		v, e := decodeUint256(word)
		{{- end}}
		if e != nil {
			return result, fmt.Errorf("decoding indexed parameter {{$input.Name}}: %w", e)
		}
		result.{{$input.Name | title}} = v
		{{- else if eq $input.Type.TypeName "bool"}}
		v, e := decodeBool(word)
		if e != nil {
			return result, fmt.Errorf("decoding indexed parameter {{$input.Name}}: %w", e)
		}
		result.{{$input.Name | title}} = v
		{{- else if eq $input.Type.TypeName "uint8"}}
		v, e := decodeUint8(word)
		if e != nil {
			return result, fmt.Errorf("decoding indexed parameter {{$input.Name}}: %w", e)
		}
		result.{{$input.Name | title}} = v
		{{- else if eq $input.Type.TypeName "uint16"}}
		v, e := decodeUint16(word)
		if e != nil {
			return result, fmt.Errorf("decoding indexed parameter {{$input.Name}}: %w", e)
		}
		result.{{$input.Name | title}} = v
		{{- else if eq $input.Type.TypeName "uint32"}}
		v, e := decodeUint32(word)
		if e != nil {
			return result, fmt.Errorf("decoding indexed parameter {{$input.Name}}: %w", e)
		}
		result.{{$input.Name | title}} = v
		{{- else if eq $input.Type.TypeName "uint64"}}
		v, e := decodeUint64(word)
		if e != nil {
			return result, fmt.Errorf("decoding indexed parameter {{$input.Name}}: %w", e)
		}
		result.{{$input.Name | title}} = v
		{{- else if eq $input.Type.TypeName "int8"}}
		raw, e := decodeInt64(word)
		if e != nil {
			return result, fmt.Errorf("decoding indexed parameter {{$input.Name}}: %w", e)
		}
		result.{{$input.Name | title}} = int8(raw)
		{{- else if eq $input.Type.TypeName "int16"}}
		raw, e := decodeInt64(word)
		if e != nil {
			return result, fmt.Errorf("decoding indexed parameter {{$input.Name}}: %w", e)
		}
		result.{{$input.Name | title}} = int16(raw)
		{{- else if eq $input.Type.TypeName "int32"}}
		raw, e := decodeInt64(word)
		if e != nil {
			return result, fmt.Errorf("decoding indexed parameter {{$input.Name}}: %w", e)
		}
		result.{{$input.Name | title}} = int32(raw)
		{{- else if eq $input.Type.TypeName "int64"}}
		v, e := decodeInt64(word)
		if e != nil {
			return result, fmt.Errorf("decoding indexed parameter {{$input.Name}}: %w", e)
		}
		result.{{$input.Name | title}} = v
		{{- else if eq $input.Type.TypeName "Address"}}
		v, e := decodeAddress(word)
		if e != nil {
			return result, fmt.Errorf("decoding indexed parameter {{$input.Name}}: %w", e)
		}
		result.{{$input.Name | title}} = v
		{{- else if eq $input.Type.TypeName "Hash"}}
		v, e := decodeHash(word)
		if e != nil {
			return result, fmt.Errorf("decoding indexed parameter {{$input.Name}}: %w", e)
		}
		result.{{$input.Name | title}} = v
		{{- else if eq $input.Type.TypeName "[1]byte"}}
		v, e := decodeBytes1(word)
		if e != nil {
			return result, fmt.Errorf("decoding indexed parameter {{$input.Name}}: %w", e)
		}
		result.{{$input.Name | title}} = v
		{{- else if eq $input.Type.TypeName "[32]byte"}}
		v, e := decodeBytes32(word)
		if e != nil {
			return result, fmt.Errorf("decoding indexed parameter {{$input.Name}}: %w", e)
		}
		result.{{$input.Name | title}} = v
		{{- else}}
		// Indexed {{$input.Type.TypeName}} is stored as a keccak hash in the topic
		// and cannot be recovered; leaving {{$input.Name | title}} at its zero value.
		_ = word
		{{- end}}
	}
	topicIndex++
	{{- end}}
	{{- end}}
	{{- end}}
	return result, nil
}

// decodeImpl contains the actual decode logic
func (e *{{.Name}}EventDecoder) decodeImpl(data []byte) ({{.Struct.Name}}, error) {
	// Decode event parameters (only non-indexed parameters are in data)
	var result {{.Struct.Name}}
	{{- $allSupported := true}}
	{{- $hasNonIndexed := false}}
	{{- range .Inputs}}{{- if not .Indexed}}{{- $hasNonIndexed = true}}{{- if not (decoderSupports .Type)}}{{- $allSupported = false}}{{- end}}{{- end}}{{- end}}
	{{- if not $allSupported}}
	return result, errors.New("unsupported parameter type in {{.Name}} event")
	{{- else}}
	{{- if $hasNonIndexed}}
	offset := 0
	{{- range $i, $input := .Inputs}}
	{{- if not $input.Indexed}}
	{{- template "decodeValue" (dict "Data" "data" "Field" (printf "result.%s" ($input.Name | title)) "T" $input.Type "Ctx" (printf "event parameter %s" $input.Name) "Zero" "result")}}
	{{- end}}
	{{- end}}
	{{- end}}
	return result, nil
	{{- end}}
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