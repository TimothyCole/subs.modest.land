package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type token struct {
	// Successful
	AccessToken  string   `json:"access_token,omitempty"`
	RefreshToken string   `json:"refresh_token,omitempty"`
	Scopes       []string `json:"scopes,omitempty"`
	ExpiresIn    int      `json:"expires_in,omitempty"`

	// Error
	Status  int    `json:"status,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

func getToken(body []byte) (*token, error) {
	var s = new(token)
	err := json.Unmarshal(body, &s)
	return s, err
}

// TwitchAuthorize handles the callback from Twitch when a user authenticates
func TwitchAuthorize(w http.ResponseWriter, r *http.Request) {
	var code = r.URL.Query().Get("code")

	session, err := CheckAuthorization(w, r, true, true)
	if err != nil {
		return
	}

	if code == "" {
		w.Write([]byte(`Authorization Denied`))
		return
	}

	validate, err := url.Parse("https://api.twitch.tv/kraken/oauth2/token")
	if err != nil {
		panic(err)
	}
	query := validate.Query()
	query.Set("client_id", os.Getenv("TWITCH_CLIENT_ID"))
	query.Set("client_secret", os.Getenv("TWITCH_CLIENT_SECRET"))
	query.Set("code", code)
	query.Set("grant_type", "authorization_code")
	query.Set("redirect_uri", os.Getenv("TWITCH_CLIENT_REDIRECT"))
	validate.RawQuery = query.Encode()

	res, err := http.Post(validate.String(), "application/json", nil)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	data, err := getToken([]byte(body))

	if err == nil {
		session.Twitch.Auth.AccessToken = data.AccessToken
		session.Twitch.Auth.RefreshToken = data.RefreshToken
		session.Twitch.Auth.ExpiresIn = data.ExpiresIn
		session.Update()
	}

	w.Write([]byte(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<meta http-equiv="X-UA-Compatible" content="ie=edge">
			<title>Modest Land Twitch Subscribers</title>
		</head>
		<body>
			<p>You can now close this window.</p>
			<script>window.close();</script>
		</body>
		</html>
	`))
}
