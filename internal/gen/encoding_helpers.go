// SPDX-License-Identifier: MIT

package gen

// encodingHelpersTemplate contains the ABI encoding helper functions
const encodingHelpersTemplate = `// ABI Encoding Implementation

// encodeUint256 encodes a uint256 value to 32 bytes (big-endian)
func encodeUint256(val interface{}) ([]byte, error) {
	result := make([]byte, 32)
	switch v := val.(type) {
	case *big.Int:
		if v.Sign() < 0 {
			return nil, errors.New("negative values not supported for uint256")
		}
		if v.BitLen() > 256 {
			return nil, errors.New("value too large for uint256")
		}
		v.FillBytes(result)
		return result, nil
	case uint64:
		big.NewInt(0).SetUint64(v).FillBytes(result)
		return result, nil
	case int64:
		if v < 0 {
			return nil, errors.New("negative values not supported for uint256")
		}
		big.NewInt(v).FillBytes(result)
		return result, nil
	case int:
		if v < 0 {
			return nil, errors.New("negative values not supported for uint256")
		}
		big.NewInt(int64(v)).FillBytes(result)
		return result, nil
	default:
		return nil, fmt.Errorf("unsupported type for uint256: %T", v)
	}
}

// encodeInt256 encodes a signed 256-bit integer to 32 bytes using two's complement.
// Valid range: [-2^255, 2^255-1].
func encodeInt256(val interface{}) ([]byte, error) {
	result := make([]byte, 32)
	switch v := val.(type) {
	case *big.Int:
		if v.Sign() >= 0 {
			// Positive: valid range [0, 2^255-1] → BitLen must be ≤ 255.
			if v.BitLen() > 255 {
				return nil, errors.New("value too large for int256")
			}
			v.FillBytes(result)
		} else {
			// Negative: valid range [-2^255, -1].
			// abs(-2^255) has BitLen == 256, which is the boundary.
			abs := new(big.Int).Neg(v)
			minNeg := new(big.Int).Lsh(big.NewInt(1), 255) // 2^255
			if abs.Cmp(minNeg) > 0 {
				return nil, errors.New("value too small for int256")
			}
			// Two's-complement: compute 2^256 + v = 2^256 - abs(v).
			mask := new(big.Int).Lsh(big.NewInt(1), 256)
			new(big.Int).Add(mask, v).FillBytes(result)
		}
		return result, nil
	case int64:
		return encodeInt256(big.NewInt(v))
	case int:
		return encodeInt256(big.NewInt(int64(v)))
	default:
		return nil, fmt.Errorf("unsupported type for int256: %T", v)
	}
}

// encodeAddress encodes an address to 32 bytes (zero-padded)
func encodeAddress(addr Address) ([]byte, error) {
	result := make([]byte, 32)
	copy(result[12:32], addr[:])
	return result, nil
}

// encodeBool encodes a boolean to 32 bytes
func encodeBool(val bool) ([]byte, error) {
	result := make([]byte, 32)
	if val {
		result[31] = 1
	}
	return result, nil
}

// encodeBytes encodes dynamic bytes 
func encodeBytes(data []byte) ([]byte, error) {
	// Length (32 bytes) + data (padded to multiple of 32 bytes)
	length := len(data)
	lengthBytes, err := encodeUint256(uint64(length))
	if err != nil {
		return nil, err
	}
	
	// Pad data to multiple of 32 bytes
	paddedLength := ((length + 31) / 32) * 32
	paddedData := make([]byte, paddedLength)
	copy(paddedData, data)
	
	return append(lengthBytes, paddedData...), nil
}

// encodeString encodes a string as dynamic bytes
func encodeString(str string) ([]byte, error) {
	return encodeBytes([]byte(str))
}

// encodeFixedBytes encodes fixed-size bytes (e.g., bytes32) as a single static
// 32-byte word: the value is left-aligned and right-padded with zeros.
func encodeFixedBytes(val []byte, size int) ([]byte, error) {
	if size < 1 || size > 32 {
		return nil, fmt.Errorf("invalid fixed bytes size: %d", size)
	}
	if len(val) != size {
		return nil, fmt.Errorf("fixed bytes length mismatch: got %d, want %d", len(val), size)
	}
	result := make([]byte, 32)
	copy(result, val)
	return result, nil
}

// encodeArg encodes a single ABI argument, returning its encoded bytes and
// whether it is a dynamic type (which gets a 32-byte offset pointer in the head
// and its data in the tail). It is the per-argument core shared by Pack.
func encodeArg(arg any) ([]byte, bool, error) {
	switch v := arg.(type) {
	case *big.Int:
		if v.Sign() < 0 {
			d, err := encodeInt256(v)
			return d, false, err
		}
		d, err := encodeUint256(v)
		return d, false, err
	case uint8:
		d, err := encodeUint256(uint64(v))
		return d, false, err
	case uint16:
		d, err := encodeUint256(uint64(v))
		return d, false, err
	case uint32:
		d, err := encodeUint256(uint64(v))
		return d, false, err
	case uint64:
		d, err := encodeUint256(v)
		return d, false, err
	case int8:
		d, err := encodeInt256(big.NewInt(int64(v)))
		return d, false, err
	case int16:
		d, err := encodeInt256(big.NewInt(int64(v)))
		return d, false, err
	case int32:
		d, err := encodeInt256(big.NewInt(int64(v)))
		return d, false, err
	case int64:
		d, err := encodeInt256(big.NewInt(v))
		return d, false, err
	case Address:
		d, err := encodeAddress(v)
		return d, false, err
	case bool:
		d, err := encodeBool(v)
		return d, false, err
	case string:
		d, err := encodeString(v)
		return d, true, err
	case []byte:
		d, err := encodeBytes(v)
		return d, true, err
	case Hash:
		d, err := encodeFixedBytes(v[:], 32)
		return d, false, err
	case [32]byte:
		d, err := encodeFixedBytes(v[:], 32)
		return d, false, err
	default:
		rt := reflect.TypeOf(arg)
		if rt != nil && rt.Kind() == reflect.Array {
			// Fixed-size byte array bytesN (1 <= N <= 32): one static word.
			if rt.Elem().Kind() == reflect.Uint8 && rt.Len() <= 32 {
				rv := reflect.ValueOf(arg)
				b := make([]byte, rv.Len())
				reflect.Copy(reflect.ValueOf(b), rv)
				d, err := encodeFixedBytes(b, len(b))
				return d, false, err
			}
			// Fixed-size array [N]T of static elements: N inline static words.
			rv := reflect.ValueOf(arg)
			var out []byte
			for i := 0; i < rv.Len(); i++ {
				elemData, elemDynamic, err := encodeArg(rv.Index(i).Interface())
				if err != nil {
					return nil, false, fmt.Errorf("encoding array element %d: %w", i, err)
				}
				if elemDynamic {
					return nil, false, fmt.Errorf("unsupported dynamic element in fixed-size array: %T", arg)
				}
				out = append(out, elemData...)
			}
			return out, false, nil
		}
		return nil, false, fmt.Errorf("unsupported argument type: %T", arg)
	}
}`