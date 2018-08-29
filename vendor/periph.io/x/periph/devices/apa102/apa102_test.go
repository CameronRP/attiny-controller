// Copyright 2016 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package apa102

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"testing"

	"periph.io/x/periph/conn/conntest"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/conn/spi/spitest"
)

func TestRamp(t *testing.T) {
	// Tests a few known values.
	data := []struct {
		input    uint8
		expected uint16
	}{
		{0x00, 0x0000},
		{0x01, 0x0001},
		{0x02, 0x0002},
		{0x03, 0x0003},
		{0x04, 0x0004},
		{0x05, 0x0005},
		{0x06, 0x0006},
		{0x07, 0x0007},
		{0x08, 0x0008},
		{0x09, 0x0009},
		{0x0A, 0x000A},
		{0x0B, 0x000B},
		{0x0C, 0x000C},
		{0x0D, 0x000D},
		{0x0E, 0x000E},
		{0x0F, 0x000F},
		{0x10, 0x0010},
		{0x11, 0x0011},
		{0x12, 0x0012},
		{0x13, 0x0013},
		{0x14, 0x0014},
		{0x15, 0x0015},
		{0x16, 0x0016},
		{0x17, 0x0017},
		{0x18, 0x0018},
		{0x19, 0x0019},
		{0x1A, 0x001A},
		{0x1B, 0x001B},
		{0x1C, 0x001C},
		{0x1D, 0x001D},
		{0x1E, 0x001E},
		{0x1F, 0x001F},
		{0x20, 0x0020},
		{0x21, 0x0021},
		{0x22, 0x0022},
		{0x23, 0x0023},
		{0x24, 0x0024},
		{0x25, 0x0025},
		{0x26, 0x0026},
		{0x27, 0x0027},
		{0x28, 0x0028},
		{0x29, 0x0029},
		{0x2A, 0x002A},
		{0x2B, 0x002B},
		{0x2C, 0x002C},
		{0x2D, 0x002D},
		{0x2E, 0x002E},
		{0x2F, 0x002F},
		{0x30, 0x0030},
		{0x31, 0x0031},
		{0x32, 0x0032},
		{0x33, 0x0033},
		{0x34, 0x0034},
		{0x35, 0x0035},
		{0x36, 0x0036},
		{0x37, 0x0037},
		{0x38, 0x0038},
		{0x39, 0x0039},
		{0x3A, 0x003A},
		{0x3B, 0x003B},
		{0x3C, 0x003C},
		{0x3D, 0x003D},
		{0x3E, 0x003E},
		{0x3F, 0x003F},
		{0x40, 0x0040},
		{0x41, 0x0041},
		{0x42, 0x0042},
		{0x43, 0x0043},
		{0x44, 0x0044},
		{0x45, 0x0045},
		{0x46, 0x0046},
		{0x47, 0x0047},
		{0x48, 0x0048},
		{0x49, 0x0049},
		{0x4A, 0x004A},
		{0x4B, 0x004B},
		{0x4C, 0x004C},
		{0x4D, 0x004D},
		{0x4E, 0x004E},
		{0x4F, 0x004F},
		{0x50, 0x004F},
		{0x51, 0x004F},
		{0x52, 0x004F},
		{0x53, 0x004F},
		{0x54, 0x004F},
		{0x55, 0x004F},
		{0x56, 0x004F},
		{0x57, 0x0050},
		{0x58, 0x0050},
		{0x59, 0x0050},
		{0x5A, 0x0051},
		{0x5B, 0x0051},
		{0x5C, 0x0052},
		{0x5D, 0x0053},
		{0x5E, 0x0054},
		{0x5F, 0x0055},
		{0x60, 0x0056},
		{0x61, 0x0057},
		{0x62, 0x0059},
		{0x63, 0x005A},
		{0x64, 0x005C},
		{0x65, 0x005E},
		{0x66, 0x0060},
		{0x67, 0x0063},
		{0x68, 0x0065},
		{0x69, 0x0068},
		{0x6A, 0x006B},
		{0x6B, 0x006E},
		{0x6C, 0x0072},
		{0x6D, 0x0075},
		{0x6E, 0x0079},
		{0x6F, 0x007E},
		{0x70, 0x0082},
		{0x71, 0x0087},
		{0x72, 0x008C},
		{0x73, 0x0092},
		{0x74, 0x0098},
		{0x75, 0x009E},
		{0x76, 0x00A4},
		{0x77, 0x00AB},
		{0x78, 0x00B2},
		{0x79, 0x00B9},
		{0x7A, 0x00C1},
		{0x7B, 0x00C9},
		{0x7C, 0x00D2},
		{0x7D, 0x00DA},
		{0x7E, 0x00E4},
		{0x7F, 0x00ED},
		{0x80, 0x00F8},
		{0x81, 0x0102},
		{0x82, 0x010D},
		{0x83, 0x0119},
		{0x84, 0x0124},
		{0x85, 0x0131},
		{0x86, 0x013E},
		{0x87, 0x014B},
		{0x88, 0x0159},
		{0x89, 0x0167},
		{0x8A, 0x0176},
		{0x8B, 0x0185},
		{0x8C, 0x0195},
		{0x8D, 0x01A5},
		{0x8E, 0x01B6},
		{0x8F, 0x01C7},
		{0x90, 0x01D9},
		{0x91, 0x01EC},
		{0x92, 0x01FF},
		{0x93, 0x0212},
		{0x94, 0x0226},
		{0x95, 0x023B},
		{0x96, 0x0251},
		{0x97, 0x0267},
		{0x98, 0x027D},
		{0x99, 0x0294},
		{0x9A, 0x02AC},
		{0x9B, 0x02C5},
		{0x9C, 0x02DE},
		{0x9D, 0x02F8},
		{0x9E, 0x0312},
		{0x9F, 0x032E},
		{0xA0, 0x034A},
		{0xA1, 0x0366},
		{0xA2, 0x0384},
		{0xA3, 0x03A2},
		{0xA4, 0x03C0},
		{0xA5, 0x03E0},
		{0xA6, 0x0400},
		{0xA7, 0x0421},
		{0xA8, 0x0443},
		{0xA9, 0x0465},
		{0xAA, 0x0489},
		{0xAB, 0x04AC},
		{0xAC, 0x04D1},
		{0xAD, 0x04F7},
		{0xAE, 0x051D},
		{0xAF, 0x0545},
		{0xB0, 0x056D},
		{0xB1, 0x0596},
		{0xB2, 0x05C0},
		{0xB3, 0x05EA},
		{0xB4, 0x0616},
		{0xB5, 0x0642},
		{0xB6, 0x066F},
		{0xB7, 0x069D},
		{0xB8, 0x06CC},
		{0xB9, 0x06FC},
		{0xBA, 0x072D},
		{0xBB, 0x075F},
		{0xBC, 0x0792},
		{0xBD, 0x07C6},
		{0xBE, 0x07FA},
		{0xBF, 0x0830},
		{0xC0, 0x0866},
		{0xC1, 0x089E},
		{0xC2, 0x08D6},
		{0xC3, 0x090F},
		{0xC4, 0x094A},
		{0xC5, 0x0985},
		{0xC6, 0x09C2},
		{0xC7, 0x09FF},
		{0xC8, 0x0A3E},
		{0xC9, 0x0A7D},
		{0xCA, 0x0ABE},
		{0xCB, 0x0B00},
		{0xCC, 0x0B42},
		{0xCD, 0x0B86},
		{0xCE, 0x0BCB},
		{0xCF, 0x0C11},
		{0xD0, 0x0C58},
		{0xD1, 0x0CA1},
		{0xD2, 0x0CEA},
		{0xD3, 0x0D34},
		{0xD4, 0x0D80},
		{0xD5, 0x0DCD},
		{0xD6, 0x0E1B},
		{0xD7, 0x0E6A},
		{0xD8, 0x0EBA},
		{0xD9, 0x0F0B},
		{0xDA, 0x0F5E},
		{0xDB, 0x0FB2},
		{0xDC, 0x1007},
		{0xDD, 0x105D},
		{0xDE, 0x10B4},
		{0xDF, 0x110D},
		{0xE0, 0x1167},
		{0xE1, 0x11C2},
		{0xE2, 0x121F},
		{0xE3, 0x127C},
		{0xE4, 0x12DB},
		{0xE5, 0x133C},
		{0xE6, 0x139D},
		{0xE7, 0x1400},
		{0xE8, 0x1464},
		{0xE9, 0x14CA},
		{0xEA, 0x1530},
		{0xEB, 0x1599},
		{0xEC, 0x1602},
		{0xED, 0x166D},
		{0xEE, 0x16D9},
		{0xEF, 0x1747},
		{0xF0, 0x17B6},
		{0xF1, 0x1826},
		{0xF2, 0x1898},
		{0xF3, 0x190B},
		{0xF4, 0x197F},
		{0xF5, 0x19F5},
		{0xF6, 0x1A6D},
		{0xF7, 0x1AE5},
		{0xF8, 0x1B60},
		{0xF9, 0x1BDB},
		{0xFA, 0x1C58},
		{0xFB, 0x1CD7},
		{0xFC, 0x1D57},
		{0xFD, 0x1DD9},
		{0xFE, 0x1E5C},
		{0xFF, 0x1EE1},
	}
	if false {
		for i := 0; i <= 255; i++ {
			fmt.Printf("{0x%02X, 0x%04X},\n", i, ramp(uint8(i), maxOut))
		}
	}
	for i, line := range data {
		if i != int(line.input) || line.expected != ramp(line.input, maxOut) {
			t.Fail()
		}
	}
	if 0x00 != ramp(0x00, 0xFF) {
		t.Fail()
	}
	if 0x21 != ramp(0x7F, 0xFF) {
		t.Fail()
	}
	if 0xFF != ramp(0xFF, 0xFF) {
		t.Fail()
	}
}

