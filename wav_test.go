package main

import (
	gowav "github.com/youpy/go-wav"
	"io/ioutil"
	"os"
	"testing"
)

func TestWrite(t *testing.T) {
	outfile, err := ioutil.TempFile("/tmp", "outfile")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		outfile.Close()
		os.Remove(outfile.Name())
	}()

	var numSamples uint32 = 2
	var numChannels uint16 = 2
	var sampleRate uint32 = 44100
	var bitsPerSample uint16 = 16

	writer := gowav.NewWriter(outfile, numSamples, numChannels, sampleRate, bitsPerSample)
	samples := make([]gowav.Sample, numSamples)

	samples[0].Values[0] = 32767
	samples[0].Values[1] = -32768
	samples[1].Values[0] = 123
	samples[1].Values[1] = -123

	err = writer.WriteSamples(samples)
	if err != nil {
		t.Fatal(err)
	}

	outfile.Close()
	file, err := os.Open(outfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		file.Close()
		os.Remove(outfile.Name())
	}()

	reader := gowav.NewReader(file)
	if err != nil {
		t.Fatal(err)
	}

	fmt, err := reader.Format()
	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, int(fmt.AudioFormat), gowav.AudioFormatPCM)
	assertEqual(t, fmt.NumChannels, numChannels)
	assertEqual(t, fmt.SampleRate, sampleRate)
	assertEqual(t, fmt.ByteRate, sampleRate*4)
	assertEqual(t, fmt.BlockAlign, numChannels*(bitsPerSample/8))
	assertEqual(t, fmt.BitsPerSample, bitsPerSample)

	samples, err = reader.ReadSamples()
	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, len(samples), 2)
	assertEqual(t, samples[0].Values[0], 32767)
	assertEqual(t, samples[0].Values[1], -32768)
	assertEqual(t, samples[1].Values[0], 123)
	assertEqual(t, samples[1].Values[1], -123)
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}
