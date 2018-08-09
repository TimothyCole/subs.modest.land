package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
)

var (
	sessions      = make(map[SessionID]SessionBody)
	sessionsMutex = &sync.Mutex{}
)

func main() {
	router := mux.NewRouter()

	var static = http.StripPrefix("/assets", http.FileServer(http.Dir("./build")))
	router.PathPrefix("/assets").Handler(static)

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		temp, _ := template.ParseFiles("./build/index.html")
		temp.Execute(w, struct {
			TwitchClientID       string
			TwitchClientRedirect string
		}{
			TwitchClientID:       os.Getenv("TWITCH_CLIENT_ID"),
			TwitchClientRedirect: os.Getenv("TWITCH_CLIENT_REDIRECT"),
		})
	})

	// Remove later just to see total sessions in the map
	go func() {
		sTimer := time.NewTicker(time.Minute * 5)
		for range sTimer.C {
			fmt.Println("Logged in Sessions " + strconv.Itoa(len(sessions)))
		}
	}()

	router.HandleFunc("/authorize", TwitchAuthorize).Methods("GET")
	router.HandleFunc("/me", Me).Methods("GET")
	router.HandleFunc("/logout", Logout).Methods("DELETE")
	router.HandleFunc("/whitelist", mcWhitelistGET).Methods("GET")
	router.HandleFunc("/whitelist", mcWhitelistPOST).Methods("POST")

	fmt.Println("Starting server on :80")
	panic(http.ListenAndServe(":80", router))
}
