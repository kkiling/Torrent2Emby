package rutracker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func (api *API) saveCookies(filename string) error {
	api.logger.Debugf("Save cookies")

	cookiesPath := filepath.Join(api.cookiesDir, filename)
	file, err := os.Create(cookiesPath)
	if err != nil {
		return fmt.Errorf("failed to create cookies file: %v", err)
	}
	defer file.Close()

	cookies := api.httpClient.Jar.Cookies(api.baseAPIUrl)
	data, err := json.Marshal(cookies)
	if err != nil {
		return fmt.Errorf("failed to marshal cookies: %w", err)
	}

	// Записываем данные в файл
	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("failed to write cookies to file: %w", err)
	}

	return nil
}

func (api *API) loadCookies(filename string) (bool, error) {
	api.logger.Debugf("Load cookies")

	cookiesPath := filepath.Join(api.cookiesDir, filename)
	data, err := os.ReadFile(cookiesPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to read cookies file: %w", err)
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

func (api *API) removeCookies(filename string) error {
	api.logger.Debugf("Remove cookies")
	cookiesPath := filepath.Join(api.cookiesDir, filename)
	// Проверяем существует ли файл
	if _, err := os.Stat(cookiesPath); os.IsNotExist(err) {
		api.logger.Debugf("Cookies file does not exist: %s", cookiesPath)
		return fmt.Errorf("cookies file does not exist")
	}

	// Удаляем файл
	err := os.Remove(cookiesPath)
	if err != nil {
		return fmt.Errorf("failed to remove cookies file: %w", err)
	}

	api.logger.Debugf("Successfully removed cookies file: %s", cookiesPath)
	return nil
}

func (api *API) tryLogin() (bool, error) {
	api.logger.Debugf("Try login")
	loginData := url.Values{
		"login_username": {api.username},
		"login_password": {api.password},
		"login":          {"вход"},
	}

	loginUrl := api.baseAPIUrl.String() + "login.php"
	resp, err := api.httpClient.PostForm(loginUrl, loginData)
	if err != nil {
		return false, fmt.Errorf("login request failed: %v", handleRequestError(api.logger, err))
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK && resp.Request.URL.Path == "/forum/index.php" {
		if err = api.saveCookies(cookeFile); err != nil {
			return false, fmt.Errorf("failed to save cookies: %v", err)
		}
		return true, nil
	} else if resp.StatusCode != http.StatusOK {
		return false, handleStatusCodeError(api.logger, resp)
	}

	return false, nil
}

func (api *API) login() error {
	if isLoad, err := api.loadCookies(cookeFile); err != nil {
		return fmt.Errorf("failed to load cookies: %v", err)
	} else if isLoad {
		return nil
	}

	if success, err := api.tryLogin(); err != nil {
		return err
	} else if !success {
		if errRemove := api.removeCookies(cookeFile); errRemove != nil {
			api.logger.Errorf("failed to remove cookies: %v", errRemove)
		}
		return AuthenticationFailedErr
	}

	return nil
}
