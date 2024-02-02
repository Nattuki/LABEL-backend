package handler

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	IDToken      string `json:"id_token"`
}

const (
	lenState        = 30
	lenCodeVerifier = 64
)

var (
	ClientID string
)

func HandleLogin(c echo.Context) error {
	sess, err := session.Get("LABEL_session", c)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "Failed to get the session.")
	}

	state, err := RandomString(lenState)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "Failed to generate the state.")
	}
	sess.Values["state"] = state

	codeVerifier, err := RandomString(lenCodeVerifier)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "Failed to generate the codeVerifier.")
	}
	sess.Values["code_verifier"] = codeVerifier
	sess.Save(c.Request(), c.Response())

	b := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b[:])

	u, err := url.Parse("https://q.trap.jp/api/v3/oauth2/authorize")
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "Failed to parse the url.")
	}
	q := u.Query()
	q.Add("response_type", "code")
	q.Add("client_id", ClientID)
	q.Add("state", state)
	q.Add("code_challenge", codeChallenge)
	q.Add("code_challenge_method", "S256")
	u.RawQuery = q.Encode()
	return c.Redirect(http.StatusFound, u.String())
}

func HandleCallback(c echo.Context) error {
	sess, err := session.Get("LABEL_session", c)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "Failed to get the session.")
	}
	u := c.Request().URL

	if u.Query().Get("state") != sess.Values["state"] {
		return c.String(http.StatusInternalServerError, "Invalid state.")
	}

	codeVerifier := sess.Values["code_verifier"].(string)
	sess.Values["state"] = ""
	sess.Values["code_verifier"] = ""
	sess.Save(c.Request(), c.Response())

	code := u.Query().Get("code")

	log.Println(code)
	log.Println(codeVerifier)

	q := url.Values{}
	q.Add("grant_type", "authorization_code")
	q.Add("client_id", ClientID)
	q.Add("code", code)
	q.Add("code_verifier", codeVerifier)
	req, err := http.NewRequest(http.MethodPost, "https://q.trap.jp/api/v3/oauth2/token", strings.NewReader(q.Encode()))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to write into the new request.")
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	log.Println(resp)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get the response from the authorization server.")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to read the body.")
	}
	var token TokenResponse
	json.Unmarshal(body, &token)
	sess.Values["access_token"] = token.AccessToken
	sess.Values["token_type"] = token.TokenType
	sess.Values["expires_in"] = token.ExpiresIn
	sess.Values["refresh_token"] = token.RefreshToken
	sess.Values["scope"] = token.Scope
	sess.Values["id_token"] = token.IDToken
	sess.Save(c.Request(), c.Response())

	log.Println(sess.Values)
	return c.Redirect(http.StatusFound, "http://localhost:3000/")
}

func RandomString(length int) (string, error) {
	if length < 0 {
		return "", fmt.Errorf("cannot generate random string of negative length %d", length)
	}
	var s strings.Builder
	for s.Len() < length {
		r, err := rand.Int(rand.Reader, big.NewInt(1<<60))
		if err != nil {
			return "", err
		}
		s.WriteString(fmt.Sprintf("%015x", r))
	}
	return s.String()[:length], nil
}
