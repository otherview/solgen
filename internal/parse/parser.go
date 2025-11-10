// SPDX-License-Identifier: MIT

package parse

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/otherview/solgen/internal/types"
)

// structRegistry holds struct definitions collected during parsing
type structRegistry struct {
	structs map[string]types.Struct // key: struct name, value: struct definition
}

// newStructRegistry creates a new struct registry
func newStructRegistry() *structRegistry {
	return &structRegistry{
		structs: make(map[string]types.Struct),
	}
}

// registerStruct adds a struct definition to the registry
func (r *structRegistry) registerStruct(structName string, abiType abi.Type) {
	if structName == "" || structName == "AnonymousTuple" {
		return // Don't register anonymous tuples
	}
	
	// Don't re-register if already exists
	if _, exists := r.structs[structName]; exists {
		return
	}
	
	// Convert tuple elements to struct fields
	var fields []types.StructField
	for i, elemType := range abiType.TupleElems {
		goType, err := mapSolidityToGoTypeWithRegistry(*elemType, r)
		if err != nil {
			continue // Skip problematic fields for now
		}
		
		fieldName := "Field" + fmt.Sprintf("%d", i+1) // Default field name
		if i < len(abiType.TupleRawNames) && abiType.TupleRawNames[i] != "" {
			fieldName = exportIdentifier(abiType.TupleRawNames[i])
		}
		
		fields = append(fields, types.StructField{
			Name:    fieldName,
			Type:    goType,
			JSONTag: strings.ToLower(fieldName),
		})
	}
	
	r.structs[structName] = types.Struct{
		Name:   structName,
		Fields: fields,
	}
}

// getAllStructs returns all registered structs as a slice
func (r *structRegistry) getAllStructs() []types.Struct {
	var structs []types.Struct
	for _, s := range r.structs {
		structs = append(structs, s)
	}
	// Sort structs for deterministic output
	sort.Slice(structs, func(i, j int) bool {
		return structs[i].Name < structs[j].Name
	})
	return structs
}

// ResultWithVersion converts solc compilation result with version info
func ResultWithVersion(result *types.CompileResult, solcVersion string) ([]*types.Contract, error) {
	var contracts []*types.Contract
	nameCollisions := make(map[string][]string) // package name -> contract names

	// First pass: collect all contracts and check for package name collisions
	for sourceFile, sourceContracts := range result.Contracts {
		for contractName := range sourceContracts {
			pkgName := sanitizePackageName(contractName)
			nameCollisions[pkgName] = append(nameCollisions[pkgName], fmt.Sprintf("%s:%s", sourceFile, contractName))
		}
	}

	// Check for collisions
	for pkgName, contractNames := range nameCollisions {
		if len(contractNames) > 1 {
			return nil, fmt.Errorf("package name collision for %q: contracts %v would generate the same package name", pkgName, contractNames)
		}
	}

	// Second pass: parse contracts
	for sourceFile, sourceContracts := range result.Contracts {
		for contractName, contractResult := range sourceContracts {
			contract, err := parseContract(sourceFile, contractName, contractResult)
			if err != nil {
				return nil, fmt.Errorf("parsing contract %s:%s: %w", sourceFile, contractName, err)
			}
			contract.SolcVersion = solcVersion
			contracts = append(contracts, contract)
		}
	}

	// Sort contracts for deterministic output
	sort.Slice(contracts, func(i, j int) bool {
		if contracts[i].SourceFile != contracts[j].SourceFile {
			return contracts[i].SourceFile < contracts[j].SourceFile
		}
		return contracts[i].Name < contracts[j].Name
	})

	return contracts, nil
}

