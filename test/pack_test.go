// SPDX-License-Identifier: MIT

package test

import (
	"bytes"
	"math/big"
	"strings"
	"testing"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"

	dt "github.com/otherview/solgen/test/data/datatypes/packable/datatypespackable"
)

// referenceABI parses the generated fixture's ABI with go-ethereum so we can
// cross-check solgen's calldata against a reference encoder, byte-for-byte.
func referenceABI(t *testing.T) ethabi.ABI {
	t.Helper()
	parsed, err := ethabi.JSON(strings.NewReader(dt.ABI()))
	if err != nil {
		t.Fatalf("parsing reference ABI: %v", err)
	}
	return parsed
}

// TestPack_FixedBytes verifies that bytesN arguments encode as a single static
// 32-byte word (left-aligned, right-padded), matching go-ethereum's encoder.
func TestPack_FixedBytes(t *testing.T) {
	ref := referenceABI(t)

	var b32 [32]byte
	for i := range b32 {
		b32[i] = byte(i + 1)
	}
	var asHash dt.Hash = b32

	tests := []struct {
		name     string
		method   string
		generate func() (dt.HexData, error)
		refArgs  []any
	}{
		{
			name:     "bytes32 as [32]byte",
			method:   "echoBytes32",
			generate: func() (dt.HexData, error) { return dt.Methods().EchoBytes32Method().Pack(b32) },
			refArgs:  []any{b32},
		},
		{
			name:     "bytes32 as Hash",
			method:   "echoBytes32",
			generate: func() (dt.HexData, error) { return dt.Methods().EchoBytes32Method().Pack(asHash) },
			refArgs:  []any{b32},
		},
		{
			name:   "bytes4 right-pads within word",
			method: "echoBytes4",
			generate: func() (dt.HexData, error) {
				return dt.Methods().EchoBytes4Method().Pack([4]byte{0xde, 0xad, 0xbe, 0xef})
			},
			refArgs: []any{[4]byte{0xde, 0xad, 0xbe, 0xef}},
		},
		{
			name:     "bytes1",
			method:   "echoBytes1",
			generate: func() (dt.HexData, error) { return dt.Methods().EchoBytes1Method().Pack([1]byte{0xab}) },
			refArgs:  []any{[1]byte{0xab}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.generate()
			if err != nil {
				t.Fatalf("Pack: %v", err)
			}
			want, err := ref.Pack(tc.method, tc.refArgs...)
			if err != nil {
				t.Fatalf("reference Pack: %v", err)
			}
			if !bytes.Equal(got.Bytes(), want) {
				t.Fatalf("calldata mismatch\n got: %x\nwant: %x", got.Bytes(), want)
			}
		})
	}
}

// TestPack_FixedBytesHardcoded checks the exact byte layout against hand-computed
// vectors, independent of any reference encoder.
func TestPack_FixedBytesHardcoded(t *testing.T) {
	ref := referenceABI(t)
	selBytes4 := ref.Methods["echoBytes4"].ID // 4-byte selector

	// bytes4(0xdeadbeef): value left-aligned, right-padded to a full 32-byte word.
	want := append([]byte{}, selBytes4...)
	word := make([]byte, 32)
	copy(word, []byte{0xde, 0xad, 0xbe, 0xef})
	want = append(want, word...)

	got, err := dt.Methods().EchoBytes4Method().Pack([4]byte{0xde, 0xad, 0xbe, 0xef})
	if err != nil {
		t.Fatalf("Pack: %v", err)
	}
	if !bytes.Equal(got.Bytes(), want) {
		t.Fatalf("calldata mismatch\n got: %x\nwant: %x", got.Bytes(), want)
	}
}

