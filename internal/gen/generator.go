// SPDX-License-Identifier: MIT

package gen

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"sort"
	"strconv"
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
	tmpl, err := template.New("contract").Funcs(g.templateFuncs()).Parse(contractTemplate)
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

// calculateImports determines which imports are needed for the contract
func (g *Generator) calculateImports(contract *types.Contract) []string {
	importSet := make(map[string]bool)
	
	// Always needed imports for the simplified template
	importSet["fmt"] = true
	importSet["github.com/ethereum/go-ethereum/common"] = true

	// Check if we need math/big - only include if it appears in struct fields
	needsBigInt := false
	checkGoType := func(goType types.GoType) {
		if goType.Import != "" && goType.Import != "math/big" {
			importSet[goType.Import] = true
		}
		if goType.Import == "math/big" {
			needsBigInt = true
		}
	}

	// Check method structs only (not individual method parameters)
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

	// Check event structs
	for _, event := range contract.Events {
		if event.Struct != nil {
			for _, field := range event.Struct.Fields {
				checkGoType(field.Type)
			}
		}
	}

	// Check error structs
	for _, err := range contract.Errors {
		if err.Struct != nil {
			for _, field := range err.Struct.Fields {
				checkGoType(field.Type)
			}
		}
	}

	// Check constructor struct
	if contract.Constructor != nil && contract.Constructor.InputStruct != nil {
		for _, field := range contract.Constructor.InputStruct.Fields {
			checkGoType(field.Type)
		}
	}

	if needsBigInt {
		importSet["math/big"] = true
	}

	// Convert to sorted slice
	var imports []string
	for imp := range importSet {
		imports = append(imports, imp)
	}
	
	sort.Strings(imports)
	return imports
}

// templateFuncs returns template helper functions
func (g *Generator) templateFuncs() template.FuncMap {
	return template.FuncMap{
		"formatGoType": g.formatGoType,
		"quote":        strconv.Quote,
		"lower":        strings.ToLower,
		"title":        g.titleCase,
		"join":         strings.Join,
		"add":          func(a, b int) int { return a + b },
		"default":      func(def, val string) string { if val == "" { return def }; return val },
	}
}

// formatGoType formats a GoType for use in generated code
func (g *Generator) formatGoType(goType types.GoType) string {
	return goType.TypeName
}

// titleCase provides a simple title case conversion
func (g *Generator) titleCase(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// TemplateData holds data for template rendering
type TemplateData struct {
	Contract *types.Contract
	Imports  []string
}