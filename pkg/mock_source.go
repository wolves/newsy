package pubsub

import (
	"context"
	"time"
)

type MockSource struct {
	SrcName  string
	Interval time.Duration
	Articles Articles
}

func (src *MockSource) Name() string {
	return src.SrcName
}

func (src *MockSource) Listen(ctx context.Context) <-chan Article {
	ch := make(chan Article)

	go func() {
		tick := time.NewTicker(src.Interval)
		defer tick.Stop()

		i := 0
		for i < len(src.Articles) {
			select {
			case <-ctx.Done():
				close(ch)
				return
			case <-tick.C:
				id := ArticleID(i + 1)
				a := src.Articles[i]
				a.ID = id
				a.Source = src.Name()

				// a := Article{ID: ArticleID(i), Title: "Mock Article", Topics: Topics{"mock topic"}, Source: src.Name()}
				ch <- a
			}
			i++
		}
	}()

	return ch
}
