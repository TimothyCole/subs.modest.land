package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"regexp"

	"cloud.google.com/go/datastore"
)

type mcUser struct {
	Minecraft string `datastore:"minecraft" json:"minecraft"`
}

// Get the mc name for the logged in user
func mcWhitelistGET(w http.ResponseWriter, r *http.Request) {
	session, err := CheckAuthorization(w, r, false, false)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, tu := session.twitchUser()
	if len(tu.Data) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user := getMC(tu)

	resp, _ := json.Marshal(user)
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func getMC(tu *TwitchUsers) *mcUser {
	var ctx = context.Background()
	client, err := datastore.NewClient(ctx, "timcole-me")
	if err != nil {
		panic(err)
	}

	key := datastore.NameKey("Minecraft", tu.Data[0].ID, nil)
	user := &mcUser{}
	client.Get(ctx, key, user)

	return user
}

// Set the mc name for the logged in user
func mcWhitelistPOST(w http.ResponseWriter, r *http.Request) {
	session, err := CheckAuthorization(w, r, false, false)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, tu := session.twitchUser()
	if len(tu.Data) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Already Whitelisted?
	user := getMC(tu)
	if user.Minecraft != "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	fmt.Println(user)

	var ctx = context.Background()
	client, err := datastore.NewClient(ctx, "timcole-me")
	if err != nil {
		panic(err)
	}

	// Clean Name?
	username := r.FormValue("username")
	reg := regexp.MustCompile(`^[a-zA-Z0-9_]*$`)
	matches := reg.Find([]byte(username))
	if len(matches) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	checks := jobRunner(session, tu)
	hasAccess := false
	for _, check := range checks {
		if check.Type == "subbed" && check.CreatedAt != "" {
			hasAccess = true
		}
	}
	if !hasAccess {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	testKey := datastore.NameKey("Minecraft", tu.Data[0].ID, nil)
	entry := mcUser{Minecraft: username}
	fmt.Println(entry)

	if _, err := client.Put(ctx, testKey, &entry); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println(err)
		return
	}

	command := exec.Command(`/home/coletimmy930/whitelist.sh`, username, tu.Data[0].DisplayName)
	out, err := command.Output()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))

	resp, _ := json.Marshal(entry)
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
