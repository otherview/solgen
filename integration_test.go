// SPDX-License-Identifier: MIT

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/otherview/solgen/internal/gen"
	"github.com/otherview/solgen/internal/parse"
	"github.com/otherview/solgen/internal/types"
)


func TestIntegration_SimpleToken(t *testing.T) {
	// Test the new pipeline architecture using Docker containers
	if !isDockerAvailable(t) {
		t.Skip("Docker is not available")
	}

	// Create temporary directories
	tempDir := t.TempDir()
	outputDir := filepath.Join(tempDir, "generated")

	// Step 1: Use solc docker container to generate combined JSON
	combinedJSON, err := runSolcDocker("testdata/contracts/SimpleToken.sol")
	if err != nil {
		t.Fatalf("solc compilation failed: %v", err)
	}

	// Step 2: Process the combined JSON using internal functions
	contracts, err := processCombinedJSON(combinedJSON)
	if err != nil {
		t.Fatalf("processing combined JSON failed: %v", err)
	}

	if len(contracts) != 1 {
		t.Fatalf("expected 1 contract, got %d", len(contracts))
	}

	contract := contracts[0]

	// Validate parsed contract
	if contract.Name != "SimpleToken" {
		t.Errorf("expected contract name 'SimpleToken', got %q", contract.Name)
	}

	if contract.PackageName != "simpletoken" {
		t.Errorf("expected package name 'simpletoken', got %q", contract.PackageName)
	}

	// Check that we have expected methods
	expectedMethods := []string{"transfer", "approve", "transferFrom", "mint", "getBalance", "multiTransfer"}
	actualMethods := make(map[string]bool)
	for _, method := range contract.Methods {
		actualMethods[method.Name] = true
	}

	for _, expected := range expectedMethods {
		if !actualMethods[expected] {
			t.Errorf("expected method %q not found", expected)
		}
	}

	// Check that we have expected events
	expectedEvents := []string{"Transfer", "Approval"}
	actualEvents := make(map[string]bool)
	for _, event := range contract.Events {
		actualEvents[event.Name] = true
	}

	for _, expected := range expectedEvents {
		if !actualEvents[expected] {
			t.Errorf("expected event %q not found", expected)
		}
	}

	// Check that we have expected errors
	expectedErrors := []string{"InsufficientBalance", "InsufficientAllowance"}
	actualErrors := make(map[string]bool)
	for _, err := range contract.Errors {
		actualErrors[err.Name] = true
	}

	for _, expected := range expectedErrors {
		if !actualErrors[expected] {
			t.Errorf("expected error %q not found", expected)
		}
	}

	// Generate Go code
	generator := gen.NewGenerator(outputDir)
	if err := generator.Generate(contracts); err != nil {
		t.Fatalf("code generation failed: %v", err)
	}

	// Validate generated files
	generatedFile := filepath.Join(outputDir, "simpletoken", "simpletoken.go")
	if _, err := os.Stat(generatedFile); os.IsNotExist(err) {
		t.Fatalf("generated file %s does not exist", generatedFile)
	}

	// Read and validate generated content
	content, err := os.ReadFile(generatedFile)
	if err != nil {
		t.Fatalf("failed to read generated file: %v", err)
	}

	contentStr := string(content)

	// Validate basic structure
	expectedContents := []string{
		"package simpletoken",
		"func ABI() string",
		"func HexBytecode() string",
		"func HexDeployedBytecode() string",
		"func Methods() MethodRegistry",
		"func Events() EventRegistry",
		"func Errors() ErrorRegistry",
		"type TransferEvent struct",
		"type InsufficientBalanceError struct",
		"TransferMethod()",
		"GetTransferEvent()",
		"GetInsufficientBalanceError()",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("generated file should contain %q", expected)
		}
	}

	// Validate that the generated Go code compiles
	if err := testGeneratedCode(t, outputDir); err != nil {
		t.Errorf("generated code compilation failed: %v", err)
	}
}

