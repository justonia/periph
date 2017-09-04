// Copyright 2017 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package nrzled

import (
	"errors"
	"image"
	"image/color"
	"time"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpiostream"
	"periph.io/x/periph/devices"
)

// NRZ converts a 8 bit byte into the NRZ encoded 24 bits.
func NRZ(b byte) uint32 {
	// The stream is 1x01x01x01x01x01x01x01x0 with the x bits being the bits from
	// `b` in reverse order.
	out := uint32(0x924924)
	out |= uint32(b&0x80) << (3*7 + 1 - 7)
	out |= uint32(b&0x40) << (3*6 + 1 - 6)
	out |= uint32(b&0x20) << (3*5 + 1 - 5)
	out |= uint32(b&0x10) << (3*4 + 1 - 4)
	out |= uint32(b&0x08) << (3*3 + 1 - 3)
	out |= uint32(b&0x04) << (3*2 + 1 - 2)
	out |= uint32(b&0x02) << (3*1 + 1 - 1)
	out |= uint32(b&0x01) << (3*0 + 1 - 0)
	return out
}

// Dev is a handle to the LED strip.
type Dev struct {
	//p         gpiostream.PinStreamer
	numLights      int
	dedicatedWhite bool
	b              gpiostream.BitStream
}

// Halt turns the lights off.
func (d *Dev) Halt() error {
	return errors.New("nrzled: not implemented")
}

// ColorModel implements devices.Display. There's no surprise, it is
// color.NRGBAModel.
func (d *Dev) ColorModel() color.Model {
	return color.NRGBAModel
}

// Bounds implements devices.Display. Min is guaranteed to be {0, 0}.
func (d *Dev) Bounds() image.Rectangle {
	return image.Rectangle{Max: image.Point{X: d.numLights, Y: 1}}
}

// Draw implements devices.Display.
//
// Using something else than image.NRGBA is 10x slower and is not recommended.
// When using image.NRGBA, the alpha channel is ignored in RGB mode and used as
// White channel in RGBW mode.
func (d *Dev) Draw(r image.Rectangle, src image.Image, sp image.Point) {
	r = r.Intersect(d.Bounds())
	srcR := src.Bounds()
	srcR.Min = srcR.Min.Add(sp)
	if dX := r.Dx(); dX < srcR.Dx() {
		srcR.Max.X = srcR.Min.X + dX
	}
	if dY := r.Dy(); dY < srcR.Dy() {
		srcR.Max.Y = srcR.Min.Y + dY
	}
	// TODO(maruel): Uh?
	//rasterImg(d.buf, r, src, srcR)
	//_, _ = d.s.Write(d.buf)
}

// Write accepts a stream of raw RGB/RGBW pixels and sends it as NRZ encoded
// stream.
func (d *Dev) Write(pixels []byte) (int, error) {
	if len(pixels)%3 != 0 {
		return 0, errLength
	}
	raster3(d.b.Bits, pixels)
	//raster4(d.b.Bits, pixels)
	/*
		if err := d.p.Stream(&d.b); err != nil {
			return 0, err
		}
		return len(pixels), nil
	*/
	return 0, errors.New("nrzled: not implemented")
}

// New opens a handle to a WS2811/WS2812/WS2812b/SK6812.
//
// For WS2812b and SK6812, speed should be 800000, for others, speed should be
// 400000.
func New(p gpio.PinOut, numLights, speed int, dedicatedWhite bool) (*Dev, error) {
	/*
		s, ok := p.(gpiostream.PinStreamer)
		if !ok {
			return nil, errors.New("sk6812rgbw: pin must implement gpiostream.PinStreamer")
		}
	*/
	if speed == 0 {
		return nil, errors.New("nrzled: specify the speed")
	}
	// It is more space effective to use gpiostream.Bits than
	// gpiostream.EdgeStream.
	return &Dev{
		//p:         s,
		numLights:      numLights,
		dedicatedWhite: dedicatedWhite,
		b: gpiostream.BitStream{
			Res: time.Second / time.Duration(speed),
			// Each bit is encoded on 3 bits.
			// Each LED is 24 bits, stored in 8 bits integers.
			Bits: make(gpiostream.Bits, numLights*3*3),
		},
	}, nil
}

//

var errLength = errors.New("nrzled: invalid RGB stream length")

// raster3 converts a RGB input stream into a binary output stream as it must be
// sent over the GPIO pin.
//
// `in` is RGB 24 bits. Each bit is encoded over 3 bits so the length of `out`
// must be 3x as large as `in`.
//
// The encoding is NRZ: https://en.wikipedia.org/wiki/Non-return-to-zero
func raster3(out, in []byte) {
	for i := 0; i < len(in); i += 3 {
		// Input is RGB in 24 bits.
		// Encoded output format is GRB as 72 bits (24 * 3).
		g := NRZ(in[i+1])
		out[3*(i+0)+0] = byte(g >> 16)
		out[3*(i+0)+1] = byte(g >> 8)
		out[3*(i+0)+2] = byte(g)
		r := NRZ(in[i+0])
		out[3*(i+1)+0] = byte(r >> 16)
		out[3*(i+1)+1] = byte(r >> 8)
		out[3*(i+1)+2] = byte(r)
		b := NRZ(in[i+2])
		out[3*(i+2)+0] = byte(b >> 16)
		out[3*(i+2)+1] = byte(b >> 8)
		out[3*(i+2)+2] = byte(b)
	}
}

// raster4 converts a RGBW input stream into a binary output stream as it must
// be sent over the GPIO pin.
//
// `in` is RGB 32 bits. Each bit is encoded over 3 bits so the length of `out`
// must be 3x as large as `in`.
//
// The encoding is NRZ: https://en.wikipedia.org/wiki/Non-return-to-zero
func raster4(out, in []byte) {
	for i := 0; i < len(in); i += 4 {
		// Input is RGBW in 24 bits.
		// Encoded output format is GRBW as 96 bits (24 * 4).
		g := NRZ(in[i+1])
		out[4*(i+0)+0] = byte(g >> 16)
		out[4*(i+0)+1] = byte(g >> 8)
		out[4*(i+0)+2] = byte(g)
		r := NRZ(in[i+0])
		out[4*(i+1)+0] = byte(r >> 16)
		out[4*(i+1)+1] = byte(r >> 8)
		out[4*(i+1)+2] = byte(r)
		b := NRZ(in[i+2])
		out[4*(i+2)+0] = byte(b >> 16)
		out[4*(i+2)+1] = byte(b >> 8)
		out[4*(i+2)+2] = byte(b)
		w := NRZ(in[i+3])
		out[4*(i+3)+0] = byte(w >> 16)
		out[4*(i+3)+1] = byte(w >> 8)
		out[4*(i+3)+2] = byte(w)
	}
}

var _ devices.Display = &Dev{}
