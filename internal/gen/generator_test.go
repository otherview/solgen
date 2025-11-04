// SPDX-License-Identifier: MIT

package gen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/otherview/solgen/internal/types"
)

func TestGenerator_calculateImports(t *testing.T) {
	contract := &types.Contract{
		Name:        "TestContract",
		PackageName: "testcontract",
		Methods: []types.Method{
			{
				InputStruct: &types.Struct{
					Fields: []types.StructField{
						{Type: types.GoTypeBigInt},
						{Type: types.GoTypeAddress},
					},
				},
				OutputStruct: &types.Struct{
					Fields: []types.StructField{
						{Type: types.GoTypeBool},
					},
				},
			},
		},
		Events: []types.Event{
			{
				Struct: &types.Struct{
					Fields: []types.StructField{
						{Type: types.GoTypeBigInt},
					},
				},
			},
		},
	}

	generator := NewGenerator("/tmp")
	imports := generator.calculateImports(contract)

	// Check required imports are present
	requiredImports := []string{
		"math/big",
		"github.com/ethereum/go-ethereum/common",
		"fmt",
	}

	for _, required := range requiredImports {
		found := false
		for _, imp := range imports {
			if imp == required {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected import %q not found in %v", required, imports)
		}
	}
}

func TestGenerator_formatGoType(t *testing.T) {
	generator := NewGenerator("/tmp")

	tests := []struct {
		goType types.GoType
		want   string
	}{
		{types.GoTypeBool, "bool"},
		{types.GoTypeString, "string"},
		{types.GoTypeBigInt, "*big.Int"},
		{types.GoTypeAddress, "common.Address"},
		{types.GoType{TypeName: "[]uint256"}, "[]uint256"},
		{types.GoType{TypeName: "[32]byte"}, "[32]byte"},
	}

	for _, tt := range tests {
		t.Run(tt.goType.TypeName, func(t *testing.T) {
			got := generator.formatGoType(tt.goType)
			if got != tt.want {
				t.Errorf("formatGoType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerator_generateContractPackage(t *testing.T) {
	tempDir := t.TempDir()
	generator := NewGenerator(tempDir)

	// Create a simple contract for testing
	contract := &types.Contract{
		Name:             "SimpleTest",
		SourceFile:       "SimpleTest.sol",
		PackageName:      "simpletest",
		SolcVersion:      "0.8.20",
		ABIJson:          `[{"type":"function","name":"test","inputs":[],"outputs":[]}]`,
		Bytecode:         "0x1234",
		DeployedBytecode: "0x5678",
		Methods: []types.Method{
			{
				Name:      "test",
				Signature: "test()",
				Selector:  "0x12345678",
				// ABI field omitted as it's not used in the simplified template
				Inputs:    []types.Parameter{},
				Outputs:   []types.Parameter{},
			},
		},
		Events: []types.Event{},
		Errors: []types.ContractError{},
	}

	err := generator.generateContractPackage(contract)
	if err != nil {
		t.Fatalf("generateContractPackage failed: %v", err)
	}

	// Check that the package directory was created
	pkgDir := filepath.Join(tempDir, "simpletest")
	if _, err := os.Stat(pkgDir); os.IsNotExist(err) {
		t.Errorf("package directory %s was not created", pkgDir)
	}

	// Check that the Go file was created
	goFile := filepath.Join(pkgDir, "simpletest.go")
	if _, err := os.Stat(goFile); os.IsNotExist(err) {
		t.Errorf("Go file %s was not created", goFile)
	}

	// Read and validate the generated file
	content, err := os.ReadFile(goFile)
	if err != nil {
		t.Fatalf("failed to read generated file: %v", err)
	}

	contentStr := string(content)

	// Check for basic structure
	if !strings.Contains(contentStr, "package simpletest") {
		t.Error("generated file should contain package declaration")
	}
	if !strings.Contains(contentStr, "func ABI() string") {
		t.Error("generated file should contain ABI() function")
	}
	if !strings.Contains(contentStr, "func HexBytecode() string") {
		t.Error("generated file should contain HexBytecode() function")
	}
	if !strings.Contains(contentStr, "func Methods() MethodRegistry") {
		t.Error("generated file should contain Methods() function")
	}
	if !strings.Contains(contentStr, "TestMethod()") {
		t.Error("generated file should contain method function")
	}
}