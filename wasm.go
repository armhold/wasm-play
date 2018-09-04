package main

import "log"

// to compile:
//
// GOARCH=wasm GOOS=js go build -o example.wasm wasm.go
//
func main() {
	log.Printf("hello world!")
}
