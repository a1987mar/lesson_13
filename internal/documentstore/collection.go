package documentstore

import (
	"errors"
	"les12/pkg/err"
	"log/slog"
	"sort"
	"sync"
)

type Collection struct {
	mu        sync.RWMutex
	documents map[string]Document
	config    CollectionConfig
	indexes   map[string]*Index
}

type Index struct {
	Field      string
	Data       map[string]map[string]struct{} // fieldValue -> set of document IDs
	SortedKeys []string                       // cache відсортованих ключів для швидкого запиту
}

type DTOCollection struct {
	Documents map[string]Document `json:"documents,omitempty"`
	Config    CollectionConfig    `json:"config"`
}

type QueryParams struct {
	Desc     bool    // Визначає в якому порядку повертати дані
	MinValue *string // Визначає мінімальне значення поля для фільтрації
	MaxValue *string // Визначає максимальне значення поля для фільтрації
}

func (s *Collection) Query(fieldName string, params QueryParams) ([]Document, error) {
	index, ok := s.indexes[fieldName]
	if !ok {
		return nil, errors.New("index does not exist")
	}
	keys := index.SortedKeys
	if params.Desc {
		reversed := make([]string, len(keys))
		copy(reversed, keys)
		sort.Sort(sort.Reverse(sort.StringSlice(reversed)))
		keys = reversed
	}

	var result []Document

	for _, key := range keys {
		if params.MinValue != nil && key < *params.MinValue {
			continue
		}
		if params.MaxValue != nil && key > *params.MaxValue {
			continue
		}

		for id := range index.Data[key] {
			if doc, ok := s.documents[id]; ok {
				result = append(result, doc)
			}
		}
	}
	return result, nil
}

func (s *Collection) CreateIndex(fieldName string) error {

	if _, exists := s.indexes[fieldName]; exists {
		return errors.New("index already exists")
	}
	index := &Index{
		Field:      fieldName,
		Data:       make(map[string]map[string]struct{}),
		SortedKeys: []string{},
	}

	for id, doc := range s.documents {
		field, ok := doc.Fields[fieldName]
		if !ok || field.Type != DocumentFieldTypeString {
			continue
		}
		val := field.Value.(string)
		if _, exists := index.Data[val]; !exists {
			index.Data[val] = map[string]struct{}{}
			index.SortedKeys = append(index.SortedKeys, val)
		}
		index.Data[val][id] = struct{}{}
	}

	sort.Strings((index.SortedKeys))

	if s.indexes == nil {
		s.indexes = map[string]*Index{}
	}
	s.indexes[fieldName] = index
	return nil
}

func (s *Collection) DeleteIndex(fieldName string) error {
	if _, exists := s.indexes[fieldName]; !exists {
		return errors.New("index does not exist")
	}
	delete(s.indexes, fieldName)
	return nil
}

func (s *Collection) ToDto() DTOCollection {
	return DTOCollection{
		Documents: s.documents,
		Config:    s.config,
	}
}

type CollectionConfig struct {
	PrimaryKey string `json:"cgg"`
}

func (s *Collection) Put(doc Document) error {
	// Потрібно перевірити що документ містить поле `{cfg.PrimaryKey}` типу `string`

	keyFilds, ok := doc.Fields[s.config.PrimaryKey]
	if !ok {
		slog.Error("error: Document must contain a key field")
		return err.ErrUnsupportedDocumentField
	}

	if keyFilds.Type != DocumentFieldTypeString {
		slog.Error("error: Key field must be of type string")
		return err.ErrUnsupportedDocumentField
	}
	keyValue, ok := keyFilds.Value.(string)
	if !ok {
		slog.Error("Error: Key field value is not a string")
		return err.ErrUnsupportedDocumentField
	}
	if s.documents == nil {
		s.documents = map[string]Document{}
	}
	s.mu.Lock()
	s.documents[keyValue] = doc
	slog.Info("document added")
	s.mu.Unlock()
	return nil

}

func (s *Collection) Get(key string) (*Document, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if doc, exists := s.documents[key]; exists {
		return &doc, nil
	}
	slog.Info("document not found")
	return nil, err.ErrDocumentNotFound
}

func (s *Collection) Delete(key string) bool {
	if _, exists := s.documents[key]; exists {
		s.mu.Lock()
		delete(s.documents, key)
		s.mu.Unlock()
		slog.Info("document delete")
		return true
	}
	return false
}

func (s *Collection) List() []Document {
	sList := make([]Document, 0, len(s.documents))
	s.mu.RLock()
	for _, v := range s.documents {
		sList = append(sList, v)
	}
	s.mu.RUnlock()
	return sList
}
