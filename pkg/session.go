package newsy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type session struct {
	Backup   io.Writer `json:"-"`
	Articles Articles  `json:"articles"`
	Topics   Topics    `json:"topics"`

	loaded    bool
	timestamp time.Time

	sync.RWMutex
}

func Restore(r io.Reader) (*session, error) {
	dec := json.NewDecoder(r)
	sess := &session{}
	if err := dec.Decode(&sess); err != nil {
		if err != io.EOF {
			return nil, err
		}
	}
	sess.loaded = true

	return sess, nil
}

func (s *session) backup() error {
	if s == nil {
		return fmt.Errorf("session is nil")
	}
	fmt.Println("STARTING BACKUP")

	// check for backup location
	if s.Backup == nil {

		// REFACTOR: create/init on if it doesn't exist
		// (move to using something like Datastore.set?)
		bak, err := os.Create("newsy_db.json")
		if err != nil {
			return err
		}

		// Set established backup location
		s.Backup = bak
	}

	s.Lock()
	// Backup the things
	err := json.NewEncoder(s.Backup).Encode(s)
	s.Unlock()
	if err != nil {
		return fmt.Errorf("session backup failure: %v", err)
	}

	if closer, ok := s.Backup.(io.Closer); ok {
		if err := closer.Close(); err != nil {
			return err
		}
	}

	s.Lock()
	s.timestamp = time.Now()
	s.Unlock()

	fmt.Println("session backup complete")
	return nil
}

func (s *session) startAutoBackup(interval time.Duration) {
	time.Sleep(interval)
	if s.Backup == nil {
		s.Backup = &bytes.Buffer{}
	}

	s.backup()
}
