package main

import (
	"flag"
	"fmt"
	"gophoact/pkg/adding"
	"gophoact/pkg/editing"
	"gophoact/pkg/jobqueue"
	"gophoact/pkg/storage"
	"gophoact/pkg/viewing"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
)


func importFromDir(add adding.Service, jq jobqueue.Service, es editing.Service, srcPath string) error {
	files, err := ioutil.ReadDir(srcPath)
	if err != nil  {
		log.Printf("Error reading src dir %v", srcPath)

		return err
	}
	for _, f := range files {
		if f.IsDir(){
			importFromDir(add, jq, es, srcPath + "/" + f.Name())
			continue
		}
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
	indexPath := flag.String("index", "", "file path")
	sourcePath := flag.String("source", "", "file path")

	flag.Parse()
	var view viewing.Service

	s, err := storage.NewDbStorage(*dbPath)
	if err != nil {
		log.Printf("error opening Db, dbPath is needed")
		flag.PrintDefaults()
		os.Exit(0)
	}
	defer s.CloseDb()

	fs := storage.NewFileStorage(*filePath)
	if err != nil {
		fmt.Printf("error fiel")
		log.Panic(err)
		
	}
	is, err := storage.NewIndexStorage(*indexPath)
	defer is.CloseIndex()

	if err != nil {
		log.Printf("error opening index, indexPath is needed")
		flag.PrintDefaults()
		os.Exit(0)
	}
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
		m, err := view.ListAll(0, 30)
		if err != nil {
			fmt.Print(err)
		}
		for _, e := range m {
			log.Printf("%v", e)
		}

	case "import-dir":
		es := editing.NewService(s, fs, is)
		jq := jobqueue.NewService(es)
		add := adding.NewService(s, fs, jq)
		err := importFromDir(add, jq, es, *sourcePath)
		defer jq.CloseQueue()
		if err != nil {
			log.Panic(err)
		}

	case "detect-mimetype":
		view = viewing.NewService(s, fs)
		es := editing.NewService(s, fs, is)
		jq := jobqueue.NewService(es)
		defer jq.CloseQueue()
		items, err := view.ListAll(0, 30)
		if err != nil {
			log.Panic(err)
		}
		for _, entry := range items {
			log.Printf("%v, %v", entry.ID, entry.Mimetype)
			if entry.Mimetype.Extension == "" {
				j := jobqueue.NewMimetypeJob(entry.ID, es)
				j.Execute()
			}
		}

	case "resize-images":
		view = viewing.NewService(s, fs)
		es := editing.NewService(s, fs, is)
		jq := jobqueue.NewService(es)
		defer jq.CloseQueue()
		items, err := view.ListAll(0, 130)
		if err != nil {
			log.Panic(err)
		}
		for _, entry := range items {
			log.Printf("%v", entry.ID)
			j := jobqueue.NewImageResizeJob(entry.ID, es)
			j.Execute()
		}

	case "create-chksum":
		view = viewing.NewService(s, fs)
		es := editing.NewService(s, fs, is)
		jq := jobqueue.NewService(es)
		defer jq.CloseQueue()
		items, err := view.ListAll(0, 128)
		if err != nil {
			log.Panic(err)
		}
		for _, entry := range items {
			log.Printf("%v", entry.ID)
			j := jobqueue.NewChkSumJob(entry.ID, es)
			j.Execute()
		}



	case "search":
		//es := editing.NewService(s, fs, is)
		result, err := is.FindDocuments("Versions:*")
		
		index := is.GetIndex()
		c, err := index.DocCount()
		log.Printf("%v", index)
		log.Printf("%v", c)

		ii, kvs, err := index.Advanced()
		log.Printf("%v",ii)
		log.Printf("%v", kvs)

		ir, err := ii.Reader()
		defer ir.Close()

		fields, err  := ir.Fields()
		log.Printf("%v", fields)

		fd, err := ir.FieldDict("Versions")
		log.Printf("%v", fd)
	
		for {
			de, _ := fd.Next()
			if de == nil{
				break
			}
			log.Printf("%v", de)
		}
		log.Printf("%v", result)
		
		if err != nil {
			log.Panic(err)
		}
		
	default:
		fmt.Printf("Nothing to do")
		fmt.Printf("Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	} 

}