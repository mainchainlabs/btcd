// Copyright (c) 2025 The mainchainlabs developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package txscript

import (
	"testing"
)

// TestM1ProposeSidechainSerialization tests serialization and deserialization of M1ProposeSidechain
func TestM1ProposeSidechainSerialization(t *testing.T) {
	original := &M1ProposeSidechain{
		SidechainNumber: SidechainNumber(1),
		Description:     SidechainDescription("TestSidechain description data"),
	}

	// Test serialization
	data, err := original.Serialize()
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	// Test deserialization
	deserialized, err := DeserializeM1ProposeSidechain(data)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	// Verify all fields match
	if original.SidechainNumber != deserialized.SidechainNumber {
		t.Errorf("SidechainNumber mismatch: got %d, want %d", deserialized.SidechainNumber, original.SidechainNumber)
	}
	if string(original.Description) != string(deserialized.Description) {
		t.Errorf("Description mismatch: got %s, want %s", string(deserialized.Description), string(original.Description))
	}
}

// TestM2AckSidechainSerialization tests serialization and deserialization of M2AckSidechain
func TestM2AckSidechainSerialization(t *testing.T) {
	original := &M2AckSidechain{
		SidechainNumber: SidechainNumber(1),
		DescriptionHash: [32]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20},
	}

	// Test serialization
	data, err := original.Serialize()
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	// Test deserialization
	deserialized, err := DeserializeM2AckSidechain(data)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	// Verify all fields match
	if original.SidechainNumber != deserialized.SidechainNumber {
		t.Errorf("SidechainNumber mismatch: got %d, want %d", deserialized.SidechainNumber, original.SidechainNumber)
	}
	if original.DescriptionHash != deserialized.DescriptionHash {
		t.Errorf("DescriptionHash mismatch: got %x, want %x", deserialized.DescriptionHash, original.DescriptionHash)
	}
}

// TestM3ProposeBundleSerialization tests serialization and deserialization of M3ProposeBundle
func TestM3ProposeBundleSerialization(t *testing.T) {
	original := &M3ProposeBundle{
		SidechainNumber: SidechainNumber(1),
		BundleTxid:      [32]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20},
	}

	// Test serialization
	data, err := original.Serialize()
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	// Test deserialization
	deserialized, err := DeserializeM3ProposeBundle(data)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	// Verify all fields match
	if original.SidechainNumber != deserialized.SidechainNumber {
		t.Errorf("SidechainNumber mismatch: got %d, want %d", deserialized.SidechainNumber, original.SidechainNumber)
	}
	if original.BundleTxid != deserialized.BundleTxid {
		t.Errorf("BundleTxid mismatch: got %x, want %x", deserialized.BundleTxid, original.BundleTxid)
	}
}

// TestM4AckBundlesSerialization tests serialization and deserialization of M4AckBundles
func TestM4AckBundlesSerialization(t *testing.T) {
	original := &M4AckBundles{
		Variant: M4_ONE_BYTE,
		Upvotes: []byte{0x01, 0x02, 0x03, 0x04},
	}

	// Test serialization
	data, err := original.Serialize()
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	// Verify serialized data contains expected tag and variant
	expectedTag := M4_TAG[:]
	if len(data) < len(expectedTag)+1 {
		t.Fatalf("Serialized data too short")
	}

	for i, b := range expectedTag {
		if data[i] != b {
			t.Errorf("Tag mismatch at position %d: got %x, want %x", i, data[i], b)
		}
	}

	if data[len(expectedTag)] != byte(original.Variant) {
		t.Errorf("Variant mismatch: got %x, want %x", data[len(expectedTag)], byte(original.Variant))
	}
}

// TestM7BmmAcceptSerialization tests serialization and deserialization of M7BmmAccept
func TestM7BmmAcceptSerialization(t *testing.T) {
	original := &M7BmmAccept{
		SidechainNumber:    SidechainNumber(1),
		SidechainBlockHash: BmmCommitment{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20},
	}

	// Test serialization
	data, err := original.Serialize()
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	// Test deserialization
	deserialized, err := DeserializeM7BmmAccept(data)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	// Verify all fields match
	if original.SidechainNumber != deserialized.SidechainNumber {
		t.Errorf("SidechainNumber mismatch: got %d, want %d", deserialized.SidechainNumber, original.SidechainNumber)
	}
	if original.SidechainBlockHash != deserialized.SidechainBlockHash {
		t.Errorf("SidechainBlockHash mismatch: got %x, want %x", deserialized.SidechainBlockHash, original.SidechainBlockHash)
	}
}

// TestM8BmmRequestSerialization tests serialization and deserialization of M8BmmRequest
func TestM8BmmRequestSerialization(t *testing.T) {
	original := &M8BmmRequest{
		SidechainNumber:        SidechainNumber(1),
		SidechainBlockHash:     BmmCommitment{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20},
		PrevMainchainBlockHash: [32]byte{0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f, 0x40},
	}

	// Test serialization
	data, err := original.Serialize()
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	// Test deserialization
	deserialized, err := DeserializeM8BmmRequest(data)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	// Verify all fields match
	if original.SidechainNumber != deserialized.SidechainNumber {
		t.Errorf("SidechainNumber mismatch: got %d, want %d", deserialized.SidechainNumber, original.SidechainNumber)
	}
	if original.SidechainBlockHash != deserialized.SidechainBlockHash {
		t.Errorf("SidechainBlockHash mismatch: got %x, want %x", deserialized.SidechainBlockHash, original.SidechainBlockHash)
	}
	if original.PrevMainchainBlockHash != deserialized.PrevMainchainBlockHash {
		t.Errorf("PrevMainchainBlockHash mismatch: got %x, want %x", deserialized.PrevMainchainBlockHash, original.PrevMainchainBlockHash)
	}
}

