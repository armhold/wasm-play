package main

import (
	"fmt"
	"log"
	"syscall/js"
	"time"
)

// to compile:
//
// GOARCH=wasm GOOS=js go build -o example.wasm wasm.go
//
func main() {
	log.Printf("hello world!")

	tickChan := time.NewTicker(time.Millisecond * 1000).C

	setupAudio()

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

func setupAudio() {
	log.Printf("enter setupAudio()...")

	// https://developers.google.com/web/fundamentals/media/recording-audio/
	navigator := js.Global().Get("navigator")
	log.Printf("navigator: %+v", navigator)

	mediaDevices := navigator.Get("mediaDevices")
	log.Printf("mediaDevices: %+v", mediaDevices)

	userMedia := mediaDevices.Call("getUserMedia",
		map[string]interface{}{
			"audio": true,
			"video": false,
		})

	log.Printf("userMedia: %+v", userMedia)

	cb := js.NewCallback(func(args []js.Value) {
		log.Printf("got callback: %+v", args)

		mediaStream := args[0]
		log.Printf("mediaStream: %+v", mediaStream)

		handleSuccess(mediaStream)
	})

	userMedia.Call("then", cb)
	//navigator.mediaDevices.getUserMedia({ audio: true, video: false }).then(handleSuccess);

	log.Printf("exit setupAudio()")
}

func handleSuccess(stream js.Value) {
	log.Printf("handleSuccess: stream: %+v", stream)

	ac := js.Global().Get("AudioContext").New()
	log.Printf("AudioContext: %+v", ac)

	source := ac.Call("createMediaStreamSource", stream)

	log.Printf("source: %+v", source)

	processor := ac.Call("createScriptProcessor", 1024, 1, 1)
	log.Printf("processor: %+v", processor)

	source.Call("connect", processor)
	result := processor.Call("connect", ac.Get("destination"))
	log.Printf("result: %+v", result)

	processor.Set("onaudioprocess", js.NewCallback(func(args []js.Value) {
		// Do something with the data, i.e Convert this to WAV
		e := args[0]
		inputBuffer := e.Get("inputBuffer")
		log.Printf("audio data: %+v", inputBuffer)
	}))

}
