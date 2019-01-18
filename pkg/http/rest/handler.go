package rest

import (
	"time"
	"strconv"
	"io"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"html/template"
	"gophoact/pkg/adding"
	"gophoact/pkg/viewing"
	"runtime"
)

// CreateRouter return http router
func CreateRouter(a adding.Service, v viewing.Service) *mux.Router {
	r := mux.NewRouter()
	apiHandler(r, a, v)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.HandleFunc("/", indexHandler).Methods("GET")
	return r
}

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	log.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	log.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	log.Printf("\tSys = %v MiB", bToMb(m.Sys))
	log.Printf("\tNumGC = %v\n", m.NumGC)
}


func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}


// LogMiddleware log every request by middleware
func LogMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
		log.Println(fmt.Sprintf("%s %s %s", r.RemoteAddr, r.Method, r.URL))
		printMemUsage()
	})
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

func apiHandler(r *mux.Router, a adding.Service, v viewing.Service)  {
	api := r.PathPrefix("/api", ).Subrouter()
	api.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "0.0.1")
	})
	api.HandleFunc("/file/upload", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad request")
	}).Methods("GET")
	api.HandleFunc("/file/upload", uploadFile(a)).Methods("POST")
	api.HandleFunc("/items", getAll(v)).Methods("GET")
	api.HandleFunc("/items/{id}", getItem(v)).Methods("GET")
	api.HandleFunc("/items/{id}/file", getFile(v)).Methods("GET")
	api.HandleFunc("/items/{id}/file/{size}", getFileBySize(v)).Methods("GET")
}

func getAll(v viewing.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := v.ListAll()	
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", err)
		}
		b, err := json.Marshal(items)
		fmt.Fprintf(w,"%s", b)
	}
}

func uploadFile(a adding.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		mpf, mph, err := r.FormFile("file")
		defer mpf.Close()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error uploading file")
			log.Fatal(err)
			return
		}
		var mpr io.Reader = mpf
		_, err = a.AddMedia(&mpr, mph)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error uploading file")
			log.Fatal(err)
			return
		}
	}
}

func getItem(v viewing.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	sid := params["id"]
	id, err := strconv.ParseUint(sid, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", err)
	}
	var item *viewing.Media
	item, err = v.GetByID(id)
	if err != nil {
	 	w.WriteHeader(http.StatusInternalServerError)
	 	fmt.Fprintf(w, "%v", err)
	}
	b, err := json.Marshal(item)
	if err != nil {
	 	w.WriteHeader(http.StatusInternalServerError)
	 	fmt.Fprintf(w, "%v", err)
	}
	fmt.Fprintf(w,"%s", b)
	}
} 

func getFile(v viewing.Service) func(w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		sid := params["id"]
		id, err := strconv.ParseUint(sid, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", err)
	   	}	
		   
		fop, err := v.GetFileByID(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", err)
		}
		defer fop.Close()
		
		http.ServeContent(w,r, "image", time.Now(), fop)
	}
} 

func getFileBySize(v viewing.Service) func(w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		sid := params["id"]
		id, err := strconv.ParseUint(sid, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", err)
	   	}	
		   
		fop, err := v.GetFileByID(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", err)
		}
		defer fop.Close()
		
		http.ServeContent(w,r, "image", time.Now(), fop)
	}
} 


