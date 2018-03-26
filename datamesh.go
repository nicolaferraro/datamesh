package main

import (
	"flag"
	mesh "github.com/nicolaferraro/datamesh/mesh"
	"log"
)

func main() {
	dir := flag.String("dir", "./data", "Data directory for the data mesh")
	flag.Parse()

	_, err := mesh.NewMesh(*dir)
	if err != nil {
		log.Fatal("Cannot start data mesh: ", err)
	}

	log.Println("Data mesh started...")

}