// TestM1ProposeSidechainErrorHandling tests error conditions for M1ProposeSidechain
func TestM1ProposeSidechainErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		expectError bool
	}{
		{
			name:        "Empty data",
			data:        []byte{},
			expectError: true,
		},
		{
			name:        "Wrong tag",
			data:        []byte{0x00, 0x00, 0x00, 0x00, 1},
			expectError: true,
		},
		{
			name:        "Incomplete data",
			data:        []byte{M1_TAG[0], M1_TAG[1], M1_TAG[2]}, // Missing last tag byte
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DeserializeM1ProposeSidechain(tt.data)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// TestM2AckSidechainErrorHandling tests error conditions for M2AckSidechain
func TestM2AckSidechainErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		expectError bool
	}{
		{
			name:        "Empty data",
			data:        []byte{},
			expectError: true,
		},
		{
			name:        "Wrong tag",
			data:        []byte{0x00, 0x00, 0x00, 0x00, 1, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20},
			expectError: true,
		},
		{
			name:        "Incomplete data",
			data:        []byte{M2_TAG[0], M2_TAG[1], M2_TAG[2], M2_TAG[3], 1}, // Missing hash
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DeserializeM2AckSidechain(tt.data)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// TestDrivechainStructures tests the basic structure definitions
func TestDrivechainStructures(t *testing.T) {
	// Test M1ProposeSidechain structure
	m1 := &M1ProposeSidechain{
		SidechainNumber: SidechainNumber(1),
		Description:     SidechainDescription("Test sidechain description"),
	}

	if m1.SidechainNumber != SidechainNumber(1) {
		t.Errorf("Expected SidechainNumber to be 1, got %d", m1.SidechainNumber)
	}

	// Test M2AckSidechain structure
	m2 := &M2AckSidechain{
		SidechainNumber: SidechainNumber(1),
		DescriptionHash: [32]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20},
	}

	if m2.SidechainNumber != SidechainNumber(1) {
		t.Errorf("Expected SidechainNumber to be 1, got %d", m2.SidechainNumber)
	}

	// Test M3ProposeBundle structure
	m3 := &M3ProposeBundle{
		SidechainNumber: SidechainNumber(1),
		BundleTxid:      [32]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20},
	}

	if m3.SidechainNumber != SidechainNumber(1) {
		t.Errorf("Expected SidechainNumber to be 1, got %d", m3.SidechainNumber)
	}

	// Test M4AckBundles structure
	m4 := &M4AckBundles{
		Variant: M4_ONE_BYTE,
		Upvotes: []byte{0x01, 0x02, 0x03},
	}

	if m4.Variant != M4_ONE_BYTE {
		t.Errorf("Expected Variant to be M4_ONE_BYTE, got %d", m4.Variant)
	}

	// Test M7BmmAccept structure
	m7 := &M7BmmAccept{
		SidechainNumber:    SidechainNumber(1),
		SidechainBlockHash: BmmCommitment{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20},
	}

	if m7.SidechainNumber != SidechainNumber(1) {
		t.Errorf("Expected SidechainNumber to be 1, got %d", m7.SidechainNumber)
	}

	// Test M8BmmRequest structure
	m8 := &M8BmmRequest{
		SidechainNumber:        SidechainNumber(1),
		SidechainBlockHash:     BmmCommitment{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20},
		PrevMainchainBlockHash: [32]byte{0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f, 0x40},
	}

	if m8.SidechainNumber != SidechainNumber(1) {
		t.Errorf("Expected SidechainNumber to be 1, got %d", m8.SidechainNumber)
	}
}

// Benchmark tests for serialization performance
func BenchmarkM1ProposeSidechainSerialization(b *testing.B) {
	m1 := &M1ProposeSidechain{
		SidechainNumber: SidechainNumber(1),
		Description:     SidechainDescription("Test sidechain description data"),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = m1.Serialize()
	}
}

func BenchmarkM2AckSidechainSerialization(b *testing.B) {
	m2 := &M2AckSidechain{
		SidechainNumber: SidechainNumber(1),
		DescriptionHash: [32]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = m2.Serialize()
	}
}

func BenchmarkM3ProposeBundleSerialization(b *testing.B) {
	m3 := &M3ProposeBundle{
		SidechainNumber: SidechainNumber(1),
		BundleTxid:      [32]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = m3.Serialize()
	}
}

func BenchmarkM7BmmAcceptSerialization(b *testing.B) {
	m7 := &M7BmmAccept{
		SidechainNumber:    SidechainNumber(1),
		SidechainBlockHash: BmmCommitment{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = m7.Serialize()
	}
}
