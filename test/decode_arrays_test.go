// SPDX-License-Identifier: MIT

package test

import (
	"bytes"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"

	dt "github.com/otherview/solgen/test/data/datatypes/packable/datatypespackable"
)

// packReturn encodes the named method's output values using go-ethereum, giving
// a reference ABI return blob to feed the generated decoder.
func packReturn(t *testing.T, method string, vals ...any) []byte {
	t.Helper()
	ref := referenceABI(t)
	data, err := ref.Methods[method].Outputs.Pack(vals...)
	if err != nil {
		t.Fatalf("packing %s outputs: %v", method, err)
	}
	return data
}

// TestPack_FixedArrays verifies fixed-size array arguments encode identically to
// go-ethereum (inline static words, no length prefix or offset pointer).
func TestPack_FixedArrays(t *testing.T) {
	ref := referenceABI(t)

	var h0, h1, h2, h3 [32]byte
	for i := range h0 {
		h0[i], h1[i], h2[i], h3[i] = byte(i), byte(i+1), byte(i+2), byte(i+3)
	}
	a0 := ethcommon.HexToAddress("0x1111111111111111111111111111111111111111")
	a1 := ethcommon.HexToAddress("0x2222222222222222222222222222222222222222")

	tests := []struct {
		name    string
		method  string
		gen     func() (dt.HexData, error)
		refArgs []any
	}{
		{
			name:   "uint256[3]",
			method: "echoUint256FixedArray",
			gen: func() (dt.HexData, error) {
				return dt.Methods().EchoUint256FixedArrayMethod().Pack([3]*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)})
			},
			refArgs: []any{[3]*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}},
		},
		{
			name:   "address[2]",
			method: "echoAddressFixedArray",
			gen: func() (dt.HexData, error) {
				return dt.Methods().EchoAddressFixedArrayMethod().Pack([2]dt.Address{dt.Address(a0), dt.Address(a1)})
			},
			refArgs: []any{[2]ethcommon.Address{a0, a1}},
		},
		{
			name:   "bytes32[4]",
			method: "echoBytes32FixedArray",
			gen: func() (dt.HexData, error) {
				return dt.Methods().EchoBytes32FixedArrayMethod().Pack([4][32]byte{h0, h1, h2, h3})
			},
			refArgs: []any{[4][32]byte{h0, h1, h2, h3}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.gen()
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

// TestDecode_FixedArrayReturn verifies a fixed-size array return value decodes.
func TestDecode_FixedArrayReturn(t *testing.T) {
	want := [3]*big.Int{big.NewInt(10), big.NewInt(20), big.NewInt(30)}
	data := packReturn(t, "echoUint256FixedArray", want)

	got, err := dt.Methods().EchoUint256FixedArrayMethod().Decode(data)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	for i := range want {
		if got[i].Cmp(want[i]) != 0 {
			t.Errorf("element %d: got %v, want %v", i, got[i], want[i])
		}
	}
}

// TestDecode_DynamicArrayReturn verifies dynamic array returns, including the
// signed int256[] case (which must use the signed element decoder).
func TestDecode_DynamicArrayReturn(t *testing.T) {
	t.Run("uint256[]", func(t *testing.T) {
		want := []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}
		data := packReturn(t, "echoUint256Array", want)
		got, err := dt.Methods().EchoUint256ArrayMethod().Decode(data)
		if err != nil {
			t.Fatalf("Decode: %v", err)
		}
		if len(got) != len(want) {
			t.Fatalf("len: got %d want %d", len(got), len(want))
		}
		for i := range want {
			if got[i].Cmp(want[i]) != 0 {
				t.Errorf("element %d: got %v want %v", i, got[i], want[i])
			}
		}
	})

	t.Run("int256[] negative", func(t *testing.T) {
		want := []*big.Int{big.NewInt(-100), big.NewInt(200), big.NewInt(-300)}
		data := packReturn(t, "echoInt256Array", want)
		got, err := dt.Methods().EchoInt256ArrayMethod().Decode(data)
		if err != nil {
			t.Fatalf("Decode: %v", err)
		}
		for i := range want {
			if got[i].Cmp(want[i]) != 0 {
				t.Errorf("element %d: got %v want %v (signed decode)", i, got[i], want[i])
			}
		}
	})
}

// TestDecode_MultiReturnWithString validates the multi-return decoder, including
// a dynamic string read via an offset pointer.
func TestDecode_MultiReturnWithString(t *testing.T) {
	addr := ethcommon.HexToAddress("0x1234567890123456789012345678901234567890")
	data := packReturn(t, "getBasicTypes", true, big.NewInt(12345), big.NewInt(-6789), addr, "hello world")

	got, err := dt.Methods().GetBasicTypesMethod().Decode(data)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if !got.BoolVal {
		t.Error("BoolVal: want true")
	}
	if got.UintVal.Cmp(big.NewInt(12345)) != 0 {
		t.Errorf("UintVal: got %v", got.UintVal)
	}
	if got.IntVal.Cmp(big.NewInt(-6789)) != 0 {
		t.Errorf("IntVal: got %v", got.IntVal)
	}
	if !bytes.Equal(got.AddressVal[:], addr.Bytes()) {
		t.Errorf("AddressVal: got %x want %x", got.AddressVal[:], addr.Bytes())
	}
	if got.StringVal != "hello world" {
		t.Errorf("StringVal: got %q", got.StringVal)
	}
}

// TestDecodeLog_DynamicAndArrayParams verifies the event decoder handles a
// non-indexed dynamic string and a non-indexed uint256[] (issue #2), with an
// indexed uint256 from topics.
func TestDecodeLog_DynamicAndArrayParams(t *testing.T) {
	ref := referenceABI(t)
	data, err := ref.Events["Logged"].Inputs.NonIndexed().Pack("hello", []*big.Int{big.NewInt(7), big.NewInt(8)})
	if err != nil {
		t.Fatalf("packing event data: %v", err)
	}
	var sig, idTopic [32]byte
	copy(sig[:], ref.Events["Logged"].ID.Bytes())
	big.NewInt(42).FillBytes(idTopic[:])

	got, err := dt.Events().LoggedEventDecoder().DecodeLog([][32]byte{sig, idTopic}, data)
	if err != nil {
		t.Fatalf("DecodeLog: %v", err)
	}
	if got.Id == nil || got.Id.Cmp(big.NewInt(42)) != 0 {
		t.Errorf("Id: got %v want 42", got.Id)
	}
	if got.Label != "hello" {
		t.Errorf("Label: got %q", got.Label)
	}
	if len(got.Values) != 2 || got.Values[0].Cmp(big.NewInt(7)) != 0 || got.Values[1].Cmp(big.NewInt(8)) != 0 {
		t.Errorf("Values: got %v", got.Values)
	}
}

// TestDecode_ErrorParams verifies the error decoder handles a dynamic string +
// uint256 (BadValue) and a uint256[] (BadList), reading dynamic data via offset
// pointers.
func TestDecode_ErrorParams(t *testing.T) {
	ref := referenceABI(t)

	t.Run("BadValue string+uint256", func(t *testing.T) {
		body, err := ref.Errors["BadValue"].Inputs.Pack("oops", big.NewInt(7))
		if err != nil {
			t.Fatalf("packing error data: %v", err)
		}
		data := append([]byte{0, 0, 0, 0}, body...) // 4-byte selector + params
		got, err := dt.Errors().BadValueError().Decode(data)
		if err != nil {
			t.Fatalf("Decode: %v", err)
		}
		if got.Reason != "oops" {
			t.Errorf("Reason: got %q", got.Reason)
		}
		if got.Code.Cmp(big.NewInt(7)) != 0 {
			t.Errorf("Code: got %v", got.Code)
		}
	})

	t.Run("BadList uint256[]", func(t *testing.T) {
		body, err := ref.Errors["BadList"].Inputs.Pack([]*big.Int{big.NewInt(1), big.NewInt(2)})
		if err != nil {
			t.Fatalf("packing error data: %v", err)
		}
		data := append([]byte{0, 0, 0, 0}, body...)
		got, err := dt.Errors().BadListError().Decode(data)
		if err != nil {
			t.Fatalf("Decode: %v", err)
		}
		if len(got.Values) != 2 || got.Values[0].Cmp(big.NewInt(1)) != 0 || got.Values[1].Cmp(big.NewInt(2)) != 0 {
			t.Errorf("Values: got %v", got.Values)
		}
	})
}

// TestDataTypesFixtureCompiles ensures the full DataTypesContract reference
// fixture (its own module) builds — it previously did not compile at all.
func TestDataTypesFixtureCompiles(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping compile test in short mode")
	}
	dir, err := filepath.Abs(filepath.Join("..", "test", "data", "datatypes", "generated", "datatypescontract"))
	if err != nil {
		t.Fatalf("resolving fixture dir: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "datatypescontract.go")); err != nil {
		t.Fatalf("fixture missing: %v", err)
	}
	if err := goBuildDir(t, dir); err != nil {
		t.Fatalf("datatypes fixture does not compile: %v", err)
	}
}