func TestRampMonotonic(t *testing.T) {
	// Ensures the ramp is 100% monotonically increasing and without bumps.
	lastValue := uint16(0)
	lastDelta := uint16(0)
	for in := uint32(0); in <= 255; in++ {
		out := ramp(uint8(in), maxOut)
		if out < lastValue {
			t.Fatalf("f(%d) = %d; f(%d) = %d", in-1, lastValue, in, out)
		}
		if out > maxOut {
			t.Fatalf("f(%d) = %d", in, out)
		}

		if out-lastValue+1 < lastDelta {
			t.Errorf("f(%d)=%d  f(%d)=%d  f(%d)=%d  Deltas: '%d+1 < %d' but should be '>='",
				in-2, ramp(uint8(in-2), maxOut), in-1, ramp(uint8(in-1), maxOut), in, ramp(uint8(in), maxOut), out-lastValue, lastDelta)
		}
		lastDelta = out - lastValue
		lastValue = out
	}
}

func TestDevEmpty(t *testing.T) {
	buf := bytes.Buffer{}
	o := DefaultOpts
	o.NumPixels = 0
	d, _ := New(spitest.NewRecordRaw(&buf), &o)
	if n, err := d.Write([]byte{}); n != 0 || err != nil {
		t.Fatalf("%d %v", n, err)
	}
	if expected := []byte{0x0, 0x0, 0x0, 0x0, 0xFF}; !bytes.Equal(expected, buf.Bytes()) {
		t.Fatalf("\ngot:  %#02v\nwant: %#02v\n", buf.Bytes(), expected)
	}
	if s := d.String(); s != "APA102{I:255, T:5000K, GPWM:true, 0LEDs, recordraw}" {
		t.Fatal(s)
	}
}

