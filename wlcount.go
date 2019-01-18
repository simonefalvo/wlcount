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

	// Get documents' filenames
	fileNames := readArgs()

	// Get addresses of available workers
	workers := getWorkers()
	w := len(workers) // number of available workers

	// Try to connect to the available workers using HTTP protocol
	clients := make([]*rpc.Client, w)
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

	for _, fileName := range fileNames {

		fmt.Printf("Counting file %s\n", fileName)

		// Open the file
		file, err := os.Open(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				log.Fatal(err)
			}
		}()

		// Split the file into w chunks
		chunks := wlutils.SplitFile(file, w)

		// Call remote procedures asynchronously
		mapReplies := make([]map[int][]string, w) // result array of the map function
		done := make(chan *rpc.Call, w)           // Async RPC call channel, sometimes does not work if cap(done)==1
		fmt.Println("Mapping...")
		for i := 0; i < w; i++ {
			clients[i].Go("MapReduce.Map", chunks[i], &mapReplies[i], done)
		}
		mergedMap := make(map[int][]string) // aggregated mapping results
		for i := 0; i < w; i++ {
			call := <-done
			if call.Error != nil {
				log.Fatal("Error in MapReduce.Map: ", call.Error.Error())
			}
			tempMap := *(call.Reply.(*map[int][]string))
			//fmt.Println("Reducer result: ", tempMap)
			// aggregate several worker replies
			wlutils.MergeMaps(&mergedMap, tempMap)
		}
		//fmt.Println("Merged Map: ", mergedMap)

		nk := 0 // number of wordLengths
		var wordLengths []int
		for k := range mergedMap {
			wordLengths = append(wordLengths, k)
			nk++
		}
		sort.Ints(wordLengths)
		fmt.Println("Sorted word lengths:", wordLengths)

		m := nk / w
		r := nk % w
		// fmt.Printf("m = %d, r = %d\n", m, r)
		var kRed int // number of wordLengths per reducer
		println("Call Reduce")
		reducedLengths := make([]map[int]int, w)
		for i := 0; i < w; i++ {
			// Evaluate how many key-lengths the i-th worker has to work on
			if i < r {
				kRed = m + 1
			} else {
				kRed = m
			}
			// Initialize/clear map to reduce
			reduceMap := make(map[int][]string)
			for j := 0; j < kRed; j++ {
				k := wordLengths[j]
				reduceMap[k] = mergedMap[k]
			}
			// Re-slice key set
			wordLengths = wordLengths[kRed:]
			//fmt.Println("Reducer Map:", reduceMap)
			clients[i].Go("MapReduce.Reduce", reduceMap, &reducedLengths[i], done)
		}

		resultMap := make(map[int]int)
		for i := 0; i < w; i++ {
			call := <-done
			if call.Error != nil {
				log.Fatal("Error in MapReduce.Reduce: ", call.Error.Error())
			}
			workerReduction := *(call.Reply.(*map[int]int))
			//fmt.Println("Reduced map:", workerReduction)
			for k, v := range workerReduction {
				resultMap[k] = v
			}
		}
		printResult(resultMap)
	}
}

// Read command line arguments
func readArgs() []string {
	if len(os.Args) <= 1 {
		log.Fatal("No input file detected.\nUsage: ./wlcount file1 [file2 ...]")
	}
	return os.Args[1:]
}

// Read the list of the available workers on the configuration file address.config
// and return the string format addresses
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

	// Read lines of the file
	workers, err = wlutils.ScanStrings(file, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("workers:", workers)

	return workers
}

func printResult(result map[int]int) {

	fmt.Println("---------------")
	fmt.Println("Length | Count ")
	fmt.Println("-------+-------")

	var wordLengths []int
	for k := range result {
		wordLengths = append(wordLengths, k)
	}
	sort.Ints(wordLengths)
	for _, k := range wordLengths {
		fmt.Printf("%6d | %d\n", k, result[k])
	}

	fmt.Println("---------------")
}
