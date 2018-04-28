package goturbopfor

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestDecode(t *testing.T) {
	for _, test := range []struct {
		name  string
		input []byte
		want  []uint32
	}{
		{
			name:  "bitpack only",
			input: []byte{0x07, 0xaa, 0x9c, 0xf6, 0x0e},
			want:  []uint32{0x2a, 0x39, 0x5a, 0x77},
		},

		// TODO: bitpack > 8 bits
		// TODO: bitpack with a larger number of values

		{
			name:  "Bitpack large exception",
			input: []byte{0x84, 0x1a, 0x0, 0x8, 0x2c, 0xf7, 0xac, 0x2, 0x97, 0x43, 0x15, 0x73, 0x13, 0xe2},
			want:  []uint32{7, 9, 3, 4, 5, 1, 3, 7, 3, 1, 2, 718238414},
		},

		{
			name:  "constant",
			input: []byte{0xc8, 0x89},
			want:  []uint32{0x89},
		},

		{
			name:  "PFOR exceptions",
			input: []byte{0x44, 0x1, 0x97, 0x43, 0x15, 0x73, 0x13, 0xe2, 0xf, 0xb},
			want:  []uint32{7, 9, 3, 4, 5, 1, 3, 7, 3, 1, 2, 254},
		},

		{
			name:  "PFOR large exceptions",
			input: []byte{0x44, 0x1, 0x97, 0x43, 0x15, 0x73, 0x13, 0x62, 0xb3, 0xe, 0xb},
			want:  []uint32{7, 9, 3, 4, 5, 1, 3, 7, 3, 1, 2, 11254},
		},

		// TODO: move these to files, too
		{
			name:  "trigram 0",
			input: []byte{0x41, 0x14, 0xf2, 0xff, 0xff, 0xff, 0xfd, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfb, 0xff, 0xff, 0xff, 0xfb, 0xff, 0xff, 0xff, 0xf9, 0xff, 0xff, 0xff, 0xfc, 0xff, 0xff, 0xff, 0x48, 0xb4, 0x13, 0xc1, 0x89, 0x1, 0x1, 0x2, 0x1, 0x1, 0x1, 0x1, 0xe3, 0x1c, 0xb8, 0xcd, 0xf1, 0x59, 0xd, 0x2, 0xbc, 0x0, 0x1, 0xb1, 0x1a, 0x2, 0x1, 0xf1, 0x2d, 0x11, 0x0, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x18, 0x24, 0x41, 0x1, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x4, 0xe6, 0xc1, 0x1, 0x41, 0x24, 0xff, 0xff, 0xbf, 0xfb, 0xff, 0xff, 0xff, 0x7b, 0xff, 0xff, 0xff, 0x76, 0xff, 0xff, 0xdf, 0xfb, 0xff, 0xff, 0xdf, 0xfb, 0xff, 0xff, 0xdf, 0xf3, 0xff, 0xff, 0x9f, 0xf1, 0xff, 0xff, 0xcf, 0xf9, 0xb3, 0x7f, 0xbe, 0x59, 0xf1, 0x67, 0x47, 0xe3, 0xb5, 0x1e, 0xf1, 0xc5, 0x83, 0x1, 0xbb, 0xa0, 0x1, 0x12, 0x10, 0x2, 0x1, 0x1, 0xa, 0xd1, 0x3, 0xb9, 0xa3, 0x2, 0x1, 0x1, 0x5, 0x1, 0x2, 0x1, 0x1, 0xb6, 0x86, 0x3, 0xb4, 0x3f, 0xe4, 0x3, 0x33, 0xb8, 0x71, 0xf1, 0xc5, 0x3, 0x1, 0x1, 0x1, 0xe9, 0x94, 0xa7, 0xaa, 0xab, 0xac, 0xad, 0xae, 0xaf, 0xb0, 0xb6, 0xbc, 0xc2, 0xce, 0xcf, 0xd0, 0xd1, 0xd3, 0xd4, 0xd5, 0xd6, 0xd7, 0xd8, 0xd9, 0xda, 0xdb, 0xdd, 0xde, 0xdf, 0xe0, 0xe5, 0xf4, 0xf5, 0xf6, 0xf7, 0xf8, 0xf9, 0xfa, 0x83, 0xe, 0x8, 0xb0, 0x19, 0x3f, 0x24, 0x0, 0x80, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x24, 0x2, 0x1c, 0x80, 0x83, 0xff, 0x61, 0x0, 0x0, 0x0, 0x0, 0x0, 0x10, 0x0, 0x0, 0x50, 0x0, 0x10, 0xc2, 0xa5, 0x30, 0x1e, 0x58, 0x2, 0x3e, 0x48, 0xea, 0xe5, 0xdf, 0x98, 0x13, 0x3, 0x42, 0x91, 0x21, 0xb8, 0x8, 0x0, 0x2, 0x40, 0xb3, 0x11, 0x0, 0x6c, 0x76, 0xda, 0xc1, 0x5a, 0x12, 0x0, 0x4, 0x0, 0x2, 0x40, 0x0, 0x10, 0x0, 0x18, 0x6a, 0xd7, 0x4, 0xdc, 0x52, 0x3, 0x0, 0xc2, 0xda, 0xc5, 0x7f, 0x71, 0xd4, 0x35, 0x3, 0x81, 0x80, 0x41, 0xa1, 0x29, 0x6c, 0x4, 0x24, 0xc3, 0x4f, 0xd6, 0xd, 0xf0, 0x2, 0xaa, 0x1a, 0x49, 0x16, 0x26, 0xc9, 0x49, 0x92, 0xec, 0x49, 0x4c, 0xf0, 0x25, 0x7f, 0xcd, 0x95, 0x28, 0x49, 0x19, 0x1b, 0xe6, 0x49, 0x61, 0x8e, 0x24, 0x79, 0x61, 0xb2, 0x24, 0x4a, 0x71, 0x92, 0xf8, 0xf, 0x33, 0x26, 0x91, 0xb0, 0xb2, 0x24, 0x4f, 0x68, 0xa8, 0x36, 0x8e, 0x85, 0xb2, 0x24, 0x89, 0x82, 0xa3, 0xc4, 0x56, 0x92, 0xba, 0x9, 0x4e, 0xe2, 0xb2, 0x29, 0x4c, 0x62, 0x33, 0x2e, 0xc9, 0xee, 0x27, 0x49, 0x92, 0x24, 0x25, 0x49, 0x92, 0x24, 0x25, 0x49, 0x92, 0x24, 0x24, 0x49, 0x92, 0x24, 0x25, 0x49, 0x90, 0x30, 0x25, 0x49, 0x94, 0x24, 0x26, 0x49, 0x92, 0x2c, 0x24, 0x49, 0x92, 0x24, 0x83, 0xd, 0x0, 0x42, 0xee, 0x17, 0x0, 0x18, 0xa5, 0x5, 0x20, 0x83, 0xb7, 0x25, 0x75, 0xe2, 0xc0, 0x9c, 0x49, 0xf2, 0x61, 0x82, 0x32, 0x60, 0x3f, 0x45, 0x0, 0x49, 0x92, 0x24, 0x61, 0x92, 0x50, 0x49, 0x9c, 0x19, 0x79, 0x98, 0x24, 0x51, 0x92, 0x10},
			want:  []uint32{0x90, 0x1, 0x1, 0x1, 0x1, 0x1, 0x789, 0x2274, 0x3, 0x2, 0x5, 0x3, 0x3, 0x3, 0x2, 0x659a, 0x10fc, 0x9c15, 0x5, 0x1763, 0x2, 0x196, 0x4, 0x1, 0x2, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0xa3bd, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x9, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x660, 0x1, 0x1, 0x1c15, 0x11030, 0x66cc, 0x3c, 0x188ec, 0x2, 0x16a2, 0x1, 0x1, 0x1, 0x1, 0x1, 0x2, 0x1, 0x1, 0x1, 0x1, 0x1, 0x25, 0x1, 0x1, 0x1, 0x1, 0x1, 0x20, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x4, 0x2, 0x2, 0x14, 0x1, 0x4168, 0x12a8, 0x4, 0x2, 0x2, 0xb, 0x3, 0x4, 0x3, 0x1, 0x2, 0xc6e, 0x7, 0x7e1, 0x1, 0x1, 0x1, 0x1, 0x6769, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x67, 0x1045, 0x88ed, 0x3, 0x3, 0x2, 0x728a, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x4, 0x1085, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x14bb, 0xf1c, 0x4, 0x4b6, 0x41f1, 0x1, 0x1, 0xbd4f, 0x6ff4, 0x1, 0x1, 0x1, 0x2733, 0x1019, 0x3228, 0x5c12, 0x15, 0x17, 0x1, 0x1, 0x1, 0x1, 0x366f, 0x1, 0x1, 0x8, 0x3, 0x1, 0x4, 0x1, 0x3, 0x1, 0x4, 0x1, 0x1, 0x1, 0x1, 0x3, 0x1, 0x2, 0x1, 0x1, 0x1, 0xecde, 0xed1, 0x7, 0x1, 0x1, 0x7, 0x1, 0x1, 0x7, 0x1, 0x1, 0x7, 0x1, 0x1, 0x1, 0x2, 0x7, 0x1, 0x1, 0x7, 0x1, 0x1, 0x7, 0x1, 0x1, 0x7, 0x1, 0x1, 0x1, 0x5, 0x1, 0x1, 0x4, 0x1, 0x1, 0x4, 0x1, 0x1, 0x5, 0x1, 0x1, 0x3, 0x3, 0x2, 0x3, 0x2, 0x3, 0x3, 0x3, 0x4, 0x1, 0x5, 0x1, 0x1, 0x3, 0x3, 0x4, 0x1, 0x1, 0x5, 0x1, 0x1, 0x2, 0x2, 0x3, 0x1, 0x1, 0x4b59, 0x1, 0x6, 0x8, 0x1, 0x1, 0x1, 0xf, 0x6, 0x1, 0x6, 0x6, 0x4, 0x1, 0x2, 0x1, 0x11, 0x9, 0xa, 0x1, 0x1, 0x1, 0x2, 0x1, 0x6, 0x2, 0x1, 0x1, 0x1, 0xd433, 0x26b8, 0x5b84, 0x2, 0x1, 0x1, 0x1, 0x1, 0x1af, 0x18403, 0x2ed6, 0x2ff8, 0xea38, 0x669, 0x40e, 0x2836, 0x14d6, 0x8df, 0x2, 0x3, 0x1, 0x3, 0x1923, 0xc9fc, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x6e8, 0x2, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x5e4, 0x1, 0xd553, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x28c4, 0x1, 0x1, 0x1, 0x1, 0x16c, 0x2, 0x1, 0x641, 0x5b79, 0x3a96, 0x1, 0x38b, 0x339e, 0xf248, 0x1309, 0xca7, 0x7ec1, 0x4, 0x451, 0x1, 0x1, 0x1, 0x1, 0x2, 0x1, 0x1, 0x1, 0x1, 0x4},
		},

		{
			name:  "trigram1",
			input: []byte{0x43, 0x1e, 0xd8, 0xff, 0xff, 0xff, 0xd1, 0xff, 0xff, 0xef, 0xf9, 0xff, 0xff, 0xff, 0xc9, 0xff, 0xff, 0xff, 0xf9, 0xff, 0xff, 0xff, 0xfb, 0xff, 0xff, 0xff, 0xf8, 0xff, 0xff, 0xff, 0xfb, 0xff, 0xff, 0xff, 0xff, 0xff, 0xe7, 0x98, 0xff, 0xff, 0xe7, 0x96, 0xff, 0xff, 0x17, 0x15, 0xff, 0xff, 0xc7, 0x9e, 0xff, 0xff, 0x47, 0x9a, 0xff, 0xff, 0xa7, 0xc4, 0xff, 0xff, 0x97, 0x53, 0xff, 0xff, 0xd4, 0xf3, 0xa4, 0x53, 0x92, 0x24, 0x66, 0xc9, 0x92, 0x24, 0x80, 0x49, 0x92, 0x25, 0x38, 0x4b, 0x12, 0x2d, 0x38, 0x49, 0x92, 0x24, 0xbc, 0x49, 0x32, 0x25, 0x92, 0x49, 0x92, 0x24, 0x58, 0x49, 0x92, 0x24, 0x12, 0xb5, 0x8f, 0x1, 0xbd, 0x5, 0xdd, 0x8a, 0x3, 0xe3, 0x7d, 0x7, 0xe1, 0x6c, 0xb3, 0x33, 0xba, 0xd1, 0x1, 0xb1, 0xdc, 0xbe, 0x39, 0xc3, 0x82, 0xbe, 0xa0, 0xb6, 0xdb, 0xde, 0xcc, 0xed, 0x74, 0xb2, 0x2a, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x3, 0xb9, 0xba, 0x0, 0x5, 0x6, 0xa, 0xb, 0x49, 0x87, 0x88, 0x89, 0x8b, 0x8c, 0x90, 0x96, 0x98, 0x9b, 0x9d, 0xa6, 0xa7, 0xa9, 0xaa, 0xab, 0xac, 0xad, 0xae, 0xaf, 0xb0, 0xb1, 0xb2, 0xb6, 0xc3, 0x82, 0x10, 0x0, 0x1, 0x2, 0x4, 0x0, 0x0, 0x0, 0x0, 0xff, 0xb1, 0x1, 0x1, 0x0, 0x1, 0x0, 0xfb, 0x3, 0x8, 0x31, 0xdc, 0xd6, 0xc5, 0x17, 0x7f, 0x32, 0x4, 0x0, 0x4, 0x0, 0xbd, 0x1, 0x3, 0x0, 0x5, 0x41, 0x5b, 0x0, 0x14, 0x27, 0xf8, 0x71, 0x1, 0x0, 0x55, 0x55, 0x58, 0x65, 0xe9, 0x56, 0x45, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x36, 0x1b, 0x65, 0x93, 0x0},
			want:  []uint32{0x90, 0x1, 0x1, 0x1, 0x1, 0x2a03, 0x8, 0x3, 0x3, 0x2, 0x65b7, 0x169d9, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x1d, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x19974, 0x3c, 0x188ec, 0x2, 0x1720, 0x5410, 0x4, 0x2, 0x2, 0xb, 0x3, 0x4, 0x3, 0x1, 0x2, 0xc6e, 0x7, 0x6f54, 0x3, 0x2, 0x999f, 0x5, 0x728a, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x4, 0x3465, 0x173ef, 0x1, 0x1e92d, 0xed8, 0x9, 0x9, 0x9, 0xc, 0x9, 0x9, 0x9, 0x8, 0x6, 0x6, 0x7, 0x1c, 0x6, 0x5, 0x3, 0x4, 0x1, 0x1, 0x5, 0x4, 0x2, 0x3, 0x1, 0x1, 0x4b5b, 0x1, 0x1, 0x1, 0x1, 0x2, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x3, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x3, 0x1, 0x1, 0x1, 0x1, 0x3, 0x2, 0x1, 0x2, 0x1, 0x1, 0x1, 0x1, 0x1, 0x3, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x4, 0x2, 0x1, 0x1, 0x1, 0x1, 0x2, 0x1, 0x1, 0x6, 0x2, 0x3, 0x2, 0x1, 0x1, 0x1, 0x1, 0x1, 0xfec, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0xc422, 0x35b71, 0x5f17, 0xc9fc, 0x13, 0x12, 0x6f5, 0xc, 0x10415, 0x1, 0x2, 0x1, 0x16f, 0x9c50, 0x1, 0x1c7e2, 0x4},
		},

		{
			name:  "trigram2",
			input: []byte{0x43, 0x28, 0x70, 0xf2, 0xff, 0xff, 0x69, 0xf2, 0xff, 0xff, 0x51, 0xf4, 0xff, 0xff, 0x71, 0xfe, 0xff, 0xff, 0x61, 0xfe, 0xff, 0xff, 0x4e, 0xfe, 0xff, 0xff, 0x55, 0xfe, 0xff, 0xff, 0x48, 0xfe, 0xff, 0xff, 0xfb, 0xff, 0xff, 0xb5, 0xff, 0xff, 0xff, 0x91, 0xff, 0xff, 0xff, 0x29, 0xff, 0xff, 0xff, 0xe5, 0xff, 0xff, 0xff, 0x95, 0xff, 0xff, 0xff, 0x97, 0xff, 0xff, 0x7f, 0x97, 0xff, 0xff, 0x3f, 0x98, 0x85, 0x80, 0x9d, 0x24, 0xe6, 0x49, 0x92, 0x24, 0x31, 0xfb, 0x9c, 0x44, 0x3c, 0x49, 0x9b, 0x24, 0x4, 0x28, 0x9d, 0x2c, 0x24, 0x89, 0x93, 0x30, 0x14, 0xe8, 0x92, 0x24, 0x24, 0xc9, 0x13, 0x86, 0x12, 0xb5, 0x8e, 0x1, 0xd5, 0x91, 0x32, 0xc4, 0xc9, 0x3, 0xe3, 0x7d, 0xe1, 0x74, 0xbd, 0xb5, 0x1, 0xb2, 0xd9, 0xd0, 0x71, 0xbe, 0xa0, 0xb4, 0xf8, 0xbb, 0x7, 0xd5, 0xf7, 0xb7, 0x38, 0xe6, 0x8a, 0xb2, 0x2a, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x3, 0xb9, 0xc2, 0x5, 0xf1, 0xf8, 0x43, 0xb1, 0x9d, 0xb5, 0x55, 0xb7, 0x2b, 0xc9, 0x8e, 0x2, 0xcc, 0x98, 0xb5, 0x67, 0xc4, 0x78, 0xe9, 0x64, 0x0, 0x5, 0x7, 0x9, 0xb, 0x1a, 0x58, 0x96, 0x97, 0x99, 0x9d, 0xa3, 0xa8, 0xaa, 0xb2, 0xb3, 0xb6, 0xb8, 0xb9, 0xba, 0xbc, 0xbe, 0xc0, 0xc2, 0xc4, 0xc6, 0xc8, 0xd2, 0xd3, 0xd4, 0xd6, 0xd7, 0xd8, 0xd9, 0xda, 0xdb, 0xdc, 0xef, 0xf4, 0xf5},
			want:  []uint32{0x90, 0x1, 0x1, 0x1, 0x1, 0x29fe, 0x5, 0x8, 0x6, 0x12a15, 0x2, 0x196, 0x4, 0x1, 0x2, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0xa3d2, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x1d, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x19975, 0x18928, 0x2, 0x6b30, 0x4, 0x2, 0x2, 0xb, 0x3, 0x4, 0x3, 0x1, 0x2, 0x1456, 0x1, 0x1, 0x1, 0x1, 0x10113, 0x5, 0x728a, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x254c, 0x55c7, 0x1, 0x1, 0x12d45, 0x1, 0x374c, 0x1b1df, 0xed9, 0x1, 0x8, 0x1, 0x8, 0x1, 0x8, 0x1, 0xb, 0x1, 0x8, 0x1, 0x8, 0x1, 0x8, 0x1, 0x7, 0x1, 0x5, 0x1, 0x5, 0x1, 0x6, 0x1, 0x1b, 0x4b9d, 0x2c, 0x6, 0x4254b, 0xa77, 0x2836, 0x36e1, 0xc9fe, 0x15, 0xe24e, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x28c4, 0x1, 0x1, 0x1, 0x1, 0xa14b, 0x1c8ac, 0x1, 0x1, 0x1, 0x1, 0x2, 0x1, 0x1, 0x1, 0x1, 0x4},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			buffer := make([]uint32, len(test.want), len(test.want)+32)
			padded := make([]byte, len(test.input)+32)
			copy(padded, test.input)
			//padded = padded[:len(test.input)]
			num := P4ndec256v32(test.input, buffer)
			if got, want := num, len(test.input); got != want {
				t.Fatalf("read: got %d, want %d", got, want)
			}
			if got, want := len(buffer), len(test.want); got != want {
				t.Fatalf("len: got %d, want %d", got, want)
			}
			mismatch := 0
			for i := 0; i < len(buffer); i++ {
				if buffer[i] != test.want[i] {
					mismatch++
				}
			}
			if mismatch > 0 {
				t.Logf("%d values don’t match", mismatch)
			}
			if !reflect.DeepEqual(buffer, test.want) {
				t.Fatalf("got %x, want %x", buffer, test.want)
			}
		})
	}
}

