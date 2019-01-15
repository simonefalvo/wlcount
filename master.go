package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strings"
)

func main() {

	const n = 2;  // number of workers
	var maps [n]map[int][]string	// maps will store the RPC results
	s := "hello gophers! How are you?\nciao \"geomidi\", come state?"
	var c = make(chan *rpc.Call, 1)
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
		fmt.Println(call.Reply)
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
