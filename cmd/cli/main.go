package main

import (
	"io"
	"os"
	"mime/multipart"
	"io/ioutil"
	"gophoact/pkg/jobqueue"
	"gophoact/pkg/editing"
	"gophoact/pkg/adding"
	"gophoact/pkg/storage"
	"gophoact/pkg/viewing"
	"fmt"
	"log"
	"flag"
	)


func importFromDir(add adding.Service, jq jobqueue.Service, es editing.Service, srcPath string) error {
	files, err := ioutil.ReadDir(srcPath)
	if err != nil  {
		log.Printf("Error reading src dir")
		return err
	}
	for _, f := range files {
		log.Printf("%s", f.Name())
		mph := multipart.FileHeader{
			Filename: f.Name(), 
			Size: f.Size(),
		}
		mpf, err := os.Open(srcPath + "/" + f.Name())
		defer mpf.Close()

		var mpr io.Reader = mpf
	
		id, err := add.AddMedia(&mpr, &mph)
		log.Printf("added %d", id)
		if err != nil {
			return err
		}
		jq.GenerateMimetypeJob(id)
		
	}
	return nil
}

func main() {
	action :=  flag.String("action", "info", "action to do")
	dbPath := flag.String("dbPath", "", "db path")
	obejctID := flag.Uint64("objId", 1, "objectId")
	filePath := flag.String("storage", "", "file path")
	sourcePath := flag.String("source", "", "file path")
	
	flag.Parse()

	var view viewing.Service
	s, err := storage.NewDbStorage(*dbPath)
	fs := storage.NewFileStorage(*filePath)
	log.Print(*dbPath)
	if err != nil {
		fmt.Printf("error")
		log.Panic(err)
		
	}
	defer s.CloseDb()


	switch *action {
	case "info":
		view = viewing.NewService(s, fs)
		m, err := view.GetByID(*obejctID)
		if err != nil {
			fmt.Print(err)
		}
		fmt.Printf("%v", m)

	case "info-all":
		view = viewing.NewService(s, fs)
		m, err := view.ListAll()
		if err != nil {
			fmt.Print(err)
		}
		for _, e := range m {
			log.Printf("%v", e)
		}

	case "import-dir":
		es := editing.NewService(s, fs)
		jq := jobqueue.NewService(es)
		add := adding.NewService(s, fs, jq)
		err := importFromDir(add, jq, es, *sourcePath)
		defer jq.CloseQueue()
		if err != nil {
			log.Panic(err)
		}

	case "runjob":
		

	default:
		fmt.Printf("Nothing to do")
	} 

}