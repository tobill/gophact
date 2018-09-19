package api

import (
	"encoding/json"
	"gophoact/backend"
	"fmt"
	"log"
	"html/template"
	"net/http"
	"github.com/gorilla/mux"
)

// Environment - runtime env
var  Environment = "development"

// CreateRouter - create router and add all handlers
func CreateRouter() *mux.Router {
	r := mux.NewRouter()
	apiHandler(r)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.HandleFunc("/", indexHandler).Methods("GET")
	return r
}

// IndexPage struct
type IndexPage struct {
    JsHostURL string
}


const devHost = "http://localhost:5000"
const prodHost = ""

func indexHandler(w http.ResponseWriter, r *http.Request)  {
	ip := &IndexPage{JsHostURL: devHost}
	t, err := template.ParseFiles("./static/index.html") 	
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}
	t.Execute(w, ip)
}

func apiHandler(r *mux.Router)  {
	api := r.PathPrefix("/api", ).Subrouter()
	api.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "0.0.1")
	})
	api.HandleFunc("/file/upload", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad request")
	}).Methods("GET")
	api.HandleFunc("/file/upload", uploadFile).Methods("POST")
	api.HandleFunc("/items", getAll).Methods("GET")

}

func getAll(w http.ResponseWriter, r *http.Request) {
	items, err := backend.LoadAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", err)
	}
	b, err := json.Marshal(items)
	fmt.Fprintf(w,"%s", b)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	mf := backend.CreateEmpty()
	fd, err := mf.OpenForWrite()

	defer fd.Close()
	if err != nil {
	 	w.WriteHeader(http.StatusBadRequest)
	 	fmt.Fprintf(w, "Error uploading file")
	 	log.Fatal(err)
	 	return
	}
	
	mpf, mph, err := r.FormFile("file")
	defer mpf.Close()
	
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error uploading file")
		log.Print(err)
		return
	}

	mf.Filename = mph.Filename
	mf.Size = uint64(mph.Size)
	err = mf.Save()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", err)
	}
	//io.Copy(fd, mpf)
}