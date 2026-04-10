// SPDX-License-Identifier: MIT

package gen

// structDecodersTemplate generates struct decoder functions.
//
// Each decoder uses the standard ABI head-tail layout:
//   - Phase 1 (head): static fields are decoded inline (32 bytes each);
//     dynamic fields (string, []byte, slices) occupy a 32-byte relative-offset
//     slot in the head and are decoded in Phase 2.
//   - Phase 2 (tail): dynamic fields are decoded at their absolute positions
//     (baseOffset + relative_offset).
//
// Nested struct types are treated as static in the outer head; if the nested
// struct itself has dynamic fields, its own decoder handles those internally.
const structDecodersTemplate = `{{/* Generate struct decoders for all structs */}}
{{- range .Contract.Structs}}
// decode{{.Name}} decodes a {{.Name}} struct from ABI-encoded data starting at offset.
// It returns the decoded value and the offset of the next field after the head section.
func decode{{.Name}}(data []byte, offset int) ({{.Name}}, int, error) {
	var result {{.Name}}
	baseOffset := offset
	headOffset := offset

	// ── Phase 1: head section ──────────────────────────────────────────────────
	// Static fields are decoded inline. Dynamic fields (string, []byte, slices)
	// store a relative-offset pointer here; the actual data is in the tail.
	{{- $structName := .Name}}
	{{- range $i, $field := .Fields}}
	{{- if isDynamic $field.Type}}
	// {{$field.Name}} (dynamic): read relative-offset pointer
	if len(data) < headOffset+32 {
		return result, 0, errors.New("insufficient data for {{$structName}}.{{$field.Name}} offset pointer")
	}
	relOff{{$i}}, errRO{{$i}} := decodeUint256(data[headOffset : headOffset+32])
	if errRO{{$i}} != nil {
		return result, 0, fmt.Errorf("reading {{$structName}}.{{$field.Name}} offset: %w", errRO{{$i}})
	}
	if !relOff{{$i}}.IsUint64() {
		return result, 0, errors.New("{{$structName}}.{{$field.Name}} offset too large")
	}
	absOff{{$i}} := baseOffset + int(relOff{{$i}}.Uint64())
	headOffset += 32
	{{- else if eq $field.Type.TypeName "*big.Int"}}
	// {{$field.Name}} (static *big.Int)
	if len(data) < headOffset+32 {
		return result, 0, errors.New("insufficient data for {{$structName}}.{{$field.Name}}")
	}
	{{- if $field.Type.IsSigned}}
	val{{$i}}, err{{$i}} := decodeInt256(data[headOffset : headOffset+32])
	{{- else}}
	val{{$i}}, err{{$i}} := decodeUint256(data[headOffset : headOffset+32])
	{{- end}}
	if err{{$i}} != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{$field.Name}}: %w", err{{$i}})
	}
	result.{{$field.Name}} = val{{$i}}
	headOffset += 32
	{{- else if eq $field.Type.TypeName "uint64"}}
	if len(data) < headOffset+32 {
		return result, 0, errors.New("insufficient data for {{$structName}}.{{$field.Name}}")
	}
	val{{$i}}, err{{$i}} := decodeUint64(data[headOffset : headOffset+32])
	if err{{$i}} != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{$field.Name}}: %w", err{{$i}})
	}
	result.{{$field.Name}} = val{{$i}}
	headOffset += 32
	{{- else if eq $field.Type.TypeName "uint32"}}
	if len(data) < headOffset+32 {
		return result, 0, errors.New("insufficient data for {{$structName}}.{{$field.Name}}")
	}
	val{{$i}}, err{{$i}} := decodeUint32(data[headOffset : headOffset+32])
	if err{{$i}} != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{$field.Name}}: %w", err{{$i}})
	}
	result.{{$field.Name}} = val{{$i}}
	headOffset += 32
	{{- else if eq $field.Type.TypeName "uint16"}}
	if len(data) < headOffset+32 {
		return result, 0, errors.New("insufficient data for {{$structName}}.{{$field.Name}}")
	}
	val{{$i}}, err{{$i}} := decodeUint16(data[headOffset : headOffset+32])
	if err{{$i}} != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{$field.Name}}: %w", err{{$i}})
	}
	result.{{$field.Name}} = val{{$i}}
	headOffset += 32
	{{- else if eq $field.Type.TypeName "uint8"}}
	if len(data) < headOffset+32 {
		return result, 0, errors.New("insufficient data for {{$structName}}.{{$field.Name}}")
	}
	val{{$i}}, err{{$i}} := decodeUint8(data[headOffset : headOffset+32])
	if err{{$i}} != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{$field.Name}}: %w", err{{$i}})
	}
	result.{{$field.Name}} = val{{$i}}
	headOffset += 32
	{{- else if eq $field.Type.TypeName "int64"}}
	if len(data) < headOffset+32 {
		return result, 0, errors.New("insufficient data for {{$structName}}.{{$field.Name}}")
	}
	val{{$i}}, err{{$i}} := decodeInt64(data[headOffset : headOffset+32])
	if err{{$i}} != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{$field.Name}}: %w", err{{$i}})
	}
	result.{{$field.Name}} = val{{$i}}
	headOffset += 32
	{{- else if eq $field.Type.TypeName "int32"}}
	if len(data) < headOffset+32 {
		return result, 0, errors.New("insufficient data for {{$structName}}.{{$field.Name}}")
	}
	raw{{$i}}, err{{$i}} := decodeInt64(data[headOffset : headOffset+32])
	if err{{$i}} != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{$field.Name}}: %w", err{{$i}})
	}
	result.{{$field.Name}} = int32(raw{{$i}})
	headOffset += 32
	{{- else if eq $field.Type.TypeName "int16"}}
	if len(data) < headOffset+32 {
		return result, 0, errors.New("insufficient data for {{$structName}}.{{$field.Name}}")
	}
	raw{{$i}}, err{{$i}} := decodeInt64(data[headOffset : headOffset+32])
	if err{{$i}} != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{$field.Name}}: %w", err{{$i}})
	}
	result.{{$field.Name}} = int16(raw{{$i}})
	headOffset += 32
	{{- else if eq $field.Type.TypeName "int8"}}
	if len(data) < headOffset+32 {
		return result, 0, errors.New("insufficient data for {{$structName}}.{{$field.Name}}")
	}
	raw{{$i}}, err{{$i}} := decodeInt64(data[headOffset : headOffset+32])
	if err{{$i}} != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{$field.Name}}: %w", err{{$i}})
	}
	result.{{$field.Name}} = int8(raw{{$i}})
	headOffset += 32
	{{- else if eq $field.Type.TypeName "bool"}}
	if len(data) < headOffset+32 {
		return result, 0, errors.New("insufficient data for {{$structName}}.{{$field.Name}}")
	}
	val{{$i}}, err{{$i}} := decodeBool(data[headOffset : headOffset+32])
	if err{{$i}} != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{$field.Name}}: %w", err{{$i}})
	}
	result.{{$field.Name}} = val{{$i}}
	headOffset += 32
	{{- else if eq $field.Type.TypeName "Address"}}
	if len(data) < headOffset+32 {
		return result, 0, errors.New("insufficient data for {{$structName}}.{{$field.Name}}")
	}
	val{{$i}}, err{{$i}} := decodeAddress(data[headOffset : headOffset+32])
	if err{{$i}} != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{$field.Name}}: %w", err{{$i}})
	}
	result.{{$field.Name}} = val{{$i}}
	headOffset += 32
	{{- else if eq $field.Type.TypeName "Hash"}}
	if len(data) < headOffset+32 {
		return result, 0, errors.New("insufficient data for {{$structName}}.{{$field.Name}}")
	}
	val{{$i}}, err{{$i}} := decodeHash(data[headOffset : headOffset+32])
	if err{{$i}} != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{$field.Name}}: %w", err{{$i}})
	}
	result.{{$field.Name}} = val{{$i}}
	headOffset += 32
	{{- else if eq $field.Type.TypeName "[1]byte"}}
	if len(data) < headOffset+32 {
		return result, 0, errors.New("insufficient data for {{$structName}}.{{$field.Name}}")
	}
	val{{$i}}, err{{$i}} := decodeBytes1(data[headOffset : headOffset+32])
	if err{{$i}} != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{$field.Name}}: %w", err{{$i}})
	}
	result.{{$field.Name}} = val{{$i}}
	headOffset += 32
	{{- else if eq $field.Type.TypeName "[32]byte"}}
	if len(data) < headOffset+32 {
		return result, 0, errors.New("insufficient data for {{$structName}}.{{$field.Name}}")
	}
	val{{$i}}, err{{$i}} := decodeBytes32(data[headOffset : headOffset+32])
	if err{{$i}} != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{$field.Name}}: %w", err{{$i}})
	}
	result.{{$field.Name}} = val{{$i}}
	headOffset += 32
	{{- else}}
	// {{$field.Name}}: nested struct type {{$field.Type.TypeName}} — treated as static
	{{- range $struct := $.Contract.Structs}}
	{{- if eq $struct.Name $field.Type.TypeName}}
	val{{$i}}, nextOff{{$i}}, err{{$i}} := decode{{$struct.Name}}(data, headOffset)
	if err{{$i}} != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{$field.Name}}: %w", err{{$i}})
	}
	result.{{$field.Name}} = val{{$i}}
	headOffset = nextOff{{$i}}
	{{- end}}
	{{- end}}
	{{- end}}
	{{- end}}

	// ── Phase 2: tail section — decode dynamic fields at their absolute offsets ─
	{{- range $i, $field := .Fields}}
	{{- if isDynamic $field.Type}}
	{{- if eq $field.Type.TypeName "string"}}
	str{{$i}}, _, errStr{{$i}} := decodeString(data, absOff{{$i}})
	if errStr{{$i}} != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{$field.Name}}: %w", errStr{{$i}})
	}
	result.{{$field.Name}} = str{{$i}}
	{{- else if eq $field.Type.TypeName "[]byte"}}
	bytes{{$i}}, _, errB{{$i}} := decodeBytes(data, absOff{{$i}})
	if errB{{$i}} != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{$field.Name}}: %w", errB{{$i}})
	}
	result.{{$field.Name}} = bytes{{$i}}
	{{- else if eq $field.Type.TypeName "[]*big.Int"}}
	elems{{$i}}, _, errArr{{$i}} := decodeArray(data, absOff{{$i}}, decodeUint256ArrayElement)
	if errArr{{$i}} != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{$field.Name}}: %w", errArr{{$i}})
	}
	arr{{$i}} := make([]*big.Int, len(elems{{$i}}))
	for j, elem := range elems{{$i}} {
		arr{{$i}}[j] = elem.(*big.Int)
	}
	result.{{$field.Name}} = arr{{$i}}
	{{- else if eq $field.Type.TypeName "[]uint64"}}
	elems{{$i}}, _, errArr{{$i}} := decodeArray(data, absOff{{$i}}, func(d []byte) (interface{}, error) { return decodeUint64(d) })
	if errArr{{$i}} != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{$field.Name}}: %w", errArr{{$i}})
	}
	arr{{$i}} := make([]uint64, len(elems{{$i}}))
	for j, elem := range elems{{$i}} {
		arr{{$i}}[j] = elem.(uint64)
	}
	result.{{$field.Name}} = arr{{$i}}
	{{- else if eq $field.Type.TypeName "[]Address"}}
	elems{{$i}}, _, errArr{{$i}} := decodeArray(data, absOff{{$i}}, decodeAddressArrayElement)
	if errArr{{$i}} != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{$field.Name}}: %w", errArr{{$i}})
	}
	arr{{$i}} := make([]Address, len(elems{{$i}}))
	for j, elem := range elems{{$i}} {
		arr{{$i}}[j] = elem.(Address)
	}
	result.{{$field.Name}} = arr{{$i}}
	{{- else if eq $field.Type.TypeName "[]bool"}}
	elems{{$i}}, _, errArr{{$i}} := decodeArray(data, absOff{{$i}}, decodeBoolArrayElement)
	if errArr{{$i}} != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{$field.Name}}: %w", errArr{{$i}})
	}
	arr{{$i}} := make([]bool, len(elems{{$i}}))
	for j, elem := range elems{{$i}} {
		arr{{$i}}[j] = elem.(bool)
	}
	result.{{$field.Name}} = arr{{$i}}
	{{- else if .Type.IsSlice}}
	// Struct-slice field: {{$field.Type.TypeName}}
	lenBig{{$i}}, errLen{{$i}} := decodeUint256(data[absOff{{$i}} : absOff{{$i}}+32])
	if errLen{{$i}} != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{$field.Name}} length: %w", errLen{{$i}})
	}
	if !lenBig{{$i}}.IsUint64() {
		return result, 0, errors.New("{{$structName}}.{{$field.Name}} array length too large")
	}
	length{{$i}} := int(lenBig{{$i}}.Uint64())
	elemOff{{$i}} := absOff{{$i}} + 32
	{{- $outerContract := $.Contract}}
	{{- $fieldType := $field.Type.TypeName}}
	{{- range $struct := $outerContract.Structs}}
	{{- if eq (print "[]" $struct.Name) $fieldType}}
	slice{{$i}} := make({{$fieldType}}, length{{$i}})
	for k := 0; k < length{{$i}}; k++ {
		elem{{$i}}, nextElemOff{{$i}}, errElem{{$i}} := decode{{$struct.Name}}(data, elemOff{{$i}})
		if errElem{{$i}} != nil {
			return result, 0, fmt.Errorf("decoding {{$structName}}.{{$field.Name}}[%d]: %w", k, errElem{{$i}})
		}
		slice{{$i}}[k] = elem{{$i}}
		elemOff{{$i}} = nextElemOff{{$i}}
	}
	result.{{$field.Name}} = slice{{$i}}
	{{- end}}
	{{- end}}
	{{- else}}
	return result, 0, errors.New("unsupported dynamic field type {{$field.Type.TypeName}} in {{$structName}}.{{$field.Name}}")
	{{- end}}
	{{- end}}
	{{- end}}

	_ = baseOffset
	return result, headOffset, nil
}
{{- end}}`

