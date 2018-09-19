package api

import (
	"encoding/json"
	"gophoact/backend"
	"io"
	"path/filepath"
	"mime/multipart"
	"bytes"
	"os"
	"io/ioutil"
	"net/http"
	"fmt"
	"log"
	"net/http/httptest"
	"testing"
)

var dbPath = "../db"

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(fmt.Sprintf("%s %s %s", r.RemoteAddr, r.Method, r.URL))
		handler.ServeHTTP(w, r)
	})
}

func TestAPIInfo(t *testing.T) {
	ts := httptest.NewServer(logRequest(CreateRouter()))
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
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path)) 
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
	ts := httptest.NewServer(logRequest(CreateRouter()))
	backend.Connect(dbPath)
	defer backend.Close()
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

	var items []backend.MediaFile
	err = json.Unmarshal(resBody, &items)
	if err != nil {
		log.Printf("%s", resBody)
		t.Error(err)
	}
	if len(items) == 0 {
		t.Error("no items returned")
	}

	log.Printf("%v", items)

	
}

func TestAPIUploadFile(t *testing.T) {
	ts := httptest.NewServer(logRequest(CreateRouter()))
	backend.Connect(dbPath)
	defer backend.Close()
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