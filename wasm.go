// +build wasm

package main

import (
	"fmt"
	"log"
	"syscall/js"
	"time"
)

// to compile:
//
// GOARCH=wasm GOOS=js go build -o example.wasm
//
func main() {
	log.Printf("hello world!")

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
