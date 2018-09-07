package main

import (
	"fmt"
	"log"
)

var (
	count = 0
)

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