func TestConnectFail(t *testing.T) {
	if d, err := New(&configFail{}, &DefaultOpts); d != nil || err == nil {
		t.Fatal("Connect() call have failed")
	}
}

func TestDevLen(t *testing.T) {
	buf := bytes.Buffer{}
	o := DefaultOpts
	o.NumPixels = 1
	d, _ := New(spitest.NewRecordRaw(&buf), &o)
	if n, err := d.Write([]byte{0}); n != 0 || err == nil {
		t.Fatalf("%d %v", n, err)
	}
	if expected := []byte{}; !bytes.Equal(expected, buf.Bytes()) {
		t.Fatalf("\ngot:  %#02v\nwant: %#02v\n", buf.Bytes(), expected)
	}
}

var writeTests = []struct {
	name   string
	pixels []byte
	want   []byte
	opts   Opts
}{
	{
		name: "Temperature",
		pixels: ToRGB([]color.NRGBA{
			{0xFF, 0xFF, 0xFF, 0x00},
			{0xFE, 0xFE, 0xFE, 0x00},
			{0xF0, 0xF0, 0xF0, 0x00},
			{0x80, 0x80, 0x80, 0x00},
			{0x80, 0x00, 0x00, 0x00},
			{0x00, 0x80, 0x00, 0x00},
			{0x00, 0x00, 0x80, 0x00},
			{0x00, 0x00, 0x10, 0x00},
			{0x00, 0x00, 0x01, 0x00},
			{0x00, 0x00, 0x00, 0x00},
		}),
		want: []byte{
			0x00, 0x00, 0x00, 0x00,
			0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFB, 0xFB, 0xFB,
			0xFF, 0xC4, 0xC4, 0xC4,
			0xE1, 0xF8, 0xF8, 0xF8,
			0xE1, 0x00, 0x00, 0xF8,
			0xE1, 0x00, 0xF8, 0x00,
			0xE1, 0xF8, 0x00, 0x00,
			0xE1, 0x10, 0x00, 0x00,
			0xE1, 0x01, 0x00, 0x00,
			0xE1, 0x00, 0x00, 0x00,
			0xFF,
		},
		opts: Opts{
			Intensity:   255,
			Temperature: NeutralTemp,
		},
	},
	{
		name: "Intensity",
		pixels: ToRGB([]color.NRGBA{
			{0xFF, 0xFF, 0xFF, 0x00},
			{0xFE, 0xFE, 0xFE, 0x00},
			{0xF0, 0xF0, 0xF0, 0x00},
			{0x80, 0x80, 0x80, 0x00},
			{0x80, 0x00, 0x00, 0x00},
			{0x00, 0x80, 0x00, 0x00},
			{0x00, 0x00, 0x80, 0x00},
			{0x00, 0x00, 0x10, 0x00},
			{0x00, 0x00, 0x01, 0x00},
			{0x00, 0x00, 0x00, 0x00},
		}),
		want: []byte{
			0x00, 0x00, 0x00, 0x00,
			0xFF, 0x7F, 0x7F, 0x7F,
			0xFF, 0x7D, 0x7D, 0x7D,
			0xFF, 0x67, 0x67, 0x67,
			0xE2, 0x9B, 0x9B, 0x9B,
			0xE2, 0x00, 0x00, 0x9B,
			0xE2, 0x00, 0x9B, 0x00,
			0xE2, 0x9B, 0x00, 0x00,
			0xE1, 0x10, 0x00, 0x00,
			0xE1, 0x01, 0x00, 0x00,
			0xE1, 0x00, 0x00, 0x00,
			0xFF,
		},
		opts: Opts{
			Intensity:   127,
			Temperature: NeutralTemp,
		},
	},
	{
		name: "PassThru",
		pixels: ToRGB([]color.NRGBA{
			{0xFF, 0xFF, 0xFF, 0x00},
			{0xFE, 0xFE, 0xFE, 0x00},
			{0xF0, 0xF0, 0xF0, 0x00},
			{0x80, 0x80, 0x80, 0x00},
			{0x80, 0x00, 0x00, 0x00},
			{0x00, 0x80, 0x00, 0x00},
			{0x00, 0x00, 0x80, 0x00},
			{0x00, 0x00, 0x10, 0x00},
			{0x00, 0x00, 0x01, 0x00},
			{0x00, 0x00, 0x00, 0x00},
		}),
		want: []byte{
			0x00, 0x00, 0x00, 0x00,
			0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFE, 0xFE, 0xFE,
			0xFF, 0xF0, 0xF0, 0xF0,
			0xFF, 0x80, 0x80, 0x80,
			0xFF, 0x00, 0x00, 0x80,
			0xFF, 0x00, 0x80, 0x00,
			0xFF, 0x80, 0x00, 0x00,
			0xFF, 0x10, 0x00, 0x00,
			0xFF, 0x01, 0x00, 0x00,
			0xFF, 0x00, 0x00, 0x00,
			0xFF,
		},
		opts: PassThruOpts,
	},
	{
		name: "DisableGlobalPWM - Intensity",
		pixels: ToRGB([]color.NRGBA{
			{0xFF, 0xFF, 0xFF, 0x00},
			{0xFE, 0xFE, 0xFE, 0x00},
			{0xF0, 0xF0, 0xF0, 0x00},
			{0x80, 0x80, 0x80, 0x00},
			{0x80, 0x00, 0x00, 0x00},
			{0x00, 0x80, 0x00, 0x00},
			{0x00, 0x00, 0x80, 0x00},
			{0x00, 0x00, 0x10, 0x00},
			{0x00, 0x00, 0x01, 0x00},
			{0x00, 0x00, 0x00, 0x00},
		}),
		want: []byte{
			0x00, 0x00, 0x00, 0x00,
			0xff, 0x40, 0x40, 0x40,
			0xff, 0x40, 0x40, 0x40,
			0xff, 0x3c, 0x3c, 0x3c,
			0xff, 0x20, 0x20, 0x20,
			0xff, 0x00, 0x00, 0x20,
			0xff, 0x00, 0x20, 0x00,
			0xff, 0x20, 0x00, 0x00,
			0xff, 0x04, 0x00, 0x00,
			0xff, 0x00, 0x00, 0x00,
			0xff, 0x00, 0x00, 0x00,
			0xff,
		},
		opts: Opts{
			Intensity:        64,
			Temperature:      NeutralTemp,
			DisableGlobalPWM: true,
		},
	},
	{
		name: "DisableGlobalPWM - Temp",
		pixels: ToRGB([]color.NRGBA{
			{0xFF, 0xFF, 0xFF, 0x00},
			{0xFE, 0xFE, 0xFE, 0x00},
			{0xF0, 0xF0, 0xF0, 0x00},
			{0x80, 0x80, 0x80, 0x00},
			{0x80, 0x00, 0x00, 0x00},
			{0x00, 0x80, 0x00, 0x00},
			{0x00, 0x00, 0x80, 0x00},
			{0x00, 0x00, 0x10, 0x00},
			{0x00, 0x00, 0x01, 0x00},
			{0x00, 0x00, 0x00, 0x00},
		}),
		want: []byte{
			0x00, 0x00, 0x00, 0x00,
			0xff, 0xd5, 0xe8, 0xff,
			0xff, 0xd4, 0xe7, 0xfe,
			0xff, 0xc8, 0xda, 0xf0,
			0xff, 0x6b, 0x74, 0x80,
			0xff, 0x00, 0x00, 0x80,
			0xff, 0x00, 0x74, 0x00,
			0xff, 0x6b, 0x00, 0x00,
			0xff, 0x0d, 0x00, 0x00,
			0xff, 0x01, 0x00, 0x00,
			0xff, 0x00, 0x00, 0x00,
			0xff,
		},
		opts: Opts{
			Intensity:        255,
			Temperature:      5000,
			DisableGlobalPWM: true,
		},
	},
	{
		name: "Intensity and temperature",
		pixels: func() []byte {
			var p []byte
			for i := 0; i < 16*3; i++ {
				p = append(p, uint8(i<<2))
			}
			return p
		}(),
		want: expectedi250t5000,
		opts: Opts{
			Intensity:   250,
			Temperature: 5000,
		},
	},
}