func TestDecodeFromFile(t *testing.T) {
	for _, fn := range []string{
		"trigram_592137",
	} {
		t.Run(fn, func(t *testing.T) {
			wantb, err := ioutil.ReadFile("testdata/" + fn + ".want")
			if err != nil {
				t.Fatal(err)
			}
			want := make([]uint32, len(wantb)/4)
			if err := binary.Read(bytes.NewReader(wantb), binary.LittleEndian, want); err != nil {
				t.Fatal(err)
			}
			input, err := ioutil.ReadFile("testdata/" + fn + ".input")
			if err != nil {
				t.Fatal(err)
			}

			buffer := make([]uint32, len(want), len(want)+32)
			padded := make([]byte, len(input)+32)
			copy(padded, input)
			//padded = padded[:len(input)]
			num := P4ndec256v32(input, buffer)
			if got, want := num, len(input); got != want {
				t.Fatalf("read: got %d, want %d", got, want)
			}
			if got, want := len(buffer), len(want); got != want {
				t.Fatalf("len: got %d, want %d", got, want)
			}
			mismatch := 0
			for i := 0; i < len(buffer); i++ {
				if buffer[i] != want[i] {
					t.Logf("index %d: got %x, want %x", i, buffer[i], want[i])
					mismatch++
				}
			}
			if mismatch > 0 {
				t.Fatalf("%d values don’t match", mismatch)
			}
			if !reflect.DeepEqual(buffer, want) {
				//t.Fatalf("got %x, want %x", buffer, want)
			}
		})
	}
}