// parseContract parses a single contract from solc output
func parseContract(sourceFile, contractName string, result types.ContractResult) (*types.Contract, error) {
	// Parse ABI
	parsedABI, err := abi.JSON(strings.NewReader(string(result.ABI)))
	if err != nil {
		return nil, fmt.Errorf("parsing ABI: %w", err)
	}

	// Create struct registry to collect struct definitions
	registry := newStructRegistry()

	contract := &types.Contract{
		Name:             contractName,
		SourceFile:       sourceFile,
		PackageName:      sanitizePackageName(contractName),
		ABIJson:          string(result.ABI),
		Bytecode:         types.HexData(prefixHex(result.EVM.Bytecode.Object)),
		DeployedBytecode: types.HexData(prefixHex(result.EVM.DeployedBytecode.Object)),
	}

	// Parse methods
	methods, err := parseMethodsWithRegistry(parsedABI, result.EVM.MethodIdentifiers, registry)
	if err != nil {
		return nil, fmt.Errorf("parsing methods: %w", err)
	}
	contract.Methods = methods

	// Parse events
	events, err := parseEventsWithRegistry(parsedABI, registry)
	if err != nil {
		return nil, fmt.Errorf("parsing events: %w", err)
	}
	contract.Events = events

	// Parse errors
	errors, err := parseErrors(parsedABI)
	if err != nil {
		return nil, fmt.Errorf("parsing errors: %w", err)
	}
	contract.Errors = errors

	// Parse constructor
	constructor := parseConstructor(parsedABI, result.EVM.Bytecode.LinkReferences)
	contract.Constructor = constructor

	// Add all collected struct definitions
	contract.Structs = registry.getAllStructs()

	return contract, nil
}

// parseMethodsWithRegistry extracts and processes contract methods using struct registry
func parseMethodsWithRegistry(parsedABI abi.ABI, methodIds map[string]string, registry *structRegistry) ([]types.Method, error) {
	var methods []types.Method
	methodNames := make(map[string]int) // track name collisions

	// First pass: count method names for overload detection
	for _, method := range parsedABI.Methods {
		methodNames[method.Name]++
	}

	// Second pass: create method descriptors
	for _, method := range parsedABI.Methods {
		selector := methodIds[method.Sig]
		if selector == "" {
			return nil, fmt.Errorf("missing method identifier for %s", method.Sig)
		}

		// Generate method name with overload suffix if needed
		methodName := method.Name
		if methodNames[method.Name] > 1 {
			methodName = generateOverloadName(method.Name, method.Sig, selector)
		}

		// Parse inputs and outputs with registry
		inputs, err := parseParametersWithRegistry(method.Inputs, false, registry)
		if err != nil {
			return nil, fmt.Errorf("parsing inputs for method %s: %w", method.Sig, err)
		}

		outputs, err := parseParametersWithRegistry(method.Outputs, false, registry)
		if err != nil {
			return nil, fmt.Errorf("parsing outputs for method %s: %w", method.Sig, err)
		}

		// Create input/output structs if needed
		var inputStruct, outputStruct *types.Struct

		if len(inputs) > 1 {
			inputStruct = &types.Struct{
				Name:   exportIdentifier(methodName) + "Input",
				Fields: parametersToFields(inputs),
			}
		}

		if len(outputs) > 1 {
			outputStruct = &types.Struct{
				Name:   exportIdentifier(methodName) + "Output",
				Fields: parametersToFields(outputs),
			}
		}

		methods = append(methods, types.Method{
			Name:         methodName,
			Signature:    method.Sig,
			Selector:     types.HexData("0x" + selector),
			Inputs:       inputs,
			Outputs:      outputs,
			InputStruct:  inputStruct,
			OutputStruct: outputStruct,
		})
	}

	// Sort methods for deterministic output
	sort.Slice(methods, func(i, j int) bool {
		if methods[i].Name != methods[j].Name {
			return methods[i].Name < methods[j].Name
		}
		return methods[i].Signature < methods[j].Signature
	})

	return methods, nil
}

// parseMethods extracts and processes contract methods
func parseMethods(parsedABI abi.ABI, methodIds map[string]string) ([]types.Method, error) {
	var methods []types.Method
	methodNames := make(map[string]int) // track name collisions

	// First pass: count method names for overload detection
	for _, method := range parsedABI.Methods {
		methodNames[method.Name]++
	}

	// Second pass: create method descriptors
	for _, method := range parsedABI.Methods {
		selector := methodIds[method.Sig]
		if selector == "" {
			return nil, fmt.Errorf("missing method identifier for %s", method.Sig)
		}

		// Generate method name with overload suffix if needed
		methodName := method.Name
		if methodNames[method.Name] > 1 {
			methodName = generateOverloadName(method.Name, method.Sig, selector)
		}

		// Parse inputs and outputs
		inputs, err := parseParameters(method.Inputs, false)
		if err != nil {
			return nil, fmt.Errorf("parsing inputs for method %s: %w", method.Sig, err)
		}

		outputs, err := parseParameters(method.Outputs, false)
		if err != nil {
			return nil, fmt.Errorf("parsing outputs for method %s: %w", method.Sig, err)
		}

		// Create input/output structs if needed
		var inputStruct, outputStruct *types.Struct

		if len(inputs) > 1 {
			inputStruct = &types.Struct{
				Name:   exportIdentifier(methodName) + "Input",
				Fields: parametersToFields(inputs),
			}
		}

		if len(outputs) > 1 {
			outputStruct = &types.Struct{
				Name:   exportIdentifier(methodName) + "Output",
				Fields: parametersToFields(outputs),
			}
		}

		methods = append(methods, types.Method{
			Name:         methodName,
			Signature:    method.Sig,
			Selector:     types.HexData(prefixHex(selector)),
			Inputs:       inputs,
			Outputs:      outputs,
			InputStruct:  inputStruct,
			OutputStruct: outputStruct,
		})
	}

	// Sort methods for deterministic output
	sort.Slice(methods, func(i, j int) bool {
		if methods[i].Name != methods[j].Name {
			return methods[i].Name < methods[j].Name
		}
		return methods[i].Signature < methods[j].Signature
	})

	return methods, nil
}

