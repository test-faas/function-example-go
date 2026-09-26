package main

// Deliberate fault fixture: becomes Ready, then stops listening while the process stays
// alive. No crash and no restart, so the pod remains Running and Ready and the activator
// keeps routing to it — every request after the cutoff fails to dial.

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("still listening\n"))
	})

	ln, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		time.Sleep(40 * time.Second)
		log.Println("closing the listener; process stays alive")
		_ = ln.Close()
	}()

	log.Printf("listening on 0.0.0.0:%s", port)
	_ = http.Serve(ln, mux)
	log.Println("no longer listening, but still running")
	select {} // keep the process up so nothing restarts
}
