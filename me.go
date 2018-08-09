package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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

var channels = map[string]string{
	"ModestTim":     "51684790",
	"Jamie254":      "54406241",
	"JamiePineLive": "48234453",
	"Ashturbate":    "46887603",
}

func (session *SessionBody) twitchUser() ([]byte, *TwitchUsers) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", session.Twitch.Auth.AccessToken))
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return nil, nil
	}

	me, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	var s = new(TwitchUsers)
	err = json.Unmarshal([]byte(me), &s)
	if err != nil {
		fmt.Println("whoops:", err)
	}

	return me, s
}

// Me just returns data about the authed user
func Me(w http.ResponseWriter, r *http.Request) {
	session, err := CheckAuthorization(w, r, false, false)
	if err != nil {
		return
	}

	me, s := session.twitchUser()
	if me == nil && s == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(s.Data) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(me))
		return
	}

	checks := jobRunner(session, s)

	jsonResp, _ := json.Marshal(struct {
		Me     *TwitchUsers     `json:"me"`
		Checks []*checkResponse `json:"checks"`
	}{
		Me:     s,
		Checks: checks,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonResp))
}

func jobRunner(session SessionBody, s *TwitchUsers) []*checkResponse {
	var done = make(chan bool)
	var jobs = make(chan *checkResponse, len(channels)*2)
	var checks []*checkResponse

	go func() {
		i := 0
		for {
			job := <-jobs
			checks = append(checks, job)

			if i++; i == len(channels)*2 {
				close(jobs)
				done <- true
				return
			}
		}
	}()

	for chann := range channels {
		for _, check := range []string{"subbed", "followed"} {
			go func(session SessionBody, s *TwitchUsers, chann, check string) {
				jobs <- session.checkStatus(check, chann, s.Data[0].ID)
			}(session, s, chann, check)
		}
	}

	<-done
	return checks
}

type checkResponse struct {
	CreatedAt   string `json:"created_at,omitempty"`
	Error       string `json:"error,omitempty"`
	ChannelName string `json:"channel_name,omitempty"`
	ChannelID   string `json:"channel_id,omitempty"`
	Type        string `json:"type,omitempty"`
}

func (session SessionBody) checkStatus(check string, channel, user string) *checkResponse {
	var endpoint = ""

	switch check {
	case "subbed":
		endpoint = "https://api.twitch.tv/kraken/users/" + user + "/subscriptions/" + channels[channel]
		break
	case "followed":
		endpoint = "https://api.twitch.tv/kraken/users/" + user + "/follows/channels/" + channels[channel]
		break
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint, nil)
	req.Header.Add("Accept", "application/vnd.twitchtv.v5+json")
	req.Header.Add("Client-ID", os.Getenv("TWITCH_CLIENT_ID"))
	req.Header.Add("Authorization", fmt.Sprintf("OAuth %s", session.Twitch.Auth.AccessToken))
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}

	me, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	var s = new(checkResponse)
	err = json.Unmarshal([]byte(me), &s)
	if err != nil {
		fmt.Println("whoops:", err)
	}

	s.ChannelID = channels[channel]
	s.ChannelName = channel
	s.Type = check

	return s
}
