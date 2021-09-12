package web

import (
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/saarmornel/reading-list/repo"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func Signup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := &repo.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	user.Password = string(hashedPassword)
	if _, err = repo.CreateUser(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sessionToken := uuid.NewV4().String()

	err = repo.CreateSession(sessionToken, user.Username)
	if err != nil {
		// If there is an error in setting the cache, return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(120 * time.Second),
	})
}

func Signin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := &repo.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	storedUser, err := repo.GetUser(user.Username)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if storedUser == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

	sessionToken := uuid.NewV4().String()

	err = repo.CreateSession(sessionToken, user.Username)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(120 * time.Second),
	})
}

func Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Expires: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	})
}

func Auth(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		c, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				// If the cookie is not set, return an unauthorized status
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			// For any other type of error, return a bad request status
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sessionToken := c.Value

		username, err := repo.GetSession(sessionToken)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if username == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "Username", username)
		h(w, r.WithContext(ctx), ps)
	}
}
