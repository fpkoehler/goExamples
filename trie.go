package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/dghubble/trie"
)

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage(memprofile string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)

	if memprofile != "" {
		f, err := os.Create(memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}
}

var totalBytes int64

func main() {
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
	var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
	var trieApproach = flag.Bool("trie", false, "use trie approach")
	var mapApproach = flag.Bool("map", false, "use map approach")

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error can get cwd")
	}

	count := 0

	if *trieApproach {
		/* Trie approach */
		pathTrie := trie.NewPathTrie()

		count = 0
		walkFn := func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				// fmt.Println(path)
				pathTrie.Put(path, info)
				count++
			}
			return nil
		}
		timeStart := time.Now()
		filepath.Walk(dir, walkFn)
		fmt.Println("Trie: Time adding", count, "files:", time.Now().Sub(timeStart))

		trieWalkFn := func(key string, value interface{}) error {
			//fmt.Println(key, value)
			//fmt.Println(key)
			totalBytes += value.(os.FileInfo).Size()
			return nil
		}

		PrintMemUsage(*memprofile)

		totalBytes = 0
		timeStart = time.Now()
		pathTrie.Walk(trieWalkFn)
		fmt.Println("Trie: total bytes", totalBytes, "Time getting files:", time.Now().Sub(timeStart))
	}

	if *mapApproach {
		/* Map approach */
		fileInfos := make(map[string]os.FileInfo)
		count = 0
		walkFn := func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				// fmt.Println(path)
				fileInfos[path] = info
				count++
			}
			return nil
		}
		timeStart := time.Now()
		filepath.Walk(dir, walkFn)
		fmt.Println("Map: Time adding", count, "files:", time.Now().Sub(timeStart))

		PrintMemUsage(*memprofile)

		totalBytes = 0
		timeStart = time.Now()
		for _, fileInfo := range fileInfos {
			totalBytes += fileInfo.Size()
		}
		fmt.Println("Map: total bytes", totalBytes, "Time getting files:", time.Now().Sub(timeStart))
	}
}
