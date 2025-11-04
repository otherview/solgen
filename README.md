# solgen ⚡

[![Test](https://github.com/otherview/solgen/workflows/Test/badge.svg)](https://github.com/otherview/solgen/actions/workflows/test.yml)
[![Release](https://github.com/otherview/solgen/workflows/Release/badge.svg)](https://github.com/otherview/solgen/actions/workflows/release.yml)
[![Edge Build](https://github.com/otherview/solgen/workflows/Edge%20Build/badge.svg)](https://github.com/otherview/solgen/actions/workflows/edge.yml)

Lightning-fast Go bindings generator for Solidity contracts. Zero dependencies, maximum simplicity.

## ✨ Features

- 🔄 **Pipeline-first**: Reads `solc` output, writes clean Go packages
- 🚀 **Zero blockchain deps**: Pure Go, no external libraries
- 📦 **One package per contract**: Clean, isolated bindings
- 🎯 **Simple API**: Contract metadata + utility functions
- 🔧 **Method overloads**: Smart naming for overloaded functions
- ⚠️ **Custom errors**: Full Solidity error support
- 📡 **Events**: Complete event type definitions

## 🚀 Quick Start

```bash
# Docker (recommended)
docker pull otherview/solgen

# Or install locally
go install github.com/otherview/solgen/cmd/solgen@latest
```

> ⚠️ **Docker Image Update**: The official Solidity compiler Docker image has moved from `ethereum/solc` (deprecated) to `ghcr.io/argotorg/solc`. All examples below use the new official image.

## 💡 Usage

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

### 🎯 Simple API

```go
package mycontract

// Contract metadata
func ABI() string                 // Contract ABI JSON
func HexBytecode() string         // Creation bytecode
func HexDeployedBytecode() string // Runtime bytecode

// Method info
func TransferMethod() MethodInfo {
  return MethodInfo{
    Name:      "transfer",
    Signature: "transfer(address,uint256)",
    Selector:  "0xa9059cbb",
  }
}

// Event info  
func GetTransferEvent() EventInfo {
  return EventInfo{
    Name:  "Transfer", 
    Topic: common.HexToHash("0xddf252ad..."),
  }
}

// Registry access
methods := Methods()  // MethodRegistry{Transfer: MethodInfo{...}}
events := Events()    // EventRegistry{Transfer: EventInfo{...}}
errors := Errors()    // ErrorRegistry{InsufficientBalance: ErrorInfo{...}}

// Type definitions
type TransferEvent struct {
  From  common.Address `json:"from"`
  To    common.Address `json:"to"`
  Value *big.Int       `json:"value"`
}

// Utility functions
data, err := HexToBytes("0x1234")    // Parse hex
hex := BytesToHex(data)              // Format hex
data = MustHexToBytes("0x1234")      // Parse hex (panic on error)
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
| `T[]` | `[]T` | `[]common.Address{...}` |
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