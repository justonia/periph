// Copyright 2016 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package nrzled is a driver for LEDs WS2811/WS2812/WS2812b and compatible
// devices like SK6812 and UCS1903 that uses a single wire NRZ encoded
// communication protocol.
//
// For LED with a dedicated white channel, often called RGBW or RGBWW (warm
// white), use package sk6812rgbw instead.
//
// Datasheet
//
// This directory contains datasheets for WS2812, WS2812b, UCS190x and various
// SK6812.
//
// https://github.com/cpldcpu/light_ws2812/tree/master/Datasheets
//
// UCS1903 datasheet
//
// http://www.bestlightingbuy.com/pdf/UCS1903%20datasheet.pdf
package nrzled
