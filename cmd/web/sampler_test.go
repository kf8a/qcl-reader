package main

import (
	"math"
	"testing"
)

func TestSamplerCreation(t *testing.T) {
	s := newSampler()
	if s == nil {
		t.Error("Can not create sampler")
	}
}

func TestFit(t *testing.T) {
	var fitTest = []struct {
		input []Coordinate
		slope float64
		r     float64
	}{
		{input: []Coordinate{{1, 1}, {2, 2}, {3, 3}}, slope: 1.0, r: 1.0},
		{input: []Coordinate{{1, 2}, {2, 4}, {3, 6}}, slope: 2.0, r: 1.0},
		{input: []Coordinate{{1, 1}, {2, 1}, {3, 5}}, slope: 2.0, r: 0.866025},
	}
	for _, fits := range fitTest {
		slope, r := Fit(fits.input)
		if fits.slope != slope {
			t.Errorf("expected Fit to return a slope of %f but it returned %f", fits.slope, slope)
		}
		if math.Abs(fits.r-r) > 0.0001 {
			t.Errorf("expected Fit to return an r of %f but it returned %f", fits.r, r)
		}
	}
}

func TestWorkflow(t *testing.T) {
	// create samples
	s := newSampler()
	// feed some data in
	go s.Sample()
	// compute fit
}
