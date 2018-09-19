package main


import (
	"gophoact/api"
	"fmt"
	"log"
	"net/http"
	"os"
)

// DbFilePath - the path of dbfiles
const DbFilePath = "../db"


func main() {
	api.Environment = os.Getenv("APPENV")
	if api.Environment == "" {
		api.Environment = "development"
	}
	log.Println(fmt.Sprintf("running in %s ", api.Environment))
	r := api.CreateRouter();
	err := http.ListenAndServe(":8080", logRequest(r))
	if err != nil {
		fmt.Printf("error")
	}
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(fmt.Sprintf("%s %s %s", r.RemoteAddr, r.Method, r.URL))
		handler.ServeHTTP(w, r)
	})
}
