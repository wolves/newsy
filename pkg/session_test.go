package newsy

import (
	"bytes"
	"testing"
	"time"
)

func newTestSession(t testing.TB, bak *bytes.Buffer) *session {
	t.Helper()

	return &session{
		Backup:   bak,
		Articles: Articles{},
		Topics:   Topics{},
	}
}

func Test_Session(t *testing.T) {
	t.Parallel()
	mockTopics := Topics{"testing", "go"}
	mockArticles := Articles{{ID: 135}}
	mockStoreJson := []byte(`{"topics":["testing","go"],"articles":[{"id": 135}]}`)

	t.Run("loads state from backup file", func(t *testing.T) {
		bb := bytes.NewBuffer(mockStoreJson)
		sess, err := Restore(bb)
		assertNoError(t, err)

		if !sess.loaded {
			t.Fatalf("expected state to be loaded, got %v", sess.loaded)
		}

		gotTop := len(sess.Topics)
		expTop := len(mockTopics)
		assertLenMatch(t, gotTop, expTop)

		gotArt := len(sess.Articles)
		expArt := len(mockArticles)
		assertLenMatch(t, gotArt, expArt)
	})

	t.Run("saves state to backup file", func(t *testing.T) {
		bb := &bytes.Buffer{}
		s := newTestSession(t, bb)
		s.Topics = mockTopics
		s.Articles = mockArticles

		err := s.backup()
		if err != nil {
			t.Fatalf("unexpected error while saving state: ERR %v", err)
		}

		act := bb.Bytes()
		exp := mockStoreJson

		if bytes.Contains(act, exp) {
			t.Fatalf("got %v, expected %v", exp, act)
		}
	})

	t.Run("session autobackup", func(t *testing.T) {
		orig := time.Now()
		bb := &bytes.Buffer{}
		s := newTestSession(t, bb)
		s.loaded = true
		s.timestamp = orig
		i := time.Millisecond * 50

		s.startAutoBackup(i)

		// FIX: Change this to have autosave use a channel and capture that?
		time.Sleep(2 * i)

		s.RLock()
		expTime := s.timestamp.After(orig)
		s.RUnlock()

		if !expTime {
			t.Fatalf("expected updated autosave timestamp, got time.After: %v", expTime)
		}
	})
}
