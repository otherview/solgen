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
	if len(data) < offset+32 {
		return "", errors.New("insufficient data for offset pointer")
	}
	ptrBig, err := decodeUint256(data[offset : offset+32])
	if err != nil {
		return "", fmt.Errorf("decoding offset pointer: %w", err)
	}
	if !ptrBig.IsUint64() {
		return "", errors.New("offset pointer too large")
	}
	result, _, err := decodeString(data, int(ptrBig.Uint64()))
	return result, err
	{{- else if eq $output.Type.TypeName "[]byte"}}
	if len(data) < offset+32 {
		return nil, errors.New("insufficient data for offset pointer")
	}
	ptrBig, err := decodeUint256(data[offset : offset+32])
	if err != nil {
		return nil, fmt.Errorf("decoding offset pointer: %w", err)
	}
	if !ptrBig.IsUint64() {
		return nil, errors.New("offset pointer too large")
	}
	result, _, err := decodeBytes(data, int(ptrBig.Uint64()))
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
	if len(data) < offset+32 {
		return nil, errors.New("insufficient data for offset pointer")
	}
	ptrBig, err := decodeUint256(data[offset : offset+32])
	if err != nil {
		return nil, fmt.Errorf("decoding offset pointer: %w", err)
	}
	if !ptrBig.IsUint64() {
		return nil, errors.New("offset pointer too large")
	}
	elems, _, err := decodeArray(data, int(ptrBig.Uint64()), decodeUint256ArrayElement)
	if err != nil {
		return nil, err
	}
	result := make([]*big.Int, len(elems))
	for i, elem := range elems {
		result[i] = elem.(*big.Int)
	}
	return result, nil
	{{- else if eq $output.Type.TypeName "[]uint64"}}
	if len(data) < offset+32 {
		return nil, errors.New("insufficient data for offset pointer")
	}
	ptrBig, err := decodeUint256(data[offset : offset+32])
	if err != nil {
		return nil, fmt.Errorf("decoding offset pointer: %w", err)
	}
	if !ptrBig.IsUint64() {
		return nil, errors.New("offset pointer too large")
	}
	elems, _, err := decodeArray(data, int(ptrBig.Uint64()), func(d []byte) (interface{}, error) { return decodeUint64(d) })
	if err != nil {
		return nil, err
	}
	result := make([]uint64, len(elems))
	for i, elem := range elems {
		result[i] = elem.(uint64)
	}
	return result, nil
	{{- else if eq $output.Type.TypeName "[]Address"}}
	if len(data) < offset+32 {
		return nil, errors.New("insufficient data for offset pointer")
	}
	ptrBig, err := decodeUint256(data[offset : offset+32])
	if err != nil {
		return nil, fmt.Errorf("decoding offset pointer: %w", err)
	}
	if !ptrBig.IsUint64() {
		return nil, errors.New("offset pointer too large")
	}
	elems, _, err := decodeArray(data, int(ptrBig.Uint64()), decodeAddressArrayElement)
	if err != nil {
		return nil, err
	}
	result := make([]Address, len(elems))
	for i, elem := range elems {
		result[i] = elem.(Address)
	}
	return result, nil
	{{- else if eq $output.Type.TypeName "[]bool"}}
	if len(data) < offset+32 {
		return nil, errors.New("insufficient data for offset pointer")
	}
	ptrBig, err := decodeUint256(data[offset : offset+32])
	if err != nil {
		return nil, fmt.Errorf("decoding offset pointer: %w", err)
	}
	if !ptrBig.IsUint64() {
		return nil, errors.New("offset pointer too large")
	}
	elems, _, err := decodeArray(data, int(ptrBig.Uint64()), decodeBoolArrayElement)
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
	{{- $needsValUint8 := false}}
	{{- $needsValUint16 := false}}
	{{- $needsValUint32 := false}}
	{{- $needsValUint64 := false}}
	{{- $needsValInt8 := false}}
	{{- $needsValInt16 := false}}
	{{- $needsValInt32 := false}}
	{{- $needsValInt64 := false}}
	{{- $needsValHash := false}}
	{{- $needsValBytes1 := false}}
	{{- $needsValBytes32 := false}}
	{{- $needsValString := false}}
	{{- $needsValBytes := false}}
	{{- $needsNextOffset := false}}
	{{- range .Outputs}}
		{{- if eq .Type.TypeName "*big.Int"}}{{- $needsVal = true}}{{- end}}
		{{- if eq .Type.TypeName "Address"}}{{- $needsValAddr = true}}{{- end}}
		{{- if eq .Type.TypeName "bool"}}{{- $needsValBool = true}}{{- end}}
		{{- if eq .Type.TypeName "uint8"}}{{- $needsValUint8 = true}}{{- end}}
		{{- if eq .Type.TypeName "uint16"}}{{- $needsValUint16 = true}}{{- end}}
		{{- if eq .Type.TypeName "uint32"}}{{- $needsValUint32 = true}}{{- end}}
		{{- if eq .Type.TypeName "uint64"}}{{- $needsValUint64 = true}}{{- end}}
		{{- if eq .Type.TypeName "int8"}}{{- $needsValInt8 = true}}{{- end}}
		{{- if eq .Type.TypeName "int16"}}{{- $needsValInt16 = true}}{{- end}}
		{{- if eq .Type.TypeName "int32"}}{{- $needsValInt32 = true}}{{- end}}
		{{- if eq .Type.TypeName "int64"}}{{- $needsValInt64 = true}}{{- end}}
		{{- if eq .Type.TypeName "Hash"}}{{- $needsValHash = true}}{{- end}}
		{{- if eq .Type.TypeName "[1]byte"}}{{- $needsValBytes1 = true}}{{- end}}
		{{- if eq .Type.TypeName "[32]byte"}}{{- $needsValBytes32 = true}}{{- end}}
		{{- if eq .Type.TypeName "string"}}{{- $needsValString = true}}{{- end}}
		{{- if eq .Type.TypeName "[]byte"}}{{- $needsValBytes = true}}{{- end}}
		{{- if .Type.IsStruct}}{{- $needsNextOffset = true}}{{- end}}
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
	{{- if $needsValUint8}}
	var valUint8 uint8
	{{- end}}
	{{- if $needsValUint16}}
	var valUint16 uint16
	{{- end}}
	{{- if $needsValUint32}}
	var valUint32 uint32
	{{- end}}
	{{- if $needsValUint64}}
	var valUint64 uint64
	{{- end}}
	{{- if $needsValInt8}}
	var valInt8 int8
	{{- end}}
	{{- if $needsValInt16}}
	var valInt16 int16
	{{- end}}
	{{- if $needsValInt32}}
	var valInt32 int32
	{{- end}}
	{{- if $needsValInt64}}
	var valInt64 int64
	{{- end}}
	{{- if $needsValHash}}
	var valHash Hash
	{{- end}}
	{{- if $needsValBytes1}}
	var valBytes1 [1]byte
	{{- end}}
	{{- if $needsValBytes32}}
	var valBytes32 [32]byte
	{{- end}}
	{{- if $needsValString}}
	var valString string
	{{- end}}
	{{- if $needsValBytes}}
	var valBytes []byte
	{{- end}}
	{{- if $needsNextOffset}}
	var nextOffset int
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
	{{- else}}
	val, err = decodeUint256(data[offset:offset+32])
	{{- end}}
	if err != nil {
		return result, fmt.Errorf("decoding return value {{$i}}: %w", err)
	}
	result.{{$output.Name | title}} = val
	offset += 32
	{{- else if eq $output.Type.TypeName "uint8"}}
	if len(data) < offset+32 {
		return result, errors.New("insufficient data for return value {{$i}}")
	}
	valUint8, err = decodeUint8(data[offset:offset+32])
	if err != nil {
		return result, fmt.Errorf("decoding return value {{$i}}: %w", err)
	}
	result.{{$output.Name | title}} = valUint8
	offset += 32
	{{- else if eq $output.Type.TypeName "uint16"}}
	if len(data) < offset+32 {
		return result, errors.New("insufficient data for return value {{$i}}")
	}
	valUint16, err = decodeUint16(data[offset:offset+32])
	if err != nil {
		return result, fmt.Errorf("decoding return value {{$i}}: %w", err)
	}
	result.{{$output.Name | title}} = valUint16
	offset += 32
	{{- else if eq $output.Type.TypeName "uint32"}}
	if len(data) < offset+32 {
		return result, errors.New("insufficient data for return value {{$i}}")
	}
	valUint32, err = decodeUint32(data[offset:offset+32])
	if err != nil {
		return result, fmt.Errorf("decoding return value {{$i}}: %w", err)
	}
	result.{{$output.Name | title}} = valUint32
	offset += 32
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
	{{- else if eq $output.Type.TypeName "int8"}}
	if len(data) < offset+32 {
		return result, errors.New("insufficient data for return value {{$i}}")
	}
	{
		raw, e := decodeInt64(data[offset:offset+32])
		if e != nil {
			return result, fmt.Errorf("decoding return value {{$i}}: %w", e)
		}
		valInt8 = int8(raw)
	}
	result.{{$output.Name | title}} = valInt8
	offset += 32
	{{- else if eq $output.Type.TypeName "int16"}}
	if len(data) < offset+32 {
		return result, errors.New("insufficient data for return value {{$i}}")
	}
	{
		raw, e := decodeInt64(data[offset:offset+32])
		if e != nil {
			return result, fmt.Errorf("decoding return value {{$i}}: %w", e)
		}
		valInt16 = int16(raw)
	}
	result.{{$output.Name | title}} = valInt16
	offset += 32
	{{- else if eq $output.Type.TypeName "int32"}}
	if len(data) < offset+32 {
		return result, errors.New("insufficient data for return value {{$i}}")
	}
	{
		raw, e := decodeInt64(data[offset:offset+32])
		if e != nil {
			return result, fmt.Errorf("decoding return value {{$i}}: %w", e)
		}
		valInt32 = int32(raw)
	}
	result.{{$output.Name | title}} = valInt32
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
	{{- else if eq $output.Type.TypeName "Hash"}}
	if len(data) < offset+32 {
		return result, errors.New("insufficient data for return value {{$i}}")
	}
	valHash, err = decodeHash(data[offset:offset+32])
	if err != nil {
		return result, fmt.Errorf("decoding return value {{$i}}: %w", err)
	}
	result.{{$output.Name | title}} = valHash
	offset += 32
	{{- else if eq $output.Type.TypeName "[1]byte"}}
	if len(data) < offset+32 {
		return result, errors.New("insufficient data for return value {{$i}}")
	}
	valBytes1, err = decodeBytes1(data[offset:offset+32])
	if err != nil {
		return result, fmt.Errorf("decoding return value {{$i}}: %w", err)
	}
	result.{{$output.Name | title}} = valBytes1
	offset += 32
	{{- else if eq $output.Type.TypeName "[32]byte"}}
	if len(data) < offset+32 {
		return result, errors.New("insufficient data for return value {{$i}}")
	}
	valBytes32, err = decodeBytes32(data[offset:offset+32])
	if err != nil {
		return result, fmt.Errorf("decoding return value {{$i}}: %w", err)
	}
	result.{{$output.Name | title}} = valBytes32
	offset += 32
	{{- else if eq $output.Type.TypeName "[]*big.Int"}}
	if len(data) < offset+32 {
		return result, errors.New("insufficient data for offset pointer in return value {{$i}}")
	}
	{
		ptrBig{{$i}}, e := decodeUint256(data[offset : offset+32])
		if e != nil {
			return result, fmt.Errorf("decoding offset pointer in return value {{$i}}: %w", e)
		}
		if !ptrBig{{$i}}.IsUint64() {
			return result, errors.New("offset pointer too large in return value {{$i}}")
		}
		elems, _, e2 := decodeArray(data, int(ptrBig{{$i}}.Uint64()), decodeUint256ArrayElement)
		if e2 != nil {
			return result, fmt.Errorf("decoding return value {{$i}}: %w", e2)
		}
		arr{{$i}} := make([]*big.Int, len(elems))
		for j, elem := range elems {
			arr{{$i}}[j] = elem.(*big.Int)
		}
		result.{{$output.Name | title}} = arr{{$i}}
	}
	offset += 32
	{{- else if eq $output.Type.TypeName "[]uint64"}}
	if len(data) < offset+32 {
		return result, errors.New("insufficient data for offset pointer in return value {{$i}}")
	}
	{
		ptrBig{{$i}}, e := decodeUint256(data[offset : offset+32])
		if e != nil {
			return result, fmt.Errorf("decoding offset pointer in return value {{$i}}: %w", e)
		}
		if !ptrBig{{$i}}.IsUint64() {
			return result, errors.New("offset pointer too large in return value {{$i}}")
		}
		elems, _, e2 := decodeArray(data, int(ptrBig{{$i}}.Uint64()), func(d []byte) (interface{}, error) { return decodeUint64(d) })
		if e2 != nil {
			return result, fmt.Errorf("decoding return value {{$i}}: %w", e2)
		}
		arr{{$i}} := make([]uint64, len(elems))
		for j, elem := range elems {
			arr{{$i}}[j] = elem.(uint64)
		}
		result.{{$output.Name | title}} = arr{{$i}}
	}
	offset += 32
	{{- else if eq $output.Type.TypeName "[]Address"}}
	if len(data) < offset+32 {
		return result, errors.New("insufficient data for offset pointer in return value {{$i}}")
	}
	{
		ptrBig{{$i}}, e := decodeUint256(data[offset : offset+32])
		if e != nil {
			return result, fmt.Errorf("decoding offset pointer in return value {{$i}}: %w", e)
		}
		if !ptrBig{{$i}}.IsUint64() {
			return result, errors.New("offset pointer too large in return value {{$i}}")
		}
		elems, _, e2 := decodeArray(data, int(ptrBig{{$i}}.Uint64()), decodeAddressArrayElement)
		if e2 != nil {
			return result, fmt.Errorf("decoding return value {{$i}}: %w", e2)
		}
		arr{{$i}} := make([]Address, len(elems))
		for j, elem := range elems {
			arr{{$i}}[j] = elem.(Address)
		}
		result.{{$output.Name | title}} = arr{{$i}}
	}
	offset += 32
	{{- else if eq $output.Type.TypeName "[]bool"}}
	if len(data) < offset+32 {
		return result, errors.New("insufficient data for offset pointer in return value {{$i}}")
	}
	{
		ptrBig{{$i}}, e := decodeUint256(data[offset : offset+32])
		if e != nil {
			return result, fmt.Errorf("decoding offset pointer in return value {{$i}}: %w", e)
		}
		if !ptrBig{{$i}}.IsUint64() {
			return result, errors.New("offset pointer too large in return value {{$i}}")
		}
		elems, _, e2 := decodeArray(data, int(ptrBig{{$i}}.Uint64()), decodeBoolArrayElement)
		if e2 != nil {
			return result, fmt.Errorf("decoding return value {{$i}}: %w", e2)
		}
		arr{{$i}} := make([]bool, len(elems))
		for j, elem := range elems {
			arr{{$i}}[j] = elem.(bool)
		}
		result.{{$output.Name | title}} = arr{{$i}}
	}
	offset += 32
	{{- else if eq $output.Type.TypeName "string"}}
	if len(data) < offset+32 {
		return result, errors.New("insufficient data for offset pointer in return value {{$i}}")
	}
	{
		ptrBig{{$i}}, e := decodeUint256(data[offset : offset+32])
		if e != nil {
			return result, fmt.Errorf("decoding offset pointer in return value {{$i}}: %w", e)
		}
		if !ptrBig{{$i}}.IsUint64() {
			return result, errors.New("offset pointer too large in return value {{$i}}")
		}
		var e2 error
		valString, _, e2 = decodeString(data, int(ptrBig{{$i}}.Uint64()))
		if e2 != nil {
			return result, fmt.Errorf("decoding return value {{$i}}: %w", e2)
		}
	}
	result.{{$output.Name | title}} = valString
	offset += 32
	{{- else if eq $output.Type.TypeName "[]byte"}}
	if len(data) < offset+32 {
		return result, errors.New("insufficient data for offset pointer in return value {{$i}}")
	}
	{
		ptrBig{{$i}}, e := decodeUint256(data[offset : offset+32])
		if e != nil {
			return result, fmt.Errorf("decoding offset pointer in return value {{$i}}: %w", e)
		}
		if !ptrBig{{$i}}.IsUint64() {
			return result, errors.New("offset pointer too large in return value {{$i}}")
		}
		var e2 error
		valBytes, _, e2 = decodeBytes(data, int(ptrBig{{$i}}.Uint64()))
		if e2 != nil {
			return result, fmt.Errorf("decoding return value {{$i}}: %w", e2)
		}
	}
	result.{{$output.Name | title}} = valBytes
	offset += 32
	{{- else if $output.Type.IsStruct}}
	{{- range $.Contract.Structs}}
	{{- if eq .Name $output.Type.TypeName}}
	{
		var sv {{.Name}}
		sv, nextOffset, err = decode{{.Name}}(data, offset)
		if err != nil {
			return result, fmt.Errorf("decoding return value {{$i}}: %w", err)
		}
		result.{{$output.Name | title}} = sv
		offset = nextOffset
	}
	{{- end}}
	{{- end}}
	{{- else if $output.Type.IsSlice}}
	{{- $elemType := slice $output.Type.TypeName 2}}
	{{- range $.Contract.Structs}}
	{{- if eq .Name $elemType}}
	// Struct array: read offset pointer then decode at pointed location
	if len(data) < offset+32 {
		return result, errors.New("insufficient data for array offset in return value {{$i}}")
	}
	{
		ptrVal{{$i}}, e := decodeUint256(data[offset:offset+32])
		if e != nil {
			return result, fmt.Errorf("decoding array offset in return value {{$i}}: %w", e)
		}
		if !ptrVal{{$i}}.IsUint64() {
			return result, errors.New("array offset too large in return value {{$i}}")
		}
		arrOff{{$i}} := int(ptrVal{{$i}}.Uint64())
		if len(data) < arrOff{{$i}}+32 {
			return result, errors.New("insufficient data for array length in return value {{$i}}")
		}
		lenVal{{$i}}, e2 := decodeUint256(data[arrOff{{$i}}:arrOff{{$i}}+32])
		if e2 != nil {
			return result, fmt.Errorf("decoding array length in return value {{$i}}: %w", e2)
		}
		if !lenVal{{$i}}.IsUint64() {
			return result, errors.New("array length too large in return value {{$i}}")
		}
		length{{$i}} := int(lenVal{{$i}}.Uint64())
		arr{{$i}} := make({{$output.Type.TypeName}}, length{{$i}})
		elemOff{{$i}} := arrOff{{$i}} + 32
		for j := 0; j < length{{$i}}; j++ {
			var elem {{.Name}}
			elem, elemOff{{$i}}, err = decode{{.Name}}(data, elemOff{{$i}})
			if err != nil {
				return result, fmt.Errorf("decoding array element %d in return value {{$i}}: %w", j, err)
			}
			arr{{$i}}[j] = elem
		}
		result.{{$output.Name | title}} = arr{{$i}}
	}
	offset += 32
	{{- end}}
	{{- end}}
	{{- else}}
	return result, errors.New("unsupported multi-return type: {{$output.Type.TypeName}}")
	{{- end}}
	{{- end}}
	return result, nil
{{- end}}
}
{{- end}}
{{- end}}`