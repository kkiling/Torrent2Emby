package videocontent

type TVShowID struct {
	ID           uint64
	SeasonNumber int
}

type ContentID struct {
	MovieID *uint64
	TVShow  *TVShowID
}
