package newsy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"
)

type Session struct {
	Articles Articles `json:"articles"`
	Topics   Topics   `json:"topics"`

	loaded    bool
	timestamp time.Time
	sync.RWMutex
}

func (s *Session) Load(data []byte) error {
	if err := json.Unmarshal(data, s); err != nil {
		return err
	}
	s.loaded = true
	fmt.Println("LOADED ===", s.loaded)

	return nil
}

func (s *Session) Loaded() bool {
	s.RLock()
	defer s.RUnlock()
	return s.loaded
}

func (s *Session) Save(storage io.Writer) error {
	s.Lock()
	data, err := json.Marshal(s)
	s.Unlock()
	if err != nil {
		return err
	}

	if _, err := storage.Write(data); err != nil {
		return err
	}

	s.Lock()
	s.timestamp = time.Now()
	s.Unlock()

	return nil
}

func (s *Session) StartAutoSave(interval time.Duration) {
	time.Sleep(interval)
	// FIX: Make this an actual write to the file
	bb := &bytes.Buffer{}

	s.Save(bb)
}
