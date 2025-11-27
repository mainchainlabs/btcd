// Copyright (c) 2025 The Mainchain Labs
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package txscript

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/wire"
)

// OP_DRIVECHAIN is defined as OP_NOP5 (0xB4) in the drivechain protocol
const OP_DRIVECHAIN = OP_NOP5

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

// SidechainDeclaration represents the structured data inside M1 description
// Format: 1-byte version (0), 1-byte title length, title bytes, description bytes, 32-byte hash_id_1, 20-byte hash_id_2
type SidechainDeclaration struct {
	Title       string
	Description string
	HashID1     [32]byte
	HashID2     [20]byte
}

// ParseSidechainDeclarationError represents errors during sidechain declaration parsing
type ParseSidechainDeclarationError struct {
	Err     error
	Message string
}

func (e *ParseSidechainDeclarationError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Err.Error()
}

func (e *ParseSidechainDeclarationError) Unwrap() error {
	return e.Err
}

// DrivechainMessage represents a basic drivechain message
type DrivechainMessage struct {
	MessageType uint8
}

// CoinbaseMessage represents a drivechain message in a coinbase transaction
// This is an interface that wraps all message types (M1, M2, M3, M4, M7)
type CoinbaseMessage interface {
	MessageType() uint8
	Serialize() ([]byte, error)
}

// Implement CoinbaseMessage for each message type
func (m1 *M1ProposeSidechain) MessageType() uint8 { return MSG_SIDECHAIN_PROPOSAL }
func (m2 *M2AckSidechain) MessageType() uint8     { return MSG_SIDECHAIN_ACK }
func (m3 *M3ProposeBundle) MessageType() uint8    { return MSG_BUNDLE_PROPOSAL }
func (m4 *M4AckBundles) MessageType() uint8       { return MSG_BUNDLE_ACK }
func (m7 *M7BmmAccept) MessageType() uint8        { return MSG_BMM_ACCEPT }

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
	Variant    M4Variant
	Upvotes    []byte   // For OneByte variant (raw bytes)
	UpvotesU16 []uint16 // For TwoBytes variant (little-endian u16 values)
}

const (
	M4_ABSTAIN_ONE_BYTE  = 0xFF
	M4_ABSTAIN_TWO_BYTES = 0xFFFF
	M4_ALARM_ONE_BYTE    = 0xFE
	M4_ALARM_TWO_BYTES   = 0xFFFE
)

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
// This parses OP_RETURN scripts from coinbase transactions
func ParseDrivechainMessage(script []byte) (CoinbaseMessage, error) {
	// Parse OP_RETURN script
	data, err := ParseOPReturnScript(script)
	if err != nil {
		return nil, err
	}

	// Parse the coinbase message from the data
	return ParseCoinbaseMessage(data)
}

// ParseCoinbaseMessage parses a coinbase message from raw data bytes
// Uses tag-based detection to identify message type
func ParseCoinbaseMessage(data []byte) (CoinbaseMessage, error) {
	if len(data) < 4 {
		return nil, errors.New("data too short for drivechain message")
	}

	// Check tags to identify message type
	if len(data) >= 4 && data[0] == M1_TAG[0] && data[1] == M1_TAG[1] && data[2] == M1_TAG[2] && data[3] == M1_TAG[3] {
		return DeserializeM1ProposeSidechain(data)
	}
	if len(data) >= 4 && data[0] == M2_TAG[0] && data[1] == M2_TAG[1] && data[2] == M2_TAG[2] && data[3] == M2_TAG[3] {
		return DeserializeM2AckSidechain(data)
	}
	if len(data) >= 4 && data[0] == M3_TAG[0] && data[1] == M3_TAG[1] && data[2] == M3_TAG[2] && data[3] == M3_TAG[3] {
		return DeserializeM3ProposeBundle(data)
	}
	if len(data) >= 4 && data[0] == M4_TAG[0] && data[1] == M4_TAG[1] && data[2] == M4_TAG[2] && data[3] == M4_TAG[3] {
		return DeserializeM4AckBundles(data)
	}
	if len(data) >= 4 && data[0] == M7_TAG[0] && data[1] == M7_TAG[1] && data[2] == M7_TAG[2] && data[3] == M7_TAG[3] {
		return DeserializeM7BmmAccept(data)
	}

	return nil, errors.New("unknown drivechain message tag")
}