// parseEventsWithRegistry extracts and processes contract events using struct registry
func parseEventsWithRegistry(parsedABI abi.ABI, registry *structRegistry) ([]types.Event, error) {
	var events []types.Event

	for _, event := range parsedABI.Events {
		// Calculate event topic (hash of signature)
		topic := common.BytesToHash(crypto.Keccak256([]byte(event.Sig)))

		// Parse event inputs with registry
		inputs, err := parseParametersWithRegistry(event.Inputs, true, registry)
		if err != nil {
			return nil, fmt.Errorf("parsing inputs for event %s: %w", event.Sig, err)
		}

		// Create event struct
		eventStruct := &types.Struct{
			Name:   event.Name + "Event",
			Fields: parametersToFields(inputs),
		}

		// Convert common.Hash to types.Hash
		var typesHash types.Hash
		copy(typesHash[:], topic[:])
		
		events = append(events, types.Event{
			Name:   event.Name,
			Topic:  typesHash,
			Inputs: inputs,
			Struct: eventStruct,
		})
	}

	return events, nil
}

// parseEvents extracts and processes contract events
func parseEvents(parsedABI abi.ABI) ([]types.Event, error) {
	var events []types.Event

	for _, event := range parsedABI.Events {
		// Calculate event topic (hash of signature)
		topic := common.BytesToHash(crypto.Keccak256([]byte(event.Sig)))

		// Parse event inputs
		inputs, err := parseParameters(event.Inputs, true)
		if err != nil {
			return nil, fmt.Errorf("parsing inputs for event %s: %w", event.Sig, err)
		}

		// Create event struct
		eventStruct := &types.Struct{
			Name:   event.Name + "Event",
			Fields: parametersToFields(inputs),
		}

		// Convert common.Hash to types.Hash
		var typesHash types.Hash
		copy(typesHash[:], topic[:])
		
		events = append(events, types.Event{
			Name:   event.Name,
			Topic:  typesHash,
			Inputs: inputs,
			Struct: eventStruct,
		})
	}

	// Sort events for deterministic output
	sort.Slice(events, func(i, j int) bool {
		return events[i].Name < events[j].Name
	})

	return events, nil
}

// parseErrors extracts and processes contract errors
func parseErrors(parsedABI abi.ABI) ([]types.ContractError, error) {
	var errors []types.ContractError

	for _, abiError := range parsedABI.Errors {
		// Calculate error selector (first 4 bytes of signature hash)
		selector := common.BytesToHash(crypto.Keccak256([]byte(abiError.Sig))).Hex()[:10]

		// Parse error inputs
		inputs, err := parseParameters(abiError.Inputs, false)
		if err != nil {
			return nil, fmt.Errorf("parsing inputs for error %s: %w", abiError.Sig, err)
		}

		// Create error struct
		errorStruct := &types.Struct{
			Name:   abiError.Name + "Error",
			Fields: parametersToFields(inputs),
		}

		errors = append(errors, types.ContractError{
			Name:      abiError.Name,
			Signature: abiError.Sig,
			Selector:  types.HexData(selector),
			Inputs:    inputs,
			Struct:    errorStruct,
		})
	}

	// Sort errors for deterministic output
	sort.Slice(errors, func(i, j int) bool {
		return errors[i].Name < errors[j].Name
	})

	return errors, nil
}

