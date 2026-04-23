package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const claimPasswordHash = "ph"

var authPassword string

func initAuth() {
	authPassword = os.Getenv("TODO_PASSWORD")
}

func passwordHash(p string) string {
	h := sha256.Sum256([]byte(p))
	return hex.EncodeToString(h[:])
}

func newAuthToken() (string, error) {
	pass := authPassword
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		claimPasswordHash: passwordHash(pass),
		"exp":             time.Now().Add(8 * time.Hour).Unix(),
	})
	return t.SignedString([]byte(pass))
}

func validateTokenString(tokenString string) error {
	pass := authPassword
	if pass == "" {
		return nil
	}
	t, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(pass), nil
	})
	if err != nil {
		return err
	}
	mc, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid claims")
	}
	ph, _ := mc[claimPasswordHash].(string)
	if ph != passwordHash(pass) {
		return errors.New("token password mismatch")
	}
	return nil
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorStatus(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}
	var body struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErrorStatus(w, http.StatusBadRequest, err)
		return
	}
	want := authPassword
	if want == "" {
		writeErrorStatus(w, http.StatusUnauthorized, errors.New("аутентификация не настроена"))
		return
	}
	if body.Password != want {
		writeErrorStatus(w, http.StatusUnauthorized, errors.New("Неверный пароль"))
		return
	}
	token, err := newAuthToken()
	if err != nil {
		writeErrorStatus(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, map[string]string{"token": token})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if authPassword == "" {
			next(w, r)
			return
		}
		cookie, err := r.Cookie("token")
		if err != nil || cookie.Value == "" {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}
		if err := validateTokenString(strings.TrimSpace(cookie.Value)); err != nil {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
