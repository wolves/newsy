package pubsub

import "testing"

func Test_Source(t *testing.T) {
	t.Parallel()

	t.Run("Report", func(t *testing.T) {
		// should be able to publish stories for any category, or categories, they wish to define.
		// a := Article{Id: 1, Title: "Test Title", Topics: Topics{"testing"}}
		//   targ := make(chan<- Article)
		// s := Source{Name: "The Tester"}

		// err := s.Report(t,a)
		// assertNoErr(t, err)

	})

	t.Run("Schedule method", func(t *testing.T) {
		// should be free to deliver stories as frequently, or as infrequently, as they wish.
	})

	t.Run("canceling", func(t *testing.T) {
		// should be able to self-cancel.
		// should not be effected by the removal of another news source.
	})

	t.Run("stopping", func(t *testing.T) {
		// should be stopped by the news service when it is stopped.
	})

	// should not be effected by the removal of a subscriber.
}
