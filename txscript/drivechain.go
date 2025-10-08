// Copyright (c) 2025 The mainchainlabs developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package txscript

import (
	"errors"
)

// Drivechain message types as defined in BIP300/BIP301
const (
	MSG_SIDECHAIN_PROPOSAL = 1 // M1 - Sidechain proposal
	MSG_SIDECHAIN_ACK      = 2 // M2 - Sidechain acknowledgment
	MSG_BUNDLE_PROPOSAL    = 3 // M3 - Bundle proposal
	MSG_BUNDLE_ACK         = 4 // M4 - Bundle acknowledgment
	MSG_DEPOSIT            = 5 // M5 - Deposit transaction
	MSG_WITHDRAWAL         = 6 // M6 - Withdrawal transaction
)

// DrivechainMessage represents a basic drivechain message
// This is a minimal structure that will be expanded as we implement
// the full drivechain functionality
type DrivechainMessage struct {
	MessageType uint8
	// Additional fields will be added in future steps
}

// ParseDrivechainMessage parses a drivechain message from script bytes
// This is a basic implementation that will be expanded
func ParseDrivechainMessage(script []byte) (*DrivechainMessage, error) {
	// Basic validation - script must be at least 2 bytes
	if len(script) < 2 {
		return nil, errors.New("script too short for drivechain message")
	}

	// First byte should be OP_DRIVECHAIN (will be defined in next step)
	// For now, we'll use a placeholder
	if script[0] != 0xE0 { // OP_DRIVECHAIN placeholder
		return nil, errors.New("not a drivechain script")
	}

	// Second byte should be a valid message type
	msgType := script[1]
	if msgType < MSG_SIDECHAIN_PROPOSAL || msgType > MSG_WITHDRAWAL {
		return nil, errors.New("invalid drivechain message type")
	}

	msg := &DrivechainMessage{
		MessageType: msgType,
	}

	return msg, nil
}

// IsValidMessageType checks if a message type is valid
func IsValidMessageType(msgType uint8) bool {
	return msgType >= MSG_SIDECHAIN_PROPOSAL && msgType <= MSG_WITHDRAWAL
}

// GetMessageTypeName returns a human-readable name for the message type
func GetMessageTypeName(msgType uint8) string {
	switch msgType {
	case MSG_SIDECHAIN_PROPOSAL:
		return "Sidechain Proposal (M1)"
	case MSG_SIDECHAIN_ACK:
		return "Sidechain Acknowledgment (M2)"
	case MSG_BUNDLE_PROPOSAL:
		return "Bundle Proposal (M3)"
	case MSG_BUNDLE_ACK:
		return "Bundle Acknowledgment (M4)"
	case MSG_DEPOSIT:
		return "Deposit Transaction (M5)"
	case MSG_WITHDRAWAL:
		return "Withdrawal Transaction (M6)"
	default:
		return "Unknown Message Type"
	}
}
