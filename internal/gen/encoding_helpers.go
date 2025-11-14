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

// encodeInt256 encodes a signed 256-bit integer to 32 bytes using two's complement
func encodeInt256(val interface{}) ([]byte, error) {
	result := make([]byte, 32)
	switch v := val.(type) {
	case *big.Int:
		// Check if value fits in 256 bits (considering sign)
		if v.BitLen() >= 256 {
			return nil, errors.New("value too large for int256")
		}
		
		if v.Sign() >= 0 {
			// Positive number - same as uint256
			v.FillBytes(result)
		} else {
			// Negative number - use two's complement
			// Create a 256-bit mask (all 1s)
			mask := new(big.Int).Lsh(big.NewInt(1), 256)
			mask.Sub(mask, big.NewInt(1))
			
			// Get absolute value, subtract 1, XOR with mask
			abs := new(big.Int).Neg(v)
			abs.Sub(abs, big.NewInt(1))
			abs.Xor(abs, mask)
			abs.FillBytes(result)
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
}`