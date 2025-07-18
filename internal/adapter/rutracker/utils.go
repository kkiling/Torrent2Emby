package rutracker

import (
	"fmt"
	"io"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
)

func readerDocument(body io.Reader) (*goquery.Document, error) {
	utf8Reader, err := charset.NewReader(body, "text/html")
	if err != nil {
		return nil, fmt.Errorf("failed to create charset reader: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(utf8Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to create document: %v", err)
	}

	return doc, nil
}

func emptyTorrentResponse() *TorrentResponse {
	return &TorrentResponse{
		Results:      []Torrent{},
		Page:         1,
		TotalResults: 1,
	}
}
