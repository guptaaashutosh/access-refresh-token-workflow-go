package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(id string) (string,error){
     loggedInToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
	 jwt.MapClaims{
		"id":id,
		"exp":time.Now().Add(time.Hour * 24).Unix(),
	 })
	
	 tokenString,err:=loggedInToken.SignedString([]byte(os.Getenv("SECRET")))

	 fmt.Println(tokenString)

	 if err != nil {
		return "", err
		}
	
	 return tokenString, nil

}