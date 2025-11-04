// SPDX-License-Identifier: MIT

package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/otherview/solgen/internal/gen"
)

func TestIntegration_SimpleToken(t *testing.T) {
	// Test the new pipeline architecture using Docker containers
	if !isDockerAvailable(t) {
		t.Skip("Docker is not available")
	}

	// Prepare test/out/integration directory (relative to project root)
	outputDir := "../test/out/integration/simpletoken"
	if err := os.RemoveAll("../test/out/integration"); err != nil {
		t.Fatalf("failed to clean integration output directory: %v", err)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("failed to create integration output directory: %v", err)
	}

	// Step 1: Use solc docker container to generate combined JSON
	combinedJSON, err := runSolcDocker("data/contracts/SimpleToken.sol")
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
		"var Bytecode = HexData(",
		"var DeployedBytecode = HexData(",
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

	// Prepare test/out/integration directory (relative to project root)
	outputDir := "../test/out/integration/cli"
	if err := os.RemoveAll("../test/out/integration/cli"); err != nil {
		t.Fatalf("failed to clean CLI output directory: %v", err)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("failed to create CLI output directory: %v", err)
	}

	// Test CLI using pipeline approach
	// Build the binary first
	binaryPath := filepath.Join("../test/out/integration", "solgen")
	absBinaryPath, _ := filepath.Abs(binaryPath)
	
	// Ensure the directory exists for the binary
	if err := os.MkdirAll(filepath.Dir(absBinaryPath), 0755); err != nil {
		t.Fatalf("failed to create binary directory: %v", err)
	}
	
	buildCmd := exec.Command("go", "build", "-o", absBinaryPath, "./cmd/solgen")
	// Set working directory to project root (one level up from test/)
	projectRoot, _ := filepath.Abs("..")
	buildCmd.Dir = projectRoot
	
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build solgen binary: %v\nOutput: %s", err, string(output))
	}

	// Step 1: Generate combined JSON using solc
	combinedJSON, err := runSolcDocker("data/contracts/SimpleToken.sol")
	if err != nil {
		t.Fatalf("solc compilation failed: %v", err)
	}

	// Step 2: Test solgen CLI with the combined JSON  
	absOutputDir, _ := filepath.Abs(outputDir)
	cmd := exec.Command(absBinaryPath, "--out", absOutputDir)
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

