package main

import (
	"log"

	"github.com/jmelowry/kiosk/menu"
)


func main() {
	if err := menu.Start(); err != nil {
		log.Fatal(err)
	}
}
