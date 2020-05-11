package pkg

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

const redirectURL = "http://127.0.0.1:12345"

// MindAuth is an application that can be plugged into mind
// using authentication based on oauth2.
type MindAuth interface {
	// GetAuthURL fetches the URL to authorize the service.
	GetAuthURL(state string) string

	// GetToken fetches the token to access the API token.
	GetToken(code string) (string, error)
}

func serveAndFetchToken(auth MindAuth, state string) (string, error) {
	s := &http.Server{Addr: ":12345", Handler: nil}
	ctx, cancel := context.WithCancel(context.Background())

	var token string

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		st := r.URL.Query().Get("state")
		if st != state {
			return
		}

		tkn, err := auth.GetToken(code)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Sorry, authorization failed"))
			logrus.WithError(err).Errorln("Cannot authorize app")
			cancel()
			return
		}

		w.WriteHeader(200)
		w.Write([]byte("Authorization successful"))
		logrus.Infoln("Authorization successful")

		token = tkn
		cancel()
	})

	errChan := make(chan error)
	go func(server *http.Server, err chan<- error) {
		err <- s.ListenAndServe()
	}(s, errChan)

	select {
	case <-ctx.Done():
		// See if the server shutsdown in this amount of time
		shutDwnCtx, shutDwnCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutDwnCancel()
		err := s.Shutdown(shutDwnCtx)
		if err != nil {
			// Force shutdown if cannot gracefully
			return token, s.Close()
		}
		// Wait for server to close and it should by now
		return token, nil

	case err := <-errChan:
		return token, err
	}
}

// RunAuthCmd runs the auth command for a MindAuth type.
func RunAuthCmd(auth MindAuth, tokenToUpdate *string) {
	state := randomToken()
	logrus.WithField("url", auth.GetAuthURL(state)).Infoln("Visit the given URL")

	token, err := serveAndFetchToken(auth, state)
	if err != nil {
		logrus.WithError(err).Fatalln("Error closing server")
	}

	*tokenToUpdate = token
}

func randomToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}
