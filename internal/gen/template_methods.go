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

// methodDecodersTemplate generates method decode functions
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
	// Single return value - use unified decoding approach
	offset := 0
	{{- $output := index .Outputs 0}}
	{{- if eq $output.Type.TypeName "*big.Int"}}
	if len(data) < offset+32 {
		return nil, errors.New("insufficient data for return value")
	}
	{{- if $output.Type.IsSigned}}
	return decodeInt256(data[offset:offset+32])
	{{- else}}
	return decodeUint256(data[offset:offset+32])
	{{- end}}
	{{- else if eq $output.Type.TypeName "uint64"}}
	if len(data) < offset+32 {
		return 0, errors.New("insufficient data for return value")
	}
	return decodeUint64(data[offset:offset+32])
	{{- else if eq $output.Type.TypeName "uint8"}}
	if len(data) < offset+32 {
		return 0, errors.New("insufficient data for return value")
	}
	return decodeUint8(data[offset:offset+32])
	{{- else if eq $output.Type.TypeName "uint16"}}
	if len(data) < offset+32 {
		return 0, errors.New("insufficient data for return value")
	}
	return decodeUint16(data[offset:offset+32])
	{{- else if eq $output.Type.TypeName "uint32"}}
	if len(data) < offset+32 {
		return 0, errors.New("insufficient data for return value")
	}
	return decodeUint32(data[offset:offset+32])
	{{- else if eq $output.Type.TypeName "int64"}}
	if len(data) < offset+32 {
		return 0, errors.New("insufficient data for return value")
	}
	return decodeInt64(data[offset:offset+32])
	{{- else if eq $output.Type.TypeName "int8"}}
	if len(data) < offset+32 {
		return 0, errors.New("insufficient data for return value")
	}
	val, err := decodeInt64(data[offset:offset+32])
	if err != nil {
		return 0, err
	}
	return int8(val), nil
	{{- else if eq $output.Type.TypeName "int16"}}
	if len(data) < offset+32 {
		return 0, errors.New("insufficient data for return value")
	}
	val, err := decodeInt64(data[offset:offset+32])
	if err != nil {
		return 0, err
	}
	return int16(val), nil
	{{- else if eq $output.Type.TypeName "int32"}}
	if len(data) < offset+32 {
		return 0, errors.New("insufficient data for return value")
	}
	val, err := decodeInt64(data[offset:offset+32])
	if err != nil {
		return 0, err
	}
	return int32(val), nil
	{{- else if eq $output.Type.TypeName "bool"}}
	if len(data) < offset+32 {
		return false, errors.New("insufficient data for return value")
	}
	return decodeBool(data[offset:offset+32])
	{{- else if eq $output.Type.TypeName "Address"}}
	if len(data) < offset+32 {
		return Address{}, errors.New("insufficient data for return value")
	}
	return decodeAddress(data[offset:offset+32])
	{{- else if eq $output.Type.TypeName "Hash"}}
	if len(data) < offset+32 {
		return Hash{}, errors.New("insufficient data for return value")
	}
	return decodeHash(data[offset:offset+32])
	{{- else if eq $output.Type.TypeName "string"}}
	result, _, err := decodeString(data, offset)
	return result, err
	{{- else if eq $output.Type.TypeName "[]byte"}}
	result, _, err := decodeBytes(data, offset)
	return result, err
	{{- else if eq $output.Type.TypeName "[1]byte"}}
	if len(data) < offset+32 {
		return [1]byte{}, errors.New("insufficient data for return value")
	}
	return decodeBytes1(data[offset:offset+32])
	{{- else if eq $output.Type.TypeName "[32]byte"}}
	if len(data) < offset+32 {
		return [32]byte{}, errors.New("insufficient data for return value")
	}
	return decodeBytes32(data[offset:offset+32])
	{{- else if eq $output.Type.TypeName "[]*big.Int"}}
	// Handle []*big.Int array
	var elems []interface{}
	var nextOffset int
	elems, nextOffset, err = decodeArray(data, offset, decodeUint256ArrayElement)
	if err != nil {
		return nil, err
	}
	result := make([]*big.Int, len(elems))
	for i, elem := range elems {
		result[i] = elem.(*big.Int)
	}
	return result, nil
	{{- else if eq $output.Type.TypeName "[]uint64"}}
	// Handle []uint64 array
	var elems []interface{}
	var nextOffset int
	elems, nextOffset, err = decodeArray(data, offset, func(d []byte) (interface{}, error) { return decodeUint64(d) })
	if err != nil {
		return nil, err
	}
	result := make([]uint64, len(elems))
	for i, elem := range elems {
		result[i] = elem.(uint64)
	}
	return result, nil
	{{- else if eq $output.Type.TypeName "[]Address"}}
	// Handle []Address array
	var elems []interface{}
	var nextOffset int
	elems, nextOffset, err = decodeArray(data, offset, decodeAddressArrayElement)
	if err != nil {
		return nil, err
	}
	result := make([]Address, len(elems))
	for i, elem := range elems {
		result[i] = elem.(Address)
	}
	return result, nil
	{{- else if eq $output.Type.TypeName "[]bool"}}
	// Handle []bool array
	var elems []interface{}
	var nextOffset int
	elems, nextOffset, err = decodeArray(data, offset, decodeBoolArrayElement)
	if err != nil {
		return nil, err
	}
	result := make([]bool, len(elems))
	for i, elem := range elems {
		result[i] = elem.(bool)
	}
	return result, nil
	{{- else}}
	// Handle struct types
	{{- range $.Contract.Structs}}
	{{- if eq .Name $output.Type.TypeName}}
	result, _, err := decode{{.Name}}(data, offset)
	return result, err
	{{- end}}
	{{- end}}
	// Handle struct array types
	{{- if and $output.Type.IsSlice (gt (len $output.Type.TypeName) 2)}}
	{{- $elemType := slice $output.Type.TypeName 2}}
	{{- range $.Contract.Structs}}
	{{- if eq .Name $elemType}}
	// Read offset pointer to array data
	if len(data) < offset+32 {
		return nil, errors.New("insufficient data for array offset pointer")
	}
	arrayOffset, err := decodeUint256(data[offset:offset+32])
	if err != nil {
		return nil, fmt.Errorf("decoding array offset pointer: %w", err)
	}
	if !arrayOffset.IsUint64() {
		return nil, errors.New("array offset too large")
	}
	arrayOffsetInt := int(arrayOffset.Uint64())
	
	// Read array length at the offset location
	if len(data) < arrayOffsetInt+32 {
		return nil, errors.New("insufficient data for array length")
	}
	val, err := decodeUint256(data[arrayOffsetInt:arrayOffsetInt+32])
	if err != nil {
		return nil, fmt.Errorf("decoding array length: %w", err)
	}
	if !val.IsUint64() {
		return nil, errors.New("array length too large")
	}
	length := int(val.Uint64())
	offset = arrayOffsetInt + 32
	
	result := make({{$output.Type.TypeName}}, length)
	for i := 0; i < length; i++ {
		var elem {{.Name}}
		var nextOffset int
		elem, nextOffset, err = decode{{.Name}}(data, offset)
		if err != nil {
			return nil, fmt.Errorf("decoding array element %d: %w", i, err)
		}
		result[i] = elem
		offset = nextOffset
	}
	return result, nil
	{{- end}}
	{{- end}}
	{{- end}}
	return {{formatGoType $output.Type}}{}, errors.New("unsupported return type: {{$output.Type.TypeName}}")
	{{- end}}
{{- else}}
	// Multiple return values - return as struct
	var result {{.Name | title}}Result
	{{- $needsVal := false}}
	{{- $needsValAddr := false}}
	{{- $needsValBool := false}}
	{{- $needsValUint64 := false}}
	{{- $needsValInt64 := false}}
	{{- $needsValString := false}}
	{{- $needsValBytes := false}}
	{{- range .Outputs}}
		{{- if eq .Type.TypeName "*big.Int"}}
			{{- $needsVal = true}}
		{{- end}}
		{{- if eq .Type.TypeName "Address"}}
			{{- $needsValAddr = true}}
		{{- end}}
		{{- if eq .Type.TypeName "bool"}}
			{{- $needsValBool = true}}
		{{- end}}
		{{- if eq .Type.TypeName "uint64"}}
			{{- $needsValUint64 = true}}
		{{- end}}
		{{- if eq .Type.TypeName "int64"}}
			{{- $needsValInt64 = true}}
		{{- end}}
		{{- if eq .Type.TypeName "string"}}
			{{- $needsValString = true}}
		{{- end}}
		{{- if eq .Type.TypeName "[]byte"}}
			{{- $needsValBytes = true}}
		{{- end}}
	{{- end}}
	{{- if $needsVal}}
	var val *big.Int
	{{- end}}
	{{- if $needsValAddr}}
	var valAddr Address
	{{- end}}
	{{- if $needsValBool}}
	var valBool bool
	{{- end}}
	{{- if $needsValUint64}}
	var valUint64 uint64
	{{- end}}
	{{- if $needsValInt64}}
	var valInt64 int64
	{{- end}}
	{{- if $needsValString}}
	var valString string
	{{- end}}
	{{- if $needsValBytes}}
	var valBytes []byte
	{{- end}}
	var err error
	offset := 0
	{{- range $i, $output := .Outputs}}
	{{- if eq $output.Type.TypeName "*big.Int"}}
	if len(data) < offset+32 {
		return result, errors.New("insufficient data for return value {{$i}}")
	}
	{{- if $output.Type.IsSigned}}
	val, err = decodeInt256(data[offset:offset+32])
	if err != nil {
		return result, fmt.Errorf("decoding return value {{$i}}: %w", err)
	}
	result.{{$output.Name | title}} = val
	offset += 32
	{{- else}}
	val, err = decodeUint256(data[offset:offset+32])
	if err != nil {
		return result, fmt.Errorf("decoding return value {{$i}}: %w", err)
	}
	result.{{$output.Name | title}} = val
	offset += 32
	{{- end}}
	{{- else if eq $output.Type.TypeName "uint64"}}
	if len(data) < offset+32 {
		return result, errors.New("insufficient data for return value {{$i}}")
	}
	valUint64, err = decodeUint64(data[offset:offset+32])
	if err != nil {
		return result, fmt.Errorf("decoding return value {{$i}}: %w", err)
	}
	result.{{$output.Name | title}} = valUint64
	offset += 32
	{{- else if eq $output.Type.TypeName "int64"}}
	if len(data) < offset+32 {
		return result, errors.New("insufficient data for return value {{$i}}")
	}
	valInt64, err = decodeInt64(data[offset:offset+32])
	if err != nil {
		return result, fmt.Errorf("decoding return value {{$i}}: %w", err)
	}
	result.{{$output.Name | title}} = valInt64
	offset += 32
	{{- else if eq $output.Type.TypeName "bool"}}
	if len(data) < offset+32 {
		return result, errors.New("insufficient data for return value {{$i}}")
	}
	valBool, err = decodeBool(data[offset:offset+32])
	if err != nil {
		return result, fmt.Errorf("decoding return value {{$i}}: %w", err)
	}
	result.{{$output.Name | title}} = valBool
	offset += 32
	{{- else if eq $output.Type.TypeName "Address"}}
	if len(data) < offset+32 {
		return result, errors.New("insufficient data for return value {{$i}}")
	}
	valAddr, err = decodeAddress(data[offset:offset+32])
	if err != nil {
		return result, fmt.Errorf("decoding return value {{$i}}: %w", err)
	}
	result.{{$output.Name | title}} = valAddr
	offset += 32
	{{- else if eq $output.Type.TypeName "[]*big.Int"}}
	// Handle []*big.Int array
	var elems []interface{}
	var nextOffset int
	elems, nextOffset, err = decodeArray(data, offset, decodeUint256ArrayElement)
	if err != nil {
		return result, fmt.Errorf("decoding return value {{$i}}: %w", err)
	}
	bigIntArray := make([]*big.Int, len(elems))
	for j, elem := range elems {
		bigIntArray[j] = elem.(*big.Int)
	}
	result.{{$output.Name | title}} = bigIntArray
	offset = nextOffset
	{{- else if eq $output.Type.TypeName "[]uint64"}}
	// Handle []uint64 array
	var elems []interface{}
	var nextOffset int
	elems, nextOffset, err = decodeArray(data, offset, func(d []byte) (interface{}, error) { return decodeUint64(d) })
	if err != nil {
		return result, fmt.Errorf("decoding return value {{$i}}: %w", err)
	}
	uint64Array := make([]uint64, len(elems))
	for j, elem := range elems {
		uint64Array[j] = elem.(uint64)
	}
	result.{{$output.Name | title}} = uint64Array
	offset = nextOffset
	{{- else if eq $output.Type.TypeName "[]Address"}}
	// Handle []Address array
	var elems []interface{}
	var nextOffset int
	elems, nextOffset, err = decodeArray(data, offset, decodeAddressArrayElement)
	if err != nil {
		return result, fmt.Errorf("decoding return value {{$i}}: %w", err)
	}
	addressArray := make([]Address, len(elems))
	for j, elem := range elems {
		addressArray[j] = elem.(Address)
	}
	result.{{$output.Name | title}} = addressArray
	offset = nextOffset
	{{- else if eq $output.Type.TypeName "[]bool"}}
	// Handle []bool array
	var elems []interface{}
	var nextOffset int
	elems, nextOffset, err = decodeArray(data, offset, decodeBoolArrayElement)
	if err != nil {
		return result, fmt.Errorf("decoding return value {{$i}}: %w", err)
	}
	boolArray := make([]bool, len(elems))
	for j, elem := range elems {
		boolArray[j] = elem.(bool)
	}
	result.{{$output.Name | title}} = boolArray
	offset = nextOffset
	{{- else if eq $output.Type.TypeName "string"}}
	// Handle string
	var nextOffset int
	valString, nextOffset, err = decodeString(data, offset)
	if err != nil {
		return result, fmt.Errorf("decoding return value {{$i}}: %w", err)
	}
	result.{{$output.Name | title}} = valString
	offset = nextOffset
	{{- else if eq $output.Type.TypeName "[]byte"}}
	// Handle []byte
	var nextOffset int
	valBytes, nextOffset, err = decodeBytes(data, offset)
	if err != nil {
		return result, fmt.Errorf("decoding return value {{$i}}: %w", err)
	}
	result.{{$output.Name | title}} = valBytes
	offset = nextOffset
	{{- else}}
	// Handle struct types in multi-return
	{{- range $.Contract.Structs}}
	{{- if eq .Name $output.Type.TypeName}}
	var structVal {{.Name}}
	var nextOffset int
	structVal, nextOffset, err = decode{{.Name}}(data, offset)
	if err != nil {
		return result, fmt.Errorf("decoding return value {{$i}}: %w", err)
	}
	result.{{$output.Name | title}} = structVal
	offset = nextOffset
	{{- end}}
	{{- end}}
	// Handle struct array types in multi-return
	{{- if and $output.Type.IsSlice (gt (len $output.Type.TypeName) 2)}}
	{{- $elemType := slice $output.Type.TypeName 2}}
	{{- range $.Contract.Structs}}
	{{- if eq .Name $elemType}}
	if len(data) < offset+32 {
		return result, errors.New("insufficient data for array length in return value {{$i}}")
	}
	val, err := decodeUint256(data[offset:offset+32])
	if err != nil {
		return result, fmt.Errorf("decoding array length in return value {{$i}}: %w", err)
	}
	if !val.IsUint64() {
		return result, errors.New("array length too large in return value {{$i}}")
	}
	length := int(val.Uint64())
	offset += 32
	
	structArray := make({{$output.Type.TypeName}}, length)
	for j := 0; j < length; j++ {
		var elem {{.Name}}
		var nextOffset int
		elem, nextOffset, err = decode{{.Name}}(data, offset)
		if err != nil {
			return result, fmt.Errorf("decoding array element %d in return value {{$i}}: %w", j, err)
		}
		structArray[j] = elem
		offset = nextOffset
	}
	result.{{$output.Name | title}} = structArray
	{{- end}}
	{{- end}}
	{{- else}}
	return result, errors.New("unsupported multi-return type: {{$output.Type.TypeName}}")
	{{- end}}
	{{- end}}
	{{- end}}
	return result, nil
{{- end}}
}
{{- end}}
{{- end}}`