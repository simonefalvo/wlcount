package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/smvfal/wlcount/mapreduce"
)

func main() {

	word := new(mapreduce.Word)

	// Register a new rpc server and the struct created above.
	server := rpc.NewServer()
	err := server.RegisterName("MapReduce", word)
	if err != nil {
		log.Fatal("Format of service Arith is not correct: ", err)
	}
	// Register an HTTP handler for RPC messages on rpcPath, and a debugging handler on debugPath
	server.HandleHTTP("/", "/debug")

	// Listen for incoming messages on port 1234
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("Listen error: ", e)
	}

	// Start go's http server on socket specified by l
	err = http.Serve(l, nil)
	if err != nil {
		log.Fatal("Serve error: ", err)
	}
}