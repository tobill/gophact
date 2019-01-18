package main


import (
	"gophoact/pkg/jobqueue"
	"gophoact/pkg/http/rest"
	"gophoact/pkg/adding"
	"gophoact/pkg/viewing"
	"gophoact/pkg/storage"
	"gophoact/pkg/editing"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"context"
	"time"
)

var  Environment = "development"

var filePath = "./data/"
var dbPath = "./db/data.db"
var indexPath = "./index/"

var (
	port            = 8080 
	shutdownTimeout = 5
)


func main() {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt)

	Environment = os.Getenv("APPENV")
	if Environment == "" {
		Environment = "development"
	}
	var adder adding.Service 
	var view viewing.Service
	s, err := storage.NewDbStorage(dbPath)
	defer s.CloseDb()

	if err != nil {
		fmt.Printf("error")
		log.Panic(err)
		
	}
	fs := storage.NewFileStorage(filePath)
	log.Println(fmt.Sprintf("running in %s ", Environment))
	is, err := storage.NewIndexStorage(indexPath) 
	defer is.CloseIndex()

	if err != nil {		
		fmt.Printf("error")
		log.Panic(err)
	}
	e := editing.NewService(s, fs, is)
	jq := jobqueue.NewService(e)
	adder = adding.NewService(s, fs, jq)
	view = viewing.NewService(s, fs)


	r := rest.CreateRouter(adder, view);
	if err != nil {
		fmt.Printf("error")
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: rest.LogMiddleware(r),
	}
	go func() {
		log.Printf("Listening on http://0.0.0.0:%d\n", 8080)

		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-stop

	log.Println("\nShutting down the server...")
	err = s.CloseDb()
	jq.CloseQueue()
	if err != nil  { log.Print(err) }

	ctx, canceld := context.WithTimeout(context.Background(), 5*time.Second)
	defer canceld()
	srv.Shutdown(ctx)
    log.Println("Server gracefully stopped")
}

