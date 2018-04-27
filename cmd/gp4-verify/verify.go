package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"goturbopfor"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"turbopfor"

	"golang.org/x/exp/mmap"
)

var (
	idx = flag.String("idx", "", "")
)

type Trigram uint32

// A MetaEntry defines the position within the index of the data associated with
// a trigram.
type MetaEntry struct {
	Trigram    Trigram
	Entries    uint32 // number of entries (excluding padding)
	OffsetCtrl int64  // control bytes offset within the corresponding .ctrl file
	OffsetEnc  int64  // streaming vbyte deltas offset within the corresponding .data file
}

const metaEntrySize = 24 // (encoding/binary).Size(MetaEntry)

func logic(idx string) error {
	log.Printf("verifying index %q", idx)
	md, err := mmap.Open(filepath.Join(idx, "posting.docid.turbopfor"))
	if err != nil {
		return err
	}
	defer md.Close()
	mf, err := os.Open(filepath.Join(idx, "posting.docid.meta"))
	if err != nil {
		return err
	}
	defer mf.Close()
	st, err := mf.Stat()
	if err != nil {
		return err
	}
	var meta, prev MetaEntry
	for i := 0; i < int(st.Size())/metaEntrySize; i++ {
		if err := binary.Read(mf, binary.LittleEndian, &meta); err != nil {
			return err
		}
		if i == 0 {
			prev = meta
			continue
		}
		buffer := make([]byte, meta.OffsetEnc-prev.OffsetEnc+32, meta.OffsetEnc-prev.OffsetEnc+32)
		if _, err := md.ReadAt(buffer, prev.OffsetEnc); err != nil {
			return err
		}

		// defense copy, just in case
		buf2 := make([]byte, len(buffer), cap(buffer))
		copy(buf2, buffer)

		deltas := make([]uint32, prev.Entries, prev.Entries+32)
		turbopfor.P4ndec256v32(buffer, deltas)

		deltas2 := make([]uint32, prev.Entries, prev.Entries+32)
		goturbopfor.P4ndec256v32(buf2, deltas2)

		if !reflect.DeepEqual(deltas, deltas2) {
			log.Printf("read %d bytes at %d for trigram %v: %v", meta.OffsetEnc-prev.OffsetEnc, prev.OffsetEnc, prev.Trigram, deltas)
			log.Printf("go: %v", deltas2)
			fmt.Printf("input: %#v,\n", buf2[:len(buffer)-32])
			fmt.Printf("want: %#v,\n", deltas)
			break
		}
		prev = meta
	}
	return nil
}

func main() {
	flag.Parse()
	if err := logic(*idx); err != nil {
		log.Fatal(err)
	}
}
