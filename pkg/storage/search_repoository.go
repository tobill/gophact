package storage

import (
	"log"
	"gophoact/pkg/editing"
	"github.com/blevesearch/bleve"
	// "github.com/blevesearch/bleve/analysis/analyzer/keyword"
	// "github.com/blevesearch/bleve/analysis/analyzer/simple"
	// "github.com/blevesearch/bleve/mapping"
)

//IndexStorage type for search repository
type IndexStorage struct {
	index bleve.Index
}


// func buildIndexMapping() (mapping.IndexMapping, error) {
// 	// a generic reusable mapping for keyword text
// 	keywordFieldMapping := bleve.NewTextFieldMapping()
// 	keywordFieldMapping.Analyzer = keyword.Name

// 	simpleFieldMapping := bleve.NewTextFieldMapping()
// 	simpleFieldMapping.Analyzer = simple.Name

// 	mediaMapping := bleve.NewDocumentMapping()

// 	mediaMapping.AddFieldMappingsAt("Mimetype", keywordFieldMapping)
// 	mediaMapping.AddFieldMappingsAt("Filename", keywordFieldMapping)
// 	mediaMapping.AddFieldMappingsAt("Size", keywordFieldMapping)
// 	mediaMapping.AddFieldMappingsAt("Key", keywordFieldMapping)
// 	mediaMapping.AddFieldMappingsAt("Versions", keywordFieldMapping)

// 	indexMapping := bleve.NewIndexMapping()
// 	indexMapping.AddDocumentMapping("media", mediaMapping)

// 	indexMapping.TypeField = "type"
// 	return indexMapping, nil
// }

//NewIndexStorage create new indexstorage object
func NewIndexStorage(indexPath string) (*IndexStorage, error) {
	index, err := bleve.Open(indexPath)
	if err == bleve.ErrorIndexMetaMissing {
		indexMapping := bleve.NewIndexMapping()
		index, err = bleve.New(indexPath, indexMapping)
		log.Print(err)
		if err != nil {
			log.Print(err)
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	sr := IndexStorage{}
	sr.index = index
	return &sr, nil
}

//AddDocument to search index
func (s *IndexStorage) AddDocument(m *editing.Media) (error) {
	log.Printf("%v", m)
	err := s.index.Index(m.Key, m)
	return err
}

//FindDocument find in index
func (s *IndexStorage) FindDocuments(searchWord string) ([]string, error) {
	query := bleve.NewQueryStringQuery(searchWord)
    search := bleve.NewSearchRequest(query)
	searchResults, err := s.index.Search(search)
	hits := searchResults.Hits
	var docs []string	
	for doc := range hits {
		docs = append(docs, hits[doc].ID)
	}
	return docs, err
}

func (s *IndexStorage) CloseIndex() (error) {
	return s.index.Close()
}

//GetIndex return search index
func (s *IndexStorage) GetIndex() (bleve.Index) {
	return s.index
}