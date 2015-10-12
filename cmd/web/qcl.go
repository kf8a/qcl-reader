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

	log.Println("read called")
	cs, co2 := q.setup(test)

	for {
		data := <-cs
		co2_data := <-co2

		data.CO2_ppm = co2_data.CO2

		sample, err := json.Marshal(data)
		if err != nil {
			log.Print(err)
		} else {
			publish("measurement", sample)
		}

		select {
		case c := <-q.register:
			log.Println("registering connection")
			log.Println(c)
			q.connections[c] = true
		case c := <-q.unregister:
			log.Println("unregistering connection")
			log.Println(c)
			log.Println(q.connections)
			q.connections[c] = false
		default:
			log.Println("Current connections")
			log.Println(q.connections)
			for c := range q.connections {
				log.Println("processing")
				log.Println(q.connections[c])
				if !q.connections[c] {
					log.Println("delete before")
					delete(q.connections, c)
					close(c.send)
					continue
				}
				log.Println(q.connections[c])
				select {
				case c.send <- []byte(sample):
					log.Println("send")
				default:
					log.Println("delete")
					delete(q.connections, c)
					close(c.send)
				}
			}
		}
	}
}
