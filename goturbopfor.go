package goturbopfor

import (
	"encoding/binary"
)

// bitpacking 7 bit uses little endian:
// 1 0101010   10 011100
// n current   nn current
//
// → 0101010, 0111001, xxxxx10, …
//
// longer values:
// 1120345    (100010001100001011001)
//                          01011001
//                  00011000
//         101 10001
//00001011
func bitunpack32(input []byte, output []uint32, nbits byte) (read int) {
	orig := len(input)
	var rbits byte // remaining bits
	var acc uint32 // accumulator
	for op := 0; op < len(output); {
		if rbits < nbits {
			// shift in one more byte
			acc |= uint32(input[0]) << rbits
			input = input[1:]
			rbits += 8
		}
		if rbits >= nbits {
			output[op] = acc & ((1 << nbits) - 1)
			op++
			acc >>= nbits
			rbits -= nbits
		}
	}
	return orig - len(input)
}

func bitunpack256v32(input []byte, output []uint32, nbits byte) (read int) {
	orig := len(input)
	var bits uint
	var acc [8]uint64 // accumulator
	for op := 0; op < len(output); {
		if bits < uint(nbits) {
			// read 8 more uint32s
			for i := 0; i < 8; i++ {
				acc[i] |= uint64(binary.LittleEndian.Uint32(input)) << bits
				input = input[4:]
			}
			bits += 32
		}
		if bits >= uint(nbits) {
			for i := 0; i < 8; i++ {
				output[op] = uint32(acc[i] & ((1 << nbits) - 1))
				op++
				acc[i] >>= nbits
			}
			bits -= uint(nbits)
		}
	}
	return orig - len(input)
}

// vbdec32 fills output from input, decoding variable byte uint32s.
//
// The variable byte encoding is similar to:
// https://sqlite.org/src4/doc/trunk/www/varint.wiki
// https://github.com/stoklund/varint
//
// [       0-       176] are stored in 1 byte (as-is)
// [     177-     16560] are stored in 2 bytes, with the highest 6 bits added to 177
// [   16561-    540848] are stored in 3 bytes, with the highest 3 bits added to 241
// [  540849-  16777215] are stored in 4 bytes, with 0 added to 249
// [16777216-4294967295] are stored in 5 bytes, with 1 added to 249
//
// An overflow marker will be used to signal that encoding the
// values would be less space-efficient than simply copying them
// (e.g. if all values require 5 bytes).
func vbdec32(input []byte, output []uint32) (read int) {
	before := len(input)
	if input[0] == 0xff {
		// overflow, memcpy the data as-is:
		input = input[1:]
		for op := 0; op < len(output); op++ {
			output[op] = binary.LittleEndian.Uint32(input)
			input = input[4:]
		}
		return before - len(input)
	}
	for op := 0; op < len(output); op++ {
		x := uint32(input[0])
		input = input[1:]
		if x < 177 {
		} else if x < 241 {
			x = uint32(input[0]) +
				((x - 177) << 8) +
				177
			input = input[1:]
		} else if x < 249 {
			x = (uint32(input[0]) << 0) +
				(uint32(input[1]) << 8) +
				((x - 241) << 16) +
				16561
			input = input[2:]
		} else {
			_b := x - 249 // _b in [0, 1]
			x = binary.LittleEndian.Uint32(input) & (((1 << (8 * _b)) << 24) - 1)
			input = input[3+_b:]
		}
		output[op] = x
	}
	return before - len(input)
}

var (
	// bitpacked values (no exceptions)
	blockBitpacking = [2]byte{0, 0}

	// exception presence bitmap + bitpacked exception values
	blockBitpackingExceptions = [2]byte{1, 0}

	// variable byte encoded exception values and exception index bytes
	blockBitpackingVBExceptions = [2]byte{0, 1}

	// constant value for entire block
	blockConstant = [2]byte{1, 1}
)

type decoder struct {
	bitunpack func(input []byte, output []uint32, b byte) int
}

var (
	// v256 is a decoder which operates on 256 uint32s.
	v256 = decoder{bitunpack: bitunpack256v32}

	// remainder is a decoder which handles the remaining (<256) uint32s.
	remainder = decoder{bitunpack: bitunpack32}
)

// p4dec32 decodes one block of TurboPFor-encoded 32 bit ints
func (d *decoder) p4dec32(input []byte, output []uint32) (read int) {
	if len(output) == 0 {
		return 0
	}
	before := len(input)            // for returning read bytes
	b, input := input[0], input[1:] // block header
	blockType := [2]byte{
		(b & 0x80) >> 7, // first bit
		(b & 0x40) >> 6, // second bit
	}
	b &= ^byte(0x80 | 0x40) // for bitpacking, b is the number of bits
	switch blockType {
	case blockConstant:
		padded := make([]byte, binary.Size(uint32(0)))
		copy(padded, input)
		u := binary.LittleEndian.Uint32(padded)
		if b < 32 {
			u &= ((1 << b) - 1)
		}
		for i := 0; i < len(output); i++ {
			output[i] = u
		}
		return 1 + (int(b)+7)/8

	case blockBitpacking:
		return 1 + d.bitunpack(input, output, b)

	case blockBitpackingExceptions:
		bx, input := input[0], input[1:]
		n := len(output)

		exmap := input
		nex := 0 // number of exceptions
		for i := 0; i < n; i++ {
			if exmap[i/8]&(1<<uint(i%8)) != 0 {
				nex++
			}
		}
		input = input[(n+7)/8:]

		exceptions := make([]uint32, nex)
		input = input[bitunpack32(input, exceptions, bx):]
		input = input[d.bitunpack(input, output, b):]

		for i := 0; i < n; i++ {
			if exmap[i/8]&(1<<uint(i%8)) != 0 {
				output[i] += exceptions[0] << b
				exceptions = exceptions[1:]
			}
		}

		return before - len(input)

	default: // blockBitpackingVBExceptions
		nex, input := int(input[0]), input[1:] // number of exceptions
		input = input[d.bitunpack(input, output, b):]

		exceptions := make([]uint32, nex)
		input = input[vbdec32(input, exceptions):]
		for i := 0; i < nex; i++ {
			output[input[i]] |= exceptions[i] << b
		}
		return before - len(input) + nex
	}
}

// P4ndec256v32 fills output from input, decoding 256 uint32s at a time.
//
// Note that different decoding algorithms are used for the last block, if that
// block does not contain 256 uint32s.
func P4ndec256v32(input []byte, output []uint32) (read int) {
	before := len(input)
	for len(output) >= 256 {
		input = input[v256.p4dec32(input, output[:256]):]
		output = output[256:]
	}
	return before - len(input) + remainder.p4dec32(input, output)
}
