// SPDX-License-Identifier: MIT

package parse

import (
	"encoding/json"
	"testing"

	"github.com/otherview/solgen/internal/types"
)

func TestParseResult_Golden(t *testing.T) {
	// This is a golden test using a pre-compiled result from SimpleToken.sol
	// This allows us to test the parsing without requiring Docker
	goldenResult := &types.CompileResult{
		Contracts: map[string]map[string]types.ContractResult{
			"SimpleToken.sol": {
				"SimpleToken": types.ContractResult{
					ABI: json.RawMessage(`[
						{
							"type": "constructor",
							"inputs": [
								{"name": "_name", "type": "string"},
								{"name": "_symbol", "type": "string"},
								{"name": "_totalSupply", "type": "uint256"}
							]
						},
						{
							"type": "function",
							"name": "transfer",
							"inputs": [
								{"name": "to", "type": "address"},
								{"name": "value", "type": "uint256"}
							],
							"outputs": [
								{"name": "", "type": "bool"}
							],
							"stateMutability": "nonpayable"
						},
						{
							"type": "function", 
							"name": "balanceOf",
							"inputs": [
								{"name": "", "type": "address"}
							],
							"outputs": [
								{"name": "", "type": "uint256"}
							],
							"stateMutability": "view"
						},
						{
							"type": "event",
							"name": "Transfer",
							"inputs": [
								{"name": "from", "type": "address", "indexed": true},
								{"name": "to", "type": "address", "indexed": true},
								{"name": "value", "type": "uint256", "indexed": false}
							]
						},
						{
							"type": "error",
							"name": "InsufficientBalance",
							"inputs": [
								{"name": "account", "type": "address"},
								{"name": "requested", "type": "uint256"},
								{"name": "available", "type": "uint256"}
							]
						}
					]`),
					EVM: types.EVMResult{
						Bytecode: types.BytecodeResult{
							Object: "608060405234801561001057600080fd5b506040516108013803806108018339818101604052810190610032919061018b565b",
						},
						DeployedBytecode: types.BytecodeResult{
							Object: "608060405234801561001057600080fd5b50600436106100575760003560e01c8063095ea7b31461005c57806318160ddd1461008c",
						},
						MethodIdentifiers: map[string]string{
							"transfer(address,uint256)": "a9059cbb",
							"balanceOf(address)":        "70a08231",
						},
					},
				},
			},
		},
	}

	contracts, err := ResultWithVersion(goldenResult, "0.8.20")
	if err != nil {
		t.Fatalf("ParseResult failed: %v", err)
	}

	if len(contracts) != 1 {
		t.Fatalf("expected 1 contract, got %d", len(contracts))
	}

	contract := contracts[0]

	// Validate contract metadata
	if contract.Name != "SimpleToken" {
		t.Errorf("expected name 'SimpleToken', got %q", contract.Name)
	}
	if contract.PackageName != "simpletoken" {
		t.Errorf("expected package name 'simpletoken', got %q", contract.PackageName)
	}
	if contract.SourceFile != "SimpleToken.sol" {
		t.Errorf("expected source file 'SimpleToken.sol', got %q", contract.SourceFile)
	}
	if contract.SolcVersion != "0.8.20" {
		t.Errorf("expected solc version '0.8.20', got %q", contract.SolcVersion)
	}
	if contract.Bytecode != "0x608060405234801561001057600080fd5b506040516108013803806108018339818101604052810190610032919061018b565b" {
		t.Errorf("bytecode mismatch")
	}

	// Validate methods
	if len(contract.Methods) != 2 {
		t.Fatalf("expected 2 methods, got %d", len(contract.Methods))
	}

	transferMethod := findMethod(contract.Methods, "transfer")
	if transferMethod == nil {
		t.Fatal("transfer method not found")
	}
	if transferMethod.Signature != "transfer(address,uint256)" {
		t.Errorf("expected transfer signature 'transfer(address,uint256)', got %q", transferMethod.Signature)
	}
	if transferMethod.Selector != "0xa9059cbb" {
		t.Errorf("expected transfer selector '0xa9059cbb', got %q", transferMethod.Selector)
	}
	if len(transferMethod.Inputs) != 2 {
		t.Errorf("expected 2 transfer inputs, got %d", len(transferMethod.Inputs))
	}
	if len(transferMethod.Outputs) != 1 {
		t.Errorf("expected 1 transfer output, got %d", len(transferMethod.Outputs))
	}

	// Validate events
	if len(contract.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(contract.Events))
	}

	transferEvent := contract.Events[0]
	if transferEvent.Name != "Transfer" {
		t.Errorf("expected event name 'Transfer', got %q", transferEvent.Name)
	}
	if len(transferEvent.Inputs) != 3 {
		t.Errorf("expected 3 event inputs, got %d", len(transferEvent.Inputs))
	}

	// Check indexed fields
	if !transferEvent.Inputs[0].Indexed || !transferEvent.Inputs[1].Indexed {
		t.Error("first two event inputs should be indexed")
	}
	if transferEvent.Inputs[2].Indexed {
		t.Error("third event input should not be indexed")
	}

	// Validate errors
	if len(contract.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(contract.Errors))
	}

	insufficientBalanceError := contract.Errors[0]
	if insufficientBalanceError.Name != "InsufficientBalance" {
		t.Errorf("expected error name 'InsufficientBalance', got %q", insufficientBalanceError.Name)
	}
	if len(insufficientBalanceError.Inputs) != 3 {
		t.Errorf("expected 3 error inputs, got %d", len(insufficientBalanceError.Inputs))
	}

	// Validate constructor
	if contract.Constructor == nil {
		t.Fatal("constructor should not be nil")
	}
	if len(contract.Constructor.Inputs) != 3 {
		t.Errorf("expected 3 constructor inputs, got %d", len(contract.Constructor.Inputs))
	}

	// Validate struct generation
	if transferMethod.InputStruct == nil {
		t.Error("transfer method should have input struct")
	} else {
		if transferMethod.InputStruct.Name != "TransferInput" {
			t.Errorf("expected input struct name 'TransferInput', got %q", transferMethod.InputStruct.Name)
		}
		if len(transferMethod.InputStruct.Fields) != 2 {
			t.Errorf("expected 2 input struct fields, got %d", len(transferMethod.InputStruct.Fields))
		}
	}

	if transferEvent.Struct == nil {
		t.Error("transfer event should have struct")
	} else {
		if transferEvent.Struct.Name != "TransferEvent" {
			t.Errorf("expected event struct name 'TransferEvent', got %q", transferEvent.Struct.Name)
		}
	}
}

func findMethod(methods []types.Method, name string) *types.Method {
	for i := range methods {
		if methods[i].Name == name {
			return &methods[i]
		}
	}
	return nil
}
