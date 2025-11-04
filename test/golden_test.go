// SPDX-License-Identifier: MIT

package test

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/otherview/solgen/internal/gen"
)

// updateGolden flag can be set to update golden files
var updateGolden = flag.Bool("update-golden", false, "update golden files")

func TestGolden_SimpleContract(t *testing.T) {
	// Simple contract for golden file testing
	input := `{
		"contracts": {
			"SimpleContract.sol:SimpleContract": {
				"abi": [
					{
						"type": "constructor",
						"inputs": [{"name": "initialValue", "type": "uint256"}]
					},
					{
						"type": "function", 
						"name": "getValue",
						"inputs": [],
						"outputs": [{"name": "", "type": "uint256"}],
						"stateMutability": "view"
					},
					{
						"type": "function",
						"name": "setValue", 
						"inputs": [{"name": "newValue", "type": "uint256"}],
						"outputs": [],
						"stateMutability": "nonpayable"
					},
					{
						"type": "event",
						"name": "ValueChanged",
						"inputs": [
							{"name": "oldValue", "type": "uint256", "indexed": false},
							{"name": "newValue", "type": "uint256", "indexed": false}
						]
					},
					{
						"type": "error",
						"name": "InvalidValue", 
						"inputs": [{"name": "provided", "type": "uint256"}]
					}
				],
				"bin": "0x608060405234801561001057600080fd5b5060405161012c38038061012c833981810160405281019061003291906100a4565b80600081905550506100d1565b600080fd5b6000819050919050565b61005a81610047565b811461006557600080fd5b50565b60008151905061007781610051565b92915050565b6000602082840312156100935761009261004257600080fd5b5b60006100a184828501610068565b91505092915050565b604c806100e06000396000f3fe608060405200",
				"bin-runtime": "0x6080604052348015600f57600080fd5b506004361060325760003560e01c806320965255146037578063552410771460005b600080fd5b60005460405190815260200160405180910390f35b6000819055565b600080fd5b6000819050919050565b605c81604f565b8114606657600080fd5b50565b600081359050607a81605556565b92915050565b600060208284031215609357609260004a565b5b6000609f84828501606d565b9150509291505056fea2646970667358221220",
				"metadata": "{\"compiler\":{\"version\":\"0.8.20+commit.a1b79de6\"},\"language\":\"Solidity\",\"output\":{\"abi\":[],\"devdoc\":{\"kind\":\"dev\",\"methods\":{},\"version\":1},\"userdoc\":{\"kind\":\"user\",\"methods\":{},\"version\":1}},\"settings\":{\"compilationTarget\":{\"SimpleContract.sol\":\"SimpleContract\"},\"evmVersion\":\"london\",\"libraries\":{},\"metadata\":{\"bytecodeHash\":\"ipfs\"},\"optimizer\":{\"enabled\":true,\"runs\":200},\"remappings\":[]},\"sources\":{\"SimpleContract.sol\":{\"keccak256\":\"0x\",\"urls\":[\"bzz-raw://\",\"dweb:/ipfs/\"]}},\"version\":1}",
				"hashes": {
					"getValue()": "20965255",
					"setValue(uint256)": "55241077"
				}
			}
		}
	}`

	testGoldenFile(t, "simple_contract", input)
}

