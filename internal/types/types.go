// SPDX-License-Identifier: MIT

package types

import (
	"encoding/hex"
	"encoding/json"
	"strings"
)

// Address represents a 20-byte Ethereum address
type Address [20]byte

// String returns the hex string representation of the address
func (a Address) String() string {
	return "0x" + hex.EncodeToString(a[:])
}

// Hash represents a 32-byte hash
type Hash [32]byte

// String returns the hex string representation of the hash
func (h Hash) String() string {
	return "0x" + hex.EncodeToString(h[:])
}

// Bytes returns the hash as a byte slice
func (h Hash) Bytes() []byte {
	return h[:]
}

// HexData represents hex-encoded byte data with convenient access methods
type HexData string

// Hex returns the hex string representation
func (h HexData) Hex() string {
	return string(h)
}

// Bytes returns the decoded bytes from the hex string
// Panics if the hex string is invalid (should not happen in generated code)
func (h HexData) Bytes() []byte {
	hexStr := string(h)
	if hexStr == "" {
		return nil
	}
	
	// Remove 0x prefix if present
	if strings.HasPrefix(hexStr, "0x") {
		hexStr = hexStr[2:]
	}
	
	if hexStr == "" {
		return nil
	}
	
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		panic("invalid hex string in generated code: " + string(h))
	}
	return data
}

// CompileConfig holds configuration for Solidity compilation
type CompileConfig struct {
	Inputs       []string // Input files, dirs, or globs
	Output       string   // Output directory
	Includes     []string // Include directories
	Optimize     bool     // Enable optimizer
	OptimizeRuns int      // Optimizer runs
	EVMVersion   string   // EVM version
	ViaIR        bool     // Via IR compilation
}

// CompileResult represents the standard JSON output from solc
type CompileResult struct {
	Contracts map[string]map[string]ContractResult `json:"contracts"`
	Errors    []CompileError                       `json:"errors,omitempty"`
	Sources   map[string]SourceResult              `json:"sources"`
}

// CompileError represents a compilation error or warning
type CompileError struct {
	Component        string `json:"component"`
	ErrorCode        string `json:"errorCode,omitempty"`
	FormattedMessage string `json:"formattedMessage"`
	Message          string `json:"message"`
	Severity         string `json:"severity"`
	Type             string `json:"type"`
}

// ContractResult holds solc output for a single contract
type ContractResult struct {
	ABI json.RawMessage `json:"abi"`
	EVM EVMResult       `json:"evm"`
}

// EVMResult holds EVM-related compilation output
type EVMResult struct {
	Bytecode         BytecodeResult            `json:"bytecode"`
	DeployedBytecode BytecodeResult            `json:"deployedBytecode"`
	MethodIdentifiers map[string]string        `json:"methodIdentifiers"`
}

// BytecodeResult holds bytecode and link references
type BytecodeResult struct {
	Object         string                    `json:"object"`
	LinkReferences map[string]map[string][]LinkRef `json:"linkReferences"`
}

// LinkRef represents a library link reference
type LinkRef struct {
	Start  int `json:"start"`
	Length int `json:"length"`
}

// SourceResult holds source-level compilation info
type SourceResult struct {
	ID  int    `json:"id"`
	AST interface{} `json:"ast,omitempty"`
}

// Contract represents a parsed contract ready for code generation
type Contract struct {
	Name             string
	SourceFile       string
	PackageName      string
	SolcVersion      string
	ABIJson          string
	Bytecode         HexData
	DeployedBytecode HexData
	Methods          []Method
	Events           []Event
	Errors           []ContractError
	Constructor      *Constructor
	Structs          []Struct // Standalone struct definitions
}

// Method represents a contract method
type Method struct {
	Name         string
	Signature    string
	Selector     HexData
	Inputs       []Parameter
	Outputs      []Parameter
	InputStruct  *Struct
	OutputStruct *Struct
}

// Event represents a contract event
type Event struct {
	Name    string
	Topic   Hash
	Inputs  []Parameter
	Struct  *Struct
}

// ContractError represents a custom contract error
type ContractError struct {
	Name      string
	Signature string
	Selector  HexData
	Inputs    []Parameter
	Struct    *Struct
}

// Constructor represents a contract constructor
type Constructor struct {
	Signature      string
	Inputs         []Parameter
	InputStruct    *Struct
	LinkReferences map[string][]LinkRef
}

// Parameter represents a method/event/error parameter
type Parameter struct {
	Name    string
	Type    GoType
	Indexed bool // for events
}

// Struct represents a generated Go struct
type Struct struct {
	Name   string
	Fields []StructField
}

// StructField represents a field in a generated struct
type StructField struct {
	Name    string
	Type    GoType
	JSONTag string
}

// GoType represents a Go type mapping
type GoType struct {
	Import   string // import path if needed
	TypeName string // Go type name
	IsSlice  bool   // for dynamic arrays
	IsPtr    bool   // for big.Int
}

// CombinedJSON represents the structure of solc --combined-json output
type CombinedJSON struct {
	Contracts map[string]CombinedContract `json:"contracts"`
	Version   string                      `json:"version,omitempty"`
}

// CombinedContract represents a single contract in combined JSON output
type CombinedContract struct {
	ABI        json.RawMessage   `json:"abi"`
	Bin        string            `json:"bin"`
	BinRuntime string            `json:"bin-runtime"`
	Hashes     map[string]string `json:"hashes,omitempty"`
	DevDoc     interface{}       `json:"devdoc,omitempty"`
	UserDoc    interface{}       `json:"userdoc,omitempty"`
}



// Common Go types
var (
	GoTypeBool         = GoType{TypeName: "bool"}
	GoTypeString       = GoType{TypeName: "string"}
	GoTypeBytes        = GoType{TypeName: "[]byte"}
	GoTypeBigInt       = GoType{Import: "math/big", TypeName: "*big.Int", IsPtr: true}
	GoTypeAddress      = GoType{TypeName: "Address"}
	GoTypeHash         = GoType{TypeName: "Hash"}
	GoTypeUint8        = GoType{TypeName: "uint8"}
	GoTypeUint16       = GoType{TypeName: "uint16"}
	GoTypeUint32       = GoType{TypeName: "uint32"}
	GoTypeUint64       = GoType{TypeName: "uint64"}
	GoTypeInt8         = GoType{TypeName: "int8"}
	GoTypeInt16        = GoType{TypeName: "int16"}
	GoTypeInt32        = GoType{TypeName: "int32"}
	GoTypeInt64        = GoType{TypeName: "int64"}
)