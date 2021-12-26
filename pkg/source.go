package pubsub

import "context"

// Source provides Articles to the Broker
type Source interface {

	// Returns the name of the source
	Name() string

	// Listen takes a context for cancellation
	// and returns a listen-only channel of
	// of type Article published by the Source
	Listen(ctx context.Context) <-chan Article
}

// func (m *Source) Report() {
// 	for _, a := range m.articles {
// 		fmt.Printf("Source publishing: %v\n", a)
// 		go func(art Article) {
// 			m.ch <- art
// 		}(a)
// 	}
// }

// INFO: Abilities
// - Register with the main service Register/Attach
// - Unregister with the main service Unregister/Detach
// - Do its work to process the payload
// - Deliver payload(articles) Report
//
// type MockSource struct {
// 	// name     string
// 	pubRate  time.Duration
// 	topics   Topics
// 	articles Articles
// 	ch       chan Article
// }
//
// func (m *MockSource) Report() {
//
// 	for _, a := range m.articles {
// 		time.Sleep(m.pubRate)
//
// 		if a.Topics == nil {
// 			a.Topics = m.topics
// 		}
//
// 		fmt.Printf("Mock Source publishing: %v\n", a)
//
// 		go func(art Article) {
// 			m.ch <- art
// 		}(a)
// 	}
//
// }
