package main

import (
	"fmt"
	"github.com/smvfal/wlcount/wlutils"
	"log"
	"net/rpc"
	"os"
	"sort"
)

func main() {

	var (
		reduceMap      map[int][]string
		tempMap        map[int][]string
		reducedMap     map[int]int
		call           *rpc.Call
		keys           []int // word lengths
		chunks         []string
		workers        []string
		clients        []*rpc.Client
		mappedLengths  []map[int][]string
		reducedLengths []map[int]int
		w              int // number of available workers
		files          []*os.File
	)

	// Get documents' filenames
	fileNames := readArgs()
	fmt.Println("Submitted file names:", fileNames)

	// Open input files
	files, err := openFiles(fileNames)
	if err != nil {
		log.Fatal(err)
	}
	// TODO: close the files

	// Get addresses of available workers
	workers = getWorkers()
	w = len(workers)

	// Initialize workers support data structures
	clients = make([]*rpc.Client, w)

	// Try to connect to the available workers using HTTP protocol
	for i := 0; i < w; i++ {
		client, err := rpc.DialHTTP("tcp", workers[i])
		if err != nil {
			log.Fatal("Error in dialing: ", err)
			// TODO: handle dialing error
		}
		defer func() {
			if err := client.Close(); err != nil {
				log.Fatal(err)
			}
		}()
		clients[i] = client
	}

	for i, file := range files {
		fmt.Printf("Counting file %s\n", fileNames[i])

		// Split the file into chunks
		chunks = wlutils.SplitFile(file, w)

		// Call remote procedures asynchronously
		println("Call Map")
		mappedLengths = make([]map[int][]string, w)
		done := make(chan *rpc.Call, w) // Async RPC call channel, sometimes does not work if cap(done)==1
		for i := 0; i < w; i++ {
			clients[i].Go("MapReduce.Map", chunks[i], &mappedLengths[i], done)
		}
		mergedMap := make(map[int][]string)
		for i := 0; i < w; i++ {
			call = <-done
			if call.Error != nil {
				log.Fatal("Error in MapReduce.Map: ", call.Error.Error())
			}
			tempMap = *(call.Reply.(*map[int][]string))
			//fmt.Println("Reducer result: ", tempMap)
			mergeMaps(&mergedMap, tempMap)
		}
		//fmt.Println("Merged Map: ", mergedMap)

		nk := 0 // number of keys
		for k := range mergedMap {
			keys = append(keys, k)
			nk++
		}

		fmt.Println("Number of keys:", nk)
		fmt.Println("Keys:", keys)
		sort.Ints(keys)
		fmt.Println("Sorted keys:", keys)

		m := nk / w
		r := nk % w
		fmt.Printf("m = %d, r = %d\n", m, r)
		var kRed int // number of keys per reducer
		println("Call Reduce")
		reducedLengths = make([]map[int]int, w)
		for i := 0; i < w; i++ {
			if i < r {
				kRed = m + 1
			} else {
				kRed = m
			}
			// initialize/clear map
			reduceMap = make(map[int][]string)
			for j := 0; j < kRed; j++ {
				k := keys[j]
				reduceMap[k] = mergedMap[k]
			}
			// reslice key set
			keys = keys[kRed:]
			//fmt.Println("Reducer Map:", reduceMap)
			clients[i].Go("MapReduce.Reduce", reduceMap, &reducedLengths[i], done)
		}

		resultMap := make(map[int]int)
		for i := 0; i < w; i++ {
			call = <-done
			if call.Error != nil {
				log.Fatal("Error in MapReduce.Reduce: ", call.Error.Error())
			}
			reducedMap = *(call.Reply.(*map[int]int))
			//fmt.Println("Reduced map:", reducedMap)
			for k, v := range reducedMap {
				resultMap[k] = v
			}
		}
		printResult(resultMap)
	}
}

func readArgs() []string {
	if len(os.Args) <= 1 {
		log.Fatal("No input file detected.\nUsage: wlcount file1 [file2 ...]")
	}
	return os.Args[1:]
}

func mergeMaps(dst *map[int][]string, temp map[int][]string) {
	for k, l := range temp {
		(*dst)[k] = append((*dst)[k], l...)
	}
}

func getWorkers() []string {

	fmt.Println("getting workers..")
	var workers []string

	file, err := os.Open("address.config")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	workers, err = wlutils.ScanStrings(file, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("workers:", workers)

	return workers
}

func openFiles(fileNames []string) ([]*os.File, error) {
	var files []*os.File
	for _, name := range fileNames {
		file, err := os.Open(name)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

func printResult(result map[int]int) {

	fmt.Println("---------------")
	fmt.Println("Length | Count ")
	fmt.Println("-------+-------")

	var keys []int
	for k := range result {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		fmt.Printf("%6d | %d\n", k, result[k])
	}

	fmt.Println("---------------")
}
