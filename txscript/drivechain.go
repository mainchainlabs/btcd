// Copyright (c) 2025 The Mainchain Labs
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
	MSG_BMM_ACCEPT         = 7 // M7 - BMM acceptance
	MSG_BMM_REQUEST        = 8 // M8 - BMM request
)

var (
	M1_TAG = [4]byte{0xD5, 0xE0, 0xC4, 0xAF}
	M2_TAG = [4]byte{0xD6, 0xE1, 0xC5, 0xDF}
	M3_TAG = [4]byte{0xD4, 0x5A, 0xA9, 0x43}
	M4_TAG = [4]byte{0xD7, 0x7D, 0x17, 0x76}
	M7_TAG = [4]byte{0xD1, 0x61, 0x73, 0x68}
	M8_TAG = [3]byte{0x00, 0xBF, 0x00}
)

// SidechainNumber represents a sidechain identifier (u8)
type SidechainNumber uint8

// SidechainDescription represents raw bytes for sidechain description
type SidechainDescription []byte

// BmmCommitment represents a BMM commitment (32 bytes)
type BmmCommitment [32]byte

// DrivechainMessage represents a basic drivechain message
type DrivechainMessage struct {
	MessageType uint8
}

// M1ProposeSidechain (M1) represents a sidechain proposal message
type M1ProposeSidechain struct {
	SidechainNumber SidechainNumber
	Description     SidechainDescription
}

// M2AckSidechain (M2) represents a sidechain acknowledgment message
type M2AckSidechain struct {
	SidechainNumber SidechainNumber
	DescriptionHash [32]byte // sha256d::Hash
}

// M3ProposeBundle (M3) represents a bundle proposal message
type M3ProposeBundle struct {
	SidechainNumber SidechainNumber
	BundleTxid      [32]byte
}

// M4AckBundles (M4) represents bundle acknowledgment variants
type M4AckBundles struct {
	Variant M4Variant
	Upvotes []byte // For OneByte and TwoBytes variants
}

// M4Variant represents the different M4 acknowledgment types
type M4Variant uint8

const (
	M4_REPEAT_PREVIOUS M4Variant = 0x00
	M4_ONE_BYTE        M4Variant = 0x01
	M4_TWO_BYTES       M4Variant = 0x02
	M4_LEADING_BY_50   M4Variant = 0x03
)

// M7BmmAccept (M7) represents BMM acceptance message
type M7BmmAccept struct {
	SidechainNumber    SidechainNumber
	SidechainBlockHash BmmCommitment
}

// M8BmmRequest (M8) represents BMM request message
type M8BmmRequest struct {
	SidechainNumber        SidechainNumber
	SidechainBlockHash     BmmCommitment
	PrevMainchainBlockHash [32]byte
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
	if msgType < MSG_SIDECHAIN_PROPOSAL || msgType > MSG_BMM_REQUEST {
		return nil, errors.New("invalid drivechain message type")
	}

	msg := &DrivechainMessage{
		MessageType: msgType,
	}

	return msg, nil
}

// IsValidMessageType checks if a message type is valid
func IsValidMessageType(msgType uint8) bool {
	return msgType >= MSG_SIDECHAIN_PROPOSAL && msgType <= MSG_BMM_REQUEST
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
	case MSG_BMM_ACCEPT:
		return "BMM Accept (M7)"
	case MSG_BMM_REQUEST:
		return "BMM Request (M8)"
	default:
		return "Unknown Message Type"
	}
}

// Serialize serializes an M1ProposeSidechain to bytes
func (m1 *M1ProposeSidechain) Serialize() ([]byte, error) {
	buf := make([]byte, 0, 4+1+len(m1.Description))

	// Tag
	buf = append(buf, M1_TAG[:]...)

	// Sidechain number
	buf = append(buf, byte(m1.SidechainNumber))

	// Description (raw bytes)
	buf = append(buf, m1.Description...)

	return buf, nil
}

// DeserializeM1ProposeSidechain deserializes bytes to an M1ProposeSidechain
func DeserializeM1ProposeSidechain(data []byte) (*M1ProposeSidechain, error) {
	if len(data) < 5 {
		return nil, errors.New("data too short for M1 message")
	}

	// Check tag
	if data[0] != M1_TAG[0] || data[1] != M1_TAG[1] || data[2] != M1_TAG[2] || data[3] != M1_TAG[3] {
		return nil, errors.New("invalid M1 tag")
	}

	data = data[4:] // Skip tag

	if len(data) < 1 {
		return nil, errors.New("data too short for sidechain number")
	}

	m1 := &M1ProposeSidechain{
		SidechainNumber: SidechainNumber(data[0]),
	}
	data = data[1:]

	// Description is the rest of the data
	m1.Description = make([]byte, len(data))
	copy(m1.Description, data)

	return m1, nil
}

// Serialize serializes an M2AckSidechain to bytes
func (m2 *M2AckSidechain) Serialize() ([]byte, error) {
	buf := make([]byte, 0, 4+1+32)

	// Tag
	buf = append(buf, M2_TAG[:]...)

	// Sidechain number
	buf = append(buf, byte(m2.SidechainNumber))

	// Description hash
	buf = append(buf, m2.DescriptionHash[:]...)

	return buf, nil
}

