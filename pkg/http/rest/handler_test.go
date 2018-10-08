package rest

import (
	"io"
	"path/filepath"
	"mime/multipart"
	"bytes"
	"os"
	"encoding/json"
	"io/ioutil"
	"gophoact/pkg/storage"
	"gophoact/pkg/viewing"
	"gophoact/pkg/adding"
	"log"
	"fmt"
	"testing"
	"net/http/httptest"
	"net/http"
)

const testDbPath = "../../../testdb"
const testFilepath = "../../../testdata"
const testFile = "../../../sampledata/TESTIMG.JPG"

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(fmt.Sprintf("%s %s %s", r.RemoteAddr, r.Method, r.URL))
		handler.ServeHTTP(w, r)
	})
}

func TestGetItem(t *testing.T) {
	var adder adding.Service 
	var view viewing.Service
	s, err := storage.NewDbStorage(testDbPath)
	defer s.CloseDb()
	if err != nil {
		t.Fatal(err)
	}

	fs := storage.NewFileStorage(testFilepath)
	adder = adding.NewService(s, fs)
	view = viewing.NewService(s, fs)
	r := CreateRouter(adder, view);
	id := 0
	ts := httptest.NewServer(logRequest(r))
	defer ts.Close()

	url := fmt.Sprintf("%s/api/items/%d", ts.URL, id)
	res, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	resBody, err  := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	var item viewing.Media
	err = json.Unmarshal(resBody, &item)
	if err != nil {
		log.Printf("%s", resBody)
		t.Error(err)
	}
}

func TestAPIInfo(t *testing.T) {
	var adder adding.Service 
	var view viewing.Service
	s, err := storage.NewDbStorage(testDbPath)
	defer s.CloseDb()
	if err != nil {
		t.Fatal(err)
	}
	fs := storage.NewFileStorage(testFilepath)
	adder = adding.NewService(s, fs)
	view = viewing.NewService(s, fs)
	r := CreateRouter(adder, view);
	ts := httptest.NewServer(logRequest(r))
	defer ts.Close()
	url := fmt.Sprintf("%v/api/info", ts.URL)
	res, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	resBody, err  := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("%s", resBody)
}

func generateUploadrequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(testFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(testFile)) 
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	err = writer.Close()

	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, err
}


func TestAPIGetFirstItems(t *testing.T) {
	var adder adding.Service 
	var view viewing.Service
	s, err := storage.NewDbStorage(testDbPath)
	defer s.CloseDb()

	if err != nil {
		t.Fatal(err)
	}
	fs := storage.NewFileStorage(testFilepath)
	adder = adding.NewService(s, fs)
	view = viewing.NewService(s, fs)
	r := CreateRouter(adder, view);
	ts := httptest.NewServer(logRequest(r))
	defer ts.Close()

	url := fmt.Sprintf("%v/api/items", ts.URL)
	res, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	resBody, err  := ioutil.ReadAll(res.Body)
	if err != nil {
			t.Fatal(err)
	}
	var items []viewing.Media
	err = json.Unmarshal(resBody, &items)
	if err != nil {
		t.Error(err)
	}
	if len(items) == 0 {
		t.Error("no items returned")
	}

	log.Printf("%v", items)
}

func TestAPIUploadFile(t *testing.T) {
	var adder adding.Service 
	var view viewing.Service
	s, err := storage.NewDbStorage(testDbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer s.CloseDb()
	fs := storage.NewFileStorage(testFilepath)
	adder = adding.NewService(s, fs)
	view = viewing.NewService(s, fs)
	r := CreateRouter(adder, view);
	ts := httptest.NewServer(logRequest(r))

	defer ts.Close()
	testfilePath, err := filepath.Abs("../testdata/TESTIMG.JPG")
	if err != nil {
		t.Fatal(err)
	}

	extraparams := map[string]string{}
	url := fmt.Sprintf("%v/api/file/upload", ts.URL)

	req, err := generateUploadrequest(url, extraparams, "file", testfilePath)
	if err != nil {
		t.Fatal(err)
	}
	
    client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode == 404 {
		t.Fatal("Error endpoint not dfound")
	}

	resBody, err  := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("%s", resBody)
}