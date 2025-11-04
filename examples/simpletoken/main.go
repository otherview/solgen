package main

import (
	"fmt"
	"log"
	"math/big"

	// Imports for blockchain integration (used in commented examples)
	// "context"
	// "crypto/ecdsa"
	// "github.com/ethereum/go-ethereum/accounts/abi/bind"
	// "github.com/ethereum/go-ethereum/common"
	// "github.com/ethereum/go-ethereum/core/types"
	// "github.com/ethereum/go-ethereum/crypto"
	// "github.com/ethereum/go-ethereum/ethclient"

	// Import the generated bindings - zero external dependencies!
	"simpletoken-example/generated/simpletoken"
)

func main() {
	fmt.Println("🚀 SimpleToken Example - Generated Go Bindings")
	fmt.Println("==============================================")

	// Example 1: Basic Contract Information
	fmt.Println("\n📋 1. Contract Metadata")
	demonstrateMetadata()

	// Example 2: Method Call Encoding
	fmt.Println("\n📦 2. Method Call Encoding")
	demonstrateMethodEncoding()

	// Example 3: Return Value Decoding
	fmt.Println("\n📤 3. Return Value Decoding")
	demonstrateReturnDecoding()

	// Example 4: Event Log Decoding
	fmt.Println("\n📊 4. Event Log Decoding")
	demonstrateEventDecoding()

	// Example 5: Custom Error Decoding
	fmt.Println("\n⚠️  5. Custom Error Decoding")
	demonstrateErrorDecoding()

	// Example 6: Real Blockchain Integration (commented out - requires node)
	fmt.Println("\n🌐 6. Blockchain Integration Example")
	demonstrateBlockchainIntegration()
}

// Example 1: Access contract metadata and ABI
func demonstrateMetadata() {
	// Get the contract ABI as JSON string
	abi := simpletoken.ABI()
	fmt.Printf("✅ Contract ABI loaded (%d chars)\n", len(abi))

	// Access bytecode for deployment
	bytecode := simpletoken.Bytecode.Hex()
	fmt.Printf("✅ Deployment bytecode ready (%d chars)\n", len(bytecode))

	// Get runtime bytecode
	runtimeBytecode := simpletoken.DeployedBytecode.Hex()
	fmt.Printf("✅ Runtime bytecode available (%d chars)\n", len(runtimeBytecode))
}

// Example 2: Encode method calls for transactions
func demonstrateMethodEncoding() {
	// Create method instances using the clean chaining API
	transferMethod := simpletoken.Methods().TransferMethod()
	approveMethod := simpletoken.Methods().ApproveMethod()

	// Create an address and amount for our examples
	recipient := simpletoken.AddressFromHex("742d35Cc6634C0532925a3b8c0b56D39C3F6C842")
	amount := big.NewInt(1000000000000000000) // 1 ETH in wei

	// Encode a transfer method call
	transferCalldata, err := transferMethod.Pack(recipient, amount)
	if err != nil {
		log.Fatal("Failed to encode transfer:", err)
	}
	fmt.Printf("✅ Transfer call encoded: %s\n", transferCalldata.Hex())

	// Encode an approve method call  
	approveCalldata, err := approveMethod.Pack(recipient, amount)
	if err != nil {
		log.Fatal("Failed to encode approve:", err)
	}
	fmt.Printf("✅ Approve call encoded: %s\n", approveCalldata.Hex())

	// You can also use MustPack if you're confident about the inputs
	balanceCalldata := simpletoken.Methods().BalanceOfMethod().MustPack(recipient)
	fmt.Printf("✅ BalanceOf call encoded: %s\n", balanceCalldata.Hex())
}

// Example 3: Decode method return values
func demonstrateReturnDecoding() {
	// Simulate return data from blockchain calls
	// This would typically come from eth_call responses

	// Example: balanceOf returns uint256 -> *big.Int
	balanceReturn := "0x0000000000000000000000000000000000000000000000000de0b6b3a7640000" // 1 ETH
	balanceBytes := simpletoken.HexData(balanceReturn).Bytes()
	
	balance := simpletoken.Methods().BalanceOfMethod().MustDecode(balanceBytes)
	fmt.Printf("✅ Decoded balance: %s wei (%.2f ETH)\n", balance.String(), weiToEth(balance))

	// Example: transfer returns bool
	transferReturn := "0x0000000000000000000000000000000000000000000000000000000000000001" // true
	transferBytes := simpletoken.HexData(transferReturn).Bytes()
	
	success := simpletoken.Methods().TransferMethod().MustDecode(transferBytes)
	fmt.Printf("✅ Transfer success: %t\n", success)

	// Example: name() returns string
	nameReturn := "0x0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000b53696d706c65546f6b656e000000000000000000000000000000000000000000" // "SimpleToken"
	nameBytes := simpletoken.HexData(nameReturn).Bytes()
	
	tokenName := simpletoken.Methods().NameMethod().MustDecode(nameBytes)
	fmt.Printf("✅ Token name: %s\n", tokenName)
}

