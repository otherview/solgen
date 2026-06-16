// SPDX-License-Identifier: MIT

package gen

// decodeValueTemplate defines the shared "decodeValue" template used by the
// multi-return method decoder, the event decoder, and the error decoder. It
// emits self-contained statements that decode one ABI value from a data section
// and assign it to a struct field.
//
// Context dict keys:
//   - Data:  the []byte variable holding the ABI section ("data" / "errorData")
//   - Field: the assignment target ("result.Foo")
//   - T:     the value's types.GoType
//   - Ctx:   a human-readable label for error messages
//   - Zero:  the value returned alongside an error (always "result")
//
// All error/temp variables inside the emitted blocks are block-local; the only
// shared variable referenced is offset (advanced by the caller's loop). This
// keeps the generated decoder free of unused err/nextOffset declarations.
const decodeValueTemplate = `
{{- define "decodeValue" -}}
{{- if eq .T.TypeName "*big.Int"}}
if len({{.Data}}) < offset+32 {
	return {{.Zero}}, errors.New("insufficient data for {{.Ctx}}")
}
{
	{{- if .T.IsSigned}}
	v, e := decodeInt256({{.Data}}[offset : offset+32])
	{{- else}}
	v, e := decodeUint256({{.Data}}[offset : offset+32])
	{{- end}}
	if e != nil {
		return {{.Zero}}, fmt.Errorf("decoding {{.Ctx}}: %w", e)
	}
	{{.Field}} = v
}
offset += 32
{{- else if eq .T.TypeName "uint8"}}
if len({{.Data}}) < offset+32 {
	return {{.Zero}}, errors.New("insufficient data for {{.Ctx}}")
}
{
	v, e := decodeUint8({{.Data}}[offset : offset+32])
	if e != nil {
		return {{.Zero}}, fmt.Errorf("decoding {{.Ctx}}: %w", e)
	}
	{{.Field}} = v
}
offset += 32
{{- else if eq .T.TypeName "uint16"}}
if len({{.Data}}) < offset+32 {
	return {{.Zero}}, errors.New("insufficient data for {{.Ctx}}")
}
{
	v, e := decodeUint16({{.Data}}[offset : offset+32])
	if e != nil {
		return {{.Zero}}, fmt.Errorf("decoding {{.Ctx}}: %w", e)
	}
	{{.Field}} = v
}
offset += 32
{{- else if eq .T.TypeName "uint32"}}
if len({{.Data}}) < offset+32 {
	return {{.Zero}}, errors.New("insufficient data for {{.Ctx}}")
}
{
	v, e := decodeUint32({{.Data}}[offset : offset+32])
	if e != nil {
		return {{.Zero}}, fmt.Errorf("decoding {{.Ctx}}: %w", e)
	}
	{{.Field}} = v
}
offset += 32
{{- else if eq .T.TypeName "uint64"}}
if len({{.Data}}) < offset+32 {
	return {{.Zero}}, errors.New("insufficient data for {{.Ctx}}")
}
{
	v, e := decodeUint64({{.Data}}[offset : offset+32])
	if e != nil {
		return {{.Zero}}, fmt.Errorf("decoding {{.Ctx}}: %w", e)
	}
	{{.Field}} = v
}
offset += 32
{{- else if eq .T.TypeName "int8"}}
if len({{.Data}}) < offset+32 {
	return {{.Zero}}, errors.New("insufficient data for {{.Ctx}}")
}
{
	raw, e := decodeInt64({{.Data}}[offset : offset+32])
	if e != nil {
		return {{.Zero}}, fmt.Errorf("decoding {{.Ctx}}: %w", e)
	}
	{{.Field}} = int8(raw)
}
offset += 32
{{- else if eq .T.TypeName "int16"}}
if len({{.Data}}) < offset+32 {
	return {{.Zero}}, errors.New("insufficient data for {{.Ctx}}")
}
{
	raw, e := decodeInt64({{.Data}}[offset : offset+32])
	if e != nil {
		return {{.Zero}}, fmt.Errorf("decoding {{.Ctx}}: %w", e)
	}
	{{.Field}} = int16(raw)
}
offset += 32
{{- else if eq .T.TypeName "int32"}}
if len({{.Data}}) < offset+32 {
	return {{.Zero}}, errors.New("insufficient data for {{.Ctx}}")
}
{
	raw, e := decodeInt64({{.Data}}[offset : offset+32])
	if e != nil {
		return {{.Zero}}, fmt.Errorf("decoding {{.Ctx}}: %w", e)
	}
	{{.Field}} = int32(raw)
}
offset += 32
{{- else if eq .T.TypeName "int64"}}
if len({{.Data}}) < offset+32 {
	return {{.Zero}}, errors.New("insufficient data for {{.Ctx}}")
}
{
	v, e := decodeInt64({{.Data}}[offset : offset+32])
	if e != nil {
		return {{.Zero}}, fmt.Errorf("decoding {{.Ctx}}: %w", e)
	}
	{{.Field}} = v
}
offset += 32
{{- else if eq .T.TypeName "bool"}}
if len({{.Data}}) < offset+32 {
	return {{.Zero}}, errors.New("insufficient data for {{.Ctx}}")
}
{
	v, e := decodeBool({{.Data}}[offset : offset+32])
	if e != nil {
		return {{.Zero}}, fmt.Errorf("decoding {{.Ctx}}: %w", e)
	}
	{{.Field}} = v
}
offset += 32
{{- else if eq .T.TypeName "Address"}}
if len({{.Data}}) < offset+32 {
	return {{.Zero}}, errors.New("insufficient data for {{.Ctx}}")
}
{
	v, e := decodeAddress({{.Data}}[offset : offset+32])
	if e != nil {
		return {{.Zero}}, fmt.Errorf("decoding {{.Ctx}}: %w", e)
	}
	{{.Field}} = v
}
offset += 32
{{- else if eq .T.TypeName "Hash"}}
if len({{.Data}}) < offset+32 {
	return {{.Zero}}, errors.New("insufficient data for {{.Ctx}}")
}
{
	v, e := decodeHash({{.Data}}[offset : offset+32])
	if e != nil {
		return {{.Zero}}, fmt.Errorf("decoding {{.Ctx}}: %w", e)
	}
	{{.Field}} = v
}
offset += 32
{{- else if eq .T.TypeName "[1]byte"}}
if len({{.Data}}) < offset+32 {
	return {{.Zero}}, errors.New("insufficient data for {{.Ctx}}")
}
{
	v, e := decodeBytes1({{.Data}}[offset : offset+32])
	if e != nil {
		return {{.Zero}}, fmt.Errorf("decoding {{.Ctx}}: %w", e)
	}
	{{.Field}} = v
}
offset += 32
{{- else if eq .T.TypeName "[32]byte"}}
if len({{.Data}}) < offset+32 {
	return {{.Zero}}, errors.New("insufficient data for {{.Ctx}}")
}
{
	v, e := decodeBytes32({{.Data}}[offset : offset+32])
	if e != nil {
		return {{.Zero}}, fmt.Errorf("decoding {{.Ctx}}: %w", e)
	}
	{{.Field}} = v
}
offset += 32
{{- else if eq .T.TypeName "string"}}
{
	ptr, e := readOffset({{.Data}}, offset)
	if e != nil {
		return {{.Zero}}, fmt.Errorf("decoding {{.Ctx}}: %w", e)
	}
	s, _, e2 := decodeString({{.Data}}, ptr)
	if e2 != nil {
		return {{.Zero}}, fmt.Errorf("decoding {{.Ctx}}: %w", e2)
	}
	{{.Field}} = s
}
offset += 32
{{- else if eq .T.TypeName "[]byte"}}
{
	ptr, e := readOffset({{.Data}}, offset)
	if e != nil {
		return {{.Zero}}, fmt.Errorf("decoding {{.Ctx}}: %w", e)
	}
	b, _, e2 := decodeBytes({{.Data}}, ptr)
	if e2 != nil {
		return {{.Zero}}, fmt.Errorf("decoding {{.Ctx}}: %w", e2)
	}
	{{.Field}} = b
}
offset += 32
{{- else if .T.IsStruct}}
{
	v, no, e := decode{{.T.TypeName}}({{.Data}}, offset)
	if e != nil {
		return {{.Zero}}, fmt.Errorf("decoding {{.Ctx}}: %w", e)
	}
	{{.Field}} = v
	offset = no
}
{{- else if and .T.IsSlice .T.ElemType.IsStruct}}
{
	arrayOffset, e := readOffset({{.Data}}, offset)
	if e != nil {
		return {{.Zero}}, fmt.Errorf("decoding {{.Ctx}}: %w", e)
	}
	length, e2 := readOffset({{.Data}}, arrayOffset)
	if e2 != nil {
		return {{.Zero}}, fmt.Errorf("decoding {{.Ctx}} length: %w", e2)
	}
	arr := make({{.T.TypeName}}, length)
	elemOff := arrayOffset + 32
	for i := 0; i < length; i++ {
		v, no, e3 := decode{{.T.ElemType.TypeName}}({{.Data}}, elemOff)
		if e3 != nil {
			return {{.Zero}}, fmt.Errorf("decoding {{.Ctx}} element %d: %w", i, e3)
		}
		arr[i] = v
		elemOff = no
	}
	{{.Field}} = arr
}
offset += 32
{{- else if and .T.IsArray (staticWordDecoder .T.ElemType)}}
{
	tmp, e := decodeStaticFixedArray({{.Data}}, offset, {{.T.ArraySize}}, {{staticWordDecoder .T.ElemType}})
	if e != nil {
		return {{.Zero}}, fmt.Errorf("decoding {{.Ctx}}: %w", e)
	}
	var arr {{.T.TypeName}}
	copy(arr[:], tmp)
	{{.Field}} = arr
}
offset += 32
{{- else if and .T.IsSlice (staticWordDecoder .T.ElemType)}}
{
	s, e := decodeStaticSliceAt({{.Data}}, offset, {{staticWordDecoder .T.ElemType}})
	if e != nil {
		return {{.Zero}}, fmt.Errorf("decoding {{.Ctx}}: %w", e)
	}
	{{.Field}} = s
}
offset += 32
{{- end -}}
{{- end -}}
`