func TestWrites(t *testing.T) {
	for _, tt := range writeTests {
		buf := bytes.Buffer{}
		tt.opts.NumPixels = len(tt.pixels) / 3
		d, _ := New(spitest.NewRecordRaw(&buf), &tt.opts)
		n, err := d.Write(tt.pixels)
		if err != nil {
			t.Fatal(err)
		}
		if n != len(tt.pixels) {
			t.Fatalf("%s: Got %d bytes result, want %d", tt.name, n, len(tt.pixels)*3)
		}
		if !bytes.Equal(buf.Bytes(), tt.want) {
			t.Fatalf("%s:\ngot:  %#02v\nwant: %#02v\n", tt.name, buf.Bytes(), tt.want)
		}
	}
}

func TestDevColor(t *testing.T) {
	if (&Dev{}).ColorModel() != color.NRGBAModel {
		t.Fail()
	}
}

func TestDevLong(t *testing.T) {
	buf := bytes.Buffer{}
	colors := make([]color.NRGBA, 256)
	o := DefaultOpts
	o.NumPixels = len(colors)
	d, _ := New(spitest.NewRecordRaw(&buf), &o)
	if n, err := d.Write(ToRGB(colors)); n != len(colors)*3 || err != nil {
		t.Fatalf("%d %v", n, err)
	}
	expected := make([]byte, 4*(256+1)+17)
	for i := 0; i < 256; i++ {
		expected[4+4*i] = 0xE1
	}
	trailer := expected[4*257:]
	for i := range trailer {
		trailer[i] = 0xFF
	}
	if !bytes.Equal(expected, buf.Bytes()) {
		t.Fatalf("\ngot:  %#02v\nwant: %#02v\n", buf.Bytes(), expected)
	}
}

