package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// TwitchUsers is a stuct for the helix users endpoint
type TwitchUsers struct {
	Data []struct {
		ID              string `json:"id"`
		Login           string `json:"login"`
		DisplayName     string `json:"display_name"`
		Type            string `json:"type"`
		BroadcasterType string `json:"broadcaster_type"`
		Description     string `json:"description"`
		ProfileImageURL string `json:"profile_image_url"`
	} `json:"data"`
}

// Me just returns data about the authed user
func Me(w http.ResponseWriter, r *http.Request) {
	session, err := CheckAuthorization(w, r, false, false)
	if err != nil {
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", session.Twitch.Auth.AccessToken))
	res, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}
