// SPDX-License-Identifier: MIT

package gen

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/otherview/solgen/internal/types"
)

// Generator handles Go code generation from parsed contracts
type Generator struct {
	outputDir string
}

// NewGenerator creates a new code generator
func NewGenerator(outputDir string) *Generator {
	return &Generator{
		outputDir: outputDir,
	}
}

// Generate creates Go packages for all contracts
func (g *Generator) Generate(contracts []*types.Contract) error {
	// Ensure output directory exists
	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	// Generate package for each contract
	for _, contract := range contracts {
		if err := g.generateContractPackage(contract); err != nil {
			return fmt.Errorf("generating package for contract %s: %w", contract.Name, err)
		}
	}

	return nil
}

// generateContractPackage creates a single Go package for a contract
func (g *Generator) generateContractPackage(contract *types.Contract) error {
	// Create package directory
	pkgDir := filepath.Join(g.outputDir, contract.PackageName)
	if err := os.MkdirAll(pkgDir, 0755); err != nil {
		return fmt.Errorf("creating package directory: %w", err)
	}

	// Generate the main package file
	filePath := filepath.Join(pkgDir, contract.PackageName+".go")
	
	// Render template
	content, err := g.renderContract(contract)
	if err != nil {
		return fmt.Errorf("rendering contract template: %w", err)
	}

	// Format the generated Go code
	formatted, err := format.Source([]byte(content))
	if err != nil {
		// If formatting fails, write unformatted code for debugging
		fmt.Printf("Warning: failed to format generated code for %s: %v\n", contract.Name, err)
		formatted = []byte(content)
	}

	// Write to file
	if err := os.WriteFile(filePath, formatted, 0644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}

// renderContract renders the Go code for a contract using templates
func (g *Generator) renderContract(contract *types.Contract) (string, error) {
	tmpl, err := template.New("contract").Funcs(templateFuncs()).Parse(contractTemplate)
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	var buf strings.Builder
	data := &TemplateData{
		Contract: contract,
		Imports:  g.calculateImports(contract),
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return buf.String(), nil
}

// calculateImports determines which imports are needed beyond the base set.
// The base template already unconditionally imports:
//
//	encoding/hex, errors, fmt, math/big, strings
//
// This function only collects additional imports required by custom GoType
// packages that fall outside that base set.
func (g *Generator) calculateImports(contract *types.Contract) []string {
	// Imports already present in the hardcoded template header.
	baseImports := map[string]bool{
		"encoding/hex": true,
		"errors":       true,
		"fmt":          true,
		"math/big":     true,
		"strings":      true,
	}

	importSet := make(map[string]bool)
	checkGoType := func(goType types.GoType) {
		if goType.Import != "" && !baseImports[goType.Import] {
			importSet[goType.Import] = true
		}
	}

	for _, method := range contract.Methods {
		if method.InputStruct != nil {
			for _, field := range method.InputStruct.Fields {
				checkGoType(field.Type)
			}
		}
		if method.OutputStruct != nil {
			for _, field := range method.OutputStruct.Fields {
				checkGoType(field.Type)
			}
		}
	}
	for _, event := range contract.Events {
		if event.Struct != nil {
			for _, field := range event.Struct.Fields {
				checkGoType(field.Type)
			}
		}
	}
	for _, err := range contract.Errors {
		if err.Struct != nil {
			for _, field := range err.Struct.Fields {
				checkGoType(field.Type)
			}
		}
	}
	if contract.Constructor != nil && contract.Constructor.InputStruct != nil {
		for _, field := range contract.Constructor.InputStruct.Fields {
			checkGoType(field.Type)
		}
	}

	var imports []string
	for imp := range importSet {
		imports = append(imports, imp)
	}
	sort.Strings(imports)
	return imports
}