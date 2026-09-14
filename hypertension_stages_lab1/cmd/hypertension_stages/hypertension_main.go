package main

import (
	"log"

	"hypertensionstages/internal/api"
)

func main() {
	log.Println("Hypertension stages application started")
	api.StartHypertensionServer()
	log.Println("Hypertension stages application stopped")
}
