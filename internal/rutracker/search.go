package rutracker

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/kkiling/torrent2emby/internal/apierr"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

func (api *Api) SearchTorrents(query string) (*TorrentResponse, error) {
	if err := api.login(); err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	api.logger.Debugf("Search torrents: %s", query)

	searchURL := api.baseAPIUrl.String() + "tracker.php?nm=" + url.QueryEscape(query)
	resp, err := api.httpClient.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to search torrents: %w", apierr.HandleRequestError(api.logger, err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apierr.HandleStatusCodeError(api.logger, resp)
	}

	// Создаем reader с правильной кодировкой
	doc, err := readerDocument(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create document: %v", err)
	}

	table := doc.Find("table#tor-tbl")
	if table.Length() == 0 {
		return emptyTorrentResponse(), nil
	}

	firstRow := table.Find("tbody tr").First()
	if firstRow.Find("td").First().Text() == "Не найдено" {
		return emptyTorrentResponse(), nil
	}

	var results []Torrent

	table.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
		cols := row.Find("td")
		if cols.Length() < 10 {
			return
		}

		forum := cols.Eq(2).Text()
		titleLink := cols.Eq(3).Find("a").First()
		title := titleLink.Text()
		href := api.baseAPIUrl.String() + titleLink.AttrOr("href", "")
		author := cols.Eq(4).Text()
		size := cols.Eq(5).Text()
		seeds := cols.Eq(6).Text()
		leeches := cols.Eq(7).Text()
		downloads := cols.Eq(8).Text()
		addedDate := cols.Eq(9).Text()

		results = append(results, Torrent{
			Title:     strings.TrimSpace(title),
			Href:      href,
			Forum:     strings.TrimSpace(forum),
			Author:    strings.TrimSpace(author),
			Size:      strings.TrimSpace(size),
			Seeds:     strings.TrimSpace(seeds),
			Leeches:   strings.TrimSpace(leeches),
			Downloads: strings.TrimSpace(downloads),
			AddedDate: strings.TrimSpace(addedDate),
		})
	})

	// Сортировка по downloads (в порядке убывания)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Downloads > results[j].Downloads
	})

	return &TorrentResponse{
		Results:      results,
		Page:         1,
		TotalResults: len(results),
	}, nil
}
