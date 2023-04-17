package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/ksrnnb/saml-impl/kvs"
	"github.com/labstack/echo/v4"
)

const sessionKey = "session_id"

const sessionLength = 32

func Set(c echo.Context, key string, value string) error {
	ck, err := c.Cookie(sessionKey)
	var sid string
	if err == nil {
		sid = ck.Value
	} else {
		if errors.Is(err, http.ErrNoCookie) {
			sid, err = generateSessionID()
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	ss := getSessionStore(sid)
	ss[key] = value
	kvs.Set(sid, ss)

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

	ss := getSessionStore(sid)

	// TODO: key で判断するのではなく flash メソッドとかをつくる
	if key == "success" || key == "error" {
		v := ss[key]
		delete(ss, key)
		return v, nil
	}
	return ss[key], nil
}

func Clear(c echo.Context) error {
	sid, err := getSessionID(c)
	if err != nil {
		return err
	}
	if sid == "" {
		return nil
	}
	kvs.Set(sid, nil)
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

func getSessionStore(sid string) map[string]string {
	sessionStore := kvs.Get(sid)
	if sessionStore == nil {
		sessionStore = map[string]string{}
	}
	return sessionStore.(map[string]string)
}

func generateSessionID() (string, error) {
	randomBytes := make([]byte, sessionLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	sessionID := base64.URLEncoding.EncodeToString(randomBytes)
	return sessionID, nil
}
