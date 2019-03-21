// Copyright 2017-2019 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package fs

import (
	"testing"

	config "github.com/ipfs/go-ipfs-config"
	peer "github.com/libp2p/go-libp2p-peer"
)

func TestPeerSignature(t *testing.T) {
	PeerID := "14V8BhkDzLXGFxfay4uMn9LMGC3XMCfWv7irKLp3S1xHxwt8o"
	PrivKey := "CAESQPD8Zz2r82zYFGljoYrRXHb3azw/6Z+xdtkZozgGOTZv5zZWxuQua7i8/+EeFb6Ps3uN6CQ/AaT9WOXsSqCo5ZA="
	Data := "Hello, Data to be signed"

	cfg := config.Identity{PeerID: PeerID, PrivKey: PrivKey}
	sk, err := cfg.DecodePrivateKey("passphrase todo!")
	if err != nil {
		t.Fatal(err)
	}
	if sk.GetPublic().Type() != 1 {
		t.Fatalf("Decoded key type: %d, expected 1 (ED25519)", sk.GetPublic().Type())
	}
	id2, err := peer.IDFromEd25519PublicKey(sk.GetPublic())
	if err != nil {
		t.Fatal(err)
	}
	if id2.Pretty() != PeerID {
		t.Fatalf("Decoded public key in config does not match Peer ID: %s, expected %s", id2.Pretty(), PeerID)
	}
	signed, err := sk.Sign([]byte(Data))
	if err != nil {
		t.Fatal(err)
	}
	// fmt.Println("signed hex:", hex.EncodeToString(signed))
	id, err := peer.IDB58Decode(PeerID)
	if err != nil {
		t.Fatal(err)
	}
	pk, err := id.ExtractEd25519PublicKey()
	if err != nil {
		t.Fatal(err)
	}
	verified, err := pk.Verify([]byte(Data), signed)
	if err != nil {
		t.Fatal(err)
	}
	if !verified {
		t.Fatal("Signature not verified")
	}
}
