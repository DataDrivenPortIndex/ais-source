package main

import (
	"fmt"
	"log"

	"github.com/DataDrivenPortIndex/ais-source/internal/source"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	aisSource, err := source.NewAisStreamSource()
	if err != nil {
		log.Fatal(err)
	}
	defer aisSource.Close()

	for aisMessage := range aisSource.Read() {
		fmt.Println(aisMessage)
	}
}
