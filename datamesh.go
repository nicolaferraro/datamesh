package main

import (
	"flag"
	mesh "github.com/nicolaferraro/datamesh/mesh"
	"log"
	"fmt"
	"os"
)

func main() {
	dir := flag.String("dir", "./data", "Data directory. Default: \".data\"")
	port := flag.Int("port", 6543, "Server port. Default: 6543")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Printf("Error. Syntax: datamesh [server|client] <options>\n")
		os.Exit(1)
	}
	context := flag.Arg(0)
	if context != "server" && context != "client" {
		fmt.Printf("Unknown context: '%s'. Allowed [server|client]\n", context)
		os.Exit(1)
	}

	log.Printf("configured dir: '%s'", *dir)
	log.Printf("configured port: '%d'", *port)

	if context == "server" {
		msh, err := mesh.NewMesh(*dir, *port)
		if err != nil {
			log.Fatal("cannot initialize data mesh: ", err)
		}

		log.Printf("data mesh started on port %d\n", *port)

		if err := msh.Start(); err != nil {
			log.Fatal("data mesh error: ", err)
		}
	} else if context == "client" {

	}

}
