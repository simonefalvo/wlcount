package main

import (
	"github.com/smvfal/wlcount/mapreduce"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

func main() {

	mr := new(mapreduce.MapReduce)

	// Register a new rpc server and the struct created above.
	server := rpc.NewServer()
	err := server.RegisterName("MapReduce", mr)
	if err != nil {
		log.Fatal("Format of service MapReduce is not correct: ", err)
	}
	// Register an HTTP handler for RPC messages on rpcPath, and a debugging handler on debugPath
	server.HandleHTTP("/", "/debug")

	// Listen for incoming messages on automatic port
	l, e := net.Listen("tcp", ":0")
	if e != nil {
		log.Fatal("Listen error: ", e)
	}

	// Register the address to a configuration file
	registerAddress(l.Addr())

	// Start go's http server on socket specified by l
	err = http.Serve(l, nil)
	if err != nil {
		log.Fatal("Serve error: ", err)
	}
}

func registerAddress(addr net.Addr) {
	f, err := os.OpenFile("address.config", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err = f.WriteString(addr.String() + "\n"); err != nil {
		log.Fatal(err)
	}
	if err = f.Close(); err != nil {
		log.Fatal(err)
	}
}
