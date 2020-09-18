package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/thecubic/godexrcvr"
)

var (
	debug = flag.Bool("debug", false, "enable debugging messages")
)

func usageQuit() {
	fmt.Println("usage: gdr-info <command> [args]")
	os.Exit(1)
}

func main() {
	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	command := flag.Arg(0)
	if command == "" {
		usageQuit()
	}

	if command == "hello" {
		fmt.Printf("syncbyte: %v\n", godexrcvr.SyncByte)
	}
}
