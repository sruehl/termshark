// Copyright 2019-2022 Graham Clark. All rights reserved.  Use of this source
// code is governed by the MIT license that can be found in the LICENSE
// file.

//go:build windows
// +build windows

package system

func CloseDescriptor(fd int) {
}

func FileRegularOrLink(filename string) bool {
	return true
}
