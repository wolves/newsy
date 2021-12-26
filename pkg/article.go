package newsy

import (
	"bytes"
	"fmt"
	"strings"
)

type ArticleID int

// Article struct defines an article that is published by a Newsy.Source
type Article struct {
	ID     ArticleID `json:"id,omitempty"`
	Title  string    `json:"title,omitempty"`
	Topics Topics    `json:"topics,omitempty"`
	Source string    `json:"source,omitempty"`
}

// Articles is a list of type Article
type Articles []Article

func (a Article) String() (string, error) {
	if err := validateArticle(a); err != nil {
		return "", err
	}

	id := a.ID
	src := a.Source

	bb := &bytes.Buffer{}
	fmt.Fprintf(bb, "ID: %d\n", id)
	fmt.Fprintf(bb, "Title: %s\n", a.Title)
	fmt.Fprintf(bb, "Topics:\n")
	for _, t := range a.Topics {
		fmt.Fprintf(bb, "\t- %s\n", t)
	}
	fmt.Fprintf(bb, "Source: %s\n", src)

	s := bb.String()
	return strings.TrimSpace(s), nil
}

func validateArticle(a Article) error {
	if a.ID == 0 {
		return ErrInvalidArticle("an article requires an ID")
	}
	if len(a.Title) == 0 {
		return ErrInvalidArticle("an article requires a Title")
	}
	if len(a.Source) == 0 {
		return ErrInvalidArticle("an article requires a Source")
	}
	if len(a.Topics) == 0 {
		return ErrInvalidArticle("an article requires at least 1 Topic")
	}

	return nil
}