func TestVbdec32(t *testing.T) {
	for _, test := range []struct {
		name  string
		input []byte
		want  []uint32
	}{
		{
			name:  "first1", // first value fitting in 1 byte
			input: []byte{0x00},
			want:  []uint32{0},
		},

		{
			name:  "last1", // last value fitting in 1 byte
			input: []byte{0xb0},
			want:  []uint32{176},
		},

		{
			name:  "first2", // first value fitting in 2 bytes
			input: []byte{0xb1, 0x00},
			want:  []uint32{177},
		},

		{
			name:  "last2", // last value fitting in 2 bytes
			input: []byte{0xf0, 0xff},
			want:  []uint32{16560},
		},

		{
			name:  "first3", // first value fitting in 3 bytes
			input: []byte{0xf1, 0x00, 0x00},
			want:  []uint32{16561},
		},

		{
			name:  "last3", // last value fitting in 3 bytes
			input: []byte{0xf8, 0xff, 0xff},
			want:  []uint32{540848},
		},

		{
			name:  "first4", // first value fitting in 4 bytes
			input: []byte{0xf9, 0xb1, 0x40, 0x08},
			want:  []uint32{540849},
		},

		{
			name:  "last4", // last value fitting in 4 bytes
			input: []byte{0xf9, 0xff, 0xff, 0xff},
			want:  []uint32{16777215},
		},

		{
			name:  "first5", // first value fitting in 5 bytes (overflow)
			input: []byte{0xff, 0x00, 0x00, 0x00, 0x01},
			want:  []uint32{16777216},
		},

		{
			name:  "last5", // last value fitting in 5 bytes (overflow)
			input: []byte{0xff, 0xff, 0xff, 0xff, 0xff},
			want:  []uint32{4294967295},
		},

		{
			name:  "multi5", // multiple values, exercising the 5 bytes
			input: []byte{0x00, 0x00, 0x00, 0xfa, 0xff, 0xff, 0xff, 0xff},
			want:  []uint32{0, 0, 0, 4294967295},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			padded := make([]byte, len(test.input)*4)
			copy(padded, test.input)
			output := make([]uint32, len(test.want))
			read := vbdec32(padded, output)
			if got, want := read, len(test.input); got != want {
				t.Fatalf("vbdec32 read %d, want %d", got, want)
			}
			if got, want := output, test.want; !reflect.DeepEqual(got, want) {
				t.Fatalf("vbdec32: got %d, want %d", got, want)
			}
		})
	}
}

