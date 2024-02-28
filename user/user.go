package user

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type UserInformation struct {
	UserId string `json:"id"`
}

func GetName(AccessToken string) string {
	req, err := http.NewRequest(http.MethodGet, "https://q.trap.jp/api/v3/users/me", nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("Authorization", "Bearer "+AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	var userInformation UserInformation
	json.Unmarshal(body, &userInformation)

	return userInformation.UserId
}

func GetIcon(AccessToken string) string {
	req, err := http.NewRequest(http.MethodGet, "https://q.trap.jp/api/v3/users/me/icon", nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("Authorization", "Bearer "+AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	base64Data := base64.StdEncoding.EncodeToString(body)
	log.Println("base64data:")
	log.Println(base64Data)
	return base64Data
}