func TestIntegration_CLI(t *testing.T) {
	if !isDockerAvailable(t) {
		t.Skip("Docker is not available")
	}

	tempDir := t.TempDir()
	outputDir := filepath.Join(tempDir, "generated")

	// Test CLI using pipeline approach
	// Build the binary first
	binaryPath := filepath.Join(tempDir, "solgen")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/solgen")

	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build solgen binary: %v", err)
	}

	// Step 1: Generate combined JSON using solc
	combinedJSON, err := runSolcDocker("testdata/contracts/SimpleToken.sol")
	if err != nil {
		t.Fatalf("solc compilation failed: %v", err)
	}

	// Step 2: Test solgen CLI with the combined JSON
	cmd := exec.Command(binaryPath, "--out", outputDir)
	cmd.Stdin = strings.NewReader(string(combinedJSON))

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("solgen command failed: %v\nOutput: %s", err, string(output))
	}

	// Verify output directory was created
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		t.Errorf("output directory %s was not created", outputDir)
	}

	// Verify generated files exist
	generatedFile := filepath.Join(outputDir, "simpletoken", "simpletoken.go")
	if _, err := os.Stat(generatedFile); os.IsNotExist(err) {
		t.Errorf("generated file %s does not exist", generatedFile)
	}
}

// runSolcDocker runs solc in a Docker container and returns the combined JSON output
func runSolcDocker(contractPath string) ([]byte, error) {
	// Get absolute path for mounting
	abs, err := filepath.Abs(contractPath)
	if err != nil {
		return nil, err
	}
	dir := filepath.Dir(abs)
	file := filepath.Base(abs)

	cmd := exec.Command("docker", "run", "--rm", "-v", dir+":/sources",
		"ghcr.io/argotorg/solc:0.8.20",
		"--combined-json", "abi,bin,bin-runtime,hashes",
		"--optimize", "--optimize-runs", "200",
		"/sources/"+file)

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return output, nil
}

func isDockerAvailable(t *testing.T) bool {
	cmd := exec.Command("docker", "version")
	err := cmd.Run()
	return err == nil
}

func testGeneratedCode(t *testing.T, outputDir string) error {
	// Create a go.mod for the generated code
	goModContent := `module generated-test

go 1.21

require (
	github.com/ethereum/go-ethereum v1.13.5
)
`
	goModPath := filepath.Join(outputDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte(goModContent), 0644); err != nil {
		return err
	}

	// Run go mod tidy to resolve dependencies
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = outputDir
	if output, err := tidyCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go mod tidy failed: %v\nOutput: %s", err, string(output))
	}

	// Try to build the generated code
	buildCmd := exec.Command("go", "build", "./...")
	buildCmd.Dir = outputDir
	output, err := buildCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go build failed: %v\nOutput: %s", err, string(output))
	}
	return nil
}

// processCombinedJSON converts combined JSON to contracts (simplified version of main.go logic)
func processCombinedJSON(data []byte) ([]*types.Contract, error) {
	var combined types.CombinedJSON
	if err := json.Unmarshal(data, &combined); err != nil {
		return nil, fmt.Errorf("parsing combined JSON: %w", err)
	}

	// Convert to standard JSON format
	result, err := convertCombinedToStandard(combined)
	if err != nil {
		return nil, fmt.Errorf("converting to standard format: %w", err)
	}

	// Parse using existing parser
	contracts, err := parse.ResultWithVersion(result, "0.8.20")
	if err != nil {
		return nil, fmt.Errorf("parsing contracts: %w", err)
	}

	return contracts, nil
}

// convertCombinedToStandard converts combined JSON to standard CompileResult format
func convertCombinedToStandard(combined types.CombinedJSON) (*types.CompileResult, error) {
	result := &types.CompileResult{
		Contracts: make(map[string]map[string]types.ContractResult),
	}

	for key, contract := range combined.Contracts {
		// Parse the key format "filename.sol:ContractName"
		parts := strings.Split(key, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid contract key format: %s", key)
		}
		filename := parts[0]
		contractName := parts[1]

		if result.Contracts[filename] == nil {
			result.Contracts[filename] = make(map[string]types.ContractResult)
		}

		contractResult := types.ContractResult{
			ABI: contract.ABI,
			EVM: types.EVMResult{
				Bytecode: types.BytecodeResult{
					Object: contract.Bin,
				},
				DeployedBytecode: types.BytecodeResult{
					Object: contract.BinRuntime,
				},
			},
		}

		// Add method identifiers if available
		if contract.Hashes != nil {
			contractResult.EVM.MethodIdentifiers = contract.Hashes
		}

		result.Contracts[filename][contractName] = contractResult
	}

	return result, nil
}
