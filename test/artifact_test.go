// SPDX-License-Identifier: MIT

package test

import (
	"os"
	"testing"

	"github.com/otherview/solgen/internal/gen"
)

func TestArtifact_OutputManagement(t *testing.T) {
	// Simple test to verify artifact output management
	simpleJSON := `{
		"contracts": {
			"Test.sol:Test": {
				"abi": [{"type": "function", "name": "test", "inputs": [], "outputs": []}],
				"bin": "0x1234",
				"bin-runtime": "0x5678",
				"hashes": {"test()": "12345678"}
			}
		}
	}`

	contracts, err := processCombinedJSON([]byte(simpleJSON))
	if err != nil {
		t.Fatalf("failed to process JSON: %v", err)
	}

	// Test artifact output to test/out/decode (relative to project root)
	outputDir := "../test/out/decode"
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

	// Check that files were created
	expectedFile := outputDir + "/test/test.go"
	if _, err := os.Stat(expectedFile); err != nil {
		t.Fatalf("expected file %s was not created: %v", expectedFile, err)
	}

	// Read and verify file exists
	content, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("failed to read generated file: %v", err)
	}

	if len(content) == 0 {
		t.Error("generated file is empty")
	}

	t.Logf("✅ Artifact successfully created at: %s", expectedFile)
	t.Logf("📁 File size: %d bytes", len(content))
	t.Logf("📂 Artifacts will remain in test/out/ for inspection")
}