// parseConstructor extracts constructor information
func parseConstructor(parsedABI abi.ABI, linkRefs map[string]map[string][]types.LinkRef) *types.Constructor {
	constructor := parsedABI.Constructor
	if constructor.Type != abi.Constructor {
		return nil
	}

	inputs, err := parseParameters(constructor.Inputs, false)
	if err != nil {
		// Log error but don't fail, constructor is optional
		return nil
	}

	var inputStruct *types.Struct
	if len(inputs) > 1 {
		inputStruct = &types.Struct{
			Name:   "ConstructorInput",
			Fields: parametersToFields(inputs),
		}
	}

	// Convert link references
	linkReferences := make(map[string][]types.LinkRef)
	for _, fileRefs := range linkRefs {
		for libName, refs := range fileRefs {
			for _, ref := range refs {
				linkReferences[libName] = append(linkReferences[libName], types.LinkRef{
					Start:  ref.Start,
					Length: ref.Length,
				})
			}
		}
	}

	return &types.Constructor{
		Signature:      constructor.Sig,
		Inputs:         inputs,
		InputStruct:    inputStruct,
		LinkReferences: linkReferences,
	}
}

// parseParametersWithRegistry converts ABI arguments to our parameter model using struct registry
func parseParametersWithRegistry(args abi.Arguments, allowIndexed bool, registry *structRegistry) ([]types.Parameter, error) {
	var params []types.Parameter

	for i, arg := range args {
		goType, err := mapSolidityToGoTypeWithRegistry(arg.Type, registry)
		if err != nil {
			return nil, fmt.Errorf("mapping type %s: %w", arg.Type.String(), err)
		}

		name := arg.Name
		if name == "" {
			name = fmt.Sprintf("Field%d", i+1) // 1-based indexing
		}

		params = append(params, types.Parameter{
			Name:    sanitizeIdentifier(name),
			Type:    goType,
			Indexed: allowIndexed && arg.Indexed,
		})
	}

	return params, nil
}

// parseParameters converts ABI arguments to our parameter model
func parseParameters(args abi.Arguments, allowIndexed bool) ([]types.Parameter, error) {
	var params []types.Parameter

	for i, arg := range args {
		goType, err := mapSolidityToGoType(arg.Type)
		if err != nil {
			return nil, fmt.Errorf("mapping type %s: %w", arg.Type.String(), err)
		}

		name := arg.Name
		if name == "" {
			name = fmt.Sprintf("Field%d", i+1) // 1-based indexing
		}

		params = append(params, types.Parameter{
			Name:    sanitizeIdentifier(name),
			Type:    goType,
			Indexed: allowIndexed && arg.Indexed,
		})
	}

	return params, nil
}

// parametersToFields converts parameters to struct fields
func parametersToFields(params []types.Parameter) []types.StructField {
	var fields []types.StructField

	for _, param := range params {
		jsonTag := strings.ToLower(param.Name)
		fields = append(fields, types.StructField{
			Name:    exportIdentifier(param.Name),
			Type:    param.Type,
			JSONTag: jsonTag,
		})
	}

	return fields
}

// mapSolidityToGoType maps Solidity types to Go types
func mapSolidityToGoType(abiType abi.Type) (types.GoType, error) {
	switch abiType.T {
	case abi.BoolTy:
		return types.GoTypeBool, nil
	case abi.StringTy:
		return types.GoTypeString, nil
	case abi.BytesTy:
		return types.GoTypeBytes, nil
	case abi.AddressTy:
		return types.GoTypeAddress, nil
	case abi.HashTy:
		return types.GoTypeHash, nil

	case abi.UintTy:
		if abiType.Size <= 64 {
			return mapUintType(abiType.Size), nil
		}
		return types.GoTypeBigInt, nil

	case abi.IntTy:
		if abiType.Size <= 64 {
			return mapIntType(abiType.Size), nil
		}
		return types.GoTypeBigInt, nil

	case abi.FixedBytesTy:
		return types.GoType{
			TypeName: fmt.Sprintf("[%d]byte", abiType.Size),
		}, nil

	case abi.SliceTy:
		elemType, err := mapSolidityToGoType(*abiType.Elem)
		if err != nil {
			return types.GoType{}, fmt.Errorf("mapping slice element type: %w", err)
		}
		return types.GoType{
			Import:   elemType.Import,
			TypeName: "[]" + elemType.TypeName,
			IsSlice:  true,
		}, nil

	case abi.ArrayTy:
		elemType, err := mapSolidityToGoType(*abiType.Elem)
		if err != nil {
			return types.GoType{}, fmt.Errorf("mapping array element type: %w", err)
		}
		return types.GoType{
			Import:   elemType.Import,
			TypeName: fmt.Sprintf("[%d]%s", abiType.Size, elemType.TypeName),
		}, nil

	case abi.TupleTy:
		// This function should not be called directly for TupleTy when we need struct registration
		// Use mapSolidityToGoTypeWithRegistry instead
		structName := extractStructName(abiType.TupleRawName)
		if structName == "" {
			structName = "AnonymousTuple" // fallback for truly anonymous tuples
		}
		return types.GoType{
			TypeName: structName,
		}, nil

	default:
		return types.GoType{}, fmt.Errorf("unsupported ABI type: %s", abiType.String())
	}
}

