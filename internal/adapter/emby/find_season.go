package emby

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/kkiling/torrent2emby/internal/adapter/apierr"
)

// GetCatalogInfo возвращает список файлов в торренте по его хешу
func (api *API) GetCatalogInfo(path string) (*CatalogInfo, error) {
	queryParams := url.Values{}
	queryParams.Add("api_key", api.apiKey)
	queryParams.Add("Recursive", "true")
	queryParams.Add("Path", path)
	queryParams.Add("Fields", "Path")

	getUrl := fmt.Sprintf("%s/emby/Items?%s", api.baseAPIUrl.String(), queryParams.Encode())
	resp, err := api.httpClient.Get(getUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to emby items %w", apierr.HandleRequestError(api.logger, err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apierr.HandleStatusCodeError(api.logger, resp)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result struct {
		Items []struct {
			Path     string `json:"Path"`
			Name     string `json:"Name"`
			ID       string `json:"Id"`
			IsFolder bool   `json:"IsFolder"`
			Type     string `json:"Type"`
		}
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, apierr.ContentNotFound
	}
	if len(result.Items) > 1 {
		return nil, fmt.Errorf("multiple items found")
	}

	id, err := strconv.ParseUint(result.Items[0].ID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse item id: %w", err)
	}
	return &CatalogInfo{
		Path:     result.Items[0].Path,
		Name:     result.Items[0].Name,
		ID:       id,
		IsFolder: result.Items[0].IsFolder,
		Type: func(t string) TypeCatalog {
			switch t {
			case "Series":
				return SeriesTypeCatalog
			case "Season":
				return SeasonTypeCatalog
			default:
				return UnknownTypeCatalog
			}
		}(result.Items[0].Type),
	}, nil
}
