package main

import "net/http"

// Logout removes the session from the sessions list
func Logout(w http.ResponseWriter, r *http.Request) {
	session, err := CheckAuthorization(w, r, false, false)
	if err != nil {
		return
	}

	session.Delete()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": 200}`))
}