func Test_Bitunpack256v32(t *testing.T) {
	for _, test := range []struct {
		name       string
		input      []byte
		want       []uint32
		bits       byte
		exceptions []uint32
		bb         []byte
	}{
		{
			name:       "trigram0",
			input:      []byte{0x49, 0x16, 0x26, 0xc9, 0x49, 0x92, 0xec, 0x49, 0x4c, 0xf0, 0x25, 0x7f, 0xcd, 0x95, 0x28, 0x49, 0x19, 0x1b, 0xe6, 0x49, 0x61, 0x8e, 0x24, 0x79, 0x61, 0xb2, 0x24, 0x4a, 0x71, 0x92, 0xf8, 0xf, 0x33, 0x26, 0x91, 0xb0, 0xb2, 0x24, 0x4f, 0x68, 0xa8, 0x36, 0x8e, 0x85, 0xb2, 0x24, 0x89, 0x82, 0xa3, 0xc4, 0x56, 0x92, 0xba, 0x9, 0x4e, 0xe2, 0xb2, 0x29, 0x4c, 0x62, 0x33, 0x2e, 0xc9, 0xee, 0x27, 0x49, 0x92, 0x24, 0x25, 0x49, 0x92, 0x24, 0x25, 0x49, 0x92, 0x24, 0x24, 0x49, 0x92, 0x24, 0x25, 0x49, 0x90, 0x30, 0x25, 0x49, 0x94, 0x24, 0x26, 0x49, 0x92, 0x2c, 0x24, 0x49, 0x92, 0x24},
			bits:       3,
			exceptions: []uint32{0x210, 0x297, 0x1e3, 0x96, 0x83e, 0x17a9, 0xdfe, 0x4e6, 0x203, 0x645, 0xb82, 0x2, 0x2, 0x6cd, 0x1, 0x1d9b, 0x1da, 0x96b, 0x1, 0x1, 0x2, 0x1, 0x1, 0x1a86, 0x4d7, 0xb70, 0x35, 0x3080, 0x5da, 0x5ff, 0x1d47, 0xcd, 0x81, 0x506, 0x29a, 0x11b, 0x324, 0x193f, 0xdd, 0xbc, 0x1aaa},
			bb:         []byte{0x8, 0xb0, 0x19, 0x3f, 0x24, 0x0, 0x80, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x24, 0x2, 0x1c, 0x80, 0x83, 0xff, 0x61, 0x0, 0x0, 0x0, 0x0, 0x0, 0x10, 0x0, 0x0, 0x50, 0x0, 0x10, 0xc2, 0xa5, 0x30, 0x1e, 0x58, 0x2, 0x3e, 0x48, 0xea, 0xe5, 0xdf, 0x98, 0x13, 0x3, 0x42, 0x91, 0x21, 0xb8, 0x8, 0x0, 0x2, 0x40, 0xb3, 0x11, 0x0, 0x6c, 0x76, 0xda, 0xc1, 0x5a, 0x12, 0x0, 0x4, 0x0, 0x2, 0x40, 0x0, 0x10, 0x0, 0x18, 0x6a, 0xd7, 0x4, 0xdc, 0x52, 0x3, 0x0, 0xc2, 0xda, 0xc5, 0x7f, 0x71, 0xd4, 0x35, 0x3, 0x81, 0x80, 0x41, 0xa1, 0x29, 0x6c, 0x4, 0x24, 0xc3, 0x4f, 0xd6, 0xd, 0xf0, 0x2, 0xaa, 0x1a, 0x49, 0x16, 0x26, 0xc9, 0x49, 0x92, 0xec, 0x49, 0x4c, 0xf0, 0x25, 0x7f, 0xcd, 0x95, 0x28, 0x49, 0x19, 0x1b, 0xe6, 0x49, 0x61, 0x8e, 0x24, 0x79, 0x61, 0xb2, 0x24, 0x4a, 0x71, 0x92, 0xf8, 0xf, 0x33, 0x26, 0x91, 0xb0, 0xb2, 0x24, 0x4f, 0x68, 0xa8, 0x36, 0x8e, 0x85, 0xb2, 0x24, 0x89, 0x82, 0xa3, 0xc4, 0x56, 0x92, 0xba, 0x9, 0x4e, 0xe2, 0xb2, 0x29, 0x4c, 0x62, 0x33, 0x2e, 0xc9, 0xee, 0x27, 0x49, 0x92, 0x24, 0x25, 0x49, 0x92, 0x24, 0x25, 0x49, 0x92, 0x24, 0x24, 0x49, 0x92, 0x24, 0x25, 0x49, 0x90, 0x30, 0x25, 0x49, 0x94, 0x24, 0x26, 0x49, 0x92, 0x2c, 0x24, 0x49, 0x92, 0x24, 0x83, 0xd, 0x0, 0x42, 0xee, 0x17, 0x0, 0x18, 0xa5, 0x5, 0x20, 0x83, 0xb7, 0x25, 0x75, 0xe2, 0xc0, 0x9c, 0x49, 0xf2, 0x61, 0x82, 0x32, 0x60, 0x3f, 0x45, 0x0, 0x49, 0x92, 0x24, 0x61, 0x92, 0x50, 0x49, 0x9c, 0x19, 0x79, 0x98, 0x24, 0x51, 0x92, 0x10},
			want:       []uint32{0x1, 0x1, 0x4, 0x1085, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x14bb, 0xf1c, 0x4, 0x4b6, 0x41f1, 0x1, 0x1, 0xbd4f, 0x6ff4, 0x1, 0x1, 0x1, 0x2733, 0x1019, 0x3228, 0x5c12, 0x15, 0x17, 0x1, 0x1, 0x1, 0x1, 0x366f, 0x1, 0x1, 0x8, 0x3, 0x1, 0x4, 0x1, 0x3, 0x1, 0x4, 0x1, 0x1, 0x1, 0x1, 0x3, 0x1, 0x2, 0x1, 0x1, 0x1, 0xecde, 0xed1, 0x7, 0x1, 0x1, 0x7, 0x1, 0x1, 0x7, 0x1, 0x1, 0x7, 0x1, 0x1, 0x1, 0x2, 0x7, 0x1, 0x1, 0x7, 0x1, 0x1, 0x7, 0x1, 0x1, 0x7, 0x1, 0x1, 0x1, 0x5, 0x1, 0x1, 0x4, 0x1, 0x1, 0x4, 0x1, 0x1, 0x5, 0x1, 0x1, 0x3, 0x3, 0x2, 0x3, 0x2, 0x3, 0x3, 0x3, 0x4, 0x1, 0x5, 0x1, 0x1, 0x3, 0x3, 0x4, 0x1, 0x1, 0x5, 0x1, 0x1, 0x2, 0x2, 0x3, 0x1, 0x1, 0x4b59, 0x1, 0x6, 0x8, 0x1, 0x1, 0x1, 0xf, 0x6, 0x1, 0x6, 0x6, 0x4, 0x1, 0x2, 0x1, 0x11, 0x9, 0xa, 0x1, 0x1, 0x1, 0x2, 0x1, 0x6, 0x2, 0x1, 0x1, 0x1, 0xd433, 0x26b8, 0x5b84, 0x2, 0x1, 0x1, 0x1, 0x1, 0x1af, 0x18403, 0x2ed6, 0x2ff8, 0xea38, 0x669, 0x40e, 0x2836, 0x14d6, 0x8df, 0x2, 0x3, 0x1, 0x3, 0x1923, 0xc9fc, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x6e8, 0x2, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x5e4, 0x1, 0xd553, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1},
		},
	} {
		output := make([]uint32, len(test.want))
		got := _bitunpack256v32(test.input, output, test.bits, test.exceptions, test.bb)
		if want := len(test.input); got != want {
			t.Fatalf("_bitunpack256v32 did not consume the entire input: got %d, want %d", got, want)
		}
		if got, want := len(output), len(test.want); got != want {
			t.Fatalf("len: got %d, want %d", got, want)
		}
		mismatch := 0
		for i := 0; i < len(output); i++ {
			if output[i] != test.want[i] {
				mismatch++
			}
		}
		if mismatch > 0 {
			t.Logf("%d values don’t match", mismatch)
		}
		if !reflect.DeepEqual(output, test.want) {
			t.Fatalf("got %x, want %x", output, test.want)
		}
	}
}
