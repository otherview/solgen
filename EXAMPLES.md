# Advanced Examples & Integration Guide

This document provides comprehensive examples for using solgen in various environments and workflows.

> ⚠️ **Docker Image Update**: All examples below use the new official `ghcr.io/argotorg/solc` image. The old `ethereum/solc` image still works but is deprecated.

## 🚀 Quick Examples

### Standard mode (recommended)
```bash
# Latest stable version
sh -c "docker run --rm -v $(pwd):/sources ghcr.io/argotorg/solc:stable \
  --combined-json abi,bin,bin-runtime,hashes \
  --optimize --optimize-runs 200 \
  /sources/contracts/MyToken.sol | \
  docker run --rm -i -v $(pwd):/sources otherview/solgen \
  --out /sources/generated"

# Pinned version (production recommended)
sh -c "docker run --rm -v $(pwd):/sources ghcr.io/argotorg/solc:0.8.20 \
  --combined-json abi,bin,bin-runtime,hashes \
  --optimize --optimize-runs 200 \
  /sources/contracts/MyToken.sol | \
  docker run --rm -i -v $(pwd):/sources otherview/solgen \
  --out /sources/generated"
```

### Minimum mode (contract info only)
```bash
# Fast development builds - just contract ABI and method hashes
sh -c "docker run --rm -v $(pwd):/sources ghcr.io/argotorg/solc:stable \
  --combined-json abi,hashes \
  /sources/contracts/MyToken.sol | \
  docker run --rm -i -v $(pwd):/sources otherview/solgen \
  --out /sources/generated"
```

### Multiple contracts with imports
```bash
sh -c "docker run --rm -v $(pwd):/sources ghcr.io/argotorg/solc:stable \
  --combined-json abi,bin,bin-runtime,hashes \
  --optimize --optimize-runs 200 \
  --base-path /sources \
  --include-path /sources/node_modules \
  /sources/contracts/*.sol | \
  docker run --rm -i -v $(pwd):/sources otherview/solgen \
  --out /sources/generated --verbose"
```

## 🔧 Advanced Compilation

### Production build (high optimization)
```bash
# Pinned version for reproducible builds
sh -c "docker run --rm -v $(pwd):/sources ghcr.io/argotorg/solc:0.8.20 \
  --combined-json abi,bin,bin-runtime,hashes \
  --optimize --optimize-runs 10000 \
  --evm-version shanghai \
  --via-ir \
  /sources/contracts/*.sol | \
  docker run --rm -i -v $(pwd):/sources otherview/solgen \
  --out /sources/generated"
```

### Debug build (no optimization, minimum mode)
```bash
# Debug with just contract info (faster compilation, smaller output)
sh -c "docker run --rm -v $(pwd):/sources ghcr.io/argotorg/solc:stable \
  --combined-json abi,hashes \
  /sources/contracts/*.sol | \
  docker run --rm -i -v $(pwd):/sources otherview/solgen \
  --out /sources/generated --verbose"
```

### Custom EVM version
```bash
sh -c "docker run --rm -v $(pwd):/sources ghcr.io/argotorg/solc:0.8.20 \
  --combined-json abi,bin,bin-runtime,hashes \
  --optimize --optimize-runs 1000 \
  --evm-version paris \
  /sources/contracts/MyContract.sol | \
  docker run --rm -i -v $(pwd):/sources otherview/solgen \
  --out /sources/generated"
```

## 🔄 CI/CD Integration

### GitHub Actions
```yaml
name: Generate Solidity Bindings

on:
  push:
    paths: ['contracts/**']
  pull_request:
    paths: ['contracts/**']

jobs:
  generate:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Generate Go bindings
      run: |
        docker run --rm \
          -v ${{ github.workspace }}:/sources \
          ghcr.io/argotorg/solc:0.8.20 \
          --combined-json abi,bin,bin-runtime,hashes \
          --optimize --optimize-runs 200 \
          --base-path /sources \
          /sources/contracts/*.sol | \
        docker run --rm -i \
          -v ${{ github.workspace }}:/sources \
          otherview/solgen \
          --out /sources/generated
    
    - name: Commit generated files
      run: |
        git config --local user.email "action@github.com"
        git config --local user.name "GitHub Action"
        git add generated/
        git commit -m "Update generated bindings" || exit 0
        git push
```

