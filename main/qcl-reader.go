package main

import (
	"flag"
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"log"
	qcl "qcl-reader"
)

func main() {
	var test bool
	flag.BoolVar(&test, "test", false, "use a random number generator instead of a live feed")
	flag.Parse()

	myqcl := qcl.QCL{}
	socket, err := zmq.NewSocket(zmq.PUB)
	if err != nil {
		log.Fatal(err)
	}
	defer socket.Close()
	socket.Bind("tcp://*:5550")

	cs := make(chan string)
	go myqcl.Sampler(test, cs)
	for {
		sample := <-cs
		fmt.Println(sample)
		socket.Send(sample, 0)
	}
}