func TestDevWrite_Long(t *testing.T) {
	buf := bytes.Buffer{}
	o := DefaultOpts
	o.NumPixels = 1
	d, _ := New(spitest.NewRecordRaw(&buf), &o)
	if n, err := d.Write([]byte{0, 0, 0, 1, 1, 1}); n != 0 || err == nil {
		t.Fatal(n, err)
	}
}

// expectedi250t5000 is the expected output for multiple test cases. Each test case
// use a completely different code path so make sure each code path results in
// the exact same output.
var expectedi250t5000 = []byte{
	0x00, 0x00, 0x00, 0x00, 0xE1, 0x08, 0x04, 0x00, 0xE1, 0x14, 0x10, 0xC, 0xE1,
	0x20, 0x1C, 0x18, 0xE1, 0x2C, 0x28, 0x24, 0xE1, 0x38, 0x34, 0x30, 0xE1, 0x41,
	0x40, 0x3C, 0xE1, 0x44, 0x47, 0x48, 0xE1, 0x53, 0x4C, 0x4E, 0xE1, 0x78, 0x62,
	0x56, 0xE1, 0xBD, 0x94, 0x73, 0xE2, 0x95, 0x77, 0x5A, 0xE2, 0xE5, 0xBD, 0x94,
	0xE4, 0xAA, 0x92, 0x77, 0xE4, 0xF3, 0xD7, 0xB8, 0xFF, 0x2B, 0x28, 0x23, 0xFF,
	0x3A, 0x36, 0x32, 0xFF, 0xFF,
}

// expectedi250t6500 is the default color temperature.
var expectedi250t6500 = []byte{
	0x00, 0x00, 0x00, 0x00, 0xE1, 0x08, 0x04, 0x00, 0xE1, 0x14, 0x10, 0x0C, 0xE1,
	0x20, 0x1C, 0x18, 0xE1, 0x2C, 0x28, 0x24, 0xE1, 0x38, 0x34, 0x30, 0xE1, 0x44,
	0x40, 0x3C, 0xE1, 0x4E, 0x4C, 0x48, 0xE1, 0x52, 0x4F, 0x4E, 0xE1, 0x66, 0x5C,
	0x56, 0xE1, 0x9A, 0x84, 0x73, 0xE1, 0xFB, 0xD4, 0xB4, 0xE2, 0xCB, 0xAE, 0x94,
	0xE4, 0xA0, 0x8A, 0x77, 0xE4, 0xF0, 0xD2, 0xB8, 0xFF, 0x2D, 0x28, 0x23, 0xFF,
	0x3E, 0x38, 0x32, 0xFF, 0xFF,
}

// expectedi250raw is using DisableGlobalPWM = true.
var expectedi250raw = []byte{
	0x00, 0x00, 0x00, 0x00, 0xFF, 0x08, 0x04, 0x00, 0xFF, 0x14, 0x10, 0x0C, 0xFF,
	0x1F, 0x1B, 0x18, 0xFF, 0x2B, 0x27, 0x23, 0xFF, 0x37, 0x33, 0x2F, 0xFF, 0x43,
	0x3F, 0x3B, 0xFF, 0x4E, 0x4B, 0x47, 0xFF, 0x5A, 0x56, 0x52, 0xFF, 0x66, 0x62,
	0x5E, 0xFF, 0x72, 0x6E, 0x6A, 0xFF, 0x7D, 0x7A, 0x76, 0xFF, 0x89, 0x85, 0x81,
	0xFF, 0x95, 0x91, 0x8D, 0xFF, 0xA1, 0x9D, 0x99, 0xFF, 0xAD, 0xA9, 0xA5, 0xFF,
	0xB8, 0xB4, 0xB0, 0xFF, 0xFF,
}

