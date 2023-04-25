package session

import (
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func SetFlash(c echo.Context, key string, value string) {
	c.SetCookie(
		&http.Cookie{
			Name:     key,
			Value:    base64.StdEncoding.EncodeToString([]byte(value)),
			HttpOnly: true,
			Path:     "/",
		},
	)
}

func GetFlash(c echo.Context, key string) (string, error) {
	ck, err := c.Cookie(key)
	var v []byte
	if errors.Is(err, http.ErrNoCookie) {
		return "", nil
	}

	if err == nil {
		v, err = base64.StdEncoding.DecodeString(ck.Value)
		if err != nil {
			return "", err
		}
	}

	c.SetCookie(
		&http.Cookie{
			Name:     key,
			Value:    ck.Value,
			HttpOnly: true,
			Path:     "/",
			MaxAge:   -1,
		},
	)
	return string(v), nil
}