func TestGolden_ComplexContract(t *testing.T) {
	// More complex contract with multiple types
	input := `{
		"contracts": {
			"ComplexContract.sol:ComplexContract": {
			"abi": [
				{
					"type": "function",
					"name": "complexFunction",
					"inputs": [
						{"name": "addresses", "type": "address[]"},
						{"name": "amounts", "type": "uint256[]"},
						{"name": "data", "type": "bytes"},
						{"name": "flag", "type": "bool"}
					],
					"outputs": [
						{"name": "success", "type": "bool"},
						{"name": "results", "type": "uint256[]"}
					],
					"stateMutability": "nonpayable"
				},
				{
					"type": "function",
					"name": "getMapping",
					"inputs": [{"name": "key", "type": "bytes32"}],
					"outputs": [{"name": "value", "type": "string"}],
					"stateMutability": "view"
				},
				{
					"type": "event",
					"name": "ComplexEvent", 
					"inputs": [
						{"name": "user", "type": "address", "indexed": true},
						{"name": "data", "type": "bytes", "indexed": false},
						{"name": "timestamp", "type": "uint256", "indexed": true}
					]
				},
				{
					"type": "error",
					"name": "ComplexError",
					"inputs": [
						{"name": "reason", "type": "string"},
						{"name": "code", "type": "uint256"}
					]
				}
			],
			"bin": "0x608060405234801561001057600080fd5b50610abc806100206000396000f3fe",
			"bin-runtime": "0x6080604052348015600f57600080fd5b50600436106100365760003560e01c8063abcd123414603a5780634567890114603f565b5b600080fd5b005b005b600080fd5b6000819050919050565b60558160048565b8114605f57600080fd5b50565b6000813590506070816050565b92915050565b6000602082840312156088576087600b565b5b600060948482850160635b915050929150505056fea264697066735822",
			"metadata": "{\"compiler\":{\"version\":\"0.8.20\"}}",
			"hashes": {
				"complexFunction(address[],uint256[],bytes,bool)": "abcd1234", 
				"getMapping(bytes32)": "45678901"
			}
		}
	}
}`

	testGoldenFile(t, "complex_contract", input)
}

func TestGolden_MultipleContracts(t *testing.T) {
	// Test with multiple contracts in same file
	input := `{
		"contracts": {
			"MultiContract.sol:ContractA": {
			"abi": [
				{
					"type": "function",
					"name": "functionA", 
					"inputs": [],
					"outputs": [{"name": "", "type": "uint256"}],
					"stateMutability": "pure"
				}
			],
			"bin": "0x608060405234801561001057600080fd5b50610123",
			"bin-runtime": "0x608060405234801561001057600080fd5b50610456",
			"metadata": "{}",
			"hashes": {"functionA()": "aaaaaaaa"}
		},
		"MultiContract.sol:ContractB": {
			"abi": [
				{
					"type": "function",
					"name": "functionB",
					"inputs": [{"name": "param", "type": "string"}],
					"outputs": [{"name": "", "type": "bytes32"}],
					"stateMutability": "pure"  
				}
			],
			"bin": "0x608060405234801561001057600080fd5b50610789",
			"bin-runtime": "0x608060405234801561001057600080fd5b50610abc", 
			"metadata": "{}",
			"hashes": {"functionB(string)": "bbbbbbbb"}
		}
	}
}`

	testGoldenFile(t, "multi_contract", input)
}

