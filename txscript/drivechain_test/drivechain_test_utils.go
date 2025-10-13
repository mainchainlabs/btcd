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

// CreateTestM1ProposeSidechain creates a test M1 sidechain proposal
func CreateTestM1ProposeSidechain() *txscript.M1ProposeSidechain {
	return &txscript.M1ProposeSidechain{
		SidechainNumber: txscript.SidechainNumber(1),
		Description:     txscript.SidechainDescription("Test sidechain description data"),
	}
}

// CreateTestM2AckSidechain creates a test M2 sidechain acknowledgment
func CreateTestM2AckSidechain() *txscript.M2AckSidechain {
	return &txscript.M2AckSidechain{
		SidechainNumber: txscript.SidechainNumber(1),
		DescriptionHash: [32]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20},
	}
}

// CreateTestM3ProposeBundle creates a test M3 bundle proposal
func CreateTestM3ProposeBundle() *txscript.M3ProposeBundle {
	return &txscript.M3ProposeBundle{
		SidechainNumber: txscript.SidechainNumber(1),
		BundleTxid:      [32]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20},
	}
}

// CreateTestM4AckBundles creates a test M4 bundle acknowledgment
func CreateTestM4AckBundles() *txscript.M4AckBundles {
	return &txscript.M4AckBundles{
		Variant: txscript.M4_ONE_BYTE,
		Upvotes: []byte{0x01, 0x02, 0x03, 0x04},
	}
}

// CreateTestM7BmmAccept creates a test M7 BMM accept
func CreateTestM7BmmAccept() *txscript.M7BmmAccept {
	return &txscript.M7BmmAccept{
		SidechainNumber:    txscript.SidechainNumber(1),
		SidechainBlockHash: txscript.BmmCommitment{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20},
	}
}

// CreateTestM8BmmRequest creates a test M8 BMM request
func CreateTestM8BmmRequest() *txscript.M8BmmRequest {
	return &txscript.M8BmmRequest{
		SidechainNumber:        txscript.SidechainNumber(1),
		SidechainBlockHash:     txscript.BmmCommitment{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20},
		PrevMainchainBlockHash: [32]byte{0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f, 0x40},
	}
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
		txscript.MSG_BMM_ACCEPT,
		txscript.MSG_BMM_REQUEST,
	}
}

// InvalidMessageTypes returns invalid message types for error testing
func InvalidMessageTypes() []uint8 {
	return []uint8{0, 7, 99, 255}
}