var drawTests = []struct {
	name string
	img  image.Image
	want []byte
	opts Opts
}{
	{
		name: "Draw NRGBA",
		img: func() image.Image {
			im := image.NewNRGBA(image.Rect(0, 0, 16, 1))
			for i := 0; i < 16; i++ {
				// Test all intensity code paths. Confirm that alpha is ignored.
				im.Pix[4*i] = uint8((3 * i) << 2)
				im.Pix[4*i+1] = uint8((3*i + 1) << 2)
				im.Pix[4*i+2] = uint8((3*i + 2) << 2)
				im.Pix[4*i+3] = 0
			}
			return im
		}(),
		want: expectedi250t5000,
		opts: Opts{
			NumPixels:   16,
			Intensity:   250,
			Temperature: 5000,
		},
	},
	{
		name: "Draw NRGBA Wide",
		img: func() image.Image {
			im := image.NewNRGBA(image.Rect(0, 0, 17, 2))
			for x := 0; x < 16; x++ {
				// Test all intensity code paths. Confirm that alpha is ignored.
				im.SetNRGBA(x, 0, color.NRGBA{uint8((3 * x) << 2), uint8((3*x + 1) << 2), uint8((3*x + 2) << 2), 0})
			}
			return im
		}(),
		want: expectedi250t6500,
		opts: Opts{
			NumPixels:   16,
			Intensity:   250,
			Temperature: NeutralTemp,
		},
	},
	{
		name: "Draw NRGBA no global PWM",
		img: func() image.Image {
			im := image.NewNRGBA(image.Rect(0, 0, 16, 1))
			for i := 0; i < 16; i++ {
				// Test all intensity code paths. Confirm that alpha is ignored.
				im.Pix[4*i] = uint8((3 * i) << 2)
				im.Pix[4*i+1] = uint8((3*i + 1) << 2)
				im.Pix[4*i+2] = uint8((3*i + 2) << 2)
				im.Pix[4*i+3] = 0
			}
			return im
		}(),
		want: expectedi250raw,
		opts: Opts{
			NumPixels:        16,
			Temperature:      NeutralTemp,
			Intensity:        250,
			DisableGlobalPWM: true,
		},
	},
	{
		name: "Draw RGBA",
		img: func() image.Image {
			im := image.NewRGBA(image.Rect(0, 0, 16, 1))
			for i := 0; i < 16; i++ {
				im.Pix[4*i] = uint8((3 * i) << 2)
				im.Pix[4*i+1] = uint8((3*i + 1) << 2)
				im.Pix[4*i+2] = uint8((3*i + 2) << 2)
				im.Pix[4*i+3] = 0xFF
			}
			return im
		}(),
		want: expectedi250t5000,
		opts: Opts{
			NumPixels:   16,
			Intensity:   250,
			Temperature: 5000,
		},
	},
	{
		// Just something that doesn't have a fast path
		name: "Draw NRGBA64",
		img: func() image.Image {
			im := image.NewNRGBA64(image.Rect(0, 0, 16, 1))
			for i := 0; i < 16; i++ {
				im.Set(i, 0, color.NRGBA64{
					R: uint16(((3 * i) << 10)),
					G: uint16(((3*i + 1) << 10)),
					B: uint16(((3*i + 2) << 10)),
					A: 0xFFFF,
				})
			}
			return im
		}(),
		want: expectedi250t5000,
		opts: Opts{
			NumPixels:   16,
			Intensity:   250,
			Temperature: 5000,
		},
	},
	{
		name: "Draw RGBA no global PWM",
		img: func() image.Image {
			im := image.NewRGBA(image.Rect(0, 0, 16, 1))
			for i := 0; i < 16; i++ {
				im.Pix[4*i] = uint8((3 * i) << 2)
				im.Pix[4*i+1] = uint8((3*i + 1) << 2)
				im.Pix[4*i+2] = uint8((3*i + 2) << 2)
				im.Pix[4*i+3] = 0xFF
			}
			return im
		}(),
		want: expectedi250raw,
		opts: Opts{
			NumPixels:        16,
			Temperature:      NeutralTemp,
			Intensity:        250,
			DisableGlobalPWM: true,
		},
	},
}

func TestDraws(t *testing.T) {
	for _, tt := range drawTests {
		buf := bytes.Buffer{}
		d, _ := New(spitest.NewRecordRaw(&buf), &tt.opts)
		if err := d.Draw(d.Bounds(), tt.img, image.Point{}); err != nil {
			t.Fatalf("%s: %v", tt.name, err)
		}
		if !bytes.Equal(buf.Bytes(), tt.want) {
			t.Fatalf("%s:\ngot:  %#02v\nwant: %#02v\n", tt.name, buf.Bytes(), tt.want)
		}
	}
}

