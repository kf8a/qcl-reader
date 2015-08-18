package main

import (
	qclReader "github.com/kf8a/qclreader"
)

type qcl struct {
	connections map[*connection]bool
	register    chan *connection
	unregister  chan *connection
	host        string
}

func newQcl(hostName string) *qcl {
	return &qcl{
		connections: make(map[*connection]bool),
		register:    make(chan *connection),
		unregister:  make(chan *connection),
		host:        hostName,
	}
}

func (q *qcl) read(test bool) {
	myqcl := qclReader.QCL{}

	cs := make(chan string)
	go myqcl.Sampler(test, cs)

	for {
		data := <-cs
		select {
		case c := <-q.register:
			q.connections[c] = true
		case c := <-q.unregister:
			q.connections[c] = false
		default:
			for c := range q.connections {
				select {
				case c.send <- []byte(data):
				default:
					delete(q.connections, c)
					close(c.send)
				}
			}
		}
	}
}
