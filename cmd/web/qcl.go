package main

import (
	"encoding/json"
	"github.com/kf8a/li820"
	qclReader "github.com/kf8a/qclreader"
	"log"
)

type qcl struct {
	connections map[*connection]bool
	register    chan *connection
	unregister  chan *connection
}

func newQcl() *qcl {
	return &qcl{
		connections: make(map[*connection]bool),
		register:    make(chan *connection),
		unregister:  make(chan *connection),
	}
}

func (q *qcl) setup(test bool) (cs chan qclReader.Datum, co2 chan li820.Datum) {
	myqcl := qclReader.QCL{}

	cs = make(chan qclReader.Datum)
	go myqcl.Sampler(test, cs)

	mylicor := li820.NewLicor("li820", "qcl", "/dev/ttyUSB1")
	co2 = make(chan li820.Datum)
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
		data := <-cs
		co2_data := <-co2

		data.CO2_ppm = co2_data.CO2

		sample, err := json.Marshal(data)
		if err != nil {
			log.Print(err)
		}

		select {
		case c := <-q.register:
			q.connections[c] = true
		case c := <-q.unregister:
			q.connections[c] = false
		default:
			for c := range q.connections {
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
