// SPDX-License-Identifier: MIT

package parse

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func TestStructArraySupport(t *testing.T) {
	// Test data simulating a contract with struct arrays
	abiJSON := `[
		{
			"inputs": [
				{
					"components": [
						{"internalType": "uint256", "name": "id", "type": "uint256"},
						{"internalType": "address", "name": "wallet", "type": "address"},
						{"internalType": "bool", "name": "active", "type": "bool"}
					],
					"internalType": "struct TestContract.User[]",
					"name": "users",
					"type": "tuple[]"
				}
			],
			"name": "processUsers",
			"outputs": [
				{
					"components": [
						{"internalType": "uint256", "name": "id", "type": "uint256"},
						{"internalType": "address", "name": "wallet", "type": "address"},
						{"internalType": "bool", "name": "active", "type": "bool"}
					],
					"internalType": "struct TestContract.User[]",
					"name": "",
					"type": "tuple[]"
				}
			],
			"stateMutability": "pure",
			"type": "function"
		}
	]`

	// Parse ABI
	parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		t.Fatalf("failed to parse ABI: %v", err)
	}

	// Create struct registry
	registry := newStructRegistry()

	// Test parseMethodsWithRegistry
	methodIds := map[string]string{
		"processUsers((uint256,address,bool)[])": "12345678",
	}

	methods, err := parseMethodsWithRegistry(parsedABI, methodIds, registry)
	if err != nil {
		t.Fatalf("parseMethodsWithRegistry failed: %v", err)
	}

	// Test that struct arrays are properly handled
	if len(methods) != 1 {
		t.Fatalf("expected 1 method, got %d", len(methods))
	}

	method := methods[0]
	if method.Name != "processUsers" {
		t.Errorf("expected method name 'processUsers', got '%s'", method.Name)
	}

	// Test input parameter type
	if len(method.Inputs) != 1 {
		t.Fatalf("expected 1 input, got %d", len(method.Inputs))
	}

	input := method.Inputs[0]
	if input.Type.TypeName != "[]User" {
		t.Errorf("expected input type '[]User', got '%s'", input.Type.TypeName)
	}

	if !input.Type.IsSlice {
		t.Error("input type should be marked as slice")
	}

	// Test output parameter type
	if len(method.Outputs) != 1 {
		t.Fatalf("expected 1 output, got %d", len(method.Outputs))
	}

	output := method.Outputs[0]
	if output.Type.TypeName != "[]User" {
		t.Errorf("expected output type '[]User', got '%s'", output.Type.TypeName)
	}

	// Test struct registration
	structs := registry.getAllStructs()
	if len(structs) != 1 {
		t.Fatalf("expected 1 registered struct, got %d", len(structs))
	}

	userStruct := structs[0]
	if userStruct.Name != "User" {
		t.Errorf("expected struct name 'User', got '%s'", userStruct.Name)
	}

	if len(userStruct.Fields) != 3 {
		t.Fatalf("expected 3 struct fields, got %d", len(userStruct.Fields))
	}

	// Verify field types
	expectedFields := []struct {
		name     string
		typeName string
	}{
		{"Id", "*big.Int"},
		{"Wallet", "Address"},
		{"Active", "bool"},
	}

	for i, expected := range expectedFields {
		if userStruct.Fields[i].Name != expected.name {
			t.Errorf("field %d: expected name '%s', got '%s'", i, expected.name, userStruct.Fields[i].Name)
		}
		if userStruct.Fields[i].Type.TypeName != expected.typeName {
			t.Errorf("field %d: expected type '%s', got '%s'", i, expected.typeName, userStruct.Fields[i].Type.TypeName)
		}
	}
}

func TestStructNameExtraction(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"TestContractUser", "User"},
		{"struct TestContract.User", "User"},
		{"MyContractCompany", "Company"}, // extract struct name from compound
		{"", ""},
	}

	for _, tc := range testCases {
		result := extractStructName(tc.input)
		if result != tc.expected {
			t.Errorf("extractStructName(%q): expected %q, got %q", tc.input, tc.expected, result)
		}
	}
}

func TestNoStructArrayRegression(t *testing.T) {
	// Ensure we don't break non-struct array types
	registry := newStructRegistry()

	// Test basic types don't get mis-registered
	basicTypes := []abi.Type{
		{T: abi.UintTy, Size: 256},
		{T: abi.AddressTy},
		{T: abi.BoolTy},
	}

	for _, abiType := range basicTypes {
		goType, err := mapSolidityToGoTypeWithRegistry(abiType, registry)
		if err != nil {
			t.Errorf("unexpected error for basic type: %v", err)
		}

		// Verify the type name is correct and no structs were registered
		switch abiType.T {
		case abi.UintTy:
			if goType.TypeName != "*big.Int" {
				t.Errorf("expected '*big.Int', got '%s'", goType.TypeName)
			}
		case abi.AddressTy:
			if goType.TypeName != "Address" {
				t.Errorf("expected 'Address', got '%s'", goType.TypeName)
			}
		case abi.BoolTy:
			if goType.TypeName != "bool" {
				t.Errorf("expected 'bool', got '%s'", goType.TypeName)
			}
		}
	}

	// No structs should be registered
	structs := registry.getAllStructs()
	if len(structs) != 0 {
		t.Errorf("expected no structs registered, got %d", len(structs))
	}
}