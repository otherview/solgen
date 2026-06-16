// SPDX-License-Identifier: MIT

package gen

// decodingHelpersTemplate contains the ABI decoding helper functions
const decodingHelpersTemplate = `// ABI Decoding Implementation

// decodeUint256 decodes a uint256 from 32 bytes to *big.Int
func decodeUint256(data []byte) (*big.Int, error) {
	if len(data) < 32 {
		return nil, errors.New("insufficient data for uint256")
	}
	return new(big.Int).SetBytes(data[:32]), nil
}

// decodeInt256 decodes a signed 256-bit integer from 32 bytes
func decodeInt256(data []byte) (*big.Int, error) {
	if len(data) < 32 {
		return nil, errors.New("insufficient data for int256")
	}
	
	result := new(big.Int).SetBytes(data[:32])
	
	// Check if negative (MSB is set)
	if data[0]&0x80 != 0 {
		// Convert from two's complement
		// Create mask with all bits set for 256-bit number
		mask := new(big.Int).Lsh(big.NewInt(1), 256)
		mask.Sub(mask, big.NewInt(1))
		
		// XOR with mask and add 1 to get absolute value
		result.Xor(result, mask)
		result.Add(result, big.NewInt(1))
		result.Neg(result)
	}
	
	return result, nil
}

// decodeAddress decodes an address from 32 bytes
func decodeAddress(data []byte) (Address, error) {
	if len(data) < 32 {
		return Address{}, errors.New("insufficient data for address")
	}
	var addr Address
	copy(addr[:], data[12:32])
	return addr, nil
}

// decodeBool decodes a boolean from 32 bytes
func decodeBool(data []byte) (bool, error) {
	if len(data) < 32 {
		return false, errors.New("insufficient data for bool")
	}
	return data[31] != 0, nil
}

// decodeBytes decodes dynamic bytes
func decodeBytes(data []byte, offset int) ([]byte, int, error) {
	if len(data) < offset+32 {
		return nil, 0, errors.New("insufficient data for bytes length")
	}
	lengthBig, err := decodeUint256(data[offset : offset+32])
	if err != nil {
		return nil, 0, fmt.Errorf("decoding bytes length: %w", err)
	}
	if !lengthBig.IsUint64() {
		return nil, 0, errors.New("bytes length too large")
	}
	length := int(lengthBig.Uint64())
	if len(data) < offset+32+length {
		return nil, 0, errors.New("insufficient data for bytes content")
	}
	result := make([]byte, length)
	copy(result, data[offset+32:offset+32+length])
	// Calculate next offset (padded to 32 bytes)
	paddedLength := ((length + 31) / 32) * 32
	return result, offset + 32 + paddedLength, nil
}

// decodeFixedBytes decodes fixed-size bytes (e.g., bytes32)
func decodeFixedBytes(data []byte, size int) ([]byte, error) {
	if len(data) < 32 {
		return nil, errors.New("insufficient data for fixed bytes")
	}
	if size > 32 {
		return nil, errors.New("fixed bytes size too large")
	}
	result := make([]byte, size)
	copy(result, data[:size])
	return result, nil
}

// decode various fixed-size byte arrays
func decodeBytes1(data []byte) ([1]byte, error) {
	bytes, err := decodeFixedBytes(data, 1)
	if err != nil {
		return [1]byte{}, err
	}
	var result [1]byte
	copy(result[:], bytes)
	return result, nil
}

func decodeBytes32(data []byte) ([32]byte, error) {
	bytes, err := decodeFixedBytes(data, 32)
	if err != nil {
		return [32]byte{}, err
	}
	var result [32]byte
	copy(result[:], bytes)
	return result, nil
}

// decodeArray decodes dynamic arrays 
func decodeArray(data []byte, offset int, elemDecoder func([]byte) (interface{}, error)) ([]interface{}, int, error) {
	if len(data) < offset+32 {
		return nil, 0, errors.New("insufficient data for array length")
	}
	
	lengthBig, err := decodeUint256(data[offset : offset+32])
	if err != nil {
		return nil, 0, fmt.Errorf("decoding array length: %w", err)
	}
	if !lengthBig.IsUint64() {
		return nil, 0, errors.New("array length too large")
	}
	length := int(lengthBig.Uint64())
	
	currentOffset := offset + 32
	result := make([]interface{}, length)
	
	for i := 0; i < length; i++ {
		if len(data) < currentOffset+32 {
			return nil, 0, fmt.Errorf("insufficient data for array element %d", i)
		}
		elem, err := elemDecoder(data[currentOffset : currentOffset+32])
		if err != nil {
			return nil, 0, fmt.Errorf("decoding array element %d: %w", i, err)
		}
		result[i] = elem
		currentOffset += 32
	}
	
	return result, currentOffset, nil
}

// Array element decoders (internal use)
func decodeUint256ArrayElement(data []byte) (interface{}, error) {
	return decodeUint256(data)
}

func decodeInt256ArrayElement(data []byte) (interface{}, error) {
	return decodeInt256(data)
}

func decodeAddressArrayElement(data []byte) (interface{}, error) {
	return decodeAddress(data)
}

func decodeBoolArrayElement(data []byte) (interface{}, error) {
	return decodeBool(data)
}

// readOffset reads a 32-byte ABI offset/length word at head position pos and
// returns it as an int.
func readOffset(data []byte, pos int) (int, error) {
	if len(data) < pos+32 {
		return 0, errors.New("insufficient data for offset pointer")
	}
	p, err := decodeUint256(data[pos : pos+32])
	if err != nil {
		return 0, fmt.Errorf("decoding offset pointer: %w", err)
	}
	if !p.IsUint64() {
		return 0, errors.New("offset pointer too large")
	}
	return int(p.Uint64()), nil
}

// decodeStaticSliceAt decodes a dynamic array of static (32-byte) elements. pos
// is the head position holding the offset pointer to the array data; dec decodes
// a single element from its 32-byte word.
func decodeStaticSliceAt[T any](data []byte, pos int, dec func([]byte) (T, error)) ([]T, error) {
	arrPos, err := readOffset(data, pos)
	if err != nil {
		return nil, err
	}
	length, err := readOffset(data, arrPos)
	if err != nil {
		return nil, fmt.Errorf("decoding array length: %w", err)
	}
	out := make([]T, length)
	cur := arrPos + 32
	for i := 0; i < length; i++ {
		if len(data) < cur+32 {
			return nil, fmt.Errorf("insufficient data for array element %d", i)
		}
		v, err := dec(data[cur : cur+32])
		if err != nil {
			return nil, fmt.Errorf("decoding array element %d: %w", i, err)
		}
		out[i] = v
		cur += 32
	}
	return out, nil
}

// decodeStaticFixedArray decodes size consecutive static (32-byte) elements
// starting at pos into a slice; the caller copies it into the fixed-size [N]T
// value. dec decodes a single element from its 32-byte word.
func decodeStaticFixedArray[T any](data []byte, pos, size int, dec func([]byte) (T, error)) ([]T, error) {
	out := make([]T, size)
	for i := 0; i < size; i++ {
		if len(data) < pos+(i+1)*32 {
			return nil, fmt.Errorf("insufficient data for fixed array element %d", i)
		}
		v, err := dec(data[pos+i*32 : pos+(i+1)*32])
		if err != nil {
			return nil, fmt.Errorf("decoding fixed array element %d: %w", i, err)
		}
		out[i] = v
	}
	return out, nil
}

// decodeUint8 decodes a uint8 from 32 bytes
func decodeUint8(data []byte) (uint8, error) {
	if len(data) < 32 {
		return 0, errors.New("insufficient data for uint8")
	}
	// Verify upper bytes are zero
	for i := 0; i < 31; i++ {
		if data[i] != 0 {
			return 0, errors.New("invalid uint8 encoding")
		}
	}
	return data[31], nil
}

// decodeUint16 decodes a uint16 from 32 bytes
func decodeUint16(data []byte) (uint16, error) {
	if len(data) < 32 {
		return 0, errors.New("insufficient data for uint16")
	}
	// Verify upper bytes are zero
	for i := 0; i < 30; i++ {
		if data[i] != 0 {
			return 0, errors.New("invalid uint16 encoding")
		}
	}
	return uint16(data[30])<<8 | uint16(data[31]), nil
}

// decodeUint32 decodes a uint32 from 32 bytes
func decodeUint32(data []byte) (uint32, error) {
	if len(data) < 32 {
		return 0, errors.New("insufficient data for uint32")
	}
	// Verify upper bytes are zero
	for i := 0; i < 28; i++ {
		if data[i] != 0 {
			return 0, errors.New("invalid uint32 encoding")
		}
	}
	var result uint32
	for i := 28; i < 32; i++ {
		result = (result << 8) | uint32(data[i])
	}
	return result, nil
}

// decodeUint64 decodes a uint64 from 32 bytes  
func decodeUint64(data []byte) (uint64, error) {
	if len(data) < 32 {
		return 0, errors.New("insufficient data for uint64")
	}
	// Check if value exceeds uint64 range
	for i := 0; i < 24; i++ {
		if data[i] != 0 {
			return 0, errors.New("value exceeds uint64 range")
		}
	}
	var result uint64
	for i := 24; i < 32; i++ {
		result = (result << 8) | uint64(data[i])
	}
	return result, nil
}

// decodeInt64 decodes an int64 from 32 bytes (ABI sign-extended big-endian).
func decodeInt64(data []byte) (int64, error) {
	if len(data) < 32 {
		return 0, errors.New("insufficient data for int64")
	}

	// ABI sign-extension: bytes 0-23 must all match the sign byte
	// (0x00 for non-negative, 0xFF for negative).
	isNegative := data[0]&0x80 != 0
	expectedByte := byte(0)
	if isNegative {
		expectedByte = 0xFF
	}
	for i := 0; i < 24; i++ {
		if data[i] != expectedByte {
			return 0, errors.New("value exceeds int64 range")
		}
	}

	// Assemble the int64 from the last 8 bytes.
	// Because data[24..31] already hold the correct two's-complement
	// representation, no further sign extension is needed.
	var result int64
	for i := 24; i < 32; i++ {
		result = (result << 8) | int64(data[i])
	}
	return result, nil
}

// decodeHash decodes a 32-byte hash
func decodeHash(data []byte) (Hash, error) {
	if len(data) < 32 {
		return Hash{}, errors.New("insufficient data for hash")
	}
	var hash Hash
	copy(hash[:], data[:32])
	return hash, nil
}

// decodeString decodes a string from dynamic bytes
func decodeString(data []byte, offset int) (string, int, error) {
	bytes, nextOffset, err := decodeBytes(data, offset)
	if err != nil {
		return "", 0, err
	}
	return string(bytes), nextOffset, nil
}`