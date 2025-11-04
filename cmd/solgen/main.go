// SPDX-License-Identifier: MIT

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/otherview/solgen/internal/gen"
	"github.com/otherview/solgen/internal/parse"
	"github.com/otherview/solgen/internal/types"
	"github.com/spf13/cobra"
)

type ProcessFlags struct {
	Output  string
	Verbose bool
}


func main() {
	if err := rootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func rootCmd() *cobra.Command {
	flags := &ProcessFlags{}

	cmd := &cobra.Command{
		Use:     "solgen",
		Short:   "Solidity to Go code generator",
		Long:    "A code generator that reads solc combined JSON output and generates Go packages.",
		Version: "0.1.0",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runProcessJSON(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Output, "out", "", "Output directory for generated Go packages")
	cmd.Flags().BoolVarP(&flags.Verbose, "verbose", "v", false, "Verbose output")

	cmd.MarkFlagRequired("out")

	return cmd
}

func runProcessJSON(flags *ProcessFlags) error {
	// Validate output directory
	if flags.Output == "" {
		return fmt.Errorf("output directory cannot be empty")
	}
	if err := os.MkdirAll(flags.Output, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Read combined JSON from stdin
	jsonData, err := io.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("reading from stdin: %w", err)
	}

	if len(jsonData) == 0 {
		return fmt.Errorf("no JSON data provided on stdin")
	}

	// Parse combined JSON
	var combinedJSON types.CombinedJSON
	if err := json.Unmarshal(jsonData, &combinedJSON); err != nil {
		return fmt.Errorf("parsing combined JSON: %w", err)
	}

	if len(combinedJSON.Contracts) == 0 {
		return fmt.Errorf("no contracts found in JSON output")
	}

	// Convert combined JSON to standard format
	standardResult, err := convertCombinedToStandard(combinedJSON, flags.Verbose)
	if err != nil {
		return fmt.Errorf("converting JSON format: %w", err)
	}

	// Extract solc version, fallback to unknown if not available
	solcVersion := combinedJSON.Version
	if solcVersion == "" {
		solcVersion = "unknown"
	}

	// Parse compilation result (reuse existing logic)
	contracts, err := parse.ResultWithVersion(standardResult, solcVersion)
	if err != nil {
		return fmt.Errorf("parsing failed: %w", err)
	}

	// Generate Go packages (reuse existing logic)
	generator := gen.NewGenerator(flags.Output)
	if err := generator.Generate(contracts); err != nil {
		return fmt.Errorf("code generation failed: %w", err)
	}

	fmt.Printf("Successfully generated %d contract packages in %s\n", len(contracts), flags.Output)
	return nil
}

// convertCombinedToStandard converts combined JSON format to standard JSON format.
// This conversion layer provides compatibility with the existing parser infrastructure
// and allows for potential future support of solc's --standard-json format.
func convertCombinedToStandard(combinedJSON types.CombinedJSON, verbose bool) (*types.CompileResult, error) {
	result := &types.CompileResult{
		Contracts: make(map[string]map[string]types.ContractResult),
	}

	for contractKey, contract := range combinedJSON.Contracts {
		// Parse contract key format: "filename.sol:ContractName"
		parts := strings.SplitN(contractKey, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid contract key format: %s (expected 'file.sol:ContractName')", contractKey)
		}

		filename := parts[0]
		contractName := parts[1]

		if verbose {
			fmt.Printf("Processing contract: %s in file: %s\n", contractName, filename)
		}

		// Create standard format structure
		if result.Contracts[filename] == nil {
			result.Contracts[filename] = make(map[string]types.ContractResult)
		}

		result.Contracts[filename][contractName] = types.ContractResult{
			ABI: contract.ABI,
			EVM: types.EVMResult{
				Bytecode: types.BytecodeResult{
					Object: contract.Bin,
				},
				DeployedBytecode: types.BytecodeResult{
					Object: contract.BinRuntime,
				},
				MethodIdentifiers: contract.Hashes,
			},
		}
	}

	return result, nil
}
