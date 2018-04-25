package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/nicolaferraro/datamesh/service"
	"context"
	"github.com/nicolaferraro/datamesh/protobuf"
	"github.com/golang/glog"
)

func main() {
	dir := flag.String("dir", "./data", "Data directory. Default: \".data\"")
	port := flag.Int("port", 6543, "Server port. Default: 6543")
	host := flag.String("host", "localhost", "Host of the server (for client commands). Default: \"localhost\"")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Printf("Error. Syntax: datamesh [server|client] <options>\n")
		os.Exit(1)
	}
	contexName := flag.Arg(0)
	if contexName != "server" && contexName != "client" {
		fmt.Printf("Unknown contexName: '%s'. Allowed [server|client]\n", contexName)
		os.Exit(1)
	}

	glog.V(1).Infof("Configured dir: '%s'\n", *dir)
	glog.V(1).Infof("Configured port: '%d'\n", *port)
	glog.V(1).Infof("Configured host: '%s'\n", *host)

	if contexName == "server" {
		ctx, cancel := context.WithCancel(context.Background())
		msh, err := service.NewMesh(ctx, *dir, *port)
		if err != nil {
			glog.Fatal("Cannot initialize data mesh: ", err)
		}

		glog.Infof("Data Mesh started on port %d\n", *port)

		if err := msh.Start(); err != nil {
			glog.Fatal("Data Mesh error: ", err)
		}
		cancel()
	} else if contexName == "client" {
		if flag.NArg() < 2 {
			fmt.Printf("Error. Syntax: datamesh client <action> <options>\n")
			os.Exit(1)
		}
		action := flag.Arg(1)
		if action == "push" {
			if flag.NArg() != 4 {
				fmt.Printf("Error. Syntax: datamesh client push <eventName> <eventPayload>\n")
				os.Exit(1)
			}
			event := protobuf.Event{
				Name: flag.Arg(2),
				Payload: []byte(flag.Arg(3)),
			}
			client, err := service.NewDataMeshClientConnection(*host, *port)
			if err != nil {
				fmt.Printf("Error. Cannot create client connection: %v\n", err)
				os.Exit(1)
			}
			ctx := context.Background()
			defer ctx.Done()
			_, err = client.Push(ctx, &event)
			if err != nil {
				fmt.Printf("Error. Cannot push event: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Event pushed successfully.\n")

		} else {
			fmt.Printf("Error. Unknown action %s\n", action)
			os.Exit(1)
		}

	} else {
		fmt.Printf("Error. Unknown contexName %s\n", contexName)
		os.Exit(1)
	}

}