// testGoldenFile is a helper that processes input and compares with golden file
func testGoldenFile(t *testing.T, testName, input string) {
	// Process the combined JSON to get contracts
	contracts, err := processCombinedJSON([]byte(input))
	if err != nil {
		t.Fatalf("processCombinedJSON failed: %v", err)
	}

	// Prepare test/out/golden directory (relative to project root)
	outputDir := filepath.Join("..", "test", "out", "golden", testName)
	if err := os.RemoveAll(outputDir); err != nil {
		t.Fatalf("failed to clean output directory: %v", err)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("failed to create output directory: %v", err)
	}

	// Generate Go code for each contract
	generator := gen.NewGenerator(outputDir)
	
	if err := generator.Generate(contracts); err != nil {
		t.Fatalf("code generation failed: %v", err)
	}

	// For each generated contract, compare with golden file
	for _, contract := range contracts {
		packageDir := filepath.Join(outputDir, contract.PackageName)
		generatedFile := filepath.Join(packageDir, contract.PackageName+".go")

		// Read generated content
		generatedContent, err := os.ReadFile(generatedFile)
		if err != nil {
			t.Fatalf("failed to read generated file %s: %v", generatedFile, err)
		}

		// Normalize line endings and whitespace
		generated := normalizeContent(string(generatedContent))

		// Golden file path (relative to project root)
		goldenFile := filepath.Join("..", "test", "data", "golden", testName+"_"+contract.PackageName, contract.PackageName+".go")

		if *updateGolden {
			// Create golden directory
			goldenDir := filepath.Dir(goldenFile)
			if err := os.MkdirAll(goldenDir, 0755); err != nil {
				t.Fatalf("failed to create golden directory %s: %v", goldenDir, err)
			}
			// Update golden file
			if err := os.WriteFile(goldenFile, []byte(generated), 0644); err != nil {
				t.Fatalf("failed to update golden file %s: %v", goldenFile, err)
			}
			t.Logf("Updated golden file: %s", goldenFile)
			continue // Continue to next contract instead of returning
		}

		// Read golden file
		goldenContent, err := os.ReadFile(goldenFile) 
		if err != nil {
			t.Fatalf("failed to read golden file %s: %v (run with -update-golden to create)", goldenFile, err)
		}

		golden := normalizeContent(string(goldenContent))

		// Compare
		if generated != golden {
			t.Errorf("Generated content for %s does not match golden file %s", contract.Name, goldenFile)
			t.Logf("Generated:\n%s", generated)
			t.Logf("Golden:\n%s", golden)
			t.Logf("Run with -update-golden to update the golden file")
		}

		// Test that the golden file compiles
		if err := testGoldenCompiles(t, goldenFile); err != nil {
			t.Errorf("Golden file %s does not compile: %v", goldenFile, err)
		}
	}

	// If we were updating golden files, we're done
	if *updateGolden {
		return
	}
}

// normalizeContent normalizes content for comparison
func normalizeContent(content string) string {
	// Normalize line endings
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")
	
	// Trim trailing whitespace from each line
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
	}
	
	// Remove trailing empty lines
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	
	return strings.Join(lines, "\n") + "\n"
}

// testGoldenCompiles tests that a golden file compiles
func testGoldenCompiles(t *testing.T, goldenFile string) error {
	// Create a temporary directory for testing compilation
	tempDir := t.TempDir()
	goldenDir := filepath.Dir(goldenFile)
	packageName := filepath.Base(goldenDir)

	// Create package directory
	testPkgDir := filepath.Join(tempDir, packageName)
	if err := os.MkdirAll(testPkgDir, 0755); err != nil {
		return err
	}

	// Copy golden file to test directory
	testFile := filepath.Join(testPkgDir, packageName+".go")
	content, err := os.ReadFile(goldenFile)
	if err != nil {
		return err
	}
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		return err
	}

	// Test compilation
	return testGeneratedCode(t, tempDir)
}

// Test that verifies generated code compiles
func TestGolden_CompileGenerated(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping compilation test in short mode")
	}

	// Use simple contract for compilation test
	input := `{
		"contracts": {
			"TestCompile.sol:TestCompile": {
			"abi": [
				{
					"type": "function",
					"name": "test",
					"inputs": [],
					"outputs": [{"name": "", "type": "uint256"}],
					"stateMutability": "pure"
				}
			],
			"bin": "0x608060405234801561001057600080fd5b50",
			"bin-runtime": "0x6080604052348015600f57600080fd5b50",
			"metadata": "{}",
			"hashes": {"test()": "12345678"}
		}
	}
}`

	contracts, err := processCombinedJSON([]byte(input))
	if err != nil {
		t.Fatalf("processCombinedJSON failed: %v", err)
	}

	// Prepare test/out/compile directory (relative to project root)
	outputDir := "../test/out/compile"
	if err := os.RemoveAll(outputDir); err != nil {
		t.Fatalf("failed to clean output directory: %v", err)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("failed to create output directory: %v", err)
	}

	generator := gen.NewGenerator(outputDir)
	
	if err := generator.Generate(contracts); err != nil {
		t.Fatalf("code generation failed: %v", err)
	}

	// Test that generated code compiles
	if err := testGeneratedCode(t, outputDir); err != nil {
		t.Errorf("generated code compilation failed: %v", err)
	}
}