var offsetDrawWant = []byte{
	0x00, 0x00, 0x00, 0x00,
	0xE1, 0x89, 0x79, 0x6B,
	0xE1, 0x9A, 0x88, 0x75,
	0xE1, 0xAD, 0x98, 0x82,
	0xE1, 0xC2, 0xAB, 0x92,
	0xE1, 0xDA, 0xC0, 0xA4,
	0xE1, 0xF5, 0xD9, 0xB9,
	0xE2, 0x89, 0x7A, 0x69,
	0xE2, 0x9A, 0x8A, 0x76,
	0xE2, 0xAC, 0x9B, 0x86,
	0xE2, 0xC0, 0xAE, 0x98,
	0xE2, 0xD5, 0xC3, 0xAC,
	0xE2, 0xED, 0xDA, 0xC2,
	0xE4, 0x83, 0x7A, 0x6E,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0xFF,
}

var offsetDrawTests = []struct {
	name   string
	img    image.Image
	point  image.Point
	offset image.Rectangle
	want   []byte
	opts   Opts
}{
	{
		name: "Offset Draw NRGBA",
		img: func() image.Image {
			im := image.NewNRGBA(image.Rect(0, 0, 16, 4))
			for x := 0; x < 16; x++ {
				for y := 0; y < 4; y++ {
					i := (y*16 + x) * 3
					im.Set(x, y, color.RGBA{R: uint8(i + 1), G: uint8(i + 2), B: uint8(i + 3), A: 0xFF})
				}
			}
			return im
		}(),
		point:  image.Point{X: 3, Y: 2},
		offset: image.Rect(0, 0, 16, 1),
		want:   offsetDrawWant,
		opts: Opts{
			NumPixels:   15,
			Intensity:   255,
			Temperature: 5000,
		},
	},
	{
		name: "Both Offset Draw NRGBA",
		img: func() image.Image {
			im := image.NewNRGBA(image.Rect(0, 0, 16, 4))
			for x := 0; x < 16; x++ {
				for y := 0; y < 4; y++ {
					i := (y*16 + x) * 3
					im.Set(x, y, color.RGBA{R: uint8(i + 1), G: uint8(i + 2), B: uint8(i + 3), A: 0xFF})
				}
			}
			return im
		}(),
		point:  image.Point{X: 3, Y: 2},
		offset: image.Rect(2, 0, 16, 1),
		want: []byte{
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0xE1, 0x89, 0x79, 0x6B,
			0xE1, 0x9A, 0x88, 0x75,
			0xE1, 0xAD, 0x98, 0x82,
			0xE1, 0xC2, 0xAB, 0x92,
			0xE1, 0xDA, 0xC0, 0xA4,
			0xE1, 0xF5, 0xD9, 0xB9,
			0xE2, 0x89, 0x7A, 0x69,
			0xE2, 0x9A, 0x8A, 0x76,
			0xE2, 0xAC, 0x9B, 0x86,
			0xE2, 0xC0, 0xAE, 0x98,
			0xE2, 0xD5, 0xC3, 0xAC,
			0xE2, 0xED, 0xDA, 0xC2,
			0xE4, 0x83, 0x7A, 0x6E,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0xFF, 0xFF,
		},
		opts: Opts{
			NumPixels:   17,
			Intensity:   255,
			Temperature: 5000,
		},
	},
}

func TestOffsetDraws(t *testing.T) {
	for _, tt := range offsetDrawTests {
		buf := bytes.Buffer{}
		d, _ := New(spitest.NewRecordRaw(&buf), &tt.opts)
		if err := d.Draw(tt.offset, tt.img, tt.point); err != nil {
			t.Fatalf("%s: %v", tt.name, err)
		}
		if !bytes.Equal(buf.Bytes(), tt.want) {
			t.Fatalf("%s:\ngot:  %#02v\nwant: %#02v\n", tt.name, buf.Bytes(), tt.want)
		}
	}
}

func TestHalt(t *testing.T) {
	s := spitest.Playback{
		Playback: conntest.Playback{
			Ops: []conntest.IO{
				{W: []byte{0x0, 0x0, 0x0, 0x0, 0xe1, 0x0, 0x0, 0x0, 0xe1, 0x0, 0x0, 0x0, 0xe1, 0x0, 0x0, 0x0, 0xe1, 0x0, 0x0, 0x0, 0xff}},
			},
		},
	}
	o := DefaultOpts
	o.NumPixels = 4
	o.Temperature = 5000
	d, _ := New(&s, &o)
	if err := d.Halt(); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestInit(t *testing.T) {
	// Catch the "maxB == maxG" line.
	l := lut{}
	l.init(255, 6000, true)
	if equalUint16(l.r[:], l.g[:]) || !equalUint16(l.g[:], l.b[:]) {
		t.Fatal("test case is for only when maxG == maxB but maxR != maxG")
	}
}

//

type genColor func(int) [3]byte

func benchmarkWrite(b *testing.B, o Opts, length int, f genColor) {
	var pixels []byte
	for i := 0; i < length; i++ {
		c := f(i)
		pixels = append(pixels, c[:]...)
	}
	o.NumPixels = length
	b.ReportAllocs()
	d, _ := New(spitest.NewRecordRaw(ioutil.Discard), &o)
	_, _ = d.Write(pixels[:])
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = d.Write(pixels[:])
	}
}

func BenchmarkWriteWhite(b *testing.B) {
	o := DefaultOpts
	o.Intensity = 250
	benchmarkWrite(b, o, 150, func(i int) [3]byte { return [3]byte{0xFF, 0xFF, 0xFF} })
}

