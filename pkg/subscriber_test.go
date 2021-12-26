package newsy

import "testing"

func Test_Subscriber(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		sub    *subscription
		topics Topics
		exp    bool
	}{
		{name: "match without a subscription", topics: Topics{"tdd"}},
		{name: "match without topics", sub: &subscription{}, topics: Topics{}},
		{
			name: "successful",
			sub: &subscription{
				Topics: Topics{"go"},
			},
			topics: Topics{"go"},
			exp:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.sub.Match(Article{Topics: tt.topics})

			if got != tt.exp {
				t.Fatalf("expected %v, got %v", tt.sub, tt.topics)
			}
		})
	}
}
