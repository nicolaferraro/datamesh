package main

import (
	"flag"
	mesh "github.com/nicolaferraro/datamesh/mesh"
	"log"
)

func main() {
	dir := flag.String("dir", "./data", "Data directory. Default: \".data\"")
	port := flag.Int("port", 6543, "Server port. Default: 6543")
	flag.Parse()

	log.Printf("configured dir: '%s'", *dir)
	log.Printf("configured port: '%d'", *port)

	msh, err := mesh.NewMesh(*dir, *port)
	if err != nil {
		log.Fatal("cannot initialize data mesh: ", err)
	}

	log.Printf("data mesh started on port %d", *port)

	if err := msh.Start(); err != nil {
		log.Fatal("data mesh error: ", err)
	}

}
