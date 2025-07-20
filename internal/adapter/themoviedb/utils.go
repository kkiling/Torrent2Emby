package themoviedb

import (
	"fmt"
	"time"
)

func (api *API) getImage(path string) *Image {
	if path == "" {
		return nil
	}
	urlImg := api.baseImgUrl.String()
	return &Image{
		ID:       path,
		W92:      fmt.Sprintf("%s/w92%s", urlImg, path),
		W154:     fmt.Sprintf("%s/w154%s", urlImg, path),
		W185:     fmt.Sprintf("%s/w185%s", urlImg, path),
		W342:     fmt.Sprintf("%s/w342%s", urlImg, path),
		W500:     fmt.Sprintf("%s/w500%s", urlImg, path),
		W780:     fmt.Sprintf("%s/w780%s", urlImg, path),
		Original: fmt.Sprintf("%s/original%s", urlImg, path),
	}
}

func parseDate(dateStr string) time.Time {
	const layout = "2006-01-02"
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}
	}
	return t
}
