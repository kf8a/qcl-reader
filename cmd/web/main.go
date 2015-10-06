package main

import (
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	// "github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"os"
	"time"
)

type connection struct {
	ws        *websocket.Conn
	send      chan []byte
	q         *qcl
	recording string
}

var connections = make(map[string]*connection)

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

var store = sessions.NewCookieStore([]byte("qcl-error-code"))

//QclHandler handles a new connection, creates and registers a new connection to the QCL reader
func QclHandler(q *qcl, w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "qcl-session")
	if user_id, ok := session.Values["user_id"].(string); ok {

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		c := &connection{send: make(chan []byte), ws: ws, q: q}
		connections[user_id] = c

		c.q.register <- c
		defer func() { c.q.unregister <- c }()
		c.reader()
	} else {
		log.Println("ERROR: no user")
		return
	}
}

type Recording struct {
	At       time.Time
	Type     string
	SampleId string
	UserId   string
}

func RecordHandler(w http.ResponseWriter, r *http.Request) {
	// generate a uuid and save it in the session
	// set the flag that this session is now recording
	session, _ := store.Get(r, "qcl-session")
	if user_id, ok := session.Values["user_id"].(string); ok {
		sample_id := uuid.NewV4().String()
		session.Values["sample_id"] = sample_id

		// send out a start recording message with the user id and the sample_id and treatment and height
		data := &Recording{
			At:       time.Now(),
			Type:     "start",
			SampleId: sample_id,
			UserId:   user_id,
		}
		sample, err := json.Marshal(data)
		if err != nil {
			log.Print(err)
		} else {
			publish("control", sample)
		}

	} else {
		log.Println("ERROR: Record Handler no user")
	}

}

func SaveDataHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "qcl-session")
	if user_id, ok := session.Values["user_id"].(string); ok {
		if user_id == "" {
			log.Println("ERROR: SaveDataHandler no user")
		}
		if sample_id, ok := session.Values["session_id"].(string); ok {
			// connection.send({'user_id': user_id, 'sample_id': sample_id, 'event': 'stopRecording'})
			data := &Recording{
				At:       time.Now(),
				Type:     "stop",
				SampleId: sample_id,
				UserId:   user_id,
			}
			sample, err := json.Marshal(data)
			if err != nil {
				log.Print(err)
			} else {
				publish("control", sample)
			}
			log.Println(sample_id)
		}
	}

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

func MyServeFileHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "qcl-session")
		if _, ok := session.Values["user_id"].(string); !ok {
			user_id := uuid.NewV4().String()
			session.Values["user_id"] = user_id
		}
		err := session.Save(r, w)
		if err != nil {
			log.Println(err)
		}
		h.ServeHTTP(w, r)
	})
}

func main() {
	var test bool
	flag.BoolVar(&test, "test", false, "use a random number generator instead of a live feed")
	flag.Parse()

	instrument := newQcl()
	go instrument.read(test)

	r := mux.NewRouter()
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		QclHandler(instrument, w, r)
	})
	r.HandleFunc("/save", SaveDataHandler)
	r.HandleFunc("/record", RecordHandler)
	fileHandler := http.FileServer(http.Dir("./public/"))
	r.PathPrefix("/").Handler(MyServeFileHandler(fileHandler))
	// r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))
	http.Handle("/", r)
	// http.ListenAndServe(":80", nil)
	http.ListenAndServe("127.0.0.1:8080", nil)
}
