package utils

import (
	"sync"

	"github.com/ReanSn0w/gokit/pkg/tool"
)

func New() *Storage {
	return &Storage{
		mx:      sync.Mutex{},
		storage: make(map[string]state),
	}
}

type Storage struct {
	mx      sync.Mutex
	storage map[string]state
}

func (s *Storage) New(url string) string {
	s.mx.Lock()
	defer s.mx.Unlock()

	id := tool.NewID()
	s.storage[id] = state{URL: url}
	return id
}

func (s *Storage) Fire(state string) (string, bool) {
	s.mx.Lock()
	defer s.mx.Unlock()

	url, ok := s.storage[state]
	if !ok {
		return "", false
	}

	delete(s.storage, state)
	return url.URL, true
}

type state struct {
	URL string
}
