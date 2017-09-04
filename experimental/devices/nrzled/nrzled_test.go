// Copyright 2017 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package nrzled

import (
	"bytes"
	"strconv"
	"testing"
)

func TestNRZ(t *testing.T) {
	data := []struct {
		in       byte
		expected uint32
	}{
		{0x00, 0x924924},
		{0x01, 0x924926},
		{0x02, 0x924934},
		{0x04, 0x9249A4},
		{0x08, 0x924D24},
		{0x10, 0x926924},
		{0x20, 0x934924},
		{0x40, 0x9A4924},
		{0x80, 0xD24924},
		{0xFD, 0xDB6DA6},
		{0xFE, 0xDB6DB4},
		{0xFF, 0xDB6DB6},
	}
	for i, line := range data {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if actual := NRZ(line.in); line.expected != actual {
				t.Fatalf("NRZ(%X): 0x%X != 0x%X", line.in, line.expected, actual)
			}
		})
	}
}

func TestRaster(t *testing.T) {
	// Input length must be multiple of 3.
	data := []byte{
		// 24 bits per pixel in RGB
		0, 1, 2,
		0xFD, 0xFE, 0xFF,
	}
	expected := []byte{
		// 72 bits per pixel in GRB
		0x92, 0x49, 0x26, 0x92, 0x49, 0x24, 0x92, 0x49, 0x34,
		0xdb, 0x6d, 0xb4, 0xdb, 0x6d, 0xa6, 0xdb, 0x6d, 0xb6,
	}
	actual := make([]byte, len(expected))
	raster3(actual, data)
	if !bytes.Equal(expected, actual) {
		t.Fatalf("expected %#v != actual %#v", expected, actual)
	}
}
