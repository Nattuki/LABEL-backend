package handler

import (
	"net/http"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	IDToken     string `json:"id_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

var (
	ClientID string
)

func OAuthHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
}
