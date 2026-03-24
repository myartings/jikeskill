package tokens

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type TokenData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Store struct {
	mu   sync.RWMutex
	path string
	data *TokenData
}

func NewStore(path string) *Store {
	if path == "" {
		path = "tokens.json"
	}
	s := &Store{path: path}
	s.load()
	return s
}

func (s *Store) Get() *TokenData {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.data == nil {
		return nil
	}
	cp := *s.data
	return &cp
}

func (s *Store) Save(data *TokenData) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = data
	return s.persist()
}

func (s *Store) Delete() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = nil
	return os.Remove(s.path)
}

func (s *Store) load() {
	raw, err := os.ReadFile(s.path)
	if err != nil {
		return
	}
	var data TokenData
	if json.Unmarshal(raw, &data) == nil {
		s.data = &data
	}
}

func (s *Store) persist() error {
	dir := filepath.Dir(s.path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	raw, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, raw, 0600)
}
