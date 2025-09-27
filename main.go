// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/iotzf/tsvc/pkg/shutdown"
)

func main() {
	fmt.Println("t svc main endpoint")

	// register shutdown hook, do something before shutdown
	// default signals are SIGINT and SIGTERM
	// you can add more signals by WithSignals method
	// e.g. shutdown.NewHook().WithSignals(syscall.SIGHUP, syscall.SIGQUIT).Close(...)
	// if you want to customize signals
	// otherwise, just use shutdown.NewHook().Close(...)
	shutdown.NewHook().Close(
		func() {
			fmt.Println("do something before shutdown")
		})
}
