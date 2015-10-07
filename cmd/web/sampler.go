package main

import (
	// "fmt"
	qclReader "github.com/kf8a/qclreader"
	"github.com/montanaflynn/stats"
	"math"
	"time"
)

type sampler struct {
	measurement chan qclReader.Datum
	control     chan *connection
}

func newSampler() *sampler {
	return &sampler{
		control: make(chan *connection),
	}
}

func (s *sampler) Sample() {
	//register with qcl

	data := make([]stats.Coordinate, 100)
	startTime := time.Now()

	for {
		datum := <-s.measurement
		c := stats.Coordinate{float64(datum.Time.Sub(startTime)), datum.N2O_dry_ppm}
		data = append(data, c)
		_, _ = stats.LinearRegression(data)
	}
}

type Coordinate struct {
	X, Y float64
}

func Fit(data []Coordinate) (slope, r float64) {

	var sum_x, sum_y, sum_xx, sum_xy, sum_yy float64
	for _, datum := range data {
		sum_x = sum_x + datum.X
		sum_y = sum_y + datum.Y
		sum_xx = sum_xx + (datum.X * datum.X)
		sum_xy = sum_xy + (datum.X * datum.Y)
		sum_yy = sum_yy + (datum.Y * datum.Y)
	}
	count := float64(len(data))
	slope = (count*sum_xy - sum_x*sum_y) / (count*sum_xx - sum_x*sum_x)
	// b := (sum_y / count) - (slope*sum_x)/count

	// mean_y := sum_y / count

	r = (count*sum_xy - sum_x*sum_y) / math.Sqrt((count*sum_xx-sum_x*sum_x)*(count*sum_yy-sum_y*sum_y))

	return slope, r
}
