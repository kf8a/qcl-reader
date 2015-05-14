package main

import (
	"flag"
	"fmt"
	qcl "github.com/kf8a/qclreader"
)

func main() {
	var test bool
	flag.BoolVar(&test, "test", false, "use a random number generator instead of a live feed")
	flag.Parse()

	myqcl := qcl.QCL{}

	cs := make(chan string)
	go myqcl.Sampler(test, cs)
	for {
		sample := <-cs
		fmt.Println(sample)
	}
}
