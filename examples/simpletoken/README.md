# SimpleToken Example

This example demonstrates how to use the generated Go bindings from solgen with a real Solidity contract.

## What's Included

- `contracts/SimpleToken.sol` - A complete ERC20-like token contract
- `generated/simpletoken/` - Generated Go bindings (zero external dependencies!)
- `main.go` - Comprehensive usage examples with detailed comments
- `go.mod` - Module setup with ethclient for blockchain integration

## Features Demonstrated

### 1. **Contract Metadata** 📋
- Access ABI JSON
- Deployment and runtime bytecode
- Contract information

### 2. **Method Call Encoding** 📦
- Pack method arguments for transactions
- Type-safe parameter encoding
- Clean chaining API: `Methods().TransferMethod().Pack(...)`

### 3. **Return Value Decoding** 📤
- Decode eth_call responses
- Type-safe return values
- Handle different return types (uint256, bool, string)

### 4. **Event Log Decoding** 📊
- Parse blockchain event logs
- Extract event parameters
- Clean event API: `Events().TransferEvent().MustDecode(...)`

### 5. **Custom Error Decoding** ⚠️
- Decode reverted transaction errors
- Parse custom error parameters
- Error API: `Errors().InsufficientBalanceError().MustDecode(...)`

### 6. **Blockchain Integration** 🌐
- Real ethereum client integration
- Transaction sending
- Event filtering
- Complete workflow examples

## Running the Example

```bash
# Generate the bindings (already done)
docker run --rm -v $(pwd):/sources ghcr.io/argotorg/solc:stable \
  --combined-json abi,bin,bin-runtime,hashes \
  --optimize /sources/contracts/SimpleToken.sol | \
  ../../solgen --out generated

# Install dependencies
go mod tidy

# Run the example
go run main.go
```

## Key Highlights

### Clean, Type-Safe API
```go
// Pack method calls
transferData := simpletoken.Methods().TransferMethod().Pack(recipient, amount)

// Decode return values  
balance := simpletoken.Methods().BalanceOfMethod().MustDecode(returnData)

// Parse events
event := simpletoken.Events().TransferEvent().MustDecode(logData)

// Handle errors
error := simpletoken.Errors().InsufficientBalanceError().MustDecode(revertData)
```

### Zero External Dependencies
The generated bindings are completely self-contained:
- No go-ethereum imports in generated code
- Custom Address and Hash types
- Built-in ABI encoding/decoding
- Ready for any Ethereum client library

### Production Ready
- Type-safe operations
- Comprehensive error handling
- Memory efficient
- Clean, readable generated code

## Generated Code Structure

```
generated/simpletoken/
└── simpletoken.go              # Complete bindings
    ├── Contract metadata       # ABI, bytecode
    ├── Custom types            # Address, Hash, HexData
    ├── Method encoders         # Pack arguments
    ├── Return decoders         # Decode responses  
    ├── Event decoders          # Parse logs
    ├── Error decoders          # Handle reverts
    └── Registry API            # Clean method/event/error access
```

This example shows everything you need to integrate solgen-generated bindings into your Go applications!