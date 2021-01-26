// Copyright 2017-2021 Digital Asset Exchange Limited. All rights reserved.
// Use of this source code is governed by BSD-3-Clause "New" or "Revised"
// License (BSD-3-Clause) that can be found in the LICENSE file.

package main

import "testing"

func TestLocalIpDetection(t *testing.T) {

	authLocal := NewAuthTrustLocal()
	expectedLocal := []string{
		"127.0.0.1", "127.0.0.9",
		"192.168.0.1", "192.168.0.101",
		"172.17.0.1", "172.20.0.100",
	}
	expectedRemote := []string{
		"8.8.8.8", "81.19.74.1", "216.58.215.68",
	}
	for _, localAddr := range expectedLocal {
		if !authLocal.IsLocal(localAddr) {
			t.Error(localAddr + " should be local, but detected as remote")
		}
	}

	for _, remoteAddr := range expectedRemote {
		if authLocal.IsLocal(remoteAddr) {
			t.Error(remoteAddr + " should be remote, but detected as remote")
		}
	}
}
