package pubsub

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

type Broker struct {
	*Session

	subs []*subscription
	errs chan error

	stopped  bool
	cancel   context.CancelFunc
	stopOnce sync.Once
	sync.RWMutex
}

func NewBroker() *Broker {
	return &Broker{
		Session: &Session{},
		errs:    make(chan error),
		stopped: true,
	}
}

func (b *Broker) Start(ctx context.Context, interval time.Duration) (context.Context, error) {
	ctx, cancel := context.WithCancel(ctx)

	b.Lock()
	b.stopped = false
	b.Unlock()

	b.Lock()
	b.cancel = cancel
	b.Unlock()

	go func(ctx context.Context) {
		<-ctx.Done()

		cancel()

		b.Stop()
	}(ctx)

	if b.Loaded() {
		s, err := fetchBackup("newsy_backup.json")
		if err != nil {
			return ctx, err
		}
		b.Load(s)
	}

	if b.Loaded() {
		go b.StartAutoSave(interval)
	}

	return ctx, nil
}

func (b *Broker) Stop() {
	b.RLock()
	if b.stopped {
		b.RUnlock()
		return
	}
	b.RUnlock()

	b.stopOnce.Do(func() {
		// Save backup
		b.Lock()
		defer b.Unlock()

		b.cancel()
		b.stopped = true

		for _, sub := range b.subs {
			close(sub.Ch)
		}

		if b.errs != nil {
			close(b.errs)
		}
	})
}

func (b *Broker) Subscribe(ctx context.Context, topics ...Topic) <-chan Article {
	sub := &subscription{
		Topics: topics,
		Ch:     make(chan Article),
		ID:     subId(len(b.subs) + 1),
	}

	b.Lock()
	// PERF: Come back and benchmark using this vs map with subId key
	b.subs = append(b.subs, sub)
	defer b.Unlock()

	go func() {
		<-ctx.Done()
		b.unsubscribe(sub)
	}()

	return sub.Ch
}

// NOTE: Add should probably take in the context returned from Start for cancellation propagation
func (b *Broker) Add(source Source) context.Context {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	ch := source.Listen(ctx)

	go func() {
		for {
			select {
			case <-ctx.Done():
				cancel()
				b.Stop()
			case a := <-ch:
				b.dispatch(a)
			}
		}
	}()

	return ctx
}

func (b *Broker) Errors() chan error {
	b.Lock()
	defer b.Unlock()

	if b.errs == nil {
		b.errs = make(chan error)
	}

	return b.errs
}

func (b *Broker) Search(ids ...int) (Articles, error) {
	res := Articles{}

	for _, i := range ids {
		for _, a := range b.Articles {
			if ArticleID(i) == a.ID {
				res = append(res, a)
			}
		}
	}
	return res, nil
}

func (b *Broker) dispatch(a Article) {
	b.Lock()
	defer b.Unlock()

	for _, sub := range b.subs {
		if !sub.Match(a) {
			fmt.Printf("No Matching Subs for Article: %+v\n", a)
			continue
		}

		fmt.Printf("Dispatching Article: %+v\n", a)
		sub.Ch <- a
	}
}

func (b *Broker) unsubscribe(usub *subscription) {
	b.RLock()
	if b.subs == nil {
		b.RUnlock()
		return
	}

	subs := b.subs
	b.RUnlock()

	// PERF: Lookup if this is more or less perf than map
	// and if recreating slice to remove one sub is faster than other ways
	newSubs := make([]*subscription, 0, len(subs)-1)
	for _, s := range b.subs {
		if s != usub {
			newSubs = append(newSubs, s)
		}
	}

	b.Lock()
	b.subs = newSubs
	b.Unlock()
}

func fetchBackup(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {

		if errors.Is(err, os.ErrNotExist) {

			b := []byte(`{}`)
			os.WriteFile(filename, []byte(`{}`), 0666)
			// n.stateLoaded = true
			return b, nil
		}
		return nil, err
	}
	return data, nil
}
