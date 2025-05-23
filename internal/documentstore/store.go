package documentstore

import (
	"encoding/json"
	"fmt"
	"les12/pkg/err"
	"log/slog"
	"os"
	"strings"
)

type Store struct {
	collections map[string]*Collection
}

func NewStore() *Store {
	return &Store{
		collections: make(map[string]*Collection),
	}
}

type DTOStore struct {
	Collections map[string]DTOCollection `json:"collections"`
}

func (s *Store) ToDto() DTOStore {
	dtoCollections := make(map[string]DTOCollection, len(s.collections))
	for name, coll := range s.collections {
		dtoCollections[name] = coll.ToDto()
	}
	return DTOStore{
		Collections: dtoCollections,
	}
}

func (s *Store) CreateCollection(name, id string) (*Collection, error) {
	// Створюємо нову колекцію і повертаємо `true` якщо колекція була створена
	// Якщо ж колекція вже створеня то повертаємо `false` та nil
	if _, exists := s.collections[name]; exists {
		return nil, err.ErrCollectionAlreadyExists
	}
	coll := &Collection{
		config: CollectionConfig{
			PrimaryKey: id,
		}}
	s.collections[name] = coll
	// slog.Info("collection added")

	return coll, nil
}

func (s *Store) GetCollection(name string) (*Collection, error) {
	if colect, ok := s.collections[name]; ok {
		return colect, nil
	}
	slog.Error("collection not found")
	return nil, err.ErrCollectionNotFound
}

func (s *Store) DeleteCollection(name string) bool {
	if _, ok := s.collections[name]; ok {
		delete(s.collections, name)
		// slog.Info("collection delete - %s")
		return true
	}
	return false
}

func NewStoreFromDump(dump []byte) (*Store, error) {
	// Функція повинна створити та проініціалізувати новий `Store`
	// зі всіма колекціями та даними з вхідного дампу.
	var s Store
	if err := json.Unmarshal(dump, &s); err != nil {
		return nil, err
	}
	if len(s.collections) == 0 {
		// slog.Info("collection not added")
		return nil, err.ErrNotFound
	}
	return &s, nil
}

func (s *Store) Dump() ([]byte, error) {
	// Методи повинен віддати дамп нашого стору в який включені дані про колекції та документ
	sToJson, err := json.MarshalIndent(s, " ", "")
	if err != nil {
		return nil, err
	}
	return sToJson, nil
}

func NewStoreFromFile(filename string) (*Store, error) {
	// Робить те ж саме що і функція `NewStoreFromDump`, але сам дамп має діставатись з файлу

	fileString := strings.Builder{}
	fileString.WriteString(filename + ".json")

	dump, err := os.ReadFile(fileString.String())
	if err != nil {
		slog.Error("file not read")
		return nil, err
	}
	// slog.Info("file read successfully " + fileString.String())
	s := NewStore()
	var dto DTOStore
	if err := json.Unmarshal(dump, &dto); err != nil {

		return nil, err
	}
	for name, dtoColl := range dto.Collections {
		coll := &Collection{
			documents: dtoColl.Documents,
			config:    dtoColl.Config,
		}
		s.collections[name] = coll
	}
	if len(s.collections) == 0 {
		slog.Error("no collections found in store from file")
		return nil, fmt.Errorf("no collections in store")
	}
	return s, nil
}

func (s *Store) DumpToFile(filename string) error {
	// Робить те ж саме що і метод  `Dump`, але записує у файл замість того щоб повертати сам дамп
	sDump, err := s.Dump()
	if err != nil {

		fmt.Println(err)
	}

	fileString := strings.Builder{}
	fileString.WriteString(filename + ".json")

	return os.WriteFile(fileString.String(), sDump, 0644)
}
