package themoviedb

// SearchQuery contains tvshowlibrary parameters
type SearchQuery struct {
	Language Language
	Query    string `validate:"min=3"`
	Page     int    `validate:"gte=1"`
	PerPage  int    `validate:"gte=1,lte=20"`
}

// MovieSearchResponse contains movie tvshowlibrary results
type MovieSearchResponse struct {
	Page         int
	TotalResults int
	Results      []MovieShort
}

// TVShowSearchResponse contains TV show tvshowlibrary results
type TVShowSearchResponse struct {
	Page         int
	TotalPages   int
	TotalResults int
	Results      []TVShowShort
}
