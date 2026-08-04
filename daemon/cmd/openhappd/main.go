package main

import (
	"log"
	"os"

	"github.com/valercha/OpenHapp/daemon/internal/state"
	"github.com/valercha/OpenHapp/daemon/internal/version"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("openhappd %s starting", version.Version)

	st := state.New(version.Version)

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version":
			log.Println(version.Version)
			return
		case "status":
			log.Printf("status: %+v", st.Snapshot())
			return
		}
	}

	st.Start()
	log.Printf("daemon skeleton is running: %+v", st.Snapshot())
}
