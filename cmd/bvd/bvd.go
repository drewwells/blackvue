package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/drewwells/blackvue"
)

func main() {

	msg := "please select a command: status, sync"
	if len(os.Args) == 1 {
		log.Fatal(msg)
	}

	switch os.Args[1] {
	case "status":
		status()
	case "sync":
		sync()
	default:
		log.Fatal(msg)
	}
}

func status() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: status ip_of_dashcam save_directory")
	}
	bv := blackvue.New(os.Args[2])

	abs, err := filepath.Abs(os.Args[3])
	if err != nil {
		log.Fatal(err)
	}

	sum, err := bv.Status(abs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", sum)
}

func sync() {

	if len(os.Args) < 4 {
		log.Fatal("Usage: sync ip_of_dashcam save_directory")
	}

	bv := blackvue.New(os.Args[2])

	abs, err := filepath.Abs(os.Args[3])
	if err != nil {
		log.Fatal(err)
	}

	if err := bv.Sync(abs); err != nil {
		log.Fatal(err)
	}
}
