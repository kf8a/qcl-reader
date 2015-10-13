package main

import (
	"encoding/gob"
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

var store = sessions.NewCookieStore([]byte("qcl-error-code"))

//QclHandler handles a new connection, creates and registers a new connection to the QCL reader
func QclHandler(q *qcl, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	c := &connection{send: make(chan []byte), ws: ws, q: q}

	c.q.register <- c
	// defer func() { c.q.unregister <- c }()
	c.reader()
	c.q.unregister <- c
}

type Recording struct {
	StartedAt     time.Time
	EndedAt       time.Time
	Canceled      bool
	SampleId      string
	UserId        string
	data          chan []byte
	Plot          string
	ChamberHeight float64
	FluxData      map[string]interface{}
}

type Chamber struct {
	Plot   string
	Height float64
}

func RecordHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "qcl-session")
	if user_id, ok := session.Values["user_id"].(string); ok {

		// generate a uuid and save it in the session
		// set the flag that this session is now recording
		sample_id := uuid.NewV4().String()
		session.Values["sample_id"] = sample_id
		// send out a start recording message with the user id and the sample_id and treatment and height

		// get form fields
		decoder := json.NewDecoder(r.Body)
		var chamber Chamber
		err := decoder.Decode(&chamber)
		if err != nil {
			log.Fatal(err)
		}

		recording := &Recording{
			StartedAt:     time.Now(),
			SampleId:      sample_id,
			UserId:        user_id,
			Plot:          chamber.Plot,
			ChamberHeight: chamber.Height,
		}
		session.Values["recording"] = recording
		sample, err := json.Marshal(recording)
		if err != nil {
			log.Print(err)
		} else {
			publish("control", sample)
		}

	} else {
		log.Println("ERROR: Record Handler no user")
	}
	session.Save(r, w)
}

func CancelHandler(w http.ResponseWriter, r *http.Request) {
	var recording Recording
	session, _ := store.Get(r, "qcl-session")
	if user_id, ok := session.Values["user_id"].(string); ok {
		if user_id == "" {
			log.Println("ERROR: SaveDataHandler no user")
		}
		recording, ok = session.Values["recording"].(Recording)
		if ok {
			recording.EndedAt = time.Now()
			recording.Canceled = true

			session.Values["recording"] = recording
			sample, err := json.Marshal(recording)
			if err != nil {
				log.Fatal(err)
			} else {
				publish("control", sample)
			}
		} else {
			log.Fatal(ok)
		}
	}
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func SaveDataHandler(w http.ResponseWriter, r *http.Request) {
	var recording Recording
	session, _ := store.Get(r, "qcl-session")
	if user_id, ok := session.Values["user_id"].(string); ok {
		if user_id == "" {
			log.Println("ERROR: SaveDataHandler no user")
		}
		recording, ok = session.Values["recording"].(Recording)
		if ok {
			recording.EndedAt = time.Now()
			recording.Canceled = false

			sample, err := json.Marshal(recording)
			if err != nil {
				log.Fatal(err)
			} else {
				publish("control", sample)
			}
		} else {
			log.Fatal("NO recording found")
		}
	} else {
		log.Fatal("NO user found")
	}

	decoder := json.NewDecoder(r.Body)

	var data map[string]interface{}
	if err := decoder.Decode(&data); err != nil {
		log.Println(err)
		return
	}
	recording.FluxData = data

	// save data
	f, err := os.OpenFile("data.json", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	if err := encoder.Encode(&recording); err != nil {
		log.Fatal(err)
	}
	// if err := encoder.Encode(&data); err != nil {
	// 	log.Fatal(err)
	// }
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

func init() {
	gob.Register(Recording{})
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
	r.HandleFunc("/cancel", CancelHandler)
	fileHandler := http.FileServer(http.Dir("./public/"))
	r.PathPrefix("/").Handler(MyServeFileHandler(fileHandler))
	// r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))
	http.Handle("/", r)
	// http.ListenAndServe(":8080", nil)
	http.ListenAndServe("127.0.0.1:8080", nil)
}
