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
		{"BMM Accept", MSG_BMM_ACCEPT, 7},
		{"BMM Request", MSG_BMM_REQUEST, 8},
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
		{"Valid: BMM Accept", MSG_BMM_ACCEPT, true},
		{"Valid: BMM Request", MSG_BMM_REQUEST, true},
		{"Invalid: Zero", 0, false},
		{"Invalid: Between valid range", 5, false},
		{"Invalid: Too high", 9, false},
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
		{"BMM Accept", MSG_BMM_ACCEPT, "BMM Accept (M7)"},
		{"BMM Request", MSG_BMM_REQUEST, "BMM Request (M8)"},
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

// TestParseDrivechainMessage tests the OP_RETURN message parsing functionality
func TestParseDrivechainMessage(t *testing.T) {
	// Create valid M1 message script
	m1 := &M1ProposeSidechain{
		SidechainNumber: SidechainNumber(1),
		Description:     SidechainDescription("test"),
	}
	m1Data, _ := m1.Serialize()
	m1Script, _ := NullDataScript(m1Data)

	// Create valid M2 message script
	m2 := &M2AckSidechain{
		SidechainNumber: SidechainNumber(1),
		DescriptionHash: [32]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20},
	}
	m2Data, _ := m2.Serialize()
	m2Script, _ := NullDataScript(m2Data)

	tests := []struct {
		name        string
		script      []byte
		expectError bool
		expectedType uint8
	}{
		{
			name:        "Valid M1 Sidechain Proposal",
			script:      m1Script,
			expectError: false,
			expectedType: MSG_SIDECHAIN_PROPOSAL,
		},
		{
			name:        "Valid M2 Sidechain Acknowledgment",
			script:      m2Script,
			expectError: false,
			expectedType: MSG_SIDECHAIN_ACK,
		},
		{
			name:        "Script too short",
			script:      []byte{OP_RETURN},
			expectError: true,
			expectedType: 0,
		},
		{
			name:        "Empty script",
			script:      []byte{},
			expectError: true,
			expectedType: 0,
		},
		{
			name:        "Not OP_RETURN",
			script:      []byte{0x01, 0x02, 0x03},
			expectError: true,
			expectedType: 0,
		},
		{
			name:        "Invalid message tag",
			script:      []byte{OP_RETURN, 0x04, 0x00, 0x00, 0x00, 0x00},
			expectError: true,
			expectedType: 0,
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
				} else if msg.MessageType() != tt.expectedType {
					t.Errorf("Expected message type %d, got %d", tt.expectedType, msg.MessageType())
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
	m1 := &M1ProposeSidechain{
		SidechainNumber: SidechainNumber(1),
		Description:     SidechainDescription("test"),
	}
	m1Data, _ := m1.Serialize()
	script, _ := NullDataScript(m1Data)

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
