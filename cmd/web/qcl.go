package main

import (
	"encoding/json"
	"github.com/kf8a/li820"
	qclReader "github.com/kf8a/qclreader"
	"log"
)

type qcl struct {
	connections    map[*connection]bool
	recordings     map[string]*connection
	register       chan *connection
	unregister     chan *connection
	dataConnection map[*connection]qclReader.Datum
}

func newQcl() *qcl {
	return &qcl{
		connections: make(map[*connection]bool),
		recordings:  make(map[string]*connection),
		register:    make(chan *connection),
		unregister:  make(chan *connection),
	}
}

func (q *qcl) setup(test bool) (cs chan qclReader.Datum, co2 chan li820.Datum) {
	myqcl := qclReader.QCL{}

	cs = make(chan qclReader.Datum, 10)
	go myqcl.Sampler(test, cs, "/dev/qcl")

	mylicor := li820.NewLicor("li820", "qcl", "/dev/licor")
	co2 = make(chan li820.Datum, 10)
	if test {
		go mylicor.TestSampler(co2)
	} else {
		go mylicor.Sampler(co2)
	}
	return
}

func (q *qcl) read(test bool) {

	cs, co2 := q.setup(test)

	for {
		// log.Println("before data")
		data := <-cs
		// log.Println("after qcl")
		co2_data := <-co2
		// log.Println("after licor")

		data.CO2_ppm = co2_data.CO2

		// log.Println(data)

		sample, err := json.Marshal(data)
		if err != nil {
			log.Print(err)
		} else {
			publish("measurement", sample)
		}

		select {
		case c := <-q.register:
			q.connections[c] = true
		case c := <-q.unregister:
			q.connections[c] = false
		default:
			for c := range q.connections {
				if !q.connections[c] {
					delete(q.connections, c)
					continue
				}
				select {
				case c.send <- []byte(sample):
				default:
					delete(q.connections, c)
					close(c.send)
				}
			}
		}
	}
}
