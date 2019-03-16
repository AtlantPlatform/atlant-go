// Package peer implements an object used to represent peers in the ipfs network.
package peer

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"

	ic "github.com/libp2p/go-libp2p-crypto"
	b58 "github.com/mr-tron/base58/base58"
	mh "github.com/multiformats/go-multihash"
)

// Code is type from deprecated go-multicodec-packed
type Code uint64

const mcEd25519Pub = Code(0xed)

func mcSplitPrefix(data []byte) (Code, []byte) {
	c, n := binary.Uvarint(data)
	return Code(c), data[n:]
}

var (
	// ErrEmptyPeerID is an error for empty peer ID.
	ErrEmptyPeerID = errors.New("empty peer ID")
	// ErrNoPublicKey is an error for peer IDs that don't embed public keys
	ErrNoPublicKey = errors.New("public key is not embedded in peer ID")
)

// AdvancedEnableInlining enables automatically inlining keys shorter than
// 42 bytes into the peer ID (using the "identity" multihash function).
//
// WARNING: This flag will likely be set to false in the future and eventually
// be removed in favor of using a hash function specified by the key itself.
// See: https://github.com/libp2p/specs/issues/138
//
// DO NOT change this flag unless you know what you're doing.
//
// This currently defaults to true for backwards compatibility but will likely
// be set to false by default when an upgrade path is determined.
var AdvancedEnableInlining = true

const maxInlineKeyLength = 42

// ID is a libp2p peer identity.
type ID string

// Pretty returns a b58-encoded string of the ID
func (id ID) Pretty() string {
	return IDB58Encode(id)
}

// Loggable returns a pretty peerID string in loggable JSON format
func (id ID) Loggable() map[string]interface{} {
	return map[string]interface{}{
		"peerID": id.Pretty(),
	}
}

// String prints out the peer.
//
// TODO(brian): ensure correctness at ID generation and
// enforce this by only exposing functions that generate
// IDs safely. Then any peer.ID type found in the
// codebase is known to be correct.
func (id ID) String() string {
	pid := id.Pretty()
	if len(pid) <= 10 {
		return fmt.Sprintf("<peer.ID %s>", pid)
	}
	return fmt.Sprintf("<peer.ID %s*%s>", pid[:2], pid[len(pid)-6:])
}

// MatchesPrivateKey tests whether this ID was derived from sk
func (id ID) MatchesPrivateKey(sk ic.PrivKey) bool {
	return id.MatchesPublicKey(sk.GetPublic())
}

// MatchesPublicKey tests whether this ID was derived from pk
func (id ID) MatchesPublicKey(pk ic.PubKey) bool {
	oid, err := IDFromEd25519PublicKey(pk)
	if err != nil {
		return false
	}
	return oid == id
}

// ErrMultihashDecode - error of decoding multihash
var ErrMultihashDecode = errors.New("unable to decode multihash")

// ErrMultihashCodec - error of multihash coded
var ErrMultihashCodec = errors.New("unexpected multihash codec")

// ErrMultihashLength - error of multihash length
var ErrMultihashLength = errors.New("unexpected multihash length")

// ErrCodePrefix - error of Ed25519 prefix code
var ErrCodePrefix = errors.New("unexpected code prefix")

// ExtractEd25519PublicKey attempts to extract the public key from an ID
func (id ID) ExtractEd25519PublicKey() (ic.PubKey, error) {
	// ed25519 pubkey identity format
	// <identity mc><length (2 + 32 = 34)><ed25519-pub mc><ed25519 pubkey>
	// <0x00       ><0x22                ><0xed01        ><ed25519 pubkey>

	var nilPubKey ic.PubKey

	// Decode multihash
	decoded, err := mh.Decode([]byte(id))
	if err != nil {
		return nilPubKey, ErrMultihashDecode
	}

	// Check ID multihash codec
	if decoded.Code != mh.ID {
		return nilPubKey, ErrMultihashCodec
	}

	// Check multihash length
	if decoded.Length != 2+32 {
		return nilPubKey, ErrMultihashLength
	}

	// Split prefix
	code, pubKeyBytes := mcSplitPrefix(decoded.Digest)

	// Check ed25519 code
	if code != mcEd25519Pub {
		return nilPubKey, ErrCodePrefix
	}

	// Unmarshall public key
	pubKey, err := ic.UnmarshalEd25519PublicKey(pubKeyBytes)
	if err != nil {
		// Should never occur because of the check decoded.Length != 2+32
		return nilPubKey, fmt.Errorf("Unexpected error unmarshalling Ed25519 public key")
	}

	return pubKey, nil
}

