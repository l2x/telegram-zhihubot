package main

import (
	"flag"
	"log"
)

func main() {
	var file string
	flag.StringVar(&file, "c", "", "config file")
	flag.Parse()

	if err := initConfig(file); err != nil {
		log.Fatal(err)
	}

	if err := botRun(); err != nil {
		log.Fatal(err)
	}
}
