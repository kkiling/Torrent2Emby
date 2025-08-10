package apierr

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/kkiling/goplatform/log"
)

var (
	// NotAuthorizedErr в момент запроса оказалось что клиент не авторизован, скорее всего протухли куки
	// Что бы исправить достаточно еще раз попробовать выполнить запрос
	NotAuthorizedErr = fmt.Errorf("not authorized")
	// AuthenticationFailedErr во время попытки логина на сайте произошла ошибка
	// Либо указан неверный пароль, либо сайт просит ввести капчу
	AuthenticationFailedErr = fmt.Errorf("authentication failed")
	// ServiceUnavailableErr сервис недоступен, можно попробовать повторить попытку
	ServiceUnavailableErr = fmt.Errorf("service unavailable")
	// ContentNotFound контент не был найден
	ContentNotFound = fmt.Errorf("content not found")
)

func HandleStatusCodeError(log log.Logger, resp *http.Response) error {
	log.Errorf("Response status %s: %s", resp.Request.URL.Path, resp.Status)
	if resp.StatusCode == http.StatusNotFound {
		return ContentNotFound
	} else if resp.StatusCode == http.StatusUnauthorized {
		return NotAuthorizedErr
	} else if resp.StatusCode == 522 || resp.StatusCode == 521 || resp.StatusCode == http.StatusForbidden {
		return ServiceUnavailableErr
	}
	return fmt.Errorf("tvshowlibrary failed with status code: %d", resp.StatusCode)
}

func HandleRequestError(log log.Logger, err error) error {
	log.Errorf("Request failed: %v", err)
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return ServiceUnavailableErr
	}
	if strings.Contains(err.Error(), "connection reset by peer") {
		return ServiceUnavailableErr
	}
	if strings.Contains(err.Error(), "connection refused") {
		return ServiceUnavailableErr
	}
	if strings.Contains(err.Error(), "TLS handshake timeout") {
		return ServiceUnavailableErr
	}

	return err
}
