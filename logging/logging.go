// Copyright 2017-2019 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

// Package logging contains various logging helpers.
package logging

import (
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

func WithFn(fields ...log.Fields) log.Fields {
	if len(fields) > 0 && fields[0] != nil {
		result := copyFields(fields[0])
		result["Fn"] = getCallerName()
		return result
	}
	return log.Fields{
		"Fn": getCallerName(),
	}
}

func WithMore(fields log.Fields, add log.Fields) log.Fields {
	fields = copyFields(fields)
	for k, v := range add {
		fields[k] = v
	}
	return fields
}

func copyFields(fields log.Fields) log.Fields {
	ff := make(log.Fields, len(fields))
	for k, v := range fields {
		ff[k] = v
	}
	return ff
}

func FnName() string {
	return getCallerName()
}

func getCallerName() string {
	pc, _, _, _ := runtime.Caller(2)
	fullName := runtime.FuncForPC(pc).Name()
	parts := strings.Split(fullName, "/")
	nameParts := strings.Split(parts[len(parts)-1], ".")
	return nameParts[len(nameParts)-1]
}
