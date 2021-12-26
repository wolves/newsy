package newsy

// ErrInvalidSource is returned when the news `Source` is invalid due to an invalid attribute/field
type ErrInvalidSource string

// func (e ErrInvalidSource) Error() string {}

// ErrInvalidArticle is returned when an `Article` is invalid due to an invalid attribute/field
type ErrInvalidArticle string

func (e ErrInvalidArticle) Error() string {
	return string(e)
}

// ErrInvalidCategory is returned when a `Category` is created with an incorrect type
type ErrInvalidCategory string

// func (e ErrInvalidCategory) Error() string {}

// ErrNewsyServiceStopped is returned if a Newsy component attempts to interact
// with it after it has been stopped
type ErrNewsyServiceStopped struct{}

func (ErrNewsyServiceStopped) Error() string {
	return "service is stopped"
}
