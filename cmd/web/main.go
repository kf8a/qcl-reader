package main

import (
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"time"
	// "log/syslog"
	"net/http"
	"os"
)

type connection struct {
	ws   *websocket.Conn
	send chan []byte
	q    *qcl
}

func (c *connection) reader() {
	for message := range c.send {
		err := c.ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println(err)
			return
		}
	}
	c.ws.Close()
}

var upgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

func QclHandler(q *qcl, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := &connection{send: make(chan []byte), ws: ws, q: q}
	c.q.register <- c
	defer func() { c.q.unregister <- c }()
	c.reader()
}

type Datum struct {
	Time  time.Time
	Value float64
}

func SaveDataHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var data map[string]interface{}
	if err := decoder.Decode(&data); err != nil {
		log.Println(err)
		return
	}
	// save data
	f, err := os.OpenFile("data.json", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	if err := encoder.Encode(&data); err != nil {
		log.Println(err)
		return
	}
	f.Sync()

	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	var test bool
	flag.BoolVar(&test, "test", false, "use a random number generator instead of a live feed")
	flag.Parse()

	var host = "127.0.0.1"
	instrument := newQcl("tcp://" + host + ":5550")
	go instrument.read(test)

	r := mux.NewRouter()
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		QclHandler(instrument, w, r)
	})
	r.HandleFunc("/save", SaveDataHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))
	http.Handle("/", r)
	http.ListenAndServe("127.0.0.1:8080", nil)
}