// Example 4: Decode event logs
func demonstrateEventDecoding() {
	// Simulate event log data from blockchain
	// Transfer event: Transfer(address indexed from, address indexed to, uint256 value)
	
	// The event data contains only non-indexed parameters (value)
	eventData := "0x0000000000000000000000000000000000000000000000000de0b6b3a7640000" // 1 ETH
	eventBytes := simpletoken.HexData(eventData).Bytes()

	// Decode the Transfer event
	transferEvent := simpletoken.Events().TransferEvent().MustDecode(eventBytes)
	fmt.Printf("✅ Transfer event decoded:\n")
	fmt.Printf("   Value: %s wei (%.2f ETH)\n", transferEvent.Value.String(), weiToEth(transferEvent.Value))

	// Note: indexed parameters (from, to) would come from log topics, not data
	// You would typically get these from: log.Topics[1] and log.Topics[2]
	
	// Approval event example
	approvalData := "0x00000000000000000000000000000000000000000000021e19e0c9bab2400000" // 10000 ETH
	approvalBytes := simpletoken.HexData(approvalData).Bytes()

	approvalEvent := simpletoken.Events().ApprovalEvent().MustDecode(approvalBytes)
	fmt.Printf("✅ Approval event decoded:\n")
	fmt.Printf("   Value: %s wei (%.2f ETH)\n", approvalEvent.Value.String(), weiToEth(approvalEvent.Value))
}

// Example 5: Decode custom errors
func demonstrateErrorDecoding() {
	// Simulate error data that would come from a reverted transaction
	// InsufficientBalance(address account, uint256 requested, uint256 available)
	
	errorData := "0xdb42144d" + // 4-byte error selector
		"000000000000000000000000742d35cc6634c0532925a3b8c0b56d39c3f6c842" + // account address
		"0000000000000000000000000000000000000000000000000de0b6b3a7640000" + // requested: 1 ETH
		"00000000000000000000000000000000000000000000000006f05b59d3b20000" // available: 0.5 ETH

	errorBytes := simpletoken.HexData(errorData).Bytes()

	// Decode the custom error
	insufficientError := simpletoken.Errors().InsufficientBalanceError().MustDecode(errorBytes)
	fmt.Printf("✅ InsufficientBalance error decoded:\n")
	fmt.Printf("   Account: %s\n", insufficientError.Account.String())
	fmt.Printf("   Requested: %s wei (%.2f ETH)\n", insufficientError.Requested.String(), weiToEth(insufficientError.Requested))
	fmt.Printf("   Available: %s wei (%.2f ETH)\n", insufficientError.Available.String(), weiToEth(insufficientError.Available))
}

// Example 6: Real blockchain integration using ethclient
func demonstrateBlockchainIntegration() {
	fmt.Println("// This example shows how to integrate with a real Ethereum node")
	fmt.Println("// Uncomment and modify the connection details to try it live")
	
	fmt.Println(`
// Connect to Ethereum node
client, err := ethclient.Dial("https://mainnet.infura.io/v3/YOUR_PROJECT_ID")
if err != nil {
    log.Fatal("Failed to connect:", err)
}
defer client.Close()

// Contract address (example)
contractAddr := common.HexToAddress("0x...")

// 1. Read contract state
balanceCalldata := simpletoken.Methods().BalanceOfMethod().MustPack(
    simpletoken.AddressFromHex("742d35Cc6634C0532925a3b8c0b56D39C3F6C842"),
)

result, err := client.CallContract(context.Background(), ethereum.CallMsg{
    To:   &contractAddr,
    Data: balanceCalldata.Bytes(),
}, nil)
if err != nil {
    log.Fatal("Call failed:", err)
}

balance := simpletoken.Methods().BalanceOfMethod().MustDecode(result)
fmt.Printf("Balance: %s ETH\n", weiToEth(balance))

// 2. Send transaction
privateKey, _ := crypto.GenerateKey() // Use your real key
auth, _ := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1))

transferCalldata := simpletoken.Methods().TransferMethod().MustPack(
    recipient, big.NewInt(1000000000000000000),
)

tx := types.NewTransaction(
    nonce, contractAddr, big.NewInt(0), gasLimit, gasPrice, transferCalldata.Bytes(),
)

signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
if err != nil {
    log.Fatal("Failed to sign:", err)
}

err = client.SendTransaction(context.Background(), signedTx)
if err != nil {
    log.Fatal("Failed to send:", err)
}

// 3. Listen for events
logs, err := client.FilterLogs(context.Background(), ethereum.FilterQuery{
    Addresses: []common.Address{contractAddr},
    Topics:    [][]common.Hash{{simpletoken.Events().TransferEvent().Topic}},
})

for _, log := range logs {
    transferEvent := simpletoken.Events().TransferEvent().MustDecode(log.Data)
    fmt.Printf("Transfer: %s ETH\n", weiToEth(transferEvent.Value))
}`)

	fmt.Println("\n✅ Integration patterns documented above")
}

// Utility function to convert wei to ETH
func weiToEth(wei *big.Int) float64 {
	eth := new(big.Float).SetInt(wei)
	eth.Quo(eth, big.NewFloat(1e18))
	result, _ := eth.Float64()
	return result
}