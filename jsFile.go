package main

import (
	_ "embed"
	"os"
)

//go:embed goload.js
var goloadJs []byte

func init() {
	_, err := os.Stat("goload.js")
	if os.IsNotExist(err) {
		wErr := os.WriteFile("goload.js", goloadJs, 0644)
		if wErr != nil {
			panic(wErr)
		}
	}
}