// DeserializeM2AckSidechain deserializes bytes to an M2AckSidechain
func DeserializeM2AckSidechain(data []byte) (*M2AckSidechain, error) {
	if len(data) < 37 {
		return nil, errors.New("data too short for M2 message")
	}

	// Check tag
	if data[0] != M2_TAG[0] || data[1] != M2_TAG[1] || data[2] != M2_TAG[2] || data[3] != M2_TAG[3] {
		return nil, errors.New("invalid M2 tag")
	}

	data = data[4:] // Skip tag

	m2 := &M2AckSidechain{
		SidechainNumber: SidechainNumber(data[0]),
	}
	data = data[1:]

	// Description hash
	copy(m2.DescriptionHash[:], data[:32])

	return m2, nil
}

// Serialize serializes an M3ProposeBundle to bytes
func (m3 *M3ProposeBundle) Serialize() ([]byte, error) {
	buf := make([]byte, 0, 4+1+32)

	// Tag
	buf = append(buf, M3_TAG[:]...)

	// Sidechain number
	buf = append(buf, byte(m3.SidechainNumber))

	// Bundle txid
	buf = append(buf, m3.BundleTxid[:]...)

	return buf, nil
}

// DeserializeM3ProposeBundle deserializes bytes to an M3ProposeBundle
func DeserializeM3ProposeBundle(data []byte) (*M3ProposeBundle, error) {
	if len(data) < 37 {
		return nil, errors.New("data too short for M3 message")
	}

	// Check tag
	if data[0] != M3_TAG[0] || data[1] != M3_TAG[1] || data[2] != M3_TAG[2] || data[3] != M3_TAG[3] {
		return nil, errors.New("invalid M3 tag")
	}

	data = data[4:] // Skip tag

	m3 := &M3ProposeBundle{
		SidechainNumber: SidechainNumber(data[0]),
	}
	data = data[1:]

	// Bundle txid
	copy(m3.BundleTxid[:], data[:32])

	return m3, nil
}

// Serialize serializes an M4AckBundles to bytes
func (m4 *M4AckBundles) Serialize() ([]byte, error) {
	buf := make([]byte, 0, 4+1+len(m4.Upvotes))

	// Tag
	buf = append(buf, M4_TAG[:]...)

	// Variant
	buf = append(buf, byte(m4.Variant))

	// Upvotes (for OneByte and TwoBytes variants)
	buf = append(buf, m4.Upvotes...)

	return buf, nil
}

// Serialize serializes an M7BmmAccept to bytes
func (m7 *M7BmmAccept) Serialize() ([]byte, error) {
	buf := make([]byte, 0, 4+1+32)

	// Tag
	buf = append(buf, M7_TAG[:]...)

	// Sidechain number
	buf = append(buf, byte(m7.SidechainNumber))

	// Sidechain block hash
	buf = append(buf, m7.SidechainBlockHash[:]...)

	return buf, nil
}

// DeserializeM7BmmAccept deserializes bytes to an M7BmmAccept
func DeserializeM7BmmAccept(data []byte) (*M7BmmAccept, error) {
	if len(data) < 37 {
		return nil, errors.New("data too short for M7 message")
	}

	// Check tag
	if data[0] != M7_TAG[0] || data[1] != M7_TAG[1] || data[2] != M7_TAG[2] || data[3] != M7_TAG[3] {
		return nil, errors.New("invalid M7 tag")
	}

	data = data[4:] // Skip tag

	m7 := &M7BmmAccept{
		SidechainNumber: SidechainNumber(data[0]),
	}
	data = data[1:]

	// Sidechain block hash
	copy(m7.SidechainBlockHash[:], data[:32])

	return m7, nil
}

// Serialize serializes an M8BmmRequest to bytes
func (m8 *M8BmmRequest) Serialize() ([]byte, error) {
	buf := make([]byte, 0, 3+1+32+32)

	// Tag
	buf = append(buf, M8_TAG[:]...)

	// Sidechain number
	buf = append(buf, byte(m8.SidechainNumber))

	// Sidechain block hash
	buf = append(buf, m8.SidechainBlockHash[:]...)

	// Previous mainchain block hash
	buf = append(buf, m8.PrevMainchainBlockHash[:]...)

	return buf, nil
}

// DeserializeM8BmmRequest deserializes bytes to an M8BmmRequest
func DeserializeM8BmmRequest(data []byte) (*M8BmmRequest, error) {
	if len(data) < 68 {
		return nil, errors.New("data too short for M8 message")
	}

	// Check tag
	if data[0] != M8_TAG[0] || data[1] != M8_TAG[1] || data[2] != M8_TAG[2] {
		return nil, errors.New("invalid M8 tag")
	}

	data = data[3:] // Skip tag

	m8 := &M8BmmRequest{
		SidechainNumber: SidechainNumber(data[0]),
	}
	data = data[1:]

	// Sidechain block hash
	copy(m8.SidechainBlockHash[:], data[:32])
	data = data[32:]

	// Previous mainchain block hash
	copy(m8.PrevMainchainBlockHash[:], data[:32])

	return m8, nil
}
