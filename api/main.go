package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Super secret info after being authorized\n")
}

func isAuthorized(endpoint func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("invalid signing method")
				}
				aud := "jwtgo.io"
				iss := "issuer"
				checkAudience := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
				if !checkAudience {
					return nil, fmt.Errorf("invalid audience")
				}
				checkIssuer := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
				if !checkIssuer {
					return nil, fmt.Errorf("invalid issuer")
				}
				return MySigningKey, nil
			})
			if err != nil {
				fmt.Fprint(w, err.Error())
				return
			}
			if token.Valid {
				endpoint(w, r)
			}
		} else {
			fmt.Fprintf(w, "No Token Found")
		}
	})
}

var MySigningKey = []byte(os.Getenv("SECRET_KEY"))

func HandleRequests() {
	http.Handle("/", isAuthorized(homePage))
	log.Fatal(http.ListenAndServe(":9001", nil))
}

func main() {
	fmt.Println("server")
	HandleRequests()
}