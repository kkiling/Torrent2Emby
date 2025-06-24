package videocontent

type CreateNewVideoContentParams struct {
	MediaID MediaID
}

type ChangeSearchQueryParams struct {
	Query string
}

type ChoseTorrentParams struct {
	Href string
}

type ChoseFileMatchesParams struct {
	// TODO: доделать
}
