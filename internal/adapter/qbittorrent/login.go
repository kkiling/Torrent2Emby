package qbittorrent

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/kkiling/torrent2emby/internal/adapter/apierr"
)

// saveCookies сохраняет куки в файл
func (api *Api) saveCookies() error {
	api.logger.Debugf("Save cookies")

	cookiesPath := filepath.Join(api.cookiesDir, cookeFile)
	file, err := os.Create(cookiesPath)
	if err != nil {
		return fmt.Errorf("failed to create cookies prepare: %v", err)
	}
	defer file.Close()

	cookies := api.httpClient.Jar.Cookies(api.baseAPIUrl)
	data, err := json.Marshal(cookies)
	if err != nil {
		return fmt.Errorf("failed to marshal cookies: %w", err)
	}

	// Записываем данные в файл
	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("failed to write cookies to prepare: %w", err)
	}

	return nil
}

func (api *Api) loadCookies() (bool, error) {
	api.logger.Debugf("Load cookies")

	cookiesPath := filepath.Join(api.cookiesDir, cookeFile)
	data, err := os.ReadFile(cookiesPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to read cookies prepare: %w", err)
	}

	var cookies []*http.Cookie
	if err := json.Unmarshal(data, &cookies); err != nil {
		return false, fmt.Errorf("failed to unmarshal cookies: %w", err)
	}

	if len(cookies) > 0 {
		api.httpClient.Jar.SetCookies(api.baseAPIUrl, cookies)
		return true, nil
	}

	return false, nil
}

func (api *Api) tryLogin() (bool, error) {
	api.logger.Debugf("Try login")
	form := url.Values{}
	form.Set("username", api.username)
	form.Set("password", api.password)

	loginUrl := api.baseAPIUrl.String() + "/api/v2/auth/login"
	resp, err := api.httpClient.PostForm(loginUrl, form)
	if err != nil {
		return false, fmt.Errorf("login request failed: %v", apierr.HandleRequestError(api.logger, err))
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	} else if resp.StatusCode != http.StatusOK {
		return false, apierr.HandleStatusCodeError(api.logger, resp)
	}

	return false, nil
}

func (api *Api) removeCookies() error {
	api.logger.Debugf("Remove cookies")
	cookiesPath := filepath.Join(api.cookiesDir, cookeFile)
	// Проверяем существует ли файл
	if _, err := os.Stat(cookiesPath); os.IsNotExist(err) {
		return nil
	}

	// Удаляем файл
	err := os.Remove(cookiesPath)
	if err != nil {
		return fmt.Errorf("failed to remove cookies prepare: %w", err)
	}

	api.logger.Debugf("Successfully removed cookies prepare: %s", cookiesPath)
	return nil
}

func (api *Api) login() error {
	if isLoad, err := api.loadCookies(); err != nil {
		return fmt.Errorf("failed to load cookies: %v", err)
	} else if isLoad {
		return nil
	}

	if success, err := api.tryLogin(); err != nil {
		return fmt.Errorf("tryLogin: %w", err)
	} else if success {
		if err = api.saveCookies(); err != nil {
			return fmt.Errorf("failed to save cookies: %v", err)
		}
		return nil
	}

	if errRemove := api.removeCookies(); errRemove != nil {
		api.logger.Errorf("failed to remove cookies: %v", errRemove)
	}

	return apierr.AuthenticationFailedErr
}
