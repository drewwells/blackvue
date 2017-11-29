package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/drewwells/blackvue"
)

func main() {
	fmt.Println(os.Args)
	if len(os.Args) < 3 {
		log.Fatal("Usage: ip_of_dashcam save_directory")
	}

	bv := blackvue.New(os.Args[1])

	abs, err := filepath.Abs(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	if err := bv.Sync(abs); err != nil {
		log.Fatal(err)
	}
}
