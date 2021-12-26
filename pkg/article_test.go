package pubsub

import "testing"

func Test_Article(t *testing.T) {
	t.Parallel()

	validArticle := Article{ID: ArticleID(1), Title: "Test Article", Topics: Topics{"Testing"}, Source: "TDD"}
	noIDArticle := Article{Title: "Test Article", Topics: Topics{"Testing"}, Source: "TDD"}
	noTitleArticle := Article{ID: ArticleID(1), Topics: Topics{"Testing"}, Source: "TDD"}
	noTopicArticle := Article{ID: ArticleID(1), Title: "Test Article", Source: "TDD"}
	noSourceArticle := Article{ID: ArticleID(1), Title: "Test Article", Topics: Topics{}}

	tests := []struct {
		name   string
		art    Article
		expErr bool
	}{
		{name: "valid article", art: validArticle, expErr: false},
		{name: "invalid ID article", art: noIDArticle, expErr: true},
		{name: "invalid title article", art: noTitleArticle, expErr: true},
		{name: "invalid topic article", art: noTopicArticle, expErr: true},
		{name: "invalid source article", art: noSourceArticle, expErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.art.String()
			if (err != nil) != tt.expErr {
				t.Fatalf("got %v, expected %v", err, tt.expErr)
			}
		})
	}

	// has an identifier
	// has a defined title
	// is valid with and without a defined author

	// has valid content
	// has one or more Topics
	// has a valid source association
}
