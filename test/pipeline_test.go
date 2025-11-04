// SPDX-License-Identifier: MIT

package test

import (
	"encoding/json"
	"testing"

	"github.com/otherview/solgen/internal/types"
)

func TestProcessCombinedJSON(t *testing.T) {
	// Mock combined JSON output from solc
	combinedJSONStr := `{
		"contracts": {
			"SimpleToken.sol:SimpleToken": {
				"abi": [
					{
						"type": "function",
						"name": "transfer",
						"inputs": [
							{"name": "to", "type": "address"},
							{"name": "amount", "type": "uint256"}
						],
						"outputs": [
							{"name": "", "type": "bool"}
						]
					},
					{
						"type": "event",
						"name": "Transfer",
						"inputs": [
							{"name": "from", "type": "address", "indexed": true},
							{"name": "to", "type": "address", "indexed": true},
							{"name": "value", "type": "uint256", "indexed": false}
						]
					}
				],
				"bin": "0x608060405234801561001057600080fd5b50...",
				"bin-runtime": "0x608060405234801561001057600080fd5b50...",
				"hashes": {
					"transfer(address,uint256)": "a9059cbb"
				}
			}
		},
		"version": "0.8.20+commit.a1b79de6.Linux.g++"
	}`

	// Test processing the combined JSON
	contracts, err := processCombinedJSON([]byte(combinedJSONStr))
	if err != nil {
		t.Fatalf("processCombinedJSON failed: %v", err)
	}

	// Validate results
	if len(contracts) != 1 {
		t.Fatalf("expected 1 contract, got %d", len(contracts))
	}

	contract := contracts[0]

	// Check basic properties
	if contract.Name != "SimpleToken" {
		t.Errorf("expected contract name 'SimpleToken', got %q", contract.Name)
	}

	if contract.PackageName != "simpletoken" {
		t.Errorf("expected package name 'simpletoken', got %q", contract.PackageName)
	}

	if contract.SourceFile != "SimpleToken.sol" {
		t.Errorf("expected source file 'SimpleToken.sol', got %q", contract.SourceFile)
	}

	// Check methods
	if len(contract.Methods) != 1 {
		t.Fatalf("expected 1 method, got %d", len(contract.Methods))
	}

	method := contract.Methods[0]
	if method.Name != "transfer" {
		t.Errorf("expected method name 'transfer', got %q", method.Name)
	}

	if method.Selector != "0xa9059cbb" {
		t.Errorf("expected selector '0xa9059cbb', got %q", method.Selector)
	}

	// Check events
	if len(contract.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(contract.Events))
	}

	event := contract.Events[0]
	if event.Name != "Transfer" {
		t.Errorf("expected event name 'Transfer', got %q", event.Name)
	}

	// Check bytecode
	if contract.Bytecode == "" {
		t.Error("expected non-empty bytecode")
	}

	if contract.DeployedBytecode == "" {
		t.Error("expected non-empty deployed bytecode")
	}
}

func TestConvertCombinedToStandard(t *testing.T) {
	// Create test combined JSON
	combined := types.CombinedJSON{
		Contracts: map[string]types.CombinedContract{
			"Test.sol:TestContract": {
				ABI:        json.RawMessage(`[{"type":"function","name":"test"}]`),
				Bin:        "0x1234",
				BinRuntime: "0x5678",
				Hashes: map[string]string{
					"test()": "12345678",
				},
			},
		},
	}

	// Convert to standard format
	result, err := convertCombinedToStandard(combined)
	if err != nil {
		t.Fatalf("convertCombinedToStandard failed: %v", err)
	}

	// Validate structure
	if len(result.Contracts) != 1 {
		t.Fatalf("expected 1 source file, got %d", len(result.Contracts))
	}

	sourceContracts, exists := result.Contracts["Test.sol"]
	if !exists {
		t.Fatal("expected 'Test.sol' in contracts")
	}

	if len(sourceContracts) != 1 {
		t.Fatalf("expected 1 contract in Test.sol, got %d", len(sourceContracts))
	}

	contract, exists := sourceContracts["TestContract"]
	if !exists {
		t.Fatal("expected 'TestContract' in Test.sol contracts")
	}

	// Validate contract properties
	if string(contract.ABI) != `[{"type":"function","name":"test"}]` {
		t.Errorf("unexpected ABI: %s", string(contract.ABI))
	}

	if contract.EVM.Bytecode.Object != "0x1234" {
		t.Errorf("expected bytecode '0x1234', got %q", contract.EVM.Bytecode.Object)
	}

	if contract.EVM.DeployedBytecode.Object != "0x5678" {
		t.Errorf("expected deployed bytecode '0x5678', got %q", contract.EVM.DeployedBytecode.Object)
	}

	// Metadata field removed from contract structure

	// Check method identifiers
	expectedHash := "12345678"
	actualHash, exists := contract.EVM.MethodIdentifiers["test()"]
	if !exists {
		t.Error("expected 'test()' in method identifiers")
	} else if actualHash != expectedHash {
		t.Errorf("expected hash %q, got %q", expectedHash, actualHash)
	}
}

func TestInvalidCombinedJSON(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{
			name: "invalid JSON",
			data: `{invalid json}`,
		},
		{
			name: "invalid contract key format",
			data: `{"contracts": {"InvalidKey": {"abi": [], "bin": "0x", "bin-runtime": "0x"}}}`,
		},
		{
			name: "empty data",
			data: `{"contracts": {}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := processCombinedJSON([]byte(tt.data))
			if err == nil && tt.name != "empty data" {
				t.Error("expected error for invalid input")
			}
		})
	}
}

// Test with the actual combined JSON types from main.go
func TestCombinedJSONTypes(t *testing.T) {
	// Test that our types match what solc outputs
	jsonData := []byte(`{
		"contracts": {
			"contract.sol:Contract": {
				"abi": [{"type": "function", "name": "test"}],
				"bin": "0x1234",
				"bin-runtime": "0x5678", 
				"hashes": {"test()": "abcd1234"}
			}
		},
		"version": "0.8.20+commit.a1b79de6.Linux.g++"
	}`)

	var combined types.CombinedJSON
	if err := json.Unmarshal(jsonData, &combined); err != nil {
		t.Fatalf("failed to unmarshal combined JSON: %v", err)
	}

	contract := combined.Contracts["contract.sol:Contract"]
	
	// Verify ABI is parsed as JSON
	var abiArray []interface{}
	if err := json.Unmarshal(contract.ABI, &abiArray); err != nil {
		t.Errorf("ABI should be valid JSON: %v", err)
	}

	// Verify other fields
	if contract.Bin != "0x1234" {
		t.Errorf("expected bin 0x1234, got %s", contract.Bin)
	}

	if contract.BinRuntime != "0x5678" {
		t.Errorf("expected bin-runtime 0x5678, got %s", contract.BinRuntime)
	}

	if len(contract.Hashes) != 1 {
		t.Errorf("expected 1 hash, got %d", len(contract.Hashes))
	}
}