### Makefile Integration
```makefile
.PHONY: generate-bindings clean-bindings generate-dev generate-prod

SOLC_VERSION ?= 0.8.20
SOLC_IMAGE ?= ghcr.io/argotorg/solc
CONTRACTS_DIR = contracts
GENERATED_DIR = generated

# Development build (minimum mode - faster compilation)
generate-dev:
	sh -c "docker run --rm -v $(PWD):/sources $(SOLC_IMAGE):stable \
		--combined-json abi,hashes \
		/sources/$(CONTRACTS_DIR)/*.sol | \
		docker run --rm -i -v $(PWD):/sources otherview/solgen \
		--out /sources/$(GENERATED_DIR) --verbose"

# Standard build (recommended)
generate-bindings:
	sh -c "docker run --rm -v $(PWD):/sources $(SOLC_IMAGE):$(SOLC_VERSION) \
		--combined-json abi,bin,bin-runtime,hashes \
		--optimize --optimize-runs 200 \
		--base-path /sources \
		/sources/$(CONTRACTS_DIR)/*.sol | \
		docker run --rm -i -v $(PWD):/sources otherview/solgen \
		--out /sources/$(GENERATED_DIR)"

# Production build (high optimization)
generate-prod:
	sh -c "docker run --rm -v $(PWD):/sources $(SOLC_IMAGE):$(SOLC_VERSION) \
		--combined-json abi,bin,bin-runtime,hashes \
		--optimize --optimize-runs 10000 \
		--evm-version shanghai \
		--via-ir \
		/sources/$(CONTRACTS_DIR)/*.sol | \
		docker run --rm -i -v $(PWD):/sources otherview/solgen \
		--out /sources/$(GENERATED_DIR)"

clean-bindings:
	rm -rf $(GENERATED_DIR)
```

### Docker Compose
```yaml
version: '3.8'
services:
  generate-bindings:
    image: ethereum/solc:0.8.20
    volumes:
      - .:/sources
    command: |
      sh -c "docker run --rm -v /sources:/sources ghcr.io/argotorg/solc:0.8.20 \
        --combined-json abi,bin,bin-runtime,hashes \
        --optimize --optimize-runs 200 \
        /sources/contracts/*.sol | \
        docker run --rm -i -v /sources:/sources otherview/solgen \
        --out /sources/generated"
```

## 🛠️ Shell Functions & Aliases

Add these to your `~/.bashrc` or `~/.zshrc`:

### Basic function
```bash
solgen-pipeline() {
  local solc_version=${1:-stable}
  local input_pattern=${2:-"contracts/*.sol"}
  local output_dir=${3:-"generated"}
  
  sh -c "docker run --rm -v $(pwd):/sources ghcr.io/argotorg/solc:$solc_version \
    --combined-json abi,bin,bin-runtime,hashes \
    --optimize --optimize-runs 200 \
    /sources/$input_pattern | \
    docker run --rm -i -v $(pwd):/sources otherview/solgen \
    --out /sources/$output_dir"
}

# Usage: solgen-pipeline stable "contracts/MyContract.sol" "generated"
# Usage: solgen-pipeline 0.8.20 "contracts/MyContract.sol" "generated"
```

### Advanced function with options
```bash
solgen-advanced() {
  local solc_version=${1:-stable}
  local input_pattern=${2:-"contracts/*.sol"}
  local output_dir=${3:-"generated"}
  local optimize_runs=${4:-200}
  local evm_version=${5:-""}
  local includes=${6:-""}
  
  local solc_cmd="docker run --rm -v $(pwd):/sources ghcr.io/argotorg/solc:$solc_version \
    --combined-json abi,bin,bin-runtime,hashes \
    --optimize --optimize-runs $optimize_runs"
  
  if [ ! -z "$evm_version" ]; then
    solc_cmd="$solc_cmd --evm-version $evm_version"
  fi
  
  if [ ! -z "$includes" ]; then
    solc_cmd="$solc_cmd --base-path /sources --include-path /sources/$includes"
  fi
  
  solc_cmd="$solc_cmd /sources/$input_pattern"
  
  sh -c "$solc_cmd | \
    docker run --rm -i -v $(pwd):/sources otherview/solgen \
    --out /sources/$output_dir --verbose"
}

# Usage: solgen-advanced stable "contracts/*.sol" "generated" 1000 "paris" "node_modules"
# Usage: solgen-advanced 0.8.20 "contracts/*.sol" "generated" 1000 "paris" "node_modules"
```

## 🌍 Platform-Specific Examples

