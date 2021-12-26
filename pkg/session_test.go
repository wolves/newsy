package pubsub

import (
	"bytes"
	"testing"
)

func Test_Session(t *testing.T) {
	t.Parallel()
	mockTopics := Topics{"testing", "go"}
	mockArticles := Articles{{ID: 135}}
	mockStoreJson := []byte(`{"topics":["testing","go"],"articles":[{"id": 135}]}`)

	t.Run("loads state from backup file", func(t *testing.T) {
		s := &Session{}

		err := s.Load(mockStoreJson)
		assertNoError(t, err)

		if !s.Loaded() {
			t.Fatalf("expected state to be loaded, got %v", s.loaded)
		}

		gotTop := len(s.Topics)
		expTop := len(mockTopics)
		assertLenMatch(t, gotTop, expTop)

		gotArt := len(s.Articles)
		expArt := len(mockArticles)
		assertLenMatch(t, gotArt, expArt)
	})

	t.Run("saves state to backup file", func(t *testing.T) {
		s := &Session{
			Topics:   mockTopics,
			Articles: mockArticles,
		}

		bb := &bytes.Buffer{}

		err := s.Save(bb)
		if err != nil {
			t.Fatalf("unexpected error while saving state: ERR %v", err)
		}

		act := bb.Bytes()
		exp := mockStoreJson

		if bytes.Contains(act, exp) {
			t.Fatalf("got %v, expected %v", exp, act)
		}
	})
}
