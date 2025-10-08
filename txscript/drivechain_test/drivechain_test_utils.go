// Copyright (c) 2025 The mainchainlabs developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package drivechain_test

import (
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

// CreateTestDrivechainScript creates a test drivechain script with the given message type
func CreateTestDrivechainScript(msgType uint8) []byte {
	return []byte{0xE0, msgType} // OP_DRIVECHAIN + message type
}

// CreateTestDrivechainTransaction creates a test transaction with a drivechain output
func CreateTestDrivechainTransaction(msgType uint8, value int64) *btcutil.Tx {
	script := CreateTestDrivechainScript(msgType)

	txOut := &wire.TxOut{
		Value:    value,
		PkScript: script,
	}

	tx := &wire.MsgTx{
		Version: 1,
		TxOut:   []*wire.TxOut{txOut},
	}

	return btcutil.NewTx(tx)
}

// CreateTestDrivechainInput creates a test transaction input with drivechain script
func CreateTestDrivechainInput(msgType uint8) *wire.TxIn {
	script := CreateTestDrivechainScript(msgType)

	return &wire.TxIn{
		PreviousOutPoint: wire.OutPoint{
			Hash:  [32]byte{}, // Zero hash for test
			Index: 0,
		},
		SignatureScript: script,
		Sequence:        wire.MaxTxInSequenceNum,
	}
}

// ValidateDrivechainMessage validates that a parsed message matches expected values
func ValidateDrivechainMessage(msg *txscript.DrivechainMessage, expectedType uint8) error {
	if msg == nil {
		return fmt.Errorf("message is nil")
	}

	if msg.MessageType != expectedType {
		return fmt.Errorf("expected message type %d, got %d", expectedType, msg.MessageType)
	}

	return nil
}

// AllValidMessageTypes returns all valid drivechain message types for testing
func AllValidMessageTypes() []uint8 {
	return []uint8{
		txscript.MSG_SIDECHAIN_PROPOSAL,
		txscript.MSG_SIDECHAIN_ACK,
		txscript.MSG_BUNDLE_PROPOSAL,
		txscript.MSG_BUNDLE_ACK,
		txscript.MSG_DEPOSIT,
		txscript.MSG_WITHDRAWAL,
	}
}

// InvalidMessageTypes returns invalid message types for error testing
func InvalidMessageTypes() []uint8 {
	return []uint8{0, 7, 99, 255}
}
