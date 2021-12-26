package pubsub

import (
	"context"
	"testing"
	"time"
)

func Test_Broker(t *testing.T) {
	t.Parallel()

	t.Run("Start method", func(t *testing.T) {
		b := NewBroker()
		saveInterval := time.Second * 1

		startService(t, b, saveInterval)
		assertStarted(t, b)

		// if n.Session.saveInterval != saveInterval {
		// 	t.Fatalf("got save interval %v, expected %v", n.Session.saveInterval, saveInterval)
		// }

		// FIX: Just randomly calling this again feels weird esp with no err check
		// should be able to be started multiple times
		startService(t, b, saveInterval)
	})

	t.Run("Stop method", func(t *testing.T) {
		t.Run("can be called by the service", func(t *testing.T) {
			b := NewBroker()
			startService(t, b, time.Minute*1)

			// should be able to be stopped by the news service itself.
			b.Stop()
			if !b.stopped {
				t.Fatal("expected Newsy to be stopped")
			}
		})
		// should be able to be stopped multiple times.
		t.Run("can be called while already stopped", func(t *testing.T) {
			b := NewBroker()

			// should be able to be stopped by the news service itself.
			b.Stop()
			assertStopped(t, b)
		})
		// should be able to be stopped by the end user.
		t.Run("can be called by the end-user", func(t *testing.T) {
			b := NewBroker()
			startService(t, b, time.Second*1)
			assertStarted(t, b)

			b.Stop()
			assertStopped(t, b)
		})
		// should stop all sources and subscribers when it is stopped.

		// should save the state of the news service to the backup file, in JSON format, when it is stopped.

		// should not be able to be stopped by the news sources.
		// should not be able to be stopped by the subscribers.
	})

	t.Run("Subscribe", func(t *testing.T) {
		b := NewBroker()
		// ctx := startService(t, n, time.Second*1)

		ts := Topic("testing")

		// should be able to subscribe to the news service and receive news stories for one or more topics.
		_ = b.Subscribe(context.Background(), ts)

		if len(b.subs) != 1 {
			t.Fatal("expected topic and chan to be added to subs")
		}

		for _, sub := range b.subs {
			got := sub.Topics[0]
			want := ts

			if got != want {
				t.Fatalf("expected topic %v, got %v", got, want)
			}
		}
	})

	t.Run("Add", func(t *testing.T) {
		b := NewBroker()
		ctx := startService(t, b, time.Second)

		arts := Articles{
			{Title: "Art 1", Topics: Topics{"tdd"}},
			{Title: "Art 2", Topics: Topics{"go"}},
			{Title: "Art 3", Topics: Topics{"neovim"}},
		}

		m := &MockSource{
			SrcName:  "The Mockery",
			Interval: time.Millisecond * 10,
			Articles: arts,
		}

		b.Add(m)
		ts := Topics{"tdd", "go"}
		ch := b.Subscribe(context.Background(), ts...)

		var got []Article
		exp := len(ts)

		go func(ch <-chan Article) {
			for a := range ch {
				b.Lock()
				got = append(got, a)
				b.Unlock()

				b.RLock()
				done := len(got) >= exp
				b.RUnlock()
				if done {
					b.Stop()
				}
			}
		}(ch)

		<-ctx.Done()

		if len(got) != exp {
			t.Errorf("got %v, expected: %v", len(got), exp)
		}
	})

	t.Run("incrementally saves session", func(t *testing.T) {
		orig := time.Now()
		s := &Session{loaded: true, timestamp: orig}
		b := &Broker{
			// subs:    make(map[Topic][]chan Article),
			errs:    make(chan error),
			stopped: true,
			Session: s,
		}

		i := time.Millisecond * 50
		b.Start(context.Background(), i)

		// FIX: Change this to have autosave use a channel and capture that?
		time.Sleep(2 * i)

		s.RLock()
		expTime := s.timestamp.After(orig)
		s.RUnlock()

		if !expTime {
			t.Fatalf("expected updated autosave timestamp, got time.After: %v", expTime)
		}
	})

	t.Run("searching article id/s", func(t *testing.T) {
		articles := Articles{}
		ids := []int{135, 246, 975}
		for _, id := range ids {
			articles = append(articles, Article{ID: ArticleID(id)})
		}
		b := NewBroker()
		b.Articles = articles
		tests := []struct {
			name        string
			articles    Articles
			ids         []int
			resultCount int
			expErr      bool
		}{
			{name: "returns matching articles", articles: b.Articles, ids: ids[:2], resultCount: 2, expErr: false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := b.Search(tt.ids...)
				if (err != nil) != tt.expErr {
					t.Fatal(err)
				}

				exp := tt.resultCount

				if len(got) != exp {
					t.Fatalf("got %d search results, expected %d", len(got), exp)
				}
			})
		}
	})

	t.Run("Errors method", func(t *testing.T) {
	})

	t.Run("unsub removes subscription", func(t *testing.T) {
		b := NewBroker()
		sub1 := &subscription{ID: subId(1), Topics: Topics{"tdd"}}
		sub2 := &subscription{ID: subId(2), Topics: Topics{"tdd"}}
		sub3 := &subscription{ID: subId(3), Topics: Topics{"tdd"}}
		origSubs := subscriptions{sub1, sub2, sub3}

		b.subs = origSubs

		b.unsubscribe(sub2)
		if len(b.subs) >= len(origSubs) {
			t.Fatal("expected a subscription to be removed")
		}

		for _, s := range b.subs {
			if s.ID == subId(2) {
				t.Fatal("expected unsubscribe to remove correct sub from broker subscriptions")
			}
		}
	})
}

// Test Assertion Helpers
func assertNoError(t testing.TB, e error) {
	t.Helper()

	if e != nil {
		t.Fatalf("got unexpected error %v", e)
	}
}

func assertStarted(t testing.TB, b *Broker) {
	t.Helper()
	if b.stopped {
		t.Fatal("expected manager to be started and stopped field to be false")
	}
}

func assertStopped(t testing.TB, b *Broker) {
	t.Helper()

	b.RLock()
	if !b.stopped {
		b.RUnlock()
		t.Fatal("expected Newsy to be stopped")
	}
	b.RUnlock()
}

func assertLenMatch(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Fatalf("expected topics len: %v, got: %v", got, want)
	}
}

// Test Service Helpers
func startService(t testing.TB, b *Broker, si time.Duration) context.Context {
	t.Helper()

	ctx, err := b.Start(context.Background(), si)
	assertNoError(t, err)

	return ctx
}