// TestPack_MixedStaticArgs verifies a method taking multiple static bytes32 args
// (resolve(bytes32,bool,bytes32)) packs to three inline 32-byte words after the
// selector, with no offset pointers.
func TestPack_MixedStaticArgs(t *testing.T) {
	ref := referenceABI(t)

	var id, evidence [32]byte
	for i := range id {
		id[i] = byte(0xa0 + i)
		evidence[i] = byte(0xb0 + i)
	}

	got, err := dt.Methods().ResolveMethod().Pack(dt.Hash(id), true, dt.Hash(evidence))
	if err != nil {
		t.Fatalf("Pack: %v", err)
	}
	want, err := ref.Pack("resolve", id, true, evidence)
	if err != nil {
		t.Fatalf("reference Pack: %v", err)
	}
	if !bytes.Equal(got.Bytes(), want) {
		t.Fatalf("calldata mismatch\n got: %x\nwant: %x", got.Bytes(), want)
	}

	// Head layout: 4-byte selector + exactly three 32-byte static words.
	if len(got.Bytes()) != 4+3*32 {
		t.Fatalf("expected selector + 3 static words (%d bytes), got %d", 4+3*32, len(got.Bytes()))
	}
}

// TestPack_DynamicBytesRegression ensures []byte still encodes as dynamic bytes
// (offset pointer + length + padded data), unchanged by the fixed-bytes fix.
func TestPack_DynamicBytesRegression(t *testing.T) {
	ref := referenceABI(t)

	payload := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	got, err := dt.Methods().EchoBytesMethod().Pack(payload)
	if err != nil {
		t.Fatalf("Pack: %v", err)
	}
	want, err := ref.Pack("echoBytes", payload)
	if err != nil {
		t.Fatalf("reference Pack: %v", err)
	}
	if !bytes.Equal(got.Bytes(), want) {
		t.Fatalf("calldata mismatch\n got: %x\nwant: %x", got.Bytes(), want)
	}

	// Sanity: head must hold a 32-byte offset pointer equal to 0x20 (32).
	body := got.Bytes()[4:]
	offset := new(big.Int).SetBytes(body[:32])
	if offset.Cmp(big.NewInt(32)) != 0 {
		t.Fatalf("expected dynamic offset pointer 32, got %s", offset)
	}
}

// TestDecodeLog_IndexedParams verifies that DecodeLog fills indexed fields
// (including an indexed bytes32) from topics and non-indexed fields from data,
// in ABI order. Event: EscrowFunded(bytes32 indexed agreementId,
// address indexed payer, address payee, uint256 amount).
func TestDecodeLog_IndexedParams(t *testing.T) {
	ref := referenceABI(t)

	// Non-indexed params (payee address, amount uint256) — both static.
	payee := ethcommon.HexToAddress("0x1234567890123456789012345678901234567890")
	amount := big.NewInt(1_000_000)
	data, err := ref.Events["EscrowFunded"].Inputs.NonIndexed().Pack(payee, amount)
	if err != nil {
		t.Fatalf("packing event data: %v", err)
	}

	// Indexed topics: [sig, agreementId (bytes32), payer (address)].
	var sig, idTopic, payerTopic [32]byte
	copy(sig[:], ref.Events["EscrowFunded"].ID.Bytes())
	for i := range idTopic {
		idTopic[i] = byte(0xc0 + i)
	}
	payer := ethcommon.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	copy(payerTopic[12:], payer.Bytes()) // addresses are right-aligned in a word
	topics := [][32]byte{sig, idTopic, payerTopic}

	result, err := dt.Events().EscrowFundedEventDecoder().DecodeLog(topics, data)
	if err != nil {
		t.Fatalf("DecodeLog: %v", err)
	}

	if !bytes.Equal(result.AgreementId[:], idTopic[:]) {
		t.Errorf("AgreementId: got %x, want %x", result.AgreementId[:], idTopic[:])
	}
	if !bytes.Equal(result.Payer[:], payer.Bytes()) {
		t.Errorf("Payer: got %x, want %x", result.Payer[:], payer.Bytes())
	}
	if !bytes.Equal(result.Payee[:], payee.Bytes()) {
		t.Errorf("Payee: got %x, want %x", result.Payee[:], payee.Bytes())
	}
	if result.Amount == nil || result.Amount.Cmp(amount) != 0 {
		t.Errorf("Amount: got %v, want %v", result.Amount, amount)
	}

	// A log missing the indexed topics must error rather than panic.
	if _, err := dt.Events().EscrowFundedEventDecoder().DecodeLog([][32]byte{sig}, data); err == nil {
		t.Error("expected error for missing indexed topics, got nil")
	}
}