// mapUintType maps uint sizes to Go types
func mapUintType(size int) types.GoType {
	switch size {
	case 8:
		return types.GoTypeUint8
	case 16:
		return types.GoTypeUint16
	case 32:
		return types.GoTypeUint32
	case 64:
		return types.GoTypeUint64
	default:
		return types.GoTypeBigInt
	}
}

// mapIntType maps int sizes to Go types
func mapIntType(size int) types.GoType {
	switch size {
	case 8:
		return types.GoTypeInt8
	case 16:
		return types.GoTypeInt16
	case 32:
		return types.GoTypeInt32
	case 64:
		return types.GoTypeInt64
	default:
		return types.GoTypeBigInt
	}
}

// mapSolidityToGoTypeWithRegistry maps Solidity types to Go types and registers structs
func mapSolidityToGoTypeWithRegistry(abiType abi.Type, registry *structRegistry) (types.GoType, error) {
	switch abiType.T {
	case abi.SliceTy:
		elemType, err := mapSolidityToGoTypeWithRegistry(*abiType.Elem, registry)
		if err != nil {
			return types.GoType{}, fmt.Errorf("mapping slice element type: %w", err)
		}
		return types.GoType{
			Import:   elemType.Import,
			TypeName: "[]" + elemType.TypeName,
			IsSlice:  true,
		}, nil
	case abi.ArrayTy:
		elemType, err := mapSolidityToGoTypeWithRegistry(*abiType.Elem, registry)
		if err != nil {
			return types.GoType{}, fmt.Errorf("mapping array element type: %w", err)
		}
		return types.GoType{
			Import:   elemType.Import,
			TypeName: fmt.Sprintf("[%d]%s", abiType.Size, elemType.TypeName),
		}, nil
	case abi.TupleTy:
		// Extract struct name and register the struct definition
		structName := extractStructName(abiType.TupleRawName)
		if structName == "" {
			structName = "AnonymousTuple" // fallback for truly anonymous tuples
		}
		
		
		// Register this struct type for generation
		if registry != nil {
			registry.registerStruct(structName, abiType)
		}
		
		return types.GoType{
			TypeName: structName,
		}, nil
	default:
		// For non-composite types, use the original mapping function
		return mapSolidityToGoType(abiType)
	}
}

// extractStructName extracts a clean struct name from the raw tuple name
// Examples: 
//   "struct TestStructArray.User" -> "User"
//   "TestStructArrayUser" -> "User" (from TupleRawName format)
//   "TestContractUser" -> "User"
//   "struct MyContract.Company" -> "Company"
//   "" -> "" (anonymous tuple)
func extractStructName(rawName string) string {
	if rawName == "" {
		return ""
	}
	
	// Remove "struct " prefix if present
	if strings.HasPrefix(rawName, "struct ") {
		rawName = rawName[7:]
	}
	
	// Split on "." and take the last part (the actual struct name)
	parts := strings.Split(rawName, ".")
	if len(parts) > 1 {
		return exportIdentifier(parts[len(parts)-1])
	}
	
	// Handle TupleRawName format like "TestContractUser" -> "User"
	// Pattern: find the last capital letter that starts the struct name
	// This handles cases like "TestContractUser" -> "User", "MyContractCompany" -> "Company"
	for i := len(rawName) - 1; i > 0; i-- {
		if rawName[i] >= 'A' && rawName[i] <= 'Z' {
			// Found a capital letter, check if it's likely the start of the struct name
			// Simple heuristic: if it's not the first char and the previous isn't uppercase
			if i > 0 && rawName[i-1] >= 'a' && rawName[i-1] <= 'z' {
				return exportIdentifier(rawName[i:])
			}
		}
	}
	
	// For now, just use the full name as fallback
	return exportIdentifier(rawName)
}

