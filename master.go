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

	const n = 2                  // number of workers
	var maps [n]map[int][]string // maps will store the RPC results
	var reducedMaps [n]map[int]int
	var mergedMap = make(map[int][]string)
	var resultMap = make(map[int]int)
	var reduceMap map[int][]string
	var tempMap map[int][]string
	var reducedMap map[int]int
	s := "Hello gophers! How are you?\nCiao \"geomidi\", come state?"
	var c = make(chan *rpc.Call, 1)
	var keys []int // lenghts
	var call *rpc.Call
	var chunks []string

	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		chunks = append(chunks, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading string:", err)
	}

	// Try to connect to localhost:1234 using HTTP protocol (the port on which RPC server is listening)
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Error in dialing: ", err)
	}
	defer client.Close()

	// Call remote procedures asynchronously
	for i := 0; i < n; i++ {
		client.Go("MapReduce.Map", chunks[i], &maps[i], c)
	}
	for i := 0; i < n; i++ {
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

	m := nk / n
	r := nk % n
	fmt.Printf("m = %d, r = %d\n", m, r)
	var kred int // number of keys per reducer
	for i := 0; i < m; i++ {
		if i < r {
			kred = m + 1
		} else {
			kred = m
		}
		// initialize/clear map
		reduceMap = make(map[int][]string)
		for j := 0; j < kred; j++ {
			k := keys[j]
			reduceMap[k] = mergedMap[k]
		}
		// reslice key set
		keys = keys[kred :]
		fmt.Println("Reducer Map:", reduceMap)
		client.Go("MapReduce.Reduce", reduceMap, &reducedMaps[i], c)
	}

	for i := 0; i < n; i++ {
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

/*
type Call struct {
        ServiceMethod string      // The name of the service and method to call.
        Args          interface{} // The argument to the function (*struct).
        Reply         interface{} // The reply from the function (*struct).
        Error         error       // After completion, the error status.
        Done          chan *Call  // Strobes when call is complete.
}
*/
