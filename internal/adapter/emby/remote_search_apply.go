package emby

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/kkiling/torrent2emby/internal/adapter/apierr"
)

// RemoteSearchApply устанавливает метаданные для сериала
func (api *API) RemoteSearchApply(embyID, theMovieDBID uint64) error {
	queryParams := url.Values{}
	queryParams.Add("api_key", api.apiKey)
	queryParams.Add("ReplaceAllImages", "true")

	url := fmt.Sprintf("%s/Items/RemoteSearch/Apply/%d?%s", api.baseAPIUrl.String(), embyID, queryParams.Encode())

	// Используем theMovieDBID параметр вместо хардкодного значения
	payload := strings.NewReader(fmt.Sprintf(`{
        "ProviderIds": {
            "Tmdb": "%d"
        }
    }`, theMovieDBID))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Устанавливаем заголовки как в curl
	req.Header.Set("accept", "*/*")
	req.Header.Set("Content-Type", "application/json")

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to emby items: %w", apierr.HandleRequestError(api.logger, err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apierr.HandleStatusCodeError(api.logger, resp)
	}

	return nil
}
