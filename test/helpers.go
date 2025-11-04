// SPDX-License-Identifier: MIT

package test

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/otherview/solgen/internal/parse"
	"github.com/otherview/solgen/internal/types"
)

// processCombinedJSON parses combined JSON format and returns contracts
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

// testGeneratedCode verifies that generated code compiles without errors
func testGeneratedCode(t *testing.T, outputDir string) error {
	// Create a go.mod for the generated code
	goModContent := `module generated-test

go 1.21
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