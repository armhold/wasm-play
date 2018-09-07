// +build wasm

package main

import (
	"log"
	"syscall/js"
)

type AudioContext struct {
	mediaDevices js.Value
	userMedia    js.Value
	mediaStream  js.Value
	audioHandler AudioHandler
}

type AudioHandler func(audioData []int)

func NewAudioContext(audioHandler AudioHandler) (*AudioContext, error) {
	// from https://developers.google.com/web/fundamentals/media/recording-audio/
	//
	// simulate:
	// navigator.mediaDevices.getUserMedia({ audio: true, video: false }).then(handleSuccess);

	navigator := js.Global().Get("navigator")
	log.Printf("navigator: %+v", navigator)

	ac := &AudioContext{audioHandler: audioHandler}

	ac.mediaDevices = navigator.Get("mediaDevices")
	log.Printf("mediaDevices: %+v", ac.mediaDevices)

	ac.userMedia = ac.mediaDevices.Call("getUserMedia",
		map[string]interface{}{
			"audio": true,
			"video": false,
		})

	log.Printf("userMedia: %+v", ac.userMedia)

	cb := js.NewCallback(func(args []js.Value) {
		log.Printf("got callback: %+v", args)

		ac.mediaStream = args[0]
		log.Printf("mediaStream: %+v", ac.mediaStream)

		ac.handleSuccess(ac.mediaStream)
	})

	ac.userMedia.Call("then", cb)

	log.Printf("exit setupAudio()")

	return ac, nil
}

func (ac *AudioContext) handleSuccess(stream js.Value) {
	log.Printf("handleSuccess: stream: %+v", stream)

	c := js.Global().Get("AudioContext").New()
	log.Printf("AudioContext: %+v", ac)

	source := c.Call("createMediaStreamSource", stream)

	log.Printf("source: %+v", source)

	processor := c.Call("createScriptProcessor", 1024, 1, 1)
	log.Printf("processor: %+v", processor)

	source.Call("connect", processor)
	result := processor.Call("connect", c.Get("destination"))
	log.Printf("result: %+v", result)

	processor.Set("onaudioprocess", js.NewCallback(func(args []js.Value) {
		// TODO: do something with the data, e.g. convert to WAV
		e := args[0]
		inputBuffer := e.Get("inputBuffer")

		d := inputBuffer.Call("getChannelData", 0)

		data := make([]int, inputBuffer.Length())

		for i := 0; i < inputBuffer.Length(); i++ {
			data[i] = d.Index(i).Int()
		}

		ac.audioHandler(data)
	}))
}
