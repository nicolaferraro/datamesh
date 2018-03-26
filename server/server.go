package server

import (
	"net/http"
	"strconv"
	"io/ioutil"
	"github.com/nicolaferraro/datamesh/utils"
	"log"
)

type Server struct {
	port		int
	observer	utils.MessageObserver
}

func NewServer(port int, observer utils.MessageObserver) *Server {
	return &Server{
		port: port,
		observer: observer,
	}
}

func (srv *Server) Start() error {
	http.HandleFunc("/", srv.handler)
	return http.ListenAndServe("0.0.0.0:" + strconv.Itoa(srv.port), nil)
}

func (srv *Server) handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(400)
	} else {
		bytes, err := ioutil.ReadAll(r.Body)
		if err == nil {
			err = srv.observer.Accept(bytes)
			if err != nil {
				log.Printf("error while saving HTTP message: %s\n", err)
				w.WriteHeader(500)
			}
		} else {
			log.Println("error while reading HTTP message")
			w.WriteHeader(500)
		}
	}
}