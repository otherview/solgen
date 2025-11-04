// SPDX-License-Identifier: MIT

package parse

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/otherview/solgen/internal/types"
)

func TestSanitizePackageName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"SimpleToken", "simpletoken"},
		{"MyContract_V2", "mycontractv2"},
		{"ERC20Token", "erc20token"},
		{"123Contract", "contract123contract"},
		{"", "contract"},
		{"TestContract-2.0", "testcontract20"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := sanitizePackageName(tt.input)
			if got != tt.want {
				t.Errorf("sanitizePackageName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestSanitizeIdentifier(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"validName", "validName"},
		{"_validName", "_validName"},
		{"123invalid", "_23invalid"},
		{"", "Field"},
		{"with-spaces", "with_spaces"},
		{"special@chars", "special_chars"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := sanitizeIdentifier(tt.input)
			if got != tt.want {
				t.Errorf("sanitizeIdentifier(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestMapSolidityToGoType(t *testing.T) {
	tests := []struct {
		name    string
		abiType abi.Type
		want    types.GoType
		wantErr bool
	}{
		{
			name:    "bool",
			abiType: abi.Type{T: abi.BoolTy},
			want:    types.GoTypeBool,
		},
		{
			name:    "string",
			abiType: abi.Type{T: abi.StringTy},
			want:    types.GoTypeString,
		},
		{
			name:    "address",
			abiType: abi.Type{T: abi.AddressTy},
			want:    types.GoTypeAddress,
		},
		{
			name:    "uint8",
			abiType: abi.Type{T: abi.UintTy, Size: 8},
			want:    types.GoTypeUint8,
		},
		{
			name:    "uint256",
			abiType: abi.Type{T: abi.UintTy, Size: 256},
			want:    types.GoTypeBigInt,
		},
		{
			name:    "bytes32",
			abiType: abi.Type{T: abi.FixedBytesTy, Size: 32},
			want:    types.GoType{TypeName: "[32]byte"},
		},
		{
			name:    "bytes",
			abiType: abi.Type{T: abi.BytesTy},
			want:    types.GoTypeBytes,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapSolidityToGoType(tt.abiType)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.TypeName != tt.want.TypeName {
				t.Errorf("mapSolidityToGoType() TypeName = %v, want %v", got.TypeName, tt.want.TypeName)
			}
			if got.Import != tt.want.Import {
				t.Errorf("mapSolidityToGoType() Import = %v, want %v", got.Import, tt.want.Import)
			}
		})
	}
}

func TestGenerateOverloadName(t *testing.T) {
	tests := []struct {
		baseName  string
		signature string
		selector  string
		want      string
	}{
		{
			baseName:  "transfer",
			signature: "transfer(address,uint256)",
			selector:  "0xa9059cbb",
			want:      "transfer_Address_Uint256",
		},
		{
			baseName:  "foo",
			signature: "foo()",
			selector:  "0x12345678",
			want:      "foo_NoArgs",
		},
		{
			baseName:  "complex",
			signature: "complex(uint256[],address[],bool)",
			selector:  "0xabcdef12",
			want:      "complex_Uint256Array_AddressArray_Bool",
		},
	}

	for _, tt := range tests {
		t.Run(tt.baseName, func(t *testing.T) {
			got := generateOverloadName(tt.baseName, tt.signature, tt.selector)
			if got != tt.want {
				t.Errorf("generateOverloadName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalizeTypeForNaming(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"uint256", "Uint256"},
		{"address", "Address"},
		{"bool", "Bool"},
		{"string", "String"},
		{"bytes", "Bytes"},
		{"bytes32", "Bytes32"},
		{"uint256[]", "Uint256Array"},
		{"address[10]", "AddressFixedArray"},
		{"int128", "Int128"},
		{"customType", "CustomType"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := normalizeTypeForNaming(tt.input)
			if got != tt.want {
				t.Errorf("normalizeTypeForNaming(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}