// Copyright 2017-21 Digital Asset Exchange Limited. All rights reserved.
// Use of this source code is governed by BSD 3-Clause "New" or "Revised"
// License (BSD 3) that can be found in the LICENSE file.

package fs

import (
	"strings"
	"testing"
)

func TestNewPrivateKey(t *testing.T) {
	buf, err := NewPrivateKey()
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(string(buf), "\n")
	if len(lines) != 3 {
		t.Fatal("NewPrivateKey: expected 3 lines")
	}
	if lines[0] != "/key/swarm/psk/1.0.0/" {
		t.Fatal("NewPrivateKey: expected type to be /key/swarm/psk/1.0.0/")
	}
	if lines[1] != "/base16/" {
		t.Fatal("NewPrivateKey: expected encoding of /base16/")
	}
	if len(lines[2]) == 0 {
		t.Fatal("NewPrivateKey: expected valid private key, received: " + lines[2])
	}
}
