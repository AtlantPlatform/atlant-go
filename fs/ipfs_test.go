// Copyright 2017-2019 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package fs

import (
	"encoding/hex"
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

func TestVerifyDataSignature(t *testing.T) {
	NodeID := "14V8Bf9kL1KpZjDrM5JJmA4aT4J2heXQJbh28vmKVqaZbAMuk"
	Signature, errSignature := hex.DecodeString("8ba5575d65a784f2284cf22190a926724f94cc1b7aaa0a61b7a96ab971f6c0d6ed3b01b20e4da30f5e1449c453662461e8df41c829fbb69bae2dc67a3aa50303")
	if errSignature != nil {
		t.Fatal(errSignature)
	}
	SignedData := "100f40031109da31157a0111290aff3031443758595948024e4645504b4b48364148514a52475935034e59ff516d535045556a4c0458686f5263434b5a5531534a3350514e75395251624578735a355466424234513f7054327169440000"
	data, err := hex.DecodeString(SignedData)
	if err != nil {
		t.Fatal(err)
	}
	ok, err := VerifyDataSignature(NodeID, string(Signature), []byte(data))
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("Signed Data was not verified")
	}
	_ = ok
}
