// +build wasm

package main

import (
	"fmt"
	"log"
	"os"
	"syscall/js"
	"time"
)

var (
	count = 0
)

// to compile:
//
// GOARCH=wasm GOOS=js go build -o example.wasm
//
func main() {
	log.Printf("hello world!")

	tickChan := time.NewTicker(time.Millisecond * 1000).C

	outfile, err := os.Create("/Users/armhold/audio.wav")
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()

	_, err = NewAudioContext(processAudio)
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
