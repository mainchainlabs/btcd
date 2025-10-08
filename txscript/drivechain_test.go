// Copyright (c) 2025 The mainchainlabs developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package txscript

import (
	"testing"
)

// TestDrivechainMessageTypes tests that all message type constants are defined correctly
func TestDrivechainMessageTypes(t *testing.T) {
	tests := []struct {
		name     string
		msgType  uint8
		expected uint8
	}{
		{"Sidechain Proposal", MSG_SIDECHAIN_PROPOSAL, 1},
		{"Sidechain Acknowledgment", MSG_SIDECHAIN_ACK, 2},
		{"Bundle Proposal", MSG_BUNDLE_PROPOSAL, 3},
		{"Bundle Acknowledgment", MSG_BUNDLE_ACK, 4},
		{"Deposit Transaction", MSG_DEPOSIT, 5},
		{"Withdrawal Transaction", MSG_WITHDRAWAL, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.msgType != tt.expected {
				t.Errorf("Expected %s to be %d, got %d", tt.name, tt.expected, tt.msgType)
			}
		})
	}
}

// TestIsValidMessageType tests the message type validation function
func TestIsValidMessageType(t *testing.T) {
	tests := []struct {
		name     string
		msgType  uint8
		expected bool
	}{
		{"Valid: Sidechain Proposal", MSG_SIDECHAIN_PROPOSAL, true},
		{"Valid: Sidechain Acknowledgment", MSG_SIDECHAIN_ACK, true},
		{"Valid: Bundle Proposal", MSG_BUNDLE_PROPOSAL, true},
		{"Valid: Bundle Acknowledgment", MSG_BUNDLE_ACK, true},
		{"Valid: Deposit Transaction", MSG_DEPOSIT, true},
		{"Valid: Withdrawal Transaction", MSG_WITHDRAWAL, true},
		{"Invalid: Zero", 0, false},
		{"Invalid: Too high", 7, false},
		{"Invalid: Very high", 255, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidMessageType(tt.msgType)
			if result != tt.expected {
				t.Errorf("IsValidMessageType(%d) = %v, expected %v", tt.msgType, result, tt.expected)
			}
		})
	}
}

// TestGetMessageTypeName tests the message type name function
func TestGetMessageTypeName(t *testing.T) {
	tests := []struct {
		name     string
		msgType  uint8
		expected string
	}{
		{"Sidechain Proposal", MSG_SIDECHAIN_PROPOSAL, "Sidechain Proposal (M1)"},
		{"Sidechain Acknowledgment", MSG_SIDECHAIN_ACK, "Sidechain Acknowledgment (M2)"},
		{"Bundle Proposal", MSG_BUNDLE_PROPOSAL, "Bundle Proposal (M3)"},
		{"Bundle Acknowledgment", MSG_BUNDLE_ACK, "Bundle Acknowledgment (M4)"},
		{"Deposit Transaction", MSG_DEPOSIT, "Deposit Transaction (M5)"},
		{"Withdrawal Transaction", MSG_WITHDRAWAL, "Withdrawal Transaction (M6)"},
		{"Unknown Type", 99, "Unknown Message Type"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetMessageTypeName(tt.msgType)
			if result != tt.expected {
				t.Errorf("GetMessageTypeName(%d) = %s, expected %s", tt.msgType, result, tt.expected)
			}
		})
	}
}

// TestParseDrivechainMessage tests the basic message parsing functionality
func TestParseDrivechainMessage(t *testing.T) {
	tests := []struct {
		name        string
		script      []byte
		expectError bool
		expectedMsg *DrivechainMessage
	}{
		{
			name:        "Valid Sidechain Proposal",
			script:      []byte{0xE0, MSG_SIDECHAIN_PROPOSAL}, // OP_DRIVECHAIN + M1
			expectError: false,
			expectedMsg: &DrivechainMessage{MessageType: MSG_SIDECHAIN_PROPOSAL},
		},
		{
			name:        "Valid Sidechain Acknowledgment",
			script:      []byte{0xE0, MSG_SIDECHAIN_ACK}, // OP_DRIVECHAIN + M2
			expectError: false,
			expectedMsg: &DrivechainMessage{MessageType: MSG_SIDECHAIN_ACK},
		},
		{
			name:        "Valid Bundle Proposal",
			script:      []byte{0xE0, MSG_BUNDLE_PROPOSAL}, // OP_DRIVECHAIN + M3
			expectError: false,
			expectedMsg: &DrivechainMessage{MessageType: MSG_BUNDLE_PROPOSAL},
		},
		{
			name:        "Valid Bundle Acknowledgment",
			script:      []byte{0xE0, MSG_BUNDLE_ACK}, // OP_DRIVECHAIN + M4
			expectError: false,
			expectedMsg: &DrivechainMessage{MessageType: MSG_BUNDLE_ACK},
		},
		{
			name:        "Valid Deposit Transaction",
			script:      []byte{0xE0, MSG_DEPOSIT}, // OP_DRIVECHAIN + M5
			expectError: false,
			expectedMsg: &DrivechainMessage{MessageType: MSG_DEPOSIT},
		},
		{
			name:        "Valid Withdrawal Transaction",
			script:      []byte{0xE0, MSG_WITHDRAWAL}, // OP_DRIVECHAIN + M6
			expectError: false,
			expectedMsg: &DrivechainMessage{MessageType: MSG_WITHDRAWAL},
		},
		{
			name:        "Script too short",
			script:      []byte{0xE0},
			expectError: true,
			expectedMsg: nil,
		},
		{
			name:        "Empty script",
			script:      []byte{},
			expectError: true,
			expectedMsg: nil,
		},
		{
			name:        "Wrong opcode",
			script:      []byte{0x01, MSG_SIDECHAIN_PROPOSAL},
			expectError: true,
			expectedMsg: nil,
		},
		{
			name:        "Invalid message type",
			script:      []byte{0xE0, 99},
			expectError: true,
			expectedMsg: nil,
		},
		{
			name:        "Message type too high",
			script:      []byte{0xE0, 255},
			expectError: true,
			expectedMsg: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := ParseDrivechainMessage(tt.script)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if msg != nil {
					t.Errorf("Expected nil message but got %+v", msg)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if msg == nil {
					t.Errorf("Expected message but got nil")
				} else if msg.MessageType != tt.expectedMsg.MessageType {
					t.Errorf("Expected message type %d, got %d", tt.expectedMsg.MessageType, msg.MessageType)
				}
			}
		})
	}
}

// TestDrivechainMessageStructure tests the basic structure of DrivechainMessage
func TestDrivechainMessageStructure(t *testing.T) {
	msg := &DrivechainMessage{
		MessageType: MSG_SIDECHAIN_PROPOSAL,
	}

	if msg.MessageType != MSG_SIDECHAIN_PROPOSAL {
		t.Errorf("Expected MessageType to be %d, got %d", MSG_SIDECHAIN_PROPOSAL, msg.MessageType)
	}
}

// Benchmark tests for performance
func BenchmarkParseDrivechainMessage(b *testing.B) {
	script := []byte{0xE0, MSG_SIDECHAIN_PROPOSAL}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseDrivechainMessage(script)
	}
}

func BenchmarkIsValidMessageType(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = IsValidMessageType(MSG_SIDECHAIN_PROPOSAL)
	}
}
