package qclreader

import (
	"encoding/csv"
	serial "github.com/tarm/serial"
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

func (qcl QCL) RandomSampler(cs chan Datum) {
	for {
		time.Sleep(1 * time.Second)

		datum := Datum{
			ObsTime:     time.Now(),
			Time:        time.Now(),
			CH4_ppm:     rand.Float64(),
			H2O_ppm:     rand.Float64(),
			N2O_ppm:     rand.Float64(),
			N2O_dry_ppm: rand.Float64(),
			CH4_dry_ppm: rand.Float64(),
		}

		cs <- datum
	}
}

//Sampler is a convenience function to allow selection of test or real samplers
func (qcl QCL) Sampler(test bool, cs chan Datum, port string) {
	if test {
		go qcl.RandomSampler(cs)
	} else {
		go qcl.RealSampler(cs, port)
	}
}

func (qcl QCL) RealSampler(cs chan Datum, connection_string string) {
	c := serial.Config{Name: connection_string, Baud: 9600}

	for {
		port, err := serial.OpenPort(&c)

		if err != nil {
			log.Println(err)
			time.Sleep(2 * time.Second)
			continue
		}
		defer port.Close()
		qcl.port = port

		for {
			reader := csv.NewReader(qcl.port)
			line, err := reader.Read()
			if err != nil {
				log.Println(err)
				continue
			}

			log.Println(line)

			if len(line) < 10 {
				log.Println("short line", line)
				continue
			}
			datum := Datum{
				ObsTime:     time.Now(),
				Time:        qcl.parseTime(line[0]),
				CH4_ppm:     qcl.parseFloat(line[1]),
				H2O_ppm:     qcl.parseFloat(line[3]),
				N2O_ppm:     qcl.parseFloat(line[5]),
				CO2_ppm:     qcl.parseFloat("0"),
				N2O_dry_ppm: qcl.parseFloat(line[7]),
				CH4_dry_ppm: qcl.parseFloat(line[9]),
			}

			cs <- datum
		}
	}
}
