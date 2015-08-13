package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

const (
	TOKEN_VALID_FOR       = time.Hour * 72
	TOKEN_INDVALID_BEFORE = time.Minute * 15
)

var (
	ALGORITHM = jwt.SigningMethodRS256
)

func tokenCreationHandler(db *sql.DB, private_key []byte) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		info := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{}
		err := decoder.Decode(&info)
		if err != nil {
			RespondBadRequest(w, "Invalid request format")
			return
		}
		if info.Username == "" {
			RespondBadRequest(w, "Username is required")
			return
		}

		user_id, err := authenticate(db, info.Username, info.Password)
		if err != nil {
			RespondBadRequest(w, "Invalid credentials")
			return
		}

		token, err := createToken(user_id, private_key)
		if err != nil {
			RespondBadRequest(w, "Token could not be created")
			return
		}
		tokenFormatter := struct {
			Token string `json:"token"`
		}{
			Token: token,
		}

		RespondSuccess(w, tokenFormatter)
	}

}

func tokenVerificationHandler(public_key []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		info := struct {
			Token string `json:"token"`
		}{}
		err := decoder.Decode(&info)
		if err != nil {
			RespondBadRequest(w, "Invalid request format")
			return
		}
		err = verifyToken(info.Token, public_key)
		if err != nil {
			RespondBadRequest(w, "Token is invalid or could not be processed")
			return
		}

		RespondSuccess(w, nil)
	}

}

func verifyToken(tokenStr string, public_key []byte) error {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return public_key, nil
	})

	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("Token is not valid")
	}

	if token.Header["alg"] != ALGORITHM.Alg() {
		return errors.New("Algorithm not supported")
	}

	return nil
}

func authenticate(db *sql.DB, username, password string) (int, error) {
	var id int
	var hash []byte
	err := db.QueryRow("SELECT id, password from users where username=?", username).Scan(&id, &hash)
	if err != nil {
		return -1, err
	}

	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	return id, err
}

func createToken(id int, private_key []byte) (string, error) {
	token := jwt.New(ALGORITHM)

	token.Claims["user_id"] = id
	token.Claims["exp"] = time.Now().Add(TOKEN_VALID_FOR).Unix()
	token.Claims["iat"] = time.Now().Unix()
	token.Claims["nbf"] = time.Now().Add(-TOKEN_INDVALID_BEFORE).Unix()

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(private_key)
	return tokenString, err
}
