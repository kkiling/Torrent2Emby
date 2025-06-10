package search

import "time"

type Image struct {
	W92      string
	W154     string
	W185     string
	W342     string
	W500     string
	W780     string
	Original string
}

// MovieShort contains basic movie information
type MovieShort struct {
	ID            uint64
	Title         string
	OriginalTitle string
	Overview      string
	PosterPath    *Image
	ReleaseDate   time.Time
	VoteAverage   float64
	VoteCount     int
	Popularity    float64
}
