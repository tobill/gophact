package main


import (
	"gophoact/pkg/http/rest"
	"gophoact/pkg/adding"
	"gophoact/pkg/viewing"
	"gophoact/pkg/storage"
	"fmt"
	"log"
	"net/http"
	"os"
)

var  Environment = "development"

var filePath = "./data/"
var dbPath = "./db/"


func main() {
	Environment = os.Getenv("APPENV")
	if Environment == "" {
		Environment = "development"
	}
	var adder adding.Service 
	var view viewing.Service
	s, err := storage.NewDbStorage(dbPath)
	fs := storage.NewFileStorage(filePath)
	defer s.CloseDb()
	log.Println(fmt.Sprintf("running in %s ", Environment))
	adder = adding.NewService(s, fs)
	view = viewing.NewService(s, fs)
	r := rest.CreateRouter(adder, view);
	err = http.ListenAndServe(":8080", rest.LogMiddleware(r))
	if err != nil {
		fmt.Printf("error")
	}
}

