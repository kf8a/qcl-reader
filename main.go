package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	zmq "github.com/pebbe/zmq4"
	serial "github.com/tarm/goserial"
	"io"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type QCL struct {
	port io.ReadWriteCloser
}

type Datum struct {
	ObsTime     time.Time `json:"obs_time"`
	Time        time.Time `json:"time"`
	CH4_ppm     float64
	H2O_ppm     float64
	N2O_ppm     float64
	CO2_ppm     float64
	CH4_dry_ppm float64
	N2O_dry_ppm float64
}

func (qcl QCL) parseFloat(value string) float64 {
	number, err := strconv.ParseFloat(strings.Trim(value, " "), 64)
	if err != nil {
		log.Print(err)
		return 0
	} else {
		return number
	}
}

func (qcl QCL) parseTime(value string) time.Time {
	loc, _ := time.LoadLocation("America/Detroit")
	layout := "2006/01/02 15:04:05"
	datetime, err := time.ParseInLocation(layout, strings.Trim(value, " "), loc)
	if err != nil {
		log.Print(err)
		return time.Now()
	} else {
		return datetime
	}
}

func (qcl QCL) RandomSample() string {
	time.Sleep(1 * time.Second)

	datum := Datum{
		ObsTime:     time.Now(),
		Time:        time.Now(),
		CH4_ppm:     rand.Float64(),
		H2O_ppm:     rand.Float64(),
		N2O_ppm:     rand.Float64(),
		CO2_ppm:     rand.Float64(),
		N2O_dry_ppm: rand.Float64(),
		CH4_dry_ppm: rand.Float64(),
	}
	b, err := json.Marshal(datum)
	if err != nil {
		fmt.Println("error:", err)
	}

	return string(b)
}

func (qcl QCL) Sample() string {
	c := serial.Config{Name: "/dev/ttyUSB0", Baud: 9600}
	port, err := serial.OpenPort(&c)
	qcl.port = port
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(qcl.port)
	line, err := reader.Read()
	if err != nil {
		log.Fatal(err)
	}
	datum := Datum{
		ObsTime:     time.Now(),
		Time:        qcl.parseTime(line[0]),
		CH4_ppm:     qcl.parseFloat(line[1]),
		H2O_ppm:     qcl.parseFloat(line[3]),
		N2O_ppm:     qcl.parseFloat(line[5]),
		N2O_dry_ppm: qcl.parseFloat(line[7]),
		CH4_dry_ppm: qcl.parseFloat(line[9]),
	}
	b, err := json.Marshal(datum)
	if err != nil {
		fmt.Println("error:", err)
	}

	return string(b)
}

func (qcl QCL) parse(data string) string {
	return data
}

func main() {
	var test bool
	flag.BoolVar(&test, "test", false, "use a random number generator instead of a live feed")
	flag.Parse()

	qcl := QCL{}
	socket, err := zmq.NewSocket(zmq.PUB)
	if err != nil {
		log.Fatal(err)
	}
	defer socket.Close()
	socket.Bind("tcp://*:5550")

	sampler := qcl.Sample
	if test {
		sampler = qcl.RandomSample
	}
	for {
		sample := sampler()
		log.Print(sample)
		socket.Send(sample, 0)
	}
}
