package themoviedb

// SearchQuery contains search parameters
type SearchQuery struct {
	Language Language
	Query    string `validate:"min=3"`
	Page     int    `validate:"gte=1"`
	PerPage  int    `validate:"gte=1,lte=20"`
}

// MovieSearchResponse contains movie search results
type MovieSearchResponse struct {
	Page         int
	TotalResults int
	Results      []MovieShort
}

// TVShowSearchResponse contains TV show search results
type TVShowSearchResponse struct {
	Page         int
	TotalPages   int
	TotalResults int
	Results      []TVShowShort
}
