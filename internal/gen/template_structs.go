// SPDX-License-Identifier: MIT

package gen

// structDecodersTemplate generates struct decoder functions
const structDecodersTemplate = `{{/* Generate struct decoders for all structs */}}
{{- range .Contract.Structs}}
// decode{{.Name}} decodes a {{.Name}} struct from ABI-encoded data
func decode{{.Name}}(data []byte, offset int) ({{.Name}}, int, error) {
	var result {{.Name}}
	{{- $needsVal := false}}
	{{- $needsValAddr := false}}
	{{- $needsValHash := false}}
	{{- $needsValBool := false}}
	{{- $needsValStr := false}}
	{{- $needsValBytes := false}}
	{{- $needsValUint64 := false}}
	{{- $needsValUint32 := false}}
	{{- $needsValUint16 := false}}
	{{- $needsValUint8 := false}}
	{{- $needsValInt64 := false}}
	{{- $needsValBytes1 := false}}
	{{- $needsValBytes32 := false}}
	{{- range .Fields}}
		{{- if or (eq .Type.TypeName "*big.Int") (and .Type.IsSlice (or (eq .Type.TypeName "[]*big.Int") (hasPrefix .Type.TypeName "[]")))}}
			{{- $needsVal = true}}
		{{- end}}
		{{- if eq .Type.TypeName "Address"}}
			{{- $needsValAddr = true}}
		{{- end}}
		{{- if eq .Type.TypeName "Hash"}}
			{{- $needsValHash = true}}
		{{- end}}
		{{- if eq .Type.TypeName "bool"}}
			{{- $needsValBool = true}}
		{{- end}}
		{{- if eq .Type.TypeName "string"}}
			{{- $needsValStr = true}}
		{{- end}}
		{{- if eq .Type.TypeName "[]byte"}}
			{{- $needsValBytes = true}}
		{{- end}}
		{{- if eq .Type.TypeName "uint64"}}
			{{- $needsValUint64 = true}}
		{{- end}}
		{{- if eq .Type.TypeName "uint32"}}
			{{- $needsValUint32 = true}}
		{{- end}}
		{{- if eq .Type.TypeName "uint16"}}
			{{- $needsValUint16 = true}}
		{{- end}}
		{{- if eq .Type.TypeName "uint8"}}
			{{- $needsValUint8 = true}}
		{{- end}}
		{{- if or (eq .Type.TypeName "int64") (eq .Type.TypeName "int8") (eq .Type.TypeName "int16") (eq .Type.TypeName "int32")}}
			{{- $needsValInt64 = true}}
		{{- end}}
		{{- if eq .Type.TypeName "[1]byte"}}
			{{- $needsValBytes1 = true}}
		{{- end}}
		{{- if eq .Type.TypeName "[32]byte"}}
			{{- $needsValBytes32 = true}}
		{{- end}}
	{{- end}}
	{{- if $needsVal}}
	var val *big.Int
	{{- end}}
	{{- if $needsValAddr}}
	var valAddr Address
	{{- end}}
	{{- if $needsValHash}}
	var valHash Hash
	{{- end}}
	{{- if $needsValBool}}
	var valBool bool
	{{- end}}
	{{- if $needsValStr}}
	var valStr string
	{{- end}}
	{{- if $needsValBytes}}
	var valBytes []byte
	{{- end}}
	{{- if $needsValUint64}}
	var valUint64 uint64
	{{- end}}
	{{- if $needsValUint32}}
	var valUint32 uint32
	{{- end}}
	{{- if $needsValUint16}}
	var valUint16 uint16
	{{- end}}
	{{- if $needsValUint8}}
	var valUint8 uint8
	{{- end}}
	{{- if $needsValInt64}}
	var valInt64 int64
	{{- end}}
	{{- if $needsValBytes1}}
	var valBytes1 [1]byte
	{{- end}}
	{{- if $needsValBytes32}}
	var valBytes32 [32]byte
	{{- end}}
	var err error
	currentOffset := offset
	{{- $structName := .Name}}
	{{- range .Fields}}
	{{- if eq .Type.TypeName "*big.Int"}}
	if len(data) < currentOffset+32 {
		return result, 0, errors.New("insufficient data for {{$structName}}.{{.Name}}")
	}
	{{- if .Type.IsSigned}}
	val, err = decodeInt256(data[currentOffset:currentOffset+32])
	if err != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{.Name}}: %w", err)
	}
	result.{{.Name}} = val
	{{- else}}
	val, err = decodeUint256(data[currentOffset:currentOffset+32])
	if err != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{.Name}}: %w", err)
	}
	result.{{.Name}} = val
	{{- end}}
	currentOffset += 32
	{{- else if eq .Type.TypeName "uint64"}}
	if len(data) < currentOffset+32 {
		return result, 0, errors.New("insufficient data for {{$structName}}.{{.Name}}")
	}
	valUint64, err = decodeUint64(data[currentOffset:currentOffset+32])
	if err != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{.Name}}: %w", err)
	}
	result.{{.Name}} = valUint64
	currentOffset += 32
	{{- else if eq .Type.TypeName "uint8"}}
	if len(data) < currentOffset+32 {
		return result, 0, errors.New("insufficient data for {{$structName}}.{{.Name}}")
	}
	valUint8, err = decodeUint8(data[currentOffset:currentOffset+32])
	if err != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{.Name}}: %w", err)
	}
	result.{{.Name}} = valUint8
	currentOffset += 32
	{{- else if eq .Type.TypeName "uint16"}}
	if len(data) < currentOffset+32 {
		return result, 0, errors.New("insufficient data for {{$structName}}.{{.Name}}")
	}
	valUint16, err = decodeUint16(data[currentOffset:currentOffset+32])
	if err != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{.Name}}: %w", err)
	}
	result.{{.Name}} = valUint16
	currentOffset += 32
	{{- else if eq .Type.TypeName "uint32"}}
	if len(data) < currentOffset+32 {
		return result, 0, errors.New("insufficient data for {{$structName}}.{{.Name}}")
	}
	valUint32, err = decodeUint32(data[currentOffset:currentOffset+32])
	if err != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{.Name}}: %w", err)
	}
	result.{{.Name}} = valUint32
	currentOffset += 32
	{{- else if eq .Type.TypeName "int64"}}
	if len(data) < currentOffset+32 {
		return result, 0, errors.New("insufficient data for {{$structName}}.{{.Name}}")
	}
	valInt64, err = decodeInt64(data[currentOffset:currentOffset+32])
	if err != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{.Name}}: %w", err)
	}
	result.{{.Name}} = valInt64
	currentOffset += 32
	{{- else if eq .Type.TypeName "int8"}}
	if len(data) < currentOffset+32 {
		return result, 0, errors.New("insufficient data for {{$structName}}.{{.Name}}")
	}
	valInt64, err = decodeInt64(data[currentOffset:currentOffset+32])
	if err != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{.Name}}: %w", err)
	}
	result.{{.Name}} = int8(valInt64)
	currentOffset += 32
	{{- else if eq .Type.TypeName "int16"}}
	if len(data) < currentOffset+32 {
		return result, 0, errors.New("insufficient data for {{$structName}}.{{.Name}}")
	}
	valInt64, err = decodeInt64(data[currentOffset:currentOffset+32])
	if err != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{.Name}}: %w", err)
	}
	result.{{.Name}} = int16(valInt64)
	currentOffset += 32
	{{- else if eq .Type.TypeName "int32"}}
	if len(data) < currentOffset+32 {
		return result, 0, errors.New("insufficient data for {{$structName}}.{{.Name}}")
	}
	valInt64, err = decodeInt64(data[currentOffset:currentOffset+32])
	if err != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{.Name}}: %w", err)
	}
	result.{{.Name}} = int32(valInt64)
	currentOffset += 32
	{{- else if eq .Type.TypeName "bool"}}
	if len(data) < currentOffset+32 {
		return result, 0, errors.New("insufficient data for {{$structName}}.{{.Name}}")
	}
	valBool, err = decodeBool(data[currentOffset:currentOffset+32])
	if err != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{.Name}}: %w", err)
	}
	result.{{.Name}} = valBool
	currentOffset += 32
	{{- else if eq .Type.TypeName "Address"}}
	if len(data) < currentOffset+32 {
		return result, 0, errors.New("insufficient data for {{$structName}}.{{.Name}}")
	}
	valAddr, err = decodeAddress(data[currentOffset:currentOffset+32])
	if err != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{.Name}}: %w", err)
	}
	result.{{.Name}} = valAddr
	currentOffset += 32
	{{- else if eq .Type.TypeName "Hash"}}
	if len(data) < currentOffset+32 {
		return result, 0, errors.New("insufficient data for {{$structName}}.{{.Name}}")
	}
	valHash, err = decodeHash(data[currentOffset:currentOffset+32])
	if err != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{.Name}}: %w", err)
	}
	result.{{.Name}} = valHash
	currentOffset += 32
	{{- else if eq .Type.TypeName "string"}}
	var nextOffset int
	valStr, nextOffset, err = decodeString(data, currentOffset)
	if err != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{.Name}}: %w", err)
	}
	result.{{.Name}} = valStr
	currentOffset = nextOffset
	{{- else if eq .Type.TypeName "[]byte"}}
	var nextOffset int
	valBytes, nextOffset, err = decodeBytes(data, currentOffset)
	if err != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{.Name}}: %w", err)
	}
	result.{{.Name}} = valBytes
	currentOffset = nextOffset
	{{- else if eq .Type.TypeName "[1]byte"}}
	if len(data) < currentOffset+32 {
		return result, 0, errors.New("insufficient data for {{$structName}}.{{.Name}}")
	}
	valBytes1, err = decodeBytes1(data[currentOffset:currentOffset+32])
	if err != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{.Name}}: %w", err)
	}
	result.{{.Name}} = valBytes1
	currentOffset += 32
	{{- else if eq .Type.TypeName "[32]byte"}}
	if len(data) < currentOffset+32 {
		return result, 0, errors.New("insufficient data for {{$structName}}.{{.Name}}")
	}
	valBytes32, err = decodeBytes32(data[currentOffset:currentOffset+32])
	if err != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{.Name}}: %w", err)
	}
	result.{{.Name}} = valBytes32
	currentOffset += 32
	{{- else if and .Type.IsSlice (eq .Type.TypeName "[]*big.Int")}}
	var elems []interface{}
	var nextOffset int
	elems, nextOffset, err = decodeArray(data, currentOffset, decodeUint256ArrayElement)
	if err != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{.Name}}: %w", err)
	}
	result.{{.Name}} = make([]*big.Int, len(elems))
	for i, elem := range elems {
		result.{{.Name}}[i] = elem.(*big.Int)
	}
	currentOffset = nextOffset
	{{- else if and .Type.IsSlice (eq .Type.TypeName "[]uint64")}}
	var elems []interface{}
	var nextOffset int
	elems, nextOffset, err = decodeArray(data, currentOffset, func(d []byte) (interface{}, error) { return decodeUint64(d) })
	if err != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{.Name}}: %w", err)
	}
	result.{{.Name}} = make([]uint64, len(elems))
	for i, elem := range elems {
		result.{{.Name}}[i] = elem.(uint64)
	}
	currentOffset = nextOffset
	{{- else if and .Type.IsSlice (eq .Type.TypeName "[]Address")}}
	var elems []interface{}
	var nextOffset int
	elems, nextOffset, err = decodeArray(data, currentOffset, decodeAddressArrayElement)
	if err != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{.Name}}: %w", err)
	}
	result.{{.Name}} = make([]Address, len(elems))
	for i, elem := range elems {
		result.{{.Name}}[i] = elem.(Address)
	}
	currentOffset = nextOffset
	{{- else if .Type.IsSlice}}
	// Handle struct array field: {{.Type.TypeName}}
	if len(data) < currentOffset+32 {
		return result, 0, errors.New("insufficient data for struct array length in {{$structName}}.{{.Name}}")
	}
	val, err = decodeUint256(data[currentOffset:currentOffset+32])
	if err != nil {
		return result, 0, fmt.Errorf("decoding {{$structName}}.{{.Name}} length: %w", err)
	}
	if !val.IsUint64() {
		return result, 0, errors.New("struct array length too large in {{$structName}}.{{.Name}}")
	}
	length := int(val.Uint64())
	currentOffset += 32
	
	elemTypeName := "{{.Type.TypeName}}"[2:] // Remove "[]" prefix
	{{- $outerContract := $.Contract}}
	{{- $fieldName := .Name}}
	{{- $fieldType := .Type.TypeName}}
	{{- range $struct := $outerContract.Structs}}
	if elemTypeName == "{{$struct.Name}}" {
		result.{{$fieldName}} = make({{$fieldType}}, length)
		for i := 0; i < length; i++ {
			var elem {{$struct.Name}}
			var nextOffsetStruct int
			elem, nextOffsetStruct, err = decode{{$struct.Name}}(data, currentOffset)
			if err != nil {
				return result, 0, fmt.Errorf("decoding {{$structName}}.{{$fieldName}}[%d]: %w", i, err)
			}
			result.{{$fieldName}}[i] = elem
			currentOffset = nextOffsetStruct
		}
	}
	{{- end}}
	{{- else}}
	return result, 0, errors.New("unsupported struct field type {{.Type.TypeName}} in {{$structName}}.{{.Name}}")
	{{- end}}
	{{- end}}
	return result, currentOffset, nil
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