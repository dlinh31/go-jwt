package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var MySigningKey = []byte(os.Getenv("SECRET_KEY"))

func GetJWT() (string, error){
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["client"] = "linh"
	claims["aud"] = "jwtgo.io"
	claims["iss"] = "issuer"
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix()
	tokenString, err := token.SignedString(MySigningKey)
	if err != nil {
		log.Printf("something went wrong: %s", err.Error())
	}
	return tokenString, nil
}

func Index(w http.ResponseWriter, r *http.Request){
	validToken, err := GetJWT()
	if err != nil {
		fmt.Println("Failed to generate JWT")
	}
	fmt.Fprint(w, string(validToken))
	
}

func HandleRequests(){
	http.HandleFunc("/", Index)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func main(){
	HandleRequests()
}