### macOS with Apple Silicon
```bash
# Force x86_64 architecture for consistency
sh -c "docker run --platform linux/amd64 --rm -v $(pwd):/sources ghcr.io/argotorg/solc:0.8.20 \
  --combined-json abi,bin,bin-runtime,hashes \
  --optimize --optimize-runs 200 \
  /sources/contracts/*.sol | \
  docker run --platform linux/amd64 --rm -i -v $(pwd):/sources otherview/solgen \
  --out /sources/generated"
```

### Windows PowerShell
```powershell
# PowerShell syntax
docker run --rm -v ${PWD}:/sources ghcr.io/argotorg/solc:0.8.20 `
  --combined-json abi,bin,bin-runtime,hashes `
  --optimize --optimize-runs 200 `
  /sources/contracts/*.sol | `
  docker run --rm -i -v ${PWD}:/sources otherview/solgen `
  --out /sources/generated
```

### Linux with permission handling
```bash
# Run with current user to avoid permission issues
sh -c "docker run --rm --user $(id -u):$(id -g) -v $(pwd):/sources ghcr.io/argotorg/solc:0.8.20 \
  --combined-json abi,bin,bin-runtime,hashes \
  --optimize --optimize-runs 200 \
  /sources/contracts/*.sol | \
  docker run --rm -i --user $(id -u):$(id -g) -v $(pwd):/sources otherview/solgen \
  --out /sources/generated"
```

## 🔍 Debugging & Testing

### Step-by-step debugging
```bash
echo "=== Testing solc compilation ==="
docker run --rm -v $(pwd):/sources ghcr.io/argotorg/solc:0.8.20 \
  --combined-json abi,bin,bin-runtime,hashes \
  --optimize --optimize-runs 200 \
  /sources/contracts/MyContract.sol > debug-output.json

echo "=== Checking JSON validity ==="
cat debug-output.json | jq . > /dev/null && echo "✅ Valid JSON" || echo "❌ Invalid JSON"

echo "=== Inspecting contracts ==="
cat debug-output.json | jq '.contracts | keys'

echo "=== Testing solgen processing ==="
cat debug-output.json | docker run --rm -i -v $(pwd):/sources otherview/solgen \
  --out /sources/debug-generated --verbose
```

### Minimal test case
```bash
# Test with minimal contract (minimum mode)
echo 'contract Test { function test() public {} }' > test.sol

sh -c "docker run --rm -v $(pwd):/sources ghcr.io/argotorg/solc:stable \
  --combined-json abi,hashes \
  /sources/test.sol | \
  docker run --rm -i -v $(pwd):/sources otherview/solgen \
  --out /sources/test-output --verbose"

# Check results
ls -la test-output/
rm test.sol
```

### Performance testing
```bash
# Time the entire pipeline
time sh -c "docker run --rm -v $(pwd):/sources ghcr.io/argotorg/solc:0.8.20 \
  --combined-json abi,bin,bin-runtime,hashes \
  --optimize --optimize-runs 200 \
  /sources/contracts/*.sol | \
  docker run --rm -i -v $(pwd):/sources otherview/solgen \
  --out /sources/generated"
```

## 📋 Best Practices

### ✅ Do's

**Choose the right mode for your needs:**
```bash
# ✅ Standard mode (recommended) - includes bytecode functions
--combined-json abi,bin,bin-runtime,hashes

# ✅ Minimum mode - just contract info (faster, smaller output)
--combined-json abi,hashes
```

**Use absolute paths in containers:**
```bash
# ✅ Correct
/sources/contracts/MyContract.sol
```

**Pin specific versions:**
```bash
# ✅ Reproducible builds
ethereum/solc:0.8.20
otherview/solgen:v1.0.0
```

**Handle multi-line commands properly:**
```bash
# ✅ Option 1: Use sh -c with quotes
sh -c "command1 | command2"

# ✅ Option 2: Use backslashes
docker run --rm \
  -v $(pwd):/sources \
  ethereum/solc:0.8.20
```

### ❌ Don'ts

**Missing required fields:**
```bash
# ❌ Missing hashes will cause failures
--combined-json abi,bin,bin-runtime

# ❌ Missing abi will cause failures
--combined-json bin,bin-runtime,hashes
```

**Relative paths in containers:**
```bash
# ❌ May not work reliably
contracts/MyContract.sol
```

**Unversioned images:**
```bash
# ❌ Non-reproducible
ethereum/solc:latest
```

## 🆚 Comparison: Traditional vs Pipeline

### Traditional solc workflow
```bash
# Step 1: Compile to separate files
docker run --rm -v $(pwd):/sources ghcr.io/argotorg/solc:0.8.20 \
  --evm-version paris --optimize --optimize-runs 200 \
  -o /sources/compiled --abi --bin-runtime /sources/staker.sol

# Result: Creates separate files
# - compiled/staker.abi
# - compiled/staker.bin-runtime

# Step 2: Manually create bindings (complex, limited functionality)
abigen --abi compiled/staker.abi --bin compiled/staker.bin-runtime --pkg staker --out staker.go
```

### solgen pipeline
```bash
# Single integrated command
sh -c "docker run --rm -v $(pwd):/sources ghcr.io/argotorg/solc:0.8.20 \
  --combined-json abi,bin,bin-runtime,hashes \
  --optimize --optimize-runs 200 \
  --evm-version paris \
  /sources/staker.sol | \
  docker run --rm -i -v $(pwd):/sources otherview/solgen \
  --out /sources/generated"

# Result: Complete Go package with rich typed bindings
# - generated/staker/staker.go (includes ABI, bytecode, typed methods, events, errors, utilities)
```

## 🚨 Troubleshooting

### Pipeline failures
```bash
# 1. Test solc separately
docker run --rm -v $(pwd):/sources ghcr.io/argotorg/solc:0.8.20 \
  --combined-json abi,bin,bin-runtime,hashes \
  /sources/contracts/MyContract.sol

# 2. Validate JSON output
docker run --rm -v $(pwd):/sources ghcr.io/argotorg/solc:0.8.20 \
  --combined-json abi,bin,bin-runtime,hashes \
  /sources/contracts/MyContract.sol | jq .

# 3. Test solgen separately
echo '{"contracts":{"test.sol:Test":{"abi":[],"bin":"","bin-runtime":"","hashes":{}}}}' | \
  docker run --rm -i -v $(pwd):/sources otherview/solgen --out /tmp/test
```

### File permissions (Docker)
```bash
# Fix ownership after generation
sudo chown -R $(id -u):$(id -g) ./generated

# Or run with current user from start
--user $(id -u):$(id -g)
```

### Common compilation errors
- ✅ Verify Solidity files compile with specified solc version
- ✅ Use absolute container paths: `/sources/contracts/File.sol`
- ✅ Include required paths: `--base-path` and `--include-path`
- ✅ Always include `hashes` in combined JSON output
- ✅ Use minimum mode (`abi,hashes`) for development and testing
- ✅ Use standard mode (`abi,bin,bin-runtime,hashes`) for production deployments

## 🔄 Migration from ethereum/solc

### Quick Migration Guide

If you're currently using `ethereum/solc`, simply replace the image name:

```bash
# Old (deprecated)
ethereum/solc:0.8.20

# New (official)
ghcr.io/argotorg/solc:0.8.20   # Pinned version
ghcr.io/argotorg/solc:stable    # Latest stable
```

### Version Mapping

| ethereum/solc | ghcr.io/argotorg/solc | Notes |
|---------------|------------------------|-------|
| `ethereum/solc:0.8.20` | `ghcr.io/argotorg/solc:0.8.20` | Same version |
| `ethereum/solc:latest` | `ghcr.io/argotorg/solc:stable` | Latest stable |
| Not available | `ghcr.io/argotorg/solc:nightly` | Development builds |

### Example Migration

```bash
# Before
sh -c "docker run --rm -v $(pwd):/src ethereum/solc:0.8.20 \
  --combined-json abi,bin,bin-runtime,hashes /src/contracts/*.sol | \
  docker run --rm -i -v $(pwd):/src otherview/solgen --out /src/generated"

# After (pinned version - recommended for CI/production)
sh -c "docker run --rm -v $(pwd):/src ghcr.io/argotorg/solc:0.8.20 \
  --combined-json abi,bin,bin-runtime,hashes /src/contracts/*.sol | \
  docker run --rm -i -v $(pwd):/src otherview/solgen --out /src/generated"

# After (stable version - recommended for development)
sh -c "docker run --rm -v $(pwd):/src ghcr.io/argotorg/solc:stable \
  --combined-json abi,bin,bin-runtime,hashes /src/contracts/*.sol | \
  docker run --rm -i -v $(pwd):/src otherview/solgen --out /src/generated"
```

---

💡 **Need help?** Check the main [README.md](README.md) for basic usage or open an issue on GitHub.