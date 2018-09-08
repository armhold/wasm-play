package main

import (
	"github.com/go-audio/wav"
	"log"
	"math"
	"os"
	"testing"
)

func TestWrite(t *testing.T) {
	output := "/Users/armhold/out.wav"
	var frequency float64 = 440
	var lengthInSeconds float64 = 5

	f, err := os.Create(output)
	if err != nil {
		log.Fatalf("error creating %s: %s", output, err)
	}
	defer f.Close()

	const sampleRate = 48000
	wavOut := wav.NewEncoder(f, sampleRate, 16, 1, 1)
	numSamples := int(sampleRate * lengthInSeconds)
	defer wavOut.Close()

	for i := 0; i < numSamples; i++ {
		fv := math.Sin(float64(i) / sampleRate * frequency * 2 * math.Pi)
		v := uint16(fv * 32767)
		wavOut.WriteFrame(v)
	}
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}