// ExtractPublicKey attempts to extract the public key from an ID
//
// This method returns ErrNoPublicKey if the peer ID looks valid but it can't extract
// the public key.
func (id ID) ExtractPublicKey() (ic.PubKey, error) {
	decoded, err := mh.Decode([]byte(id))
	if err != nil {
		return nil, err
	}
	if decoded.Code != mh.ID {
		return nil, ErrNoPublicKey
	}
	pk, err := ic.UnmarshalPublicKey(decoded.Digest)
	if err != nil {
		return nil, err
	}
	return pk, nil
}

// Validate check if ID is empty or not
func (id ID) Validate() error {
	if id == ID("") {
		return ErrEmptyPeerID
	}

	return nil
}

// IDFromString cast a string to ID type, and validate
// the id to make sure it is a multihash.
func IDFromString(s string) (ID, error) {
	if _, err := mh.Cast([]byte(s)); err != nil {
		return ID(""), err
	}
	return ID(s), nil
}

// IDFromBytes cast a string to ID type, and validate
// the id to make sure it is a multihash.
func IDFromBytes(b []byte) (ID, error) {
	if _, err := mh.Cast(b); err != nil {
		return ID(""), err
	}
	return ID(b), nil
}

// IDB58Decode returns a b58-decoded Peer
func IDB58Decode(s string) (ID, error) {
	m, err := mh.FromB58String(s)
	if err != nil {
		return "", err
	}
	return ID(m), err
}

// IDB58Encode returns b58-encoded string
func IDB58Encode(id ID) string {
	return b58.Encode([]byte(id))
}

// IDHexDecode returns a hex-decoded Peer
func IDHexDecode(s string) (ID, error) {
	m, err := mh.FromHexString(s)
	if err != nil {
		return "", err
	}
	return ID(m), err
}

// IDHexEncode returns hex-encoded string
func IDHexEncode(id ID) string {
	return hex.EncodeToString([]byte(id))
}

// IDFromPublicKey returns the Peer ID corresponding to pk
func IDFromPublicKey(pk ic.PubKey) (ID, error) {
	b, err := pk.Bytes()
	if err != nil {
		return "", err
	}
	var alg uint64 = mh.SHA2_256
	if AdvancedEnableInlining && len(b) <= maxInlineKeyLength {
		alg = mh.ID
	}
	hash, _ := mh.Sum(b, alg, -1)
	return ID(hash), nil
}

// IDFromEd25519PublicKey returns the Peer ID corresponding to Id25519 pk
func IDFromEd25519PublicKey(pk ic.PubKey) (ID, error) {
	b, err := pk.Bytes()
	if err != nil {
		return "", err
	}

	// Build the ed25519 public key multi-codec
	Ed25519PubMultiCodec := make([]byte, 2)
	binary.PutUvarint(Ed25519PubMultiCodec, uint64(mcEd25519Pub))

	hash, err := mh.Sum(append(Ed25519PubMultiCodec, b[len(b)-32:]...), mh.ID, 34)
	if err != nil {
		return "", err
	}
	return ID(hash), nil
}

// IDFromPrivateKey returns the Peer ID corresponding to sk
func IDFromPrivateKey(sk ic.PrivKey) (ID, error) {
	return IDFromEd25519PublicKey(sk.GetPublic())
}

// IDSlice for sorting peers
type IDSlice []ID

func (es IDSlice) Len() int           { return len(es) }
func (es IDSlice) Swap(i, j int)      { es[i], es[j] = es[j], es[i] }
func (es IDSlice) Less(i, j int) bool { return string(es[i]) < string(es[j]) }
