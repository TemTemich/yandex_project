package main

import (
	"log"
	"orchestrator/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}

}
