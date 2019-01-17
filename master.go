package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"sort"
	"strings"
)

func main() {

	var (
		s              = "Hello gophers! How are you?\nCiao \"geomidi\", come state?"
		mergedMap      = make(map[int][]string)
		resultMap      = make(map[int]int)
		reduceMap      map[int][]string
		tempMap        map[int][]string
		reducedMap     map[int]int
		call           *rpc.Call
		c              = make(chan *rpc.Call, 1) // Async RPC call channel
		keys           []int                     // word lengths
		chunks         []string
		workers        []string
		clients        []*rpc.Client
		mappedLengths  []map[int][]string
		reducedLengths []map[int]int
		w              int // number of available workers
	)

	// TODO: get filename from stdin
	// TODO: open input file

	// Get addresses of available workers
	workers = getWorkers()
	w = len(workers)

	// Initialize workers support data structures
	clients = make([]*rpc.Client, w)
	mappedLengths = make([]map[int][]string, w)
	reducedLengths = make([]map[int]int, w)

	// Try to connect to the available workers using HTTP protocol
	for i := 0; i < w; i++ {
		client, err := rpc.DialHTTP("tcp", workers[i])
		if err != nil {
			log.Fatal("Error in dialing: ", err)
		}
		clients[i] = client
		defer clients[i].Close()
	}

	// TODO: Split file into chunks
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		chunks = append(chunks, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading string:", err)
	}

	// Call remote procedures asynchronously
	for i := 0; i < w; i++ {
		clients[i].Go("MapReduce.Map", chunks[i], &mappedLengths[i], c)
	}
	for i := 0; i < w; i++ {
		call = <-c
		if call.Error != nil {
			log.Fatal("Error in MapReduce.Map: ", call.Error.Error())
		}
		tempMap = *(call.Reply.(*map[int][]string))
		fmt.Println("Reducer result: ", tempMap)
		mergeMaps(&mergedMap, tempMap)
	}
	fmt.Println("Merged Map: ", mergedMap)

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
	for i := 0; i < m; i++ {
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
		fmt.Println("Reducer Map:", reduceMap)
		clients[i].Go("MapReduce.Reduce", reduceMap, &reducedLengths[i], c)
	}

	for i := 0; i < w; i++ {
		call = <-c
		if call.Error != nil {
			log.Fatal("Error in MapReduce.Reduce: ", call.Error.Error())
		}
		reducedMap = *(call.Reply.(*map[int]int))
		fmt.Println("Reduced map:", reducedMap)
		for k, v := range reducedMap {
			resultMap[k] = v
		}
	}
	fmt.Println("Result Map: ", resultMap)
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
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		workers = append(workers, scanner.Text())
		fmt.Println(workers)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading workers addresses:", err)
	}

	return workers
}

/*
type Call struct {
        ServiceMethod string      // The name of the service and method to call.
        Args          interface{} // The argument to the function (*struct).
        Reply         interface{} // The reply from the function (*struct).
        Error         error       // After completion, the error status.
        Done          chan *Call  // Strobes when call is complete.
}
*/
