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

/*
func TestVerifyDataSignature(t *testing.T) {
	NodeID := "14V8BjQJ2E2xyQPC3FCFA5Aa9cMxBGpCkeDTPsYLxDqb7Xfne"
	Signature := "0cb19a5136395d295aaf8823721150abc28647391686cac6a85a8a45b9ecf76f5ba776360871eefec151df97c529ef83b29098352b49b7a1321a2383f4b94e01"
	SignedData := "101440031109da31157a0131297a01ff3031434254423459023033534e385358544350384441515650033030ff516d5043436f7550046d4e6e636a4b6f753139594d3741505850707276754a635357434a42575a66713f5048556d3470ff516d5770696863440463785943746f325979703667535a504862474464337a4e6546657161645445673f554737315a63"
	data, err := hex.DecodeString(SignedData)
	if err != nil {
		t.Fatal(err)
	}
	ok, err := VerifyDataSignature(NodeID, string(Signature), []byte(data))
	if err != nil {
		t.Fatal(err)
	}
	// TODO: uncommnet when
	if !ok {
		t.Fatal("Signed Data was not verified")
	}
	_ = ok
}
*/
