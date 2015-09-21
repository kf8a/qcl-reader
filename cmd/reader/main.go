package main

import (
	"encoding/json"
	"flag"
	"fmt"
	qcl "github.com/kf8a/qclreader"
	"log"
)

func main() {
	var test bool
	flag.BoolVar(&test, "test", false, "use a random number generator instead of a live feed")
	flag.Parse()

	myqcl := qcl.QCL{}

	cs := make(chan qcl.Datum)
	go myqcl.Sampler(test, cs)
	for {
		data := <-cs
		sample, err := json.Marshal(data)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(sample))
	}
}