// ParseOPReturnScript extracts data from an OP_RETURN script
// Format: OP_RETURN <data>
func ParseOPReturnScript(script []byte) ([]byte, error) {
	if len(script) < 1 {
		return nil, errors.New("script too short")
	}

	// Check for OP_RETURN
	if script[0] != OP_RETURN {
		return nil, errors.New("not an OP_RETURN script")
	}

	// Single OP_RETURN (no data)
	if len(script) == 1 {
		return nil, errors.New("OP_RETURN script has no data")
	}

	// Parse the data push after OP_RETURN
	// Script version 0 is the only supported version
	tokenizer := MakeScriptTokenizer(0, script[1:])
	if !tokenizer.Next() {
		return nil, errors.New("failed to parse OP_RETURN data")
	}

	// Get the data
	data := tokenizer.Data()
	if len(data) == 0 {
		return nil, errors.New("OP_RETURN data is empty")
	}

	// Ensure we consumed all data
	if !tokenizer.Done() {
		return nil, errors.New("extra data after OP_RETURN push")
	}

	return data, nil
}

// CreateOPReturnOutput creates an OP_RETURN script output with the given data
func CreateOPReturnOutput(data []byte) ([]byte, error) {
	return NullDataScript(data)
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
	switch m4.Variant {
	case M4_ONE_BYTE:
		buf = append(buf, m4.Upvotes...)
	case M4_TWO_BYTES:
		for _, upvote := range m4.UpvotesU16 {
			upvoteBytes := make([]byte, 2)
			binary.LittleEndian.PutUint16(upvoteBytes, upvote)
			buf = append(buf, upvoteBytes...)
		}
	}

	return buf, nil
}