// structDefinitionsTemplate generates struct type definitions
const structDefinitionsTemplate = `{{/* Generate event structs */}}
{{- range .Contract.Events}}

// {{.Struct.Name}} represents the {{.Name}} event
type {{.Struct.Name}} struct {
{{- range .Struct.Fields}}
	{{.Name}} {{formatGoType .Type}} ` + "`" + `json:"{{.JSONTag}}"` + "`" + `
{{- end}}
}
{{- end}}

{{/* Generate error structs */}}
{{- range .Contract.Errors}}

// {{.Struct.Name}} represents the {{.Name}} custom error
type {{.Struct.Name}} struct {
{{- range .Struct.Fields}}
	{{.Name}} {{formatGoType .Type}} ` + "`" + `json:"{{.JSONTag}}"` + "`" + `
{{- end}}
}
{{- end}}

{{/* Generate standalone structs */}}
{{- range .Contract.Structs}}

// {{.Name}} represents a Solidity struct
type {{.Name}} struct {
{{- range .Fields}}
	{{.Name}} {{formatGoType .Type}} ` + "`" + `json:"{{.JSONTag}}"` + "`" + `
{{- end}}
}
{{- end}}

{{/* Generate input/output structs for methods */}}
{{- range .Contract.Methods}}
{{- if .InputStruct}}

// {{.InputStruct.Name}} represents inputs for method {{.Name}}
type {{.InputStruct.Name}} struct {
{{- range .InputStruct.Fields}}
	{{.Name}} {{formatGoType .Type}} ` + "`" + `json:"{{.JSONTag}}"` + "`" + `
{{- end}}
}
{{- end}}

{{- if .OutputStruct}}

// {{.OutputStruct.Name}} represents outputs for method {{.Name}}
type {{.OutputStruct.Name}} struct {
{{- range .OutputStruct.Fields}}
	{{.Name}} {{formatGoType .Type}} ` + "`" + `json:"{{.JSONTag}}"` + "`" + `
{{- end}}
}
{{- end}}
{{- end}}

{{/* Generate constructor struct if needed */}}
{{- if and .Contract.Constructor .Contract.Constructor.InputStruct}}

// {{.Contract.Constructor.InputStruct.Name}} represents constructor inputs
type {{.Contract.Constructor.InputStruct.Name}} struct {
{{- range .Contract.Constructor.InputStruct.Fields}}
	{{.Name}} {{formatGoType .Type}} ` + "`" + `json:"{{.JSONTag}}"` + "`" + `
{{- end}}
}
{{- end}}

{{/* Generate custom result structs for methods with multiple return values */}}
{{- range .Contract.Methods}}
{{- if gt (len .Outputs) 1}}

// {{.Name | title}}Result represents the return values for {{.Name}} method
type {{.Name | title}}Result struct {
{{- range .Outputs}}
	{{.Name | title}} {{formatGoType .Type}} ` + "`" + `json:"{{.Name | lower}}"` + "`" + `
{{- end}}
}
{{- end}}
{{- end}}`