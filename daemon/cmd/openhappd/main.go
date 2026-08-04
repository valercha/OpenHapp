package main

import (
	"log"
	"os"
)

const version = "0.1.0-dev"

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("openhappd %s starting", version)

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version":
			log.Println(version)
			return
		case "status":
			log.Println("status: not implemented")
			return
		}
	}

	log.Println("daemon skeleton is running")
}
