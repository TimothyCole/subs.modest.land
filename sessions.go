package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	signingKey = os.Getenv("SIGNING_KEY")
)

// SessionID represents a valid session ID
type SessionID string

// NewSessionID creates and returns a new session ID
func NewSessionID() (SessionID, error) {
	buf := make([]byte, 64)

	_, err := rand.Read(buf[:32])
	if err != nil {
		return SessionID(""), err
	}

	mac := hmac.New(sha256.New, []byte(signingKey))
	mac.Write(buf[:32])
	sig := mac.Sum(nil)
	copy(buf[32:], sig)

	return SessionID(base64.URLEncoding.EncodeToString(buf)), nil
}

// SessionValidateID checks to see if the session id is valid with signingKey
func SessionValidateID(id string) (SessionID, error) {
	buf, err := base64.URLEncoding.DecodeString(id)
	if err != nil {
		return "", err
	}

	if len(buf) < 64 {
		return SessionID(""), errors.New("Session ID Invalid")
	}

	mac := hmac.New(sha256.New, []byte(signingKey))
	mac.Write(buf[:32])
	messageMAC := mac.Sum(nil)
	if !hmac.Equal(messageMAC, buf[32:]) {
		return SessionID(""), errors.New("Session ID Invalid")
	}
	return SessionID(id), nil
}

// SessionBody is the stucture of the Sessions body
type SessionBody struct {
	Type   string    `json:"type"`
	Token  SessionID `json:"token"`
	Twitch struct {
		Auth struct {
			AccessToken  string `json:"access_token,omitempty"`
			RefreshToken string `json:"refresh_token,omitempty"`
			ExpiresIn    int    `json:"expires_in,omitempty"`
		} `json:"auth,omitempty"`
	} `json:"twitch,omitempty"`
}

// CheckAuthorization uses gin.Context to check for the Authorization header and make sure it's valid. If nil the create a new one and make it valid.
func CheckAuthorization(w http.ResponseWriter, r *http.Request, createOnNil bool, allowState bool) (SessionBody, error) {
	authorization := r.Header.Get("Authorization")

	state := r.URL.Query().Get("state")
	if state != "" && allowState {
		authorization = fmt.Sprintf("Session %s", state)
	}

	if authorization == "" && createOnNil {
		session, err := NewSessionID()

		body := SessionBody{
			Type:  "Session",
			Token: session,
		}

		sessionsMutex.Lock()
		sessions[session] = body
		sessionsMutex.Unlock()

		cookie := &http.Cookie{
			Name:    "modestguard",
			Value:   session.String(),
			Expires: time.Now().Add(365 * 24 * time.Hour),
			Path:    "/",
		}
		http.SetCookie(w, cookie)

		return body, err
	}

	if strings.Contains(authorization, "Session ") {
		sessionHeader := strings.Replace(authorization, "Session ", "", 1)
		sessionID, err := SessionValidateID(sessionHeader)
		if err != nil {
			return SessionForbidden(w)
		}

		body := sessions[sessionID]

		if err != nil {
			return SessionForbidden(w)
		}

		return body, err
	}

	return SessionForbidden(w)
}

// SessionForbidden just sends a Forbidden message back for unauthed users
func SessionForbidden(w http.ResponseWriter) (SessionBody, error) {
	w.Header().Set("Content-Type", "application/json")

	var out = struct {
		Status int    `json:"status"`
		Error  string `json:"error"`
	}{
		Status: http.StatusForbidden,
		Error:  "Forbidden",
	}
	w.WriteHeader(http.StatusForbidden)
	jsonOut, _ := json.Marshal(out)
	w.Write(jsonOut)

	return SessionBody{}, errors.New("Forbidden Session ID")
}

// Update pushes session updates to sessions list
func (sb SessionBody) Update() SessionBody {
	sessionsMutex.Lock()
	sessions[sb.Token] = sb
	sessionsMutex.Unlock()

	return sb
}

// Delete removes session from sessions list
func (sb SessionBody) Delete() SessionBody {
	sessionsMutex.Lock()
	delete(sessions, sb.Token)
	sessionsMutex.Unlock()

	return SessionBody{}
}

func (sid SessionID) String() string {
	return string(sid)
}
