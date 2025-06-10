package search

type MovieSearchParams struct {
	Query string
}

type MovieSearchResult struct {
	Results []MovieShort
}

type TVShowSearchParams struct {
	Query string
}

type TVShowSearchResult struct {
	Results []TVShowShort
}

type GetTVShowParams struct {
	TVShowID uint64
}

type GetTVShowResult struct {
	TVShowInfo TVShowSearch
}

type GetSeasonEpisodesParams struct {
	TVShowID     uint64
	SeasonNumber int
}

type GetSeasonEpisodesResult struct {
	Episodes []Episode
}

type AddTVShowToLibraryParams struct {
	TVShowID uint64
}

type GetTVShowsFromLibraryParams struct {
}

type GetTVShowsFromLibraryResult struct {
	Results []TVShowShort
}
