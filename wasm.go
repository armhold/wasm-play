// +build wasm

package main

import (
	"fmt"
	"log"
	"syscall/js"
	"time"
)

var (
	count           = 0
	recordingButton js.Value
)

// to compile:
//
// GOARCH=wasm GOOS=js go build -o example.wasm
//
func main() {
	log.Printf("hello world!")

	recordingButton = js.Global().Get("document").Call("getElementById", "recordingButton")

	recordingButton.Set("onclick", js.NewCallback(func(args []js.Value) {
		s := recordingButton.Get("innerHTML").String()
		log.Printf("toggle got: %s", s)

		if s == "Record" {
			startRecording()
		} else {
			stopRecording()
		}
	}))

	tickChan := time.NewTicker(time.Millisecond * 1000).C

	_, err := NewAudioContext(processAudio)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {

		case <-tickChan:
			log.Println("Ticker ticked")

			t := time.Now()
			ts := t.Format("Mon Jan _2 15:04:05 2006")

			s := fmt.Sprintf("<p>Time is: %v</p>", ts)
			el := js.Global().Get("document").Call("getElementById", "foo")
			el.Set("innerHTML", s)
		}
	}
}

func startRecording() {
	recordingButton.Set("innerHTML", "Stop")
}

func stopRecording() {
	recordingButton.Set("innerHTML", "Record")
}

func processAudio(audioData []int) {
	log.Printf("processAudio()...")

	count++

	if count%100 == 0 {
		// just print out a few bytes
		maxToPrint := 128
		line := ""
		sep := ""
		for i := 0; i < maxToPrint; i++ {
			v := audioData[i]
			line = fmt.Sprintf("%s%s%d", line, sep, v)
			sep = ", "
		}
		log.Print(line)
	}
}