// DeserializeM4AckBundles deserializes bytes to an M4AckBundles
func DeserializeM4AckBundles(data []byte) (*M4AckBundles, error) {
	if len(data) < 5 {
		return nil, errors.New("data too short for M4 message")
	}

	// Check tag
	if data[0] != M4_TAG[0] || data[1] != M4_TAG[1] || data[2] != M4_TAG[2] || data[3] != M4_TAG[3] {
		return nil, errors.New("invalid M4 tag")
	}

	data = data[4:] // Skip tag

	if len(data) < 1 {
		return nil, errors.New("data too short for M4 variant")
	}

	variant := M4Variant(data[0])
	data = data[1:]

	m4 := &M4AckBundles{
		Variant: variant,
	}

	switch variant {
	case M4_REPEAT_PREVIOUS:
		// No upvotes
		return m4, nil

	case M4_ONE_BYTE:
		// Rest of data is upvotes as bytes
		m4.Upvotes = make([]byte, len(data))
		copy(m4.Upvotes, data)
		return m4, nil

	case M4_TWO_BYTES:
		// Parse pairs of bytes as little-endian u16
		if len(data)%2 != 0 {
			return nil, errors.New("M4 TwoBytes variant requires even number of bytes")
		}
		m4.UpvotesU16 = make([]uint16, 0, len(data)/2)
		for i := 0; i < len(data); i += 2 {
			upvote := binary.LittleEndian.Uint16(data[i : i+2])
			m4.UpvotesU16 = append(m4.UpvotesU16, upvote)
		}
		return m4, nil

	case M4_LEADING_BY_50:
		// No upvotes
		return m4, nil

	default:
		return nil, errors.New("invalid M4 variant")
	}
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
// This expects data AFTER OP_RETURN and length byte have been parsed
// Format: M8_TAG, sidechain_number, sidechain_block_hash, prev_mainchain_block_hash
func DeserializeM8BmmRequest(data []byte) (*M8BmmRequest, error) {
	const (
		headerLength                 = 3 // M8_TAG length
		sidechainNumberLength        = 1
		sidechainBlockHashLength     = 32
		prevMainchainBlockHashLength = 32
	)

	m8BmmRequestLength := headerLength + sidechainNumberLength + sidechainBlockHashLength + prevMainchainBlockHashLength

	if len(data) < m8BmmRequestLength {
		return nil, errors.New("data too short for M8 message")
	}

	// Check tag
	if data[0] != M8_TAG[0] || data[1] != M8_TAG[1] || data[2] != M8_TAG[2] {
		return nil, errors.New("invalid M8 tag")
	}

	data = data[headerLength:] // Skip tag

	m8 := &M8BmmRequest{
		SidechainNumber: SidechainNumber(data[0]),
	}
	data = data[sidechainNumberLength:]

	// Sidechain block hash
	copy(m8.SidechainBlockHash[:], data[:sidechainBlockHashLength])
	data = data[sidechainBlockHashLength:]

	// Previous mainchain block hash
	copy(m8.PrevMainchainBlockHash[:], data[:prevMainchainBlockHashLength])

	return m8, nil
}

// ParseM8Transaction parses an M8 BMM request from a transaction's first output
// M8 format in script: OP_RETURN (0x6A), length (0x44), M8_TAG, sidechain_number, sidechain_block_hash, prev_mainchain_block_hash
func ParseM8Transaction(tx *wire.MsgTx) (*M8BmmRequest, error) {
	if len(tx.TxOut) == 0 {
		return nil, errors.New("transaction has no outputs")
	}

	// Get first output
	output := tx.TxOut[0]
	script := output.PkScript

	// M8 has special format: OP_RETURN, length byte, then data
	if len(script) < 2 {
		return nil, errors.New("script too short")
	}

	if script[0] != OP_RETURN {
		return nil, errors.New("not an OP_RETURN script")
	}

	// M8 includes length byte (0x44) after OP_RETURN
	const m8Length = 0x44
	if script[1] != m8Length {
		return nil, fmt.Errorf("invalid M8 length byte: expected 0x44, got 0x%02x", script[1])
	}

	// Extract data after OP_RETURN and length byte
	if len(script) < 2+int(m8Length) {
		return nil, errors.New("script too short for M8 data")
	}

	data := script[2:] // Skip OP_RETURN and length byte

	return DeserializeM8BmmRequest(data)
}

// ParseSidechainDeclaration parses a SidechainDeclaration from raw bytes
// Format: 1-byte version (0), 1-byte title length, title bytes, description bytes (calculated), 32-byte hash_id_1, 20-byte hash_id_2
// Description length = total_length - (1 + 1 + title_len + 32 + 20)
func ParseSidechainDeclaration(data []byte) (*SidechainDeclaration, error) {
	const (
		versionLength     = 1
		titleLengthLength = 1
		hashID1Length     = 32
		hashID2Length     = 20
		version0          = 0
	)

	if len(data) < versionLength+titleLengthLength+hashID1Length+hashID2Length {
		return nil, &ParseSidechainDeclarationError{
			Err:     errors.New("data too short"),
			Message: "failed to deserialize sidechain declaration",
		}
	}

	// Check version
	if data[0] != version0 {
		return nil, &ParseSidechainDeclarationError{
			Err:     fmt.Errorf("unknown version: %d", data[0]),
			Message: "unknown sidechain declaration version",
		}
	}

	data = data[versionLength:]

	// Read title length
	titleLen := int(data[0])
	data = data[titleLengthLength:]

	if len(data) < titleLen {
		return nil, &ParseSidechainDeclarationError{
			Err:     errors.New("data too short for title"),
			Message: "failed to deserialize sidechain declaration",
		}
	}

	// Read title
	titleBytes := data[:titleLen]
	title := string(titleBytes)
	data = data[titleLen:]

	// Calculate description length
	// Description length = total_length - (version + title_len_byte + title + hash_id_1 + hash_id_2)
	descriptionLength := len(data) - hashID1Length - hashID2Length
	if descriptionLength < 0 {
		return nil, &ParseSidechainDeclarationError{
			Err:     errors.New("data too short for description"),
			Message: "failed to deserialize sidechain declaration",
		}
	}

	// Read description
	descriptionBytes := data[:descriptionLength]
	description := string(descriptionBytes)
	data = data[descriptionLength:]

	// Read hash_id_1
	if len(data) < hashID1Length {
		return nil, &ParseSidechainDeclarationError{
			Err:     errors.New("data too short for hash_id_1"),
			Message: "failed to deserialize sidechain declaration",
		}
	}
	var hashID1 [32]byte
	copy(hashID1[:], data[:hashID1Length])
	data = data[hashID1Length:]

	// Read hash_id_2
	if len(data) < hashID2Length {
		return nil, &ParseSidechainDeclarationError{
			Err:     errors.New("data too short for hash_id_2"),
			Message: "failed to deserialize sidechain declaration",
		}
	}
	var hashID2 [20]byte
	copy(hashID2[:], data[:hashID2Length])

	return &SidechainDeclaration{
		Title:       title,
		Description: description,
		HashID1:     hashID1,
		HashID2:     hashID2,
	}, nil
}

// SerializeSidechainDeclaration serializes a SidechainDeclaration to bytes
func SerializeSidechainDeclaration(decl *SidechainDeclaration) (SidechainDescription, error) {
	const version = 0

	titleBytes := []byte(decl.Title)
	if len(titleBytes) > 255 {
		return nil, errors.New("title too long (max 255 bytes)")
	}

	descriptionBytes := []byte(decl.Description)

	buf := make([]byte, 0, 1+1+len(titleBytes)+len(descriptionBytes)+32+20)

	// Version
	buf = append(buf, version)

	// Title length
	buf = append(buf, byte(len(titleBytes)))

	// Title
	buf = append(buf, titleBytes...)

	// Description
	buf = append(buf, descriptionBytes...)

	// Hash ID 1
	buf = append(buf, decl.HashID1[:]...)

	// Hash ID 2
	buf = append(buf, decl.HashID2[:]...)

	return SidechainDescription(buf), nil
}

// CreateSidechainProposal creates an M1ProposeSidechain from a SidechainNumber and SidechainDeclaration
// This serializes the declaration into the SidechainDescription format expected by M1
func CreateSidechainProposal(sidechainNumber SidechainNumber, declaration *SidechainDeclaration) (*M1ProposeSidechain, SidechainDescription, error) {
	description, err := SerializeSidechainDeclaration(declaration)
	if err != nil {
		return nil, nil, err
	}

	m1 := &M1ProposeSidechain{
		SidechainNumber: sidechainNumber,
		Description:     description,
	}

	return m1, description, nil
}

// ParseOPDrivechain parses an OP_DRIVECHAIN script and extracts the sidechain number
// Format: OP_DRIVECHAIN (0xB4), OP_PUSHBYTES_1 (0x01), sidechain_number (1 byte), OP_TRUE (0x51)
// OP_PUSHBYTES_1 is the same as OP_DATA_1 (0x01)
func ParseOPDrivechain(script []byte) (SidechainNumber, error) {
	if len(script) < 4 {
		return 0, errors.New("script too short for OP_DRIVECHAIN")
	}

	// Check for OP_DRIVECHAIN
	if script[0] != OP_DRIVECHAIN {
		return 0, errors.New("not an OP_DRIVECHAIN script")
	}

	// Check for OP_PUSHBYTES_1 (OP_DATA_1 = 0x01)
	if script[1] != OP_DATA_1 {
		return 0, errors.New("invalid OP_DRIVECHAIN format: expected OP_PUSHBYTES_1 (OP_DATA_1)")
	}

	// OP_DATA_1 means 1 byte of data follows
	if len(script) < 4 {
		return 0, errors.New("script too short for OP_DRIVECHAIN data")
	}

	sidechainNumber := SidechainNumber(script[2])

	// Check for OP_TRUE
	if script[3] != OP_TRUE {
		return 0, errors.New("invalid OP_DRIVECHAIN format: expected OP_TRUE")
	}

	return sidechainNumber, nil
}

// CreateOPDrivechainScript creates an OP_DRIVECHAIN script for the specified sidechain
// Format: OP_DRIVECHAIN (0xB4), OP_DATA_1 (0x01), sidechain_number, OP_TRUE (0x51)
// Note: We must use OP_DATA_1 explicitly, not OP_1-OP_16, even for small numbers
func CreateOPDrivechainScript(sidechainNumber SidechainNumber) ([]byte, error) {
	// Manually construct to ensure OP_DATA_1 is used (not OP_1-OP_16 for values 1-16)
	script := make([]byte, 4)
	script[0] = OP_DRIVECHAIN
	script[1] = OP_DATA_1
	script[2] = byte(sidechainNumber)
	script[3] = OP_TRUE
	return script, nil
}

// CreateM5DepositOutput creates a deposit output (M5) with OP_DRIVECHAIN script
// Takes: SidechainNumber, old CTIP amount, deposit amount
// Returns: TxOut with OP_DRIVECHAIN script and value = old_ctip_amount + deposit_amount
func CreateM5DepositOutput(sidechainNumber SidechainNumber, oldCtipAmount int64, depositAmount int64) (*wire.TxOut, error) {
	script, err := CreateOPDrivechainScript(sidechainNumber)
	if err != nil {
		return nil, err
	}

	return &wire.TxOut{
		Value:    oldCtipAmount + depositAmount,
		PkScript: script,
	}, nil
}
