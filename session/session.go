package session

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

var sessionStore map[string]map[string]string

const sessionKey = "session_id"

func init() {
	if sessionStore == nil {
		sessionStore = map[string]map[string]string{}
	}
}

func Set(c echo.Context, key string, value string) error {
	ck, err := c.Cookie(sessionKey)
	var sid string
	if err == nil {
		sid = ck.Value
	} else {
		if errors.Is(err, http.ErrNoCookie) {
			sid = sessionID()
		} else {
			return err
		}
	}

	sessionStore[sid] = map[string]string{key: value}
	c.SetCookie(
		&http.Cookie{
			Name:  sessionKey,
			Value: sid,
		},
	)
	return nil
}

func Get(c echo.Context, key string) (string, error) {
	sid, err := getSessionID(c)
	if err != nil {
		return "", err
	}
	if sid == "" {
		return "", nil
	}
	return sessionStore[sid][key], nil
}

func Clear(c echo.Context) error {
	sid, err := getSessionID(c)
	if err != nil {
		return err
	}
	if sid == "" {
		return nil
	}
	sessionStore[sid] = map[string]string{}
	return nil
}

func getSessionID(c echo.Context) (string, error) {
	ck, err := c.Cookie(sessionKey)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return "", nil
		}
		return "", err
	}
	return ck.Value, nil
}

func sessionID() string {
	// TODO: secure random value
	return "dummy_session_id"
}
