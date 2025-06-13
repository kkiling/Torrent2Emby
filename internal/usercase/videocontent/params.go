package videocontent

type CreateNewVideoContentParams struct {
	MediaID MediaID
}

type GenerateSearchQueryParams struct {
	MediaID MediaID
}

type SearchTorrentParams struct {
	SearchQuery string
}

type ChangeSearchQueryParams struct {
	Query string
}

type ChoseTorrentParams struct {
	Href string
}

type GetMagnetLinkParams struct {
	Href string
}

type CreateTorrentParams struct {
	MediaID MediaID
	Magnet  string
}
