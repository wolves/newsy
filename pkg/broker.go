package newsy

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"sync"
	"time"
)

type Broker struct {
	Logger io.Writer

	*session
	subs []*subscription
	errs chan error

	running  bool
	stopOnce sync.Once
	cancel   context.CancelFunc
	sync.RWMutex
}

func (b *Broker) Start(ctx context.Context, interval time.Duration) (context.Context, error) {
	if b == nil {
		// TODO: ErrNilBroker (name?)
		return nil, fmt.Errorf("broker is nil")
	}

	b.log("broker:\tstarting\n")

	ctx, cancel := context.WithCancel(ctx)
	b.Lock()
	b.cancel = cancel
	b.running = true
	b.Unlock()

	go func(ctx context.Context) {
		<-ctx.Done()
		cancel()
		b.Stop()
	}(ctx)

	if b.session == nil {
		if err := b.initSession(); err != nil {
			return nil, err
		}
	}

	if b.session != nil && b.loaded {
		go b.startAutoBackup(interval)
	}

	b.log("broker:\tstarted\n")
	return ctx, nil
}

// func (b *Broker) backup() error {
// 	b.log("broker:\tstarting backup\n")
// 	f, err := os.Create("newsy_backup.json")
// 	if err != nil {
// 		if errors.Is(err, fs.ErrNotExist) {
// 			return nil
// 		}
// 		return err
// 	}
//
// 	b.session.backup(f)
// 	return nil
// }

func (b *Broker) Stop() {
	if b == nil {
		return
	}

	b.RLock()
	if !b.running {
		b.RUnlock()
		return
	}
	b.RUnlock()

	b.stopOnce.Do(func() {
		b.log("broker:\tstopping\n")
		b.backup()

		b.Lock()
		defer b.Unlock()

		if b.cancel != nil {
			b.log("broker:\tcancelling context\n")
			b.cancel()
		}
		b.running = false

		for _, sub := range b.subs {
			b.log("broker:\tclosing sub channels\n")
			close(sub.Ch)
		}

		if b.errs != nil {
			close(b.errs)
		}
	})
	b.log("broker:\tstopped\n")
}

// NOTE: Add should probably take in the context returned from Start for cancellation propagation
func (b *Broker) Add(source Source) context.Context {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	if b == nil {
		cancel()
		return ctx
	}

	b.log("broker:\tsource added: [%s]\n", source.Name())
	ch := source.Listen(ctx)

	go func() {
		for {
			select {
			case <-ctx.Done():
				cancel()
				b.log("broker:\tsource stopped: [%s]\n", source.Name())
				b.Stop()
				return
			case a := <-ch:
				b.dispatch(a)
			}
		}
	}()

	return ctx
}

func (b *Broker) Subscribe(ctx context.Context, topics ...Topic) <-chan Article {
	if b == nil {
		return nil
	}

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

func (b *Broker) Errors() chan error {
	if b == nil {
		return nil
	}

	b.Lock()
	defer b.Unlock()

	if b.errs == nil {
		b.errs = make(chan error)
	}

	return b.errs
}

func (b *Broker) initSession() error {
	// TODO: Add check for user defined file with default fallback
	f, err := os.Open("newsy_db.json")
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			b.session = &session{}
			return nil
		}
		return err
	}

	sess, err := Restore(f)
	if err != nil {
		return err
	}

	b.session = sess
	return nil
}

// WIP: Swap Stdout after dev
func (b *Broker) log(msg string, a ...interface{}) {
	// w := io.Discard

	if b != nil && b.Logger == nil {
		b.Logger = os.Stdout
	}

	w := b.Logger

	fmt.Fprintf(w, msg, a...)
}

func (b *Broker) dispatch(a Article) {
	if b == nil {
		return
	}

	b.Lock()
	defer b.Unlock()
	fmt.Printf("Article:\t%v", a)

	for _, sub := range b.subs {
		if !sub.Match(a) {
			continue
		}

		sub.Ch <- a
	}
}

func (b *Broker) unsubscribe(usub *subscription) {
	if b == nil {
		return
	}

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