// generateOverloadName creates a unique method name for overloaded functions
func generateOverloadName(baseName, signature, selector string) string {
	// Extract parameter types from signature: "foo(uint256,address)" -> ["uint256", "address"]
	start := strings.Index(signature, "(")
	end := strings.Index(signature, ")")
	if start == -1 || end == -1 || end <= start {
		// Fallback to selector-based naming
		return fmt.Sprintf("%s__%s", baseName, selector[2:])
	}

	paramStr := signature[start+1 : end]
	if paramStr == "" {
		return baseName + "_NoArgs"
	}

	// Split and normalize parameter types
	params := strings.Split(paramStr, ",")
	var normalizedParams []string
	for _, param := range params {
		param = strings.TrimSpace(param)
		normalized := normalizeTypeForNaming(param)
		normalizedParams = append(normalizedParams, normalized)
	}

	candidate := fmt.Sprintf("%s_%s", baseName, strings.Join(normalizedParams, "_"))

	// If still too complex, fall back to selector
	if len(candidate) > 50 {
		return fmt.Sprintf("%s__%s", baseName, selector[2:])
	}

	return candidate
}

// normalizeTypeForNaming converts Solidity types to naming-friendly strings
func normalizeTypeForNaming(typeName string) string {
	// Handle arrays
	if strings.HasSuffix(typeName, "[]") {
		base := strings.TrimSuffix(typeName, "[]")
		return normalizeTypeForNaming(base) + "Array"
	}

	// Handle fixed arrays
	if strings.Contains(typeName, "[") && strings.Contains(typeName, "]") {
		base := typeName[:strings.Index(typeName, "[")]
		return normalizeTypeForNaming(base) + "FixedArray"
	}

	// Common type mappings
	switch typeName {
	case "uint256":
		return "Uint256"
	case "address":
		return "Address"
	case "bool":
		return "Bool"
	case "string":
		return "String"
	case "bytes":
		return "Bytes"
	default:
		// Handle uintN, intN, bytesN
		if strings.HasPrefix(typeName, "uint") {
			if size := typeName[4:]; size != "" {
				return "Uint" + size
			}
		}
		if strings.HasPrefix(typeName, "int") {
			if size := typeName[3:]; size != "" {
				return "Int" + size
			}
		}
		if strings.HasPrefix(typeName, "bytes") {
			if size := typeName[5:]; size != "" {
				return "Bytes" + size
			}
		}

		// Capitalize first letter for other types
		return exportIdentifier(typeName)
	}
}

// Utility functions

// sanitizePackageName converts contract names to valid Go package names
func sanitizePackageName(name string) string {
	var result strings.Builder

	for _, r := range strings.ToLower(name) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			result.WriteRune(r)
		}
	}

	pkg := result.String()
	if pkg == "" || (pkg[0] >= '0' && pkg[0] <= '9') {
		pkg = "contract" + pkg
	}

	return pkg
}

// sanitizeIdentifier converts names to valid Go identifiers
func sanitizeIdentifier(name string) string {
	if name == "" {
		return "Field"
	}

	var result strings.Builder
	first := true

	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (!first && r >= '0' && r <= '9') || r == '_' {
			result.WriteRune(r)
		} else {
			result.WriteRune('_')
		}
		first = false
	}

	id := result.String()
	if id == "" || (id[0] >= '0' && id[0] <= '9') {
		id = "Field_" + id
	}

	return id
}

// exportIdentifier makes an identifier exportable (uppercase first letter)
func exportIdentifier(name string) string {
	if name == "" {
		return "Field"
	}
	// Handle names starting with underscore
	if strings.HasPrefix(name, "_") {
		name = name[1:] // Remove leading underscore
		if name == "" {
			return "Field"
		}
	}
	return strings.ToUpper(name[:1]) + name[1:]
}

// prefixHex adds 0x prefix if not present
func prefixHex(hex string) string {
	if hex == "" {
		return ""
	}
	if strings.HasPrefix(hex, "0x") {
		return hex
	}
	return "0x" + hex
}
