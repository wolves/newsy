package pubsub

// Category is a string that is used to define the category of an Article
type Topic string

// Topics is a list of type Category
type Topics []Topic

func (tops Topics) Match(query ...Topic) bool {
	for _, t := range tops {
		for _, qt := range query {
			if t == qt {
				return true
			}
		}
	}
	return false
}
