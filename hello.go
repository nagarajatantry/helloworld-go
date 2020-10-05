package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

var ra *rand.Rand
var latencies = []time.Duration{80,50,70,75,70,50,50,50,80,60,60,80,90,100,120,150,190,350,150,250,250,250,240,240,
	80,50,70,75,70,50,50,50,80,60,60,80,90,100,120,150,190,350,150,250,250,250,240,240,
	80,50,70,75,70,50,50,50,80,60,60,80,90,100,120,150,190,350,150,250,250,250,240,240,
	80,50,70,75,70,50,50,50,80,60,60,80,90,100,120,150,190,350,150,250,250,250,240,240,
	80,50,70,75,70,50,50,50,80,60,60,80,90,100,120,150,190,350,150,250,250,250,240,240,2000,
	80,50,70,75,70,50,50,50,80,60,60,80,90,100,120,150,190,350,150,250,250,250,240,240,
	80,50,70,75,70,50,50,50,80,60,60,80,90,100,120,150,190,350,150,250,250,250,240,240,
	80,50,70,75,70,50,50,50,80,60,60,80,90,100,120,150,190,350,150,250,250,250,240,240,
	80,50,70,75,70,50,50,50,80,60,60,80,90,100,120,150,190,350,150,250,250,250,240,240,3000,
	80,50,70,75,70,50,50,50,80,60,60,80,90,100,120,150,190,350,150,250,250,250,240,240,
	80,50,70,75,70,50,50,50,80,60,60,80,90,100,120,150,190,350,150,250,250,250,240,240,
	80,50,70,75,70,50,50,50,80,60,60,80,90,100,120,150,190,350,150,250,250,250,240,240,10000,
	80,50,70,75,70,50,50,50,80,60,60,80,90,100,120,150,190,350,150,250,250,250,240,240,
	80,50,70,75,70,50,50,50,80,60,60,80,90,100,120,150,190,350,150,250,250,250,240,240,1000}
var size int
func init(){
	ra =rand.New(rand.NewSource(time.Now().UnixNano()))
	size = len(latencies)
}

func trace(w http.ResponseWriter, r *http.Request) {
	timeToSleep := latencies[ra.Intn(size)]
	time.Sleep(timeToSleep * time.Millisecond)
	random := r.URL.Query().Get("random")
	if random == "" && ra.Intn(120000) != 500 {
		w.Write([]byte("{}"))
		return
	}
	log.Printf("Protocol:%s \n",r.Proto)
	log.Printf("header:%s \n", r.Header)
	w.Write([]byte("{}"))
}
func handler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	name := query.Get("name")
	if name == "" {
		name = "world"
	}
	w.Write([]byte(fmt.Sprintf("hello, %s\n", name)))
}

func main() {
	// Create Server and Route Handlers
	r := mux.NewRouter()

	r.HandleFunc("/greet", handler)
	r.HandleFunc("/v1/trace", trace)
	r.HandleFunc("/", handler)
	srv := &http.Server{
		Handler:      r,
		Addr:         ":9098",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,

	}

	// Start Server
	go func() {
		log.Println("Starting Server")
		/*
			openssl req -new -newkey rsa:2048 -nodes -sha256 -days 730 -x509 \
			-subj '/CN=localhost/O=Intuit Inc./C=US' \
			-keyout key.pem \
			-out cert.pem
		*/
		if err := srv.ListenAndServeTLS("cert.pem", "key.pem"); err != nil {
			log.Fatal(err)
		}

	}()

	// Graceful Shutdown
	waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}
