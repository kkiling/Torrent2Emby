package tvshowlibrary

type MovieSearchParams struct {
	Query string
}

type TVShowSearchParams struct {
	Query string
}

type TVShowSearchResult struct {
	Items []TVShowShort
}

type GetTVShowParams struct {
	TVShowID uint64
}

type GetTVShowResult struct {
	Result *TVShow
}

type GetSeasonEpisodesParams struct {
	TVShowID     uint64
	SeasonNumber uint8
}

type GetSeasonEpisodesResult struct {
	Items []Episode
}

type GetTVShowsFromLibraryParams struct {
}

type GetTVShowsFromLibraryResult struct {
	Items []TVShowShort
}
