package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"goturbopfor"
	"goturbopfor/internal/mmap"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"sync"
	"sync/atomic"
	"time"
	"turbopfor"
	// go get github.com/Debian/dcs/internal/turbopfor
	// ln -s github.com/Debian/dcs/internal/turbopfor turbopfor
	// </hackerman>
)

var (
	idx = flag.String("idx", "", "")
)

type Trigram uint32

// A MetaEntry specifies the offset (in the corresponding data file) and number
// of entries for each trigram.
type MetaEntry struct {
	Trigram    Trigram
	Entries    uint32 // number of entries (excluding padding)
	OffsetData int64  // delta offset within the corresponding .data or .turbopfor file
}

var encoding = binary.LittleEndian

func (me *MetaEntry) Unmarshal(b []byte) {
	me.Trigram = Trigram(encoding.Uint32(b))
	me.Entries = encoding.Uint32(b[4:])
	me.OffsetData = int64(encoding.Uint64(b[8:]))
}

// metaEntrySize is (encoding/binary).Size(&MetaEntry{}), which Go 1.11 does not
// turn into a compile-time constant yet.
const metaEntrySize = 16

var (
	totalBytes   uint64
	totalEntries uint64
)

func init() {
	go func() {
		var last uint64
		for range time.Tick(1 * time.Second) {
			val := atomic.LoadUint64(&totalBytes) - last
			log.Printf("%d bytes/s", val)
			last = val
		}
	}()
}

var bufPool = sync.Pool{
	New: func() interface{} {
		// TODO: get average size
		return make([]uint32, 1000)
	},
}

func logic(idx string) error {
	log.Printf("verifying index %q", idx)
	md, err := mmap.Open(filepath.Join(idx, "posting.docid.turbopfor"))
	if err != nil {
		return err
	}
	defer md.Close()
	mf, err := mmap.Open(filepath.Join(idx, "posting.docid.meta"))
	if err != nil {
		return err
	}
	defer mf.Close()

	work := make(chan MetaEntry, 4*runtime.NumCPU())
	var wg sync.WaitGroup
	for i := 0; i < 2*runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for prev := range work {
				n := prev.Entries // for convenience
				deltas := bufPool.Get().([]uint32)
				if cap(deltas) < 2*(int(n)+32) {
					deltas = make([]uint32, n, 2*(int(n)+32))
				}
				turbopfor.P4ndec256v32(md.Data[prev.OffsetData:], deltas[:n])

				// measure TurboPFor decoding speed:
				// bufPool.Put(deltas)
				// continue
				deltas = deltas[:n]

				deltas2 := make([]uint32, prev.Entries, prev.Entries+32)
				goturbopfor.P4ndec256v32(md.Data[prev.OffsetData:], deltas2)

				if !reflect.DeepEqual(deltas, deltas2) {
					log.Printf("read %d bytes at %d for trigram %v: %v (len: %d)", -1 /*meta.OffsetData-prev.OffsetData*/, prev.OffsetData, prev.Trigram, deltas, prev.Entries)
					log.Printf("go: %v", deltas2)
					//fmt.Printf("input: %#v,\n", buffer[:len(buffer)-32])
					fmt.Printf("want: %#v,\n", deltas)
					// dump to turn into a testcase
					prefix := fmt.Sprintf("/home/michael/go/src/goturbopfor/testdata/trigram_%d", prev.Trigram)
					want, err := os.Create(prefix + ".want")
					if err != nil {
						log.Fatal(err)
					}
					defer want.Close()
					if err := binary.Write(want, binary.LittleEndian, deltas); err != nil {
						log.Fatal(err)
					}

					input, err := os.Create(prefix + ".input")
					if err != nil {
						log.Fatal(err)
					}
					defer input.Close()
					// if _, err := input.Write(buffer[:len(buffer)-32]); err != nil {
					// 	log.Fatal(err)
					// }

					os.Exit(1)
				}

				bufPool.Put(deltas)
			}
		}()
	}

	max := int(len(mf.Data))/metaEntrySize - 2 /* corrupt? */
	var meta, prev MetaEntry
	for i := 0; i < max; i++ {
		meta.Unmarshal(mf.Data[i*metaEntrySize:])
		if i == 0 {
			prev = meta
			continue
		}
		blockLen := meta.OffsetData - prev.OffsetData
		atomic.AddUint64(&totalBytes, uint64(blockLen+metaEntrySize))
		atomic.AddUint64(&totalEntries, uint64(prev.Entries))
		work <- prev
		prev = meta
	}
	close(work)
	wg.Wait()
	return nil
}

// Global flags (not command-specific)
var cpuprofile, memprofile, listen, traceFn string

func init() {
	// TODO: remove in favor of running as a test
	flag.StringVar(&cpuprofile, "cpuprofile", "", "")
	flag.StringVar(&memprofile, "memprofile", "", "write memory profile to this file")
	flag.StringVar(&listen, "listen", "", "speak HTTP on this [host]:port if non-empty")
	flag.StringVar(&traceFn, "trace", "", "create runtime/trace file")
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatalf("usage: %v <Debian Code Search index>", os.Args[0])
	}

	if listen != "" {
		go func() {
			if err := http.ListenAndServe(listen, nil); err != nil {
				log.Fatal(err)
			}
		}()
	}

	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if memprofile != "" {
		defer func() {
			f, err := os.Create(memprofile)
			if err != nil {
				log.Fatal(err)
			}
			runtime.GC()
			pprof.WriteHeapProfile(f)
			f.Close()
		}()
	}

	if traceFn != "" {
		f, err := os.Create(traceFn)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		if err := trace.Start(f); err != nil {
			log.Fatal(err)
		}
		defer trace.Stop()
	}

	start := time.Now()
	if err := logic(flag.Arg(0)); err != nil {
		log.Fatal(err)
	}
	val := atomic.LoadUint64(&totalBytes)
	log.Printf("rate: %.2f bytes/s", float64(val)/float64(time.Since(start).Nanoseconds())*float64(time.Second))
	log.Printf("total: %d entries", atomic.LoadUint64(&totalEntries))
}
