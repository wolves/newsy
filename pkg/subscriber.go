package newsy

type subId int

type subscription struct {
	ID     subId
	Topics Topics

	Ch chan Article
}

type subscriptions []*subscription

func (s *subscription) Match(a Article) bool {
	if s == nil {
		return false
	}
	return s.Topics.Match(a.Topics...)
}
