// Copyright 2017-2021 Digital Asset Exchange Limited. All rights reserved.
// Use of this source code is governed by BSD-3-Clause "New" or "Revised"
// License (BSD-3-Clause) that can be found in the LICENSE file.

package state

import (
	"fmt"
	"runtime"
)

func orPanic(err error, finalizers ...func()) {
	if err != nil {
		for _, fn := range finalizers {
			fn()
		}
		panic(err)
	}
}

func checkErr(err *error) {
	if v := recover(); v != nil {
		*err = fmt.Errorf("%+v", v)
	}
}

func checkErrStack(err *error) {
	if v := recover(); v != nil {
		stack := make([]byte, 32*1024)
		n := runtime.Stack(stack, false)
		switch event := v.(type) {
		case error:
			*err = fmt.Errorf("%s\n%s", event.Error(), stack[:n])
		default:
			*err = fmt.Errorf("%+v %s", v, stack[:n])
		}
	}
}
