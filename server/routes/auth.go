package routes

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"bcDashboard/services"

	"github.com/lib/pq"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func SignUp(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := services.IsAuthenticated(r)
		if err == nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("Already logged in"))
			return
		}
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		user := User{}
		if err := decoder.Decode(&user); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}
		hash, err := services.GenerateHash(user.Password)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}
		tx, err := db.Begin()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}
		_, err = tx.Exec(`
			INSERT INTO "users"("username", "password")
			VALUES($1, $2);
		`, user.Username, hash)
		if err != nil {
			tx.Rollback()
			if err, ok := err.(*pq.Error); ok {
				if err.Code.Name() == "unique_violation" {
					w.WriteHeader(http.StatusNotAcceptable)
					w.Write([]byte("Username already exists"))
					return
				}
			}
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}
		err = tx.Commit()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}
		token, err := services.GenerateJwtToken(user.Username)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Can't generate a token. Please try again a few minutes later."))
			return
		}
		http.SetCookie(w, &http.Cookie{
			Value:    token,
			Expires:  time.Now().Add(7 * 24 * time.Hour),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
			Name:     "token",
			Path:     "/",
		})
		w.Write([]byte("Successfully created the user"))
	})
}

func Login(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := services.IsAuthenticated(r)
		if err == nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("Already logged in"))
			return
		}
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		var user User
		if err := decoder.Decode(&user); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}
		result, err := db.Query(`SELECT "password" FROM "users" WHERE "username" = $1;`, user.Username)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("Username or password is wrong"))
				return
			}
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}
		defer result.Close()
		var encodedHash string
		for result.Next() {
			scanErr := result.Scan(&encodedHash)
			if scanErr != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte(scanErr.Error()))
				return
			}
		}
		isMatch, err := services.VerifyPassword(user.Password, encodedHash)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}
		if !isMatch {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("Not a match"))
			return
		}
		token, err := services.GenerateJwtToken(user.Username)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}
		http.SetCookie(w, &http.Cookie{
			Value:    token,
			Expires:  time.Now().Add(7 * 24 * time.Hour),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
			Name:     "token",
			Path:     "/",
		})
		w.Write([]byte("Successfully logged in"))
	})
}

func Verify() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name, err := services.IsAuthenticated(r)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte(err.Error()))
			return
		}
		token, err := services.GenerateJwtToken(name)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}
		http.SetCookie(w, &http.Cookie{
			Value:    token,
			Expires:  time.Now().Add(7 * 24 * time.Hour),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
			Name:     "token",
			Path:     "/",
		})
		w.Write([]byte("Successfully refreshed your token."))
	})
}