func BenchmarkWriteDim(b *testing.B) {
	o := DefaultOpts
	o.Intensity = 250
	benchmarkWrite(b, o, 150, func(i int) [3]byte { return [3]byte{0x01, 0x01, 0x01} })
}

func BenchmarkWriteBlack(b *testing.B) {
	o := DefaultOpts
	o.Intensity = 250
	benchmarkWrite(b, o, 150, func(i int) [3]byte { return [3]byte{0x0, 0x0, 0x0} })
}

func genColorfulPixel(x int) [3]byte {
	i := x * 3
	return [3]byte{uint8(i) + uint8(i>>8),
		uint8(i+1) + uint8(i+1>>8),
		uint8(i+2) + uint8(i+2>>8),
	}
}

func BenchmarkWriteColorful(b *testing.B) {
	o := DefaultOpts
	o.Intensity = 250
	o.Temperature = 5000
	benchmarkWrite(b, o, 150, genColorfulPixel)
}

func BenchmarkWriteColorfulPassThru(b *testing.B) {
	o := PassThruOpts
	o.Intensity = 250
	benchmarkWrite(b, o, 150, genColorfulPixel)
}

func BenchmarkWriteColorfulVariation(b *testing.B) {
	// Continuously vary the lookup tables.
	b.ReportAllocs()
	pixels := [256 * 3]byte{}
	for i := range pixels {
		pixels[i] = uint8(i) + uint8(i>>8)
	}
	o := DefaultOpts
	o.NumPixels = len(pixels) / 3
	o.Intensity = 250
	o.Temperature = 5000
	d, _ := New(spitest.NewRecordRaw(ioutil.Discard), &o)
	_, _ = d.Write(pixels[:])
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Intensity = uint8(i)
		d.Temperature = uint16((3000 + i) & 0x1FFF)
		_, _ = d.Write(pixels[:])
	}
}

func benchmarkDraw(b *testing.B, o Opts, img draw.Image, f genColor) {
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			pix := f(x)
			c := color.NRGBA{R: pix[0], G: pix[1], B: pix[2], A: 255}
			img.Set(x, y, c)
		}
	}
	o.NumPixels = img.Bounds().Max.X
	b.ReportAllocs()
	d, _ := New(spitest.NewRecordRaw(ioutil.Discard), &o)
	r := d.Bounds()
	p := image.Point{}
	if err := d.Draw(r, img, p); err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := d.Draw(r, img, p); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDrawNRGBAColorful(b *testing.B) {
	o := DefaultOpts
	o.Intensity = 250
	o.Temperature = 5000
	benchmarkDraw(b, o, image.NewNRGBA(image.Rect(0, 0, 150, 1)), genColorfulPixel)
}

func BenchmarkDrawNRGBAColorfulPassThru(b *testing.B) {
	o := PassThruOpts
	o.Intensity = 250
	benchmarkDraw(b, o, image.NewNRGBA(image.Rect(0, 0, 150, 1)), genColorfulPixel)
}

func BenchmarkDrawNRGBAWhite(b *testing.B) {
	o := DefaultOpts
	o.Intensity = 250
	o.Temperature = 5000
	benchmarkDraw(b, o, image.NewNRGBA(image.Rect(0, 0, 150, 1)), func(i int) [3]byte { return [3]byte{0xFF, 0xFF, 0xFF} })
}

func BenchmarkDrawRGBAColorful(b *testing.B) {
	o := DefaultOpts
	o.Intensity = 250
	o.Temperature = 5000
	benchmarkDraw(b, o, image.NewRGBA(image.Rect(0, 0, 256, 1)), genColorfulPixel)
}

func BenchmarkDrawRGBAColorfulPassThru(b *testing.B) {
	o := PassThruOpts
	o.Intensity = 250
	benchmarkDraw(b, o, image.NewRGBA(image.Rect(0, 0, 256, 1)), genColorfulPixel)
}

func BenchmarkDrawSlowpath(b *testing.B) {
	// Should be an image type that doesn't have a fast path
	img := image.NewCMYK(image.Rect(0, 0, 150, 1))
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			pix := genColorfulPixel(x)
			c := color.CMYK{C: pix[0], M: pix[1], Y: pix[2], K: 0xFF}
			img.Set(x, y, c)
		}
	}
	o := DefaultOpts
	o.NumPixels = img.Bounds().Max.X
	b.ReportAllocs()
	d, _ := New(spitest.NewRecordRaw(ioutil.Discard), &o)
	r := d.Bounds()
	p := image.Point{}
	if err := d.Draw(r, img, p); err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := d.Draw(r, img, p); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHalt(b *testing.B) {
	b.ReportAllocs()
	d, _ := New(spitest.NewRecordRaw(ioutil.Discard), &DefaultOpts)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Halt()
	}
}

//

type configFail struct {
	spitest.Record
}

func (c *configFail) Connect(f physic.Frequency, mode spi.Mode, bits int) (spi.Conn, error) {
	return nil, errors.New("injected error")
}

func equalUint16(a, b []uint16) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}