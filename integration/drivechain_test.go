// Copyright (c) 2025 The mainchainlabs developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package integration

import (
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/integration/rpctest"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

// TestDrivechainBasicIntegration tests basic drivechain functionality in regtest
func TestDrivechainBasicIntegration(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a regtest harness
	harness, err := rpctest.New(&chaincfg.RegressionNetParams, nil, nil, "")
	if err != nil {
		t.Fatalf("Unable to create test harness: %v", err)
	}
	defer harness.TearDown()

	// Start the harness
	if err := harness.SetUp(false, 0); err != nil {
		t.Fatalf("Unable to setup test harness: %v", err)
	}

	// Generate some blocks to have funds
	_, err = harness.Client.Generate(101)
	if err != nil {
		t.Fatalf("Unable to generate blocks: %v", err)
	}

	// Test that we can create a basic drivechain message
	t.Run("CreateDrivechainMessage", func(t *testing.T) {
		script := []byte{0xE0, txscript.MSG_SIDECHAIN_PROPOSAL} // OP_DRIVECHAIN + M1

		msg, err := txscript.ParseDrivechainMessage(script)
		if err != nil {
			t.Fatalf("Failed to parse drivechain message: %v", err)
		}

		if msg.MessageType != txscript.MSG_SIDECHAIN_PROPOSAL {
			t.Errorf("Expected message type %d, got %d", txscript.MSG_SIDECHAIN_PROPOSAL, msg.MessageType)
		}
	})

	// Test that we can create a transaction with drivechain script
	t.Run("CreateDrivechainTransaction", func(t *testing.T) {
		// Create a simple drivechain script
		script := []byte{0xE0, txscript.MSG_SIDECHAIN_PROPOSAL}

		// Create a transaction output with this script
		txOut := &wire.TxOut{
			Value:    1000, // 1000 satoshis
			PkScript: script,
		}

		// Create a transaction
		tx := &wire.MsgTx{
			Version: 1,
			TxOut:   []*wire.TxOut{txOut},
		}

		// Convert to btcutil.Tx
		btcTx := btcutil.NewTx(tx)

		// For now, we just verify the transaction structure is valid
		// In future steps, we'll test actual submission to mempool
		if btcTx == nil {
			t.Fatal("Failed to create btcutil.Tx")
		}

		t.Logf("Created drivechain transaction: %s", btcTx.Hash())
	})
}

// TestDrivechainMessageTypes tests all message types in integration context
func TestDrivechainMessageTypes(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	messageTypes := []struct {
		name     string
		msgType  uint8
		expected string
	}{
		{"Sidechain Proposal", txscript.MSG_SIDECHAIN_PROPOSAL, "Sidechain Proposal (M1)"},
		{"Sidechain Acknowledgment", txscript.MSG_SIDECHAIN_ACK, "Sidechain Acknowledgment (M2)"},
		{"Bundle Proposal", txscript.MSG_BUNDLE_PROPOSAL, "Bundle Proposal (M3)"},
		{"Bundle Acknowledgment", txscript.MSG_BUNDLE_ACK, "Bundle Acknowledgment (M4)"},
		{"Deposit Transaction", txscript.MSG_DEPOSIT, "Deposit Transaction (M5)"},
		{"Withdrawal Transaction", txscript.MSG_WITHDRAWAL, "Withdrawal Transaction (M6)"},
	}

	for _, mt := range messageTypes {
		t.Run(mt.name, func(t *testing.T) {
			// Test message type validation
			if !txscript.IsValidMessageType(mt.msgType) {
				t.Errorf("Message type %d should be valid", mt.msgType)
			}

			// Test message type name
			name := txscript.GetMessageTypeName(mt.msgType)
			if name != mt.expected {
				t.Errorf("Expected name %s, got %s", mt.expected, name)
			}

			// Test parsing
			script := []byte{0xE0, mt.msgType}
			msg, err := txscript.ParseDrivechainMessage(script)
			if err != nil {
				t.Errorf("Failed to parse message type %d: %v", mt.msgType, err)
			}

			if msg.MessageType != mt.msgType {
				t.Errorf("Expected message type %d, got %d", mt.msgType, msg.MessageType)
			}
		})
	}
}

// TestDrivechainErrorHandling tests error conditions in integration context
func TestDrivechainErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	errorTests := []struct {
		name        string
		script      []byte
		expectError bool
	}{
		{"Empty script", []byte{}, true},
		{"Too short", []byte{0xE0}, true},
		{"Wrong opcode", []byte{0x01, txscript.MSG_SIDECHAIN_PROPOSAL}, true},
		{"Invalid message type", []byte{0xE0, 99}, true},
		{"Message type too high", []byte{0xE0, 255}, true},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := txscript.ParseDrivechainMessage(tt.script)

			if tt.expectError && err == nil {
				t.Errorf("Expected error for %s but got none", tt.name)
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for %s: %v", tt.name, err)
			}
		})
	}
}

// BenchmarkDrivechainParsing benchmarks drivechain message parsing
func BenchmarkDrivechainParsing(b *testing.B) {
	script := []byte{0xE0, txscript.MSG_SIDECHAIN_PROPOSAL}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = txscript.ParseDrivechainMessage(script)
	}
}
