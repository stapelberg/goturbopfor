package goturbopfor

import (
	"encoding/binary"
	"log"
	"math/bits"
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
	var rbits uint // remaining bits
	var acc uint32 // accumulator
	for op := 0; op < len(output); {
		if rbits < uint(nbits) {
			// read one more byte
			acc |= uint32(input[0]) << rbits
			input = input[1:]
			rbits += 8
		}
		if rbits >= uint(nbits) {
			output[op] = acc & ((1 << nbits) - 1)
			op++
			acc >>= nbits
			rbits -= uint(nbits)
		}
	}
	return orig - len(input)
}

func vbdec32(input []byte, output []uint32) (read int) {
	if input[0] == 0xff {
		log.Fatal("memcpy not implemented")
		// TODO: memcpy
	}
	before := len(input)
	for op := 0; op < len(output); op++ {
		x := uint32(input[0])
		input = input[1:]
		if x < 177 {
		} else if x < 241 {
			x = (x << 8) + uint32(input[0]) - 0xb04f
			input = input[1:]
		} else if x < 249 {
			x = (uint32(input[0]) + (uint32(input[1]) << 8)) +
				((x - 241) << 16) +
				16561
			input = input[2:]
		} else {
			_b := x - 249
			x = (uint32(input[0]) + (uint32(input[1]) << 8) + (uint32(input[2]) << 16) + (uint32(input[3]) << 24)) &
				(((1 << (8 * _b)) << 24) - 1)
			input = input[3+_b:]
		}
		output[op] = x
	}
	return before - len(input)
}

const blockSize = 256 // number of ints per block

func _p4dec32(input []byte, output []uint32, b, bx byte) (read int) {
	before := len(input)

	// Copy the exception bitmap into bb and count ones (i.e. exceptions).
	var bb [blockSize / 64]uint64
	var i int
	var num int
	n := len(output)
	for i = 0; i < n/64; i++ {
		bb[i] = binary.LittleEndian.Uint64(input)
		num += bits.OnesCount64(bb[i])
		input = input[8:]
	}
	if rest := n % 64; rest != 0 {
		bb[i] = binary.LittleEndian.Uint64(input) & ((1 << uint(rest)) - 1)
		num += bits.OnesCount64(bb[i])
		input = input[(rest+7)/8:]
	}

	exceptions := make([]uint32, num)
	input = input[bitunpack32(input, exceptions, bx):]
	input = input[bitunpack32(input, output, b):]

	op := 0
	k := 0
	for i := 0; i < (n+63)/64; i++ { // 64 bits of exception bitmap at a time
		for u := bb[i]; u != 0; u &= u - 1 { // bit-wise
			output[op+bits.TrailingZeros64(u)] += exceptions[k] << b
			k++
		}
		op += 64
	}

	return before - len(input)
}

var (
	blockBitpacking             = [2]byte{0, 0}
	blockBitpackingExceptions   = [2]byte{1, 0}
	blockBitpackingVBExceptions = [2]byte{0, 1}
	blockConstant               = [2]byte{1, 1}
)

// p4dec32 decodes one block of TurboPFor-encoded 32 bit ints
func p4dec32(input []byte, output []uint32) (read int) {
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
		return 1 + bitunpack32(input, output, b)

	case blockBitpackingExceptions:
		bx, input := input[0], input[1:]
		return before - len(input) + _p4dec32(input, output, b, bx)

	default: // blockBitpackingVBExceptions
		n, input := int(input[0]), input[1:] // number of exceptions
		input = input[bitunpack32(input, output, b):]

		exceptions := make([]uint32, n)
		input = input[vbdec32(input, exceptions):]
		for i := 0; i < n; i++ {
			output[input[i]] |= exceptions[i] << b
		}
		return before - len(input) + n
	}
}

func bitunpack256v32(input []byte, output []uint32, b byte) (read int) {
	orig := len(input)
	var bits uint
	var acc [8]uint64 // accumulator
	for op := 0; op < len(output); {
		if bits < uint(b) {
			// read 8 more bytes
			for i := 0; i < 8; i++ {
				acc[i] |= uint64(binary.LittleEndian.Uint32(input)) << bits
				input = input[4:]
			}
			bits += 32
		}
		if bits >= uint(b) {
			for i := 0; i < 8; i++ {
				output[op] = uint32(acc[i] & ((1 << b) - 1))
				op++
				acc[i] >>= b
			}
			bits -= uint(b)
		}
	}
	return orig - len(input)
}

func _bitunpack256v32(input []byte, output []uint32, b byte, exceptions []uint32, bb []byte) (read int) {
	orig := len(input)
	var rbits uint
	var acc [8]uint64 // accumulator
	for op := 0; op < len(output); {
		if rbits < uint(b) {
			// read 8 more bytes
			for i := 0; i < 8; i++ {
				acc[i] |= uint64(binary.LittleEndian.Uint32(input)) << rbits
				input = input[4:]
			}
			rbits += 32
		}
		if rbits >= uint(b) {
			xm := bb[0]
			bb = bb[1:]
			exc := 0
			for i := 0; i < 8; i++ {
				ov := acc[i] & ((1 << b) - 1)
				if xm&(1<<uint(i)) != 0 {
					ov |= uint64(exceptions[exc] << b)
					exc++
				}
				output[op] = uint32(ov)
				op++
				acc[i] >>= b
			}
			exceptions = exceptions[exc:]
			rbits -= uint(b)
		}
	}
	return orig - len(input)
}

func _p4dec256v32(input []byte, output []uint32, b, bx byte) (read int) {
	before := len(input)

	pb := input // TODO: document
	exceptions := make([]uint32, 256+64)
	var num int
	for i := 0; i < 32; i += 8 {
		num += bits.OnesCount64(binary.LittleEndian.Uint64(input))
		input = input[8:]
	}
	input = input[bitunpack32(input, exceptions[:num], bx):]
	input = input[_bitunpack256v32(input, output, b, exceptions[:num], pb):]
	return before - len(input)
}

// p4dec256v32 fills output from input, decoding 256 uint32s.
func p4dec256v32(input []byte, output []uint32) (read int) {
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
		return 1 + bitunpack256v32(input, output, b)

	case blockBitpackingExceptions:
		bx, input := input[0], input[1:]
		return before - len(input) + _p4dec256v32(input, output, b, bx)

	default: // blockBitpackingVBExceptions
		n, input := int(input[0]), input[1:] // number of exceptions
		input = input[bitunpack256v32(input, output, b):]

		exceptions := make([]uint32, n)
		input = input[vbdec32(input, exceptions):]
		for i := 0; i < n; i++ {
			output[input[i]] |= exceptions[i] << b
		}
		return before - len(input) + n
	}
}

// P4ndec256v32 fills output from input, decoding 256 uint32s at a time.
//
// Note that different decoding algorithms are used for the last block, if that
// block does not contain 256 uint32s.
func P4ndec256v32(input []byte, output []uint32) (read int) {
	before := len(input)
	for len(output) >= 256 {
		input = input[p4dec256v32(input, output[:256]):]
		output = output[256:]
	}
	return before - len(input) + p4dec32(input, output)
}
