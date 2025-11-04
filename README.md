# solgen ⚡

[![Test](https://github.com/otherview/solgen/workflows/Test/badge.svg)](https://github.com/otherview/solgen/actions/workflows/test.yml)
[![Release](https://github.com/otherview/solgen/workflows/Release/badge.svg)](https://github.com/otherview/solgen/actions/workflows/release.yml)
[![Edge Build](https://github.com/otherview/solgen/workflows/Edge%20Build/badge.svg)](https://github.com/otherview/solgen/actions/workflows/edge.yml)

Lightning-fast Go bindings generator for Solidity contracts. **Zero dependencies, maximum simplicity.**

Generate clean, type-safe Go code from Solidity contracts in seconds. No external dependencies in generated bindings - just pure Go ready for production.

## ✨ Features

- 🚀 **Zero Dependencies**: Generated bindings are completely self-contained
- 🎯 **Type-Safe API**: Clean chaining API with compile-time safety
- 📦 **One Package per Contract**: Isolated, clean Go packages
- ⚡ **Production Ready**: Built-in ABI encoding/decoding, no external libs
- 🔧 **Method Overloads**: Smart naming for overloaded functions  
- ⚠️ **Custom Errors**: Full Solidity error support with type-safe decoding
- 📊 **Event Logs**: Complete event parsing with structured data
- 🔄 **Pipeline-First**: Reads `solc` output, writes clean Go code

## 💻 Usage Examples

Clean, type-safe bindings that work with any Ethereum client:

### ✨ Basic Usage
```go
import "yourproject/generated/simpletoken"

// Pack method calls for transactions
transferData := simpletoken.Methods().TransferMethod().Pack(recipient, amount)
approveData := simpletoken.Methods().ApproveMethod().MustPack(spender, amount)

// Decode return values from eth_call  
balance := simpletoken.Methods().BalanceOfMethod().MustDecode(returnData)
success := simpletoken.Methods().TransferMethod().MustDecode(returnData)
tokenName := simpletoken.Methods().NameMethod().MustDecode(returnData)
```

### 📊 Event & Error Handling
```go
// Decode event logs
transferEvent := simpletoken.Events().TransferEvent().MustDecode(logData)
fmt.Printf("Transfer: %s to %s, amount: %s ETH\n", 
    transferEvent.From, transferEvent.To, weiToEth(transferEvent.Value))

// Handle custom errors from reverted transactions  
if revertData != nil {
    error := simpletoken.Errors().InsufficientBalanceError().MustDecode(revertData)
    fmt.Printf("Error: insufficient balance - requested %s, available %s\n",
        weiToEth(error.Requested), weiToEth(error.Available))
}
```

### 🌐 Blockchain Integration
```go
import "github.com/ethereum/go-ethereum/ethclient"

// Works seamlessly with ethclient or any Ethereum library
client, _ := ethclient.Dial("https://mainnet.infura.io/v3/YOUR_PROJECT_ID")
contractAddr := common.HexToAddress("0x...")

// Read contract state
callData := simpletoken.Methods().BalanceOfMethod().MustPack(userAddress)
result, _ := client.CallContract(ctx, ethereum.CallMsg{
    To: &contractAddr, Data: callData.Bytes(),
}, nil)
balance := simpletoken.Methods().BalanceOfMethod().MustDecode(result)

// Send transactions  
tx := types.NewTransaction(nonce, contractAddr, big.NewInt(0), gasLimit, gasPrice, 
    simpletoken.Methods().TransferMethod().MustPack(recipient, amount).Bytes())
```

### 🔗 Zero Dependencies
Generated bindings are completely self-contained:
```go
// No external imports needed in generated code!
type Address [20]byte  // Custom address type
type Hash [32]byte     // Custom hash type  
type HexData string    // Convenient hex handling

// Built-in ABI encoding/decoding
// Type-safe method calls
// Clean error handling
```

> 🚀 **[Complete Example](examples/simpletoken/)** - See the full working example with detailed comments showing all features in action!



## 🚀 Quick Start

```bash
# Docker (recommended)
docker pull otherview/solgen

# Or install locally
go install github.com/otherview/solgen/cmd/solgen@latest
```

> ⚠️ **Docker Image Update**: The official Solidity compiler Docker image has moved from `ethereum/solc` (deprecated) to `ghcr.io/argotorg/solc`. All examples below use the new official image.

### ⚡ One-liner magic
```bash
# Standard: Contract info + bytecode (recommended)
solc --combined-json abi,bin,bin-runtime,hashes contracts/*.sol | \
  solgen --out generated

# Minimum: Just contract info (no bytecode functions)
solc --combined-json abi,hashes contracts/*.sol | \
  solgen --out generated
```

### 🐳 Docker pipeline
```bash
# Latest stable version (recommended for development)
sh -c "docker run --rm -v $(pwd):/src ghcr.io/argotorg/solc:stable \
  --combined-json abi,bin,bin-runtime,hashes \
  --optimize /src/contracts/*.sol | \
  docker run --rm -i -v $(pwd):/src otherview/solgen --out /src/generated"

# Pinned version (recommended for production)
sh -c "docker run --rm -v $(pwd):/src ghcr.io/argotorg/solc:0.8.20 \
  --combined-json abi,bin,bin-runtime,hashes \
  --optimize /src/contracts/*.sol | \
  docker run --rm -i -v $(pwd):/src otherview/solgen --out /src/generated"
```

### 🔄 go:generate integration
```go
// Using stable version
//go:generate sh -c "docker run --rm -v $(pwd):/src ghcr.io/argotorg/solc:stable --combined-json abi,bin,bin-runtime,hashes /src/contracts/*.sol | docker run --rm -i -v $(pwd):/src otherview/solgen --out /src/generated"

// Or with local solc installation
//go:generate sh -c "solc --combined-json abi,bin,bin-runtime,hashes contracts/*.sol | solgen --out generated"
```

> 📚 **More examples**: See [EXAMPLES.md](EXAMPLES.md) for advanced usage, CI/CD integration, and platform-specific examples

### ⚙️ Options

**solgen**
- `--out` (required): Output directory
- `--verbose`: Detailed output

**solc** (required fields)
- 🎯 **Minimum**: `--combined-json abi,hashes` (contract info only)
- ⚡ **Standard**: `--combined-json abi,bin,bin-runtime,hashes` (+ bytecode functions) 
- 🔧 **Options**: `--optimize`, `--optimize-runs 200`

**Docker Images**
- ✅ **Current**: `ghcr.io/argotorg/solc:stable` (latest) or `ghcr.io/argotorg/solc:0.8.20` (pinned)
- ❌ **Deprecated**: `ethereum/solc:0.8.20` (still works but not recommended)

## 📁 Generated Code

Each contract gets a clean Go package:

```
generated/
  mycontract/
    mycontract.go    # Clean, dependency-free bindings
```

### 🔄 Type Mapping

| Solidity | Go | Example |
|----------|----|---------| 
| `bool` | `bool` | `true` |
| `string` | `string` | `"hello"` |
| `address` | `common.Address` | `0x742d...` |
| `bytes` | `[]byte` | `[]byte{0x12, 0x34}` |
| `bytes32` | `[32]byte` | `[32]byte{...}` |
| `uint256` | `*big.Int` | `big.NewInt(123)` |
| `uint64` | `uint64` | `uint64(123)` |
| `int256` | `*big.Int` | `big.NewInt(-123)` |
| `T[]` | `[]T` | `[]Address{...}` |
| `T[N]` | `[N]T` | `[3]uint256{...}` |

## 📋 Requirements

- **Docker**: For `ghcr.io/argotorg/solc` + `otherview/solgen` containers
- **Local**: Go 1.21+ and `solc` binary

## 🔄 Migration from ethereum/solc

If you're currently using the deprecated `ethereum/solc` image:

```bash
# Old (deprecated)
docker run --rm ethereum/solc:0.8.20 --combined-json abi,bin,bin-runtime,hashes

# New (recommended)
docker run --rm ghcr.io/argotorg/solc:0.8.20 --combined-json abi,bin,bin-runtime,hashes
# or
docker run --rm ghcr.io/argotorg/solc:stable --combined-json abi,bin,bin-runtime,hashes
```

The functionality is identical - just replace the image name.

## 🛠️ Development

```bash
# Run tests
go test ./...

# Build from source  
go build ./cmd/solgen
```

## 📄 License

MIT License - see [LICENSE](LICENSE) file.

## 🤝 Contributing

1. Fork it
2. Create your feature branch
3. Add tests
4. Submit a pull request

## 🔧 Troubleshooting  

**Pipeline fails?** Test solc and solgen separately with `--verbose`

**Docker permissions?** Use `--user $(id -u):$(id -g)` or `chown` after generation

**Package conflicts?** Rename contracts - package names are derived from contract names (lowercase, alphanumeric only)

> 🔍 **Detailed troubleshooting**: See [EXAMPLES.md](EXAMPLES.md) for step-by-step debugging, platform-specific issues, and advanced solutions