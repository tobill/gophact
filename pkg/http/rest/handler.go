package rest

import (
	"io"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"html/template"
	"gophoact/pkg/adding"
	"gophoact/pkg/viewing"
)

// CreateRouter return http router
func CreateRouter(a adding.Service, v viewing.Service) *mux.Router {
	r := mux.NewRouter()
	apiHandler(r, a, v)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.HandleFunc("/", indexHandler).Methods("GET")
	return r
}

// LogMiddleware log every request by middleware
func LogMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(fmt.Sprintf("%s %s %s", r.RemoteAddr, r.Method, r.URL))
		handler.ServeHTTP(w, r)
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
	api.HandleFunc("/items/{id}", getItem).Methods("GET")
	api.HandleFunc("/items/{id}/file", getFile).Methods("GET")
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
		err = a.AddMedia(&mpr, mph)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error uploading file")
			log.Fatal(err)
			return
		}
	}
}

func getItem(w http.ResponseWriter, r *http.Request) {
	// v := mux.Vars(r)
	// mf := backend.MediaFile{}
	// sid := v["ID"]
	// id, err := strconv.ParseUint(sid, 10, 64)
	// err = mf.Load(id)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	fmt.Fprintf(w, "%v", err)
	// }
	// b, err := json.Marshal(mf)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	fmt.Fprintf(w, "%v", err)
	// }
	// fmt.Fprintf(w,"%s", b)
	
} 

func getFile(w http.ResponseWriter, r *http.Request) {
	
} 
